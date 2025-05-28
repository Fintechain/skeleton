package component

import (
	"errors"
	"testing"

	"github.com/fintechain/skeleton/internal/domain/component/mocks"
)

func TestLifecycleAwareComponentStates(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create a base component
	baseComp := NewBaseComponent("test-id", "Test Component", TypeBasic)

	// Create lifecycle-aware component
	lifecycleComp := NewLifecycleAwareComponentWithOptions(LifecycleAwareComponentOptions{
		Base:   baseComp,
		Logger: mockLogger,
	})

	// Check initial state
	if lifecycleComp.State() != StateCreated {
		t.Errorf("Initial state should be StateCreated, got %s", lifecycleComp.State())
	}

	// Test state change
	lifecycleComp.SetState(StateInitializing)
	if lifecycleComp.State() != StateInitializing {
		t.Errorf("State should be StateInitializing, got %s", lifecycleComp.State())
	}

	// Test callback
	var oldStateReceived, newStateReceived LifecycleState
	callbackCalled := false

	lifecycleComp.OnStateChange(func(oldState, newState LifecycleState) {
		callbackCalled = true
		oldStateReceived = oldState
		newStateReceived = newState
	})

	lifecycleComp.SetState(StateInitialized)

	if !callbackCalled {
		t.Error("State change callback was not called")
	}

	if oldStateReceived != StateInitializing {
		t.Errorf("Old state in callback should be StateInitializing, got %s", oldStateReceived)
	}

	if newStateReceived != StateInitialized {
		t.Errorf("New state in callback should be StateInitialized, got %s", newStateReceived)
	}
}

func TestLifecycleAwareComponentInitialize(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create a base component
	baseComp := NewBaseComponent("test-id", "Test Component", TypeBasic)

	// Create lifecycle-aware component
	lifecycleComp := NewLifecycleAwareComponentWithOptions(LifecycleAwareComponentOptions{
		Base:   baseComp,
		Logger: mockLogger,
	})

	// Initialize the component
	ctx := &mockContext{}
	err := lifecycleComp.Initialize(ctx)
	if err != nil {
		t.Errorf("Failed to initialize component: %s", err)
	}

	// Check state after initialization
	if lifecycleComp.State() != StateInitialized {
		t.Errorf("State after initialization should be StateInitialized, got %s", lifecycleComp.State())
	}
}

func TestLifecycleAwareComponentInitializeFailure(t *testing.T) {
	// Create a failing component
	failingComp := &mockFailingComponent{id: "failing-comp"}

	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create lifecycle-aware component with failing base
	lifecycleComp := NewLifecycleAwareComponentWithOptions(LifecycleAwareComponentOptions{
		Base:   failingComp,
		Logger: mockLogger,
	})

	// Initialize the component (should fail)
	ctx := &mockContext{}
	err := lifecycleComp.Initialize(ctx)
	if err == nil {
		t.Error("Expected initialization to fail, but it succeeded")
	}

	// Check state after failed initialization
	if lifecycleComp.State() != StateFailed {
		t.Errorf("State after failed initialization should be StateFailed, got %s", lifecycleComp.State())
	}
}

func TestLifecycleAwareComponentDispose(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create a base component
	baseComp := NewBaseComponent("test-id", "Test Component", TypeBasic)

	// Create lifecycle-aware component
	lifecycleComp := NewLifecycleAwareComponentWithOptions(LifecycleAwareComponentOptions{
		Base:   baseComp,
		Logger: mockLogger,
	})

	// Initialize first
	ctx := &mockContext{}
	lifecycleComp.Initialize(ctx)

	// Dispose the component
	err := lifecycleComp.Dispose()
	if err != nil {
		t.Errorf("Failed to dispose component: %s", err)
	}

	// Check state after disposal
	if lifecycleComp.State() != StateDisposed {
		t.Errorf("State after disposal should be StateDisposed, got %s", lifecycleComp.State())
	}
}

func TestLifecycleAwareComponentActivation(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create a base component
	baseComp := NewBaseComponent("test-id", "Test Component", TypeBasic)

	// Create lifecycle-aware component
	lifecycleComp := NewLifecycleAwareComponentWithOptions(LifecycleAwareComponentOptions{
		Base:   baseComp,
		Logger: mockLogger,
	})

	// Initialize first
	ctx := &mockContext{}
	lifecycleComp.Initialize(ctx)

	// Activate the component
	lifecycleComp.Activate()

	// Check state after activation
	if lifecycleComp.State() != StateActive {
		t.Errorf("State after activation should be StateActive, got %s", lifecycleComp.State())
	}

	// Deactivate the component
	lifecycleComp.Deactivate()

	// Check state after deactivation
	if lifecycleComp.State() != StateInitialized {
		t.Errorf("State after deactivation should be StateInitialized, got %s", lifecycleComp.State())
	}

	// Test activate in wrong state
	lifecycleComp.SetState(StateCreated)
	lifecycleComp.Activate()

	// Should not change state
	if lifecycleComp.State() != StateCreated {
		t.Errorf("State should remain StateCreated after invalid activation, got %s", lifecycleComp.State())
	}

	// Test deactivate in wrong state
	lifecycleComp.SetState(StateInitialized)
	lifecycleComp.Deactivate()

	// Should not change state
	if lifecycleComp.State() != StateInitialized {
		t.Errorf("State should remain StateInitialized after invalid deactivation, got %s", lifecycleComp.State())
	}
}

// mockFailingComponent is a component that fails to initialize and dispose
type mockFailingComponent struct {
	id string
}

func (m *mockFailingComponent) ID() string {
	return m.id
}

func (m *mockFailingComponent) Name() string {
	return "Failing Component"
}

func (m *mockFailingComponent) Type() ComponentType {
	return TypeBasic
}

func (m *mockFailingComponent) Metadata() Metadata {
	return make(Metadata)
}

func (m *mockFailingComponent) Initialize(ctx Context) error {
	return errors.New("mock initialization failure")
}

func (m *mockFailingComponent) Dispose() error {
	return errors.New("mock disposal failure")
}

func TestNewLifecycleAwareComponent(t *testing.T) {
	// Test the legacy constructor
	baseComp := NewBaseComponent("test-id", "Test Component", TypeBasic)

	// Create component using legacy constructor
	lifecycleComp := NewLifecycleAwareComponent(baseComp)

	if lifecycleComp == nil {
		t.Error("NewLifecycleAwareComponent returned nil")
	}

	// Verify the base component was set correctly
	if lifecycleComp.ID() != "test-id" {
		t.Errorf("Component has wrong ID: got %s, expected test-id", lifecycleComp.ID())
	}

	if lifecycleComp.Name() != "Test Component" {
		t.Errorf("Component has wrong name: got %s, expected Test Component", lifecycleComp.Name())
	}

	// Verify initial state
	if lifecycleComp.State() != StateCreated {
		t.Errorf("Initial state should be StateCreated, got %s", lifecycleComp.State())
	}

	// Verify we can use the component normally
	lifecycleComp.SetState(StateInitialized)
	if lifecycleComp.State() != StateInitialized {
		t.Errorf("State should be StateInitialized after setting, got %s", lifecycleComp.State())
	}

	// Test basic lifecycle operations
	ctx := &mockContext{}
	err := lifecycleComp.Initialize(ctx)
	if err != nil {
		t.Errorf("Failed to initialize component: %s", err)
	}

	lifecycleComp.Activate()
	if lifecycleComp.State() != StateActive {
		t.Errorf("State should be StateActive after activating, got %s", lifecycleComp.State())
	}
}

func TestLifecycleAwareComponentDisposeFailure(t *testing.T) {
	// Create a failing component
	failingComp := &mockFailingComponent{id: "failing-comp"}

	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create lifecycle-aware component with failing base
	lifecycleComp := NewLifecycleAwareComponentWithOptions(LifecycleAwareComponentOptions{
		Base:   failingComp,
		Logger: mockLogger,
	})

	// Dispose the component (should fail)
	err := lifecycleComp.Dispose()
	if err == nil {
		t.Error("Expected dispose to fail, but it succeeded")
	}

	// Check state after failed disposal
	if lifecycleComp.State() != StateFailed {
		t.Errorf("State after failed disposal should be StateFailed, got %s", lifecycleComp.State())
	}
}
