package component

import (
	"testing"

	"github.com/fintechain/skeleton/internal/domain/component/mocks"
)

func TestFactoryCreate(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create factory with mock dependencies
	factory := NewFactory(DefaultFactoryOptions{
		Logger: mockLogger,
	})

	// Create a component configuration
	config := ComponentConfig{
		ID:           "test-id",
		Name:         "Test Component",
		Type:         TypeBasic,
		Description:  "Test component description",
		Dependencies: []string{"dep1", "dep2"},
		Properties: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		},
	}

	// Create a component
	comp, err := factory.Create(config)
	if err != nil {
		t.Errorf("Failed to create component: %s", err)
	}

	// Verify component properties
	if comp.ID() != "test-id" {
		t.Errorf("Component has wrong ID: got %s, expected test-id", comp.ID())
	}

	if comp.Name() != "Test Component" {
		t.Errorf("Component has wrong name: got %s, expected Test Component", comp.Name())
	}

	if comp.Type() != TypeBasic {
		t.Errorf("Component has wrong type: got %s, expected %s", comp.Type(), TypeBasic)
	}

	// Verify metadata
	metadata := comp.Metadata()

	description, ok := metadata["description"]
	if !ok || description != "Test component description" {
		t.Errorf("Component metadata missing or has wrong description: %v", description)
	}

	dependencies, ok := metadata["dependencies"]
	if !ok {
		t.Error("Component metadata missing dependencies")
	} else {
		deps, ok := dependencies.([]string)
		if !ok || len(deps) != 2 || deps[0] != "dep1" || deps[1] != "dep2" {
			t.Errorf("Component metadata has wrong dependencies: %v", dependencies)
		}
	}

	val1, ok := metadata["key1"]
	if !ok || val1 != "value1" {
		t.Errorf("Component metadata missing or has wrong key1: %v", val1)
	}

	val2, ok := metadata["key2"]
	if !ok || val2 != 42 {
		t.Errorf("Component metadata missing or has wrong key2: %v", val2)
	}
}

func TestFactoryCreateInvalidType(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create factory with mock dependencies
	factory := NewFactory(DefaultFactoryOptions{
		Logger: mockLogger,
	})

	// Create a component configuration with an unregistered type
	config := ComponentConfig{
		ID:          "test-id",
		Name:        "Test Component",
		Type:        "invalid-type",
		Description: "Test component description",
	}

	// Try to create a component with an invalid type
	_, err := factory.Create(config)
	if err == nil {
		t.Error("Expected error when creating component with invalid type, but got nil")
	}

	// Check if error is the expected type
	if !IsComponentError(err, ErrInvalidComponent) {
		t.Errorf("Expected invalid component error, but got: %s", err)
	}
}

func TestFactoryRegisterTypeCreator(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create factory with mock dependencies
	factory := NewFactory(DefaultFactoryOptions{
		Logger: mockLogger,
	})

	// Register a custom component type creator
	factory.RegisterTypeCreator("custom-type", func(config ComponentConfig) (Component, error) {
		component := NewBaseComponent(config.ID, config.Name, "custom-type")
		component.SetMetadata("custom", true)
		return component, nil
	})

	// Create a component with the custom type
	config := ComponentConfig{
		ID:   "custom-id",
		Name: "Custom Component",
		Type: "custom-type",
	}

	comp, err := factory.Create(config)
	if err != nil {
		t.Errorf("Failed to create component with custom type: %s", err)
	}

	// Verify component properties
	if comp.ID() != "custom-id" {
		t.Errorf("Component has wrong ID: got %s, expected custom-id", comp.ID())
	}

	if comp.Type() != "custom-type" {
		t.Errorf("Component has wrong type: got %s, expected custom-type", comp.Type())
	}

	// Verify custom metadata
	metadata := comp.Metadata()
	customFlag, ok := metadata["custom"]
	if !ok || customFlag != true {
		t.Errorf("Component metadata missing or has wrong custom flag: %v", customFlag)
	}
}

func TestCreateFactory(t *testing.T) {
	// Test the factory method for backward compatibility
	factory := CreateFactory()

	if factory == nil {
		t.Error("CreateFactory returned nil")
	}

	// Verify we can use the factory normally
	config := ComponentConfig{
		ID:   "test-id",
		Name: "Test Component",
		Type: TypeBasic,
	}

	comp, err := factory.Create(config)
	if err != nil {
		t.Errorf("Failed to create component using CreateFactory-created factory: %s", err)
	}

	if comp.ID() != "test-id" {
		t.Errorf("Component created by CreateFactory has wrong ID: got %s, expected test-id", comp.ID())
	}
}
