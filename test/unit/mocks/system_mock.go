package mocks

import (
	"fmt"
	"sync"

	"github.com/fintechain/skeleton/pkg/config"
	"github.com/fintechain/skeleton/pkg/context"
	"github.com/fintechain/skeleton/pkg/event"
	"github.com/fintechain/skeleton/pkg/plugin"
	"github.com/fintechain/skeleton/pkg/registry"
	"github.com/fintechain/skeleton/pkg/storage"
	"github.com/fintechain/skeleton/pkg/system"
)

// MockSystem provides a configurable mock implementation of the system.System interface.
// It supports behavior configuration, error injection, call tracking, and state verification
// for comprehensive testing of components that depend on system functionality.
type MockSystem struct {
	mu sync.RWMutex

	// System resources
	registry      registry.Registry
	pluginManager plugin.PluginManager
	eventBus      event.EventBus
	configuration config.Configuration
	store         storage.MultiStore

	// System state
	isRunning     bool
	isInitialized bool

	// Behavior configuration
	shouldFail   bool
	failureError string

	// Call tracking
	callCount map[string]int
	lastCalls map[string][]interface{}

	// Operation results
	operationResults map[string]interface{}
	operationErrors  map[string]error
}

// NewMockSystem creates a new configurable system mock.
func NewMockSystem() *MockSystem {
	factory := NewFactory()
	return &MockSystem{
		registry:         factory.RegistryInterface(),
		pluginManager:    nil, // Placeholder since plugin manager mock is not implemented yet
		eventBus:         factory.EventBusInterface(),
		configuration:    nil, // Placeholder since configuration mock is not implemented yet
		store:            factory.MultiStoreInterface(),
		isRunning:        true,
		isInitialized:    true,
		callCount:        make(map[string]int),
		lastCalls:        make(map[string][]interface{}),
		operationResults: make(map[string]interface{}),
		operationErrors:  make(map[string]error),
	}
}

// System Interface Implementation

// Registry returns the system registry.
func (m *MockSystem) Registry() registry.Registry {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.trackCall("Registry")
	return m.registry
}

// PluginManager returns the system plugin manager.
func (m *MockSystem) PluginManager() plugin.PluginManager {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.trackCall("PluginManager")
	return m.pluginManager
}

// EventBus returns the system event bus.
func (m *MockSystem) EventBus() event.EventBus {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.trackCall("EventBus")
	return m.eventBus
}

// Configuration returns the system configuration.
func (m *MockSystem) Configuration() config.Configuration {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.trackCall("Configuration")
	return m.configuration
}

// Store returns the system multi-store.
func (m *MockSystem) Store() storage.MultiStore {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.trackCall("Store")
	return m.store
}

// ExecuteOperation executes an operation by ID.
func (m *MockSystem) ExecuteOperation(ctx context.Context, operationID string, input interface{}) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("ExecuteOperation", ctx, operationID, input)

	if m.shouldFail {
		return nil, fmt.Errorf("%s", m.getFailureError("ExecuteOperation"))
	}

	if err, exists := m.operationErrors[operationID]; exists {
		return nil, err
	}

	if result, exists := m.operationResults[operationID]; exists {
		return result, nil
	}

	return fmt.Sprintf("result-for-%s", operationID), nil
}

// StartService starts a service by ID.
func (m *MockSystem) StartService(ctx context.Context, serviceID string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("StartService", ctx, serviceID)

	if m.shouldFail {
		return fmt.Errorf("%s", m.getFailureError("StartService"))
	}

	return nil
}

// StopService stops a service by ID.
func (m *MockSystem) StopService(ctx context.Context, serviceID string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("StopService", ctx, serviceID)

	if m.shouldFail {
		return fmt.Errorf("%s", m.getFailureError("StopService"))
	}

	return nil
}

// IsRunning returns true if the system is running.
func (m *MockSystem) IsRunning() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.trackCall("IsRunning")
	return m.isRunning
}

// IsInitialized returns true if the system is initialized.
func (m *MockSystem) IsInitialized() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.trackCall("IsInitialized")
	return m.isInitialized
}

// Mock Configuration Methods

// SetRegistry sets the system registry.
func (m *MockSystem) SetRegistry(registry registry.Registry) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.registry = registry
}

// SetPluginManager sets the system plugin manager.
func (m *MockSystem) SetPluginManager(pluginManager plugin.PluginManager) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pluginManager = pluginManager
}

// SetEventBus sets the system event bus.
func (m *MockSystem) SetEventBus(eventBus event.EventBus) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.eventBus = eventBus
}

// SetConfiguration sets the system configuration.
func (m *MockSystem) SetConfiguration(config interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Store as interface{} since we don't have the config mock yet
}

// SetStore sets the system multi-store.
func (m *MockSystem) SetStore(store storage.MultiStore) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.store = store
}

// SetRunning sets the system running state.
func (m *MockSystem) SetRunning(running bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.isRunning = running
}

// SetInitialized sets the system initialized state.
func (m *MockSystem) SetInitialized(initialized bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.isInitialized = initialized
}

// SetOperationResult sets the result for a specific operation.
func (m *MockSystem) SetOperationResult(operationID string, result interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.operationResults[operationID] = result
}

// SetOperationError sets the error for a specific operation.
func (m *MockSystem) SetOperationError(operationID string, err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.operationErrors[operationID] = err
}

// SetShouldFail configures the mock to fail all operations.
func (m *MockSystem) SetShouldFail(fail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail = fail
}

// SetFailureError sets the error message for failed operations.
func (m *MockSystem) SetFailureError(err string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failureError = err
}

// State Verification Methods

// GetCallCount returns the number of times a method was called.
func (m *MockSystem) GetCallCount(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.callCount[method]
}

// GetLastCall returns the parameters of the last call to a method.
func (m *MockSystem) GetLastCall(method string) []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastCalls[method]
}

// Reset clears all mock state and configuration.
func (m *MockSystem) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	factory := NewFactory()
	m.registry = factory.RegistryInterface()
	m.pluginManager = nil // Placeholder since plugin manager mock is not implemented yet
	m.eventBus = factory.EventBusInterface()
	m.configuration = nil // Placeholder since configuration mock is not implemented yet
	m.store = factory.MultiStoreInterface()
	m.isRunning = true
	m.isInitialized = true
	m.shouldFail = false
	m.failureError = ""
	m.callCount = make(map[string]int)
	m.lastCalls = make(map[string][]interface{})
	m.operationResults = make(map[string]interface{})
	m.operationErrors = make(map[string]error)
}

// Helper Methods

// trackCall records a method call for verification.
func (m *MockSystem) trackCall(method string, args ...interface{}) {
	m.callCount[method]++
	m.lastCalls[method] = args
}

// getFailureError returns the configured failure error or a default.
func (m *MockSystem) getFailureError(method string) string {
	if m.failureError != "" {
		return m.failureError
	}
	return fmt.Sprintf("mock_system.%s_failed", method)
}

// SystemMockBuilder provides a fluent interface for configuring system mocks.
type SystemMockBuilder struct {
	mock *MockSystem
}

// NewSystemMockBuilder creates a new system mock builder.
func NewSystemMockBuilder() *SystemMockBuilder {
	return &SystemMockBuilder{
		mock: NewMockSystem(),
	}
}

// WithRegistry sets the system registry.
func (b *SystemMockBuilder) WithRegistry(registry registry.Registry) *SystemMockBuilder {
	b.mock.SetRegistry(registry)
	return b
}

// WithPluginManager sets the system plugin manager.
func (b *SystemMockBuilder) WithPluginManager(pluginManager plugin.PluginManager) *SystemMockBuilder {
	b.mock.SetPluginManager(pluginManager)
	return b
}

// WithEventBus sets the system event bus.
func (b *SystemMockBuilder) WithEventBus(eventBus event.EventBus) *SystemMockBuilder {
	b.mock.SetEventBus(eventBus)
	return b
}

// WithConfiguration sets the configuration for the mock system.
func (b *SystemMockBuilder) WithConfiguration(config interface{}) *SystemMockBuilder {
	// Store as interface{} since we don't have the config mock yet
	return b
}

// WithStore sets the system multi-store.
func (b *SystemMockBuilder) WithStore(store storage.MultiStore) *SystemMockBuilder {
	b.mock.SetStore(store)
	return b
}

// WithRunning sets the system running state.
func (b *SystemMockBuilder) WithRunning(running bool) *SystemMockBuilder {
	b.mock.SetRunning(running)
	return b
}

// WithInitialized sets the system initialized state.
func (b *SystemMockBuilder) WithInitialized(initialized bool) *SystemMockBuilder {
	b.mock.SetInitialized(initialized)
	return b
}

// WithOperationResult sets the result for a specific operation.
func (b *SystemMockBuilder) WithOperationResult(operationID string, result interface{}) *SystemMockBuilder {
	b.mock.SetOperationResult(operationID, result)
	return b
}

// WithOperationError sets the error for a specific operation.
func (b *SystemMockBuilder) WithOperationError(operationID string, err error) *SystemMockBuilder {
	b.mock.SetOperationError(operationID, err)
	return b
}

// WithFailure configures the mock to fail operations.
func (b *SystemMockBuilder) WithFailure(fail bool) *SystemMockBuilder {
	b.mock.SetShouldFail(fail)
	return b
}

// WithFailureError sets the error message for failed operations.
func (b *SystemMockBuilder) WithFailureError(err string) *SystemMockBuilder {
	b.mock.SetFailureError(err)
	return b
}

// Build returns the configured mock system.
func (b *SystemMockBuilder) Build() system.System {
	return b.mock
}
