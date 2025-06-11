// Package component provides interfaces and types for the component system.
package component

// Registry manages component registration and lookup.
// Provides centralized component directory with type-safe access.
type Registry interface {
	// Register adds a component to the registry
	Register(component Component) error

	// Get retrieves a component by ID
	Get(id ComponentID) (Component, error)

	// GetByType retrieves all components of a given type
	GetByType(typ ComponentType) ([]Component, error)

	// Find returns all components that match the given predicate function
	Find(predicate func(Component) bool) ([]Component, error)

	// Has checks if a component exists
	Has(id ComponentID) bool

	// List returns all registered component IDs
	List() []ComponentID

	// Unregister removes a component from the registry
	Unregister(id ComponentID) error

	// Count returns the number of registered components
	Count() int

	// Clear removes all components from the registry
	Clear() error
}
