package mocks

import (
	"github.com/ebanfa/skeleton/internal/domain/component"
)

// MockPlugin implements the plugin.Plugin interface for testing
type MockPlugin struct {
	// Function fields for controlling behavior
	IDFunc         func() string
	VersionFunc    func() string
	LoadFunc       func(ctx component.Context, registry component.Registry) error
	UnloadFunc     func(ctx component.Context) error
	ComponentsFunc func() []component.Component

	// Call tracking for verification
	IDCalls         int
	VersionCalls    int
	LoadCalls       []component.Registry
	UnloadCalls     int
	ComponentsCalls int
}

// ID mocks the ID method
func (m *MockPlugin) ID() string {
	m.IDCalls++
	if m.IDFunc != nil {
		return m.IDFunc()
	}
	return "mock-plugin-id"
}

// Version mocks the Version method
func (m *MockPlugin) Version() string {
	m.VersionCalls++
	if m.VersionFunc != nil {
		return m.VersionFunc()
	}
	return "1.0.0"
}

// Load mocks the Load method
func (m *MockPlugin) Load(ctx component.Context, registry component.Registry) error {
	m.LoadCalls = append(m.LoadCalls, registry)
	if m.LoadFunc != nil {
		return m.LoadFunc(ctx, registry)
	}
	return nil
}

// Unload mocks the Unload method
func (m *MockPlugin) Unload(ctx component.Context) error {
	m.UnloadCalls++
	if m.UnloadFunc != nil {
		return m.UnloadFunc(ctx)
	}
	return nil
}

// Components mocks the Components method
func (m *MockPlugin) Components() []component.Component {
	m.ComponentsCalls++
	if m.ComponentsFunc != nil {
		return m.ComponentsFunc()
	}
	return []component.Component{}
}
