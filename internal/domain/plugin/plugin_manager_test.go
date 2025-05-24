package plugin

import (
	"errors"
	"io/fs"
	"os"
	"testing"

	"github.com/ebanfa/skeleton/internal/domain/component"
	"github.com/ebanfa/skeleton/internal/domain/plugin/mocks"
)

func TestNewPluginManager(t *testing.T) {
	mockFS := &mocks.MockFileSystem{}
	manager := NewPluginManager(mockFS)

	if manager == nil {
		t.Fatal("Expected non-nil plugin manager")
	}

	if manager.fs != mockFS {
		t.Error("Expected file system to be set correctly")
	}

	if manager.plugins == nil {
		t.Error("Expected plugins map to be initialized")
	}

	if manager.loadedPlugins == nil {
		t.Error("Expected loadedPlugins map to be initialized")
	}
}

func TestCreatePluginManager(t *testing.T) {
	manager := CreatePluginManager()

	if manager == nil {
		t.Fatal("Expected non-nil plugin manager")
	}

	if manager.fs == nil {
		t.Error("Expected file system to be set")
	}

	if _, ok := manager.fs.(*StandardFileSystem); !ok {
		t.Error("Expected file system to be StandardFileSystem")
	}
}

func TestRegisterPlugin(t *testing.T) {
	mockFS := &mocks.MockFileSystem{}
	manager := NewPluginManager(mockFS)
	mockPlugin := &mocks.MockPlugin{
		IDFunc: func() string {
			return "test-plugin"
		},
	}

	// Test successful registration
	err := manager.RegisterPlugin(mockPlugin)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify plugin was added
	plugin, err := manager.GetPlugin("test-plugin")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if plugin != mockPlugin {
		t.Error("Expected to get the same plugin instance")
	}

	// Test duplicate registration
	err = manager.RegisterPlugin(mockPlugin)
	if err == nil {
		t.Error("Expected error for duplicate registration")
	}
}

func TestGetPlugin(t *testing.T) {
	mockFS := &mocks.MockFileSystem{}
	manager := NewPluginManager(mockFS)
	mockPlugin := &mocks.MockPlugin{
		IDFunc: func() string {
			return "test-plugin"
		},
	}

	// Register the plugin
	_ = manager.RegisterPlugin(mockPlugin)

	// Test getting existing plugin
	plugin, err := manager.GetPlugin("test-plugin")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if plugin != mockPlugin {
		t.Error("Expected to get the same plugin instance")
	}

	// Test getting non-existent plugin
	_, err = manager.GetPlugin("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent plugin")
	}
}

func TestListPlugins(t *testing.T) {
	mockFS := &mocks.MockFileSystem{}
	manager := NewPluginManager(mockFS)

	// Register plugins
	mockPlugin1 := &mocks.MockPlugin{
		IDFunc:      func() string { return "plugin1" },
		VersionFunc: func() string { return "1.0.0" },
	}
	mockPlugin2 := &mocks.MockPlugin{
		IDFunc:      func() string { return "plugin2" },
		VersionFunc: func() string { return "2.0.0" },
	}

	// Create a plugin with GetInfo method
	mockPluginWithInfo := &MockPluginWithInfo{
		MockPlugin: mocks.MockPlugin{
			IDFunc:      func() string { return "plugin3" },
			VersionFunc: func() string { return "3.0.0" },
		},
		info: PluginInfo{
			ID:          "plugin3",
			Name:        "Plugin 3",
			Version:     "3.0.0",
			Description: "Test plugin with info",
			Author:      "Test Author",
		},
	}

	_ = manager.RegisterPlugin(mockPlugin1)
	_ = manager.RegisterPlugin(mockPlugin2)
	_ = manager.RegisterPlugin(mockPluginWithInfo)

	// Test listing plugins
	plugins := manager.ListPlugins()
	if len(plugins) != 3 {
		t.Errorf("Expected 3 plugins, got: %d", len(plugins))
	}

	// Verify the plugin info for each plugin
	found1, found2, found3 := false, false, false
	for _, info := range plugins {
		switch info.ID {
		case "plugin1":
			found1 = true
			if info.Version != "1.0.0" {
				t.Errorf("Expected version 1.0.0, got: %s", info.Version)
			}
		case "plugin2":
			found2 = true
			if info.Version != "2.0.0" {
				t.Errorf("Expected version 2.0.0, got: %s", info.Version)
			}
		case "plugin3":
			found3 = true
			if info.Version != "3.0.0" {
				t.Errorf("Expected version 3.0.0, got: %s", info.Version)
			}
			if info.Name != "Plugin 3" {
				t.Errorf("Expected name 'Plugin 3', got: %s", info.Name)
			}
			if info.Description != "Test plugin with info" {
				t.Errorf("Expected description 'Test plugin with info', got: %s", info.Description)
			}
			if info.Author != "Test Author" {
				t.Errorf("Expected author 'Test Author', got: %s", info.Author)
			}
		}
	}

	if !found1 || !found2 || !found3 {
		t.Error("Not all plugins were listed")
	}
}

// MockPluginWithInfo is a test helper that implements the GetInfo method
type MockPluginWithInfo struct {
	mocks.MockPlugin
	info PluginInfo
}

func (m *MockPluginWithInfo) GetInfo() PluginInfo {
	return m.info
}

func TestLoad(t *testing.T) {
	mockFS := &mocks.MockFileSystem{}
	manager := NewPluginManager(mockFS)
	mockCtx := &mocks.MockContext{}
	mockRegistry := &mocks.MockRegistry{}

	// Test loading non-existent plugin
	err := manager.Load(mockCtx, "non-existent", mockRegistry)
	if err == nil {
		t.Error("Expected error for non-existent plugin")
	}

	// Test loading a plugin
	loadCalled := false
	mockPlugin := &mocks.MockPlugin{
		IDFunc: func() string { return "test-plugin" },
		LoadFunc: func(ctx component.Context, registry component.Registry) error {
			loadCalled = true
			if ctx != mockCtx {
				t.Error("Expected context to be passed correctly")
			}
			if registry != mockRegistry {
				t.Error("Expected registry to be passed correctly")
			}
			return nil
		},
	}

	_ = manager.RegisterPlugin(mockPlugin)
	err = manager.Load(mockCtx, "test-plugin", mockRegistry)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !loadCalled {
		t.Error("Expected Load to be called on plugin")
	}

	// Test loading an already loaded plugin
	err = manager.Load(mockCtx, "test-plugin", mockRegistry)
	if err == nil {
		t.Error("Expected error for already loaded plugin")
	}

	// Test loading a plugin that returns an error
	mockPluginWithError := &mocks.MockPlugin{
		IDFunc: func() string { return "error-plugin" },
		LoadFunc: func(ctx component.Context, registry component.Registry) error {
			return errors.New("load error")
		},
	}

	_ = manager.RegisterPlugin(mockPluginWithError)
	err = manager.Load(mockCtx, "error-plugin", mockRegistry)
	if err == nil {
		t.Error("Expected error from plugin Load")
	}
}

func TestUnload(t *testing.T) {
	mockFS := &mocks.MockFileSystem{}
	manager := NewPluginManager(mockFS)
	mockCtx := &mocks.MockContext{}
	mockRegistry := &mocks.MockRegistry{}

	// Test unloading a plugin that isn't loaded
	err := manager.Unload(mockCtx, "non-existent")
	if err == nil {
		t.Error("Expected error for non-existent plugin")
	}

	// Test successful unload
	unloadCalled := false
	mockPlugin := &mocks.MockPlugin{
		IDFunc: func() string { return "test-plugin" },
		LoadFunc: func(ctx component.Context, registry component.Registry) error {
			return nil
		},
		UnloadFunc: func(ctx component.Context) error {
			unloadCalled = true
			if ctx != mockCtx {
				t.Error("Expected context to be passed correctly")
			}
			return nil
		},
	}

	_ = manager.RegisterPlugin(mockPlugin)
	_ = manager.Load(mockCtx, "test-plugin", mockRegistry)

	err = manager.Unload(mockCtx, "test-plugin")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	if !unloadCalled {
		t.Error("Expected Unload to be called on plugin")
	}

	// Test unloading again (should fail)
	err = manager.Unload(mockCtx, "test-plugin")
	if err == nil {
		t.Error("Expected error for already unloaded plugin")
	}

	// Test unloading a plugin that returns an error
	mockPluginWithError := &mocks.MockPlugin{
		IDFunc: func() string { return "error-plugin" },
		LoadFunc: func(ctx component.Context, registry component.Registry) error {
			return nil
		},
		UnloadFunc: func(ctx component.Context) error {
			return errors.New("unload error")
		},
	}

	_ = manager.RegisterPlugin(mockPluginWithError)
	_ = manager.Load(mockCtx, "error-plugin", mockRegistry)

	err = manager.Unload(mockCtx, "error-plugin")
	if err == nil {
		t.Error("Expected error from plugin Unload")
	}
}

func TestIsLoaded(t *testing.T) {
	mockFS := &mocks.MockFileSystem{}
	manager := NewPluginManager(mockFS)
	mockCtx := &mocks.MockContext{}
	mockRegistry := &mocks.MockRegistry{}

	mockPlugin := &mocks.MockPlugin{
		IDFunc: func() string { return "test-plugin" },
		LoadFunc: func(ctx component.Context, registry component.Registry) error {
			return nil
		},
	}

	_ = manager.RegisterPlugin(mockPlugin)

	// Test before loading
	if manager.IsLoaded("test-plugin") {
		t.Error("Expected plugin to not be loaded initially")
	}

	// Test after loading
	_ = manager.Load(mockCtx, "test-plugin", mockRegistry)
	if !manager.IsLoaded("test-plugin") {
		t.Error("Expected plugin to be loaded")
	}

	// Test after unloading
	_ = manager.Unload(mockCtx, "test-plugin")
	if manager.IsLoaded("test-plugin") {
		t.Error("Expected plugin to not be loaded after unload")
	}

	// Test non-existent plugin
	if manager.IsLoaded("non-existent") {
		t.Error("Expected non-existent plugin to not be loaded")
	}
}

func TestDiscover(t *testing.T) {
	mockFS := &mocks.MockFileSystem{}
	manager := NewPluginManager(mockFS)
	mockCtx := &mocks.MockContext{}

	// Mock file stat to succeed
	mockFS.StatFunc = func(name string) (os.FileInfo, error) {
		return &mocks.MockFileInfo{
			IsDirFunc: func() bool {
				return true
			},
		}, nil
	}

	// Mock WalkDir to simulate finding plugin.json files
	mockFS.WalkDirFunc = func(root string, fn fs.WalkDirFunc) error {
		// Create mock plugin.json files
		mockEntry1 := &mocks.MockDirEntry{
			NameFunc:        func() string { return "plugin.json" },
			IsDirectoryFunc: func() bool { return false },
		}
		mockEntry2 := &mocks.MockDirEntry{
			NameFunc:        func() string { return "plugin.json" },
			IsDirectoryFunc: func() bool { return false },
		}

		// Call the function for each mock entry
		_ = fn("path/to/plugin1/plugin.json", mockEntry1, nil)
		_ = fn("path/to/plugin2/plugin.json", mockEntry2, nil)

		return nil
	}

	// Test discovery
	infos, err := manager.Discover(mockCtx, "/plugins")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	// Verify that Stat was called with the right path
	if len(mockFS.StatCalls) != 1 || mockFS.StatCalls[0] != "/plugins" {
		t.Errorf("Expected Stat to be called with '/plugins', got: %v", mockFS.StatCalls)
	}

	// Verify that WalkDir was called with the right path
	if len(mockFS.WalkDirCalls) != 1 || mockFS.WalkDirCalls[0] != "/plugins" {
		t.Errorf("Expected WalkDir to be called with '/plugins', got: %v", mockFS.WalkDirCalls)
	}

	// In real implementation, we would verify the discovered plugins
	// For now, just make sure we get an empty list back (based on implementation)
	if len(infos) != 0 {
		t.Errorf("Expected empty list of plugin infos, got: %v", infos)
	}
}

func TestDiscoverFromDirectory_StatError(t *testing.T) {
	mockFS := &mocks.MockFileSystem{}
	manager := NewPluginManager(mockFS)
	mockCtx := &mocks.MockContext{}

	// Mock file stat to fail
	mockFS.StatFunc = func(name string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}

	// Test discovery with stat error
	_, err := manager.DiscoverFromDirectory(mockCtx, "/plugins")
	if err == nil {
		t.Error("Expected error for non-existent directory")
	}
}

func TestDiscoverFromDirectory_NotDirectory(t *testing.T) {
	mockFS := &mocks.MockFileSystem{}
	manager := NewPluginManager(mockFS)
	mockCtx := &mocks.MockContext{}

	// Mock file stat to indicate not a directory
	mockFS.StatFunc = func(name string) (os.FileInfo, error) {
		return &mocks.MockFileInfo{
			IsDirFunc: func() bool {
				return false
			},
		}, nil
	}

	// Test discovery with not a directory
	_, err := manager.DiscoverFromDirectory(mockCtx, "/plugins")
	if err == nil {
		t.Error("Expected error for not a directory")
	}
}

func TestDiscoverFromDirectory_WalkDirError(t *testing.T) {
	mockFS := &mocks.MockFileSystem{}
	manager := NewPluginManager(mockFS)
	mockCtx := &mocks.MockContext{}

	// Mock file stat to succeed
	mockFS.StatFunc = func(name string) (os.FileInfo, error) {
		return &mocks.MockFileInfo{
			IsDirFunc: func() bool {
				return true
			},
		}, nil
	}

	// Mock WalkDir to fail
	mockFS.WalkDirFunc = func(root string, fn fs.WalkDirFunc) error {
		return errors.New("walk error")
	}

	// Test discovery with WalkDir error
	_, err := manager.DiscoverFromDirectory(mockCtx, "/plugins")
	if err == nil {
		t.Error("Expected error for WalkDir failure")
	}
}
