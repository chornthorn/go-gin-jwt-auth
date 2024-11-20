package controller

import (
	"github.com/gin-gonic/gin"
	"jwt-auth-app/services"
	"jwt-auth-app/types"
	"jwt-auth-app/utils"
	"net/http"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

func (ac *AuthController) Register(c *gin.Context) {
	var req types.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: err.Error(),
		})
		return
	}

	response, err := ac.authService.Register(&req)
	if err != nil {
		status, errResponse := utils.GetErrorResponse(err)
		c.JSON(status, errResponse)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (ac *AuthController) Login(c *gin.Context) {
	var req types.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: err.Error(),
		})
		return
	}

	response, err := ac.authService.Login(&req)
	if err != nil {
		status, errResponse := utils.GetErrorResponse(err)
		c.JSON(status, errResponse)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (ac *AuthController) RefreshToken(c *gin.Context) {
	var input types.RefreshTokenInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Code:    "INVALID_INPUT",
			Message: err.Error(),
		})
		return
	}

	// Validate refresh token
	metadata, err := utils.ValidateRefreshToken(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, types.ErrorResponse{
			Code:    "INVALID_REFRESH_TOKEN",
			Message: "Invalid refresh token",
		})
		return
	}

	// Generate new token pair
	tokenPair, err := utils.GenerateTokenPair(metadata.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, types.ErrorResponse{
			Code:    "TOKEN_GENERATION_FAILED",
			Message: "Failed to generate new tokens",
		})
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}
