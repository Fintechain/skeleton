// Package service provides functionality for services in the system.
package service

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
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
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	Status() ServiceStatus
}

// ServiceConfig defines the configuration for creating a service.
type ServiceConfig struct {
	component.ComponentConfig
	// Service-specific configuration properties can be added here
}

// NewServiceConfig creates a new ServiceConfig with the given parameters.
// This ensures the component type is set to TypeService.
func NewServiceConfig(id, name, description string) ServiceConfig {
	return ServiceConfig{
		ComponentConfig: component.NewComponentConfig(id, name, component.TypeService, description),
	}
}

// ServiceFactory creates services from configuration.
type ServiceFactory interface {
	component.Factory

	// CreateService creates a service from the given configuration.
	CreateService(config ServiceConfig) (Service, error)
}

// Common error codes for service operations
const (
	ErrServiceStart    = "service.start_failed"
	ErrServiceStop     = "service.stop_failed"
	ErrServiceNotFound = "service.not_found"
)
