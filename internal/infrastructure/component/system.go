package component

import (
	"errors"
	"sync"
	"sync/atomic"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
)

// System implements the component.System interface.
type System struct {
	registry component.Registry
	running  atomic.Bool
	mu       sync.RWMutex
}

// NewSystem creates a new system with the provided registry.
func NewSystem(registry component.Registry) *System {
	return &System{
		registry: registry,
	}
}

// Registry returns the component registry.
func (s *System) Registry() component.Registry {
	return s.registry
}

// ExecuteOperation executes a registered operation component with the given input.
func (s *System) ExecuteOperation(ctx context.Context, operationID component.ComponentID, input component.Input) (component.Output, error) {
	comp, err := s.registry.Get(operationID)
	if err != nil {
		return component.Output{}, err
	}

	operation, ok := comp.(component.Operation)
	if !ok {
		return component.Output{}, errors.New(component.ErrInvalidComponentType)
	}

	return operation.Execute(ctx, input)
}

// StartService starts a registered service component.
func (s *System) StartService(ctx context.Context, serviceID component.ComponentID) error {
	comp, err := s.registry.Get(serviceID)
	if err != nil {
		return err
	}

	service, ok := comp.(component.Service)
	if !ok {
		return errors.New(component.ErrInvalidComponentType)
	}

	return service.Start(ctx)
}

// StopService stops a running service component gracefully.
func (s *System) StopService(ctx context.Context, serviceID component.ComponentID) error {
	comp, err := s.registry.Get(serviceID)
	if err != nil {
		return err
	}

	service, ok := comp.(component.Service)
	if !ok {
		return errors.New(component.ErrInvalidComponentType)
	}

	return service.Stop(ctx)
}

// Start initializes and starts the entire system.
func (s *System) Start(ctx context.Context) error {
	if s.running.Load() {
		return nil // Already running
	}

	s.running.Store(true)
	return nil
}

// Stop gracefully shuts down the entire system.
func (s *System) Stop(ctx context.Context) error {
	if !s.running.Load() {
		return nil // Already stopped
	}

	s.running.Store(false)
	return nil
}

// IsRunning returns whether the system is currently running.
func (s *System) IsRunning() bool {
	return s.running.Load()
}
