// Code generated by MockGen. DO NOT EDIT.
// Source: golangwithgin/internal/domain (interfaces: TaskProcessor)

// Package mocks is a generated GoMock package.
package mocks

import (
	domain "golangwithgin/internal/domain"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTaskProcessor is a mock of TaskProcessor interface.
type MockTaskProcessor struct {
	ctrl     *gomock.Controller
	recorder *MockTaskProcessorMockRecorder
}

// MockTaskProcessorMockRecorder is the mock recorder for MockTaskProcessor.
type MockTaskProcessorMockRecorder struct {
	mock *MockTaskProcessor
}

// NewMockTaskProcessor creates a new mock instance.
func NewMockTaskProcessor(ctrl *gomock.Controller) *MockTaskProcessor {
	mock := &MockTaskProcessor{ctrl: ctrl}
	mock.recorder = &MockTaskProcessorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskProcessor) EXPECT() *MockTaskProcessorMockRecorder {
	return m.recorder
}

// Process mocks base method.
func (m *MockTaskProcessor) Process(arg0 *domain.Task) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Process", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Process indicates an expected call of Process.
func (mr *MockTaskProcessorMockRecorder) Process(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Process", reflect.TypeOf((*MockTaskProcessor)(nil).Process), arg0)
}

// Shutdown mocks base method.
func (m *MockTaskProcessor) Shutdown() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Shutdown")
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockTaskProcessorMockRecorder) Shutdown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockTaskProcessor)(nil).Shutdown))
}
