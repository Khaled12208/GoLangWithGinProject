package server

import (
	"fmt"
	"golangwithgin/config"
	"golangwithgin/internal/app/handlers"
	"golangwithgin/internal/app/middlewares"
	"golangwithgin/internal/app/routes/v1"
	"golangwithgin/internal/repository/mysql"
	"golangwithgin/internal/service"
	"golangwithgin/pkg/database"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	Router *gin.Engine
	config *config.Config
	logger *logrus.Logger
}

// New creates a new server instance
func New(cfg *config.Config, logger *logrus.Logger) *Server {
	router := gin.Default()

	// Initialize database
	db, err := database.NewMySQLDB(cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize repositories
	userRepo := mysql.NewUserRepository(db)
	taskRepo := mysql.NewTaskRepository(db)

	// Initialize task processor
	taskProcessor := service.NewTaskProcessor()

	// Initialize services
	userService := service.NewUserService(userRepo, cfg.JWT.Secret)
	taskService := service.NewTaskService(taskRepo, taskProcessor)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	taskHandler := handlers.NewTaskHandler(taskService)

	// Initialize middlewares
	authMiddleware := middlewares.NewAuthMiddleware(cfg.JWT.Secret)

	// Setup routes
	v1.SetupRoutes(router, userHandler, taskHandler, authMiddleware)

	return &Server{
		Router: router,
		config: cfg,
		logger: logger,
	}
}

// Start starts the server
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.config.Server.Port)
	s.logger.Infof("Starting server on %s", addr)
	return s.Router.Run(addr)
}

func (s *Server) Shutdown() {
	// No need to implement shutdown logic as the server is managed by gin
}

// GetRouter returns the server's router instance
func (s *Server) GetRouter() *gin.Engine {
	return s.Router
} 