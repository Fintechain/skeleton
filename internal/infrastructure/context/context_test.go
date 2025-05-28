package context

import (
	"context"
	"testing"
	"time"

	"github.com/fintechain/skeleton/internal/domain/component"
)

func TestNewContext(t *testing.T) {
	stdCtx := context.Background()
	ctx := NewContext(stdCtx)

	if ctx == nil {
		t.Fatal("NewContext returned nil")
	}
}

func TestWrapContext(t *testing.T) {
	stdCtx := context.Background()
	ctx := WrapContext(stdCtx)

	if ctx == nil {
		t.Fatal("WrapContext returned nil")
	}
}

func TestContextValue(t *testing.T) {
	const key = "test-key"
	const value = "test-value"

	stdCtx := context.WithValue(context.Background(), key, value)
	ctx := NewContext(stdCtx)

	result := ctx.Value(key)
	if result != value {
		t.Errorf("Expected value '%v', got '%v'", value, result)
	}
}

func TestContextWithValue(t *testing.T) {
	const key = "test-key"
	const value = "test-value"

	ctx := Background()
	ctx = ctx.WithValue(key, value)

	result := ctx.Value(key)
	if result != value {
		t.Errorf("Expected value '%v', got '%v'", value, result)
	}
}

func TestBackgroundContext(t *testing.T) {
	ctx := Background()

	if ctx == nil {
		t.Fatal("Background returned nil")
	}

	select {
	case <-ctx.Done():
		t.Error("Background context should not be done")
	default:
		// Expected - not done
	}
}

func TestTODOContext(t *testing.T) {
	ctx := TODO()

	if ctx == nil {
		t.Fatal("TODO returned nil")
	}

	select {
	case <-ctx.Done():
		t.Error("TODO context should not be done")
	default:
		// Expected - not done
	}
}

func TestWithCancel(t *testing.T) {
	parent := Background()
	ctx, cancel := WithCancel(parent)

	if ctx == nil {
		t.Fatal("WithCancel returned nil context")
	}

	// Not cancelled yet
	select {
	case <-ctx.Done():
		t.Error("Context should not be done before cancellation")
	default:
		// Expected - not done
	}

	// Cancel the context
	cancel()

	// Should be cancelled now
	select {
	case <-ctx.Done():
		// Expected - done
	case <-time.After(100 * time.Millisecond):
		t.Error("Context was not cancelled after calling cancel")
	}

	if ctx.Err() != context.Canceled {
		t.Errorf("Expected error %v, got %v", context.Canceled, ctx.Err())
	}
}

func TestWithTimeout(t *testing.T) {
	parent := Background()
	ctx, cancel := WithTimeout(parent, 50*time.Millisecond)
	defer cancel()

	if ctx == nil {
		t.Fatal("WithTimeout returned nil context")
	}

	// Not timed out yet
	select {
	case <-ctx.Done():
		t.Error("Context should not be done immediately")
	default:
		// Expected - not done
	}

	// Wait for timeout
	time.Sleep(100 * time.Millisecond)

	// Should be done now
	select {
	case <-ctx.Done():
		// Expected - done
	default:
		t.Error("Context should be done after timeout")
	}

	if ctx.Err() != context.DeadlineExceeded {
		t.Errorf("Expected error %v, got %v", context.DeadlineExceeded, ctx.Err())
	}
}

func TestWithDeadline(t *testing.T) {
	parent := Background()
	deadline := time.Now().Add(50 * time.Millisecond)
	ctx, cancel := WithDeadline(parent, deadline)
	defer cancel()

	if ctx == nil {
		t.Fatal("WithDeadline returned nil context")
	}

	// Check deadline
	d, ok := ctx.Deadline()
	if !ok {
		t.Error("Context should have a deadline")
	}

	if d.Sub(deadline) > time.Millisecond {
		t.Errorf("Expected deadline near %v, got %v", deadline, d)
	}

	// Wait for deadline
	time.Sleep(100 * time.Millisecond)

	// Should be done now
	select {
	case <-ctx.Done():
		// Expected - done
	default:
		t.Error("Context should be done after deadline")
	}
}

// TestFallbackWithNonStdContext tests the fallback behavior with non-StdContext implementations
func TestFallbackWithNonStdContext(t *testing.T) {
	// Create a mock implementation of component.Context that is not a StdContext
	mockCtx := &mockContext{}

	// Test WithCancel
	ctx1, cancel1 := WithCancel(mockCtx)
	if ctx1 == nil {
		t.Error("WithCancel should create a new context when given a non-StdContext")
	}
	cancel1()

	// Test WithTimeout
	ctx2, cancel2 := WithTimeout(mockCtx, 50*time.Millisecond)
	if ctx2 == nil {
		t.Error("WithTimeout should create a new context when given a non-StdContext")
	}
	cancel2()

	// Test WithDeadline
	ctx3, cancel3 := WithDeadline(mockCtx, time.Now().Add(50*time.Millisecond))
	if ctx3 == nil {
		t.Error("WithDeadline should create a new context when given a non-StdContext")
	}
	cancel3()
}

// mockContext is a simple implementation of component.Context for testing the fallback cases
type mockContext struct{}

func (m *mockContext) Value(key interface{}) interface{} {
	return nil
}

func (m *mockContext) WithValue(key, value interface{}) component.Context {
	return m
}

func (m *mockContext) Deadline() (time.Time, bool) {
	return time.Time{}, false
}

func (m *mockContext) Done() <-chan struct{} {
	ch := make(chan struct{})
	return ch
}

func (m *mockContext) Err() error {
	return nil
}

// Component is a local interface only for testing to avoid import cycles
type Component interface {
	Value(key interface{}) interface{}
	WithValue(key, value interface{}) Component
	Deadline() (time.Time, bool)
	Done() <-chan struct{}
	Err() error
}
