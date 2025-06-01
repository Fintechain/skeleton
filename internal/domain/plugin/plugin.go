// Package plugin provides functionality for plugins in the system.
package plugin

import (
	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/registry"
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
// It composes Identifiable to inherit ID, Name, Description, and Version methods.
type Plugin interface {
	registry.Identifiable

	// Lifecycle
	Load(ctx context.Context, registrar registry.Registry) error
	Unload(ctx context.Context) error
}

// PluginManager handles plugin discovery and lifecycle.
type PluginManager interface {
	// Discovery
	Discover(ctx context.Context, location string) ([]PluginInfo, error)

	// Lifecycle
	Load(ctx context.Context, id string, registrar registry.Registry) error
	Unload(ctx context.Context, id string) error

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
