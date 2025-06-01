// Package operation provides concrete implementations for operation factory functionality.
package operation

import (
	"fmt"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/operation"
)

// OperationFactory provides a concrete implementation of the OperationFactory interface.
type OperationFactory struct {
	componentFactory component.Factory
}

// NewOperationFactory creates a new operation factory instance.
// This constructor accepts a component factory interface dependency for component creation.
func NewOperationFactory(componentFactory component.Factory) operation.OperationFactory {
	if componentFactory == nil {
		return nil
	}

	return &OperationFactory{
		componentFactory: componentFactory,
	}
}

// Create creates a component from the given configuration.
// This method implements the Factory interface by delegating to the component factory.
func (f *OperationFactory) Create(config component.ComponentConfig) (component.Component, error) {
	if f.componentFactory == nil {
		return nil, fmt.Errorf("component factory is required for operation creation")
	}

	// Ensure the component type is set to Operation if not already set
	if config.Type != component.TypeOperation {
		config.Type = component.TypeOperation
	}

	// Delegate to the component factory for basic component creation
	baseComponent, err := f.componentFactory.Create(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create base component for operation: %w", err)
	}

	// Wrap the base component in an operation implementation
	op := NewOperation(baseComponent)
	if op == nil {
		return nil, fmt.Errorf("failed to create operation from component")
	}

	return op, nil
}

// CreateOperation creates an operation from the given operation configuration.
// This method provides operation-specific creation functionality.
func (f *OperationFactory) CreateOperation(config operation.OperationConfig) (operation.Operation, error) {
	if f.componentFactory == nil {
		return nil, fmt.Errorf("component factory is required for operation creation")
	}

	// Validate operation configuration
	if config.ID == "" {
		return nil, fmt.Errorf("operation ID cannot be empty")
	}
	if config.Name == "" {
		return nil, fmt.Errorf("operation name cannot be empty")
	}

	// Ensure the component type is set to Operation
	config.ComponentConfig.Type = component.TypeOperation

	// Create the base component using the component factory
	baseComponent, err := f.componentFactory.Create(config.ComponentConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create base component for operation '%s': %w", config.ID, err)
	}

	// Create the operation by wrapping the base component
	op := NewOperation(baseComponent)
	if op == nil {
		return nil, fmt.Errorf("failed to create operation '%s' from component", config.ID)
	}

	return op, nil
}

// CreateOperationWithValidation creates an operation with additional validation.
// This method provides extended creation functionality with custom validation.
func (f *OperationFactory) CreateOperationWithValidation(config operation.OperationConfig, validator func(operation.OperationConfig) error) (operation.Operation, error) {
	if validator != nil {
		if err := validator(config); err != nil {
			return nil, fmt.Errorf("operation configuration validation failed: %w", err)
		}
	}

	return f.CreateOperation(config)
}
