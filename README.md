# Fintechain Skeleton Framework

> **Modern, Domain-Driven Application Framework for Go**

A production-ready framework built on **Domain-Driven Design** and **Clean Architecture** principles, featuring automatic dependency injection, pluggable components, and dual runtime modes for both long-running services and CLI applications.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](#)

## 🚀 Quick Start

Get a web service running in under 30 seconds:

```go
package main

import (
    "github.com/fintechain/skeleton/pkg/fx"
    "github.com/fintechain/skeleton/pkg/plugin"
)

func main() {
    err := fx.StartDaemon(
        fx.WithPlugins(&WebServerPlugin{Port: 8080}),
    )
    if err != nil {
        panic(err)
    }
}
```

**Installation:**
```bash
go mod init myapp
go get github.com/fintechain/skeleton
```

## 🎯 Why Skeleton Framework?

### The Problem
Building production Go applications often means:
- ❌ Manual dependency wiring and lifecycle management
- ❌ Scattered configuration and service discovery
- ❌ Tight coupling between business logic and infrastructure
- ❌ Difficult testing due to hard dependencies

### The Solution
```go
// Before: Manual wiring, error-prone setup
config := loadConfig()
logger := newLogger(config)
db := newDatabase(config, logger)
server := newWebServer(config, db, logger)
// ... 50+ lines of manual setup

// After: Declarative, automatic, reliable
fx.StartDaemon(
    fx.WithPlugins(&ConfigPlugin{}, &DatabasePlugin{}, &WebServerPlugin{}),
)
```

### Key Benefits

| Feature | Traditional Approach | Skeleton Framework |
|---------|---------------------|-------------------|
| **Dependency Injection** | Manual, error-prone | Automatic, type-safe |
| **Architecture** | Coupled, monolithic | Clean, domain-driven |
| **Testing** | Complex mocking | Built-in test support |
| **Lifecycle** | Manual start/stop | Automatic management |
| **Configuration** | Scattered | Centralized, hierarchical |
| **Deployment** | Single mode | Daemon + CLI modes |

## 🏗️ Core Concepts

### Component System
Everything is a **Component** with unified lifecycle management:

```go
type Component interface {
    // Identity methods
    ID() ComponentID
    Name() string
    Description() string
    Version() string
    
    // Component-specific methods
    Type() ComponentType
    Metadata() Metadata
    Initialize(ctx context.Context, system System) error
    Dispose() error
}
```

**Three Component Types:**
- **Components**: Basic entities (database connections, config loaders)
- **Operations**: Executable tasks with input/output (calculations, transformations)
- **Services**: Long-running processes (web servers, message processors)

### Dual Runtime Modes

**🔄 Daemon Mode** - Long-running services:
```go
fx.StartDaemon(fx.WithPlugins(&WebServerPlugin{}))
```

**⚡ Command Mode** - CLI tools and batch processing:
```go
result, err := fx.ExecuteCommand("process-data", inputData,
    fx.WithPlugins(&ProcessorPlugin{}))
```

### Plugin Architecture
Extend functionality through focused, testable plugins:

```go
type CalculatorPlugin struct {
    *component.BaseService
}

func (p *CalculatorPlugin) Initialize(ctx context.Context, system component.System) error {
    registry := system.Registry()
    return registry.Register(&AddOperation{})
}
```

## 📋 Usage Examples

### Web Service Application
```go
func main() {
    fx.StartDaemon(
        fx.WithPlugins(
            &ConfigPlugin{},
            &DatabasePlugin{URL: "postgres://localhost/myapp"},
            &WebServerPlugin{Port: 8080},
            &MetricsPlugin{},
        ),
    )
}
```

### CLI Tool
```go
func main() {
    operation := os.Args[1]
    input := parseArgs(os.Args[2:])
    
    result, err := fx.ExecuteCommand(operation, input,
        fx.WithPlugins(&DataProcessorPlugin{}),
    )
    
    fmt.Printf("Result: %v\n", result)
}
```

### Testing
```go
func TestCalculator(t *testing.T) {
    result, err := fx.ExecuteCommand("add", map[string]any{
        "a": 5, "b": 3,
    }, fx.WithPlugins(&CalculatorPlugin{}))
    
    require.NoError(t, err)
    assert.Equal(t, 8.0, result["result"])
}
```

## 🏛️ Architecture

Built on **Domain-Driven Design** and **Clean Architecture** principles:

```
┌─────────────────────────────────────────────────────────────┐
│                    Your Application                         │
├─────────────────────────────────────────────────────────────┤
│                  Public API (pkg/)                         │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────────────┐   │
│  │   fx/   │ │runtime/ │ │component│ │    plugin/      │   │
│  └─────────┘ └─────────┘ └─────────┘ └─────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│                Infrastructure Layer                        │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────────────┐   │
│  │ Storage │ │ Events  │ │ Config  │ │    Logging      │   │
│  └─────────┘ └─────────┘ └─────────┘ └─────────────────┘   │
├─────────────────────────────────────────────────────────────┤
│                   Domain Layer                             │
│           (Interfaces, Business Logic)                     │
└─────────────────────────────────────────────────────────────┘
```

**Key Principles:**
- **Domain Independence**: Business logic doesn't depend on infrastructure
- **Dependency Inversion**: Infrastructure implements domain interfaces
- **Single Responsibility**: Each component has one clear purpose
- **Testability**: Easy mocking through interface-based design

## 📚 Documentation

### 🎓 Getting Started
- **[Public API Guide](pkg/README.md)** - Complete API reference and usage patterns
- **[Examples](examples/README.md)** - Traditional vs Modern approaches with working code
- **[FX Integration Guide](pkg/fx/README.md)** - Deep dive into dependency injection

### 🏗️ Architecture & Development
- **[Domain Layer Guide](internal/domain/README.md)** - DDD principles and interfaces
- **[Testing Framework](test/unit/README.md)** - Comprehensive testing patterns and mocks

### 🔌 Advanced Topics
- **[Plugin Development](pkg/plugin/README.md)** - Building and distributing plugins
- **[Storage Systems](internal/domain/README.md#storage-system)** - Multi-backend storage abstraction
- **[Event-Driven Architecture](internal/domain/README.md#event-system)** - Pub/sub messaging patterns

## 🛠️ Development

### Prerequisites
- Go 1.21 or higher
- Git

### Setup
```bash
git clone https://github.com/fintechain/skeleton.git
cd skeleton
go mod download
```

### Running Tests
```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run specific test suites
go test ./test/unit/infrastructure/...
go test ./test/unit/pkg/...
```

### Project Structure
```
├── pkg/                    # Public API
├── internal/
│   ├── domain/            # Domain interfaces and business logic
│   ├── infrastructure/    # Concrete implementations
│   └── fx/               # FX integration logic
├── examples/              # Working examples and tutorials
├── test/                  # Comprehensive test suite
└── docs/                  # Additional documentation
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md) for details.

### Quick Contribution Guide
1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Test** your changes (`go test ./...`)
4. **Commit** your changes (`git commit -m 'Add amazing feature'`)
5. **Push** to the branch (`git push origin feature/amazing-feature`)
6. **Open** a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Documentation**: Start with the [Public API Guide](pkg/README.md)
- **Examples**: Check out [working examples](examples/README.md)
- **Issues**: [GitHub Issues](https://github.com/fintechain/skeleton/issues)
- **Discussions**: [GitHub Discussions](https://github.com/fintechain/skeleton/discussions)

---

**Built with ❤️ by the Fintechain Team**

*Skeleton Framework - Where Clean Architecture meets Developer Productivity* 