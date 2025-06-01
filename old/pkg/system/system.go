// Package system provides the main entry point for system lifecycle management.
// This package focuses solely on starting and configuring the Fintechain Skeleton system
// with various options through a functional options pattern.
//
// Example usage:
//
//	err := system.StartSystem(
//	    system.WithConfig(config),
//	    system.WithPlugins(plugins),
//	    system.WithRegistry(registry),
//	)
package system

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/plugin"
	"github.com/fintechain/skeleton/internal/domain/storage"
	"github.com/fintechain/skeleton/internal/infrastructure/event"
	"github.com/fintechain/skeleton/internal/infrastructure/system"
)

// StartSystem starts the system with the given options
func StartSystem(options ...Option) error {
	config := &system.SystemConfig{}
	for _, option := range options {
		option(config)
	}
	return system.StartWithFx(config)
}

// Option is a functional option for configuring the system
type Option func(*system.SystemConfig)

// WithConfig sets the system configuration
func WithConfig(config *InternalConfig) Option {
	return func(sc *system.SystemConfig) {
		sc.Config = config
	}
}

// WithPlugins sets the plugins to load
func WithPlugins(plugins []plugin.Plugin) Option {
	return func(sc *system.SystemConfig) {
		sc.Plugins = plugins
	}
}

// WithRegistry sets the component registry
func WithRegistry(registry component.Registry) Option {
	return func(sc *system.SystemConfig) {
		sc.Registry = registry
	}
}

// WithPluginManager sets the plugin manager
func WithPluginManager(pluginMgr plugin.PluginManager) Option {
	return func(sc *system.SystemConfig) {
		sc.PluginMgr = pluginMgr
	}
}

// WithEventBus sets the event bus
func WithEventBus(eventBus event.EventBus) Option {
	return func(sc *system.SystemConfig) {
		sc.EventBus = eventBus
	}
}

// WithMultiStore sets the multistore
func WithMultiStore(multiStore storage.MultiStore) Option {
	return func(sc *system.SystemConfig) {
		sc.MultiStore = multiStore
	}
}
