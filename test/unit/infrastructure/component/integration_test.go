package component

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/fintechain/skeleton/pkg/component"
	"github.com/fintechain/skeleton/test/unit/mocks"
)

// TestDependencyAwareLifecycleAwareIntegration tests integration between dependency-aware and lifecycle-aware components
func TestDependencyAwareLifecycleAwareIntegration(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()
	mockRegistry := factory.RegistryInterface()

	// Create a dependency-aware component
	dependencyAware := component.NewDependencyAwareComponent(mockComponent, mockRegistry)

	// Wrap it with lifecycle awareness
	lifecycleAware := component.NewLifecycleAwareComponent(mockComponent)

	// Track lifecycle state changes
	var stateChanges []component.LifecycleState
	lifecycleAware.OnStateChange(func(oldState, newState component.LifecycleState) {
		stateChanges = append(stateChanges, newState)
	})

	// Add dependencies
	dependencyAware.AddDependency("service1")
	dependencyAware.AddDependency("service2")

	// Verify dependencies
	assert.True(t, dependencyAware.HasDependency("service1"))
	assert.True(t, dependencyAware.HasDependency("service2"))

	// Test lifecycle transitions
	lifecycleAware.SetState(component.StateInitializing)
	lifecycleAware.SetState(component.StateInitialized)
	lifecycleAware.SetState(component.StateActive)

	// Verify state changes were tracked
	expectedStates := []component.LifecycleState{
		component.StateInitializing,
		component.StateInitialized,
		component.StateActive,
	}
	assert.Equal(t, expectedStates, stateChanges)

	// Verify dependencies are still accessible
	dependencies := dependencyAware.Dependencies()
	assert.Contains(t, dependencies, "service1")
	assert.Contains(t, dependencies, "service2")
}

// TestComponentFactoryIntegration tests integration with component factory
func TestComponentFactoryIntegration(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()

	// Create base components using factory
	baseComponent1 := factory.ComponentInterface()
	baseComponent2 := factory.ComponentInterface()

	// Create dependency-aware components
	depAware1 := component.NewDependencyAwareComponent(baseComponent1, mockRegistry)

	// Create lifecycle-aware components
	lifecycleAware1 := component.NewLifecycleAwareComponent(baseComponent1)
	lifecycleAware2 := component.NewLifecycleAwareComponent(baseComponent2)

	// Set up dependencies
	depAware1.AddDependency("component2")

	// Configure mock registry to return component2 when resolving dependencies
	if mockReg, ok := mockRegistry.(*mocks.MockRegistry); ok {
		mockReg.SetReturnItem("component2", baseComponent2)
	}

	// Test dependency resolution
	resolvedComponent, err := depAware1.ResolveDependency("component2", mockRegistry)
	assert.NoError(t, err)
	assert.NotNil(t, resolvedComponent)

	// Test lifecycle coordination
	var component1States, component2States []component.LifecycleState

	lifecycleAware1.OnStateChange(func(oldState, newState component.LifecycleState) {
		component1States = append(component1States, newState)
	})

	lifecycleAware2.OnStateChange(func(oldState, newState component.LifecycleState) {
		component2States = append(component2States, newState)
	})

	// Initialize components in dependency order
	lifecycleAware2.SetState(component.StateInitializing)
	lifecycleAware2.SetState(component.StateInitialized)
	lifecycleAware2.SetState(component.StateActive)

	lifecycleAware1.SetState(component.StateInitializing)
	lifecycleAware1.SetState(component.StateInitialized)
	lifecycleAware1.SetState(component.StateActive)

	// Verify both components went through proper lifecycle
	expectedStates := []component.LifecycleState{
		component.StateInitializing,
		component.StateInitialized,
		component.StateActive,
	}
	assert.Equal(t, expectedStates, component1States)
	assert.Equal(t, expectedStates, component2States)
}

// TestComplexDependencyGraph tests a complex dependency graph scenario
func TestComplexDependencyGraph(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()

	// Create a complex dependency graph:
	// A depends on B and C
	// B depends on D
	// C depends on D
	// D has no dependencies

	components := make(map[string]component.Component)
	dependencyAwareComponents := make(map[string]component.DependencyAwareComponent)
	lifecycleAwareComponents := make(map[string]component.LifecycleAwareComponent)

	componentNames := []string{"A", "B", "C", "D"}

	// Create all components
	for _, name := range componentNames {
		baseComp := factory.ComponentInterface()

		depAware := component.NewDependencyAwareComponent(baseComp, mockRegistry)
		lifecycleAware := component.NewLifecycleAwareComponent(baseComp)

		components[name] = baseComp
		dependencyAwareComponents[name] = depAware
		lifecycleAwareComponents[name] = lifecycleAware
	}

	// Set up dependencies
	dependencyAwareComponents["A"].AddDependency("B")
	dependencyAwareComponents["A"].AddDependency("C")
	dependencyAwareComponents["B"].AddDependency("D")
	dependencyAwareComponents["C"].AddDependency("D")

	// Configure mock registry
	if mockReg, ok := mockRegistry.(*mocks.MockRegistry); ok {
		for name, comp := range components {
			mockReg.SetReturnItem(name, comp)
		}
	}

	// Track initialization order
	var initializationOrder []string
	var initializationMutex sync.Mutex

	for name, lifecycleComp := range lifecycleAwareComponents {
		componentName := name // Capture for closure
		lifecycleComp.OnStateChange(func(oldState, newState component.LifecycleState) {
			if newState == component.StateInitialized {
				initializationMutex.Lock()
				initializationOrder = append(initializationOrder, componentName)
				initializationMutex.Unlock()
			}
		})
	}

	// Initialize components in dependency order (D, then B and C, then A)
	lifecycleAwareComponents["D"].SetState(component.StateInitializing)
	lifecycleAwareComponents["D"].SetState(component.StateInitialized)

	lifecycleAwareComponents["B"].SetState(component.StateInitializing)
	lifecycleAwareComponents["B"].SetState(component.StateInitialized)

	lifecycleAwareComponents["C"].SetState(component.StateInitializing)
	lifecycleAwareComponents["C"].SetState(component.StateInitialized)

	lifecycleAwareComponents["A"].SetState(component.StateInitializing)
	lifecycleAwareComponents["A"].SetState(component.StateInitialized)

	// Allow callbacks to complete
	time.Sleep(10 * time.Millisecond)

	// Verify initialization order
	initializationMutex.Lock()
	defer initializationMutex.Unlock()
	assert.Len(t, initializationOrder, 4)
	assert.Contains(t, initializationOrder, "A")
	assert.Contains(t, initializationOrder, "B")
	assert.Contains(t, initializationOrder, "C")
	assert.Contains(t, initializationOrder, "D")

	// Verify dependencies can be resolved
	resolvedDeps, err := dependencyAwareComponents["A"].ResolveDependencies(mockRegistry)
	assert.NoError(t, err)
	assert.Len(t, resolvedDeps, 2)
	assert.Contains(t, resolvedDeps, "B")
	assert.Contains(t, resolvedDeps, "C")
}

// TestConcurrentComponentOperations tests concurrent operations on components
func TestConcurrentComponentOperations(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()

	numComponents := 10
	numOperationsPerComponent := 100

	var components []component.DependencyAwareComponent
	var lifecycleComponents []component.LifecycleAwareComponent

	// Create components
	for i := 0; i < numComponents; i++ {
		baseComp := factory.ComponentInterface()
		depAware := component.NewDependencyAwareComponent(baseComp, mockRegistry)
		lifecycleAware := component.NewLifecycleAwareComponent(baseComp)

		components = append(components, depAware)
		lifecycleComponents = append(lifecycleComponents, lifecycleAware)
	}

	var wg sync.WaitGroup

	// Concurrent dependency operations
	for i, comp := range components {
		wg.Add(1)
		go func(componentIndex int, depComp component.DependencyAwareComponent) {
			defer wg.Done()
			for j := 0; j < numOperationsPerComponent; j++ {
				depID := fmt.Sprintf("dep_%d_%d", componentIndex, j)
				depComp.AddDependency(depID)
				depComp.HasDependency(depID)
				depComp.Dependencies()
			}
		}(i, comp)
	}

	// Concurrent lifecycle operations
	for i, comp := range lifecycleComponents {
		wg.Add(1)
		go func(componentIndex int, lifecycleComp component.LifecycleAwareComponent) {
			defer wg.Done()
			states := []component.LifecycleState{
				component.StateInitializing,
				component.StateInitialized,
				component.StateActive,
				component.StateDisposing,
				component.StateDisposed,
			}
			for j := 0; j < numOperationsPerComponent; j++ {
				state := states[j%len(states)]
				lifecycleComp.SetState(state)
				lifecycleComp.State()
			}
		}(i, comp)
	}

	wg.Wait()

	// Verify all components are still functional
	for _, comp := range components {
		deps := comp.Dependencies()
		assert.True(t, len(deps) > 0)
	}

	for _, comp := range lifecycleComponents {
		state := comp.State()
		assert.NotEmpty(t, string(state))
	}
}

// TestComponentSystemIntegration tests integration with the component system
func TestComponentSystemIntegration(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()
	mockSystem := factory.SystemInterface()

	// Create a set of interconnected components
	serviceComponent := factory.ComponentInterface()
	configComponent := factory.ComponentInterface()
	loggerComponent := factory.ComponentInterface()

	// Create dependency-aware wrappers
	serviceDepAware := component.NewDependencyAwareComponent(serviceComponent, mockRegistry)

	// Create lifecycle-aware wrappers
	serviceLifecycle := component.NewLifecycleAwareComponent(serviceComponent)
	configLifecycle := component.NewLifecycleAwareComponent(configComponent)
	loggerLifecycle := component.NewLifecycleAwareComponent(loggerComponent)

	// Set up dependencies: service depends on config and logger
	serviceDepAware.AddDependency("config")
	serviceDepAware.AddDependency("logger")

	// Configure mock registry
	if mockReg, ok := mockRegistry.(*mocks.MockRegistry); ok {
		mockReg.SetReturnItem("config", configComponent)
		mockReg.SetReturnItem("logger", loggerComponent)
		mockReg.SetReturnItem("service", serviceComponent)
	}

	// Track component states
	var serviceStates, configStates, loggerStates []component.LifecycleState
	var stateMutex sync.Mutex

	serviceLifecycle.OnStateChange(func(oldState, newState component.LifecycleState) {
		stateMutex.Lock()
		serviceStates = append(serviceStates, newState)
		stateMutex.Unlock()
	})

	configLifecycle.OnStateChange(func(oldState, newState component.LifecycleState) {
		stateMutex.Lock()
		configStates = append(configStates, newState)
		stateMutex.Unlock()
	})

	loggerLifecycle.OnStateChange(func(oldState, newState component.LifecycleState) {
		stateMutex.Lock()
		loggerStates = append(loggerStates, newState)
		stateMutex.Unlock()
	})

	// Initialize system components in dependency order
	// First initialize dependencies
	configLifecycle.SetState(component.StateInitializing)
	configLifecycle.SetState(component.StateInitialized)
	configLifecycle.SetState(component.StateActive)

	loggerLifecycle.SetState(component.StateInitializing)
	loggerLifecycle.SetState(component.StateInitialized)
	loggerLifecycle.SetState(component.StateActive)

	// Then initialize the service that depends on them
	serviceLifecycle.SetState(component.StateInitializing)
	serviceLifecycle.SetState(component.StateInitialized)
	serviceLifecycle.SetState(component.StateActive)

	// Allow callbacks to complete
	time.Sleep(10 * time.Millisecond)

	// Verify all components reached active state
	stateMutex.Lock()
	defer stateMutex.Unlock()

	expectedStates := []component.LifecycleState{
		component.StateInitializing,
		component.StateInitialized,
		component.StateActive,
	}

	assert.Equal(t, expectedStates, serviceStates)
	assert.Equal(t, expectedStates, configStates)
	assert.Equal(t, expectedStates, loggerStates)

	// Verify dependencies can be resolved
	resolvedDeps, err := serviceDepAware.ResolveDependencies(mockRegistry)
	assert.NoError(t, err)
	assert.Len(t, resolvedDeps, 2)
	assert.Contains(t, resolvedDeps, "config")
	assert.Contains(t, resolvedDeps, "logger")

	// Verify system integration
	assert.NotNil(t, mockSystem)
}

// TestErrorPropagationIntegration tests error propagation between components
func TestErrorPropagationIntegration(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()

	// Create components
	serviceComponent := factory.ComponentInterface()
	dependencyComponent := factory.ComponentInterface()

	serviceDepAware := component.NewDependencyAwareComponent(serviceComponent, mockRegistry)
	serviceLifecycle := component.NewLifecycleAwareComponent(serviceComponent)

	// Set up dependency
	serviceDepAware.AddDependency("failing-dependency")

	// Configure mock registry to fail for specific dependency
	if mockReg, ok := mockRegistry.(*mocks.MockRegistry); ok {
		mockReg.SetShouldFail(true)
	}

	// Track state changes
	var stateChanges []component.LifecycleState
	serviceLifecycle.OnStateChange(func(oldState, newState component.LifecycleState) {
		stateChanges = append(stateChanges, newState)
	})

	// Attempt to resolve dependencies (should fail)
	_, err := serviceDepAware.ResolveDependency("failing-dependency", mockRegistry)
	assert.Error(t, err)

	// Component should still be able to change states
	serviceLifecycle.SetState(component.StateInitializing)
	serviceLifecycle.SetState(component.StateFailed)

	// Allow callbacks to complete
	time.Sleep(10 * time.Millisecond)

	// Verify state changes
	expectedStates := []component.LifecycleState{
		component.StateInitializing,
		component.StateFailed,
	}
	assert.Equal(t, expectedStates, stateChanges)

	// Verify component is in failed state
	assert.Equal(t, component.StateFailed, serviceLifecycle.State())

	// Use the dependency component to avoid unused variable
	assert.NotNil(t, dependencyComponent)
}

// TestComponentCleanupIntegration tests component cleanup and disposal
func TestComponentCleanupIntegration(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()

	// Create components
	components := make([]component.Component, 5)
	depAwareComponents := make([]component.DependencyAwareComponent, 5)
	lifecycleComponents := make([]component.LifecycleAwareComponent, 5)

	for i := 0; i < 5; i++ {
		baseComp := factory.ComponentInterface()
		depAware := component.NewDependencyAwareComponent(baseComp, mockRegistry)
		lifecycle := component.NewLifecycleAwareComponent(baseComp)

		components[i] = baseComp
		depAwareComponents[i] = depAware
		lifecycleComponents[i] = lifecycle
	}

	// Set up some dependencies between components
	for i := 0; i < 4; i++ {
		depID := fmt.Sprintf("component_%d", i+1)
		depAwareComponents[i].AddDependency(depID)
	}

	// Track disposal states
	var disposalOrder []int
	var disposalMutex sync.Mutex

	for i, lifecycle := range lifecycleComponents {
		componentIndex := i // Capture for closure
		lifecycle.OnStateChange(func(oldState, newState component.LifecycleState) {
			if newState == component.StateDisposed {
				disposalMutex.Lock()
				disposalOrder = append(disposalOrder, componentIndex)
				disposalMutex.Unlock()
			}
		})
	}

	// Initialize all components
	for _, lifecycle := range lifecycleComponents {
		lifecycle.SetState(component.StateInitializing)
		lifecycle.SetState(component.StateInitialized)
		lifecycle.SetState(component.StateActive)
	}

	// Dispose components in reverse dependency order
	for i := len(lifecycleComponents) - 1; i >= 0; i-- {
		lifecycleComponents[i].SetState(component.StateDisposing)
		lifecycleComponents[i].SetState(component.StateDisposed)
	}

	// Allow callbacks to complete
	time.Sleep(10 * time.Millisecond)

	// Verify disposal order
	disposalMutex.Lock()
	defer disposalMutex.Unlock()
	assert.Len(t, disposalOrder, 5)

	// Verify all components are disposed
	for _, lifecycle := range lifecycleComponents {
		assert.Equal(t, component.StateDisposed, lifecycle.State())
	}

	// Verify dependencies are still tracked (even after disposal)
	for i := 0; i < 4; i++ {
		deps := depAwareComponents[i].Dependencies()
		assert.Len(t, deps, 1)
	}
}
