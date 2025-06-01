// Package registry provides a generic registry for identifiable items.
package registry

// Identifiable represents any item that can be stored in the registry.
// This minimal interface breaks the circular dependency with Component.
// It includes the core identity properties that all registry items should have.
type Identifiable interface {
	ID() string
	Name() string
	Description() string
	Version() string
}

// ComponentType represents the type of a component.
// This minimal interface breaks the circular dependency with Component.
type IdentifiableType string

const (
	// Basic component types
	TypeBasic       IdentifiableType = "basic"
	TypeOperation   IdentifiableType = "operation"
	TypeService     IdentifiableType = "service"
	TypeSystem      IdentifiableType = "system"
	TypeApplication IdentifiableType = "application"
)

// ComponentConfig defines the configuration for creating a component.
// This minimal interface breaks the circular dependency with Component.
type IdentifiableConfig struct {
	ID           string                 // Unique identifier
	Name         string                 // Descriptive name
	Type         IdentifiableType       // Component type
	Description  string                 // Optional description
	Version      string                 // Component version (defaults to "1.0.0" if not specified)
	Dependencies []string               // Component dependencies
	Properties   map[string]interface{} // Component-specific properties
}

// Factory creates components from configuration.
type IdentifiableFactory interface {
	Create(config IdentifiableConfig) (Identifiable, error)
}

// Registry provides a generic storage mechanism for identifiable items.
// It doesn't depend on any domain-specific interfaces, making it pure infrastructure.
type Registry interface {
	// Register stores an item in the registry.
	Register(item Identifiable) error

	// Get retrieves an item by its ID.
	Get(id string) (Identifiable, error)

	// List returns all registered items.
	List() []Identifiable

	// Remove removes an item from the registry.
	Remove(id string) error

	// Has checks if an item with the given ID exists.
	Has(id string) bool

	// Count returns the number of registered items.
	Count() int

	// Clear removes all items from the registry.
	Clear()
}

// Common registry error codes
const (
	ErrItemNotFound      = "registry.item_not_found"
	ErrItemAlreadyExists = "registry.item_already_exists"
	ErrInvalidItem       = "registry.invalid_item"
)
