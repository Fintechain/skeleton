package operation

import (
	"errors"
	"testing"

	"github.com/ebanfa/skeleton/internal/domain/component"
)

// testOperation is a simple operation implementation for testing pipelines
type testOperation struct {
	component.Component
	executeFunc func(ctx component.Context, input Input) (Output, error)
}

func newTestOperation(comp component.Component, executeFunc func(ctx component.Context, input Input) (Output, error)) *testOperation {
	return &testOperation{
		Component:   comp,
		executeFunc: executeFunc,
	}
}

func (m *testOperation) Execute(ctx component.Context, input Input) (Output, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, input)
	}
	return nil, nil
}

func TestNewPipeline(t *testing.T) {
	// Create a test component
	testComp := &testComponent{
		id:   "base-op-id",
		name: "Test Component",
	}

	// Create a base operation
	baseOp := CreateDefaultOperation(testComp)

	// Test with empty operations
	pipeline := NewPipeline(DefaultPipelineOptions{
		BaseOperation: baseOp,
	})

	if len(pipeline.Operations()) != 0 {
		t.Errorf("Expected 0 operations, got %d", len(pipeline.Operations()))
	}

	// Test with provided operations
	op1 := newTestOperation(testComp, nil)
	op2 := newTestOperation(testComp, nil)
	operations := []Operation{op1, op2}

	pipeline = NewPipeline(DefaultPipelineOptions{
		BaseOperation: baseOp,
		Operations:    operations,
	})

	if len(pipeline.Operations()) != 2 {
		t.Errorf("Expected 2 operations, got %d", len(pipeline.Operations()))
	}
}

func TestCreatePipeline(t *testing.T) {
	// Create a test component
	testComp := &testComponent{
		id:   "base-op-id",
		name: "Test Component",
	}

	// Create a base operation
	baseOp := CreateDefaultOperation(testComp)

	// Test the factory method for backward compatibility
	pipeline := CreatePipeline(baseOp)

	if pipeline.ID() != "base-op-id" {
		t.Errorf("Expected ID 'base-op-id', got '%s'", pipeline.ID())
	}

	if len(pipeline.Operations()) != 0 {
		t.Errorf("Expected 0 operations, got %d", len(pipeline.Operations()))
	}
}

func TestDefaultPipeline_AddOperation(t *testing.T) {
	// Create a test component
	testComp := &testComponent{
		id:   "test-id",
		name: "Test Component",
	}

	// Create a base operation
	baseOp := CreateDefaultOperation(testComp)

	// Create a pipeline
	pipeline := NewPipeline(DefaultPipelineOptions{
		BaseOperation: baseOp,
	})

	// Create operations to add
	op1 := newTestOperation(testComp, nil)
	op2 := newTestOperation(testComp, nil)

	// Add operations
	pipeline.AddOperation(op1)
	pipeline.AddOperation(op2)

	// Check operations were added correctly
	ops := pipeline.Operations()
	if len(ops) != 2 {
		t.Errorf("Expected 2 operations, got %d", len(ops))
	}
}

func TestDefaultPipeline_Execute(t *testing.T) {
	// Create test cases
	tests := []struct {
		name          string
		setupPipeline func() Pipeline
		input         Input
		expected      Output
		expectError   bool
	}{
		{
			name: "Empty pipeline",
			setupPipeline: func() Pipeline {
				testComp := &testComponent{
					id:   "test-id",
					name: "Test Component",
				}
				baseOp := CreateDefaultOperation(testComp)
				return CreatePipeline(baseOp)
			},
			input:       "test-input",
			expected:    nil,
			expectError: false,
		},
		{
			name: "Successful pipeline execution",
			setupPipeline: func() Pipeline {
				testComp := &testComponent{
					id:   "test-id",
					name: "Test Component",
				}
				baseOp := CreateDefaultOperation(testComp)
				pipeline := CreatePipeline(baseOp)

				op1 := newTestOperation(testComp, func(ctx component.Context, input Input) (Output, error) {
					return "op1:" + input.(string), nil
				})

				op2 := newTestOperation(testComp, func(ctx component.Context, input Input) (Output, error) {
					return input.(string) + ":op2", nil
				})

				pipeline.AddOperation(op1)
				pipeline.AddOperation(op2)

				return pipeline
			},
			input:       "test",
			expected:    "op1:test:op2",
			expectError: false,
		},
		{
			name: "Operation error",
			setupPipeline: func() Pipeline {
				testComp := &testComponent{
					id:   "error-op",
					name: "Error Component",
				}
				baseOp := CreateDefaultOperation(testComp)
				pipeline := CreatePipeline(baseOp)

				op1 := newTestOperation(testComp, func(ctx component.Context, input Input) (Output, error) {
					return "op1:" + input.(string), nil
				})

				op2 := newTestOperation(testComp, func(ctx component.Context, input Input) (Output, error) {
					return nil, errors.New("operation error")
				})

				pipeline.AddOperation(op1)
				pipeline.AddOperation(op2)

				return pipeline
			},
			input:       "test",
			expected:    nil,
			expectError: true,
		},
	}

	// Run the tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pipeline := tt.setupPipeline()
			ctx := &testContext{}

			output, err := pipeline.Execute(ctx, tt.input)

			if tt.expectError && err == nil {
				t.Error("Expected error, got nil")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if !tt.expectError && output != tt.expected {
				t.Errorf("Expected output '%v', got '%v'", tt.expected, output)
			}
		})
	}
}

func TestPipelineBuilder(t *testing.T) {
	// Create a test component
	testComp := &testComponent{
		id:   "base-op-id",
		name: "Test Component",
	}

	// Create a base operation
	baseOp := CreateDefaultOperation(testComp)

	// Test different constructors
	// 1. Test NewPipelineBuilder with empty operations
	t.Run("NewPipelineBuilder empty", func(t *testing.T) {
		builder := NewPipelineBuilder(PipelineBuilderOptions{
			BaseOperation: baseOp,
		})

		pipeline := builder.Build()
		if len(pipeline.Operations()) != 0 {
			t.Errorf("Expected 0 operations, got %d", len(pipeline.Operations()))
		}
	})

	// 2. Test NewPipelineBuilder with provided operations
	t.Run("NewPipelineBuilder with operations", func(t *testing.T) {
		op1 := newTestOperation(testComp, nil)
		op2 := newTestOperation(testComp, nil)
		operations := []Operation{op1, op2}

		builder := NewPipelineBuilder(PipelineBuilderOptions{
			BaseOperation: baseOp,
			Operations:    operations,
		})

		pipeline := builder.Build()
		if len(pipeline.Operations()) != 2 {
			t.Errorf("Expected 2 operations, got %d", len(pipeline.Operations()))
		}
	})

	// 3. Test CreatePipelineBuilder
	t.Run("CreatePipelineBuilder", func(t *testing.T) {
		builder := CreatePipelineBuilder(baseOp)

		pipeline := builder.Build()
		if len(pipeline.Operations()) != 0 {
			t.Errorf("Expected 0 operations, got %d", len(pipeline.Operations()))
		}

		if pipeline.ID() != "base-op-id" {
			t.Errorf("Expected ID 'base-op-id', got '%s'", pipeline.ID())
		}
	})

	// 4. Test AddOperation chaining
	t.Run("AddOperation chaining", func(t *testing.T) {
		builder := CreatePipelineBuilder(baseOp)

		op1 := newTestOperation(testComp, nil)
		op2 := newTestOperation(testComp, nil)

		// Add operations with chaining
		result := builder.AddOperation(op1).AddOperation(op2)

		// Build and verify
		pipeline := result.Build()
		if len(pipeline.Operations()) != 2 {
			t.Errorf("Expected 2 operations, got %d", len(pipeline.Operations()))
		}
	})

	// 5. Test full pipeline building and execution
	t.Run("Pipeline building and execution", func(t *testing.T) {
		builder := CreatePipelineBuilder(baseOp)

		op1 := newTestOperation(testComp, func(ctx component.Context, input Input) (Output, error) {
			return "op1:" + input.(string), nil
		})

		op2 := newTestOperation(testComp, func(ctx component.Context, input Input) (Output, error) {
			return input.(string) + ":op2", nil
		})

		// Add operations and build
		pipeline := builder.AddOperation(op1).AddOperation(op2).Build()

		// Test execution
		ctx := &testContext{}
		output, err := pipeline.Execute(ctx, "test")

		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		expected := "op1:test:op2"
		if output != expected {
			t.Errorf("Expected output '%v', got '%v'", expected, output)
		}
	})
}
