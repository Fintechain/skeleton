# Runtime Development Guide

A comprehensive guide to building applications with the Fintechain Skeleton framework runtime package - your **application entry point** for orchestrating plugins, services, and operations.

## 🎯 What is the Runtime Package?

The runtime package is **your application's main entry point** that orchestrates the entire Fintechain Skeleton framework. Think of it as the **conductor of an orchestra** - it doesn't play the music (that's what your plugins do), but it coordinates everything to work together harmoniously.

### Why Does Runtime Exist?

**Before Runtime** (Manual Setup):
```go
// You had to wire everything manually
registry := component.NewRegistry()
config := config.NewMemoryConfiguration()
logger := logging.NewNoOpLogger()
eventBus := event.NewEventBus()
pluginManager := plugin.NewManager()

runtime, err := runtime.NewRuntime(registry, config, pluginManager, eventBus, logger)
if err != nil {
    return err
}

ctx := context.NewContext()
if err := runtime.Start(ctx); err != nil {
    return err
}
defer runtime.Stop(ctx)

// Manual signal handling, error handling, lifecycle management...
```

**With Runtime Package** (Automatic):
```go
// Runtime handles everything automatically
import "github.com/fintechain/skeleton/pkg/runtime"

runtime.StartDaemon(
    runtime.WithPlugins(webPlugin, dbPlugin),
)
```

### What Runtime Does For You

1. **🏗️ Dependency Injection**: Automatically creates and wires all framework services (logger, config, event bus, etc.)
2. **🔄 Lifecycle Management**: Handles initialization, startup, shutdown, and cleanup in the correct order
3. **🔌 Plugin Orchestration**: Loads your plugins and ensures they can find each other
4. **⚡ Mode Selection**: Runs your app as either a long-running daemon or a one-shot command
5. **🛡️ Error Handling**: Provides graceful error handling and recovery
6. **📡 Signal Handling**: Automatically handles shutdown signals (CTRL+C, SIGTERM)

## 🏗️ How Applications Work in Skeleton Framework

### The Big Picture

```
Your Application
       ↓
Runtime Package (pkg/runtime)
       ↓
Framework Core (automatic dependency injection)
├── Configuration Service
├── Logging Service  
├── Event Bus Service
├── Component Registry
└── Plugin Manager
       ↓
Your Plugins (loaded by runtime)
├── Plugin A → registers components
├── Plugin B → registers components  
└── Plugin C → registers components
       ↓
Components (services & operations)
├── Web Server (service)
├── Database (service)
├── Calculator (operation)
└── File Processor (operation)
```

### Application Types

The runtime supports two fundamental application patterns:

#### 🔄 **Daemon Applications** - Long-Running Services
- **Examples**: Web servers, API services, background processors, monitoring systems
- **Behavior**: Starts up, runs continuously, handles requests/events until shutdown
- **Use Case**: When your application needs to stay running and respond to external requests

#### ⚡ **Command Applications** - Execute and Exit  
- **Examples**: CLI tools, batch processors, data migration scripts, calculators
- **Behavior**: Starts up, executes a specific task, returns result, exits immediately
- **Use Case**: When your application performs a specific task and then terminates

## 📚 Runtime API Reference

### Core Functions

#### StartDaemon
```go
func StartDaemon(opts ...Option) error
```
Starts a long-running daemon application that blocks until shutdown (SIGINT/SIGTERM).

#### ExecuteCommand
```go
func ExecuteCommand(operationID string, input map[string]interface{}, opts ...Option) (map[string]interface{}, error)
```
Executes a specific operation and returns immediately.

**Parameters:**
- `operationID`: String ID of the registered operation component
- `input`: Input data map (use `map[string]interface{}{}` for no input)
- `opts`: Configuration options (plugins, FX options)

**Returns:**
- `map[string]interface{}`: Operation results
- `error`: Execution error if any

#### StartDaemonWithSignalHandling
```go
func StartDaemonWithSignalHandling(signals []os.Signal, opts ...Option) error
```
Starts a daemon with custom signal handling beyond the default SIGINT/SIGTERM.

```go
import (
    "os"
    "syscall"
)

err := runtime.StartDaemonWithSignalHandling(
    []os.Signal{syscall.SIGTERM, syscall.SIGHUP}, 
    runtime.WithPlugins(myPlugin),
)
```

### Configuration Types

#### Option and Config
```go
type Option func(*Config)

type Config struct {
    Plugins      []plugin.Plugin
    ExtraOptions []fx.Option
}
```

#### WithPlugins
```go
func WithPlugins(plugins ...plugin.Plugin) Option
```
Adds plugins to be loaded at startup.

#### WithOptions
```go
func WithOptions(options ...fx.Option) Option
```
Adds custom FX dependency injection options.

### RuntimeEnvironment Type
```go
type RuntimeEnvironment = domainRuntime.RuntimeEnvironment
```
Provides access to framework services within components. Available in plugin Initialize methods via `system.(runtime.RuntimeEnvironment)`.

## 🚀 Building Your First Application

### Daemon Application Example

Let's build a simple web API server:

```go
package main

import (
    "log"
    "github.com/fintechain/skeleton/pkg/runtime"
)

func main() {
    // Start a daemon application
    err := runtime.StartDaemon(
        runtime.WithPlugins(
            NewWebServerPlugin(),    // Provides HTTP server
            NewDatabasePlugin(),     // Provides database access
            NewAPIPlugin(),          // Provides API endpoints
        ),
    )
    if err != nil {
        log.Fatal("Failed to start application:", err)
    }
    // Application runs until CTRL+C or SIGTERM
}
```

**What happens when you run this:**

1. **Framework Startup**: Runtime creates logger, config, event bus, registry, plugin manager
2. **Plugin Loading**: Each plugin gets `Initialize()` called to register their components
3. **Service Startup**: All services (like web server, database) get `Start()` called
4. **Running State**: Application handles requests, processes events, runs background tasks
5. **Graceful Shutdown**: On CTRL+C, all services get `Stop()` called in reverse order

### Command Application Example

Let's build a data processing CLI tool:

```go
package main

import (
    "fmt"
    "log"
    "os"
    "github.com/fintechain/skeleton/pkg/runtime"
)

func main() {
    if len(os.Args) < 2 {
        log.Fatal("Usage: processor <input-file>")
    }
    
    inputFile := os.Args[1]
    
    // Execute a command and get result
    result, err := runtime.ExecuteCommand("process-file", 
        map[string]interface{}{
            "inputFile": inputFile,
            "format":    "json",
        },
        runtime.WithPlugins(NewFileProcessorPlugin(), NewValidatorPlugin()),
    )
    if err != nil {
        log.Fatal("Processing failed:", err)
    }
    
    fmt.Printf("Processed %d records\n", result["recordCount"])
    fmt.Printf("Output file: %s\n", result["outputFile"])
}
```

**What happens when you run this:**

1. **Framework Startup**: Runtime creates all framework services (but doesn't start long-running services)
2. **Plugin Loading**: Plugins register their operations (like "process-file")
3. **Operation Execution**: The "process-file" operation is found and executed with your input
4. **Result Return**: Operation returns results, application prints them and exits
5. **Quick Cleanup**: Framework cleans up and terminates

## 🔧 Configuration and Customization

### Default Configuration (Zero Setup)

The runtime provides sensible defaults so you can get started immediately:

```go
// This works out of the box - no configuration needed!
runtime.StartDaemon(runtime.WithPlugins(myPlugin))
```

**What you get automatically:**
- Memory-based configuration (perfect for development)
- Console logging with structured output
- In-memory event bus for plugin communication
- Component registry for service discovery
- Plugin manager for lifecycle coordination

### Custom Configuration

When you need custom configuration (database connections, API keys, etc.):

```go
import (
    "go.uber.org/fx"
    "github.com/fintechain/skeleton/pkg/runtime"
    "github.com/fintechain/skeleton/internal/domain/config"
    infraConfig "github.com/fintechain/skeleton/internal/infrastructure/config"
)

func createAppConfig() config.Configuration {
    // Create configuration with your app's settings
    settings := map[string]interface{}{
        "app.name":        "My Application",
        "app.port":        8080,
        "database.host":   "localhost",
        "database.port":   5432,
        "api.key":         os.Getenv("API_KEY"),
        "log.level":       "info",
    }
    return infraConfig.NewMemoryConfigurationWithData(settings)
}

func main() {
    err := runtime.StartDaemon(
        runtime.WithPlugins(myWebPlugin, myDatabasePlugin),
        runtime.WithOptions(
            // Replace default config with your custom config
            fx.Replace(fx.Annotate(createAppConfig, fx.As(new(config.Configuration)))),
        ),
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

### Advanced Customization

For advanced scenarios, you can customize any part of the framework:

```go
func main() {
    err := runtime.StartDaemon(
        runtime.WithPlugins(myPlugin),
        runtime.WithOptions(
            // Add your own services to the dependency injection container
            fx.Provide(func() *MyExternalAPI {
                return &MyExternalAPI{
                    APIKey: os.Getenv("EXTERNAL_API_KEY"),
                    BaseURL: "https://api.example.com",
                }
            }),
            
            // Add lifecycle hooks for your services
            fx.Invoke(func(lc fx.Lifecycle, api *MyExternalAPI) {
                lc.Append(fx.Hook{
                    OnStart: func(ctx context.Context) error {
                        return api.Connect()
                    },
                    OnStop: func(ctx context.Context) error {
                        return api.Disconnect()
                    },
                })
            }),
            
            // Replace framework services with your implementations
            fx.Replace(fx.Annotate(myCustomLogger, fx.As(new(logging.LoggerService)))),
        ),
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

## 🔌 Working with Plugins

### What Are Plugins?

Plugins are **packages of functionality** that you load into your application. Each plugin can provide:

- **Services**: Long-running components (web servers, database connections, background workers)
- **Operations**: Executable tasks (calculations, data transformations, file processing)
- **Both**: A plugin can provide multiple services and operations

### Loading Plugins

```go
runtime.StartDaemon(
    runtime.WithPlugins(
        webServerPlugin,     // Provides HTTP server service
        databasePlugin,      // Provides database connection service  
        calculatorPlugin,    // Provides math operations
        fileProcessorPlugin, // Provides file processing operations
    ),
)
```

### Plugin Communication

Plugins automatically discover each other through the framework:

```go
// In your plugin's component
func (c *MyComponent) Initialize(ctx context.Context, system component.System) error {
    // Store runtime reference to access framework services
    c.runtime = system.(runtime.RuntimeEnvironment)
    
    // Access other components through the registry
    registry := system.Registry()
    
    // Find database service provided by database plugin
    dbComp, err := registry.Get("database-service")
    if err == nil {
        c.database = dbComp.(*DatabaseService)
    }
    
    // Access framework services
    logger := c.runtime.Logger()
    config := c.runtime.Configuration()
    eventBus := c.runtime.EventBus()
    
    logger.Info("Component initialized with dependencies")
    return nil
}
```

## 🎯 Common Application Patterns

### Pattern 1: Web API Server

```go
func main() {
    err := runtime.StartDaemon(
        runtime.WithPlugins(
            NewWebServerPlugin(8080),
            NewDatabasePlugin("postgres", "connection-string"),
            NewAuthPlugin(),
            NewAPIPlugin(), // Registers HTTP routes
        ),
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

### Pattern 2: Background Worker

```go
func main() {
    err := runtime.StartDaemon(
        runtime.WithPlugins(
            NewMessageQueuePlugin("rabbitmq://localhost"),
            NewDatabasePlugin("postgres", "connection-string"),
            NewWorkerPlugin(), // Processes messages from queue
        ),
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

### Pattern 3: CLI Data Processor

```go
func main() {
    inputFile := os.Args[1]
    outputFile := os.Args[2]
    
    result, err := runtime.ExecuteCommand("transform-data",
        map[string]interface{}{
            "input":  inputFile,
            "output": outputFile,
            "format": "csv",
        },
        runtime.WithPlugins(
            NewFileIOPlugin(),
            NewDataTransformPlugin(),
            NewValidationPlugin(),
        ),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Transformed %d records\n", result["count"])
}
```

### Pattern 4: Microservice with Multiple Protocols

```go
func main() {
    err := runtime.StartDaemon(
        runtime.WithPlugins(
            NewHTTPServerPlugin(8080),      // REST API
            NewGRPCServerPlugin(9090),      // gRPC API  
            NewMessageQueuePlugin("amqp://localhost"), // Async messaging
            NewDatabasePlugin("postgres", "connection-string"),
            NewBusinessLogicPlugin(),       // Your core functionality
        ),
    )
    if err != nil {
        log.Fatal(err)
    }
}
```

## 🔄 Application Lifecycle

### Daemon Mode Lifecycle

```
1. Framework Initialization
   ├── Create logger, config, event bus, registry, plugin manager
   └── Set up dependency injection container

2. Plugin Loading  
   ├── Call plugin.Initialize() for each plugin
   ├── Plugins register their services and operations
   └── Framework resolves dependencies between components

3. Service Startup
   ├── Call service.Start() for all services
   ├── Services begin accepting requests/connections
   └── Background workers start processing

4. Running State
   ├── Handle HTTP requests, process messages, run scheduled tasks
   ├── Components communicate through events and direct calls
   └── Framework monitors health and handles errors

5. Graceful Shutdown (on CTRL+C or SIGTERM)
   ├── Stop accepting new requests
   ├── Finish processing current requests
   ├── Call service.Stop() for all services (in reverse order)
   └── Clean up resources and exit
```

### Command Mode Lifecycle

```
1. Framework Initialization
   ├── Create logger, config, event bus, registry, plugin manager
   └── Set up dependency injection container

2. Plugin Loading
   ├── Call plugin.Initialize() for each plugin  
   ├── Plugins register their operations (services are NOT started)
   └── Framework resolves dependencies between components

3. Operation Execution
   ├── Find the requested operation in the registry
   ├── Execute operation with provided input data
   └── Return operation result

4. Quick Cleanup
   ├── Clean up framework resources
   └── Exit with result or error code
```

## 🛠️ Debugging and Troubleshooting

### Common Issues

**Plugin Not Found**:
```go
// Make sure plugin is loaded
runtime.StartDaemon(
    runtime.WithPlugins(myPlugin), // ← Plugin must be in this list
)
```

**Operation Not Found**:
```go
// Make sure operation is registered by a plugin
func (p *MyPlugin) Initialize(ctx context.Context, system component.System) error {
    registry := system.Registry()
    return registry.Register(NewMyOperation()) // ← Operation must be registered
}
```

**Service Won't Start**:
```go
// Check service dependencies in Initialize()
func (s *MyService) Initialize(ctx context.Context, system component.System) error {
    // Make sure all dependencies are available
    registry := system.Registry()
    dep, err := registry.Get("required-dependency")
    if err != nil {
        return fmt.Errorf("missing dependency: %w", err)
    }
    s.dependency = dep
    return nil
}
```

### Enabling Debug Logging

```go
func createDebugConfig() config.Configuration {
    settings := map[string]interface{}{
        "log.level": "debug", // Enable debug logging
    }
    return infraConfig.NewMemoryConfigurationWithData(settings)
}

runtime.StartDaemon(
    runtime.WithPlugins(myPlugin),
    runtime.WithOptions(
        fx.Replace(fx.Annotate(createDebugConfig, fx.As(new(config.Configuration)))),
    ),
)
```

## 🎯 Best Practices

### ✅ Do

1. **Choose the right mode**: Use daemon for servers, command for CLI tools
2. **Load plugins in logical order**: Dependencies first, then dependents
3. **Handle errors gracefully**: Always check errors from runtime functions
4. **Use configuration**: Don't hardcode values, use the config service
5. **Let runtime manage lifecycle**: Don't manually start/stop services
6. **Keep main() simple**: Put complex logic in plugins, not main()

### ❌ Don't

1. **Mix modes**: Don't try to run daemon and command patterns together
2. **Skip error handling**: Runtime functions return meaningful errors
3. **Manage lifecycle manually**: Let runtime handle start/stop sequences
4. **Hardcode dependencies**: Use the registry to find other components
5. **Block in main()**: Runtime handles blocking and signal handling
6. **Ignore shutdown signals**: Runtime handles CTRL+C automatically

## 🚀 Next Steps

1. **Start Simple**: Begin with a basic daemon or command application
2. **Add Plugins**: Create plugins for your specific functionality
3. **Configure as Needed**: Add custom configuration when you need it
4. **Scale Up**: Add more plugins and services as your application grows
5. **Deploy**: Use the daemon mode for production deployments

---

**Remember**: The runtime package is your **application orchestrator**. It handles all the framework complexity so you can focus on building your business logic in plugins. For plugin development details, see the [Plugin Development Guide](PLUGIN_DEVELOPMENT_GUIDE.md). For service and operation implementation, see the [Service Operations Development Guide](SERVICE_OPERATIONS_DEVELOPMENT_GUIDE.md).