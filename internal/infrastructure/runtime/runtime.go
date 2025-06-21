// Package runtime provides the runtime environment implementation.
package runtime

import (
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/config"
	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/event"
	"github.com/fintechain/skeleton/internal/domain/logging"
	"github.com/fintechain/skeleton/internal/domain/plugin"
)

// Error constants
const (
	ErrNilRegistry      = "runtime.nil_registry"
	ErrNilConfiguration = "runtime.nil_configuration"
	ErrNilPluginManager = "runtime.nil_plugin_manager"
	ErrNilEventBus      = "runtime.nil_event_bus"
	ErrNilLogger        = "runtime.nil_logger"
)

// Runtime implements the RuntimeEnvironment interface directly.
type Runtime struct {
	// Direct dependency injection
	registry      component.Registry
	config        config.Configuration
	pluginManager plugin.PluginManager
	eventBus      event.EventBusService
	logger        logging.LoggerService

	// State
	running atomic.Bool
}

// NewRuntime creates a new runtime environment with direct dependency injection.
func NewRuntime(
	registry component.Registry,
	config config.Configuration,
	pluginManager plugin.PluginManager,
	eventBus event.EventBusService,
	logger logging.LoggerService,
) (*Runtime, error) {
	if registry == nil {
		return nil, errors.New(ErrNilRegistry)
	}
	if config == nil {
		return nil, errors.New(ErrNilConfiguration)
	}
	if pluginManager == nil {
		return nil, errors.New(ErrNilPluginManager)
	}
	if eventBus == nil {
		return nil, errors.New(ErrNilEventBus)
	}
	if logger == nil {
		return nil, errors.New(ErrNilLogger)
	}

	return &Runtime{
		registry:      registry,
		config:        config,
		pluginManager: pluginManager,
		eventBus:      eventBus,
		logger:        logger,
	}, nil
}

// Registry returns the component registry.
func (r *Runtime) Registry() component.Registry {
	return r.registry
}

// ExecuteOperation executes a registered operation component with the given input.
func (r *Runtime) ExecuteOperation(ctx context.Context, operationID component.ComponentID, input component.Input) (component.Output, error) {
	comp, err := r.registry.Get(operationID)
	if err != nil {
		return component.Output{}, fmt.Errorf("operation not found: %w", err)
	}

	operation, ok := comp.(component.Operation)
	if !ok {
		return component.Output{}, fmt.Errorf("component %s is not an operation", operationID)
	}

	return operation.Execute(ctx, input)
}

// StartService starts a registered service component.
func (r *Runtime) StartService(ctx context.Context, serviceID component.ComponentID) error {
	comp, err := r.registry.Get(serviceID)
	if err != nil {
		return fmt.Errorf("service not found: %w", err)
	}

	service, ok := comp.(component.Service)
	if !ok {
		return fmt.Errorf("component %s is not a service", serviceID)
	}

	return service.Start(ctx)
}

// StopService stops a running service component gracefully.
func (r *Runtime) StopService(ctx context.Context, serviceID component.ComponentID) error {
	comp, err := r.registry.Get(serviceID)
	if err != nil {
		return fmt.Errorf("service not found: %w", err)
	}

	service, ok := comp.(component.Service)
	if !ok {
		return fmt.Errorf("component %s is not a service", serviceID)
	}

	return service.Stop(ctx)
}

// Start initializes and starts the entire system.
func (r *Runtime) Start(ctx context.Context) error {
	if r.running.Load() {
		return nil // Already running
	}

	// Start core services
	if err := r.eventBus.Start(ctx); err != nil {
		return fmt.Errorf("failed to start event bus: %w", err)
	}

	if err := r.pluginManager.Start(ctx); err != nil {
		return fmt.Errorf("failed to start plugin manager: %w", err)
	}

	if err := r.logger.Start(ctx); err != nil {
		return fmt.Errorf("failed to start logger: %w", err)
	}

	r.running.Store(true)
	return nil
}

// Stop gracefully shuts down the entire system.
func (r *Runtime) Stop(ctx context.Context) error {
	if !r.running.Load() {
		return nil // Already stopped
	}

	// Stop core services in reverse order
	if err := r.logger.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop logger: %w", err)
	}

	if err := r.pluginManager.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop plugin manager: %w", err)
	}

	if err := r.eventBus.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop event bus: %w", err)
	}

	r.running.Store(false)
	return nil
}

// IsRunning returns whether the system is currently running.
func (r *Runtime) IsRunning() bool {
	return r.running.Load()
}

// PluginManager returns the system's plugin manager service.
func (r *Runtime) PluginManager() plugin.PluginManager {
	return r.pluginManager
}

// EventBus returns the system's event bus service.
func (r *Runtime) EventBus() event.EventBusService {
	return r.eventBus
}

// Logger returns the system's logger.
func (r *Runtime) Logger() logging.Logger {
	return r.logger
}

// LoadPlugins loads multiple plugins into the system.
func (r *Runtime) LoadPlugins(ctx context.Context, plugins []plugin.Plugin) error {
	// Add each plugin to the manager
	for _, p := range plugins {
		// Add plugin to manager
		if err := r.pluginManager.Add(p.ID(), p); err != nil {
			return fmt.Errorf("failed to add plugin %s: %w", p.ID(), err)
		}
	}

	// Let the plugin manager initialize all added plugins
	if err := r.pluginManager.Initialize(ctx, r); err != nil {
		return fmt.Errorf("failed to initialize plugin manager: %w", err)
	}

	return nil
}

// Configuration returns the runtime environment configuration.
func (r *Runtime) Configuration() config.Configuration {
	return r.config
}
