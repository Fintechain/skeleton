# Fintechain Skeleton Framework - Public API Guide

Welcome to the **Fintechain Skeleton Framework** - a modern, domain-driven application framework built for scalability, testability, and developer productivity.

## 🚀 Quick Start

Get up and running in under 5 minutes:

### Option 1: Modern Builder API (Recommended)

```go
package main

import (
    "github.com/fintechain/skeleton/pkg/runtime"
    "github.com/fintechain/skeleton/pkg/plugin"
)

func main() {
    // Start a long-running service
    err := runtime.NewBuilder().
        WithPlugins(myWebServerPlugin, myDatabasePlugin).
        BuildDaemon()
    if err != nil {
        panic(err)
    }
}
```

### Option 2: Traditional Options Pattern

```go
package main

import (
    "github.com/fintechain/skeleton/pkg/runtime"
    "github.com/fintechain/skeleton/pkg/context"
)

func main() {
    runtime, err := runtime.NewRuntimeWithOptions(
        runtime.WithPlugins(myWebServerPlugin, myDatabasePlugin),
    )
    if err != nil {
        panic(err)
    }
    
    ctx := context.NewContext()
    if err := runtime.Start(ctx); err != nil {
        panic(err)
    }
    
    // Note: No Wait() method - use signal handling for blocking
    select {} // Block indefinitely
}
```

## 📦 Package Overview

The public API is organized into focused packages that abstract away infrastructure complexity:

```
pkg/
├── runtime/     🔥 Modern Builder API and runtime management
├── component/   🏗️ Component system abstractions
├── context/     🔄 Application context handling
├── event/       📡 Event system integration
└── plugin/      🔌 Plugin system interfaces
```

### 🔥 `pkg/runtime` - Modern Builder API

**Use When**: Building new applications (recommended approach)

```go
// Daemon mode - long-running services
runtime.NewBuilder().WithPlugins(webServer, database).BuildDaemon()

// Command mode - execute and exit
result, err := runtime.NewBuilder().
    WithPlugins(calculatorPlugin).
    BuildCommand("calculate-total", inputData)
```

**Benefits**:
- ✅ Simple, explicit dependency injection
- ✅ Easy debugging and testing
- ✅ No complex framework knowledge required
- ✅ Built-in lifecycle management
- ✅ Custom dependency support

### ⚡ Traditional Options Pattern

**Use When**: Migrating existing code or need explicit control

```go
runtime, err := runtime.NewRuntimeWithOptions(
    runtime.WithConfiguration(config),
    runtime.WithPlugins(plugin1, plugin2),
)
```

**Benefits**:
- ✅ Explicit dependency management
- ✅ Step-by-step initialization control
- ✅ Familiar options pattern
- ✅ Easy to debug and understand

## 🎯 Core Concepts

### Runtime Modes

The framework supports two fundamental execution patterns:

#### 🔄 Daemon Mode
For **long-running applications** that provide continuous services:

```go
// Web servers, message processors, monitoring services
err := runtime.NewBuilder().
    WithPlugins(webServerPlugin, messageProcessorPlugin).
    BuildDaemon()
```

**Characteristics**:
- Starts all services and keeps them running
- Handles graceful shutdown on signals (SIGINT, SIGTERM)
- Ideal for servers, workers, and background processes

#### ⚡ Command Mode
For **short-lived applications** that execute specific tasks:

```go
// CLI tools, batch processors, one-time calculations
result, err := runtime.NewBuilder().
    WithPlugins(dataProcessorPlugin).
    BuildCommand("process-data", inputData)
```

**Characteristics**:
- Initializes dependencies but doesn't start long-running services
- Executes operation and exits immediately
- Ideal for CLI commands and batch jobs

### Component System

Everything in the framework is a **Component** with a unified lifecycle:

```go
type Component interface {
    // Identity methods (from Identifiable interface)
    ID() ComponentID
    Name() string
    Description() string
    Version() string
    
    // Component-specific methods
    Type() ComponentType
    Metadata() Metadata
    Initialize(ctx context.Context, system System) error
    Dispose() error
}
```

#### Component Types

| Type | Purpose | Examples |
|------|---------|----------|
| `Component` | Basic managed entities | Database connections, config loaders |
| `Operation` | Executable tasks with input/output | Calculations, data transformations |
| `Service` | Long-running background processes | Web servers, message queues |

### Plugin System

Plugins are the primary way to extend functionality:

```go
type Plugin interface {
    component.Service  // Extends Service with lifecycle methods
    
    Author() string
    PluginType() PluginType
}

// Example plugin implementation
type CalculatorPlugin struct {
    *component.BaseService
}

func (p *CalculatorPlugin) Author() string {
    return "My Team"
}

func (p *CalculatorPlugin) PluginType() plugin.PluginType {
    return plugin.TypeProcessor
}

func (p *CalculatorPlugin) Initialize(ctx context.Context, system component.System) error {
    // Register components with the system
    registry := system.Registry()
    return registry.Register(&AddOperation{})
}
```

## 🛠️ Usage Patterns

### Pattern 1: Web Service Application

```go
// main.go
func main() {
    webPlugin := &WebServerPlugin{Port: 8080}
    dbPlugin := &DatabasePlugin{URL: "postgres://..."}
    
    err := runtime.NewBuilder().
        WithPlugins(webPlugin, dbPlugin).
        BuildDaemon()
    if err != nil {
        log.Fatal(err)
    }
}
```

### Pattern 2: CLI Tool

```go
// cmd/calculator/main.go
func main() {
    if len(os.Args) < 4 {
        fmt.Println("Usage: calculator <operation> <a> <b>")
        os.Exit(1)
    }
    
    operation := os.Args[1]
    a, _ := strconv.ParseFloat(os.Args[2], 64)
    b, _ := strconv.ParseFloat(os.Args[3], 64)
    
    result, err := runtime.NewBuilder().
        WithPlugins(&CalculatorPlugin{}).
        BuildCommand(operation, map[string]interface{}{
            "a": a,
            "b": b,
        })
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Result: %v\n", result["result"])
}
```

### Pattern 3: Batch Processor

```go
func processFiles(inputDir, outputDir string) error {
    _, err := runtime.NewBuilder().
        WithPlugins(
            &FileProcessorPlugin{},
            &CompressionPlugin{},
            &ValidationPlugin{},
        ).
        BuildCommand("process-directory", map[string]interface{}{
            "inputDir":  inputDir,
            "outputDir": outputDir,
        })
    return err
}
```

### Pattern 4: Testing Components

```go
func TestCalculatorOperation(t *testing.T) {
    // Use the options pattern for precise test control
    runtime, err := runtime.NewRuntimeWithOptions(
        runtime.WithPlugins(&CalculatorPlugin{}),
    )
    require.NoError(t, err)
    
    ctx := context.NewContext()
    result, err := runtime.ExecuteOperation(ctx, "add", component.Input{
        Data: map[string]interface{}{"a": 5, "b": 3},
        Metadata: map[string]string{"test": "true"},
    })
    
    require.NoError(t, err)
    assert.Equal(t, 8.0, result.Data)
}
```

### Component Package (`pkg/component`)

```go
type ComponentID string
type ComponentType string

const (
    TypeComponent   = "component"
    TypeOperation   = "operation"
    TypeService     = "service"
    TypeBasic       = "basic"
    TypeSystem      = "system"
    TypeApplication = "application"
)

type Input struct {
    Data any
    Metadata map[string]string
}

type Output struct {
    Data any
}
```

## 🚀 Migration Guide

### From Manual Setup to Builder API

#### Before (Manual)
```go
// Lots of manual wiring
config := loadConfig()
logger := newLogger(config)
db := newDatabase(config, logger)
server := newWebServer(config, db, logger)

// Manual lifecycle management
db.Connect()
server.Start()
// ... handle shutdown manually
```

#### After (Builder API)
```go
// Clean, declarative setup
runtime.NewBuilder().
    WithPlugins(
        &ConfigPlugin{},
        &LoggerPlugin{},
        &DatabasePlugin{},
        &WebServerPlugin{},
    ).
    BuildDaemon()
```

### From Options Pattern to Builder API

#### Before (Options Pattern)
```go
runtime, err := runtime.NewRuntimeWithOptions(
    runtime.WithPlugins(&DatabasePlugin{}, &WebServerPlugin{}),
)

ctx := context.NewContext()
runtime.Start(ctx)
// Note: Manual signal handling needed for blocking
```

#### After (Builder API)
```go
runtime.NewBuilder().
    WithPlugins(&DatabasePlugin{}, &WebServerPlugin{}).
    BuildDaemon()
```

## 📚 Next Steps

- **New to the framework?** Check out the [examples](../examples/README.md) directory
- **Want to build plugins?** See the [Plugin Development Guide](../docs/PLUGIN_DEVELOPMENT_GUIDE.md)
- **Building components?** Review the [Service Operations Development Guide](../docs/SERVICE_OPERATIONS_DEVELOPMENT_GUIDE.md)
- **Need advanced patterns?** Review the [Runtime Development Guide](../docs/RUNTIME_DEVELOPMENT_GUIDE.md)

## 🤝 Best Practices

### ✅ Do

- **Use Builder API** for new applications
- **Keep plugins focused** - one responsibility per plugin
- **Handle errors gracefully** - always check for initialization errors
- **Use context correctly** - pass the framework's domain context
- **Test components** - use the builder pattern for precise test control

### ❌ Don't

- **Mix runtime modes** - don't start services in command mode
- **Ignore lifecycle** - always implement proper Initialize/Dispose
- **Hardcode dependencies** - use the registry for component discovery
- **Forget error handling** - especially during startup and shutdown
- **Block in Initialize** - keep initialization fast, defer heavy work to Start()

---

**Framework Version**: v1.0.0  
**Documentation Updated**: 2024  
**Need Help?** Check the examples or open an issue on GitHub