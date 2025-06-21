# Fintechain Skeleton Framework

> **Modern Go Framework for Domain-Driven Applications**

A production-ready framework built on **Domain-Driven Design** and **Clean Architecture** principles, featuring a clean Builder API, dual runtime modes, and a powerful plugin system.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.24-blue.svg)](https://golang.org/)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Version](https://img.shields.io/badge/version-0.4.0-orange.svg)](CHANGELOG.md)

## 🚀 Quick Start

Get a service running in 30 seconds:

```go
package main

import (
    "log"
    "github.com/fintechain/skeleton/pkg/runtime"
)

func main() {
    // Start a long-running daemon
    err := runtime.NewBuilder().
        WithPlugins(myPlugin).
        BuildDaemon()
    if err != nil {
        log.Fatal(err)
    }
}
```

**Installation:**
```bash
go mod init myapp
go get github.com/fintechain/skeleton@v0.4.0
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

// After: Clean Builder API
runtime.NewBuilder().
    WithPlugins(webPlugin, dbPlugin).
    BuildDaemon()
```

## 🏗️ Core Concepts

### Builder Pattern API

**Simple and Intuitive:**
```go
// Build a daemon application
runtime.NewBuilder().
    WithPlugins(plugin1, plugin2).
    WithConfig(customConfig).
    WithLogger(customLogger).
    BuildDaemon()

// Build a command application
result, err := runtime.NewBuilder().
    WithPlugins(plugin1, plugin2).
    BuildCommand("operation-id", inputData)
```

### Dual Runtime Modes

**🔄 Daemon Mode** - Long-running services:
```go
runtime.NewBuilder().
    WithPlugins(webServerPlugin).
    BuildDaemon()
```

**⚡ Command Mode** - Execute and exit:
```go
result, err := runtime.NewBuilder().
    WithPlugins(calculatorPlugin).
    BuildCommand("calculate", inputData)
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
    err := runtime.NewBuilder().
        WithPlugins(
            webserver.NewPlugin(8080),
            database.NewPlugin("postgres://localhost/myapp"),
        ).
        BuildDaemon()
    if err != nil {
        log.Fatal(err)
    }
}
```

### CLI Tool
```go
func main() {
    if len(os.Args) < 2 {
        log.Fatal("Usage: app <input-file>")
    }
    
    result, err := runtime.NewBuilder().
        WithPlugins(processor.NewPlugin()).
        BuildCommand("process-file", map[string]interface{}{
            "input":  os.Args[1],
            "format": "json",
        })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Processed %d records\n", result["count"])
}
```

### Testing
```go
func TestCalculator(t *testing.T) {
    result, err := runtime.NewBuilder().
        WithPlugins(calculator.NewPlugin()).
        BuildCommand("add", map[string]interface{}{
            "a": 5, 
            "b": 3,
        })
    
    require.NoError(t, err)
    assert.Equal(t, 8.0, result["result"])
}
```

### Custom Configuration
```go
func createAppConfig() config.Configuration {
    settings := map[string]interface{}{
        "app.name":        "My Application",
        "app.port":        8080,
        "database.host":   "localhost",
        "database.port":   5432,
        "log.level":       "info",
    }
    return infraConfig.NewMemoryConfigurationWithData(settings)
}

func main() {
    err := runtime.NewBuilder().
        WithPlugins(myWebPlugin, myDatabasePlugin).
        WithConfig(createAppConfig()).
        BuildDaemon()
    if err != nil {
        log.Fatal(err)
    }
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
- **Builder Pattern**: Clean, fluent API for application construction
- **Plugin-Based**: Extend functionality through focused plugins

## 📦 Package Structure

```
pkg/                    # Public API
├── runtime/           # Application runtime and Builder API
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
runtime.NewBuilder().
    WithPlugins(myPlugin).
    BuildDaemon()
```

### Custom Configuration
```go
import (
    "github.com/fintechain/skeleton/internal/domain/config"
    infraConfig "github.com/fintechain/skeleton/internal/infrastructure/config"
)

func createConfig() config.Configuration {
    return infraConfig.NewMemoryConfigurationWithData(map[string]interface{}{
        "app.port": 8080,
        "db.host":  "localhost",
    })
}

runtime.NewBuilder().
    WithPlugins(myPlugin).
    WithConfig(createConfig()).
    BuildDaemon()
```

### Advanced Customization
```go
// Custom logger, event bus, and configuration
customConfig := infraConfig.NewMemoryConfigurationWithData(settings)
customLogger := infraLogging.NewConsoleLogger()
customEventBus := infraEvent.NewInMemoryEventBus()

runtime.NewBuilder().
    WithPlugins(myPlugin).
    WithConfig(customConfig).
    WithLogger(customLogger).
    WithEventBus(customEventBus).
    BuildDaemon()
```

## 🚀 What's New in v0.4.0

### ✅ **NEW: Builder Pattern API**
- Clean, fluent API: `runtime.NewBuilder().WithPlugins().BuildDaemon()`
- Custom dependency injection: `WithConfig()`, `WithLogger()`, `WithEventBus()`
- Simplified application construction

### ❌ **REMOVED: FX Dependency Injection**
- Eliminated complex FX dependency wiring
- Removed `runtime.StartDaemon()` and `runtime.ExecuteCommand()` legacy functions
- Cleaner, more predictable dependency management

### 🔧 **IMPROVED: Developer Experience**
- Simpler API with better error messages
- Faster startup times without FX overhead
- More intuitive component lifecycle management

## 🧪 Development

### Prerequisites
- Go 1.24 or higher
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
go run examples/complete-app/main.go daemon
go run examples/complete-app/main.go command calculate
```

## 📚 Documentation

### 🎓 Getting Started
- **[Runtime Development Guide](docs/RUNTIME_DEVELOPMENT_GUIDE.md)** - Application runtime and Builder API
- **[Examples](examples/README.md)** - Working code examples

### 🏗️ Development Guides
- **[Plugin Development Guide](docs/PLUGIN_DEVELOPMENT_GUIDE.md)** - Building plugins
- **[Service & Operations Guide](docs/SERVICE_OPERATIONS_DEVELOPMENT_GUIDE.md)** - Component development
- **[System API Reference](docs/SYSTEM_API_REFERENCE.md)** - Complete API documentation

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

- **Documentation**: Start with [Runtime Development Guide](docs/RUNTIME_DEVELOPMENT_GUIDE.md)
- **Examples**: Check out [working examples](examples/README.md)
- **Issues**: [GitHub Issues](https://github.com/fintechain/skeleton/issues)
- **Discussions**: [GitHub Discussions](https://github.com/fintechain/skeleton/discussions)

---

**Built with ❤️ by the Fintechain Team**

*Skeleton Framework - Where Clean Architecture meets Developer Productivity* 