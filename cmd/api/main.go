package main

import (
	"golangwithgin/config"
	"golangwithgin/internal/app/server"
	"golangwithgin/pkg/logger"

	_ "golangwithgin/docs" // Import generated Swagger docs

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           GolangWithGin API
// @version         1.0
// @description     A RESTful API server using Go and Gin framework
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8889
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Initialize logger
	log := logger.New()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Create and start server
	srv := server.New(cfg, log)
	
	// Add Swagger documentation route
	srv.Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := srv.Start(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
} 