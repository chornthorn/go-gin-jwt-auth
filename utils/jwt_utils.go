package utils

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"jwt-auth-app/config"
	"jwt-auth-app/types"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTKeys struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

type TokenManager struct {
	accessKeys  JWTKeys
	refreshKeys JWTKeys
	config      *config.JWTConfig
}

var tokenManager *TokenManager

func InitializeJWTManager(cfg *config.JWTConfig) error {
	manager := &TokenManager{config: cfg}

	// Initialize access token keys
	accessKeys, err := loadKeys(cfg.AccessToken.PrivateKeyPath, cfg.AccessToken.PublicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load access token keys: %w", err)
	}
	manager.accessKeys = accessKeys

	// Initialize refresh token keys
	refreshKeys, err := loadKeys(cfg.RefreshToken.PrivateKeyPath, cfg.RefreshToken.PublicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load refresh token keys: %w", err)
	}
	manager.refreshKeys = refreshKeys

	tokenManager = manager
	return nil
}

func loadKeys(privateKeyPath, publicKeyPath string) (JWTKeys, error) {
	// Load private key
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return JWTKeys{}, fmt.Errorf("failed to read private key: %w", err)
	}

	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	if privateKeyBlock == nil {
		return JWTKeys{}, errors.New("failed to decode private key PEM block")
	}

	var privateKey *rsa.PrivateKey

	// Try PKCS8 first
	privateKeyParsed, err := x509.ParsePKCS8PrivateKey(privateKeyBlock.Bytes)
	if err == nil {
		var ok bool
		privateKey, ok = privateKeyParsed.(*rsa.PrivateKey)
		if !ok {
			return JWTKeys{}, errors.New("private key is not RSA key")
		}
	} else {
		// Fallback to PKCS1
		privateKey, err = x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
		if err != nil {
			return JWTKeys{}, fmt.Errorf("failed to parse private key: %w", err)
		}
	}

	// Load public key
	publicKeyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return JWTKeys{}, fmt.Errorf("failed to read public key: %w", err)
	}

	publicKeyBlock, _ := pem.Decode(publicKeyBytes)
	if publicKeyBlock == nil {
		return JWTKeys{}, errors.New("failed to decode public key PEM block")
	}

	parsedPublicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return JWTKeys{}, fmt.Errorf("failed to parse public key: %w", err)
	}

	publicKey, ok := parsedPublicKey.(*rsa.PublicKey)
	if !ok {
		return JWTKeys{}, errors.New("public key is not RSA key")
	}

	return JWTKeys{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func (tm *TokenManager) GenerateTokenPair(userID uint) (*types.TokenPair, error) {
	// Generate access token
	accessToken, err := tm.generateToken(userID, types.AccessToken, tm.accessKeys, tm.config.AccessToken.ExpirationTime)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := tm.generateToken(userID, types.RefreshToken, tm.refreshKeys, tm.config.RefreshToken.ExpirationTime)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &types.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (tm *TokenManager) generateToken(userID uint, tokenType types.TokenType, keys JWTKeys, expiration time.Duration) (string, error) {
	now := time.Now()
	claims := &types.CustomClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			Issuer:    tm.config.Issuer,
		},
		UserID:    userID,
		TokenType: tokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(keys.privateKey)
}

func (tm *TokenManager) ValidateToken(tokenString string, tokenType types.TokenType) (*types.TokenMetadata, error) {
	keys := tm.accessKeys
	if tokenType == types.RefreshToken {
		keys = tm.refreshKeys
	}

	token, err := jwt.ParseWithClaims(tokenString, &types.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return keys.publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*types.CustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	if claims.TokenType != tokenType {
		return nil, errors.New("invalid token type")
	}

	return &types.TokenMetadata{
		UserID:    claims.UserID,
		TokenType: claims.TokenType,
		IssuedAt:  claims.IssuedAt.Unix(),
		ExpiresAt: claims.ExpiresAt.Unix(),
	}, nil
}

// GenerateTokenPair Helper functions to expose the functionality
func GenerateTokenPair(userID uint) (*types.TokenPair, error) {
	return tokenManager.GenerateTokenPair(userID)
}

func ValidateAccessToken(tokenString string) (*types.TokenMetadata, error) {
	return tokenManager.ValidateToken(tokenString, types.AccessToken)
}

func ValidateRefreshToken(tokenString string) (*types.TokenMetadata, error) {
	return tokenManager.ValidateToken(tokenString, types.RefreshToken)
}
