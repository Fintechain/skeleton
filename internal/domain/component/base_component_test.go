package component

import (
	"testing"

	"github.com/fintechain/skeleton/internal/domain/component/mocks"
)

func TestBaseComponentCreation(t *testing.T) {
	// Test legacy constructor
	comp1 := NewBaseComponent("test-id-1", "Test Component 1", TypeBasic)

	if comp1.ID() != "test-id-1" {
		t.Errorf("Component has wrong ID: got %s, expected test-id-1", comp1.ID())
	}

	if comp1.Name() != "Test Component 1" {
		t.Errorf("Component has wrong name: got %s, expected Test Component 1", comp1.Name())
	}

	if comp1.Type() != TypeBasic {
		t.Errorf("Component has wrong type: got %s, expected %s", comp1.Type(), TypeBasic)
	}

	// Test new constructor with dependency injection
	mockLogger := mocks.NewMockLogger()
	comp2 := NewBaseComponentWithOptions(BaseComponentOptions{
		ID:     "test-id-2",
		Name:   "Test Component 2",
		Type:   TypeService,
		Logger: mockLogger,
	})

	if comp2.ID() != "test-id-2" {
		t.Errorf("Component has wrong ID: got %s, expected test-id-2", comp2.ID())
	}

	if comp2.Name() != "Test Component 2" {
		t.Errorf("Component has wrong name: got %s, expected Test Component 2", comp2.Name())
	}

	if comp2.Type() != TypeService {
		t.Errorf("Component has wrong type: got %s, expected %s", comp2.Type(), TypeService)
	}
}

func TestBaseComponentMetadata(t *testing.T) {
	// Create a component with dependencies
	mockLogger := mocks.NewMockLogger()
	component := NewBaseComponentWithOptions(BaseComponentOptions{
		ID:     "test-id",
		Name:   "Test Component",
		Type:   TypeBasic,
		Logger: mockLogger,
	})

	// Initially metadata should be empty
	metadata := component.Metadata()
	if len(metadata) != 0 {
		t.Errorf("Initial metadata should be empty, but has %d items", len(metadata))
	}

	// Set metadata
	component.SetMetadata("key1", "value1")
	component.SetMetadata("key2", 42)

	// Check metadata
	metadata = component.Metadata()
	if len(metadata) != 2 {
		t.Errorf("Expected 2 metadata items, but got %d", len(metadata))
	}

	val1, ok := metadata["key1"]
	if !ok || val1 != "value1" {
		t.Errorf("Metadata missing or has wrong key1: %v", val1)
	}

	val2, ok := metadata["key2"]
	if !ok || val2 != 42 {
		t.Errorf("Metadata missing or has wrong key2: %v", val2)
	}

	// Update metadata
	component.SetMetadata("key1", "updated")

	// Check updated metadata
	metadata = component.Metadata()
	val1, ok = metadata["key1"]
	if !ok || val1 != "updated" {
		t.Errorf("Metadata not updated correctly for key1: %v", val1)
	}
}

func TestBaseComponentLifecycle(t *testing.T) {
	// Create a component with dependencies
	mockLogger := mocks.NewMockLogger()
	component := NewBaseComponentWithOptions(BaseComponentOptions{
		ID:     "test-id",
		Name:   "Test Component",
		Type:   TypeBasic,
		Logger: mockLogger,
	})

	// Initialize should succeed (no-op in base component)
	ctx := &mockContext{}
	err := component.Initialize(ctx)
	if err != nil {
		t.Errorf("Initialize should succeed, but got error: %s", err)
	}

	// Dispose should succeed (no-op in base component)
	err = component.Dispose()
	if err != nil {
		t.Errorf("Dispose should succeed, but got error: %s", err)
	}
}
