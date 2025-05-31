package v1

import (
	"github.com/gin-gonic/gin"
	"golangwithgin/internal/app/handlers"
	"golangwithgin/internal/app/middlewares"
)

func SetupUserRoutes(
	router *gin.RouterGroup,
	userHandler *handlers.UserHandler,
	jwtSecret string,
) {
	// Create auth middleware
	authMiddleware := middlewares.NewAuthMiddleware(jwtSecret)

	// Public routes
	router.POST("/register", userHandler.Register)
	router.POST("/login", userHandler.Login)

	// Protected routes
	protected := router.Group("/")
	protected.Use(authMiddleware.AuthRequired)
	{
		protected.GET("/user", userHandler.GetUser)
		protected.PUT("/user", userHandler.UpdateUser)
	}
} 