# Builder API Examples

This directory contains examples demonstrating the new **Builder API** for the Fintechain Skeleton framework. The Builder API replaces FX dependency injection with a simple, explicit builder pattern that provides the same functionality without the complexity.

## 🎯 Key Advantages of Builder API

### ✅ **Simple and Explicit**
- No FX knowledge required
- Clear dependency flow
- Easy debugging and testing
- Compile-time errors for missing dependencies

### ✅ **Custom Dependency Support**
- Easy injection of custom configurations
- Custom logger implementations
- Custom event bus services
- No complex FX provider patterns

### ✅ **Backward Compatible**
- Existing FX-based code continues to work
- Gradual migration path
- Same runtime behavior

## 🚀 Running the Examples

### Basic Usage

```bash
# Daemon mode - long-running service
go run . daemon

# Command mode - execute operation and exit
go run . command

# Custom dependencies - demonstrates dependency injection
go run . custom
```

### Example Output

**Command Mode:**
```
=== Builder API Command Mode Example ===
[Fintechain] Loading 1 plugins...
[Fintechain] Executing command: test-operation
[Fintechain] Command completed successfully
Command result: map[operation_id:test-operation processed_message:Processed: Hello from Builder API! processed_number:84 status:success timestamp:2024-01-15T10:30:45Z]
```

**Custom Dependencies:**
```
=== Builder API Custom Dependencies Examples ===
This demonstrates the key advantage of the Builder API:
Easy injection of custom configurations, loggers, and other services

=== Custom Configuration Example ===
Result with custom config: map[...] 
App name from config: Custom Builder App
App version from config: 2.0.0

=== Custom Logger Example ===
Result with custom logger: map[...]
Custom logger was used for all framework logging

=== Full Custom Dependencies Example ===
Result with full customization: map[...]
Environment: production
Advanced features enabled: true
```

## 📁 Files Overview

### `main.go`
Main entry point that demonstrates three usage modes:
- **Daemon mode**: Long-running service with signal handling
- **Command mode**: Execute operation and exit
- **Custom mode**: Demonstrates custom dependency injection

### `test_plugin.go`
Simple test plugin that provides:
- **TestService**: Demonstrates service lifecycle management
- **TestOperation**: Demonstrates stateless operation execution
- Follows the plugin-as-orchestrator pattern

### `custom.go`
Advanced examples showing custom dependency injection:
- Custom configuration with application-specific settings
- Custom logger implementation
- Full customization with multiple dependencies

## 🔧 Builder API Usage Patterns

### Basic Usage (Zero Configuration)

```go
// Daemon mode with defaults
err := runtime.NewBuilder().
    WithPlugins(myPlugin).
    BuildDaemon()

// Command mode with defaults  
result, err := runtime.NewBuilder().
    WithPlugins(myPlugin).
    BuildCommand("my-operation", inputData)
```

### Custom Configuration

```go
// Create custom configuration
config := infraConfig.NewMemoryConfigurationWithData(map[string]interface{}{
    "app.name": "My Application",
    "app.port": 8080,
    "database.host": "localhost",
})

// Use with builder
err := runtime.NewBuilder().
    WithConfig(config).
    WithPlugins(myPlugin).
    BuildDaemon()
```

### Custom Logger

```go
// Create custom logger
logger := myCustomLogger()

// Use with builder
result, err := runtime.NewBuilder().
    WithLogger(logger).
    WithPlugins(myPlugin).
    BuildCommand("my-operation", inputData)
```

### Multiple Custom Dependencies

```go
// Full customization
err := runtime.NewBuilder().
    WithConfig(customConfig).
    WithLogger(customLogger).
    WithEventBus(customEventBus).
    WithPlugins(plugin1, plugin2).
    BuildDaemon()
```

## 🆚 Comparison with FX API

### FX-based API (Legacy)

```go
// Complex FX setup for custom dependencies
err := runtime.StartDaemon(
    runtime.WithOptions(
        fx.Replace(fx.Annotate(createCustomConfig, fx.As(new(config.Configuration)))),
        fx.Replace(fx.Annotate(createCustomLogger, fx.As(new(logging.LoggerService)))),
    ),
    runtime.WithPlugins(myPlugin),
)
```

**Problems:**
- Complex FX provider signatures
- Hard to debug dependency injection failures
- Requires FX knowledge
- Magic dependency resolution

### Builder API (Recommended)

```go
// Simple builder setup for custom dependencies
err := runtime.NewBuilder().
    WithConfig(createCustomConfig()).
    WithLogger(createCustomLogger()).
    WithPlugins(myPlugin).
    BuildDaemon()
```

**Benefits:**
- Explicit dependency injection
- Easy to understand and debug
- No FX knowledge required
- Clear error messages

## 🧪 Testing Custom Dependencies

The custom examples demonstrate how to verify that custom dependencies are working:

1. **Custom Configuration**: Verify that custom values are accessible
2. **Custom Logger**: Verify that custom logger receives framework messages
3. **Full Customization**: Verify that all custom dependencies work together

## 📚 Next Steps

1. **Try the examples**: Run all three modes to see the Builder API in action
2. **Create your own plugin**: Use the test plugin as a template
3. **Add custom dependencies**: Experiment with custom configurations and loggers
4. **Migrate existing code**: Gradually replace FX-based code with Builder API

## 🔗 Related Documentation

- [Builder API Implementation Plan](../../docs/BUILDER_API_IMPLEMENTATION_PLAN.md)
- [Plugin Development Guide](../../docs/PLUGIN_DEVELOPMENT_GUIDE.md)
- [Runtime Development Guide](../../docs/RUNTIME_DEVELOPMENT_GUIDE.md)

---

**The Builder API provides the same functionality as FX dependency injection with none of the complexity!** 