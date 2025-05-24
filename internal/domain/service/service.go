// Package service provides functionality for services in the system.
package service

import (
	"github.com/ebanfa/skeleton/internal/domain/component"
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

// Common error codes for service operations
const (
	ErrServiceStart    = "service.start_failed"
	ErrServiceStop     = "service.stop_failed"
	ErrServiceNotFound = "service.not_found"
)
