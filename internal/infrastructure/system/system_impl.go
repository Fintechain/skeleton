package system

import (
	"errors"
	"fmt"
	"sync"

	"github.com/fintechain/skeleton/internal/domain/config"
	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/event"
	"github.com/fintechain/skeleton/internal/domain/operation"
	"github.com/fintechain/skeleton/internal/domain/plugin"
	"github.com/fintechain/skeleton/internal/domain/registry"
	"github.com/fintechain/skeleton/internal/domain/service"
	"github.com/fintechain/skeleton/internal/domain/storage"
	"github.com/fintechain/skeleton/internal/domain/system"
)

// DefaultSystem provides a concrete implementation of the System interface.
type DefaultSystem struct {
	registry      registry.Registry
	pluginManager plugin.PluginManager
	eventBus      event.EventBus
	configuration config.Configuration
	store         storage.MultiStore

	// State management
	mu          sync.RWMutex
	initialized bool
	running     bool
}

// NewSystem creates a new System instance with the provided dependencies.
// This constructor accepts all resource interface dependencies for testability.
func NewSystem(
	registry registry.Registry,
	pluginManager plugin.PluginManager,
	eventBus event.EventBus,
	configuration config.Configuration,
	store storage.MultiStore,
) system.System {
	return &DefaultSystem{
		registry:      registry,
		pluginManager: pluginManager,
		eventBus:      eventBus,
		configuration: configuration,
		store:         store,
		initialized:   true, // System is initialized when all dependencies are provided
		running:       true, // System is running when created with valid dependencies
	}
}

// Registry returns the system's registry.
func (s *DefaultSystem) Registry() registry.Registry {
	return s.registry
}

// PluginManager returns the system's plugin manager.
func (s *DefaultSystem) PluginManager() plugin.PluginManager {
	return s.pluginManager
}

// EventBus returns the system's event bus.
func (s *DefaultSystem) EventBus() event.EventBus {
	return s.eventBus
}

// Configuration returns the system's configuration.
func (s *DefaultSystem) Configuration() config.Configuration {
	return s.configuration
}

// Store returns the system's multi-store.
func (s *DefaultSystem) Store() storage.MultiStore {
	return s.store
}

// ExecuteOperation executes an operation by ID with the given context and input.
func (s *DefaultSystem) ExecuteOperation(ctx context.Context, operationID string, input interface{}) (interface{}, error) {
	if !s.IsInitialized() {
		return nil, errors.New(system.ErrSystemNotInitialized)
	}

	if !s.IsRunning() {
		return nil, errors.New(system.ErrSystemNotStarted)
	}

	// Get the operation from the registry
	item, err := s.registry.Get(operationID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", system.ErrOperationNotFound, err)
	}

	// Cast to operation interface
	op, ok := item.(operation.Operation)
	if !ok {
		return nil, fmt.Errorf("%s: item %s is not an operation", system.ErrOperationNotFound, operationID)
	}

	// Execute the operation
	result, err := op.Execute(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", system.ErrOperationFailed, err)
	}

	return result, nil
}

// StartService starts a service by ID.
func (s *DefaultSystem) StartService(ctx context.Context, serviceID string) error {
	if !s.IsInitialized() {
		return errors.New(system.ErrSystemNotInitialized)
	}

	if !s.IsRunning() {
		return errors.New(system.ErrSystemNotStarted)
	}

	// Get the service from the registry
	item, err := s.registry.Get(serviceID)
	if err != nil {
		return fmt.Errorf("%s: %w", system.ErrServiceNotFound, err)
	}

	// Cast to service interface
	svc, ok := item.(service.Service)
	if !ok {
		return fmt.Errorf("%s: item %s is not a service", system.ErrServiceNotFound, serviceID)
	}

	// Start the service
	if err := svc.Start(ctx); err != nil {
		return fmt.Errorf("%s: %w", system.ErrServiceStart, err)
	}

	return nil
}

// StopService stops a service by ID.
func (s *DefaultSystem) StopService(ctx context.Context, serviceID string) error {
	if !s.IsInitialized() {
		return errors.New(system.ErrSystemNotInitialized)
	}

	if !s.IsRunning() {
		return errors.New(system.ErrSystemNotStarted)
	}

	// Get the service from the registry
	item, err := s.registry.Get(serviceID)
	if err != nil {
		return fmt.Errorf("%s: %w", system.ErrServiceNotFound, err)
	}

	// Cast to service interface
	svc, ok := item.(service.Service)
	if !ok {
		return fmt.Errorf("%s: item %s is not a service", system.ErrServiceNotFound, serviceID)
	}

	// Stop the service
	if err := svc.Stop(ctx); err != nil {
		return fmt.Errorf("%s: %w", system.ErrServiceStop, err)
	}

	return nil
}

// IsRunning returns whether the system is running.
func (s *DefaultSystem) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

// IsInitialized returns whether the system is initialized.
func (s *DefaultSystem) IsInitialized() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.initialized
}
