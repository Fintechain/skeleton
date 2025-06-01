package system

import (
	"context"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	// Current Skeleton Framework API imports

	"github.com/fintechain/skeleton/pkg/config"
	frameworkContext "github.com/fintechain/skeleton/pkg/context"
	"github.com/fintechain/skeleton/pkg/event"
	"github.com/fintechain/skeleton/pkg/logging"
	"github.com/fintechain/skeleton/pkg/plugin"
	"github.com/fintechain/skeleton/pkg/registry"
	"github.com/fintechain/skeleton/pkg/storage"
	"github.com/fintechain/skeleton/pkg/system"
)

// SystemConfig holds all the configuration for the system
type SystemConfig struct {
	Config     *Config
	Plugins    []plugin.Plugin
	Registry   registry.Registry
	PluginMgr  plugin.PluginManager
	EventBus   event.EventBus
	MultiStore storage.MultiStore
	Logger     logging.Logger
}

// Config holds the basic system configuration
type Config struct {
	ServiceID     string                   `json:"serviceId"`
	StorageConfig storage.MultiStoreConfig `json:"storage"`
}

// applyDefaults creates default implementations for any that weren't provided
func applyDefaults(config *SystemConfig) *SystemConfig {
	if config == nil {
		config = &SystemConfig{}
	}

	if config.Config == nil {
		config.Config = &Config{
			ServiceID: "system",
			StorageConfig: storage.MultiStoreConfig{
				RootPath:      "./data",
				DefaultEngine: "memory",
				EngineConfigs: make(map[string]storage.Config),
			},
		}
	}

	if config.Registry == nil {
		config.Registry = registry.NewRegistry()
	}

	if config.PluginMgr == nil {
		// Note: PluginManager requires a filesystem dependency
		// For now, we'll create a nil plugin manager and handle this in the provider
		config.PluginMgr = nil
	}

	if config.EventBus == nil {
		config.EventBus = event.NewEventBus()
	}

	if config.MultiStore == nil {
		config.MultiStore = storage.NewMultiStore()
	}

	if config.Logger == nil {
		config.Logger = logging.NewLogger()
	}

	return config
}

// providePluginManager creates a plugin manager with proper dependencies
func providePluginManager() plugin.PluginManager {
	// TODO: Implement a proper filesystem interface for plugin manager
	// For now, return nil and handle gracefully
	return nil
}

// provideConfiguration creates a configuration from the Config
func provideConfiguration(cfg *Config) config.Configuration {
	// Create memory configuration source with the config values
	values := map[string]interface{}{
		"system.serviceId":             cfg.ServiceID,
		"system.storage.rootPath":      cfg.StorageConfig.RootPath,
		"system.storage.defaultEngine": cfg.StorageConfig.DefaultEngine,
	}

	memorySource := config.NewMemoryConfigurationSource(values)
	return config.NewConfiguration(memorySource)
}

// provideSystem creates the system with all dependencies
func provideSystem(
	registry registry.Registry,
	pluginMgr plugin.PluginManager,
	eventBus event.EventBus,
	configuration config.Configuration,
	multiStore storage.MultiStore,
) system.System {
	return system.NewSystem(
		registry,
		pluginMgr,
		eventBus,
		configuration,
		multiStore,
	)
}

// initializeAndStart initializes the system and registers plugins
func initializeAndStart(sys system.System, plugins []plugin.Plugin, logger logging.Logger) error {
	ctx := frameworkContext.NewContext()

	// Register all plugins if plugin manager is available
	if pluginMgr := sys.PluginManager(); pluginMgr != nil {
		for _, plugin := range plugins {
			if err := pluginMgr.Load(ctx, plugin.ID(), sys.Registry()); err != nil {
				logger.Error("Failed to load plugin: %s - %v", plugin.ID(), err)
				return err
			}
		}
	}

	logger.Info("System initialized and started successfully")
	return nil
}

// StartWithFx is the public API for starting a system with fx
// This abstracts fx from the client and provides a simple interface
func StartWithFx(config *SystemConfig) error {
	app := createFxApp(config)
	if err := app.Err(); err != nil {
		return err
	}
	app.Run() // This blocks until signal
	return nil
}

// StartWithFxAndContext starts the system using the fx framework with context control
// This is useful for testing where you need to control the system lifecycle
func StartWithFxAndContext(ctx context.Context, config *SystemConfig) (*fx.App, error) {
	app := createFxApp(config)
	if err := app.Err(); err != nil {
		return nil, err
	}

	// Start the application
	startCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if err := app.Start(startCtx); err != nil {
		return nil, err
	}

	// Return the app so the caller can control its lifecycle
	return app, nil
}

// createFxApp creates the fx application with the given configuration
func createFxApp(config *SystemConfig) *fx.App {
	// Apply defaults
	config = applyDefaults(config)

	return fx.New(
		fx.Provide(
			func() *Config { return config.Config },
			func() []plugin.Plugin { return config.Plugins },
			func() registry.Registry { return config.Registry },
			func() event.EventBus { return config.EventBus },
			func() storage.MultiStore { return config.MultiStore },
			func() logging.Logger { return config.Logger },
			providePluginManager,
			provideConfiguration,
			provideSystem,
		),
		fx.Invoke(initializeAndStart),
		fx.WithLogger(func() fxevent.Logger {
			return fxevent.NopLogger
		}),
	)
}
