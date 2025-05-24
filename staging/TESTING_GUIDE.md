# Component System Testing Guide

This document outlines our testing strategy for the redesigned component system, covering both unit and integration tests.

## Testing Philosophy

1. **Comprehensive Coverage**: Aim for high test coverage of all critical code paths
2. **Test Behaviors, Not Implementation**: Focus on testing the behavior of components, not their implementation details
3. **Hierarchical Approach**: Test from the smallest units up to full integration
4. **Edge Cases**: Explicitly test error conditions and edge cases 

## Test Organization

### Unit Tests

- **Location**: In the same package as the code being tested (`*_test.go` files)
- **Naming**: `TestXxx` for function tests, `TestTypeXxx` for methods on types
- **Focus**: Individual components and their methods in isolation
- **Access**: Can test unexported methods and internal state

### Integration Tests

- **Location**: In a separate `/tests` or `/integration` directory
- **Naming**: `TestIntXxx` prefix to distinguish from unit tests
- **Focus**: Interaction between components and subsystems
- **Access**: Only use exported APIs, simulating real client usage

## Test Helpers and Utilities

### Mock Components

Create simple mock implementations of core interfaces:

```go
// MockComponent implements Component interface for testing
type MockComponent struct {
    id        string
    name      string
    initCount int
    disposed  bool
}

func NewMockComponent(id, name string) *MockComponent {
    return &MockComponent{id: id, name: name}
}

// Implement Component interface methods...
```

### Test Registry

```go
// TestRegistry provides a simple Registry implementation for testing
func NewTestRegistry() Registry {
    // Create a minimal Registry implementation for tests
}
```

### Event Recorder

```go
// EventRecorder captures events for testing
type EventRecorder struct {
    events []Event
    mu     sync.Mutex
}

func (r *EventRecorder) RecordEvent(e Event) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.events = append(r.events, e)
}

func (r *EventRecorder) Events() []Event {
    r.mu.Lock()
    defer r.mu.Unlock()
    return append([]Event{}, r.events...)
}
```

### Context Utilities

```go
// TestContext creates a context with timeout for tests
func TestContext(t *testing.T) (Context, func()) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    return WrapContext(ctx), cancel
}
```

## Testing Specific Components

### Component Tests

- Test all lifecycle methods (Initialize, Dispose)
- Verify proper error propagation
- Test metadata handling
- Test proper event emission on lifecycle changes

### Registry Tests

- Test component registration and unregistration
- Test component lookup by ID, type, and metadata
- Test duplicate registration handling
- Test dependency resolution
- Test bulk operations (initialize all, dispose all)

### Factory Tests

- Test component creation from various configurations
- Test configuration validation
- Test error handling for invalid configurations
- Test type-specific factory methods

### Service Tests

- Test all service lifecycle methods (Start, Stop)
- Test state transitions and invalid transitions
- Test status reporting
- Test health check functionality
- Test concurrency safety

### Operation Tests

- Test execution with various inputs
- Test context handling
- Test error propagation
- Test cancellation and timeouts
- Test input validation

### Event Bus Tests

- Test publication to various topics
- Test subscription and unsubscription
- Test synchronous and asynchronous event handling
- Test event delivery under load
- Test proper cleanup of resources

### Plugin Tests

- Test plugin discovery
- Test plugin loading and unloading
- Test plugin component registration
- Test isolation between plugins
- Test version compatibility

## Integration Test Scenarios

1. **Component Lifecycle Integration**
   - Register multiple components with dependencies
   - Initialize all components
   - Verify proper initialization order
   - Dispose all components
   - Verify proper disposal order

2. **Service Orchestration**
   - Register multiple services with dependencies
   - Start all services
   - Verify proper start order
   - Verify service interactions
   - Stop all services
   - Verify proper stop order

3. **Operation Pipeline**
   - Create a pipeline of operations
   - Execute the pipeline with sample input
   - Verify correct output and proper operation chaining
   - Test error propagation through the pipeline

4. **Plugin System Integration**
   - Discover and load multiple plugins
   - Verify plugin components are registered
   - Test interactions between plugin components
   - Unload plugins and verify proper cleanup

5. **Event-Driven Integration**
   - Set up components that interact via events
   - Trigger events and verify proper propagation
   - Test complex event chains and reactions

## Test Coverage Goals

- **Core Interfaces**: 100% coverage
- **Base Implementations**: 95%+ coverage
- **Utility Functions**: 90%+ coverage
- **Edge Cases**: Explicit tests for all error conditions
- **Concurrency**: Dedicated tests for race conditions

## Test Implementation Process

1. Start with core interfaces and base implementations
2. Implement test helpers and utilities
3. Write unit tests for each component type
4. Create integration tests for component interactions
5. Add performance and stress tests for critical components
6. Continuously run tests as part of development workflow

## Benchmarking

For performance-critical components:

```go
func BenchmarkRegistry_Get(b *testing.B) {
    // Setup test registry with many components
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        // Benchmark component lookup
    }
}
```

## Test Documentation

- Document test coverage metrics
- Document performance benchmarks
- Document test patterns and utilities for developer reference
- Update this guide as testing requirements evolve 