// Package storage provides interfaces and types for the storage system.
package storage

import (
	"time"
)

// Storage event topics
const (
	// TopicStoreCreated is triggered when a store is created.
	TopicStoreCreated = "store.created"

	// TopicStoreDeleted is triggered when a store is deleted.
	TopicStoreDeleted = "store.deleted"

	// TopicStoreClosed is triggered when a store is closed.
	TopicStoreClosed = "store.closed"

	// TopicVersionSaved is triggered when a version is saved.
	TopicVersionSaved = "store.version.saved"

	// TopicVersionLoaded is triggered when a version is loaded.
	TopicVersionLoaded = "store.version.loaded"

	// TopicTransactionBegin is triggered when a transaction begins.
	TopicTransactionBegin = "store.transaction.begin"

	// TopicTransactionCommit is triggered when a transaction is committed.
	TopicTransactionCommit = "store.transaction.commit"

	// TopicTransactionRollback is triggered when a transaction is rolled back.
	TopicTransactionRollback = "store.transaction.rollback"
)

// CreateStoreEventPayload creates a storage event payload with the given store name and additional data.
func CreateStoreEventPayload(storeName, engineName string, additionalData map[string]interface{}) map[string]interface{} {
	payload := map[string]interface{}{
		"storeName":  storeName,
		"engineName": engineName,
		"timestamp":  time.Now(),
	}

	// Merge additional data if provided
	if additionalData != nil {
		for k, v := range additionalData {
			payload[k] = v
		}
	}

	return payload
}
