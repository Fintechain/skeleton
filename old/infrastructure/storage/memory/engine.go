// Package memory provides an in-memory implementation of the storage engine.
package memory

import (
	"fmt"
	"sync"

	"github.com/fintechain/skeleton/internal/domain/storage"
	"github.com/fintechain/skeleton/internal/infrastructure/logging"
)

// Engine implements storage.Engine for in-memory storage.
type Engine struct {
	logger logging.Logger
	mutex  sync.RWMutex
	stores map[string]*Store
}

// NewEngine creates a new memory storage engine.
// Follows constructor injection pattern by requiring all dependencies.
func NewEngine(logger logging.Logger) storage.Engine {
	if logger == nil {
		panic("logger dependency cannot be nil")
	}

	return &Engine{
		logger: logger,
		stores: make(map[string]*Store),
	}
}

// Name returns the engine identifier.
func (e *Engine) Name() string {
	return "memory"
}

// Create creates a new store instance.
func (e *Engine) Create(name, path string, config storage.Config) (storage.Store, error) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	// Check if store already exists
	if _, exists := e.stores[name]; exists {
		e.logger.Warn("Attempted to create store that already exists: %s", name)
		return nil, storage.WrapError(storage.ErrStoreExists, fmt.Sprintf("memory store %s already exists", name))
	}

	options := parseConfig(config)
	store := NewStore(name, path, options, e.logger)
	e.stores[name] = store

	e.logger.Debug("Created memory store: %s", name)
	return store, nil
}

// Open opens an existing store.
func (e *Engine) Open(name, path string) (storage.Store, error) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	// Check if store exists
	store, exists := e.stores[name]
	if !exists {
		return nil, storage.WrapError(storage.ErrStoreNotFound, fmt.Sprintf("memory store %s not found", name))
	}

	if store.IsClosed() {
		return nil, storage.WrapError(storage.ErrStoreClosed, fmt.Sprintf("memory store %s is closed", name))
	}

	e.logger.Debug("Opened memory store: %s", name)
	return store, nil
}

// Capabilities returns what this engine supports.
func (e *Engine) Capabilities() storage.Capabilities {
	return storage.Capabilities{
		Transactions: true,
		Versioning:   true,
		RangeQueries: true,
		Persistence:  false,
		Compression:  false,
	}
}

// Options for memory store configuration
type Options struct {
	MaxSize     int64
	MaxVersions int
}

// parseConfig converts the generic config to memory-specific options
func parseConfig(config storage.Config) Options {
	options := Options{
		MaxSize:     -1, // unlimited
		MaxVersions: 100,
	}

	if config == nil {
		return options
	}

	if maxSize, ok := config[storage.ConfigCacheSize].(int64); ok {
		options.MaxSize = maxSize
	} else if maxSizeInt, ok := config[storage.ConfigCacheSize].(int); ok {
		options.MaxSize = int64(maxSizeInt)
	}

	if maxVersions, ok := config[storage.ConfigMaxVersions].(int); ok {
		options.MaxVersions = maxVersions
	}

	return options
}
