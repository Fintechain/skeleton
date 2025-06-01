// Package operation provides functionality for operations in the system.
package operation

import (
	"github.com/fintechain/skeleton/internal/domain/component"
)

// Input represents the input for an operation.
type Input interface{}

// Output represents the output from an operation.
type Output interface{}

// Operation is a specialized component that executes discrete units of work.
type Operation interface {
	component.Component

	// Execute runs the operation with the given context and input.
	Execute(ctx component.Context, input Input) (Output, error)
}

// OperationConfig defines the configuration for creating an operation.
type OperationConfig struct {
	component.ComponentConfig
	// Operation-specific configuration properties
}

// OperationFactory creates operations from configuration.
type OperationFactory interface {
	component.Factory

	// CreateOperation creates an operation from the given configuration.
	CreateOperation(config OperationConfig) (Operation, error)
}

// DefaultOperationFactory provides a standard implementation of OperationFactory.
type DefaultOperationFactory struct {
	factory component.Factory
}

// NewOperationFactory creates a new operation factory.
func NewOperationFactory() OperationFactory {
	return &DefaultOperationFactory{
		factory: component.CreateFactory(),
	}
}

// Create creates a component from configuration (implements component.Factory).
func (f *DefaultOperationFactory) Create(config component.ComponentConfig) (component.Component, error) {
	return f.factory.Create(config)
}

// CreateOperation creates an operation from the given configuration.
func (f *DefaultOperationFactory) CreateOperation(config OperationConfig) (Operation, error) {
	// Create a component from the configuration
	comp, err := f.factory.Create(config.ComponentConfig)
	if err != nil {
		return nil, component.NewError(
			"operation.factory.create_failed",
			"failed to create component for operation",
			err,
		).WithDetail("config", config)
	}

	// Create a default operation with the component
	operation := CreateDefaultOperation(comp)
	return operation, nil
}

// NewOperation creates a new operation with the given options.
func NewOperation(options DefaultOperationOptions) Operation {
	return NewDefaultOperation(options)
}

// CreateOperation creates an operation from a component (convenience function).
func CreateOperation(comp component.Component) Operation {
	return CreateDefaultOperation(comp)
}

// Common error codes for operation execution
const (
	ErrOperationExecution = "operation.execution_failed"
	ErrInvalidInput       = "operation.invalid_input"
	ErrOperationTimeout   = "operation.timeout"
)
