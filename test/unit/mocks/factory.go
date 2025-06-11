package mocks

import (
	"github.com/fintechain/skeleton/internal/domain/config"
	"github.com/fintechain/skeleton/internal/domain/storage"
)

// Factory provides centralized access to all mock implementations.
// This factory ensures consistent mock creation across all test suites
// and provides a single point of configuration for mock behavior.
type Factory struct{}

// NewFactory creates a new mock factory instance.
func NewFactory() *Factory {
	return &Factory{}
}

// Domain Interface Mocks

// ComponentInterface returns a mock implementation of the Component interface.
func (f *Factory) ComponentInterface() *MockComponent {
	return &MockComponent{}
}

// RegistryInterface returns a mock implementation of the Registry interface.
func (f *Factory) RegistryInterface() *MockRegistry {
	return &MockRegistry{}
}

// SystemInterface returns a mock implementation of the System interface.
func (f *Factory) SystemInterface() *MockSystem {
	return &MockSystem{}
}

// FactoryInterface returns a mock implementation of the Factory interface.
func (f *Factory) FactoryInterface() *MockFactory {
	return &MockFactory{}
}

// OperationInterface returns a mock implementation of the Operation interface.
func (f *Factory) OperationInterface() *MockOperation {
	return &MockOperation{}
}

// ServiceInterface returns a mock implementation of the Service interface.
func (f *Factory) ServiceInterface() *MockService {
	return &MockService{}
}

// PluginInterface returns a mock implementation of the Plugin interface.
func (f *Factory) PluginInterface() *MockPlugin {
	return &MockPlugin{}
}

// PluginManagerInterface returns a mock implementation of the PluginManager interface.
func (f *Factory) PluginManagerInterface() *MockPluginManager {
	return &MockPluginManager{}
}

// EventBusInterface returns a mock implementation of the EventBus interface.
func (f *Factory) EventBusInterface() *MockEventBus {
	return &MockEventBus{}
}

// EventBusServiceInterface returns a mock implementation of the EventBusService interface.
func (f *Factory) EventBusServiceInterface() *MockEventBusService {
	return &MockEventBusService{}
}

// Storage Mocks

// Storage returns a mock storage builder for fluent configuration.
func (f *Factory) Storage() *StorageBuilder {
	return NewStorageBuilder()
}

// MultiStoreInterface returns a mock implementation of the MultiStore interface.
func (f *Factory) MultiStoreInterface() *MockMultiStore {
	return &MockMultiStore{}
}

// StorageEngineInterface returns a mock implementation of the storage Engine interface.
func (f *Factory) StorageEngineInterface() *MockEngine {
	return &MockEngine{}
}

// TransactionalInterface returns a mock implementation of the Transactional interface.
func (f *Factory) TransactionalInterface() *MockTransactional {
	return &MockTransactional{}
}

// TransactionInterface returns a mock implementation of the Transaction interface.
func (f *Factory) TransactionInterface() *MockTransaction {
	return &MockTransaction{}
}

// Configuration Mocks

// Config returns a mock configuration builder for fluent configuration.
func (f *Factory) Config() *ConfigBuilder {
	return NewConfigBuilder()
}

// ConfigurationInterface returns a mock implementation of the Configuration interface.
func (f *Factory) ConfigurationInterface() *MockConfiguration {
	return &MockConfiguration{}
}

// ConfigurationSource returns a mock implementation of the ConfigurationSource interface.
func (f *Factory) ConfigurationSource() *MockConfigurationSource {
	return &MockConfigurationSource{}
}

// Context Mocks

// ContextInterface returns a mock implementation of the Context interface.
func (f *Factory) ContextInterface() *MockContext {
	return &MockContext{}
}

// Logging Mocks

// LoggerInterface returns a mock implementation of the Logger interface.
func (f *Factory) LoggerInterface() *MockLogger {
	return &MockLogger{}
}

// LoggerServiceInterface returns a mock implementation of the LoggerService interface.
func (f *Factory) LoggerServiceInterface() *MockLoggerService {
	return &MockLoggerService{}
}

// Event Mocks

// SubscriptionInterface returns a mock implementation of the Subscription interface.
func (f *Factory) SubscriptionInterface() *MockSubscription {
	return &MockSubscription{}
}

// StorageBuilder provides fluent configuration for mock storage.
type StorageBuilder struct {
	data map[string][]byte
	name string
}

// NewStorageBuilder creates a new storage builder.
func NewStorageBuilder() *StorageBuilder {
	return &StorageBuilder{
		data: make(map[string][]byte),
	}
}

// WithKeyValue adds a key-value pair to the mock storage.
func (b *StorageBuilder) WithKeyValue(key, value string) *StorageBuilder {
	b.data[key] = []byte(value)
	return b
}

// WithName sets the name for the mock storage.
func (b *StorageBuilder) WithName(name string) *StorageBuilder {
	b.name = name
	return b
}

// Build creates the configured mock storage.
func (b *StorageBuilder) Build() storage.Store {
	mock := &MockStore{}

	// Configure mock behavior based on builder state
	for key, value := range b.data {
		keyBytes := []byte(key)
		mock.On("Get", keyBytes).Return(value, nil).Maybe()
		mock.On("Has", keyBytes).Return(true, nil).Maybe()
	}

	if b.name != "" {
		mock.On("Name").Return(b.name).Maybe()
	}

	return mock
}

// ConfigBuilder provides fluent configuration for mock configuration.
type ConfigBuilder struct {
	values map[string]interface{}
}

// NewConfigBuilder creates a new config builder.
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		values: make(map[string]interface{}),
	}
}

// WithString adds a string configuration value.
func (b *ConfigBuilder) WithString(key, value string) *ConfigBuilder {
	b.values[key] = value
	return b
}

// WithInt adds an integer configuration value.
func (b *ConfigBuilder) WithInt(key string, value int) *ConfigBuilder {
	b.values[key] = value
	return b
}

// WithBool adds a boolean configuration value.
func (b *ConfigBuilder) WithBool(key string, value bool) *ConfigBuilder {
	b.values[key] = value
	return b
}

// Build creates the configured mock configuration.
func (b *ConfigBuilder) Build() config.Configuration {
	mock := &MockConfiguration{}

	// Configure mock behavior based on builder state
	for key, value := range b.values {
		switch v := value.(type) {
		case string:
			mock.On("GetString", key).Return(v).Maybe()
		case int:
			mock.On("GetInt", key).Return(v, nil).Maybe()
		case bool:
			mock.On("GetBool", key).Return(v, nil).Maybe()
		}
	}

	return mock
}
