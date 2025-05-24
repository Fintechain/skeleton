// Package storage provides concrete implementations of storage interfaces.
package storage

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/ebanfa/skeleton/internal/domain/storage"
	"github.com/ebanfa/skeleton/internal/infrastructure/event"
	"github.com/ebanfa/skeleton/internal/infrastructure/logging"
	"github.com/ebanfa/skeleton/internal/infrastructure/storage/memory"
)

// DefaultMultiStore implements the domain MultiStore interface.
type DefaultMultiStore struct {
	stores        map[string]storage.Store
	engines       map[string]storage.Engine
	defaultEngine string
	rootPath      string
	logger        logging.Logger
	eventBus      event.EventBus
	mutex         sync.RWMutex
}

// NewMultiStore creates a new MultiStore instance with the given configuration.
func NewMultiStore(config *storage.MultiStoreConfig, logger logging.Logger, eventBus event.EventBus) storage.MultiStore {
	if logger == nil {
		panic("logger dependency cannot be nil")
	}

	if config == nil {
		config = &storage.MultiStoreConfig{
			RootPath:      "./data",
			DefaultEngine: "memory",
			EngineConfigs: make(map[string]storage.Config),
		}
	}

	ms := &DefaultMultiStore{
		stores:        make(map[string]storage.Store),
		engines:       make(map[string]storage.Engine),
		defaultEngine: config.DefaultEngine,
		rootPath:      config.RootPath,
		logger:        logger,
		eventBus:      eventBus,
	}

	// Register built-in engines
	ms.registerBuiltInEngines()

	return ms
}

// registerBuiltInEngines registers the built-in storage engines.
func (ms *DefaultMultiStore) registerBuiltInEngines() {
	// Register memory engine
	memoryEngine := memory.NewEngine(ms.logger)
	if err := ms.RegisterEngine(memoryEngine); err != nil {
		ms.logger.Error("Failed to register memory engine: %v", err)
	}
}

// GetStore retrieves a store by name.
func (ms *DefaultMultiStore) GetStore(name string) (storage.Store, error) {
	if name == "" {
		return nil, storage.WrapError(storage.ErrInvalidConfig, "store name cannot be empty")
	}

	ms.mutex.RLock()
	store, exists := ms.stores[name]
	ms.mutex.RUnlock()

	if !exists {
		return nil, storage.WrapError(storage.ErrStoreNotFound, fmt.Sprintf("store '%s' not found", name))
	}

	ms.logger.Debug("Retrieved store: %s", name)
	return store, nil
}

// CreateStore creates a new store with the given name and engine.
func (ms *DefaultMultiStore) CreateStore(name, engineName string, config storage.Config) error {
	if name == "" {
		return storage.WrapError(storage.ErrInvalidConfig, "store name cannot be empty")
	}

	// Use default engine if none specified
	if engineName == "" {
		engineName = ms.defaultEngine
	}

	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	// Check if store already exists
	if _, exists := ms.stores[name]; exists {
		return storage.WrapError(storage.ErrStoreExists, fmt.Sprintf("store '%s' already exists", name))
	}

	// Get the engine
	engine, exists := ms.engines[engineName]
	if !exists {
		return storage.WrapError(storage.ErrEngineNotFound, fmt.Sprintf("engine '%s' not found", engineName))
	}

	// Create store path
	storePath := filepath.Join(ms.rootPath, name)

	// Create the store
	store, err := engine.Create(name, storePath, config)
	if err != nil {
		return storage.WrapError(err, fmt.Sprintf("failed to create store '%s'", name))
	}

	// Register the store
	ms.stores[name] = store

	ms.logger.Info("Created store '%s' using engine '%s'", name, engineName)

	// Publish store created event
	if ms.eventBus != nil {
		payload := storage.CreateStoreEventPayload(name, engineName, nil)
		ms.eventBus.Publish(storage.TopicStoreCreated, payload)
	}

	return nil
}

// DeleteStore removes a store by name.
func (ms *DefaultMultiStore) DeleteStore(name string) error {
	if name == "" {
		return storage.WrapError(storage.ErrInvalidConfig, "store name cannot be empty")
	}

	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	store, exists := ms.stores[name]
	if !exists {
		return storage.WrapError(storage.ErrStoreNotFound, fmt.Sprintf("store '%s' not found", name))
	}

	// Get engine name for event payload
	engineName := "unknown"
	for _, engine := range ms.engines {
		// Try to match by checking if this engine could have created this store
		// This is a best-effort approach since we don't store engine info with stores
		if engine.Name() != "" {
			engineName = engine.Name()
			break
		}
	}

	// Close the store first
	if err := store.Close(); err != nil {
		ms.logger.Warn("Failed to close store '%s' during deletion: %v", name, err)
	}

	// Remove from registry
	delete(ms.stores, name)

	ms.logger.Info("Deleted store: %s", name)

	// Publish store deleted event
	if ms.eventBus != nil {
		payload := storage.CreateStoreEventPayload(name, engineName, nil)
		ms.eventBus.Publish(storage.TopicStoreDeleted, payload)
	}

	return nil
}

// ListStores returns the names of all stores.
func (ms *DefaultMultiStore) ListStores() []string {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	names := make([]string, 0, len(ms.stores))
	for name := range ms.stores {
		names = append(names, name)
	}

	return names
}

// StoreExists checks if a store exists.
func (ms *DefaultMultiStore) StoreExists(name string) bool {
	if name == "" {
		return false
	}

	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	_, exists := ms.stores[name]
	return exists
}

// CloseAll closes all stores.
func (ms *DefaultMultiStore) CloseAll() error {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	var lastError error

	for name, store := range ms.stores {
		if err := store.Close(); err != nil {
			ms.logger.Error("Failed to close store '%s': %v", name, err)
			lastError = err
		} else {
			// Publish store closed event for each successfully closed store
			if ms.eventBus != nil {
				payload := storage.CreateStoreEventPayload(name, "unknown", nil)
				ms.eventBus.Publish(storage.TopicStoreClosed, payload)
			}
		}
	}

	// Clear the stores map
	ms.stores = make(map[string]storage.Store)

	if lastError != nil {
		return storage.WrapError(lastError, "failed to close some stores")
	}

	ms.logger.Info("Closed all stores")
	return nil
}

// SetDefaultEngine sets the default engine to use when no engine is specified.
func (ms *DefaultMultiStore) SetDefaultEngine(engine string) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	ms.defaultEngine = engine
	ms.logger.Debug("Set default engine to: %s", engine)
}

// GetDefaultEngine returns the current default engine.
func (ms *DefaultMultiStore) GetDefaultEngine() string {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	return ms.defaultEngine
}

// RegisterEngine registers a new engine.
func (ms *DefaultMultiStore) RegisterEngine(engine storage.Engine) error {
	if engine == nil {
		return storage.WrapError(storage.ErrInvalidConfig, "engine cannot be nil")
	}

	engineName := engine.Name()
	if engineName == "" {
		return storage.WrapError(storage.ErrInvalidConfig, "engine name cannot be empty")
	}

	ms.mutex.Lock()
	defer ms.mutex.Unlock()

	if _, exists := ms.engines[engineName]; exists {
		return storage.WrapError(storage.ErrStoreExists, fmt.Sprintf("engine '%s' already registered", engineName))
	}

	ms.engines[engineName] = engine
	ms.logger.Info("Registered engine: %s", engineName)
	return nil
}

// ListEngines returns the names of all registered engines.
func (ms *DefaultMultiStore) ListEngines() []string {
	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	names := make([]string, 0, len(ms.engines))
	for name := range ms.engines {
		names = append(names, name)
	}

	return names
}

// GetEngine retrieves an engine by name.
func (ms *DefaultMultiStore) GetEngine(name string) (storage.Engine, error) {
	if name == "" {
		return nil, storage.WrapError(storage.ErrInvalidConfig, "engine name cannot be empty")
	}

	ms.mutex.RLock()
	defer ms.mutex.RUnlock()

	engine, exists := ms.engines[name]
	if !exists {
		return nil, storage.WrapError(storage.ErrEngineNotFound, fmt.Sprintf("engine '%s' not found", name))
	}

	return engine, nil
}
