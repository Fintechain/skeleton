package system

import (
	"context"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/ebanfa/skeleton/internal/domain/component"
	"github.com/ebanfa/skeleton/internal/domain/plugin"
	"github.com/ebanfa/skeleton/internal/domain/storage"
	"github.com/ebanfa/skeleton/internal/domain/system"
	"github.com/ebanfa/skeleton/internal/infrastructure/config"
	infraContext "github.com/ebanfa/skeleton/internal/infrastructure/context"
	"github.com/ebanfa/skeleton/internal/infrastructure/event"
	infraEvent "github.com/ebanfa/skeleton/internal/infrastructure/event"
	"github.com/ebanfa/skeleton/internal/infrastructure/logging"
	infraStorage "github.com/ebanfa/skeleton/internal/infrastructure/storage"
)

// SystemConfig holds all the configuration for the system
type SystemConfig struct {
	Config     *Config
	Plugins    []plugin.Plugin
	Registry   component.Registry
	PluginMgr  plugin.PluginManager
	EventBus   event.EventBus
	MultiStore storage.MultiStore
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
			},
		}
	}

	if config.Registry == nil {
		config.Registry = component.CreateRegistry()
	}

	if config.PluginMgr == nil {
		config.PluginMgr = plugin.CreatePluginManager()
	}

	if config.EventBus == nil {
		config.EventBus = event.CreateEventBus()
	}

	if config.MultiStore == nil {
		// Create a simple memory-based multistore
		config.MultiStore = createDefaultMultiStore(config.Config.StorageConfig)
	}

	return config
}

// createDefaultMultiStore creates a proper multistore implementation
func createDefaultMultiStore(config storage.MultiStoreConfig) storage.MultiStore {
	logger := logging.CreateStandardLogger(logging.Info)
	eventBus := infraEvent.NewEventBus()
	return infraStorage.NewMultiStore(&config, logger, eventBus)
}

// provideSystemService creates the system service with all dependencies
func provideSystemService(
	config *Config,
	registry component.Registry,
	pluginMgr plugin.PluginManager,
	eventBus event.EventBus,
	multiStore storage.MultiStore,
	plugins []plugin.Plugin,
) (*DefaultSystemService, error) {
	// Convert to domain config
	domainConfig := &system.SystemServiceConfig{
		ServiceID:        config.ServiceID,
		EnableOperations: true,
		EnableServices:   true,
		EnablePlugins:    true,
		EnableEventLog:   true,
		StorageConfig:    config.StorageConfig,
	}

	// Create factory
	logger := logging.CreateStandardLogger(logging.Info)
	factory := NewFactory(
		registry,
		pluginMgr,
		eventBus,
		config.CreateConfiguration(),
		multiStore,
		logger,
	)

	// Create system service
	svc, err := factory.CreateSystemService(domainConfig)
	if err != nil {
		return nil, err
	}

	// Cast to concrete type for fx
	return svc.(*DefaultSystemService), nil
}

// initializeAndStart registers plugins and starts the system
func initializeAndStart(sys *DefaultSystemService, plugins []plugin.Plugin) error {
	ctx := infraContext.Background()

	// Initialize the system
	if err := sys.Initialize(ctx); err != nil {
		return err
	}

	// Register all plugins - cast to concrete type to access RegisterPlugin
	if defaultPluginMgr, ok := sys.PluginManager().(*plugin.DefaultPluginManager); ok {
		for _, plugin := range plugins {
			if err := defaultPluginMgr.RegisterPlugin(plugin); err != nil {
				return err
			}
		}
	}

	// Start the system
	return sys.Start(ctx)
}

// CreateConfiguration creates a configuration from the Config
func (c *Config) CreateConfiguration() config.Configuration {
	cfg := config.CreateDefaultConfig()
	cfg.Set("system.serviceId", c.ServiceID)
	cfg.Set("system.storage.rootPath", c.StorageConfig.RootPath)
	cfg.Set("system.storage.defaultEngine", c.StorageConfig.DefaultEngine)
	return cfg
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
			func() component.Registry { return config.Registry },
			func() plugin.PluginManager { return config.PluginMgr },
			func() event.EventBus { return config.EventBus },
			func() storage.MultiStore { return config.MultiStore },
			provideSystemService,
		),
		fx.Invoke(initializeAndStart),
		fx.WithLogger(func() fxevent.Logger {
			return fxevent.NopLogger
		}),
	)
}
