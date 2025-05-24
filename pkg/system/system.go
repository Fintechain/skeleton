package system

import (
	"github.com/ebanfa/skeleton/internal/domain/component"
	"github.com/ebanfa/skeleton/internal/domain/plugin"
	"github.com/ebanfa/skeleton/internal/domain/storage"
	"github.com/ebanfa/skeleton/internal/infrastructure/event"
	"github.com/ebanfa/skeleton/internal/infrastructure/system"
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
func WithConfig(config *system.Config) Option {
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
