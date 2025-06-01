// Package operation provides functionality for operations in the system.
package operation

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/operation"
	operationImpl "github.com/fintechain/skeleton/internal/infrastructure/operation"
)

// Re-export operation interfaces
type Operation = operation.Operation
type OperationFactory = operation.OperationFactory

// Re-export operation types
type Input = operation.Input
type Output = operation.Output
type OperationConfig = operation.OperationConfig

// Re-export constructor function
var NewOperationConfig = operation.NewOperationConfig

// Factory functions

// NewOperation creates a new operation instance from a base component.
func NewOperation(baseComponent component.Component) Operation {
	return operationImpl.NewOperation(baseComponent)
}

// NewOperationFactory creates a new operation factory instance.
func NewOperationFactory(componentFactory component.Factory) OperationFactory {
	return operationImpl.NewOperationFactory(componentFactory)
}
