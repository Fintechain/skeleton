// Package component provides interfaces and types for the component system.
package component

import (
	"github.com/fintechain/skeleton/internal/domain/context"
)

// Service represents a long-running component that can be started and stopped.
// Services maintain persistent state and run continuously until explicitly stopped.
type Service interface {
	Component

	// Start begins the service operation
	Start(ctx context.Context) error

	// Stop ends the service operation
	Stop(ctx context.Context) error

	// IsRunning returns true if the service is currently running
	IsRunning() bool

	// Status returns the current service status
	Status() ServiceStatus
}

// ServiceStatus represents the current state of a service
type ServiceStatus string

const (
	StatusStopped  ServiceStatus = "stopped"
	StatusStarting ServiceStatus = "starting"
	StatusRunning  ServiceStatus = "running"
	StatusStopping ServiceStatus = "stopping"
	StatusError    ServiceStatus = "error"
)
