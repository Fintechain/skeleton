// Package operation provides functionality for operations in the system.
package operation

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
)

// Input represents the input for an operation.
type Input interface{}

// Output represents the output from an operation.
type Output interface{}

// Operation is a specialized component that executes discrete units of work.
type Operation interface {
	component.Component

	// Execute runs the operation with the given context and input.
	Execute(ctx context.Context, input Input) (Output, error)
}

// OperationConfig defines the configuration for creating an operation.
type OperationConfig struct {
	component.ComponentConfig
	// Operation-specific configuration properties can be added here
}

// NewOperationConfig creates a new OperationConfig with the given parameters.
// This ensures the component type is set to TypeOperation.
func NewOperationConfig(id, name, description string) OperationConfig {
	return OperationConfig{
		ComponentConfig: component.NewComponentConfig(id, name, component.TypeOperation, description),
	}
}

// OperationFactory creates operations from configuration.
type OperationFactory interface {
	component.Factory

	// CreateOperation creates an operation from the given configuration.
	CreateOperation(config OperationConfig) (Operation, error)
}
