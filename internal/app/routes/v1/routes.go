package v1

import (
	"github.com/gin-gonic/gin"
	"golangwithgin/internal/app/handlers"
	"golangwithgin/internal/app/middlewares"
)

func SetupRoutes(
	router *gin.Engine,
	userHandler *handlers.UserHandler,
	taskHandler *handlers.TaskHandler,
	authMiddleware *middlewares.AuthMiddleware,
) {
	v1 := router.Group("/api/v1")
	{
		// Public routes
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)

		// Protected routes
		protected := v1.Group("/")
		protected.Use(authMiddleware.AuthRequired)
		{
			// User routes
			protected.GET("/user", userHandler.GetUser)
			protected.PUT("/user", userHandler.UpdateUser)

			// Task routes
			protected.POST("/tasks", taskHandler.CreateTask)
			protected.GET("/tasks", taskHandler.GetAllTasks)
			protected.GET("/tasks/:id", taskHandler.GetTask)
		}
	}
} 