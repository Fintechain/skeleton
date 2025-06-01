package component

import (
	"fmt"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/registry"
)

// ComponentFactory provides a concrete implementation of the Factory interface.
type ComponentFactory struct {
	registry registry.Registry
}

// NewFactory creates a new component factory instance.
// This constructor accepts a registry interface dependency for component registration and lookup.
func NewFactory(registry registry.Registry) component.Factory {
	return &ComponentFactory{
		registry: registry,
	}
}

// Create creates a component from the given configuration.
// Supports different component types: Basic, Operation, Service, System, Application.
func (f *ComponentFactory) Create(config component.ComponentConfig) (component.Component, error) {
	// Validate configuration
	if config.ID == "" {
		return nil, fmt.Errorf("component ID cannot be empty")
	}
	if config.Name == "" {
		return nil, fmt.Errorf("component name cannot be empty")
	}

	// Check if component already exists in registry
	if f.registry != nil && f.registry.Has(config.ID) {
		return nil, fmt.Errorf("component with ID '%s' already exists", config.ID)
	}

	// Create base component based on type
	var comp component.Component
	switch component.ComponentType(config.Type) {
	case component.TypeBasic, component.TypeOperation, component.TypeService, component.TypeSystem, component.TypeApplication:
		// For now, all types use the same base component implementation
		// Specific type implementations can be added later
		comp = NewBaseComponent(config)
	default:
		return nil, fmt.Errorf("unsupported component type: %v", config.Type)
	}

	// Handle component configuration and property injection
	if len(config.Properties) > 0 {
		// Properties are already included in the component's metadata through the base implementation
		// Additional property injection logic can be added here if needed
	}

	// Register component in registry if registry is available
	if f.registry != nil {
		if err := f.registry.Register(comp); err != nil {
			return nil, fmt.Errorf("failed to register component: %w", err)
		}
	}

	return comp, nil
}
