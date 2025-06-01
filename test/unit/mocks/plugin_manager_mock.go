package mocks

import (
	"fmt"
	"sync"

	"github.com/fintechain/skeleton/pkg/context"
	"github.com/fintechain/skeleton/pkg/plugin"
	"github.com/fintechain/skeleton/pkg/registry"
)

// MockPluginManager provides a configurable mock implementation of the plugin.PluginManager interface.
type MockPluginManager struct {
	mu sync.RWMutex

	// Plugin storage
	loadedPlugins map[string]plugin.Plugin
	pluginInfos   map[string]plugin.PluginInfo

	// Behavior configuration
	shouldFail   bool
	failureError string

	// Call tracking
	callCount map[string]int
	lastCalls map[string][]interface{}
}

// NewMockPluginManager creates a new configurable plugin manager mock.
func NewMockPluginManager() *MockPluginManager {
	return &MockPluginManager{
		loadedPlugins: make(map[string]plugin.Plugin),
		pluginInfos:   make(map[string]plugin.PluginInfo),
		callCount:     make(map[string]int),
		lastCalls:     make(map[string][]interface{}),
	}
}

// PluginManager Interface Implementation

// Discover discovers plugins at the specified location.
func (m *MockPluginManager) Discover(ctx context.Context, location string) ([]plugin.PluginInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("Discover", ctx, location)

	if m.shouldFail {
		return nil, fmt.Errorf("%s", m.getFailureError("Discover"))
	}

	infos := make([]plugin.PluginInfo, 0, len(m.pluginInfos))
	for _, info := range m.pluginInfos {
		infos = append(infos, info)
	}

	return infos, nil
}

// Load loads a plugin by ID.
func (m *MockPluginManager) Load(ctx context.Context, id string, registrar registry.Registry) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("Load", ctx, id, registrar)

	if m.shouldFail {
		return fmt.Errorf("%s", m.getFailureError("Load"))
	}

	if _, exists := m.loadedPlugins[id]; exists {
		return nil // Already loaded
	}

	// Create a mock plugin
	mockPlugin := &MockPlugin{
		id:          id,
		name:        fmt.Sprintf("Plugin-%s", id),
		description: fmt.Sprintf("Mock plugin %s", id),
		version:     "1.0.0",
	}
	m.loadedPlugins[id] = mockPlugin

	return nil
}

// Unload unloads a plugin by ID.
func (m *MockPluginManager) Unload(ctx context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("Unload", ctx, id)

	if m.shouldFail {
		return fmt.Errorf("%s", m.getFailureError("Unload"))
	}

	delete(m.loadedPlugins, id)
	return nil
}

// ListPlugins returns information about all loaded plugins.
func (m *MockPluginManager) ListPlugins() []plugin.PluginInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("ListPlugins")

	infos := make([]plugin.PluginInfo, 0, len(m.pluginInfos))
	for _, info := range m.pluginInfos {
		infos = append(infos, info)
	}

	return infos
}

// GetPlugin retrieves a loaded plugin by ID.
func (m *MockPluginManager) GetPlugin(id string) (plugin.Plugin, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("GetPlugin", id)

	if m.shouldFail {
		return nil, fmt.Errorf("%s", m.getFailureError("GetPlugin"))
	}

	if p, exists := m.loadedPlugins[id]; exists {
		return p, nil
	}

	return nil, fmt.Errorf("plugin not found: %s", id)
}

// Mock Configuration Methods

// AddPluginInfo adds plugin information to the mock manager.
func (m *MockPluginManager) AddPluginInfo(info plugin.PluginInfo) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pluginInfos[info.ID] = info
}

// AddPlugin adds a plugin to the mock manager.
func (m *MockPluginManager) AddPlugin(id string, p plugin.Plugin) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.loadedPlugins[id] = p
}

// SetShouldFail configures the mock to fail operations.
func (m *MockPluginManager) SetShouldFail(fail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail = fail
}

// SetFailureError sets the error message for failed operations.
func (m *MockPluginManager) SetFailureError(err string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failureError = err
}

// State Verification Methods

// GetCallCount returns the number of times a method was called.
func (m *MockPluginManager) GetCallCount(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.callCount[method]
}

// GetLastCall returns the parameters of the last call to a method.
func (m *MockPluginManager) GetLastCall(method string) []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastCalls[method]
}

// Reset clears all mock state and configuration.
func (m *MockPluginManager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.loadedPlugins = make(map[string]plugin.Plugin)
	m.pluginInfos = make(map[string]plugin.PluginInfo)
	m.shouldFail = false
	m.failureError = ""
	m.callCount = make(map[string]int)
	m.lastCalls = make(map[string][]interface{})
}

// Helper Methods

// trackCall records a method call for verification.
func (m *MockPluginManager) trackCall(method string, args ...interface{}) {
	m.callCount[method]++
	m.lastCalls[method] = args
}

// getFailureError returns the configured failure error or a default.
func (m *MockPluginManager) getFailureError(method string) string {
	if m.failureError != "" {
		return m.failureError
	}
	return fmt.Sprintf("mock_plugin_manager.%s_failed", method)
}

// PluginManagerMockBuilder provides a fluent interface for configuring plugin manager mocks.
type PluginManagerMockBuilder struct {
	mock *MockPluginManager
}

// NewPluginManagerMockBuilder creates a new plugin manager mock builder.
func NewPluginManagerMockBuilder() *PluginManagerMockBuilder {
	return &PluginManagerMockBuilder{
		mock: NewMockPluginManager(),
	}
}

// WithPluginInfo adds plugin information to the mock manager.
func (b *PluginManagerMockBuilder) WithPluginInfo(info plugin.PluginInfo) *PluginManagerMockBuilder {
	b.mock.AddPluginInfo(info)
	return b
}

// WithPlugin adds a plugin to the mock manager.
func (b *PluginManagerMockBuilder) WithPlugin(id string, p plugin.Plugin) *PluginManagerMockBuilder {
	b.mock.AddPlugin(id, p)
	return b
}

// WithFailure configures the mock to fail operations.
func (b *PluginManagerMockBuilder) WithFailure(fail bool) *PluginManagerMockBuilder {
	b.mock.SetShouldFail(fail)
	return b
}

// WithFailureError sets the error message for failed operations.
func (b *PluginManagerMockBuilder) WithFailureError(err string) *PluginManagerMockBuilder {
	b.mock.SetFailureError(err)
	return b
}

// Build returns the configured mock plugin manager.
func (b *PluginManagerMockBuilder) Build() plugin.PluginManager {
	return b.mock
}
