package registry

import (
	"fmt"
	"sync"

	"github.com/fintechain/skeleton/internal/domain/registry"
)

// DefaultRegistry provides a thread-safe implementation of the Registry interface.
type DefaultRegistry struct {
	items map[string]registry.Identifiable
	mu    sync.RWMutex
}

// NewRegistry creates a new registry instance with minimal dependencies.
// This constructor accepts no dependencies to keep it simple and focused.
func NewRegistry() registry.Registry {
	return &DefaultRegistry{
		items: make(map[string]registry.Identifiable),
	}
}

// Register stores an item in the registry.
func (r *DefaultRegistry) Register(item registry.Identifiable) error {
	if item == nil {
		return fmt.Errorf(registry.ErrInvalidItem)
	}

	id := item.ID()
	if id == "" {
		return fmt.Errorf(registry.ErrInvalidItem)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; exists {
		return fmt.Errorf(registry.ErrItemAlreadyExists)
	}

	r.items[id] = item
	return nil
}

// Get retrieves an item by its ID.
func (r *DefaultRegistry) Get(id string) (registry.Identifiable, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	item, exists := r.items[id]
	if !exists {
		return nil, fmt.Errorf(registry.ErrItemNotFound)
	}

	return item, nil
}

// List returns all registered items.
func (r *DefaultRegistry) List() []registry.Identifiable {
	r.mu.RLock()
	defer r.mu.RUnlock()

	items := make([]registry.Identifiable, 0, len(r.items))
	for _, item := range r.items {
		items = append(items, item)
	}

	return items
}

// Remove removes an item from the registry.
func (r *DefaultRegistry) Remove(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.items[id]; !exists {
		return fmt.Errorf(registry.ErrItemNotFound)
	}

	delete(r.items, id)
	return nil
}

// Has checks if an item with the given ID exists.
func (r *DefaultRegistry) Has(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.items[id]
	return exists
}

// Count returns the number of registered items.
func (r *DefaultRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.items)
}

// Clear removes all items from the registry.
func (r *DefaultRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.items = make(map[string]registry.Identifiable)
}
