// Package component provides the core interfaces and types for the component system.
package component

import (
	"time"
)

// ComponentType represents the type of a component.
type ComponentType string

const (
	// Basic component types
	TypeBasic       ComponentType = "basic"
	TypeOperation   ComponentType = "operation"
	TypeService     ComponentType = "service"
	TypeSystem      ComponentType = "system"
	TypeApplication ComponentType = "application"
)

// Metadata is a map of key-value pairs for component metadata.
type Metadata map[string]interface{}

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

// Component is the fundamental building block of the system.
type Component interface {
	// Identity
	ID() string
	Name() string
	Type() ComponentType

	// Metadata
	Metadata() Metadata

	// Lifecycle
	Initialize(ctx Context) error
	Dispose() error
}
