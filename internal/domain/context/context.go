package context

import (
	"time"
)

// Context represents the execution context for components.
type Context interface {
	// Value retrieves a value from the context.
	Value(key interface{}) interface{}

	// WithValue adds a value to the context and returns a new context.
	WithValue(key, value interface{}) Context

	// Deadline returns the deadline for the context, if any.
	Deadline() (time.Time, bool)

	// Done returns a channel that's closed when the context is done.
	Done() <-chan struct{}

	// Err returns the error why the context was canceled, if any.
	Err() error
}
