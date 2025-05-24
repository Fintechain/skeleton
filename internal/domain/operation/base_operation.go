package operation

import (
	"github.com/ebanfa/skeleton/internal/domain/component"
)

// BaseOperationOptions contains options for creating a BaseOperation.
type BaseOperationOptions struct {
	Component component.Component
}

// BaseOperation provides a basic implementation of the Operation interface.
type BaseOperation struct {
	component.Component
}

// NewBaseOperation creates a new base operation with the given component using dependency injection.
func NewBaseOperation(options BaseOperationOptions) *BaseOperation {
	return &BaseOperation{
		Component: options.Component,
	}
}

// CreateBaseOperation is a factory method for backward compatibility.
func CreateBaseOperation(comp component.Component) *BaseOperation {
	return NewBaseOperation(BaseOperationOptions{
		Component: comp,
	})
}

// Execute runs the operation with the given context and input.
// This base implementation just returns an error, subclasses should override it.
func (o *BaseOperation) Execute(ctx component.Context, input Input) (Output, error) {
	return nil, component.NewError(
		ErrOperationExecution,
		"base operation does not implement Execute",
		nil,
	).WithDetail("operation_id", o.ID())
}
