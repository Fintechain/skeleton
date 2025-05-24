package mocks

import (
	"github.com/ebanfa/skeleton/internal/domain/component"
)

type PluginLoadCall struct {
	Ctx      component.Context
	Registry component.Registry
}

// MockPlugin is a mock implementation of plugin.Plugin for testing
type MockPlugin struct {
	// Function fields for customizing behavior
	IDFunc         func() string
	VersionFunc    func() string
	LoadFunc       func(component.Context, component.Registry) error
	UnloadFunc     func(component.Context) error
	ComponentsFunc func() []component.Component

	// Call tracking
	IDCalls         int
	VersionCalls    int
	LoadCalls       []PluginLoadCall
	UnloadCalls     []component.Context
	ComponentsCalls int

	// State
	PluginID         string
	PluginVersion    string
	PluginComponents []component.Component
	IsLoaded         bool
}

// NewMockPlugin creates a new mock plugin
func NewMockPlugin(id, version string) *MockPlugin {
	return &MockPlugin{
		PluginID:      id,
		PluginVersion: version,
	}
}

// ID implements plugin.Plugin
func (m *MockPlugin) ID() string {
	m.IDCalls++
	if m.IDFunc != nil {
		return m.IDFunc()
	}
	return m.PluginID
}

// Version implements plugin.Plugin
func (m *MockPlugin) Version() string {
	m.VersionCalls++
	if m.VersionFunc != nil {
		return m.VersionFunc()
	}
	return m.PluginVersion
}

// Load implements plugin.Plugin
func (m *MockPlugin) Load(ctx component.Context, registry component.Registry) error {
	m.LoadCalls = append(m.LoadCalls, PluginLoadCall{Ctx: ctx, Registry: registry})
	if m.LoadFunc != nil {
		err := m.LoadFunc(ctx, registry)
		if err == nil {
			m.IsLoaded = true
		}
		return err
	}
	m.IsLoaded = true
	return nil
}

// Unload implements plugin.Plugin
func (m *MockPlugin) Unload(ctx component.Context) error {
	m.UnloadCalls = append(m.UnloadCalls, ctx)
	if m.UnloadFunc != nil {
		err := m.UnloadFunc(ctx)
		if err == nil {
			m.IsLoaded = false
		}
		return err
	}
	m.IsLoaded = false
	return nil
}

// Components implements plugin.Plugin
func (m *MockPlugin) Components() []component.Component {
	m.ComponentsCalls++
	if m.ComponentsFunc != nil {
		return m.ComponentsFunc()
	}
	return m.PluginComponents
}

// SetComponents is a helper method for testing (not part of the interface)
func (m *MockPlugin) SetComponents(components []component.Component) {
	m.PluginComponents = components
}
