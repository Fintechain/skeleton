package plugin

import (
	"fmt"
	"sync"

	"github.com/fintechain/skeleton/internal/domain/component"
)

// DefaultPlugin provides a standard implementation of the Plugin interface.
type DefaultPlugin struct {
	id          string
	name        string
	version     string
	description string
	author      string
	components  []component.Component
	metadata    map[string]interface{}
	isLoaded    bool
	mu          sync.RWMutex
}

// NewDefaultPlugin creates a new default plugin.
func NewDefaultPlugin(
	id string,
	version string,
	name string,
	description string,
	author string,
) *DefaultPlugin {
	return &DefaultPlugin{
		id:          id,
		version:     version,
		name:        name,
		description: description,
		author:      author,
		components:  make([]component.Component, 0),
		metadata:    make(map[string]interface{}),
		isLoaded:    false,
	}
}

// ID returns the plugin's ID.
func (p *DefaultPlugin) ID() string {
	return p.id
}

// Version returns the plugin's version.
func (p *DefaultPlugin) Version() string {
	return p.version
}

// Components returns the components provided by this plugin.
func (p *DefaultPlugin) Components() []component.Component {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Return a copy to prevent modification
	components := make([]component.Component, len(p.components))
	copy(components, p.components)
	return components
}

// AddComponent adds a component to the plugin.
func (p *DefaultPlugin) AddComponent(comp component.Component) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.components = append(p.components, comp)
}

// GetInfo returns the plugin information.
func (p *DefaultPlugin) GetInfo() PluginInfo {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Copy metadata
	metadata := make(map[string]interface{})
	for k, v := range p.metadata {
		metadata[k] = v
	}

	return PluginInfo{
		ID:          p.id,
		Name:        p.name,
		Version:     p.version,
		Description: p.description,
		Author:      p.author,
		Metadata:    metadata,
	}
}

// Load registers the plugin and its components with the registry.
func (p *DefaultPlugin) Load(ctx component.Context, registry component.Registry) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.isLoaded {
		return fmt.Errorf("plugin %s is already loaded", p.id)
	}

	// Register each component
	for _, comp := range p.components {
		err := registry.Register(comp)
		if err != nil {
			// If registration fails, unregister any components we've already registered
			for _, c := range p.components {
				if c.ID() == comp.ID() {
					break
				}
				_ = registry.Unregister(c.ID())
			}
			return component.NewError(
				ErrPluginLoad,
				"failed to register component",
				err,
			).WithDetail("plugin_id", p.id).WithDetail("component_id", comp.ID())
		}

		// Initialize the component
		err = comp.Initialize(ctx)
		if err != nil {
			// If initialization fails, unregister the component and any previous ones
			_ = registry.Unregister(comp.ID())
			for _, c := range p.components {
				if c.ID() == comp.ID() {
					break
				}
				_ = registry.Unregister(c.ID())
			}
			return component.NewError(
				ErrPluginLoad,
				"failed to initialize component",
				err,
			).WithDetail("plugin_id", p.id).WithDetail("component_id", comp.ID())
		}
	}

	p.isLoaded = true
	return nil
}

// Unload unregisters the plugin and its components from the registry.
func (p *DefaultPlugin) Unload(ctx component.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if !p.isLoaded {
		return fmt.Errorf("plugin %s is not loaded", p.id)
	}

	// Unregister components in reverse order
	for i := len(p.components) - 1; i >= 0; i-- {
		comp := p.components[i]

		// Dispose the component
		err := comp.Dispose()
		if err != nil {
			return component.NewError(
				ErrPluginUnload,
				"failed to dispose component",
				err,
			).WithDetail("plugin_id", p.id).WithDetail("component_id", comp.ID())
		}
	}

	p.isLoaded = false
	return nil
}

// IsLoaded returns whether the plugin is loaded.
func (p *DefaultPlugin) IsLoaded() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.isLoaded
}

// SetMetadata sets a metadata key-value pair.
func (p *DefaultPlugin) SetMetadata(key string, value interface{}) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.metadata[key] = value
}
