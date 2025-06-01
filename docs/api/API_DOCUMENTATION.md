# Skeleton Framework API Documentation

This document provides comprehensive documentation for the Skeleton Framework API. The framework is designed around dependency injection, component lifecycle management, and modular architecture patterns.

**Implementation Status**: ✅ **COMPLETE** - All infrastructure implementations are available and production-ready.

## Table of Contents

1. [Framework Overview](#framework-overview)
2. [Core Concepts](#core-concepts)
3. [Registry System](#registry-system)
4. [Component System](#component-system)
5. [Operation System](#operation-system)
6. [Service System](#service-system)
7. [System Integration](#system-integration)
8. [Context Management](#context-management)
9. [Configuration System](#configuration-system)
10. [Event System](#event-system)
11. [Storage System](#storage-system)
12. [Plugin System](#plugin-system)
13. [Logging System](#logging-system)
14. [Error Handling](#error-handling)
15. [Factory Functions](#factory-functions)
16. [Usage Examples](#usage-examples)

## Framework Overview

The Skeleton Framework is a **fully implemented** dependency injection and component management framework for Go applications. It provides:

- **Component-based architecture** with lifecycle management
- **Dependency injection** through a central registry
- **Operation patterns** for discrete units of work
- **Service patterns** for long-running functionality
- **Plugin system** for extensibility
- **Event-driven communication** between components
- **Multi-store persistence** layer with advanced features
- **Configuration management** with type safety

**All components have concrete implementations available** in the infrastructure layer, making the framework ready for production use.

## Core Concepts

### Identifiable Pattern

All framework entities implement the `Identifiable` interface, providing consistent identity management:

```go
type Identifiable interface {
    ID() string          // Unique identifier
    Name() string        // Human-readable name
    Description() string // Optional description
    Version() string     // Semantic version
}
```

### Component Types

The framework supports five component types:

- `TypeBasic` - Simple components with basic lifecycle
- `TypeOperation` - Components that execute discrete work units
- `TypeService` - Components providing ongoing functionality
- `TypeSystem` - System-level components
- `TypeApplication` - Application-level components

## Registry System

The Registry is the foundation of the dependency injection system, providing thread-safe storage and retrieval of components.

### Registry Interface

```go
type Registry interface {
    // Storage operations
    Register(item Identifiable) error
    Get(id string) (Identifiable, error)
    Remove(id string) error
    
    // Query operations
    Has(id string) bool
    Count() int
    List() []Identifiable
    
    // Maintenance
    Clear()
}
```

### Factory Function

```go
// NewRegistry creates a new thread-safe registry instance
func NewRegistry() Registry
```

### Identifiable Configuration

Registry items are configured using the `IdentifiableConfig` structure:

```go
type IdentifiableConfig struct {
    ID           string                 // Unique identifier
    Name         string                 // Descriptive name
    Type         IdentifiableType       // Component type
    Description  string                 // Optional description
    Version      string                 // Component version
    Dependencies []string               // Component dependencies
    Properties   map[string]interface{} // Component-specific properties
}
```

### Identifiable Factory

Components are created through factories that implement the `IdentifiableFactory` interface:

```go
type IdentifiableFactory interface {
    Create(config IdentifiableConfig) (Identifiable, error)
}
```

### Error Constants

```go
const (
    ErrItemNotFound      = "registry.item_not_found"
    ErrItemAlreadyExists = "registry.item_already_exists"
    ErrInvalidItem       = "registry.invalid_item"
)
```

## Component System

Components are the fundamental building blocks, extending the `Identifiable` interface with lifecycle management.

### Component Interface

```go
type Component interface {
    registry.Identifiable
    
    // Component properties
    Type() ComponentType
    Metadata() Metadata
    
    // Lifecycle management
    Initialize(ctx context.Context, system sys.System) error
    Dispose() error
}
```

### Component Type System

The component system uses a unified type system with the registry:

```go
// ComponentType is an alias to registry.IdentifiableType
type ComponentType = registry.IdentifiableType

// Component type constants (re-exported from registry)
const (
    TypeBasic       = registry.TypeBasic
    TypeOperation   = registry.TypeOperation
    TypeService     = registry.TypeService
    TypeSystem      = registry.TypeSystem
    TypeApplication = registry.TypeApplication
)
```

### Component Configuration

Components are configured using `ComponentConfig`, which **composes** `IdentifiableConfig` to inherit all core identity properties:

```go
// ComponentConfig composes IdentifiableConfig to inherit core identity properties
type ComponentConfig struct {
    registry.IdentifiableConfig
    // Component-specific configuration properties can be added here in the future
}

// Constructor function for convenient creation
func NewComponentConfig(id, name string, componentType ComponentType, description string) ComponentConfig
```

**Key Benefits of this Composition Pattern:**
- **Single source of truth** for identity properties (ID, Name, Description, Version, Dependencies, Properties)
- **Automatic inheritance** of new fields added to `IdentifiableConfig`
- **Type system unification** between registry and component domains
- **Consistent behavior** across the framework
- **No duplication** of configuration fields

This composition ensures that `ComponentConfig` automatically inherits all properties from `IdentifiableConfig`, including:
- `ID` - Unique identifier
- `Name` - Human-readable name  
- `Type` - Component type (unified with registry types)
- `Description` - Optional description
- `Version` - Semantic version
- `Dependencies` - Component dependencies
- `Properties` - Component-specific properties

### Factory Functions

```go
// NewBaseComponent creates a basic component implementation
func NewBaseComponent(config ComponentConfig) Component

// NewFactory creates a component factory with registry dependency
func NewFactory(registry Registry) Factory

// NewDependencyAwareComponent wraps a component with dependency management
func NewDependencyAwareComponent(base Component, registry Registry) DependencyAwareComponent

// NewLifecycleAwareComponent wraps a component with lifecycle management
func NewLifecycleAwareComponent(base Component) LifecycleAwareComponent
```

### Component Factory

```go
type Factory interface {
    Create(config ComponentConfig) (Component, error)
}
```

### Dependency-Aware Components

Components can declare and resolve dependencies:

```go
type DependencyAwareComponent interface {
    // Dependency management
    Dependencies() []string
    AddDependency(id string)
    RemoveDependency(id string)
    HasDependency(id string) bool
    
    // Dependency resolution
    ResolveDependency(id string, registrar registry.Registry) (Component, error)
    ResolveDependencies(registrar registry.Registry) (map[string]Component, error)
}
```

### Lifecycle-Aware Components

Components can manage their lifecycle state with callbacks:

```go
type LifecycleAwareComponent interface {
    State() LifecycleState
    SetState(state LifecycleState)
    OnStateChange(callback func(oldState, newState LifecycleState))
}
```

#### Lifecycle States

```go
const (
    StateCreated      LifecycleState = "created"      // Component created
    StateInitializing LifecycleState = "initializing" // Initializing
    StateInitialized  LifecycleState = "initialized"  // Successfully initialized
    StateActive       LifecycleState = "active"       // Active and operational
    StateDisposing    LifecycleState = "disposing"    // Being disposed
    StateDisposed     LifecycleState = "disposed"     // Disposed
    StateFailed       LifecycleState = "failed"       // Failed during operation
)
```

## Operation System

Operations represent discrete units of work that can be executed with input and produce output.

### Operation Interface

```go
type Operation interface {
    component.Component
    
    // Execute the operation
    Execute(ctx context.Context, input Input) (Output, error)
}

// Type aliases for flexibility
type Input interface{}
type Output interface{}
```

### Operation Configuration

```go
type OperationConfig struct {
    component.ComponentConfig
    // Operation-specific configuration properties can be added here
}

// Constructor function that ensures proper component type
func NewOperationConfig(id, name, description string) OperationConfig
```

The `NewOperationConfig` constructor automatically sets the component type to `TypeOperation`, ensuring consistency.

### Factory Functions

```go
// NewOperation creates an operation with component dependency
func NewOperation(component Component) Operation

// NewOperationFactory creates an operation factory with component factory dependency
func NewOperationFactory(componentFactory Factory) OperationFactory
```

### Operation Factory

```go
type OperationFactory interface {
    component.Factory
    
    // Create operation from specific config
    CreateOperation(config OperationConfig) (Operation, error)
}
```

## Service System

Services provide ongoing functionality with start/stop lifecycle management.

### Service Interface

```go
type Service interface {
    component.Component
    
    // Service lifecycle
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    Status() ServiceStatus
}
```

### Service Status

```go
const (
    StatusStopped  ServiceStatus = "stopped"  // Service is stopped
    StatusStarting ServiceStatus = "starting" // Service is starting
    StatusRunning  ServiceStatus = "running"  // Service is running
    StatusStopping ServiceStatus = "stopping" // Service is stopping
    StatusFailed   ServiceStatus = "failed"   // Service has failed
)
```

### Service Configuration

```go
type ServiceConfig struct {
    component.ComponentConfig
    // Service-specific configuration properties can be added here
}

// Constructor function that ensures proper component type
func NewServiceConfig(id, name, description string) ServiceConfig
```

The `NewServiceConfig` constructor automatically sets the component type to `TypeService`, ensuring consistency.

### Factory Functions

```go
// NewService creates a service with component dependency
func NewService(component Component) Service

// NewServiceFactory creates a service factory with component factory dependency
func NewServiceFactory(componentFactory Factory) ServiceFactory
```

### Service Factory

```go
type ServiceFactory interface {
    component.Factory
    
    // Create service from specific config
    CreateService(config ServiceConfig) (Service, error)
}
```

### Error Constants

```go
const (
    ErrServiceStart    = "service.start_failed"
    ErrServiceStop     = "service.stop_failed"
    ErrServiceNotFound = "service.not_found"
)
```

## System Integration

The System interface provides centralized access to all framework resources and operations.

### System Interface

```go
type System interface {
    // Resource access
    Registry() registry.Registry
    PluginManager() plugin.PluginManager
    EventBus() event.EventBus
    Configuration() config.Configuration
    Store() storage.MultiStore
    
    // System operations (using IDs, not component interfaces)
    // Note: ctx parameter uses the framework's context.Context interface (not Go's standard context)
    ExecuteOperation(ctx context.Context, operationID string, input interface{}) (interface{}, error)
    StartService(ctx context.Context, serviceID string) error
    StopService(ctx context.Context, serviceID string) error
    
    // System state
    IsRunning() bool
    IsInitialized() bool
}
```

### Factory Function

```go
// NewSystem creates a system with all required dependencies
func NewSystem(
    registry Registry,
    pluginManager PluginManager,
    eventBus EventBus,
    configuration Configuration,
    store MultiStore,
) System
```

### Error Constants

```go
const (
    ErrSystemNotInitialized = "system.not_initialized"
    ErrSystemNotStarted     = "system.not_started"
    ErrOperationNotFound    = "system.operation_not_found"
    ErrOperationFailed      = "system.operation_failed"
    ErrServiceNotFound      = "system.service_not_found"
    ErrServiceStart         = "system.service_start_failed"
    ErrServiceStop          = "system.service_stop_failed"
)
```

## Context Management

The framework provides its own context interface for execution management, defined in `skeleton/internal/domain/context`. **This is NOT Go's standard context package** - it's the framework's own context interface designed for component lifecycle and execution management.

### Framework Context Interface

```go
// Framework's own context interface (not Go's standard context)
type Context interface {
    // Value management
    Value(key interface{}) interface{}
    WithValue(key, value interface{}) Context
    
    // Deadline and cancellation
    Deadline() (time.Time, bool)
    Done() <-chan struct{}
    Err() error
}
```

### Factory Functions

```go
// NewContext creates a new framework context
func NewContext() Context

// WrapContext wraps a standard Go context into framework context
func WrapContext(ctx Context) Context
```

**Important Notes:**
- **Framework Context**: All system operations use the framework's `context.Context` interface
- **Not Go Context**: This is NOT the standard Go `context.Context` package
- **Consistent Usage**: All component initialization, system operations, and lifecycle methods use this framework context
- **Import Path**: Import from `github.com/fintechain/skeleton/pkg/context`

### Context Usage in System Operations

All system operations consistently use the framework context:

```go
// All these methods use framework context.Context (not Go's standard context)
ExecuteOperation(ctx context.Context, operationID string, input interface{}) (interface{}, error)
StartService(ctx context.Context, serviceID string) error
StopService(ctx context.Context, serviceID string) error

// Component lifecycle also uses framework context
Initialize(ctx context.Context, system System) error
```

## Configuration System

Type-safe configuration management with support for multiple data types.

### Configuration Interface

```go
type Configuration interface {
    // String values
    GetString(key string) string
    GetStringDefault(key, defaultValue string) string
    
    // Integer values
    GetInt(key string) (int, error)
    GetIntDefault(key string, defaultValue int) int
    
    // Boolean values
    GetBool(key string) (bool, error)
    GetBoolDefault(key string, defaultValue bool) bool
    
    // Duration values
    GetDuration(key string) (time.Duration, error)
    GetDurationDefault(key string, defaultValue time.Duration) time.Duration
    
    // Object deserialization
    GetObject(key string, result interface{}) error
    
    // Key existence
    Exists(key string) bool
}
```

### Factory Function

```go
// NewConfiguration creates a configuration with source dependencies
func NewConfiguration(sources ...ConfigurationSource) Configuration
```

### Configuration Source

```go
type ConfigurationSource interface {
    LoadConfig() error
    GetValue(key string) (interface{}, bool)
}
```

### Error Constants

```go
const (
    ErrConfigNotFound   = "config.not_found"
    ErrConfigWrongType  = "config.wrong_type"
    ErrConfigLoadFailed = "config.load_failed"
)
```

## Event System

Publish-subscribe event system for component communication.

### Event Structure

```go
type Event struct {
    Topic   string                 // Event topic/type
    Source  string                 // Component that generated the event
    Time    time.Time              // When the event occurred
    Payload map[string]interface{} // Event data
}
```

### Event Bus Interface

```go
type EventBus interface {
    // Publication
    Publish(topic string, data interface{})
    
    // Subscription
    Subscribe(topic string, handler EventHandler) Subscription
    SubscribeAsync(topic string, handler EventHandler) Subscription
    
    // Control
    WaitAsync()
}

type EventHandler func(event *Event)
```

### Factory Function

```go
// NewEventBus creates a new event bus
func NewEventBus() EventBus
```

### Subscription Interface

```go
type Subscription interface {
    Cancel()        // Cancel the subscription
    Topic() string  // Get the subscription topic
}
```

## Storage System

Multi-store persistence layer with pluggable storage engines and advanced features.

### MultiStore Interface

```go
type MultiStore interface {
    // Store management
    GetStore(name string) (Store, error)
    CreateStore(name, engine string, config Config) error
    DeleteStore(name string) error
    ListStores() []string
    StoreExists(name string) bool
    
    // Bulk operations
    CloseAll() error
    
    // Engine configuration
    SetDefaultEngine(engine string)
    GetDefaultEngine() string
    
    // Engine management
    RegisterEngine(engine Engine) error
    ListEngines() []string
    GetEngine(name string) (Engine, error)
}
```

### Factory Function

```go
// NewMultiStore creates a new multi-store instance
func NewMultiStore() MultiStore
```

### Store Interface

```go
type Store interface {
    // Basic CRUD operations
    Get(key []byte) ([]byte, error)
    Set(key, value []byte) error
    Delete(key []byte) error
    Has(key []byte) (bool, error)
    
    // Iteration over all key-value pairs
    Iterate(fn func(key, value []byte) bool) error
    
    // Resource cleanup
    Close() error
    
    // Store metadata
    Name() string
    Path() string
}
```

### Transaction Support

```go
type Transactional interface {
    BeginTx() (Transaction, error)
    SupportsTransactions() bool
}

type Transaction interface {
    Store // Embed Store interface
    
    // Transaction control
    Commit() error
    Rollback() error
    IsActive() bool
}
```

### Engine Interface

```go
type Engine interface {
    Name() string
    Create(name, path string, config Config) (Store, error)
    Open(name, path string) (Store, error)
    Capabilities() Capabilities
}
```

### Engine Capabilities

```go
type Capabilities struct {
    Transactions  bool // Supports atomic transactions
    Versioning    bool // Supports versioning/snapshots
    RangeQueries  bool // Supports efficient range queries
    Persistence   bool // Persists data to disk
    Compression   bool // Supports data compression
}
```

### Range Query Support

```go
type RangeQueryable interface {
    IterateRange(start, end []byte, ascending bool, fn func(key, value []byte) bool) error
    SupportsRangeQueries() bool
}
```

### Storage Configuration

```go
type Config map[string]interface{}
```

### Storage Error Constants

The storage system uses **string constants** (consistent with all other framework domains):

```go
const (
    ErrKeyNotFound     = "storage.key_not_found"
    ErrStoreNotFound   = "storage.store_not_found"
    ErrStoreClosed     = "storage.store_closed"
    ErrStoreExists     = "storage.store_exists"
    ErrEngineNotFound  = "storage.engine_not_found"
    ErrTxNotActive     = "storage.transaction_not_active"
    ErrTxReadOnly      = "storage.transaction_read_only"
    ErrVersionNotFound = "storage.version_not_found"
    ErrInvalidConfig   = "storage.invalid_config"
)
```

**Error Handling Pattern**: Use standard string comparison for storage errors, consistent with other framework domains:

```go
// Check for specific storage errors using string comparison
if err != nil && err.Error() == storage.ErrKeyNotFound {
    // Handle key not found
}
```

## Plugin System

Dynamic plugin loading and management for framework extensibility.

### Plugin Interface

```go
type Plugin interface {
    registry.Identifiable
    
    // Plugin lifecycle
    Load(ctx context.Context, registrar registry.Registry) error
    Unload(ctx context.Context) error
}
```

### Plugin Manager Interface

```go
type PluginManager interface {
    // Discovery
    Discover(ctx context.Context, location string) ([]PluginInfo, error)
    
    // Lifecycle
    Load(ctx context.Context, id string, registrar registry.Registry) error
    Unload(ctx context.Context, id string) error
    
    // Information
    ListPlugins() []PluginInfo
    GetPlugin(id string) (Plugin, error)
}
```

### Factory Function

```go
// NewPluginManager creates a plugin manager with filesystem dependency
func NewPluginManager(filesystem FileSystem) PluginManager
```

### Plugin Info

```go
type PluginInfo struct {
    ID          string                 // Unique identifier
    Name        string                 // Human-readable name
    Version     string                 // Semantic version
    Description string                 // Plugin description
    Author      string                 // Plugin author/maintainer
    Metadata    map[string]interface{} // Additional metadata
}
```

### Error Constants

```go
const (
    ErrPluginNotFound  = "plugin.not_found"
    ErrPluginLoad      = "plugin.load_failed"
    ErrPluginUnload    = "plugin.unload_failed"
    ErrPluginDiscovery = "plugin.discovery_failed"
)
```

## Logging System

Structured logging interface for framework-wide logging capabilities.

### Logger Interface

```go
type Logger interface {
    Debug(msg string, args ...interface{})
    Info(msg string, args ...interface{})
    Warn(msg string, args ...interface{})
    Error(msg string, args ...interface{})
}
```

### Factory Function

```go
// NewLogger creates a new logger instance
func NewLogger() Logger
```

### Error Constants

```go
const (
    ErrLoggerNotAvailable = "logging.logger_not_available"
    ErrInvalidLogLevel    = "logging.invalid_log_level"
)
```

## Error Handling

The framework uses structured error handling with **consistent string constants** across all domains.

### Error Code Patterns

All framework domains use **string constants** for error identification:

- **Registry errors**: `registry.*` (e.g., `registry.item_not_found`)
- **Component errors**: `component.*` (e.g., `component.invalid_config`)
- **Service errors**: `service.*` (e.g., `service.start_failed`)
- **System errors**: `system.*` (e.g., `system.not_initialized`)
- **Config errors**: `config.*` (e.g., `config.not_found`)
- **Plugin errors**: `plugin.*` (e.g., `plugin.load_failed`)
- **Storage errors**: `storage.*` (e.g., `storage.key_not_found`)

### Error Handling Best Practices

1. **Use defined error constants** for consistent error identification
2. **Use string comparison** for error checking across all domains
3. **Wrap errors** with additional context when propagating
4. **Consistent pattern**: All domains use string constants, checked via `err.Error() == ErrorConstant`

### Error Checking Pattern

```go
// Standard error checking pattern across all framework domains
if err != nil && err.Error() == registry.ErrItemNotFound {
    // Handle registry item not found
}

if err != nil && err.Error() == storage.ErrKeyNotFound {
    // Handle storage key not found
}

if err != nil && err.Error() == service.ErrServiceStart {
    // Handle service start failure
}
```

## Factory Functions

All infrastructure implementations are accessible through factory functions in the public API:

### Core Infrastructure
```go
// Registry
func NewRegistry() Registry

// Context
func NewContext() Context
func WrapContext(ctx Context) Context

// Components
func NewBaseComponent(config ComponentConfig) Component
func NewFactory(registry Registry) Factory
func NewDependencyAwareComponent(base Component, registry Registry) DependencyAwareComponent
func NewLifecycleAwareComponent(base Component) LifecycleAwareComponent
```

### Specialized Components
```go
// Operations
func NewOperation(component Component) Operation
func NewOperationFactory(componentFactory Factory) OperationFactory

// Services
func NewService(component Component) Service
func NewServiceFactory(componentFactory Factory) ServiceFactory
```

### System Integration
```go
// System
func NewSystem(
    registry Registry,
    pluginManager PluginManager,
    eventBus EventBus,
    configuration Configuration,
    store MultiStore,
) System

// Plugin Management
func NewPluginManager(filesystem FileSystem) PluginManager

// Configuration
func NewConfiguration(sources ...ConfigurationSource) Configuration
```

### Supporting Infrastructure
```go
// Event System
func NewEventBus() EventBus

// Storage
func NewMultiStore() MultiStore

// Logging
func NewLogger() Logger
```

**Note**: All factory functions correspond to concrete implementations in the infrastructure layer and are ready for production use.

## Usage Examples

### Basic Component Registration

```go
import "github.com/fintechain/skeleton/pkg/registry"
import "github.com/fintechain/skeleton/pkg/component"

// Create registry
registry := registry.NewRegistry()

// Create component
config := component.NewComponentConfig("my-component", "My Component", component.TypeBasic, "A basic component")
component := component.NewBaseComponent(config)

// Register component
err := registry.Register(component)
if err != nil {
    // Handle registration error
}
```

### System Setup with Dependencies

```go
import (
    "github.com/fintechain/skeleton/pkg/registry"
    "github.com/fintechain/skeleton/pkg/system"
    "github.com/fintechain/skeleton/pkg/event"
    "github.com/fintechain/skeleton/pkg/storage"
    "github.com/fintechain/skeleton/pkg/config"
    "github.com/fintechain/skeleton/pkg/plugin"
)

// Create all dependencies
registry := registry.NewRegistry()
eventBus := event.NewEventBus()
store := storage.NewMultiStore()
configuration := config.NewConfiguration() // with sources
pluginManager := plugin.NewPluginManager(filesystem) // with filesystem

// Create system
system := system.NewSystem(registry, pluginManager, eventBus, configuration, store)
```

### Operation Execution

```go
// Get operation from registry
op, err := registry.Get("my-operation")
if err != nil {
    // Handle not found
}

// Cast to operation
operation := op.(operation.Operation)

// Execute operation (using framework context)
ctx := context.NewContext()
result, err := operation.Execute(ctx, inputData)
if err != nil {
    // Handle execution error
}
```

### Service Lifecycle

```go
// Get service from registry
svc, err := registry.Get("my-service")
if err != nil {
    // Handle not found
}

// Cast to service
service := svc.(service.Service)

// Start service (using framework context)
ctx := context.NewContext()
err = service.Start(ctx)
if err != nil {
    // Handle start error
}

// Check status
status := service.Status()
if status == service.StatusRunning {
    // Service is running
}
```

### Storage Operations with Engine Registration

```go
// Create multistore
multiStore := storage.NewMultiStore()

// Register a memory engine (implementation needed)
memoryEngine := NewMemoryEngine() // Custom implementation
err := multiStore.RegisterEngine(memoryEngine)

// Create store
err = multiStore.CreateStore("my-store", "memory", storage.Config{})

// Get store and perform operations
store, err := multiStore.GetStore("my-store")
if err == nil {
    err = store.Set([]byte("key"), []byte("value"))
    value, err := store.Get([]byte("key"))
}
```

## Public API Access

All these interfaces are re-exported through the public API in `skeleton/pkg/`:

```go
import (
    "github.com/fintechain/skeleton/pkg/component"
    "github.com/fintechain/skeleton/pkg/registry"
    "github.com/fintechain/skeleton/pkg/operation"
    "github.com/fintechain/skeleton/pkg/service"
    "github.com/fintechain/skeleton/pkg/system"
    "github.com/fintechain/skeleton/pkg/storage"
    "github.com/fintechain/skeleton/pkg/config"
    "github.com/fintechain/skeleton/pkg/event"
    "github.com/fintechain/skeleton/pkg/plugin"
    "github.com/fintechain/skeleton/pkg/logging"
    "github.com/fintechain/skeleton/pkg/context"
)
```

This provides a clean separation between the internal domain model and the public API surface, with **complete infrastructure implementations** available for all components. 