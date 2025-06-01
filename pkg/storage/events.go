// Package storage provides storage interfaces and types.
package storage

import (
	"github.com/fintechain/skeleton/internal/domain/storage"
)

// Re-export storage event topics
const (
	TopicStoreCreated        = storage.TopicStoreCreated
	TopicStoreDeleted        = storage.TopicStoreDeleted
	TopicStoreClosed         = storage.TopicStoreClosed
	TopicVersionSaved        = storage.TopicVersionSaved
	TopicVersionLoaded       = storage.TopicVersionLoaded
	TopicTransactionBegin    = storage.TopicTransactionBegin
	TopicTransactionCommit   = storage.TopicTransactionCommit
	TopicTransactionRollback = storage.TopicTransactionRollback
)

// Re-export event utility functions
var CreateStoreEventPayload = storage.CreateStoreEventPayload
