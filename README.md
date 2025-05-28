# Skeleton Framework

[![Go Version](https://img.shields.io/badge/Go-1.24.2-blue.svg)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/fintechain/skeleton)
[![Coverage](https://img.shields.io/badge/coverage-85%25-green.svg)](https://github.com/fintechain/skeleton)
[![Go Report Card](https://img.shields.io/badge/go%20report-A+-brightgreen.svg)](https://goreportcard.com/report/github.com/fintechain/skeleton)
[![Documentation](https://img.shields.io/badge/docs-available-blue.svg)](https://github.com/fintechain/skeleton/tree/master/docs)
[![Release](https://img.shields.io/badge/release-v1.0.0-blue.svg)](https://github.com/fintechain/skeleton/releases)

A modular, component-based framework for building scalable Go applications with dependency injection, plugin architecture, and comprehensive storage abstractions.

## Features

- **Component System**: Clean, interface-based component architecture with lifecycle management
- **Dependency Injection**: Built-in FX integration with clean functional options API
- **Plugin Architecture**: Dynamic plugin loading and management
- **Storage Abstraction**: Multi-backend storage system with versioning and transactions
- **Event-Driven**: Decoupled communication through event bus
- **Production Ready**: Comprehensive testing, logging, and error handling

## Quick Start

### Installation

```bash
go mod init your-project
go get github.com/fintechain/skeleton
```

### Basic Usage

```go
package main

import (
    "github.com/fintechain/skeleton/pkg/system"
)

func main() {
    // Start the system with default configuration
    if err := system.StartSystem(); err != nil {
        panic(err)
    }
}
```

### With Custom Configuration

```go
package main

import (
    "github.com/fintechain/skeleton/pkg/system"
    "github.com/fintechain/skeleton/internal/domain/system"
)

func main() {
    config := &system.Config{
        ServiceID:        "my-service",
        EnableOperations: true,
        EnableServices:   true,
        EnablePlugins:    true,
        EnableEventLog:   true,
    }

    err := system.StartSystem(
        system.WithConfig(config),
    )
    if err != nil {
        panic(err)
    }
}
```

## Architecture Overview

The framework is built around several core concepts:

- **Components**: Building blocks with identity, lifecycle, and metadata
- **Registry**: Central repository for component discovery and management
- **Services**: Long-running components with start/stop lifecycle
- **Operations**: Discrete units of work with inputs and outputs
- **Plugins**: Containers for extending the system dynamically
- **Storage**: Multi-backend storage with versioning and transactions

## Development

### Prerequisites

- Go 1.24.2 or later
- Make (for build automation)

### Building

```bash
# Build all binaries
make build

# Build specific components
make build-server
make build-client
make build-fx-example
```

### Testing

```bash
# Run all tests
make test

# Run specific test suites
make test-unit
make test-integration
make test-component
make test-storage

# Generate coverage report
make coverage
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run full development cycle
make dev
```

### Running Examples

```bash
# Run the FX integration example
make run-fx-example

# Run the server
make run-server

# Run the client
make run-client
```

## Available Make Targets

The project includes a comprehensive Makefile with the following target categories:

- **Build**: `build`, `build-server`, `build-client`, `build-fx-example`, `install`
- **Test**: `test`, `test-unit`, `test-integration`, `test-component`, `test-storage`
- **Coverage**: `coverage`, `coverage-unit`, `coverage-integration`, `coverage-show`
- **Quality**: `lint`, `lint-fix`, `fmt`, `vet`, `mod-tidy`
- **Development**: `run-fx-example`, `run-server`, `run-client`, `dev`, `ci`
- **Tools**: `install-tools`, `check-tools`, `mocks`, `bench`
- **Cleanup**: `clean`, `clean-mocks`
- **Release**: `version`, `tag`, `docker-build`

Run `make help` to see all available targets with descriptions.

## Project Structure

```
skeleton/
├── cmd/                    # Application entry points
│   ├── fx-example/        # FX integration example
│   ├── server/            # Server application
│   └── client/            # Client application
├── pkg/                   # Public API packages
│   └── system/            # System startup API
├── internal/              # Private implementation
│   ├── domain/            # Domain models and interfaces
│   └── infrastructure/    # Infrastructure implementations
├── test/                  # Test suites
│   ├── integration/       # Integration tests
│   └── unit/              # Unit tests
├── docs/                  # Documentation
├── configs/               # Configuration files
├── deployments/           # Deployment configurations
└── scripts/               # Build and utility scripts
```

## Documentation

Comprehensive documentation is available in the `docs/` directory:

- **[Component System Reference](docs/COMPONENT_SYSTEM_REFERENCE.md)**: Complete guide to the component architecture
- **[FX Integration Implementation](docs/FX_INTEGRATION_IMPLEMENTATION.md)**: Dependency injection and system startup
- **[Storage System Implementation](docs/STORAGE_SYSTEM_IMPLEMENTATION_PLAN.md)**: Multi-backend storage architecture

## Key Components

### Component System

The framework provides a flexible component system with:

- **Base Components**: Foundation for all components with identity and lifecycle
- **Registry**: Component discovery and dependency management
- **Factory**: Component creation from configuration
- **Lifecycle Management**: Proper initialization and cleanup

### Storage System

Multi-backend storage abstraction supporting:

- **Multiple Engines**: In-memory, file-based, LevelDB, IAVL tree
- **Transactions**: Atomic operations across multiple stores
- **Versioning**: Immutable snapshots and rollback capabilities
- **Range Queries**: Efficient iteration over key ranges

### Plugin Architecture

Dynamic plugin system featuring:

- **Plugin Manager**: Registration and lifecycle management
- **Component Discovery**: Automatic component registration from plugins
- **Isolation**: Clean separation between plugin and core system

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test`)
5. Run linting (`make lint`)
6. Commit your changes (`git commit -m 'Add amazing feature'`)
7. Push to the branch (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## Development Workflow

```bash
# Install development tools
make install-tools

# Full development cycle
make dev

# CI pipeline (what runs in CI)
make ci

# Clean build artifacts
make clean
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Support

For questions, issues, or contributions:

- Create an issue on GitHub
- Check the documentation in the `docs/` directory
- Review the examples in the `cmd/` directory 