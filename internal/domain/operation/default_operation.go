package operation

import (
	"github.com/ebanfa/skeleton/internal/domain/component"
)

// DefaultOperationOptions contains options for creating a DefaultOperation.
type DefaultOperationOptions struct {
	Component   component.Component
	ExecuteFunc func(ctx component.Context, input Input) (Output, error)
}

// DefaultOperation provides a standard implementation of the Operation interface.
type DefaultOperation struct {
	*BaseOperation
	executeFunc func(ctx component.Context, input Input) (Output, error)
}

// NewDefaultOperation creates a new default operation with dependency injection.
func NewDefaultOperation(options DefaultOperationOptions) *DefaultOperation {
	return &DefaultOperation{
		BaseOperation: NewBaseOperation(BaseOperationOptions{
			Component: options.Component,
		}),
		executeFunc: options.ExecuteFunc,
	}
}

// CreateDefaultOperation is a factory method for backward compatibility.
func CreateDefaultOperation(comp component.Component) *DefaultOperation {
	return NewDefaultOperation(DefaultOperationOptions{
		Component: comp,
	})
}

// WithExecuteFunc sets the execution function for this operation.
func (o *DefaultOperation) WithExecuteFunc(fn func(ctx component.Context, input Input) (Output, error)) *DefaultOperation {
	o.executeFunc = fn
	return o
}

// Execute runs the operation with the given context and input.
func (o *DefaultOperation) Execute(ctx component.Context, input Input) (Output, error) {
	// If we have an execute function, use it
	if o.executeFunc != nil {
		return o.executeFunc(ctx, input)
	}

	// Otherwise, use the base operation's execute function
	return o.BaseOperation.Execute(ctx, input)
}

// MapOperation creates an operation that maps input to output using a transform function.
func MapOperation(id string, transform func(input Input) (Output, error)) Operation {
	// Create a base component for this operation
	comp := component.CreateDefaultComponent(id, "Map: "+id, component.TypeOperation, "Maps input to output using a transform function")

	// Create the operation with the execute function
	return NewDefaultOperation(DefaultOperationOptions{
		Component: comp,
		ExecuteFunc: func(ctx component.Context, input Input) (Output, error) {
			// Apply the transform function to the input
			return transform(input)
		},
	})
}

// FilterOperation creates an operation that filters input based on a predicate function.
func FilterOperation(id string, predicate func(input Input) bool) Operation {
	// Create a base component for this operation
	comp := component.CreateDefaultComponent(id, "Filter: "+id, component.TypeOperation, "Filters input based on a predicate function")

	// Create the operation with the execute function
	return NewDefaultOperation(DefaultOperationOptions{
		Component: comp,
		ExecuteFunc: func(ctx component.Context, input Input) (Output, error) {
			// If the predicate returns true, pass the input through
			if predicate(input) {
				return input, nil
			}

			// Otherwise, return nil
			return nil, nil
		},
	})
}

// AsyncOperation wraps an operation to execute asynchronously.
func AsyncOperation(id string, wrapped Operation) Operation {
	// Create a base component for this operation
	comp := component.CreateDefaultComponent(id, "Async: "+id, component.TypeOperation, "Executes an operation asynchronously")

	// Create the operation with the execute function
	return NewDefaultOperation(DefaultOperationOptions{
		Component: comp,
		ExecuteFunc: func(ctx component.Context, input Input) (Output, error) {
			// Create a channel for the result
			resultCh := make(chan struct {
				output Output
				err    error
			})

			// Execute the wrapped operation in a goroutine
			go func() {
				output, err := wrapped.Execute(ctx, input)
				resultCh <- struct {
					output Output
					err    error
				}{output, err}
			}()

			// Wait for either the result or context cancellation
			select {
			case result := <-resultCh:
				return result.output, result.err
			case <-ctx.Done():
				return nil, component.NewError(
					ErrOperationTimeout,
					"operation was canceled",
					ctx.Err(),
				).WithDetail("operation_id", id)
			}
		},
	})
}
