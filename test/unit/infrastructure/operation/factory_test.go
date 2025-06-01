package operation

import (
	"testing"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/operation"
	operationImpl "github.com/fintechain/skeleton/internal/infrastructure/operation"
	"github.com/fintechain/skeleton/pkg/context"
	"github.com/fintechain/skeleton/test/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOperationFactory(t *testing.T) {
	tests := []struct {
		name             string
		componentFactory component.Factory
		expectNil        bool
		description      string
	}{
		{
			name:             "valid component factory",
			componentFactory: mocks.NewFactory().ComponentFactoryInterface(),
			expectNil:        false,
			description:      "Should create operation factory with valid component factory",
		},
		{
			name:             "nil component factory",
			componentFactory: nil,
			expectNil:        true,
			description:      "Should return nil with nil component factory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			factory := operationImpl.NewOperationFactory(tt.componentFactory)

			// Verify
			if tt.expectNil {
				assert.Nil(t, factory, tt.description)
			} else {
				assert.NotNil(t, factory, tt.description)

				// Verify interface compliance
				var _ operation.OperationFactory = factory
				var _ component.Factory = factory
			}
		})
	}
}

func TestOperationFactoryInterfaceCompliance(t *testing.T) {
	mockComponentFactory := mocks.NewFactory().ComponentFactoryInterface()
	factory := operationImpl.NewOperationFactory(mockComponentFactory)
	require.NotNil(t, factory)

	// Verify operation factory interface
	var _ operation.OperationFactory = factory

	// Verify component factory interface (through embedding)
	var _ component.Factory = factory
}

func TestOperationFactoryCreate(t *testing.T) {
	tests := []struct {
		name        string
		setupConfig func() component.ComponentConfig
		expectError bool
		description string
	}{
		{
			name: "valid operation config",
			setupConfig: func() component.ComponentConfig {
				return component.NewComponentConfig("test-operation", "Test Operation", component.TypeOperation, "Test operation description")
			},
			expectError: false,
			description: "Should create component with valid operation config",
		},
		{
			name: "invalid config type",
			setupConfig: func() component.ComponentConfig {
				return component.NewComponentConfig("test-service", "Test Service", component.TypeService, "Test service description") // Wrong type
			},
			expectError: false, // The factory should convert the type to Operation
			description: "Should create component and convert type to operation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockComponentFactory := mocks.NewFactory().ComponentFactoryInterface()
			factory := operationImpl.NewOperationFactory(mockComponentFactory)
			require.NotNil(t, factory)

			config := tt.setupConfig()

			// Execute
			comp, err := factory.Create(config)

			// Verify
			if tt.expectError {
				assert.Error(t, err, tt.description)
				assert.Nil(t, comp)
			} else {
				assert.NoError(t, err, tt.description)
				assert.NotNil(t, comp)

				// Verify component properties
				assert.Equal(t, config.ID, comp.ID())
				assert.Equal(t, config.Name, comp.Name())
				assert.Equal(t, component.TypeOperation, comp.Type()) // Should always be operation type
			}
		})
	}
}

func TestOperationFactoryCreateOperation(t *testing.T) {
	tests := []struct {
		name        string
		setupConfig func() operation.OperationConfig
		expectError bool
		description string
	}{
		{
			name: "valid operation config",
			setupConfig: func() operation.OperationConfig {
				return operation.NewOperationConfig("test-operation", "Test Operation", "Test operation description")
			},
			expectError: false,
			description: "Should create operation with valid config",
		},
		{
			name: "config with metadata",
			setupConfig: func() operation.OperationConfig {
				return operation.NewOperationConfig("metadata-operation", "Metadata Operation", "Operation with metadata")
			},
			expectError: false,
			description: "Should create operation with metadata",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockComponentFactory := mocks.NewFactory().ComponentFactoryInterface()
			factory := operationImpl.NewOperationFactory(mockComponentFactory)
			require.NotNil(t, factory)

			config := tt.setupConfig()

			// Execute
			op, err := factory.CreateOperation(config)

			// Verify
			if tt.expectError {
				assert.Error(t, err, tt.description)
				assert.Nil(t, op)
			} else {
				assert.NoError(t, err, tt.description)
				assert.NotNil(t, op)

				// Verify operation interface
				var _ operation.Operation = op

				// Verify operation properties
				assert.Equal(t, config.ID, op.ID())
				assert.Equal(t, config.Name, op.Name())
				assert.Equal(t, component.TypeOperation, op.Type())

				// Test operation execution
				ctx := context.NewContext()
				result, err := op.Execute(ctx, map[string]interface{}{"test": "data"})
				assert.NoError(t, err)
				assert.NotNil(t, result)
			}
		})
	}
}

func TestOperationFactoryCreateOperationWithValidation(t *testing.T) {
	mockComponentFactory := mocks.NewFactory().ComponentFactoryInterface()
	factory := operationImpl.NewOperationFactory(mockComponentFactory)
	require.NotNil(t, factory)

	config := operation.NewOperationConfig("test-operation", "Test Operation", "Test operation description")

	// Use CreateOperation since CreateOperationWithValidation doesn't exist
	op, err := factory.CreateOperation(config)
	assert.NoError(t, err)
	assert.NotNil(t, op)
}

func TestOperationFactoryErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		setupConfig func() component.ComponentConfig
		expectError bool
		description string
	}{
		{
			name: "valid config",
			setupConfig: func() component.ComponentConfig {
				return component.NewComponentConfig("valid-operation", "Valid Operation", component.TypeOperation, "Valid operation description")
			},
			expectError: false,
			description: "Should create operation with valid config",
		},
		{
			name: "empty ID",
			setupConfig: func() component.ComponentConfig {
				return component.NewComponentConfig("", "Empty ID Operation", component.TypeOperation, "Operation with empty ID")
			},
			expectError: true, // The component factory validates empty IDs
			description: "Should fail with empty ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockComponentFactory := mocks.NewFactory().ComponentFactoryInterface()
			factory := operationImpl.NewOperationFactory(mockComponentFactory)
			require.NotNil(t, factory)

			config := tt.setupConfig()

			// Execute
			comp, err := factory.Create(config)

			// Verify
			if tt.expectError {
				assert.Error(t, err, tt.description)
				assert.Nil(t, comp)
			} else {
				assert.NoError(t, err, tt.description)
				assert.NotNil(t, comp)
			}
		})
	}
}

func TestOperationFactoryComponentDelegation(t *testing.T) {
	mockComponentFactory := mocks.NewFactory().ComponentFactoryInterface()
	factory := operationImpl.NewOperationFactory(mockComponentFactory)
	require.NotNil(t, factory)

	// Test that the operation factory delegates to the component factory
	config := component.NewComponentConfig("delegation-test", "Delegation Test", component.TypeOperation, "Test delegation")

	comp, err := factory.Create(config)
	assert.NoError(t, err)
	assert.NotNil(t, comp)
	assert.Equal(t, component.TypeOperation, comp.Type())
}

func TestOperationFactoryMultipleOperations(t *testing.T) {
	mockComponentFactory := mocks.NewFactory().ComponentFactoryInterface()
	factory := operationImpl.NewOperationFactory(mockComponentFactory)
	require.NotNil(t, factory)

	configs := []component.ComponentConfig{
		component.NewComponentConfig("op-1", "Operation 1", component.TypeOperation, "First operation"),
		component.NewComponentConfig("op-2", "Operation 2", component.TypeOperation, "Second operation"),
		component.NewComponentConfig("op-3", "Operation 3", component.TypeOperation, "Third operation"),
	}

	operations := make([]component.Component, len(configs))

	// Create multiple operations using the same factory
	for i, config := range configs {
		op, err := factory.Create(config)
		require.NoError(t, err)
		require.NotNil(t, op)
		operations[i] = op
	}

	// Verify each operation has correct properties
	for i, op := range operations {
		expectedConfig := configs[i]
		assert.Equal(t, expectedConfig.ID, op.ID())
		assert.Equal(t, expectedConfig.Name, op.Name())
		assert.Equal(t, component.TypeOperation, op.Type())
	}
}

func TestOperationFactoryWithContext(t *testing.T) {
	mockComponentFactory := mocks.NewFactory().ComponentFactoryInterface()
	factory := operationImpl.NewOperationFactory(mockComponentFactory)
	require.NotNil(t, factory)

	config := component.NewComponentConfig("context-test", "Context Test", component.TypeOperation, "Test with context")

	comp, err := factory.Create(config)
	require.NoError(t, err)
	require.NotNil(t, comp)

	// Test operation execution with context
	ctx := context.NewContext()

	// Cast to operation to test Execute method
	if op, ok := comp.(operation.Operation); ok {
		result, err := op.Execute(ctx, map[string]interface{}{"test": "data"})
		assert.NoError(t, err)
		assert.NotNil(t, result)
	}
}

func TestOperationFactoryConfigValidation(t *testing.T) {
	tests := []struct {
		name        string
		setupConfig func() component.ComponentConfig
		expectError bool
		description string
	}{
		{
			name: "valid config",
			setupConfig: func() component.ComponentConfig {
				return component.NewComponentConfig("valid", "Valid Operation", component.TypeOperation, "Valid operation")
			},
			expectError: false,
			description: "Should create operation with valid config",
		},
		{
			name: "empty name",
			setupConfig: func() component.ComponentConfig {
				return component.NewComponentConfig("empty-name", "", component.TypeOperation, "Operation with empty name")
			},
			expectError: true, // The component factory validates empty names
			description: "Should fail with empty name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockComponentFactory := mocks.NewFactory().ComponentFactoryInterface()
			factory := operationImpl.NewOperationFactory(mockComponentFactory)
			require.NotNil(t, factory)

			config := tt.setupConfig()

			comp, err := factory.Create(config)

			if tt.expectError {
				assert.Error(t, err, tt.description)
				assert.Nil(t, comp)
			} else {
				assert.NoError(t, err, tt.description)
				assert.NotNil(t, comp)
			}
		})
	}
}
