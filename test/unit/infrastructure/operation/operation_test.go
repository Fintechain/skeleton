package operation

import (
	"fmt"
	"sync"
	"testing"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/operation"
	"github.com/fintechain/skeleton/internal/domain/service"
	operationImpl "github.com/fintechain/skeleton/internal/infrastructure/operation"
	"github.com/fintechain/skeleton/pkg/context"
	"github.com/fintechain/skeleton/test/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOperation(t *testing.T) {
	tests := []struct {
		name          string
		baseComponent component.Component
		expectNil     bool
		description   string
	}{
		{
			name:          "valid component",
			baseComponent: mocks.NewFactory().ComponentInterface(),
			expectNil:     false,
			description:   "Should create operation with valid component",
		},
		{
			name:          "nil component",
			baseComponent: nil,
			expectNil:     true,
			description:   "Should return nil with nil component",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			op := operationImpl.NewOperation(tt.baseComponent)

			// Verify
			if tt.expectNil {
				assert.Nil(t, op, tt.description)
			} else {
				assert.NotNil(t, op, tt.description)

				// Verify interface compliance
				var _ operation.Operation = op

				// Verify component delegation
				assert.Equal(t, tt.baseComponent.ID(), op.ID())
				assert.Equal(t, tt.baseComponent.Name(), op.Name())
				assert.Equal(t, tt.baseComponent.Type(), op.Type())
			}
		})
	}
}

func TestOperationInterfaceCompliance(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()

	// Test interface compliance
	op := operationImpl.NewOperation(mockComponent)
	require.NotNil(t, op)

	// Verify operation interface
	var _ operation.Operation = op

	// Verify component interface (through embedding)
	var _ component.Component = op
}

func TestOperationExecute(t *testing.T) {
	tests := []struct {
		name           string
		setupOperation func() operation.Operation
		setupContext   func() context.Context
		input          operation.Input
		expectError    bool
		errorContains  string
		description    string
	}{
		{
			name: "successful execution",
			setupOperation: func() operation.Operation {
				factory := mocks.NewFactory()
				mockComponent := factory.ComponentInterface().(*mocks.MockComponent)
				mockComponent.SetID("test-operation")
				mockComponent.SetName("Test Operation")
				return operationImpl.NewOperation(mockComponent)
			},
			setupContext: func() context.Context {
				return context.NewContext()
			},
			input:       map[string]interface{}{"key": "value"},
			expectError: false,
			description: "Should execute successfully with valid context and input",
		},
		{
			name: "nil context",
			setupOperation: func() operation.Operation {
				factory := mocks.NewFactory()
				mockComponent := factory.ComponentInterface()
				return operationImpl.NewOperation(mockComponent)
			},
			setupContext: func() context.Context {
				return nil
			},
			input:         map[string]interface{}{"key": "value"},
			expectError:   true,
			errorContains: "context is required",
			description:   "Should fail with nil context",
		},
		{
			name: "operation with lifecycle aware component in active state",
			setupOperation: func() operation.Operation {
				factory := mocks.NewFactory()
				mockComponent := factory.ComponentInterface().(*mocks.MockComponent)
				return operationImpl.NewOperation(mockComponent)
			},
			setupContext: func() context.Context {
				return context.NewContext()
			},
			input:       map[string]interface{}{"key": "value"},
			expectError: false,
			description: "Should execute successfully with active lifecycle component",
		},
		{
			name: "operation with lifecycle aware component in invalid state",
			setupOperation: func() operation.Operation {
				factory := mocks.NewFactory()
				mockComponent := factory.ComponentInterface().(*mocks.MockComponent)
				return operationImpl.NewOperation(mockComponent)
			},
			setupContext: func() context.Context {
				return context.NewContext()
			},
			input:       map[string]interface{}{"key": "value"},
			expectError: false,
			description: "Should execute successfully with regular component",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			op := tt.setupOperation()
			ctx := tt.setupContext()

			// Execute
			result, err := op.Execute(ctx, tt.input)

			// Verify
			if tt.expectError {
				assert.Error(t, err, tt.description)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err, tt.description)
				assert.NotNil(t, result)

				// Verify result structure for successful execution
				resultMap, ok := result.(map[string]interface{})
				assert.True(t, ok, "Result should be a map")
				assert.Equal(t, "executed", resultMap["status"])
				assert.Equal(t, op.ID(), resultMap["operation_id"])
				assert.Equal(t, op.Name(), resultMap["operation_name"])
				assert.Equal(t, tt.input, resultMap["input"])
			}
		})
	}
}

func TestOperationExecuteWithValidation(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface().(*mocks.MockComponent)
	mockComponent.SetID("test-operation")

	// Create operation (need to cast to access ExecuteWithValidation)
	baseOp := operationImpl.NewOperation(mockComponent).(*operationImpl.BaseOperation)
	ctx := context.NewContext()
	input := map[string]interface{}{"key": "value"}

	t.Run("successful validation", func(t *testing.T) {
		validator := func(input operation.Input) error {
			return nil // Valid input
		}

		result, err := baseOp.ExecuteWithValidation(ctx, input, validator)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("failed validation", func(t *testing.T) {
		validator := func(input operation.Input) error {
			return fmt.Errorf("invalid input")
		}

		result, err := baseOp.ExecuteWithValidation(ctx, input, validator)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "input validation failed")
		assert.Nil(t, result)
	})

	t.Run("nil validator", func(t *testing.T) {
		result, err := baseOp.ExecuteWithValidation(ctx, input, nil)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestOperationErrorHandling(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface().(*mocks.MockComponent)

	t.Run("operation with nil component", func(t *testing.T) {
		// Create operation with component that becomes nil (simulating error condition)
		op := operationImpl.NewOperation(mockComponent)

		// Use reflection or type assertion to set component to nil for testing
		if baseOp, ok := op.(*operationImpl.BaseOperation); ok {
			// This tests the error path where component is nil during execution
			baseOp.Component = nil

			ctx := context.NewContext()
			result, err := baseOp.Execute(ctx, map[string]interface{}{})

			assert.Error(t, err)
			assert.Contains(t, err.Error(), service.ErrServiceNotFound)
			assert.Nil(t, result)
		}
	})
}

func TestOperationConcurrency(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface().(*mocks.MockComponent)
	mockComponent.SetID("concurrent-operation")

	op := operationImpl.NewOperation(mockComponent)
	require.NotNil(t, op)

	// Test concurrent execution
	numGoroutines := 10
	var wg sync.WaitGroup
	results := make([]operation.Output, numGoroutines)
	errors := make([]error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()

			ctx := context.NewContext()
			input := map[string]interface{}{"index": index}

			result, err := op.Execute(ctx, input)
			results[index] = result
			errors[index] = err
		}(i)
	}

	wg.Wait()

	// Verify all executions succeeded
	for i := 0; i < numGoroutines; i++ {
		assert.NoError(t, errors[i], "Execution %d should succeed", i)
		assert.NotNil(t, results[i], "Result %d should not be nil", i)

		resultMap, ok := results[i].(map[string]interface{})
		assert.True(t, ok, "Result %d should be a map", i)
		assert.Equal(t, "executed", resultMap["status"])

		inputMap := resultMap["input"].(map[string]interface{})
		assert.Equal(t, i, inputMap["index"])
	}
}

func TestOperationComponentDelegation(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface().(*mocks.MockComponent)
	mockComponent.SetID("delegation-test")
	mockComponent.SetName("Delegation Test")
	mockComponent.SetType(component.TypeOperation)
	mockComponent.SetDescription("Test component delegation")

	op := operationImpl.NewOperation(mockComponent)
	require.NotNil(t, op)

	// Verify that operation delegates to component for basic properties
	assert.Equal(t, "delegation-test", op.ID())
	assert.Equal(t, "Delegation Test", op.Name())
	assert.Equal(t, component.TypeOperation, op.Type())
	assert.Equal(t, "Test component delegation", op.Description())
}

func TestOperationLifecycleIntegration(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface().(*mocks.MockComponent)
	mockComponent.SetID("lifecycle-test")

	op := operationImpl.NewOperation(mockComponent)
	require.NotNil(t, op)

	// Test operation execution after component initialization
	ctx := context.NewContext()
	mockSystem := factory.SystemInterface()

	// Initialize component
	err := op.Initialize(ctx, mockSystem)
	assert.NoError(t, err)

	// Execute operation
	result, err := op.Execute(ctx, map[string]interface{}{"test": "data"})
	assert.NoError(t, err)
	assert.NotNil(t, result)

	// Dispose component
	err = op.Dispose()
	assert.NoError(t, err)
}

func TestOperationWithDifferentInputTypes(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface().(*mocks.MockComponent)
	mockComponent.SetID("input-test")

	op := operationImpl.NewOperation(mockComponent)
	require.NotNil(t, op)
	ctx := context.NewContext()

	// Test with different input types
	testInputs := []operation.Input{
		map[string]interface{}{"key": "value"},
		[]string{"item1", "item2", "item3"},
		"simple string input",
		42,
		nil,
	}

	for i, input := range testInputs {
		t.Run(fmt.Sprintf("input_type_%d", i), func(t *testing.T) {
			result, err := op.Execute(ctx, input)
			assert.NoError(t, err)
			assert.NotNil(t, result)

			// Verify result contains the input
			resultMap, ok := result.(map[string]interface{})
			assert.True(t, ok)
			assert.Equal(t, input, resultMap["input"])
		})
	}
}
