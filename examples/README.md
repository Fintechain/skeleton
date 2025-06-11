# Fintechain Skeleton Framework Examples

This directory demonstrates both **traditional** and **modern** approaches to building applications with the Fintechain Skeleton framework, focusing on **framework patterns** rather than business logic complexity.

## 📁 Example Structure

### Runnable Examples (Each in own package)

#### `fx-demo/` - Modern FX Integration
The recommended approach using Uber's FX framework for dependency injection.

**Best for**:
- New applications
- Production deployments
- When you want automatic dependency management
- Learning modern framework patterns

**Run it**:
```bash
# Daemon mode (long-running service)
go run examples/fx-demo/main.go daemon

# Command mode (execute and exit)
go run examples/fx-demo/main.go command
```

**What it demonstrates**:
- Simple calculator operation
- Basic plugin structure
- FX integration patterns
- Component lifecycle management

#### `traditional-runtime/` - Traditional Builder Pattern
The original approach using explicit builder pattern and manual lifecycle management.

**Best for**:
- Learning the framework fundamentals
- Migrating existing applications 
- When you need explicit control over initialization order
- Understanding component patterns

**Run it**:
```bash
# Daemon mode
go run examples/traditional-runtime/main.go daemon

# Command mode  
go run examples/traditional-runtime/main.go command
```

**What it demonstrates**:
- Traditional runtime creation
- Manual lifecycle management
- Component registration patterns
- Options pattern usage

#### `complete-app/` - Multi-Plugin Application
Demonstrates multiple plugins working together with simplified operations.

**Best for**:
- Understanding plugin coordination
- Learning multi-plugin patterns
- Seeing framework orchestration

**Run it**:
```bash
# Start multiple services in daemon mode
go run examples/complete-app/main.go daemon

# Execute operation and exit in command mode
go run examples/complete-app/main.go command
```

**What it demonstrates**:
- Multiple plugins coordination
- Service lifecycle in daemon mode
- Operation execution in command mode
- Framework patterns without business complexity

### Plugin Libraries

#### `plugins/webserver/` - Web Server Plugin
Simplified HTTP server plugin demonstrating framework patterns.
- **HTTPService**: Service lifecycle with framework integration
- **RouteOperation**: Simple input/output operation processing
- **WebServerPlugin**: Plugin orchestration and component management

**Key patterns shown**:
- Runtime reference storage in components
- Framework services access (logger, config)
- Plugin-managed service lifecycle
- Component registration with registry

#### `plugins/database/` - Database Plugin  
Simplified database plugin demonstrating framework patterns.
- **DatabaseConnectionService**: Service with simulated connection
- **QueryOperation**: Simple query processing operation
- **DatabasePlugin**: Plugin orchestration patterns

**Key patterns shown**:
- Component initialization by plugin
- Service start/stop management
- Operation input/output processing
- Framework services integration

### Documentation

#### `PLUGIN_DEVELOPMENT_GUIDE.md` - Framework Patterns Guide
Comprehensive guide covering:
- Plugin-as-orchestrator architecture
- Component lifecycle patterns
- Framework service access
- Runtime reference storage
- Testing approaches
- Best practices

## 🚀 Quick Start

### Try the Examples

```bash
# 1. Start with FX demo (recommended)
go run examples/fx-demo/main.go command

# 2. Try the complete application
go run examples/complete-app/main.go command

# 3. Compare with traditional runtime
go run examples/traditional-runtime/main.go command
```

### Build Your First Plugin

```bash
# 1. Create your plugin directory
mkdir -p myplugin

# 2. Copy the webserver plugin as a template
cp -r examples/plugins/webserver/* myplugin/

# 3. Modify for your needs (focus on framework patterns)
# 4. Test it with FX
```

## 🔄 Comparison: Traditional vs FX

### Runtime Creation

#### Traditional Options Pattern
```go
// Explicit options configuration
runtime, err := runtime.NewRuntimeWithOptions(
    runtime.WithPlugins(plugin1, plugin2),
)

// Manual lifecycle management
ctx := context.NewContext()
if err := runtime.Start(ctx); err != nil {
    log.Fatal(err)
}
defer runtime.Stop(ctx)
```

#### Modern FX Integration
```go
// Declarative configuration
err := fx.StartDaemon(
    fx.WithPlugins(plugin1, plugin2),
)
// Lifecycle handled automatically
```

### Component Implementation

Both approaches use the same component patterns:

#### Framework Pattern (All Examples)
```go
type MyComponent struct {
    *component.BaseComponent
    runtime runtime.RuntimeEnvironment // Store runtime reference
}

func (c *MyComponent) Initialize(ctx context.Context, system component.System) error {
    // Store runtime reference for framework services access
    c.runtime = system.(runtime.RuntimeEnvironment)
    
    // Access framework services
    logger := c.runtime.Logger()
    config := c.runtime.Configuration()
    
    return nil
}
```

## 🎯 Framework Patterns Demonstrated

### Pattern 1: Component Lifecycle

```go
// Plugin orchestrates component lifecycle
func (p *MyPlugin) Initialize(ctx context.Context, system component.System) error {
    // 1. Initialize components
    p.myService.Initialize(ctx, system)
    p.myOperation.Initialize(ctx, system)
    
    // 2. Register with registry
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

### Pattern 2: Framework Services Access

```go
func (c *MyComponent) DoWork() {
    // Access framework services through stored runtime reference
    logger := c.runtime.Logger()
    config := c.runtime.Configuration()
    eventBus := c.runtime.EventBus()
    
    logger.Info("Component doing work")
}
```

### Pattern 3: Runtime Modes

#### Daemon Flow
```
StartDaemon → Plugin.Initialize → Plugin.Start → Services Running → Signal → Plugin.Stop
```

#### Command Flow
```
ExecuteCommand → Plugin.Initialize → Operation.Execute → Cleanup → Exit
```

## 🎯 Key Simplifications

### What Examples Focus On:
- ✅ **Framework Patterns**: Component lifecycle, plugin orchestration
- ✅ **Runtime Integration**: Proper framework services access
- ✅ **Component Types**: Service, Operation, Component distinctions
- ✅ **Plugin Coordination**: Multiple plugins working together
- ✅ **Lifecycle Management**: Initialize → Start → Stop → Dispose

### What Examples Avoid:
- ❌ **Business Logic**: Real HTTP servers, database connections
- ❌ **Complex Operations**: JSON parsing, SQL validation
- ❌ **Infrastructure Code**: Actual networking, file I/O
- ❌ **Error Complexity**: Business-specific error handling
- ❌ **External Dependencies**: Real databases, web frameworks

## 📚 Next Steps

- **New to the framework?** Start with `fx-demo/main.go`
- **Want to understand plugins?** Read `PLUGIN_DEVELOPMENT_GUIDE.md`
- **Building your first plugin?** Copy and modify `plugins/webserver/`
- **Need advanced patterns?** Review the [Domain Architecture Guide](../internal/domain/README.md)

## 🤝 Best Practices from Examples

### ✅ Do

- **Store runtime reference** in all components
- **Use plugin-as-orchestrator** pattern
- **Access framework services** through runtime reference
- **Keep operations simple** - focus on input/output transformation
- **Let plugins manage** service lifecycle

### ❌ Don't

- **Implement real infrastructure** in examples
- **Mix business logic** with framework patterns
- **Access system directly** - use runtime reference
- **Start services manually** - let plugins manage lifecycle
- **Complicate operations** - keep them simple and focused

---

**Framework Version**: v1.0.0  
**Documentation Updated**: 2024  
**Focus**: Framework patterns over business logic 