package domain

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

// TokenResponse represents a successful login response
type TokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// SwaggerUserResponse represents the user data returned to clients for Swagger documentation
type SwaggerUserResponse struct {
	ID        uint   `json:"id" example:"1"`
	Username  string `json:"username" example:"johndoe"`
	Email     string `json:"email" example:"john@example.com"`
	CreatedAt string `json:"created_at" example:"2025-05-31T15:04:05Z"`
	UpdatedAt string `json:"updated_at" example:"2025-05-31T15:04:05Z"`
}

// SwaggerTask represents a task in the system for Swagger documentation
type SwaggerTask struct {
	ID          uint   `json:"id" example:"1"`
	Title       string `json:"title" example:"Process Data"`
	Description string `json:"description" example:"Process the uploaded data file"`
	Status      string `json:"status" example:"pending"`
	CreatedAt   string `json:"created_at" example:"2025-05-31T15:04:05Z"`
	UpdatedAt   string `json:"updated_at" example:"2025-05-31T15:04:05Z"`
} 