package plugin

import (
	"errors"
	"testing"

	"github.com/ebanfa/skeleton/internal/domain/component"
	"github.com/ebanfa/skeleton/internal/domain/plugin/mocks"
)

func TestNewDefaultPlugin(t *testing.T) {
	plugin := NewDefaultPlugin(
		"test-plugin",
		"1.0.0",
		"Test Plugin",
		"Test Description",
		"Test Author",
	)

	if plugin == nil {
		t.Fatal("Expected non-nil plugin")
	}

	if plugin.id != "test-plugin" {
		t.Errorf("Expected ID to be 'test-plugin', got: %s", plugin.id)
	}

	if plugin.version != "1.0.0" {
		t.Errorf("Expected version to be '1.0.0', got: %s", plugin.version)
	}

	if plugin.name != "Test Plugin" {
		t.Errorf("Expected name to be 'Test Plugin', got: %s", plugin.name)
	}

	if plugin.description != "Test Description" {
		t.Errorf("Expected description to be 'Test Description', got: %s", plugin.description)
	}

	if plugin.author != "Test Author" {
		t.Errorf("Expected author to be 'Test Author', got: %s", plugin.author)
	}

	if plugin.components == nil {
		t.Error("Expected components to be initialized")
	}

	if plugin.metadata == nil {
		t.Error("Expected metadata to be initialized")
	}

	if plugin.isLoaded {
		t.Error("Expected plugin to not be loaded initially")
	}
}

func TestDefaultPlugin_ID(t *testing.T) {
	plugin := NewDefaultPlugin("test-plugin", "1.0.0", "Test Plugin", "", "")
	if plugin.ID() != "test-plugin" {
		t.Errorf("Expected ID to be 'test-plugin', got: %s", plugin.ID())
	}
}

func TestDefaultPlugin_Version(t *testing.T) {
	plugin := NewDefaultPlugin("test-plugin", "1.0.0", "Test Plugin", "", "")
	if plugin.Version() != "1.0.0" {
		t.Errorf("Expected version to be '1.0.0', got: %s", plugin.Version())
	}
}

func TestDefaultPlugin_Components(t *testing.T) {
	plugin := NewDefaultPlugin("test-plugin", "1.0.0", "Test Plugin", "", "")

	// Initially, components should be empty
	components := plugin.Components()
	if len(components) != 0 {
		t.Errorf("Expected no components initially, got: %d", len(components))
	}

	// Add a component
	mockComp := &mocks.MockComponent{
		IDFunc: func() string { return "comp1" },
	}
	plugin.AddComponent(mockComp)

	// Check that the component was added
	components = plugin.Components()
	if len(components) != 1 {
		t.Errorf("Expected 1 component after adding, got: %d", len(components))
	}
	if components[0] != mockComp {
		t.Error("Expected component to be the one we added")
	}

	// Verify that components are copied and not shared
	componentsRef1 := plugin.Components()
	componentsRef2 := plugin.Components()
	if &componentsRef1[0] == &componentsRef2[0] {
		t.Error("Expected components to be copied, not shared by reference")
	}
}

func TestDefaultPlugin_AddComponent(t *testing.T) {
	plugin := NewDefaultPlugin("test-plugin", "1.0.0", "Test Plugin", "", "")

	// Add multiple components
	mockComp1 := &mocks.MockComponent{
		IDFunc: func() string { return "comp1" },
	}
	mockComp2 := &mocks.MockComponent{
		IDFunc: func() string { return "comp2" },
	}

	plugin.AddComponent(mockComp1)
	plugin.AddComponent(mockComp2)

	// Check that both components were added
	components := plugin.Components()
	if len(components) != 2 {
		t.Errorf("Expected 2 components after adding, got: %d", len(components))
	}

	found1, found2 := false, false
	for _, comp := range components {
		if comp.ID() == "comp1" {
			found1 = true
		}
		if comp.ID() == "comp2" {
			found2 = true
		}
	}

	if !found1 || !found2 {
		t.Error("Not all components were found")
	}
}

func TestDefaultPlugin_GetInfo(t *testing.T) {
	plugin := NewDefaultPlugin(
		"test-plugin",
		"1.0.0",
		"Test Plugin",
		"Test Description",
		"Test Author",
	)

	// Add some metadata
	plugin.SetMetadata("key1", "value1")
	plugin.SetMetadata("key2", 42)

	// Get the plugin info
	info := plugin.GetInfo()

	// Verify the info fields
	if info.ID != "test-plugin" {
		t.Errorf("Expected ID to be 'test-plugin', got: %s", info.ID)
	}
	if info.Name != "Test Plugin" {
		t.Errorf("Expected name to be 'Test Plugin', got: %s", info.Name)
	}
	if info.Version != "1.0.0" {
		t.Errorf("Expected version to be '1.0.0', got: %s", info.Version)
	}
	if info.Description != "Test Description" {
		t.Errorf("Expected description to be 'Test Description', got: %s", info.Description)
	}
	if info.Author != "Test Author" {
		t.Errorf("Expected author to be 'Test Author', got: %s", info.Author)
	}

	// Verify metadata
	if val, ok := info.Metadata["key1"]; !ok || val != "value1" {
		t.Errorf("Expected metadata key1 to be 'value1', got: %v", val)
	}
	if val, ok := info.Metadata["key2"]; !ok || val != 42 {
		t.Errorf("Expected metadata key2 to be 42, got: %v", val)
	}

	// Verify metadata is copied and not shared
	plugin.SetMetadata("key3", "value3")
	if _, ok := info.Metadata["key3"]; ok {
		t.Error("Expected metadata to be copied, not shared by reference")
	}
}

func TestDefaultPlugin_SetMetadata(t *testing.T) {
	plugin := NewDefaultPlugin("test-plugin", "1.0.0", "Test Plugin", "", "")

	// Set some metadata
	plugin.SetMetadata("key1", "value1")
	plugin.SetMetadata("key2", 42)
	plugin.SetMetadata("key3", true)

	// Get the info and verify the metadata
	info := plugin.GetInfo()

	expectedMetadata := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
		"key3": true,
	}

	for key, expectedValue := range expectedMetadata {
		if actualValue, ok := info.Metadata[key]; !ok || actualValue != expectedValue {
			t.Errorf("Expected metadata %s to be %v, got: %v", key, expectedValue, actualValue)
		}
	}

	// Test overwriting a key
	plugin.SetMetadata("key1", "new-value")
	info = plugin.GetInfo()
	if val, ok := info.Metadata["key1"]; !ok || val != "new-value" {
		t.Errorf("Expected metadata key1 to be 'new-value', got: %v", val)
	}
}

func TestDefaultPlugin_IsLoaded(t *testing.T) {
	plugin := NewDefaultPlugin("test-plugin", "1.0.0", "Test Plugin", "", "")

	// Initially not loaded
	if plugin.IsLoaded() {
		t.Error("Expected plugin to not be loaded initially")
	}

	// Set it to loaded
	plugin.isLoaded = true
	if !plugin.IsLoaded() {
		t.Error("Expected plugin to be loaded after setting isLoaded to true")
	}

	// Set it back to not loaded
	plugin.isLoaded = false
	if plugin.IsLoaded() {
		t.Error("Expected plugin to not be loaded after setting isLoaded to false")
	}
}

func TestDefaultPlugin_Load(t *testing.T) {
	plugin := NewDefaultPlugin("test-plugin", "1.0.0", "Test Plugin", "", "")
	mockCtx := &mocks.MockContext{}
	mockRegistry := &mocks.MockRegistry{}

	// Test loading an empty plugin
	err := plugin.Load(mockCtx, mockRegistry)
	if err != nil {
		t.Errorf("Expected no error loading empty plugin, got: %v", err)
	}
	if !plugin.isLoaded {
		t.Error("Expected plugin to be marked as loaded")
	}

	// Test loading an already loaded plugin
	err = plugin.Load(mockCtx, mockRegistry)
	if err == nil {
		t.Error("Expected error loading already loaded plugin")
	}

	// Reset the plugin
	plugin.isLoaded = false

	// Add a component
	registerCalled := false
	initializeCalled := false
	mockComp := &mocks.MockComponent{
		IDFunc: func() string { return "comp1" },
		InitializeFunc: func(ctx component.Context) error {
			initializeCalled = true
			if ctx != mockCtx {
				t.Error("Expected context to be passed correctly")
			}
			return nil
		},
	}
	mockRegistry.RegisterFunc = func(comp component.Component) error {
		registerCalled = true
		if comp != mockComp {
			t.Error("Expected component to be passed correctly")
		}
		return nil
	}
	plugin.AddComponent(mockComp)

	// Test loading with components
	err = plugin.Load(mockCtx, mockRegistry)
	if err != nil {
		t.Errorf("Expected no error loading plugin with components, got: %v", err)
	}
	if !registerCalled {
		t.Error("Expected registry Register to be called")
	}
	if !initializeCalled {
		t.Error("Expected component Initialize to be called")
	}
	if !plugin.isLoaded {
		t.Error("Expected plugin to be marked as loaded")
	}
}

func TestDefaultPlugin_Load_RegisterError(t *testing.T) {
	plugin := NewDefaultPlugin("test-plugin", "1.0.0", "Test Plugin", "", "")
	mockCtx := &mocks.MockContext{}
	mockRegistry := &mocks.MockRegistry{}

	// Add a component
	mockComp := &mocks.MockComponent{
		IDFunc: func() string { return "comp1" },
	}
	plugin.AddComponent(mockComp)

	// Make Register return an error
	mockRegistry.RegisterFunc = func(comp component.Component) error {
		return errors.New("register error")
	}

	// Test loading with Register error
	err := plugin.Load(mockCtx, mockRegistry)
	if err == nil {
		t.Error("Expected error loading plugin with Register error")
	}
	if plugin.isLoaded {
		t.Error("Expected plugin to not be marked as loaded")
	}
}

func TestDefaultPlugin_Load_InitializeError(t *testing.T) {
	plugin := NewDefaultPlugin("test-plugin", "1.0.0", "Test Plugin", "", "")
	mockCtx := &mocks.MockContext{}
	mockRegistry := &mocks.MockRegistry{}

	// Add a component that fails to initialize
	mockComp := &mocks.MockComponent{
		IDFunc: func() string { return "comp1" },
		InitializeFunc: func(ctx component.Context) error {
			return errors.New("initialize error")
		},
	}
	plugin.AddComponent(mockComp)

	// Make Register succeed but Initialize fail
	unregisterCalled := false
	mockRegistry.RegisterFunc = func(comp component.Component) error {
		return nil
	}
	mockRegistry.UnregisterFunc = func(id string) error {
		unregisterCalled = true
		if id != "comp1" {
			t.Errorf("Expected unregister ID to be 'comp1', got: %s", id)
		}
		return nil
	}

	// Test loading with Initialize error
	err := plugin.Load(mockCtx, mockRegistry)
	if err == nil {
		t.Error("Expected error loading plugin with Initialize error")
	}
	if plugin.isLoaded {
		t.Error("Expected plugin to not be marked as loaded")
	}
	if !unregisterCalled {
		t.Error("Expected registry Unregister to be called")
	}
}

func TestDefaultPlugin_Unload(t *testing.T) {
	plugin := NewDefaultPlugin("test-plugin", "1.0.0", "Test Plugin", "", "")
	mockCtx := &mocks.MockContext{}

	// Test unloading a plugin that isn't loaded
	err := plugin.Unload(mockCtx)
	if err == nil {
		t.Error("Expected error unloading plugin that isn't loaded")
	}

	// Mark as loaded and add components
	plugin.isLoaded = true
	disposeCalled1 := false
	disposeCalled2 := false
	mockComp1 := &mocks.MockComponent{
		IDFunc: func() string { return "comp1" },
		DisposeFunc: func() error {
			disposeCalled1 = true
			return nil
		},
	}
	mockComp2 := &mocks.MockComponent{
		IDFunc: func() string { return "comp2" },
		DisposeFunc: func() error {
			disposeCalled2 = true
			return nil
		},
	}
	plugin.AddComponent(mockComp1)
	plugin.AddComponent(mockComp2)

	// Test unloading with components
	err = plugin.Unload(mockCtx)
	if err != nil {
		t.Errorf("Expected no error unloading loaded plugin, got: %v", err)
	}
	if plugin.isLoaded {
		t.Error("Expected plugin to be marked as not loaded")
	}
	if !disposeCalled1 {
		t.Error("Expected component 1 Dispose to be called")
	}
	if !disposeCalled2 {
		t.Error("Expected component 2 Dispose to be called")
	}
}

func TestDefaultPlugin_Unload_DisposeError(t *testing.T) {
	plugin := NewDefaultPlugin("test-plugin", "1.0.0", "Test Plugin", "", "")
	mockCtx := &mocks.MockContext{}

	// Mark as loaded and add a component that fails to dispose
	plugin.isLoaded = true
	mockComp := &mocks.MockComponent{
		IDFunc: func() string { return "comp1" },
		DisposeFunc: func() error {
			return errors.New("dispose error")
		},
	}
	plugin.AddComponent(mockComp)

	// Test unloading with Dispose error
	err := plugin.Unload(mockCtx)
	if err == nil {
		t.Error("Expected error unloading plugin with Dispose error")
	}
	if !plugin.isLoaded {
		t.Error("Expected plugin to still be marked as loaded")
	}
}
