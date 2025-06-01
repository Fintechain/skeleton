package mocks

import (
	"time"

	"github.com/fintechain/skeleton/internal/domain/component"
)

// MockContext implements the component.Context interface for testing
type MockContext struct {
	// Function fields for controlling behavior
	ValueFunc     func(key interface{}) interface{}
	WithValueFunc func(key, value interface{}) component.Context
	DeadlineFunc  func() (time.Time, bool)
	DoneFunc      func() <-chan struct{}
	ErrFunc       func() error

	// Call tracking for verification
	ValueCalls     []interface{}
	WithValueCalls []struct {
		Key   interface{}
		Value interface{}
	}
	DeadlineCalls int
	DoneCalls     int
	ErrCalls      int
}

// Value mocks the Value method
func (m *MockContext) Value(key interface{}) interface{} {
	m.ValueCalls = append(m.ValueCalls, key)
	if m.ValueFunc != nil {
		return m.ValueFunc(key)
	}
	return nil
}

// WithValue mocks the WithValue method
func (m *MockContext) WithValue(key, value interface{}) component.Context {
	m.WithValueCalls = append(m.WithValueCalls, struct {
		Key   interface{}
		Value interface{}
	}{key, value})
	if m.WithValueFunc != nil {
		return m.WithValueFunc(key, value)
	}
	return m
}

// Deadline mocks the Deadline method
func (m *MockContext) Deadline() (time.Time, bool) {
	m.DeadlineCalls++
	if m.DeadlineFunc != nil {
		return m.DeadlineFunc()
	}
	return time.Time{}, false
}

// Done mocks the Done method
func (m *MockContext) Done() <-chan struct{} {
	m.DoneCalls++
	if m.DoneFunc != nil {
		return m.DoneFunc()
	}
	// Return a never-closing channel by default
	return make(chan struct{})
}

// Err mocks the Err method
func (m *MockContext) Err() error {
	m.ErrCalls++
	if m.ErrFunc != nil {
		return m.ErrFunc()
	}
	return nil
}
