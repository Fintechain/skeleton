package component

import (
	"testing"

	"github.com/fintechain/skeleton/internal/domain/component"
	infraComponent "github.com/fintechain/skeleton/internal/infrastructure/component"
	infraContext "github.com/fintechain/skeleton/internal/infrastructure/context"
	"github.com/fintechain/skeleton/test/unit/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewBaseOperation(t *testing.T) {
	config := component.ComponentConfig{
		ID:          "test-operation",
		Name:        "Test Operation",
		Type:        component.TypeOperation,
		Description: "Test operation description",
		Version:     "1.0.0",
	}

	operation := infraComponent.NewBaseOperation(config)
	assert.NotNil(t, operation)

	// Verify interface compliance
	var _ component.Operation = operation
	var _ component.Component = operation

	// Test basic properties
	assert.Equal(t, component.ComponentID("test-operation"), operation.ID())
	assert.Equal(t, "Test Operation", operation.Name())
	assert.Equal(t, component.TypeOperation, operation.Type())
	assert.Equal(t, "1.0.0", operation.Version())
}

func TestBaseOperationExecute(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "test-operation",
		Name: "Test Operation",
		Type: component.TypeOperation,
	}

	operation := infraComponent.NewBaseOperation(config)
	ctx := infraContext.NewContext()

	// Test execute with input
	input := component.Input{
		Data:     "test-data",
		Metadata: map[string]string{"key": "value"},
	}

	output, err := operation.Execute(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, output)

	// Base implementation should return the input data as output
	assert.Equal(t, input.Data, output.Data)
}

func TestBaseOperationExecuteWithNilData(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "test-operation",
		Name: "Test Operation",
		Type: component.TypeOperation,
	}

	operation := infraComponent.NewBaseOperation(config)
	ctx := infraContext.NewContext()

	// Test execute with nil data
	input := component.Input{
		Data:     nil,
		Metadata: map[string]string{"key": "value"},
	}

	output, err := operation.Execute(ctx, input)
	assert.NoError(t, err)
	assert.NotNil(t, output)
	assert.Nil(t, output.Data)
}

func TestBaseOperationExecuteWithDifferentDataTypes(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "test-operation",
		Name: "Test Operation",
		Type: component.TypeOperation,
	}

	operation := infraComponent.NewBaseOperation(config)
	ctx := infraContext.NewContext()

	tests := []struct {
		name string
		data interface{}
	}{
		{"string data", "test string"},
		{"integer data", 42},
		{"boolean data", true},
		{"map data", map[string]interface{}{"key": "value"}},
		{"slice data", []string{"item1", "item2"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := component.Input{
				Data:     tt.data,
				Metadata: map[string]string{"type": tt.name},
			}

			output, err := operation.Execute(ctx, input)
			assert.NoError(t, err)
			assert.NotNil(t, output)
			assert.Equal(t, tt.data, output.Data)
		})
	}
}

func TestBaseOperationInitializeAndDispose(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "test-operation",
		Name: "Test Operation",
		Type: component.TypeOperation,
	}

	operation := infraComponent.NewBaseOperation(config)
	factory := mocks.NewFactory()
	mockSystem := factory.SystemInterface()
	ctx := infraContext.NewContext()

	// Test initialization
	err := operation.Initialize(ctx, mockSystem)
	assert.NoError(t, err)

	// Test disposal
	err = operation.Dispose()
	assert.NoError(t, err)
}
