// Package context provides concrete implementations for the context domain interfaces.
package context

import (
	"errors"
	"sync"
	"time"

	"github.com/fintechain/skeleton/internal/domain/context"
)

// DomainContext implements the domain Context interface with thread-safe value storage,
// deadline management, and cancellation support.
type DomainContext struct {
	values   map[interface{}]interface{}
	deadline time.Time
	done     chan struct{}
	err      error
	mu       sync.RWMutex
}

// NewContext creates a new domain context instance.
// This is the primary constructor for creating context instances.
func NewContext() *DomainContext {
	return &DomainContext{
		values: make(map[interface{}]interface{}),
		done:   make(chan struct{}),
	}
}

// NewContextWithDeadline creates a new domain context with a deadline.
// The context will be automatically cancelled when the deadline is reached.
func NewContextWithDeadline(deadline time.Time) *DomainContext {
	ctx := &DomainContext{
		values:   make(map[interface{}]interface{}),
		deadline: deadline,
		done:     make(chan struct{}),
	}

	// Start deadline monitoring if deadline is set and in the future
	if !deadline.IsZero() && deadline.After(time.Now()) {
		go ctx.monitorDeadline()
	}

	return ctx
}

// NewContextWithTimeout creates a new domain context with a timeout duration.
// The context will be automatically cancelled after the timeout duration.
func NewContextWithTimeout(timeout time.Duration) *DomainContext {
	if timeout <= 0 {
		// Return immediately cancelled context for non-positive timeout
		ctx := &DomainContext{
			values: make(map[interface{}]interface{}),
			done:   make(chan struct{}),
			err:    errors.New(context.ErrContextDeadlineExceeded),
		}
		close(ctx.done)
		return ctx
	}

	return NewContextWithDeadline(time.Now().Add(timeout))
}

// Value retrieves a value from the context by key.
// Returns nil if the key is not found.
func (c *DomainContext) Value(key interface{}) interface{} {
	if key == nil {
		return nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.values[key]
}

// WithValue creates a new context with an additional key-value pair.
// The new context inherits all values from the parent context.
func (c *DomainContext) WithValue(key, value interface{}) context.Context {
	if key == nil {
		return c // Return same context if key is nil
	}

	c.mu.RLock()
	// Copy existing values
	newValues := make(map[interface{}]interface{}, len(c.values)+1)
	for k, v := range c.values {
		newValues[k] = v
	}
	deadline := c.deadline
	err := c.err
	c.mu.RUnlock()

	// Add new value
	newValues[key] = value

	newCtx := &DomainContext{
		values:   newValues,
		deadline: deadline,
		done:     make(chan struct{}),
		err:      err,
	}

	// If parent context is already cancelled, cancel the new context too
	select {
	case <-c.done:
		close(newCtx.done)
	default:
		// Start deadline monitoring if deadline is set and in the future
		if !deadline.IsZero() && deadline.After(time.Now()) && err == nil {
			go newCtx.monitorDeadline()
		}
	}

	return newCtx
}

// Deadline returns the deadline for this context, if any.
// Returns zero time and false if no deadline is set.
func (c *DomainContext) Deadline() (time.Time, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.deadline.IsZero() {
		return time.Time{}, false
	}
	return c.deadline, true
}

// Done returns a channel that's closed when the context is cancelled or times out.
// This channel can be used in select statements for cancellation handling.
func (c *DomainContext) Done() <-chan struct{} {
	return c.done
}

// Err returns the error that caused the context to be cancelled.
// Returns nil if the context is not cancelled.
func (c *DomainContext) Err() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.err
}

// Cancel manually cancels the context with a cancellation error.
// This is useful for explicit cancellation scenarios.
func (c *DomainContext) Cancel() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.err != nil {
		return // Already cancelled
	}

	c.err = errors.New(context.ErrContextCanceled)
	close(c.done)
}

// monitorDeadline runs in a goroutine to monitor deadline expiration.
// It automatically cancels the context when the deadline is reached.
func (c *DomainContext) monitorDeadline() {
	c.mu.RLock()
	deadline := c.deadline
	c.mu.RUnlock()

	if deadline.IsZero() {
		return
	}

	timer := time.NewTimer(time.Until(deadline))
	defer timer.Stop()

	select {
	case <-timer.C:
		c.mu.Lock()
		if c.err == nil { // Only cancel if not already cancelled
			c.err = errors.New(context.ErrContextDeadlineExceeded)
			close(c.done)
		}
		c.mu.Unlock()
	case <-c.done:
		// Context was cancelled before deadline
		return
	}
}

// IsCancelled returns true if the context has been cancelled.
// This is a convenience method for checking cancellation status.
func (c *DomainContext) IsCancelled() bool {
	select {
	case <-c.done:
		return true
	default:
		return false
	}
}

// WrapContext creates a domain context from a standard Go context.
// This provides a bridge between Go's standard context and the domain context.
func WrapContext(stdCtx interface{}) context.Context {
	// For now, create a new domain context
	// In a full implementation, this could bridge to Go's standard context
	return NewContext()
}
