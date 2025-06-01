package component

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fintechain/skeleton/pkg/component"
	"github.com/fintechain/skeleton/pkg/context"
	"github.com/fintechain/skeleton/test/unit/mocks"
)

// TestNewBaseComponent tests the base component constructor function
func TestNewBaseComponent(t *testing.T) {
	config := component.NewComponentConfig(
		"test-component",
		"Test Component",
		component.TypeBasic,
		"A test component for unit testing",
	)

	comp := component.NewBaseComponent(config)

	assert.NotNil(t, comp)
	assert.Equal(t, "test-component", comp.ID())
	assert.Equal(t, "Test Component", comp.Name())
	assert.Equal(t, component.TypeBasic, comp.Type())
	assert.Equal(t, "A test component for unit testing", comp.Description())
}

// TestBaseComponentInterfaceCompliance verifies the implementation satisfies the domain interface
func TestBaseComponentInterfaceCompliance(t *testing.T) {
	config := component.NewComponentConfig(
		"test-component",
		"Test Component",
		component.TypeBasic,
		"A test component",
	)

	// Verify interface compliance
	var _ component.Component = component.NewBaseComponent(config)
}

// TestBaseComponentProperties tests component property access
func TestBaseComponentProperties(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		compName    string
		compType    component.ComponentType
		description string
	}{
		{
			name:        "basic component",
			id:          "basic-comp",
			compName:    "Basic Component",
			compType:    component.TypeBasic,
			description: "A basic component",
		},
		{
			name:        "operation component",
			id:          "op-comp",
			compName:    "Operation Component",
			compType:    component.TypeOperation,
			description: "An operation component",
		},
		{
			name:        "service component",
			id:          "svc-comp",
			compName:    "Service Component",
			compType:    component.TypeService,
			description: "A service component",
		},
		{
			name:        "system component",
			id:          "sys-comp",
			compName:    "System Component",
			compType:    component.TypeSystem,
			description: "A system component",
		},
		{
			name:        "application component",
			id:          "app-comp",
			compName:    "Application Component",
			compType:    component.TypeApplication,
			description: "An application component",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := component.NewComponentConfig(tt.id, tt.compName, tt.compType, tt.description)
			comp := component.NewBaseComponent(config)

			assert.Equal(t, tt.id, comp.ID())
			assert.Equal(t, tt.compName, comp.Name())
			assert.Equal(t, tt.compType, comp.Type())
			assert.Equal(t, tt.description, comp.Description())
		})
	}
}

// TestBaseComponentMetadata tests component metadata functionality
func TestBaseComponentMetadata(t *testing.T) {
	config := component.NewComponentConfig(
		"test-component",
		"Test Component",
		component.TypeBasic,
		"A test component",
	)

	comp := component.NewBaseComponent(config)
	metadata := comp.Metadata()

	// Just verify that metadata is not nil (implementation may return empty map)
	assert.NotNil(t, metadata)
}

// TestBaseComponentLifecycle tests component lifecycle methods
func TestBaseComponentLifecycle(t *testing.T) {
	config := component.NewComponentConfig(
		"test-component",
		"Test Component",
		component.TypeBasic,
		"A test component",
	)

	comp := component.NewBaseComponent(config)
	factory := mocks.NewFactory()
	mockSystem := factory.SystemInterface()

	// Use framework context for lifecycle operations
	ctx := context.NewContext()

	// Test initialization
	err := comp.Initialize(ctx, mockSystem)
	assert.NoError(t, err)

	// Test disposal (no context parameter)
	err = comp.Dispose()
	assert.NoError(t, err)
}

// TestBaseComponentWithDifferentConfigs tests component creation with various configurations
func TestBaseComponentWithDifferentConfigs(t *testing.T) {
	testCases := []struct {
		name        string
		id          string
		compName    string
		compType    component.ComponentType
		description string
		expectValid bool
	}{
		{
			name:        "valid basic component",
			id:          "valid-basic",
			compName:    "Valid Basic",
			compType:    component.TypeBasic,
			description: "Valid basic component",
			expectValid: true,
		},
		{
			name:        "valid operation component",
			id:          "valid-operation",
			compName:    "Valid Operation",
			compType:    component.TypeOperation,
			description: "Valid operation component",
			expectValid: true,
		},
		{
			name:        "valid service component",
			id:          "valid-service",
			compName:    "Valid Service",
			compType:    component.TypeService,
			description: "Valid service component",
			expectValid: true,
		},
		{
			name:        "component with empty description",
			id:          "empty-desc",
			compName:    "Empty Description",
			compType:    component.TypeBasic,
			description: "",
			expectValid: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := component.NewComponentConfig(tc.id, tc.compName, tc.compType, tc.description)
			comp := component.NewBaseComponent(config)

			if tc.expectValid {
				assert.NotNil(t, comp)
				assert.Equal(t, tc.id, comp.ID())
				assert.Equal(t, tc.compName, comp.Name())
				assert.Equal(t, tc.compType, comp.Type())
				assert.Equal(t, tc.description, comp.Description())
			}
		})
	}
}

// TestBaseComponentVersion tests component version functionality
func TestBaseComponentVersion(t *testing.T) {
	config := component.NewComponentConfig(
		"test-component",
		"Test Component",
		component.TypeBasic,
		"A test component",
	)

	comp := component.NewBaseComponent(config)

	// Test default version (if available)
	version := comp.Version()
	assert.NotEmpty(t, version) // Assuming components have a default version
}

// TestBaseComponentMultipleInstances tests creating multiple component instances
func TestBaseComponentMultipleInstances(t *testing.T) {
	configs := []component.ComponentConfig{
		component.NewComponentConfig("comp-1", "Component 1", component.TypeBasic, "First component"),
		component.NewComponentConfig("comp-2", "Component 2", component.TypeOperation, "Second component"),
		component.NewComponentConfig("comp-3", "Component 3", component.TypeService, "Third component"),
	}

	components := make([]component.Component, len(configs))

	// Create multiple components
	for i, config := range configs {
		components[i] = component.NewBaseComponent(config)
		assert.NotNil(t, components[i])
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

// TestBaseComponentInitializationWithMockSystem tests component initialization with mock system
func TestBaseComponentInitializationWithMockSystem(t *testing.T) {
	config := component.NewComponentConfig(
		"test-component",
		"Test Component",
		component.TypeBasic,
		"A test component",
	)

	comp := component.NewBaseComponent(config)
	factory := mocks.NewFactory()
	mockSystem := factory.SystemInterface()

	// Configure mock system behavior if needed
	// (This would depend on the mock implementation)

	ctx := context.NewContext()

	// Test initialization with mock system
	err := comp.Initialize(ctx, mockSystem)
	assert.NoError(t, err)

	// Test that component can be disposed after initialization
	err = comp.Dispose()
	assert.NoError(t, err)
}

// TestBaseComponentContextOperations tests component operations with framework context
func TestBaseComponentContextOperations(t *testing.T) {
	config := component.NewComponentConfig(
		"test-component",
		"Test Component",
		component.TypeBasic,
		"A test component",
	)

	comp := component.NewBaseComponent(config)
	factory := mocks.NewFactory()
	mockSystem := factory.SystemInterface()

	// Test with context containing values
	ctx := context.NewContext()
	ctx = ctx.WithValue("test-key", "test-value")

	// Test initialization with context containing values
	err := comp.Initialize(ctx, mockSystem)
	assert.NoError(t, err)

	// Test disposal with context
	err = comp.Dispose()
	assert.NoError(t, err)
}
