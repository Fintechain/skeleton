// Package component provides component interfaces and types.
package component

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/registry"
	componentImpl "github.com/fintechain/skeleton/internal/infrastructure/component"
)

// Re-export dependency-aware component interface
type DependencyAwareComponent = component.DependencyAwareComponent

// Factory functions

// NewDependencyAwareComponent creates a new dependency-aware component.
func NewDependencyAwareComponent(base Component, registry registry.Registry) DependencyAwareComponent {
	return componentImpl.NewDependencyAwareComponent(base, registry)
}
