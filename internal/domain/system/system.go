// Package system provides the core system interfaces and types.
package system

import (
	"github.com/ebanfa/skeleton/internal/domain/component"
	"github.com/ebanfa/skeleton/internal/domain/plugin"
	"github.com/ebanfa/skeleton/internal/domain/service"
	"github.com/ebanfa/skeleton/internal/domain/storage"
	"github.com/ebanfa/skeleton/internal/infrastructure/config"
	"github.com/ebanfa/skeleton/internal/infrastructure/event"
)

// SystemService serves as the central coordinating service of the application
type SystemService interface {
	service.Service

	// Core components access
	Registry() component.Registry
	PluginManager() plugin.PluginManager
	EventBus() event.EventBus
	Configuration() config.Configuration
	Store() storage.MultiStore

	// Operations
	ExecuteOperation(ctx component.Context, operationID string, input interface{}) (interface{}, error)
	StartService(ctx component.Context, serviceID string) error
	StopService(ctx component.Context, serviceID string) error
}

// SystemServiceConfig defines configuration for SystemService
type SystemServiceConfig struct {
	// Service identity
	ServiceID string `json:"serviceId" default:"system"`

	// Feature flags
	EnableOperations bool `json:"enableOperations" default:"true"`
	EnableServices   bool `json:"enableServices" default:"true"`
	EnablePlugins    bool `json:"enablePlugins" default:"true"`
	EnableEventLog   bool `json:"enableEventLog" default:"true"`

	// Storage configuration
	StorageConfig storage.MultiStoreConfig `json:"storage"`
}

// SystemOperationInput represents the input for a system operation
type SystemOperationInput struct {
	// The actual input data (usually deserializable to a specific type)
	Data interface{} `json:"data"`

	// Metadata about the operation call
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// SystemOperationOutput represents the output from a system operation
type SystemOperationOutput struct {
	// The actual output data
	Data interface{} `json:"data"`

	// Metadata about the operation result
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Core system event topics
const (
	// System lifecycle events
	TopicSystemInitialized = "system.initialized"
	TopicSystemStarted     = "system.started"
	TopicSystemStopped     = "system.stopped"

	// Operation events
	TopicOperationExecuted = "system.operation.executed"
	TopicOperationFailed   = "system.operation.failed"

	// Service events
	TopicServiceStarted = "system.service.started"
	TopicServiceStopped = "system.service.stopped"
	TopicServiceFailed  = "system.service.failed"
)

// Common error codes for system operations
const (
	ErrSystemNotInitialized = "system.not_initialized"
	ErrSystemNotStarted     = "system.not_started"
	ErrOperationNotFound    = "system.operation_not_found"
	ErrOperationFailed      = "system.operation_failed"
	ErrServiceNotFound      = "system.service_not_found"
	ErrServiceStart         = "system.service_start_failed"
	ErrServiceStop          = "system.service_stop_failed"
)
