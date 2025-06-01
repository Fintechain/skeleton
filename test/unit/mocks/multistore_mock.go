package mocks

import (
	"fmt"
	"sync"

	"github.com/fintechain/skeleton/pkg/storage"
)

// MockMultiStore provides a configurable mock implementation of the storage.MultiStore interface.
type MockMultiStore struct {
	mu sync.RWMutex

	// Storage data
	stores        map[string]storage.Store
	defaultEngine string
	engines       map[string]storage.Engine

	// Behavior configuration
	shouldFail   bool
	failureError string

	// Call tracking
	callCount map[string]int
	lastCalls map[string][]interface{}
}

// NewMockMultiStore creates a new configurable multi-store mock.
func NewMockMultiStore() *MockMultiStore {
	return &MockMultiStore{
		stores:        make(map[string]storage.Store),
		engines:       make(map[string]storage.Engine),
		defaultEngine: "memory",
		callCount:     make(map[string]int),
		lastCalls:     make(map[string][]interface{}),
	}
}

// MultiStore Interface Implementation

// GetStore retrieves a store by name.
func (ms *MockMultiStore) GetStore(name string) (storage.Store, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	ms.trackCall("GetStore", name)

	if ms.shouldFail {
		return nil, fmt.Errorf("%s", ms.getFailureError("GetStore"))
	}

	if store, exists := ms.stores[name]; exists {
		return store, nil
	}

	return nil, fmt.Errorf("store not found: %s", name)
}

// CreateStore creates a new store with the given name and engine.
func (ms *MockMultiStore) CreateStore(name, engine string, config storage.Config) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.trackCall("CreateStore", name, engine, config)

	if ms.shouldFail {
		return fmt.Errorf("%s", ms.getFailureError("CreateStore"))
	}

	if _, exists := ms.stores[name]; exists {
		return fmt.Errorf("store already exists: %s", name)
	}

	// Create a mock store
	mockStore := NewMockStore()
	mockStore.SetName(name)
	ms.stores[name] = mockStore
	return nil
}

// DeleteStore removes a store by name.
func (ms *MockMultiStore) DeleteStore(name string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.trackCall("DeleteStore", name)

	if ms.shouldFail {
		return fmt.Errorf("%s", ms.getFailureError("DeleteStore"))
	}

	if _, exists := ms.stores[name]; !exists {
		return fmt.Errorf("store not found: %s", name)
	}

	delete(ms.stores, name)
	return nil
}

// ListStores returns all store names.
func (m *MockMultiStore) ListStores() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("ListStores")

	names := make([]string, 0, len(m.stores))
	for name := range m.stores {
		names = append(names, name)
	}
	return names
}

// StoreExists checks if a store exists.
func (m *MockMultiStore) StoreExists(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("StoreExists", name)

	_, exists := m.stores[name]
	return exists
}

// CloseAll closes all stores.
func (ms *MockMultiStore) CloseAll() error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.trackCall("CloseAll")

	if ms.shouldFail {
		return fmt.Errorf("%s", ms.getFailureError("CloseAll"))
	}

	for _, store := range ms.stores {
		if err := store.Close(); err != nil {
			return err
		}
	}
	return nil
}

// SetDefaultEngine sets the default engine.
func (m *MockMultiStore) SetDefaultEngine(engine string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("SetDefaultEngine", engine)
	m.defaultEngine = engine
}

// GetDefaultEngine returns the current default engine.
func (m *MockMultiStore) GetDefaultEngine() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("GetDefaultEngine")
	return m.defaultEngine
}

// RegisterEngine registers a new engine.
func (ms *MockMultiStore) RegisterEngine(engine storage.Engine) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.trackCall("RegisterEngine", engine)

	if ms.shouldFail {
		return fmt.Errorf("%s", ms.getFailureError("RegisterEngine"))
	}

	if engine == nil {
		return fmt.Errorf("engine cannot be nil")
	}

	name := engine.Name()
	if _, exists := ms.engines[name]; exists {
		return fmt.Errorf("engine already exists: %s", name)
	}

	ms.engines[name] = engine
	return nil
}

// ListEngines returns the names of all registered engines.
func (m *MockMultiStore) ListEngines() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("ListEngines")

	names := make([]string, 0, len(m.engines))
	for name := range m.engines {
		names = append(names, name)
	}
	return names
}

// GetEngine retrieves an engine by name.
func (ms *MockMultiStore) GetEngine(name string) (storage.Engine, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	ms.trackCall("GetEngine", name)

	if ms.shouldFail {
		return nil, fmt.Errorf("%s", ms.getFailureError("GetEngine"))
	}

	if engine, exists := ms.engines[name]; exists {
		return engine, nil
	}

	return nil, fmt.Errorf("engine not found: %s", name)
}

// Mock Configuration Methods

// AddStore adds a store to the mock manager.
func (m *MockMultiStore) AddStore(name string, store storage.Store) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.stores[name] = store
}

// AddEngine adds an engine to the mock manager.
func (m *MockMultiStore) AddEngine(name string, engine storage.Engine) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.engines[name] = engine
}

// SetShouldFail configures the mock to fail operations.
func (m *MockMultiStore) SetShouldFail(fail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail = fail
}

// SetFailureError sets the error message for failed operations.
func (m *MockMultiStore) SetFailureError(err string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failureError = err
}

// State Verification Methods

// GetCallCount returns the number of times a method was called.
func (m *MockMultiStore) GetCallCount(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.callCount[method]
}

// GetLastCall returns the parameters of the last call to a method.
func (m *MockMultiStore) GetLastCall(method string) []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastCalls[method]
}

// Reset clears all mock state and configuration.
func (m *MockMultiStore) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stores = make(map[string]storage.Store)
	m.engines = make(map[string]storage.Engine)
	m.defaultEngine = "memory"
	m.shouldFail = false
	m.failureError = ""
	m.callCount = make(map[string]int)
	m.lastCalls = make(map[string][]interface{})
}

// Helper Methods

// trackCall records a method call for verification.
func (m *MockMultiStore) trackCall(method string, args ...interface{}) {
	m.callCount[method]++
	m.lastCalls[method] = args
}

// getFailureError returns the configured failure error or a default.
func (m *MockMultiStore) getFailureError(method string) string {
	if m.failureError != "" {
		return m.failureError
	}
	return fmt.Sprintf("mock multistore %s failed", method)
}

// MockStore provides a configurable mock implementation of the storage.Store interface.
type MockStore struct {
	mu sync.RWMutex

	// Storage data
	data map[string][]byte
	name string
	path string

	// State
	closed bool

	// Behavior configuration
	shouldFail   bool
	failureError string

	// Call tracking
	callCount map[string]int
	lastCalls map[string][]interface{}
}

// NewMockStore creates a new configurable store mock.
func NewMockStore() *MockStore {
	return &MockStore{
		data:      make(map[string][]byte),
		callCount: make(map[string]int),
		lastCalls: make(map[string][]interface{}),
	}
}

// Store Interface Implementation

// Get retrieves a value by key.
func (s *MockStore) Get(key []byte) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.trackCall("Get", key)

	if s.shouldFail {
		return nil, fmt.Errorf("%s", s.getFailureError("Get"))
	}

	if value, exists := s.data[string(key)]; exists {
		// Return a copy to prevent external modification
		result := make([]byte, len(value))
		copy(result, value)
		return result, nil
	}

	return nil, fmt.Errorf("key not found")
}

// Set stores a key-value pair.
func (s *MockStore) Set(key, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.trackCall("Set", key, value)

	if s.shouldFail {
		return fmt.Errorf("%s", s.getFailureError("Set"))
	}

	// Store a copy to prevent external modification
	valueCopy := make([]byte, len(value))
	copy(valueCopy, value)
	s.data[string(key)] = valueCopy
	return nil
}

// Delete removes a key-value pair.
func (s *MockStore) Delete(key []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.trackCall("Delete", key)

	if s.shouldFail {
		return fmt.Errorf("%s", s.getFailureError("Delete"))
	}

	delete(s.data, string(key))
	return nil
}

// Has checks if a key exists.
func (s *MockStore) Has(key []byte) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.trackCall("Has", key)

	if s.shouldFail {
		return false, fmt.Errorf("%s", s.getFailureError("Has"))
	}

	_, exists := s.data[string(key)]
	return exists, nil
}

// Iterate iterates over all key-value pairs.
func (s *MockStore) Iterate(fn func(key, value []byte) bool) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.trackCall("Iterate", fn)

	if s.shouldFail {
		return fmt.Errorf("%s", s.getFailureError("Iterate"))
	}

	for k, v := range s.data {
		if !fn([]byte(k), v) {
			break
		}
	}
	return nil
}

// Close closes the store.
func (s *MockStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.trackCall("Close")

	if s.shouldFail {
		return fmt.Errorf("%s", s.getFailureError("Close"))
	}

	s.closed = true
	return nil
}

// Name returns the store name.
func (s *MockStore) Name() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.trackCall("Name")
	return s.name
}

// Path returns the store path.
func (s *MockStore) Path() string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.trackCall("Path")
	return s.path
}

// Mock Configuration Methods

// SetKeyValue sets a key-value pair in the mock store.
func (s *MockStore) SetKeyValue(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = []byte(value)
}

// SetName sets the store name.
func (s *MockStore) SetName(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.name = name
}

// SetPath sets the store path.
func (s *MockStore) SetPath(path string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.path = path
}

// Helper Methods

// trackCall records a method call for verification.
func (s *MockStore) trackCall(method string, args ...interface{}) {
	s.callCount[method]++
	s.lastCalls[method] = args
}

// getFailureError returns the configured failure error or a default.
func (s *MockStore) getFailureError(method string) string {
	if s.failureError != "" {
		return s.failureError
	}
	return fmt.Sprintf("mock store %s failed", method)
}

// MultiStoreMockBuilder provides a fluent interface for configuring multistore mocks.
type MultiStoreMockBuilder struct {
	mock *MockMultiStore
}

// NewMultiStoreMockBuilder creates a new multistore mock builder.
func NewMultiStoreMockBuilder() *MultiStoreMockBuilder {
	return &MultiStoreMockBuilder{
		mock: NewMockMultiStore(),
	}
}

// WithStore adds a store to the multistore.
func (b *MultiStoreMockBuilder) WithStore(name string, store storage.Store) *MultiStoreMockBuilder {
	b.mock.AddStore(name, store)
	return b
}

// WithEngine adds an engine to the multistore.
func (b *MultiStoreMockBuilder) WithEngine(name string, engine storage.Engine) *MultiStoreMockBuilder {
	b.mock.AddEngine(name, engine)
	return b
}

// WithDefaultEngine sets the default engine.
func (b *MultiStoreMockBuilder) WithDefaultEngine(engine string) *MultiStoreMockBuilder {
	b.mock.SetDefaultEngine(engine)
	return b
}

// WithFailure configures the mock to fail operations.
func (b *MultiStoreMockBuilder) WithFailure(fail bool) *MultiStoreMockBuilder {
	b.mock.SetShouldFail(fail)
	return b
}

// WithFailureError sets the error message for failed operations.
func (b *MultiStoreMockBuilder) WithFailureError(err string) *MultiStoreMockBuilder {
	b.mock.SetFailureError(err)
	return b
}

// Build returns the configured multistore mock.
func (b *MultiStoreMockBuilder) Build() storage.MultiStore {
	return b.mock
}

// StoreMockBuilder provides a fluent interface for configuring store mocks.
type StoreMockBuilder struct {
	mock *MockStore
}

// NewStoreMockBuilder creates a new store mock builder.
func NewStoreMockBuilder() *StoreMockBuilder {
	return &StoreMockBuilder{
		mock: NewMockStore(),
	}
}

// WithKeyValue adds a key-value pair to the store.
func (b *StoreMockBuilder) WithKeyValue(key, value string) *StoreMockBuilder {
	b.mock.SetKeyValue(key, value)
	return b
}

// WithName sets the store name.
func (b *StoreMockBuilder) WithName(name string) *StoreMockBuilder {
	b.mock.SetName(name)
	return b
}

// WithPath sets the store path.
func (b *StoreMockBuilder) WithPath(path string) *StoreMockBuilder {
	b.mock.SetPath(path)
	return b
}

// WithFailure configures the mock to fail operations.
func (b *StoreMockBuilder) WithFailure(fail bool) *StoreMockBuilder {
	b.mock.shouldFail = fail
	return b
}

// WithFailureError sets the error message for failed operations.
func (b *StoreMockBuilder) WithFailureError(err string) *StoreMockBuilder {
	b.mock.failureError = err
	return b
}

// Build returns the configured store mock.
func (b *StoreMockBuilder) Build() storage.Store {
	return b.mock
}
