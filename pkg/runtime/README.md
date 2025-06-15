# Runtime Package

Modern dependency injection for the Fintechain Skeleton Framework using **Uber's FX** framework.

## 🎯 Why FX?

FX provides **automatic dependency injection** and **lifecycle management**, eliminating manual service wiring and reducing boilerplate code.

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

### After (Automatic with FX)
```go
// Automatic dependency injection and lifecycle management
import "github.com/fintechain/skeleton/pkg/runtime"

runtime.StartDaemon(
    runtime.WithPlugins(webPlugin, dbPlugin),
)
```

## 🏗️ How FX Works

FX creates and manages the framework runtime automatically:

```
Your Code:
runtime.StartDaemon(runtime.WithPlugins(myPlugin))
    ↓
FX CoreModule:
├── Configuration (memory-based default)
├── Logger (structured logging)
├── EventBus (pub/sub events)
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
    err := runtime.StartDaemon(
        runtime.WithPlugins(
            webserver.NewWebServerPlugin(8080),
            database.NewDatabasePlugin("postgres", "connection-string"),
        ),
    )
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
    result, err := runtime.ExecuteCommand("database-query", 
        map[string]interface{}{
            "query": "SELECT * FROM users WHERE active = true",
            "limit": 100,
        },
        runtime.WithPlugins(
            database.NewDatabasePlugin("postgres", "connection-string"),
        ),
    )
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
runtime.StartDaemon(runtime.WithPlugins(myPlugin))
```

### Custom Configuration

You can provide your own configuration implementation using FX options:

```go
import (
    "go.uber.org/fx"
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

runtime.StartDaemon(
    runtime.WithPlugins(myPlugin),
    runtime.WithOptions(
        fx.Replace(fx.Annotate(createCustomConfig, fx.As(new(config.Configuration)))),
    ),
)
```

## 🔧 Advanced Usage

### WithOptions - Escape Hatch

Use `WithOptions` to customize the FX dependency injection container:

```go
import (
    "go.uber.org/fx"
    "github.com/fintechain/skeleton/pkg/runtime"
)

runtime.StartDaemon(
    runtime.WithPlugins(myPlugin),
    runtime.WithOptions(
        // Override framework services
        fx.Replace(fx.Annotate(myCustomLogger, fx.As(new(logging.LoggerService)))),
        
        // Add custom providers
        fx.Provide(func() *MyExternalService {
            return &MyExternalService{APIKey: os.Getenv("API_KEY")}
        }),
        
        // Add lifecycle hooks
        fx.Invoke(func(lc fx.Lifecycle, service *MyExternalService) {
            lc.Append(fx.Hook{
                OnStart: func(ctx context.Context) error {
                    return service.Connect()
                },
                OnStop: func(ctx context.Context) error {
                    return service.Disconnect()
                },
            })
        }),
    ),
)
```

### Custom Signal Handling

```go
import (
    "os"
    "syscall"
    "github.com/fintechain/skeleton/pkg/runtime"
)

runtime.StartDaemonWithSignalHandling(
    []os.Signal{syscall.SIGTERM, syscall.SIGUSR1},
    runtime.WithPlugins(myPlugin),
)
```

## 🔀 Custom Providers

Replace any framework service with your custom implementation using FX dependency injection.

### Basic Pattern

```go
runtime.StartDaemon(
    runtime.WithPlugins(myPlugin),
    runtime.WithOptions(
        // Replace framework services
        fx.Replace(fx.Annotate(myCustomLogger, fx.As(new(logging.LoggerService)))),
        fx.Replace(fx.Annotate(myCustomConfig, fx.As(new(config.Configuration)))),
        fx.Replace(fx.Annotate(myCustomEventBus, fx.As(new(event.EventBusService)))),
        
        // Add external services
        fx.Provide(createDatabaseConnection),
        fx.Provide(createAPIClient),
    ),
)
```

### Replaceable Services

- **Logger** - `logging.LoggerService` interface
- **Configuration** - `config.Configuration` interface  
- **EventBus** - `event.EventBusService` interface
- **Registry** - `component.Registry` interface
- **PluginManager** - `plugin.PluginManager` interface

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
    result, err := runtime.ExecuteCommand("my-operation", 
        map[string]interface{}{
            "name":  "test",
            "value": 10.0,
        },
        runtime.WithPlugins(NewMyPlugin()),
        runtime.WithOptions(
            fx.Replace(fx.Annotate(func() config.Configuration { return testConfig }, fx.As(new(config.Configuration)))),
        ),
    )
    
    require.NoError(t, err)
    assert.Equal(t, 20.0, result["processed_value"])
}
```

### Simple Testing

```go
func TestMyOperationSimple(t *testing.T) {
    result, err := runtime.ExecuteCommand("my-operation", map[string]interface{}{
        "name":  "test",
        "value": 10.0,
    }, runtime.WithPlugins(NewMyPlugin()))
    
    require.NoError(t, err)
    assert.Equal(t, 20.0, result["processed_value"])
}
```

## 📈 Best Practices

### ✅ Do

1. **Store runtime reference** in all components
2. **Use BaseService/BaseOperation** embedding
3. **Access framework services** through runtime reference  
4. **Keep operations simple** - focus on input/output transformation
5. **Use memory configuration** for testing with custom data
6. **Handle configuration** with defaults

### ❌ Don't

1. **Access system directly** - always use runtime reference
2. **Skip runtime reference storage** - it's required for framework services
3. **Complicate operations** - keep them simple and focused
4. **Import wrong packages** - use `pkg/runtime`, not `fx`
