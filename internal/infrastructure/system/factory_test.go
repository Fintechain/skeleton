package system

import (
	"testing"

	"github.com/ebanfa/skeleton/internal/domain/system"
	"github.com/ebanfa/skeleton/internal/infrastructure/system/mocks"
)

func TestNewFactory(t *testing.T) {
	// Test creating a new factory with all dependencies
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockMultiStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	factory := NewFactory(
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockMultiStore,
		mockLogger,
	)

	if factory == nil {
		t.Error("Expected non-nil factory")
	}

	// Verify all dependencies are set
	if factory.registry != mockRegistry {
		t.Error("Expected registry to be set correctly")
	}

	if factory.pluginManager != mockPluginManager {
		t.Error("Expected plugin manager to be set correctly")
	}

	if factory.eventBus != mockEventBus {
		t.Error("Expected event bus to be set correctly")
	}

	if factory.configuration != mockConfig {
		t.Error("Expected configuration to be set correctly")
	}

	if factory.multiStore != mockMultiStore {
		t.Error("Expected multi-store to be set correctly")
	}

	if factory.logger != mockLogger {
		t.Error("Expected logger to be set correctly")
	}
}

func TestNewFactory_WithNilLogger(t *testing.T) {
	// Test creating a factory with nil logger (should create default logger)
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockMultiStore := mocks.NewMockMultiStore()

	factory := NewFactory(
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockMultiStore,
		nil, // nil logger
	)

	if factory == nil {
		t.Error("Expected non-nil factory")
	}

	if factory.logger == nil {
		t.Error("Expected logger to be created when nil is passed")
	}

	// Verify other dependencies are set
	if factory.registry != mockRegistry {
		t.Error("Expected registry to be set correctly")
	}

	if factory.pluginManager != mockPluginManager {
		t.Error("Expected plugin manager to be set correctly")
	}

	if factory.eventBus != mockEventBus {
		t.Error("Expected event bus to be set correctly")
	}

	if factory.configuration != mockConfig {
		t.Error("Expected configuration to be set correctly")
	}

	if factory.multiStore != mockMultiStore {
		t.Error("Expected multi-store to be set correctly")
	}
}

func TestFactory_CreateSystemService(t *testing.T) {
	// Test creating a system service
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockMultiStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	factory := NewFactory(
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockMultiStore,
		mockLogger,
	)

	config := &system.SystemServiceConfig{
		ServiceID:        "test-service",
		EnableOperations: true,
		EnableServices:   false,
		EnablePlugins:    true,
		EnableEventLog:   false,
	}

	service, err := factory.CreateSystemService(config)

	if err != nil {
		t.Errorf("CreateSystemService() error = %v", err)
	}

	if service == nil {
		t.Error("Expected non-nil service")
	}

	// Verify service ID
	if service.ID() != "test-service" {
		t.Errorf("Expected service ID 'test-service', got %s", service.ID())
	}

	// Verify dependencies are set correctly
	if service.Registry() != mockRegistry {
		t.Error("Expected registry to be set correctly")
	}

	if service.PluginManager() != mockPluginManager {
		t.Error("Expected plugin manager to be set correctly")
	}

	if service.EventBus() != mockEventBus {
		t.Error("Expected event bus to be set correctly")
	}

	if service.Store() != mockMultiStore {
		t.Error("Expected multi-store to be set correctly")
	}

	// Verify metadata is set correctly
	metadata := service.Metadata()
	if metadata == nil {
		t.Error("Expected non-nil metadata")
	}

	if enableOps, exists := metadata["enableOperations"]; !exists || enableOps != true {
		t.Error("Expected enableOperations to be true in metadata")
	}

	if enableSvc, exists := metadata["enableServices"]; !exists || enableSvc != false {
		t.Error("Expected enableServices to be false in metadata")
	}

	if enablePlugins, exists := metadata["enablePlugins"]; !exists || enablePlugins != true {
		t.Error("Expected enablePlugins to be true in metadata")
	}

	if enableEventLog, exists := metadata["enableEventLog"]; !exists || enableEventLog != false {
		t.Error("Expected enableEventLog to be false in metadata")
	}
}

func TestFactory_CreateSystemService_DefaultConfig(t *testing.T) {
	// Test creating a system service with default configuration values
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockMultiStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	factory := NewFactory(
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockMultiStore,
		mockLogger,
	)

	config := &system.SystemServiceConfig{
		ServiceID: "default-service",
		// All other fields will be zero values (false for bools)
	}

	service, err := factory.CreateSystemService(config)

	if err != nil {
		t.Errorf("CreateSystemService() error = %v", err)
	}

	if service == nil {
		t.Error("Expected non-nil service")
	}

	// Verify service ID
	if service.ID() != "default-service" {
		t.Errorf("Expected service ID 'default-service', got %s", service.ID())
	}

	// Verify metadata is set correctly with default values
	metadata := service.Metadata()
	if metadata == nil {
		t.Error("Expected non-nil metadata")
	}

	if enableOps, exists := metadata["enableOperations"]; !exists || enableOps != false {
		t.Error("Expected enableOperations to be false in metadata")
	}

	if enableSvc, exists := metadata["enableServices"]; !exists || enableSvc != false {
		t.Error("Expected enableServices to be false in metadata")
	}

	if enablePlugins, exists := metadata["enablePlugins"]; !exists || enablePlugins != false {
		t.Error("Expected enablePlugins to be false in metadata")
	}

	if enableEventLog, exists := metadata["enableEventLog"]; !exists || enableEventLog != false {
		t.Error("Expected enableEventLog to be false in metadata")
	}
}

func TestFactory_CreateSystemService_LoggerCalled(t *testing.T) {
	// Test that the logger is called when creating a system service
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockMultiStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	factory := NewFactory(
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockMultiStore,
		mockLogger,
	)

	config := &system.SystemServiceConfig{
		ServiceID:        "logged-service",
		EnableOperations: true,
		EnableServices:   true,
		EnablePlugins:    true,
		EnableEventLog:   true,
	}

	service, err := factory.CreateSystemService(config)

	if err != nil {
		t.Errorf("CreateSystemService() error = %v", err)
	}

	if service == nil {
		t.Error("Expected non-nil service")
	}

	// Verify logger was called
	if len(mockLogger.InfoCalls) == 0 {
		t.Error("Expected logger.Info to be called")
	}

	// Check if the log message format contains expected text
	if len(mockLogger.InfoCalls) > 0 {
		logFormat := mockLogger.InfoCalls[0]
		expectedFormat := "Created system service with ID: %s"
		if logFormat != expectedFormat {
			t.Errorf("Expected log format '%s', got '%s'", expectedFormat, logFormat)
		}
	}

	// Also check LogCalls for the actual arguments
	if len(mockLogger.LogCalls) > 0 {
		logCall := mockLogger.LogCalls[0]
		if len(logCall.Args) > 0 {
			if logCall.Args[0] != "logged-service" {
				t.Errorf("Expected log message to contain service ID 'logged-service', got %v", logCall.Args[0])
			}
		}
	}
}

func TestFactory_CreateSystemService_MultipleServices(t *testing.T) {
	// Test creating multiple services with the same factory
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockMultiStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	factory := NewFactory(
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockMultiStore,
		mockLogger,
	)

	// Create first service
	config1 := &system.SystemServiceConfig{
		ServiceID:        "service-1",
		EnableOperations: true,
	}

	service1, err1 := factory.CreateSystemService(config1)

	if err1 != nil {
		t.Errorf("CreateSystemService() error = %v", err1)
	}

	if service1 == nil {
		t.Error("Expected non-nil service1")
	}

	// Create second service
	config2 := &system.SystemServiceConfig{
		ServiceID:      "service-2",
		EnableServices: true,
	}

	service2, err2 := factory.CreateSystemService(config2)

	if err2 != nil {
		t.Errorf("CreateSystemService() error = %v", err2)
	}

	if service2 == nil {
		t.Error("Expected non-nil service2")
	}

	// Verify services are different instances
	if service1 == service2 {
		t.Error("Expected different service instances")
	}

	// Verify service IDs are different
	if service1.ID() == service2.ID() {
		t.Error("Expected different service IDs")
	}

	if service1.ID() != "service-1" {
		t.Errorf("Expected service1 ID 'service-1', got %s", service1.ID())
	}

	if service2.ID() != "service-2" {
		t.Errorf("Expected service2 ID 'service-2', got %s", service2.ID())
	}

	// Verify both services share the same dependencies
	if service1.Registry() != service2.Registry() {
		t.Error("Expected services to share the same registry")
	}

	if service1.PluginManager() != service2.PluginManager() {
		t.Error("Expected services to share the same plugin manager")
	}

	if service1.EventBus() != service2.EventBus() {
		t.Error("Expected services to share the same event bus")
	}

	if service1.Store() != service2.Store() {
		t.Error("Expected services to share the same multi-store")
	}
}
