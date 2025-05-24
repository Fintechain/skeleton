// Package operation provides functionality for operations in the system.
package operation

import (
	"github.com/ebanfa/skeleton/internal/domain/component"
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

// Common error codes for operation execution
const (
	ErrOperationExecution = "operation.execution_failed"
	ErrInvalidInput       = "operation.invalid_input"
	ErrOperationTimeout   = "operation.timeout"
)
