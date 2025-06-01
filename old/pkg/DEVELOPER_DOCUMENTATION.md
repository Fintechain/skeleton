# Fintechain Skeleton Framework - Developer Documentation

## Overview

The Fintechain Skeleton Framework provides a comprehensive set of packages for building modular, component-based applications. This documentation covers the public APIs available in the `skeleton/pkg` directory, which serves as the primary interface for developers using the framework.

## Table of Contents

1. [Architecture Overview](#architecture-overview)
2. [Core Packages](#core-packages)
3. [Package Reference](#package-reference)
4. [Usage Examples](#usage-examples)
5. [Error Handling](#error-handling)
6. [Best Practices](#best-practices)

## Architecture Overview

The skeleton framework follows a layered architecture with clear separation between public APIs (`pkg/`) and internal implementations (`internal/`). The framework is built around the concept of components, which are the fundamental building blocks of the system.

### Key Concepts

- **Components**: Basic building blocks that can be composed to create complex systems
- **Services**: Long-running components that provide ongoing functionality
- **Operations**: Discrete units of work that execute specific tasks
- **Plugins**: Containers for components that extend the system
- **Events**: Asynchronous communication mechanism between components
- **Storage**: Persistent data management with multiple backend support

## Core Packages

### 1. Component System (`pkg/component`)

The component package provides the foundation for all other packages in the framework.

#### Key Types

```go
// Core interfaces
type Component interface
type Context interface
type Registry interface
type Factory interface

// Dependency management
type DependencyAware interface
type DependencyAwareComponent interface

// Base implementations
type BaseComponent interface
type DefaultFactory interface
```

#### Component Types

- `TypeBasic`: Basic components
- `TypeOperation`: Operation components
- `TypeService`: Service components
- `TypeSystem`: System-level components
- `TypeApplication`: Application-level components

#### Key Functions

```go
// Component creation
func NewRegistry() Registry
func NewBaseComponent(id, name string, componentType ComponentType) Component
func NewDependencyAwareComponent(base Component, dependencies []string) DependencyAware
func NewFactory() Factory

// Context management
func NewContext(ctx stdctx.Context) Context
func WithCancel(parent Context) (Context, func())
func WithTimeout(parent Context, timeout time.Duration) (Context, func())
```

### 2. System Management (`pkg/system`)

The system package provides the main entry point for system lifecycle management using a functional options pattern.

#### Key Types

```go
type SystemService interface
type Config struct {
    ServiceID     string
    StorageConfig storage.MultiStoreConfig
}
```

#### System Startup

```go
func StartSystem(options ...Option) error

// Configuration options
func WithConfig(config *InternalConfig) Option
func WithPlugins(plugins []plugin.Plugin) Option
func WithRegistry(registry component.Registry) Option
func WithPluginManager(pluginMgr plugin.PluginManager) Option
func WithEventBus(eventBus event.EventBus) Option
func WithMultiStore(multiStore storage.MultiStore) Option
```

#### Error Codes

- `ErrSystemStartup`: System startup failures
- `ErrSystemShutdown`: System shutdown failures
- `ErrInvalidConfig`: Invalid configuration
- `ErrDependencyMissing`: Missing dependencies

### 3. Storage System (`pkg/storage`)

The storage package provides a unified interface for data persistence with support for multiple backends.

#### Core Interfaces

```go
type Store interface           // Basic key-value operations
type MultiStore interface     // Multiple named stores
type Engine interface         // Storage backend implementation
type Transaction interface    // Atomic operations
type Transactional interface  // Transaction support
type Versioned interface      // Versioning/snapshots
type RangeQueryable interface // Range queries
```

#### Configuration

```go
type MultiStoreConfig struct {
    RootPath      string
    DefaultEngine string
}

type DefaultStoreConfig struct {
    Engine  string
    Path    string
    Options Config
}
```

#### Key Functions

```go
// Store creation
func NewMultiStore(config *MultiStoreConfig, logger logging.Logger, eventBus event.EventBus) MultiStore
func NewMemoryEngine(logger logging.Logger) Engine
func CreateStore(name, engine, path string) (Store, error)

// Convenience operations
func Get(store Store, key []byte) ([]byte, error)
func Set(store Store, key, value []byte) error
func Delete(store Store, key []byte) error
func WithTransaction(store Store, fn func(Transaction) error) error
```

### 4. Service Management (`pkg/service`)

Services are specialized components that provide ongoing functionality.

#### Key Types

```go
type Service interface
type ServiceFactory interface
type HealthCheck interface
type ServiceStatus int
type HealthStatus int
```

#### Service Status Constants

- `StatusStopped`: Service is stopped
- `StatusStarting`: Service is starting
- `StatusRunning`: Service is running
- `StatusStopping`: Service is stopping
- `StatusFailed`: Service has failed

#### Health Status Constants

- `HealthStatusUnknown`: Health status unknown
- `HealthStatusHealthy`: Service is healthy
- `HealthStatusUnhealthy`: Service is unhealthy
- `HealthStatusDegraded`: Service is degraded

#### Key Functions

```go
func NewServiceFactory() ServiceFactory
func NewService(name, serviceType string) (Service, error)
func IsHealthy(svc Service) bool
func WaitForStatus(ctx component.Context, svc Service, expectedStatus ServiceStatus) bool
```

### 5. Operation Management (`pkg/operation`)

Operations are specialized components that execute discrete units of work.

#### Key Types

```go
type Operation interface
type OperationFactory interface
type Input interface
type Output interface
type OperationConfig struct
```

#### Key Functions

```go
func NewOperationFactory() OperationFactory
func NewOperation(name, operationType string) (Operation, error)
func Execute(ctx component.Context, op Operation, input Input) (Output, error)
func ValidateInput(input Input, requiredFields ...string) error
func CreateOutput(data interface{}) Output
func CreateErrorOutput(message string, cause error) Output
```

### 6. Event System (`pkg/event`)

The event system provides asynchronous communication between components.

#### Key Types

```go
type EventBus interface
type EventHandler interface
type Subscription interface
type Event struct {
    Topic   string
    Source  string
    Time    time.Time
    Payload map[string]interface{}
}
```

#### Key Functions

```go
func NewEventBus() EventBus
func NewEvent(topic, source string, payload map[string]interface{}) *Event
func Publish(bus EventBus, topic string, data interface{})
func Subscribe(bus EventBus, topic string, handler EventHandler) Subscription
func SubscribeAsync(bus EventBus, topic string, handler EventHandler) Subscription
```

### 7. Plugin System (`pkg/plugin`)

Plugins are containers for components that extend the system.

#### Key Types

```go
type Plugin interface
type PluginManager interface
type PluginInfo struct
```

#### Key Functions

```go
func NewPluginManager() PluginManager
```

#### Error Codes

- `ErrPluginNotFound`: Plugin not found
- `ErrPluginLoad`: Plugin load failed
- `ErrPluginUnload`: Plugin unload failed
- `ErrPluginDiscovery`: Plugin discovery failed
- `ErrPluginConflict`: Plugin conflict
- `ErrInvalidPlugin`: Invalid plugin

### 8. Configuration Management (`pkg/config`)

The configuration package provides access to application configuration values.

#### Key Types

```go
type Configuration interface
type ConfigurationSource interface
```

#### Key Functions

```go
func NewConfiguration(name string, source ...ConfigurationSource) Configuration
func NewDefaultConfiguration() Configuration

// Value retrieval
func GetString(cfg Configuration, key string) string
func GetStringDefault(cfg Configuration, key, defaultValue string) string
func GetInt(cfg Configuration, key string) (int, error)
func GetIntDefault(cfg Configuration, key string, defaultValue int) int
func GetBool(cfg Configuration, key string) (bool, error)
func GetBoolDefault(cfg Configuration, key string, defaultValue bool) bool
func GetDuration(cfg Configuration, key string) (time.Duration, error)
func GetDurationDefault(cfg Configuration, key string, defaultValue time.Duration) time.Duration
func GetObject(cfg Configuration, key string, result interface{}) error
func Exists(cfg Configuration, key string) bool
```

### 9. Cryptographic Utilities (`pkg/crypto`)

The crypto package provides cryptographic utilities for the framework.

#### Hash Functions

```go
func HashSHA256(input string) string
func HashSHA256Bytes(input []byte) string
func HashSHA1(input string) string
func HashSHA512(input string) string
func HashMD5(input string) string // Note: MD5 is cryptographically broken
```

#### HMAC Functions

```go
func HMACSHA256(key, input string) string
func HMACSHA256Bytes(key, input []byte) string
func HMACSHA512(key, input string) string
func VerifyHMAC(key, input, expectedHMAC string) bool
func VerifyHMACBytes(key, input, expectedHMAC []byte) bool
```

#### Random Generation

```go
func GenerateRandomBytes(length int) ([]byte, error)
func GenerateRandomString(length int) (string, error)
func GenerateRandomHex(byteLength int) (string, error)
```

#### Encoding Utilities

```go
func EncodeBase64(input string) string
func DecodeBase64(input string) (string, error)
func EncodeHex(input string) string
func DecodeHex(input string) (string, error)
```

#### Hash Utilities

```go
func HashWithAlgorithm(algorithm, input string) (string, error)
```

### 10. Validation Utilities (`pkg/validation`)

The validation package provides common validation utilities.

#### Key Types

```go
type Validator interface {
    Validate(value interface{}) error
}
type ValidatorFunc func(value interface{}) error
```

#### Validation Functions

```go
func Required(fieldName string, value interface{}) error
func MinLength(fieldName string, value string, minLen int) error
func Email(fieldName string, value string) error
func Regex(fieldName string, value string, pattern string) error
```

#### Error Codes

- `ErrRequired`: Required field validation
- `ErrMinLength`: Minimum length validation
- `ErrInvalidFormat`: Invalid format validation
- `ErrInvalidPattern`: Invalid pattern validation

### 11. ID Generation (`pkg/id`)

The ID package provides ID generation utilities.

#### Key Types

```go
type IDGeneratorInterface interface {
    GenerateID() (string, error)
}

type ProcessIDGenerator struct {
    prefix string
}
```

#### Key Functions

```go
func NewProcessIDGenerator(prefix string) *ProcessIDGenerator
func (gen *ProcessIDGenerator) GenerateID() (string, error)
```

### 12. Logging (`pkg/logging`)

The logging package provides public access to logging utilities.

#### Key Types

```go
type Logger interface
type LogLevel int
type LoggerBackend interface
type StandardLogger interface
type StandardLoggerOptions struct
```

#### Log Levels

- `Debug`: Detailed troubleshooting
- `Info`: General operational information
- `Warn`: Warnings
- `Error`: Errors
- `Fatal`: Fatal errors that cause program exit

#### Key Functions

```go
func NewStandardLogger(options StandardLoggerOptions) *StandardLogger
func CreateStandardLogger(level LogLevel) *StandardLogger
func CreateStandardLoggerWithWriter(level LogLevel, writer io.Writer) *StandardLogger
```

## Usage Examples

### Basic System Setup

```go
package main

import (
    "github.com/fintechain/skeleton/pkg/system"
    "github.com/fintechain/skeleton/pkg/storage"
    "github.com/fintechain/skeleton/pkg/logging"
    "github.com/fintechain/skeleton/pkg/event"
)

func main() {
    // Create logger
    logger := logging.CreateStandardLogger(logging.Info)
    
    // Create event bus
    eventBus := event.NewEventBus()
    
    // Create storage configuration
    storageConfig := &storage.MultiStoreConfig{
        RootPath:      "./data",
        DefaultEngine: "memory",
    }
    
    // Create multistore
    multiStore := storage.NewMultiStore(storageConfig, logger, eventBus)
    
    // Start system
    err := system.StartSystem(
        system.WithEventBus(eventBus),
        system.WithMultiStore(multiStore),
    )
    if err != nil {
        logger.Error("Failed to start system: %v", err)
        return
    }
}
```

### Creating and Using a Service

```go
package main

import (
    "github.com/fintechain/skeleton/pkg/service"
    "github.com/fintechain/skeleton/pkg/component"
)

func main() {
    // Create a service
    svc, err := service.NewService("api-server", "http")
    if err != nil {
        panic(err)
    }
    
    // Start the service
    ctx := component.Background()
    err = svc.Start(ctx)
    if err != nil {
        panic(err)
    }
    
    // Check if service is healthy
    if service.IsHealthy(svc) {
        fmt.Println("Service is healthy")
    }
    
    // Wait for service to be running
    if service.WaitForStatus(ctx, svc, service.StatusRunning) {
        fmt.Println("Service is now running")
    }
}
```

### Storage Operations

```go
package main

import (
    "github.com/fintechain/skeleton/pkg/storage"
    "github.com/fintechain/skeleton/pkg/logging"
    "github.com/fintechain/skeleton/pkg/event"
)

func main() {
    // Create a store
    store, err := storage.CreateStore("user-data", "memory", "./data")
    if err != nil {
        panic(err)
    }
    
    // Store data
    err = storage.Set(store, []byte("user:123"), []byte(`{"name": "John", "email": "john@example.com"}`))
    if err != nil {
        panic(err)
    }
    
    // Retrieve data
    data, err := storage.Get(store, []byte("user:123"))
    if err != nil {
        if storage.IsStorageError(err, storage.ErrKeyNotFound) {
            fmt.Println("User not found")
        } else {
            panic(err)
        }
    } else {
        fmt.Printf("User data: %s\n", string(data))
    }
    
    // Use transactions
    err = storage.WithTransaction(store, func(tx storage.Transaction) error {
        tx.Set([]byte("user:124"), []byte(`{"name": "Jane"}`))
        tx.Set([]byte("user:125"), []byte(`{"name": "Bob"}`))
        return nil
    })
    if err != nil {
        panic(err)
    }
}
```

### Event System Usage

```go
package main

import (
    "github.com/fintechain/skeleton/pkg/event"
    "fmt"
)

func main() {
    // Create event bus
    bus := event.NewEventBus()
    
    // Subscribe to events
    subscription := event.Subscribe(bus, "user.created", func(data interface{}) {
        fmt.Printf("User created: %v\n", data)
    })
    defer subscription.Unsubscribe()
    
    // Publish an event
    event.Publish(bus, "user.created", map[string]interface{}{
        "id":    "123",
        "name":  "John Doe",
        "email": "john@example.com",
    })
}
```

### Validation Example

```go
package main

import (
    "github.com/fintechain/skeleton/pkg/validation"
    "fmt"
)

func validateUser(name, email, password string) error {
    // Validate required fields
    if err := validation.Required("name", name); err != nil {
        return err
    }
    
    if err := validation.Required("email", email); err != nil {
        return err
    }
    
    // Validate email format
    if err := validation.Email("email", email); err != nil {
        return err
    }
    
    // Validate password length
    if err := validation.MinLength("password", password, 8); err != nil {
        return err
    }
    
    return nil
}

func main() {
    err := validateUser("John Doe", "john@example.com", "password123")
    if err != nil {
        fmt.Printf("Validation error: %v\n", err)
    } else {
        fmt.Println("User data is valid")
    }
}
```

### Cryptographic Operations

```go
package main

import (
    "github.com/fintechain/skeleton/pkg/crypto"
    "fmt"
)

func main() {
    // Hash data
    hash := crypto.HashSHA256("hello world")
    fmt.Printf("SHA256 hash: %s\n", hash)
    
    // Generate HMAC
    hmac := crypto.HMACSHA256("secret-key", "hello world")
    fmt.Printf("HMAC: %s\n", hmac)
    
    // Generate random data
    randomBytes, err := crypto.GenerateRandomBytes(32)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Random bytes: %x\n", randomBytes)
    
    // Generate random string
    randomString, err := crypto.GenerateRandomString(16)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Random string: %s\n", randomString)
    
    // Encode/decode base64
    encoded := crypto.EncodeBase64("hello world")
    decoded, err := crypto.DecodeBase64(encoded)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Encoded: %s, Decoded: %s\n", encoded, decoded)
}
```

## Error Handling

The framework uses a structured error handling approach with domain-specific error codes. All packages provide:

1. **Error Types**: Each package defines its own `Error` type that wraps the component error system
2. **Error Codes**: Predefined constants for common error conditions
3. **Error Creation**: `NewError()` functions for creating structured errors
4. **Error Checking**: `Is*Error()` functions for checking specific error types

### Error Structure

```go
type Error struct {
    Code    string
    Message string
    Cause   error
    Details map[string]string
}
```

### Error Handling Pattern

```go
// Creating errors
err := storage.NewError(storage.ErrKeyNotFound, "user not found", nil).
    WithDetail("userId", "123")

// Checking errors
if storage.IsStorageError(err, storage.ErrKeyNotFound) {
    // Handle key not found specifically
} else if storage.IsStorageError(err, storage.ErrStoreNotFound) {
    // Handle store not found
}
```

## Best Practices

### 1. Component Design

- Keep components focused on a single responsibility
- Use dependency injection for component dependencies
- Implement proper lifecycle management (Initialize/Dispose)
- Use contexts for cancellation and timeouts

### 2. Error Handling

- Always check for errors and handle them appropriately
- Use structured errors with meaningful codes and messages
- Add context details to errors for debugging
- Don't ignore errors - log them at minimum

### 3. Storage Usage

- Use transactions for atomic operations
- Handle `ErrKeyNotFound` gracefully
- Choose appropriate storage engines for your use case
- Configure proper paths and options

### 4. Event System

- Use meaningful topic names
- Handle events asynchronously when possible
- Unsubscribe from events when no longer needed
- Consider event ordering and delivery guarantees

### 5. Configuration

- Use environment-specific configurations
- Validate configuration values
- Provide sensible defaults
- Document configuration options

### 6. Security

- Use cryptographically secure random generation
- Validate all inputs
- Use appropriate hash algorithms (avoid MD5)
- Implement proper authentication and authorization

### 7. Logging

- Use appropriate log levels
- Include context in log messages
- Avoid logging sensitive information
- Configure log output appropriately for environment

### 8. Testing

- Write unit tests for all components
- Use dependency injection for testability
- Mock external dependencies
- Test error conditions

This documentation provides a comprehensive overview of the Fintechain Skeleton Framework's public APIs. For more detailed information about specific implementations, refer to the source code and internal documentation. 