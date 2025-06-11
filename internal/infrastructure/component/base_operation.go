package component

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
)

// BaseOperation provides common operation functionality that can be embedded
// in concrete operation implementations.
type BaseOperation struct {
	*BaseComponent
}

// NewBaseOperation creates a new base operation with the provided configuration.
func NewBaseOperation(config component.ComponentConfig) *BaseOperation {
	config.Type = component.TypeOperation
	return &BaseOperation{
		BaseComponent: NewBaseComponent(config),
	}
}

// Type returns the operation component type.
func (o *BaseOperation) Type() component.ComponentType {
	return component.TypeOperation
}

// Initialize prepares the operation for use within the system.
func (o *BaseOperation) Initialize(ctx context.Context, system component.System) error {
	return o.BaseComponent.Initialize(ctx, system)
}

// Execute runs the operation with the given context and input.
// Base implementation returns input as output - override in concrete implementations.
func (o *BaseOperation) Execute(ctx context.Context, input component.Input) (component.Output, error) {
	return component.Output{Data: input.Data}, nil
}
