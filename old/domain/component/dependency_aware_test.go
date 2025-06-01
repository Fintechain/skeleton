package component

import (
	"testing"

	"github.com/fintechain/skeleton/internal/domain/component/mocks"
)

func TestDependencyAwareComponentBasics(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create a base component
	baseComp := NewBaseComponent("test-id", "Test Component", TypeBasic)

	// Create dependency-aware component with initial dependencies
	initialDeps := []string{"dep1", "dep2"}
	depAwareComp := NewDependencyAwareComponentWithOptions(DependencyAwareComponentOptions{
		Base:         baseComp,
		Dependencies: initialDeps,
		Logger:       mockLogger,
	})

	// Check initial dependencies
	deps := depAwareComp.Dependencies()
	if len(deps) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(deps))
	}

	if deps[0] != "dep1" || deps[1] != "dep2" {
		t.Errorf("Dependencies not set correctly: %v", deps)
	}

	// Check has dependency
	if !depAwareComp.HasDependency("dep1") {
		t.Error("Component should have dependency 'dep1'")
	}

	if depAwareComp.HasDependency("non-existent") {
		t.Error("Component should not have dependency 'non-existent'")
	}
}

func TestDependencyAwareComponentAddRemove(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create a base component
	baseComp := NewBaseComponent("test-id", "Test Component", TypeBasic)

	// Create dependency-aware component with no initial dependencies
	depAwareComp := NewDependencyAwareComponentWithOptions(DependencyAwareComponentOptions{
		Base:         baseComp,
		Dependencies: []string{},
		Logger:       mockLogger,
	})

	// Add dependencies
	depAwareComp.AddDependency("dep1")
	depAwareComp.AddDependency("dep2")

	// Check dependencies were added
	deps := depAwareComp.Dependencies()
	if len(deps) != 2 {
		t.Errorf("Expected 2 dependencies after adding, got %d", len(deps))
	}

	// Adding the same dependency again should be a no-op
	depAwareComp.AddDependency("dep1")
	deps = depAwareComp.Dependencies()
	if len(deps) != 2 {
		t.Errorf("Expected still 2 dependencies after adding duplicate, got %d", len(deps))
	}

	// Remove a dependency
	depAwareComp.RemoveDependency("dep1")
	deps = depAwareComp.Dependencies()
	if len(deps) != 1 || deps[0] != "dep2" {
		t.Errorf("Expected only 'dep2' after removal, got: %v", deps)
	}

	// Removing a non-existent dependency should be a no-op
	depAwareComp.RemoveDependency("non-existent")
	deps = depAwareComp.Dependencies()
	if len(deps) != 1 {
		t.Errorf("Expected still 1 dependency after removing non-existent, got %d", len(deps))
	}
}

func TestDependencyAwareComponentResolve(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create a registry
	registry := NewRegistry(DefaultRegistryOptions{
		Logger: mockLogger,
	})

	// Create components for registry
	comp1 := NewBaseComponent("dep1", "Dependency 1", TypeBasic)
	comp2 := NewBaseComponent("dep2", "Dependency 2", TypeService)

	// Register components
	registry.Register(comp1)
	registry.Register(comp2)

	// Create a dependency-aware component
	baseComp := NewBaseComponent("test-id", "Test Component", TypeBasic)
	depAwareComp := NewDependencyAwareComponentWithOptions(DependencyAwareComponentOptions{
		Base:         baseComp,
		Dependencies: []string{"dep1", "dep2"},
		Logger:       mockLogger,
	})

	// Test resolving a single dependency
	resolvedDep, err := depAwareComp.ResolveDependency("dep1", registry)
	if err != nil {
		t.Errorf("Failed to resolve dependency: %s", err)
	}

	if resolvedDep.ID() != "dep1" {
		t.Errorf("Resolved wrong dependency: got %s, expected dep1", resolvedDep.ID())
	}

	// Test resolving a non-existent dependency
	_, err = depAwareComp.ResolveDependency("non-existent", registry)
	if err == nil {
		t.Error("Expected error resolving non-existent dependency, but got nil")
	}

	if !IsComponentError(err, ErrDependencyNotFound) {
		t.Errorf("Expected dependency not found error, but got: %s", err)
	}

	// Test resolving all dependencies
	allDeps, err := depAwareComp.ResolveDependencies(registry)
	if err != nil {
		t.Errorf("Failed to resolve all dependencies: %s", err)
	}

	if len(allDeps) != 2 {
		t.Errorf("Expected 2 resolved dependencies, got %d", len(allDeps))
	}

	// Remove a dependency from registry and try to resolve again
	registry.Unregister("dep2")

	_, err = depAwareComp.ResolveDependencies(registry)
	if err == nil {
		t.Error("Expected error resolving all dependencies with missing dep, but got nil")
	}
}

func TestNewDependencyAwareComponent(t *testing.T) {
	// Test the legacy constructor
	baseComp := NewBaseComponent("test-id", "Test Component", TypeBasic)
	initialDeps := []string{"dep1", "dep2"}

	// Create component using legacy constructor
	depAwareComp := NewDependencyAwareComponent(baseComp, initialDeps)

	if depAwareComp == nil {
		t.Error("NewDependencyAwareComponent returned nil")
	}

	// Verify the base component was set correctly
	if depAwareComp.ID() != "test-id" {
		t.Errorf("Component has wrong ID: got %s, expected test-id", depAwareComp.ID())
	}

	if depAwareComp.Name() != "Test Component" {
		t.Errorf("Component has wrong name: got %s, expected Test Component", depAwareComp.Name())
	}

	// Verify dependencies were set correctly
	deps := depAwareComp.Dependencies()
	if len(deps) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(deps))
	}

	if deps[0] != "dep1" || deps[1] != "dep2" {
		t.Errorf("Dependencies not set correctly: %v", deps)
	}

	// Verify we can use the component normally
	if !depAwareComp.HasDependency("dep1") {
		t.Error("Component should have dependency 'dep1'")
	}

	// Add a new dependency
	depAwareComp.AddDependency("dep3")
	if !depAwareComp.HasDependency("dep3") {
		t.Error("Component should have dependency 'dep3' after adding it")
	}
}

func TestDependencyAwareResolveDependencyErrors(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create a registry
	registry := NewRegistry(DefaultRegistryOptions{
		Logger: mockLogger,
	})

	// Create a dependency-aware component
	baseComp := NewBaseComponent("test-id", "Test Component", TypeBasic)
	depAwareComp := NewDependencyAwareComponentWithOptions(DependencyAwareComponentOptions{
		Base:         baseComp,
		Dependencies: []string{"dep1"},
		Logger:       mockLogger,
	})

	// Test case: Dependency not registered in registry
	_, err := depAwareComp.ResolveDependency("dep1", registry)
	if err == nil {
		t.Error("Expected error resolving unregistered dependency, but got nil")
	}

	// Verify it's a dependency not found error
	if !IsComponentError(err, ErrDependencyNotFound) {
		t.Errorf("Expected dependency not found error, but got: %s", err)
	}
}
