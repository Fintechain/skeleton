# Minimal Storage System Design Specification

## 1. Overview

### 1.1 Purpose
This document provides a complete specification for implementing a simplified storage system that replaces the overly complex multi-layered design with a clean, minimal interface that's easy to implement and extend.

### 1.2 Design Principles
- **Simplicity First**: Core interface has minimal required methods
- **Optional Features**: Advanced features through interface composition
- **Easy Implementation**: Backends only implement what they support
- **Standard Go Patterns**: Use standard Go error handling and interfaces
- **Plugin Architecture**: Simple engine registration system

### 1.3 Architecture Overview
```
┌─────────────────────────────────────────────────────────┐
│                   Application Layer                      │
└───────────────────────┬─────────────────────────────────┘
                        │
┌───────────────────────▼─────────────────────────────────┐
│                   MultiStore                            │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────┐ │
│  │   Store A   │ │   Store B   │ │      Store C        │ │
│  │ (memory)    │ │ (leveldb)   │ │ (file+versioned)    │ │
│  └─────────────┘ └─────────────┘ └─────────────────────┘ │
└─────────────────────────────────────────────────────────┘
┌─────────────────────────────────────────────────────────┐
│                Engine Factory                           │
│  ┌─────────────┐ ┌─────────────┐ ┌─────────────────────┐ │
│  │MemoryEngine │ │LevelDBEngine│ │    FileEngine       │ │
│  └─────────────┘ └─────────────┘ └─────────────────────┘ │
└─────────────────────────────────────────────────────────┘
```

## 2. Core Interfaces

All interfaces are defined in the domain layer (`internal/domain/storage/`) following the project's domain-driven design principles.

### 2.1 Store Interface
The primary storage interface that all backends must implement.

```go
package storage

import "io"

// Store defines the core storage operations that all backends must implement
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

**Implementation Notes:**
- All concrete implementations go in `internal/infrastructure/storage/`
- `Get` returns `ErrKeyNotFound` if key doesn't exist
- `Set` overwrites existing values
- `Delete` is idempotent (no error if key doesn't exist)
- `Has` is more efficient than `Get` for existence checks
- `Iterate` calls `fn` for each key-value pair, stops if `fn` returns false
- `Close` releases resources and makes store unusable
- `Name` returns the store identifier
- `Path` returns the storage path/location

### 2.2 Optional Interfaces

#### 2.2.1 Transactional Interface
For stores that support atomic transactions.

```go
// Transactional interface for stores that support transactions
type Transactional interface {
    // BeginTx starts a new transaction
    BeginTx() (Transaction, error)
    
    // SupportsTransactions returns true if this store supports transactions
    SupportsTransactions() bool
}

// Transaction represents an atomic set of operations
type Transaction interface {
    Store // transactions support all store operations
    
    // Commit makes all changes permanent
    Commit() error
    
    // Rollback discards all changes
    Rollback() error
    
    // IsActive returns true if transaction is still active
    IsActive() bool
}
```

#### 2.2.2 Versioned Interface
For stores that support versioning/snapshots.

```go
// Versioned interface for stores that support versioning
type Versioned interface {
    // SaveVersion creates a new immutable version of the store
    SaveVersion() (version int64, hash []byte, error)
    
    // LoadVersion loads a specific version of the store
    LoadVersion(version int64) error
    
    // ListVersions returns all available versions
    ListVersions() []int64
    
    // CurrentVersion returns the current version number
    CurrentVersion() int64
    
    // SupportsVersioning returns true if this store supports versioning
    SupportsVersioning() bool
}
```

#### 2.2.3 RangeQueryable Interface
For stores that support efficient range queries.

```go
// RangeQueryable interface for stores that support range queries
type RangeQueryable interface {
    // IterateRange iterates over keys in the specified range
    IterateRange(start, end []byte, ascending bool, fn func(key, value []byte) bool) error
    
    // SupportsRangeQueries returns true if this store supports range queries
    SupportsRangeQueries() bool
}
```

### 2.3 Engine Interface
For implementing storage backend plugins.

```go
// Engine interface for storage backend implementations
type Engine interface {
    // Name returns the engine identifier (e.g., "memory", "leveldb", "file")
    Name() string
    
    // Create creates a new store instance
    Create(name, path string, config Config) (Store, error)
    
    // Open opens an existing store
    Open(name, path string) (Store, error)
    
    // Capabilities returns what this engine supports
    Capabilities() Capabilities
}

// Capabilities describes what features an engine supports
type Capabilities struct {
    Transactions  bool
    Versioning    bool
    RangeQueries  bool
    Persistence   bool
    Compression   bool
}

// Config holds configuration parameters for store creation
type Config map[string]interface{}
```

## 3. MultiStore Implementation

### 3.1 MultiStore Interface

```go
// MultiStore manages multiple named stores
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

### 3.2 MultiStore Implementation

The MultiStore interface is defined in the domain layer, but implemented in the infrastructure layer:

**Domain Layer** (`internal/domain/storage/multistore.go`):
```go
// MultiStore manages multiple named stores
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

**Infrastructure Layer** (`internal/infrastructure/storage/multistore_impl.go`):
```go
// DefaultMultiStore implements the domain MultiStore interface
type DefaultMultiStore struct {
    stores        map[string]storage.Store
    engines       map[string]storage.Engine
    defaultEngine string
    rootPath      string
    mutex         sync.RWMutex
}

// NewMultiStore creates a new MultiStore instance
func NewMultiStore(rootPath string) storage.MultiStore {
    ms := &DefaultMultiStore{
        stores:        make(map[string]storage.Store),
        engines:       make(map[string]storage.Engine),
        defaultEngine: "memory",
        rootPath:      rootPath,
    }
    
    // Register built-in engines from infrastructure layer
    ms.RegisterEngine(memory.NewEngine())
    ms.RegisterEngine(file.NewEngine())
    
    return ms
}
```

## 4. Error Handling

### 4.1 Standard Errors

```go
package storage

import "errors"

// Standard storage errors
var (
    ErrKeyNotFound    = errors.New("key not found")
    ErrStoreNotFound  = errors.New("store not found")
    ErrStoreClosed    = errors.New("store is closed")
    ErrStoreExists    = errors.New("store already exists")
    ErrEngineNotFound = errors.New("engine not found")
    ErrTxNotActive    = errors.New("transaction not active")
    ErrTxReadOnly     = errors.New("transaction is read-only")
    ErrVersionNotFound = errors.New("version not found")
    ErrInvalidConfig  = errors.New("invalid configuration")
)

// Error wrapping helpers
func WrapError(err error, context string) error {
    if err == nil {
        return nil
    }
    return fmt.Errorf("%s: %w", context, err)
}

func IsKeyNotFound(err error) bool {
    return errors.Is(err, ErrKeyNotFound)
}

func IsStoreNotFound(err error) bool {
    return errors.Is(err, ErrStoreNotFound)
}
```

## 5. Built-in Engine Implementations

### 5.1 Memory Engine

**Infrastructure Layer** (`internal/infrastructure/storage/memory/engine.go`):
```go
package memory

import (
    "github.com/fintechain/skeleton/internal/domain/storage"
)

// Engine implements storage.Engine for in-memory storage
type Engine struct{}

func NewEngine() storage.Engine {
    return &Engine{}
}

func (e *Engine) Name() string {
    return "memory"
}

func (e *Engine) Capabilities() storage.Capabilities {
    return storage.Capabilities{
        Transactions:  true,
        Versioning:    true,
        RangeQueries:  true,
        Persistence:   false,
        Compression:   false,
    }
}

func (e *Engine) Create(name, path string, config storage.Config) (storage.Store, error) {
    return NewStore(name, path), nil
}

func (e *Engine) Open(name, path string) (storage.Store, error) {
    // Memory stores can't be "opened" from disk
    return nil, storage.WrapError(storage.ErrStoreNotFound, "memory stores are not persistent")
}
```

### 5.2 File Engine

**Infrastructure Layer** (`internal/infrastructure/storage/file/engine.go`):
```go
package file

import (
    "github.com/fintechain/skeleton/internal/domain/storage"
)

// Engine implements storage.Engine for file-based storage
type Engine struct{}

func NewEngine() storage.Engine {
    return &Engine{}
}

func (e *Engine) Name() string {
    return "file"
}

func (e *Engine) Capabilities() storage.Capabilities {
    return storage.Capabilities{
        Transactions:  false,
        Versioning:    true,
        RangeQueries:  false,
        Persistence:   true,
        Compression:   false,
    }
}

func (e *Engine) Create(name, path string, config storage.Config) (storage.Store, error) {
    return NewStore(name, path)
}

func (e *Engine) Open(name, path string) (storage.Store, error) {
    return OpenStore(name, path)
}
```

## 6. Implementation Guidelines

### 6.1 Store Implementation Pattern

**Infrastructure Layer** (`internal/infrastructure/storage/memory/store.go`):
```go
package memory

import (
    "sync"
    "github.com/fintechain/skeleton/internal/domain/storage"
)

// Store implements storage.Store for in-memory storage
type Store struct {
    name     string
    path     string
    data     map[string][]byte
    closed   bool
    mutex    sync.RWMutex
    
    // Optional features
    txSupport      bool
    versionSupport bool
}

func NewStore(name, path string) storage.Store {
    return &Store{
        name:           name,
        path:           path,
        data:           make(map[string][]byte),
        closed:         false,
        txSupport:      true,
        versionSupport: true,
    }
}

func (s *Store) Get(key []byte) ([]byte, error) {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    
    if s.closed {
        return nil, storage.ErrStoreClosed
    }
    
    value, exists := s.data[string(key)]
    if !exists {
        return nil, storage.ErrKeyNotFound
    }
    
    // Return copy to prevent external modification
    result := make([]byte, len(value))
    copy(result, value)
    return result, nil
}

func (s *Store) Set(key, value []byte) error {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    if s.closed {
        return storage.ErrStoreClosed
    }
    
    // Store copy to prevent external modification
    valueCopy := make([]byte, len(value))
    copy(valueCopy, value)
    s.data[string(key)] = valueCopy
    
    return nil
}

func (s *Store) Name() string {
    return s.name
}

func (s *Store) Path() string {
    return s.path
}
```

### 6.2 Transaction Implementation Pattern

**Infrastructure Layer** (`internal/infrastructure/storage/memory/transaction.go`):
```go
package memory

import (
    "github.com/fintechain/skeleton/internal/domain/storage"
)

type Transaction struct {
    store    *Store
    changes  map[string][]byte
    deletes  map[string]bool
    active   bool
    readOnly bool
}

func (s *Store) BeginTx() (storage.Transaction, error) {
    if !s.txSupport {
        return nil, storage.WrapError(storage.ErrTxNotActive, "transactions not supported")
    }
    
    return &Transaction{
        store:   s,
        changes: make(map[string][]byte),
        deletes: make(map[string]bool),
        active:  true,
    }, nil
}

func (tx *Transaction) Commit() error {
    if !tx.active {
        return storage.ErrTxNotActive
    }
    
    tx.store.mutex.Lock()
    defer tx.store.mutex.Unlock()
    
    // Apply all changes atomically
    for key, value := range tx.changes {
        tx.store.data[key] = value
    }
    
    for key := range tx.deletes {
        delete(tx.store.data, key)
    }
    
    tx.active = false
    return nil
}
```

### 6.3 Version Implementation Pattern

```go
type VersionedStore struct {
    ExampleStore
    versions map[int64]map[string][]byte
    currentVersion int64
    maxVersions int
}

func (s *VersionedStore) SaveVersion() (int64, []byte, error) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    
    if s.closed {
        return 0, nil, ErrStoreClosed
    }
    
    // Create snapshot
    s.currentVersion++
    snapshot := make(map[string][]byte)
    for k, v := range s.data {
        snapshot[k] = append([]byte(nil), v...)
    }
    
    s.versions[s.currentVersion] = snapshot
    
    // Cleanup old versions if needed
    s.cleanupOldVersions()
    
    // Calculate hash
    hash := s.calculateHash(snapshot)
    
    return s.currentVersion, hash, nil
}
```

## 7. Configuration

### 7.1 MultiStore Configuration

```go
// MultiStoreConfig defines configuration for MultiStore
type MultiStoreConfig struct {
    RootPath      string            `json:"rootPath"`
    DefaultEngine string            `json:"defaultEngine"`
    EngineConfigs map[string]Config `json:"engineConfigs"`
}

// Example configuration
var exampleConfig = MultiStoreConfig{
    RootPath:      "/data/stores",
    DefaultEngine: "leveldb",
    EngineConfigs: map[string]Config{
        "leveldb": {
            "cache_size":    64 * 1024 * 1024, // 64MB
            "write_buffer":  16 * 1024 * 1024, // 16MB
            "compression":   "snappy",
        },
        "file": {
            "format":        "json",
            "max_versions":  100,
        },
    },
}
```

### 7.2 Store-Specific Configuration

```go
// Example engine-specific configs
type LevelDBConfig struct {
    CacheSize   int    `json:"cache_size"`
    WriteBuffer int    `json:"write_buffer"`
    Compression string `json:"compression"`
}

type FileConfig struct {
    Format      string `json:"format"`       // "json", "gob", "msgpack"
    MaxVersions int    `json:"max_versions"`
    Compression bool   `json:"compression"`
}

type MemoryConfig struct {
    MaxSize     int64 `json:"max_size"`
    MaxVersions int   `json:"max_versions"`
}
```

## 8. Testing Strategy

### 8.1 Interface Compliance Tests

```go
// TestStoreCompliance tests that a store implementation satisfies the Store interface
func TestStoreCompliance(t *testing.T, createStore func() Store) {
    store := createStore()
    defer store.Close()
    
    // Test basic operations
    key := []byte("test-key")
    value := []byte("test-value")
    
    // Test Set/Get
    err := store.Set(key, value)
    require.NoError(t, err)
    
    retrieved, err := store.Get(key)
    require.NoError(t, err)
    require.Equal(t, value, retrieved)
    
    // Test Has
    exists, err := store.Has(key)
    require.NoError(t, err)
    require.True(t, exists)
    
    // Test Delete
    err = store.Delete(key)
    require.NoError(t, err)
    
    // Test key not found
    _, err = store.Get(key)
    require.True(t, IsKeyNotFound(err))
}
```

### 8.2 Optional Feature Tests

```go
// TestTransactionalCompliance tests transaction support
func TestTransactionalCompliance(t *testing.T, store Store) {
    tx, ok := store.(Transactional)
    if !ok {
        t.Skip("Store does not support transactions")
    }
    
    transaction, err := tx.BeginTx()
    require.NoError(t, err)
    require.True(t, transaction.IsActive())
    
    // Test transaction operations
    key := []byte("tx-key")
    value := []byte("tx-value")
    
    err = transaction.Set(key, value)
    require.NoError(t, err)
    
    // Value should not be visible outside transaction
    _, err = store.Get(key)
    require.True(t, IsKeyNotFound(err))
    
    // Commit and verify
    err = transaction.Commit()
    require.NoError(t, err)
    require.False(t, transaction.IsActive())
    
    retrieved, err := store.Get(key)
    require.NoError(t, err)
    require.Equal(t, value, retrieved)
}
```

## 9. Package Structure

Following the project's established domain/infrastructure separation:

```
# Domain Layer - Interfaces and Domain Logic
internal/domain/storage/
├── store.go              # Core Store interface and types
├── multistore.go         # MultiStore interface  
├── engine.go             # Engine interface and registry
├── transaction.go        # Transaction interfaces
├── version.go            # Versioning interfaces
├── config.go             # Configuration types
├── errors.go             # Error definitions
├── events.go             # Storage event definitions
└── testing/
    ├── compliance.go     # Interface compliance tests
    └── helpers.go        # Test utilities

# Infrastructure Layer - Concrete Implementations
internal/infrastructure/storage/
├── multistore_impl.go    # DefaultMultiStore implementation
├── factory.go            # Engine factory implementation
├── registry.go           # Engine registry implementation
├── memory/
│   ├── engine.go         # Memory engine implementation
│   ├── store.go          # Memory store implementation
│   ├── transaction.go    # Memory transaction implementation
│   └── store_test.go
├── file/
│   ├── engine.go         # File engine implementation
│   ├── store.go          # File store implementation
│   ├── versioned.go      # File versioning implementation
│   └── store_test.go
├── leveldb/
│   ├── engine.go         # LevelDB engine implementation
│   ├── store.go          # LevelDB store implementation
│   ├── transaction.go    # LevelDB transaction implementation
│   └── store_test.go
└── badger/
    ├── engine.go         # Badger engine implementation
    ├── store.go          # Badger store implementation
    ├── transaction.go    # Badger transaction implementation
    └── store_test.go

# Examples and Documentation
examples/storage/
├── basic_usage.go        # Basic usage examples
├── advanced.go           # Advanced feature examples
└── migration.go          # Migration examples
```

## 10. Implementation Checklist

### Phase 1: Core Interfaces (Week 1)
- [ ] Define Store interface in `store.go`
- [ ] Define optional interfaces (Transactional, Versioned, RangeQueryable)
- [ ] Define Engine interface in `engine.go`
- [ ] Create error definitions in `errors.go`
- [ ] Create configuration types in `config.go`

### Phase 2: MultiStore (Week 1)
- [ ] Implement MultiStore interface in `multistore.go`
- [ ] Implement DefaultMultiStore
- [ ] Add engine registration and management
- [ ] Create basic tests

### Phase 3: Memory Engine (Week 1-2)
- [ ] Implement MemoryEngine in `engines/memory/engine.go`
- [ ] Implement MemoryStore in `engines/memory/store.go`
- [ ] Add transaction support
- [ ] Add versioning support
- [ ] Create comprehensive tests

### Phase 4: File Engine (Week 2)
- [ ] Implement FileEngine in `engines/file/engine.go`
- [ ] Implement FileStore in `engines/file/store.go`
- [ ] Add JSON serialization
- [ ] Add versioning support
- [ ] Create comprehensive tests

### Phase 5: Testing Framework (Week 2)
- [ ] Create compliance tests in `testing/compliance.go`
- [ ] Create test utilities in `testing/helpers.go`
- [ ] Ensure all engines pass compliance tests
- [ ] Add performance benchmarks

### Phase 6: Documentation and Examples (Week 3)
- [ ] Create usage examples in `examples/`
- [ ] Write comprehensive README
- [ ] Document configuration options
- [ ] Create migration guide

### Phase 7: Advanced Engines (Week 3-4)
- [ ] Implement LevelDB engine (if needed)
- [ ] Implement Badger engine (if needed)
- [ ] Add plugin system for external engines
- [ ] Performance optimization

## 11. Usage Examples

### 11.1 Basic Usage

```go
import (
    "github.com/fintechain/skeleton/internal/domain/storage"
    multistore "github.com/fintechain/skeleton/internal/infrastructure/storage"
)

// Create MultiStore
ms := multistore.NewMultiStore("/data")

// Create stores
err := ms.CreateStore("users", "memory", nil)
err = ms.CreateStore("sessions", "file", storage.Config{
    "format": "json",
    "max_versions": 10,
})

// Use stores
userStore, err := ms.GetStore("users")
err = userStore.Set([]byte("user:123"), []byte(`{"name":"john"}`))

value, err := userStore.Get([]byte("user:123"))
```

### 11.2 Transaction Usage

```go
store, err := ms.GetStore("users")
if tx, ok := store.(storage.Transactional); ok {
    transaction, err := tx.BeginTx()
    if err != nil {
        return err
    }
    defer transaction.Rollback() // Rollback if not committed
    
    err = transaction.Set([]byte("key1"), []byte("value1"))
    err = transaction.Set([]byte("key2"), []byte("value2"))
    
    if err != nil {
        return err // Transaction will be rolled back
    }
    
    return transaction.Commit()
}
```

### 11.3 Versioning Usage

```go
store, err := ms.GetStore("data")
if vs, ok := store.(storage.Versioned); ok {
    // Make some changes
    store.Set([]byte("key"), []byte("value1"))
    
    // Save version
    version1, hash1, err := vs.SaveVersion()
    
    // Make more changes
    store.Set([]byte("key"), []byte("value2"))
    version2, hash2, err := vs.SaveVersion()
    
    // Load previous version
    err = vs.LoadVersion(version1)
    value, _ := store.Get([]byte("key")) // Returns "value1"
}
```

This specification provides a complete blueprint for implementing the minimal storage system. The design emphasizes simplicity while maintaining extensibility through optional interfaces and a clean plugin architecture.