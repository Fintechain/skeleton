// Package storage provides interfaces and types for the storage system.
package storage

// MultiStore manages multiple named stores.
// It provides a central registry for all stores and their engines.
type MultiStore interface {
	// Store management

	// GetStore retrieves a store by name.
	// Returns ErrStoreNotFound if the store doesn't exist.
	GetStore(name string) (Store, error)

	// CreateStore creates a new store with the given name and engine.
	// The config parameter contains engine-specific configuration options.
	// Returns ErrStoreExists if the store already exists.
	CreateStore(name, engine string, config Config) error

	// DeleteStore removes a store by name.
	// Returns ErrStoreNotFound if the store doesn't exist.
	DeleteStore(name string) error

	// ListStores returns the names of all stores.
	// Returns an empty slice if no stores exist.
	ListStores() []string

	// StoreExists checks if a store exists.
	StoreExists(name string) bool

	// Bulk operations

	// CloseAll closes all stores.
	// This should be called when shutting down the application.
	CloseAll() error

	// Configuration

	// SetDefaultEngine sets the default engine to use when no engine is specified.
	SetDefaultEngine(engine string)

	// GetDefaultEngine returns the current default engine.
	GetDefaultEngine() string

	// Engine management

	// RegisterEngine registers a new engine.
	// Returns an error if an engine with that name is already registered.
	RegisterEngine(engine Engine) error

	// ListEngines returns the names of all registered engines.
	ListEngines() []string

	// GetEngine retrieves an engine by name.
	// Returns ErrEngineNotFound if the engine doesn't exist.
	GetEngine(name string) (Engine, error)
}
