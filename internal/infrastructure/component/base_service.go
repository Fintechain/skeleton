package component

import (
	"sync"
	"sync/atomic"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
)

// BaseService provides common service functionality that can be embedded
// in concrete service implementations.
type BaseService struct {
	*BaseComponent
	status  component.ServiceStatus
	running atomic.Bool
	mu      sync.RWMutex
}

// NewBaseService creates a new base service with the provided configuration.
func NewBaseService(config component.ComponentConfig) *BaseService {
	config.Type = component.TypeService
	return &BaseService{
		BaseComponent: NewBaseComponent(config),
		status:        component.StatusStopped,
	}
}

// Type returns the service component type.
func (s *BaseService) Type() component.ComponentType {
	return component.TypeService
}

// Initialize prepares the service for use within the system.
func (s *BaseService) Initialize(ctx context.Context, system component.System) error {
	return s.BaseComponent.Initialize(ctx, system)
}

// Start begins the service operation.
// Base implementation just sets status - override in concrete implementations.
func (s *BaseService) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running.Load() {
		return nil // Already running
	}

	s.status = component.StatusStarting
	s.running.Store(true)
	s.status = component.StatusRunning
	return nil
}

// Stop ends the service operation.
// Base implementation just sets status - override in concrete implementations.
func (s *BaseService) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running.Load() {
		return nil // Already stopped
	}

	s.status = component.StatusStopping
	s.running.Store(false)
	s.status = component.StatusStopped
	return nil
}

// IsRunning returns true if the service is currently running.
func (s *BaseService) IsRunning() bool {
	return s.running.Load()
}

// Status returns the current service status.
func (s *BaseService) Status() component.ServiceStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status
}

// SetStatus updates the service status (protected method for subclasses).
func (s *BaseService) SetStatus(status component.ServiceStatus) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status = status
}

// SetRunning updates the running state (protected method for subclasses).
func (s *BaseService) SetRunning(running bool) {
	s.running.Store(running)
}
