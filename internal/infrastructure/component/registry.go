package component

import (
	"errors"
	"sync"

	"github.com/fintechain/skeleton/internal/domain/component"
)

// Registry implements the component.Registry interface.
type Registry struct {
	components map[component.ComponentID]component.Component
	mu         sync.RWMutex
}

// NewRegistry creates a new component registry.
func NewRegistry() *Registry {
	return &Registry{
		components: make(map[component.ComponentID]component.Component),
	}
}

// Register adds a component to the registry.
func (r *Registry) Register(comp component.Component) error {
	if comp == nil {
		return errors.New(component.ErrInvalidItem)
	}

	id := comp.ID()
	if id == "" {
		return errors.New(component.ErrInvalidItem)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.components[id]; exists {
		return errors.New(component.ErrItemAlreadyExists)
	}

	r.components[id] = comp
	return nil
}

// Get retrieves a component by ID.
func (r *Registry) Get(id component.ComponentID) (component.Component, error) {
	if id == "" {
		return nil, errors.New(component.ErrInvalidItem)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	comp, exists := r.components[id]
	if !exists {
		return nil, errors.New(component.ErrItemNotFound)
	}

	return comp, nil
}

// GetByType retrieves all components of a given type.
func (r *Registry) GetByType(typ component.ComponentType) ([]component.Component, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []component.Component
	for _, comp := range r.components {
		if comp.Type() == typ {
			result = append(result, comp)
		}
	}

	return result, nil
}

// Find returns all components that match the given predicate function.
func (r *Registry) Find(predicate func(component.Component) bool) ([]component.Component, error) {
	if predicate == nil {
		return nil, errors.New(component.ErrInvalidItem)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []component.Component
	for _, comp := range r.components {
		if predicate(comp) {
			result = append(result, comp)
		}
	}

	return result, nil
}

// Has checks if a component exists.
func (r *Registry) Has(id component.ComponentID) bool {
	if id == "" {
		return false
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.components[id]
	return exists
}

// List returns all registered component IDs.
func (r *Registry) List() []component.ComponentID {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]component.ComponentID, 0, len(r.components))
	for id := range r.components {
		ids = append(ids, id)
	}

	return ids
}

// Unregister removes a component from the registry.
func (r *Registry) Unregister(id component.ComponentID) error {
	if id == "" {
		return errors.New(component.ErrInvalidItem)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.components[id]; !exists {
		return errors.New(component.ErrItemNotFound)
	}

	delete(r.components, id)
	return nil
}

// Count returns the number of registered components.
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.components)
}

// Clear removes all components from the registry.
func (r *Registry) Clear() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.components = make(map[component.ComponentID]component.Component)
	return nil
}
