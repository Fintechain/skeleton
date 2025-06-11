// Package plugin provides plugin-related interfaces and types for the Fintechain Skeleton framework.
package plugin

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
)

// PluginType represents the type of plugin
type PluginType string

// Plugin type constants
const (
	TypeExtension   PluginType = "extension"   // Extends core functionality
	TypeIntegration PluginType = "integration" // Integrates with external systems
	TypeMiddleware  PluginType = "middleware"  // Provides middleware functionality
	TypeConnector   PluginType = "connector"   // Connects to external services
	TypeProcessor   PluginType = "processor"   // Processes data/events
	TypeAdapter     PluginType = "adapter"     // Adapts interfaces/protocols
)

// Plugin represents a dynamically loadable component.
// This is the core plugin interface without lifecycle management.
type Plugin interface {
	component.Service

	// Plugin author/maintainer
	Author() string

	// Plugin type
	PluginType() PluginType
}

// PluginManager manages plugin lifecycle and discovery.
// This is the core plugin manager interface without service lifecycle.
type PluginManager interface {
	component.Service

	// Plugin registry operations
	Add(pluginID component.ComponentID, plugin Plugin) error
	Remove(pluginID component.ComponentID) error

	// Plugin execution
	StartPlugin(ctx context.Context, pluginID component.ComponentID) error
	StopPlugin(ctx context.Context, pluginID component.ComponentID) error

	// Plugin queries
	GetPlugin(pluginID component.ComponentID) (Plugin, error)
	ListPlugins() []component.ComponentID
}
