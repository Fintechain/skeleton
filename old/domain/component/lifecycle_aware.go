package component

import (
	"github.com/fintechain/skeleton/internal/infrastructure/logging"
)

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

// LifecycleAware is an interface for components that want to manage their lifecycle.
type LifecycleAware interface {
	// State returns the current lifecycle state of the component.
	State() LifecycleState

	// SetState sets the lifecycle state of the component.
	SetState(state LifecycleState)

	// OnStateChange registers a callback to be called when the state changes.
	OnStateChange(callback func(oldState, newState LifecycleState))
}

// LifecycleAwareComponent is a component that is aware of its lifecycle.
type LifecycleAwareComponent struct {
	Component
	state     LifecycleState
	callbacks []func(oldState, newState LifecycleState)
	logger    logging.Logger
}

// LifecycleAwareComponentOptions contains options for creating a LifecycleAwareComponent.
type LifecycleAwareComponentOptions struct {
	Base   Component
	Logger logging.Logger
}

// NewLifecycleAwareComponent creates a new lifecycle-aware component.
func NewLifecycleAwareComponent(base Component) *LifecycleAwareComponent {
	return &LifecycleAwareComponent{
		Component: base,
		state:     StateCreated,
		callbacks: make([]func(oldState, newState LifecycleState), 0),
	}
}

// NewLifecycleAwareComponentWithOptions creates a new lifecycle-aware component with injected dependencies.
func NewLifecycleAwareComponentWithOptions(options LifecycleAwareComponentOptions) *LifecycleAwareComponent {
	return &LifecycleAwareComponent{
		Component: options.Base,
		state:     StateCreated,
		callbacks: make([]func(oldState, newState LifecycleState), 0),
		logger:    options.Logger,
	}
}

// State returns the current lifecycle state of the component.
func (c *LifecycleAwareComponent) State() LifecycleState {
	return c.state
}

// SetState sets the lifecycle state of the component.
func (c *LifecycleAwareComponent) SetState(state LifecycleState) {
	oldState := c.state
	c.state = state

	if c.logger != nil {
		c.logger.Debug("Component %s (%s) changing state: %s -> %s",
			c.Component.Name(), c.Component.ID(), oldState, state)
	}

	// Notify state change callbacks
	for _, callback := range c.callbacks {
		callback(oldState, state)
	}
}

// OnStateChange registers a callback to be called when the state changes.
func (c *LifecycleAwareComponent) OnStateChange(callback func(oldState, newState LifecycleState)) {
	c.callbacks = append(c.callbacks, callback)
}

// Initialize initializes the component with lifecycle awareness.
func (c *LifecycleAwareComponent) Initialize(ctx Context) error {
	// Set state to initializing
	c.SetState(StateInitializing)

	// Initialize the base component
	err := c.Component.Initialize(ctx)
	if err != nil {
		c.SetState(StateFailed)
		if c.logger != nil {
			c.logger.Error("Failed to initialize component %s (%s): %s",
				c.Component.Name(), c.Component.ID(), err)
		}
		return err
	}

	// Set state to initialized
	c.SetState(StateInitialized)
	return nil
}

// Dispose disposes the component with lifecycle awareness.
func (c *LifecycleAwareComponent) Dispose() error {
	// Set state to disposing
	c.SetState(StateDisposing)

	// Dispose the base component
	err := c.Component.Dispose()
	if err != nil {
		c.SetState(StateFailed)
		if c.logger != nil {
			c.logger.Error("Failed to dispose component %s (%s): %s",
				c.Component.Name(), c.Component.ID(), err)
		}
		return err
	}

	// Set state to disposed
	c.SetState(StateDisposed)
	return nil
}

// Activate sets the component state to active.
func (c *LifecycleAwareComponent) Activate() {
	if c.state == StateInitialized {
		c.SetState(StateActive)
	} else if c.logger != nil {
		c.logger.Warn("Cannot activate component %s (%s): not in initialized state (current: %s)",
			c.Component.Name(), c.Component.ID(), c.state)
	}
}

// Deactivate sets the component state from active to initialized.
func (c *LifecycleAwareComponent) Deactivate() {
	if c.state == StateActive {
		c.SetState(StateInitialized)
	} else if c.logger != nil {
		c.logger.Warn("Cannot deactivate component %s (%s): not in active state (current: %s)",
			c.Component.Name(), c.Component.ID(), c.state)
	}
}
