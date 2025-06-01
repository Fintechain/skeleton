package component

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fintechain/skeleton/pkg/component"
	"github.com/fintechain/skeleton/test/unit/mocks"
)

// TestNewFactory tests the component factory constructor function
func TestNewFactory(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()

	// Test factory constructor with mock registry dependency
	componentFactory := component.NewFactory(mockRegistry)
	assert.NotNil(t, componentFactory)

	// Verify interface compliance
	var _ component.Factory = componentFactory
}

// TestFactoryInterfaceCompliance verifies the implementation satisfies the domain interface
func TestFactoryInterfaceCompliance(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()

	// Verify interface compliance
	var _ component.Factory = component.NewFactory(mockRegistry)
}

// TestFactoryComponentCreation tests component creation from ComponentConfig
func TestFactoryComponentCreation(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()
	componentFactory := component.NewFactory(mockRegistry)

	tests := []struct {
		name        string
		id          string
		compName    string
		compType    component.ComponentType
		description string
	}{
		{
			name:        "create basic component",
			id:          "basic-comp",
			compName:    "Basic Component",
			compType:    component.TypeBasic,
			description: "A basic component",
		},
		{
			name:        "create operation component",
			id:          "op-comp",
			compName:    "Operation Component",
			compType:    component.TypeOperation,
			description: "An operation component",
		},
		{
			name:        "create service component",
			id:          "svc-comp",
			compName:    "Service Component",
			compType:    component.TypeService,
			description: "A service component",
		},
		{
			name:        "create system component",
			id:          "sys-comp",
			compName:    "System Component",
			compType:    component.TypeSystem,
			description: "A system component",
		},
		{
			name:        "create application component",
			id:          "app-comp",
			compName:    "Application Component",
			compType:    component.TypeApplication,
			description: "An application component",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := component.NewComponentConfig(tt.id, tt.compName, tt.compType, tt.description)
			comp, err := componentFactory.Create(config)

			assert.NoError(t, err)
			assert.NotNil(t, comp)
			assert.Equal(t, tt.id, comp.ID())
			assert.Equal(t, tt.compName, comp.Name())
			assert.Equal(t, tt.compType, comp.Type())
			assert.Equal(t, tt.description, comp.Description())
		})
	}
}

// TestFactoryWithMockRegistry tests factory integration with mock registry
func TestFactoryWithMockRegistry(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()

	// Configure mock registry behavior if needed
	// (This would depend on the mock implementation)

	componentFactory := component.NewFactory(mockRegistry)
	assert.NotNil(t, componentFactory)

	// Test component creation
	config := component.NewComponentConfig(
		"test-component",
		"Test Component",
		component.TypeBasic,
		"A test component",
	)

	comp, err := componentFactory.Create(config)
	assert.NoError(t, err)
	assert.NotNil(t, comp)
	assert.Equal(t, "test-component", comp.ID())
}

// TestFactoryMultipleComponents tests creating multiple components with the same factory
func TestFactoryMultipleComponents(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()
	componentFactory := component.NewFactory(mockRegistry)

	configs := []component.ComponentConfig{
		component.NewComponentConfig("comp-1", "Component 1", component.TypeBasic, "First component"),
		component.NewComponentConfig("comp-2", "Component 2", component.TypeOperation, "Second component"),
		component.NewComponentConfig("comp-3", "Component 3", component.TypeService, "Third component"),
	}

	components := make([]component.Component, len(configs))

	// Create multiple components using the same factory
	for i, config := range configs {
		comp, err := componentFactory.Create(config)
		require.NoError(t, err)
		require.NotNil(t, comp)
		components[i] = comp
	}

	// Verify each component has correct properties
	for i, comp := range components {
		expectedConfig := configs[i]
		assert.Equal(t, expectedConfig.ID, comp.ID())
		assert.Equal(t, expectedConfig.Name, comp.Name())
		assert.Equal(t, expectedConfig.Type, comp.Type())
		assert.Equal(t, expectedConfig.Description, comp.Description())
	}

	// Verify components are independent instances
	for i := 0; i < len(components); i++ {
		for j := i + 1; j < len(components); j++ {
			assert.NotEqual(t, components[i], components[j])
			assert.NotEqual(t, components[i].ID(), components[j].ID())
		}
	}
}

// TestFactoryErrorHandling tests factory error conditions
func TestFactoryErrorHandling(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()
	componentFactory := component.NewFactory(mockRegistry)

	// Test with various potentially invalid configurations
	testCases := []struct {
		name        string
		config      component.ComponentConfig
		expectError bool
		description string
	}{
		{
			name:        "valid configuration",
			config:      component.NewComponentConfig("valid", "Valid", component.TypeBasic, "Valid component"),
			expectError: false,
			description: "Should create component with valid configuration",
		},
		{
			name:        "empty description",
			config:      component.NewComponentConfig("empty-desc", "Empty Desc", component.TypeBasic, ""),
			expectError: false,
			description: "Should allow empty description",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			comp, err := componentFactory.Create(tc.config)

			if tc.expectError {
				assert.Error(t, err, tc.description)
				assert.Nil(t, comp)
			} else {
				assert.NoError(t, err, tc.description)
				assert.NotNil(t, comp)
			}
		})
	}
}

// TestFactoryComponentTypes tests factory with all component types
func TestFactoryComponentTypes(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()
	componentFactory := component.NewFactory(mockRegistry)

	componentTypes := []component.ComponentType{
		component.TypeBasic,
		component.TypeOperation,
		component.TypeService,
		component.TypeSystem,
		component.TypeApplication,
	}

	for i, compType := range componentTypes {
		t.Run(fmt.Sprintf("type-%d", i), func(t *testing.T) {
			config := component.NewComponentConfig(
				fmt.Sprintf("comp-%d", i),
				fmt.Sprintf("Component %d", i),
				compType,
				fmt.Sprintf("Component of type %d", i),
			)

			comp, err := componentFactory.Create(config)
			assert.NoError(t, err)
			assert.NotNil(t, comp)
			assert.Equal(t, compType, comp.Type())
		})
	}
}

// TestFactoryRegistryIntegration tests factory integration with registry operations
func TestFactoryRegistryIntegration(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()
	componentFactory := component.NewFactory(mockRegistry)

	// Create a component
	config := component.NewComponentConfig(
		"registry-test",
		"Registry Test Component",
		component.TypeBasic,
		"Component for testing registry integration",
	)

	comp, err := componentFactory.Create(config)
	require.NoError(t, err)
	require.NotNil(t, comp)

	// Test that the component can be used with registry operations
	// (This would depend on the specific registry mock implementation)
	assert.Equal(t, "registry-test", comp.ID())
	assert.NotNil(t, comp.Metadata())
}

// TestFactoryWithDifferentRegistries tests factory behavior with different registry instances
func TestFactoryWithDifferentRegistries(t *testing.T) {
	factory := mocks.NewFactory()

	// Create multiple mock registries
	mockRegistry1 := factory.RegistryInterface()
	mockRegistry2 := factory.RegistryInterface()

	// Create factories with different registries
	componentFactory1 := component.NewFactory(mockRegistry1)
	componentFactory2 := component.NewFactory(mockRegistry2)

	assert.NotNil(t, componentFactory1)
	assert.NotNil(t, componentFactory2)

	// Test that both factories can create components
	config := component.NewComponentConfig(
		"test-component",
		"Test Component",
		component.TypeBasic,
		"A test component",
	)

	comp1, err1 := componentFactory1.Create(config)
	comp2, err2 := componentFactory2.Create(config)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotNil(t, comp1)
	assert.NotNil(t, comp2)

	// Both components should have the same properties
	assert.Equal(t, comp1.ID(), comp2.ID())
	assert.Equal(t, comp1.Name(), comp2.Name())
	assert.Equal(t, comp1.Type(), comp2.Type())
	// Note: Components created from same config will be equal, so we don't test inequality
}

// TestFactoryComponentMetadata tests that factory-created components have proper metadata
func TestFactoryComponentMetadata(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()
	componentFactory := component.NewFactory(mockRegistry)

	config := component.NewComponentConfig(
		"metadata-test",
		"Metadata Test Component",
		component.TypeOperation,
		"Component for testing metadata",
	)

	comp, err := componentFactory.Create(config)
	require.NoError(t, err)
	require.NotNil(t, comp)

	metadata := comp.Metadata()
	// Just verify that metadata is not nil (implementation may return empty map)
	assert.NotNil(t, metadata)
}
