# System and RuntimeEnvironment API Reference

Complete API reference for the core system interfaces used in plugin and component development.

## 🎯 Interface Overview

### System vs RuntimeEnvironment

**`System`** - Basic system operations:
- Component registry access
- Operation execution
- Service lifecycle management
- System-wide start/stop control

**`RuntimeEnvironment`** - Enhanced system access:
- All `System` methods
- Direct access to framework services (Logger, Configuration, EventBus, PluginManager)
- Batch plugin loading

### ⚠️ Important: Context Types

**All framework interfaces use Skeleton's context, NOT Go's standard context:**

```go
// ✅ CORRECT - Framework context
import "github.com/fintechain/skeleton/pkg/context"

func (c *MyComponent) Initialize(ctx context.Context, system component.System) error {
    // ctx is skeleton context.Context, not Go's context.Context
}

// ❌ WRONG - Go's standard context  
import "context"

func (c *MyComponent) Initialize(ctx context.Context, system component.System) error {
    // This will cause type mismatch errors!
}
```

**Key Differences:**
- **Skeleton Context**: `WithValue()` returns `Context` (framework type)
- **Go Context**: `WithValue()` returns `context.Context` (Go standard type)
- **Skeleton Context**: Custom error constants (`ErrContextCanceled`, `ErrContextDeadlineExceeded`)
- **Go Context**: Standard error values (`context.Canceled`, `context.DeadlineExceeded`)

**Always import**: `"github.com/fintechain/skeleton/pkg/context"`

## 🔄 Context Interface

```go
type Context interface {
    Value(key interface{}) interface{}
    WithValue(key, value interface{}) Context
    Deadline() (time.Time, bool)
    Done() <-chan struct{}
    Err() error
}
```

**Import Path**: `github.com/fintechain/skeleton/pkg/context`

```go
import "github.com/fintechain/skeleton/pkg/context"

var ctx context.Context = context.NewContext()
```

### Key Methods

#### `Value(key interface{}) interface{}`
Retrieves a value from the context by key. Returns nil if key not found.

```go
value := ctx.Value("my-key")
if value != nil {
    myData := value.(string)
}
```

#### `WithValue(key, value interface{}) Context`
Creates a new context with an additional key-value pair. **Returns framework Context type**.

```go
newCtx := ctx.WithValue("component-id", "my-component")
newCtx = newCtx.WithValue("request-id", "req-123")
```

#### `Deadline() (time.Time, bool)`
Returns the deadline for this context, if any. Second return value indicates if deadline is set.

```go
deadline, hasDeadline := ctx.Deadline()
if hasDeadline {
    timeLeft := time.Until(deadline)
}
```

#### `Done() <-chan struct{}`
Returns a channel that's closed when context is cancelled or times out. Use in select statements.

```go
select {
case <-ctx.Done():
    // Context was cancelled or timed out
    return ctx.Err()
case result := <-workChannel:
    // Work completed before cancellation
    return nil
}
```

#### `Err() error`
Returns the error that caused context cancellation. Nil if not cancelled.

```go
if err := ctx.Err(); err != nil {
    if strings.Contains(err.Error(), "context.context_canceled") {
        // Handle cancellation
    } else if strings.Contains(err.Error(), "context.context_deadline_exceeded") {
        // Handle timeout
    }
}
```

### Context Creation Functions

#### `NewContext() Context`
Creates a new framework context instance.

```go
ctx := context.NewContext()
```

#### `NewContextWithTimeout(timeout time.Duration) Context`
Creates a context that automatically cancels after the timeout duration.

```go
ctx := context.NewContextWithTimeout(30 * time.Second)
defer func() {
    if ctx.Err() != nil {
        log.Println("Operation timed out")
    }
}()
```

#### `NewContextWithDeadline(deadline time.Time) Context`
Creates a context that automatically cancels at the specified deadline.

```go
deadline := time.Now().Add(1 * time.Hour)
ctx := context.NewContextWithDeadline(deadline)
```

## 📚 System Interface

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

**Import Path**: `github.com/fintechain/skeleton/pkg/component`

```go
import "github.com/fintechain/skeleton/pkg/component"

var system component.System
```

### Methods

#### `Registry() Registry`
Returns the component registry for discovering and accessing registered components.

```go
registry := system.Registry()
component, err := registry.Get("my-component-id")
```

#### `ExecuteOperation(ctx, operationID, input) (Output, error)`
Executes a registered operation component with the given input.

**Parameters:**
- `ctx context.Context` - Execution context
- `operationID ComponentID` - ID of the operation to execute
- `input Input` - Input data for the operation

**Returns:**
- `Output` - Operation results
- `error` - Execution error if any

```go
result, err := system.ExecuteOperation(ctx, "calculate-sum", component.Input{
    Data: map[string]interface{}{
        "numbers": []float64{1, 2, 3, 4, 5},
    },
})
```

#### `StartService(ctx, serviceID) error`
Starts a registered service component.

```go
err := system.StartService(ctx, "web-server")
```

#### `StopService(ctx, serviceID) error`
Stops a running service component gracefully.

```go
err := system.StopService(ctx, "web-server")
```

#### `Start(ctx) error`
Initializes and starts the entire system. Idempotent - safe to call multiple times.

#### `Stop(ctx) error`
Gracefully shuts down the entire system. Idempotent - safe to call multiple times.

#### `IsRunning() bool`
Returns whether the system is currently running.

## 🚀 RuntimeEnvironment Interface

```go
type RuntimeEnvironment interface {
    component.System  // All System methods above
    
    PluginManager() plugin.PluginManager
    EventBus() event.EventBusService
    Logger() logging.Logger
    Configuration() config.Configuration
    LoadPlugins(ctx context.Context, plugins []plugin.Plugin) error
}
```

**Import Path**: `github.com/fintechain/skeleton/pkg/runtime`

```go
import "github.com/fintechain/skeleton/pkg/runtime"

var env runtime.RuntimeEnvironment
```

### Additional Methods

#### `Logger() logging.Logger`
Returns the system's logger for structured logging.

```go
logger := runtime.Logger()
logger.Info("Component started", "component_id", "my-component")
```

#### `Configuration() config.Configuration`
Returns the system's configuration service.

```go
config := runtime.Configuration()
port := config.GetIntDefault("server.port", 8080)
```

#### `EventBus() event.EventBusService`
Returns the system's event bus for publish-subscribe messaging.

```go
eventBus := runtime.EventBus()
eventBus.PublishAsync(&event.Event{
    Topic: "component.started",
    Source: "my-component",
    Payload: map[string]interface{}{"status": "ok"},
})
```

#### `PluginManager() plugin.PluginManager`
Returns the system's plugin manager for plugin lifecycle operations.

```go
pluginManager := runtime.PluginManager()
err := pluginManager.StartPlugin(ctx, "my-plugin")
```

#### `LoadPlugins(ctx, plugins) error`
Loads multiple plugins into the system for batch plugin loading.

```go
err := runtime.LoadPlugins(ctx, []plugin.Plugin{plugin1, plugin2})
```

## 🗂️ Registry Interface

```go
type Registry interface {
    Register(component Component) error
    Get(id ComponentID) (Component, error)
    GetByType(typ ComponentType) ([]Component, error)
    Find(predicate func(Component) bool) ([]Component, error)
    Has(id ComponentID) bool
    List() []ComponentID
    Unregister(id ComponentID) error
    Count() int
    Clear() error
}
```

**Import Path**: `github.com/fintechain/skeleton/pkg/component`

```go
import "github.com/fintechain/skeleton/pkg/component"

var registry component.Registry
```

### Key Methods

#### `Get(id) (Component, error)`
Retrieves a component by ID.

```go
comp, err := registry.Get("database-service")
if err == nil {
    dbService := comp.(*DatabaseService)
}
```

#### `GetByType(typ) ([]Component, error)`
Retrieves all components of a specific type.

```go
services, err := registry.GetByType(component.TypeService)
```

#### `Has(id) bool`
Checks if a component exists.

```go
if registry.Has("optional-service") {
    // Use the service
}
```

## ⚙️ Configuration Interface

```go
type Configuration interface {
    GetString(key string) string
    GetStringDefault(key, defaultValue string) string
    GetInt(key string) (int, error)
    GetIntDefault(key string, defaultValue int) int
    GetBool(key string) (bool, error)
    GetBoolDefault(key string, defaultValue bool) bool
    GetDuration(key string) (time.Duration, error)
    GetDurationDefault(key string, defaultValue time.Duration) time.Duration
    GetObject(key string, result interface{}) error
    Exists(key string) bool
}
```

**Import Path**: `github.com/fintechain/skeleton/pkg/config`

```go
import "github.com/fintechain/skeleton/pkg/config"

var config config.Configuration
```

### Usage Examples

```go
config := runtime.Configuration()

// String values
host := config.GetStringDefault("database.host", "localhost")

// Numeric values
port := config.GetIntDefault("database.port", 5432)
timeout := config.GetDurationDefault("request.timeout", 30*time.Second)

// Boolean values
enabled := config.GetBoolDefault("feature.enabled", false)

// Check existence
if config.Exists("optional.setting") {
    value := config.GetString("optional.setting")
}
```

## 📝 Logger Interface

```go
type Logger interface {
    Debug(msg string, args ...interface{})
    Info(msg string, args ...interface{})
    Warn(msg string, args ...interface{})
    Error(msg string, args ...interface{})
}
```

**Import Path**: `github.com/fintechain/skeleton/pkg/logging`

```go
import "github.com/fintechain/skeleton/pkg/logging"

// Can work with Logger types directly
var logger logging.Logger = runtime.Logger()
```

### Structured Logging

```go
logger := runtime.Logger()

// Key-value pairs
logger.Info("User login", 
    "user_id", "123",
    "ip_address", "192.168.1.1",
    "timestamp", time.Now())

// Maps
logger.Error("Database error", map[string]interface{}{
    "error_code": 500,
    "query": "SELECT * FROM users",
    "duration": "2.5s",
})
```

## 🔌 PluginManager Interface

```go
type PluginManager interface {
    component.Service  // All Service methods (Start, Stop, etc.)
    
    Add(pluginID ComponentID, plugin Plugin) error
    Remove(pluginID ComponentID) error
    StartPlugin(ctx context.Context, pluginID ComponentID) error
    StopPlugin(ctx context.Context, pluginID ComponentID) error
    GetPlugin(pluginID ComponentID) (Plugin, error)
    ListPlugins() []ComponentID
}
```

**Import Path**: `github.com/fintechain/skeleton/pkg/plugin`

```go
import "github.com/fintechain/skeleton/pkg/plugin"

var pluginManager plugin.PluginManager
```

### Key Methods

#### `Add(pluginID, plugin) error`
Adds a plugin to the manager's registry.

```go
pluginManager := runtime.PluginManager()
err := pluginManager.Add("my-plugin", myPlugin)
```

#### `StartPlugin(ctx, pluginID) error`
Starts a specific plugin by ID.

```go
err := pluginManager.StartPlugin(ctx, "web-plugin")
```

#### `StopPlugin(ctx, pluginID) error`
Stops a specific plugin gracefully.

```go
err := pluginManager.StopPlugin(ctx, "web-plugin")
```

#### `GetPlugin(pluginID) (Plugin, error)`
Retrieves a plugin by ID.

```go
plugin, err := pluginManager.GetPlugin("my-plugin")
if err == nil {
    // Use the plugin
}
```

#### `ListPlugins() []ComponentID`
Returns all registered plugin IDs.

```go
pluginIDs := pluginManager.ListPlugins()
for _, id := range pluginIDs {
    logger.Info("Found plugin", "plugin_id", id)
}
```

### Usage Examples

```go
pluginManager := runtime.PluginManager()

// Register and start a plugin
err := pluginManager.Add("auth-plugin", authPlugin)
if err != nil {
    return err
}

err = pluginManager.StartPlugin(ctx, "auth-plugin")
if err != nil {
    return err
}

// List all plugins
plugins := pluginManager.ListPlugins()
logger.Info("Active plugins", "count", len(plugins))

// Stop and remove a plugin
err = pluginManager.StopPlugin(ctx, "auth-plugin")
err = pluginManager.Remove("auth-plugin")
```

## 📡 EventBus Interface

```go
type EventBus interface {
    Publish(event *Event) error
    PublishAsync(event *Event) error
    Subscribe(eventType string, handler EventHandler) Subscription
    SubscribeAsync(eventType string, handler EventHandler) Subscription
    WaitAsync()
}
```

**Import Path**: `github.com/fintechain/skeleton/pkg/event`

```go
import "github.com/fintechain/skeleton/pkg/event"

var eventBus event.EventBus
```

### Event Structure

```go
type Event struct {
    Topic   string                 // Event type/category
    Source  string                 // Component that generated the event
    Time    time.Time             // When the event occurred
    Payload map[string]interface{} // Event data
}
```

### Usage Examples

```go
eventBus := runtime.EventBus()

// Publishing events
event := &event.Event{
    Topic: "user.login",
    Source: "auth-service",
    Time: time.Now(),
    Payload: map[string]interface{}{
        "user_id": "123",
        "success": true,
    },
}
eventBus.PublishAsync(event)

// Subscribing to events
subscription := eventBus.SubscribeAsync("user.login", func(e *event.Event) {
    logger.Info("User login event received", "user_id", e.Payload["user_id"])
})

// Cancel subscription when done
defer subscription.Cancel()
```

## 🎯 Common Usage Patterns

### Component Initialization

```go
func (c *MyComponent) Initialize(ctx context.Context, system component.System) error {
    // Store runtime reference
    c.runtime = system.(runtime.RuntimeEnvironment)
    
    // Access framework services
    logger := c.runtime.Logger()
    config := c.runtime.Configuration()
    
    // Access other components
    registry := system.Registry()
    if registry.Has("database-service") {
        dbComp, _ := registry.Get("database-service")
        c.database = dbComp.(*DatabaseService)
    }
    
    logger.Info("Component initialized", "component_id", c.ID())
    return nil
}
```

### Service Discovery

```go
// Find all services
registry := system.Registry()
services, err := registry.GetByType(component.TypeService)

// Find specific component
if registry.Has("cache-service") {
    cache, _ := registry.Get("cache-service")
}

// Find components by criteria
webServices, err := registry.Find(func(comp component.Component) bool {
    return strings.Contains(string(comp.ID()), "web")
})
```

### Configuration-Driven Behavior

```go
config := runtime.Configuration()

// Feature flags
if config.GetBoolDefault("features.caching", false) {
    // Enable caching
}

// Timeouts and limits
maxRetries := config.GetIntDefault("http.max_retries", 3)
timeout := config.GetDurationDefault("http.timeout", 30*time.Second)

// Environment-specific settings
env := config.GetStringDefault("environment", "development")
if env == "production" {
    // Production-specific configuration
}
```

---

**Note**: This reference covers the core interfaces used in plugin and component development. For complete implementation examples, see the [Plugin Development Guide](PLUGIN_DEVELOPMENT_GUIDE.md) and [Service Operations Development Guide](SERVICE_OPERATIONS_DEVELOPMENT_GUIDE.md).