// Package context provides request-scoped context management for Fintechain Skeleton.
package context

import (
	"time"
)

// Context represents the execution context for components and operations.
type Context interface {
	// Value retrieves a value from the context by key.
	Value(key interface{}) interface{}

	// WithValue creates a new context with an additional key-value pair.
	WithValue(key, value interface{}) Context

	// Deadline returns the deadline for this context, if any.
	Deadline() (time.Time, bool)

	// Done returns a channel that's closed when the context is cancelled or times out.
	Done() <-chan struct{}

	// Err returns the error that caused the context to be cancelled.
	Err() error
}
