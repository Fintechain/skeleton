package context

import (
	"sync"
	"time"

	"github.com/fintechain/skeleton/internal/domain/context"
)

// FrameworkContext provides a concrete implementation of the Context interface.
type FrameworkContext struct {
	values      map[interface{}]interface{}
	deadline    time.Time
	hasDeadline bool
	done        chan struct{}
	err         error
	mu          sync.RWMutex
	cancel      func()
}

// NewContext creates a new framework context instance.
// This constructor accepts no dependencies to keep it simple and focused.
func NewContext() context.Context {
	return &FrameworkContext{
		values: make(map[interface{}]interface{}),
		done:   make(chan struct{}),
	}
}

// WrapContext wraps an existing context, preserving its values and state.
// This allows integration with existing context instances.
func WrapContext(ctx context.Context) context.Context {
	if ctx == nil {
		return NewContext()
	}

	// If it's already our type, return as-is
	if frameworkCtx, ok := ctx.(*FrameworkContext); ok {
		return frameworkCtx
	}

	// Create new context and copy values if possible
	newCtx := &FrameworkContext{
		values: make(map[interface{}]interface{}),
		done:   make(chan struct{}),
	}

	// Copy deadline if available
	if deadline, ok := ctx.Deadline(); ok {
		newCtx.deadline = deadline
		newCtx.hasDeadline = true
	}

	// Copy error if available
	if err := ctx.Err(); err != nil {
		newCtx.err = err
		close(newCtx.done)
	}

	return newCtx
}

// Value retrieves a value from the context.
func (c *FrameworkContext) Value(key interface{}) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.values[key]
}

// WithValue adds a value to the context and returns a new context.
func (c *FrameworkContext) WithValue(key, value interface{}) context.Context {
	c.mu.RLock()
	// Create a copy of existing values
	newValues := make(map[interface{}]interface{}, len(c.values)+1)
	for k, v := range c.values {
		newValues[k] = v
	}
	c.mu.RUnlock()

	// Add the new value
	newValues[key] = value

	// Create new context with copied values
	newCtx := &FrameworkContext{
		values:      newValues,
		deadline:    c.deadline,
		hasDeadline: c.hasDeadline,
		done:        make(chan struct{}),
		err:         c.err,
	}

	// If parent is already done, mark new context as done
	if c.isDone() {
		close(newCtx.done)
		newCtx.err = c.err
	}

	return newCtx
}

// Deadline returns the deadline for the context, if any.
func (c *FrameworkContext) Deadline() (time.Time, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.deadline, c.hasDeadline
}

// Done returns a channel that's closed when the context is done.
func (c *FrameworkContext) Done() <-chan struct{} {
	return c.done
}

// Err returns the error why the context was canceled, if any.
func (c *FrameworkContext) Err() error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.err
}

// isDone checks if the context is done (internal helper)
func (c *FrameworkContext) isDone() bool {
	select {
	case <-c.done:
		return true
	default:
		return false
	}
}

// WithTimeout creates a context with a timeout.
// This is a convenience function for creating contexts with deadlines.
func WithTimeout(parent context.Context, timeout time.Duration) (context.Context, func()) {
	return WithDeadline(parent, time.Now().Add(timeout))
}

// WithDeadline creates a context with a deadline.
// This is a convenience function for creating contexts with specific deadlines.
func WithDeadline(parent context.Context, deadline time.Time) (context.Context, func()) {
	if parent == nil {
		parent = NewContext()
	}

	ctx := parent.WithValue("__deadline__", deadline).(*FrameworkContext)
	ctx.deadline = deadline
	ctx.hasDeadline = true

	// Create a new done channel for this context
	ctx.done = make(chan struct{})

	// Set up cancellation
	cancel := func() {
		ctx.mu.Lock()
		defer ctx.mu.Unlock()

		if ctx.err == nil {
			ctx.err = &ContextCancelledError{}
			close(ctx.done)
		}
	}

	ctx.cancel = cancel

	// Set up deadline timer
	if time.Now().Before(deadline) {
		timer := time.AfterFunc(time.Until(deadline), func() {
			ctx.mu.Lock()
			defer ctx.mu.Unlock()

			if ctx.err == nil {
				ctx.err = &ContextDeadlineExceededError{}
				close(ctx.done)
			}
		})

		// Wrap cancel to also stop the timer
		originalCancel := cancel
		cancel = func() {
			timer.Stop()
			originalCancel()
		}
	}

	return ctx, cancel
}

// ContextCancelledError represents a context cancellation error.
type ContextCancelledError struct{}

func (e *ContextCancelledError) Error() string {
	return "context cancelled"
}

// ContextDeadlineExceededError represents a context deadline exceeded error.
type ContextDeadlineExceededError struct{}

func (e *ContextDeadlineExceededError) Error() string {
	return "context deadline exceeded"
}
