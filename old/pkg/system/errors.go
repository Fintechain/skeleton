package system

import (
	"github.com/fintechain/skeleton/internal/domain/component"
)

// Error represents a domain-specific error from the component system.
// This is re-exported from the internal component package for public use.
type Error = component.Error

// NewError creates a new component error with the given code, message, and optional cause.
// This function provides structured error handling with:
// - Error codes (string constants for categorization)
// - Human-readable messages
// - Error chaining/wrapping support
// - Additional metadata through the WithDetail method
//
// Example usage:
//
//	err := system.NewError("system.startup_failed", "Failed to initialize system", nil)
//	err.WithDetail("component", "storage").WithDetail("reason", "connection_timeout")
func NewError(code, message string, cause error) *Error {
	return component.NewError(code, message, cause)
}

// IsComponentError checks if an error is a component error with the given code.
// This is useful for error handling and classification.
//
// Example usage:
//
//	if system.IsComponentError(err, "system.startup_failed") {
//	    // Handle system startup error
//	}
func IsComponentError(err error, code string) bool {
	return component.IsComponentError(err, code)
}

// System-specific error codes
const (
	// System startup and configuration errors
	ErrSystemStartup     = "system.startup_failed"
	ErrSystemShutdown    = "system.shutdown_failed"
	ErrInvalidConfig     = "system.invalid_config"
	ErrDependencyMissing = "system.dependency_missing"
)
