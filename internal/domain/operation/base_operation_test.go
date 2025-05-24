package operation

import (
	"strings"
	"testing"
	"time"

	"github.com/ebanfa/skeleton/internal/domain/component"
)

// testComponent implements component.Component for testing
type testComponent struct {
	id   string
	name string
}

func (c *testComponent) ID() string                         { return c.id }
func (c *testComponent) Name() string                       { return c.name }
func (c *testComponent) Type() component.ComponentType      { return component.TypeOperation }
func (c *testComponent) Metadata() component.Metadata       { return component.Metadata{} }
func (c *testComponent) Initialize(component.Context) error { return nil }
func (c *testComponent) Dispose() error                     { return nil }

// testContext implements component.Context for testing
type testContext struct {
	doneChannel  chan struct{}
	contextError error
}

func (c *testContext) Value(interface{}) interface{}           { return nil }
func (c *testContext) Deadline() (deadline time.Time, ok bool) { return time.Time{}, false }
func (c *testContext) Done() <-chan struct{} {
	if c.doneChannel == nil {
		return nil
	}
	return c.doneChannel
}
func (c *testContext) Err() error { return c.contextError }
func (c *testContext) WithValue(key, value interface{}) component.Context {
	return c // For tests, we just return the same context
}

func TestBaseOperation_New(t *testing.T) {
	// Create a test component
	testComp := &testComponent{
		id:   "test-id",
		name: "Test Component",
	}

	// Create a base operation with the test component using the options struct
	baseOp := NewBaseOperation(BaseOperationOptions{
		Component: testComp,
	})

	// Verify the component is set correctly
	if baseOp.ID() != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", baseOp.ID())
	}

	if baseOp.Name() != "Test Component" {
		t.Errorf("Expected Name 'Test Component', got '%s'", baseOp.Name())
	}
}

func TestBaseOperation_Create(t *testing.T) {
	// Create a test component
	testComp := &testComponent{
		id:   "test-id",
		name: "Test Component",
	}

	// Test the factory method for backward compatibility
	baseOp := CreateBaseOperation(testComp)

	// Verify the component is set correctly
	if baseOp.ID() != "test-id" {
		t.Errorf("Expected ID 'test-id', got '%s'", baseOp.ID())
	}
}

func TestBaseOperation_Execute(t *testing.T) {
	// Create a test component
	testComp := &testComponent{
		id:   "test-id",
		name: "Test Component",
	}

	// Create a base operation
	baseOp := NewBaseOperation(BaseOperationOptions{
		Component: testComp,
	})

	// Execute should return an error for base operation
	ctx := &testContext{}
	_, err := baseOp.Execute(ctx, "test-input")

	// Base operation should return an error
	if err == nil {
		t.Error("Expected error from base operation Execute, got nil")
	}

	// Check if the error message contains the expected error code
	errMsg := err.Error()
	if !strings.Contains(errMsg, ErrOperationExecution) {
		t.Errorf("Expected error message to contain '%s', got '%s'",
			ErrOperationExecution, errMsg)
	}

	// The base implementation error message contains "base operation does not implement Execute"
	// rather than the operation ID
	if !strings.Contains(errMsg, "base operation does not implement Execute") {
		t.Errorf("Expected error message to contain 'base operation does not implement Execute', got '%s'", errMsg)
	}
}
