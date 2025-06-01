package component

import (
	"errors"
	"testing"
)

func TestErrorCreation(t *testing.T) {
	// Create a simple error
	err := NewError(ErrComponentNotFound, "component not found", nil)

	// Check error code and message
	if err.Code != ErrComponentNotFound {
		t.Errorf("Wrong error code: got %s, expected %s", err.Code, ErrComponentNotFound)
	}

	if err.Message != "component not found" {
		t.Errorf("Wrong error message: got %s, expected 'component not found'", err.Message)
	}

	// Check details are initially empty
	if len(err.Details) != 0 {
		t.Errorf("Initial details should be empty, but has %d items", len(err.Details))
	}

	// Check cause is nil
	if err.Cause != nil {
		t.Errorf("Cause should be nil, but got: %v", err.Cause)
	}

	// Check string representation
	expectedStr := "component.not_found: component not found"
	if err.Error() != expectedStr {
		t.Errorf("Wrong error string: got '%s', expected '%s'", err.Error(), expectedStr)
	}
}

func TestErrorWithCause(t *testing.T) {
	// Create an error with a cause
	cause := errors.New("underlying error")
	err := NewError(ErrComponentCreation, "failed to create component", cause)

	// Check cause
	if err.Cause != cause {
		t.Errorf("Wrong cause: got %v, expected %v", err.Cause, cause)
	}

	// Check string representation includes cause
	expectedStr := "component.creation_failed: failed to create component (cause: underlying error)"
	if err.Error() != expectedStr {
		t.Errorf("Wrong error string: got '%s', expected '%s'", err.Error(), expectedStr)
	}

	// Check unwrapping
	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Errorf("Unwrapped error should be the cause, got %v", unwrapped)
	}
}

func TestErrorWithDetails(t *testing.T) {
	// Create an error with details
	err := NewError(ErrDependencyNotFound, "dependency not found", nil)
	err.WithDetail("component_id", "comp1").WithDetail("dependency_id", "dep1")

	// Check details
	if len(err.Details) != 2 {
		t.Errorf("Expected 2 details, got %d", len(err.Details))
	}

	compID, ok := err.Details["component_id"]
	if !ok || compID != "comp1" {
		t.Errorf("Detail component_id missing or wrong: %v", compID)
	}

	depID, ok := err.Details["dependency_id"]
	if !ok || depID != "dep1" {
		t.Errorf("Detail dependency_id missing or wrong: %v", depID)
	}
}

func TestIsComponentError(t *testing.T) {
	// Create component errors
	err1 := NewError(ErrComponentNotFound, "component not found", nil)
	err2 := NewError(ErrComponentExists, "component already exists", nil)

	// Create a non-component error
	err3 := errors.New("generic error")

	// Test IsComponentError function
	if !IsComponentError(err1, ErrComponentNotFound) {
		t.Error("IsComponentError should return true for matching code")
	}

	if IsComponentError(err1, ErrComponentExists) {
		t.Error("IsComponentError should return false for non-matching code")
	}

	if IsComponentError(err2, ErrComponentNotFound) {
		t.Error("IsComponentError should return false for non-matching code")
	}

	if IsComponentError(err3, ErrComponentNotFound) {
		t.Error("IsComponentError should return false for non-component error")
	}

	if IsComponentError(nil, ErrComponentNotFound) {
		t.Error("IsComponentError should return false for nil error")
	}
}
