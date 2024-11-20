package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"jwt-auth-app/config"
	"jwt-auth-app/model"
	"jwt-auth-app/types"
	"jwt-auth-app/utils"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService() *AuthService {
	return &AuthService{
		db: config.DB,
	}
}

func (s *AuthService) Register(req *types.RegisterRequest) (*types.AuthResponse, error) {
	// Check if user exists
	var existingUser model.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, utils.ErrUserExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, utils.ErrInternalServer
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, utils.ErrInternalServer
	}

	// Create user
	user := model.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, utils.ErrInternalServer
	}

	// Generate tokens
	tokens, err := utils.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, utils.ErrInternalServer
	}

	return &types.AuthResponse{
		User: types.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
		Token: *tokens,
	}, nil
}

func (s *AuthService) Login(req *types.LoginRequest) (*types.AuthResponse, error) {
	// Find user
	var user model.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrInvalidCredentials
		}
		return nil, utils.ErrInternalServer
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, utils.ErrInvalidCredentials
	}

	// Generate tokens
	tokens, err := utils.GenerateTokenPair(user.ID)
	if err != nil {
		return nil, utils.ErrInternalServer
	}

	return &types.AuthResponse{
		User: types.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		},
		Token: *tokens,
	}, nil
}
