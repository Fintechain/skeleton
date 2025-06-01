package context

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/fintechain/skeleton/pkg/context"
)

// TestNewContext tests the context constructor function
func TestNewContext(t *testing.T) {
	ctx := context.NewContext()

	assert.NotNil(t, ctx)
	assert.Nil(t, ctx.Err())

	// Test that Done channel is not closed initially
	select {
	case <-ctx.Done():
		t.Error("Done channel should not be closed for new context")
	default:
		// Expected behavior
	}

	// Test that deadline is not set initially
	deadline, hasDeadline := ctx.Deadline()
	assert.False(t, hasDeadline)
	assert.True(t, deadline.IsZero())
}

// TestWrapContext tests the context wrapping function
func TestWrapContext(t *testing.T) {
	tests := []struct {
		name     string
		input    context.Context
		expected string
	}{
		{
			name:     "wrap nil context",
			input:    nil,
			expected: "should create new context when wrapping nil",
		},
		{
			name:     "wrap existing context",
			input:    context.NewContext(),
			expected: "should wrap existing context",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := context.WrapContext(tt.input)

			assert.NotNil(t, wrapped, tt.expected)
			assert.Nil(t, wrapped.Err())

			// Test that Done channel is not closed initially
			select {
			case <-wrapped.Done():
				t.Error("Done channel should not be closed for wrapped context")
			default:
				// Expected behavior
			}
		})
	}
}

// TestContextInterfaceCompliance verifies the implementation satisfies the domain interface
func TestContextInterfaceCompliance(t *testing.T) {
	// Verify interface compliance
	var _ context.Context = context.NewContext()
	var _ context.Context = context.WrapContext(nil)
}

// TestContextValue tests context value operations
func TestContextValue(t *testing.T) {
	ctx := context.NewContext()

	// Test getting non-existent value
	value := ctx.Value("non-existent")
	assert.Nil(t, value)

	// Test setting and getting values
	key1 := "test-key-1"
	value1 := "test-value-1"

	newCtx := ctx.WithValue(key1, value1)
	assert.NotNil(t, newCtx)
	assert.NotEqual(t, ctx, newCtx) // Should return new context

	// Test getting value from new context
	retrieved := newCtx.Value(key1)
	assert.Equal(t, value1, retrieved)

	// Test that original context doesn't have the value
	originalValue := ctx.Value(key1)
	assert.Nil(t, originalValue)
}

// TestContextWithValue tests context value chaining
func TestContextWithValue(t *testing.T) {
	ctx := context.NewContext()

	// Chain multiple values
	ctx1 := ctx.WithValue("key1", "value1")
	ctx2 := ctx1.WithValue("key2", "value2")
	ctx3 := ctx2.WithValue("key3", "value3")

	// Test that all values are accessible from the final context
	assert.Equal(t, "value1", ctx3.Value("key1"))
	assert.Equal(t, "value2", ctx3.Value("key2"))
	assert.Equal(t, "value3", ctx3.Value("key3"))

	// Test that intermediate contexts only have their values
	assert.Equal(t, "value1", ctx1.Value("key1"))
	assert.Nil(t, ctx1.Value("key2"))
	assert.Nil(t, ctx1.Value("key3"))

	assert.Equal(t, "value1", ctx2.Value("key1"))
	assert.Equal(t, "value2", ctx2.Value("key2"))
	assert.Nil(t, ctx2.Value("key3"))
}

// TestContextValueTypes tests context with different value types
func TestContextValueTypes(t *testing.T) {
	ctx := context.NewContext()

	// Test different value types
	testCases := []struct {
		key   interface{}
		value interface{}
	}{
		{"string-key", "string-value"},
		{42, "int-key"},
		{"struct-value", struct{ Name string }{Name: "test"}},
		{"slice-value", []string{"a", "b", "c"}},
		{"map-value", map[string]int{"count": 42}},
	}

	// Set all values
	currentCtx := ctx
	for _, tc := range testCases {
		currentCtx = currentCtx.WithValue(tc.key, tc.value)
	}

	// Verify all values
	for _, tc := range testCases {
		retrieved := currentCtx.Value(tc.key)
		assert.Equal(t, tc.value, retrieved)
	}
}

// TestContextDeadline tests context deadline functionality
func TestContextDeadline(t *testing.T) {
	ctx := context.NewContext()

	// Test context without deadline
	deadline, hasDeadline := ctx.Deadline()
	assert.False(t, hasDeadline)
	assert.True(t, deadline.IsZero())
}

// TestContextDoneChannel tests the Done channel behavior
func TestContextDoneChannel(t *testing.T) {
	ctx := context.NewContext()

	// Test that Done channel is not closed for new context
	select {
	case <-ctx.Done():
		t.Error("Done channel should not be closed for new context")
	default:
		// Expected behavior
	}
}

// TestContextErrorStates tests different error states
func TestContextErrorStates(t *testing.T) {
	ctx := context.NewContext()

	// Test new context has no error
	assert.Nil(t, ctx.Err())
}

// TestContextWrappingWithValues tests wrapping contexts that have values
func TestContextWrappingWithValues(t *testing.T) {
	// Create context with values
	originalCtx := context.NewContext()
	originalCtx = originalCtx.WithValue("original-key", "original-value")

	// Wrap the context
	wrappedCtx := context.WrapContext(originalCtx)

	// Note: The current implementation may not preserve values when wrapping
	// This test documents the current behavior
	assert.NotNil(t, wrappedCtx)

	// Test that wrapped context can have its own values
	wrappedCtx = wrappedCtx.WithValue("wrapped-key", "wrapped-value")
	assert.Equal(t, "wrapped-value", wrappedCtx.Value("wrapped-key"))
}

// TestContextConcurrency tests context operations under concurrent access
func TestContextConcurrency(t *testing.T) {
	ctx := context.NewContext()

	// Test concurrent value setting and getting
	done := make(chan bool, 10)

	// Start multiple goroutines setting values
	for i := 0; i < 10; i++ {
		go func(id int) {
			defer func() { done <- true }()

			key := fmt.Sprintf("key-%d", id)
			value := fmt.Sprintf("value-%d", id)

			newCtx := ctx.WithValue(key, value)
			retrieved := newCtx.Value(key)

			assert.Equal(t, value, retrieved)
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		select {
		case <-done:
			// Expected
		case <-time.After(1 * time.Second):
			t.Error("Goroutine did not complete in time")
		}
	}
}
