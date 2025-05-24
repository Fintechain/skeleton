package operation

import (
	"github.com/ebanfa/skeleton/internal/domain/component"
)

// Pipeline represents a chain of operations that execute in sequence.
type Pipeline interface {
	Operation

	// AddOperation adds an operation to the pipeline.
	AddOperation(operation Operation)

	// Operations returns the operations in the pipeline.
	Operations() []Operation
}

// DefaultPipelineOptions contains options for creating a DefaultPipeline.
type DefaultPipelineOptions struct {
	BaseOperation Operation
	Operations    []Operation
}

// DefaultPipeline is a basic implementation of the Pipeline interface.
type DefaultPipeline struct {
	Operation
	operations []Operation
}

// NewPipeline creates a new operation pipeline with dependency injection.
func NewPipeline(options DefaultPipelineOptions) *DefaultPipeline {
	// Initialize the operations slice based on input or create a new one
	operations := options.Operations
	if operations == nil {
		operations = make([]Operation, 0)
	}

	return &DefaultPipeline{
		Operation:  options.BaseOperation,
		operations: operations,
	}
}

// CreatePipeline is a factory method for backward compatibility.
func CreatePipeline(baseOperation Operation) *DefaultPipeline {
	return NewPipeline(DefaultPipelineOptions{
		BaseOperation: baseOperation,
		Operations:    make([]Operation, 0),
	})
}

// AddOperation adds an operation to the pipeline.
func (p *DefaultPipeline) AddOperation(operation Operation) {
	p.operations = append(p.operations, operation)
}

// Operations returns the operations in the pipeline.
func (p *DefaultPipeline) Operations() []Operation {
	return p.operations
}

// Execute executes the pipeline by running each operation in sequence.
func (p *DefaultPipeline) Execute(ctx component.Context, input Input) (Output, error) {
	var currentInput Input = input
	var currentOutput Output
	var err error

	// Execute each operation in sequence, passing the output of one as input to the next
	for _, op := range p.operations {
		currentOutput, err = op.Execute(ctx, currentInput)
		if err != nil {
			return nil, component.NewError(
				ErrOperationExecution,
				"pipeline operation failed",
				err,
			).WithDetail("operation_id", op.ID())
		}

		// Use the output as input for the next operation
		currentInput = currentOutput
	}

	return currentOutput, nil
}

// PipelineBuilderOptions contains options for creating a PipelineBuilder.
type PipelineBuilderOptions struct {
	BaseOperation Operation
	Operations    []Operation
}

// PipelineBuilder helps construct pipelines.
type PipelineBuilder struct {
	baseOperation Operation
	operations    []Operation
}

// NewPipelineBuilder creates a new pipeline builder with dependency injection.
func NewPipelineBuilder(options PipelineBuilderOptions) *PipelineBuilder {
	// Initialize the operations slice based on input or create a new one
	operations := options.Operations
	if operations == nil {
		operations = make([]Operation, 0)
	}

	return &PipelineBuilder{
		baseOperation: options.BaseOperation,
		operations:    operations,
	}
}

// CreatePipelineBuilder is a factory method for backward compatibility.
func CreatePipelineBuilder(baseOperation Operation) *PipelineBuilder {
	return NewPipelineBuilder(PipelineBuilderOptions{
		BaseOperation: baseOperation,
		Operations:    make([]Operation, 0),
	})
}

// AddOperation adds an operation to the pipeline.
func (b *PipelineBuilder) AddOperation(operation Operation) *PipelineBuilder {
	b.operations = append(b.operations, operation)
	return b
}

// Build constructs the pipeline.
func (b *PipelineBuilder) Build() Pipeline {
	pipeline := NewPipeline(DefaultPipelineOptions{
		BaseOperation: b.baseOperation,
		Operations:    b.operations,
	})
	return pipeline
}
