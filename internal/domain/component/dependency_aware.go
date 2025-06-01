package component

import (
	"github.com/fintechain/skeleton/internal/domain/registry"
)

// DependencyAwareComponent is an interface for components that have dependencies on other components.
type DependencyAwareComponent interface {
	// Dependencies returns the IDs of components this component depends on.
	Dependencies() []string

	// AddDependency adds a component dependency.
	AddDependency(id string)

	// RemoveDependency removes a component dependency.
	RemoveDependency(id string)

	// HasDependency checks if this component depends on the component with the given ID.
	HasDependency(id string) bool

	// ResolveDependency resolves a dependency to a component instance.
	ResolveDependency(id string, registrar registry.Registry) (Component, error)

	// ResolveDependencies resolves all dependencies to component instances.
	ResolveDependencies(registrar registry.Registry) (map[string]Component, error)
}
