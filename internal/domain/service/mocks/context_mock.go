package mocks

import (
	"time"

	"github.com/fintechain/skeleton/internal/domain/component"
)

// MockContext implements the component.Context interface for testing
type MockContext struct {
	// Fields to control mock behavior
	ValueFunc     func(key interface{}) interface{}
	WithValueFunc func(key, val interface{}) component.Context
	DeadlineFunc  func() (time.Time, bool)
	DoneFunc      func() <-chan struct{}
	ErrFunc       func() error

	// Fields to track method calls
	ValueCalls     []interface{}
	WithValueCalls []struct{ Key, Value interface{} }
	DeadlineCalls  int
	DoneCalls      int
	ErrCalls       int
}

// Value retrieves a value from the context
func (m *MockContext) Value(key interface{}) interface{} {
	m.ValueCalls = append(m.ValueCalls, key)
	if m.ValueFunc != nil {
		return m.ValueFunc(key)
	}
	return nil
}

// WithValue adds a value to the context and returns a new context
func (m *MockContext) WithValue(key, value interface{}) component.Context {
	m.WithValueCalls = append(m.WithValueCalls, struct{ Key, Value interface{} }{key, value})
	if m.WithValueFunc != nil {
		return m.WithValueFunc(key, value)
	}
	return m
}

// Deadline returns the deadline for the context, if any
func (m *MockContext) Deadline() (time.Time, bool) {
	m.DeadlineCalls++
	if m.DeadlineFunc != nil {
		return m.DeadlineFunc()
	}
	return time.Time{}, false
}

// Done returns a channel that's closed when the context is done
func (m *MockContext) Done() <-chan struct{} {
	m.DoneCalls++
	if m.DoneFunc != nil {
		return m.DoneFunc()
	}
	// Return a nil channel by default
	return nil
}

// Err returns the error why the context was canceled, if any
func (m *MockContext) Err() error {
	m.ErrCalls++
	if m.ErrFunc != nil {
		return m.ErrFunc()
	}
	return nil
}

// NewMockContext creates a new mock context with default behavior
func NewMockContext() *MockContext {
	return &MockContext{
		ValueCalls:     make([]interface{}, 0),
		WithValueCalls: make([]struct{ Key, Value interface{} }, 0),
	}
}
