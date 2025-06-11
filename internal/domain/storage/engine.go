// Package storage provides interfaces and types for the storage system.
package storage

// Capabilities describes the optional features supported by a storage engine.
type Capabilities struct {
	// Transactions indicates whether the engine supports atomic transactions.
	Transactions bool

	// Versioning indicates whether the engine supports data versioning or snapshots.
	Versioning bool

	// RangeQueries indicates whether the engine supports efficient range queries.
	RangeQueries bool

	// Persistence indicates whether the engine persists data to durable storage.
	Persistence bool

	// Compression indicates whether the engine supports data compression.
	Compression bool
}

// RangeQueryable defines the interface for stores that support efficient range queries.
type RangeQueryable interface {
	// IterateRange iterates over keys within the specified range.
	IterateRange(start, end []byte, ascending bool, fn func(key, value []byte) bool) error

	// SupportsRangeQueries returns whether this store implementation supports range queries.
	SupportsRangeQueries() bool
}
