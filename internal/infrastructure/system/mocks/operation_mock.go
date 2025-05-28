package mocks

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/operation"
)

// MockOperation is a mock implementation of the Operation interface
type MockOperation struct {
	*MockComponent
	ExecuteFunc func(ctx component.Context, input operation.Input) (operation.Output, error)
}

// NewMockOperation creates a new MockOperation instance
func NewMockOperation(id, name string) *MockOperation {
	return &MockOperation{
		MockComponent: NewMockComponent(id, name, component.TypeOperation),
		ExecuteFunc:   func(ctx component.Context, input operation.Input) (operation.Output, error) { return nil, nil },
	}
}

// Execute implements the Operation interface
func (m *MockOperation) Execute(ctx component.Context, input operation.Input) (operation.Output, error) {
	if m.ExecuteFunc != nil {
		return m.ExecuteFunc(ctx, input)
	}
	return nil, nil
}
