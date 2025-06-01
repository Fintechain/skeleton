// Package memory provides an in-memory implementation of the storage engine.
package memory

import (
	"github.com/fintechain/skeleton/internal/domain/storage"
)

// Transaction implements storage.Transaction for in-memory storage.
type Transaction struct {
	store    *Store
	changes  map[string][]byte // Pending changes
	deletes  map[string]bool   // Keys to delete
	active   bool
	readOnly bool
}

// BeginTx starts a new transaction.
// Explicitly returns the domain interface for dependency inversion.
func (s *Store) BeginTx() (storage.Transaction, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	if s.closed {
		return nil, storage.ErrStoreClosed
	}

	return &Transaction{
		store:    s,
		changes:  make(map[string][]byte),
		deletes:  make(map[string]bool),
		active:   true,
		readOnly: false,
	}, nil
}

// SupportsTransactions returns true as memory store supports transactions.
func (s *Store) SupportsTransactions() bool {
	return true
}

// Get retrieves the value for the given key within the transaction.
func (tx *Transaction) Get(key []byte) ([]byte, error) {
	if !tx.active {
		return nil, storage.ErrTxNotActive
	}

	keyStr := string(key)

	// Check if the key is in the transaction's changes
	if value, exists := tx.changes[keyStr]; exists {
		// Return a copy to prevent external modification
		result := make([]byte, len(value))
		copy(result, value)
		return result, nil
	}

	// Check if the key is scheduled for deletion
	if tx.deletes[keyStr] {
		return nil, storage.ErrKeyNotFound
	}

	// Otherwise get from the store
	return tx.store.Get(key)
}

// Set stores the value for the given key within the transaction.
func (tx *Transaction) Set(key, value []byte) error {
	if !tx.active {
		return storage.ErrTxNotActive
	}

	if tx.readOnly {
		return storage.ErrTxReadOnly
	}

	keyStr := string(key)

	// Store a copy to prevent external modification
	valueCopy := make([]byte, len(value))
	copy(valueCopy, value)
	tx.changes[keyStr] = valueCopy

	// Remove from deletes if it was there
	delete(tx.deletes, keyStr)

	return nil
}

// Delete removes the key-value pair within the transaction.
func (tx *Transaction) Delete(key []byte) error {
	if !tx.active {
		return storage.ErrTxNotActive
	}

	if tx.readOnly {
		return storage.ErrTxReadOnly
	}

	keyStr := string(key)

	// Mark for deletion
	tx.deletes[keyStr] = true

	// Remove from changes if it was there
	delete(tx.changes, keyStr)

	return nil
}

// Has checks if a key exists within the transaction.
func (tx *Transaction) Has(key []byte) (bool, error) {
	if !tx.active {
		return false, storage.ErrTxNotActive
	}

	keyStr := string(key)

	// Check if the key is in the transaction's changes
	if _, exists := tx.changes[keyStr]; exists {
		return true, nil
	}

	// Check if the key is scheduled for deletion
	if tx.deletes[keyStr] {
		return false, nil
	}

	// Otherwise check the store
	return tx.store.Has(key)
}

// Iterate calls fn for each key-value pair within the transaction.
func (tx *Transaction) Iterate(fn func(key, value []byte) bool) error {
	if !tx.active {
		return storage.ErrTxNotActive
	}

	// First collect keys that are not deleted
	storeKeys := make(map[string]bool)

	// Get keys from the store
	err := tx.store.Iterate(func(key, value []byte) bool {
		keyStr := string(key)
		// Skip if scheduled for deletion
		if !tx.deletes[keyStr] {
			storeKeys[keyStr] = true
		}
		return true
	})

	if err != nil {
		return err
	}

	// Add keys from changes
	for keyStr := range tx.changes {
		storeKeys[keyStr] = true
	}

	// Now iterate over all keys
	for keyStr := range storeKeys {
		key := []byte(keyStr)
		value, err := tx.Get(key)
		if err != nil {
			continue // Skip if error (should not happen)
		}

		if !fn(key, value) {
			break
		}
	}

	return nil
}

// Commit makes all changes permanent.
func (tx *Transaction) Commit() error {
	if !tx.active {
		return storage.ErrTxNotActive
	}

	tx.store.mutex.Lock()
	defer tx.store.mutex.Unlock()

	if tx.store.closed {
		return storage.ErrStoreClosed
	}

	// Apply all changes atomically
	for keyStr, value := range tx.changes {
		tx.store.data[keyStr] = value

		// Update total size
		if tx.store.options.MaxSize > 0 {
			if existing, exists := tx.store.data[keyStr]; exists {
				tx.store.totalSize -= int64(len(existing))
			}
			tx.store.totalSize += int64(len(value))
		}
	}

	// Apply all deletes
	for keyStr := range tx.deletes {
		if existing, exists := tx.store.data[keyStr]; exists {
			tx.store.totalSize -= int64(len(existing))
		}
		delete(tx.store.data, keyStr)
	}

	tx.active = false
	return nil
}

// Rollback discards all changes.
func (tx *Transaction) Rollback() error {
	if !tx.active {
		return storage.ErrTxNotActive
	}

	tx.active = false
	tx.changes = nil
	tx.deletes = nil
	return nil
}

// IsActive returns true if transaction is still active.
func (tx *Transaction) IsActive() bool {
	return tx.active
}

// Close implements the storage.Store Close method but does nothing for transactions.
func (tx *Transaction) Close() error {
	return tx.Rollback()
}

// Name returns the store identifier.
func (tx *Transaction) Name() string {
	return tx.store.Name()
}

// Path returns the storage path/location.
func (tx *Transaction) Path() string {
	return tx.store.Path()
}
