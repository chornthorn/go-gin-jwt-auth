package services

import (
	"errors"
	"gorm.io/gorm"
	"jwt-auth-app/config"
	"jwt-auth-app/model"
	"jwt-auth-app/types"
	"jwt-auth-app/utils"
)

type UsersService struct {
	DB *gorm.DB
}

func NewUsersService() *UsersService {
	return &UsersService{DB: config.DB}
}

// GetUserByID retrieves a user by their ID
func (s *UsersService) GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	if err := s.DB.First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrUserNotFound
		}
		return nil, utils.ErrInternalServer
	}
	return &user, nil
}

// GetAuthenticatedUser converts a model.User to types.AuthenticatedUser
func (s *UsersService) GetAuthenticatedUser(user *model.User) *types.AuthenticatedUser {
	return &types.AuthenticatedUser{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
}

// ValidateAndGetUser validates the token and returns the user
func (s *UsersService) ValidateAndGetUser(token string) (*types.AuthenticatedUser, *types.TokenMetadata, error) {
	// Validate token
	tokenMetadata, err := utils.ValidateAccessToken(token)
	if err != nil {
		return nil, nil, utils.ErrInvalidToken
	}

	// Get user from database
	user, err := s.GetUserByID(tokenMetadata.UserID)
	if err != nil {
		return nil, nil, err
	}

	authenticatedUser := s.GetAuthenticatedUser(user)
	return authenticatedUser, tokenMetadata, nil
}

func (s *UsersService) UpdateProfile(userID uint, req *types.UpdateProfileRequest) (*types.UserResponse, error) {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	// Update only the fields that are provided
	if req.Name != "" {
		user.Name = req.Name
	}
	// Add more fields as needed

	if err := s.DB.Save(user).Error; err != nil {
		return nil, utils.ErrInternalServer
	}

	return &types.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}, nil
}
