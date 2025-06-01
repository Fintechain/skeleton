# Unit Testing Framework Implementation Prompts

This document contains 6 focused prompts for implementing comprehensive unit tests for the Skeleton Framework infrastructure implementations. These prompts are designed to create a robust testing framework that emphasizes mock-driven dependency injection and comprehensive coverage of all infrastructure components.

**Prerequisites**: These unit test implementations depend on the infrastructure implementations being completed first. The infrastructure layer (`skeleton/internal/infrastructure/`) must contain concrete implementations before comprehensive unit tests can be created.

## Prompt 1: Mock Framework Foundation (Priority 1)

**Objective**: Create a centralized mock framework and core infrastructure mocks for dependency injection testing

**Context**: 
- The Skeleton Framework uses domain interfaces defined in `skeleton/internal/domain/` for all dependencies
- All infrastructure tests must use domain interface mocks for complete isolation
- The framework has its own context interface defined in `skeleton/internal/domain/context/context.go` and re-exported through `skeleton/pkg/context/context.go` (this is NOT Go's standard context)
- Mock framework should provide consistent, configurable behavior for all test scenarios

**Implementation Requirements**:

1. **Central Mock Factory** (`skeleton/test/unit/mocks/factory.go`)
   - Create `Factory` struct for centralized mock creation and management
   - Provide `NewFactory() *Factory` constructor function
   - Implement mock creators for all domain interfaces:
     - `RegistryInterface() registry.Registry`
     - `ComponentInterface() component.Component` 
     - `ContextInterface() context.Context` (framework context, not Go standard context)
     - `SystemInterface() system.System`
     - Additional interface creators as needed

2. **Core Infrastructure Mocks** (`skeleton/test/unit/mocks/`)
   - **Registry Mock** (`registry_mock.go`) - Configurable registry interface implementation
   - **Component Mock** (`component_mock.go`) - Configurable component interface implementation  
   - **Context Mock** (`context_mock.go`) - Framework context interface implementation
   - Each mock should support:
     - Configurable return values for different test scenarios
     - Error injection capabilities
     - Call tracking and verification
     - State management for test verification

3. **Mock Configuration Capabilities**
   - **Behavior Configuration**: Set specific return values and error conditions
   - **State Verification**: Track method calls, parameters, and internal state
   - **Thread Safety**: Ensure mocks are thread-safe where real implementations require it
   - **Error Injection**: Support for testing error conditions and edge cases

**Key Mock Patterns**:
```go
// Central factory for consistent mock creation
type Factory struct { /* centralized mock factory */ }
func NewFactory() *Factory

// Configurable registry mock
type MockRegistry struct {
    items map[string]registry.Identifiable
    shouldFail bool
    callCount map[string]int
}

// Framework context mock (not Go standard context)
type MockContext struct {
    values map[interface{}]interface{}
    deadline time.Time
    done chan struct{}
    err error
}
```

**Testing Framework Integration**:
- All mocks implement exact domain interface signatures
- Consistent mock creation patterns across all test suites
- Support for both positive and negative test scenarios
- Integration with standard Go testing patterns

**Files to Create**:
- `skeleton/test/unit/mocks/factory.go`
- `skeleton/test/unit/mocks/registry_mock.go`
- `skeleton/test/unit/mocks/component_mock.go`
- `skeleton/test/unit/mocks/context_mock.go`

---

## Prompt 2: Core Infrastructure Tests (Priority 2)

**Objective**: Implement comprehensive tests for core infrastructure implementations (Registry, Context, Base Component)

**Context**:
- Tests target concrete implementations in `skeleton/internal/infrastructure/`
- All dependencies must be mocked using domain interfaces from the mock factory
- Framework context interface from `skeleton/pkg/context` is used throughout (not Go's standard context)
- Error handling uses string constants consistently across all domains
- Factory functions follow the pattern: `NewFactory(registry Registry) Factory`

**Implementation Requirements**:

1. **Registry Infrastructure Tests** (`skeleton/test/unit/infrastructure/registry/registry_test.go`)
   - Test `NewRegistry()` constructor function
   - Verify interface compliance: `var _ registry.Registry = registry.NewRegistry()`
   - Test all registry operations: Register, Get, Remove, Has, Count, List, Clear
   - Test thread-safety with concurrent operations
   - Test error conditions using string constants: `registry.ErrItemNotFound`, `registry.ErrItemAlreadyExists`, `registry.ErrInvalidItem`

2. **Context Infrastructure Tests** (`skeleton/test/unit/infrastructure/context/context_test.go`)
   - Test `NewContext()` and `WrapContext()` constructor functions
   - Test framework context operations: Value, WithValue, Deadline, Done, Err
   - Test context cancellation and timeout behavior
   - Test context chaining and value inheritance
   - Verify framework context interface compliance (not Go's standard context)

3. **Base Component Infrastructure Tests** (`skeleton/test/unit/infrastructure/component/`)
   - **Base Component Tests** (`base_component_test.go`)
     - Test `NewBaseComponent(config ComponentConfig) Component` constructor
     - Test component properties: ID, Name, Type, Description, Version
     - Test component lifecycle: Initialize, Dispose
     - Test with mock system and framework context dependencies
   - **Factory Tests** (`factory_test.go`)
     - Test `NewFactory(registry Registry) Factory` constructor with mock registry
     - Test component creation from ComponentConfig
     - Test factory integration with registry dependency

**Key Testing Patterns**:
```go
// Constructor function testing with dependency injection
func TestNewFactory(t *testing.T) {
    factory := mocks.NewFactory()
    mockRegistry := factory.RegistryInterface()
    
    // Test actual factory function signature
    componentFactory := component.NewFactory(mockRegistry)
    assert.NotNil(t, componentFactory)
    
    // Verify interface compliance
    var _ component.Factory = componentFactory
}

// Framework context usage (not Go standard context)
func TestComponentInitialization(t *testing.T) {
    // Use framework context from pkg/context
    ctx := context.NewContext()
    
    config := component.NewComponentConfig("test", "Test Component", component.TypeBasic, "Test")
    comp := component.NewBaseComponent(config)
    
    mockSystem := mocks.NewFactory().SystemInterface()
    err := comp.Initialize(ctx, mockSystem)
    assert.NoError(t, err)
}

// Error handling with string constants
func TestErrorHandling(t *testing.T) {
    registry := registry.NewRegistry()
    
    _, err := registry.Get("non-existent")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), registry.ErrItemNotFound)
}
```

**Files to Create**:
- `skeleton/test/unit/infrastructure/registry/registry_test.go`
- `skeleton/test/unit/infrastructure/context/context_test.go`
- `skeleton/test/unit/infrastructure/component/base_component_test.go`
- `skeleton/test/unit/infrastructure/component/factory_test.go`

---

## Prompt 3: Component System Tests (Priority 3)

**Objective**: Implement tests for advanced component functionality (Dependency-Aware, Lifecycle-Aware, Integration)

**Context**:
- Builds on core infrastructure tests from Prompt 2
- Tests advanced component wrappers and lifecycle management
- All dependencies mocked using domain interfaces
- Framework context from `skeleton/pkg/context` used consistently
- Factory functions: `NewDependencyAwareComponent(base Component, registry Registry)`, `NewLifecycleAwareComponent(base Component)`

**Implementation Requirements**:

1. **Dependency-Aware Component Tests** (`skeleton/test/unit/infrastructure/component/dependency_aware_test.go`)
   - Test `NewDependencyAwareComponent(base Component, registry Registry)` constructor with mocks
   - Test dependency management: Dependencies, AddDependency, RemoveDependency, HasDependency
   - Test dependency resolution: ResolveDependency, ResolveDependencies
   - Test circular dependency detection and prevention
   - Use mock registry for dependency resolution testing

2. **Lifecycle-Aware Component Tests** (`skeleton/test/unit/infrastructure/component/lifecycle_aware_test.go`)
   - Test `NewLifecycleAwareComponent(base Component)` constructor with mock component
   - Test lifecycle state management: State, SetState
   - Test state change callbacks: OnStateChange
   - Test thread-safe state transitions
   - Test all lifecycle states: Created, Initializing, Initialized, Active, Disposing, Disposed, Failed

3. **Component Integration Tests** (`skeleton/test/unit/infrastructure/component/integration_test.go`)
   - Test component factory integration with dependency-aware and lifecycle-aware wrappers
   - Test complex component creation scenarios
   - Test component interaction patterns through interfaces
   - Use framework context for all component operations

**Advanced Mock Requirements**:
```go
// Additional mocks for component system testing
type MockSystem struct { /* system interface mock */ }
type MockConfiguration struct { /* config interface mock */ }
type MockEventBus struct { /* event bus interface mock */ }
```

**Key Testing Patterns**:
```go
// Dependency resolution testing
func TestDependencyResolution(t *testing.T) {
    factory := mocks.NewFactory()
    mockRegistry := factory.RegistryInterface()
    mockComponent := factory.ComponentInterface()
    
    // Configure mock registry behavior
    mockRegistry.SetReturnItem("dependency-id", mockComponent)
    
    // Test dependency-aware component
    depAware := component.NewDependencyAwareComponent(mockComponent, mockRegistry)
    resolved, err := depAware.ResolveDependency("dependency-id", mockRegistry)
    
    assert.NoError(t, err)
    assert.NotNil(t, resolved)
}

// Lifecycle state management testing
func TestLifecycleStateManagement(t *testing.T) {
    factory := mocks.NewFactory()
    mockComponent := factory.ComponentInterface()
    
    lifecycleAware := component.NewLifecycleAwareComponent(mockComponent)
    
    // Test state transitions
    lifecycleAware.SetState(component.StateInitializing)
    assert.Equal(t, component.StateInitializing, lifecycleAware.State())
}
```

**Files to Create**:
- `skeleton/test/unit/infrastructure/component/dependency_aware_test.go`
- `skeleton/test/unit/infrastructure/component/lifecycle_aware_test.go`
- `skeleton/test/unit/infrastructure/component/integration_test.go`
- `skeleton/test/unit/mocks/system_mock.go`
- `skeleton/test/unit/mocks/config_mock.go`
- `skeleton/test/unit/mocks/event_bus_mock.go`

---

## Prompt 4: Operation and Service Tests (Priority 4)

**Objective**: Implement tests for Operation and Service infrastructure implementations

**Context**:
- Operations and Services extend Component interface with specialized functionality
- Factory functions: `NewOperation(component Component)`, `NewOperationFactory(componentFactory Factory)`
- Service factory: `NewService(component Component)`, `NewServiceFactory(componentFactory Factory)`
- Framework context from `skeleton/pkg/context` used for all execution methods
- Error handling uses string constants: `service.ErrServiceStart`, `service.ErrServiceStop`, `service.ErrServiceNotFound`

**Implementation Requirements**:

1. **Operation Infrastructure Tests** (`skeleton/test/unit/infrastructure/operation/`)
   - **Operation Tests** (`operation_test.go`)
     - Test `NewOperation(component Component)` constructor with mock component
     - Test `Execute(ctx context.Context, input Input) (Output, error)` method
     - Test operation interface compliance verification
     - Test error handling in operation execution
     - Use framework context for all execution testing
   - **Factory Tests** (`factory_test.go`)
     - Test `NewOperationFactory(componentFactory Factory)` constructor with mock factory
     - Test operation creation from OperationConfig
     - Test factory integration with component system

2. **Service Infrastructure Tests** (`skeleton/test/unit/infrastructure/service/`)
   - **Service Tests** (`service_test.go`)
     - Test `NewService(component Component)` constructor with mock component
     - Test service lifecycle: Start, Stop, Status
     - Test thread-safe status transitions
     - Test all service status states: Stopped, Starting, Running, Stopping, Failed
     - Use framework context for Start/Stop operations
   - **Factory Tests** (`factory_test.go`)
     - Test `NewServiceFactory(componentFactory Factory)` constructor with mock factory
     - Test service creation from ServiceConfig
     - Test factory integration with component system

**Key Testing Patterns**:
```go
// Operation execution testing with framework context
func TestOperationExecution(t *testing.T) {
    factory := mocks.NewFactory()
    mockComponent := factory.ComponentInterface()
    
    operation := operation.NewOperation(mockComponent)
    
    // Use framework context (not Go standard context)
    ctx := context.NewContext()
    input := map[string]interface{}{"key": "value"}
    
    result, err := operation.Execute(ctx, input)
    assert.NoError(t, err)
    assert.NotNil(t, result)
}

// Service lifecycle testing with framework context
func TestServiceLifecycle(t *testing.T) {
    factory := mocks.NewFactory()
    mockComponent := factory.ComponentInterface()
    
    service := service.NewService(mockComponent)
    
    // Use framework context for lifecycle operations
    ctx := context.NewContext()
    
    // Test start
    err := service.Start(ctx)
    assert.NoError(t, err)
    assert.Equal(t, service.StatusRunning, service.Status())
    
    // Test stop
    err = service.Stop(ctx)
    assert.NoError(t, err)
    assert.Equal(t, service.StatusStopped, service.Status())
}

// Error handling with string constants
func TestServiceErrorHandling(t *testing.T) {
    // Test service error conditions
    err := someServiceOperation()
    assert.Error(t, err)
    assert.Contains(t, err.Error(), service.ErrServiceStart)
}
```

**Files to Create**:
- `skeleton/test/unit/infrastructure/operation/operation_test.go`
- `skeleton/test/unit/infrastructure/operation/factory_test.go`
- `skeleton/test/unit/infrastructure/service/service_test.go`
- `skeleton/test/unit/infrastructure/service/factory_test.go`

---

## Prompt 5: System Integration Tests (Priority 5)

**Objective**: Implement tests for system-level infrastructure that coordinates all framework resources

**Context**:
- System interface provides centralized access to all framework resources
- Factory function: `NewSystem(registry Registry, pluginManager PluginManager, eventBus EventBus, configuration Configuration, store MultiStore)`
- All system operations use framework context from `skeleton/pkg/context`
- Plugin and configuration implementations with their respective factory functions
- Error handling uses string constants consistently

**Implementation Requirements**:

1. **System Infrastructure Tests** (`skeleton/test/unit/infrastructure/system/system_test.go`)
   - Test `NewSystem()` constructor with all mock dependencies
   - Test resource access methods: Registry(), PluginManager(), EventBus(), Configuration(), Store()
   - Test system operations: ExecuteOperation(), StartService(), StopService()
   - Test system state: IsRunning(), IsInitialized()
   - Use framework context for all system operations
   - Test error handling with string constants: `system.ErrSystemNotInitialized`, `system.ErrOperationNotFound`, etc.

2. **Plugin Manager Tests** (`skeleton/test/unit/infrastructure/plugin/plugin_manager_test.go`)
   - Test `NewPluginManager(filesystem FileSystem)` constructor with mock filesystem
   - Test plugin lifecycle: Discover(), Load(), Unload()
   - Test plugin information: ListPlugins(), GetPlugin()
   - Test integration with Registry for plugin components
   - Test error handling with string constants: `plugin.ErrPluginNotFound`, `plugin.ErrPluginLoad`, etc.

3. **Configuration Tests** (`skeleton/test/unit/infrastructure/config/config_test.go`)
   - Test `NewConfiguration(sources ...ConfigurationSource)` constructor with mock sources
   - Test all getter methods: GetString(), GetInt(), GetBool(), GetDuration(), GetObject()
   - Test default value methods: GetStringDefault(), GetIntDefault(), etc.
   - Test key existence checking: Exists()
   - Test error handling with string constants: `config.ErrConfigNotFound`, `config.ErrConfigWrongType`, etc.

**System Integration Patterns**:
```go
// System creation with all mock dependencies
func TestSystemWithMockDependencies(t *testing.T) {
    factory := mocks.NewFactory()
    
    mockRegistry := factory.RegistryInterface()
    mockPluginManager := factory.PluginManagerInterface()
    mockEventBus := factory.EventBusInterface()
    mockConfig := factory.ConfigurationInterface()
    mockStore := factory.MultiStoreInterface()
    
    // Test system creation with actual factory function signature
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

// System operations with framework context
func TestSystemOperations(t *testing.T) {
    // Use framework context for all system operations
    ctx := context.NewContext()
    
    // Test operation execution
    result, err := system.ExecuteOperation(ctx, "operation-id", inputData)
    assert.NoError(t, err)
    
    // Test service management
    err = system.StartService(ctx, "service-id")
    assert.NoError(t, err)
}
```

**Additional Mock Requirements**:
- `skeleton/test/unit/mocks/plugin_manager_mock.go`
- `skeleton/test/unit/mocks/config_source_mock.go`

**Files to Create**:
- `skeleton/test/unit/infrastructure/system/system_test.go`
- `skeleton/test/unit/infrastructure/plugin/plugin_manager_test.go`
- `skeleton/test/unit/infrastructure/config/config_test.go`
- `skeleton/test/unit/mocks/plugin_manager_mock.go`
- `skeleton/test/unit/mocks/config_source_mock.go`

---

## Prompt 6: Supporting Infrastructure Tests (Priority 6)

**Objective**: Implement tests for Event System, Storage, and Logging infrastructure

**Context**:
- Supporting systems enhance framework capabilities with event handling, persistence, and logging
- Factory functions: `NewEventBus()`, `NewMultiStore()`, `NewLogger()`
- Storage system uses string constants for error handling (consistent with all other domains)
- Framework context used where applicable
- Thread-safety requirements for concurrent operations

**Implementation Requirements**:

1. **Event System Tests** (`skeleton/test/unit/infrastructure/event/event_test.go`)
   - Test `NewEventBus()` constructor
   - Test event publishing: Publish()
   - Test event subscription: Subscribe(), SubscribeAsync()
   - Test subscription management: Cancel(), Topic()
   - Test thread-safety for concurrent publish/subscribe operations
   - Test asynchronous event handling: WaitAsync()

2. **Storage System Tests** (`skeleton/test/unit/infrastructure/storage/`)
   - **MultiStore Tests** (`multistore_test.go`)
     - Test `NewMultiStore()` constructor
     - Test store management: GetStore(), CreateStore(), DeleteStore(), ListStores(), StoreExists()
     - Test engine management: RegisterEngine(), ListEngines(), GetEngine()
     - Test bulk operations: CloseAll()
   - **Engine Tests** (`engine_test.go`)
     - Test memory engine implementation
     - Test file engine implementation
     - Test engine capabilities and configuration
   - **Store Tests** (`store_test.go`)
     - Test basic CRUD operations: Get, Set, Delete, Has
     - Test store iteration capabilities
     - Test storage error handling with string constants: `storage.ErrKeyNotFound`, `storage.ErrStoreNotFound`, etc.

3. **Logging System Tests** (`skeleton/test/unit/infrastructure/logging/logging_test.go`)
   - Test `NewLogger()` constructor
   - Test logging methods: Debug(), Info(), Warn(), Error()
   - Test structured logging capabilities
   - Test log level handling and formatting
   - Test error handling with string constants: `logging.ErrLoggerNotAvailable`, `logging.ErrInvalidLogLevel`

**Key Testing Patterns**:
```go
// Event system testing
func TestEventBusOperations(t *testing.T) {
    eventBus := event.NewEventBus()
    
    // Test subscription
    subscription := eventBus.Subscribe("test.topic", func(event *event.Event) {
        // Handle event
    })
    
    // Test publishing
    eventBus.Publish("test.topic", map[string]interface{}{"key": "value"})
    
    // Test subscription management
    subscription.Cancel()
    assert.Equal(t, "test.topic", subscription.Topic())
}

// Storage error handling with string constants (not Go error variables)
func TestStorageErrorHandling(t *testing.T) {
    store, _ := multiStore.GetStore("test-store")
    
    _, err := store.Get([]byte("non-existent-key"))
    assert.Error(t, err)
    assert.Contains(t, err.Error(), storage.ErrKeyNotFound)
}

// Thread-safety testing
func TestConcurrentEventOperations(t *testing.T) {
    eventBus := event.NewEventBus()
    
    var wg sync.WaitGroup
    numGoroutines := 10
    
    // Test concurrent publish/subscribe
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            topic := fmt.Sprintf("topic-%d", id)
            eventBus.Subscribe(topic, func(event *event.Event) {})
            eventBus.Publish(topic, map[string]interface{}{"id": id})
        }(i)
    }
    
    wg.Wait()
}
```

**Additional Mock Requirements**:
- `skeleton/test/unit/mocks/storage_mock.go`
- `skeleton/test/unit/mocks/multistore_mock.go`
- `skeleton/test/unit/mocks/engine_mock.go`
- `skeleton/test/unit/mocks/logger_mock.go`

**Files to Create**:
- `skeleton/test/unit/infrastructure/event/event_test.go`
- `skeleton/test/unit/infrastructure/storage/multistore_test.go`
- `skeleton/test/unit/infrastructure/storage/engine_test.go`
- `skeleton/test/unit/infrastructure/storage/store_test.go`
- `skeleton/test/unit/infrastructure/logging/logging_test.go`
- `skeleton/test/unit/mocks/storage_mock.go`
- `skeleton/test/unit/mocks/multistore_mock.go`
- `skeleton/test/unit/mocks/engine_mock.go`
- `skeleton/test/unit/mocks/logger_mock.go`

---

## Implementation Notes

### Framework Context Usage
- **Always use framework context**: Import from `github.com/fintechain/skeleton/pkg/context`
- **Not Go standard context**: The framework has its own context interface
- **Consistent usage**: All system operations, component initialization, and lifecycle methods use framework context
- **Context creation**: Use `context.NewContext()` and `context.WrapContext()` from the framework

### Error Handling Patterns
- **String constants across all domains**: All framework domains use string constants for error identification
- **Consistent pattern**: Check errors using `err.Error() == ErrorConstant`
- **No Go error variables**: All domains use string constants, including storage
- **Domain-specific prefixes**: `registry.*`, `component.*`, `service.*`, `system.*`, `config.*`, `plugin.*`, `storage.*`, `logging.*`

### Factory Function Signatures
- **Component Factory**: `component.NewFactory(registry Registry) Factory`
- **Operation Factory**: `operation.NewOperationFactory(componentFactory Factory) OperationFactory`
- **Service Factory**: `service.NewServiceFactory(componentFactory Factory) ServiceFactory`
- **System**: `system.NewSystem(registry, pluginManager, eventBus, configuration, store) System`

### Mock-Driven Testing Requirements
- **All dependencies mocked**: Use domain interface mocks from central factory
- **Interface compliance**: Verify all implementations satisfy domain contracts
- **Constructor testing**: Test all factory functions with mock dependencies
- **Thread safety**: Test concurrent operations where applicable
- **Error injection**: Support testing error conditions and edge cases

### Prerequisites and Dependencies
- **Infrastructure implementations required**: Unit tests depend on concrete implementations in `skeleton/internal/infrastructure/`
- **Implementation order**: Infrastructure implementations should be completed before or alongside unit test creation
- **Domain interfaces**: All domain interfaces in `skeleton/internal/domain/` are complete and ready for mocking
- **Public API**: Tests verify that public API in `skeleton/pkg/` correctly exposes infrastructure implementations

### Quality Standards
- **>95% infrastructure coverage** for all implementations
- **100% constructor coverage** for all factory functions
- **100% interface compliance** verification
- **>90% error path coverage** for all error conditions
- **Thread-safety testing** for concurrent components
- **Mock-driven isolation** with no dependencies on real external systems