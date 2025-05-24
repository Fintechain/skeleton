package component

import (
	"github.com/ebanfa/skeleton/internal/infrastructure/logging"
)

// ComponentConfig defines the configuration for creating a component.
type ComponentConfig struct {
	ID           string                 // Unique identifier
	Name         string                 // Descriptive name
	Type         ComponentType          // Component type
	Description  string                 // Optional description
	Dependencies []string               // Component dependencies
	Properties   map[string]interface{} // Component-specific properties
}

// Factory creates components from configuration.
type Factory interface {
	Create(config ComponentConfig) (Component, error)
}

// DefaultFactory provides a standard implementation of the Factory interface.
type DefaultFactory struct {
	// typeCreators maps component types to creation functions
	typeCreators map[ComponentType]func(config ComponentConfig) (Component, error)
	logger       logging.Logger
}

// DefaultFactoryOptions contains options for creating a DefaultFactory.
type DefaultFactoryOptions struct {
	Logger logging.Logger
}

// NewFactory creates a new component factory.
func NewFactory(options DefaultFactoryOptions) *DefaultFactory {
	factory := &DefaultFactory{
		typeCreators: make(map[ComponentType]func(config ComponentConfig) (Component, error)),
		logger:       options.Logger,
	}

	// Register default creator for basic components
	factory.RegisterTypeCreator(TypeBasic, func(config ComponentConfig) (Component, error) {
		component := NewBaseComponent(config.ID, config.Name, config.Type)

		// Set description in metadata if provided
		if config.Description != "" {
			component.SetMetadata("description", config.Description)
		}

		// Set dependencies in metadata if provided
		if len(config.Dependencies) > 0 {
			component.SetMetadata("dependencies", config.Dependencies)
		}

		// Add all properties to metadata
		for key, value := range config.Properties {
			component.SetMetadata(key, value)
		}

		return component, nil
	})

	return factory
}

// CreateFactory is a factory method for backward compatibility.
// Creates a factory with default dependencies.
func CreateFactory() *DefaultFactory {
	return NewFactory(DefaultFactoryOptions{
		Logger: logging.CreateStandardLogger(logging.Info),
	})
}

// RegisterTypeCreator registers a creation function for a specific component type.
func (f *DefaultFactory) RegisterTypeCreator(
	componentType ComponentType,
	creator func(config ComponentConfig) (Component, error),
) {
	f.typeCreators[componentType] = creator
}

// Create creates a component from configuration.
func (f *DefaultFactory) Create(config ComponentConfig) (Component, error) {
	// Check if we have a creator for this component type
	creator, exists := f.typeCreators[config.Type]
	if !exists {
		return nil, NewError(ErrInvalidComponent, "no creator for component type", nil).
			WithDetail("type", config.Type)
	}

	return creator(config)
}
