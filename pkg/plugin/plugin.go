// Package plugin provides plugin system interfaces and implementations.
package plugin

import (
	"github.com/fintechain/skeleton/internal/domain/plugin"
	infraPlugin "github.com/fintechain/skeleton/internal/infrastructure/plugin"
)

// Core interfaces
type Plugin = plugin.Plugin
type PluginManager = plugin.PluginManager
type PluginType = plugin.PluginType

// Plugin type constants
const (
	TypeExtension   = plugin.TypeExtension
	TypeIntegration = plugin.TypeIntegration
	TypeMiddleware  = plugin.TypeMiddleware
	TypeConnector   = plugin.TypeConnector
	TypeProcessor   = plugin.TypeProcessor
	TypeAdapter     = plugin.TypeAdapter
)

// Factory functions
var NewManager = infraPlugin.NewManager
