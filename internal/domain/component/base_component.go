package component

import (
	"sync"

	"github.com/fintechain/skeleton/internal/infrastructure/logging"
)

// BaseComponent provides a basic implementation of the Component interface.
type BaseComponent struct {
	id       string
	name     string
	compType ComponentType
	metadata Metadata
	mu       sync.RWMutex // Protects metadata
	logger   logging.Logger
}

// BaseComponentOptions contains options for creating a BaseComponent.
type BaseComponentOptions struct {
	ID     string
	Name   string
	Type   ComponentType
	Logger logging.Logger
}

// NewBaseComponent creates a new base component with the given options.
func NewBaseComponent(id, name string, compType ComponentType) *BaseComponent {
	return &BaseComponent{
		id:       id,
		name:     name,
		compType: compType,
		metadata: make(Metadata),
	}
}

// NewBaseComponentWithOptions creates a new base component with dependency injection.
func NewBaseComponentWithOptions(options BaseComponentOptions) *BaseComponent {
	return &BaseComponent{
		id:       options.ID,
		name:     options.Name,
		compType: options.Type,
		metadata: make(Metadata),
		logger:   options.Logger,
	}
}

// ID returns the component's ID.
func (c *BaseComponent) ID() string {
	return c.id
}

// Name returns the component's name.
func (c *BaseComponent) Name() string {
	return c.name
}

// Type returns the component's type.
func (c *BaseComponent) Type() ComponentType {
	return c.compType
}

// Metadata returns the component's metadata.
func (c *BaseComponent) Metadata() Metadata {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Return a copy to prevent concurrent map access
	metadata := make(Metadata, len(c.metadata))
	for k, v := range c.metadata {
		metadata[k] = v
	}

	return metadata
}

// SetMetadata sets a metadata key-value pair.
func (c *BaseComponent) SetMetadata(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metadata[key] = value
}

// Initialize initializes the component.
func (c *BaseComponent) Initialize(ctx Context) error {
	// Base implementation does nothing
	return nil
}

// Dispose releases resources used by the component.
func (c *BaseComponent) Dispose() error {
	// Base implementation does nothing
	return nil
}
