# Complete Application Example

This example demonstrates a **complete multi-plugin application** using the Fintechain Skeleton framework with both web server and database plugins working together, focusing on **framework patterns**.

## 🎯 What This Demonstrates

- **Plugin Coordination**: Multiple plugins working in the same application
- **Framework Patterns**: Component lifecycle, service management, operation execution
- **Runtime Modes**: Daemon vs command mode for different use cases
- **Simplified Operations**: Focus on framework usage rather than business complexity

## 🚀 How to Run

### 1. Daemon Mode - Long-Running Services
Starts both web server and database services:

```bash
go run examples/complete-app/main.go daemon
```

**What happens:**
- WebServer plugin initializes HTTPService and RouteOperation
- Database plugin initializes DatabaseConnectionService and QueryOperation
- Both services start and simulate running (no real servers)
- Services run until Ctrl+C with graceful shutdown
- Framework logging shows service lifecycle

### 2. Command Mode - Execute Operation and Exit
Tests operation execution without starting services:

```bash
go run examples/complete-app/main.go command
```

**What happens:**
- Both plugins initialize their components
- RouteOperation executes with sample HTTP route data
- Simple route processing returns structured output
- Application exits immediately (no services started)
- Demonstrates operation execution pattern

## 🏗️ Architecture

```
Complete Application
├── WebServerPlugin (Framework Patterns)
│   ├── HTTPService (simulated service lifecycle)
│   └── RouteOperation (simple input/output processing)
└── DatabasePlugin (Framework Patterns)
    ├── DatabaseConnectionService (simulated connection lifecycle)
    └── QueryOperation (simple query processing)
```

## 💡 Key Learning Points

### Framework Patterns Focus
- **Runtime Reference Storage**: All components store runtime reference
- **Plugin Orchestration**: Plugins initialize and register components
- **Service Lifecycle**: Start/stop management in daemon mode
- **Operation Processing**: Simple input/output transformation

### Plugin Independence
- Each plugin can work independently
- Plugins don't directly depend on each other
- Framework handles component registration and lifecycle

### Mode Separation
- **Daemon Mode**: Services start and run continuously (simulated)
- **Command Mode**: Operations execute and complete immediately
- Same plugins work in both modes with different behavior

### Component Communication
- Components access framework services through stored runtime reference
- Registry used for component discovery
- Clean separation of concerns

## 🔧 Configuration

The application demonstrates configuration access patterns:

```json
{
  "http": {
    "port": 8080,
    "host": "0.0.0.0"
  },
  "database": {
    "driver": "postgres",
    "datasource": "test://connection",
    "max_connections": 10
  }
}
```

**Configuration Usage Pattern**:
```go
// Components access config through runtime reference
config := component.runtime.Configuration()
port := config.GetIntDefault("http.port", 8080)
```

## 📊 Framework Patterns Demonstrated

### 1. Plugin-as-Orchestrator
```go
func (p *WebServerPlugin) Initialize(ctx context.Context, system component.System) error {
    // 1. Initialize components
    p.httpService.Initialize(ctx, system)
    p.routeOperation.Initialize(ctx, system)
    
    // 2. Register with registry
    registry := system.Registry()
    registry.Register(p.httpService)
    registry.Register(p.routeOperation)
    
    return nil
}
```

### 2. Runtime Reference Storage
```go
func (h *HTTPService) Initialize(ctx context.Context, system component.System) error {
    // Store runtime reference for framework services access
    h.runtime = system.(runtime.RuntimeEnvironment)
    
    // Access framework services
    logger := h.runtime.Logger()
    logger.Info("HTTP Service initialized")
    
    return nil
}
```

### 3. Service Lifecycle Management
```go
func (p *WebServerPlugin) Start(ctx context.Context) error {
    // Plugin manages service lifecycle
    return p.httpService.Start(ctx)
}
```

## 🎯 Key Simplifications

### What This Example Demonstrates:
- ✅ **Framework Patterns**: Component lifecycle, plugin orchestration
- ✅ **Multi-Plugin Coordination**: Multiple plugins working together
- ✅ **Runtime Modes**: Daemon vs command execution patterns
- ✅ **Service Management**: Start/stop lifecycle simulation
- ✅ **Operation Processing**: Simple input/output transformation

### What This Example Avoids:
- ❌ **Real Infrastructure**: No actual HTTP servers or database connections
- ❌ **Complex Business Logic**: No real web routing or SQL processing
- ❌ **External Dependencies**: No web frameworks or database drivers
- ❌ **Complex Error Handling**: Focus on framework patterns
- ❌ **Production Concerns**: No real networking or persistence

## 🧪 Testing the Example

### Test Daemon Mode
```bash
go run examples/complete-app/main.go daemon
# Watch the logs to see:
# - Plugin initialization
# - Component registration
# - Service startup
# - Graceful shutdown on Ctrl+C
```

### Test Command Mode
```bash
go run examples/complete-app/main.go command
# Watch the logs to see:
# - Plugin initialization
# - Operation execution
# - Immediate cleanup and exit
```

## 🎯 Next Steps

1. **Run both modes** to see different execution patterns
2. **Check the plugin code** in `examples/plugins/` to understand implementation
3. **Copy plugins as templates** for your own applications
4. **Read the Plugin Development Guide** for advanced patterns
5. **Focus on framework patterns** when building your own plugins

This example serves as a **framework pattern template** showing how multiple plugins coordinate while maintaining clean separation of concerns. 