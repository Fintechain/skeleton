// Package system provides the core system interface for resource access and operations.
package system

import (
	"context"

	"go.uber.org/fx"

	"github.com/fintechain/skeleton/internal/domain/config"
	domainContext "github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/event"
	"github.com/fintechain/skeleton/internal/domain/plugin"
	"github.com/fintechain/skeleton/internal/domain/registry"
	"github.com/fintechain/skeleton/internal/domain/storage"
	"github.com/fintechain/skeleton/internal/domain/system"
	systemImpl "github.com/fintechain/skeleton/internal/infrastructure/system"
)

// Re-export system interfaces
type Context = domainContext.Context
type System = system.System

// Re-export fx bootstrap types
type SystemConfig = systemImpl.SystemConfig
type Config = systemImpl.Config

// Re-export system error constants
const (
	ErrSystemNotInitialized = system.ErrSystemNotInitialized
	ErrSystemNotStarted     = system.ErrSystemNotStarted
	ErrOperationNotFound    = system.ErrOperationNotFound
	ErrOperationFailed      = system.ErrOperationFailed
	ErrServiceNotFound      = system.ErrServiceNotFound
	ErrServiceStart         = system.ErrServiceStart
	ErrServiceStop          = system.ErrServiceStop
)

// NewSystem creates a new System instance with the provided dependencies.
// This factory function provides access to the concrete system implementation.
func NewSystem(
	registry registry.Registry,
	pluginManager plugin.PluginManager,
	eventBus event.EventBus,
	configuration config.Configuration,
	store storage.MultiStore,
) System {
	return systemImpl.NewSystem(registry, pluginManager, eventBus, configuration, store)
}

// StartWithFx is the public API for starting a system with fx.
// This abstracts fx from the client and provides a simple interface.
// The function blocks until the system receives a shutdown signal.
func StartWithFx(config *SystemConfig) error {
	return systemImpl.StartWithFx(config)
}

// StartWithFxAndContext starts the system using the fx framework with context control.
// This is useful for testing where you need to control the system lifecycle.
// Returns the fx.App instance so the caller can control its lifecycle.
func StartWithFxAndContext(ctx context.Context, config *SystemConfig) (*fx.App, error) {
	return systemImpl.StartWithFxAndContext(ctx, config)
}
