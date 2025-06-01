// Package system provides the core system interface for resource access and operations.
package system

import (
	"github.com/fintechain/skeleton/internal/domain/config"
	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/event"
	"github.com/fintechain/skeleton/internal/domain/plugin"
	"github.com/fintechain/skeleton/internal/domain/registry"
	"github.com/fintechain/skeleton/internal/domain/storage"
)

// System provides access to all system resources and operations.
// Key insight: System doesn't depend on Component interface, breaking circular dependency.
type System interface {
	// Core components access
	Registry() registry.Registry
	PluginManager() plugin.PluginManager
	EventBus() event.EventBus
	Configuration() config.Configuration
	Store() storage.MultiStore

	// System operations (using IDs, not component interfaces)
	ExecuteOperation(ctx context.Context, operationID string, input interface{}) (interface{}, error)
	StartService(ctx context.Context, serviceID string) error
	StopService(ctx context.Context, serviceID string) error

	// System state
	IsRunning() bool
	IsInitialized() bool
}

// Common system error codes
const (
	ErrSystemNotInitialized = "system.not_initialized"
	ErrSystemNotStarted     = "system.not_started"
	ErrOperationNotFound    = "system.operation_not_found"
	ErrOperationFailed      = "system.operation_failed"
	ErrServiceNotFound      = "system.service_not_found"
	ErrServiceStart         = "system.service_start_failed"
	ErrServiceStop          = "system.service_stop_failed"
)
