// Package context provides execution context interfaces and types.
package context

import (
	"github.com/fintechain/skeleton/internal/domain/context"
	contextImpl "github.com/fintechain/skeleton/internal/infrastructure/context"
)

// Re-export context interface
type Context = context.Context

// Factory functions

// NewContext creates a new framework context instance.
func NewContext() Context {
	return contextImpl.NewContext()
}

// WrapContext wraps an existing context with framework context capabilities.
func WrapContext(ctx Context) Context {
	return contextImpl.WrapContext(ctx)
}
