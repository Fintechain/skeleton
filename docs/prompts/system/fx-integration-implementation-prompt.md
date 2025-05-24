# Fx Integration Implementation Prompt (Execute SECOND)

## Task Overview
You are tasked with implementing a clean public API using fx for dependency injection while completely hiding fx from client code. This is **Phase 2** and requires **Phase 1** (PluginManager integration) to be completed first.

## CRITICAL CONSTRAINTS

### Prerequisites Verification
**BEFORE starting, verify Phase 1 completion:**
- [ ] SystemService constructor accepts PluginManager
- [ ] SystemService has PluginManager() getter method
- [ ] SystemServiceConfig includes PluginManager field

### File Creation Rules
**ONLY create these specific files (and NO others):**
- `pkg/system/system.go` - Main public API
- `pkg/system/config.go` - Public configuration types  
- `pkg/system/types.go` - Type aliases and re-exports
- `internal/infrastructure/system/fx_bootstrap.go` - fx integration (hidden)

### Mandatory Analysis Phase
**BEFORE creating ANY files, you MUST:**

1. **Re-examine these packages to understand current state:**
   - `internal/domain/system/system.go` - Updated system interfaces
   - `internal/infrastructure/system/` - Updated SystemService with PluginManager
   - `internal/infrastructure/event/` - Event bus interfaces
   - `internal/infrastructure/config/` - Configuration interfaces
   - `internal/infrastructure/context/` - Context interfaces
   - `internal/infrastructure/storage/` - Storage interfaces
   - `internal/infrastructure/logging/` - Logging interfaces

2. **Understand the updated SystemService:**
   - New constructor signature with PluginManager
   - How to create default implementations
   - Existing interface patterns

## Design Specification

### Public API Design (pkg/system/system.go)
```go
package system

import "context"

// StartSystem is the main entry point - simple functional options pattern
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

// Configuration options
func WithConfig(config *Config) SystemOption
func WithPlugins(plugins ...Plugin) SystemOption

// Dependency injection options (use existing interfaces)
func WithRegistry(registry component.Registry) SystemOption
func WithPluginManager(pluginMgr plugin.PluginManager) SystemOption
func WithEventBus(eventBus event.EventBus) SystemOption
func WithMultiStore(multiStore storage.MultiStore) SystemOption
```

### Public Types (pkg/system/types.go)
```go
package system

import (
    "github.com/your-org/skeleton/internal/domain/plugin"
    "github.com/your-org/skeleton/internal/domain/storage"
)

// Re-export commonly used types for client convenience
type Plugin = plugin.Plugin

// Config represents system configuration
type Config struct {
    ServiceID     string                     `json:"serviceId"`
    StorageConfig storage.MultiStoreConfig   `json:"storage"`
}
```

### Public Config (pkg/system/config.go)
```go
package system

import (
    "github.com/your-org/skeleton/internal/domain/storage"
)

// DefaultConfig creates a default system configuration
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

### Internal fx Integration (internal/infrastructure/system/fx_bootstrap.go)
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

// systemConfig holds configuration for startup (internal type)
type systemConfig struct {
    config      *Config
    plugins     []plugin.Plugin
    registry    component.Registry
    pluginMgr   plugin.PluginManager
    eventBus    event.EventBus
    multiStore  storage.MultiStore
}

// applyDefaults creates default implementations
func (sc *systemConfig) applyDefaults() error

// startWithFx initializes and starts the system using fx
func startWithFx(config *systemConfig) error

// fx provider functions
func provideSystemService(...) (*DefaultSystemService, error)
func initializeAndStart(...) error
```

## Implementation Requirements

### 1. Use Existing Interfaces Only
- **Registry**: Use `component.Registry` interface
- **PluginManager**: Use `plugin.PluginManager` interface  
- **EventBus**: Use `event.EventBus` interface
- **MultiStore**: Use `storage.MultiStore` interface
- **DO NOT** create new interfaces

### 2. Hide fx Completely
- fx imports only in `fx_bootstrap.go`
- Public API has zero fx dependencies
- Client never sees fx types or concepts

### 3. Progressive Complexity
```go
// Simple usage
err := system.StartSystem()

// With plugins
err := system.StartSystem(
    system.WithPlugins(myPlugin1, myPlugin2),
)

// Full customization
err := system.StartSystem(
    system.WithConfig(myConfig),
    system.WithPlugins(myPlugin1, myPlugin2),
    system.WithRegistry(customRegistry),
    system.WithEventBus(customEventBus),
)
```

### 4. Default Creation Pattern
Create defaults using existing constructors:
```go
// Use existing constructors like:
component.NewDefaultRegistry()
plugin.NewDefaultPluginManager()
event.NewDefaultEventBus()
// etc.
```

### 5. Error Handling
- Use existing error patterns from the codebase
- Wrap errors with context
- Don't introduce new error types

## fx Integration Pattern

### Supply Dependencies Directly
```go
return fx.New(
    // Supply all dependencies directly
    fx.Supply(config.config),
    fx.Supply(config.registry),
    fx.Supply(config.pluginMgr),
    fx.Supply(config.eventBus),
    fx.Supply(config.multiStore),
    fx.Supply(config.plugins),
    
    // Provide the system service
    fx.Provide(provideSystemService),
    
    // Initialize and start
    fx.Invoke(initializeAndStart),
).Run()
```

### Use Updated SystemService Constructor
```go
func provideSystemService(
    config *Config,
    registry component.Registry,
    pluginMgr plugin.PluginManager,  // ← From Phase 1
    eventBus event.EventBus,
    multiStore storage.MultiStore,
) (*DefaultSystemService, error) {
    return NewDefaultSystemService(SystemServiceConfig{
        Registry:      registry,
        PluginManager: pluginMgr,  // ← From Phase 1
        EventBus:      eventBus,
        MultiStore:    multiStore,
        Config:        config,
    })
}
```

## Safety Guidelines

### Package Import Rules
- `pkg/system/` files: Only import domain interfaces, no infrastructure
- `fx_bootstrap.go`: Can import fx and infrastructure
- Follow existing import patterns

### Interface Usage
- Use interfaces exactly as defined in domain packages
- Don't create wrapper types
- Don't modify existing interfaces

### Testing Considerations
- Design should be testable (dependencies are injectable)
- Client should be able to provide test doubles
- fx integration should be isolated

## Success Criteria
- [ ] Client can call `StartSystem()` with all defaults
- [ ] Client can inject custom implementations of any dependency
- [ ] fx is completely hidden from client code
- [ ] Uses existing interfaces without modification
- [ ] Follows established project patterns
- [ ] Code compiles successfully
- [ ] Only specified files created
- [ ] Builds on Phase 1 SystemService with PluginManager

## Validation Steps
1. Verify Phase 1 completion first
2. Ensure public API compiles independently
3. Verify fx integration works with updated SystemService
4. Test default creation works for all dependencies
5. Confirm client can inject custom implementations
6. Validate error handling follows project patterns

**REMEMBER: This builds on Phase 1. The SystemService now has PluginManager integrated. Use the updated constructor signature and ensure fx provides all required dependencies including PluginManager.**