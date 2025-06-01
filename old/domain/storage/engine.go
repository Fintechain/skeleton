// Package storage provides interfaces and types for the storage system.
package storage

// Engine interface for storage backend implementations.
// Each concrete storage implementation (memory, file, etc.) must provide
// an Engine implementation that can create and open stores.
type Engine interface {
	// Name returns the engine identifier (e.g., "memory", "leveldb", "file").
	// This is used for engine registration and lookup.
	Name() string

	// Create creates a new store instance with the given name and path.
	// The config parameter contains engine-specific configuration options.
	Create(name, path string, config Config) (Store, error)

	// Open opens an existing store with the given name and path.
	// Returns ErrStoreNotFound if the store doesn't exist.
	Open(name, path string) (Store, error)

	// Capabilities returns what features this engine supports.
	// This allows clients to check for optional features before using them.
	Capabilities() Capabilities
}

// Capabilities describes what features an engine supports.
// This allows clients to check for optional capabilities before using them.
type Capabilities struct {
	// Transactions indicates if the engine supports atomic transactions.
	Transactions bool

	// Versioning indicates if the engine supports versioning/snapshots.
	Versioning bool

	// RangeQueries indicates if the engine supports efficient range queries.
	RangeQueries bool

	// Persistence indicates if the engine persists data to disk.
	Persistence bool

	// Compression indicates if the engine supports data compression.
	Compression bool
}

// RangeQueryable interface for stores that support range queries.
type RangeQueryable interface {
	// IterateRange iterates over keys in the specified range.
	// If ascending is true, iteration is from start to end.
	// If ascending is false, iteration is from end to start.
	// The fn callback works the same as in Store.Iterate.
	IterateRange(start, end []byte, ascending bool, fn func(key, value []byte) bool) error

	// SupportsRangeQueries returns true if this store supports range queries.
	SupportsRangeQueries() bool
}
