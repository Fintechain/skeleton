// Package storage provides interfaces and types for the storage system.
package storage

import (
	"errors"
	"fmt"
)

// Standard storage errors
var (
	// ErrKeyNotFound is returned when a key doesn't exist in a store
	ErrKeyNotFound = errors.New("key not found")

	// ErrStoreNotFound is returned when a store doesn't exist
	ErrStoreNotFound = errors.New("store not found")

	// ErrStoreClosed is returned when operations are performed on a closed store
	ErrStoreClosed = errors.New("store is closed")

	// ErrStoreExists is returned when creating a store that already exists
	ErrStoreExists = errors.New("store already exists")

	// ErrEngineNotFound is returned when an engine doesn't exist
	ErrEngineNotFound = errors.New("engine not found")

	// ErrTxNotActive is returned when operations are performed on a non-active transaction
	ErrTxNotActive = errors.New("transaction not active")

	// ErrTxReadOnly is returned when write operations are performed on a read-only transaction
	ErrTxReadOnly = errors.New("transaction is read-only")

	// ErrVersionNotFound is returned when a version doesn't exist
	ErrVersionNotFound = errors.New("version not found")

	// ErrInvalidConfig is returned when invalid configuration is provided
	ErrInvalidConfig = errors.New("invalid configuration")
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
	return errors.Is(err, ErrKeyNotFound)
}

// IsStoreNotFound returns true if the error indicates a store was not found.
func IsStoreNotFound(err error) bool {
	return errors.Is(err, ErrStoreNotFound)
}

// IsStoreClosed returns true if the error indicates a store is closed.
func IsStoreClosed(err error) bool {
	return errors.Is(err, ErrStoreClosed)
}

// IsStoreExists returns true if the error indicates a store already exists.
func IsStoreExists(err error) bool {
	return errors.Is(err, ErrStoreExists)
}

// IsEngineNotFound returns true if the error indicates an engine was not found.
func IsEngineNotFound(err error) bool {
	return errors.Is(err, ErrEngineNotFound)
}

// IsTxNotActive returns true if the error indicates a transaction is not active.
func IsTxNotActive(err error) bool {
	return errors.Is(err, ErrTxNotActive)
}

// IsTxReadOnly returns true if the error indicates a transaction is read-only.
func IsTxReadOnly(err error) bool {
	return errors.Is(err, ErrTxReadOnly)
}

// IsVersionNotFound returns true if the error indicates a version was not found.
func IsVersionNotFound(err error) bool {
	return errors.Is(err, ErrVersionNotFound)
}

// IsInvalidConfig returns true if the error indicates invalid configuration.
func IsInvalidConfig(err error) bool {
	return errors.Is(err, ErrInvalidConfig)
}
