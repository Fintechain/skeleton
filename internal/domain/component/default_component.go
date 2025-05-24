package component

import (
	"github.com/ebanfa/skeleton/internal/infrastructure/logging"
)

// DefaultComponent is a concrete implementation of the Component interface
// built on top of BaseComponent.
type DefaultComponent struct {
	*BaseComponent
	description string
	initialized bool
	disposed    bool
	logger      logging.Logger
}

// DefaultComponentOptions contains options for creating a DefaultComponent.
type DefaultComponentOptions struct {
	ID          string
	Name        string
	Type        ComponentType
	Description string
	Logger      logging.Logger
}

// NewDefaultComponent creates a new DefaultComponent with dependency injection.
func NewDefaultComponent(options DefaultComponentOptions) *DefaultComponent {
	return &DefaultComponent{
		BaseComponent: NewBaseComponent(options.ID, options.Name, options.Type),
		description:   options.Description,
		initialized:   false,
		disposed:      false,
		logger:        options.Logger,
	}
}

// CreateDefaultComponent is a factory method for backward compatibility.
// Creates a component with default dependencies.
func CreateDefaultComponent(id, name string, componentType ComponentType, description string) *DefaultComponent {
	return NewDefaultComponent(DefaultComponentOptions{
		ID:          id,
		Name:        name,
		Type:        componentType,
		Description: description,
		Logger:      logging.CreateStandardLogger(logging.Info),
	})
}

// Initialize initializes the component.
func (c *DefaultComponent) Initialize(ctx Context) error {
	if c.initialized {
		return nil
	}

	// Store the description in metadata
	c.SetMetadata("description", c.description)

	if c.logger != nil {
		c.logger.Debug("Initializing component %s (%s)", c.Name(), c.ID())
	}

	// Mark as initialized
	c.initialized = true
	return nil
}

// Dispose releases resources used by the component.
func (c *DefaultComponent) Dispose() error {
	if c.disposed {
		return nil
	}

	if c.logger != nil {
		c.logger.Debug("Disposing component %s (%s)", c.Name(), c.ID())
	}

	// Clean up resources (if any)

	// Mark as disposed
	c.disposed = true
	return nil
}

// IsInitialized returns whether the component has been initialized.
func (c *DefaultComponent) IsInitialized() bool {
	return c.initialized
}

// IsDisposed returns whether the component has been disposed.
func (c *DefaultComponent) IsDisposed() bool {
	return c.disposed
}

// Description returns the component's description.
func (c *DefaultComponent) Description() string {
	return c.description
}
