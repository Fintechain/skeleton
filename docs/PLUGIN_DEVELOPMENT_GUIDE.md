# Plugin Development Guide

A comprehensive guide to building plugins for the Fintechain Skeleton framework, focusing on **framework patterns** and best practices.

## 🎯 Philosophy: Framework Patterns Over Business Logic

This guide demonstrates **how to use the framework correctly** rather than implementing complex business functionality. The examples focus on:

- ✅ **Component lifecycle management**
- ✅ **Framework services integration** 
- ✅ **Plugin orchestration patterns**
- ✅ **Runtime reference storage**
- ✅ **Proper dependency injection**

## 🏗️ Plugin Architecture

### Plugin-as-Orchestrator Pattern

**Plugins don't implement functionality** - they **orchestrate components** that implement functionality.

```go
type MyPlugin struct {
    *component.BaseService
    runtime runtime.RuntimeEnvironment // Store runtime reference
    
    // Components this plugin orchestrates
    myService    *MyService
    myOperation  *MyOperation
}

func (p *MyPlugin) Initialize(ctx context.Context, system component.System) error {
    // 1. Store runtime reference
    p.runtime = system.(runtime.RuntimeEnvironment)
    
    // 2. Initialize components
    p.myService.Initialize(ctx, system)
    p.myOperation.Initialize(ctx, system)
    
    // 3. Register components
    registry := system.Registry()
    registry.Register(p.myService)
    registry.Register(p.myOperation)
    
    return nil
}

func (p *MyPlugin) Start(ctx context.Context) error {
    // Plugin manages service lifecycle
    return p.myService.Start(ctx)
}
```

## 🔧 Component Implementation Patterns

### 1. Runtime Reference Storage (Critical Pattern)

**All components must store the runtime reference** for framework services access:

```go
import (
    "github.com/fintechain/skeleton/pkg/component"
    "github.com/fintechain/skeleton/pkg/context"  // Framework context
    "github.com/fintechain/skeleton/pkg/runtime"
)

type MyComponent struct {
    *component.BaseComponent
    runtime runtime.RuntimeEnvironment // REQUIRED: Store runtime reference
}

func (c *MyComponent) Initialize(ctx context.Context, system component.System) error {
    if err := c.BaseComponent.Initialize(ctx, system); err != nil {
        return err
    }
    
    // CRITICAL: Store runtime reference for framework services
    c.runtime = system.(runtime.RuntimeEnvironment)
    
    // Access framework services to demonstrate pattern
    logger := c.runtime.Logger()
    config := c.runtime.Configuration()
    eventBus := c.runtime.EventBus()
    
    logger.Info("Component initialized", "component_id", c.ID())
    return nil
}
```

### 2. Service Components

Services are long-running components started in daemon mode:

```go
type MyService struct {
    *component.BaseService
    runtime runtime.RuntimeEnvironment // Store runtime reference
}

func NewMyService() *MyService {
    config := component.ComponentConfig{
        ID:          "my-service",
        Name:        "My Service",
        Description: "Demonstrates service patterns",
        Version:     "1.0.0",
    }
    return &MyService{
        BaseService: component.NewBaseService(config),
    }
}

func (s *MyService) Initialize(ctx context.Context, system component.System) error {
    if err := s.BaseService.Initialize(ctx, system); err != nil {
        return err
    }
    
    // Store runtime reference
    s.runtime = system.(runtime.RuntimeEnvironment)
    
    // Access framework services
    logger := s.runtime.Logger()
    logger.Info("Service initialized", "service_id", s.ID())
    
    return nil
}

func (s *MyService) Start(ctx context.Context) error {
    if err := s.BaseService.Start(ctx); err != nil {
        return err
    }
    
    // Access framework services through stored runtime reference
    logger := s.runtime.Logger()
    config := s.runtime.Configuration()
    
    // Simulate service start (focus on framework patterns)
    port := config.GetIntDefault("my_service.port", 8080)
    logger.Info("Service started", 
        "service_id", s.ID(),
        "port", port,
        "status", "running")
    
    return nil
}

func (s *MyService) Stop(ctx context.Context) error {
    // Access framework services
    logger := s.runtime.Logger()
    logger.Info("Service stopping", "service_id", s.ID())
    
    // Simulate cleanup
    logger.Info("Service stopped successfully")
    
    return s.BaseService.Stop(ctx)
}
```

### 3. Operation Components

Operations are stateless components that process input and return output:

```go
type MyOperation struct {
    *component.BaseOperation
    runtime runtime.RuntimeEnvironment // Store runtime reference
}

func NewMyOperation() *MyOperation {
    config := component.ComponentConfig{
        ID:          "my-operation",
        Name:        "My Operation",
        Description: "Demonstrates operation patterns",
        Version:     "1.0.0",
    }
    return &MyOperation{
        BaseOperation: component.NewBaseOperation(config),
    }
}

func (o *MyOperation) Initialize(ctx context.Context, system component.System) error {
    if err := o.BaseOperation.Initialize(ctx, system); err != nil {
        return err
    }
    
    // Store runtime reference
    o.runtime = system.(runtime.RuntimeEnvironment)
    
    // Access framework services
    logger := o.runtime.Logger()
    logger.Info("Operation initialized", "operation_id", o.ID())
    
    return nil
}

func (o *MyOperation) Execute(ctx context.Context, input component.Input) (component.Output, error) {
    // Access framework services through stored runtime reference
    logger := o.runtime.Logger()
    
    // Parse input (keep simple)
    data, ok := input.Data.(map[string]interface{})
    if !ok {
        return component.Output{}, fmt.Errorf("invalid input data format")
    }
    
    // Simple processing (focus on framework patterns, not business logic)
    name, _ := data["name"].(string)
    value, _ := data["value"].(float64)
    
    logger.Info("Processing operation", 
        "operation_id", o.ID(),
        "input_name", name,
        "input_value", value)
    
    // Simple transformation
    result := map[string]interface{}{
        "processed_name":  name,
        "processed_value": value * 2, // Simple transformation
        "status":         "success",
        "timestamp":      time.Now().Format(time.RFC3339),
    }
    
    return component.Output{Data: result}, nil
}
```

## 🔌 Complete Plugin Example

Here's a complete plugin following all framework patterns:

```go
package myplugin

import (
    "github.com/fintechain/skeleton/pkg/component"
    "github.com/fintechain/skeleton/pkg/context"
    "github.com/fintechain/skeleton/pkg/plugin"
    "github.com/fintechain/skeleton/pkg/runtime"
)

// MyPlugin demonstrates proper plugin patterns
type MyPlugin struct {
    *component.BaseService
    runtime   runtime.RuntimeEnvironment // Store runtime reference
    service   *MyService
    operation *MyOperation
}

func NewMyPlugin() *MyPlugin {
    config := component.ComponentConfig{
        ID:          "my-plugin",
        Name:        "My Plugin",
        Description: "Demonstrates plugin patterns",
        Version:     "1.0.0",
    }
    
    return &MyPlugin{
        BaseService: component.NewBaseService(config),
        service:     NewMyService(),
        operation:   NewMyOperation(),
    }
}

func (p *MyPlugin) Author() string {
    return "My Team"
}

func (p *MyPlugin) PluginType() plugin.PluginType {
    return plugin.TypeProcessor
}

// Initialize orchestrates component setup
func (p *MyPlugin) Initialize(ctx context.Context, system component.System) error {
    if err := p.BaseService.Initialize(ctx, system); err != nil {
        return err
    }
    
    // Store runtime reference
    p.runtime = system.(runtime.RuntimeEnvironment)
    logger := p.runtime.Logger()
    logger.Info("Initializing plugin", "plugin_id", p.ID())
    
    // 1. Initialize components this plugin provides
    if err := p.service.Initialize(ctx, system); err != nil {
        return err
    }
    
    if err := p.operation.Initialize(ctx, system); err != nil {
        return err
    }
    
    // 2. Register components with system registry
    registry := system.Registry()
    if err := registry.Register(p.service); err != nil {
        return err
    }
    
    if err := registry.Register(p.operation); err != nil {
        return err
    }
    
    logger.Info("Plugin initialized successfully",
        "components_registered", 2,
        "service_id", p.service.ID(),
        "operation_id", p.operation.ID())
    
    return nil
}

// Start manages service lifecycle (daemon mode)
func (p *MyPlugin) Start(ctx context.Context) error {
    if err := p.BaseService.Start(ctx); err != nil {
        return err
    }
    
    logger := p.runtime.Logger()
    logger.Info("Starting plugin", "plugin_id", p.ID())
    
    // Plugin's responsibility: Start services it manages
    if err := p.service.Start(ctx); err != nil {
        return err
    }
    
    logger.Info("Plugin started successfully")
    return nil
}

// Stop manages service cleanup
func (p *MyPlugin) Stop(ctx context.Context) error {
    logger := p.runtime.Logger()
    logger.Info("Stopping plugin", "plugin_id", p.ID())
    
    // Plugin's responsibility: Stop services it manages
    if err := p.service.Stop(ctx); err != nil {
        return err
    }
    
    if err := p.BaseService.Stop(ctx); err != nil {
        return err
    }
    
    logger.Info("Plugin stopped successfully")
    return nil
}
```

## 🚀 Usage Examples

### Daemon Mode (Long-Running Services)

```go
func main() {
    err := runtime.NewBuilder().
        WithPlugins(NewMyPlugin()).
        BuildDaemon()
    if err != nil {
        log.Fatal(err)
    }
    // Services run until shutdown signal
}
```

### Command Mode (Execute and Exit)

```go
func main() {
    result, err := runtime.NewBuilder().
        WithPlugins(NewMyPlugin()).
        BuildCommand("my-operation", map[string]interface{}{
            "name":  "test",
            "value": 42.0,
        })
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Result: %+v\n", result)
}
```

## 🧪 Testing Patterns

### Testing Components

```go
func TestMyOperation(t *testing.T) {
    // Use Builder API for test setup
    result, err := runtime.NewBuilder().
        WithPlugins(NewMyPlugin()).
        BuildCommand("my-operation", map[string]interface{}{
            "name":  "test",
            "value": 10.0,
        })
    
    require.NoError(t, err)
    
    assert.Equal(t, "test", result["processed_name"])
    assert.Equal(t, 20.0, result["processed_value"])
}
```

### Testing with Custom Dependencies

```go
func TestMyOperationWithCustomConfig(t *testing.T) {
    // Create custom configuration
    config := infraConfig.NewMemoryConfigurationWithData(map[string]interface{}{
        "my_component.multiplier": 3,
    })
    
    result, err := runtime.NewBuilder().
        WithPlugins(NewMyPlugin()).
        WithConfig(config).
        BuildCommand("my-operation", map[string]interface{}{
            "name":  "test",
            "value": 10.0,
        })
    
    require.NoError(t, err)
    assert.Equal(t, 30.0, result["processed_value"]) // 10 * 3
}
```

## 🎯 Framework Integration Patterns

### Configuration Access

```go
func (c *MyComponent) DoWork() {
    config := c.runtime.Configuration()
    
    // Type-safe configuration access with defaults
    timeout := config.GetDurationDefault("my_component.timeout", 30*time.Second)
    maxRetries := config.GetIntDefault("my_component.max_retries", 3)
    enabled := config.GetBoolDefault("my_component.enabled", true)
    
    if !enabled {
        return
    }
    
    // Use configuration values...
}
```

### Event System Integration

```go
func (c *MyComponent) PublishEvent() {
    eventBus := c.runtime.EventBus()
    
    event := &event.Event{
        Topic:   "my_component.work_completed",
        Source:  c.ID(),
        Time:    time.Now(),
        Payload: map[string]interface{}{
            "component_id": c.ID(),
            "status":      "success",
        },
    }
    
    eventBus.PublishAsync(event)
}

func (c *MyComponent) SubscribeToEvents() {
    eventBus := c.runtime.EventBus()
    
    subscription := eventBus.SubscribeAsync("system.shutdown", func(event *event.Event) {
        logger := c.runtime.Logger()
        logger.Info("Received shutdown event", "component_id", c.ID())
        // Handle shutdown preparation
    })
    
    // Store subscription for cleanup if needed
}
```

### Logging Patterns

```go
func (c *MyComponent) LogWithContext() {
    logger := c.runtime.Logger()
    
    // Structured logging with context
    logger.Info("Component operation started",
        "component_id", c.ID(),
        "operation", "process_data",
        "timestamp", time.Now(),
    )
    
    logger.Debug("Processing details",
        "component_id", c.ID(),
        "items_count", 42,
        "processing_time", "150ms",
    )
    
    logger.Error("Operation failed",
        "component_id", c.ID(),
        "error", "connection timeout",
        "retry_count", 3,
    )
}
```

## 🤝 Best Practices

### ✅ Do

1. **Store runtime reference** in all components
2. **Use plugin-as-orchestrator** pattern
3. **Access framework services** through runtime reference
4. **Keep operations simple** - focus on input/output transformation
5. **Let plugins manage** service lifecycle
6. **Use structured logging** with component context
7. **Handle configuration** with defaults
8. **Test components** using traditional runtime for control

### ❌ Don't

1. **Implement real infrastructure** - focus on framework patterns
2. **Mix business logic** with framework demonstration
3. **Access system directly** - always use runtime reference
4. **Start services manually** - let plugins manage lifecycle
5. **Complicate operations** - keep them simple and focused
6. **Ignore error handling** - but keep it framework-focused
7. **Hardcode values** - use configuration with defaults

## 🔄 Component Lifecycle Summary

### Plugin Lifecycle
```
Initialize → Register Components → Start Services → Stop Services → Dispose
```

### Component Lifecycle
```
Create → Initialize (store runtime) → Register → Active → Dispose
```

### Service Lifecycle (Daemon Mode)
```
Initialize → Start → Running → Stop → Dispose
```

### Operation Lifecycle (Command Mode)
```
Initialize → Execute → Dispose
```

## 📚 Next Steps

1. **Copy a plugin template** from `examples/plugins/`
2. **Modify for your needs** - keep framework patterns
3. **Test with both runtime modes** - daemon and command
4. **Focus on framework integration** - not business logic
5. **Use structured logging** throughout
6. **Handle configuration properly** with defaults

---

**Remember**: Examples demonstrate **framework usage patterns**, not business functionality. Keep your plugins simple and focused on showing how to use the framework correctly. 