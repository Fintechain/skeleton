package mocks

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/service"
)

// MockService is a mock implementation of the Service interface
type MockService struct {
	*MockComponent
	StartFunc  func(ctx component.Context) error
	StopFunc   func(ctx component.Context) error
	StatusFunc func() service.ServiceStatus
}

// NewMockService creates a new MockService instance
func NewMockService(id, name string) *MockService {
	return &MockService{
		MockComponent: NewMockComponent(id, name, component.TypeService),
		StartFunc:     func(ctx component.Context) error { return nil },
		StopFunc:      func(ctx component.Context) error { return nil },
		StatusFunc:    func() service.ServiceStatus { return "running" },
	}
}

// Start implements the Service interface
func (m *MockService) Start(ctx component.Context) error {
	if m.StartFunc != nil {
		return m.StartFunc(ctx)
	}
	return nil
}

// Stop implements the Service interface
func (m *MockService) Stop(ctx component.Context) error {
	if m.StopFunc != nil {
		return m.StopFunc(ctx)
	}
	return nil
}

// Status implements the Service interface
func (m *MockService) Status() service.ServiceStatus {
	if m.StatusFunc != nil {
		return m.StatusFunc()
	}
	return "running"
}
