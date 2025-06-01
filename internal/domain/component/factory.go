package component

import (
	"github.com/fintechain/skeleton/internal/domain/registry"
)

// ComponentConfig defines the configuration for creating a component.
// It composes IdentifiableConfig to inherit core identity properties.
type ComponentConfig struct {
	registry.IdentifiableConfig
	// Component-specific configuration properties can be added here in the future
}

// NewComponentConfig creates a new ComponentConfig with the given parameters.
// This provides a convenient constructor that ensures proper initialization.
func NewComponentConfig(id, name string, componentType ComponentType, description string) ComponentConfig {
	return ComponentConfig{
		IdentifiableConfig: registry.IdentifiableConfig{
			ID:           id,
			Name:         name,
			Type:         registry.IdentifiableType(componentType),
			Description:  description,
			Version:      "1.0.0", // Default version
			Dependencies: []string{},
			Properties:   make(map[string]interface{}),
		},
	}
}

// Factory creates components from configuration.
type Factory interface {
	Create(config ComponentConfig) (Component, error)
}
