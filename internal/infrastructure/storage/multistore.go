package storage

import (
	"errors"
	"path/filepath"
	"sync"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/storage"
	infraComponent "github.com/fintechain/skeleton/internal/infrastructure/component"
)

// MultiStore implements the MultiStoreService interface.
type MultiStore struct {
	*infraComponent.BaseService
	stores   map[string]storage.Store
	engines  map[string]storage.Engine
	mu       sync.RWMutex
	rootPath string
}

// NewMultiStore creates a new multi-store service.
func NewMultiStore(config component.ComponentConfig, rootPath string) *MultiStore {
	return &MultiStore{
		BaseService: infraComponent.NewBaseService(config),
		stores:      make(map[string]storage.Store),
		engines:     make(map[string]storage.Engine),
		rootPath:    rootPath,
	}
}

// GetStore retrieves a store by name.
func (ms *MultiStore) GetStore(name string) (storage.Store, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	store, exists := ms.stores[name]
	if !exists {
		return nil, errors.New(storage.ErrStoreNotFound)
	}

	return store, nil
}

// CreateStore creates a new store with the specified engine and configuration.
func (ms *MultiStore) CreateStore(name, engineName string, config storage.Config) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	// Check if store already exists
	if _, exists := ms.stores[name]; exists {
		return errors.New(storage.ErrStoreExists)
	}

	// Get the engine
	engine, exists := ms.engines[engineName]
	if !exists {
		return errors.New(storage.ErrEngineNotFound)
	}

	// Create store path
	storePath := filepath.Join(ms.rootPath, name)

	// Create the store
	store, err := engine.Create(name, storePath, config)
	if err != nil {
		return err
	}

	// Register the store
	ms.stores[name] = store
	return nil
}

// DeleteStore removes a store by name.
func (ms *MultiStore) DeleteStore(name string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	store, exists := ms.stores[name]
	if !exists {
		return errors.New(storage.ErrStoreNotFound)
	}

	// Close the store
	if err := store.Close(); err != nil {
		return err
	}

	// Remove from registry
	delete(ms.stores, name)
	return nil
}

// ListStores returns a list of all store names.
func (ms *MultiStore) ListStores() []string {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	names := make([]string, 0, len(ms.stores))
	for name := range ms.stores {
		names = append(names, name)
	}

	return names
}

// CloseAll closes all stores.
func (ms *MultiStore) CloseAll() error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var lastErr error
	for name, store := range ms.stores {
		if err := store.Close(); err != nil {
			lastErr = err // Continue closing others, return last error
		}
		delete(ms.stores, name)
	}

	return lastErr
}

// RegisterEngine registers a storage engine.
func (ms *MultiStore) RegisterEngine(engine storage.Engine) error {
	if engine == nil {
		return errors.New(storage.ErrInvalidConfig)
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	name := engine.Name()
	if name == "" {
		return errors.New(storage.ErrInvalidConfig)
	}

	// Check if engine already exists
	if _, exists := ms.engines[name]; exists {
		return errors.New(storage.ErrEngineExists)
	}

	ms.engines[name] = engine
	return nil
}

// Stop stops the multi-store service and closes all stores.
func (ms *MultiStore) Stop(ctx context.Context) error {
	if err := ms.CloseAll(); err != nil {
		return err
	}

	return ms.BaseService.Stop(ctx)
}
