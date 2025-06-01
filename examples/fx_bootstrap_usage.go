package main

import (
	"context"
	"fmt"
	"log"
	"time"

	// Import the updated fx bootstrap
	fxBootstrap "github.com/fintechain/skeleton/old/infrastructure/system"

	// Import current Skeleton Framework API
	"github.com/fintechain/skeleton/pkg/component"
	"github.com/fintechain/skeleton/pkg/event"
	"github.com/fintechain/skeleton/pkg/logging"
	"github.com/fintechain/skeleton/pkg/plugin"
	"github.com/fintechain/skeleton/pkg/registry"
	"github.com/fintechain/skeleton/pkg/storage"

	// Import domain interfaces for plugin implementation
	domainContext "github.com/fintechain/skeleton/internal/domain/context"
	domainRegistry "github.com/fintechain/skeleton/internal/domain/registry"
)

// ExamplePlugin demonstrates a simple plugin implementation
type ExamplePlugin struct {
	id          string
	name        string
	description string
	version     string
}

func NewExamplePlugin() *ExamplePlugin {
	return &ExamplePlugin{
		id:          "example-plugin",
		name:        "Example Plugin",
		description: "A simple example plugin for demonstration",
		version:     "1.0.0",
	}
}

func (p *ExamplePlugin) ID() string          { return p.id }
func (p *ExamplePlugin) Name() string        { return p.name }
func (p *ExamplePlugin) Description() string { return p.description }
func (p *ExamplePlugin) Version() string     { return p.version }

func (p *ExamplePlugin) Load(ctx domainContext.Context, registrar domainRegistry.Registry) error {
	log.Printf("Loading plugin: %s v%s", p.name, p.version)

	// Create and register a sample component
	config := component.NewComponentConfig(
		"example-component",
		"Example Component",
		component.TypeBasic,
		"A component created by the example plugin",
	)

	comp := component.NewBaseComponent(config)
	return registrar.Register(comp)
}

func (p *ExamplePlugin) Unload(ctx domainContext.Context) error {
	log.Printf("Unloading plugin: %s", p.name)
	return nil
}

func main() {
	fmt.Println("=== Skeleton Framework FX Bootstrap Examples ===")
	fmt.Println("This demonstrates how to use the updated FX bootstrap with the current Skeleton Framework API")

	// Example 1: Basic usage with default configuration
	fmt.Println("\n=== Example 1: Basic Usage ===")
	basicExample()

	// Example 2: Custom configuration
	fmt.Println("\n=== Example 2: Custom Configuration ===")
	customConfigExample()

	// Example 3: With plugins
	fmt.Println("\n=== Example 3: With Plugins ===")
	pluginExample()

	// Example 4: Controlled lifecycle
	fmt.Println("\n=== Example 4: Controlled Lifecycle ===")
	controlledLifecycleExample()

	// Example 5: Migration example
	migrationExample()
}

func basicExample() {
	// Use default configuration - this will create all dependencies automatically
	// Note: This would block indefinitely, so we use the context version instead
	config := &fxBootstrap.SystemConfig{
		Config: &fxBootstrap.Config{
			ServiceID: "basic-example",
			StorageConfig: storage.MultiStoreConfig{
				RootPath:      "./basic-data",
				DefaultEngine: "memory",
				EngineConfigs: make(map[string]storage.Config),
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	app, err := fxBootstrap.StartWithFxAndContext(ctx, config)
	if err != nil {
		log.Printf("Error starting basic system: %v", err)
		return
	}

	log.Println("Basic system started successfully")

	// Stop the system
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer stopCancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Printf("Error stopping basic system: %v", err)
	}
}

func customConfigExample() {
	// Create custom configuration
	config := &fxBootstrap.SystemConfig{
		Config: &fxBootstrap.Config{
			ServiceID: "my-custom-service",
			StorageConfig: storage.MultiStoreConfig{
				RootPath:      "./custom-data",
				DefaultEngine: "memory",
				EngineConfigs: map[string]storage.Config{
					"memory": {
						"cache_size":   1000,
						"max_versions": 10,
					},
				},
			},
		},
		// You can also provide custom implementations
		Logger: logging.NewLoggerWithLevel(logging.InfoLevel),
	}

	// Start with custom configuration
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app, err := fxBootstrap.StartWithFxAndContext(ctx, config)
	if err != nil {
		log.Printf("Error starting system: %v", err)
		return
	}

	log.Println("System started successfully with custom configuration")

	// Stop the system
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Printf("Error stopping system: %v", err)
	}
}

func pluginExample() {
	// Create plugins
	plugins := []plugin.Plugin{
		NewExamplePlugin(),
	}

	// Create configuration with plugins
	config := &fxBootstrap.SystemConfig{
		Config: &fxBootstrap.Config{
			ServiceID: "plugin-enabled-service",
			StorageConfig: storage.MultiStoreConfig{
				RootPath:      "./plugin-data",
				DefaultEngine: "memory",
				EngineConfigs: make(map[string]storage.Config),
			},
		},
		Plugins: plugins,
	}

	// Start system with plugins
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app, err := fxBootstrap.StartWithFxAndContext(ctx, config)
	if err != nil {
		log.Printf("Error starting system with plugins: %v", err)
		return
	}

	log.Println("System started successfully with plugins")

	// Stop the system
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Printf("Error stopping system: %v", err)
	}
}

func controlledLifecycleExample() {
	// Create a system with full control over dependencies
	registry := registry.NewRegistry()
	eventBus := event.NewEventBus()
	multiStore := storage.NewMultiStore()
	logger := logging.NewLoggerWithPrefix("CONTROLLED")

	// Create custom configuration with pre-built dependencies
	config := &fxBootstrap.SystemConfig{
		Config: &fxBootstrap.Config{
			ServiceID: "controlled-service",
			StorageConfig: storage.MultiStoreConfig{
				RootPath:      "./controlled-data",
				DefaultEngine: "memory",
				EngineConfigs: make(map[string]storage.Config),
			},
		},
		Registry:   registry,
		EventBus:   eventBus,
		MultiStore: multiStore,
		Logger:     logger,
	}

	// Pre-register some components
	exampleConfig := component.NewComponentConfig(
		"pre-registered-component",
		"Pre-registered Component",
		component.TypeBasic,
		"A component registered before system start",
	)
	comp := component.NewBaseComponent(exampleConfig)
	if err := registry.Register(comp); err != nil {
		log.Printf("Error registering component: %v", err)
		return
	}

	// Subscribe to events
	eventBus.Subscribe("system.*", func(event *event.Event) {
		logger.Info("Received event: %s from %s", event.Topic, event.Source)
	})

	// Start the system
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app, err := fxBootstrap.StartWithFxAndContext(ctx, config)
	if err != nil {
		log.Printf("Error starting controlled system: %v", err)
		return
	}

	log.Println("Controlled system started successfully")
	log.Printf("Registry contains %d items", registry.Count())

	// Stop the system
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Printf("Error stopping system: %v", err)
	}
}

// Migration example: How to migrate from old API to new API
func migrationExample() {
	fmt.Println("\n=== Migration Example ===")

	// OLD WAY (would not work with current API):
	// This is what the old code might have looked like
	/*
		oldConfig := &OldSystemConfig{
			ServiceID: "old-service",
			// ... old configuration structure
		}

		// Old factory functions that no longer exist:
		// registry := component.CreateRegistry()
		// eventBus := event.CreateEventBus()
		// pluginMgr := plugin.CreatePluginManager()
	*/

	// NEW WAY (works with current API):
	newConfig := &fxBootstrap.SystemConfig{
		Config: &fxBootstrap.Config{
			ServiceID: "migrated-service",
			StorageConfig: storage.MultiStoreConfig{
				RootPath:      "./migrated-data",
				DefaultEngine: "memory",
				EngineConfigs: make(map[string]storage.Config),
			},
		},
		// Use current API factory functions
		Registry:   registry.NewRegistry(),
		EventBus:   event.NewEventBus(),
		MultiStore: storage.NewMultiStore(),
		Logger:     logging.NewLogger(),
		// Note: PluginManager requires filesystem dependency, so we leave it nil
		PluginMgr: nil,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	app, err := fxBootstrap.StartWithFxAndContext(ctx, newConfig)
	if err != nil {
		log.Printf("Error starting migrated system: %v", err)
		return
	}

	log.Println("Successfully migrated to new API!")

	// Stop the system
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()

	if err := app.Stop(stopCtx); err != nil {
		log.Printf("Error stopping migrated system: %v", err)
	}
}
