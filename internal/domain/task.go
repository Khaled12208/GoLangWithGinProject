package domain

import "time"

// Task represents a task entity
type Task struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TaskRepository defines the interface for task persistence
type TaskRepository interface {
	Create(task *Task) error
	Update(task *Task) error
	FindByID(id uint) (*Task, error)
	FindAll() ([]*Task, error)
}

// TaskProcessor defines the interface for task processing
type TaskProcessor interface {
	Process(task *Task) error
	Shutdown()
}

// TaskService defines the interface for task business logic
type TaskService interface {
	SubmitTask(task *Task) error
	GetTaskStatus(id uint) (*Task, error)
	GetAllTasks() ([]*Task, error)
} 