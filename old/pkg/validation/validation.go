// Package validation provides common validation utilities for the skeleton framework.
package validation

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"github.com/fintechain/skeleton/internal/domain/component"
)

// ===== VALIDATION INTERFACES =====

// Validator defines the interface for validation functions.
type Validator interface {
	Validate(value interface{}) error
}

// ValidatorFunc is a function type that implements the Validator interface.
type ValidatorFunc func(value interface{}) error

// Validate implements the Validator interface for ValidatorFunc.
func (f ValidatorFunc) Validate(value interface{}) error {
	return f(value)
}

// ===== BASIC VALIDATION FUNCTIONS =====

// Required validates that a value is not nil, empty string, or zero value.
//
// Example usage:
//
//	if err := validation.Required("name", userName); err != nil {
//	    // Handle validation error
//	}
func Required(fieldName string, value interface{}) error {
	if value == nil {
		return NewError(ErrRequired, fmt.Sprintf("field '%s' is required", fieldName), nil).
			WithDetail("field", fieldName)
	}

	switch v := value.(type) {
	case string:
		if strings.TrimSpace(v) == "" {
			return NewError(ErrRequired, fmt.Sprintf("field '%s' cannot be empty", fieldName), nil).
				WithDetail("field", fieldName)
		}
	case []interface{}:
		if len(v) == 0 {
			return NewError(ErrRequired, fmt.Sprintf("field '%s' cannot be empty", fieldName), nil).
				WithDetail("field", fieldName)
		}
	}

	return nil
}

// MinLength validates that a string has at least the specified minimum length.
//
// Example usage:
//
//	if err := validation.MinLength("password", password, 8); err != nil {
//	    // Handle validation error
//	}
func MinLength(fieldName string, value string, minLen int) error {
	if len(value) < minLen {
		return NewError(ErrMinLength,
			fmt.Sprintf("field '%s' must be at least %d characters long", fieldName, minLen), nil).
			WithDetail("field", fieldName).
			WithDetail("minLength", fmt.Sprintf("%d", minLen)).
			WithDetail("actualLength", fmt.Sprintf("%d", len(value)))
	}
	return nil
}

// Email validates that a string is a valid email address.
//
// Example usage:
//
//	if err := validation.Email("email", userEmail); err != nil {
//	    // Handle validation error
//	}
func Email(fieldName string, value string) error {
	if value == "" {
		return NewError(ErrRequired, fmt.Sprintf("field '%s' is required", fieldName), nil).
			WithDetail("field", fieldName)
	}

	_, err := mail.ParseAddress(value)
	if err != nil {
		return NewError(ErrInvalidFormat,
			fmt.Sprintf("field '%s' must be a valid email address", fieldName), err).
			WithDetail("field", fieldName).
			WithDetail("value", value)
	}
	return nil
}

// Regex validates that a string matches the specified regular expression pattern.
//
// Example usage:
//
//	phonePattern := `^\+?[1-9]\d{1,14}$`
//	if err := validation.Regex("phone", phoneNumber, phonePattern); err != nil {
//	    // Handle validation error
//	}
func Regex(fieldName string, value string, pattern string) error {
	matched, err := regexp.MatchString(pattern, value)
	if err != nil {
		return NewError(ErrInvalidPattern,
			fmt.Sprintf("invalid regex pattern for field '%s'", fieldName), err).
			WithDetail("field", fieldName).
			WithDetail("pattern", pattern)
	}

	if !matched {
		return NewError(ErrInvalidFormat,
			fmt.Sprintf("field '%s' does not match required pattern", fieldName), nil).
			WithDetail("field", fieldName).
			WithDetail("value", value).
			WithDetail("pattern", pattern)
	}
	return nil
}

// ===== ERROR CONSTANTS =====

// Common validation error codes
const (
	ErrRequired       = "validation.required"
	ErrMinLength      = "validation.min_length"
	ErrInvalidFormat  = "validation.invalid_format"
	ErrInvalidPattern = "validation.invalid_pattern"
)

// ===== ERROR HANDLING =====

// Error represents a domain-specific error from the validation system.
type Error = component.Error

// NewError creates a new validation error with the given code, message, and optional cause.
func NewError(code, message string, cause error) *Error {
	return component.NewError(code, message, cause)
}

// IsValidationError checks if an error is a validation error with the given code.
func IsValidationError(err error, code string) bool {
	return component.IsComponentError(err, code)
}
