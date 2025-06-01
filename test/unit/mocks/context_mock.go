package mocks

import (
	"sync"
	"time"

	"github.com/fintechain/skeleton/pkg/context"
)

// MockContext provides a configurable mock implementation of the framework's context.Context interface.
// This is NOT Go's standard context - it's the framework's own context interface.
// It supports behavior configuration, error injection, call tracking, and state verification
// for comprehensive testing of components that depend on framework context functionality.
type MockContext struct {
	mu sync.RWMutex

	// Context state
	values   map[interface{}]interface{}
	deadline time.Time
	done     chan struct{}
	err      error

	// Behavior configuration
	shouldFail     bool
	failureError   string
	hasDeadline    bool
	forceDeadline  time.Time
	forceDone      bool
	forceError     error
	valueOverrides map[interface{}]interface{}

	// Call tracking
	callCount map[string]int
	lastCalls map[string][]interface{}

	// State verification
	valueCalls     []ValueCall
	withValueCalls []WithValueCall
	deadlineCalls  int
	doneCalls      int
	errCalls       int
}

// ValueCall represents a call to Value method.
type ValueCall struct {
	Key interface{}
}

// WithValueCall represents a call to WithValue method.
type WithValueCall struct {
	Key   interface{}
	Value interface{}
}

// NewMockContext creates a new configurable framework context mock.
func NewMockContext() *MockContext {
	return &MockContext{
		values:         make(map[interface{}]interface{}),
		done:           make(chan struct{}),
		valueOverrides: make(map[interface{}]interface{}),
		callCount:      make(map[string]int),
		lastCalls:      make(map[string][]interface{}),
	}
}

// Framework Context Interface Implementation

// Value retrieves a value from the context.
func (m *MockContext) Value(key interface{}) interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("Value", key)
	m.valueCalls = append(m.valueCalls, ValueCall{Key: key})

	// Check for value overrides first
	if value, exists := m.valueOverrides[key]; exists {
		return value
	}

	// Check actual values
	if value, exists := m.values[key]; exists {
		return value
	}

	return nil
}

// WithValue adds a value to the context and returns a new context.
func (m *MockContext) WithValue(key, value interface{}) context.Context {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("WithValue", key, value)
	m.withValueCalls = append(m.withValueCalls, WithValueCall{
		Key:   key,
		Value: value,
	})

	// Create a new mock context with the added value
	newMock := &MockContext{
		values:         make(map[interface{}]interface{}),
		done:           m.done,
		deadline:       m.deadline,
		err:            m.err,
		hasDeadline:    m.hasDeadline,
		valueOverrides: make(map[interface{}]interface{}),
		callCount:      make(map[string]int),
		lastCalls:      make(map[string][]interface{}),
	}

	// Copy existing values
	for k, v := range m.values {
		newMock.values[k] = v
	}
	for k, v := range m.valueOverrides {
		newMock.valueOverrides[k] = v
	}

	// Add the new value
	newMock.values[key] = value

	return newMock
}

// Deadline returns the deadline for the context, if any.
func (m *MockContext) Deadline() (time.Time, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("Deadline")
	m.deadlineCalls++

	if !m.forceDeadline.IsZero() {
		return m.forceDeadline, true
	}

	if m.hasDeadline {
		return m.deadline, true
	}

	return time.Time{}, false
}

// Done returns a channel that's closed when the context is done.
func (m *MockContext) Done() <-chan struct{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("Done")
	m.doneCalls++

	if m.forceDone {
		// Return a closed channel
		done := make(chan struct{})
		close(done)
		return done
	}

	return m.done
}

// Err returns the error why the context was canceled, if any.
func (m *MockContext) Err() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("Err")
	m.errCalls++

	if m.forceError != nil {
		return m.forceError
	}

	return m.err
}

// Mock Configuration Methods

// SetValue sets a value in the context.
func (m *MockContext) SetValue(key, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.values[key] = value
}

// SetValueOverride sets a value override that takes precedence over normal values.
func (m *MockContext) SetValueOverride(key, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.valueOverrides[key] = value
}

// SetDeadline sets the context deadline.
func (m *MockContext) SetDeadline(deadline time.Time, hasDeadline bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deadline = deadline
	m.hasDeadline = hasDeadline
}

// SetForceDeadline forces a specific deadline to be returned.
func (m *MockContext) SetForceDeadline(deadline time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.forceDeadline = deadline
}

// SetDone closes the done channel.
func (m *MockContext) SetDone() {
	m.mu.Lock()
	defer m.mu.Unlock()
	close(m.done)
}

// SetForceDone forces Done() to return a closed channel.
func (m *MockContext) SetForceDone(done bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.forceDone = done
}

// SetError sets the context error.
func (m *MockContext) SetError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.err = err
}

// SetForceError forces a specific error to be returned.
func (m *MockContext) SetForceError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.forceError = err
}

// SetShouldFail configures the mock to fail all operations.
func (m *MockContext) SetShouldFail(fail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail = fail
}

// SetFailureError sets the error message for failed operations.
func (m *MockContext) SetFailureError(err string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failureError = err
}

// State Verification Methods

// GetCallCount returns the number of times a method was called.
func (m *MockContext) GetCallCount(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.callCount[method]
}

// GetLastCall returns the parameters of the last call to a method.
func (m *MockContext) GetLastCall(method string) []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastCalls[method]
}

// GetValueCalls returns all calls to Value method.
func (m *MockContext) GetValueCalls() []ValueCall {
	m.mu.RLock()
	defer m.mu.RUnlock()
	calls := make([]ValueCall, len(m.valueCalls))
	copy(calls, m.valueCalls)
	return calls
}

// GetWithValueCalls returns all calls to WithValue method.
func (m *MockContext) GetWithValueCalls() []WithValueCall {
	m.mu.RLock()
	defer m.mu.RUnlock()
	calls := make([]WithValueCall, len(m.withValueCalls))
	copy(calls, m.withValueCalls)
	return calls
}

// GetDeadlineCalls returns the number of times Deadline was called.
func (m *MockContext) GetDeadlineCalls() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.deadlineCalls
}

// GetDoneCalls returns the number of times Done was called.
func (m *MockContext) GetDoneCalls() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.doneCalls
}

// GetErrCalls returns the number of times Err was called.
func (m *MockContext) GetErrCalls() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.errCalls
}

// HasValue returns true if the context has a value for the given key.
func (m *MockContext) HasValue(key interface{}) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, exists := m.valueOverrides[key]; exists {
		return true
	}

	_, exists := m.values[key]
	return exists
}

// IsDone returns true if the context is done.
func (m *MockContext) IsDone() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.forceDone {
		return true
	}

	select {
	case <-m.done:
		return true
	default:
		return false
	}
}

// HasDeadline returns true if the context has a deadline.
func (m *MockContext) HasDeadline() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if !m.forceDeadline.IsZero() {
		return true
	}

	return m.hasDeadline
}

// Reset clears all mock state and configuration.
func (m *MockContext) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.values = make(map[interface{}]interface{})
	m.deadline = time.Time{}
	m.done = make(chan struct{})
	m.err = nil
	m.shouldFail = false
	m.failureError = ""
	m.hasDeadline = false
	m.forceDeadline = time.Time{}
	m.forceDone = false
	m.forceError = nil
	m.valueOverrides = make(map[interface{}]interface{})
	m.callCount = make(map[string]int)
	m.lastCalls = make(map[string][]interface{})
	m.valueCalls = nil
	m.withValueCalls = nil
	m.deadlineCalls = 0
	m.doneCalls = 0
	m.errCalls = 0
}

// Helper Methods

// trackCall records a method call for verification.
func (m *MockContext) trackCall(method string, args ...interface{}) {
	m.callCount[method]++
	m.lastCalls[method] = args
}

// ContextMockBuilder provides a fluent interface for configuring context mocks.
type ContextMockBuilder struct {
	mock *MockContext
}

// NewContextMockBuilder creates a new context mock builder.
func NewContextMockBuilder() *ContextMockBuilder {
	return &ContextMockBuilder{
		mock: NewMockContext(),
	}
}

// WithValue adds a value to the context.
func (b *ContextMockBuilder) WithValue(key, value interface{}) *ContextMockBuilder {
	b.mock.SetValue(key, value)
	return b
}

// WithValueOverride sets a value override.
func (b *ContextMockBuilder) WithValueOverride(key, value interface{}) *ContextMockBuilder {
	b.mock.SetValueOverride(key, value)
	return b
}

// WithDeadline sets the context deadline.
func (b *ContextMockBuilder) WithDeadline(deadline time.Time) *ContextMockBuilder {
	b.mock.SetDeadline(deadline, true)
	return b
}

// WithForceDeadline forces a specific deadline.
func (b *ContextMockBuilder) WithForceDeadline(deadline time.Time) *ContextMockBuilder {
	b.mock.SetForceDeadline(deadline)
	return b
}

// WithDone marks the context as done.
func (b *ContextMockBuilder) WithDone() *ContextMockBuilder {
	b.mock.SetDone()
	return b
}

// WithForceDone forces the context to appear done.
func (b *ContextMockBuilder) WithForceDone(done bool) *ContextMockBuilder {
	b.mock.SetForceDone(done)
	return b
}

// WithError sets the context error.
func (b *ContextMockBuilder) WithError(err error) *ContextMockBuilder {
	b.mock.SetError(err)
	return b
}

// WithForceError forces a specific error.
func (b *ContextMockBuilder) WithForceError(err error) *ContextMockBuilder {
	b.mock.SetForceError(err)
	return b
}

// WithFailure configures the mock to fail operations.
func (b *ContextMockBuilder) WithFailure(fail bool) *ContextMockBuilder {
	b.mock.SetShouldFail(fail)
	return b
}

// WithFailureError sets the error message for failed operations.
func (b *ContextMockBuilder) WithFailureError(err string) *ContextMockBuilder {
	b.mock.SetFailureError(err)
	return b
}

// Build returns the configured mock context.
func (b *ContextMockBuilder) Build() context.Context {
	return b.mock
}
