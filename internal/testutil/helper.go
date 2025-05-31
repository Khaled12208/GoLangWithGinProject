package testutil

import (
	"golangwithgin/internal/domain"
	"time"
)

// CreateTestTask creates a task for testing
func CreateTestTask() *domain.Task {
	return &domain.Task{
		ID:          1,
		Title:       "Test Task",
		Description: "Test Description",
		Status:      "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// CreateTestTasks creates multiple tasks for testing
func CreateTestTasks(count int) []*domain.Task {
	tasks := make([]*domain.Task, count)
	for i := 0; i < count; i++ {
		tasks[i] = &domain.Task{
			ID:          uint(i + 1),
			Title:       "Test Task " + string(rune(i+1)),
			Description: "Test Description " + string(rune(i+1)),
			Status:      "pending",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
	}
	return tasks
} 