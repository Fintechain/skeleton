// Package operation provides concrete implementations for operation functionality.
package operation

import (
	"fmt"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/operation"
	"github.com/fintechain/skeleton/internal/domain/service"
)

// BaseOperation provides a concrete implementation of the Operation interface.
type BaseOperation struct {
	component.Component
}

// NewOperation creates a new operation instance.
// This constructor accepts a component interface dependency for composition.
func NewOperation(baseComponent component.Component) operation.Operation {
	if baseComponent == nil {
		return nil
	}

	return &BaseOperation{
		Component: baseComponent,
	}
}

// Execute runs the operation with the given context and input.
// This method provides the core operation execution functionality.
func (o *BaseOperation) Execute(ctx context.Context, input operation.Input) (operation.Output, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is required for operation execution")
	}

	// Check if the operation component is properly initialized
	if o.Component == nil {
		return nil, fmt.Errorf(service.ErrServiceNotFound + ": operation component not available")
	}

	// Validate that the operation is in a valid state for execution
	// If the component supports lifecycle awareness, check its state
	if lifecycleAware, ok := o.Component.(component.LifecycleAwareComponent); ok {
		state := lifecycleAware.State()
		if state != component.StateActive && state != component.StateInitialized {
			return nil, fmt.Errorf(service.ErrServiceStart+": operation not in executable state: %s", state)
		}
	}

	// For the base implementation, we provide a simple pass-through execution
	// Specific operation implementations would override this method with actual business logic
	result := map[string]interface{}{
		"operation_id":   o.Component.ID(),
		"operation_name": o.Component.Name(),
		"input":          input,
		"status":         "executed",
	}

	return result, nil
}

// ExecuteWithValidation provides a wrapper around Execute with additional validation.
// This method can be used by specific operation implementations for common validation patterns.
func (o *BaseOperation) ExecuteWithValidation(ctx context.Context, input operation.Input, validator func(operation.Input) error) (operation.Output, error) {
	if validator != nil {
		if err := validator(input); err != nil {
			return nil, fmt.Errorf("input validation failed: %w", err)
		}
	}

	return o.Execute(ctx, input)
}
