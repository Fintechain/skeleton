package mocks

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/plugin"
)

// MockPluginManager is a mock implementation of plugin.PluginManager for testing
type MockPluginManager struct {
	// Function fields for customizing behavior
	DiscoverFunc    func(component.Context, string) ([]plugin.PluginInfo, error)
	LoadFunc        func(component.Context, string, component.Registry) error
	UnloadFunc      func(component.Context, string) error
	ListPluginsFunc func() []plugin.PluginInfo
	GetPluginFunc   func(string) (plugin.Plugin, error)

	// Call tracking
	DiscoverCalls    []DiscoverCall
	LoadCalls        []LoadCall
	UnloadCalls      []UnloadCall
	ListPluginsCalls int
	GetPluginCalls   []string

	// State
	Plugins       map[string]plugin.Plugin
	LoadedPlugins map[string]bool
}

type DiscoverCall struct {
	Ctx      component.Context
	Location string
}

type LoadCall struct {
	Ctx      component.Context
	ID       string
	Registry component.Registry
}

type UnloadCall struct {
	Ctx component.Context
	ID  string
}

// NewMockPluginManager creates a new mock plugin manager
func NewMockPluginManager() *MockPluginManager {
	return &MockPluginManager{
		Plugins:       make(map[string]plugin.Plugin),
		LoadedPlugins: make(map[string]bool),
	}
}

// Discover implements plugin.PluginManager
func (m *MockPluginManager) Discover(ctx component.Context, location string) ([]plugin.PluginInfo, error) {
	m.DiscoverCalls = append(m.DiscoverCalls, DiscoverCall{Ctx: ctx, Location: location})
	if m.DiscoverFunc != nil {
		return m.DiscoverFunc(ctx, location)
	}
	return []plugin.PluginInfo{}, nil
}

// Load implements plugin.PluginManager
func (m *MockPluginManager) Load(ctx component.Context, id string, registry component.Registry) error {
	m.LoadCalls = append(m.LoadCalls, LoadCall{Ctx: ctx, ID: id, Registry: registry})
	if m.LoadFunc != nil {
		return m.LoadFunc(ctx, id, registry)
	}
	if m.LoadedPlugins != nil {
		m.LoadedPlugins[id] = true
	}
	return nil
}

// Unload implements plugin.PluginManager
func (m *MockPluginManager) Unload(ctx component.Context, id string) error {
	m.UnloadCalls = append(m.UnloadCalls, UnloadCall{Ctx: ctx, ID: id})
	if m.UnloadFunc != nil {
		return m.UnloadFunc(ctx, id)
	}
	if m.LoadedPlugins != nil {
		delete(m.LoadedPlugins, id)
	}
	return nil
}

// ListPlugins implements plugin.PluginManager
func (m *MockPluginManager) ListPlugins() []plugin.PluginInfo {
	m.ListPluginsCalls++
	if m.ListPluginsFunc != nil {
		return m.ListPluginsFunc()
	}
	var result []plugin.PluginInfo
	if m.Plugins != nil {
		for _, p := range m.Plugins {
			result = append(result, plugin.PluginInfo{
				ID:      p.ID(),
				Version: p.Version(),
			})
		}
	}
	return result
}

// GetPlugin implements plugin.PluginManager
func (m *MockPluginManager) GetPlugin(id string) (plugin.Plugin, error) {
	m.GetPluginCalls = append(m.GetPluginCalls, id)
	if m.GetPluginFunc != nil {
		return m.GetPluginFunc(id)
	}
	if m.Plugins != nil {
		if p, exists := m.Plugins[id]; exists {
			return p, nil
		}
	}
	return nil, component.NewError(plugin.ErrPluginNotFound, "plugin not found", nil)
}

// RegisterPlugin is a helper method for testing (not part of the interface)
func (m *MockPluginManager) RegisterPlugin(p plugin.Plugin) error {
	if m.Plugins != nil {
		m.Plugins[p.ID()] = p
	}
	return nil
}
