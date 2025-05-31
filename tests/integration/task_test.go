package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golangwithgin/config"
	"golangwithgin/internal/app/server"
	"golangwithgin/internal/domain"
	"golangwithgin/pkg/logger"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ory/dockertest/v3"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TaskIntegrationTestSuite struct {
	suite.Suite
	pool     *dockertest.Pool
	resource *dockertest.Resource
	db       *gorm.DB
	server   *server.Server
	router   *gin.Engine
	token    string
}

func TestTaskIntegrationSuite(t *testing.T) {
	suite.Run(t, new(TaskIntegrationTestSuite))
}

func (s *TaskIntegrationTestSuite) SetupSuite() {
	var err error

	// Create a new Docker pool
	s.pool, err = dockertest.NewPool("")
	s.Require().NoError(err)

	// Start MySQL container
	s.resource, err = s.pool.Run("mysql", "8.0", []string{
		"MYSQL_ROOT_PASSWORD=secret",
		"MYSQL_DATABASE=test_db",
	})
	s.Require().NoError(err)

	// Wait for MySQL to be ready
	var db *gorm.DB
	err = s.pool.Retry(func() error {
		dsn := fmt.Sprintf("root:secret@tcp(localhost:%s)/test_db?charset=utf8mb4&parseTime=True&loc=Local",
			s.resource.GetPort("3306/tcp"))
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Ping()
	})
	s.Require().NoError(err)
	s.db = db

	// Run migrations
	err = s.db.AutoMigrate(&domain.Task{}, &domain.User{})
	s.Require().NoError(err)

	// Convert port string to int
	dbPort, err := strconv.Atoi(s.resource.GetPort("3306/tcp"))
	s.Require().NoError(err)

	// Create test configuration
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port: "8888",
		},
		Database: config.DatabaseConfig{
			Host:     "localhost",
			Port:     dbPort,
			Username: "root",
			Password: "secret",
			DBName:   "test_db",
		},
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			Expiration: 24 * time.Hour,
		},
		Logger: config.LoggerConfig{
			Level: "info",
			File:  "",  // Empty string for stdout only
		},
	}

	// Initialize logger
	log := logger.New()  // Use the simple logger for tests

	// Initialize server
	s.server = server.New(cfg, log)
	s.router = s.server.GetRouter()

	// Create test user and get JWT token
	s.setupTestUser()
}

func (s *TaskIntegrationTestSuite) setupTestUser() {
	// Create test user
	user := &domain.User{
		Username: "testuser",
		Password: "testpass",
		Email:    "test@example.com",
	}

	// Register user
	body, err := json.Marshal(map[string]string{
		"username": user.Username,
		"password": user.Password,
		"email":    user.Email,
	})
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(body))
	s.server.GetRouter().ServeHTTP(w, req)
	s.Equal(http.StatusCreated, w.Code)

	// Login and get token
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
	s.server.GetRouter().ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	s.Require().NoError(err)
	s.token = response["token"]
}

func (s *TaskIntegrationTestSuite) TearDownSuite() {
	if s.resource != nil {
		s.pool.Purge(s.resource)
	}
}

func (s *TaskIntegrationTestSuite) TestTaskLifecycle() {
	// Create a new task
	task := map[string]string{
		"title":       "Integration Test Task",
		"description": "Testing task lifecycle",
	}
	body, err := json.Marshal(task)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+s.token)
	s.server.GetRouter().ServeHTTP(w, req)
	s.Equal(http.StatusAccepted, w.Code)

	var createdTask domain.Task
	err = json.Unmarshal(w.Body.Bytes(), &createdTask)
	s.Require().NoError(err)
	s.NotZero(createdTask.ID)
	s.Equal(task["title"], createdTask.Title)
	s.Equal(task["description"], createdTask.Description)
	s.Equal("pending", createdTask.Status)

	// Wait for task processing
	time.Sleep(2 * time.Second)

	// Get task status
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/tasks/%d", createdTask.ID), nil)
	req.Header.Set("Authorization", "Bearer "+s.token)
	s.server.GetRouter().ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)

	var processedTask domain.Task
	err = json.Unmarshal(w.Body.Bytes(), &processedTask)
	s.Require().NoError(err)
	s.Equal(createdTask.ID, processedTask.ID)
	s.Equal("completed", processedTask.Status)

	// Get all tasks
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+s.token)
	s.server.GetRouter().ServeHTTP(w, req)
	s.Equal(http.StatusOK, w.Code)

	var tasks []domain.Task
	err = json.Unmarshal(w.Body.Bytes(), &tasks)
	s.Require().NoError(err)
	s.NotEmpty(tasks)
	s.Contains(tasks, processedTask)
}

func (s *TaskIntegrationTestSuite) TestConcurrentTaskProcessing() {
	numTasks := 5
	tasks := make([]domain.Task, numTasks)

	// Create multiple tasks concurrently
	for i := 0; i < numTasks; i++ {
		task := map[string]string{
			"name":        fmt.Sprintf("Concurrent Task %d", i+1),
			"description": "Testing concurrent processing",
		}
		body, err := json.Marshal(task)
		s.Require().NoError(err)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
		req.Header.Set("Authorization", "Bearer "+s.token)
		s.server.GetRouter().ServeHTTP(w, req)
		s.Equal(http.StatusAccepted, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &tasks[i])
		s.Require().NoError(err)
	}

	// Wait for all tasks to be processed
	time.Sleep(3 * time.Second)

	// Verify all tasks are completed
	for _, task := range tasks {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/api/v1/tasks/%d", task.ID), nil)
		req.Header.Set("Authorization", "Bearer "+s.token)
		s.server.GetRouter().ServeHTTP(w, req)
		s.Equal(http.StatusOK, w.Code)

		var processedTask domain.Task
		err := json.Unmarshal(w.Body.Bytes(), &processedTask)
		s.Require().NoError(err)
		s.Equal("completed", processedTask.Status)
	}
} 