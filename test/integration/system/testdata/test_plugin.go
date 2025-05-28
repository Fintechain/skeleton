package testdata

import (
	"errors"
	"time"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/plugin"
)

// TestPlugin is a test plugin implementation for integration testing
type TestPlugin struct {
	id         string
	version    string
	components []component.Component
	loaded     bool
	loadError  error
	loadDelay  time.Duration
}

// NewTestPlugin creates a new test plugin with the given ID and version
func NewTestPlugin(id, version string) plugin.Plugin {
	return &TestPlugin{
		id:      id,
		version: version,
		components: []component.Component{
			NewTestComponent(id + "-component"),
		},
	}
}

// ID returns the plugin ID
func (p *TestPlugin) ID() string {
	return p.id
}

// Version returns the plugin version
func (p *TestPlugin) Version() string {
	return p.version
}

// Load loads the plugin and registers its components
func (p *TestPlugin) Load(ctx component.Context, registry component.Registry) error {
	// Simulate load delay if configured
	if p.loadDelay > 0 {
		time.Sleep(p.loadDelay)
	}

	// Return error if configured to fail
	if p.loadError != nil {
		return p.loadError
	}

	// Register components with registry
	for _, comp := range p.components {
		if err := registry.Register(comp); err != nil {
			return err
		}
	}

	p.loaded = true
	return nil
}

// Unload unloads the plugin
func (p *TestPlugin) Unload(ctx component.Context) error {
	p.loaded = false
	return nil
}

// Components returns the plugin's components
func (p *TestPlugin) Components() []component.Component {
	return p.components
}

// Test helper methods

// IsLoaded returns whether the plugin is currently loaded
func (p *TestPlugin) IsLoaded() bool {
	return p.loaded
}

// SetLoadError configures the plugin to fail loading with the given error
func (p *TestPlugin) SetLoadError(err error) {
	p.loadError = err
}

// SetLoadDelay configures the plugin to delay loading by the given duration
func (p *TestPlugin) SetLoadDelay(delay time.Duration) {
	p.loadDelay = delay
}

// CreateFailingPlugin creates a plugin that will fail to load
func CreateFailingPlugin(id, version string, errorMsg string) plugin.Plugin {
	plugin := &TestPlugin{
		id:      id,
		version: version,
		components: []component.Component{
			NewTestComponent(id + "-component"),
		},
	}
	plugin.SetLoadError(errors.New(errorMsg))
	return plugin
}

// CreateSlowPlugin creates a plugin that takes time to load
func CreateSlowPlugin(id, version string, delay time.Duration) plugin.Plugin {
	plugin := &TestPlugin{
		id:      id,
		version: version,
		components: []component.Component{
			NewTestComponent(id + "-component"),
		},
	}
	plugin.SetLoadDelay(delay)
	return plugin
}
