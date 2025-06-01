// Package service provides concrete implementations for service functionality.
package service

import (
	"fmt"
	"sync"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/service"
)

// BaseService provides a concrete implementation of the Service interface.
type BaseService struct {
	component.Component
	status service.ServiceStatus
	mu     sync.RWMutex
}

// NewService creates a new service instance.
// This constructor accepts a component interface dependency for composition.
func NewService(baseComponent component.Component) service.Service {
	if baseComponent == nil {
		return nil
	}

	return &BaseService{
		Component: baseComponent,
		status:    service.StatusStopped, // Initial status
	}
}

// Start starts the service with the given context.
// This method provides thread-safe service startup functionality.
func (s *BaseService) Start(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("context is required for service startup")
	}

	// Check if the service component is properly initialized
	if s.Component == nil {
		return fmt.Errorf(service.ErrServiceNotFound + ": service component not available")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check current status to prevent invalid transitions
	switch s.status {
	case service.StatusRunning:
		return nil // Already running, no error
	case service.StatusStarting:
		return fmt.Errorf(service.ErrServiceStart + ": service is already starting")
	case service.StatusStopping:
		return fmt.Errorf(service.ErrServiceStart + ": service is currently stopping")
	}

	// Set status to starting
	s.status = service.StatusStarting

	// Update component lifecycle state if supported
	if lifecycleAware, ok := s.Component.(component.LifecycleAwareComponent); ok {
		lifecycleAware.SetState(component.StateInitializing)
	}

	// Perform service startup logic
	// For the base implementation, we simulate successful startup
	// Specific service implementations would override this with actual startup logic
	if err := s.performStartup(ctx); err != nil {
		s.status = service.StatusFailed
		if lifecycleAware, ok := s.Component.(component.LifecycleAwareComponent); ok {
			lifecycleAware.SetState(component.StateFailed)
		}
		return fmt.Errorf(service.ErrServiceStart+": %w", err)
	}

	// Set status to running on successful startup
	s.status = service.StatusRunning
	if lifecycleAware, ok := s.Component.(component.LifecycleAwareComponent); ok {
		lifecycleAware.SetState(component.StateActive)
	}

	return nil
}

// Stop stops the service with the given context.
// This method provides thread-safe service shutdown functionality.
func (s *BaseService) Stop(ctx context.Context) error {
	if ctx == nil {
		return fmt.Errorf("context is required for service shutdown")
	}

	// Check if the service component is properly initialized
	if s.Component == nil {
		return fmt.Errorf(service.ErrServiceNotFound + ": service component not available")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Check current status to prevent invalid transitions
	switch s.status {
	case service.StatusStopped:
		return nil // Already stopped, no error
	case service.StatusStopping:
		return fmt.Errorf(service.ErrServiceStop + ": service is already stopping")
	case service.StatusStarting:
		return fmt.Errorf(service.ErrServiceStop + ": service is currently starting")
	}

	// Set status to stopping
	s.status = service.StatusStopping

	// Update component lifecycle state if supported
	if lifecycleAware, ok := s.Component.(component.LifecycleAwareComponent); ok {
		lifecycleAware.SetState(component.StateDisposing)
	}

	// Perform service shutdown logic
	// For the base implementation, we simulate successful shutdown
	// Specific service implementations would override this with actual shutdown logic
	if err := s.performShutdown(ctx); err != nil {
		s.status = service.StatusFailed
		if lifecycleAware, ok := s.Component.(component.LifecycleAwareComponent); ok {
			lifecycleAware.SetState(component.StateFailed)
		}
		return fmt.Errorf(service.ErrServiceStop+": %w", err)
	}

	// Set status to stopped on successful shutdown
	s.status = service.StatusStopped
	if lifecycleAware, ok := s.Component.(component.LifecycleAwareComponent); ok {
		lifecycleAware.SetState(component.StateDisposed)
	}

	return nil
}

// Status returns the current status of the service.
// This method provides thread-safe status access.
func (s *BaseService) Status() service.ServiceStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.status
}

// performStartup performs the actual service startup logic.
// This method can be overridden by specific service implementations.
func (s *BaseService) performStartup(ctx context.Context) error {
	// Base implementation performs no specific startup logic
	// Specific service implementations would override this method
	return nil
}

// performShutdown performs the actual service shutdown logic.
// This method can be overridden by specific service implementations.
func (s *BaseService) performShutdown(ctx context.Context) error {
	// Base implementation performs no specific shutdown logic
	// Specific service implementations would override this method
	return nil
}

// IsRunning returns true if the service is currently running.
// This is a convenience method for status checking.
func (s *BaseService) IsRunning() bool {
	return s.Status() == service.StatusRunning
}

// IsStopped returns true if the service is currently stopped.
// This is a convenience method for status checking.
func (s *BaseService) IsStopped() bool {
	return s.Status() == service.StatusStopped
}
