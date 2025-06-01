// Package plugin provides plugin system interfaces and types.
package plugin

import (
	"github.com/fintechain/skeleton/internal/domain/plugin"
	pluginImpl "github.com/fintechain/skeleton/internal/infrastructure/plugin"
)

// Re-export plugin interfaces
type Plugin = plugin.Plugin
type PluginManager = plugin.PluginManager

// Re-export plugin types
type PluginInfo = plugin.PluginInfo

// Re-export plugin error constants
const (
	ErrPluginNotFound  = plugin.ErrPluginNotFound
	ErrPluginLoad      = plugin.ErrPluginLoad
	ErrPluginUnload    = plugin.ErrPluginUnload
	ErrPluginDiscovery = plugin.ErrPluginDiscovery
)

// NewPluginManager creates a new PluginManager instance with the provided filesystem dependency.
// This factory function provides access to the concrete plugin manager implementation.
func NewPluginManager(filesystem plugin.FileSystem) PluginManager {
	return pluginImpl.NewPluginManager(filesystem)
}
