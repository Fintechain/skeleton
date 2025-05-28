package service

import (
	"sync"

	"github.com/fintechain/skeleton/internal/domain/component"
)

// BaseService provides a basic implementation of the Service interface.
type BaseService struct {
	component.Component
	status ServiceStatus
	mu     sync.RWMutex // Protects status
}

// BaseServiceOptions contains options for creating a BaseService.
type BaseServiceOptions struct {
	Component component.Component
}

// NewBaseService creates a new base service with dependency injection.
func NewBaseService(options BaseServiceOptions) *BaseService {
	return &BaseService{
		Component: options.Component,
		status:    StatusStopped,
	}
}

// CreateBaseService is a factory method for backward compatibility.
func CreateBaseService(comp component.Component) *BaseService {
	return NewBaseService(BaseServiceOptions{
		Component: comp,
	})
}

// Start starts the service.
// This base implementation just changes the status, subclasses should override it.
func (s *BaseService) Start(ctx component.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if we can start from current state
	switch s.status {
	case StatusRunning, StatusStarting:
		// Already started or starting, do nothing
		return nil
	case StatusStopping:
		// Currently stopping, can't start
		return component.NewError(
			ErrServiceStart,
			"cannot start service while it is stopping",
			nil,
		).WithDetail("service_id", s.ID())
	}

	// Change to starting state
	s.status = StatusStarting

	// Here a derived class would do actual startup work

	// If successful, change to running state
	s.status = StatusRunning

	return nil
}

// Stop stops the service.
// This base implementation just changes the status, subclasses should override it.
func (s *BaseService) Stop(ctx component.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if we can stop from current state
	switch s.status {
	case StatusStopped, StatusStopping:
		// Already stopped or stopping, do nothing
		return nil
	case StatusStarting:
		// Currently starting, can't stop cleanly
		return component.NewError(
			ErrServiceStop,
			"cannot stop service while it is starting",
			nil,
		).WithDetail("service_id", s.ID())
	}

	// Change to stopping state
	s.status = StatusStopping

	// Here a derived class would do actual shutdown work

	// If successful, change to stopped state
	s.status = StatusStopped

	return nil
}

// Status returns the current service status.
func (s *BaseService) Status() ServiceStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.status
}

// SetStatus updates the service status.
func (s *BaseService) SetStatus(status ServiceStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.status = status
}

// MarkAsFailed marks the service as failed.
func (s *BaseService) MarkAsFailed() {
	s.SetStatus(StatusFailed)
}
