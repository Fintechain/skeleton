# Fx Integration

This document explains how to use the fx integration to start a system with plugins.

## Overview

The fx integration provides a simple way to start a system using the Uber Fx dependency injection framework. It abstracts the complexity of fx from the client and provides a clean API.

## Usage

### Basic Usage

```go
package main

import (
    "log"
    
    "github.com/fintechain/skeleton/internal/domain/plugin"
    "github.com/fintechain/skeleton/internal/domain/storage"
    "github.com/fintechain/skeleton/internal/infrastructure/system"
)

func main() {
    // Create configuration
    config := &system.Config{
        ServiceID: "my-system",
        StorageConfig: storage.MultiStoreConfig{
            RootPath:      "./data",
            DefaultEngine: "memory",
        },
    }
    
    // Create plugins (optional)
    plugins := []plugin.Plugin{
        // Add your plugins here
    }
    
    // Start the system with fx
    if err := system.StartWithFx(config, plugins); err != nil {
        log.Fatalf("Failed to start system: %v", err)
    }
}
```

### Configuration

The `system.Config` struct supports the following fields:

- `ServiceID`: Unique identifier for the system service
- `StorageConfig`: Configuration for the multi-store
  - `RootPath`: Root directory for storage
  - `DefaultEngine`: Default storage engine to use (e.g., "memory")

### Plugins

Plugins must implement the `plugin.Plugin` interface:

```go
type Plugin interface {
    ID() string
    Version() string
    Load(ctx component.Context, registry component.Registry) error
    Unload(ctx component.Context) error
    Components() []component.Component
}
```

### Default Components

The fx integration automatically creates default implementations for:

- Component Registry
- Plugin Manager
- Event Bus
- Multi-Store (memory-based)

These can be overridden by providing custom implementations in the system configuration.

## Example

See `cmd/fx-example/main.go` for a complete working example.

## Architecture

The fx integration follows these principles:

1. **Abstraction**: Fx is completely hidden from the client
2. **Defaults**: Sensible defaults are provided for all components
3. **Flexibility**: Custom implementations can be provided
4. **Lifecycle**: Proper initialization and startup sequence

## Dependencies

The fx integration manages the following dependencies:

- Logger (standard logger with Info level)
- Component Registry
- Plugin Manager
- Event Bus
- Multi-Store
- System Configuration

All dependencies are automatically wired together using fx's dependency injection. 