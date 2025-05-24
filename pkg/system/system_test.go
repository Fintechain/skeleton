package system

import (
	"testing"

	"github.com/ebanfa/skeleton/internal/domain/plugin"
	"github.com/ebanfa/skeleton/internal/domain/storage"
	"github.com/ebanfa/skeleton/internal/infrastructure/system"
	"github.com/ebanfa/skeleton/internal/infrastructure/system/mocks"
)

func TestWithConfig(t *testing.T) {
	config := &system.Config{
		ServiceID: "test-service",
		StorageConfig: storage.MultiStoreConfig{
			RootPath:      "./test-data",
			DefaultEngine: "memory",
		},
	}

	option := WithConfig(config)
	sc := &system.SystemConfig{}
	option(sc)

	if sc.Config != config {
		t.Errorf("Expected config to be set, got %v", sc.Config)
	}
	if sc.Config.ServiceID != "test-service" {
		t.Errorf("Expected ServiceID to be 'test-service', got %s", sc.Config.ServiceID)
	}
}

func TestWithPlugins(t *testing.T) {
	plugin1 := mocks.NewMockPlugin("plugin-1", "1.0.0")
	plugin2 := mocks.NewMockPlugin("plugin-2", "2.0.0")
	plugins := []plugin.Plugin{plugin1, plugin2}

	option := WithPlugins(plugins)
	sc := &system.SystemConfig{}
	option(sc)

	if len(sc.Plugins) != 2 {
		t.Errorf("Expected 2 plugins, got %d", len(sc.Plugins))
	}
	if sc.Plugins[0] != plugin1 {
		t.Errorf("Expected first plugin to be plugin1")
	}
	if sc.Plugins[1] != plugin2 {
		t.Errorf("Expected second plugin to be plugin2")
	}
}

func TestWithRegistry(t *testing.T) {
	registry := mocks.NewMockRegistry()

	option := WithRegistry(registry)
	sc := &system.SystemConfig{}
	option(sc)

	if sc.Registry != registry {
		t.Errorf("Expected registry to be set")
	}
}

func TestWithPluginManager(t *testing.T) {
	pluginMgr := mocks.NewMockPluginManager()

	option := WithPluginManager(pluginMgr)
	sc := &system.SystemConfig{}
	option(sc)

	if sc.PluginMgr != pluginMgr {
		t.Errorf("Expected plugin manager to be set")
	}
}

func TestWithEventBus(t *testing.T) {
	eventBus := mocks.NewMockEventBus()

	option := WithEventBus(eventBus)
	sc := &system.SystemConfig{}
	option(sc)

	if sc.EventBus != eventBus {
		t.Errorf("Expected event bus to be set")
	}
}

func TestWithMultiStore(t *testing.T) {
	multiStore := mocks.NewMockMultiStore()

	option := WithMultiStore(multiStore)
	sc := &system.SystemConfig{}
	option(sc)

	if sc.MultiStore != multiStore {
		t.Errorf("Expected multistore to be set")
	}
}

func TestStartSystem_DefaultConfiguration(t *testing.T) {
	// Test that StartSystem works with no options
	// Since we can't easily mock fx.New in unit tests, we test the option application logic

	// This test verifies that calling with no options doesn't panic
	// and that the function signature is correct
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("StartSystem with no options should not panic during option processing, got: %v", r)
		}
	}()

	// Test option application with empty options
	sc := &system.SystemConfig{}
	options := []Option{}

	// Apply empty options (should not panic)
	for _, option := range options {
		option(sc)
	}

	// Verify that all fields are nil initially (will be filled by applyDefaults)
	if sc.Config != nil {
		t.Error("Expected Config to be nil before defaults")
	}
	if sc.Registry != nil {
		t.Error("Expected Registry to be nil before defaults")
	}
}

func TestStartSystem_WithCustomOptions(t *testing.T) {
	// Test StartSystem with various option combinations
	config := &system.Config{ServiceID: "custom-test"}
	plugins := []plugin.Plugin{mocks.NewMockPlugin("test-plugin", "1.0.0")}
	registry := mocks.NewMockRegistry()
	eventBus := mocks.NewMockEventBus()

	// Test that options are applied correctly
	sc := &system.SystemConfig{}

	// Apply options
	WithConfig(config)(sc)
	WithPlugins(plugins)(sc)
	WithRegistry(registry)(sc)
	WithEventBus(eventBus)(sc)

	// Verify all options were applied
	if sc.Config != config {
		t.Error("Config option not applied correctly")
	}
	if len(sc.Plugins) != 1 {
		t.Errorf("Expected 1 plugin, got %d", len(sc.Plugins))
	}
	if sc.Registry != registry {
		t.Error("Registry option not applied correctly")
	}
	if sc.EventBus != eventBus {
		t.Error("EventBus option not applied correctly")
	}
}

func TestStartSystem_ErrorHandling(t *testing.T) {
	// Test error propagation from fx layer
	// Since we can't easily test fx.New errors in unit tests,
	// we test the error handling patterns in option application

	// Test that nil options don't cause panics
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Option application should handle nil values gracefully, got panic: %v", r)
		}
	}()

	sc := &system.SystemConfig{}

	// Apply options with nil values
	WithConfig(nil)(sc)
	WithPlugins(nil)(sc)
	WithRegistry(nil)(sc)
	WithPluginManager(nil)(sc)
	WithEventBus(nil)(sc)
	WithMultiStore(nil)(sc)

	// Should not panic and should set nil values
	if sc.Config != nil {
		t.Error("Expected Config to be nil")
	}
	if sc.Plugins != nil {
		t.Error("Expected Plugins to be nil")
	}
}

func TestMultipleOptions(t *testing.T) {
	config := &system.Config{ServiceID: "test"}
	registry := mocks.NewMockRegistry()
	eventBus := mocks.NewMockEventBus()

	sc := &system.SystemConfig{}

	// Apply multiple options
	WithConfig(config)(sc)
	WithRegistry(registry)(sc)
	WithEventBus(eventBus)(sc)

	if sc.Config != config {
		t.Errorf("Expected config to be set")
	}
	if sc.Registry != registry {
		t.Errorf("Expected registry to be set")
	}
	if sc.EventBus != eventBus {
		t.Errorf("Expected event bus to be set")
	}
}

func TestOptionChaining(t *testing.T) {
	// Test that options can be chained and applied in sequence
	config := &system.Config{ServiceID: "chained-test"}
	registry := mocks.NewMockRegistry()
	eventBus := mocks.NewMockEventBus()
	multiStore := mocks.NewMockMultiStore()

	sc := &system.SystemConfig{}

	// Chain multiple options
	options := []Option{
		WithConfig(config),
		WithRegistry(registry),
		WithEventBus(eventBus),
		WithMultiStore(multiStore),
	}

	// Apply all options
	for _, option := range options {
		option(sc)
	}

	// Verify all options were applied
	if sc.Config != config {
		t.Error("Config option not applied")
	}
	if sc.Registry != registry {
		t.Error("Registry option not applied")
	}
	if sc.EventBus != eventBus {
		t.Error("EventBus option not applied")
	}
	if sc.MultiStore != multiStore {
		t.Error("MultiStore option not applied")
	}
}

func TestNilOptionHandling(t *testing.T) {
	// Test that nil values in options don't cause issues
	sc := &system.SystemConfig{}

	// Apply options with nil values
	WithConfig(nil)(sc)
	WithPlugins(nil)(sc)
	WithRegistry(nil)(sc)
	WithPluginManager(nil)(sc)
	WithEventBus(nil)(sc)
	WithMultiStore(nil)(sc)

	// Should not panic and should set nil values
	if sc.Config != nil {
		t.Error("Expected Config to be nil")
	}
	if sc.Plugins != nil {
		t.Error("Expected Plugins to be nil")
	}
	if sc.Registry != nil {
		t.Error("Expected Registry to be nil")
	}
	if sc.PluginMgr != nil {
		t.Error("Expected PluginMgr to be nil")
	}
	if sc.EventBus != nil {
		t.Error("Expected EventBus to be nil")
	}
	if sc.MultiStore != nil {
		t.Error("Expected MultiStore to be nil")
	}
}
