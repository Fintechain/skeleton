// Package runtime provides high-level runtime environment interfaces for the Fintechain Skeleton framework.
//
// The runtime package combines multiple core domain interfaces into a unified API
// that makes it easier for plugins and applications to interact with the system.
package runtime

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/config"
	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/event"
	"github.com/fintechain/skeleton/internal/domain/logging"
	"github.com/fintechain/skeleton/internal/domain/plugin"
)

// RuntimeEnvironment extends the component.System interface with additional
// accessors for commonly used core services. This provides a more comprehensive
// API for plugins to interact with the system.
type RuntimeEnvironment interface {
	// Embed the base System interface
	component.System

	// PluginManager returns the system's plugin manager service.
	// This provides access to plugin lifecycle management.
	// Direct access is guaranteed since the plugin manager is injected as a dependency.
	PluginManager() plugin.PluginManager

	// EventBus returns the system's event bus service.
	// This provides access to publish-subscribe event functionality.
	// Direct access is guaranteed since the event bus is injected as a dependency.
	EventBus() event.EventBusService

	// Logger returns the system's logger.
	// This provides access to structured logging capabilities.
	// Direct access is guaranteed since the logger is injected as a dependency.
	Logger() logging.Logger

	// Configuration returns the system's configuration service.
	// This provides access to application configuration.
	// Direct access is guaranteed since configuration is injected as a dependency.
	Configuration() config.Configuration

	// LoadPlugins loads multiple plugins into the system.
	// This provides batch plugin loading for efficient startup.
	LoadPlugins(ctx context.Context, plugins []plugin.Plugin) error
}
