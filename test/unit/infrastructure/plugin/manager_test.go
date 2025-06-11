package plugin

import (
	"testing"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/plugin"
	infraContext "github.com/fintechain/skeleton/internal/infrastructure/context"
	infraPlugin "github.com/fintechain/skeleton/internal/infrastructure/plugin"
	"github.com/fintechain/skeleton/test/unit/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewPluginManager(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "plugin-manager",
		Name: "Plugin Manager",
		Type: component.TypeService,
	}

	manager := infraPlugin.NewManager(config)
	assert.NotNil(t, manager)

	// Verify interface compliance
	var _ plugin.PluginManager = manager
	var _ component.Service = manager
	var _ component.Component = manager

	// Test basic properties
	assert.Equal(t, component.ComponentID("plugin-manager"), manager.ID())
	assert.Equal(t, "Plugin Manager", manager.Name())
	assert.Equal(t, component.TypeService, manager.Type())
}

func TestPluginManagerInitialState(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "plugin-manager",
		Name: "Plugin Manager",
		Type: component.TypeService,
	}

	manager := infraPlugin.NewManager(config)

	// Test initial state
	assert.False(t, manager.IsRunning())
	assert.Equal(t, component.StatusStopped, manager.Status())
}

func TestPluginManagerAdd(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "plugin-manager",
		Name: "Plugin Manager",
		Type: component.TypeService,
	}

	manager := infraPlugin.NewManager(config)
	factory := mocks.NewFactory()
	mockPlugin := factory.PluginInterface()

	// Test adding a plugin
	err := manager.Add("test-plugin", mockPlugin)
	assert.NoError(t, err)

	// Test adding duplicate plugin (should overwrite, no error)
	err = manager.Add("test-plugin", mockPlugin)
	assert.NoError(t, err)
}

func TestPluginManagerAddNilPlugin(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "plugin-manager",
		Name: "Plugin Manager",
		Type: component.TypeService,
	}

	manager := infraPlugin.NewManager(config)

	// Test adding nil plugin (current implementation allows this)
	err := manager.Add("test-plugin", nil)
	assert.NoError(t, err)
}

func TestPluginManagerGet(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "plugin-manager",
		Name: "Plugin Manager",
		Type: component.TypeService,
	}

	manager := infraPlugin.NewManager(config)
	factory := mocks.NewFactory()
	mockPlugin := factory.PluginInterface()

	// Add plugin first
	err := manager.Add("test-plugin", mockPlugin)
	assert.NoError(t, err)

	// Test getting existing plugin
	retrievedPlugin, err := manager.GetPlugin("test-plugin")
	assert.NoError(t, err)
	assert.Equal(t, mockPlugin, retrievedPlugin)

	// Test getting non-existent plugin
	_, err = manager.GetPlugin("non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "component.component_not_found")
}

func TestPluginManagerRemove(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "plugin-manager",
		Name: "Plugin Manager",
		Type: component.TypeService,
	}

	manager := infraPlugin.NewManager(config)
	factory := mocks.NewFactory()
	mockPlugin := factory.PluginInterface()

	// Add plugin first
	err := manager.Add("test-plugin", mockPlugin)
	assert.NoError(t, err)

	// Test removing existing plugin
	err = manager.Remove("test-plugin")
	assert.NoError(t, err)

	// Verify plugin is removed
	_, err = manager.GetPlugin("test-plugin")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "component.component_not_found")

	// Test removing non-existent plugin
	err = manager.Remove("non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "component.component_not_found")
}

func TestPluginManagerList(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "plugin-manager",
		Name: "Plugin Manager",
		Type: component.TypeService,
	}

	manager := infraPlugin.NewManager(config)
	factory := mocks.NewFactory()

	// Test empty list
	pluginIDs := manager.ListPlugins()
	assert.Empty(t, pluginIDs)

	// Add some plugins
	mockPlugin1 := factory.PluginInterface()
	mockPlugin2 := factory.PluginInterface()

	err := manager.Add("plugin1", mockPlugin1)
	assert.NoError(t, err)
	err = manager.Add("plugin2", mockPlugin2)
	assert.NoError(t, err)

	// Test list with plugins
	pluginIDs = manager.ListPlugins()
	assert.Len(t, pluginIDs, 2)
	assert.Contains(t, pluginIDs, component.ComponentID("plugin1"))
	assert.Contains(t, pluginIDs, component.ComponentID("plugin2"))
}

func TestPluginManagerLifecycle(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "plugin-manager",
		Name: "Plugin Manager",
		Type: component.TypeService,
	}

	manager := infraPlugin.NewManager(config)
	ctx := infraContext.NewContext()

	// Test start without plugins
	err := manager.Start(ctx)
	assert.NoError(t, err)
	assert.True(t, manager.IsRunning())
	assert.Equal(t, component.StatusRunning, manager.Status())

	// Test stop
	err = manager.Stop(ctx)
	assert.NoError(t, err)
	assert.False(t, manager.IsRunning())
	assert.Equal(t, component.StatusStopped, manager.Status())
}

func TestPluginManagerInitializeAndDispose(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "plugin-manager",
		Name: "Plugin Manager",
		Type: component.TypeService,
	}

	manager := infraPlugin.NewManager(config)
	factory := mocks.NewFactory()
	mockSystem := factory.SystemInterface()
	ctx := infraContext.NewContext()

	// Test initialization
	err := manager.Initialize(ctx, mockSystem)
	assert.NoError(t, err)

	// Test disposal
	err = manager.Dispose()
	assert.NoError(t, err)
}
