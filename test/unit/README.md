# Unit Testing Framework for Skeleton Framework

This directory contains a comprehensive unit testing framework for the Skeleton Framework infrastructure implementations. The framework emphasizes consistency, reliability, and maintainability through extensive use of mocks and standardized testing patterns.

## 🏗️ Architecture Overview

```
test/unit/
├── README.md                   # This documentation
├── mocks/                      # Centralized mock implementations
│   ├── factory.go             # Mock factory for dependency injection
│   ├── domain_mocks.go        # Domain interface mocks
│   ├── registry_mock.go       # Registry interface mocks
│   ├── component_mock.go      # Component interface mocks
│   ├── context_mock.go        # Context interface mocks
│   ├── system_mock.go         # System interface mocks
│   ├── config_mock.go         # Configuration system mocks
│   ├── logger_mock.go         # Logging system mocks
│   ├── storage_mock.go        # Storage system mocks
│   ├── multistore_mock.go     # Multi-store management mocks
│   ├── engine_mock.go         # Storage engine mocks
│   ├── event_bus_mock.go      # Event system mocks
│   └── plugin_manager_mock.go # Plugin management mocks
├── infrastructure/            # Infrastructure implementation tests
│   ├── registry/             # Registry implementation tests
│   ├── component/            # Component implementation tests
│   ├── context/              # Context implementation tests
│   ├── operation/            # Operation implementation tests
│   ├── service/              # Service implementation tests
│   ├── system/               # System implementation tests
│   ├── config/               # Configuration implementation tests
│   ├── event/                # Event system implementation tests
│   ├── storage/              # Storage implementation tests
│   ├── plugin/               # Plugin implementation tests
│   └── logging/              # Logging implementation tests
└── pkg/                       # Public API tests
    ├── registry/             # Registry package tests
    ├── component/            # Component package tests
    ├── operation/            # Operation package tests
    ├── service/              # Service package tests
    ├── system/               # System package tests
    ├── config/               # Configuration package tests
    ├── event/                # Event package tests
    ├── storage/              # Storage package tests
    ├── plugin/               # Plugin package tests
    └── logging/              # Logging package tests
```

## 🎯 Testing Philosophy

Our testing framework is built on these core principles:

### 1. **Infrastructure-First Testing**
- Test concrete implementations in `skeleton/internal/infrastructure/`
- Verify interface compliance with domain contracts
- Test dependency injection patterns and constructor functions
- Validate thread-safety and error handling

### 2. **Mock-Driven Dependency Injection**
- All dependencies are mocked using domain interfaces
- Tests are isolated and don't depend on real implementations
- Constructor functions tested with mock dependencies
- Consistent behavior across all test suites

### 3. **Comprehensive Coverage**
- Every infrastructure implementation is tested
- Factory functions and constructors are verified
- Interface compliance is validated
- Error conditions and edge cases are covered

### 4. **Clean Architecture Compliance**
- Infrastructure tests use domain interface mocks
- No circular dependencies in test setup
- Clear separation between infrastructure and domain testing

## 🔧 Mock Framework

### Central Mock Factory

The `mocks.Factory` provides centralized access to all mock implementations:

```go
import (
    "github.com/fintechain/skeleton/pkg/registry"
    "github.com/fintechain/skeleton/pkg/component"
)

func TestRegistryImplementation(t *testing.T) {
    // Create factory for consistent mock creation
    factory := mocks.NewFactory()
    
    // Test infrastructure implementation
    registry := registry.NewRegistry()
    
    // Verify interface compliance
    var _ registry.Registry = registry
    
    // Test basic operations
    item := factory.ComponentInterface()
    err := registry.Register(item)
    assert.NoError(t, err)
}
```

### Domain Interface Mocks

For testing infrastructure implementations that depend on domain interfaces:

```go
import (
    "github.com/fintechain/skeleton/pkg/component"
    "github.com/fintechain/skeleton/pkg/registry"
)

// Test component factory with mock registry dependency
func TestComponentFactory(t *testing.T) {
    factory := mocks.NewFactory()
    mockRegistry := factory.RegistryInterface()
    
    // Test infrastructure implementation
    componentFactory := component.NewFactory(mockRegistry)
    
    config := component.NewComponentConfig("test", "Test", component.TypeBasic, "Test component")
    comp, err := componentFactory.Create(config)
    
    assert.NoError(t, err)
    assert.NotNil(t, comp)
}
```

### Available Mock Implementations

#### 🏗️ Infrastructure Mocks

**Domain Registry Mock** - For testing components with registry dependencies:
```go
import "github.com/fintechain/skeleton/pkg/registry"

mockRegistry := factory.RegistryInterface()

// Configure mock behavior
mockRegistry.SetReturnItem("test-id", mockComponent)
mockRegistry.SetShouldFail(false)

// Use in infrastructure tests
componentFactory := component.NewFactory(mockRegistry)
```

**Domain Component Mock** - For testing component wrappers:
```go
import "github.com/fintechain/skeleton/pkg/component"

mockComponent := factory.ComponentInterface()

// Configure component behavior
mockComponent.SetID("test-component")
mockComponent.SetType(component.TypeBasic)

// Use in dependency-aware component tests
depAware := component.NewDependencyAwareComponent(mockComponent, mockRegistry)
```

**Domain System Mock** - For testing component initialization:
```go
import (
    "github.com/fintechain/skeleton/pkg/system"
    "github.com/fintechain/skeleton/pkg/context"
)

mockSystem := factory.SystemInterface()

// Configure system resources
mockSystem.SetRegistry(mockRegistry)
mockSystem.SetEventBus(mockEventBus)

// Test component initialization (using framework context)
ctx := context.NewContext()
err := component.Initialize(ctx, mockSystem)
assert.NoError(t, err)
```

#### 🗃️ Storage Mocks

**MockStorage** - Complete storage interface implementation:
```go
import "github.com/fintechain/skeleton/pkg/storage"

storage := factory.Storage().
    WithKeyValue("key1", "value1").
    WithKeyValue("key2", "value2").
    WithName("test-store").
    Build()

// Test storage operations
value, err := storage.Get([]byte("key1"))
assert.NoError(t, err)
assert.Equal(t, []byte("value1"), value)
```

**MockMultiStore** - Multi-store management:
```go
import "github.com/fintechain/skeleton/pkg/storage"

multiStore := factory.MultiStore()

// Test store creation
err := multiStore.CreateStore("store1", "memory", nil)
assert.NoError(t, err)

// Verify store exists
assert.True(t, multiStore.StoreExists("store1"))
```

#### 📝 Configuration Mocks

**MockConfiguration** - Type-safe configuration testing:
```go
import "github.com/fintechain/skeleton/pkg/config"

config := factory.Config().
    WithString("database.host", "localhost").
    WithInt("database.port", 5432).
    WithBool("debug.enabled", true).
    Build()

// Test configuration implementation
configImpl := config.NewConfiguration(mockSource)
assert.Equal(t, "localhost", configImpl.GetString("database.host"))
```

**MockConfigurationSource** - Configuration source testing:
```go
import "github.com/fintechain/skeleton/pkg/config"

mockSource := factory.ConfigurationSource()
mockSource.SetValue("key", "value")

// Test configuration with mock source
config := config.NewConfiguration(mockSource)
```

#### 🚌 Event System Mocks

**MockEventBus** - Event publishing and subscription:
```go
import "github.com/fintechain/skeleton/pkg/event"

eventBus := factory.EventBus()

// Test event bus implementation
eventImpl := event.NewEventBus()

// Publish and verify
eventImpl.Publish("test.topic", map[string]interface{}{"key": "value"})

events := eventBus.GetEventsByTopic("test.topic")
assert.Len(t, events, 1)
```

## 📋 Testing Patterns

### 1. Infrastructure Implementation Testing

Test concrete implementations with dependency injection:

```go
import (
    "testing"
    "github.com/fintechain/skeleton/pkg/registry"
    "github.com/stretchr/testify/assert"
)

func TestRegistryImplementation(t *testing.T) {
    tests := []struct {
        name        string
        operation   func(registry.Registry) error
        expectError bool
        description string
    }{
        {
            name: "register valid item",
            operation: func(r registry.Registry) error {
                item := mocks.NewFactory().ComponentInterface()
                return r.Register(item)
            },
            expectError: false,
            description: "Should register valid items successfully",
        },
        {
            name: "register nil item",
            operation: func(r registry.Registry) error {
                return r.Register(nil)
            },
            expectError: true,
            description: "Should reject nil items",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            registry := registry.NewRegistry()
            err := tt.operation(registry)
            
            if tt.expectError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### 2. Constructor Function Testing

Test factory functions and dependency injection:

```go
import (
    "testing"
    "github.com/fintechain/skeleton/pkg/component"
    "github.com/stretchr/testify/assert"
)

func TestNewFactory(t *testing.T) {
    factory := mocks.NewFactory()
    mockRegistry := factory.RegistryInterface()
    
    // Test constructor (note: actual function name is NewFactory, not NewComponentFactory)
    componentFactory := component.NewFactory(mockRegistry)
    assert.NotNil(t, componentFactory)
    
    // Verify interface compliance
    var _ component.Factory = componentFactory
}
```

### 3. Interface Compliance Testing

Verify implementations satisfy domain contracts:

```go
import (
    "testing"
    "github.com/fintechain/skeleton/pkg/registry"
    "github.com/fintechain/skeleton/pkg/component"
    "github.com/fintechain/skeleton/pkg/operation"
    "github.com/fintechain/skeleton/pkg/service"
    "github.com/fintechain/skeleton/pkg/system"
)

func TestInterfaceCompliance(t *testing.T) {
    // Test all infrastructure implementations
    var _ registry.Registry = registry.NewRegistry()
    var _ component.Component = component.NewBaseComponent(config)
    var _ operation.Operation = operation.NewOperation(mockComponent)
    var _ service.Service = service.NewService(mockComponent)
    var _ system.System = system.NewSystem(deps...)
}
```

### 4. Dependency Injection Testing

Test components with mock dependencies:

```go
import (
    "testing"
    "github.com/fintechain/skeleton/pkg/system"
    "github.com/fintechain/skeleton/pkg/registry"
    "github.com/fintechain/skeleton/pkg/plugin"
    "github.com/fintechain/skeleton/pkg/event"
    "github.com/fintechain/skeleton/pkg/config"
    "github.com/fintechain/skeleton/pkg/storage"
    "github.com/stretchr/testify/assert"
)

func TestSystemWithMockDependencies(t *testing.T) {
    factory := mocks.NewFactory()
    
    // Create all mock dependencies
    mockRegistry := factory.RegistryInterface()
    mockPluginManager := factory.PluginManagerInterface()
    mockEventBus := factory.EventBusInterface()
    mockConfig := factory.ConfigurationInterface()
    mockStore := factory.MultiStoreInterface()
    
    // Test system creation
    system := system.NewSystem(
        mockRegistry,
        mockPluginManager,
        mockEventBus,
        mockConfig,
        mockStore,
    )
    
    assert.NotNil(t, system)
    assert.True(t, system.IsInitialized())
}
```

### 5. Error Handling Verification

Test error conditions with proper error checking:

```go
import (
    "testing"
    "github.com/fintechain/skeleton/pkg/registry"
    "github.com/stretchr/testify/assert"
)

func TestErrorHandling(t *testing.T) {
    registry := registry.NewRegistry()
    
    // Test item not found
    _, err := registry.Get("non-existent")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), registry.ErrItemNotFound)
    
    // Test duplicate registration
    item := mocks.NewFactory().ComponentInterface()
    registry.Register(item)
    err = registry.Register(item)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), registry.ErrItemAlreadyExists)
}
```

### 6. Thread-Safety Testing

Verify concurrent operations:

```go
import (
    "fmt"
    "sync"
    "testing"
    "github.com/fintechain/skeleton/pkg/registry"
    "github.com/stretchr/testify/assert"
)

func TestConcurrentRegistryAccess(t *testing.T) {
    registry := registry.NewRegistry()
    factory := mocks.NewFactory()
    
    var wg sync.WaitGroup
    numGoroutines := 10
    
    // Concurrent registrations
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            item := factory.ComponentInterface()
            item.SetID(fmt.Sprintf("item-%d", id))
            err := registry.Register(item)
            assert.NoError(t, err)
        }(i)
    }
    
    wg.Wait()
    assert.Equal(t, numGoroutines, registry.Count())
}
```

### 7. Framework Context Usage

Test with the framework's context interface (not Go's standard context):

```go
import (
    "testing"
    "github.com/fintechain/skeleton/pkg/context"
    "github.com/fintechain/skeleton/pkg/component"
    "github.com/stretchr/testify/assert"
)

func TestComponentInitialization(t *testing.T) {
    // Use framework context, not standard Go context
    ctx := context.NewContext()
    
    // Test component initialization
    config := component.NewComponentConfig("test", "Test Component", component.TypeBasic, "Test")
    comp := component.NewBaseComponent(config)
    
    mockSystem := mocks.NewFactory().SystemInterface()
    err := comp.Initialize(ctx, mockSystem)
    assert.NoError(t, err)
}
```

## 📊 Test Coverage by Infrastructure Component

### ✅ Registry Infrastructure (`internal/infrastructure/registry`)
- **NewRegistry()**: Constructor function testing
- **Thread Safety**: Concurrent access patterns
- **Interface Compliance**: Domain registry interface
- **Error Handling**: All error conditions covered
- **CRUD Operations**: Register, get, remove, list, clear

### ✅ Component Infrastructure (`internal/infrastructure/component`)
- **Base Component**: NewBaseComponent() with configuration
- **Component Factory**: NewFactory() with registry dependency
- **Dependency-Aware**: NewDependencyAwareComponent() with resolution
- **Lifecycle-Aware**: NewLifecycleAwareComponent() with state management
- **Interface Compliance**: All component interfaces

### ✅ Context Infrastructure (`internal/infrastructure/context`)
- **Context Creation**: NewContext() and WrapContext()
- **Value Management**: WithValue() and Value() operations
- **Cancellation**: Deadline(), Done(), Err() behavior
- **Interface Compliance**: Framework context interface

### ✅ Operation Infrastructure (`internal/infrastructure/operation`)
- **Operation Creation**: NewOperation() with component dependency
- **Operation Factory**: NewOperationFactory() with factory dependency
- **Execute Method**: Input/output handling and error conditions
- **Interface Compliance**: Operation interface

### ✅ Service Infrastructure (`internal/infrastructure/service`)
- **Service Creation**: NewService() with component dependency
- **Service Factory**: NewServiceFactory() with factory dependency
- **Lifecycle Management**: Start(), Stop(), Status() operations
- **Thread-Safe Status**: Concurrent status changes
- **Interface Compliance**: Service interface

### ✅ System Infrastructure (`internal/infrastructure/system`)
- **System Creation**: NewSystem() with all dependencies
- **Resource Access**: Registry(), EventBus(), Store(), etc.
- **Operation Execution**: ExecuteOperation() with error handling
- **Service Management**: StartService(), StopService()
- **State Management**: IsRunning(), IsInitialized()

### ✅ Supporting Infrastructure
- **Event System**: NewEventBus() with pub/sub patterns
- **Storage System**: NewMultiStore() with engine management
- **Configuration**: NewConfiguration() with source dependencies
- **Plugin System**: NewPluginManager() with filesystem dependency
- **Logging System**: NewLogger() with structured logging

## 🚀 Running Tests

### Run All Infrastructure Tests
```bash
go test ./test/unit/infrastructure/... -v
```

### Run Specific Infrastructure Tests
```bash
go test ./test/unit/infrastructure/registry/... -v
go test ./test/unit/infrastructure/component/... -v
go test ./test/unit/infrastructure/system/... -v
```

### Run Public API Tests
```bash
go test ./test/unit/pkg/... -v
```

### Run with Coverage
```bash
go test ./test/unit/... -v -cover
```

### Run with Race Detection
```bash
go test ./test/unit/... -v -race
```

## 🔍 Test Quality Metrics

### Coverage Goals
- **Infrastructure Coverage**: >95% for all implementations
- **Constructor Coverage**: 100% for all factory functions
- **Interface Compliance**: 100% verification
- **Error Path Coverage**: >90% for all error conditions

### Quality Standards
- ✅ All infrastructure tests use domain interface mocks
- ✅ Constructor functions tested with dependency injection
- ✅ Interface compliance verified for all implementations
- ✅ Thread-safety tested for concurrent components
- ✅ Error handling covers all defined error constants
- ✅ No tests depend on real external systems

## 🛠️ Best Practices

### Infrastructure Testing
1. **Test constructors**: Verify all `New*()` functions work correctly
2. **Mock dependencies**: Use domain interface mocks for all dependencies
3. **Verify compliance**: Ensure implementations satisfy domain interfaces
4. **Test error paths**: Cover all error constants and conditions

### Mock Usage
1. **Use the factory**: `factory := mocks.NewFactory()`
2. **Mock interfaces**: Use domain interface mocks for dependencies
3. **Configure behavior**: Set up expected mock responses
4. **Verify interactions**: Check mock calls and state

### Test Structure
1. **Descriptive names**: Clear indication of what's being tested
2. **Table-driven tests**: Comprehensive scenario coverage
3. **Proper isolation**: Each test uses fresh mocks
4. **Clear assertions**: Specific and meaningful checks

### Dependency Injection Testing
1. **Test with mocks**: All dependencies should be mocked
2. **Verify construction**: Ensure constructors accept correct interfaces
3. **Test behavior**: Verify components work with mock dependencies
4. **Check integration**: Test how components interact through interfaces

### Context Usage
1. **Use framework context**: Import from `github.com/fintechain/skeleton/pkg/context`
2. **Not Go standard context**: The framework has its own context interface
3. **Consistent usage**: All system operations use framework context
4. **Test context operations**: Value management, cancellation, deadlines

## 📚 Contributing to Tests

When adding new infrastructure tests:

1. **Follow patterns**: Use established infrastructure testing patterns
2. **Mock dependencies**: Use domain interface mocks from factory
3. **Test constructors**: Verify all factory functions
4. **Check compliance**: Ensure interface compliance verification
5. **Cover errors**: Test all error paths and conditions
6. **Document intent**: Clear test names and descriptions
7. **Use correct imports**: Import from public API packages (`pkg/`)
8. **Use framework context**: Not Go's standard context

## 🔧 Troubleshooting

### Common Issues

**Constructor tests failing:**
- Ensure all required dependencies are provided as mocks
- Check that mock interfaces match domain interface signatures
- Verify correct factory function names (e.g., `NewFactory`, not `NewComponentFactory`)

**Interface compliance errors:**
- Verify implementation methods match domain interface exactly
- Check that all interface methods are implemented

**Mock dependency issues:**
- Use domain interface mocks, not concrete implementations
- Ensure mock factory provides correct interface types

**Context-related errors:**
- Use framework context (`pkg/context`), not Go standard context
- Import the correct context package

### Debug Tips

1. **Check interfaces**: Verify domain interface compliance
2. **Review mocks**: Ensure mocks implement correct interfaces
3. **Test isolation**: Verify tests don't share state
4. **Dependency injection**: Check constructor parameter types
5. **Import statements**: Verify correct package imports from `pkg/`

## 📝 Important Notes

### Factory Function Names
- **Component Factory**: `component.NewFactory()` (not `NewComponentFactory`)
- **Operation Factory**: `operation.NewOperationFactory()`
- **Service Factory**: `service.NewServiceFactory()`

### Context Interface
- **Framework Context**: Use `context.NewContext()` from `pkg/context`
- **Not Go Context**: The framework has its own context interface
- **System Operations**: All use framework context, not standard Go context

### Public API Structure
- **Import from `pkg/`**: All public APIs are re-exported through `pkg/` packages
- **Clean Separation**: Public API provides clean interface to internal implementations
- **Consistent Naming**: Factory functions follow consistent patterns

This testing framework provides comprehensive coverage for the Skeleton Framework's infrastructure implementations, ensuring reliability, maintainability, and adherence to Clean Architecture principles. 