# Skeleton Framework Implementation Prompts

This document contains 6 focused prompts for implementing the concrete implementations needed for the skeleton framework, based on the priority list in IMPLEMENTATION_LIST.md and the actual interfaces defined in the codebase.

## Prompt 1: Core Foundation (Priority 1)

**Objective**: Implement the critical foundation components that everything else depends on

**Context**: 
- The skeleton framework has complete interface definitions in `skeleton/internal/domain/` but is missing concrete implementations
- All concrete implementations must be created in `skeleton/internal/infrastructure/` following Clean Architecture principles
- These components are the foundation for all dependency injection and component lifecycle patterns
- The Registry interface is defined in `skeleton/internal/domain/registry/registry.go`
- Component interfaces are defined in `skeleton/internal/domain/component/`
- Context interface is defined in `skeleton/internal/domain/context/context.go`

**Implementation Requirements**:

1. **Registry Implementation** (`skeleton/internal/infrastructure/registry/registry_impl.go`)
   - Implement concrete `Registry` interface with thread-safe operations
   - Provide `NewRegistry()` constructor function that accepts minimal dependencies
   - Implement all methods: `Register()`, `Get()`, `List()`, `Remove()`, `Has()`, `Count()`, `Clear()`
   - Use proper error handling with constants: `ErrItemNotFound`, `ErrItemAlreadyExists`, `ErrInvalidItem`
   - Thread-safe implementation using appropriate synchronization

2. **Base Component Implementation** (`skeleton/internal/infrastructure/component/base_component.go`)
   - Implement `NewBaseComponent()` constructor that accepts configuration and minimal dependencies
   - Provide concrete implementation of `Component` interface
   - Implement `Identifiable` methods: `ID()`, `Name()`, `Description()`, `Version()`
   - Implement `Component` methods: `Type()`, `Metadata()`, `Initialize()`, `Dispose()`
   - Support for component metadata management
   - Follow the actual interface signature: `Initialize(ctx context.Context, system sys.System) error`

3. **Context Implementation** (`skeleton/internal/infrastructure/context/context_impl.go`)
   - Implement concrete `Context` interface
   - Provide `NewContext()` and `WrapContext()` functions
   - Implement methods: `Value()`, `WithValue()`, `Deadline()`, `Done()`, `Err()`
   - Support for context cancellation and timeouts

**Acceptance Criteria**:
- All interface methods implemented correctly
- Thread-safe Registry operations
- Proper error handling with defined error codes from domain interfaces
- Constructor functions available for public API usage
- Context implements framework-specific interface (not Go's standard context patterns)
- **Testable through dependency injection** (constructors accept interface dependencies)
- **External API compatibility maintained** (public API in `skeleton/pkg/` unchanged)
- **Simple, focused implementations** without unnecessary complexity

**Documentation Deliverables**:
- Brief explanation of constructor dependency patterns used
- Any design decisions for thread safety or error handling
- Notes on interface compatibility and public API preservation

**Files to Create**:
- `skeleton/internal/infrastructure/registry/registry_impl.go`
- `skeleton/internal/infrastructure/component/base_component.go`
- `skeleton/internal/infrastructure/context/context_impl.go`

---

## Prompt 2: Component System (Priority 2)

**Objective**: Implement component creation, dependency injection, and lifecycle management patterns

**Context**:
- Depends on Registry and base Component implementations from Prompt 1
- Factory interface is defined in `skeleton/internal/domain/component/factory.go`
- DependencyAwareComponent interface is defined in `skeleton/internal/domain/component/dependency_aware.go`
- LifecycleAwareComponent interface is defined in `skeleton/internal/domain/component/lifecycle_aware.go`
- All implementations go in `skeleton/internal/infrastructure/component/`

**Implementation Requirements**:

1. **Component Factory Implementation** (`skeleton/internal/infrastructure/component/factory_impl.go`)
   - Implement concrete `Factory` interface
   - Provide `NewFactory()` constructor that accepts registry interface dependency
   - Implement `Create(config ComponentConfig) (Component, error)` method
   - Support for different component types: Basic, Operation, Service, System, Application
   - Handle component configuration and property injection

2. **Dependency-Aware Component Implementation** (`skeleton/internal/infrastructure/component/dependency_aware_impl.go`)
   - Implement concrete `DependencyAwareComponent` interface
   - Provide `NewDependencyAwareComponent()` constructor that accepts base component and registry interfaces
   - Implement methods: `Dependencies()`, `AddDependency()`, `RemoveDependency()`, `HasDependency()`
   - Implement dependency resolution: `ResolveDependency()`, `ResolveDependencies()`
   - Support for circular dependency detection

3. **Lifecycle-Aware Component Implementation** (`skeleton/internal/infrastructure/component/lifecycle_aware_impl.go`)
   - Implement concrete `LifecycleAwareComponent` interface
   - Provide `NewLifecycleAwareComponent()` constructor that accepts base component interface
   - Implement methods: `State()`, `SetState()`, `OnStateChange()`
   - Support lifecycle states: Created, Initializing, Initialized, Active, Disposing, Disposed, Failed
   - Thread-safe state management with callbacks

**Acceptance Criteria**:
- Factory can create components from ComponentConfig
- Dependency resolution works with Registry integration
- Lifecycle state management is thread-safe
- State change callbacks are properly invoked
- Circular dependency detection prevents infinite loops
- **Testable through dependency injection** (all constructors accept interface dependencies)
- **External API compatibility maintained** (factory functions available through public API)
- **Simple, focused implementations** without over-engineering

**Documentation Deliverables**:
- Explanation of dependency injection patterns used in constructors
- Design decisions for circular dependency detection
- Notes on lifecycle state management and thread safety

**Files to Create**:
- `skeleton/internal/infrastructure/component/factory_impl.go`
- `skeleton/internal/infrastructure/component/dependency_aware_impl.go`
- `skeleton/internal/infrastructure/component/lifecycle_aware_impl.go`

---

## Prompt 3: Operation System (Priority 3)

**Objective**: Implement operation execution patterns for discrete units of work

**Context**:
- Depends on Component system from Prompts 1-2
- Operation interface is defined in `skeleton/internal/domain/operation/operation.go`
- Operations extend Component interface and add Execute method
- OperationFactory interface is defined for operation creation
- All implementations go in `skeleton/internal/infrastructure/operation/`

**Implementation Requirements**:

1. **Operation Base Implementation** (`skeleton/internal/infrastructure/operation/operation_impl.go`)
   - Implement base operation functionality
   - Provide `NewOperation()` constructor that accepts component interface dependency
   - Implement `Operation` interface extending `Component`
   - Implement `Execute(ctx context.Context, input Input) (Output, error)` method
   - Support for operation input/output handling
   - Use error constants defined in service domain: `ErrServiceStart`, `ErrServiceStop`, `ErrServiceNotFound`

2. **Operation Factory Implementation** (`skeleton/internal/infrastructure/operation/factory_impl.go`)
   - Implement concrete `OperationFactory` interface
   - Provide `NewOperationFactory()` constructor that accepts component factory interface dependency
   - Implement `Create(config ComponentConfig) (Component, error)` from Factory
   - Implement `CreateOperation(config OperationConfig) (Operation, error)`
   - Support for operation-specific configuration

**Acceptance Criteria**:
- Operations can be created from OperationConfig
- Execute method properly handles context, input, and output
- Error handling follows operation patterns
- Factory integration with component system
- Operations support component lifecycle
- **Testable through dependency injection** (constructors accept interface dependencies)
- **External API compatibility maintained** (operation creation available through public API)
- **Simple, focused implementations** without unnecessary abstraction

**Documentation Deliverables**:
- Explanation of operation execution patterns and dependency injection
- Design decisions for input/output handling
- Notes on error handling and factory integration

**Files to Create**:
- `skeleton/internal/infrastructure/operation/operation_impl.go`
- `skeleton/internal/infrastructure/operation/factory_impl.go`

---

## Prompt 4: Service System (Priority 4)

**Objective**: Implement service lifecycle management for ongoing functionality

**Context**:
- Depends on Component system from Prompts 1-2
- Service interface is defined in `skeleton/internal/domain/service/service.go`
- Services extend Component interface and add Start/Stop lifecycle
- ServiceFactory interface is defined for service creation
- Service status constants: Stopped, Starting, Running, Stopping, Failed
- All implementations go in `skeleton/internal/infrastructure/service/`

**Implementation Requirements**:

1. **Service Base Implementation** (`skeleton/internal/infrastructure/service/service_impl.go`)
   - Implement base service functionality
   - Provide `NewService()` constructor that accepts component interface dependency
   - Implement `Service` interface extending `Component`
   - Implement methods: `Start(ctx context.Context) error`, `Stop(ctx context.Context) error`, `Status() ServiceStatus`
   - Thread-safe status management
   - Support for service lifecycle state transitions

2. **Service Factory Implementation** (`skeleton/internal/infrastructure/service/factory_impl.go`)
   - Implement concrete `ServiceFactory` interface
   - Provide `NewServiceFactory()` constructor that accepts component factory interface dependency
   - Implement `Create(config ComponentConfig) (Component, error)` from Factory
   - Implement `CreateService(config ServiceConfig) (Service, error)`
   - Support for service-specific configuration

**Acceptance Criteria**:
- Services can be created from ServiceConfig
- Start/Stop lifecycle is properly managed
- Status transitions are thread-safe
- Error handling uses defined error codes: `ErrServiceStart`, `ErrServiceStop`, `ErrServiceNotFound`
- Factory integration with component system
- **Testable through dependency injection** (constructors accept interface dependencies)
- **External API compatibility maintained** (service management available through public API)
- **Simple, focused implementations** without over-complicating lifecycle management

**Documentation Deliverables**:
- Explanation of service lifecycle management and dependency injection patterns
- Design decisions for thread-safe status transitions
- Notes on error handling and factory integration

**Files to Create**:
- `skeleton/internal/infrastructure/service/service_impl.go`
- `skeleton/internal/infrastructure/service/factory_impl.go`

---

## Prompt 5: System Integration (Priority 5)

**Objective**: Implement system-level resource access and plugin management

**Context**:
- Depends on all previous implementations (Registry, Components, Operations, Services)
- System interface is defined in `skeleton/internal/domain/system/system.go`
- Plugin interfaces are defined in `skeleton/internal/domain/plugin/plugin.go`
- Configuration interface is defined in `skeleton/internal/domain/config/config.go`
- All implementations go in `skeleton/internal/infrastructure/`

**Implementation Requirements**:

1. **System Implementation** (`skeleton/internal/infrastructure/system/system_impl.go`)
   - Implement concrete `System` interface
   - Provide `NewSystem()` constructor that accepts all resource interface dependencies (registry, plugin manager, event bus, configuration, store)
   - Implement resource access methods: `Registry()`, `PluginManager()`, `EventBus()`, `Configuration()`, `Store()`
   - Implement system operations: `ExecuteOperation()`, `StartService()`, `StopService()`
   - Implement state methods: `IsRunning()`, `IsInitialized()`
   - Coordinate all system resources
   - Use error constants: `ErrSystemNotInitialized`, `ErrSystemNotStarted`, `ErrOperationNotFound`, `ErrOperationFailed`, `ErrServiceNotFound`, `ErrServiceStart`, `ErrServiceStop`

2. **Plugin Manager Implementation** (`skeleton/internal/infrastructure/plugin/plugin_impl.go`)
   - Implement concrete `PluginManager` interface
   - Provide `NewPluginManager()` constructor that accepts filesystem interface dependency
   - Implement methods: `Discover()`, `Load()`, `Unload()`, `ListPlugins()`, `GetPlugin()`
   - Support for plugin discovery and lifecycle management
   - Integration with Registry for plugin components
   - Use error constants: `ErrPluginNotFound`, `ErrPluginLoad`, `ErrPluginUnload`, `ErrPluginDiscovery`

3. **Configuration Implementation** (`skeleton/internal/infrastructure/config/config_impl.go`)
   - Implement concrete `Configuration` interface
   - Provide `NewConfiguration()` constructor that accepts configuration source interface dependencies
   - Implement getter methods: `GetString()`, `GetInt()`, `GetBool()`, `GetDuration()`, `GetObject()`
   - Implement default value methods: `GetStringDefault()`, `GetIntDefault()`, etc.
   - Implement `Exists()` method for key checking
   - Use error constants: `ErrConfigNotFound`, `ErrConfigWrongType`, `ErrConfigLoadFailed`

**Acceptance Criteria**:
- System provides access to all registered resources
- Plugin discovery and loading works correctly
- Configuration supports all data types and defaults
- Error handling uses defined error codes from domain interfaces
- System state management is consistent
- **Testable through dependency injection** (all constructors accept interface dependencies)
- **External API compatibility maintained** (system access available through public API)
- **Simple, focused implementations** without over-engineering system coordination

**Documentation Deliverables**:
- Explanation of system resource coordination and dependency injection patterns
- Design decisions for plugin discovery and management
- Notes on configuration handling and error management

**Files to Create**:
- `skeleton/internal/infrastructure/system/system_impl.go`
- `skeleton/internal/infrastructure/plugin/plugin_impl.go`
- `skeleton/internal/infrastructure/config/config_impl.go`

---

## Prompt 6: Supporting Infrastructure (Priority 6)

**Objective**: Implement event system, storage, and logging infrastructure

**Context**:
- These are supporting systems that enhance the framework capabilities
- Event interface is defined in `skeleton/internal/domain/event/event.go`
- Storage interfaces are defined in `skeleton/internal/domain/storage/`
- Logging interface is defined in `skeleton/internal/domain/logging/logging.go`
- All implementations go in `skeleton/internal/infrastructure/`

**Implementation Requirements**:

1. **Event System Implementation** (`skeleton/internal/infrastructure/event/event_impl.go`)
   - Implement concrete `EventBus` interface
   - Provide `NewEventBus()` constructor with minimal dependencies
   - Implement methods: `Publish()`, `Subscribe()`, `SubscribeAsync()`, `WaitAsync()`
   - Implement `Subscription` interface with `Cancel()` and `Topic()` methods
   - Support for synchronous and asynchronous event handling

2. **Storage System Implementation** (`skeleton/internal/infrastructure/storage/storage_impl.go`)
   - Implement concrete `MultiStore` interface
   - Provide `NewMultiStore()` constructor that accepts engine interface dependencies
   - Implement store management: `GetStore()`, `CreateStore()`, `DeleteStore()`, `ListStores()`, `StoreExists()`
   - Implement bulk operations: `CloseAll()`
   - Implement engine management: `RegisterEngine()`, `ListEngines()`, `GetEngine()`
   - Support for default engine configuration
   - Use string constants from `skeleton/internal/domain/storage/errors.go`: `ErrKeyNotFound`, `ErrStoreNotFound`, `ErrStoreClosed`, etc.

3. **Logging System Implementation** (`skeleton/internal/infrastructure/logging/logging_impl.go`)
   - Implement concrete `Logger` interface from `skeleton/internal/domain/logging/logging.go`
   - Provide `NewLogger()` constructor with minimal dependencies
   - Implement methods: `Debug()`, `Info()`, `Warn()`, `Error()`
   - Support for structured logging and log levels
   - Integration with system-wide logging patterns
   - Use error constants: `ErrLoggerNotAvailable`, `ErrInvalidLogLevel`

**Acceptance Criteria**:
- Event bus supports publish/subscribe patterns
- Storage system manages multiple stores and engines
- Logging provides structured output capabilities
- All implementations are thread-safe
- Error handling follows framework patterns (string constants for storage, string constants for others)
- **Testable through dependency injection** (constructors accept interface dependencies where applicable)
- **External API compatibility maintained** (supporting infrastructure available through public API)
- **Simple, focused implementations** without unnecessary complexity

**Documentation Deliverables**:
- Explanation of event system patterns and dependency injection
- Design decisions for storage engine management
- Notes on logging integration and error handling patterns

**Files to Create**:
- `skeleton/internal/infrastructure/event/event_impl.go`
- `skeleton/internal/infrastructure/storage/storage_impl.go`
- `skeleton/internal/infrastructure/logging/logging_impl.go`

---

## Implementation Notes

- **Clean Architecture**: All concrete implementations go in `skeleton/internal/infrastructure/`, never in `skeleton/internal/domain/`
- **Domain Interfaces**: Always implement the exact interfaces defined in `skeleton/internal/domain/`
- **Constructor Functions**: All `New*()` functions are infrastructure constructors, not domain interfaces
- **Dependency Interface Usage**: 
  - Check if dependencies already provide interfaces before creating new ones
  - Use existing interfaces from dependency packages when available
  - Accept interface dependencies in constructors for testability
  - Keep dependency injection simple and focused
- **External API Compatibility**: 
  - Maintain backward compatibility with existing public API in `skeleton/pkg/`
  - Implementation changes should be invisible to consumers
  - Ensure factory functions are available through public API
- **Testability Requirements**:
  - All implementations must be testable through dependency injection
  - Constructor functions should accept interface dependencies where applicable
  - Avoid direct instantiation of dependencies within implementations
  - Keep implementations simple to facilitate testing (separate testing guides exist)
- **Documentation Deliverables**:
  - Document dependency injection patterns used in constructors
  - Explain key design decisions and trade-offs
  - Note interface compatibility and public API preservation
  - Keep documentation brief and focused
- **Error Handling**: 
  - Storage domain uses string constants (`const ErrKeyNotFound = "storage.key_not_found"`)
  - Other domains use string constants (`const ErrServiceStart = "service.start_failed"`)
  - Use existing error constants from domain interfaces, don't create new error systems
- **Thread Safety**: Implementations should be thread-safe where the domain interface will be used concurrently
- **Dependency Order**: Each prompt builds on the previous implementations in dependency order
- **Public API**: The public API in `skeleton/pkg/` re-exports these internal implementations through factory functions
- **Simplicity**: Avoid over-engineering or unnecessary complexity in implementations 