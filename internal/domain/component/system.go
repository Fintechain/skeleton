// Package component provides interfaces and types for the component system.
package component

import (
	"github.com/fintechain/skeleton/internal/domain/context"
)

type System interface {

	// Registry returns the component registry.
	Registry() Registry

	// ExecuteOperation executes a registered operation component with the given input.
	// Returns the operation result or an error if execution fails.
	ExecuteOperation(ctx context.Context, operationID ComponentID, input Input) (Output, error)

	// StartService starts a registered service component.
	// Services run asynchronously in the background.
	StartService(ctx context.Context, serviceID ComponentID) error

	// StopService stops a running service component gracefully.
	StopService(ctx context.Context, serviceID ComponentID) error

	// Start initializes and starts the entire system.
	// Should be idempotent - safe to call multiple times.
	Start(ctx context.Context) error

	// Stop gracefully shuts down the entire system.
	// Should be idempotent - safe to call multiple times.
	Stop(ctx context.Context) error

	// IsRunning returns whether the system is currently running.
	IsRunning() bool
}
