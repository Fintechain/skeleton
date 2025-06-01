package component

import (
	"testing"
)

func TestComponentTypes(t *testing.T) {
	// Test component type constants
	if TypeBasic != "basic" {
		t.Errorf("Wrong TypeBasic value: got %s, expected 'basic'", TypeBasic)
	}

	if TypeOperation != "operation" {
		t.Errorf("Wrong TypeOperation value: got %s, expected 'operation'", TypeOperation)
	}

	if TypeService != "service" {
		t.Errorf("Wrong TypeService value: got %s, expected 'service'", TypeService)
	}

	if TypeSystem != "system" {
		t.Errorf("Wrong TypeSystem value: got %s, expected 'system'", TypeSystem)
	}

	if TypeApplication != "application" {
		t.Errorf("Wrong TypeApplication value: got %s, expected 'application'", TypeApplication)
	}
}

func TestMetadata(t *testing.T) {
	// Create metadata
	meta := Metadata{
		"string": "value",
		"int":    42,
		"bool":   true,
	}

	// Check values
	if v, ok := meta["string"]; !ok || v != "value" {
		t.Errorf("Metadata string value wrong or missing: %v", v)
	}

	if v, ok := meta["int"]; !ok || v != 42 {
		t.Errorf("Metadata int value wrong or missing: %v", v)
	}

	if v, ok := meta["bool"]; !ok || v != true {
		t.Errorf("Metadata bool value wrong or missing: %v", v)
	}
}

func TestContext(t *testing.T) {
	// Create a context
	ctx := &mockContext{}

	// Test empty context
	if v := ctx.Value("key"); v != nil {
		t.Errorf("Value from empty context should be nil, got %v", v)
	}

	// Test adding values
	ctx2 := ctx.WithValue("key1", "value1")
	if v := ctx2.Value("key1"); v != "value1" {
		t.Errorf("Context value wrong: got %v, expected 'value1'", v)
	}

	// Add another value
	ctx3 := ctx2.WithValue("key2", 42)
	if v := ctx3.Value("key1"); v != "value1" {
		t.Errorf("First context value lost: got %v, expected 'value1'", v)
	}

	if v := ctx3.Value("key2"); v != 42 {
		t.Errorf("Second context value wrong: got %v, expected 42", v)
	}

	// Test deadline (mock implementation returns zero time)
	deadline, ok := ctx3.Deadline()
	if ok || !deadline.IsZero() {
		t.Errorf("Mock context deadline should be zero time and not set, got %v, %v", deadline, ok)
	}

	// Test done channel (mock implementation returns nil)
	if ch := ctx3.Done(); ch != nil {
		t.Error("Mock context done channel should be nil")
	}

	// Test error (mock implementation returns nil)
	if err := ctx3.Err(); err != nil {
		t.Errorf("Mock context err should be nil, got %v", err)
	}
}
