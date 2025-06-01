// Package storage provides concrete implementations of the storage system.
package storage

import (
	"errors"
	"fmt"
	"sync"

	"github.com/fintechain/skeleton/internal/domain/storage"
)

// DefaultMultiStore is a concrete implementation of the MultiStore interface.
type DefaultMultiStore struct {
	mu            sync.RWMutex
	stores        map[string]storage.Store
	engines       map[string]storage.Engine
	defaultEngine string
}

// NewMultiStore creates a new MultiStore instance that accepts engine interface dependencies.
func NewMultiStore() storage.MultiStore {
	return &DefaultMultiStore{
		stores:  make(map[string]storage.Store),
		engines: make(map[string]storage.Engine),
	}
}

// GetStore retrieves a store by name.
func (ms *DefaultMultiStore) GetStore(name string) (storage.Store, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	store, exists := ms.stores[name]
	if !exists {
		return nil, errors.New(storage.ErrStoreNotFound)
	}

	return store, nil
}

// CreateStore creates a new store with the given name and engine.
func (ms *DefaultMultiStore) CreateStore(name, engineName string, config storage.Config) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// Check if store already exists
	if _, exists := ms.stores[name]; exists {
		return errors.New(storage.ErrStoreExists)
	}

	// Use default engine if none specified
	if engineName == "" {
		engineName = ms.defaultEngine
	}

	// Get the engine
	engine, exists := ms.engines[engineName]
	if !exists {
		return errors.New(storage.ErrEngineNotFound)
	}

	// Create the store using the engine
	// For simplicity, we'll use the store name as the path
	store, err := engine.Create(name, name, config)
	if err != nil {
		return fmt.Errorf("failed to create store: %w", err)
	}

	// Register the store
	ms.stores[name] = store
	return nil
}

// DeleteStore removes a store by name.
func (ms *DefaultMultiStore) DeleteStore(name string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	store, exists := ms.stores[name]
	if !exists {
		return errors.New(storage.ErrStoreNotFound)
	}

	// Close the store first
	if err := store.Close(); err != nil {
		return fmt.Errorf("failed to close store before deletion: %w", err)
	}

	// Remove from registry
	delete(ms.stores, name)
	return nil
}

// ListStores returns the names of all stores.
func (ms *DefaultMultiStore) ListStores() []string {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	names := make([]string, 0, len(ms.stores))
	for name := range ms.stores {
		names = append(names, name)
	}
	return names
}

// StoreExists checks if a store exists.
func (ms *DefaultMultiStore) StoreExists(name string) bool {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	_, exists := ms.stores[name]
	return exists
}

// CloseAll closes all stores.
func (ms *DefaultMultiStore) CloseAll() error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var errs []error
	for name, store := range ms.stores {
		if err := store.Close(); err != nil {
			errs = append(errs, fmt.Errorf("failed to close store %s: %w", name, err))
		}
	}

	// Clear the stores map
	ms.stores = make(map[string]storage.Store)

	// Return combined errors if any
	if len(errs) > 0 {
		return fmt.Errorf("errors closing stores: %v", errs)
	}

	return nil
}

// SetDefaultEngine sets the default engine to use when no engine is specified.
func (ms *DefaultMultiStore) SetDefaultEngine(engine string) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	ms.defaultEngine = engine
}

// GetDefaultEngine returns the current default engine.
func (ms *DefaultMultiStore) GetDefaultEngine() string {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	return ms.defaultEngine
}

// RegisterEngine registers a new engine.
func (ms *DefaultMultiStore) RegisterEngine(engine storage.Engine) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	name := engine.Name()
	if _, exists := ms.engines[name]; exists {
		return fmt.Errorf("engine %s already registered", name)
	}

	ms.engines[name] = engine

	// Set as default if it's the first engine
	if ms.defaultEngine == "" {
		ms.defaultEngine = name
	}

	return nil
}

// ListEngines returns the names of all registered engines.
func (ms *DefaultMultiStore) ListEngines() []string {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	names := make([]string, 0, len(ms.engines))
	for name := range ms.engines {
		names = append(names, name)
	}
	return names
}

// GetEngine retrieves an engine by name.
func (ms *DefaultMultiStore) GetEngine(name string) (storage.Engine, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	engine, exists := ms.engines[name]
	if !exists {
		return nil, errors.New(storage.ErrEngineNotFound)
	}

	return engine, nil
}
