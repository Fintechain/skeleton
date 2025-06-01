// Package mocks provides centralized mock implementations for the Skeleton Framework.
// This file contains mocks for service interfaces.
package mocks

import (
	"sync"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/service"
)

// MockService implements the service.Service interface for testing.
type MockService struct {
	MockComponent
	mu sync.RWMutex

	// Service-specific fields
	status      service.ServiceStatus
	startFunc   func(ctx context.Context) error
	stopFunc    func(ctx context.Context) error
	startError  error
	stopError   error
	startCalled bool
	stopCalled  bool
}

// NewMockService creates a new mock service with default behavior.
func NewMockService() *MockService {
	mock := &MockService{
		MockComponent: *NewMockComponent(),
		status:        service.StatusStopped,
	}
	mock.SetType(component.TypeService)
	return mock
}

// Start implements the service.Service interface.
func (m *MockService) Start(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.startCalled = true

	if m.startFunc != nil {
		err := m.startFunc(ctx)
		if err == nil {
			m.status = service.StatusRunning
		} else {
			m.status = service.StatusFailed
		}
		return err
	}

	if m.startError != nil {
		m.status = service.StatusFailed
		return m.startError
	}

	m.status = service.StatusRunning
	return nil
}

// Stop implements the service.Service interface.
func (m *MockService) Stop(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stopCalled = true

	if m.stopFunc != nil {
		err := m.stopFunc(ctx)
		if err == nil {
			m.status = service.StatusStopped
		} else {
			m.status = service.StatusFailed
		}
		return err
	}

	if m.stopError != nil {
		m.status = service.StatusFailed
		return m.stopError
	}

	m.status = service.StatusStopped
	return nil
}

// Status implements the service.Service interface.
func (m *MockService) Status() service.ServiceStatus {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.status
}

// Mock configuration methods

// SetStartFunc sets a custom function to be called when Start is invoked.
func (m *MockService) SetStartFunc(fn func(ctx context.Context) error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.startFunc = fn
}

// SetStopFunc sets a custom function to be called when Stop is invoked.
func (m *MockService) SetStopFunc(fn func(ctx context.Context) error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stopFunc = fn
}

// SetStartError sets the error to return from Start calls.
func (m *MockService) SetStartError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.startError = err
}

// SetStopError sets the error to return from Stop calls.
func (m *MockService) SetStopError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stopError = err
}

// SetStatus sets the service status.
func (m *MockService) SetStatus(status service.ServiceStatus) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.status = status
}

// Verification methods

// WasStartCalled returns true if Start was called.
func (m *MockService) WasStartCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.startCalled
}

// WasStopCalled returns true if Stop was called.
func (m *MockService) WasStopCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.stopCalled
}

// MockServiceFactory implements the service.ServiceFactory interface for testing.
type MockServiceFactory struct {
	MockComponentFactory
	mu sync.RWMutex

	// Factory-specific fields
	createServiceFunc   func(config service.ServiceConfig) (service.Service, error)
	createServiceResult service.Service
	createServiceError  error
	createServiceCalled bool
	createServiceConfig service.ServiceConfig
}

// NewMockServiceFactory creates a new mock service factory.
func NewMockServiceFactory() *MockServiceFactory {
	return &MockServiceFactory{
		MockComponentFactory: *NewMockComponentFactory(),
	}
}

// CreateService implements the service.ServiceFactory interface.
func (m *MockServiceFactory) CreateService(config service.ServiceConfig) (service.Service, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.createServiceCalled = true
	m.createServiceConfig = config

	if m.createServiceFunc != nil {
		return m.createServiceFunc(config)
	}

	if m.createServiceError != nil {
		return nil, m.createServiceError
	}

	if m.createServiceResult != nil {
		return m.createServiceResult, nil
	}

	// Default behavior: create a mock service
	mockSvc := NewMockService()
	mockSvc.SetID(config.ID)
	mockSvc.SetName(config.Name)
	mockSvc.SetDescription(config.Description)
	return mockSvc, nil
}

// Mock configuration methods

// SetCreateServiceFunc sets a custom function for CreateService.
func (m *MockServiceFactory) SetCreateServiceFunc(fn func(config service.ServiceConfig) (service.Service, error)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createServiceFunc = fn
}

// SetCreateServiceResult sets the result to return from CreateService.
func (m *MockServiceFactory) SetCreateServiceResult(result service.Service) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createServiceResult = result
}

// SetCreateServiceError sets the error to return from CreateService.
func (m *MockServiceFactory) SetCreateServiceError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createServiceError = err
}

// Verification methods

// WasCreateServiceCalled returns true if CreateService was called.
func (m *MockServiceFactory) WasCreateServiceCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.createServiceCalled
}

// GetCreateServiceConfig returns the config passed to the last CreateService call.
func (m *MockServiceFactory) GetCreateServiceConfig() service.ServiceConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.createServiceConfig
}
