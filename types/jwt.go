package types

import "github.com/golang-jwt/jwt/v5"

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	UserID    uint      `json:"user_id"`
	TokenType TokenType `json:"token_type"`
}

type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type TokenMetadata struct {
	UserID    uint
	TokenType TokenType
	IssuedAt  int64
	ExpiresAt int64
}
