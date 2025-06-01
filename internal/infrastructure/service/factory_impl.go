// Package service provides concrete implementations for service factory functionality.
package service

import (
	"fmt"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/service"
)

// ServiceFactory provides a concrete implementation of the ServiceFactory interface.
type ServiceFactory struct {
	componentFactory component.Factory
}

// NewServiceFactory creates a new service factory instance.
// This constructor accepts a component factory interface dependency for component creation.
func NewServiceFactory(componentFactory component.Factory) service.ServiceFactory {
	if componentFactory == nil {
		return nil
	}

	return &ServiceFactory{
		componentFactory: componentFactory,
	}
}

// Create creates a component from the given configuration.
// This method implements the Factory interface by delegating to the component factory.
func (f *ServiceFactory) Create(config component.ComponentConfig) (component.Component, error) {
	if f.componentFactory == nil {
		return nil, fmt.Errorf("component factory is required for service creation")
	}

	// Ensure the component type is set to Service if not already set
	if config.Type != component.TypeService {
		config.Type = component.TypeService
	}

	// Delegate to the component factory for basic component creation
	baseComponent, err := f.componentFactory.Create(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create base component for service: %w", err)
	}

	// Wrap the base component in a service implementation
	svc := NewService(baseComponent)
	if svc == nil {
		return nil, fmt.Errorf("failed to create service from component")
	}

	return svc, nil
}

// CreateService creates a service from the given service configuration.
// This method provides service-specific creation functionality.
func (f *ServiceFactory) CreateService(config service.ServiceConfig) (service.Service, error) {
	if f.componentFactory == nil {
		return nil, fmt.Errorf("component factory is required for service creation")
	}

	// Validate service configuration
	if config.ID == "" {
		return nil, fmt.Errorf("service ID cannot be empty")
	}
	if config.Name == "" {
		return nil, fmt.Errorf("service name cannot be empty")
	}

	// Ensure the component type is set to Service
	config.ComponentConfig.Type = component.TypeService

	// Create the base component using the component factory
	baseComponent, err := f.componentFactory.Create(config.ComponentConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create base component for service '%s': %w", config.ID, err)
	}

	// Create the service by wrapping the base component
	svc := NewService(baseComponent)
	if svc == nil {
		return nil, fmt.Errorf("failed to create service '%s' from component", config.ID)
	}

	return svc, nil
}

// CreateServiceWithValidation creates a service with additional validation.
// This method provides extended creation functionality with custom validation.
func (f *ServiceFactory) CreateServiceWithValidation(config service.ServiceConfig, validator func(service.ServiceConfig) error) (service.Service, error) {
	if validator != nil {
		if err := validator(config); err != nil {
			return nil, fmt.Errorf("service configuration validation failed: %w", err)
		}
	}

	return f.CreateService(config)
}

// CreateManagedService creates a service with automatic lifecycle management.
// This method provides enhanced service creation with lifecycle callbacks.
func (f *ServiceFactory) CreateManagedService(config service.ServiceConfig, onStart func() error, onStop func() error) (service.Service, error) {
	// Create the base service
	svc, err := f.CreateService(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create managed service: %w", err)
	}

	// If the service supports lifecycle awareness, register callbacks
	if baseService, ok := svc.(*BaseService); ok {
		if lifecycleAware, ok := baseService.Component.(component.LifecycleAwareComponent); ok {
			// Register state change callbacks for managed lifecycle
			lifecycleAware.OnStateChange(func(oldState, newState component.LifecycleState) {
				switch newState {
				case component.StateActive:
					if onStart != nil {
						onStart()
					}
				case component.StateDisposed:
					if onStop != nil {
						onStop()
					}
				}
			})
		}
	}

	return svc, nil
}
