// Package storage provides interfaces and types for the storage system.
package storage

// Standard storage error codes
const (
	// ErrKeyNotFound is returned when a key doesn't exist in a store
	ErrKeyNotFound = "storage.key_not_found"

	// ErrStoreNotFound is returned when a store doesn't exist
	ErrStoreNotFound = "storage.store_not_found"

	// ErrStoreClosed is returned when operations are performed on a closed store
	ErrStoreClosed = "storage.store_closed"

	// ErrStoreExists is returned when creating a store that already exists
	ErrStoreExists = "storage.store_exists"

	// ErrEngineNotFound is returned when an engine doesn't exist
	ErrEngineNotFound = "storage.engine_not_found"

	// ErrEngineExists is returned when registering an engine that already exists
	ErrEngineExists = "storage.engine_exists"

	// ErrTxNotActive is returned when operations are performed on a non-active transaction
	ErrTxNotActive = "storage.transaction_not_active"

	// ErrTxReadOnly is returned when write operations are performed on a read-only transaction
	ErrTxReadOnly = "storage.transaction_read_only"

	// ErrTxAlreadyActive is returned when starting a transaction that is already active
	ErrTxAlreadyActive = "storage.transaction_already_active"

	// ErrVersionNotFound is returned when a version doesn't exist
	ErrVersionNotFound = "storage.version_not_found"

	// ErrInvalidVersion is returned when an invalid version is provided
	ErrInvalidVersion = "storage.invalid_version"

	// ErrInvalidConfig is returned when invalid configuration is provided
	ErrInvalidConfig = "storage.invalid_config"

	// ErrStoreCorrupted is returned when a store is corrupted
	ErrStoreCorrupted = "storage.store_corrupted"

	// ErrInsufficientSpace is returned when there is insufficient storage space
	ErrInsufficientSpace = "storage.insufficient_space"

	// ErrOperationNotSupported is returned when an operation is not supported by the engine
	ErrOperationNotSupported = "storage.operation_not_supported"
)
