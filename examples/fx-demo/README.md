# FX Integration Demo

This example demonstrates the **modern FX integration** approach using Uber's FX framework for automatic dependency injection and lifecycle management.

## 🎯 What This Demonstrates

- **Modern Dependency Injection**: Using FX for automatic component wiring
- **Simplified Lifecycle**: Automatic startup/shutdown handling
- **Clean Plugin Design**: Focus on business logic, not infrastructure
- **Dual Mode Support**: Same plugin works in daemon and command modes

## 🚀 How to Run

### Daemon Mode
Starts the calculator plugin as a long-running service:

```bash
go run examples/fx-demo/main.go daemon
```

**What happens:**
- FX initializes all dependencies automatically
- Calculator plugin starts as a service
- Application runs until Ctrl+C
- Graceful shutdown with automatic cleanup

### Command Mode
Executes calculator operations and exits:

```bash
go run examples/fx-demo/main.go command
```

**What happens:**
- FX initializes dependencies
- Executes addition: 10 + 5 = 15
- Executes division: 20 / 4 = 5
- Automatic cleanup and exit

## 🏗️ Code Structure

```go
// Simple operation
type Calculator struct {
    *infraComponent.BaseOperation
}

// Simple plugin - just register components
type CalculatorPlugin struct {
    *infraComponent.BaseService
    calculator *Calculator
}

// FX handles everything else
fx.StartDaemon(fx.WithPlugins(calculatorPlugin))
```

## 💡 Key Benefits of FX Approach

### ✅ What FX Provides
- **Automatic dependency wiring** - No manual service discovery
- **Lifecycle management** - Automatic start/stop ordering  
- **Error handling** - Centralized error reporting
- **Signal handling** - Built-in graceful shutdown
- **Testing support** - Easy mocking and isolation

### 🔄 vs Traditional Approach
- **Less boilerplate** - No manual runtime creation
- **Type safety** - Compile-time dependency checking
- **Consistent patterns** - Standard FX patterns across team
- **Production ready** - Battle-tested in many companies

## 🎯 When to Use FX Approach

✅ **Perfect for:**
- New applications
- Production deployments
- Team development (consistent patterns)
- Complex dependency graphs

⚠️ **Consider Traditional for:**
- Learning the framework (more explicit)
- Migrating existing code
- Custom initialization order requirements

## 🔧 Extending This Example

```go
// Add more operations
type AdvancedCalculator struct {
    *infraComponent.BaseOperation
}

// Add to plugin
func (p *CalculatorPlugin) Initialize(ctx context.Context, system component.System) error {
    registry := system.Registry()
    registry.Register(p.calculator)
    registry.Register(&AdvancedCalculator{}) // Add more operations
    return nil
}

// Use in FX
fx.StartDaemon(fx.WithPlugins(calculatorPlugin))
```

## 📚 Next Steps

1. **Compare with traditional**: Run `examples/traditional-runtime/main.go`
2. **Try multi-plugin**: Run `examples/complete-app/main.go`
3. **Read FX guide**: Check `pkg/fx/README.md`
4. **Build your plugin**: Copy this as a template 