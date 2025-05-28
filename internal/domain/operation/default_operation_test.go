package operation

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/fintechain/skeleton/internal/domain/component"
)

func TestDefaultOperation_New(t *testing.T) {
	// Create a test component
	testComp := &testComponent{
		id:   "test-id",
		name: "Test Component",
	}

	// Test with no execute function
	op := NewDefaultOperation(DefaultOperationOptions{
		Component: testComp,
	})

	// Verify the component is set correctly
	if op.ID() != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", op.ID())
	}

	if op.Name() != "Test Component" {
		t.Errorf("Expected Name 'Test Component', got '%s'", op.Name())
	}

	// Test with execute function
	executed := false
	op = NewDefaultOperation(DefaultOperationOptions{
		Component: testComp,
		ExecuteFunc: func(ctx component.Context, input Input) (Output, error) {
			executed = true
			return "result", nil
		},
	})

	// Execute the operation
	ctx := &testContext{}
	output, err := op.Execute(ctx, "input")

	// Check the result
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output != "result" {
		t.Errorf("Expected output 'result', got '%v'", output)
	}
	if !executed {
		t.Error("Execute function was not called")
	}
}

func TestDefaultOperation_Create(t *testing.T) {
	// Create a test component
	testComp := &testComponent{
		id:   "test-id",
		name: "Test Component",
	}

	// Test the factory method for backward compatibility
	op := CreateDefaultOperation(testComp)

	// Verify the component is set correctly
	if op.ID() != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", op.ID())
	}

	// Execute should call BaseOperation.Execute and return an error
	ctx := &testContext{}
	_, err := op.Execute(ctx, "input")
	if err == nil {
		t.Error("Expected error from Execute, got nil")
	}
}

func TestDefaultOperation_WithExecuteFunc(t *testing.T) {
	// Create a test component
	testComp := &testComponent{
		id:   "test-id",
		name: "Test Component",
	}

	// Create a default operation
	op := NewDefaultOperation(DefaultOperationOptions{
		Component: testComp,
	})

	// Set an execute function
	executed := false
	op.WithExecuteFunc(func(ctx component.Context, input Input) (Output, error) {
		executed = true
		return "result", nil
	})

	// Execute the operation
	ctx := &testContext{}
	output, err := op.Execute(ctx, "input")

	// Check the result
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if output != "result" {
		t.Errorf("Expected output 'result', got '%v'", output)
	}
	if !executed {
		t.Error("Execute function was not called")
	}
}

func TestDefaultOperation_ExecuteFallback(t *testing.T) {
	// Create a test component
	testComp := &testComponent{
		id:   "test-id",
		name: "Test Component",
	}

	// Create an operation with no execute function
	op := NewDefaultOperation(DefaultOperationOptions{
		Component: testComp,
	})

	// Execute should call BaseOperation.Execute and return an error
	ctx := &testContext{}
	_, err := op.Execute(ctx, "input")
	if err == nil {
		t.Error("Expected error from Execute, got nil")
	}
}

func TestMapOperation(t *testing.T) {
	tests := []struct {
		name     string
		input    Input
		expected Output
		hasError bool
	}{
		{
			name:     "String input",
			input:    "test",
			expected: "MAPPED:test",
			hasError: false,
		},
		{
			name:     "Integer input",
			input:    42,
			expected: "MAPPED:42",
			hasError: false,
		},
		{
			name:     "Error case",
			input:    nil,
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the map operation
			mapOp := MapOperation("map-op", func(input Input) (Output, error) {
				if input == nil {
					return nil, errors.New("nil input")
				}

				// Convert input to string regardless of type
				var inputStr string
				switch v := input.(type) {
				case string:
					inputStr = v
				case int:
					inputStr = strconv.Itoa(v)
				default:
					inputStr = fmt.Sprintf("%v", v)
				}

				return "MAPPED:" + inputStr, nil
			})

			// Execute the operation
			ctx := &testContext{}
			output, err := mapOp.Execute(ctx, tt.input)

			// Check the result
			if tt.hasError && err == nil {
				t.Error("Expected error, got nil")
			}
			if !tt.hasError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if !tt.hasError && output != tt.expected {
				t.Errorf("Expected output '%v', got '%v'", tt.expected, output)
			}
		})
	}
}

func TestFilterOperation(t *testing.T) {
	tests := []struct {
		name     string
		input    Input
		predFunc func(input Input) bool
		expected Output
	}{
		{
			name:  "Pass filter",
			input: "pass",
			predFunc: func(input Input) bool {
				return input.(string) == "pass"
			},
			expected: "pass",
		},
		{
			name:  "Fail filter",
			input: "fail",
			predFunc: func(input Input) bool {
				return input.(string) == "pass"
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create the filter operation
			filterOp := FilterOperation("filter-op", tt.predFunc)

			// Execute the operation
			ctx := &testContext{}
			output, err := filterOp.Execute(ctx, tt.input)

			// Check the result
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if output != tt.expected {
				t.Errorf("Expected output '%v', got '%v'", tt.expected, output)
			}
		})
	}
}

func TestAsyncOperation(t *testing.T) {
	// Create a test component
	testComp := &testComponent{}

	// Test successful execution
	t.Run("Successful execution", func(t *testing.T) {
		innerOp := NewDefaultOperation(DefaultOperationOptions{
			Component: testComp,
			ExecuteFunc: func(ctx component.Context, input Input) (Output, error) {
				return "result", nil
			},
		})

		asyncOp := AsyncOperation("async-op", innerOp)

		ctx := &testContext{}

		output, err := asyncOp.Execute(ctx, "input")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if output != "result" {
			t.Errorf("Expected output 'result', got '%v'", output)
		}
	})

	// Test context cancellation
	t.Run("Context canceled", func(t *testing.T) {
		innerOp := NewDefaultOperation(DefaultOperationOptions{
			Component: testComp,
			ExecuteFunc: func(ctx component.Context, input Input) (Output, error) {
				// This should not be reached due to context cancellation
				time.Sleep(time.Hour)
				return "result", nil
			},
		})

		asyncOp := AsyncOperation("async-op", innerOp)

		// Create a done channel that's already closed
		doneCh := make(chan struct{})
		close(doneCh)

		ctx := &testContext{
			doneChannel:  doneCh,
			contextError: errors.New("context canceled"),
		}

		_, err := asyncOp.Execute(ctx, "input")

		if err == nil {
			t.Error("Expected error, got nil")
		}

		// Check for specific error message/code
		errMsg := err.Error()
		if !strings.Contains(errMsg, ErrOperationTimeout) {
			t.Errorf("Expected error message to contain '%s', got '%s'",
				ErrOperationTimeout, errMsg)
		}
	})

	// Test inner operation error
	t.Run("Inner operation error", func(t *testing.T) {
		innerOp := NewDefaultOperation(DefaultOperationOptions{
			Component: testComp,
			ExecuteFunc: func(ctx component.Context, input Input) (Output, error) {
				return nil, errors.New("inner error")
			},
		})

		asyncOp := AsyncOperation("async-op", innerOp)

		ctx := &testContext{}

		_, err := asyncOp.Execute(ctx, "input")

		if err == nil {
			t.Error("Expected error, got nil")
		}
		if err.Error() != "inner error" {
			t.Errorf("Expected error 'inner error', got '%v'", err.Error())
		}
	})
}
