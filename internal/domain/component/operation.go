// Package operation provides functionality for operations in the system.
package component

import (
	"github.com/fintechain/skeleton/internal/domain/context"
)

// Input represents the input for an operation.
type Input struct {
	Data any
	// Metadata is additional information about the input.
	Metadata map[string]string
}

// Output represents the output from an operation.
type Output struct {
	Data any
}

// Operation is an executable instruction that runs within the system context.
// Operations are stateless, discrete tasks that execute once and return results.
type Operation interface {
	Component

	// Execute runs the operation with the given context and input.
	Execute(ctx context.Context, input Input) (Output, error)
}
