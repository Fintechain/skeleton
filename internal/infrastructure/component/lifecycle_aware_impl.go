package component

import (
	"sync"

	"github.com/fintechain/skeleton/internal/domain/component"
)

// LifecycleAwareComponentImpl provides a concrete implementation of the LifecycleAwareComponent interface.
type LifecycleAwareComponentImpl struct {
	component.Component
	state     component.LifecycleState
	callbacks []func(oldState, newState component.LifecycleState)
	mu        sync.RWMutex
}

// NewLifecycleAwareComponent creates a new lifecycle-aware component instance.
// This constructor accepts a base component interface dependency.
func NewLifecycleAwareComponent(baseComponent component.Component) component.LifecycleAwareComponent {
	return &LifecycleAwareComponentImpl{
		Component: baseComponent,
		state:     component.StateCreated, // Initial state
		callbacks: make([]func(oldState, newState component.LifecycleState), 0),
	}
}

// State returns the current lifecycle state of the component.
func (l *LifecycleAwareComponentImpl) State() component.LifecycleState {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.state
}

// SetState sets the lifecycle state of the component.
// This method is thread-safe and triggers state change callbacks.
func (l *LifecycleAwareComponentImpl) SetState(state component.LifecycleState) {
	l.mu.Lock()
	oldState := l.state
	l.state = state

	// Copy callbacks to avoid holding lock during callback execution
	callbacks := make([]func(oldState, newState component.LifecycleState), len(l.callbacks))
	copy(callbacks, l.callbacks)
	l.mu.Unlock()

	// Execute callbacks outside of lock to prevent deadlocks
	if oldState != state {
		for _, callback := range callbacks {
			if callback != nil {
				// Execute callback in a safe manner
				func() {
					defer func() {
						// Recover from panics in callbacks to prevent crashing the component
						if r := recover(); r != nil {
							// Log the panic if logging is available, but don't crash
							// For now, we silently recover to maintain stability
						}
					}()
					callback(oldState, state)
				}()
			}
		}
	}
}

// OnStateChange registers a callback to be called when the state changes.
// Multiple callbacks can be registered and they will be called in registration order.
func (l *LifecycleAwareComponentImpl) OnStateChange(callback func(oldState, newState component.LifecycleState)) {
	if callback == nil {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	l.callbacks = append(l.callbacks, callback)
}
