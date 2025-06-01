// Package storage provides storage interfaces and types.
package storage

import (
	"github.com/fintechain/skeleton/internal/domain/storage"
	storageImpl "github.com/fintechain/skeleton/internal/infrastructure/storage"
)

// Re-export multistore interface
type MultiStore = storage.MultiStore

// NewMultiStore creates a new MultiStore instance.
// This factory function provides access to the concrete multistore implementation.
func NewMultiStore() MultiStore {
	return storageImpl.NewMultiStore()
}
