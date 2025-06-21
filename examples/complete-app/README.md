# Complete Application Example (Builder API)

This example demonstrates a complete application using the **new Builder API** for the Fintechain Skeleton framework, showcasing all major framework patterns with **simple, explicit dependency injection** instead of complex FX patterns.

## 🎯 Builder API Benefits

### ✅ **Replaced Complex FX with Simple Builder**
- **No FX knowledge required** - simple builder pattern
- **Easy debugging** - no magic dependency injection
- **Clear dependency flow** - explicit dependency setting
- **Immediate feedback** - compile-time errors for missing dependencies
- **Simple custom dependencies** - direct creation, no complex providers

### ✅ **Before vs After**

**Before (FX-based):**
```go
// Complex FX provider pattern
runtime.ExecuteCommand("operation", input,
    runtime.WithPlugins(plugin1, plugin2),
    runtime.WithOptions(
        fx.Replace(fx.Annotate(providers.NewCustomLogger, fx.As(new(logging.LoggerService)))),
        fx.Replace(fx.Annotate(providers.NewCustomConfiguration, fx.As(new(config.Configuration)))),
        fx.Replace(fx.Annotate(providers.NewCustomEventBus, fx.As(new(event.EventBusService)))),
    ),
)
```

**After (Builder API):**
```go
// Simple builder pattern
runtime.NewBuilder().
    WithConfig(createCustomConfiguration()).
    WithLogger(createCustomLogger()).
    WithEventBus(createCustomEventBus()).
    WithPlugins(plugin1, plugin2).
    BuildCommand("operation", input)
```

## Project Structure

```
examples/complete-app/
├── main.go                 # Entry point and mode selection (Builder API)
├── modes/                  # Different execution modes (Builder API)
│   └── modes.go           # Daemon, command, and custom dependency modes
├── plugins/               # Example plugins (unchanged - work with both APIs)
│   ├── database/         # Database plugin with connection and query components
│   │   ├── plugin.go     # Main plugin orchestrator
│   │   ├── connection_service.go  # Database connection service
│   │   ├── query_operation.go     # Query processing operation
│   │   └── README.md     # Database plugin documentation
│   └── webserver/        # Web server plugin with HTTP components
│       ├── plugin.go     # Main plugin orchestrator
│       ├── http_service.go       # HTTP service implementation
│       ├── route_operation.go    # Route processing operation
│       └── README.md     # Web server plugin documentation
└── README.md             # This file
```

**Note**: The `providers/` directory was **removed** because the Builder API handles custom dependencies directly without complex FX provider patterns.

## What This Example Demonstrates

### 🔌 Plugin Orchestration (Same as Before)
- **Multiple plugins working together**: Webserver + Database plugins
- **Plugin lifecycle management**: Automatic initialization, startup, and shutdown
- **Plugin communication**: Through the shared runtime environment
- **Self-contained plugins**: Each plugin includes all its components

### 🛠️ Custom Dependencies (Simplified!)
- **Custom Logger**: Direct creation, no FX providers needed
- **Custom Configuration**: Simple function call, no complex annotations
- **Custom Event Bus**: Straightforward instantiation, no FX magic

### 🔄 Component Lifecycle (Same as Before)
- **Service lifecycle**: Initialize → Start → Stop → Dispose
- **Status tracking**: Monitor service states (Stopped, Running, etc.)
- **Graceful shutdown**: Proper cleanup of all resources

### 🚀 Execution Modes (Simplified API)
- **Daemon Mode**: Long-running services using `BuildDaemon()`
- **Command Mode**: Execute operations using `BuildCommand()`
- **Custom Dependencies**: Simple direct injection using `WithXxx()` methods

## Usage

### Daemon Mode (Long-running Services)
```bash
go run examples/complete-app/main.go daemon
```

**Builder API Code:**
```go
return runtime.NewBuilder().
    WithPlugins(
        webserver.NewWebServerPlugin(8080),
        database.NewDatabasePlugin("postgres", "test://connection"),
    ).
    BuildDaemon()
```

### Command Mode (Execute and Exit)
```bash
go run examples/complete-app/main.go command
```

**Builder API Code:**
```go
result, err := runtime.NewBuilder().
    WithPlugins(
        webserver.NewWebServerPlugin(8080),
        database.NewDatabasePlugin("postgres", "test://connection"),
    ).
    BuildCommand("database-query", map[string]interface{}{
        "query": "SELECT * FROM users WHERE active = true",
    })
```

### Custom Dependencies Mode
```bash
go run examples/complete-app/main.go custom
```

**Builder API Code:**
```go
// Create custom dependencies directly (no FX complexity!)
customConfig := createCustomConfiguration()
customLogger := createCustomLogger()
customEventBus := createCustomEventBus()

// Use with simple builder pattern
result, err := runtime.NewBuilder().
    WithConfig(customConfig).
    WithLogger(customLogger).
    WithEventBus(customEventBus).
    WithPlugins(
        webserver.NewWebServerPlugin(8080),
        database.NewDatabasePlugin("postgres", "custom://connection"),
    ).
    BuildCommand("database-query", input)
```

## Key Framework Patterns (Builder API)

### 1. Simple Builder Pattern
```go
// No FX knowledge required - simple builder pattern
err := runtime.NewBuilder().
    WithPlugins(plugin1, plugin2).
    WithConfig(myConfig).
    WithLogger(myLogger).
    BuildDaemon()
```

### 2. Direct Custom Dependencies
```go
// Create dependencies directly - no complex FX providers
func createCustomConfiguration() config.Configuration {
    return infraConfig.NewMemoryConfigurationWithData(map[string]interface{}{
        "app.name": "My App",
        "app.version": "1.0.0",
    })
}

// Use directly with builder
runtime.NewBuilder().
    WithConfig(createCustomConfiguration()).
    BuildDaemon()
```

### 3. Backward Compatibility
```go
// Old FX-based code still works
runtime.StartDaemon(
    runtime.WithPlugins(plugin1, plugin2),
)

// New Builder API is recommended
runtime.NewBuilder().
    WithPlugins(plugin1, plugin2).
    BuildDaemon()
```

## Architecture Benefits

### 🏗️ Simplified Dependency Management
- **No FX complexity**: Direct dependency creation and injection
- **Clear error messages**: Compile-time errors for missing dependencies
- **Easy debugging**: No magic dependency injection to debug
- **Simple testing**: Direct dependency injection for tests

### 🔧 Maintainable Code
- **Single responsibility**: Each file has one clear purpose
- **Easy to extend**: Add new dependencies without complex FX patterns
- **Clear dependencies**: Builder pattern shows explicit dependencies
- **Self-contained**: All example code demonstrates Builder API

### 🧪 Testable Design
- **Direct injection**: Easy to inject mocks and test dependencies
- **No FX setup**: Tests don't need complex FX container setup
- **Clear boundaries**: Builder pattern makes dependencies explicit
- **Simple mocking**: Direct dependency injection enables easy mocking

## Plugin Architecture (Unchanged)

The plugins themselves **don't need any changes** - they work with both FX and Builder APIs:

### Database Plugin
- **DatabaseConnectionService**: Manages database connections and lifecycle
- **QueryOperation**: Processes database queries with simple validation
- **Plugin Orchestrator**: Coordinates components and manages their lifecycle

### Web Server Plugin
- **HTTPService**: Manages HTTP server functionality and lifecycle
- **RouteOperation**: Processes HTTP route requests with simple routing
- **Plugin Orchestrator**: Coordinates components and manages their lifecycle

## Framework Integration (Builder API)

This example showcases the **Builder API** approach to the Fintechain Skeleton framework:

- **Builder Pattern**: Simple, explicit dependency injection
- **Component Lifecycle**: Managed initialization and cleanup (same as before)
- **Plugin Architecture**: Extensible system with clean interfaces (same as before)
- **Configuration Management**: Type-safe configuration access (same as before)
- **Event-Driven Communication**: Publish-subscribe messaging (same as before)
- **Structured Logging**: Consistent logging across all components (same as before)

## Migration from FX to Builder API

### What Changed
- ✅ **Removed `providers/` directory** - Builder API handles custom dependencies directly
- ✅ **Simplified `modes.go`** - No FX imports or complex provider patterns
- ✅ **Updated documentation** - Focus on Builder API benefits
- ✅ **Direct dependency creation** - Simple functions instead of FX providers

### What Didn't Change
- ✅ **Plugins remain the same** - They work with both APIs
- ✅ **Component lifecycle** - Same Initialize/Start/Stop patterns
- ✅ **Framework services** - Same logger, config, event bus interfaces
- ✅ **Application functionality** - Same behavior, simpler code

## Next Steps

1. **Run the examples** to see Builder API in action
2. **Compare with FX examples** to see the simplification
3. **Study the custom dependency creation** to understand the pattern
4. **Create your own applications** using the Builder API
5. **Migrate existing FX code** to the simpler Builder pattern

This example demonstrates how the **Builder API makes the Fintechain Skeleton framework much simpler to use** while maintaining all the same functionality. No more FX complexity - just simple, explicit dependency injection that's easy to understand, debug, and test. 