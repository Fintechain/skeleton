// Package fx provides FX dependency injection integration for the Fintechain Skeleton framework.
//
// This package offers a clean API for creating and running Fintechain applications using
// Uber's FX framework for dependency injection and lifecycle management.
//
// Usage Examples:
//
// Daemon Mode (long-running service):
//
//	err := fx.StartDaemon(
//		fx.WithPlugins(myPlugin1, myPlugin2),
//	)
//
// Command Mode (execute and exit):
//
//	result, err := fx.ExecuteCommand("calculate-total",
//		map[string]interface{}{"items": items},
//		fx.WithPlugins(calculatorPlugin),
//	)
package fx

import (
	"os"

	"go.uber.org/fx"

	"github.com/fintechain/skeleton/internal/domain/plugin"
	internalFX "github.com/fintechain/skeleton/internal/fx"
)

// StartDaemon creates and runs a long-running daemon application.
// This function blocks until the application receives a shutdown signal.
//
// The daemon mode is ideal for:
//   - Web servers
//   - Background services
//   - Message processors
//   - Long-running workers
//
// The application will:
//  1. Build all dependencies via FX
//  2. Start the runtime environment
//  3. Load plugins if provided
//  4. Block and wait for shutdown signals (SIGINT, SIGTERM)
//  5. Gracefully shut down all services
//
// Returns an error if startup fails.
func StartDaemon(opts ...RuntimeOption) error {
	return internalFX.StartDaemon(convertOptions(opts)...)
}

// ExecuteCommand creates and runs a command-mode application.
// This function executes a specific operation and returns immediately.
//
// The command mode is ideal for:
//   - CLI commands
//   - Batch processing
//   - One-time calculations
//   - Data transformations
//
// The application will:
//  1. Build all dependencies via FX
//  2. Initialize the runtime environment (without starting long-running services)
//  3. Load plugins if provided
//  4. Execute the specified operation
//  5. Clean up and exit
//
// Parameters:
//   - operationID: The ID of the operation component to execute
//   - input: Input data for the operation
//   - opts: Runtime configuration options
//
// Returns the operation output and any error that occurred.
func ExecuteCommand(operationID string, input map[string]interface{}, opts ...RuntimeOption) (map[string]interface{}, error) {
	return internalFX.ExecuteCommand(operationID, input, convertOptions(opts)...)
}

// StartDaemonWithSignalHandling starts a daemon with custom signal handling.
// This is a convenience function that provides more control over signal handling.
//
// If no signals are provided, defaults to SIGINT and SIGTERM.
func StartDaemonWithSignalHandling(signals []os.Signal, opts ...RuntimeOption) error {
	return internalFX.StartDaemonWithSignalHandling(signals, convertOptions(opts)...)
}

// RuntimeOption configures the runtime behavior.
type RuntimeOption func(*RuntimeConfig)

// RuntimeConfig holds configuration options for runtime creation.
type RuntimeConfig struct {
	// Plugins to load at startup
	Plugins []plugin.Plugin

	// Additional FX options to include in the application
	ExtraOptions []fx.Option
}

// WithPlugins adds plugins to be loaded at startup.
//
// Example:
//
//	fx.StartDaemon(
//		fx.WithPlugins(webServerPlugin, databasePlugin),
//	)
func WithPlugins(plugins ...plugin.Plugin) RuntimeOption {
	return func(cfg *RuntimeConfig) {
		cfg.Plugins = append(cfg.Plugins, plugins...)
	}
}

// WithFXOptions adds additional FX options to the application.
// This allows advanced users to customize the DI container.
//
// Example:
//
//	fx.StartDaemon(
//		fx.WithFXOptions(
//			fx.Provide(MyCustomService),
//			fx.Invoke(MyInitFunction),
//		),
//	)
func WithFXOptions(options ...fx.Option) RuntimeOption {
	return func(cfg *RuntimeConfig) {
		cfg.ExtraOptions = append(cfg.ExtraOptions, options...)
	}
}

// convertOptions converts public RuntimeOption functions to internal ones.
func convertOptions(opts []RuntimeOption) []internalFX.RuntimeOption {
	if len(opts) == 0 {
		return nil
	}

	// Apply public options to get config
	cfg := &RuntimeConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	// Convert to internal options
	var internalOpts []internalFX.RuntimeOption

	if len(cfg.Plugins) > 0 {
		internalOpts = append(internalOpts, internalFX.WithPlugins(cfg.Plugins...))
	}

	if len(cfg.ExtraOptions) > 0 {
		internalOpts = append(internalOpts, internalFX.WithFXOptions(cfg.ExtraOptions...))
	}

	return internalOpts
}
