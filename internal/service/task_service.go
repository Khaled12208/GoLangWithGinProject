package service

import (
	"golangwithgin/internal/domain"
	"time"
)

// taskService implements the TaskService interface
type taskService struct {
	repository domain.TaskRepository
	processor  domain.TaskProcessor
}

// NewTaskService creates a new task service
func NewTaskService(repository domain.TaskRepository, processor domain.TaskProcessor) domain.TaskService {
	return &taskService{
		repository: repository,
		processor:  processor,
	}
}

func (s *taskService) SubmitTask(task *domain.Task) error {
	// Set initial task state
	task.Status = "pending"
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	// Save task to database
	if err := s.repository.Create(task); err != nil {
		return err
	}

	// Process task asynchronously
	go func() {
		// Update task status to processing
		task.Status = "processing"
		task.UpdatedAt = time.Now()
		if err := s.repository.Update(task); err != nil {
			return
		}

		// Process the task
		if err := s.processor.Process(task); err != nil {
			task.Status = "failed"
		} else {
			task.Status = "completed"
		}

		// Update task in database after processing
		task.UpdatedAt = time.Now()
		if err := s.repository.Update(task); err != nil {
			// Log error but don't return it since we're in a goroutine
			// In a production system, you'd want to handle this error properly
			return
		}
	}()

	return nil
}

func (s *taskService) GetTaskStatus(id uint) (*domain.Task, error) {
	return s.repository.FindByID(id)
}

func (s *taskService) GetAllTasks() ([]*domain.Task, error) {
	return s.repository.FindAll()
} 