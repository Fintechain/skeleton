// Package system provides the implementation of the system service.
package system

import (
	"fmt"
	"sync"
	"time"

	"github.com/ebanfa/skeleton/internal/domain/component"
	"github.com/ebanfa/skeleton/internal/domain/operation"
	"github.com/ebanfa/skeleton/internal/domain/plugin"
	"github.com/ebanfa/skeleton/internal/domain/service"
	"github.com/ebanfa/skeleton/internal/domain/storage"
	"github.com/ebanfa/skeleton/internal/domain/system"
	"github.com/ebanfa/skeleton/internal/infrastructure/config"
	"github.com/ebanfa/skeleton/internal/infrastructure/event"
	"github.com/ebanfa/skeleton/internal/infrastructure/logging"
)

// DefaultSystemService implements the SystemService interface
type DefaultSystemService struct {
	*service.BaseService
	registry      component.Registry
	pluginManager plugin.PluginManager
	eventBus      event.EventBus
	configuration config.Configuration
	multiStore    storage.MultiStore
	logger        logging.Logger
	mutex         sync.RWMutex
	initialized   bool
	started       bool
}

// NewDefaultSystemService creates a new instance of DefaultSystemService
func NewDefaultSystemService(
	id string,
	registry component.Registry,
	pluginManager plugin.PluginManager,
	eventBus event.EventBus,
	configuration config.Configuration,
	multiStore storage.MultiStore,
	logger logging.Logger,
) *DefaultSystemService {
	// Create the base component
	baseComponent := component.NewBaseComponent(id, "System Service", component.TypeSystem)

	// Create the service with the base component
	baseService := service.NewBaseService(service.BaseServiceOptions{
		Component: baseComponent,
	})

	if logger == nil {
		logger = logging.CreateStandardLogger(logging.Info)
	}

	return &DefaultSystemService{
		BaseService:   baseService,
		registry:      registry,
		pluginManager: pluginManager,
		eventBus:      eventBus,
		configuration: configuration,
		multiStore:    multiStore,
		logger:        logger,
		initialized:   false,
		started:       false,
	}
}

// Registry returns the component registry
func (s *DefaultSystemService) Registry() component.Registry {
	return s.registry
}

// PluginManager returns the plugin manager
func (s *DefaultSystemService) PluginManager() plugin.PluginManager {
	return s.pluginManager
}

// EventBus returns the event bus
func (s *DefaultSystemService) EventBus() event.EventBus {
	return s.eventBus
}

// Configuration returns the system configuration
func (s *DefaultSystemService) Configuration() config.Configuration {
	return s.configuration
}

// Store returns the multi-store
func (s *DefaultSystemService) Store() storage.MultiStore {
	return s.multiStore
}

// Initialize initializes the system service
func (s *DefaultSystemService) Initialize(ctx component.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Call the base implementation first
	if err := s.BaseService.Component.Initialize(ctx); err != nil {
		return err
	}

	// Initialize the registry
	if err := s.registry.Initialize(ctx); err != nil {
		return component.NewError(
			system.ErrSystemNotInitialized,
			"Failed to initialize component registry",
			err,
		)
	}

	// Initialize the multi-store if it implements Component
	if storeComponent, ok := s.multiStore.(component.Component); ok {
		if err := storeComponent.Initialize(ctx); err != nil {
			return component.NewError(
				system.ErrSystemNotInitialized,
				"Failed to initialize multi-store",
				err,
			)
		}
	}

	s.initialized = true

	// Release the lock before publishing events to avoid potential deadlocks
	s.mutex.Unlock()

	// Publish the system initialized event
	s.eventBus.Publish(system.TopicSystemInitialized, map[string]interface{}{
		"serviceId": s.ID(),
		"time":      time.Now(),
	})

	// Re-acquire the lock to maintain the deferred unlock semantics
	s.mutex.Lock()

	return nil
}

// Start starts the system service
func (s *DefaultSystemService) Start(ctx component.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.initialized {
		return component.NewError(
			system.ErrSystemNotInitialized,
			"System service not initialized",
			nil,
		)
	}

	// Call the base implementation first
	if err := s.BaseService.Start(ctx); err != nil {
		return err
	}

	// Start the multi-store if it implements Service
	if storeService, ok := s.multiStore.(service.Service); ok {
		if err := storeService.Start(ctx); err != nil {
			return component.NewError(
				system.ErrSystemNotStarted,
				"Failed to start multi-store service",
				err,
			)
		}
	}

	s.started = true

	// Release the lock before publishing events to avoid potential deadlocks
	s.mutex.Unlock()

	// Publish the system started event
	s.eventBus.Publish(system.TopicSystemStarted, map[string]interface{}{
		"serviceId": s.ID(),
		"time":      time.Now(),
	})

	// Re-acquire the lock to maintain the deferred unlock semantics
	s.mutex.Lock()

	return nil
}

// Stop stops the system service
func (s *DefaultSystemService) Stop(ctx component.Context) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.started {
		return component.NewError(
			system.ErrSystemNotStarted,
			"System service not started",
			nil,
		)
	}

	// Stop all registered services
	services := s.registry.FindByType(component.TypeService)
	for _, svc := range services {
		if serviceComp, ok := svc.(service.Service); ok {
			// Skip stopping ourselves
			if serviceComp.ID() == s.ID() {
				continue
			}

			// Stop the service
			if err := serviceComp.Stop(ctx); err != nil {
				// Log the error but continue stopping other services
				s.logger.Error("Error stopping service: %s - %v", serviceComp.ID(), err)
			}
		}
	}

	// Stop the multi-store if it implements Service
	if storeService, ok := s.multiStore.(service.Service); ok {
		if err := storeService.Stop(ctx); err != nil {
			return component.NewError(
				system.ErrServiceStop,
				"Failed to stop multi-store service",
				err,
			)
		}
	}

	// Call the base implementation last
	if err := s.BaseService.Stop(ctx); err != nil {
		return err
	}

	s.started = false

	// Release the lock before publishing events to avoid potential deadlocks
	s.mutex.Unlock()

	// Publish the system stopped event
	s.eventBus.Publish(system.TopicSystemStopped, map[string]interface{}{
		"serviceId": s.ID(),
		"time":      time.Now(),
	})

	// Re-acquire the lock to maintain the deferred unlock semantics
	s.mutex.Lock()

	return nil
}

// ExecuteOperation executes an operation with the specified ID
func (s *DefaultSystemService) ExecuteOperation(ctx component.Context, operationID string, input interface{}) (interface{}, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.started {
		return nil, component.NewError(
			system.ErrSystemNotStarted,
			"System service not started",
			nil,
		)
	}

	// Check if operations are enabled
	enableOperations := s.GetBoolMetadata("enableOperations", true)
	if !enableOperations {
		return nil, component.NewError(
			system.ErrOperationFailed,
			"Operations are disabled in system configuration",
			nil,
		)
	}

	// Get the operation component
	comp, err := s.registry.Get(operationID)
	if err != nil {
		return nil, component.NewError(
			system.ErrOperationNotFound,
			fmt.Sprintf("Operation not found: %s", operationID),
			err,
		)
	}

	// Check if the component is an operation
	op, ok := comp.(operation.Operation)
	if !ok {
		return nil, component.NewError(
			system.ErrOperationNotFound,
			fmt.Sprintf("Component is not an operation: %s", operationID),
			nil,
		)
	}

	// Create a typed input
	typedInput := &system.SystemOperationInput{
		Data:     input,
		Metadata: make(map[string]interface{}),
	}

	// Release the lock while executing the operation
	s.mutex.RUnlock()

	// Execute the operation
	result, err := op.Execute(ctx, typedInput)

	// Re-acquire the lock
	s.mutex.RLock()

	if err != nil {
		// Publish operation failed event
		s.eventBus.Publish(system.TopicOperationFailed, map[string]interface{}{
			"serviceId":   s.ID(),
			"operationId": operationID,
			"error":       err.Error(),
			"time":        time.Now(),
		})

		return nil, component.NewError(
			system.ErrOperationFailed,
			fmt.Sprintf("Operation execution failed: %s", operationID),
			err,
		)
	}

	// Publish operation executed event
	s.eventBus.Publish(system.TopicOperationExecuted, map[string]interface{}{
		"serviceId":   s.ID(),
		"operationId": operationID,
		"time":        time.Now(),
	})

	// Convert result to the expected output format
	output, ok := result.(*system.SystemOperationOutput)
	if !ok {
		// If the operation doesn't return the expected type, wrap it
		output = &system.SystemOperationOutput{
			Data:     result,
			Metadata: make(map[string]interface{}),
		}
	}

	return output, nil
}

// StartService starts a service with the specified ID
func (s *DefaultSystemService) StartService(ctx component.Context, serviceID string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.started {
		return component.NewError(
			system.ErrSystemNotStarted,
			"System service not started",
			nil,
		)
	}

	// Check if services are enabled
	enableServices := s.GetBoolMetadata("enableServices", true)
	if !enableServices {
		return component.NewError(
			system.ErrServiceStart,
			"Services are disabled in system configuration",
			nil,
		)
	}

	// Get the service component
	comp, err := s.registry.Get(serviceID)
	if err != nil {
		return component.NewError(
			system.ErrServiceNotFound,
			fmt.Sprintf("Service not found: %s", serviceID),
			err,
		)
	}

	// Check if the component is a service
	svc, ok := comp.(service.Service)
	if !ok {
		return component.NewError(
			system.ErrServiceNotFound,
			fmt.Sprintf("Component is not a service: %s", serviceID),
			nil,
		)
	}

	// Release the lock while starting the service
	s.mutex.RUnlock()

	// Start the service
	err = svc.Start(ctx)

	// Re-acquire the lock
	s.mutex.RLock()

	if err != nil {
		// Publish service failed event
		s.eventBus.Publish(system.TopicServiceFailed, map[string]interface{}{
			"serviceId":       s.ID(),
			"targetServiceId": serviceID,
			"action":          "start",
			"error":           err.Error(),
			"time":            time.Now(),
		})

		return component.NewError(
			system.ErrServiceStart,
			fmt.Sprintf("Failed to start service: %s", serviceID),
			err,
		)
	}

	// Publish service started event
	s.eventBus.Publish(system.TopicServiceStarted, map[string]interface{}{
		"serviceId":       s.ID(),
		"targetServiceId": serviceID,
		"time":            time.Now(),
	})

	return nil
}

// StopService stops a service with the specified ID
func (s *DefaultSystemService) StopService(ctx component.Context, serviceID string) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if !s.started {
		return component.NewError(
			system.ErrSystemNotStarted,
			"System service not started",
			nil,
		)
	}

	// Check if services are enabled
	enableServices := s.GetBoolMetadata("enableServices", true)
	if !enableServices {
		return component.NewError(
			system.ErrServiceStop,
			"Services are disabled in system configuration",
			nil,
		)
	}

	// Get the service component
	comp, err := s.registry.Get(serviceID)
	if err != nil {
		return component.NewError(
			system.ErrServiceNotFound,
			fmt.Sprintf("Service not found: %s", serviceID),
			err,
		)
	}

	// Check if the component is a service
	svc, ok := comp.(service.Service)
	if !ok {
		return component.NewError(
			system.ErrServiceNotFound,
			fmt.Sprintf("Component is not a service: %s", serviceID),
			nil,
		)
	}

	// Release the lock while stopping the service
	s.mutex.RUnlock()

	// Stop the service
	err = svc.Stop(ctx)

	// Re-acquire the lock
	s.mutex.RLock()

	if err != nil {
		// Publish service failed event
		s.eventBus.Publish(system.TopicServiceFailed, map[string]interface{}{
			"serviceId":       s.ID(),
			"targetServiceId": serviceID,
			"action":          "stop",
			"error":           err.Error(),
			"time":            time.Now(),
		})

		return component.NewError(
			system.ErrServiceStop,
			fmt.Sprintf("Failed to stop service: %s", serviceID),
			err,
		)
	}

	// Publish service stopped event
	s.eventBus.Publish(system.TopicServiceStopped, map[string]interface{}{
		"serviceId":       s.ID(),
		"targetServiceId": serviceID,
		"time":            time.Now(),
	})

	return nil
}

// GetBoolMetadata retrieves a boolean value from metadata with a default value
func (s *DefaultSystemService) GetBoolMetadata(key string, defaultValue bool) bool {
	metadata := s.Metadata()
	if metadata == nil {
		return defaultValue
	}

	value, exists := metadata[key]
	if !exists {
		return defaultValue
	}

	boolValue, ok := value.(bool)
	if !ok {
		return defaultValue
	}

	return boolValue
}
