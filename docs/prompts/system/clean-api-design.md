# Simple API Design for System Initialization

## Keep It Simple - Single Entry Point

```go
// system.go - Public API
package system

// Single function - that's it
func StartSystem(opts ...SystemOption) error

// Functional options for customization
type SystemOption func(*systemConfig)

// Basic options
func WithConfig(config *Config) SystemOption
func WithPlugins(plugins ...Plugin) SystemOption

// Custom dependency providers (optional)
func WithRegistry(provider RegistryProvider) SystemOption
func WithPluginManager(provider PluginManagerProvider) SystemOption
func WithEventBus(provider EventBusProvider) SystemOption
func WithMultiStore(provider MultiStoreProvider) SystemOption

// Provider function types
type RegistryProvider func(*Config) (Registry, error)
type PluginManagerProvider func(*Config) (PluginManager, error)
type EventBusProvider func(*Config) (EventBus, error)
type MultiStoreProvider func(*Config) (MultiStore, error)
```

## Usage Examples

```go
// Simple - everything defaults
err := StartSystem()

// With plugins
err := StartSystem(
    WithPlugins(myPlugin1, myPlugin2),
)

// With custom config
err := StartSystem(
    WithConfig(&system.Config{
        ServiceID: "my-app",
    }),
    WithPlugins(myPlugin1, myPlugin2),
)

// With custom dependencies (direct injection)
customRegistry, _ := component.NewDefaultRegistry()
customEventBus, _ := event.NewDefaultEventBus()

err := StartSystem(
    WithConfig(myConfig),
    WithPlugins(myPlugin1, myPlugin2),
    WithRegistry(customRegistry),
    WithEventBus(customEventBus),
)
```

## Complete Implementation Details

### Public API (pkg/system/system.go)

```go
package system

import "context"

// StartSystem is the main entry point to start the component system
func StartSystem(opts ...SystemOption) error {
    config := &systemConfig{}
    
    // Apply all options
    for _, opt := range opts {
        opt(config)
    }
    
    // Start with fx (hidden from client)
    return startWithFx(config)
}

// SystemOption configures the system startup
type SystemOption func(*systemConfig)

// WithConfig sets the system configuration
func WithConfig(config *Config) SystemOption {
    return func(sc *systemConfig) {
        sc.config = config
    }
}

// WithPlugins adds plugins to the system
func WithPlugins(plugins ...Plugin) SystemOption {
    return func(sc *systemConfig) {
        sc.plugins = append(sc.plugins, plugins...)
    }
}

// WithRegistry sets a custom registry implementation
func WithRegistry(registry component.Registry) SystemOption {
    return func(sc *systemConfig) {
        sc.registry = registry
    }
}

// WithPluginManager sets a custom plugin manager implementation
func WithPluginManager(pluginMgr plugin.PluginManager) SystemOption {
    return func(sc *systemConfig) {
        sc.pluginMgr = pluginMgr
    }
}

// WithEventBus sets a custom event bus implementation
func WithEventBus(eventBus event.EventBus) SystemOption {
    return func(sc *systemConfig) {
        sc.eventBus = eventBus
    }
}

// WithMultiStore sets a custom multi-store implementation
func WithMultiStore(multiStore storage.MultiStore) SystemOption {
    return func(sc *systemConfig) {
        sc.multiStore = multiStore
    }
}
```

### Internal Types (internal/domain/system/config.go)

```go
package system

import (
    "github.com/your-org/skeleton/internal/domain/component"
    "github.com/your-org/skeleton/internal/domain/plugin"
    "github.com/your-org/skeleton/internal/domain/storage"
    "github.com/your-org/skeleton/internal/infrastructure/event"
)

// systemConfig holds all configuration for system startup
type systemConfig struct {
    config      *Config
    plugins     []plugin.Plugin
    registry    component.Registry
    pluginMgr   plugin.PluginManager
    eventBus    event.EventBus
    multiStore  storage.MultiStore
}

// applyDefaults creates default implementations for any that weren't provided
func (sc *systemConfig) applyDefaults() error {
    if sc.config == nil {
        sc.config = DefaultConfig()
    }
    
    var err error
    
    if sc.registry == nil {
        sc.registry, err = component.NewDefaultRegistry()
        if err != nil {
            return err
        }
    }
    
    if sc.pluginMgr == nil {
        sc.pluginMgr, err = plugin.NewDefaultPluginManager()
        if err != nil {
            return err
        }
    }
    
    if sc.eventBus == nil {
        sc.eventBus, err = event.NewDefaultEventBus()
        if err != nil {
            return err
        }
    }
    
    if sc.multiStore == nil {
        engine := memory.NewEngine()
        sc.multiStore, err = storage.NewDefaultMultiStore(sc.config.StorageConfig, engine)
        if err != nil {
            return err
        }
    }
    
    return nil
}
```

### Internal Bootstrap (internal/infrastructure/system/fx_bootstrap.go)

```go
package system

import (
    "context"
    "go.uber.org/fx"
    
    "github.com/your-org/skeleton/internal/domain/component"
    "github.com/your-org/skeleton/internal/domain/plugin"
    "github.com/your-org/skeleton/internal/domain/storage"
    "github.com/your-org/skeleton/internal/infrastructure/event"
    "github.com/your-org/skeleton/internal/infrastructure/storage/memory"
)

// startWithFx initializes and starts the system using fx
func startWithFx(config *systemConfig) error {
    config.applyDefaults()
    
    return fx.New(
        // Supply the configuration
        fx.Supply(config),
        
        // Provide all dependencies
        fx.Provide(provideRegistry),
        fx.Provide(providePluginManager),
        fx.Provide(provideEventBus),
        fx.Provide(provideMultiStore),
        fx.Provide(provideSystemService),
        
        // Initialize and start everything
        fx.Invoke(initializeAndStart),
    ).Run()
}

// provideRegistry creates the registry using the configured provider
func provideRegistry(config *systemConfig) (component.Registry, error) {
    return config.registryProvider(config.config)
}

// providePluginManager creates the plugin manager using the configured provider
func providePluginManager(config *systemConfig) (plugin.PluginManager, error) {
    return config.pluginMgrProvider(config.config)
}

// provideEventBus creates the event bus using the configured provider
func provideEventBus(config *systemConfig) (event.EventBus, error) {
    return config.eventBusProvider(config.config)
}

// provideMultiStore creates the multi-store using the configured provider
func provideMultiStore(config *systemConfig) (storage.MultiStore, error) {
    return config.multiStoreProvider(config.config)
}

// provideSystemService creates the main system service with all dependencies
func provideSystemService(
    registry component.Registry,
    pluginMgr plugin.PluginManager,
    eventBus event.EventBus,
    multiStore storage.MultiStore,
    config *systemConfig,
) (*DefaultSystemService, error) {
    return NewDefaultSystemService(SystemServiceConfig{
        Registry:      registry,
        PluginManager: pluginMgr,
        EventBus:      eventBus,
        MultiStore:    multiStore,
        Config:        config.config,
    })
}

// initializeAndStart registers plugins and starts the system
func initializeAndStart(sys *DefaultSystemService, config *systemConfig) error {
    ctx := context.Background()
    
    // Initialize the system
    if err := sys.Initialize(ctx); err != nil {
        return err
    }
    
    // Register all plugins
    for _, plugin := range config.plugins {
        if err := sys.PluginManager().RegisterPlugin(plugin); err != nil {
            return err
        }
    }
    
    // Start the system
    return sys.Start(ctx)
}
```

### Default Creation (internal/infrastructure/system/defaults.go)

```go
package system

import (
    "github.com/your-org/skeleton/internal/domain/component"
    "github.com/your-org/skeleton/internal/domain/plugin"
    "github.com/your-org/skeleton/internal/domain/storage"
    "github.com/your-org/skeleton/internal/infrastructure/event"
    "github.com/your-org/skeleton/internal/infrastructure/storage/memory"
    "github.com/your-org/skeleton/internal/infrastructure/config"
)

// Default configuration
func DefaultConfig() *Config {
    return &Config{
        ServiceID: "system",
        StorageConfig: storage.MultiStoreConfig{
            RootPath:      "./data",
            DefaultEngine: "memory",
        },
    }
}
```

### Package Structure

```
pkg/
└── system/
    ├── system.go               // Public API - StartSystem function
    ├── config.go              // Config types
    └── types.go               // Public types and interfaces

internal/
├── domain/
│   ├── component/             // Component interfaces and implementations
│   ├── plugin/               // Plugin interfaces and implementations  
│   ├── service/              // Service interfaces and implementations
│   ├── operation/            // Operation interfaces and implementations
│   ├── storage/              // Storage interfaces
│   └── system/               // System domain types and interfaces
└── infrastructure/
    ├── system/               // System service implementation + fx integration
    │   ├── default_system_service.go
    │   ├── fx_bootstrap.go   // fx integration (hidden)
    │   ├── config.go         // systemConfig type
    │   └── defaults.go       // default providers
    ├── event/               // Event bus implementation
    ├── storage/             // Storage implementations
    ├── config/              // Configuration implementation
    └── logging/             // Logging implementation
```

### Provider Function Types (pkg/system/types.go)

```go
package system

// Re-export commonly used types for client convenience
// These are the actual interfaces from their respective packages

// Plugin represents a plugin that can be registered with the system
type Plugin = plugin.Plugin

// Config represents system configuration
type Config struct {
    ServiceID     string                     `json:"serviceId"`
    StorageConfig storage.MultiStoreConfig   `json:"storage"`
}
```