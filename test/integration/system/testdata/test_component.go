package testdata

import (
	"github.com/ebanfa/skeleton/internal/domain/component"
)

// TestComponent is a test component implementation for integration testing
type TestComponent struct {
	id            string
	name          string
	componentType component.ComponentType
	metadata      component.Metadata
	initialized   bool
	disposed      bool
}

// NewTestComponent creates a new test component with the given ID
func NewTestComponent(id string) *TestComponent {
	return &TestComponent{
		id:            id,
		name:          "Test Component " + id,
		componentType: component.TypeBasic,
		metadata:      make(component.Metadata),
	}
}

// ID returns the component ID
func (c *TestComponent) ID() string {
	return c.id
}

// Name returns the component name
func (c *TestComponent) Name() string {
	return c.name
}

// Type returns the component type
func (c *TestComponent) Type() component.ComponentType {
	return c.componentType
}

// Metadata returns the component metadata
func (c *TestComponent) Metadata() component.Metadata {
	return c.metadata
}

// Initialize initializes the component
func (c *TestComponent) Initialize(ctx component.Context) error {
	c.initialized = true
	c.metadata["initialized"] = true
	return nil
}

// Dispose disposes the component
func (c *TestComponent) Dispose() error {
	c.disposed = true
	c.metadata["disposed"] = true
	return nil
}

// Test helper methods

// IsInitialized returns whether the component is initialized
func (c *TestComponent) IsInitialized() bool {
	return c.initialized
}

// IsDisposed returns whether the component is disposed
func (c *TestComponent) IsDisposed() bool {
	return c.disposed
}

// SetMetadata sets a metadata value
func (c *TestComponent) SetMetadata(key string, value interface{}) {
	c.metadata[key] = value
}

// GetMetadata gets a metadata value
func (c *TestComponent) GetMetadata(key string) (interface{}, bool) {
	value, exists := c.metadata[key]
	return value, exists
}
