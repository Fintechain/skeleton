# Fintechain Skeleton Framework Examples

This directory contains practical examples demonstrating how to build applications with the Fintechain Skeleton framework using the modern **Builder API**.

## 🚀 Quick Start

```bash
# Run the complete application example
go run examples/complete-app/main.go command
```

## 📁 What's Included

### `builder-app/` - Builder API Examples

Focused examples demonstrating the new Builder API for simple, explicit dependency injection.

**Features:**
- Builder API patterns
- Custom dependency injection
- Daemon and command modes
- Simple, debuggable code

**Usage:**
```bash
# Daemon mode - long-running services
go run examples/builder-app/main.go daemon

# Command mode - execute operation and exit
go run examples/builder-app/main.go command

# Custom dependencies mode
go run examples/builder-app/main.go custom
```

### `complete-app/` - Complete Application Example

A full application demonstrating multiple plugins working together with the framework.

**Features:**
- Multiple plugins (webserver + database)
- Custom service providers
- Both daemon and command modes
- Plugin lifecycle management

**Usage:**
```bash
# Command mode - execute operation and exit
go run examples/complete-app/main.go command

# Daemon mode - long-running services
go run examples/complete-app/main.go daemon

# Custom providers mode
go run examples/complete-app/main.go custom
```

### Plugin Examples

Located in `complete-app/plugins/`:

- **`webserver/`** - HTTP service and routing operation
- **`database/`** - Database connection service and query operation

## 🎯 Key Patterns Demonstrated

### 1. Plugin Structure
```go
type MyPlugin struct {
    *infraComponent.BaseService
    runtime runtime.RuntimeEnvironment
}

func (p *MyPlugin) Initialize(ctx context.Context, system component.System) error {
    p.runtime = system.(runtime.RuntimeEnvironment)
    // Initialize and register components
    return nil
}
```

### 2. Service Components
```go
type MyService struct {
    *infraComponent.BaseService
    runtime runtime.RuntimeEnvironment
}

func (s *MyService) Start(ctx context.Context) error {
    logger := s.runtime.Logger()
    config := s.runtime.Configuration()
    // Service logic here
    return nil
}
```

### 3. Operation Components
```go
type MyOperation struct {
    *infraComponent.BaseOperation
    runtime runtime.RuntimeEnvironment
}

func (o *MyOperation) Execute(ctx context.Context, input component.Input) (component.Output, error) {
    // Process input and return output
    return component.Output{Data: result}, nil
}
```

## 📚 Runtime API

### Daemon Mode (Long-running services)
```go
err := runtime.NewBuilder().
    WithPlugins(myPlugin1, myPlugin2).
    BuildDaemon()
```

### Command Mode (Execute and exit)
```go
result, err := runtime.NewBuilder().
    WithPlugins(myPlugin1, myPlugin2).
    BuildCommand("operation-id", inputData)
```

### Custom Dependencies
```go
err := runtime.NewBuilder().
    WithPlugins(myPlugin).
    WithConfig(customConfig).
    WithLogger(customLogger).
    WithEventBus(customEventBus).
    BuildDaemon()
```

## 🎯 Best Practices

- **Use Builder API** for new applications
- Store `runtime.RuntimeEnvironment` reference in components
- Use `runtime.Logger()`, `runtime.Configuration()` for framework services
- Keep operations simple and focused
- Let plugins manage service lifecycle
- Focus on framework patterns, not business logic

## 📖 Next Steps

- Start with `builder-app/` to learn the Builder API
- Explore `complete-app/` to understand complex plugin interactions
- Read the [Plugin Development Guide](../docs/PLUGIN_DEVELOPMENT_GUIDE.md)
- Check the [Runtime Development Guide](../docs/RUNTIME_DEVELOPMENT_GUIDE.md)

---

**Framework Version**: v1.0.0  
**Documentation Updated**: 2024  
**Focus**: Builder API with simple, explicit dependency injection 