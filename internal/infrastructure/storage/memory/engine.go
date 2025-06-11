// Package memory provides an in-memory storage engine implementation.
package memory

import (
	"errors"

	"github.com/fintechain/skeleton/internal/domain/storage"
)

// Engine implements the storage.Engine interface for in-memory storage.
type Engine struct {
	name string
}

// NewEngine creates a new memory storage engine.
func NewEngine() *Engine {
	return &Engine{
		name: "memory",
	}
}

// Name returns the unique identifier for this storage engine.
func (e *Engine) Name() string {
	return e.name
}

// Create creates a new in-memory store instance with the specified configuration.
func (e *Engine) Create(name, path string, config storage.Config) (storage.Store, error) {
	if name == "" {
		return nil, errors.New(storage.ErrInvalidConfig)
	}

	return NewStore(name, path), nil
}

// Open opens an existing in-memory store at the specified path.
// For memory storage, this is equivalent to Create since there's no persistence.
func (e *Engine) Open(name, path string) (storage.Store, error) {
	return e.Create(name, path, nil)
}

// Capabilities returns the features supported by this storage engine.
func (e *Engine) Capabilities() storage.Capabilities {
	return storage.Capabilities{
		Transactions: false, // Memory store doesn't support transactions
		Versioning:   false, // Memory store doesn't support versioning
		RangeQueries: false, // Simplified - no range queries
		Persistence:  false, // Memory-only storage
		Compression:  false, // No compression needed for memory
	}
}
