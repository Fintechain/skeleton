// Package plugin provides functionality for plugins in the system.
package plugin

import (
	"github.com/fintechain/skeleton/internal/domain/component"
)

// PluginInfo provides metadata about a plugin.
type PluginInfo struct {
	ID          string                 // Unique identifier
	Name        string                 // Human-readable name
	Version     string                 // Semantic version
	Description string                 // Plugin description
	Author      string                 // Plugin author/maintainer
	Metadata    map[string]interface{} // Additional metadata
}

// Plugin is a container for components that extends the system.
type Plugin interface {
	// Identity
	ID() string
	Version() string

	// Lifecycle
	Load(ctx component.Context, registry component.Registry) error
	Unload(ctx component.Context) error

	// Components
	Components() []component.Component
}

// PluginManager handles plugin discovery and lifecycle.
type PluginManager interface {
	// Discovery
	Discover(ctx component.Context, location string) ([]PluginInfo, error)

	// Lifecycle
	Load(ctx component.Context, id string, registry component.Registry) error
	Unload(ctx component.Context, id string) error

	// Information
	ListPlugins() []PluginInfo
	GetPlugin(id string) (Plugin, error)
}

// Common error codes for plugin operations
const (
	ErrPluginNotFound  = "plugin.not_found"
	ErrPluginLoad      = "plugin.load_failed"
	ErrPluginUnload    = "plugin.unload_failed"
	ErrPluginDiscovery = "plugin.discovery_failed"
)
