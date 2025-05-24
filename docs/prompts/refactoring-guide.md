# Refactoring Guide: Implementing Dependency Injection & Interface Substitution

## Pattern Overview
Apply these patterns to decouple your code from external dependencies for unit testing:

## 1. Interface Substitution Pattern
- Utilize existing interfaces for external dependencies whenever possible
- Only define new interfaces when absolutely necessary (i.e., when the external dependency doesn't provide one)
- Replace direct dependency usage with interface calls
- Keep implementation files unchanged externally
```go
// Example using an existing interface
// If the external dependency already has an interface like:
type ExternalStorageInterface interface {
    Get(key []byte) ([]byte, error)
    Set(key, value []byte) error
}

// Use it directly instead of creating a new one
```

## 2. Constructor Injection Pattern
- Modify constructors to accept dependency interfaces instead of creating them directly
- Maintain backward compatibility using factory methods
```go
// New constructor with dependency injection
func NewComponent(id string, storage ExternalStorageInterface) *Component {
    return &Component{id: id, storage: storage}
}

// Factory method for backward compatibility
func CreateComponent(id string) (*Component, error) {
    realStorage := CreateRealStorageImplementation()
    return NewComponent(id, realStorage), nil
}
```

## 3. Mock Implementation Pattern
Create mocks only for external dependencies, following Go standards:
- Place mock files in the same package/directory as the tests that use them
- Use consistent naming: `*_mock.go` (e.g., `storage_mock.go`)
- Each mock should implement the corresponding dependency interface fully
- Keep mock implementations focused on test requirements

For example:
```go
// storage.go - Real implementation
package mypackage

// storage_mock.go - Mock implementation
package mypackage

type MockStorage struct {
    // Fields needed for the mock
}

// Implement all methods of the StorageInterface
func (m *MockStorage) Get(key []byte) ([]byte, error) {
    // Mock implementation
}
```

## 4. Test Implementation Pattern
- Every file with code logic should have a corresponding test file
- Focus on unit tests (not integration tests) for pure business logic
- Aim for 90% test coverage to ensure proper validation of code behavior
- Update test files to inject mock dependencies:
```go
// component_test.go
func TestComponent(t *testing.T) {
    mockStorage := &MockStorage{} // Setup the mock with test-specific behavior
    component := NewComponent("test-id", mockStorage)
    
    // Test component behavior using the mock
    result := component.DoSomething()
    
    // Assert results
    if result != expectedResult {
        t.Errorf("Expected %v, got %v", expectedResult, result)
    }
}
```

## Implementation Steps
1. Identify external dependencies in your codebase
2. Check if the dependencies already have interfaces - use them if available
3. Only define new interfaces for dependencies without existing interfaces
4. Refactor constructors to accept these interfaces 
5. Update code to use the interfaces instead of concrete implementations
6. Create mock implementations of external dependencies for testing
7. Ensure every file with logic has a corresponding test file
8. Refactor tests to inject these mock dependencies
9. Aim for 90% test coverage through comprehensive unit tests

## Benefits
- Isolate your code from external systems during testing
- Test edge cases and failure scenarios that are difficult to reproduce with real dependencies
- Improve separation of concerns
- Maintain existing API compatibility
- Ensure high test coverage with pure unit tests

This pattern focuses on using existing interfaces where possible, only creating new ones when necessary, and thoroughly testing with mocks to achieve 90% test coverage.