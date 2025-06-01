package component

import (
	"fmt"
)

// Common error codes for component operations
const (
	ErrComponentNotFound    = "component.not_found"
	ErrComponentExists      = "component.already_exists"
	ErrInvalidComponent     = "component.invalid"
	ErrComponentCreation    = "component.creation_failed"
	ErrDependencyNotFound   = "component.dependency_not_found"
	ErrInitializationFailed = "component.initialization_failed"
	ErrDisposalFailed       = "component.disposal_failed"
)

// Error represents a domain-specific error from the component system.
type Error struct {
	Code    string                 // Error code
	Message string                 // Human-readable message
	Details map[string]interface{} // Additional details
	Cause   error                  // Underlying cause
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %s)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying cause.
func (e *Error) Unwrap() error {
	return e.Cause
}

// NewError creates a new component error.
func NewError(code, message string, cause error) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Details: make(map[string]interface{}),
		Cause:   cause,
	}
}

// WithDetail adds a detail to the error and returns it.
func (e *Error) WithDetail(key string, value interface{}) *Error {
	e.Details[key] = value
	return e
}

// IsComponentError checks if an error is a component error with the given code.
func IsComponentError(err error, code string) bool {
	if err == nil {
		return false
	}

	compErr, ok := err.(*Error)
	if !ok {
		return false
	}

	return compErr.Code == code
}
