package component

import (
	"testing"
	"time"

	"github.com/ebanfa/skeleton/internal/domain/component/mocks"
)

func TestDefaultComponentInitialize(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create a DefaultComponent with mock dependencies
	component := NewDefaultComponent(DefaultComponentOptions{
		ID:          "test-id",
		Name:        "Test Component",
		Type:        TypeBasic,
		Description: "Test component description",
		Logger:      mockLogger,
	})

	// Check initial state
	if component.IsInitialized() {
		t.Error("Component should not be initialized initially")
	}

	if component.IsDisposed() {
		t.Error("Component should not be disposed initially")
	}

	// Initialize the component
	ctx := &mockContext{}
	err := component.Initialize(ctx)
	if err != nil {
		t.Errorf("Failed to initialize component: %s", err)
	}

	// Check initialized state
	if !component.IsInitialized() {
		t.Error("Component should be initialized after Initialize() call")
	}

	// Check that description was stored in metadata
	metadata := component.Metadata()
	description, ok := metadata["description"]
	if !ok || description != "Test component description" {
		t.Errorf("Component metadata missing or has wrong description: %v", description)
	}

	// Check that a log message was created
	if len(mockLogger.LogEntries) == 0 {
		t.Error("Expected log entries, but found none")
	}

	// Initialize again should be a no-op
	err = component.Initialize(ctx)
	if err != nil {
		t.Errorf("Second initialization should not return error: %s", err)
	}
}

func TestDefaultComponentDispose(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create a DefaultComponent with mock dependencies
	component := NewDefaultComponent(DefaultComponentOptions{
		ID:          "test-id",
		Name:        "Test Component",
		Type:        TypeBasic,
		Description: "Test component description",
		Logger:      mockLogger,
	})

	// Initialize the component
	ctx := &mockContext{}
	component.Initialize(ctx)

	// Clear log entries to check only dispose logs
	mockLogger.ClearLogEntries()

	// Dispose the component
	err := component.Dispose()
	if err != nil {
		t.Errorf("Failed to dispose component: %s", err)
	}

	// Check disposed state
	if !component.IsDisposed() {
		t.Error("Component should be disposed after Dispose() call")
	}

	// Check that a log message was created
	if len(mockLogger.LogEntries) == 0 {
		t.Error("Expected log entries, but found none")
	}

	// Dispose again should be a no-op
	err = component.Dispose()
	if err != nil {
		t.Errorf("Second disposal should not return error: %s", err)
	}
}

func TestDefaultComponentDescription(t *testing.T) {
	// Create a DefaultComponent with mock dependencies
	mockLogger := mocks.NewMockLogger()
	component := NewDefaultComponent(DefaultComponentOptions{
		ID:          "test-id",
		Name:        "Test Component",
		Type:        TypeBasic,
		Description: "Test component description",
		Logger:      mockLogger,
	})

	// Check description
	if component.Description() != "Test component description" {
		t.Errorf("Component has wrong description: got %s, expected Test component description",
			component.Description())
	}
}

// mockContext is a simple implementation of the Context interface for testing
type mockContext struct {
	values map[interface{}]interface{}
}

func (m *mockContext) Value(key interface{}) interface{} {
	if m.values == nil {
		return nil
	}
	return m.values[key]
}

func (m *mockContext) WithValue(key, value interface{}) Context {
	if m.values == nil {
		m.values = make(map[interface{}]interface{})
	}
	m.values[key] = value
	return m
}

func (m *mockContext) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (m *mockContext) Done() <-chan struct{} {
	return nil
}

func (m *mockContext) Err() error {
	return nil
}

func TestCreateDefaultComponent(t *testing.T) {
	// Test the factory method for backward compatibility
	component := CreateDefaultComponent("test-id", "Test Component", TypeBasic, "Test description")

	if component == nil {
		t.Error("CreateDefaultComponent returned nil")
	}

	// Verify basic properties
	if component.ID() != "test-id" {
		t.Errorf("Component has wrong ID: got %s, expected test-id", component.ID())
	}

	if component.Name() != "Test Component" {
		t.Errorf("Component has wrong name: got %s, expected Test Component", component.Name())
	}

	if component.Type() != TypeBasic {
		t.Errorf("Component has wrong type: got %s, expected %s", component.Type(), TypeBasic)
	}

	if component.Description() != "Test description" {
		t.Errorf("Component has wrong description: got %s, expected Test description", component.Description())
	}

	// Verify we can use the component normally
	ctx := &mockContext{}
	err := component.Initialize(ctx)
	if err != nil {
		t.Errorf("Failed to initialize component created with CreateDefaultComponent: %s", err)
	}

	if !component.IsInitialized() {
		t.Error("Component should be initialized after Initialize() call")
	}
}
