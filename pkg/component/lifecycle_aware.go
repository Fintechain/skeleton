// Package component provides component interfaces and types.
package component

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	componentImpl "github.com/fintechain/skeleton/internal/infrastructure/component"
)

// Re-export lifecycle-aware component interface
type LifecycleAwareComponent = component.LifecycleAwareComponent

// Re-export lifecycle types
type LifecycleState = component.LifecycleState

// Re-export lifecycle state constants
const (
	StateCreated      = component.StateCreated
	StateInitializing = component.StateInitializing
	StateInitialized  = component.StateInitialized
	StateActive       = component.StateActive
	StateDisposing    = component.StateDisposing
	StateDisposed     = component.StateDisposed
	StateFailed       = component.StateFailed
)

// Factory functions

// NewLifecycleAwareComponent creates a new lifecycle-aware component.
func NewLifecycleAwareComponent(base Component) LifecycleAwareComponent {
	return componentImpl.NewLifecycleAwareComponent(base)
}
