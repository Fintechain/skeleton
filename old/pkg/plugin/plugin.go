// Package plugin provides public APIs for the plugin system.
package plugin

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/plugin"
)

// ===== PLUGIN INTERFACES =====

// Plugin is a container for components that extends the system.
type Plugin = plugin.Plugin

// PluginManager handles plugin discovery and lifecycle.
type PluginManager = plugin.PluginManager

// ===== PLUGIN TYPES =====

// PluginInfo provides metadata about a plugin.
type PluginInfo = plugin.PluginInfo

// ===== PLUGIN ERROR CONSTANTS =====

// Common plugin error codes
const (
	ErrPluginNotFound  = "plugin.not_found"
	ErrPluginLoad      = "plugin.load_failed"
	ErrPluginUnload    = "plugin.unload_failed"
	ErrPluginDiscovery = "plugin.discovery_failed"
	ErrPluginConflict  = "plugin.conflict"
	ErrInvalidPlugin   = "plugin.invalid"
)

// ===== ERROR HANDLING =====

// Error represents a domain-specific error from the plugin system.
type Error = component.Error

// NewError creates a new plugin error with the given code, message, and optional cause.
func NewError(code, message string, cause error) *Error {
	return component.NewError(code, message, cause)
}

// IsPluginError checks if an error is a plugin error with the given code.
func IsPluginError(err error, code string) bool {
	return component.IsComponentError(err, code)
}

// ===== PLUGIN CONSTRUCTORS =====

// NewPluginManager creates a new plugin manager with default settings.
// This is the primary way to create a PluginManager instance for managing plugins.
func NewPluginManager() PluginManager {
	return plugin.CreatePluginManager()
}
