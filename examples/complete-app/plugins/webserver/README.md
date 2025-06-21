# WebServer Plugin

A simplified web server plugin demonstrating **framework patterns** for the Fintechain Skeleton framework.

## 🎯 Purpose

This plugin demonstrates **framework usage patterns**, not real web server functionality:

- ✅ **Service lifecycle management**
- ✅ **Operation input/output processing**
- ✅ **Plugin orchestration patterns**
- ✅ **Runtime reference storage**
- ✅ **Framework services integration**

## 🏗️ Architecture

### Components

#### `HTTPService` - Service Component
Demonstrates service lifecycle with framework integration.

**Key Patterns**:
- Stores runtime reference for framework services access
- Simulates HTTP server start/stop (no real server)
- Uses configuration with defaults
- Structured logging with context

```go
type HTTPService struct {
    *component.BaseService
    runtime runtime.RuntimeEnvironment // Store runtime reference
    port    int
}
```

#### `RouteOperation` - Operation Component
Demonstrates simple input/output operation processing.

**Key Patterns**:
- Stores runtime reference for framework services access
- Simple input transformation (no real HTTP processing)
- Structured output format
- Framework logging integration

```go
type RouteOperation struct {
    *component.BaseOperation
    runtime runtime.RuntimeEnvironment // Store runtime reference
}
```

#### `WebServerPlugin` - Plugin Orchestrator
Demonstrates plugin-as-orchestrator pattern.

**Key Patterns**:
- Initializes and registers components
- Manages service lifecycle
- Stores runtime reference
- Component coordination

## 🚀 Usage Examples

### Daemon Mode (Service Lifecycle)

```go
func main() {
    err := runtime.NewBuilder().
        WithPlugins(webserver.NewWebServerPlugin(8080)).
        BuildDaemon()
    if err != nil {
        log.Fatal(err)
    }
    // HTTPService starts and runs until shutdown signal
}
```

**What happens**:
1. Plugin initializes and registers HTTPService and RouteOperation
2. HTTPService.Start() is called (simulates server start)
3. Service runs until SIGINT/SIGTERM
4. HTTPService.Stop() is called (simulates graceful shutdown)

### Command Mode (Operation Execution)

```go
func main() {
    result, err := runtime.NewBuilder().
        WithPlugins(webserver.NewWebServerPlugin(8080)).
        BuildCommand("http-route", map[string]interface{}{
            "method": "GET",
            "path":   "/api/health",
        })
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Result: %+v\n", result)
}
```

**What happens**:
1. Plugin initializes and registers components (HTTPService NOT started)
2. RouteOperation.Execute() is called with input data
3. Simple route processing returns structured output
4. Application exits immediately

## 🔧 Configuration

The plugin demonstrates configuration access patterns:

```json
{
  "http": {
    "port": 8080,
    "host": "0.0.0.0"
  }
}
```

**Configuration Access Pattern**:
```go
func (h *HTTPService) Start(ctx context.Context) error {
    config := h.runtime.Configuration()
    
    // Type-safe access with defaults
    port := config.GetIntDefault("http.port", h.port)
    host := config.GetStringDefault("http.host", "0.0.0.0")
    
    // Use configuration values...
}
```

## 📊 Framework Patterns Demonstrated

### 1. Runtime Reference Storage
```go
func (h *HTTPService) Initialize(ctx context.Context, system component.System) error {
    // Store runtime reference - critical pattern
    h.runtime = system.(runtime.RuntimeEnvironment)
    
    // Access framework services
    logger := h.runtime.Logger()
    logger.Info("HTTP Service initialized", "component_id", h.ID())
    
    return nil
}
```

### 2. Plugin Orchestration
```go
func (w *WebServerPlugin) Initialize(ctx context.Context, system component.System) error {
    // 1. Initialize components
    w.httpService.Initialize(ctx, system)
    w.routeOperation.Initialize(ctx, system)
    
    // 2. Register with registry
    registry := system.Registry()
    registry.Register(w.httpService)
    registry.Register(w.routeOperation)
    
    return nil
}
```

### 3. Service Lifecycle Management
```go
func (w *WebServerPlugin) Start(ctx context.Context) error {
    // Plugin manages service lifecycle
    return w.httpService.Start(ctx)
}

func (w *WebServerPlugin) Stop(ctx context.Context) error {
    // Plugin manages service cleanup
    return w.httpService.Stop(ctx)
}
```

### 4. Simple Operation Processing
```go
func (r *RouteOperation) Execute(ctx context.Context, input component.Input) (component.Output, error) {
    // Simple input parsing
    data := input.Data.(map[string]interface{})
    method := data["method"].(string)
    path := data["path"].(string)
    
    // Simple processing (no real HTTP logic)
    response := r.processRoute(method, path)
    
    return component.Output{Data: response}, nil
}
```

## 🧪 Testing

### Component Testing
```go
func TestHTTPService(t *testing.T) {
    runtime, err := runtime.NewRuntimeWithOptions(
        runtime.WithPlugins(webserver.NewWebServerPlugin(8080)),
    )
    require.NoError(t, err)
    
    ctx := context.NewContext()
    err = runtime.Start(ctx)
    require.NoError(t, err)
    defer runtime.Stop(ctx)
    
    // Service should be running
    registry := runtime.Registry()
    service, err := registry.Get("http-server")
    require.NoError(t, err)
    
    httpService := service.(*webserver.HTTPService)
    assert.True(t, httpService.IsRunning())
}
```

### Operation Testing
```go
func TestRouteOperation(t *testing.T) {
    result, err := runtime.NewBuilder().
        WithPlugins(webserver.NewWebServerPlugin(8080)).
        BuildCommand("http-route", map[string]interface{}{
            "method": "GET",
            "path":   "/api/health",
        })
    
    require.NoError(t, err)
    assert.Equal(t, "healthy", result["status"])
}
```

## 🎯 Key Simplifications

### What This Plugin Demonstrates:
- ✅ **Framework Patterns**: Component lifecycle, plugin orchestration
- ✅ **Service Management**: Start/stop lifecycle in daemon mode
- ✅ **Operation Processing**: Simple input/output transformation
- ✅ **Configuration Access**: Type-safe config with defaults
- ✅ **Logging Integration**: Structured logging throughout

### What This Plugin Avoids:
- ❌ **Real HTTP Server**: No actual networking or HTTP handling
- ❌ **Complex Routing**: No real URL routing or middleware
- ❌ **JSON Processing**: No complex request/response parsing
- ❌ **Error Handling**: No HTTP-specific error codes
- ❌ **External Dependencies**: No web framework dependencies

## 📚 Next Steps

1. **Copy this plugin** as a template for your own plugins
2. **Modify the components** for your specific needs
3. **Keep the framework patterns** - runtime reference, lifecycle management
4. **Focus on your business logic** - not framework demonstration
5. **Test both modes** - daemon and command execution

---

**Remember**: This plugin demonstrates **framework usage patterns**. For real web server functionality, implement actual HTTP handling while maintaining these framework patterns. 