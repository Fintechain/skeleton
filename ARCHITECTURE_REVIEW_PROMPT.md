# Fintechain Skeleton Architecture Review Prompt

Use this prompt to conduct comprehensive reviews of the Fintechain skeleton architecture to ensure it maintains its quality, design principles, and provides a robust foundation for application development.

## Review Instructions

When reviewing the Fintechain skeleton, systematically evaluate each section below. For each area, provide:
- ✅ **Strengths**: What is working well
- ⚠️ **Issues**: Problems or areas of concern  
- 🔧 **Recommendations**: Specific improvements needed
- 📊 **Score**: Rate 1-5 (5 = excellent, 1 = needs major work)

## 1. Package Organization & Structure

### Evaluation Criteria:
- **Package Independence**: Each package should be self-contained with minimal cross-dependencies
- **Focused Responsibility**: Each package should have a single, clear purpose
- **Logical Hierarchy**: Package structure should reflect architectural layers
- **No Monolithic Packages**: Avoid large packages that handle multiple concerns

### Review Questions:
1. Are packages organized by domain/function rather than technical layers?
2. Can each package be understood and modified independently?
3. Are there any circular dependencies between packages?
4. Do package names clearly indicate their purpose?
5. Is the separation between `pkg/` (public API) and `internal/` (implementation) maintained?

### Expected Structure:
```
pkg/                    # Public API
├── system/            # System lifecycle only
├── component/         # Component management
├── service/           # Service lifecycle with factories
├── operation/         # Operations with factories
├── plugin/            # Plugin system
├── config/            # Configuration management
├── event/             # Event system
├── storage/           # Storage abstractions
├── logging/           # Logging utilities
├── id/                # ID generation
└── crypto/            # Cryptographic utilities

internal/              # Implementation details
├── domain/            # Domain logic
├── infrastructure/    # Infrastructure implementations
└── application/       # Application services
```

## 2. Factory Pattern Implementation

### Evaluation Criteria:
- **Consistent Creation**: All major components use factory patterns
- **Configuration-Driven**: Factories accept configuration objects
- **Interface Compliance**: Factories implement standard interfaces
- **Error Handling**: Proper error handling in factory methods

### Review Questions:
1. Do service and operation packages provide factory constructors?
2. Are factory interfaces consistent across packages?
3. Do factories support both configuration-based and convenience creation?
4. Are factory implementations properly integrated with component factories?
5. Do factories provide clear error messages with proper error codes?

### Required Factory Methods:
```go
// Service Package
func NewServiceFactory() ServiceFactory
func (f *ServiceFactory) CreateService(config ServiceConfig) (Service, error)
func NewService(name, type string) (Service, error)

// Operation Package  
func NewOperationFactory() OperationFactory
func (f *OperationFactory) CreateOperation(config OperationConfig) (Operation, error)
func NewOperation(name, type string) (Operation, error)
```

## 3. Interface Design & Consistency

### Evaluation Criteria:
- **Interface-Driven**: Major functionality exposed through interfaces
- **Minimal Interfaces**: Interfaces should be focused and minimal
- **Consistent Patterns**: Similar interfaces across packages
- **Proper Abstractions**: Interfaces hide implementation details

### Review Questions:
1. Are all major components defined as interfaces?
2. Do interfaces follow the "accept interfaces, return structs" principle?
3. Are interface methods focused and cohesive?
4. Do similar components have similar interface patterns?
5. Are there any leaky abstractions?

### Core Interface Patterns:
```go
// Component Pattern
type Component interface {
    ID() string
    Name() string
    Type() ComponentType
    Initialize(Context) error
    Dispose(Context) error
}

// Factory Pattern
type Factory interface {
    Create(ComponentConfig) (Component, error)
}

// Lifecycle Pattern (Services)
type Service interface {
    Component
    Start(Context) error
    Stop(Context) error
    Status() ServiceStatus
}
```

## 4. Error Handling & Type Safety

### Evaluation Criteria:
- **Structured Errors**: Domain-specific error types with codes
- **Error Chaining**: Support for wrapping underlying errors
- **Type Safety**: Strong typing throughout the system
- **Consistent Patterns**: Uniform error handling across packages

### Review Questions:
1. Does each package define domain-specific error codes?
2. Are errors properly wrapped and chained?
3. Do packages provide error checking utilities?
4. Are error messages descriptive and actionable?
5. Is type safety maintained without excessive type assertions?

### Required Error Patterns:
```go
// Error Constants
const (
    ErrComponentNotFound = "component.not_found"
    ErrServiceStart     = "service.start_failed"
)

// Error Creation
func NewError(code, message string, cause error) *Error

// Error Checking
func IsComponentError(err error, code string) bool
```

## 5. Documentation & Examples

### Evaluation Criteria:
- **Comprehensive Documentation**: Clear package and API documentation
- **Usage Examples**: Practical examples for all major features
- **Architecture Documentation**: Clear explanation of design decisions
- **Migration Guides**: Help for breaking changes

### Review Questions:
1. Is each package documented with purpose and usage examples?
2. Are all public APIs documented with Go doc comments?
3. Does the README provide comprehensive usage patterns?
4. Are design principles clearly explained?
5. Is there a migration guide for breaking changes?

### Required Documentation:
- Package-level documentation with examples
- README with comprehensive usage patterns
- CHANGELOG with detailed version history
- Architecture review prompt (this document)

## 6. Testing & Quality Assurance

### Evaluation Criteria:
- **Test Coverage**: Comprehensive test coverage for all packages
- **Test Organization**: Tests organized by package and functionality
- **Mock Support**: Interfaces support easy mocking
- **Integration Tests**: End-to-end testing capabilities

### Review Questions:
1. Does each package have comprehensive unit tests?
2. Are interfaces easily mockable for testing?
3. Are there integration tests for major workflows?
4. Do tests cover error conditions and edge cases?
5. Is there a clear testing strategy documented?

## 7. Extensibility & Plugin Architecture

### Evaluation Criteria:
- **Plugin Support**: Clear plugin architecture
- **Extension Points**: Well-defined extension mechanisms
- **Backward Compatibility**: Changes don't break existing code
- **Customization**: Easy to customize and extend

### Review Questions:
1. Can new functionality be added without modifying core code?
2. Is the plugin system well-designed and documented?
3. Are there clear extension points for customization?
4. Do changes maintain backward compatibility?
5. Can components be easily replaced or extended?

## 8. Performance & Scalability

### Evaluation Criteria:
- **Efficient Patterns**: No obvious performance bottlenecks
- **Resource Management**: Proper resource cleanup
- **Scalable Design**: Architecture supports scaling
- **Memory Management**: No memory leaks or excessive allocations

### Review Questions:
1. Are there any obvious performance bottlenecks?
2. Is resource cleanup properly handled?
3. Does the architecture support horizontal scaling?
4. Are there any memory management issues?
5. Is the design efficient for high-throughput scenarios?

## 9. Security Considerations

### Evaluation Criteria:
- **Secure Defaults**: Secure configuration by default
- **Input Validation**: Proper validation of inputs
- **Error Information**: Errors don't leak sensitive information
- **Dependency Security**: Dependencies are secure and up-to-date

### Review Questions:
1. Are there secure defaults for all configurations?
2. Is input validation comprehensive and consistent?
3. Do error messages avoid leaking sensitive information?
4. Are dependencies regularly updated and security-scanned?
5. Are there any obvious security vulnerabilities?

## 10. Maintainability & Code Quality

### Evaluation Criteria:
- **Code Organization**: Well-organized and readable code
- **Naming Conventions**: Consistent and clear naming
- **Code Duplication**: Minimal code duplication
- **Complexity**: Manageable complexity levels

### Review Questions:
1. Is the code well-organized and easy to navigate?
2. Are naming conventions consistent and descriptive?
3. Is there minimal code duplication?
4. Are complex areas properly documented?
5. Is the overall complexity manageable?

## Review Checklist

Use this checklist to ensure comprehensive review:

### Package Structure
- [ ] Package organization follows domain boundaries
- [ ] No circular dependencies
- [ ] Clear separation of public API and implementation
- [ ] Focused package responsibilities

### Factory Patterns
- [ ] Service factory implementation complete
- [ ] Operation factory implementation complete
- [ ] Factory interfaces consistent
- [ ] Configuration-driven creation supported

### Interface Design
- [ ] Major components are interface-driven
- [ ] Interfaces are minimal and focused
- [ ] Consistent patterns across packages
- [ ] No leaky abstractions

### Error Handling
- [ ] Domain-specific error codes defined
- [ ] Error wrapping and chaining supported
- [ ] Error checking utilities provided
- [ ] Type safety maintained

### Documentation
- [ ] Package documentation complete
- [ ] API documentation comprehensive
- [ ] Usage examples provided
- [ ] Architecture decisions documented

### Testing
- [ ] Unit test coverage adequate
- [ ] Integration tests present
- [ ] Interfaces easily mockable
- [ ] Error conditions tested

### Extensibility
- [ ] Plugin architecture functional
- [ ] Extension points well-defined
- [ ] Backward compatibility maintained
- [ ] Customization supported

### Quality
- [ ] Performance considerations addressed
- [ ] Security best practices followed
- [ ] Code quality high
- [ ] Maintainability excellent

## Overall Assessment

After completing the review, provide:

1. **Overall Score**: 1-5 rating for the skeleton's quality
2. **Key Strengths**: Top 3 architectural strengths
3. **Critical Issues**: Any issues that must be addressed
4. **Improvement Priorities**: Top 3 areas for improvement
5. **Recommendation**: Ready for use / Needs work / Major revision required

## Success Criteria

The skeleton should achieve:
- ✅ **Score 4+ in all areas**: Consistently high quality
- ✅ **No critical issues**: No blocking problems
- ✅ **Clear documentation**: Comprehensive and accurate
- ✅ **Extensible design**: Easy to extend and customize
- ✅ **Production ready**: Suitable for building real applications

Use this prompt regularly to maintain architectural quality and ensure the skeleton continues to provide a robust foundation for application development. 