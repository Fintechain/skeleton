# Fintechain Skeleton Framework Examples

This directory contains practical examples demonstrating how to build applications with the Fintechain Skeleton framework.

## 🚀 Quick Start

```bash
# Run the complete application example
go run examples/complete-app/main.go command
```

## 📁 What's Included

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
err := runtime.StartDaemon(
    runtime.WithPlugins(myPlugin1, myPlugin2),
)
```

### Command Mode (Execute and exit)
```go
result, err := runtime.ExecuteCommand("operation-id", inputData,
    runtime.WithPlugins(myPlugin1, myPlugin2),
)
```

## 🎯 Best Practices

- Store `runtime.RuntimeEnvironment` reference in components
- Use `runtime.Logger()`, `runtime.Configuration()` for framework services
- Keep operations simple and focused
- Let plugins manage service lifecycle
- Focus on framework patterns, not business logic

## 📖 Next Steps

- Explore `complete-app/` to understand the framework
- Read the [Plugin Development Guide](../docs/development/PLUGIN_DEVELOPMENT_GUIDE.md)
- Check the [Domain Architecture Guide](../internal/domain/README.md)

---

**Framework Version**: v1.0.0  
**Documentation Updated**: 2024  
**Focus**: FX-based runtime with framework patterns 