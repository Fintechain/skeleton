// Package storage provides interfaces and types for the storage system.
package storage

import (
	"fmt"
)

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

	// ErrTxNotActive is returned when operations are performed on a non-active transaction
	ErrTxNotActive = "storage.transaction_not_active"

	// ErrTxReadOnly is returned when write operations are performed on a read-only transaction
	ErrTxReadOnly = "storage.transaction_read_only"

	// ErrVersionNotFound is returned when a version doesn't exist
	ErrVersionNotFound = "storage.version_not_found"

	// ErrInvalidConfig is returned when invalid configuration is provided
	ErrInvalidConfig = "storage.invalid_config"
)

// WrapError wraps an error with additional context.
// If err is nil, it returns nil.
func WrapError(err error, context string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}

// IsKeyNotFound returns true if the error indicates a key was not found.
func IsKeyNotFound(err error) bool {
	return err != nil && err.Error() == ErrKeyNotFound
}

// IsStoreNotFound returns true if the error indicates a store was not found.
func IsStoreNotFound(err error) bool {
	return err != nil && err.Error() == ErrStoreNotFound
}

// IsStoreClosed returns true if the error indicates a store is closed.
func IsStoreClosed(err error) bool {
	return err != nil && err.Error() == ErrStoreClosed
}

// IsStoreExists returns true if the error indicates a store already exists.
func IsStoreExists(err error) bool {
	return err != nil && err.Error() == ErrStoreExists
}

// IsEngineNotFound returns true if the error indicates an engine was not found.
func IsEngineNotFound(err error) bool {
	return err != nil && err.Error() == ErrEngineNotFound
}

// IsTxNotActive returns true if the error indicates a transaction is not active.
func IsTxNotActive(err error) bool {
	return err != nil && err.Error() == ErrTxNotActive
}

// IsTxReadOnly returns true if the error indicates a transaction is read-only.
func IsTxReadOnly(err error) bool {
	return err != nil && err.Error() == ErrTxReadOnly
}

// IsVersionNotFound returns true if the error indicates a version was not found.
func IsVersionNotFound(err error) bool {
	return err != nil && err.Error() == ErrVersionNotFound
}

// IsInvalidConfig returns true if the error indicates invalid configuration.
func IsInvalidConfig(err error) bool {
	return err != nil && err.Error() == ErrInvalidConfig
}
