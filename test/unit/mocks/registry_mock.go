package mocks

import (
	"fmt"
	"sync"

	"github.com/fintechain/skeleton/pkg/registry"
)

// MockRegistry provides a configurable mock implementation of the registry.Registry interface.
// It supports behavior configuration, error injection, call tracking, and state verification
// for comprehensive testing of components that depend on registry functionality.
type MockRegistry struct {
	mu sync.RWMutex

	// Storage for registered items
	items map[string]registry.Identifiable

	// Behavior configuration
	shouldFail    bool
	failureError  string
	returnItems   map[string]registry.Identifiable
	forceNotFound map[string]bool
	forceExists   map[string]bool

	// Call tracking
	callCount map[string]int
	lastCalls map[string][]interface{}

	// State verification
	registerCalls []registry.Identifiable
	getCalls      []string
	removeCalls   []string
}

// NewMockRegistry creates a new configurable registry mock.
func NewMockRegistry() *MockRegistry {
	return &MockRegistry{
		items:         make(map[string]registry.Identifiable),
		returnItems:   make(map[string]registry.Identifiable),
		forceNotFound: make(map[string]bool),
		forceExists:   make(map[string]bool),
		callCount:     make(map[string]int),
		lastCalls:     make(map[string][]interface{}),
		registerCalls: make([]registry.Identifiable, 0),
		getCalls:      make([]string, 0),
		removeCalls:   make([]string, 0),
	}
}

// Registry Interface Implementation

// Register stores an item in the registry.
func (m *MockRegistry) Register(item registry.Identifiable) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("Register", item)
	m.registerCalls = append(m.registerCalls, item)

	if m.shouldFail {
		return fmt.Errorf("%s", m.getFailureError("Register"))
	}

	if item == nil {
		return fmt.Errorf(registry.ErrInvalidItem)
	}

	if _, exists := m.items[item.ID()]; exists {
		return fmt.Errorf(registry.ErrItemAlreadyExists)
	}

	m.items[item.ID()] = item
	return nil
}

// Get retrieves an item by its ID.
func (m *MockRegistry) Get(id string) (registry.Identifiable, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("Get", id)
	m.getCalls = append(m.getCalls, id)

	if m.shouldFail {
		return nil, fmt.Errorf("%s", m.getFailureError("Get"))
	}

	// Check for forced not found
	if m.forceNotFound[id] {
		return nil, fmt.Errorf(registry.ErrItemNotFound)
	}

	// Check for configured return items
	if item, exists := m.returnItems[id]; exists {
		return item, nil
	}

	// Check actual storage
	if item, exists := m.items[id]; exists {
		return item, nil
	}

	return nil, fmt.Errorf(registry.ErrItemNotFound)
}

// List returns all registered items.
func (m *MockRegistry) List() []registry.Identifiable {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("List")

	items := make([]registry.Identifiable, 0, len(m.items))
	for _, item := range m.items {
		items = append(items, item)
	}

	// Add configured return items
	for _, item := range m.returnItems {
		items = append(items, item)
	}

	return items
}

// Remove removes an item from the registry.
func (m *MockRegistry) Remove(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("Remove", id)
	m.removeCalls = append(m.removeCalls, id)

	if m.shouldFail {
		return fmt.Errorf("%s", m.getFailureError("Remove"))
	}

	if !m.Has(id) {
		return fmt.Errorf(registry.ErrItemNotFound)
	}

	delete(m.items, id)
	return nil
}

// Has checks if an item with the given ID exists.
func (m *MockRegistry) Has(id string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("Has", id)

	// Check for forced existence
	if forced, exists := m.forceExists[id]; exists {
		return forced
	}

	// Check configured return items
	if _, exists := m.returnItems[id]; exists {
		return true
	}

	// Check actual storage
	_, exists := m.items[id]
	return exists
}

// Count returns the number of registered items.
func (m *MockRegistry) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("Count")

	return len(m.items) + len(m.returnItems)
}

// Clear removes all items from the registry.
func (m *MockRegistry) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("Clear")

	m.items = make(map[string]registry.Identifiable)
	m.returnItems = make(map[string]registry.Identifiable)
}

// Mock Configuration Methods

// SetShouldFail configures the mock to fail all operations.
func (m *MockRegistry) SetShouldFail(fail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail = fail
}

// SetFailureError sets the error message for failed operations.
func (m *MockRegistry) SetFailureError(err string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failureError = err
}

// SetReturnItem configures the mock to return a specific item for a given ID.
func (m *MockRegistry) SetReturnItem(id string, item registry.Identifiable) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.returnItems[id] = item
}

// SetForceNotFound configures the mock to always return "not found" for a specific ID.
func (m *MockRegistry) SetForceNotFound(id string, notFound bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.forceNotFound[id] = notFound
}

// SetForceExists configures the mock to always return "exists" for a specific ID.
func (m *MockRegistry) SetForceExists(id string, exists bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.forceExists[id] = exists
}

// State Verification Methods

// GetCallCount returns the number of times a method was called.
func (m *MockRegistry) GetCallCount(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.callCount[method]
}

// GetLastCall returns the parameters of the last call to a method.
func (m *MockRegistry) GetLastCall(method string) []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastCalls[method]
}

// GetRegisterCalls returns all items that were registered.
func (m *MockRegistry) GetRegisterCalls() []registry.Identifiable {
	m.mu.RLock()
	defer m.mu.RUnlock()
	calls := make([]registry.Identifiable, len(m.registerCalls))
	copy(calls, m.registerCalls)
	return calls
}

// GetGetCalls returns all IDs that were requested via Get.
func (m *MockRegistry) GetGetCalls() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	calls := make([]string, len(m.getCalls))
	copy(calls, m.getCalls)
	return calls
}

// GetRemoveCalls returns all IDs that were requested for removal.
func (m *MockRegistry) GetRemoveCalls() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	calls := make([]string, len(m.removeCalls))
	copy(calls, m.removeCalls)
	return calls
}

// Reset clears all mock state and configuration.
func (m *MockRegistry) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items = make(map[string]registry.Identifiable)
	m.returnItems = make(map[string]registry.Identifiable)
	m.forceNotFound = make(map[string]bool)
	m.forceExists = make(map[string]bool)
	m.callCount = make(map[string]int)
	m.lastCalls = make(map[string][]interface{})
	m.registerCalls = make([]registry.Identifiable, 0)
	m.getCalls = make([]string, 0)
	m.removeCalls = make([]string, 0)
	m.shouldFail = false
	m.failureError = ""
}

// Helper Methods

// trackCall records a method call for verification.
func (m *MockRegistry) trackCall(method string, args ...interface{}) {
	m.callCount[method]++
	m.lastCalls[method] = args
}

// getFailureError returns the configured failure error or a default.
func (m *MockRegistry) getFailureError(method string) string {
	if m.failureError != "" {
		return m.failureError
	}
	return fmt.Sprintf("mock_registry.%s_failed", method)
}

// RegistryMockBuilder provides a fluent interface for configuring registry mocks.
type RegistryMockBuilder struct {
	mock *MockRegistry
}

// NewRegistryMockBuilder creates a new registry mock builder.
func NewRegistryMockBuilder() *RegistryMockBuilder {
	return &RegistryMockBuilder{
		mock: NewMockRegistry(),
	}
}

// WithItem adds an item to the mock registry.
func (b *RegistryMockBuilder) WithItem(item registry.Identifiable) *RegistryMockBuilder {
	b.mock.items[item.ID()] = item
	return b
}

// WithReturnItem configures the mock to return a specific item for a given ID.
func (b *RegistryMockBuilder) WithReturnItem(id string, item registry.Identifiable) *RegistryMockBuilder {
	b.mock.SetReturnItem(id, item)
	return b
}

// WithFailure configures the mock to fail operations.
func (b *RegistryMockBuilder) WithFailure(fail bool) *RegistryMockBuilder {
	b.mock.SetShouldFail(fail)
	return b
}

// WithFailureError sets the error message for failed operations.
func (b *RegistryMockBuilder) WithFailureError(err string) *RegistryMockBuilder {
	b.mock.SetFailureError(err)
	return b
}

// Build returns the configured mock registry.
func (b *RegistryMockBuilder) Build() registry.Registry {
	return b.mock
}
