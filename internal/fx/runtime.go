// Package fx provides runtime mode implementations using FX dependency injection.
package fx

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/fx"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/plugin"
	infraContext "github.com/fintechain/skeleton/internal/infrastructure/context"
	infraRuntime "github.com/fintechain/skeleton/internal/infrastructure/runtime"
)

// RuntimeMode represents the execution mode for the application.
type RuntimeMode int

const (
	// DaemonMode runs the application as a long-running service that
	// starts up, runs indefinitely, and waits for shutdown signals.
	DaemonMode RuntimeMode = iota

	// CommandMode runs the application to execute a specific operation
	// and then exits immediately after completion.
	CommandMode
)

// RuntimeConfig holds configuration options for runtime creation.
type RuntimeConfig struct {
	// Plugins to load at startup
	Plugins []plugin.Plugin

	// Additional FX options to include in the application
	ExtraOptions []fx.Option
}

// RuntimeOption configures the runtime behavior.
type RuntimeOption func(*RuntimeConfig)

// WithPlugins adds plugins to be loaded at startup.
func WithPlugins(plugins ...plugin.Plugin) RuntimeOption {
	return func(cfg *RuntimeConfig) {
		cfg.Plugins = append(cfg.Plugins, plugins...)
	}
}

// WithFXOptions adds additional FX options to the application.
// This allows advanced users to customize the DI container.
func WithFXOptions(options ...fx.Option) RuntimeOption {
	return func(cfg *RuntimeConfig) {
		cfg.ExtraOptions = append(cfg.ExtraOptions, options...)
	}
}

// StartDaemon creates and runs a long-running daemon application.
// This function blocks until the application receives a shutdown signal.
//
// The daemon mode:
// 1. Builds all dependencies via FX
// 2. Starts the runtime environment via FX lifecycle hooks
// 3. Loads plugins if provided
// 4. Blocks and waits for shutdown signals (SIGINT, SIGTERM)
// 5. Gracefully shuts down all services
//
// Returns an error if startup fails.
func StartDaemon(opts ...RuntimeOption) error {
	cfg := &RuntimeConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	// Build FX options
	options := []fx.Option{
		CoreModule,
	}

	// Add plugin loading if plugins are provided
	if len(cfg.Plugins) > 0 {
		options = append(options, fx.Invoke(func(runtime *infraRuntime.Runtime) error {
			ctx := infraContext.NewContext()
			return runtime.LoadPlugins(ctx, cfg.Plugins)
		}))
	}

	// Add lifecycle management for daemon mode
	options = append(options, fx.Invoke(func(lc fx.Lifecycle, runtime *infraRuntime.Runtime) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				fmt.Println("[Fintechain] Starting daemon mode...")
				domainCtx := infraContext.NewContext()
				if err := runtime.Start(domainCtx); err != nil {
					return fmt.Errorf("failed to start runtime: %w", err)
				}
				fmt.Println("[Fintechain] Daemon started successfully")
				return nil
			},
			OnStop: func(ctx context.Context) error {
				fmt.Println("[Fintechain] Stopping daemon...")
				domainCtx := infraContext.NewContext()
				if err := runtime.Stop(domainCtx); err != nil {
					return fmt.Errorf("failed to stop runtime: %w", err)
				}
				fmt.Println("[Fintechain] Daemon stopped successfully")
				return nil
			},
		})
	}))

	// Add extra options
	options = append(options, cfg.ExtraOptions...)

	// Create and run the FX application
	app := fx.New(options...)
	app.Run()
	return nil
}

// ExecuteCommand creates and runs a command-mode application.
// This function executes a specific operation and returns immediately.
//
// The command mode:
// 1. Builds all dependencies via FX
// 2. Initializes the runtime environment (but doesn't start long-running services)
// 3. Loads plugins if provided
// 4. Executes the specified operation
// 5. Cleans up and exits
//
// Parameters:
//   - operationID: The ID of the operation component to execute
//   - input: Input data for the operation
//   - opts: Runtime configuration options
//
// Returns the operation output and any error that occurred.
func ExecuteCommand(operationID string, input map[string]interface{}, opts ...RuntimeOption) (map[string]interface{}, error) {
	cfg := &RuntimeConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	var result map[string]interface{}
	var execError error

	// Build FX options
	options := []fx.Option{
		CoreModule,
	}

	// Add plugin loading if plugins are provided
	if len(cfg.Plugins) > 0 {
		options = append(options, fx.Invoke(func(runtime *infraRuntime.Runtime) error {
			ctx := infraContext.NewContext()
			return runtime.LoadPlugins(ctx, cfg.Plugins)
		}))
	}

	// Add command execution logic
	options = append(options, fx.Invoke(func(runtime *infraRuntime.Runtime) error {
		ctx := infraContext.NewContext()

		fmt.Printf("[Fintechain] Executing command: %s\n", operationID)

		// Initialize runtime without starting services
		// (services are not started in command mode to avoid long-running processes)

		// Execute the operation
		operationInput := component.Input{
			Data:     input,
			Metadata: map[string]string{"mode": "command"},
		}

		output, err := runtime.ExecuteOperation(ctx, component.ComponentID(operationID), operationInput)
		if err != nil {
			execError = fmt.Errorf("operation execution failed: %w", err)
			return execError
		}

		// Type assert the output data
		if outputData, ok := output.Data.(map[string]interface{}); ok {
			result = outputData
		} else {
			// Wrap single values in a map
			result = map[string]interface{}{"result": output.Data}
		}
		fmt.Printf("[Fintechain] Command completed successfully\n")
		return nil
	}))

	// Add extra options
	options = append(options, cfg.ExtraOptions...)

	// Create FX application
	app := fx.New(options...)

	// Start and stop immediately (don't run indefinitely)
	if err := app.Start(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to start command application: %w", err)
	}

	// Stop the application
	if err := app.Stop(context.Background()); err != nil {
		return result, fmt.Errorf("failed to stop command application: %w", err)
	}

	// Return execution error if any occurred during invoke
	if execError != nil {
		return nil, execError
	}

	return result, nil
}

// StartDaemonWithSignalHandling starts a daemon with custom signal handling.
// This is a convenience function that provides more control over signal handling.
func StartDaemonWithSignalHandling(signals []os.Signal, opts ...RuntimeOption) error {
	if len(signals) == 0 {
		signals = []os.Signal{syscall.SIGINT, syscall.SIGTERM}
	}

	cfg := &RuntimeConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	// Build FX options
	options := []fx.Option{
		CoreModule,
	}

	// Add plugin loading if plugins are provided
	if len(cfg.Plugins) > 0 {
		options = append(options, fx.Invoke(func(runtime *infraRuntime.Runtime) error {
			ctx := infraContext.NewContext()
			return runtime.LoadPlugins(ctx, cfg.Plugins)
		}))
	}

	// Custom lifecycle management with signal handling
	options = append(options, fx.Invoke(func(lc fx.Lifecycle, runtime *infraRuntime.Runtime) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				fmt.Println("[Fintechain] Starting daemon with custom signal handling...")
				domainCtx := infraContext.NewContext()
				if err := runtime.Start(domainCtx); err != nil {
					return fmt.Errorf("failed to start runtime: %w", err)
				}

				// Set up signal handling
				sigChan := make(chan os.Signal, 1)
				signal.Notify(sigChan, signals...)

				go func() {
					sig := <-sigChan
					fmt.Printf("[Fintechain] Received signal: %v, initiating shutdown...\n", sig)
					// The FX shutdown will be handled by the application
				}()

				fmt.Println("[Fintechain] Daemon started with signal handling")
				return nil
			},
			OnStop: func(ctx context.Context) error {
				fmt.Println("[Fintechain] Stopping daemon...")
				domainCtx := infraContext.NewContext()
				if err := runtime.Stop(domainCtx); err != nil {
					return fmt.Errorf("failed to stop runtime: %w", err)
				}
				fmt.Println("[Fintechain] Daemon stopped successfully")
				return nil
			},
		})
	}))

	// Add extra options
	options = append(options, cfg.ExtraOptions...)

	// Create and run the FX application
	app := fx.New(options...)
	app.Run()
	return nil
}
