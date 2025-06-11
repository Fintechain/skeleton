// Package runtime provides runtime environment builder for easy system setup.
package runtime

import (
	"os"

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
)

// RuntimeOption configures the runtime environment.
type RuntimeOption func(*runtimeConfig)

// runtimeConfig holds the configuration for building a runtime.
type runtimeConfig struct {
	plugins       []plugin.Plugin
	registry      component.Registry
	pluginManager plugin.PluginManager
	eventBus      event.EventBusService
	logger        logging.LoggerService
	configuration config.Configuration
}

// NewRuntimeWithOptions creates a new runtime environment with the provided options.
func NewRuntimeWithOptions(opts ...RuntimeOption) (*Runtime, error) {
	cfg := &runtimeConfig{}

	// Apply all options
	for _, opt := range opts {
		opt(cfg)
	}

	// Create defaults for any missing services
	cfg.applyDefaults()

	// Create runtime with direct dependencies
	runtime, err := NewRuntime(
		cfg.registry,
		cfg.configuration,
		cfg.pluginManager,
		cfg.eventBus,
		cfg.logger,
	)
	if err != nil {
		return nil, err
	}

	// Load plugins if any were provided
	if len(cfg.plugins) > 0 {
		ctx := infraContext.NewContext()
		if err := runtime.LoadPlugins(ctx, cfg.plugins); err != nil {
			return nil, err
		}
	}

	return runtime, nil
}

// WithPlugins sets the plugins to load at startup.
func WithPlugins(plugins ...plugin.Plugin) RuntimeOption {
	return func(cfg *runtimeConfig) {
		cfg.plugins = append(cfg.plugins, plugins...)
	}
}

// WithRegistry sets a custom registry implementation.
func WithRegistry(registry component.Registry) RuntimeOption {
	return func(cfg *runtimeConfig) {
		cfg.registry = registry
	}
}

// WithPluginManager sets a custom plugin manager implementation.
func WithPluginManager(pluginManager plugin.PluginManager) RuntimeOption {
	return func(cfg *runtimeConfig) {
		cfg.pluginManager = pluginManager
	}
}

// WithEventBus sets a custom event bus implementation.
func WithEventBus(eventBus event.EventBusService) RuntimeOption {
	return func(cfg *runtimeConfig) {
		cfg.eventBus = eventBus
	}
}

// WithLogger sets a custom logger implementation.
func WithLogger(logger logging.LoggerService) RuntimeOption {
	return func(cfg *runtimeConfig) {
		cfg.logger = logger
	}
}

// WithConfiguration sets a custom configuration implementation.
func WithConfiguration(configuration config.Configuration) RuntimeOption {
	return func(cfg *runtimeConfig) {
		cfg.configuration = configuration
	}
}

// applyDefaults creates default implementations for any missing services.
func (cfg *runtimeConfig) applyDefaults() {
	if cfg.registry == nil {
		cfg.registry = infraComponent.NewRegistry()
	}

	if cfg.pluginManager == nil {
		pmConfig := component.ComponentConfig{
			ID:   "plugin_manager",
			Name: "Plugin Manager",
			Type: component.TypeService,
		}
		cfg.pluginManager = infraPlugin.NewManager(pmConfig)
	}

	if cfg.eventBus == nil {
		ebConfig := component.ComponentConfig{
			ID:   "event_bus",
			Name: "Event Bus",
			Type: component.TypeService,
		}
		cfg.eventBus = infraEvent.NewEventBus(ebConfig)
	}

	if cfg.logger == nil {
		logConfig := component.ComponentConfig{
			ID:   "logger",
			Name: "Logger",
			Type: component.TypeService,
		}
		noOpLogger := infraLogging.NewNoOpLogger()
		logger, err := infraLogging.NewLogger(logConfig, noOpLogger)
		if err != nil {
			panic("failed to create default logger: " + err.Error())
		}
		cfg.logger = logger
	}

	if cfg.configuration == nil {
		cfg.configuration = createDefaultConfiguration()
	}
}

// createDefaultConfiguration creates a default configuration with graceful error handling.
func createDefaultConfiguration() config.Configuration {
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

	return composite
}

// tryCreateFileSource creates a file configuration source if the file exists.
func tryCreateFileSource(path string) config.ConfigurationSource {
	if _, err := os.Stat(path); err == nil {
		return infraConfig.NewFileSource(path)
	}
	return nil
}
