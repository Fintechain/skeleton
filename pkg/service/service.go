// Package service provides functionality for services in the system.
package service

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/service"
	serviceImpl "github.com/fintechain/skeleton/internal/infrastructure/service"
)

// Re-export service interfaces
type Service = service.Service
type ServiceFactory = service.ServiceFactory

// Re-export service types
type ServiceStatus = service.ServiceStatus
type ServiceConfig = service.ServiceConfig

// Re-export service status constants
const (
	StatusStopped  = service.StatusStopped
	StatusStarting = service.StatusStarting
	StatusRunning  = service.StatusRunning
	StatusStopping = service.StatusStopping
	StatusFailed   = service.StatusFailed
)

// Re-export service error constants
const (
	ErrServiceStart    = service.ErrServiceStart
	ErrServiceStop     = service.ErrServiceStop
	ErrServiceNotFound = service.ErrServiceNotFound
)

// Re-export constructor function
var NewServiceConfig = service.NewServiceConfig

// Factory functions

// NewService creates a new service instance from a base component.
func NewService(baseComponent component.Component) Service {
	return serviceImpl.NewService(baseComponent)
}

// NewServiceFactory creates a new service factory instance.
func NewServiceFactory(componentFactory component.Factory) ServiceFactory {
	return serviceImpl.NewServiceFactory(componentFactory)
}
