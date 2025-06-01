// Package storage provides storage interfaces and types.
package storage

import (
	"github.com/fintechain/skeleton/internal/domain/storage"
)

// Re-export storage error constants
const (
	ErrKeyNotFound     = storage.ErrKeyNotFound
	ErrStoreNotFound   = storage.ErrStoreNotFound
	ErrStoreClosed     = storage.ErrStoreClosed
	ErrStoreExists     = storage.ErrStoreExists
	ErrEngineNotFound  = storage.ErrEngineNotFound
	ErrTxNotActive     = storage.ErrTxNotActive
	ErrTxReadOnly      = storage.ErrTxReadOnly
	ErrVersionNotFound = storage.ErrVersionNotFound
	ErrInvalidConfig   = storage.ErrInvalidConfig
)

// Re-export error utility functions
var (
	WrapError         = storage.WrapError
	IsKeyNotFound     = storage.IsKeyNotFound
	IsStoreNotFound   = storage.IsStoreNotFound
	IsStoreClosed     = storage.IsStoreClosed
	IsStoreExists     = storage.IsStoreExists
	IsEngineNotFound  = storage.IsEngineNotFound
	IsTxNotActive     = storage.IsTxNotActive
	IsTxReadOnly      = storage.IsTxReadOnly
	IsVersionNotFound = storage.IsVersionNotFound
	IsInvalidConfig   = storage.IsInvalidConfig
)
