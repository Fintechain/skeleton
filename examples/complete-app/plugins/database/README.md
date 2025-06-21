# Database Plugin

A simplified database plugin demonstrating **framework patterns** for the Fintechain Skeleton framework.

## 🎯 Purpose

This plugin demonstrates **framework usage patterns**, not real database functionality:

- ✅ **Service lifecycle management**
- ✅ **Operation input/output processing**
- ✅ **Plugin orchestration patterns**
- ✅ **Runtime reference storage**
- ✅ **Framework services integration**

## 🏗️ Architecture

### Components

#### `DatabaseConnectionService` - Service Component
Demonstrates service lifecycle with framework integration.

**Key Patterns**:
- Stores runtime reference for framework services access
- Simulates database connection (no real database)
- Uses configuration with defaults
- Structured logging with context

```go
type DatabaseConnectionService struct {
    *component.BaseService
    runtime    runtime.RuntimeEnvironment // Store runtime reference
    driverName string
    dataSource string
    connected  bool
}
```

#### `QueryOperation` - Operation Component
Demonstrates simple query processing operation.

**Key Patterns**:
- Stores runtime reference for framework services access
- Simple query processing (no real SQL)
- Structured output format
- Framework logging integration

```go
type QueryOperation struct {
    *component.BaseOperation
    runtime runtime.RuntimeEnvironment // Store runtime reference
}
```

#### `DatabasePlugin` - Plugin Orchestrator
Demonstrates plugin-as-orchestrator pattern.

**Key Patterns**:
- Initializes and registers components
- Manages service lifecycle
- Stores runtime reference
- Component coordination

## 🚀 Usage Examples

### Daemon Mode (Service Lifecycle)

```go
func main() {
    err := runtime.NewBuilder().
        WithPlugins(database.NewDatabasePlugin("postgres", "test://connection")).
        BuildDaemon()
    if err != nil {
        log.Fatal(err)
    }
    // DatabaseConnectionService starts and runs until shutdown signal
}
```

**What happens**:
1. Plugin initializes and registers DatabaseConnectionService and QueryOperation
2. DatabaseConnectionService.Start() is called (simulates connection)
3. Service runs until SIGINT/SIGTERM
4. DatabaseConnectionService.Stop() is called (simulates cleanup)

### Command Mode (Operation Execution)

```go
func main() {
    result, err := runtime.NewBuilder().
        WithPlugins(database.NewDatabasePlugin("postgres", "test://connection")).
        BuildCommand("database-query", map[string]interface{}{
            "type": "validate",
            "sql":  "SELECT id, name FROM users",
        })
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Result: %+v\n", result)
}
```

**What happens**:
1. Plugin initializes and registers components (DatabaseConnectionService NOT started)
2. QueryOperation.Execute() is called with input data
3. Simple query processing returns structured output
4. Application exits immediately

## 🔧 Configuration

The plugin demonstrates configuration access patterns:

```json
{
  "database": {
    "driver": "postgres",
    "datasource": "postgres://localhost/myapp",
    "max_connections": 10
  }
}
```

**Configuration Access Pattern**:
```go
func (d *DatabaseConnectionService) Start(ctx context.Context) error {
    config := d.runtime.Configuration()
    
    // Type-safe access with defaults
    driverName := config.GetStringDefault("database.driver", d.driverName)
    maxConns := config.GetIntDefault("database.max_connections", 10)
    
    // Use configuration values...
}
```

## 📊 Framework Patterns Demonstrated

### 1. Runtime Reference Storage
```go
func (d *DatabaseConnectionService) Initialize(ctx context.Context, system component.System) error {
    // Store runtime reference - critical pattern
    d.runtime = system.(runtime.RuntimeEnvironment)
    
    // Access framework services
    logger := d.runtime.Logger()
    logger.Info("Database Connection Service initialized", "component_id", d.ID())
    
    return nil
}
```

### 2. Plugin Orchestration
```go
func (d *DatabasePlugin) Initialize(ctx context.Context, system component.System) error {
    // 1. Initialize components
    d.connectionService.Initialize(ctx, system)
    d.queryOperation.Initialize(ctx, system)
    
    // 2. Register with registry
    registry := system.Registry()
    registry.Register(d.connectionService)
    registry.Register(d.queryOperation)
    
    return nil
}
```

### 3. Service Lifecycle Management
```go
func (d *DatabasePlugin) Start(ctx context.Context) error {
    // Plugin manages service lifecycle
    return d.connectionService.Start(ctx)
}

func (d *DatabasePlugin) Stop(ctx context.Context) error {
    // Plugin manages service cleanup
    return d.connectionService.Stop(ctx)
}
```

### 4. Simple Operation Processing
```go
func (q *QueryOperation) Execute(ctx context.Context, input component.Input) (component.Output, error) {
    // Simple input parsing
    data := input.Data.(map[string]interface{})
    queryType := data["type"].(string)
    querySQL := data["sql"].(string)
    
    // Simple processing (no real SQL logic)
    result := q.processQuery(queryType, querySQL)
    
    return component.Output{Data: result}, nil
}
```

## 🧪 Testing

### Component Testing
```go
func TestDatabaseConnectionService(t *testing.T) {
    runtime, err := runtime.NewRuntimeWithOptions(
        runtime.WithPlugins(database.NewDatabasePlugin("postgres", "test://connection")),
    )
    require.NoError(t, err)
    
    ctx := context.NewContext()
    err = runtime.Start(ctx)
    require.NoError(t, err)
    defer runtime.Stop(ctx)
    
    // Service should be running
    registry := runtime.Registry()
    service, err := registry.Get("database-connection")
    require.NoError(t, err)
    
    dbService := service.(*database.DatabaseConnectionService)
    assert.True(t, dbService.IsConnected())
}
```

### Operation Testing
```go
func TestQueryOperation(t *testing.T) {
    result, err := runtime.NewBuilder().
        WithPlugins(database.NewDatabasePlugin("postgres", "test://connection")).
        BuildCommand("database-query", map[string]interface{}{
            "type": "validate",
            "sql":  "SELECT id FROM users",
        })
    
    require.NoError(t, err)
    assert.Equal(t, true, result["valid"])
}
```

## 🎯 Key Simplifications

### What This Plugin Demonstrates:
- ✅ **Framework Patterns**: Component lifecycle, plugin orchestration
- ✅ **Service Management**: Start/stop lifecycle in daemon mode
- ✅ **Operation Processing**: Simple input/output transformation
- ✅ **Configuration Access**: Type-safe config with defaults
- ✅ **Logging Integration**: Structured logging throughout

### What This Plugin Avoids:
- ❌ **Real Database**: No actual database connections or SQL execution
- ❌ **Complex SQL**: No real SQL parsing or validation
- ❌ **Connection Pooling**: No real connection management
- ❌ **Transaction Handling**: No database transaction logic
- ❌ **External Dependencies**: No database driver dependencies

## 📚 Next Steps

1. **Copy this plugin** as a template for your own plugins
2. **Modify the components** for your specific needs
3. **Keep the framework patterns** - runtime reference, lifecycle management
4. **Focus on your business logic** - not framework demonstration
5. **Test both modes** - daemon and command execution

---

**Remember**: This plugin demonstrates **framework usage patterns**. For real database functionality, implement actual database handling while maintaining these framework patterns. 