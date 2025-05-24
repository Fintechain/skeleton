package mocks

import (
	"github.com/ebanfa/skeleton/internal/domain/storage"
)

// MockMultiStore is a mock implementation of storage.MultiStore for testing
type MockMultiStore struct {
	// Function fields for customizing behavior
	GetStoreFunc         func(string) (storage.Store, error)
	CreateStoreFunc      func(string, string, storage.Config) error
	DeleteStoreFunc      func(string) error
	ListStoresFunc       func() []string
	StoreExistsFunc      func(string) bool
	CloseAllFunc         func() error
	SetDefaultEngineFunc func(string)
	GetDefaultEngineFunc func() string
	RegisterEngineFunc   func(storage.Engine) error
	ListEnginesFunc      func() []string
	GetEngineFunc        func(string) (storage.Engine, error)

	// Call tracking
	GetStoreCalls         []string
	CreateStoreCalls      []CreateStoreCall
	DeleteStoreCalls      []string
	ListStoresCalls       int
	StoreExistsCalls      []string
	CloseAllCalls         int
	SetDefaultEngineCalls []string
	GetDefaultEngineCalls int
	RegisterEngineCalls   []storage.Engine
	ListEnginesCalls      int
	GetEngineCalls        []string

	// State
	Stores        map[string]storage.Store
	Engines       map[string]storage.Engine
	DefaultEngine string
}

type CreateStoreCall struct {
	Name   string
	Engine string
	Config storage.Config
}

// NewMockMultiStore creates a new mock multistore
func NewMockMultiStore() *MockMultiStore {
	return &MockMultiStore{
		Stores:        make(map[string]storage.Store),
		Engines:       make(map[string]storage.Engine),
		DefaultEngine: "memory",
	}
}

// GetStore implements storage.MultiStore
func (m *MockMultiStore) GetStore(name string) (storage.Store, error) {
	m.GetStoreCalls = append(m.GetStoreCalls, name)
	if m.GetStoreFunc != nil {
		return m.GetStoreFunc(name)
	}
	if store, exists := m.Stores[name]; exists {
		return store, nil
	}
	return nil, storage.WrapError(storage.ErrStoreNotFound, "store not found")
}

// CreateStore implements storage.MultiStore
func (m *MockMultiStore) CreateStore(name, engine string, config storage.Config) error {
	m.CreateStoreCalls = append(m.CreateStoreCalls, CreateStoreCall{Name: name, Engine: engine, Config: config})
	if m.CreateStoreFunc != nil {
		return m.CreateStoreFunc(name, engine, config)
	}
	if _, exists := m.Stores[name]; exists {
		return storage.WrapError(storage.ErrStoreExists, "store already exists")
	}
	// Create a mock store
	m.Stores[name] = &MockStore{StoreName: name}
	return nil
}

// DeleteStore implements storage.MultiStore
func (m *MockMultiStore) DeleteStore(name string) error {
	m.DeleteStoreCalls = append(m.DeleteStoreCalls, name)
	if m.DeleteStoreFunc != nil {
		return m.DeleteStoreFunc(name)
	}
	if _, exists := m.Stores[name]; !exists {
		return storage.WrapError(storage.ErrStoreNotFound, "store not found")
	}
	delete(m.Stores, name)
	return nil
}

// ListStores implements storage.MultiStore
func (m *MockMultiStore) ListStores() []string {
	m.ListStoresCalls++
	if m.ListStoresFunc != nil {
		return m.ListStoresFunc()
	}
	var names []string
	for name := range m.Stores {
		names = append(names, name)
	}
	return names
}

// StoreExists implements storage.MultiStore
func (m *MockMultiStore) StoreExists(name string) bool {
	m.StoreExistsCalls = append(m.StoreExistsCalls, name)
	if m.StoreExistsFunc != nil {
		return m.StoreExistsFunc(name)
	}
	_, exists := m.Stores[name]
	return exists
}

// CloseAll implements storage.MultiStore
func (m *MockMultiStore) CloseAll() error {
	m.CloseAllCalls++
	if m.CloseAllFunc != nil {
		return m.CloseAllFunc()
	}
	return nil
}

// SetDefaultEngine implements storage.MultiStore
func (m *MockMultiStore) SetDefaultEngine(engine string) {
	m.SetDefaultEngineCalls = append(m.SetDefaultEngineCalls, engine)
	if m.SetDefaultEngineFunc != nil {
		m.SetDefaultEngineFunc(engine)
		return
	}
	m.DefaultEngine = engine
}

// GetDefaultEngine implements storage.MultiStore
func (m *MockMultiStore) GetDefaultEngine() string {
	m.GetDefaultEngineCalls++
	if m.GetDefaultEngineFunc != nil {
		return m.GetDefaultEngineFunc()
	}
	return m.DefaultEngine
}

// RegisterEngine implements storage.MultiStore
func (m *MockMultiStore) RegisterEngine(engine storage.Engine) error {
	m.RegisterEngineCalls = append(m.RegisterEngineCalls, engine)
	if m.RegisterEngineFunc != nil {
		return m.RegisterEngineFunc(engine)
	}
	m.Engines[engine.Name()] = engine
	return nil
}

// ListEngines implements storage.MultiStore
func (m *MockMultiStore) ListEngines() []string {
	m.ListEnginesCalls++
	if m.ListEnginesFunc != nil {
		return m.ListEnginesFunc()
	}
	var names []string
	for name := range m.Engines {
		names = append(names, name)
	}
	return names
}

// GetEngine implements storage.MultiStore
func (m *MockMultiStore) GetEngine(name string) (storage.Engine, error) {
	m.GetEngineCalls = append(m.GetEngineCalls, name)
	if m.GetEngineFunc != nil {
		return m.GetEngineFunc(name)
	}
	if engine, exists := m.Engines[name]; exists {
		return engine, nil
	}
	return nil, storage.WrapError(storage.ErrEngineNotFound, "engine not found")
}

// MockStore is a simple mock implementation of storage.Store
type MockStore struct {
	StoreName string
	Data      map[string][]byte
}

// Get implements storage.Store
func (s *MockStore) Get(key []byte) ([]byte, error) {
	if s.Data == nil {
		s.Data = make(map[string][]byte)
	}
	if value, exists := s.Data[string(key)]; exists {
		return value, nil
	}
	return nil, storage.WrapError(storage.ErrKeyNotFound, "key not found")
}

// Set implements storage.Store
func (s *MockStore) Set(key, value []byte) error {
	if s.Data == nil {
		s.Data = make(map[string][]byte)
	}
	s.Data[string(key)] = value
	return nil
}

// Delete implements storage.Store
func (s *MockStore) Delete(key []byte) error {
	if s.Data == nil {
		s.Data = make(map[string][]byte)
	}
	delete(s.Data, string(key))
	return nil
}

// Has implements storage.Store
func (s *MockStore) Has(key []byte) (bool, error) {
	if s.Data == nil {
		s.Data = make(map[string][]byte)
	}
	_, exists := s.Data[string(key)]
	return exists, nil
}

// Iterate implements storage.Store
func (s *MockStore) Iterate(fn func(key, value []byte) bool) error {
	if s.Data == nil {
		s.Data = make(map[string][]byte)
	}
	for k, v := range s.Data {
		if !fn([]byte(k), v) {
			break
		}
	}
	return nil
}

// Close implements storage.Store
func (s *MockStore) Close() error {
	return nil
}

// Name implements storage.Store
func (s *MockStore) Name() string {
	return s.StoreName
}

// Path implements storage.Store
func (s *MockStore) Path() string {
	return "/mock/path/" + s.StoreName
}
