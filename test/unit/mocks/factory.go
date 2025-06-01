// Package mocks provides centralized mock implementations for the Skeleton Framework.
// This package implements the mock-driven dependency injection pattern for comprehensive
// unit testing of infrastructure implementations.
package mocks

import (
	"github.com/fintechain/skeleton/pkg/component"
	"github.com/fintechain/skeleton/pkg/context"
	"github.com/fintechain/skeleton/pkg/event"
	"github.com/fintechain/skeleton/pkg/plugin"
	"github.com/fintechain/skeleton/pkg/registry"
	"github.com/fintechain/skeleton/pkg/storage"
	"github.com/fintechain/skeleton/pkg/system"
)

// Factory provides centralized creation of all mock types for consistent testing.
// It ensures that all mocks are created with proper configuration and can be
// easily managed across test suites.
type Factory struct {
	// Configuration for default mock behavior
	defaultFailureRate float64
	enableCallTracking bool
	threadSafe         bool
}

// NewFactory creates a new mock factory with default configuration.
func NewFactory() *Factory {
	return &Factory{
		defaultFailureRate: 0.0,
		enableCallTracking: true,
		threadSafe:         true,
	}
}

// Registry Interface Methods

// RegistryInterface creates a registry mock with default configuration for
// testing registry operations. The mock provides comprehensive functionality for
// configurable behavior for testing different scenarios.
func (f *Factory) RegistryInterface() registry.Registry {
	return f.CreateRegistryMock()
}

// Component Interface Methods

// ComponentInterface creates a component mock with default configuration for
// testing component lifecycle operations. The mock provides comprehensive functionality for
// configurable behavior for testing component interactions.
func (f *Factory) ComponentInterface() component.Component {
	return f.CreateComponentMock()
}

// ComponentFactoryInterface creates a component factory mock with default configuration for
// testing component factory operations. The mock provides comprehensive functionality for
// configurable behavior for testing component creation.
func (f *Factory) ComponentFactoryInterface() component.Factory {
	return f.CreateComponentFactoryMock()
}

// Context Interface Methods

// ContextInterface creates a context mock with default configuration for
// testing framework context operations. The mock provides comprehensive functionality for
// The mock implements the context.Context interface from the framework.
func (f *Factory) ContextInterface() context.Context {
	return f.CreateContextMock()
}

// System Interface Methods

// SystemInterface creates a system mock with default configuration for
// testing system operations. The mock provides comprehensive functionality for
// configurable behavior for testing system operations.
func (f *Factory) SystemInterface() system.System {
	return f.CreateSystemMock()
}

// EventBus Interface Methods

// EventBusInterface creates an event bus mock with default configuration for
// testing event bus operations. The mock provides comprehensive functionality for
// publish/subscribe behavior mocking.
func (f *Factory) EventBusInterface() event.EventBus {
	return f.CreateEventBusMock()
}

// Storage Interface Methods

// MultiStoreInterface creates a multi-store mock with default configuration for
// testing multi-store operations. The mock provides comprehensive functionality for
// storage operation behavior mocking.
func (f *Factory) MultiStoreInterface() storage.MultiStore {
	return f.CreateMultiStoreMock()
}

// Builder Interface Methods

// Component creates a component mock builder for fluent configuration.
// This provides a fluent interface for configuring component mock behavior.
func (f *Factory) Component() *ComponentMockBuilder {
	return f.CreateComponentMockBuilder()
}

// Registry creates a registry mock builder for fluent configuration.
// This provides a fluent interface for configuring registry mock behavior.
func (f *Factory) Registry() *RegistryMockBuilder {
	return f.CreateRegistryMockBuilder()
}

// Context creates a context mock builder for fluent configuration.
// This provides a fluent interface for configuring framework context mock behavior.
func (f *Factory) Context() *ContextMockBuilder {
	return f.CreateContextMockBuilder()
}

// System creates a system mock builder for fluent configuration.
// This provides a fluent interface for configuring system mock behavior.
func (f *Factory) System() *SystemMockBuilder {
	return f.CreateSystemMockBuilder()
}

// EventBus creates an event bus mock builder for fluent configuration.
// This provides a fluent interface for configuring event bus mock behavior.
func (f *Factory) EventBus() *EventBusMockBuilder {
	return f.CreateEventBusMockBuilder()
}

// Registry Mocks

// CreateRegistryMock creates a new registry mock with optional configuration.
func (f *Factory) CreateRegistryMock() registry.Registry {
	mock := NewMockRegistry()
	f.applyDefaultConfig(mock)
	return mock
}

// CreateRegistryMockBuilder creates a new registry mock builder for fluent configuration.
func (f *Factory) CreateRegistryMockBuilder() *RegistryMockBuilder {
	return NewRegistryMockBuilder()
}

// Component Mocks

// CreateComponentMock creates a new component mock with optional configuration.
func (f *Factory) CreateComponentMock() component.Component {
	mock := NewMockComponent()
	f.applyDefaultConfig(mock)
	return mock
}

// CreateComponentFactoryMock creates a new component factory mock with optional configuration.
func (f *Factory) CreateComponentFactoryMock() component.Factory {
	mock := NewMockComponentFactory()
	f.applyDefaultConfig(mock)
	return mock
}

// CreateComponentMockBuilder creates a new component mock builder for fluent configuration.
func (f *Factory) CreateComponentMockBuilder() *ComponentMockBuilder {
	return NewComponentMockBuilder()
}

// Context Mocks

// CreateContextMock creates a new context mock with optional configuration.
func (f *Factory) CreateContextMock() context.Context {
	mock := NewMockContext()
	f.applyDefaultConfig(mock)
	return mock
}

// CreateContextMockBuilder creates a new context mock builder for fluent configuration.
func (f *Factory) CreateContextMockBuilder() *ContextMockBuilder {
	return NewContextMockBuilder()
}

// System Mocks

// CreateSystemMock creates a new system mock with optional configuration.
func (f *Factory) CreateSystemMock() system.System {
	mock := NewMockSystem()
	f.applyDefaultConfig(mock)
	return mock
}

// CreateSystemMockBuilder creates a new system mock builder for fluent configuration.
func (f *Factory) CreateSystemMockBuilder() *SystemMockBuilder {
	return NewSystemMockBuilder()
}

// EventBus Mocks

// CreateEventBusMock creates a new event bus mock with optional configuration.
func (f *Factory) CreateEventBusMock() event.EventBus {
	mock := NewMockEventBus()
	f.applyDefaultConfig(mock)
	return mock
}

// CreateEventBusMockBuilder creates a new event bus mock builder for fluent configuration.
func (f *Factory) CreateEventBusMockBuilder() *EventBusMockBuilder {
	return NewEventBusMockBuilder()
}

// Storage Mocks

// CreateStorageMock creates a new storage mock with optional configuration.
func (f *Factory) CreateStorageMock() storage.Store {
	mock := NewMockStore()
	f.applyDefaultConfig(mock)
	return mock
}

// CreateMultiStoreMock creates a new multi-store mock with optional configuration.
func (f *Factory) CreateMultiStoreMock() storage.MultiStore {
	mock := NewMockMultiStore()
	f.applyDefaultConfig(mock)
	return mock
}

// Plugin Interface Methods

// PluginInterface creates a plugin mock with default configuration for
// testing plugin operations. The mock provides comprehensive functionality for
// configurable behavior for testing plugin lifecycle.
func (f *Factory) PluginInterface() plugin.Plugin {
	return f.CreatePluginMock()
}

// PluginManagerInterface creates a plugin manager mock with default configuration for
// testing plugin manager operations. The mock provides comprehensive functionality for
// plugin discovery and lifecycle management.
func (f *Factory) PluginManagerInterface() plugin.PluginManager {
	return f.CreatePluginManagerMock()
}

// Plugin Mocks

// CreatePluginMock creates a new plugin mock with optional configuration.
func (f *Factory) CreatePluginMock() plugin.Plugin {
	mock := NewMockPlugin()
	f.applyDefaultConfig(mock)
	return mock
}

// CreatePluginManagerMock creates a new plugin manager mock with optional configuration.
func (f *Factory) CreatePluginManagerMock() plugin.PluginManager {
	mock := NewMockPluginManager()
	f.applyDefaultConfig(mock)
	return mock
}

// Factory Configuration

// SetDefaultFailureRate sets the default failure rate for all created mocks.
func (f *Factory) SetDefaultFailureRate(rate float64) *Factory {
	f.defaultFailureRate = rate
	return f
}

// SetCallTracking enables or disables call tracking for all created mocks.
func (f *Factory) SetCallTracking(enabled bool) *Factory {
	f.enableCallTracking = enabled
	return f
}

// SetThreadSafe enables or disables thread safety for all created mocks.
func (f *Factory) SetThreadSafe(enabled bool) *Factory {
	f.threadSafe = enabled
	return f
}

// Helper Methods

// applyDefaultConfig applies the factory's default configuration to a mock.
func (f *Factory) applyDefaultConfig(mock interface{}) {
	// This is a placeholder for applying default configuration
	// Each mock type would need to implement a common interface
	// for configuration if we want to apply defaults uniformly
}

// Reset clears all factory configuration and returns to defaults.
func (f *Factory) Reset() {
	f.defaultFailureRate = 0.0
	f.enableCallTracking = true
	f.threadSafe = true
}
