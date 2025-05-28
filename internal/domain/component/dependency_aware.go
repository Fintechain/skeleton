package component

import (
	"github.com/fintechain/skeleton/internal/infrastructure/logging"
)

// DependencyAware is an interface for components that have dependencies on other components.
type DependencyAware interface {
	// Dependencies returns the IDs of components this component depends on.
	Dependencies() []string

	// AddDependency adds a component dependency.
	AddDependency(id string)

	// RemoveDependency removes a component dependency.
	RemoveDependency(id string)

	// HasDependency checks if this component depends on the component with the given ID.
	HasDependency(id string) bool

	// ResolveDependency resolves a dependency to a component instance.
	ResolveDependency(id string, registry Registry) (Component, error)

	// ResolveDependencies resolves all dependencies to component instances.
	ResolveDependencies(registry Registry) (map[string]Component, error)
}

// DependencyAwareComponent is a component that is aware of its dependencies.
type DependencyAwareComponent struct {
	Component
	dependencies []string
	logger       logging.Logger
}

// DependencyAwareComponentOptions contains options for creating a DependencyAwareComponent.
type DependencyAwareComponentOptions struct {
	Base         Component
	Dependencies []string
	Logger       logging.Logger
}

// NewDependencyAwareComponent creates a new dependency-aware component.
func NewDependencyAwareComponent(base Component, dependencies []string) *DependencyAwareComponent {
	return &DependencyAwareComponent{
		Component:    base,
		dependencies: dependencies,
	}
}

// NewDependencyAwareComponentWithOptions creates a new dependency-aware component with injected dependencies.
func NewDependencyAwareComponentWithOptions(options DependencyAwareComponentOptions) *DependencyAwareComponent {
	return &DependencyAwareComponent{
		Component:    options.Base,
		dependencies: options.Dependencies,
		logger:       options.Logger,
	}
}

// Dependencies returns the IDs of components this component depends on.
func (c *DependencyAwareComponent) Dependencies() []string {
	return c.dependencies
}

// AddDependency adds a component dependency.
func (c *DependencyAwareComponent) AddDependency(id string) {
	// Check if the dependency already exists
	if c.HasDependency(id) {
		return
	}

	if c.logger != nil {
		c.logger.Debug("Adding dependency %s to component %s", id, c.ID())
	}

	c.dependencies = append(c.dependencies, id)
}

// RemoveDependency removes a component dependency.
func (c *DependencyAwareComponent) RemoveDependency(id string) {
	for i, dep := range c.dependencies {
		if dep == id {
			if c.logger != nil {
				c.logger.Debug("Removing dependency %s from component %s", id, c.ID())
			}

			// Remove the dependency by replacing it with the last one
			// and reducing the slice length (more efficient than creating a new slice)
			lastIdx := len(c.dependencies) - 1
			c.dependencies[i] = c.dependencies[lastIdx]
			c.dependencies = c.dependencies[:lastIdx]
			return
		}
	}
}

// HasDependency checks if this component depends on the component with the given ID.
func (c *DependencyAwareComponent) HasDependency(id string) bool {
	for _, dep := range c.dependencies {
		if dep == id {
			return true
		}
	}
	return false
}

// ResolveDependency resolves a dependency to a component instance.
func (c *DependencyAwareComponent) ResolveDependency(id string, registry Registry) (Component, error) {
	if !c.HasDependency(id) {
		errMsg := "not a dependency of this component"
		if c.logger != nil {
			c.logger.Error(errMsg+" (component_id=%s, dependency_id=%s)", c.ID(), id)
		}

		return nil, NewError(
			ErrDependencyNotFound,
			errMsg,
			nil,
		).WithDetail("component_id", c.ID()).WithDetail("dependency_id", id)
	}

	dep, err := registry.Get(id)
	if err != nil {
		errMsg := "dependency not found in registry"
		if c.logger != nil {
			c.logger.Error(errMsg+" (component_id=%s, dependency_id=%s): %s", c.ID(), id, err)
		}

		return nil, NewError(
			ErrDependencyNotFound,
			errMsg,
			err,
		).WithDetail("component_id", c.ID()).WithDetail("dependency_id", id)
	}

	return dep, nil
}

// ResolveDependencies resolves all dependencies to component instances.
func (c *DependencyAwareComponent) ResolveDependencies(registry Registry) (map[string]Component, error) {
	resolved := make(map[string]Component)

	for _, id := range c.dependencies {
		comp, err := registry.Get(id)
		if err != nil {
			errMsg := "dependency not found in registry"
			if c.logger != nil {
				c.logger.Error(errMsg+" (component_id=%s, dependency_id=%s): %s", c.ID(), id, err)
			}

			return nil, NewError(
				ErrDependencyNotFound,
				errMsg,
				err,
			).WithDetail("component_id", c.ID()).WithDetail("dependency_id", id)
		}
		resolved[id] = comp
	}

	return resolved, nil
}
