// Package storage provides storage interfaces and types.
package storage

import (
	"github.com/fintechain/skeleton/internal/domain/storage"
)

// Re-export storage engine interfaces
type Engine = storage.Engine
type RangeQueryable = storage.RangeQueryable

// Re-export storage engine types
type Capabilities = storage.Capabilities
