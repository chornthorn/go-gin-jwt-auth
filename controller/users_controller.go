package controller

import (
	"github.com/gin-gonic/gin"
	"jwt-auth-app/middleware"
	"jwt-auth-app/services"
	"jwt-auth-app/types"
	"jwt-auth-app/utils"
	"net/http"
)

type UserController struct {
	usersService *services.UsersService
}

func NewUserController() *UserController {
	return &UserController{
		usersService: services.NewUsersService(),
	}
}

func (uc *UserController) GetProfile(c *gin.Context) {
	// Get authenticated user from context
	authUser, err := middleware.GetAuthUser(c)
	if err != nil {
		status, errResponse := utils.GetErrorResponse(err)
		c.JSON(status, errResponse)
		return
	}

	// Get full user profile from service
	user, err := uc.usersService.GetUserByID(authUser.ID)
	if err != nil {
		status, errResponse := utils.GetErrorResponse(err)
		c.JSON(status, errResponse)
		return
	}

	userResponse := types.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}

	c.JSON(http.StatusOK, gin.H{
		"profile": userResponse,
	})
}

func (uc *UserController) UpdateProfile(c *gin.Context) {
	var req types.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, types.ErrorResponse{
			Code:    "INVALID_REQUEST",
			Message: "Invalid request body",
		})
		return
	}

	authUser, err := middleware.GetAuthUser(c)
	if err != nil {
		status, errResponse := utils.GetErrorResponse(err)
		c.JSON(status, errResponse)
		return
	}

	updatedUser, err := uc.usersService.UpdateProfile(authUser.ID, &req)
	if err != nil {
		status, errResponse := utils.GetErrorResponse(err)
		c.JSON(status, errResponse)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profile": updatedUser,
	})
}
