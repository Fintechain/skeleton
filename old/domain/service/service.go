// Package service provides functionality for services in the system.
package service

import (
	"github.com/fintechain/skeleton/internal/domain/component"
)

// ServiceStatus represents the status of a service.
type ServiceStatus string

const (
	// Service status constants
	StatusStopped  ServiceStatus = "stopped"
	StatusStarting ServiceStatus = "starting"
	StatusRunning  ServiceStatus = "running"
	StatusStopping ServiceStatus = "stopping"
	StatusFailed   ServiceStatus = "failed"
)

// Service is a specialized component providing ongoing functionality.
type Service interface {
	component.Component

	// Service lifecycle
	Start(ctx component.Context) error
	Stop(ctx component.Context) error
	Status() ServiceStatus
}

// ServiceConfig defines the configuration for creating a service.
type ServiceConfig struct {
	component.ComponentConfig
	// Service-specific configuration properties
}

// ServiceFactory creates services from configuration.
type ServiceFactory interface {
	component.Factory

	// CreateService creates a service from the given configuration.
	CreateService(config ServiceConfig) (Service, error)
}

// DefaultServiceFactory provides a standard implementation of ServiceFactory.
type DefaultServiceFactory struct {
	factory component.Factory
}

// NewServiceFactory creates a new service factory.
func NewServiceFactory() ServiceFactory {
	return &DefaultServiceFactory{
		factory: component.CreateFactory(),
	}
}

// Create creates a component from configuration (implements component.Factory).
func (f *DefaultServiceFactory) Create(config component.ComponentConfig) (component.Component, error) {
	return f.factory.Create(config)
}

// CreateService creates a service from the given configuration.
func (f *DefaultServiceFactory) CreateService(config ServiceConfig) (Service, error) {
	// Create a component from the configuration
	comp, err := f.factory.Create(config.ComponentConfig)
	if err != nil {
		return nil, component.NewError(
			"service.factory.create_failed",
			"failed to create component for service",
			err,
		).WithDetail("config", config)
	}

	// Create a default service with the component
	service := CreateDefaultService(comp)
	return service, nil
}

// NewService creates a new service with the given options.
func NewService(options DefaultServiceOptions) Service {
	return NewDefaultService(options)
}

// CreateService creates a service from a component (convenience function).
func CreateService(comp component.Component) Service {
	return CreateDefaultService(comp)
}

// Common error codes for service operations
const (
	ErrServiceStart    = "service.start_failed"
	ErrServiceStop     = "service.stop_failed"
	ErrServiceNotFound = "service.not_found"
)
