# Complete Application Example

This example demonstrates a complete application using the Fintechain Skeleton framework, showcasing all major framework patterns in a well-organized, self-contained structure.

## Project Structure

```
examples/complete-app/
├── main.go                 # Entry point and mode selection
├── modes/                  # Different execution modes
│   └── modes.go           # Daemon, command, and custom provider modes
├── providers/             # Custom service implementations
│   ├── logger.go         # Custom logger implementation
│   ├── config.go         # Custom configuration implementation
│   └── eventbus.go       # Custom event bus implementation
├── plugins/               # Example plugins (self-contained)
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

## What This Example Demonstrates

### 🔌 Plugin Orchestration
- **Multiple plugins working together**: Webserver + Database plugins
- **Plugin lifecycle management**: Automatic initialization, startup, and shutdown
- **Plugin communication**: Through the shared runtime environment
- **Self-contained plugins**: Each plugin includes all its components

### 🛠️ Custom Providers
- **Custom Logger**: Structured logging with timestamps and prefixes
- **Custom Configuration**: In-memory configuration with type-safe access
- **Custom Event Bus**: Async event publishing with subscription management

### 🔄 Component Lifecycle
- **Service lifecycle**: Initialize → Start → Stop → Dispose
- **Status tracking**: Monitor service states (Stopped, Running, etc.)
- **Graceful shutdown**: Proper cleanup of all resources

### 🚀 Execution Modes
- **Daemon Mode**: Long-running services (web servers, background workers)
- **Command Mode**: Execute operations and exit (CLI commands, batch processing)
- **Custom Providers**: Replace framework services with custom implementations

## Usage

### Daemon Mode (Long-running Services)
```bash
go run examples/complete-app/main.go daemon
```

This mode:
- Starts multiple plugins as long-running services
- Uses default framework providers
- Blocks until shutdown signal (SIGINT/SIGTERM)
- Demonstrates typical server/service applications

### Command Mode (Execute and Exit)
```bash
go run examples/complete-app/main.go command
```

This mode:
- Executes a specific operation
- Returns results immediately
- Exits after completion
- Demonstrates CLI tools and batch processing

### Custom Providers Mode
```bash
go run examples/complete-app/main.go custom
```

This mode:
- Uses custom implementations of framework services
- Shows how to replace logger, config, and event bus
- Demonstrates advanced FX dependency injection
- Perfect for specialized requirements

## Key Framework Patterns

### 1. Plugin Development
```go
// Plugins implement the plugin.Plugin interface
type MyPlugin struct {
    // plugin configuration
}

func (p *MyPlugin) Initialize(ctx context.Context, system component.System) error {
    // Setup plugin resources
}
```

### 2. Custom Service Providers
```go
// Custom services implement framework interfaces
func NewCustomLogger() logging.LoggerService {
    return &CustomLogger{
        prefix: "[CUSTOM]",
        status: component.StatusStopped,
    }
}

// Replace framework services using FX
runtime.WithOptions(
    fx.Replace(NewCustomLogger()),
)
```

### 3. Runtime Modes
```go
// Daemon mode - long-running services
runtime.StartDaemon(
    runtime.WithPlugins(plugin1, plugin2),
)

// Command mode - execute and exit
result, err := runtime.ExecuteCommand("operation-id", input,
    runtime.WithPlugins(plugin1, plugin2),
)
```

## Architecture Benefits

### 🏗️ Clean Separation of Concerns
- **main.go**: Entry point and mode selection only
- **modes/**: Execution logic separated by use case
- **providers/**: Custom implementations in dedicated files
- **plugins/**: Self-contained plugin implementations

### 🔧 Maintainable Code
- **Single responsibility**: Each file has one clear purpose
- **Easy to extend**: Add new modes, providers, or plugins without touching existing code
- **Clear dependencies**: Import structure shows relationships
- **Self-contained**: All example code is in one place

### 🧪 Testable Design
- **Isolated components**: Each provider and plugin can be tested independently
- **Mockable interfaces**: Framework interfaces enable easy mocking
- **Clear boundaries**: Separation makes unit testing straightforward

## Plugin Architecture

### Database Plugin
- **DatabaseConnectionService**: Manages database connections and lifecycle
- **QueryOperation**: Processes database queries with simple validation
- **Plugin Orchestrator**: Coordinates components and manages their lifecycle

### Web Server Plugin
- **HTTPService**: Manages HTTP server functionality and lifecycle
- **RouteOperation**: Processes HTTP route requests with simple routing
- **Plugin Orchestrator**: Coordinates components and manages their lifecycle

## Framework Integration

This example showcases the full power of the Fintechain Skeleton framework:

- **FX Dependency Injection**: Automatic wiring of services and plugins
- **Component Lifecycle**: Managed initialization and cleanup
- **Plugin Architecture**: Extensible system with clean interfaces
- **Configuration Management**: Type-safe configuration access
- **Event-Driven Communication**: Publish-subscribe messaging
- **Structured Logging**: Consistent logging across all components

## Next Steps

1. **Study the code structure** to understand the organization patterns
2. **Run different modes** to see various execution patterns
3. **Examine custom providers** to learn service replacement techniques
4. **Study the plugin implementations** to understand component orchestration
5. **Create your own plugins** following the established patterns
6. **Extend the example** with additional modes, providers, or plugins

This example serves as a comprehensive, self-contained reference for building production-ready applications with the Fintechain Skeleton framework. 