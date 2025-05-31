package routes

import (
	"golangwithgin/config"
	"golangwithgin/internal/app/handlers"
	"golangwithgin/internal/app/routes/v1"
	"golangwithgin/internal/domain"
	"golangwithgin/pkg/logger"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, cfg *config.Config, log *logger.Logger, userService domain.UserService) {
	// Create handlers
	userHandler := handlers.NewUserHandler(userService)

	// API v1 routes
	apiV1 := router.Group("/api/v1")
	v1.SetupUserRoutes(apiV1, userHandler, cfg.JWT.Secret)

	// Add more versioned routes here as needed
} 