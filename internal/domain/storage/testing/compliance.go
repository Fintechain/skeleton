// Package testing provides compliance tests for storage implementations.
package testing

import (
	"testing"

	"github.com/ebanfa/skeleton/internal/domain/storage"
	"github.com/stretchr/testify/require"
)

// TestStoreCompliance tests that a store implementation satisfies the Store interface.
func TestStoreCompliance(t *testing.T, createStore func() storage.Store) {
	store := createStore()
	defer store.Close()

	// Test basic operations
	key := []byte("test-key")
	value := []byte("test-value")

	// Test Set/Get
	err := store.Set(key, value)
	require.NoError(t, err)

	retrieved, err := store.Get(key)
	require.NoError(t, err)
	require.Equal(t, value, retrieved)

	// Test Has
	exists, err := store.Has(key)
	require.NoError(t, err)
	require.True(t, exists)

	// Test Delete
	err = store.Delete(key)
	require.NoError(t, err)

	// Test key not found
	_, err = store.Get(key)
	require.True(t, storage.IsKeyNotFound(err))

	exists, err = store.Has(key)
	require.NoError(t, err)
	require.False(t, exists)
}

// TestTransactionalCompliance tests transaction support.
func TestTransactionalCompliance(t *testing.T, store storage.Store) {
	tx, ok := store.(storage.Transactional)
	if !ok {
		t.Skip("Store does not support transactions")
	}

	transaction, err := tx.BeginTx()
	require.NoError(t, err)
	require.True(t, transaction.IsActive())

	// Test transaction operations
	key := []byte("tx-key")
	value := []byte("tx-value")

	err = transaction.Set(key, value)
	require.NoError(t, err)

	// Value should be visible inside transaction
	txValue, err := transaction.Get(key)
	require.NoError(t, err)
	require.Equal(t, value, txValue)

	// Value should not be visible outside transaction
	_, err = store.Get(key)
	require.True(t, storage.IsKeyNotFound(err))

	// Commit and verify
	err = transaction.Commit()
	require.NoError(t, err)
	require.False(t, transaction.IsActive())

	retrieved, err := store.Get(key)
	require.NoError(t, err)
	require.Equal(t, value, retrieved)
}

// TestVersionedCompliance tests versioning support.
func TestVersionedCompliance(t *testing.T, store storage.Store) {
	vs, ok := store.(storage.Versioned)
	if !ok {
		t.Skip("Store does not support versioning")
	}

	// Test versioning operations
	key := []byte("version-key")
	value1 := []byte("version-value-1")
	value2 := []byte("version-value-2")

	// Set initial value
	err := store.Set(key, value1)
	require.NoError(t, err)

	// Save first version
	version1, hash1, err := vs.SaveVersion()
	require.NoError(t, err)
	require.Greater(t, version1, int64(0))
	require.NotEmpty(t, hash1)

	// Update value
	err = store.Set(key, value2)
	require.NoError(t, err)

	// Save second version
	version2, hash2, err := vs.SaveVersion()
	require.NoError(t, err)
	require.Greater(t, version2, version1)
	require.NotEmpty(t, hash2)
	require.NotEqual(t, hash1, hash2)

	// Check current version
	currVersion := vs.CurrentVersion()
	require.Equal(t, version2, currVersion)

	// List versions
	versions := vs.ListVersions()
	require.Contains(t, versions, version1)
	require.Contains(t, versions, version2)

	// Load first version
	err = vs.LoadVersion(version1)
	require.NoError(t, err)

	// Check value is from first version
	retrieved, err := store.Get(key)
	require.NoError(t, err)
	require.Equal(t, value1, retrieved)

	// Load second version
	err = vs.LoadVersion(version2)
	require.NoError(t, err)

	// Check value is from second version
	retrieved, err = store.Get(key)
	require.NoError(t, err)
	require.Equal(t, value2, retrieved)
}
