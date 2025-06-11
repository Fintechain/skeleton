# Service and Operations Development Guide

A comprehensive guide to building services and operations for the Fintechain Skeleton framework, focusing on **framework patterns** and best practices.

## 🎯 Philosophy: Framework Patterns Over Business Logic

This guide demonstrates **how to use the framework correctly** rather than implementing complex business functionality. The examples focus on:

- ✅ **Component lifecycle management**
- ✅ **Framework services integration** 
- ✅ **Runtime reference storage**
- ✅ **Proper dependency injection**
- ✅ **Service and operation patterns**

## 🏗️ Component Architecture

### Service vs Operation Decision Matrix

**Services** are long-running components that maintain state and run continuously:
- Web servers, message processors, background workers
- Have Start/Stop lifecycle with status management
- Run in daemon mode until shutdown

**Operations** are stateless components that process input and return output:
- Data transformers, calculators, validators
- Execute once and return results
- Run in command mode and exit

```go
// Service Example: Runs continuously
type BackgroundProcessor struct {
    *component.BaseService
    runtime runtime.RuntimeEnvironment
}

// Operation Example: Execute and return
type DataTransformer struct {
    *component.BaseOperation
    runtime runtime.RuntimeEnvironment
}
```

## 🔧 Service Implementation Patterns

### 1. Runtime Reference Storage (Critical Pattern)

**All services must store the runtime reference** for framework services access:

```go
import (
    "github.com/fintechain/skeleton/pkg/component"
    "github.com/fintechain/skeleton/pkg/context"  // Framework context
    "github.com/fintechain/skeleton/pkg/runtime"
)

type MyService struct {
    *component.BaseService
    runtime runtime.RuntimeEnvironment // REQUIRED: Store runtime reference
}

func (s *MyService) Initialize(ctx context.Context, system component.System) error {
    if err := s.BaseService.Initialize(ctx, system); err != nil {
        return err
    }
    
    // CRITICAL: Store runtime reference for framework services
    s.runtime = system.(runtime.RuntimeEnvironment)
    
    // Access framework services to demonstrate pattern
    logger := s.runtime.Logger()
    config := s.runtime.Configuration()
    eventBus := s.runtime.EventBus()
    
    logger.Info("Service initialized", "service_id", s.ID())
    return nil
}
```

### 2. Service Lifecycle Management

Services follow Start → Running → Stop lifecycle with status management:

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
    enabled := config.GetBoolDefault("my_service.enabled", true)
    
    if !enabled {
        logger.Info("Service disabled by configuration")
        return nil
    }
    
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

### 3. Background Processing Pattern

Services can run background tasks while maintaining framework integration:

```go
type BackgroundService struct {
    *component.BaseService
    runtime runtime.RuntimeEnvironment
    done    chan struct{}
}

func (b *BackgroundService) Start(ctx context.Context) error {
    if err := b.BaseService.Start(ctx); err != nil {
        return err
    }
    
    logger := b.runtime.Logger()
    config := b.runtime.Configuration()
    
    // Get configuration
    interval := config.GetDurationDefault("background.interval", 30*time.Second)
    
    b.done = make(chan struct{})
    
    // Start background goroutine
    go b.backgroundWorker(interval)
    
    logger.Info("Background service started", "interval", interval)
    return nil
}

func (b *BackgroundService) Stop(ctx context.Context) error {
    logger := b.runtime.Logger()
    logger.Info("Stopping background service")
    
    // Signal background worker to stop
    if b.done != nil {
        close(b.done)
    }
    
    logger.Info("Background service stopped")
    return b.BaseService.Stop(ctx)
}

func (b *BackgroundService) backgroundWorker(interval time.Duration) {
    logger := b.runtime.Logger()
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    
    for {
        select {
        case <-b.done:
            logger.Info("Background worker stopping")
            return
        case <-ticker.C:
            b.doWork()
        }
    }
}

func (b *BackgroundService) doWork() {
    logger := b.runtime.Logger()
    logger.Info("Background work executed", "timestamp", time.Now())
    
    // Simple work simulation - focus on framework patterns
    // Real implementations would do actual work here
}
```

## ⚡ Operation Implementation Patterns

### 1. Stateless Processing Pattern

Operations process input and return output without maintaining state:

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

### 2. Input Validation Pattern

Operations should validate input and handle errors gracefully:

```go
type ValidatingOperation struct {
    *component.BaseOperation
    runtime runtime.RuntimeEnvironment
}

func (v *ValidatingOperation) Execute(ctx context.Context, input component.Input) (component.Output, error) {
    logger := v.runtime.Logger()
    
    // Validate input structure
    data, ok := input.Data.(map[string]interface{})
    if !ok {
        logger.Error("Invalid input format", "expected", "map[string]interface{}")
        return component.Output{}, fmt.Errorf("input must be a map")
    }
    
    // Validate required fields
    if _, exists := data["required_field"]; !exists {
        logger.Error("Missing required field", "field", "required_field")
        return component.Output{}, fmt.Errorf("required_field is missing")
    }
    
    // Process valid input
    logger.Info("Input validated successfully")
    
    result := map[string]interface{}{
        "validation": "passed",
        "processed":  true,
    }
    
    return component.Output{Data: result}, nil
}
```

### 3. Configuration-Driven Operation

Operations can use configuration to modify behavior:

```go
type ConfigurableOperation struct {
    *component.BaseOperation
    runtime runtime.RuntimeEnvironment
}

func (c *ConfigurableOperation) Execute(ctx context.Context, input component.Input) (component.Output, error) {
    logger := c.runtime.Logger()
    config := c.runtime.Configuration()
    
    // Get operation configuration
    multiplier := config.GetIntDefault("operation.multiplier", 1)
    prefix := config.GetStringDefault("operation.prefix", "result")
    enabled := config.GetBoolDefault("operation.enabled", true)
    
    if !enabled {
        logger.Info("Operation disabled by configuration")
        return component.Output{
            Data: map[string]interface{}{
                "status": "disabled",
            },
        }, nil
    }
    
    // Process with configuration
    data := input.Data.(map[string]interface{})
    value, _ := data["value"].(float64)
    
    result := map[string]interface{}{
        prefix:      value * float64(multiplier),
        "config":    "applied",
        "timestamp": time.Now().Format(time.RFC3339),
    }
    
    logger.Info("Operation executed with configuration",
        "multiplier", multiplier,
        "prefix", prefix)
    
    return component.Output{Data: result}, nil
}
```

## 🔄 Framework Integration Patterns

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
2. **Use BaseService/BaseOperation** embedding
3. **Access framework services** through runtime reference
4. **Keep operations simple** - focus on input/output transformation
5. **Handle configuration** with defaults
6. **Use structured logging** with component context
7. **Validate input** in operations
8. **Manage service lifecycle** properly

### ❌ Don't

1. **Implement real infrastructure** - focus on framework patterns
2. **Mix business logic** with framework demonstration
3. **Access system directly** - always use runtime reference
4. **Ignore error handling** - but keep it framework-focused
5. **Hardcode values** - use configuration with defaults
6. **Complicate operations** - keep them simple and focused
7. **Skip initialization** - always store runtime reference

## 🔄 Component Lifecycle Summary

### Service Lifecycle
```
Create → Initialize (store runtime) → Start → Running → Stop → Dispose
```

### Operation Lifecycle
```
Create → Initialize (store runtime) → Execute → Dispose
```

## 📚 Next Steps

1. **Choose component type** - Service for long-running, Operation for stateless
2. **Embed base component** - BaseService or BaseOperation
3. **Store runtime reference** - Critical for framework access
4. **Implement lifecycle methods** - Initialize, Start/Execute, Stop
5. **Use framework services** - Logger, Configuration, EventBus
6. **Keep it simple** - Focus on framework patterns, not business logic

---

**Remember**: Examples demonstrate **framework usage patterns**, not business functionality. Keep your components simple and focused on showing how to use the framework correctly. 