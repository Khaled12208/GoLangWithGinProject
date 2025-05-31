package service

import (
	"errors"
	"golangwithgin/internal/domain"
	"golangwithgin/internal/domain/mocks"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type TaskServiceTestSuite struct {
	suite.Suite
	mockCtrl       *gomock.Controller
	mockRepository *mocks.MockTaskRepository
	mockProcessor  *mocks.MockTaskProcessor
	service        domain.TaskService
}

func TestTaskServiceSuite(t *testing.T) {
	suite.Run(t, new(TaskServiceTestSuite))
}

func (s *TaskServiceTestSuite) SetupTest() {
	s.mockCtrl = gomock.NewController(s.T())
	s.mockRepository = mocks.NewMockTaskRepository(s.mockCtrl)
	s.mockProcessor = mocks.NewMockTaskProcessor(s.mockCtrl)
	s.service = NewTaskService(s.mockRepository, s.mockProcessor)
}

func (s *TaskServiceTestSuite) TearDownTest() {
	s.mockCtrl.Finish()
}

func (s *TaskServiceTestSuite) TestSubmitTask_Success() {
	task := &domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
	}

	// Expect Create call
	s.mockRepository.EXPECT().
		Create(gomock.Any()).
		DoAndReturn(func(t *domain.Task) error {
			s.Equal("pending", t.Status)
			s.NotZero(t.CreatedAt)
			s.NotZero(t.UpdatedAt)
			return nil
		})

	// Expect Update call for processing status
	s.mockRepository.EXPECT().
		Update(gomock.Any()).
		DoAndReturn(func(t *domain.Task) error {
			s.Equal("processing", t.Status)
			s.NotZero(t.UpdatedAt)
			return nil
		})

	// Expect Process call
	s.mockProcessor.EXPECT().
		Process(gomock.Any()).
		Return(nil)

	// Expect Update call for completed status
	s.mockRepository.EXPECT().
		Update(gomock.Any()).
		DoAndReturn(func(t *domain.Task) error {
			s.Equal("completed", t.Status)
			s.NotZero(t.UpdatedAt)
			return nil
		})

	err := s.service.SubmitTask(task)
	s.NoError(err)

	// Wait for goroutine to complete
	time.Sleep(100 * time.Millisecond)
}

func (s *TaskServiceTestSuite) TestSubmitTask_CreateError() {
	task := &domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
	}

	expectedErr := errors.New("create error")
	s.mockRepository.EXPECT().
		Create(gomock.Any()).
		Return(expectedErr)

	err := s.service.SubmitTask(task)
	s.Equal(expectedErr, err)
}

func (s *TaskServiceTestSuite) TestSubmitTask_ProcessError() {
	task := &domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
	}

	// Expect Create call
	s.mockRepository.EXPECT().
		Create(gomock.Any()).
		Return(nil)

	// Expect Update call for processing status
	s.mockRepository.EXPECT().
		Update(gomock.Any()).
		Return(nil)

	// Expect Process call with error
	expectedErr := errors.New("process error")
	s.mockProcessor.EXPECT().
		Process(gomock.Any()).
		Return(expectedErr)

	// Expect Update call for failed status
	s.mockRepository.EXPECT().
		Update(gomock.Any()).
		DoAndReturn(func(t *domain.Task) error {
			s.Equal("failed", t.Status)
			s.NotZero(t.UpdatedAt)
			return nil
		})

	err := s.service.SubmitTask(task)
	s.NoError(err)

	// Wait for goroutine to complete
	time.Sleep(100 * time.Millisecond)
}

func (s *TaskServiceTestSuite) TestGetTaskStatus() {
	expectedTask := &domain.Task{ID: 1, Title: "Test Task"}
	s.mockRepository.EXPECT().
		FindByID(uint(1)).
		Return(expectedTask, nil)

	task, err := s.service.GetTaskStatus(1)
	s.NoError(err)
	s.Equal(expectedTask, task)
}

func (s *TaskServiceTestSuite) TestGetAllTasks() {
	expectedTasks := []*domain.Task{
		{ID: 1, Title: "Task 1"},
		{ID: 2, Title: "Task 2"},
	}
	s.mockRepository.EXPECT().
		FindAll().
		Return(expectedTasks, nil)

	tasks, err := s.service.GetAllTasks()
	s.NoError(err)
	s.Equal(expectedTasks, tasks)
} 