// Package context provides context implementation for the component system.
package context

import (
	stdctx "context"
	"time"

	"github.com/fintechain/skeleton/internal/domain/component"
)

// StdContext adapts Go's standard context to our component.Context interface.
type StdContext struct {
	ctx stdctx.Context
}

// NewContext creates a new component context from Go's standard context.
func NewContext(ctx stdctx.Context) component.Context {
	return &StdContext{ctx: ctx}
}

// WrapContext creates a component.Context from a standard Go context.
// This is the preferred factory method for creating a context.
func WrapContext(ctx stdctx.Context) component.Context {
	return NewContext(ctx)
}

// Value retrieves a value from the context.
func (c *StdContext) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

// WithValue adds a value to the context and returns a new context.
func (c *StdContext) WithValue(key, value interface{}) component.Context {
	return &StdContext{ctx: stdctx.WithValue(c.ctx, key, value)}
}

// Deadline returns the deadline for the context, if any.
func (c *StdContext) Deadline() (time.Time, bool) {
	return c.ctx.Deadline()
}

// Done returns a channel that's closed when the context is done.
func (c *StdContext) Done() <-chan struct{} {
	return c.ctx.Done()
}

// Err returns the error why the context was canceled, if any.
func (c *StdContext) Err() error {
	return c.ctx.Err()
}

// Background returns a new context with no values or cancellation.
func Background() component.Context {
	return &StdContext{ctx: stdctx.Background()}
}

// TODO returns a new context that is never canceled.
func TODO() component.Context {
	return &StdContext{ctx: stdctx.TODO()}
}

// WithCancel returns a new context and a cancel function.
func WithCancel(parent component.Context) (component.Context, func()) {
	// Type assertion is safe because we control all context creation
	// and only create StdContext instances in this package
	if p, ok := parent.(*StdContext); ok {
		ctx, cancel := stdctx.WithCancel(p.ctx)
		return &StdContext{ctx: ctx}, cancel
	}
	// Fallback to a new background context if parent is not a StdContext
	ctx, cancel := stdctx.WithCancel(stdctx.Background())
	return &StdContext{ctx: ctx}, cancel
}

// WithTimeout returns a new context with a timeout and a cancel function.
func WithTimeout(parent component.Context, timeout time.Duration) (component.Context, func()) {
	// Type assertion is safe because we control all context creation
	// and only create StdContext instances in this package
	if p, ok := parent.(*StdContext); ok {
		ctx, cancel := stdctx.WithTimeout(p.ctx, timeout)
		return &StdContext{ctx: ctx}, cancel
	}
	// Fallback to a new background context if parent is not a StdContext
	ctx, cancel := stdctx.WithTimeout(stdctx.Background(), timeout)
	return &StdContext{ctx: ctx}, cancel
}

// WithDeadline returns a new context with a deadline and a cancel function.
func WithDeadline(parent component.Context, deadline time.Time) (component.Context, func()) {
	// Type assertion is safe because we control all context creation
	// and only create StdContext instances in this package
	if p, ok := parent.(*StdContext); ok {
		ctx, cancel := stdctx.WithDeadline(p.ctx, deadline)
		return &StdContext{ctx: ctx}, cancel
	}
	// Fallback to a new background context if parent is not a StdContext
	ctx, cancel := stdctx.WithDeadline(stdctx.Background(), deadline)
	return &StdContext{ctx: ctx}, cancel
}
