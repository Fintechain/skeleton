package mocks

import (
	"time"

	"github.com/fintechain/skeleton/internal/domain/component"
)

// MockExternalContext is a mock implementation of component.Context for testing
// It mocks external context dependencies, not contexts within the package
type MockExternalContext struct {
	// Fields to control mock behavior
	ValueFunc     func(key interface{}) interface{}
	WithValueFunc func(key, value interface{}) component.Context
	DeadlineFunc  func() (time.Time, bool)
	DoneFunc      func() <-chan struct{}
	ErrFunc       func() error

	// Track method calls for verification
	ValueCalls     []interface{}
	WithValueCalls []struct {
		Key   interface{}
		Value interface{}
	}
	DeadlineCalls int
	DoneCalls     int
	ErrCalls      int
}

// Value retrieves a value from the context
func (m *MockExternalContext) Value(key interface{}) interface{} {
	m.ValueCalls = append(m.ValueCalls, key)
	if m.ValueFunc != nil {
		return m.ValueFunc(key)
	}
	return nil
}

// WithValue adds a value to the context and returns a new context
func (m *MockExternalContext) WithValue(key, value interface{}) component.Context {
	m.WithValueCalls = append(m.WithValueCalls, struct {
		Key   interface{}
		Value interface{}
	}{key, value})
	if m.WithValueFunc != nil {
		return m.WithValueFunc(key, value)
	}
	return m
}

// Deadline returns the deadline for the context, if any
func (m *MockExternalContext) Deadline() (time.Time, bool) {
	m.DeadlineCalls++
	if m.DeadlineFunc != nil {
		return m.DeadlineFunc()
	}
	return time.Time{}, false
}

// Done returns a channel that's closed when the context is done
func (m *MockExternalContext) Done() <-chan struct{} {
	m.DoneCalls++
	if m.DoneFunc != nil {
		return m.DoneFunc()
	}
	return nil
}

// Err returns the error why the context was canceled, if any
func (m *MockExternalContext) Err() error {
	m.ErrCalls++
	if m.ErrFunc != nil {
		return m.ErrFunc()
	}
	return nil
}
