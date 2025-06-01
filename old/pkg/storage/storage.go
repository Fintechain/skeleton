// Package storage provides public APIs for the storage system.
package storage

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/storage"
	"github.com/fintechain/skeleton/internal/infrastructure/event"
	"github.com/fintechain/skeleton/internal/infrastructure/logging"
	infraStorage "github.com/fintechain/skeleton/internal/infrastructure/storage"
	"github.com/fintechain/skeleton/internal/infrastructure/storage/memory"
)

// ===== CORE STORAGE INTERFACES =====

// Store defines the core storage operations that all backends must implement.
type Store = storage.Store

// MultiStore manages multiple named stores.
type MultiStore = storage.MultiStore

// Engine interface for storage backend implementations.
type Engine = storage.Engine

// ===== ADVANCED STORAGE INTERFACES =====

// Transaction represents an atomic set of operations.
type Transaction = storage.Transaction

// Transactional interface for stores that support transactions.
type Transactional = storage.Transactional

// Versioned interface for stores that support versioning/snapshots.
type Versioned = storage.Versioned

// RangeQueryable interface for stores that support range queries.
type RangeQueryable = storage.RangeQueryable

// Capabilities describes what features an engine supports.
type Capabilities = storage.Capabilities

// ===== CONFIGURATION TYPES =====

// Config type alias for storage configuration.
type Config = storage.Config

// StoreConfig provides access to store-specific configuration.
type StoreConfig = storage.StoreConfig

// DefaultStoreConfig implements StoreConfig interface.
type DefaultStoreConfig = storage.DefaultStoreConfig

// MultiStoreConfig defines configuration for MultiStore.
type MultiStoreConfig = storage.MultiStoreConfig

// ===== STORAGE ERROR CONSTANTS =====

// Common storage error constants
const (
	ErrKeyNotFound       = "storage.key_not_found"
	ErrStoreNotFound     = "storage.store_not_found"
	ErrStoreExists       = "storage.store_exists"
	ErrEngineNotFound    = "storage.engine_not_found"
	ErrVersionNotFound   = "storage.version_not_found"
	ErrTransactionFailed = "storage.transaction_failed"
	ErrInvalidOperation  = "storage.invalid_operation"
	ErrInvalidInput      = "storage.invalid_input"
	ErrInvalidConfig     = "storage.invalid_config"
)

// ===== ERROR HANDLING =====

// Error represents a domain-specific error from the storage system.
type Error = component.Error

// NewError creates a new storage error with the given code, message, and optional cause.
func NewError(code, message string, cause error) *Error {
	return component.NewError(code, message, cause)
}

// IsStorageError checks if an error is a storage error with the given code.
func IsStorageError(err error, code string) bool {
	return component.IsComponentError(err, code)
}

// ===== STORAGE CONSTRUCTORS =====

// NewMultiStore creates a new MultiStore with the given configuration.
// This is the primary way to create a MultiStore instance for managing multiple named stores.
// Requires a logger and event bus for proper operation.
//
// Example usage:
//
//	config := &storage.MultiStoreConfig{
//	    RootPath:      "./data",
//	    DefaultEngine: "memory",
//	}
//	logger := logging.NewLogrusLogger(logging.LevelInfo)
//	eventBus := event.NewEventBus()
//	multiStore := storage.NewMultiStore(config, logger, eventBus)
func NewMultiStore(config *MultiStoreConfig, logger logging.Logger, eventBus event.EventBus) MultiStore {
	return infraStorage.NewMultiStore(config, logger, eventBus)
}

// NewMemoryEngine creates a new in-memory storage engine.
// This engine stores data in memory and is suitable for testing or temporary storage.
// Requires a logger for proper operation.
//
// Example usage:
//
//	logger := logging.NewLogrusLogger(logging.LevelInfo)
//	engine := storage.NewMemoryEngine(logger)
func NewMemoryEngine(logger logging.Logger) Engine {
	return memory.NewEngine(logger)
}

// NewDefaultStoreConfig creates a default store configuration with the given engine and path.
// This provides a convenient way to create store configuration.
//
// Example usage:
//
//	config := storage.NewDefaultStoreConfig("memory", "./data/store")
func NewDefaultStoreConfig(engine, path string) *DefaultStoreConfig {
	return &DefaultStoreConfig{
		Engine:  engine,
		Path:    path,
		Options: make(Config),
	}
}

// ===== STORAGE UTILITIES =====

// CreateStore creates a new store with the given name and configuration.
// This is a convenience function for creating stores without complex setup.
//
// Example usage:
//
//	store, err := storage.CreateStore("user-data", "memory", "./data")
//	if err != nil {
//	    // Handle error
//	}
func CreateStore(name, engine, path string) (Store, error) {
	if name == "" {
		return nil, NewError(ErrInvalidInput, "store name cannot be empty", nil)
	}
	if engine == "" {
		return nil, NewError(ErrInvalidInput, "engine type cannot be empty", nil)
	}

	// Create a basic logger for the store
	logger := logging.CreateStandardLogger(logging.Info)

	// Validate engine type
	switch engine {
	case "memory":
		// Valid engine type
	default:
		return nil, NewError(ErrEngineNotFound, "unsupported engine type", nil).
			WithDetail("engine", engine)
	}

	// Create multistore config
	config := &MultiStoreConfig{
		RootPath:      path,
		DefaultEngine: engine,
	}

	// Create event bus
	eventBus := event.NewEventBus()

	// Create multistore
	multiStore := NewMultiStore(config, logger, eventBus)

	// Get or create the named store
	store, err := multiStore.GetStore(name)
	if err != nil {
		return nil, NewError(ErrStoreNotFound, "failed to get store", err).
			WithDetail("storeName", name)
	}

	return store, nil
}

// Get is a convenience function to retrieve a value from a store with error handling.
//
// Example usage:
//
//	value, err := storage.Get(store, []byte("user:123"))
//	if err != nil {
//	    if storage.IsStorageError(err, storage.ErrKeyNotFound) {
//	        // Handle key not found
//	    } else {
//	        // Handle other errors
//	    }
//	}
func Get(store Store, key []byte) ([]byte, error) {
	if store == nil {
		return nil, NewError(ErrInvalidInput, "store cannot be nil", nil)
	}
	if len(key) == 0 {
		return nil, NewError(ErrInvalidInput, "key cannot be empty", nil)
	}

	value, err := store.Get(key)
	if err != nil {
		return nil, NewError(ErrKeyNotFound, "failed to get value", err).
			WithDetail("key", string(key))
	}

	return value, nil
}

// Set is a convenience function to store a value with error handling.
//
// Example usage:
//
//	err := storage.Set(store, []byte("user:123"), userData)
//	if err != nil {
//	    // Handle error
//	}
func Set(store Store, key, value []byte) error {
	if store == nil {
		return NewError(ErrInvalidInput, "store cannot be nil", nil)
	}
	if len(key) == 0 {
		return NewError(ErrInvalidInput, "key cannot be empty", nil)
	}

	err := store.Set(key, value)
	if err != nil {
		return NewError(ErrInvalidOperation, "failed to set value", err).
			WithDetail("key", string(key))
	}

	return nil
}

// Delete is a convenience function to delete a value with error handling.
//
// Example usage:
//
//	err := storage.Delete(store, []byte("user:123"))
//	if err != nil {
//	    // Handle error
//	}
func Delete(store Store, key []byte) error {
	if store == nil {
		return NewError(ErrInvalidInput, "store cannot be nil", nil)
	}
	if len(key) == 0 {
		return NewError(ErrInvalidInput, "key cannot be empty", nil)
	}

	err := store.Delete(key)
	if err != nil {
		return NewError(ErrInvalidOperation, "failed to delete value", err).
			WithDetail("key", string(key))
	}

	return nil
}

// WithTransaction executes a function within a transaction if the store supports it.
//
// Example usage:
//
//	err := storage.WithTransaction(store, func(tx storage.Transaction) error {
//	    // Perform multiple operations within the transaction
//	    tx.Set([]byte("key1"), []byte("value1"))
//	    tx.Set([]byte("key2"), []byte("value2"))
//	    return nil
//	})
func WithTransaction(store Store, fn func(Transaction) error) error {
	if store == nil {
		return NewError(ErrInvalidInput, "store cannot be nil", nil)
	}

	transactional, ok := store.(Transactional)
	if !ok {
		return NewError(ErrInvalidOperation, "store does not support transactions", nil)
	}

	tx, err := transactional.BeginTx()
	if err != nil {
		return NewError(ErrTransactionFailed, "failed to begin transaction", err)
	}

	// Execute the function
	if err := fn(tx); err != nil {
		// Rollback on error
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return NewError(ErrTransactionFailed, "failed to rollback transaction", rollbackErr).
				WithDetail("originalError", err.Error())
		}
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return NewError(ErrTransactionFailed, "failed to commit transaction", err)
	}

	return nil
}
