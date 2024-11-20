package middleware

import (
	"github.com/gin-gonic/gin"
	"jwt-auth-app/services"
	"jwt-auth-app/types"
	"jwt-auth-app/utils"
	"net/http"
	"strings"
)

// AuthMiddleware contains the dependencies for the auth middleware
type AuthMiddleware struct {
	usersService *services.UsersService
}

// NewAuthMiddleware creates a new auth middleware instance
func NewAuthMiddleware() *AuthMiddleware {
	return &AuthMiddleware{
		usersService: services.NewUsersService(),
	}
}

// JWT middleware verifies the access token and loads the user into the context
func (m *AuthMiddleware) JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, types.ErrorResponse{
				Code:    "UNAUTHORIZED",
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		// Validate token and get user
		authenticatedUser, tokenMetadata, err := m.usersService.ValidateAndGetUser(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, types.ErrorResponse{
				Code:    "INVALID_TOKEN",
				Message: err.Error(),
			})
			c.Abort()
			return
		}

		// Store user and token metadata in context
		c.Set(string(UserContextKey), *authenticatedUser)
		c.Set(string(TokenMetadataKey), tokenMetadata)

		c.Next()
	}
}

// RequireRole middleware checks if the user has the required role
func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation for role-based access control
		// Add this when you implement role management
		c.Next()
	}
}

// GetAuthUser helper function to get the authenticated user from context
func GetAuthUser(c *gin.Context) (*AuthenticatedUser, error) {
	user, exists := c.Get(string(UserContextKey))
	if !exists {
		return nil, utils.ErrUnauthorized
	}

	authUser, ok := user.(AuthenticatedUser)
	if !ok {
		return nil, utils.ErrUnauthorized
	}

	return &authUser, nil
}

// GetTokenMetadata helper function to get the token metadata from context
func GetTokenMetadata(c *gin.Context) (*types.TokenMetadata, error) {
	metadata, exists := c.Get(string(TokenMetadataKey))
	if !exists {
		return nil, utils.ErrUnauthorized
	}

	// Type assertion to convert interface{} to TokenMetadata
	tokenMetadata, ok := metadata.(*types.TokenMetadata)
	if !ok {
		return nil, utils.ErrUnauthorized
	}

	return tokenMetadata, nil
}

// extractToken extracts the token from the Authorization header
func extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", utils.ErrMissingAuthHeader
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", utils.ErrInvalidAuthHeader
	}

	return parts[1], nil
}
