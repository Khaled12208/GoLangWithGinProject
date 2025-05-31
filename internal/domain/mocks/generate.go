package mocks

//go:generate mockgen -destination=task_repository_mock.go -package=mocks golangwithgin/internal/domain TaskRepository
//go:generate mockgen -destination=task_processor_mock.go -package=mocks golangwithgin/internal/domain TaskProcessor
//go:generate mockgen -destination=task_service_mock.go -package=mocks golangwithgin/internal/domain TaskService 