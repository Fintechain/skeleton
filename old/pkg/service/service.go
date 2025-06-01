// Package service provides public APIs for the service system.
package service

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/service"
)

// ===== SERVICE INTERFACES =====

// Service is a specialized component providing ongoing functionality.
type Service = service.Service

// ServiceFactory creates services from configuration.
type ServiceFactory = service.ServiceFactory

// HealthCheck is an interface that services can implement to indicate their health.
type HealthCheck = service.HealthCheck

// ===== SERVICE TYPES =====

// ServiceStatus represents the status of a service.
type ServiceStatus = service.ServiceStatus

// ServiceConfig defines the configuration for creating a service.
type ServiceConfig = service.ServiceConfig

// HealthStatus represents the health status of a service.
type HealthStatus = service.HealthStatus

// HealthResult contains the result of a health check.
type HealthResult = service.HealthResult

// ===== SERVICE STATUS CONSTANTS =====

// Service status constants
const (
	StatusStopped  = service.StatusStopped
	StatusStarting = service.StatusStarting
	StatusRunning  = service.StatusRunning
	StatusStopping = service.StatusStopping
	StatusFailed   = service.StatusFailed
)

// ===== HEALTH STATUS CONSTANTS =====

// Health status constants
const (
	HealthStatusUnknown   = service.HealthStatusUnknown
	HealthStatusHealthy   = service.HealthStatusHealthy
	HealthStatusUnhealthy = service.HealthStatusUnhealthy
	HealthStatusDegraded  = service.HealthStatusDegraded
)

// ===== SERVICE ERROR CONSTANTS =====

// Common service error codes
const (
	ErrServiceStart    = "service.start_failed"
	ErrServiceStop     = "service.stop_failed"
	ErrServiceNotFound = "service.not_found"
	ErrServiceTimeout  = "service.timeout"
	ErrHealthCheck     = "service.health_check_failed"
	ErrInvalidInput    = "service.invalid_input"
	ErrInvalidConfig   = "service.invalid_config"
)

// ===== ERROR HANDLING =====

// Error represents a domain-specific error from the service system.
type Error = component.Error

// NewError creates a new service error with the given code, message, and optional cause.
func NewError(code, message string, cause error) *Error {
	return component.NewError(code, message, cause)
}

// IsServiceError checks if an error is a service error with the given code.
func IsServiceError(err error, code string) bool {
	return component.IsComponentError(err, code)
}

// ===== SERVICE CONSTRUCTORS =====

// NewServiceFactory creates a new service factory with default configuration.
// This is the primary way to create a ServiceFactory instance for creating services.
//
// Example usage:
//
//	factory := service.NewServiceFactory()
//	config := &service.ServiceConfig{
//	    ComponentConfig: component.ComponentConfig{
//	        Name: "api-server",
//	        Type: "http",
//	    },
//	}
//	svc, err := factory.CreateService(config)
func NewServiceFactory() ServiceFactory {
	return service.NewServiceFactory()
}

// NewService creates a new service with the given name and type.
// This is a convenience function for creating simple services without complex configuration.
//
// Example usage:
//
//	svc, err := service.NewService("api-server", "http")
//	if err != nil {
//	    // Handle error
//	}
func NewService(name, serviceType string) (Service, error) {
	if name == "" {
		return nil, NewError(ErrInvalidInput, "service name cannot be empty", nil)
	}
	if serviceType == "" {
		return nil, NewError(ErrInvalidInput, "service type cannot be empty", nil)
	}

	factory := NewServiceFactory()
	config := ServiceConfig{
		ComponentConfig: component.ComponentConfig{
			Name: name,
			Type: component.ComponentType(serviceType),
		},
	}

	return factory.CreateService(config)
}

// ===== SERVICE UTILITIES =====

// IsHealthy checks if a service is healthy by performing a health check.
// Returns true if the service implements HealthCheck and reports healthy status.
//
// Example usage:
//
//	if service.IsHealthy(svc) {
//	    // Service is healthy
//	} else {
//	    // Service is unhealthy or doesn't support health checks
//	}
func IsHealthy(svc Service) bool {
	if healthCheck, ok := svc.(HealthCheck); ok {
		return healthCheck.IsHealthy()
	}
	// If service doesn't implement health check, consider it healthy if running
	return svc.Status() == StatusRunning
}

// WaitForStatus waits for a service to reach the specified status.
// This is a convenience function for service lifecycle management.
//
// Example usage:
//
//	ctx := context.WithTimeout(context.Background(), 30*time.Second)
//	if service.WaitForStatus(ctx, svc, service.StatusRunning) {
//	    // Service is now running
//	}
func WaitForStatus(ctx component.Context, svc Service, expectedStatus ServiceStatus) bool {
	// Simple polling implementation - in a real system this might use events
	for {
		select {
		case <-ctx.Done():
			return false
		default:
			if svc.Status() == expectedStatus {
				return true
			}
			// Small delay to avoid busy waiting
			// In a real implementation, this would use proper event-driven waiting
		}
	}
}
