package component

import (
	"fmt"
	"sync"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/registry"
)

// DependencyAwareComponentImpl provides a concrete implementation of the DependencyAwareComponent interface.
type DependencyAwareComponentImpl struct {
	component.Component
	registry     registry.Registry
	dependencies []string
	mu           sync.RWMutex
}

// NewDependencyAwareComponent creates a new dependency-aware component instance.
// This constructor accepts a base component and registry interface dependencies.
func NewDependencyAwareComponent(baseComponent component.Component, registry registry.Registry) component.DependencyAwareComponent {
	// Initialize with dependencies from the base component if available
	var initialDeps []string
	if baseComponent != nil {
		// Check if the base component has a Dependencies method (like our BaseComponent)
		if depProvider, ok := baseComponent.(interface{ Dependencies() []string }); ok {
			initialDeps = depProvider.Dependencies()
		}
	}

	return &DependencyAwareComponentImpl{
		Component:    baseComponent,
		registry:     registry,
		dependencies: append([]string{}, initialDeps...), // Copy to avoid shared slice
	}
}

// Dependencies returns the IDs of components this component depends on.
func (d *DependencyAwareComponentImpl) Dependencies() []string {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Return a copy to prevent external modification
	deps := make([]string, len(d.dependencies))
	copy(deps, d.dependencies)
	return deps
}

// AddDependency adds a component dependency.
func (d *DependencyAwareComponentImpl) AddDependency(id string) {
	if id == "" {
		return
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	// Check if dependency already exists
	for _, dep := range d.dependencies {
		if dep == id {
			return // Already exists
		}
	}

	d.dependencies = append(d.dependencies, id)
}

// RemoveDependency removes a component dependency.
func (d *DependencyAwareComponentImpl) RemoveDependency(id string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for i, dep := range d.dependencies {
		if dep == id {
			// Remove by swapping with last element and truncating
			d.dependencies[i] = d.dependencies[len(d.dependencies)-1]
			d.dependencies = d.dependencies[:len(d.dependencies)-1]
			return
		}
	}
}

// HasDependency checks if this component depends on the component with the given ID.
func (d *DependencyAwareComponentImpl) HasDependency(id string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	for _, dep := range d.dependencies {
		if dep == id {
			return true
		}
	}
	return false
}

// ResolveDependency resolves a dependency to a component instance.
func (d *DependencyAwareComponentImpl) ResolveDependency(id string, registrar registry.Registry) (component.Component, error) {
	if registrar == nil {
		return nil, fmt.Errorf("registry is required for dependency resolution")
	}

	// Check for circular dependency
	if err := d.checkCircularDependency(id, registrar, make(map[string]bool)); err != nil {
		return nil, err
	}

	// Get the component from registry
	item, err := registrar.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve dependency '%s': %w", id, err)
	}

	// Ensure it's a component
	comp, ok := item.(component.Component)
	if !ok {
		return nil, fmt.Errorf("dependency '%s' is not a component", id)
	}

	return comp, nil
}

// ResolveDependencies resolves all dependencies to component instances.
func (d *DependencyAwareComponentImpl) ResolveDependencies(registrar registry.Registry) (map[string]component.Component, error) {
	if registrar == nil {
		return nil, fmt.Errorf("registry is required for dependency resolution")
	}

	d.mu.RLock()
	deps := make([]string, len(d.dependencies))
	copy(deps, d.dependencies)
	d.mu.RUnlock()

	resolved := make(map[string]component.Component)

	for _, depID := range deps {
		comp, err := d.ResolveDependency(depID, registrar)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve dependencies: %w", err)
		}
		resolved[depID] = comp
	}

	return resolved, nil
}

// checkCircularDependency performs circular dependency detection using depth-first search.
func (d *DependencyAwareComponentImpl) checkCircularDependency(targetID string, registrar registry.Registry, visited map[string]bool) error {
	// If we're checking our own ID, we have a circular dependency
	if d.Component != nil && targetID == d.Component.ID() {
		return fmt.Errorf("circular dependency detected: component '%s' depends on itself", targetID)
	}

	// If we've already visited this component in this path, we have a cycle
	if visited[targetID] {
		return fmt.Errorf("circular dependency detected involving component '%s'", targetID)
	}

	// Mark as visited for this path
	visited[targetID] = true
	defer func() {
		delete(visited, targetID)
	}()

	// Get the target component
	item, err := registrar.Get(targetID)
	if err != nil {
		// If component doesn't exist, no circular dependency possible
		return nil
	}

	// Check if target component is also dependency-aware
	if depAware, ok := item.(component.DependencyAwareComponent); ok {
		// Check all of its dependencies
		for _, depID := range depAware.Dependencies() {
			if err := d.checkCircularDependency(depID, registrar, visited); err != nil {
				return err
			}
		}
	}

	return nil
}
