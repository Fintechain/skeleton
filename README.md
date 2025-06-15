# Fintechain Skeleton Framework

> **Modern Go Framework for Domain-Driven Applications**

A production-ready framework built on **Domain-Driven Design** and **Clean Architecture** principles, featuring automatic dependency injection, dual runtime modes, and a powerful plugin system.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-0.2.0-orange.svg)](CHANGELOG.md)

## 🚀 Quick Start

Get a service running in 30 seconds:

```go
package main

import (
    "github.com/fintechain/skeleton/pkg/runtime"
)

func main() {
    // Start a long-running daemon
    err := runtime.StartDaemon(
        runtime.WithPlugins(myPlugin),
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
- ❌ Complex setup for both CLI tools and long-running services

### The Solution
```go
// Before: Manual wiring, complex setup
config := loadConfig()
logger := newLogger(config)
db := newDatabase(config, logger)
server := newWebServer(config, db, logger)
// ... 50+ lines of manual setup

// After: Declarative, automatic
runtime.StartDaemon(
    runtime.WithPlugins(webPlugin, dbPlugin),
)
```

## 🏗️ Core Concepts

### Dual Runtime Modes

**🔄 Daemon Mode** - Long-running services:
```go
runtime.StartDaemon(runtime.WithPlugins(webServerPlugin))
```

**⚡ Command Mode** - Execute and exit:
```go
result, err := runtime.ExecuteCommand("calculate", inputData,
    runtime.WithPlugins(calculatorPlugin))
```

### Component System
Everything is a **Component** with unified lifecycle:

```go
type Component interface {
    ID() ComponentID
    Name() string
    Initialize(ctx context.Context, system System) error
    Dispose() error
}
```

**Three Component Types:**
- **Services**: Long-running processes (web servers, workers)
- **Operations**: Stateless tasks (calculations, transformations)
- **Components**: Basic entities (database connections, config)

### Plugin Architecture
Plugins orchestrate components and provide functionality:

```go
type MyPlugin struct {
    *component.BaseService
    runtime runtime.RuntimeEnvironment
}

func (p *MyPlugin) Initialize(ctx context.Context, system component.System) error {
    p.runtime = system.(runtime.RuntimeEnvironment)
    registry := system.Registry()
    return registry.Register(myService, myOperation)
}
```

## 📋 Usage Examples

### Web Service
```go
func main() {
    runtime.StartDaemon(
        runtime.WithPlugins(
            webserver.NewPlugin(8080),
            database.NewPlugin("postgres://localhost/myapp"),
        ),
    )
}
```

### CLI Tool
```go
func main() {
    result, err := runtime.ExecuteCommand("process-file", 
        map[string]interface{}{
            "input": os.Args[1],
            "format": "json",
        },
        runtime.WithPlugins(processor.NewPlugin()),
    )
    fmt.Printf("Result: %v\n", result)
}
```

### Testing
```go
func TestCalculator(t *testing.T) {
    result, err := runtime.ExecuteCommand("add", 
        map[string]interface{}{"a": 5, "b": 3},
        runtime.WithPlugins(calculator.NewPlugin()),
    )
    
    require.NoError(t, err)
    assert.Equal(t, 8.0, result["result"])
}
```

## 🏛️ Architecture

Built on **Domain-Driven Design** and **Clean Architecture**:

```
┌─────────────────────────────────────────────────────────────┐
│                    Your Application                         │
├─────────────────────────────────────────────────────────────┤
│                  Public API (pkg/)                         │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────────────┐   │
│  │runtime/ │ │component│ │ plugin/ │ │    config/      │   │
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
- **Automatic DI**: Uber FX handles dependency injection
- **Plugin-Based**: Extend functionality through focused plugins

## 📦 Package Structure

```
pkg/                    # Public API
├── runtime/           # Application runtime and lifecycle
├── component/         # Component system abstractions
├── plugin/           # Plugin system interfaces
├── config/           # Configuration management
├── event/            # Event system
└── context/          # Application context

internal/             # Framework implementation
├── domain/          # Domain interfaces and business logic
├── infrastructure/  # Concrete implementations
└── plugins/         # Built-in plugins

examples/            # Working examples
└── complete-app/    # Full application example
```

## 🔧 Configuration

### Default (Zero Setup)
```go
// Works out of the box with sensible defaults
runtime.StartDaemon(runtime.WithPlugins(myPlugin))
```

### Custom Configuration
```go
import "go.uber.org/fx"

func createConfig() config.Configuration {
    return config.NewMemoryConfigurationWithData(map[string]interface{}{
        "app.port": 8080,
        "db.host":  "localhost",
    })
}

runtime.StartDaemon(
    runtime.WithPlugins(myPlugin),
    runtime.WithOptions(
        fx.Replace(fx.Annotate(createConfig, fx.As(new(config.Configuration)))),
    ),
)
```

## 🧪 Development

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

# Run examples
go run examples/complete-app/main.go command
```

## 📚 Documentation

### 🎓 Getting Started
- **[Runtime Guide](pkg/runtime/README.md)** - Application runtime and modes
- **[Examples](examples/README.md)** - Working code examples

### 🏗️ Development Guides
- **[Plugin Development](docs/PLUGIN_DEVELOPMENT_GUIDE.md)** - Building plugins
- **[Service & Operations](docs/SERVICE_OPERATIONS_DEVELOPMENT_GUIDE.md)** - Component development
- **[Runtime Development](docs/RUNTIME_DEVELOPMENT_GUIDE.md)** - Application architecture

### 🔧 Advanced Topics
- **[Domain Layer](internal/domain/README.md)** - DDD principles and interfaces
- **[Testing Framework](test/unit/README.md)** - Testing patterns and mocks

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guidelines](CONTRIBUTING.md).

### Quick Start
1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Test** your changes (`go test ./...`)
4. **Commit** your changes (`git commit -m 'Add amazing feature'`)
5. **Push** and **open** a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Documentation**: Start with [Runtime Guide](pkg/runtime/README.md)
- **Examples**: Check out [working examples](examples/README.md)
- **Issues**: [GitHub Issues](https://github.com/fintechain/skeleton/issues)
- **Discussions**: [GitHub Discussions](https://github.com/fintechain/skeleton/discussions)

---

**Built with ❤️ by the Fintechain Team**

*Skeleton Framework - Where Clean Architecture meets Developer Productivity* 