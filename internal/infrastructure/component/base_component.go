package component

import (
	"sync"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/system"
)

// BaseComponent provides a concrete implementation of the Component interface.
type BaseComponent struct {
	config   component.ComponentConfig
	metadata component.Metadata
	mu       sync.RWMutex
}

// NewBaseComponent creates a new base component instance.
// This constructor accepts configuration and minimal dependencies to keep it simple and focused.
func NewBaseComponent(config component.ComponentConfig) component.Component {
	return &BaseComponent{
		config:   config,
		metadata: make(component.Metadata),
	}
}

// ID returns the component's unique identifier.
func (c *BaseComponent) ID() string {
	return c.config.ID
}

// Name returns the component's human-readable name.
func (c *BaseComponent) Name() string {
	return c.config.Name
}

// Description returns the component's description.
func (c *BaseComponent) Description() string {
	return c.config.Description
}

// Version returns the component's version.
func (c *BaseComponent) Version() string {
	if c.config.Version == "" {
		return "1.0.0" // Default version
	}
	return c.config.Version
}

// Type returns the component's type.
func (c *BaseComponent) Type() component.ComponentType {
	return component.ComponentType(c.config.Type)
}

// Metadata returns a copy of the component's metadata to prevent concurrent map access.
func (c *BaseComponent) Metadata() component.Metadata {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Return a copy to prevent concurrent map access
	metadata := make(component.Metadata, len(c.metadata))
	for k, v := range c.metadata {
		metadata[k] = v
	}

	// Include configuration properties in metadata
	for k, v := range c.config.Properties {
		metadata[k] = v
	}

	// Include dependencies in metadata
	if len(c.config.Dependencies) > 0 {
		metadata["dependencies"] = c.config.Dependencies
	}

	return metadata
}

// SetMetadata sets a metadata key-value pair (thread-safe).
func (c *BaseComponent) SetMetadata(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.metadata[key] = value
}

// Initialize initializes the component with the given context and system.
// This follows the exact interface signature: Initialize(ctx context.Context, system sys.System) error
func (c *BaseComponent) Initialize(ctx context.Context, system system.System) error {
	// Set initialization metadata
	c.SetMetadata("initialized", true)
	c.SetMetadata("system_available", system != nil)

	// Basic components don't need complex initialization
	// Subclasses can override this method for more complex initialization logic
	return nil
}

// Dispose releases resources used by the component.
func (c *BaseComponent) Dispose() error {
	// Set disposal metadata
	c.SetMetadata("disposed", true)

	// Clear metadata to help with garbage collection
	c.mu.Lock()
	defer c.mu.Unlock()
	c.metadata = make(component.Metadata)

	return nil
}

// GetConfig returns the component's configuration (for subclasses).
func (c *BaseComponent) GetConfig() component.ComponentConfig {
	return c.config
}

// Dependencies returns the component's dependencies.
func (c *BaseComponent) Dependencies() []string {
	return c.config.Dependencies
}

// HasDependency checks if the component has a specific dependency.
func (c *BaseComponent) HasDependency(id string) bool {
	for _, dep := range c.config.Dependencies {
		if dep == id {
			return true
		}
	}
	return false
}
