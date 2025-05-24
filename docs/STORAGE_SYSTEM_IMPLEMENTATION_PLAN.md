# Implementation Plan: Storage and System Components Redesign

## 1. Overview

### 1.1 Purpose and Goals
This document outlines the implementation plan for redesigning the storage and system components within our new component architecture. The primary goals of this redesign are:

- Redesign the existing storage and system functionality to align with our component system
- Create clean, well-defined interfaces with clear responsibilities
- Support multiple storage backends through a common abstraction
- Implement an event-driven architecture for better decoupling
- Achieve high test coverage and maintainability
- Allow for future extensions without breaking changes

### 1.2 Architectural Principles
The redesign will follow these architectural principles:

1. **Interface Segregation**: Keep interfaces focused on specific capabilities
2. **Composition Over Inheritance**: Use composition to build complex components
3. **Dependency Injection**: Make dependencies explicit and injectable
4. **Event-Driven Communication**: Use events for loose coupling between components
5. **Clean Abstraction**: Hide implementation details behind well-defined interfaces
6. **Immutable State**: Prefer immutable state and explicit state transitions
7. **Testability**: Design for comprehensive testing

### 1.3 Timeline and Milestones

| Milestone | Description | Timeline |
|-----------|-------------|----------|
| **M1** | Core interfaces and domain components | Week 1-2 |
| **M2** | Basic implementations (in-memory, file) | Week 3-4 |
| **M3** | Advanced implementations (IAVL, LevelDB) | Week 5-6 |
| **M4** | System service integration | Week 7-8 |
| **M5** | Testing and optimization | Week 9-10 |
| **M6** | Documentation and examples | Week 11-12 |

## 2. Core Abstractions

### 2.1 Store Interface

```go
// Store defines the core storage operations that all backends must implement.
type Store interface {
    // Basic CRUD operations

    // Get retrieves the value for the given key.
    // Returns ErrKeyNotFound if the key doesn't exist.
    Get(key []byte) ([]byte, error)

    // Set stores the value for the given key.
    // Overwrites any existing value.
    Set(key, value []byte) error

    // Delete removes the key-value pair.
    // It's idempotent (no error if key doesn't exist).
    Delete(key []byte) error

    // Has checks if a key exists.
    // More efficient than Get for existence checks.
    Has(key []byte) (bool, error)

    // Iteration over all key-value pairs

    // Iterate calls fn for each key-value pair, stops if fn returns false.
    // The key and value byte slices must not be modified by fn.
    Iterate(fn func(key, value []byte) bool) error

    // Resource cleanup

    // Close releases resources and makes the store unusable.
    Close() error

    // Store metadata

    // Name returns the store identifier.
    Name() string

    // Path returns the storage path/location.
    Path() string
}
```

### 2.2 Range Query Extension

```go
// RangeQueryable interface for stores that support range queries.
type RangeQueryable interface {
    // IterateRange iterates over keys in the specified range.
    // If ascending is true, iteration is from start to end.
    // If ascending is false, iteration is from end to start.
    // The fn callback works the same as in Store.Iterate.
    IterateRange(start, end []byte, ascending bool, fn func(key, value []byte) bool) error

    // SupportsRangeQueries returns true if this store supports range queries.
    SupportsRangeQueries() bool
}
```

### 2.3 Transaction Extension

```go
// Transactional interface for stores that support transactions.
type Transactional interface {
    // BeginTx starts a new transaction.
    // Returns an error if transactions are not supported or cannot be started.
    BeginTx() (Transaction, error)

    // SupportsTransactions returns true if this store supports transactions.
    SupportsTransactions() bool
}

// Transaction represents an atomic set of operations.
// All operations within a transaction either succeed as a group or fail.
type Transaction interface {
    // Embed the Store interface - transactions support all store operations
    Store

    // Commit makes all changes permanent.
    // Returns an error if the transaction cannot be committed.
    Commit() error

    // Rollback discards all changes.
    // Returns an error if the transaction cannot be rolled back.
    Rollback() error

    // IsActive returns true if transaction is still active.
    // A transaction is active until Commit or Rollback is called.
    IsActive() bool
}
```

### 2.4 Versioning Extension

```go
// Versioned interface for stores that support versioning/snapshots.
type Versioned interface {
    // SaveVersion creates a new immutable version of the store.
    // Returns the version number and a hash that uniquely identifies the version.
    // The hash can be used for integrity verification.
    SaveVersion() (version int64, hash []byte, err error)

    // LoadVersion loads a specific version of the store.
    // Returns ErrVersionNotFound if the version doesn't exist.
    LoadVersion(version int64) error

    // ListVersions returns all available versions.
    // Returns an empty slice if no versions are available.
    ListVersions() []int64

    // CurrentVersion returns the current version number.
    // Returns 0 if no versions have been saved.
    CurrentVersion() int64

    // SupportsVersioning returns true if this store supports versioning.
    SupportsVersioning() bool
}
```

### 2.5 MultiStore Interface

```go
// MultiStore manages multiple named stores.
// It provides a central registry for all stores and their engines.
type MultiStore interface {
    // Store management

    // GetStore retrieves a store by name.
    // Returns ErrStoreNotFound if the store doesn't exist.
    GetStore(name string) (Store, error)

    // CreateStore creates a new store with the given name and engine.
    // The config parameter contains engine-specific configuration options.
    // Returns ErrStoreExists if the store already exists.
    CreateStore(name, engine string, config Config) error

    // DeleteStore removes a store by name.
    // Returns ErrStoreNotFound if the store doesn't exist.
    DeleteStore(name string) error

    // ListStores returns the names of all stores.
    // Returns an empty slice if no stores exist.
    ListStores() []string

    // StoreExists checks if a store exists.
    StoreExists(name string) bool

    // Bulk operations

    // CloseAll closes all stores.
    // This should be called when shutting down the application.
    CloseAll() error

    // Configuration

    // SetDefaultEngine sets the default engine to use when no engine is specified.
    SetDefaultEngine(engine string)

    // GetDefaultEngine returns the current default engine.
    GetDefaultEngine() string

    // Engine management

    // RegisterEngine registers a new engine.
    // Returns an error if an engine with that name is already registered.
    RegisterEngine(engine Engine) error

    // ListEngines returns the names of all registered engines.
    ListEngines() []string

    // GetEngine retrieves an engine by name.
    // Returns ErrEngineNotFound if the engine doesn't exist.
    GetEngine(name string) (Engine, error)
}
```

### 2.6 Storage Engine Interface

```go
// Engine interface for storage backend implementations.
// Each concrete storage implementation (memory, file, etc.) must provide
// an Engine implementation that can create and open stores.
type Engine interface {
    // Name returns the engine identifier (e.g., "memory", "leveldb", "file").
    // This is used for engine registration and lookup.
    Name() string

    // Create creates a new store instance with the given name and path.
    // The config parameter contains engine-specific configuration options.
    Create(name, path string, config Config) (Store, error)

    // Open opens an existing store with the given name and path.
    // Returns ErrStoreNotFound if the store doesn't exist.
    Open(name, path string) (Store, error)

    // Capabilities returns what features this engine supports.
    // This allows clients to check for optional features before using them.
    Capabilities() Capabilities
}

// Capabilities describes what features an engine supports.
// This allows clients to check for optional capabilities before using them.
type Capabilities struct {
    // Transactions indicates if the engine supports atomic transactions.
    Transactions bool

    // Versioning indicates if the engine supports versioning/snapshots.
    Versioning bool

    // RangeQueries indicates if the engine supports efficient range queries.
    RangeQueries bool

    // Persistence indicates if the engine persists data to disk.
    Persistence bool

    // Compression indicates if the engine supports data compression.
    Compression bool
}
```

### 2.7 Configuration

```go
// Config type alias for configuration options
type Config map[string]interface{}

// MultiStoreConfig defines configuration for MultiStore
type MultiStoreConfig struct {
    // RootPath is the base directory for all stores managed by this MultiStore.
    RootPath string `json:"rootPath"`

    // DefaultEngine is the engine to use when no engine is specified.
    DefaultEngine string `json:"defaultEngine"`

    // EngineConfigs contains engine-specific configuration.
    // The keys are engine names, and the values are engine-specific configs.
    EngineConfigs map[string]Config `json:"engineConfigs"`
}
```

### 2.8 System Service Integration

```go
// SystemService serves as the central coordinating service
type SystemService interface {
    service.Service
    
    // Core components access
    Registry() component.Registry
    EventBus() event.EventBus
    Configuration() config.Configuration
    Store() storage.MultiStore
    
    // Operations
    ExecuteOperation(ctx component.Context, operationID string, input interface{}) (interface{}, error)
    StartService(ctx component.Context, serviceID string) error
    StopService(ctx component.Context, serviceID string) error
}
```

## 3. Implementation Phases

### 3.1 Phase 1: Core Interfaces and Abstractions (Week 1-2)

#### Tasks:
1. Define domain interfaces in `internal/domain/storage/`:
   - Store interface (core storage operations)
   - Extension interfaces (RangeQueryable, Transactional, Versioned)
   - MultiStore interface
   - Engine interface
   - Configuration types

2. Define domain events in `internal/domain/storage/`:
   - Storage event topics (TopicStoreCreated, TopicStoreDeleted, etc.)
   - Helper functions for creating event payloads

3. Define system service interface in `internal/domain/system/`:
   - SystemService interface
   - System-related events

4. Define error types in `internal/domain/storage/`:
   - Common error constants
   - Error wrapping functions

#### Deliverables:
- Complete interface definitions with documentation
- Error constants and helper functions
- Event topics and payload creation utilities
- Configuration types for storage components

### 3.2 Phase 2: Basic Implementations (Week 3-4)

#### Tasks:
1. Implement in-memory storage in `internal/infrastructure/storage/memory/`:
   - Memory Engine implementation
   - Memory Store implementation
   - Transaction support
   - Versioning support
   - Range query support

2. Implement basic store tests:
   - Core CRUD operations
   - Transaction tests
   - Versioning tests
   - Range query tests

3. ~~Implement file-based storage engine~~ (Deferred to future phase)

4. ~~Implement basic MultiStore~~ (Deferred to future phase)

#### Deliverables:
- Working in-memory implementation
- Comprehensive test suite for the in-memory implementation

### 3.3 Phase 3: Additional Storage Engines (Week 5-6)

#### Tasks:
1. Implement file-based storage engine in `internal/infrastructure/storage/file/`:
   - FileEngine implementation
   - FileStore implementation
   - JSON serialization for storage
   - Directory-based versioning

2. Implement basic MultiStore in `internal/infrastructure/storage/`:
   - DefaultMultiStore implementation
   - Store registry and management
   - Event publication
   - Engine registration and configuration

3. ~~Implement LevelDB storage engine~~ (Deferred to future phase)
   
4. ~~Implement Badger storage engine~~ (Deferred to future phase)

#### Deliverables:
- Working file-based implementation
- MultiStore implementation for managing multiple stores
- Test suite for file storage and MultiStore

### 3.4 Phase 4: System Service Integration (Week 7-8)

#### Tasks:
1. Implement core SystemService in `internal/infrastructure/system/`:
   - DefaultSystemService implementation
   - Component registration
   - Service lifecycle management
   - Operation execution

2. Integrate MultiStore with SystemService:
   - Store initialization during system startup
   - Store management APIs
   - Event propagation between storage and system

3. Implement operation routing:
   - Operation discovery via component registry
   - Input/output conversion for operations
   - Error handling and reporting

4. Implement service management:
   - Service discovery via component registry
   - Start/stop handling
   - Feature flag checking

#### Deliverables:
- Complete SystemService implementation
- Integration between storage and system components
- Operation execution pipeline
- Service management functionality

### 3.5 Phase 5: Advanced Storage Engines (Week 9-10)

#### Tasks:
1. Implement LevelDB storage engine in `internal/infrastructure/storage/leveldb/`:
   - LevelDBEngine implementation
   - LevelDBStore implementation
   - Transaction support
   - Versioning support

2. Implement Badger storage engine in `internal/infrastructure/storage/badger/`:
   - BadgerEngine implementation
   - BadgerStore implementation
   - Transaction support using Badger's native transactions
   - TTL-based versioning

3. Implement engine discovery and registration:
   - Dynamic engine loading
   - Configuration-based engine selection
   - Engine capability detection

4. Optimize critical paths:
   - Profile and optimize store operations
   - Improve transaction performance
   - Optimize memory usage

#### Deliverables:
- LevelDB storage engine implementation
- Badger storage engine implementation
- Engine discovery and registration system
- Performance optimizations

### 3.6 Phase 6: Documentation and Examples (Week 11-12)

#### Tasks:
1. Create comprehensive API documentation:
   - Interface documentation
   - Implementation notes
   - Error handling guidelines

2. Create usage examples:
   - Basic storage operations
   - Transaction usage
   - Versioning usage
   - MultiStore configuration

3. Create migration guide:
   - From legacy storage to new storage
   - API differences
   - Breaking changes and workarounds

4. Create configuration guide:
   - Engine selection
   - Performance tuning
   - Security considerations

#### Deliverables:
- Complete API documentation
- Usage examples
- Migration guide
- Configuration guide

## 4. Detailed Component Design

### 4.1 Storage Domain Components

#### 4.1.1 Store Interface
The `Store` interface defines the core storage operations that all backends must implement:

```go
// Store defines the core storage operations that all backends must implement.
type Store interface {
    // Basic CRUD operations
    Get(key []byte) ([]byte, error)
    Set(key, value []byte) error
    Delete(key []byte) error
    Has(key []byte) (bool, error)
    
    // Iteration over all key-value pairs
    Iterate(fn func(key, value []byte) bool) error
    
    // Resource cleanup
    Close() error
    
    // Store metadata
    Name() string
    Path() string
}
```

#### 4.1.2 Extension Interfaces
The storage system provides extension interfaces for optional capabilities:

```go
// RangeQueryable interface for stores that support range queries.
type RangeQueryable interface {
    IterateRange(start, end []byte, ascending bool, fn func(key, value []byte) bool) error
    SupportsRangeQueries() bool
}

// Transactional interface for stores that support transactions.
type Transactional interface {
    BeginTx() (Transaction, error)
    SupportsTransactions() bool
}

// Versioned interface for stores that support versioning/snapshots.
type Versioned interface {
    SaveVersion() (version int64, hash []byte, err error)
    LoadVersion(version int64) error
    ListVersions() []int64
    CurrentVersion() int64
    SupportsVersioning() bool
}
```

#### 4.1.3 MultiStore Interface
The `MultiStore` interface manages multiple named stores and their engines:

```go
// MultiStore manages multiple named stores.
type MultiStore interface {
    // Store management
    GetStore(name string) (Store, error)
    CreateStore(name, engine string, config Config) error
    DeleteStore(name string) error
    ListStores() []string
    StoreExists(name string) bool
    
    // Bulk operations
    CloseAll() error
    
    // Configuration
    SetDefaultEngine(engine string)
    GetDefaultEngine() string
    
    // Engine management
    RegisterEngine(engine Engine) error
    ListEngines() []string
    GetEngine(name string) (Engine, error)
}
```

#### 4.1.4 Engine Interface
The `Engine` interface represents a storage backend implementation:

```go
// Engine interface for storage backend implementations.
type Engine interface {
    // Name returns the engine identifier (e.g., "memory", "leveldb", "file").
    // This is used for engine registration and lookup.
    Name() string

    // Create creates a new store instance with the given name and path.
    // The config parameter contains engine-specific configuration options.
    Create(name, path string, config Config) (Store, error)

    // Open opens an existing store with the given name and path.
    // Returns ErrStoreNotFound if the store doesn't exist.
    Open(name, path string) (Store, error)

    // Capabilities returns what features this engine supports.
    // This allows clients to check for optional features before using them.
    Capabilities() Capabilities
}

// Capabilities describes what features an engine supports.
// This allows clients to check for optional capabilities before using them.
type Capabilities struct {
    // Transactions indicates if the engine supports atomic transactions.
    Transactions bool

    // Versioning indicates if the engine supports versioning/snapshots.
    Versioning bool

    // RangeQueries indicates if the engine supports efficient range queries.
    RangeQueries bool

    // Persistence indicates if the engine persists data to disk.
    Persistence bool

    // Compression indicates if the engine supports data compression.
    Compression bool
}
```

#### 4.1.5 Transaction Interface
The `Transaction` interface represents an atomic set of operations:

```go
// Transaction represents an atomic set of operations.
type Transaction interface {
    // Embed the Store interface
    Store
    
    // Transaction management
    Commit() error
    Rollback() error
    IsActive() bool
}
```

### 4.2 Storage Infrastructure Components

#### 4.2.1 Memory Store Implementation
The memory store provides an in-memory implementation of the store interfaces:

```go
// Store implements storage.Store for in-memory storage.
type Store struct {
    name      string
    path      string
    options   Options
    data      map[string][]byte
    closed    bool
    mutex     sync.RWMutex
    versions  map[int64]map[string][]byte
    currVer   int64
    totalSize int64
    logger    logging.Logger
}
```

The memory store implements:
- The base `Store` interface
- `RangeQueryable` for range queries
- `Transactional` for transaction support
- `Versioned` for versioning/snapshot support

#### 4.2.2 Memory Engine Implementation
The memory engine creates and manages memory stores:

```go
// Engine implements storage.Engine for in-memory storage.
type Engine struct {
    logger logging.Logger
    mutex  sync.RWMutex
    stores map[string]*Store
}
```

The memory engine provides:
- Store creation and opening
- Capability reporting
- Configuration parsing

#### 4.2.3 Memory Transaction Implementation
Transactions in the memory store provide atomic operations:

```go
// Transaction implements storage.Transaction for in-memory storage.
type Transaction struct {
    store    *Store
    changes  map[string][]byte // Pending changes
    deletes  map[string]bool   // Keys to delete
    active   bool
    readOnly bool
}
```

The transaction implementation:
- Tracks pending changes separately from the main store
- Applies changes atomically on commit
- Discards changes on rollback

### 4.3 System Integration Components

#### 4.3.1 System Service
The `SystemService` interface provides the central coordination service:

```go
// SystemService serves as the central coordinating service
type SystemService interface {
    service.Service
    
    // Core components access
    Registry() component.Registry
    EventBus() event.EventBus
    Configuration() config.Configuration
    Store() storage.MultiStore
    
    // Operations
    ExecuteOperation(ctx component.Context, operationID string, input interface{}) (interface{}, error)
    StartService(ctx component.Context, serviceID string) error
    StopService(ctx component.Context, serviceID string) error
}
```

The system service implementation:
- Initializes and manages the registry, event bus, and storage
- Handles the system lifecycle
- Executes operations
- Manages services

#### 4.3.2 Configuration System
The configuration handling includes:

```go
// MultiStoreConfig defines configuration for MultiStore
type MultiStoreConfig struct {
    RootPath       string                   `json:"rootPath"`
    DefaultEngine  string                   `json:"defaultEngine"`
    EngineConfigs  map[string]Config        `json:"engineConfigs"`
}

// Common configuration keys
const (
    ConfigCacheSize     = "cache_size"
    ConfigCompression   = "compression"
    ConfigMaxVersions   = "max_versions"
    ConfigSyncWrites    = "sync_writes"
    ConfigFileFormat    = "format"
    ConfigReadOnly      = "read_only"
)
```

#### 4.3.3 Event Propagation
Events are published for important storage operations:

```go
// Storage event topics
const (
    TopicStoreCreated        = "store.created"
    TopicStoreDeleted        = "store.deleted"
    TopicStoreClosed         = "store.closed"
    TopicVersionSaved        = "store.version.saved"
    TopicVersionLoaded       = "store.version.loaded"
    TopicTransactionBegin    = "store.transaction.begin"
    TopicTransactionCommit   = "store.transaction.commit"
    TopicTransactionRollback = "store.transaction.rollback"
)
```

The system service handles these events and propagates them through the system event bus.

### 4.4 Extension Points

The storage system provides several extension points:

1. **New Storage Engines**: Implement the `Engine` interface to provide new storage backends.
2. **Storage Capabilities**: Implement optional interfaces for additional capabilities.
3. **Event Handlers**: Subscribe to storage events to respond to storage operations.
4. **Configuration Options**: Provide engine-specific configuration options.

## 5. Migration Strategy

### 5.1 From Legacy Storage to New Storage

The migration from legacy storage to the new storage system will involve:

1. **Adapter Pattern**:
   - Create adapter implementations to bridge between old and new interfaces
   - Ensure backward compatibility for existing code

2. **Data Migration**:
   - Develop utilities to convert data from old format to new format
   - Validate data integrity during migration

3. **Gradual Replacement**:
   - Replace components one at a time
   - Use feature flags to control migration path

4. **Parallel Operation**:
   - Run old and new systems in parallel during transition
   - Compare results to ensure consistency

### 5.2 Backward Compatibility

To maintain backward compatibility:

1. **Interface Adapters**:
   - Create adapters for legacy interfaces
   - Map old method signatures to new ones

2. **Configuration Mapping**:
   - Create mapping between old and new configuration formats
   - Provide sensible defaults for new configuration options

3. **Event Bridges**:
   - Translate between old and new event systems
   - Ensure events propagate correctly across systems

## 6. Testing Approach

### 6.1 Unit Testing

Each component will have unit tests covering:

1. **Store Operations**:
   - CRUD operations
   - Range queries
   - Iteration
   - Error conditions

2. **Transaction Handling**:
   - Transaction creation
   - Commit operations
   - Rollback operations
   - Concurrent transactions

3. **Versioning**:
   - Version creation
   - Version loading
   - Version listing
   - Version cleanup

4. **Engine Management**:
   - Store creation and opening
   - Configuration parsing
   - Capability reporting

### 6.2 Integration Testing

Integration tests will verify:

1. **System Service with Storage**:
   - Storage initialization
   - Store management
   - Event propagation

2. **MultiStore Management**:
   - Engine registration
   - Store creation and retrieval
   - Store deletion

3. **Cross-Component Interaction**:
   - Event handling
   - Error propagation
   - Configuration application

### 6.3 Performance Testing

Performance benchmarks will measure:

1. **Throughput**:
   - Operations per second
   - Batch operation performance

2. **Latency**:
   - Operation response time
   - Transaction commit time

3. **Resource Usage**:
   - Memory consumption
   - Disk usage
   - CPU utilization

## 7. Documentation

### 7.1 API Documentation

Comprehensive API documentation will include:

1. **Interface Documentation**:
   - Purpose and responsibility
   - Method descriptions
   - Parameter details
   - Return values and errors

2. **Implementation Notes**:
   - Design decisions
   - Implementation details
   - Usage considerations
   - Performance characteristics

3. **Error Handling**:
   - Error types and codes
   - Error recovery strategies
   - Best practices

### 7.2 Usage Examples

Examples will demonstrate:

1. **Basic Usage**:
   - Creating and using Store
   - Basic CRUD operations
   - Versioning operations

2. **Advanced Usage**:
   - Transactions
   - Custom store types
   - Event handling
   - Performance optimization

3. **Integration Examples**:
   - Using storage with SystemService
   - Component integration
   - Plugin integration

### 7.3 Configuration Guide

Configuration documentation will cover:

1. **Engine Selection**:
   - Available engines
   - Pros and cons
   - Configuration options

2. **Performance Tuning**:
   - Memory usage
   - Disk usage
   - Throughput optimization
   - Latency optimization

3. **Security Configuration**:
   - Access control
   - Encryption
   - Secure deletion
   - Audit logging

4. **Advanced Features**:
   - Custom serialization
   - Compression
   - Replication
   - Backup/restore 