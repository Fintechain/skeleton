// Package storage provides storage-related interfaces and types for the Fintechain Skeleton framework.
package storage

import (
	"github.com/fintechain/skeleton/internal/domain/component"
)

// MultiStore manages multiple named stores.
// This is the core multi-store interface without lifecycle management.
type MultiStore interface {
	// Store management
	GetStore(name string) (Store, error)
	CreateStore(name, engine string, config Config) error
	DeleteStore(name string) error
	ListStores() []string

	// Bulk operations
	CloseAll() error

	// Engine management
	RegisterEngine(engine Engine) error
}

// MultiStoreService provides multi-store persistence functionality as an infrastructure service.
// It combines the core multi-store functionality with service lifecycle management.
type MultiStoreService interface {
	component.Service
	MultiStore
}
