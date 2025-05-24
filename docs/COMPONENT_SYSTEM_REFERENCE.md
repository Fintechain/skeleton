# Component System Reference Documentation

## 1. Introduction

### 1.1 Purpose and Scope
This document serves as a comprehensive reference for the redesigned component system. It covers the architecture, design principles, implementation details, and usage patterns for all aspects of the system.

### 1.2 Target Audience
- **Developers**: Engineers working with or extending the component system
- **Architects**: Technical leaders making design decisions
- **Maintainers**: Team members responsible for ongoing maintenance

### 1.3 How to Use This Document
This reference is organized into logical sections that progress from high-level concepts to detailed implementation. New users should start with the System Overview and Core Domain Model sections, while developers looking for specific details can navigate directly to relevant sections.

## 2. System Overview

### 2.1 High-Level Architecture
The component system is built on a set of clean, focused interfaces with clear responsibilities:

```
┌─────────────────────────────────────────────────────────┐
│                      Applications                        │
└───────────────────────────┬─────────────────────────────┘
                            │
┌───────────────────────────▼─────────────────────────────┐
│                 Core Component System                    │
├─────────────────┬─────────────────┬─────────────────────┤
│  Components     │ Registry        │ Factory             │
└─────────────────┴─────────────────┴─────────────────────┘
┌─────────────────┬─────────────────┬─────────────────────┐
│  Operations     │ Services        │ Plugins             │
└─────────────────┴─────────────────┴─────────────────────┘
┌─────────────────┬─────────────────┬─────────────────────┐
│  Event Bus      │ Context         │ Config/Logging      │
└─────────────────┴─────────────────┴─────────────────────┘
```

### 2.2 Design Philosophy and Principles

The system follows several key design principles:

- **Composition Over Inheritance**: Components use interface composition and field composition rather than struct embedding
- **Interface Segregation**: Interfaces are small and focused on specific behaviors
- **Dependency Injection**: Components declare dependencies explicitly
- **Event-Driven Communication**: Components communicate through events when possible
- **Explicit Error Handling**: Errors include context and are domain-specific
- **Idiomatic Go**: Following Go's conventions and best practices

### 2.3 Key Concepts and Terminology

- **Component**: The fundamental building block of the system, with identity, lifecycle, and metadata
- **Operation**: A component that executes discrete units of work with inputs and outputs
- **Service**: A component providing ongoing functionality with start/stop lifecycle
- **Plugin**: A container for components that extends the system dynamically
- **Registry**: A repository that tracks all registered components
- **Factory**: A mechanism for creating components from configuration
- **Event Bus**: A system for decoupled communication between components

### 2.4 Package Organization

The system is organized into these key packages:

- **domain/component/**: Core component interfaces and types
- **domain/operation/**: Operation-specific interfaces and implementations
- **domain/service/**: Service-specific interfaces and implementations
- **domain/plugin/**: Plugin-related interfaces and implementations
- **domain/system/**: System service interfaces and types
- **infrastructure/event/**: Event bus and handling
- **infrastructure/context/**: Context management
- **infrastructure/config/**: Configuration management
- **infrastructure/logging/**: Logging facilities
- **infrastructure/system/**: System service implementation and utilities

## 3. Core Domain Model

### 3.1 Component Model

#### 3.1.1 Component Interface
The fundamental building block is the `Component` interface defined in `domain/component/component.go`:

```go
type Component interface {
    // Identity
    ID() string
    Name() string
    Type() ComponentType
    
    // Metadata
    Metadata() Metadata
    
    // Lifecycle
    Initialize(ctx Context) error
    Dispose() error
}
```

#### 3.1.2 Base Component
The `BaseComponent` provides a foundation for component implementations:

- Provides core identity properties (ID, name, type)
- Manages metadata
- Implements basic lifecycle methods
- Defined in `domain/component/base_component.go`

Usage example:
```go
base := component.NewBaseComponent("component-1", "Example Component", component.TypeBasic)
err := base.Initialize(ctx)
```

#### 3.1.3 Default Component
The `DefaultComponent` is a complete implementation of the `Component` interface:

- Embeds a `BaseComponent` for core functionality
- Provides default implementations for all methods
- Suitable for simple component needs
- Defined in `domain/component/default_component.go`

Usage example:
```go
comp := component.NewDefaultComponent("component-1", "Example Component", component.TypeBasic)
err := comp.Initialize(ctx)
```

#### 3.1.4 Lifecycle-Aware Component
The `LifecycleAwareComponent` extends the basic component with enhanced lifecycle management:

- Maintains component state (created, initialized, disposed)
- Enforces proper state transitions
- Provides state inquiry methods
- Defined in `domain/component/lifecycle_aware.go`

#### 3.1.5 Dependency-Aware Component
The `DependencyAwareComponent` adds dependency management capabilities:

- Tracks component dependencies
- Ensures dependencies are initialized before the component
- Provides methods to add and remove dependencies
- Defined in `domain/component/dependency_aware.go`

### 3.2 Registry System

#### 3.2.1 Registry Interface
The `Registry` manages component registration and discovery:

```go
type Registry interface {
    Register(Component) error
    Unregister(id string) error
    Get(id string) (Component, error)
    FindByType(componentType ComponentType) []Component
    FindByMetadata(key string, value interface{}) []Component
    Initialize(ctx Context) error
    Shutdown() error
}
```

Defined in `domain/component/registry.go`.

#### 3.2.2 Registry Implementation
The default registry implementation provides:

- Thread-safe component storage
- Component lookup by ID, type, and metadata
- Component lifecycle management
- Dependency resolution
- Defined in `domain/component/registry_impl.go`

Usage example:
```go
registry := component.NewDefaultRegistry()
registry.Register(myComponent)
comp, err := registry.Get("component-id")
```

### 3.3 Factory System

#### 3.3.1 Factory Interface
The `Factory` creates components from configuration:

```go
type Factory interface {
    Create(config ComponentConfig) (Component, error)
}
```

Defined in `domain/component/factory.go`.

#### 3.3.2 Factory Implementation
The default factory implementation:

- Creates components based on their type
- Configures components from provided configuration
- Validates configuration before component creation
- Defined in `domain/component/factory_impl.go`

Usage example:
```go
factory := component.NewDefaultFactory()
config := component.ComponentConfig{
    ID: "component-1",
    Name: "Example Component",
    Type: string(component.TypeBasic),
}
comp, err := factory.Create(config)
```

### 3.4 Metadata Handling
Components can store arbitrary metadata as key-value pairs:

```go
type Metadata map[string]interface{}
```

Metadata is used for component discovery, configuration, and dynamic behavior.

### 3.5 Error Model
The system uses a domain-specific error model defined in `domain/component/errors.go`:

```go
type Error struct {
    Code    string
    Message string
    Details map[string]interface{}
    Cause   error
}
```

Common error categories include:
- Registration errors (`component.registration_failed`)
- Initialization errors (`component.initialization_failed`)
- Dependency errors (`component.dependency_failed`)

## 4. Specialized Components

### 4.1 Operations

#### 4.1.1 Operation Model
Operations execute discrete units of work with inputs and outputs:

```go
type Operation interface {
    component.Component
    Execute(ctx component.Context, input Input) (Output, error)
}
```

Defined in `domain/operation/operation.go`.

#### 4.1.2 Base Operation
The `BaseOperation` provides foundation for operation implementations:

- Embeds a component for identity and lifecycle
- Provides structure for execution logic
- Defined in `domain/operation/base_operation.go`

#### 4.1.3 Default Operation
The `DefaultOperation` is a complete implementation of the `Operation` interface:

- Integrates with component lifecycle
- Provides input validation
- Handles execution context
- Defined in `domain/operation/default_operation.go`

Usage example:
```go
op := operation.NewDefaultOperation("op-1", "Example Operation")
result, err := op.Execute(ctx, myInput)
```

#### 4.1.4 Operation Pipelines
The system supports operation pipelines for chaining operations:

- Sequential execution of multiple operations
- Output of one operation becomes input to the next
- Error handling between pipeline stages
- Defined in `domain/operation/pipeline.go`

Usage example:
```go
pipeline := operation.NewPipeline("pipeline", []operation.Operation{op1, op2, op3})
result, err := pipeline.Execute(ctx, initialInput)
```

### 4.2 Services

#### 4.2.1 Service Model
Services provide ongoing functionality with a start/stop lifecycle:

```go
type Service interface {
    component.Component
    Start(ctx component.Context) error
    Stop(ctx component.Context) error
    Status() ServiceStatus
}
```

Defined in `domain/service/service.go`.

#### 4.2.2 Base Service
The `BaseService` provides foundation for service implementations:

- Manages service state transitions
- Provides status reporting
- Implements core service lifecycle
- Defined in `domain/service/base_service.go`

#### 4.2.3 Default Service
The `DefaultService` is a complete implementation of the `Service` interface:

- Composes a `BaseService` for core functionality
- Delegates to the base service appropriately
- Provides customizable start/stop behavior
- Defined in `domain/service/default_service.go`

Usage example:
```go
svc := service.NewDefaultService("svc-1", "Example Service")
err := svc.Start(ctx)
```

#### 4.2.4 Health Monitoring
Services support health monitoring capabilities:

- Health check registration and execution
- Health status reporting
- Automatic background health monitoring
- Defined in `domain/service/health.go`

Usage example:
```go
svc.RegisterHealthCheck("database", func(ctx component.Context) error {
    return checkDatabaseConnection()
})
```

### 4.3 Plugins

#### 4.3.1 Plugin Model
Plugins extend the system with dynamically loaded components:

```go
type Plugin interface {
    ID() string
    Version() string
    Load(ctx component.Context, registry component.Registry) error
    Unload(ctx component.Context) error
    Components() []component.Component
}
```

Defined in `domain/plugin/plugin.go`.

#### 4.3.2 Default Plugin
The `DefaultPlugin` is a standard implementation of the `Plugin` interface:

- Manages plugin identity and versioning
- Handles component registration and unregistration
- Provides lifecycle management
- Defined in `domain/plugin/default_plugin.go`

Usage example:
```go
plugin := plugin.NewDefaultPlugin("plugin-1", "1.0.0", []component.Component{comp1, comp2})
err := plugin.Load(ctx, registry)
```

#### 4.3.3 Plugin Manager
The `PluginManager` provides plugin discovery and lifecycle management:

```go
type PluginManager interface {
    Discover(location string) ([]PluginInfo, error)
    Load(id string, registry component.Registry) error
    Unload(id string) error
    ListPlugins() []PluginInfo
    GetPlugin(id string) (Plugin, error)
}
```

Defined in `domain/plugin/plugin.go` and implemented in `domain/plugin/plugin_manager.go`.

Usage example:
```go
manager := plugin.NewDefaultPluginManager()
plugins, err := manager.Discover("./plugins")
err = manager.Load("plugin-1", registry)
```

### 4.4 System Service

#### 4.4.1 System Service Model
The SystemService serves as the central coordinating service of the application:

```go
type SystemService interface {
    service.Service

    // Core components access
    Registry() component.Registry
    EventBus() event.EventBus
    Configuration() config.Configuration
    Store() storage.MultiStore

    // Operations
    ExecuteOperation(ctx component.Context, operationID string, input interface{}) (interface{}, error)
    StartService(ctx component.Context, serviceID string) error
    StopService(ctx component.Context, serviceID string) error
}
```

Defined in `domain/system/system.go`.

#### 4.4.2 System Service Configuration
The `SystemServiceConfig` defines configuration options for the SystemService:

```go
type SystemServiceConfig struct {
    // Service identity
    ServiceID string `json:"serviceId" default:"system"`

    // Feature flags
    EnableOperations bool `json:"enableOperations" default:"true"`
    EnableServices   bool `json:"enableServices" default:"true"`
    EnablePlugins    bool `json:"enablePlugins" default:"true"`
    EnableEventLog   bool `json:"enableEventLog" default:"true"`

    // Storage configuration
    StorageConfig storage.MultiStoreConfig `json:"storage"`
}
```

This configuration allows customization of system behavior with feature flags and component-specific settings.

#### 4.4.3 Default System Service
The `DefaultSystemService` provides a standard implementation of the SystemService interface:

- Manages the lifecycle of all registered components
- Coordinates component initialization, startup, and shutdown
- Handles operation execution and routing
- Manages service discovery and lifecycle
- Provides access to core system components (Registry, EventBus, etc.)
- Defined in `infrastructure/system/default_system_service.go`

#### 4.4.4 System Service Factory
The System Service Factory creates system service instances from configuration:

```go
// Factory creates system service instances
type Factory struct {
    registry      component.Registry
    eventBus      event.EventBus
    configuration config.Configuration
    multiStore    storage.MultiStore
}

// CreateSystemService creates a new system service instance
func (f *Factory) CreateSystemService(config *system.SystemServiceConfig) (system.SystemService, error)
```

The factory follows the constructor injection pattern for dependencies and configures the system service based on the provided configuration.

#### 4.4.5 System Service Builder
The System Service Builder provides a fluent interface for constructing system services:

```go
// Builder helps construct a complete system with default components
type Builder struct {
    serviceID     string
    logger        logging.Logger
    configuration config.Configuration
    registry      component.Registry
    eventBus      event.EventBus
    multiStore    storage.MultiStore
}

// Method chaining API
func (b *Builder) WithLogger(logger logging.Logger) *Builder
func (b *Builder) WithConfiguration(configuration config.Configuration) *Builder
func (b *Builder) WithRegistry(registry component.Registry) *Builder
func (b *Builder) WithEventBus(eventBus event.EventBus) *Builder
func (b *Builder) WithMultiStore(multiStore storage.MultiStore) *Builder
func (b *Builder) Build() (system.SystemService, error)
```

The builder pattern simplifies the creation and configuration of system services with a clean, readable API.

#### 4.4.6 Configuration Loader
The `ConfigLoader` loads system configuration from a configuration source:

```go
// ConfigLoader loads system service configuration
type ConfigLoader struct {
    config config.Configuration
}

// LoadSystemConfig loads the system service configuration
func (l *ConfigLoader) LoadSystemConfig() (*system.SystemServiceConfig, error)
```

This utility provides sensible defaults while allowing configuration overrides from external sources.

#### 4.4.7 System Event Topics
The system service publishes standard events for important state changes:

- `system.initialized`: Published when the system is initialized
- `system.started`: Published when the system is started
- `system.stopped`: Published when the system is stopped
- `system.operation.executed`: Published when an operation is executed successfully
- `system.operation.failed`: Published when an operation execution fails
- `system.service.started`: Published when a service is started
- `system.service.stopped`: Published when a service is stopped
- `system.service.failed`: Published when a service fails to start or stop

#### 4.4.8 System Error Codes
Common error codes for system operations:

- `system.not_initialized`: System not initialized
- `system.not_started`: System not started
- `system.operation_not_found`: Operation not found
- `system.operation_failed`: Operation execution failed
- `system.service_not_found`: Service not found
- `system.service_start_failed`: Service failed to start
- `system.service_stop_failed`: Service failed to stop

## 5. Infrastructure

### 5.1 Event System

#### 5.1.1 Event Model
Events facilitate decoupled communication between components:

```go
type Event struct {
    Topic   string
    Source  string
    Time    time.Time
    Payload map[string]interface{}
}
```

Defined in `infrastructure/event/event_bus.go`.

#### 5.1.2 Event Bus
The `EventBus` handles event publication and subscription:

```go
type EventBus interface {
    Publish(topic string, data interface{})
    Subscribe(topic string, handler EventHandler) Subscription
    SubscribeAsync(topic string, handler EventHandler) Subscription
    WaitAsync()
}
```

Implemented in `infrastructure/event/event_bus_impl.go`.

Usage example:
```go
bus := event.NewDefaultEventBus()
sub := bus.Subscribe("component.initialized", func(e *event.Event) {
    fmt.Printf("Component initialized: %s\n", e.Source)
})
bus.Publish("component.initialized", map[string]interface{}{"componentId": "comp-1"})
```

### 5.2 Context Management
The system provides context management based on Go's context.Context:

- Wraps standard Go context in a domain-specific interface
- Supports values, deadlines, and cancellation
- Provides utility functions for common operations
- Defined in `infrastructure/context/context.go`

Usage example:
```go
ctx := context.Background()
compCtx := infrastructure.WrapContext(ctx)
compCtx = compCtx.WithValue("userId", "user-123")
```

### 5.3 Configuration System
The configuration system supports component configuration:

- Provides typed configuration access
- Supports default values
- Validates configuration
- Defined in `infrastructure/config/config.go`

### 5.4 Logging Facilities
The logging system provides structured logging for the component system:

- Log levels (debug, info, warn, error)
- Contextual logging with component information
- Support for structured logging
- Defined in `infrastructure/logging/logging.go`

## 6. API Reference

### 6.1 Component API

| Method | Description |
|--------|-------------|
| `ID() string` | Returns the unique identifier of the component |
| `Name() string` | Returns the human-readable name of the component |
| `Type() ComponentType` | Returns the type of the component |
| `Metadata() Metadata` | Returns the component metadata |
| `Initialize(ctx Context) error` | Initializes the component |
| `Dispose() error` | Disposes the component and cleans up resources |

### 6.2 Operation API

| Method | Description |
|--------|-------------|
| *All Component methods* | Operations implement the Component interface |
| `Execute(ctx Context, input Input) (Output, error)` | Executes the operation with given input |

### 6.3 Service API

| Method | Description |
|--------|-------------|
| *All Component methods* | Services implement the Component interface |
| `Start(ctx Context) error` | Starts the service |
| `Stop(ctx Context) error` | Stops the service |
| `Status() ServiceStatus` | Returns the current status of the service |
| `RegisterHealthCheck(name string, check HealthCheck) error` | Registers a health check |
| `PerformHealthChecks(ctx Context) map[string]error` | Executes all health checks |

### 6.4 Plugin API

| Method | Description |
|--------|-------------|
| `ID() string` | Returns the plugin identifier |
| `Version() string` | Returns the plugin version |
| `Load(ctx Context, registry Registry) error` | Loads the plugin |
| `Unload(ctx Context) error` | Unloads the plugin |
| `Components() []Component` | Returns the components provided by the plugin |

### 6.5 System Service API

#### 6.5.1 SystemService Interface

| Method | Description |
|--------|-------------|
| *All Service methods* | SystemService implements the Service interface |
| `Registry() component.Registry` | Returns the component registry |
| `EventBus() event.EventBus` | Returns the event bus |
| `Configuration() config.Configuration` | Returns the system configuration |
| `Store() storage.MultiStore` | Returns the multi-store |
| `ExecuteOperation(ctx Context, operationID string, input interface{}) (interface{}, error)` | Executes an operation with the specified ID |
| `StartService(ctx Context, serviceID string) error` | Starts a service with the specified ID |
| `StopService(ctx Context, serviceID string) error` | Stops a service with the specified ID |

#### 6.5.2 Factory API

| Method | Description |
|--------|-------------|
| `NewFactory(registry, eventBus, configuration, multiStore)` | Creates a factory with the specified dependencies |
| `CreateSystemService(config *SystemServiceConfig) (SystemService, error)` | Creates a system service from configuration |

#### 6.5.3 Builder API

| Method | Description |
|--------|-------------|
| `NewBuilder(serviceID string)` | Creates a new builder for the specified service ID |
| `WithLogger(logger Logger)` | Sets the logger and returns the builder |
| `WithConfiguration(config Configuration)` | Sets the configuration and returns the builder |
| `WithRegistry(registry Registry)` | Sets the registry and returns the builder |
| `WithEventBus(eventBus EventBus)` | Sets the event bus and returns the builder |
| `WithMultiStore(store MultiStore)` | Sets the multi-store and returns the builder |
| `Build()` | Builds and returns the system service |

#### 6.5.4 ConfigLoader API

| Method | Description |
|--------|-------------|
| `NewConfigLoader(config Configuration)` | Creates a new config loader with the specified configuration |
| `LoadSystemConfig() (*SystemServiceConfig, error)` | Loads and returns the system configuration |

### 6.6 Infrastructure API

| Component | Key Methods |
|-----------|------------|
| EventBus | `Publish`, `Subscribe`, `SubscribeAsync`, `WaitAsync` |
| Context | `Value`, `WithValue`, `Deadline`, `Done`, `Err` |
| Config | `GetString`, `GetInt`, `GetBool`, `GetDuration` |
| Logging | `Debug`, `Info`, `Warn`, `Error` |

## 7. Usage Patterns

### 7.1 Component Creation and Registration

```go
// Create a component
comp := component.NewDefaultComponent("comp-1", "Example Component", component.TypeBasic)

// Register with registry
registry := component.NewDefaultRegistry()
err := registry.Register(comp)
if err != nil {
    return err
}

// Initialize components
err = registry.Initialize(ctx)
if err != nil {
    return err
}
```

### 7.2 Component Lifecycle Management

```go
// Create and initialize
comp := component.NewDefaultComponent("comp-1", "Example Component", component.TypeBasic)
err := comp.Initialize(ctx)
if err != nil {
    return err
}

// Use the component
// ...

// Dispose when done
err = comp.Dispose()
if err != nil {
    log.Printf("Error disposing component: %v", err)
}
```

### 7.3 Component Dependency Management

```go
// Create components
comp1 := component.NewDefaultComponent("comp-1", "Component 1", component.TypeBasic)
comp2 := component.NewDefaultComponent("comp-2", "Component 2", component.TypeBasic)

// Create dependency-aware component
depComp := component.NewDependencyAwareComponent("dep-comp", "Dependent Component", component.TypeBasic)

// Add dependencies
depComp.AddDependency(comp1)
depComp.AddDependency(comp2)

// Register with registry
registry.Register(comp1)
registry.Register(comp2)
registry.Register(depComp)

// Initialize registry (handles dependency order)
registry.Initialize(ctx)
```

### 7.4 Error Handling Patterns

```go
// Creating domain errors
err := component.NewError(
    component.ErrInitializationFailed,
    "Failed to initialize database connection",
    map[string]interface{}{"component": "database"},
    originalError,
)

// Checking error types
if component.IsErrorCode(err, component.ErrComponentNotFound) {
    // Handle component not found case
}

// Unwrapping errors
originalErr := component.UnwrapError(err)
```

### 7.5 Event-Based Communication

```go
// Create event bus
bus := event.NewDefaultEventBus()

// Subscribe to events
subscription := bus.Subscribe("component.initialized", func(e *event.Event) {
    componentID := e.Payload["componentId"].(string)
    fmt.Printf("Component initialized: %s\n", componentID)
})

// Publish events
bus.Publish("component.initialized", map[string]interface{}{
    "componentId": "comp-1",
    "timestamp": time.Now(),
})

// Cancel subscription when done
subscription.Cancel()
```

### 7.6 System Service Usage

#### 7.6.1 Using the Builder Pattern

```go
// Create a builder
builder := system.NewBuilder("system-1")

// Configure system components
builder.WithLogger(logging.CreateStandardLogger(logging.Info))
       .WithConfiguration(config.LoadConfiguration("config.yaml"))
       .WithRegistry(component.NewDefaultRegistry())
       .WithEventBus(event.NewDefaultEventBus())
       .WithMultiStore(storage.CreateMultiStore("data"))

// Build the system service
systemService, err := builder.Build()
if err != nil {
    log.Fatalf("Failed to build system: %v", err)
}

// Initialize and start the system
ctx := component.CreateContext()
err = systemService.Initialize(ctx)
err = systemService.Start(ctx)
```

#### 7.6.2 Using the Factory Pattern

```go
// Load system configuration
configLoader := system.NewConfigLoader(appConfig)
systemConfig, err := configLoader.LoadSystemConfig()
if err != nil {
    log.Fatalf("Failed to load system config: %v", err)
}

// Create factory with dependencies
factory := system.NewFactory(
    registry,
    eventBus, 
    configuration,
    multiStore,
)

// Create system service
systemService, err := factory.CreateSystemService(systemConfig)
if err != nil {
    log.Fatalf("Failed to create system service: %v", err)
}
```

#### 7.6.3 Using the System Service

```go
// Execute an operation
result, err := systemService.ExecuteOperation(ctx, "operation-id", inputData)
if err != nil {
    log.Printf("Operation failed: %v", err)
    return
}

// Process operation result
output := result.(*system.SystemOperationOutput)
processData(output.Data)

// Start a service
err = systemService.StartService(ctx, "service-id")
if err != nil {
    log.Printf("Failed to start service: %v", err)
    return
}

// Stop a service
err = systemService.StopService(ctx, "service-id")
```

#### 7.6.4 Shutdown Sequence

```go
// Graceful shutdown
err = systemService.Stop(ctx)
if err != nil {
    log.Printf("Error during system shutdown: %v", err)
}
```

## 8. Extension Points

### 8.1 Creating Custom Components

```go
type MyComponent struct {
    *component.BaseComponent
    // Custom fields
    config MyConfig
    client *http.Client
}

func NewMyComponent(id, name string, config MyConfig) *MyComponent {
    return &MyComponent{
        BaseComponent: component.NewBaseComponent(id, name, component.TypeBasic),
        config: config,
        client: &http.Client{},
    }
}

func (c *MyComponent) Initialize(ctx component.Context) error {
    // Call base implementation
    if err := c.BaseComponent.Initialize(ctx); err != nil {
        return err
    }
    
    // Custom initialization
    // ...
    
    return nil
}

func (c *MyComponent) Dispose() error {
    // Custom cleanup
    // ...
    
    // Call base implementation
    return c.BaseComponent.Dispose()
}
```

### 8.2 Creating Custom Operations

```go
type MyOperation struct {
    *operation.BaseOperation
    // Custom fields
}

func NewMyOperation(id, name string) *MyOperation {
    return &MyOperation{
        BaseOperation: operation.NewBaseOperation(id, name),
    }
}

func (o *MyOperation) Execute(ctx component.Context, input operation.Input) (operation.Output, error) {
    // Input validation
    typedInput, ok := input.(MyInput)
    if !ok {
        return nil, component.NewError(
            operation.ErrInvalidInput,
            "Expected MyInput type",
            nil,
            nil,
        )
    }
    
    // Operation implementation
    result := processInput(typedInput)
    
    return result, nil
}
```

### 8.3 Creating Custom Services

```go
type MyService struct {
    *service.BaseService
    // Custom fields
    server *http.Server
}

func NewMyService(id, name string) *MyService {
    s := &MyService{
        BaseService: service.NewBaseService(id, name),
        server: &http.Server{Addr: ":8080"},
    }
    
    // Register health check
    s.RegisterHealthCheck("server", s.checkServerHealth)
    
    return s
}

func (s *MyService) Start(ctx component.Context) error {
    if err := s.BaseService.Start(ctx); err != nil {
        return err
    }
    
    // Start HTTP server in a goroutine
    go func() {
        if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            // Log error
        }
    }()
    
    return nil
}

func (s *MyService) Stop(ctx component.Context) error {
    // Stop HTTP server
    stopCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := s.server.Shutdown(stopCtx); err != nil {
        return err
    }
    
    return s.BaseService.Stop(ctx)
}

func (s *MyService) checkServerHealth(ctx component.Context) error {
    // Implement health check
    return nil
}
```

### 8.4 Building and Deploying Plugins

```go
func CreatePlugin() plugin.Plugin {
    p := plugin.NewDefaultPlugin("my-plugin", "1.0.0", nil)
    
    // Create plugin components
    comp1 := component.NewDefaultComponent("plugin-comp-1", "Plugin Component 1", component.TypeBasic)
    comp2 := mypackage.NewMyComponent("plugin-comp-2", "Plugin Component 2")
    
    // Add components to plugin
    p.AddComponent(comp1)
    p.AddComponent(comp2)
    
    return p
}

// In main.go of plugin
func main() {
    // Export plugin creation function
    plugin.Export(CreatePlugin)
}
```

## 9. Best Practices

### 9.1 Component Design

1. **Keep components focused**: Each component should have a single responsibility
2. **Use composition**: Favor composition over complex inheritance hierarchies
3. **Design for testability**: Make components easy to test in isolation
4. **Minimize dependencies**: Keep the dependency graph manageable
5. **Document component behavior**: Make component contracts clear

### 9.2 Error Handling

1. **Use domain-specific errors**: Create errors with meaningful codes and messages
2. **Include context**: Add relevant details to error information
3. **Wrap original errors**: Preserve original error information when wrapping
4. **Proper cleanup on errors**: Ensure resources are released when errors occur
5. **Don't swallow errors**: Always handle or propagate errors appropriately

### 9.3 Configuration Management

1. **Define defaults**: Provide sensible defaults for component configuration
2. **Validate early**: Check configuration validity at component creation
3. **Use typed configuration**: Convert string configuration to appropriate types
4. **Support runtime reconfiguration**: Allow components to adapt to configuration changes
5. **Document configuration options**: Make configuration requirements clear

### 9.4 Event Usage

1. **Define clear event semantics**: Document what each event means and when it's emitted
2. **Keep events focused**: Each event should represent a single concept
3. **Handle subscription errors**: Subscribers should not crash the system
4. **Clean up subscriptions**: Cancel subscriptions when components are disposed
5. **Consider synchronous vs. asynchronous**: Choose the right event model for each use case

### 9.5 Testing

1. **Unit test all components**: Ensure each component functions correctly in isolation
2. **Test lifecycle methods**: Verify initialize and dispose work correctly
3. **Test error cases**: Ensure components handle errors appropriately
4. **Test concurrency**: Verify thread safety claims
5. **Integration test component interactions**: Test how components work together

## 10. Migration Guide

### 10.1 Migrating from Legacy Components

1. **Create adapters**: Build adapters to bridge between old and new components
2. **Gradual migration**: Replace components one at a time
3. **Use the registry**: The registry can help manage hybrid systems
4. **Test thoroughly**: Verify compatibility at each migration step
5. **Update documentation**: Keep documentation in sync with migrated components

### 10.2 Compatibility Considerations

1. **API compatibility**: Ensure new components provide equivalent functionality
2. **Error handling**: Map between old and new error models
3. **Lifecycle management**: Align component lifecycles
4. **Configuration migration**: Convert configuration formats
5. **Performance characteristics**: Verify similar performance characteristics

### 10.3 Gradual Migration Strategy

1. **Start with independent components**: Begin with components that have few dependencies
2. **Create a parallel registry**: Maintain both old and new registries during migration
3. **Bridge events**: Connect old event system to new event bus
4. **Migrate core components last**: Leave the most fundamental components until the end
5. **Validate at each step**: Maintain testing coverage throughout migration

## 11. Appendices

### 11.1 Glossary

| Term | Definition |
|------|------------|
| Component | Fundamental building block of the system |
| Operation | Component that performs discrete units of work |
| Service | Component that provides ongoing functionality |
| Plugin | Container for components that extends the system |
| Registry | System that tracks registered components |
| Factory | System that creates components from configuration |
| Lifecycle | The sequence of states a component goes through |
| Metadata | Additional information associated with a component |

### 11.2 Configuration Reference

| Configuration | Type | Default | Description |
|---------------|------|---------|-------------|
| component.id | string | required | Unique identifier for the component |
| component.name | string | required | Human-readable name for the component |
| component.type | string | required | Type of the component (basic, operation, service) |
| component.dependencies | []string | [] | IDs of components this component depends on |
| service.healthCheckInterval | duration | 30s | Interval between automatic health checks |

### 11.3 Common Error Codes

| Error Code | Description |
|------------|-------------|
| component.not_found | Component not found in registry |
| component.already_exists | Component with same ID already registered |
| component.initialization_failed | Component initialization failed |
| component.dispose_failed | Component disposal failed |
| operation.execution_failed | Operation execution failed |
| operation.invalid_input | Invalid input for operation |
| service.start_failed | Service failed to start |
| service.stop_failed | Service failed to stop |
| plugin.not_found | Plugin not found |
| plugin.load_failed | Plugin failed to load |

### 11.4 Testing Approach

The component system has comprehensive testing:

1. **Unit Tests**: All components have unit tests with >90% coverage
2. **Mock Components**: Test doubles are available for testing
3. **Test Registry**: A test registry simplifies component testing
4. **Event Recording**: Event recorders capture events for verification
5. **Test Context**: Helper functions create contexts for testing

### 11.5 Performance Considerations

Key performance aspects to consider:

1. **Component Initialization**: Component initialization is done in dependency order and can impact startup time
2. **Event Handling**: Synchronous event handling blocks until all handlers complete
3. **Registry Lookups**: Component lookups are optimized but can be a bottleneck
4. **Service Health Checks**: Health check frequency affects CPU and resource usage
5. **Plugin Loading**: Dynamic plugin loading has overhead and security implications 