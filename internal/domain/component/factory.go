package component

// ComponentConfig defines the configuration for creating a component.
// It includes all identity and component properties directly.
type ComponentConfig struct {
	ID           ComponentID            // Unique identifier
	Name         string                 // Descriptive name
	Type         ComponentType          // Component type
	Description  string                 // Optional description
	Version      string                 // Component version (defaults to "1.0.0" if not specified)
	Dependencies []ComponentID          // Component dependencies
	Properties   map[string]interface{} // Component-specific properties
}

// Factory creates components from configuration with dependency injection support.
type Factory interface {
	Create(config ComponentConfig) (Component, error)
}
