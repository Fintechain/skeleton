package mocks

import (
	"fmt"
	"sync"

	"github.com/fintechain/skeleton/pkg/storage"
)

// MockEngine provides a configurable mock implementation of the storage.Engine interface.
type MockEngine struct {
	mu sync.RWMutex

	// Engine data
	name         string
	description  string
	version      string
	stores       map[string]storage.Store
	capabilities storage.Capabilities

	// Behavior configuration
	shouldFail   bool
	failureError string

	// Call tracking
	callCount map[string]int
	lastCalls map[string][]interface{}
}

// NewMockEngine creates a new configurable engine mock.
func NewMockEngine() *MockEngine {
	return &MockEngine{
		name:        "mock-engine",
		description: "Mock storage engine for testing",
		version:     "1.0.0",
		stores:      make(map[string]storage.Store),
		capabilities: storage.Capabilities{
			Transactions: true,
			Versioning:   true,
			RangeQueries: true,
			Persistence:  false,
			Compression:  false,
		},
		callCount: make(map[string]int),
		lastCalls: make(map[string][]interface{}),
	}
}

// Engine Interface Implementation

// Name returns the engine name.
func (e *MockEngine) Name() string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	e.trackCall("Name")
	return e.name
}

// Description returns the engine description.
func (e *MockEngine) Description() string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	e.trackCall("Description")
	return e.description
}

// Version returns the engine version.
func (e *MockEngine) Version() string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	e.trackCall("Version")
	return e.version
}

// Create creates a new store with the given configuration.
func (e *MockEngine) Create(name, path string, config storage.Config) (storage.Store, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.trackCall("Create", name, path, config)

	if e.shouldFail {
		return nil, fmt.Errorf("%s", e.getFailureError("Create"))
	}

	if _, exists := e.stores[name]; exists {
		return nil, fmt.Errorf("store already exists: %s", name)
	}

	mockStore := NewMockStore()
	mockStore.SetName(name)
	e.stores[name] = mockStore
	return mockStore, nil
}

// Open opens an existing store.
func (e *MockEngine) Open(name, path string) (storage.Store, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	e.trackCall("Open", name, path)

	if e.shouldFail {
		return nil, fmt.Errorf("%s", e.getFailureError("Open"))
	}

	if store, exists := e.stores[name]; exists {
		return store, nil
	}

	return nil, fmt.Errorf("store not found: %s", name)
}

// DeleteStore deletes a store.
func (e *MockEngine) DeleteStore(name string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.trackCall("DeleteStore", name)

	if e.shouldFail {
		return fmt.Errorf("%s", e.getFailureError("DeleteStore"))
	}

	if _, exists := e.stores[name]; !exists {
		return fmt.Errorf("store not found: %s", name)
	}

	delete(e.stores, name)
	return nil
}

// ListStores returns all store names managed by this engine.
func (e *MockEngine) ListStores() []string {
	e.mu.RLock()
	defer e.mu.RUnlock()

	e.trackCall("ListStores")

	names := make([]string, 0, len(e.stores))
	for name := range e.stores {
		names = append(names, name)
	}
	return names
}

// StoreExists checks if a store exists.
func (e *MockEngine) StoreExists(name string) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	e.trackCall("StoreExists", name)

	_, exists := e.stores[name]
	return exists
}

// Close closes the engine and all its stores.
func (e *MockEngine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.trackCall("Close")

	if e.shouldFail {
		return fmt.Errorf("%s", e.getFailureError("Close"))
	}

	for _, store := range e.stores {
		if err := store.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Capabilities returns what features this engine supports.
func (e *MockEngine) Capabilities() storage.Capabilities {
	e.mu.RLock()
	defer e.mu.RUnlock()

	e.trackCall("Capabilities")
	return e.capabilities
}

// Mock Configuration Methods

// SetName sets the engine name.
func (e *MockEngine) SetName(name string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.name = name
}

// SetDescription sets the engine description.
func (e *MockEngine) SetDescription(description string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.description = description
}

// SetVersion sets the engine version.
func (e *MockEngine) SetVersion(version string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.version = version
}

// SetCapabilities sets the engine capabilities.
func (e *MockEngine) SetCapabilities(capabilities storage.Capabilities) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.capabilities = capabilities
}

// AddStore adds a store to the engine.
func (e *MockEngine) AddStore(name string, store storage.Store) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.stores[name] = store
}

// SetShouldFail configures the mock to fail operations.
func (e *MockEngine) SetShouldFail(fail bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.shouldFail = fail
}

// SetFailureError sets the error message for failed operations.
func (e *MockEngine) SetFailureError(err string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.failureError = err
}

// State Verification Methods

// GetCallCount returns the number of times a method was called.
func (e *MockEngine) GetCallCount(method string) int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.callCount[method]
}

// GetLastCall returns the parameters of the last call to a method.
func (e *MockEngine) GetLastCall(method string) []interface{} {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.lastCalls[method]
}

// GetStores returns all stores managed by this engine.
func (e *MockEngine) GetStores() map[string]storage.Store {
	e.mu.RLock()
	defer e.mu.RUnlock()

	stores := make(map[string]storage.Store)
	for name, store := range e.stores {
		stores[name] = store
	}
	return stores
}

// Reset clears all mock state and configuration.
func (e *MockEngine) Reset() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.name = "mock-engine"
	e.description = "Mock storage engine for testing"
	e.version = "1.0.0"
	e.stores = make(map[string]storage.Store)
	e.capabilities = storage.Capabilities{
		Transactions: true,
		Versioning:   true,
		RangeQueries: true,
		Persistence:  false,
		Compression:  false,
	}
	e.shouldFail = false
	e.failureError = ""
	e.callCount = make(map[string]int)
	e.lastCalls = make(map[string][]interface{})
}

// Helper Methods

// trackCall records a method call for verification.
func (e *MockEngine) trackCall(method string, args ...interface{}) {
	e.callCount[method]++
	e.lastCalls[method] = args
}

// getFailureError returns the configured failure error or a default.
func (e *MockEngine) getFailureError(method string) string {
	if e.failureError != "" {
		return e.failureError
	}
	return fmt.Sprintf("mock engine %s failed", method)
}

// EngineMockBuilder provides a fluent interface for configuring engine mocks.
type EngineMockBuilder struct {
	mock *MockEngine
}

// NewEngineMockBuilder creates a new engine mock builder.
func NewEngineMockBuilder() *EngineMockBuilder {
	return &EngineMockBuilder{
		mock: NewMockEngine(),
	}
}

// WithName sets the engine name.
func (b *EngineMockBuilder) WithName(name string) *EngineMockBuilder {
	b.mock.SetName(name)
	return b
}

// WithDescription sets the engine description.
func (b *EngineMockBuilder) WithDescription(description string) *EngineMockBuilder {
	b.mock.SetDescription(description)
	return b
}

// WithVersion sets the engine version.
func (b *EngineMockBuilder) WithVersion(version string) *EngineMockBuilder {
	b.mock.SetVersion(version)
	return b
}

// WithCapabilities sets the engine capabilities.
func (b *EngineMockBuilder) WithCapabilities(capabilities storage.Capabilities) *EngineMockBuilder {
	b.mock.SetCapabilities(capabilities)
	return b
}

// WithStore adds a store to the engine.
func (b *EngineMockBuilder) WithStore(name string, store storage.Store) *EngineMockBuilder {
	b.mock.AddStore(name, store)
	return b
}

// WithFailure configures the mock to fail operations.
func (b *EngineMockBuilder) WithFailure(fail bool) *EngineMockBuilder {
	b.mock.SetShouldFail(fail)
	return b
}

// WithFailureError sets the error message for failed operations.
func (b *EngineMockBuilder) WithFailureError(err string) *EngineMockBuilder {
	b.mock.SetFailureError(err)
	return b
}

// Build returns the configured engine mock.
func (b *EngineMockBuilder) Build() storage.Engine {
	return b.mock
}
