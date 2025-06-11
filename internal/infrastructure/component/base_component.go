package component

import (
	"sync"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
)

// BaseComponent provides common component functionality that can be embedded
// in concrete component implementations.
type BaseComponent struct {
	id          component.ComponentID
	name        string
	description string
	version     string
	metadata    component.Metadata
	systemRef   component.System
	initialized bool
	mu          sync.RWMutex
}

// NewBaseComponent creates a new base component with the provided configuration.
func NewBaseComponent(config component.ComponentConfig) *BaseComponent {
	version := config.Version
	if version == "" {
		version = "1.0.0"
	}

	return &BaseComponent{
		id:          config.ID,
		name:        config.Name,
		description: config.Description,
		version:     version,
		metadata:    config.Properties,
	}
}

// ID returns the unique identifier for this component.
func (c *BaseComponent) ID() component.ComponentID {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.id
}

// Name returns a human-readable name for this component.
func (c *BaseComponent) Name() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.name
}

// Description returns a detailed description of the component's purpose.
func (c *BaseComponent) Description() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.description
}

// Version returns the version string for this component.
func (c *BaseComponent) Version() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.version
}

// Type returns the component type classification.
func (c *BaseComponent) Type() component.ComponentType {
	return component.TypeComponent
}

// Metadata returns component metadata as key-value pairs.
func (c *BaseComponent) Metadata() component.Metadata {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.metadata
}

// Initialize prepares the component for use within the system.
func (c *BaseComponent) Initialize(ctx context.Context, system component.System) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.initialized {
		return nil
	}

	c.systemRef = system
	c.initialized = true
	return nil
}

// Dispose cleans up component resources and prepares for shutdown.
func (c *BaseComponent) Dispose() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.systemRef = nil
	c.initialized = false
	return nil
}

// system returns the system reference (protected access).
func (c *BaseComponent) system() component.System {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.systemRef
}

// IsInitialized returns whether the component has been initialized.
func (c *BaseComponent) IsInitialized() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.initialized
}
