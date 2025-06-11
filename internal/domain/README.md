# Domain Layer - Architecture & Interface Guide

The domain layer is the architectural foundation of the Fintechain Skeleton framework, implementing **Domain-Driven Design (DDD)** principles to provide clean, testable, and maintainable abstractions.

## Why This Layer Exists

The domain layer serves as the **contract layer** between your business logic and infrastructure concerns:

- **Infrastructure Independence**: Domain interfaces can be implemented by any concrete technology (memory, database, cloud services)
- **Testability**: Easy to mock and test business logic without external dependencies
- **Flexibility**: Swap implementations without changing business code
- **Clean Architecture**: Enforces dependency inversion - infrastructure depends on domain, never the reverse

## Package Architecture

```
internal/domain/
├── component/     # 🏗️  Core component system with lifecycle management
├── storage/       # 💾  Storage abstraction with multi-backend support
├── event/         # 📡  Event-driven communication system
├── config/        # ⚙️   Configuration management interfaces
├── context/       # 🔄  Application context handling
├── logging/       # 📝  Structured logging abstraction
└── runtime/       # 🚀  High-level runtime environment orchestration
```

---

## 🏗️ Component System (`component/`)

The component system is the **heart** of the framework, providing a unified way to manage application building blocks.

### Core Philosophy

**Everything is a Component** - Whether it's a database connection, a web server, or a calculation engine, it's all managed through the same lifecycle and discovery mechanisms.

### Component Types

```go
const (
    TypeComponent   = "component"  // Basic managed entity
    TypeOperation   = "operation"  // Executable instruction  
    TypeService     = "service"    // Long-running process
    
    // Legacy types (deprecated but supported)
    TypeBasic       = "basic"        // Legacy - use TypeComponent
    TypeSystem      = "system"       // Legacy - use TypeService  
    TypeApplication = "application"  // Legacy - use TypeComponent
)
```

**Why These Types?**
- **Component**: Simple entities that need lifecycle management
- **Operation**: Discrete work units with input/output (think functions with context)
- **Service**: Background processes that run continuously (web servers, message processors)
- **Legacy Types**: Maintained for backward compatibility but prefer the main three types

### Key Interfaces

#### `Component` - The Foundation
```go
type Component interface {
    // Identity
    ID() ComponentID
    Name() string
    Type() ComponentType
    
    // Lifecycle
    Initialize(ctx context.Context, system System) error
    Dispose() error
}
```

**Why This Design?**
- **Identity**: Every component is discoverable and unique
- **Lifecycle**: Predictable initialization and cleanup patterns
- **System Access**: Components get access to the broader system during initialization

#### `System` - The Orchestrator
```go
type System interface {
    Registry() Registry
    ExecuteOperation(ctx context.Context, operationID ComponentID, input Input) (Output, error)
    StartService(ctx context.Context, serviceID ComponentID) error
    StopService(ctx context.Context, serviceID ComponentID) error
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    IsRunning() bool
}
```

**Why This Design?**
- **Single Entry Point**: One place to access all system capabilities
- **Type-Safe Operations**: Different methods for different component types
- **Lifecycle Control**: Centralized start/stop coordination
- **Registry Access**: All components are discoverable through the registry

#### `Registry` - The Directory
```go
type Registry interface {
    Register(component Component) error
    Get(id ComponentID) (Component, error)
    GetByType(typ ComponentType) ([]Component, error)
    Find(predicate func(Component) bool) ([]Component, error)
}
```

**Why This Design?**
- **Dependency Injection**: Components find their dependencies through the registry
- **Type Discovery**: Find all components of a specific type
- **Flexible Queries**: Custom search logic through predicates

### Implementation Guide

#### Creating a Basic Component
```go
type DatabaseConnection struct {
    *infraComponent.BaseComponent
    connectionString string
    db              *sql.DB
}

func (d *DatabaseConnection) Initialize(ctx context.Context, system component.System) error {
    if err := d.BaseComponent.Initialize(ctx, system); err != nil {
        return err
    }
    
    // Access system services through registry
    registry := d.system().Registry()
    loggerComp, err := registry.Get("logger")
    if err == nil {
        if logger, ok := loggerComp.(logging.LoggerService); ok {
            logger.Info("Connecting to database")
        }
    }
    
    // Initialize your component
    db, err := sql.Open("postgres", d.connectionString)
    d.db = db
    return err
}
```

#### Creating an Operation
```go
type CalculateTotal struct {
    *infraComponent.BaseOperation
}

func (c *CalculateTotal) Execute(ctx context.Context, input component.Input) (component.Output, error) {
    items := input.Data.([]Item)
    total := 0.0
    for _, item := range items {
        total += item.Price
    }
    return component.Output{Data: total}, nil
}
```

#### Creating a Service
```go
type WebServer struct {
    *infraComponent.BaseService
    server *http.Server
}

func (w *WebServer) Start(ctx context.Context) error {
    if err := w.BaseService.Start(ctx); err != nil {
        return err
    }
    
    go w.server.ListenAndServe()
    return nil
}
```

---

## 💾 Storage System (`storage/`)

Provides **backend-agnostic** storage with support for multiple engines and advanced features.

### Design Philosophy

**Storage as a Service** - Storage is treated as a pluggable service that can be swapped without changing business logic.

### Key Interfaces

#### `Store` - The Foundation
```go
type Store interface {
    Get(key []byte) ([]byte, error)
    Set(key, value []byte) error
    Delete(key []byte) error
    Has(key []byte) (bool, error)
    Close() error
}
```

#### `Engine` - The Factory
```go
type Engine interface {
    Name() string
    Create(name, path string, config Config) (Store, error)
    Capabilities() Capabilities
}
```

**Why Engines?**
- **Multiple Backends**: Memory, file-based, database, cloud storage
- **Capability Reporting**: Know what features each engine supports
- **Factory Pattern**: Consistent store creation across different technologies

### Advanced Features

#### Transactions (Optional)
```go
type Transactional interface {
    BeginTransaction() (Transaction, error)
}

type Transaction interface {
    Set(key, value []byte) error
    Delete(key []byte) error
    Commit() error
    Rollback() error
}
```

#### Range Queries (Optional)
```go
type RangeQueryable interface {
    Range(start, end []byte, fn func(key, value []byte) bool) error
}
```

### Usage Example
```go
// Access storage through registry (storage is registered as a component)
registry := system.Registry()
storageComp, err := registry.Get("storage")
if err != nil {
    return err
}

multiStore := storageComp.(storage.MultiStoreService)
store, err := multiStore.GetStore("user-data")
if err != nil {
    return err
}

// Use storage
user := User{ID: "123", Name: "Alice"}
data, _ := json.Marshal(user)
store.Set([]byte("user:123"), data)
```

---

## 📡 Event System (`event/`)

Enables **loose coupling** between components through publish-subscribe messaging.

### Design Philosophy

**Event-Driven Architecture** - Components communicate through events rather than direct calls, enabling better modularity and testability.

### Key Interfaces

#### `EventBus` - The Messenger
```go
type EventBus interface {
    Publish(event *Event) error
    PublishAsync(event *Event) error
    Subscribe(eventType string, handler EventHandler) Subscription
    SubscribeAsync(eventType string, handler EventHandler) Subscription
    WaitAsync()
}
```

#### `Event` - The Message
```go
type Event struct {
    Topic   string
    Source  string
    Time    time.Time
    Payload map[string]interface{}
}
```

### Usage Example
```go
// Access event bus through RuntimeEnvironment
if runtime, ok := system.(runtime.RuntimeEnvironment); ok {
    eventBus := runtime.EventBus()
    // Subscribe to events
    subscription := eventBus.Subscribe("user.created", func(event *Event) {
        userID := event.Payload["user_id"].(string)
        // Note: You'd need to access logger the same way
    })

    // Publish events  
    event := &Event{
        Topic:   "user.created",
        Source:  "user-service",
        Time:    time.Now(),
        Payload: map[string]interface{}{"user_id": "123"},
    }
    eventBus.Publish(event)
}
```

---

## ⚙️ Configuration System (`config/`)

Provides **hierarchical configuration** with multiple sources and type-safe access.

### Design Philosophy

**Configuration Composition** - Multiple configuration sources (files, environment, memory) are composed with precedence rules.

### Key Interfaces

#### `Configuration` - The Accessor
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

#### `ConfigurationSource` - The Provider
```go
type ConfigurationSource interface {
    LoadConfig() error
    GetValue(key string) (interface{}, bool)
}
```

### Usage Example
```go
// Access configuration through RuntimeEnvironment
if runtime, ok := system.(runtime.RuntimeEnvironment); ok {
    config := runtime.Configuration()
    
    // Type-safe access with defaults
    dbHost := config.GetStringDefault("database.host", "localhost")
    dbPort := config.GetIntDefault("database.port", 5432)
    debugMode := config.GetBoolDefault("debug", false)
}
```

---

## 📝 Logging System (`logging/`)

Provides **structured logging** with multiple output targets and pluggable implementations.

### Design Philosophy

**Structured Logging** - All log entries include structured data for better searchability and analysis.

### Key Interface

#### `Logger` - The Writer
```go
type Logger interface {
    Debug(msg string, args ...interface{})
    Info(msg string, args ...interface{})
    Warn(msg string, args ...interface{})
    Error(msg string, args ...interface{})
}
```

### Usage Example
```go
// Access logger through RuntimeEnvironment
if runtime, ok := system.(runtime.RuntimeEnvironment); ok {
    logger := runtime.Logger()
    // Structured logging with key-value pairs
    logger.Info("User login", 
        "user_id", "123",
        "ip_address", "192.168.1.1",
        "timestamp", time.Now())
}
```

---

## 🔄 Context System (`context/`)

Provides **request-scoped data** and cancellation support for component operations.

### Design Philosophy

**Domain Context** - Extends Go's context with domain-specific features while maintaining compatibility.

### Key Interface

#### `Context` - The Carrier
```go
type Context interface {
    Value(key interface{}) interface{}
    WithValue(key, value interface{}) Context
    Done() <-chan struct{}
    Err() error
}
```

### Usage Example
```go
// Create context with request data
ctx := context.NewContext()
ctx = ctx.WithValue("request_id", "abc-123")
ctx = ctx.WithValue("user_id", "456")

// Pass context through operations
input := component.Input{
    Data: calculationData,
    Metadata: map[string]string{"request_id": "abc-123"},
}
result, err := system.ExecuteOperation(ctx, "calculate-total", input)
```

---

## 🚀 Runtime System (`runtime/`)

Provides **high-level orchestration** that combines all domain services into a unified runtime environment.

### Design Philosophy

**Unified Runtime** - Single entry point that provides access to all system capabilities with proper lifecycle management.

### System vs RuntimeEnvironment

The framework provides two levels of system access:

**`System` (Basic)**:
- Core lifecycle operations (start, stop, execute operations)
- Registry access for component discovery
- Basic system orchestration

**`RuntimeEnvironment` (Enhanced)**:
- Extends `System` with convenience accessors
- Direct access to core services: `Logger()`, `EventBus()`, `PluginManager()`
- Batch operations like `LoadPlugins()`
- Typically used at the application level

### Key Interface

#### `RuntimeEnvironment` - The Orchestrator
```go
type RuntimeEnvironment interface {
    component.System
    
    // Convenience accessors (direct access with dependency injection)
    PluginManager() plugin.PluginManager
    EventBus() event.EventBusService
    Logger() logging.Logger
    Configuration() config.Configuration
    
    // Batch operations
    LoadPlugins(ctx context.Context, plugins []plugin.Plugin) error
}
```

---

## Implementation Best Practices

### 1. **Embed Base Components**
Always embed `BaseComponent`, `BaseService`, or `BaseOperation` to get common functionality:
```go
type MyService struct {
    *infraComponent.BaseService  // Gets lifecycle, system access, etc.
    // Your specific fields
}
```

### 2. **Use System Access Carefully**
Access the system through the protected `system()` method only after initialization:
```go
func (s *MyService) DoWork() error {
    if !s.IsInitialized() {
        return errors.New("service not initialized")
    }
    
    // Access services through registry or RuntimeEnvironment
    if runtime, ok := s.system().(runtime.RuntimeEnvironment); ok {
        logger := runtime.Logger()
        logger.Info("Doing work")
    }
    return nil
}
```

### 3. **Handle Errors Gracefully**
Use domain error constants for consistent error handling:
```go
if err != nil {
    return fmt.Errorf("%s: %w", component.ErrServiceStartFailed, err)
}
```

### 4. **Implement Proper Cleanup**
Always clean up resources in `Dispose()`:
```go
func (s *MyService) Dispose() error {
    if s.connection != nil {
        s.connection.Close()
    }
    return s.BaseService.Dispose()
}
```

---

## Testing Your Implementations

The framework provides mock implementations for all domain interfaces:

```go
func TestMyComponent(t *testing.T) {
    factory := mocks.NewFactory()
    mockSystem := factory.SystemInterface()
    
    component := &MyComponent{}
    err := component.Initialize(context.Background(), mockSystem)
    
    assert.NoError(t, err)
    assert.True(t, component.IsInitialized())
}
```

---

## Next Steps

1. **Explore Infrastructure**: Look at `internal/infrastructure/` to see concrete implementations
2. **Check Examples**: Review `examples/` for usage patterns
3. **Build Components**: Start implementing your own components using these interfaces
4. **Run Tests**: Use the test suite to validate your implementations

The domain layer provides the foundation - everything else builds upon these contracts to create a flexible, testable, and maintainable system. 