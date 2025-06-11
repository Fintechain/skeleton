// Package storage provides interfaces and types for the storage system.
package storage

// Store defines the core storage operations that all backends must implement.
type Store interface {
	// Get retrieves the value associated with the given key.
	Get(key []byte) ([]byte, error)

	// Set stores a value for the given key.
	Set(key, value []byte) error

	// Delete removes the key-value pair for the given key.
	Delete(key []byte) error

	// Has checks whether a key exists in the store.
	Has(key []byte) (bool, error)

	// Iterate calls the provided function for each key-value pair in the store.
	Iterate(fn func(key, value []byte) bool) error

	// Close releases all resources associated with the store.
	Close() error

	// Name returns the name of this store instance.
	Name() string

	// Path returns the storage path or location for this store.
	Path() string
}

// Engine defines the factory interface for storage backend implementations.
type Engine interface {
	// Name returns the unique identifier for this storage engine.
	Name() string

	// Create creates a new store instance with the specified configuration.
	Create(name, path string, config Config) (Store, error)

	// Open opens an existing store at the specified path.
	Open(name, path string) (Store, error)

	// Capabilities returns the features supported by this storage engine.
	Capabilities() Capabilities
}
