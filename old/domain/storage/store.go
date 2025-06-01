// Package storage provides interfaces and types for the storage system.
package storage

// Store defines the core storage operations that all backends must implement.
type Store interface {
	// Basic CRUD operations

	// Get retrieves the value for the given key.
	// Returns ErrKeyNotFound if the key doesn't exist.
	Get(key []byte) ([]byte, error)

	// Set stores the value for the given key.
	// Overwrites any existing value.
	Set(key, value []byte) error

	// Delete removes the key-value pair.
	// It's idempotent (no error if key doesn't exist).
	Delete(key []byte) error

	// Has checks if a key exists.
	// More efficient than Get for existence checks.
	Has(key []byte) (bool, error)

	// Iteration over all key-value pairs

	// Iterate calls fn for each key-value pair, stops if fn returns false.
	// The key and value byte slices must not be modified by fn.
	Iterate(fn func(key, value []byte) bool) error

	// Resource cleanup

	// Close releases resources and makes the store unusable.
	Close() error

	// Store metadata

	// Name returns the store identifier.
	Name() string

	// Path returns the storage path/location.
	Path() string
}
