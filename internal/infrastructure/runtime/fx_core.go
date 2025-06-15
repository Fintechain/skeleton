// Package runtime provides FX dependency injection integration for the Fintechain Skeleton framework.
//
// This package bridges the gap between Uber's FX dependency injection framework
// and our domain-driven component system, providing clean separation between
// infrastructure concerns and business logic.
package runtime

import (
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
// This provides a default memory-based configuration for development and testing.
// For production use, users should provide their own Configuration implementation via FX.
func NewConfiguration(p ConfigurationParams) ConfigurationResult {
	// Provide a default memory configuration for development and testing
	return ConfigurationResult{
		Configuration: infraConfig.NewMemoryConfiguration(),
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

	Runtime *Runtime
}

// NewRuntimeEnvironment creates a new runtime environment using FX-provided dependencies.
func NewRuntimeEnvironment(p RuntimeEnvironmentParams) (RuntimeEnvironmentResult, error) {
	runtime, err := NewRuntime(
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
