package context

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	domainContext "github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/infrastructure/context"
)

func TestNewContext(t *testing.T) {
	ctx := context.NewContext()

	// Verify interface compliance
	var _ domainContext.Context = ctx

	// Verify initial state
	assert.NotNil(t, ctx)
	assert.NotNil(t, ctx.Done())
	assert.Nil(t, ctx.Err())
	assert.False(t, ctx.IsCancelled())

	// Verify deadline is not set
	deadline, hasDeadline := ctx.Deadline()
	assert.False(t, hasDeadline)
	assert.True(t, deadline.IsZero())
}

func TestNewContextWithDeadline(t *testing.T) {
	tests := []struct {
		name        string
		deadline    time.Time
		expectError bool
		description string
	}{
		{
			name:        "future deadline",
			deadline:    time.Now().Add(time.Hour),
			expectError: false,
			description: "Should create context with future deadline",
		},
		{
			name:        "past deadline",
			deadline:    time.Now().Add(-time.Hour),
			expectError: false,
			description: "Should create context with past deadline (no auto-cancel)",
		},
		{
			name:        "zero deadline",
			deadline:    time.Time{},
			expectError: false,
			description: "Should create context with zero deadline",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.NewContextWithDeadline(tt.deadline)

			// Verify interface compliance
			var _ domainContext.Context = ctx

			assert.NotNil(t, ctx)

			// Check deadline
			deadline, hasDeadline := ctx.Deadline()
			if tt.deadline.IsZero() {
				assert.False(t, hasDeadline)
				assert.True(t, deadline.IsZero())
			} else {
				assert.True(t, hasDeadline)
				assert.Equal(t, tt.deadline, deadline)
			}
		})
	}
}

func TestNewContextWithTimeout(t *testing.T) {
	tests := []struct {
		name            string
		timeout         time.Duration
		expectCancelled bool
		description     string
	}{
		{
			name:            "positive timeout",
			timeout:         time.Hour,
			expectCancelled: false,
			description:     "Should create context with positive timeout",
		},
		{
			name:            "zero timeout",
			timeout:         0,
			expectCancelled: true,
			description:     "Should create immediately cancelled context for zero timeout",
		},
		{
			name:            "negative timeout",
			timeout:         -time.Hour,
			expectCancelled: true,
			description:     "Should create immediately cancelled context for negative timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.NewContextWithTimeout(tt.timeout)

			// Verify interface compliance
			var _ domainContext.Context = ctx

			assert.NotNil(t, ctx)

			if tt.expectCancelled {
				assert.True(t, ctx.IsCancelled())
				assert.NotNil(t, ctx.Err())
				assert.Contains(t, ctx.Err().Error(), domainContext.ErrContextDeadlineExceeded)
			} else {
				assert.False(t, ctx.IsCancelled())
				assert.Nil(t, ctx.Err())

				// Check deadline is set correctly
				deadline, hasDeadline := ctx.Deadline()
				assert.True(t, hasDeadline)
				assert.True(t, deadline.After(time.Now()))
			}
		})
	}
}

func TestContextValue(t *testing.T) {
	tests := []struct {
		name        string
		key         interface{}
		value       interface{}
		lookupKey   interface{}
		expectFound bool
		description string
	}{
		{
			name:        "string key and value",
			key:         "test-key",
			value:       "test-value",
			lookupKey:   "test-key",
			expectFound: true,
			description: "Should store and retrieve string values",
		},
		{
			name:        "int key and value",
			key:         42,
			value:       "forty-two",
			lookupKey:   42,
			expectFound: true,
			description: "Should store and retrieve with int keys",
		},
		{
			name:        "struct key and value",
			key:         struct{ name string }{"test"},
			value:       "struct-value",
			lookupKey:   struct{ name string }{"test"},
			expectFound: true,
			description: "Should store and retrieve with struct keys",
		},
		{
			name:        "nil key",
			key:         nil,
			value:       "nil-key-value",
			lookupKey:   nil,
			expectFound: false,
			description: "Should handle nil keys gracefully",
		},
		{
			name:        "key not found",
			key:         "existing-key",
			value:       "existing-value",
			lookupKey:   "non-existent-key",
			expectFound: false,
			description: "Should return nil for non-existent keys",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.NewContext()

			// Set value if key is not nil
			if tt.key != nil {
				newCtx := ctx.WithValue(tt.key, tt.value)
				ctx = newCtx.(*context.DomainContext)
			}

			// Lookup value
			result := ctx.Value(tt.lookupKey)

			if tt.expectFound {
				assert.Equal(t, tt.value, result)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestContextWithValue(t *testing.T) {
	tests := []struct {
		name        string
		key         interface{}
		value       interface{}
		expectSame  bool
		description string
	}{
		{
			name:        "valid key and value",
			key:         "test-key",
			value:       "test-value",
			expectSame:  false,
			description: "Should create new context with value",
		},
		{
			name:        "nil key",
			key:         nil,
			value:       "test-value",
			expectSame:  true,
			description: "Should return same context for nil key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalCtx := context.NewContext()
			newCtx := originalCtx.WithValue(tt.key, tt.value)

			// Verify interface compliance
			var _ domainContext.Context = newCtx

			if tt.expectSame {
				assert.Same(t, originalCtx, newCtx)
			} else {
				assert.NotSame(t, originalCtx, newCtx)

				// Verify value is accessible in new context
				if tt.key != nil {
					assert.Equal(t, tt.value, newCtx.Value(tt.key))
				}
			}
		})
	}
}

func TestContextWithValueInheritance(t *testing.T) {
	// Create parent context with multiple values
	parentCtx := context.NewContext()
	ctx1 := parentCtx.WithValue("key1", "value1")
	ctx2 := ctx1.WithValue("key2", "value2")
	ctx3 := ctx2.WithValue("key3", "value3")

	// Verify all values are accessible
	assert.Equal(t, "value1", ctx3.Value("key1"))
	assert.Equal(t, "value2", ctx3.Value("key2"))
	assert.Equal(t, "value3", ctx3.Value("key3"))

	// Verify parent contexts still have their values
	assert.Equal(t, "value1", ctx1.Value("key1"))
	assert.Nil(t, ctx1.Value("key2"))
	assert.Nil(t, ctx1.Value("key3"))

	assert.Equal(t, "value1", ctx2.Value("key1"))
	assert.Equal(t, "value2", ctx2.Value("key2"))
	assert.Nil(t, ctx2.Value("key3"))
}

func TestContextDeadline(t *testing.T) {
	// Test context without deadline
	ctx := context.NewContext()
	deadline, hasDeadline := ctx.Deadline()
	assert.False(t, hasDeadline)
	assert.True(t, deadline.IsZero())

	// Test context with deadline
	expectedDeadline := time.Now().Add(time.Hour)
	ctxWithDeadline := context.NewContextWithDeadline(expectedDeadline)
	deadline, hasDeadline = ctxWithDeadline.Deadline()
	assert.True(t, hasDeadline)
	assert.Equal(t, expectedDeadline, deadline)
}

func TestContextDone(t *testing.T) {
	ctx := context.NewContext()

	// Verify done channel is not closed initially
	select {
	case <-ctx.Done():
		t.Fatal("Done channel should not be closed initially")
	default:
		// Expected
	}

	// Cancel context
	ctx.Cancel()

	// Verify done channel is closed after cancellation
	select {
	case <-ctx.Done():
		// Expected
	default:
		t.Fatal("Done channel should be closed after cancellation")
	}
}

func TestContextErr(t *testing.T) {
	ctx := context.NewContext()

	// Initially no error
	assert.Nil(t, ctx.Err())

	// Cancel context
	ctx.Cancel()

	// Should have cancellation error
	assert.NotNil(t, ctx.Err())
	assert.Contains(t, ctx.Err().Error(), domainContext.ErrContextCanceled)
}

func TestContextCancel(t *testing.T) {
	ctx := context.NewContext()

	// Initially not cancelled
	assert.False(t, ctx.IsCancelled())
	assert.Nil(t, ctx.Err())

	// Cancel context
	ctx.Cancel()

	// Should be cancelled
	assert.True(t, ctx.IsCancelled())
	assert.NotNil(t, ctx.Err())
	assert.Contains(t, ctx.Err().Error(), domainContext.ErrContextCanceled)

	// Multiple cancellations should be safe
	ctx.Cancel()
	assert.True(t, ctx.IsCancelled())
}

func TestContextDeadlineExpiration(t *testing.T) {
	// Create context with very short timeout
	ctx := context.NewContextWithTimeout(50 * time.Millisecond)

	// Initially not cancelled
	assert.False(t, ctx.IsCancelled())
	assert.Nil(t, ctx.Err())

	// Wait for deadline to expire
	time.Sleep(100 * time.Millisecond)

	// Should be cancelled due to deadline
	assert.True(t, ctx.IsCancelled())
	assert.NotNil(t, ctx.Err())
	assert.Contains(t, ctx.Err().Error(), domainContext.ErrContextDeadlineExceeded)
}

func TestContextWithValueFromCancelledParent(t *testing.T) {
	// Create and cancel parent context
	parentCtx := context.NewContext()
	parentCtx.Cancel()

	// Create child context from cancelled parent
	childCtx := parentCtx.WithValue("key", "value")

	// Child should also be cancelled
	assert.True(t, childCtx.(*context.DomainContext).IsCancelled())

	// But should still have the value
	assert.Equal(t, "value", childCtx.Value("key"))
}

func TestWrapContext(t *testing.T) {
	// Test wrapping with nil
	ctx := context.WrapContext(nil)
	assert.NotNil(t, ctx)

	// Verify interface compliance
	var _ domainContext.Context = ctx

	// Test wrapping with some value
	ctx2 := context.WrapContext("some-context")
	assert.NotNil(t, ctx2)

	// Verify interface compliance
	var _ domainContext.Context = ctx2
}

func TestContextThreadSafety(t *testing.T) {
	ctx := context.NewContext()

	// Test concurrent value access
	done := make(chan bool)

	// Writer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			newCtx := ctx.WithValue(i, i*2)
			ctx = newCtx.(*context.DomainContext)
		}
		done <- true
	}()

	// Reader goroutine
	go func() {
		for i := 0; i < 100; i++ {
			_ = ctx.Value(i)
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done

	// Test should complete without race conditions
}

func TestContextInterfaceCompliance(t *testing.T) {
	// Test all constructor functions return interface-compliant types
	var _ domainContext.Context = context.NewContext()
	var _ domainContext.Context = context.NewContextWithDeadline(time.Now().Add(time.Hour))
	var _ domainContext.Context = context.NewContextWithTimeout(time.Hour)
	var _ domainContext.Context = context.WrapContext(nil)
}

// Benchmark tests for performance verification
func BenchmarkContextValue(b *testing.B) {
	ctx := context.NewContext()
	ctx = ctx.WithValue("key", "value").(*context.DomainContext)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx.Value("key")
	}
}

func BenchmarkContextWithValue(b *testing.B) {
	ctx := context.NewContext()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx = ctx.WithValue(i, i).(*context.DomainContext)
	}
}

func BenchmarkContextDeadline(b *testing.B) {
	ctx := context.NewContextWithDeadline(time.Now().Add(time.Hour))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = ctx.Deadline()
	}
}
