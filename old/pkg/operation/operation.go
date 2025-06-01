// Package operation provides public APIs for the operation system.
package operation

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/operation"
)

// ===== OPERATION INTERFACES =====

// Operation is a specialized component that executes discrete units of work.
type Operation = operation.Operation

// OperationFactory creates operations from configuration.
type OperationFactory = operation.OperationFactory

// ===== OPERATION TYPES =====

// Input represents the input for an operation.
type Input = operation.Input

// Output represents the output from an operation.
type Output = operation.Output

// OperationConfig defines the configuration for creating an operation.
type OperationConfig = operation.OperationConfig

// ===== OPERATION ERROR CONSTANTS =====

// Common operation error codes
const (
	ErrOperationExecution = "operation.execution_failed"
	ErrInvalidInput       = "operation.invalid_input"
	ErrOperationTimeout   = "operation.timeout"
	ErrOperationNotFound  = "operation.not_found"
	ErrInvalidOutput      = "operation.invalid_output"
)

// ===== ERROR HANDLING =====

// Error represents a domain-specific error from the operation system.
type Error = component.Error

// NewError creates a new operation error with the given code, message, and optional cause.
func NewError(code, message string, cause error) *Error {
	return component.NewError(code, message, cause)
}

// IsOperationError checks if an error is an operation error with the given code.
func IsOperationError(err error, code string) bool {
	return component.IsComponentError(err, code)
}

// ===== OPERATION CONSTRUCTORS =====

// NewOperationFactory creates a new operation factory with default configuration.
// This is the primary way to create an OperationFactory instance for creating operations.
//
// Example usage:
//
//	factory := operation.NewOperationFactory()
//	config := operation.OperationConfig{
//	    ComponentConfig: component.ComponentConfig{
//	        Name: "data-processor",
//	        Type: "batch",
//	    },
//	}
//	op, err := factory.CreateOperation(config)
func NewOperationFactory() OperationFactory {
	return operation.NewOperationFactory()
}

// NewOperation creates a new operation with the given name and type.
// This is a convenience function for creating simple operations without complex configuration.
//
// Example usage:
//
//	op, err := operation.NewOperation("data-processor", "batch")
//	if err != nil {
//	    // Handle error
//	}
func NewOperation(name, operationType string) (Operation, error) {
	if name == "" {
		return nil, NewError(ErrInvalidInput, "operation name cannot be empty", nil)
	}
	if operationType == "" {
		return nil, NewError(ErrInvalidInput, "operation type cannot be empty", nil)
	}

	factory := NewOperationFactory()
	config := OperationConfig{
		ComponentConfig: component.ComponentConfig{
			Name: name,
			Type: component.ComponentType(operationType),
		},
	}

	return factory.CreateOperation(config)
}

// ===== OPERATION UTILITIES =====

// Execute is a convenience function to execute an operation with input validation.
// It provides additional error handling and input validation on top of the basic Execute method.
//
// Example usage:
//
//	result, err := operation.Execute(ctx, op, inputData)
//	if err != nil {
//	    // Handle execution error
//	}
func Execute(ctx component.Context, op Operation, input Input) (Output, error) {
	if op == nil {
		return nil, NewError(ErrInvalidInput, "operation cannot be nil", nil)
	}
	if ctx == nil {
		return nil, NewError(ErrInvalidInput, "context cannot be nil", nil)
	}

	// Execute the operation
	output, err := op.Execute(ctx, input)
	if err != nil {
		return nil, NewError(ErrOperationExecution, "operation execution failed", err).
			WithDetail("operationId", op.ID()).
			WithDetail("operationType", string(op.Type()))
	}

	return output, nil
}

// ValidateInput performs basic validation on operation input.
// This is a utility function that can be used by operation implementations.
//
// Example usage:
//
//	if err := operation.ValidateInput(input, "data"); err != nil {
//	    return nil, err
//	}
func ValidateInput(input Input, requiredFields ...string) error {
	if input == nil {
		return NewError(ErrInvalidInput, "input cannot be nil", nil)
	}

	// If input is a map, check for required fields
	if inputMap, ok := input.(map[string]interface{}); ok {
		for _, field := range requiredFields {
			if _, exists := inputMap[field]; !exists {
				return NewError(ErrInvalidInput, "required field missing", nil).
					WithDetail("field", field)
			}
		}
	}

	return nil
}

// CreateOutput creates a standardized output structure for operations.
// This is a utility function to help create consistent operation outputs.
//
// Example usage:
//
//	output := operation.CreateOutput(map[string]interface{}{
//	    "result": processedData,
//	    "count": len(processedData),
//	})
func CreateOutput(data interface{}) Output {
	return map[string]interface{}{
		"data":    data,
		"success": true,
	}
}

// CreateErrorOutput creates a standardized error output structure for operations.
// This is a utility function to help create consistent error outputs.
//
// Example usage:
//
//	output := operation.CreateErrorOutput("validation failed", err)
func CreateErrorOutput(message string, cause error) Output {
	output := map[string]interface{}{
		"success": false,
		"error":   message,
	}

	if cause != nil {
		output["cause"] = cause.Error()
	}

	return output
}
