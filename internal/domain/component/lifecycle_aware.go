package component

// LifecycleState represents the state of a component in its lifecycle.
type LifecycleState string

const (
	// StateCreated indicates a component has been created but not initialized.
	StateCreated LifecycleState = "created"

	// StateInitializing indicates a component is in the process of initializing.
	StateInitializing LifecycleState = "initializing"

	// StateInitialized indicates a component has been successfully initialized.
	StateInitialized LifecycleState = "initialized"

	// StateActive indicates a component is active and operational.
	StateActive LifecycleState = "active"

	// StateDisposing indicates a component is in the process of being disposed.
	StateDisposing LifecycleState = "disposing"

	// StateDisposed indicates a component has been disposed.
	StateDisposed LifecycleState = "disposed"

	// StateFailed indicates a component has failed during initialization or operation.
	StateFailed LifecycleState = "failed"
)

// LifecycleAwareComponent is an interface for components that want to manage their lifecycle.
// It extends the Component interface to include lifecycle management capabilities.
type LifecycleAwareComponent interface {
	Component // Embed the Component interface to inherit all its methods

	// State returns the current lifecycle state of the component.
	State() LifecycleState

	// SetState sets the lifecycle state of the component.
	SetState(state LifecycleState)

	// OnStateChange registers a callback to be called when the state changes.
	OnStateChange(callback func(oldState, newState LifecycleState))
}
