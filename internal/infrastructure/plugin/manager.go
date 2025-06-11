package plugin

import (
	"errors"
	"sync"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/plugin"
	infraComponent "github.com/fintechain/skeleton/internal/infrastructure/component"
)

// Manager implements the PluginManager interface.
type Manager struct {
	*infraComponent.BaseService
	plugins map[component.ComponentID]plugin.Plugin
	mu      sync.RWMutex
}

// NewManager creates a new plugin manager.
func NewManager(config component.ComponentConfig) *Manager {
	return &Manager{
		BaseService: infraComponent.NewBaseService(config),
		plugins:     make(map[component.ComponentID]plugin.Plugin),
	}
}

// Add adds a plugin to the manager.
func (m *Manager) Add(pluginID component.ComponentID, p plugin.Plugin) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.plugins[pluginID] = p
	return nil
}

// Remove removes a plugin from the manager.
func (m *Manager) Remove(pluginID component.ComponentID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.plugins[pluginID]; !exists {
		return errors.New(component.ErrComponentNotFound)
	}

	delete(m.plugins, pluginID)
	return nil
}

// StartPlugin starts a specific plugin.
func (m *Manager) StartPlugin(ctx context.Context, pluginID component.ComponentID) error {
	m.mu.RLock()
	p, exists := m.plugins[pluginID]
	m.mu.RUnlock()

	if !exists {
		return errors.New(component.ErrComponentNotFound)
	}

	return p.Start(ctx)
}

// StopPlugin stops a specific plugin.
func (m *Manager) StopPlugin(ctx context.Context, pluginID component.ComponentID) error {
	m.mu.RLock()
	p, exists := m.plugins[pluginID]
	m.mu.RUnlock()

	if !exists {
		return errors.New(component.ErrComponentNotFound)
	}

	return p.Stop(ctx)
}

// GetPlugin retrieves a plugin by ID.
func (m *Manager) GetPlugin(pluginID component.ComponentID) (plugin.Plugin, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	p, exists := m.plugins[pluginID]
	if !exists {
		return nil, errors.New(component.ErrComponentNotFound)
	}

	return p, nil
}

// ListPlugins returns all plugin IDs.
func (m *Manager) ListPlugins() []component.ComponentID {
	m.mu.RLock()
	defer m.mu.RUnlock()

	ids := make([]component.ComponentID, 0, len(m.plugins))
	for id := range m.plugins {
		ids = append(ids, id)
	}

	return ids
}

// Start starts the plugin manager and all registered plugins.
func (m *Manager) Start(ctx context.Context) error {
	if err := m.BaseService.Start(ctx); err != nil {
		return err
	}

	m.mu.RLock()
	plugins := make([]plugin.Plugin, 0, len(m.plugins))
	for _, p := range m.plugins {
		plugins = append(plugins, p)
	}
	m.mu.RUnlock()

	// Start all plugins
	for _, p := range plugins {
		if err := p.Start(ctx); err != nil {
			return err
		}
	}

	return nil
}

// Stop stops all plugins and the plugin manager.
func (m *Manager) Stop(ctx context.Context) error {
	m.mu.RLock()
	plugins := make([]plugin.Plugin, 0, len(m.plugins))
	for _, p := range m.plugins {
		plugins = append(plugins, p)
	}
	m.mu.RUnlock()

	// Stop all plugins
	for _, p := range plugins {
		if err := p.Stop(ctx); err != nil {
			return err
		}
	}

	return m.BaseService.Stop(ctx)
}

// Initialize initializes the plugin manager and all registered plugins.
func (m *Manager) Initialize(ctx context.Context, system component.System) error {
	if err := m.BaseService.Initialize(ctx, system); err != nil {
		return err
	}

	m.mu.RLock()
	plugins := make([]plugin.Plugin, 0, len(m.plugins))
	for _, p := range m.plugins {
		plugins = append(plugins, p)
	}
	m.mu.RUnlock()

	// Initialize all plugins
	for _, p := range plugins {
		if err := p.Initialize(ctx, system); err != nil {
			return err
		}
	}

	return nil
}
