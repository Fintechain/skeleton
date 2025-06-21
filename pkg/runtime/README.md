# Runtime Package

Modern, simple runtime management for the Fintechain Skeleton Framework using the **Builder API**.

## 🎯 Why Builder API?

The Builder API provides **simple, explicit dependency injection** and **lifecycle management**, eliminating complex framework knowledge and reducing debugging complexity.

### Before (Manual Setup)
```go
// Manual dependency creation and lifecycle management
registry := component.NewRegistry()
config := config.NewMemoryConfiguration()
logger := logging.NewNoOpLogger()
eventBus := event.NewEventBus()
pluginManager := plugin.NewManager()

runtime, err := runtime.NewRuntime(registry, config, pluginManager, eventBus, logger)
if err != nil {
    return err
}

ctx := context.NewContext()
if err := runtime.Start(ctx); err != nil {
    return err
}
defer runtime.Stop(ctx)

// Manual signal handling for daemon mode
// Manual error handling for each step
```

### After (Builder API)
```go
// Simple, explicit dependency injection and lifecycle management
import "github.com/fintechain/skeleton/pkg/runtime"

runtime.NewBuilder().
    WithPlugins(webPlugin, dbPlugin).
    BuildDaemon()
```

## 🏗️ How Builder API Works

The Builder API creates and manages the framework runtime with explicit control:

```
Your Code:
runtime.NewBuilder().WithPlugins(myPlugin).BuildDaemon()
    ↓
Builder Creates:
├── Configuration (memory-based default or custom)
├── Logger (structured logging or custom)
├── EventBus (pub/sub events or custom)
├── Registry (component management)
├── PluginManager (plugin lifecycle)
└── Runtime (orchestrates everything)
    ↓
Plugin Loading:
└── Calls plugin.Initialize() to register components
    ↓
Lifecycle Management:
└── Automatic Start/Stop with graceful shutdown
```

## 🚀 Runtime Modes

### Daemon Mode - Long-Running Services

**When to Use**: Web servers, background processors, message queues

```go
import "github.com/fintechain/skeleton/pkg/runtime"

func main() {
    err := runtime.NewBuilder().
        WithPlugins(
            webserver.NewWebServerPlugin(8080),
            database.NewDatabasePlugin("postgres", "connection-string"),
        ).
        BuildDaemon()
    if err != nil {
        log.Fatal("Failed to start daemon:", err)
    }
}
```

**What Happens**:
1. **Initialize**: All plugins get `Initialize()` called, register components
2. **Start**: Services get `Start()` called, begin accepting requests
3. **Run**: Application blocks, handling requests until shutdown signal
4. **Stop**: Graceful shutdown on SIGINT/SIGTERM, `Stop()` called in reverse order

### Command Mode - Execute and Exit

**When to Use**: CLI tools, batch processing, data transformations

```go
import "github.com/fintechain/skeleton/pkg/runtime"

func main() {
    result, err := runtime.NewBuilder().
        WithPlugins(
            database.NewDatabasePlugin("postgres", "connection-string"),
        ).
        BuildCommand("database-query", map[string]interface{}{
            "query": "SELECT * FROM users WHERE active = true",
            "limit": 100,
        })
    if err != nil {
        log.Fatal("Command failed:", err)
    }
    
    fmt.Printf("Result: %+v\n", result)
}
```

**What Happens**:
1. **Initialize**: Plugins register components, but services are **NOT** started
2. **Execute**: Specific operation is executed with input data
3. **Exit**: Immediate cleanup and exit

**Key Difference**: Command mode skips service startup for faster execution.

## 🔧 Configuration

### Default Configuration (Zero Setup)

The runtime automatically provides a default memory-based configuration:

```go
// This works with zero configuration
runtime.NewBuilder().WithPlugins(myPlugin).BuildDaemon()
```

### Custom Configuration

You can provide your own configuration implementation using the Builder API:

```go
import (
    "github.com/fintechain/skeleton/pkg/runtime"
    "github.com/fintechain/skeleton/internal/domain/config"
    infraConfig "github.com/fintechain/skeleton/internal/infrastructure/config"
)

func createCustomConfig() config.Configuration {
    // Create a memory configuration with custom data
    data := map[string]interface{}{
        "app.name": "MyApp",
        "app.port": 8080,
        "database.host": "localhost",
        "database.port": 5432,
    }
    return infraConfig.NewMemoryConfigurationWithData(data)
}

runtime.NewBuilder().
    WithPlugins(myPlugin).
    WithConfig(createCustomConfig()).
    BuildDaemon()
```

## 🔧 Advanced Usage

### Custom Dependencies

Use the Builder API to inject custom dependencies:

```go
import (
    "github.com/fintechain/skeleton/pkg/runtime"
    "github.com/fintechain/skeleton/internal/domain/config"
    "github.com/fintechain/skeleton/internal/domain/logging"
    "github.com/fintechain/skeleton/internal/domain/event"
    infraConfig "github.com/fintechain/skeleton/internal/infrastructure/config"
    infraLogging "github.com/fintechain/skeleton/internal/infrastructure/logging"
    infraEvent "github.com/fintechain/skeleton/internal/infrastructure/event"
)

runtime.NewBuilder().
    WithPlugins(myPlugin).
    // Custom configuration
    WithConfig(infraConfig.NewMemoryConfigurationWithData(map[string]interface{}{
        "app.name": "Custom App",
        "log.level": "debug",
    })).
    // Custom logger
    WithLogger(infraLogging.NewConsoleLogger()).
    // Custom event bus
    WithEventBus(infraEvent.NewInMemoryEventBus()).
    BuildDaemon()
```

## 🔀 Custom Dependencies

Replace any framework service with your custom implementation using the Builder API.

### Basic Pattern

```go
// Create custom dependencies
customConfig := infraConfig.NewMemoryConfigurationWithData(map[string]interface{}{
    "app.name": "Custom App",
})
customLogger := infraLogging.NewConsoleLogger()
customEventBus := infraEvent.NewInMemoryEventBus()

runtime.NewBuilder().
    WithPlugins(myPlugin).
    // Replace framework services
    WithConfig(customConfig).
    WithLogger(customLogger).
    WithEventBus(customEventBus).
    BuildDaemon()
```

### Replaceable Services

- **Logger** - `logging.LoggerService` interface
- **Configuration** - `config.Configuration` interface  
- **EventBus** - `event.EventBusService` interface

## 🧪 Testing

### Testing with Mock Configuration

For testing, create mock configurations and use the runtime package:

```go
import (
    "testing"
    "github.com/fintechain/skeleton/pkg/runtime"
    infraConfig "github.com/fintechain/skeleton/internal/infrastructure/config"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestMyOperation(t *testing.T) {
    // Create test configuration
    testConfig := infraConfig.NewMemoryConfigurationWithData(map[string]interface{}{
        "test.mode": true,
        "test.value": 42,
    })
    
    // Test operation execution
    result, err := runtime.NewBuilder().
        WithPlugins(NewMyPlugin()).
        WithConfig(testConfig).
        BuildCommand("my-operation", map[string]interface{}{
            "name":  "test",
            "value": 10.0,
        })
    
    require.NoError(t, err)
    assert.Equal(t, 20.0, result["processed_value"])
}
```

### Simple Testing

```go
func TestMyOperationSimple(t *testing.T) {
    result, err := runtime.NewBuilder().
        WithPlugins(NewMyPlugin()).
        BuildCommand("my-operation", map[string]interface{}{
            "name":  "test",
            "value": 10.0,
        })
    
    require.NoError(t, err)
    assert.Equal(t, 20.0, result["processed_value"])
}
```

## 📈 Best Practices

### ✅ Do

1. **Use Builder API** for new applications
2. **Store runtime reference** in all components
3. **Use BaseService/BaseOperation** embedding
4. **Access framework services** through runtime reference  
5. **Keep operations simple** - focus on input/output transformation
6. **Use memory configuration** for testing with custom data
7. **Handle configuration** with defaults

### ❌ Don't

1. **Use legacy API** for new projects - prefer Builder API
2. **Access system directly** - always use runtime reference
3. **Skip runtime reference storage** - it's required for framework services
4. **Complicate operations** - keep them simple and focused
5. **Mix Builder and legacy APIs** - choose one approach

---

**Framework Version**: v1.0.0  
**Documentation Updated**: 2024  
**Need Help?** Check the examples or open an issue on GitHub
