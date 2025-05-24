package system

import (
	"testing"

	"github.com/ebanfa/skeleton/internal/infrastructure/system/mocks"
)

func TestNewConfigLoader(t *testing.T) {
	// Test creating a new config loader
	mockConfig := mocks.NewMockConfiguration()
	loader := NewConfigLoader(mockConfig)

	if loader == nil {
		t.Error("Expected non-nil config loader")
	}

	if loader.config != mockConfig {
		t.Error("Expected config to be set correctly")
	}
}

func TestNewConfigLoader_WithNilConfig(t *testing.T) {
	// Test creating a config loader with nil config
	loader := NewConfigLoader(nil)

	if loader == nil {
		t.Error("Expected non-nil config loader")
	}

	if loader.config != nil {
		t.Error("Expected config to be nil")
	}
}

func TestConfigLoader_LoadSystemConfig_DefaultValues(t *testing.T) {
	// Test loading system config with default values (nil config)
	loader := NewConfigLoader(nil)
	cfg, err := loader.LoadSystemConfig()

	if err != nil {
		t.Errorf("LoadSystemConfig() error = %v", err)
	}

	if cfg == nil {
		t.Error("Expected non-nil system config")
	}

	// Verify default values
	if cfg.ServiceID != "system" {
		t.Errorf("Expected default ServiceID 'system', got %s", cfg.ServiceID)
	}

	if !cfg.EnableOperations {
		t.Error("Expected EnableOperations to be true by default")
	}

	if !cfg.EnableServices {
		t.Error("Expected EnableServices to be true by default")
	}

	if !cfg.EnablePlugins {
		t.Error("Expected EnablePlugins to be true by default")
	}

	if !cfg.EnableEventLog {
		t.Error("Expected EnableEventLog to be true by default")
	}
}

func TestConfigLoader_LoadSystemConfig_WithConfiguration(t *testing.T) {
	// Test loading system config with custom configuration
	mockConfig := mocks.NewMockConfiguration()

	// Set up mock configuration values
	mockConfig.Set("system.serviceId", "custom-service")
	mockConfig.Set("system.enableOperations", false)
	mockConfig.Set("system.enableServices", false)
	mockConfig.Set("system.enablePlugins", false)
	mockConfig.Set("system.enableEventLog", false)
	mockConfig.Set("system.storage.rootPath", "custom-data")
	mockConfig.Set("system.storage.defaultEngine", "leveldb")

	loader := NewConfigLoader(mockConfig)
	cfg, err := loader.LoadSystemConfig()

	if err != nil {
		t.Errorf("LoadSystemConfig() error = %v", err)
	}

	if cfg == nil {
		t.Error("Expected non-nil system config")
	}

	// Verify custom values
	if cfg.ServiceID != "custom-service" {
		t.Errorf("Expected ServiceID 'custom-service', got %s", cfg.ServiceID)
	}

	if cfg.EnableOperations {
		t.Error("Expected EnableOperations to be false")
	}

	if cfg.EnableServices {
		t.Error("Expected EnableServices to be false")
	}

	if cfg.EnablePlugins {
		t.Error("Expected EnablePlugins to be false")
	}

	if cfg.EnableEventLog {
		t.Error("Expected EnableEventLog to be false")
	}

	// Verify storage configuration
	if cfg.StorageConfig.RootPath != "custom-data" {
		t.Errorf("Expected RootPath 'custom-data', got %s", cfg.StorageConfig.RootPath)
	}

	if cfg.StorageConfig.DefaultEngine != "leveldb" {
		t.Errorf("Expected DefaultEngine 'leveldb', got %s", cfg.StorageConfig.DefaultEngine)
	}

	if cfg.StorageConfig.EngineConfigs == nil {
		t.Error("Expected EngineConfigs to be initialized")
	}
}

func TestConfigLoader_LoadSystemConfig_PartialConfiguration(t *testing.T) {
	// Test loading system config with partial configuration
	mockConfig := mocks.NewMockConfiguration()

	// Set only some configuration values
	mockConfig.Set("system.serviceId", "partial-service")
	mockConfig.Set("system.enableOperations", false)
	mockConfig.Set("system.storage.rootPath", "partial-data")

	loader := NewConfigLoader(mockConfig)
	cfg, err := loader.LoadSystemConfig()

	if err != nil {
		t.Errorf("LoadSystemConfig() error = %v", err)
	}

	// Verify mixed values (custom and defaults)
	if cfg.ServiceID != "partial-service" {
		t.Errorf("Expected ServiceID 'partial-service', got %s", cfg.ServiceID)
	}

	if cfg.EnableOperations {
		t.Error("Expected EnableOperations to be false")
	}

	// These should remain default values
	if !cfg.EnableServices {
		t.Error("Expected EnableServices to be true (default)")
	}

	if !cfg.EnablePlugins {
		t.Error("Expected EnablePlugins to be true (default)")
	}

	if !cfg.EnableEventLog {
		t.Error("Expected EnableEventLog to be true (default)")
	}

	// Storage config should have mixed values
	if cfg.StorageConfig.RootPath != "partial-data" {
		t.Errorf("Expected RootPath 'partial-data', got %s", cfg.StorageConfig.RootPath)
	}

	if cfg.StorageConfig.DefaultEngine != "memory" {
		t.Errorf("Expected DefaultEngine 'memory' (default), got %s", cfg.StorageConfig.DefaultEngine)
	}
}

func TestConfigLoader_LoadSystemConfig_EmptyServiceID(t *testing.T) {
	// Test loading system config with empty service ID
	mockConfig := mocks.NewMockConfiguration()
	mockConfig.Set("system.serviceId", "")

	loader := NewConfigLoader(mockConfig)
	cfg, err := loader.LoadSystemConfig()

	if err != nil {
		t.Errorf("LoadSystemConfig() error = %v", err)
	}

	// Should use default service ID when empty string is provided
	if cfg.ServiceID != "system" {
		t.Errorf("Expected default ServiceID 'system' when empty string provided, got %s", cfg.ServiceID)
	}
}

func TestConfigLoader_LoadSystemConfig_StorageDefaults(t *testing.T) {
	// Test storage configuration defaults
	mockConfig := mocks.NewMockConfiguration()

	loader := NewConfigLoader(mockConfig)
	cfg, err := loader.LoadSystemConfig()

	if err != nil {
		t.Errorf("LoadSystemConfig() error = %v", err)
	}

	// Verify storage defaults
	if cfg.StorageConfig.RootPath != "data" {
		t.Errorf("Expected default RootPath 'data', got %s", cfg.StorageConfig.RootPath)
	}

	if cfg.StorageConfig.DefaultEngine != "memory" {
		t.Errorf("Expected default DefaultEngine 'memory', got %s", cfg.StorageConfig.DefaultEngine)
	}

	if cfg.StorageConfig.EngineConfigs == nil {
		t.Error("Expected EngineConfigs to be initialized")
	}

	if len(cfg.StorageConfig.EngineConfigs) != 0 {
		t.Errorf("Expected empty EngineConfigs, got %d entries", len(cfg.StorageConfig.EngineConfigs))
	}
}

func TestSystemServiceConfig_ExtendedConfig(t *testing.T) {
	// Test the extended SystemServiceConfig struct
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockMultiStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	// Create base config
	loader := NewConfigLoader(nil)
	baseConfig, err := loader.LoadSystemConfig()
	if err != nil {
		t.Errorf("LoadSystemConfig() error = %v", err)
	}

	// Create extended config
	extendedConfig := &SystemServiceConfig{
		SystemServiceConfig: baseConfig,
		Registry:            mockRegistry,
		PluginManager:       mockPluginManager,
		EventBus:            mockEventBus,
		MultiStore:          mockMultiStore,
		Logger:              mockLogger,
	}

	// Verify all fields are set
	if extendedConfig.SystemServiceConfig != baseConfig {
		t.Error("Expected base config to be set")
	}

	if extendedConfig.Registry != mockRegistry {
		t.Error("Expected registry to be set")
	}

	if extendedConfig.PluginManager != mockPluginManager {
		t.Error("Expected plugin manager to be set")
	}

	if extendedConfig.EventBus != mockEventBus {
		t.Error("Expected event bus to be set")
	}

	if extendedConfig.MultiStore != mockMultiStore {
		t.Error("Expected multi-store to be set")
	}

	if extendedConfig.Logger != mockLogger {
		t.Error("Expected logger to be set")
	}

	// Verify base config properties are accessible
	if extendedConfig.ServiceID != "system" {
		t.Errorf("Expected ServiceID 'system', got %s", extendedConfig.ServiceID)
	}

	if !extendedConfig.EnableOperations {
		t.Error("Expected EnableOperations to be true")
	}
}
