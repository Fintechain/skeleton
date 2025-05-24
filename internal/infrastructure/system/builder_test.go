package system

import (
	"testing"

	"github.com/ebanfa/skeleton/internal/infrastructure/system/mocks"
)

func TestNewBuilder(t *testing.T) {
	// Test creating a new builder
	serviceID := "test-service"
	builder := NewBuilder(serviceID)

	if builder == nil {
		t.Error("Expected non-nil builder")
	}

	if builder.serviceID != serviceID {
		t.Errorf("Expected serviceID '%s', got '%s'", serviceID, builder.serviceID)
	}

	if builder.logger == nil {
		t.Error("Expected logger to be initialized")
	}
}

func TestBuilder_WithLogger(t *testing.T) {
	// Test setting logger
	builder := NewBuilder("test-service")
	mockLogger := mocks.NewMockLogger()

	result := builder.WithLogger(mockLogger)

	// Should return the same builder for chaining
	if result != builder {
		t.Error("Expected builder to return itself for chaining")
	}

	if builder.logger != mockLogger {
		t.Error("Expected logger to be set correctly")
	}
}

func TestBuilder_WithConfiguration(t *testing.T) {
	// Test setting configuration
	builder := NewBuilder("test-service")
	mockConfig := mocks.NewMockConfiguration()

	result := builder.WithConfiguration(mockConfig)

	// Should return the same builder for chaining
	if result != builder {
		t.Error("Expected builder to return itself for chaining")
	}

	if builder.configuration != mockConfig {
		t.Error("Expected configuration to be set correctly")
	}
}

func TestBuilder_WithRegistry(t *testing.T) {
	// Test setting registry
	builder := NewBuilder("test-service")
	mockRegistry := mocks.NewMockRegistry()

	result := builder.WithRegistry(mockRegistry)

	// Should return the same builder for chaining
	if result != builder {
		t.Error("Expected builder to return itself for chaining")
	}

	if builder.registry != mockRegistry {
		t.Error("Expected registry to be set correctly")
	}
}

func TestBuilder_WithPluginManager(t *testing.T) {
	// Test setting plugin manager
	builder := NewBuilder("test-service")
	mockPluginManager := mocks.NewMockPluginManager()

	result := builder.WithPluginManager(mockPluginManager)

	// Should return the same builder for chaining
	if result != builder {
		t.Error("Expected builder to return itself for chaining")
	}

	if builder.pluginManager != mockPluginManager {
		t.Error("Expected plugin manager to be set correctly")
	}
}

func TestBuilder_WithEventBus(t *testing.T) {
	// Test setting event bus
	builder := NewBuilder("test-service")
	mockEventBus := mocks.NewMockEventBus()

	result := builder.WithEventBus(mockEventBus)

	// Should return the same builder for chaining
	if result != builder {
		t.Error("Expected builder to return itself for chaining")
	}

	if builder.eventBus != mockEventBus {
		t.Error("Expected event bus to be set correctly")
	}
}

func TestBuilder_WithMultiStore(t *testing.T) {
	// Test setting multi-store
	builder := NewBuilder("test-service")
	mockMultiStore := mocks.NewMockMultiStore()

	result := builder.WithMultiStore(mockMultiStore)

	// Should return the same builder for chaining
	if result != builder {
		t.Error("Expected builder to return itself for chaining")
	}

	if builder.multiStore != mockMultiStore {
		t.Error("Expected multi-store to be set correctly")
	}
}

func TestBuilder_Build_Success(t *testing.T) {
	// Test successful build with all required dependencies
	builder := NewBuilder("test-service")
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockMultiStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	service, err := builder.
		WithRegistry(mockRegistry).
		WithPluginManager(mockPluginManager).
		WithEventBus(mockEventBus).
		WithConfiguration(mockConfig).
		WithMultiStore(mockMultiStore).
		WithLogger(mockLogger).
		Build()

	if err != nil {
		t.Errorf("Build() error = %v", err)
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
}

func TestBuilder_Build_MissingRegistry(t *testing.T) {
	// Test build failure when registry is missing
	builder := NewBuilder("test-service")
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockMultiStore := mocks.NewMockMultiStore()

	service, err := builder.
		WithPluginManager(mockPluginManager).
		WithEventBus(mockEventBus).
		WithConfiguration(mockConfig).
		WithMultiStore(mockMultiStore).
		Build()

	if err == nil {
		t.Error("Expected error when registry is missing")
	}

	if service != nil {
		t.Error("Expected nil service when build fails")
	}

	expectedError := "registry is required"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestBuilder_Build_MissingPluginManager(t *testing.T) {
	// Test build failure when plugin manager is missing
	builder := NewBuilder("test-service")
	mockRegistry := mocks.NewMockRegistry()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockMultiStore := mocks.NewMockMultiStore()

	service, err := builder.
		WithRegistry(mockRegistry).
		WithEventBus(mockEventBus).
		WithConfiguration(mockConfig).
		WithMultiStore(mockMultiStore).
		Build()

	if err == nil {
		t.Error("Expected error when plugin manager is missing")
	}

	if service != nil {
		t.Error("Expected nil service when build fails")
	}

	expectedError := "plugin manager is required"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestBuilder_Build_MissingEventBus(t *testing.T) {
	// Test build failure when event bus is missing
	builder := NewBuilder("test-service")
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockConfig := mocks.NewMockConfiguration()
	mockMultiStore := mocks.NewMockMultiStore()

	service, err := builder.
		WithRegistry(mockRegistry).
		WithPluginManager(mockPluginManager).
		WithConfiguration(mockConfig).
		WithMultiStore(mockMultiStore).
		Build()

	if err == nil {
		t.Error("Expected error when event bus is missing")
	}

	if service != nil {
		t.Error("Expected nil service when build fails")
	}

	expectedError := "event bus is required"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestBuilder_Build_MissingConfiguration(t *testing.T) {
	// Test build failure when configuration is missing
	builder := NewBuilder("test-service")
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockMultiStore := mocks.NewMockMultiStore()

	service, err := builder.
		WithRegistry(mockRegistry).
		WithPluginManager(mockPluginManager).
		WithEventBus(mockEventBus).
		WithMultiStore(mockMultiStore).
		Build()

	if err == nil {
		t.Error("Expected error when configuration is missing")
	}

	if service != nil {
		t.Error("Expected nil service when build fails")
	}

	expectedError := "configuration is required"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestBuilder_Build_MissingMultiStore(t *testing.T) {
	// Test build failure when multi-store is missing
	builder := NewBuilder("test-service")
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()

	service, err := builder.
		WithRegistry(mockRegistry).
		WithPluginManager(mockPluginManager).
		WithEventBus(mockEventBus).
		WithConfiguration(mockConfig).
		Build()

	if err == nil {
		t.Error("Expected error when multi-store is missing")
	}

	if service != nil {
		t.Error("Expected nil service when build fails")
	}

	expectedError := "multi-store is required"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestBuilder_MethodChaining(t *testing.T) {
	// Test that all methods can be chained together
	builder := NewBuilder("test-service")
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockMultiStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	// This should compile and work without issues
	result := builder.
		WithLogger(mockLogger).
		WithConfiguration(mockConfig).
		WithRegistry(mockRegistry).
		WithPluginManager(mockPluginManager).
		WithEventBus(mockEventBus).
		WithMultiStore(mockMultiStore)

	if result != builder {
		t.Error("Expected method chaining to return the same builder")
	}

	// Verify all dependencies are set
	if builder.logger != mockLogger {
		t.Error("Expected logger to be set")
	}
	if builder.configuration != mockConfig {
		t.Error("Expected configuration to be set")
	}
	if builder.registry != mockRegistry {
		t.Error("Expected registry to be set")
	}
	if builder.pluginManager != mockPluginManager {
		t.Error("Expected plugin manager to be set")
	}
	if builder.eventBus != mockEventBus {
		t.Error("Expected event bus to be set")
	}
	if builder.multiStore != mockMultiStore {
		t.Error("Expected multi-store to be set")
	}
}
