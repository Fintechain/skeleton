# Storage System Implementation Prompt

## Context and Role
You are an expert Go developer implementing a minimal storage system for the Skeleton project. You have been provided with a comprehensive design specification that follows clean architecture principles with domain/infrastructure separation.

## Your Mission
Implement the complete storage system according to the "Minimal Storage System Design Specification" document. The system should be production-ready, well-tested, and follow Go best practices.

## Critical Architecture Requirements

### Package Structure Compliance
**MUST FOLLOW**: The project uses strict domain/infrastructure separation:

- **Domain Layer** (`internal/domain/storage/`): Contains ONLY interfaces, types, errors, events
- **Infrastructure Layer** (`internal/infrastructure/storage/`): Contains ALL concrete implementations

```
internal/domain/storage/          # Interfaces & Domain Logic
internal/infrastructure/storage/  # Concrete Implementations
```

### Import Rules
- Infrastructure packages import domain packages
- Domain packages NEVER import infrastructure packages
- All concrete types return domain interfaces

## Implementation Phases

### Phase 1: Domain Layer Foundation
**Files to create in `internal/domain/storage/`:**

1. **`errors.go`** - All error definitions and helper functions
2. **`store.go`** - Core Store interface and related types
3. **`transaction.go`** - Transaction and Transactional interfaces
4. **`version.go`** - Versioned interface and versioning types
5. **`engine.go`** - Engine interface and capabilities
6. **`multistore.go`** - MultiStore interface definition
7. **`config.go`** - Configuration types and constants
8. **`events.go`** - Storage event definitions (if needed)

### Phase 2: Infrastructure Layer - Memory Engine
**Files to create in `internal/infrastructure/storage/memory/`:**

1. **`engine.go`** - Memory engine implementation
2. **`store.go`** - Memory store implementation
3. **`transaction.go`** - Memory transaction implementation
4. **`store_test.go`** - Comprehensive tests

### Phase 3: Infrastructure Layer - File Engine
**Files to create in `internal/infrastructure/storage/file/`:**

1. **`engine.go`** - File engine implementation
2. **`store.go`** - File store implementation
3. **`versioned.go`** - File versioning implementation
4. **`store_test.go`** - Comprehensive tests

### Phase 4: Infrastructure Layer - Core Services
**Files to create in `internal/infrastructure/storage/`:**

1. **`multistore_impl.go`** - DefaultMultiStore implementation
2. **`factory.go`** - Engine factory (if needed)
3. **`registry.go`** - Engine registry (if needed)

### Phase 5: Testing Framework
**Files to create in `internal/domain/storage/testing/`:**

1. **`compliance.go`** - Interface compliance tests
2. **`helpers.go`** - Test utilities and mocks

## Code Quality Requirements

### 1. Go Best Practices
- Use proper error handling with wrapped errors
- Include comprehensive comments for all public interfaces
- Follow Go naming conventions
- Use proper package documentation
- Implement proper resource cleanup (defer statements)

### 2. Thread Safety
- All stores must be thread-safe
- Use proper mutex patterns (RWMutex for read-heavy operations)
- Protect all shared state
- Document thread-safety guarantees

### 3. Error Handling
- Use standard Go error patterns from the specification
- Wrap errors with context using `fmt.Errorf` with `%w` verb
- Return domain-specific errors (ErrKeyNotFound, etc.)
- Never panic in normal operation

### 4. Memory Management
- Always copy byte slices when storing/returning data
- Properly close resources in defer statements
- Avoid memory leaks in long-running operations

### 5. Testing Requirements
- **Unit tests** for all implementations with >90% coverage
- **Compliance tests** that verify interface adherence
- **Integration tests** for MultiStore functionality
- **Benchmark tests** for performance-critical operations
- Use table-driven tests where appropriate

## Implementation Guidelines

### Interface Implementation Pattern
```go
// Domain interface (internal/domain/storage/store.go)
type Store interface {
    Get(key []byte) ([]byte, error)
    // ... other methods
}

// Infrastructure implementation (internal/infrastructure/storage/memory/store.go)
type Store struct {
    // implementation fields
}

// Constructor returns domain interface
func NewStore(name, path string) storage.Store {
    return &Store{
        // initialization
    }
}
```

### Error Handling Pattern
```go
// Use domain errors
if !exists {
    return nil, storage.ErrKeyNotFound
}

// Wrap errors with context
if err != nil {
    return storage.WrapError(err, "failed to read from disk")
}
```

### Optional Interface Pattern
```go
// Check for optional capabilities
if tx, ok := store.(storage.Transactional); ok {
    transaction, err := tx.BeginTx()
    // ...
}
```

## Specific Implementation Notes

### Memory Engine
- Use `map[string][]byte` for data storage
- Implement full transaction support with rollback
- Support versioning with snapshot copies
- Thread-safe with RWMutex

### File Engine
- Use JSON or GOB serialization
- Implement file-based versioning with numbered files
- Support graceful recovery from corruption
- Efficient loading/saving strategies

### MultiStore
- Manage multiple store instances by name
- Support different engines per store
- Provide engine registration and discovery
- Thread-safe store management

## Testing Strategy

### Compliance Tests
Create tests that verify ANY store implementation satisfies the interfaces:
```go
func TestStoreCompliance(t *testing.T, createStore func() storage.Store) {
    // Test all Store interface methods
}

func TestTransactionalCompliance(t *testing.T, store storage.Store) {
    // Test transaction capabilities if supported
}
```

### Performance Requirements
- Get/Set operations should be under 1ms for memory engine
- File operations should handle 10MB+ files efficiently
- Memory usage should be reasonable for large datasets
- No memory leaks in long-running tests

## Deliverables Checklist

### Domain Layer ✓
- [ ] All interfaces defined with proper documentation
- [ ] Error types and constants defined
- [ ] Configuration types defined
- [ ] No infrastructure dependencies

### Infrastructure Layer ✓
- [ ] Memory engine fully implemented with all optional interfaces
- [ ] File engine implemented with versioning support
- [ ] MultiStore implementation with engine management
- [ ] All implementations return domain interfaces

### Testing ✓
- [ ] Compliance tests for all interfaces
- [ ] Unit tests for all implementations
- [ ] Integration tests for MultiStore
- [ ] Performance benchmarks
- [ ] Test coverage >90%

### Documentation ✓
- [ ] All public interfaces documented
- [ ] Implementation notes in README
- [ ] Usage examples
- [ ] Performance characteristics documented

## Success Criteria

The implementation is successful when:

1. **All tests pass** including compliance, unit, and integration tests
2. **Interfaces are clean** and follow the specification exactly
3. **Performance is acceptable** for the intended use cases
4. **Code is maintainable** with clear structure and documentation
5. **Architecture compliance** with proper domain/infrastructure separation
6. **Thread safety** is guaranteed and tested
7. **Error handling** is robust and informative

## Getting Started

1. **Start with Phase 1** - Create all domain interfaces first
2. **Implement Memory engine** - Begin with the simplest implementation
3. **Add compliance tests** - Ensure your implementation satisfies interfaces
4. **Implement File engine** - Add persistent storage capability
5. **Build MultiStore** - Tie everything together
6. **Comprehensive testing** - Achieve high test coverage
7. **Performance optimization** - Profile and optimize critical paths

Remember: The goal is a simple, clean, highly testable storage system that can be extended with additional engines in the future. Prioritize correctness and clarity over premature optimization.

## Questions to Consider While Implementing

- Is this interface easy to implement for new engines?
- Are the error messages helpful for debugging?
- Is the thread safety model clear and consistent?
- Can this be tested effectively?
- Does this follow Go idioms and conventions?
- Is the separation between domain and infrastructure clean?

Start implementing and create production-quality Go code that follows the specification precisely!