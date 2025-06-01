// Package mocks provides centralized mock implementations for the Skeleton Framework.
// This file contains mocks for operation interfaces.
package mocks

import (
	"sync"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/operation"
)

// MockOperation implements the operation.Operation interface for testing.
type MockOperation struct {
	MockComponent
	mu sync.RWMutex

	// Operation-specific fields
	executeFunc   func(ctx context.Context, input operation.Input) (operation.Output, error)
	executeResult operation.Output
	executeError  error
	executeCalled bool
	executeInput  operation.Input
}

// NewMockOperation creates a new mock operation with default behavior.
func NewMockOperation() *MockOperation {
	mock := &MockOperation{
		MockComponent: *NewMockComponent(),
	}
	mock.SetType(component.TypeOperation)
	return mock
}

// Execute implements the operation.Operation interface.
func (m *MockOperation) Execute(ctx context.Context, input operation.Input) (operation.Output, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.executeCalled = true
	m.executeInput = input

	if m.executeFunc != nil {
		return m.executeFunc(ctx, input)
	}

	return m.executeResult, m.executeError
}

// Mock configuration methods

// SetExecuteFunc sets a custom function to be called when Execute is invoked.
func (m *MockOperation) SetExecuteFunc(fn func(ctx context.Context, input operation.Input) (operation.Output, error)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.executeFunc = fn
}

// SetExecuteResult sets the result to return from Execute calls.
func (m *MockOperation) SetExecuteResult(result operation.Output) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.executeResult = result
}

// SetExecuteError sets the error to return from Execute calls.
func (m *MockOperation) SetExecuteError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.executeError = err
}

// Verification methods

// WasExecuteCalled returns true if Execute was called.
func (m *MockOperation) WasExecuteCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.executeCalled
}

// GetExecuteInput returns the input passed to the last Execute call.
func (m *MockOperation) GetExecuteInput() operation.Input {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.executeInput
}

// MockOperationFactory implements the operation.OperationFactory interface for testing.
type MockOperationFactory struct {
	MockComponentFactory
	mu sync.RWMutex

	// Factory-specific fields
	createOperationFunc   func(config operation.OperationConfig) (operation.Operation, error)
	createOperationResult operation.Operation
	createOperationError  error
	createOperationCalled bool
	createOperationConfig operation.OperationConfig
}

// NewMockOperationFactory creates a new mock operation factory.
func NewMockOperationFactory() *MockOperationFactory {
	return &MockOperationFactory{
		MockComponentFactory: *NewMockComponentFactory(),
	}
}

// CreateOperation implements the operation.OperationFactory interface.
func (m *MockOperationFactory) CreateOperation(config operation.OperationConfig) (operation.Operation, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.createOperationCalled = true
	m.createOperationConfig = config

	if m.createOperationFunc != nil {
		return m.createOperationFunc(config)
	}

	if m.createOperationError != nil {
		return nil, m.createOperationError
	}

	if m.createOperationResult != nil {
		return m.createOperationResult, nil
	}

	// Default behavior: create a mock operation
	mockOp := NewMockOperation()
	mockOp.SetID(config.ID)
	mockOp.SetName(config.Name)
	mockOp.SetDescription(config.Description)
	return mockOp, nil
}

// Mock configuration methods

// SetCreateOperationFunc sets a custom function for CreateOperation.
func (m *MockOperationFactory) SetCreateOperationFunc(fn func(config operation.OperationConfig) (operation.Operation, error)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createOperationFunc = fn
}

// SetCreateOperationResult sets the result to return from CreateOperation.
func (m *MockOperationFactory) SetCreateOperationResult(result operation.Operation) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createOperationResult = result
}

// SetCreateOperationError sets the error to return from CreateOperation.
func (m *MockOperationFactory) SetCreateOperationError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createOperationError = err
}

// Verification methods

// WasCreateOperationCalled returns true if CreateOperation was called.
func (m *MockOperationFactory) WasCreateOperationCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.createOperationCalled
}

// GetCreateOperationConfig returns the config passed to the last CreateOperation call.
func (m *MockOperationFactory) GetCreateOperationConfig() operation.OperationConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.createOperationConfig
}
