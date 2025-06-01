package plugin

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/plugin"
	"github.com/fintechain/skeleton/internal/domain/registry"
)

// DefaultPluginManager provides a concrete implementation of the PluginManager interface.
type DefaultPluginManager struct {
	filesystem plugin.FileSystem

	// Plugin storage
	mu      sync.RWMutex
	plugins map[string]plugin.Plugin
	infos   map[string]plugin.PluginInfo
}

// NewPluginManager creates a new PluginManager instance with the provided filesystem dependency.
// This constructor accepts filesystem interface dependency for testability.
func NewPluginManager(filesystem plugin.FileSystem) plugin.PluginManager {
	return &DefaultPluginManager{
		filesystem: filesystem,
		plugins:    make(map[string]plugin.Plugin),
		infos:      make(map[string]plugin.PluginInfo),
	}
}

// Discover discovers plugins at the specified location.
func (pm *DefaultPluginManager) Discover(ctx context.Context, location string) ([]plugin.PluginInfo, error) {
	var discoveredPlugins []plugin.PluginInfo

	// Check if location exists
	_, err := pm.filesystem.Stat(location)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to access location %s: %w", plugin.ErrPluginDiscovery, location, err)
	}

	// Walk the directory tree to find plugin files
	err = pm.filesystem.WalkDir(location, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Look for plugin files (simplified discovery - looking for .so files or plugin directories)
		if strings.HasSuffix(path, ".so") || strings.HasSuffix(path, ".plugin") {
			// Extract plugin info from file/directory name
			pluginID := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

			info := plugin.PluginInfo{
				ID:          pluginID,
				Name:        pluginID, // Use ID as name for simplicity
				Version:     "1.0.0",  // Default version
				Description: fmt.Sprintf("Plugin discovered at %s", path),
				Author:      "Unknown",
				Metadata: map[string]interface{}{
					"path":     path,
					"location": location,
				},
			}

			discoveredPlugins = append(discoveredPlugins, info)

			// Store the info for later use
			pm.mu.Lock()
			pm.infos[pluginID] = info
			pm.mu.Unlock()
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("%s: failed to walk directory %s: %w", plugin.ErrPluginDiscovery, location, err)
	}

	return discoveredPlugins, nil
}

// Load loads a plugin by ID.
func (pm *DefaultPluginManager) Load(ctx context.Context, id string, registrar registry.Registry) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Check if plugin is already loaded
	if _, exists := pm.plugins[id]; exists {
		return nil // Already loaded
	}

	// Check if we have info about this plugin
	info, exists := pm.infos[id]
	if !exists {
		return fmt.Errorf("%s: plugin %s not found", plugin.ErrPluginNotFound, id)
	}

	// Create a mock plugin implementation for demonstration
	// In a real implementation, this would load the actual plugin from the file system
	mockPlugin := &MockPlugin{
		id:          info.ID,
		name:        info.Name,
		description: info.Description,
		version:     info.Version,
		loaded:      false,
	}

	// Load the plugin
	err := mockPlugin.Load(ctx, registrar)
	if err != nil {
		return fmt.Errorf("%s: failed to load plugin %s: %w", plugin.ErrPluginLoad, id, err)
	}

	// Store the loaded plugin
	pm.plugins[id] = mockPlugin

	return nil
}

// Unload unloads a plugin by ID.
func (pm *DefaultPluginManager) Unload(ctx context.Context, id string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	// Check if plugin is loaded
	pluginInstance, exists := pm.plugins[id]
	if !exists {
		return fmt.Errorf("%s: plugin %s not loaded", plugin.ErrPluginNotFound, id)
	}

	// Unload the plugin
	err := pluginInstance.Unload(ctx)
	if err != nil {
		return fmt.Errorf("%s: failed to unload plugin %s: %w", plugin.ErrPluginUnload, id, err)
	}

	// Remove from loaded plugins
	delete(pm.plugins, id)

	return nil
}

// ListPlugins returns information about all discovered plugins.
func (pm *DefaultPluginManager) ListPlugins() []plugin.PluginInfo {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	var plugins []plugin.PluginInfo
	for _, info := range pm.infos {
		plugins = append(plugins, info)
	}

	return plugins
}

// GetPlugin returns a loaded plugin by ID.
func (pm *DefaultPluginManager) GetPlugin(id string) (plugin.Plugin, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	pluginInstance, exists := pm.plugins[id]
	if !exists {
		return nil, fmt.Errorf("%s: plugin %s not loaded", plugin.ErrPluginNotFound, id)
	}

	return pluginInstance, nil
}

// MockPlugin is a simple plugin implementation for demonstration purposes.
// In a real implementation, plugins would be loaded dynamically from shared libraries.
type MockPlugin struct {
	id          string
	name        string
	description string
	version     string
	loaded      bool
}

// ID returns the plugin's unique identifier.
func (p *MockPlugin) ID() string {
	return p.id
}

// Name returns the plugin's human-readable name.
func (p *MockPlugin) Name() string {
	return p.name
}

// Description returns the plugin's description.
func (p *MockPlugin) Description() string {
	return p.description
}

// Version returns the plugin's version.
func (p *MockPlugin) Version() string {
	return p.version
}

// Load loads the plugin and registers its components.
func (p *MockPlugin) Load(ctx context.Context, registrar registry.Registry) error {
	if p.loaded {
		return nil // Already loaded
	}

	// In a real implementation, this would:
	// 1. Load the plugin binary/library
	// 2. Initialize the plugin
	// 3. Register plugin components with the registrar

	p.loaded = true
	return nil
}

// Unload unloads the plugin and cleans up resources.
func (p *MockPlugin) Unload(ctx context.Context) error {
	if !p.loaded {
		return nil // Already unloaded
	}

	// In a real implementation, this would:
	// 1. Unregister plugin components
	// 2. Clean up plugin resources
	// 3. Unload the plugin binary/library

	p.loaded = false
	return nil
}
