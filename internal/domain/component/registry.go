package component

import (
	"sync"

	"github.com/ebanfa/skeleton/internal/infrastructure/logging"
)

// Registry manages component registration and discovery.
type Registry interface {
	// Registration
	Register(Component) error
	Unregister(id string) error

	// Discovery
	Get(id string) (Component, error)
	FindByType(componentType ComponentType) []Component
	FindByMetadata(key string, value interface{}) []Component

	// Lifecycle
	Initialize(ctx Context) error
	Shutdown() error
}

// DefaultRegistry provides a standard implementation of the Registry interface.
type DefaultRegistry struct {
	components map[string]Component
	factories  map[string]Factory
	mu         sync.RWMutex // Protects components and factories
	logger     logging.Logger
}

// DefaultRegistryOptions contains options for creating a DefaultRegistry.
type DefaultRegistryOptions struct {
	Logger logging.Logger
}

// NewRegistry creates a new component registry.
func NewRegistry(options DefaultRegistryOptions) Registry {
	return &DefaultRegistry{
		components: make(map[string]Component),
		factories:  make(map[string]Factory),
		logger:     options.Logger,
	}
}

// CreateRegistry is a factory method for backward compatibility.
// Creates a registry with default dependencies.
func CreateRegistry() Registry {
	return NewRegistry(DefaultRegistryOptions{
		Logger: logging.CreateStandardLogger(logging.Info),
	})
}

// Register adds a component to the registry.
func (r *DefaultRegistry) Register(component Component) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := component.ID()

	// Check if a component with this ID already exists
	if _, exists := r.components[id]; exists {
		return NewError(ErrComponentExists, "component with ID already exists", nil).
			WithDetail("id", id)
	}

	r.components[id] = component
	return nil
}

// Unregister removes a component from the registry.
func (r *DefaultRegistry) Unregister(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if the component exists
	if _, exists := r.components[id]; !exists {
		return NewError(ErrComponentNotFound, "component not found", nil).
			WithDetail("id", id)
	}

	delete(r.components, id)
	return nil
}

// Get retrieves a component by ID.
func (r *DefaultRegistry) Get(id string) (Component, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	component, exists := r.components[id]
	if !exists {
		return nil, NewError(ErrComponentNotFound, "component not found", nil).
			WithDetail("id", id)
	}

	return component, nil
}

// FindByType finds all components of a specific type.
func (r *DefaultRegistry) FindByType(componentType ComponentType) []Component {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Component

	for _, comp := range r.components {
		if comp.Type() == componentType {
			result = append(result, comp)
		}
	}

	return result
}

// FindByMetadata finds all components with a specific metadata key-value pair.
func (r *DefaultRegistry) FindByMetadata(key string, value interface{}) []Component {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Component

	for _, comp := range r.components {
		metadata := comp.Metadata()
		if val, exists := metadata[key]; exists && val == value {
			result = append(result, comp)
		}
	}

	return result
}

// RegisterFactory registers a component factory.
func (r *DefaultRegistry) RegisterFactory(id string, factory Factory) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.factories[id]; exists {
		return NewError(ErrComponentExists, "factory with ID already exists", nil).
			WithDetail("id", id)
	}

	r.factories[id] = factory
	return nil
}

// GetFactory retrieves a factory by ID.
func (r *DefaultRegistry) GetFactory(id string) (Factory, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	factory, exists := r.factories[id]
	if !exists {
		return nil, NewError(ErrComponentNotFound, "factory not found", nil).
			WithDetail("id", id)
	}

	return factory, nil
}

// Initialize initializes all registered components.
func (r *DefaultRegistry) Initialize(ctx Context) error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// In a real implementation, we would need to sort by dependencies
	// and initialize in the correct order. This is a simplified version.
	for id, comp := range r.components {
		if err := comp.Initialize(ctx); err != nil {
			return NewError(ErrInitializationFailed, "failed to initialize component", err).
				WithDetail("id", id)
		}
	}

	return nil
}

// Shutdown disposes all registered components.
func (r *DefaultRegistry) Shutdown() error {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// In a real implementation, we would need to sort by dependencies
	// and dispose in the reverse order. This is a simplified version.
	for id, comp := range r.components {
		if err := comp.Dispose(); err != nil {
			return NewError(ErrDisposalFailed, "failed to dispose component", err).
				WithDetail("id", id)
		}
	}

	return nil
}
