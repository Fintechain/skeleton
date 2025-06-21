// Package runtime provides the Builder-based API for the Fintechain Skeleton framework.
//
// This package offers a simple builder pattern for creating and running Fintechain applications
// without the complexity of dependency injection frameworks.
//
// Builder-based Usage Examples:
//
// Daemon Mode (long-running service):
//
//	err := runtime.NewBuilder().
//		WithPlugins(myPlugin1, myPlugin2).
//		WithConfig(myConfig).
//		BuildDaemon()
//
// Command Mode (execute and exit):
//
//	result, err := runtime.NewBuilder().
//		WithPlugins(calculatorPlugin).
//		WithLogger(myLogger).
//		BuildCommand("calculate-total", map[string]interface{}{
//			"items": items,
//		})
package runtime

import (
	"github.com/fintechain/skeleton/internal/domain/config"
	"github.com/fintechain/skeleton/internal/domain/event"
	"github.com/fintechain/skeleton/internal/domain/logging"
	domainRuntime "github.com/fintechain/skeleton/internal/domain/runtime"
)

// RuntimeEnvironment provides high-level runtime environment interface.
// This extends the component.System interface with additional accessors for commonly used core services.
type RuntimeEnvironment = domainRuntime.RuntimeEnvironment

// =============================================================================
// Builder-based API (Recommended)
// =============================================================================

// StartDaemonWithBuilder creates and runs a long-running daemon application using the builder pattern.
// This function blocks until the application receives a shutdown signal.
//
// This approach provides:
// - Simple, explicit dependency injection
// - Easy debugging and testing
// - Custom dependency support
// - No complex framework dependencies
//
// Example:
//
//	err := runtime.StartDaemonWithBuilder(
//		runtime.NewBuilder().
//			WithPlugins(webPlugin, dbPlugin).
//			WithConfig(myConfig),
//	)
func StartDaemonWithBuilder(builder *RuntimeBuilder) error {
	return builder.BuildDaemon()
}

// ExecuteCommandWithBuilder creates and runs a command-mode application using the builder pattern.
// This function executes a specific operation and returns immediately.
//
// This approach provides:
// - Simple, explicit dependency injection
// - Easy debugging and testing
// - Custom dependency support
// - No complex framework dependencies
//
// Example:
//
//	result, err := runtime.ExecuteCommandWithBuilder(
//		"calculate-total",
//		map[string]interface{}{"items": items},
//		runtime.NewBuilder().
//			WithPlugins(calculatorPlugin).
//			WithLogger(myLogger),
//	)
func ExecuteCommandWithBuilder(operationID string, input map[string]interface{}, builder *RuntimeBuilder) (map[string]interface{}, error) {
	return builder.BuildCommand(operationID, input)
}

// =============================================================================
// Builder-based Convenience Functions
// =============================================================================

// WithCustomConfig creates a builder option to set a custom configuration.
// This is a convenience function for the builder API.
//
// Example:
//
//	config := infraConfig.NewMemoryConfigurationWithData(map[string]interface{}{
//		"app.name": "My App",
//		"app.port": 8080,
//	})
//	err := runtime.NewBuilder().
//		WithConfig(config).
//		WithPlugins(myPlugin).
//		BuildDaemon()
func WithCustomConfig(config config.Configuration) func(*RuntimeBuilder) *RuntimeBuilder {
	return func(builder *RuntimeBuilder) *RuntimeBuilder {
		return builder.WithConfig(config)
	}
}

// WithCustomLogger creates a builder option to set a custom logger.
// This is a convenience function for the builder API.
//
// Example:
//
//	logger := myCustomLogger()
//	err := runtime.NewBuilder().
//		WithLogger(logger).
//		WithPlugins(myPlugin).
//		BuildDaemon()
func WithCustomLogger(logger logging.LoggerService) func(*RuntimeBuilder) *RuntimeBuilder {
	return func(builder *RuntimeBuilder) *RuntimeBuilder {
		return builder.WithLogger(logger)
	}
}

// WithCustomEventBus creates a builder option to set a custom event bus.
// This is a convenience function for the builder API.
//
// Example:
//
//	eventBus := myCustomEventBus()
//	err := runtime.NewBuilder().
//		WithEventBus(eventBus).
//		WithPlugins(myPlugin).
//		BuildDaemon()
func WithCustomEventBus(eventBus event.EventBusService) func(*RuntimeBuilder) *RuntimeBuilder {
	return func(builder *RuntimeBuilder) *RuntimeBuilder {
		return builder.WithEventBus(eventBus)
	}
}
