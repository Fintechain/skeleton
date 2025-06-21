// Package runtime provides a clean builder API for the Fintechain Skeleton framework.
//
// The RuntimeBuilder replaces FX dependency injection with a simple, explicit
// builder pattern that supports custom dependency injection while maintaining
// the same functionality.
//
// Usage Examples:
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
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/config"
	"github.com/fintechain/skeleton/internal/domain/event"
	"github.com/fintechain/skeleton/internal/domain/logging"
	"github.com/fintechain/skeleton/internal/domain/plugin"
	infraComponent "github.com/fintechain/skeleton/internal/infrastructure/component"
	infraConfig "github.com/fintechain/skeleton/internal/infrastructure/config"
	infraContext "github.com/fintechain/skeleton/internal/infrastructure/context"
	infraEvent "github.com/fintechain/skeleton/internal/infrastructure/event"
	infraLogging "github.com/fintechain/skeleton/internal/infrastructure/logging"
	infraPlugin "github.com/fintechain/skeleton/internal/infrastructure/plugin"
	infraRuntime "github.com/fintechain/skeleton/internal/infrastructure/runtime"
)

// RuntimeBuilder provides a simple builder API for creating and running
// Fintechain applications without FX dependency injection complexity.
type RuntimeBuilder struct {
	plugins   []plugin.Plugin
	config    config.Configuration
	logger    logging.LoggerService
	eventBus  event.EventBusService
	registry  component.Registry
	pluginMgr plugin.PluginManager
}

// NewBuilder creates a new RuntimeBuilder with no dependencies set.
// Dependencies will be created with defaults when Build methods are called.
func NewBuilder() *RuntimeBuilder {
	return &RuntimeBuilder{}
}

// WithPlugins adds plugins to be loaded at startup.
//
// Example:
//
//	builder := runtime.NewBuilder().
//		WithPlugins(webPlugin, dbPlugin)
func (b *RuntimeBuilder) WithPlugins(plugins ...plugin.Plugin) *RuntimeBuilder {
	b.plugins = append(b.plugins, plugins...)
	return b
}

// WithConfig sets a custom configuration service.
// If not set, a default memory configuration will be used.
//
// Example:
//
//	config := infraConfig.NewMemoryConfigurationWithData(map[string]interface{}{
//		"app.name": "My App",
//		"app.port": 8080,
//	})
//	builder := runtime.NewBuilder().
//		WithConfig(config)
func (b *RuntimeBuilder) WithConfig(config config.Configuration) *RuntimeBuilder {
	b.config = config
	return b
}

// WithLogger sets a custom logger service.
// If not set, a default NoOp logger will be used.
//
// Example:
//
//	logger := myCustomLogger()
//	builder := runtime.NewBuilder().
//		WithLogger(logger)
func (b *RuntimeBuilder) WithLogger(logger logging.LoggerService) *RuntimeBuilder {
	b.logger = logger
	return b
}

// WithEventBus sets a custom event bus service.
// If not set, a default in-memory event bus will be used.
//
// Example:
//
//	eventBus := myCustomEventBus()
//	builder := runtime.NewBuilder().
//		WithEventBus(eventBus)
func (b *RuntimeBuilder) WithEventBus(eventBus event.EventBusService) *RuntimeBuilder {
	b.eventBus = eventBus
	return b
}

// WithRegistry sets a custom component registry.
// If not set, a default in-memory registry will be used.
//
// Example:
//
//	registry := myCustomRegistry()
//	builder := runtime.NewBuilder().
//		WithRegistry(registry)
func (b *RuntimeBuilder) WithRegistry(registry component.Registry) *RuntimeBuilder {
	b.registry = registry
	return b
}

// WithPluginManager sets a custom plugin manager.
// If not set, a default plugin manager will be used.
//
// Example:
//
//	pluginMgr := myCustomPluginManager()
//	builder := runtime.NewBuilder().
//		WithPluginManager(pluginMgr)
func (b *RuntimeBuilder) WithPluginManager(pluginMgr plugin.PluginManager) *RuntimeBuilder {
	b.pluginMgr = pluginMgr
	return b
}

// createDefaultDependencies creates default implementations for any dependencies
// that were not explicitly set via WithXxx methods.
func (b *RuntimeBuilder) createDefaultDependencies() error {
	// Create default configuration if not set
	if b.config == nil {
		b.config = infraConfig.NewMemoryConfiguration()
	}

	// Create default registry if not set
	if b.registry == nil {
		b.registry = infraComponent.NewRegistry()
	}

	// Create default event bus if not set
	if b.eventBus == nil {
		config := component.ComponentConfig{
			ID:   "event_bus",
			Name: "Event Bus",
			Type: component.TypeService,
		}
		b.eventBus = infraEvent.NewEventBus(config)
	}

	// Create default logger if not set
	if b.logger == nil {
		config := component.ComponentConfig{
			ID:   "logger",
			Name: "Logger",
			Type: component.TypeService,
		}
		noOpLogger := infraLogging.NewNoOpLogger()
		logger, err := infraLogging.NewLogger(config, noOpLogger)
		if err != nil {
			return fmt.Errorf("failed to create default logger: %w", err)
		}
		b.logger = logger
	}

	// Create default plugin manager if not set
	if b.pluginMgr == nil {
		config := component.ComponentConfig{
			ID:   "plugin_manager",
			Name: "Plugin Manager",
			Type: component.TypeService,
		}
		b.pluginMgr = infraPlugin.NewManager(config)
	}

	return nil
}

// createRuntime creates the runtime with all dependencies.
func (b *RuntimeBuilder) createRuntime() (*infraRuntime.Runtime, error) {
	// Ensure all dependencies are available
	if err := b.createDefaultDependencies(); err != nil {
		return nil, err
	}

	// Create runtime with dependencies using direct constructor
	runtime, err := infraRuntime.NewRuntime(
		b.registry,
		b.config,
		b.pluginMgr,
		b.eventBus,
		b.logger,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create runtime: %w", err)
	}

	return runtime, nil
}

// BuildDaemon creates and runs a long-running daemon application.
// This function blocks until the application receives a shutdown signal.
//
// The daemon mode is ideal for:
//   - Web servers
//   - Background services
//   - Message processors
//   - Long-running workers
//
// The application will:
//  1. Create dependencies (use defaults if not provided)
//  2. Create runtime using existing constructor
//  3. Load plugins if provided
//  4. Start runtime and handle signals
//  5. Block and wait for shutdown signals (SIGINT, SIGTERM)
//  6. Gracefully shut down all services
//
// Returns an error if startup fails.
func (b *RuntimeBuilder) BuildDaemon() error {
	// Create runtime
	runtime, err := b.createRuntime()
	if err != nil {
		return err
	}

	// Create context
	ctx := infraContext.NewContext()

	// Load plugins if provided
	if len(b.plugins) > 0 {
		fmt.Printf("[Fintechain] Loading %d plugins...\n", len(b.plugins))
		if err := runtime.LoadPlugins(ctx, b.plugins); err != nil {
			return fmt.Errorf("failed to load plugins: %w", err)
		}
	}

	// Start runtime
	fmt.Println("[Fintechain] Starting daemon mode...")
	if err := runtime.Start(ctx); err != nil {
		return fmt.Errorf("failed to start runtime: %w", err)
	}
	fmt.Println("[Fintechain] Daemon started successfully")

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Block until signal received
	<-sigChan
	fmt.Println("[Fintechain] Shutdown signal received")

	// Stop runtime
	fmt.Println("[Fintechain] Stopping daemon...")
	if err := runtime.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop runtime: %w", err)
	}
	fmt.Println("[Fintechain] Daemon stopped successfully")

	return nil
}

// BuildCommand creates and runs a command-mode application.
// This function executes a specific operation and returns immediately.
//
// The command mode is ideal for:
//   - CLI commands
//   - Batch processing
//   - One-time calculations
//   - Data transformations
//
// The application will:
//  1. Create dependencies (use defaults if not provided)
//  2. Create runtime using existing constructor
//  3. Load plugins if provided
//  4. Execute the specified operation
//  5. Clean up and return results
//
// Parameters:
//   - operationID: The ID of the operation component to execute
//   - input: Input data for the operation
//
// Returns the operation output and any error that occurred.
func (b *RuntimeBuilder) BuildCommand(operationID string, input map[string]interface{}) (map[string]interface{}, error) {
	// Create runtime
	runtime, err := b.createRuntime()
	if err != nil {
		return nil, err
	}

	// Create context
	ctx := infraContext.NewContext()

	// Load plugins if provided
	if len(b.plugins) > 0 {
		fmt.Printf("[Fintechain] Loading %d plugins...\n", len(b.plugins))
		if err := runtime.LoadPlugins(ctx, b.plugins); err != nil {
			return nil, fmt.Errorf("failed to load plugins: %w", err)
		}
	}

	// Initialize runtime without starting long-running services
	// (services are not started in command mode to avoid long-running processes)
	fmt.Printf("[Fintechain] Executing command: %s\n", operationID)

	// Execute the operation
	operationInput := component.Input{
		Data:     input,
		Metadata: map[string]string{"mode": "command"},
	}

	output, err := runtime.ExecuteOperation(ctx, component.ComponentID(operationID), operationInput)
	if err != nil {
		return nil, fmt.Errorf("operation execution failed: %w", err)
	}

	// Type assert the output data
	var result map[string]interface{}
	if outputData, ok := output.Data.(map[string]interface{}); ok {
		result = outputData
	} else {
		// Wrap single values in a map
		result = map[string]interface{}{"result": output.Data}
	}

	fmt.Printf("[Fintechain] Command completed successfully\n")
	return result, nil
}
