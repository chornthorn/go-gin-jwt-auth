package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"jwt-auth-app/config"
	"jwt-auth-app/controller"
	"jwt-auth-app/middleware"
	"jwt-auth-app/utils"
	"log"
	"net/http"
)

func main() {
	// Load environment variables
	config.LoadConfig()

	// Set Gin mode
	gin.SetMode(config.AppConfig.Server.GinMode)

	// Initialize JWT keys
	if err := utils.InitializeJWTManager(&config.AppConfig.JWT); err != nil {
		log.Fatal("Failed to initialize JWT keys:", err)
	}

	// Initialize Middleware
	authMiddleware := middleware.NewAuthMiddleware()

	// Initialize Controllers
	authController := controller.NewAuthController()
	userController := controller.NewUserController()

	// Create Gin router
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	// API routes
	api := r.Group("/api/v1")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authController.Register)
			auth.POST("/login", authController.Login)
			auth.POST("/refresh", authController.RefreshToken)
		}

		// Protected routes
		protected := api.Group("")
		protected.Use(authMiddleware.JWT())
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("/profile", userController.GetProfile)
				users.PUT("/profile", userController.UpdateProfile)
			}

			// Token info route
			protected.GET("/token/info", func(c *gin.Context) {
				metadata, err := middleware.GetTokenMetadata(c)
				if err != nil {
					status, errResponse := utils.GetErrorResponse(err)
					c.JSON(status, errResponse)
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"token_metadata": metadata,
				})
			})
		}
	}

	// Start server
	serverAddr := fmt.Sprintf(":%s", config.AppConfig.Server.Port)
	log.Printf("Server starting on port %s", config.AppConfig.Server.Port)
	if err := r.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
