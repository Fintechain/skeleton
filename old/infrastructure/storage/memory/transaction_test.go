package memory

import (
	"testing"

	"github.com/fintechain/skeleton/internal/domain/storage"
	"github.com/fintechain/skeleton/internal/infrastructure/storage/memory/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestStoreForTx() *Store {
	logger := mocks.NewMockLogger()
	options := Options{
		MaxSize:     -1,
		MaxVersions: 10,
	}
	return NewStore("test-store", "/test/path", options, logger)
}

func TestTransactionBasicOperations(t *testing.T) {
	store := createTestStoreForTx()
	defer store.Close()

	// Populate store with initial data
	initialData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for k, v := range initialData {
		err := store.Set([]byte(k), []byte(v))
		require.NoError(t, err)
	}

	// Start a transaction
	tx, err := store.BeginTx()
	require.NoError(t, err)
	require.NotNil(t, tx)
	require.True(t, tx.IsActive())

	// Test Get within transaction (unchanged data)
	val, err := tx.Get([]byte("key1"))
	require.NoError(t, err)
	assert.Equal(t, []byte("value1"), val)

	// Test Set within transaction
	err = tx.Set([]byte("key1"), []byte("new-value1"))
	require.NoError(t, err)

	// Test Get reflects transaction changes
	val, err = tx.Get([]byte("key1"))
	require.NoError(t, err)
	assert.Equal(t, []byte("new-value1"), val)

	// Test Store still has old value
	val, err = store.Get([]byte("key1"))
	require.NoError(t, err)
	assert.Equal(t, []byte("value1"), val)

	// Test Delete within transaction
	err = tx.Delete([]byte("key2"))
	require.NoError(t, err)

	// Test Has reflects delete
	exists, err := tx.Has([]byte("key2"))
	require.NoError(t, err)
	assert.False(t, exists)

	// Test Store still has the key
	exists, err = store.Has([]byte("key2"))
	require.NoError(t, err)
	assert.True(t, exists)

	// Test adding a new key
	err = tx.Set([]byte("key4"), []byte("value4"))
	require.NoError(t, err)

	exists, err = tx.Has([]byte("key4"))
	require.NoError(t, err)
	assert.True(t, exists)

	exists, err = store.Has([]byte("key4"))
	require.NoError(t, err)
	assert.False(t, exists)

	// Verify Name and Path match the store
	assert.Equal(t, store.Name(), tx.Name())
	assert.Equal(t, store.Path(), tx.Path())
}

func TestTransactionCommit(t *testing.T) {
	store := createTestStoreForTx()
	defer store.Close()

	// Populate store with initial data
	initialData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for k, v := range initialData {
		err := store.Set([]byte(k), []byte(v))
		require.NoError(t, err)
	}

	// Start a transaction and make changes
	tx, err := store.BeginTx()
	require.NoError(t, err)

	transactionOps := []struct {
		operation string
		key       string
		value     string
	}{
		{"set", "key1", "updated-value1"},
		{"delete", "key2", ""},
		{"set", "key4", "new-value4"},
	}

	for _, op := range transactionOps {
		switch op.operation {
		case "set":
			err = tx.Set([]byte(op.key), []byte(op.value))
			require.NoError(t, err)
		case "delete":
			err = tx.Delete([]byte(op.key))
			require.NoError(t, err)
		}
	}

	// Commit the transaction
	err = tx.Commit()
	require.NoError(t, err)
	assert.False(t, tx.IsActive())

	// Verify all changes are now visible in the store
	expectedAfterCommit := map[string]struct {
		exists bool
		value  string
	}{
		"key1": {true, "updated-value1"},
		"key2": {false, ""},
		"key3": {true, "value3"},
		"key4": {true, "new-value4"},
	}

	for k, expected := range expectedAfterCommit {
		exists, err := store.Has([]byte(k))
		require.NoError(t, err)
		assert.Equal(t, expected.exists, exists, "Key existence mismatch for %s", k)

		if expected.exists {
			val, err := store.Get([]byte(k))
			require.NoError(t, err)
			assert.Equal(t, []byte(expected.value), val, "Value mismatch for %s", k)
		}
	}

	// Test operations on committed transaction should fail
	_, err = tx.Get([]byte("key1"))
	assert.Error(t, err)
	assert.True(t, storage.IsTxNotActive(err))

	err = tx.Set([]byte("key5"), []byte("value5"))
	assert.Error(t, err)
	assert.True(t, storage.IsTxNotActive(err))

	err = tx.Delete([]byte("key3"))
	assert.Error(t, err)
	assert.True(t, storage.IsTxNotActive(err))

	_, err = tx.Has([]byte("key1"))
	assert.Error(t, err)
	assert.True(t, storage.IsTxNotActive(err))

	err = tx.Commit()
	assert.Error(t, err)
	assert.True(t, storage.IsTxNotActive(err))

	err = tx.Rollback()
	assert.Error(t, err)
	assert.True(t, storage.IsTxNotActive(err))
}

func TestTransactionRollback(t *testing.T) {
	store := createTestStoreForTx()
	defer store.Close()

	// Populate store with initial data
	initialData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for k, v := range initialData {
		err := store.Set([]byte(k), []byte(v))
		require.NoError(t, err)
	}

	// Start a transaction and make changes
	tx, err := store.BeginTx()
	require.NoError(t, err)

	// Make various changes
	err = tx.Set([]byte("key1"), []byte("changed-value1"))
	require.NoError(t, err)
	err = tx.Delete([]byte("key2"))
	require.NoError(t, err)
	err = tx.Set([]byte("key4"), []byte("new-value4"))
	require.NoError(t, err)

	// Rollback the transaction
	err = tx.Rollback()
	require.NoError(t, err)
	assert.False(t, tx.IsActive())

	// Verify store is unchanged
	for k, v := range initialData {
		exists, err := store.Has([]byte(k))
		require.NoError(t, err)
		assert.True(t, exists, "Key should still exist: %s", k)

		val, err := store.Get([]byte(k))
		require.NoError(t, err)
		assert.Equal(t, []byte(v), val, "Value should be unchanged for %s", k)
	}

	// Verify new key was not added
	exists, err := store.Has([]byte("key4"))
	require.NoError(t, err)
	assert.False(t, exists, "Key should not exist after rollback")

	// Test operations on rolled back transaction should fail
	_, err = tx.Get([]byte("key1"))
	assert.Error(t, err)
	assert.True(t, storage.IsTxNotActive(err))
}

func TestTransactionIterate(t *testing.T) {
	store := createTestStoreForTx()
	defer store.Close()

	// Populate store with initial data
	initialData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for k, v := range initialData {
		err := store.Set([]byte(k), []byte(v))
		require.NoError(t, err)
	}

	// Start a transaction and make changes
	tx, err := store.BeginTx()
	require.NoError(t, err)

	// Make various changes
	err = tx.Set([]byte("key1"), []byte("changed-value1"))
	require.NoError(t, err)
	err = tx.Delete([]byte("key2"))
	require.NoError(t, err)
	err = tx.Set([]byte("key4"), []byte("new-value4"))
	require.NoError(t, err)

	// Test iterate function
	expected := map[string]string{
		"key1": "changed-value1",
		"key3": "value3",
		"key4": "new-value4",
	}

	collected := make(map[string]string)
	err = tx.Iterate(func(key, value []byte) bool {
		collected[string(key)] = string(value)
		return true
	})
	require.NoError(t, err)

	// Verify collected matches expected
	assert.Equal(t, len(expected), len(collected))
	for k, v := range expected {
		assert.Equal(t, v, collected[k])
	}

	// Test early termination
	count := 0
	err = tx.Iterate(func(key, value []byte) bool {
		count++
		return count < 2 // Stop after processing 1 item
	})
	require.NoError(t, err)
	assert.Equal(t, 2, count) // Should process 2 items (including the one that returns false)
}

func TestTransactionClose(t *testing.T) {
	store := createTestStoreForTx()
	defer store.Close()

	// Start a transaction
	tx, err := store.BeginTx()
	require.NoError(t, err)

	// Make a change
	err = tx.Set([]byte("key1"), []byte("value1"))
	require.NoError(t, err)

	// Close should be equivalent to rollback
	err = tx.Close()
	require.NoError(t, err)
	assert.False(t, tx.IsActive())

	// Verify store is unchanged
	exists, err := store.Has([]byte("key1"))
	require.NoError(t, err)
	assert.False(t, exists)
}

func TestTransactionWithClosedStore(t *testing.T) {
	store := createTestStoreForTx()

	// Start a transaction
	tx, err := store.BeginTx()
	require.NoError(t, err)

	// Make some changes
	err = tx.Set([]byte("key1"), []byte("value1"))
	require.NoError(t, err)

	// Close the store
	err = store.Close()
	require.NoError(t, err)

	// Try to commit - should fail
	err = tx.Commit()
	assert.Error(t, err)
	assert.True(t, storage.IsStoreClosed(err))
}

func TestReadOnlyTransaction(t *testing.T) {
	// Note: Memory transaction doesn't currently support read-only mode explicitly
	// This test just documents the expected behavior if it was implemented

	store := createTestStoreForTx()
	defer store.Close()

	// Populate store with initial data
	err := store.Set([]byte("key1"), []byte("value1"))
	require.NoError(t, err)

	// Start a transaction
	tx, err := store.BeginTx()
	require.NoError(t, err)

	// Manually set to read-only since the API doesn't expose this
	txInternal := tx.(*Transaction)
	txInternal.readOnly = true

	// Read operations should work
	val, err := tx.Get([]byte("key1"))
	require.NoError(t, err)
	assert.Equal(t, []byte("value1"), val)

	exists, err := tx.Has([]byte("key1"))
	require.NoError(t, err)
	assert.True(t, exists)

	// Write operations should fail
	err = tx.Set([]byte("key2"), []byte("value2"))
	assert.Error(t, err)
	assert.True(t, storage.IsTxReadOnly(err))

	err = tx.Delete([]byte("key1"))
	assert.Error(t, err)
	assert.True(t, storage.IsTxReadOnly(err))

	// Commit should work on read-only transaction
	err = tx.Commit()
	require.NoError(t, err)
}
