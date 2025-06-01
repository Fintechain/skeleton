package component

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/fintechain/skeleton/pkg/component"
	"github.com/fintechain/skeleton/test/unit/mocks"
)

// TestNewLifecycleAwareComponent tests the constructor function
func TestNewLifecycleAwareComponent(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()

	// Test constructor with valid component
	lifecycleAware := component.NewLifecycleAwareComponent(mockComponent)

	assert.NotNil(t, lifecycleAware)

	// Verify interface compliance
	var _ component.LifecycleAwareComponent = lifecycleAware

	// Verify initial state
	assert.Equal(t, component.StateCreated, lifecycleAware.State())
}

// TestNewLifecycleAwareComponentWithNilBase tests constructor with nil base component
func TestNewLifecycleAwareComponentWithNilBase(t *testing.T) {
	// Test constructor with nil base component
	lifecycleAware := component.NewLifecycleAwareComponent(nil)

	assert.NotNil(t, lifecycleAware)

	// Should still be a valid lifecycle-aware component
	var _ component.LifecycleAwareComponent = lifecycleAware

	// Should have initial state
	assert.Equal(t, component.StateCreated, lifecycleAware.State())
}

// TestLifecycleStateManagement tests basic state management operations
func TestLifecycleStateManagement(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()

	lifecycleAware := component.NewLifecycleAwareComponent(mockComponent)

	// Test initial state
	assert.Equal(t, component.StateCreated, lifecycleAware.State())

	// Test state transitions
	lifecycleAware.SetState(component.StateInitializing)
	assert.Equal(t, component.StateInitializing, lifecycleAware.State())

	lifecycleAware.SetState(component.StateInitialized)
	assert.Equal(t, component.StateInitialized, lifecycleAware.State())

	lifecycleAware.SetState(component.StateActive)
	assert.Equal(t, component.StateActive, lifecycleAware.State())

	lifecycleAware.SetState(component.StateDisposing)
	assert.Equal(t, component.StateDisposing, lifecycleAware.State())

	lifecycleAware.SetState(component.StateDisposed)
	assert.Equal(t, component.StateDisposed, lifecycleAware.State())

	lifecycleAware.SetState(component.StateFailed)
	assert.Equal(t, component.StateFailed, lifecycleAware.State())
}

// TestLifecycleStateCallbacks tests state change callbacks
func TestLifecycleStateCallbacks(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()

	lifecycleAware := component.NewLifecycleAwareComponent(mockComponent)

	// Track callback invocations
	var callbackInvocations []component.LifecycleState
	var mu sync.Mutex

	// Register callback
	lifecycleAware.OnStateChange(func(oldState, newState component.LifecycleState) {
		mu.Lock()
		defer mu.Unlock()
		callbackInvocations = append(callbackInvocations, newState)
	})

	// Trigger state changes
	lifecycleAware.SetState(component.StateInitializing)
	lifecycleAware.SetState(component.StateInitialized)
	lifecycleAware.SetState(component.StateActive)

	// Allow callbacks to complete
	time.Sleep(10 * time.Millisecond)

	// Verify callbacks were invoked
	mu.Lock()
	defer mu.Unlock()
	assert.Len(t, callbackInvocations, 3)
	assert.Equal(t, component.StateInitializing, callbackInvocations[0])
	assert.Equal(t, component.StateInitialized, callbackInvocations[1])
	assert.Equal(t, component.StateActive, callbackInvocations[2])
}

// TestMultipleStateCallbacks tests multiple callback registration
func TestMultipleStateCallbacks(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()

	lifecycleAware := component.NewLifecycleAwareComponent(mockComponent)

	// Track callback invocations
	var callback1Count, callback2Count int
	var mu sync.Mutex

	// Register multiple callbacks
	lifecycleAware.OnStateChange(func(oldState, newState component.LifecycleState) {
		mu.Lock()
		defer mu.Unlock()
		callback1Count++
	})

	lifecycleAware.OnStateChange(func(oldState, newState component.LifecycleState) {
		mu.Lock()
		defer mu.Unlock()
		callback2Count++
	})

	// Trigger state change
	lifecycleAware.SetState(component.StateActive)

	// Allow callbacks to complete
	time.Sleep(10 * time.Millisecond)

	// Verify both callbacks were invoked
	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, 1, callback1Count)
	assert.Equal(t, 1, callback2Count)
}

// TestStateCallbackWithOldAndNewState tests callback parameters
func TestStateCallbackWithOldAndNewState(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()

	lifecycleAware := component.NewLifecycleAwareComponent(mockComponent)

	// Track state transitions
	var transitions []struct {
		oldState component.LifecycleState
		newState component.LifecycleState
	}
	var mu sync.Mutex

	// Register callback
	lifecycleAware.OnStateChange(func(oldState, newState component.LifecycleState) {
		mu.Lock()
		defer mu.Unlock()
		transitions = append(transitions, struct {
			oldState component.LifecycleState
			newState component.LifecycleState
		}{oldState, newState})
	})

	// Trigger state changes
	lifecycleAware.SetState(component.StateInitializing)
	lifecycleAware.SetState(component.StateActive)

	// Allow callbacks to complete
	time.Sleep(10 * time.Millisecond)

	// Verify callback parameters
	mu.Lock()
	defer mu.Unlock()
	assert.Len(t, transitions, 2)
	assert.Equal(t, component.StateCreated, transitions[0].oldState)
	assert.Equal(t, component.StateInitializing, transitions[0].newState)
	assert.Equal(t, component.StateInitializing, transitions[1].oldState)
	assert.Equal(t, component.StateActive, transitions[1].newState)
}

// TestConcurrentStateTransitions tests thread-safe state transitions
func TestConcurrentStateTransitions(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()

	lifecycleAware := component.NewLifecycleAwareComponent(mockComponent)

	var wg sync.WaitGroup
	numGoroutines := 10
	numTransitionsPerGoroutine := 100

	// Track callback invocations
	var callbackCount int
	var mu sync.Mutex

	// Register callback
	lifecycleAware.OnStateChange(func(oldState, newState component.LifecycleState) {
		mu.Lock()
		defer mu.Unlock()
		callbackCount++
	})

	// Available states for cycling
	states := []component.LifecycleState{
		component.StateCreated,
		component.StateInitializing,
		component.StateInitialized,
		component.StateActive,
		component.StateDisposing,
		component.StateDisposed,
		component.StateFailed,
	}

	// Concurrent state transitions
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numTransitionsPerGoroutine; j++ {
				// Cycle through states
				state := states[j%len(states)]
				lifecycleAware.SetState(state)
			}
		}(i)
	}

	// Concurrent state reads
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numTransitionsPerGoroutine; j++ {
				lifecycleAware.State()
			}
		}()
	}

	wg.Wait()

	// Allow callbacks to complete
	time.Sleep(100 * time.Millisecond)

	// Verify that operations completed without panics
	finalState := lifecycleAware.State()
	assert.Contains(t, states, finalState)

	// Verify callbacks were invoked (exact count may vary due to concurrency)
	mu.Lock()
	defer mu.Unlock()
	assert.True(t, callbackCount > 0)
}

// TestLifecycleAwareComponentInterfaceCompliance tests interface compliance
func TestLifecycleAwareComponentInterfaceCompliance(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()

	// Test interface compliance
	var _ component.LifecycleAwareComponent = component.NewLifecycleAwareComponent(mockComponent)
}

// TestLifecycleStateConstants tests state constant values
func TestLifecycleStateConstants(t *testing.T) {
	// Test that state constants have expected string values
	assert.Equal(t, "created", string(component.StateCreated))
	assert.Equal(t, "initializing", string(component.StateInitializing))
	assert.Equal(t, "initialized", string(component.StateInitialized))
	assert.Equal(t, "active", string(component.StateActive))
	assert.Equal(t, "disposing", string(component.StateDisposing))
	assert.Equal(t, "disposed", string(component.StateDisposed))
	assert.Equal(t, "failed", string(component.StateFailed))
}

// TestLifecycleStateTransitionValidation tests state transition validation
func TestLifecycleStateTransitionValidation(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()

	lifecycleAware := component.NewLifecycleAwareComponent(mockComponent)

	// Test valid transitions
	validTransitions := []struct {
		from component.LifecycleState
		to   component.LifecycleState
	}{
		{component.StateCreated, component.StateInitializing},
		{component.StateInitializing, component.StateInitialized},
		{component.StateInitialized, component.StateActive},
		{component.StateActive, component.StateDisposing},
		{component.StateDisposing, component.StateDisposed},
		{component.StateActive, component.StateFailed},
	}

	for _, transition := range validTransitions {
		lifecycleAware.SetState(transition.from)
		assert.Equal(t, transition.from, lifecycleAware.State())

		lifecycleAware.SetState(transition.to)
		assert.Equal(t, transition.to, lifecycleAware.State())
	}
}

// TestLifecycleAwareComponentWithBaseComponent tests integration with base component
func TestLifecycleAwareComponentWithBaseComponent(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()

	// Configure mock component metadata
	metadata := component.Metadata{
		"version": "1.0.0",
		"author":  "test",
	}
	if mockComp, ok := mockComponent.(*mocks.MockComponent); ok {
		mockComp.SetMetadata(metadata)
	}

	lifecycleAware := component.NewLifecycleAwareComponent(mockComponent)

	// Test that base component functionality is preserved
	retrievedMetadata := mockComponent.Metadata()
	assert.Equal(t, metadata, retrievedMetadata)

	// Test that lifecycle functionality is added
	assert.Equal(t, component.StateCreated, lifecycleAware.State())

	lifecycleAware.SetState(component.StateActive)
	assert.Equal(t, component.StateActive, lifecycleAware.State())
}

// TestLifecycleCallbackErrorHandling tests callback error handling
func TestLifecycleCallbackErrorHandling(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()

	lifecycleAware := component.NewLifecycleAwareComponent(mockComponent)

	// Register callback that panics
	lifecycleAware.OnStateChange(func(oldState, newState component.LifecycleState) {
		panic("callback error")
	})

	// Register normal callback
	var normalCallbackInvoked bool
	var mu sync.Mutex
	lifecycleAware.OnStateChange(func(oldState, newState component.LifecycleState) {
		mu.Lock()
		defer mu.Unlock()
		normalCallbackInvoked = true
	})

	// State change should not panic even if callback does
	assert.NotPanics(t, func() {
		lifecycleAware.SetState(component.StateActive)
	})

	// Allow callbacks to complete
	time.Sleep(10 * time.Millisecond)

	// Normal callback should still be invoked
	mu.Lock()
	defer mu.Unlock()
	assert.True(t, normalCallbackInvoked)
}

// TestLifecycleStateStringRepresentation tests state string representation
func TestLifecycleStateStringRepresentation(t *testing.T) {
	// Test that states have meaningful string representations
	states := []component.LifecycleState{
		component.StateCreated,
		component.StateInitializing,
		component.StateInitialized,
		component.StateActive,
		component.StateDisposing,
		component.StateDisposed,
		component.StateFailed,
	}

	expectedStrings := []string{
		"created",
		"initializing",
		"initialized",
		"active",
		"disposing",
		"disposed",
		"failed",
	}

	for i, state := range states {
		assert.Equal(t, expectedStrings[i], string(state))
	}
}

// TestLifecycleAwareComponentComplexScenario tests complex lifecycle scenarios
func TestLifecycleAwareComponentComplexScenario(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()

	lifecycleAware := component.NewLifecycleAwareComponent(mockComponent)

	// Track all state changes
	var stateHistory []component.LifecycleState
	var mu sync.Mutex

	lifecycleAware.OnStateChange(func(oldState, newState component.LifecycleState) {
		mu.Lock()
		defer mu.Unlock()
		stateHistory = append(stateHistory, newState)
	})

	// Simulate complete lifecycle
	lifecycleStates := []component.LifecycleState{
		component.StateInitializing,
		component.StateInitialized,
		component.StateActive,
		component.StateDisposing,
		component.StateDisposed,
	}

	for _, state := range lifecycleStates {
		lifecycleAware.SetState(state)
		assert.Equal(t, state, lifecycleAware.State())
	}

	// Allow callbacks to complete
	time.Sleep(10 * time.Millisecond)

	// Verify complete state history
	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, lifecycleStates, stateHistory)
}

// TestLifecycleAwareComponentMemoryUsage tests memory usage patterns
func TestLifecycleAwareComponentMemoryUsage(t *testing.T) {
	factory := mocks.NewFactory()

	// Create many lifecycle-aware components
	components := make([]component.LifecycleAwareComponent, 1000)
	for i := 0; i < 1000; i++ {
		mockComponent := factory.ComponentInterface()
		components[i] = component.NewLifecycleAwareComponent(mockComponent)
	}

	// Register callbacks on all components
	for _, comp := range components {
		comp.OnStateChange(func(oldState, newState component.LifecycleState) {
			// Simple callback
		})
	}

	// Trigger state changes
	for _, comp := range components {
		comp.SetState(component.StateActive)
	}

	// Verify all components are in expected state
	for _, comp := range components {
		assert.Equal(t, component.StateActive, comp.State())
	}

	// This test primarily ensures no memory leaks or excessive allocations
	// The actual memory usage would need to be measured with profiling tools
}
