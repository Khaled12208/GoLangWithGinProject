package integration

import (
	"bytes"
	"encoding/json"
	"golangwithgin/internal/app/handlers"
	"golangwithgin/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockUserService is a mock implementation of domain.UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Register(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) Login(username, password string) (string, error) {
	args := m.Called(username, password)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) GetByID(id uint) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserService) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserService) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

type UserHandlerTestSuite struct {
	suite.Suite
	router          *gin.Engine
	mockUserService *MockUserService
}

func TestUserHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}

func (s *UserHandlerTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	s.router = gin.New()
	s.mockUserService = new(MockUserService)
	userHandler := handlers.NewUserHandler(s.mockUserService)

	// Setup routes
	v1 := s.router.Group("/api/v1")
	{
		v1.POST("/register", userHandler.Register)
		v1.POST("/login", userHandler.Login)
		v1.GET("/users/:id", userHandler.GetUser)
		v1.PUT("/users/:id", userHandler.UpdateUser)
		v1.DELETE("/users/:id", userHandler.DeleteUser)
	}
}

func (s *UserHandlerTestSuite) TestRegister_Success() {
	registerRequest := handlers.RegisterRequest{
		Username: "testuser",
		Password: "testpass",
		Email:    "test@example.com",
	}

	expectedUser := &domain.User{
		Username: registerRequest.Username,
		Password: registerRequest.Password,
		Email:    registerRequest.Email,
	}

	s.mockUserService.On("Register", mock.MatchedBy(func(u *domain.User) bool {
		return u.Username == expectedUser.Username &&
			u.Password == expectedUser.Password &&
			u.Email == expectedUser.Email
	})).Return(nil)

	body, err := json.Marshal(registerRequest)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusCreated, w.Code)

	var response domain.UserResponse
	err = json.Unmarshal(w.Body.Bytes(), &response)
	s.Require().NoError(err)
	s.Equal(registerRequest.Username, response.Username)
	s.Equal(registerRequest.Email, response.Email)

	s.mockUserService.AssertExpectations(s.T())
}

func (s *UserHandlerTestSuite) TestRegister_ValidationError() {
	registerRequest := handlers.RegisterRequest{
		Username: "",
		Password: "testpass",
		Email:    "invalid-email",
	}

	body, err := json.Marshal(registerRequest)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusBadRequest, w.Code)
}

func (s *UserHandlerTestSuite) TestLogin_Success() {
	loginRequest := handlers.LoginRequest{
		Username: "testuser",
		Password: "testpass",
	}

	expectedToken := "test-token"
	s.mockUserService.On("Login", loginRequest.Username, loginRequest.Password).Return(expectedToken, nil)

	body, err := json.Marshal(loginRequest)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusOK, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	s.Require().NoError(err)
	s.Equal(expectedToken, response["token"])

	s.mockUserService.AssertExpectations(s.T())
}

func (s *UserHandlerTestSuite) TestLogin_InvalidCredentials() {
	loginRequest := handlers.LoginRequest{
		Username: "testuser",
		Password: "wrongpass",
	}

	s.mockUserService.On("Login", loginRequest.Username, loginRequest.Password).Return("", domain.ErrInvalidCredentials)

	body, err := json.Marshal(loginRequest)
	s.Require().NoError(err)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	s.router.ServeHTTP(w, req)

	s.Equal(http.StatusUnauthorized, w.Code)

	s.mockUserService.AssertExpectations(s.T())
} 