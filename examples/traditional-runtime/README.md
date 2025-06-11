# Traditional Runtime Demo

This example demonstrates the **traditional options pattern** approach using explicit runtime builder and manual lifecycle management.

## 🎯 What This Demonstrates

- **Explicit Control**: Manual configuration and lifecycle management
- **Options Pattern**: Builder-style runtime configuration
- **Framework Fundamentals**: Core concepts without abstractions
- **Migration Path**: How to gradually move to modern approaches

## 🚀 How to Run

### Daemon Mode
Starts the calculator plugin with manual runtime management:

```bash
go run examples/traditional-runtime/main.go daemon
```

**What happens:**
- Manual runtime creation with options
- Explicit Start() call with context
- Plugin starts as a service
- Manual Stop() call for cleanup

### Command Mode
Executes calculator operations with explicit runtime control:

```bash
go run examples/traditional-runtime/main.go command
```

**What happens:**
- Manual runtime creation and startup
- Explicit ExecuteOperation() calls
- Operations: 15 + 25 = 40, 6 * 7 = 42
- Manual runtime shutdown

## 🏗️ Code Structure

```go
// Same plugin implementation as FX
type SimpleCalculatorPlugin struct {
    *infraComponent.BaseService
    calculator *SimpleCalculator
}

// But explicit runtime management
runtime, err := runtime.NewRuntimeWithOptions(
    runtime.WithPlugins(NewSimpleCalculatorPlugin()),
)

ctx := context.NewContext()
runtime.Start(ctx)
defer runtime.Stop(ctx)

// Manual operation execution
result, err := runtime.ExecuteOperation(ctx, "simple-calculator", input)
```

## 💡 Key Characteristics

### ✅ What Traditional Provides
- **Explicit control** - You manage every step
- **Debugging clarity** - Easy to see what's happening
- **Learning value** - Understand framework internals
- **Migration support** - Gradual adoption path

### 🔄 vs FX Approach
- **More boilerplate** - Manual setup and teardown
- **Explicit errors** - Handle each step's errors
- **Custom patterns** - You control initialization order
- **Framework agnostic** - Not tied to FX patterns

## 🎯 When to Use Traditional Approach

✅ **Perfect for:**
- Learning the framework fundamentals
- Migrating existing applications
- Custom initialization requirements
- Debugging dependency issues
- Testing with precise control

⚠️ **Consider FX for:**
- New applications
- Production deployments
- Team development
- Standard patterns

## 🔧 Advanced Usage

```go
// Custom registry
mockRegistry := mocks.NewRegistry()

// Custom configuration
config := customConfig.Load()

// Precise runtime setup
runtime, err := runtime.NewRuntimeWithOptions(
    runtime.WithRegistry(mockRegistry),
    runtime.WithConfiguration(config),
    runtime.WithPlugins(plugin1, plugin2),
)

// Custom context
ctx := context.NewContextWithTimeout(30 * time.Second)

// Explicit lifecycle
if err := runtime.Start(ctx); err != nil {
    // Handle startup errors
}

// Use runtime...

if err := runtime.Stop(ctx); err != nil {
    // Handle shutdown errors
}
```

## 🧪 Testing Benefits

Traditional approach excels in testing scenarios:

```go
func TestWithPreciseControl(t *testing.T) {
    // Mock specific dependencies
    mockRegistry := mocks.NewRegistry()
    mockConfig := mocks.NewConfiguration()
    
    runtime, err := runtime.NewRuntimeWithOptions(
        runtime.WithRegistry(mockRegistry),
        runtime.WithConfiguration(mockConfig),
        runtime.WithPlugins(&TestPlugin{}),
    )
    
    // Test specific initialization steps
    ctx := context.NewContext()
    err = runtime.Start(ctx)
    require.NoError(t, err)
    
    // Test specific operations
    result, err := runtime.ExecuteOperation(ctx, "test-op", input)
    
    // Verify specific behaviors
    assert.True(t, mockRegistry.RegisterCalled)
}
```

## 📚 Next Steps

1. **Compare with FX**: Run `examples/fx-demo/main.go`
2. **See multi-plugin**: Run `examples/complete-app/main.go`
3. **Read runtime guide**: Check `pkg/runtime/README.md`
4. **Migration path**: Learn how to move to FX gradually 