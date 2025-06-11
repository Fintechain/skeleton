package component

import (
	"testing"

	"github.com/fintechain/skeleton/internal/domain/component"
	infraComponent "github.com/fintechain/skeleton/internal/infrastructure/component"
	infraContext "github.com/fintechain/skeleton/internal/infrastructure/context"
	"github.com/fintechain/skeleton/test/unit/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewBaseComponent(t *testing.T) {
	config := component.ComponentConfig{
		ID:          "test-component",
		Name:        "Test Component",
		Type:        component.TypeComponent,
		Description: "Test component description",
		Version:     "1.0.0",
		Properties:  map[string]interface{}{"key": "value"},
	}

	comp := infraComponent.NewBaseComponent(config)
	assert.NotNil(t, comp)

	// Verify interface compliance
	var _ component.Component = comp

	// Test basic properties
	assert.Equal(t, component.ComponentID("test-component"), comp.ID())
	assert.Equal(t, "Test Component", comp.Name())
	assert.Equal(t, "Test component description", comp.Description())
	assert.Equal(t, "1.0.0", comp.Version())
	assert.Equal(t, component.TypeComponent, comp.Type())
	assert.NotNil(t, comp.Metadata())
	assert.Equal(t, "value", comp.Metadata()["key"])
}

func TestNewBaseComponentWithDefaults(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "minimal-component",
		Name: "Minimal Component",
		Type: component.TypeComponent,
		// No version specified - should default to "1.0.0"
	}

	comp := infraComponent.NewBaseComponent(config)
	assert.NotNil(t, comp)

	// Test default version
	assert.Equal(t, "1.0.0", comp.Version())
	assert.Equal(t, component.ComponentID("minimal-component"), comp.ID())
	assert.Equal(t, "Minimal Component", comp.Name())
}

func TestBaseComponentInitialize(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "test-component",
		Name: "Test Component",
		Type: component.TypeComponent,
	}

	comp := infraComponent.NewBaseComponent(config)
	factory := mocks.NewFactory()
	mockSystem := factory.SystemInterface()
	ctx := infraContext.NewContext()

	// Test initialization (base implementation should be no-op)
	err := comp.Initialize(ctx, mockSystem)
	assert.NoError(t, err)
}

func TestBaseComponentDispose(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "test-component",
		Name: "Test Component",
		Type: component.TypeComponent,
	}

	comp := infraComponent.NewBaseComponent(config)

	// Test disposal (base implementation should be no-op)
	err := comp.Dispose()
	assert.NoError(t, err)
}

func TestBaseComponentMetadata(t *testing.T) {
	tests := []struct {
		name       string
		properties map[string]interface{}
		expected   component.Metadata
	}{
		{
			name:       "with properties",
			properties: map[string]interface{}{"key1": "value1", "key2": 42},
			expected:   component.Metadata{"key1": "value1", "key2": 42},
		},
		{
			name:       "nil properties",
			properties: nil,
			expected:   nil,
		},
		{
			name:       "empty properties",
			properties: map[string]interface{}{},
			expected:   component.Metadata{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := component.ComponentConfig{
				ID:         "test-component",
				Name:       "Test Component",
				Type:       component.TypeComponent,
				Properties: tt.properties,
			}

			comp := infraComponent.NewBaseComponent(config)
			metadata := comp.Metadata()

			if tt.expected == nil {
				assert.Nil(t, metadata)
			} else {
				assert.Equal(t, tt.expected, metadata)
			}
		})
	}
}
