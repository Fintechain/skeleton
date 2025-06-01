// Package storage provides interfaces and types for the storage system.
package storage

// Versioned interface for stores that support versioning/snapshots.
type Versioned interface {
	// SaveVersion creates a new immutable version of the store.
	// Returns the version number and a hash that uniquely identifies the version.
	// The hash can be used for integrity verification.
	SaveVersion() (version int64, hash []byte, err error)

	// LoadVersion loads a specific version of the store.
	// Returns ErrVersionNotFound if the version doesn't exist.
	LoadVersion(version int64) error

	// ListVersions returns all available versions.
	// Returns an empty slice if no versions are available.
	ListVersions() []int64

	// CurrentVersion returns the current version number.
	// Returns 0 if no versions have been saved.
	CurrentVersion() int64

	// SupportsVersioning returns true if this store supports versioning.
	SupportsVersioning() bool
}
