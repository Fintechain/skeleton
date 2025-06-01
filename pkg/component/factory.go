// Package component provides component interfaces and types.
package component

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/registry"
	componentImpl "github.com/fintechain/skeleton/internal/infrastructure/component"
)

// Re-export factory interface
type Factory = component.Factory

// Re-export factory types
type ComponentConfig = component.ComponentConfig

// Re-export constructor function
var NewComponentConfig = component.NewComponentConfig

// Factory functions

// NewBaseComponent creates a new base component instance.
func NewBaseComponent(config ComponentConfig) Component {
	return componentImpl.NewBaseComponent(config)
}

// NewFactory creates a new component factory instance.
func NewFactory(registry registry.Registry) Factory {
	return componentImpl.NewFactory(registry)
}
