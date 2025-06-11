// Package fx provides FX dependency injection integration for the Fintechain Skeleton framework.
//
// This package bridges the gap between Uber's FX dependency injection framework
// and our domain-driven component system, providing clean separation between
// infrastructure concerns and business logic.
package fx

import (
	"os"

	"go.uber.org/fx"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/config"
	"github.com/fintechain/skeleton/internal/domain/event"
	"github.com/fintechain/skeleton/internal/domain/logging"
	"github.com/fintechain/skeleton/internal/domain/plugin"
	infraComponent "github.com/fintechain/skeleton/internal/infrastructure/component"
	infraConfig "github.com/fintechain/skeleton/internal/infrastructure/config"
	infraEvent "github.com/fintechain/skeleton/internal/infrastructure/event"
	infraLogging "github.com/fintechain/skeleton/internal/infrastructure/logging"
	infraPlugin "github.com/fintechain/skeleton/internal/infrastructure/plugin"
	infraRuntime "github.com/fintechain/skeleton/internal/infrastructure/runtime"
)

// CoreModule provides all essential dependencies for the Fintechain runtime system.
// This module follows FX best practices by using parameter/result objects and
// providing only the dependencies that this module owns.
var CoreModule = fx.Module("fintechain-core",
	fx.Provide(
		NewConfiguration,
		NewRegistry,
		NewEventBus,
		NewLogger,
		NewPluginManager,
		NewRuntimeEnvironment,
	),
)

// ConfigurationParams holds dependencies needed to create a Configuration.
type ConfigurationParams struct {
	fx.In
	// No dependencies - Configuration is created with default sources
}

// ConfigurationResult holds the results of creating a Configuration.
type ConfigurationResult struct {
	fx.Out

	Configuration config.Configuration
}

// NewConfiguration creates a new configuration with default sources.
// This follows the existing builder pattern but provides it through FX.
func NewConfiguration(p ConfigurationParams) ConfigurationResult {
	// Use the same default configuration logic as the builder
	sources := []config.ConfigurationSource{}

	// Add file sources with graceful handling
	configFiles := []string{"config.json", "skeleton.json"}
	for _, file := range configFiles {
		if fileSource := tryCreateFileSource(file); fileSource != nil {
			sources = append(sources, fileSource)
		}
	}

	// Always add environment source (precedence: env > files)
	envSource := infraConfig.NewEnvSource("SKELETON_")
	sources = append(sources, envSource)

	// Create composite configuration
	composite := infraConfig.NewCompositeConfig(sources...)

	// Load with graceful error handling
	if err := composite.LoadConfig(); err != nil {
		// Non-fatal: continue with environment variables and defaults
		// System will work with environment variables and built-in defaults
	}

	return ConfigurationResult{
		Configuration: composite,
	}
}

// RegistryParams holds dependencies needed to create a Registry.
type RegistryParams struct {
	fx.In
	// No dependencies - Registry is a standalone component
}

// RegistryResult holds the results of creating a Registry.
type RegistryResult struct {
	fx.Out

	Registry component.Registry
}

// NewRegistry creates a new component registry.
func NewRegistry(p RegistryParams) RegistryResult {
	return RegistryResult{
		Registry: infraComponent.NewRegistry(),
	}
}

// EventBusParams holds dependencies needed to create an EventBus.
type EventBusParams struct {
	fx.In
	// No dependencies - EventBus is self-contained
}

// EventBusResult holds the results of creating an EventBus.
type EventBusResult struct {
	fx.Out

	EventBus event.EventBusService
}

// NewEventBus creates a new event bus service.
func NewEventBus(p EventBusParams) EventBusResult {
	config := component.ComponentConfig{
		ID:   "event_bus",
		Name: "Event Bus",
		Type: component.TypeService,
	}

	return EventBusResult{
		EventBus: infraEvent.NewEventBus(config),
	}
}

// LoggerParams holds dependencies needed to create a Logger.
type LoggerParams struct {
	fx.In
	// No dependencies - Logger uses NoOp implementation by default
}

// LoggerResult holds the results of creating a Logger.
type LoggerResult struct {
	fx.Out

	Logger logging.LoggerService
}

// NewLogger creates a new logger service.
func NewLogger(p LoggerParams) LoggerResult {
	config := component.ComponentConfig{
		ID:   "logger",
		Name: "Logger",
		Type: component.TypeService,
	}

	noOpLogger := infraLogging.NewNoOpLogger()
	logger, err := infraLogging.NewLogger(config, noOpLogger)
	if err != nil {
		panic("failed to create default logger: " + err.Error())
	}

	return LoggerResult{
		Logger: logger,
	}
}

// PluginManagerParams holds dependencies needed to create a PluginManager.
type PluginManagerParams struct {
	fx.In
	// No dependencies - PluginManager is self-contained
}

// PluginManagerResult holds the results of creating a PluginManager.
type PluginManagerResult struct {
	fx.Out

	PluginManager plugin.PluginManager
}

// NewPluginManager creates a new plugin manager.
func NewPluginManager(p PluginManagerParams) PluginManagerResult {
	config := component.ComponentConfig{
		ID:   "plugin_manager",
		Name: "Plugin Manager",
		Type: component.TypeService,
	}

	return PluginManagerResult{
		PluginManager: infraPlugin.NewManager(config),
	}
}

// RuntimeEnvironmentParams holds dependencies needed to create a RuntimeEnvironment.
type RuntimeEnvironmentParams struct {
	fx.In

	Registry      component.Registry
	Configuration config.Configuration
	PluginManager plugin.PluginManager
	EventBus      event.EventBusService
	Logger        logging.LoggerService
}

// RuntimeEnvironmentResult holds the results of creating a RuntimeEnvironment.
type RuntimeEnvironmentResult struct {
	fx.Out

	Runtime *infraRuntime.Runtime
}

// NewRuntimeEnvironment creates a new runtime environment using FX-provided dependencies.
func NewRuntimeEnvironment(p RuntimeEnvironmentParams) (RuntimeEnvironmentResult, error) {
	runtime, err := infraRuntime.NewRuntime(
		p.Registry,
		p.Configuration,
		p.PluginManager,
		p.EventBus,
		p.Logger,
	)
	if err != nil {
		return RuntimeEnvironmentResult{}, err
	}

	return RuntimeEnvironmentResult{
		Runtime: runtime,
	}, nil
}

// tryCreateFileSource creates a file configuration source if the file exists.
// This is a helper function copied from the existing builder implementation.
func tryCreateFileSource(path string) config.ConfigurationSource {
	// Check if file exists
	if _, err := os.Stat(path); err == nil {
		return infraConfig.NewFileSource(path)
	}
	return nil
}
