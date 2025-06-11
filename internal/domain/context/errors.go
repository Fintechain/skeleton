// Package context provides interfaces and types for the context system.
package context

// Standard context error codes
const (
	// ErrContextNotFound is returned when a context doesn't exist
	ErrContextNotFound = "context.context_not_found"

	// ErrContextCanceled is returned when a context has been canceled
	ErrContextCanceled = "context.context_canceled"

	// ErrContextDeadlineExceeded is returned when a context deadline is exceeded
	ErrContextDeadlineExceeded = "context.context_deadline_exceeded"

	// ErrInvalidContextValue is returned when an invalid context value is provided
	ErrInvalidContextValue = "context.invalid_context_value"

	// ErrContextKeyNotFound is returned when a context key doesn't exist
	ErrContextKeyNotFound = "context.context_key_not_found"

	// ErrInvalidContextConfig is returned when invalid context configuration is provided
	ErrInvalidContextConfig = "context.invalid_context_config"
)
