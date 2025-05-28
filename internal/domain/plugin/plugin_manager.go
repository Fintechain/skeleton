package plugin

import (
	"io/fs"
	"sync"

	"github.com/fintechain/skeleton/internal/domain/component"
)

// DefaultPluginManager provides a standard implementation of the PluginManager interface.
type DefaultPluginManager struct {
	plugins       map[string]Plugin
	loadedPlugins map[string]Plugin
	mu            sync.RWMutex
	fs            FileSystem
}

// NewPluginManager creates a new plugin manager with the given filesystem.
func NewPluginManager(fs FileSystem) *DefaultPluginManager {
	return &DefaultPluginManager{
		plugins:       make(map[string]Plugin),
		loadedPlugins: make(map[string]Plugin),
		fs:            fs,
	}
}

// CreatePluginManager is a factory method for backward compatibility.
// It creates a plugin manager with the standard file system implementation.
func CreatePluginManager() *DefaultPluginManager {
	return NewPluginManager(NewStandardFileSystem())
}

// RegisterPlugin registers a plugin with the manager.
func (m *DefaultPluginManager) RegisterPlugin(plugin Plugin) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := plugin.ID()
	if _, exists := m.plugins[id]; exists {
		return component.NewError(
			ErrPluginNotFound,
			"plugin with this ID already registered",
			nil,
		).WithDetail("plugin_id", id)
	}

	m.plugins[id] = plugin
	return nil
}

// Load loads a plugin and registers its components with the registry.
func (m *DefaultPluginManager) Load(ctx component.Context, id string, registry component.Registry) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the plugin exists
	plugin, exists := m.plugins[id]
	if !exists {
		return component.NewError(
			ErrPluginNotFound,
			"plugin not found",
			nil,
		).WithDetail("plugin_id", id)
	}

	// Check if the plugin is already loaded
	if _, loaded := m.loadedPlugins[id]; loaded {
		return component.NewError(
			ErrPluginLoad,
			"plugin already loaded",
			nil,
		).WithDetail("plugin_id", id)
	}

	// Load the plugin with the provided context
	err := plugin.Load(ctx, registry)
	if err != nil {
		return component.NewError(
			ErrPluginLoad,
			"failed to load plugin",
			err,
		).WithDetail("plugin_id", id)
	}

	// Mark the plugin as loaded
	m.loadedPlugins[id] = plugin
	return nil
}

// Unload unloads a plugin and unregisters its components from the registry.
func (m *DefaultPluginManager) Unload(ctx component.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if the plugin is loaded
	plugin, loaded := m.loadedPlugins[id]
	if !loaded {
		return component.NewError(
			ErrPluginUnload,
			"plugin not loaded",
			nil,
		).WithDetail("plugin_id", id)
	}

	// Unload the plugin with the provided context
	err := plugin.Unload(ctx)
	if err != nil {
		return component.NewError(
			ErrPluginUnload,
			"failed to unload plugin",
			err,
		).WithDetail("plugin_id", id)
	}

	// Remove the plugin from the loaded plugins
	delete(m.loadedPlugins, id)
	return nil
}

// GetPlugin gets a plugin by ID.
func (m *DefaultPluginManager) GetPlugin(id string) (Plugin, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugin, exists := m.plugins[id]
	if !exists {
		return nil, component.NewError(
			ErrPluginNotFound,
			"plugin not found",
			nil,
		).WithDetail("plugin_id", id)
	}

	return plugin, nil
}

// ListPlugins lists all registered plugins.
func (m *DefaultPluginManager) ListPlugins() []PluginInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	plugins := make([]PluginInfo, 0, len(m.plugins))
	for _, plugin := range m.plugins {
		// Get plugin info if the plugin supports it
		if infoProvider, ok := plugin.(interface{ GetInfo() PluginInfo }); ok {
			plugins = append(plugins, infoProvider.GetInfo())
		} else {
			// Create a basic info if not supported
			plugins = append(plugins, PluginInfo{
				ID:      plugin.ID(),
				Version: plugin.Version(),
			})
		}
	}

	return plugins
}

// IsLoaded checks if a plugin is loaded.
func (m *DefaultPluginManager) IsLoaded(id string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, loaded := m.loadedPlugins[id]
	return loaded
}

// Discover finds plugins in the specified location.
func (m *DefaultPluginManager) Discover(ctx component.Context, location string) ([]PluginInfo, error) {
	// Attempt to discover plugins from directory
	plugins, err := m.DiscoverFromDirectory(ctx, location)
	if err != nil {
		return nil, component.NewError(
			ErrPluginDiscovery,
			"failed to discover plugins",
			err,
		).WithDetail("location", location)
	}

	return plugins, nil
}

// DiscoverFromDirectory finds embedded plugins in the specified directory.
func (m *DefaultPluginManager) DiscoverFromDirectory(ctx component.Context, location string) ([]PluginInfo, error) {
	// Check if the directory exists
	info, err := m.fs.Stat(location)
	if err != nil {
		return nil, component.NewError(
			ErrPluginDiscovery,
			"failed to access directory",
			err,
		).WithDetail("location", location)
	}

	if !info.IsDir() {
		return nil, component.NewError(
			ErrPluginDiscovery,
			"not a directory",
			nil,
		).WithDetail("location", location)
	}

	// Find all plugin manifests
	var manifests []string
	err = m.fs.WalkDir(location, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Look for plugin.json files
		if !d.IsDir() && d.Name() == "plugin.json" {
			manifests = append(manifests, path)
		}

		return nil
	})

	if err != nil {
		return nil, component.NewError(
			ErrPluginDiscovery,
			"failed to scan directory",
			err,
		).WithDetail("location", location)
	}

	// In a real implementation, we would parse the manifests and load plugins
	// For now, just return an empty list
	return []PluginInfo{}, nil
}
