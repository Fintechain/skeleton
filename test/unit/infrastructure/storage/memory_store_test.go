package storage

import (
	"testing"

	"github.com/fintechain/skeleton/internal/domain/storage"
	memoryStorage "github.com/fintechain/skeleton/internal/infrastructure/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestNewMemoryStore(t *testing.T) {
	store := memoryStorage.NewStore("test-store", "/tmp/test")
	assert.NotNil(t, store)

	// Verify interface compliance
	var _ storage.Store = store

	// Test basic properties
	assert.Equal(t, "test-store", store.Name())
	assert.Equal(t, "/tmp/test", store.Path())
}

func TestMemoryStoreInitialState(t *testing.T) {
	store := memoryStorage.NewStore("test-store", "/tmp/test")

	// Test initial empty state
	exists, err := store.Has([]byte("non-existent"))
	assert.NoError(t, err)
	assert.False(t, exists)

	// Test get non-existent key
	_, err = store.Get([]byte("non-existent"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "storage.key_not_found")
}

func TestMemoryStorePutAndGet(t *testing.T) {
	store := memoryStorage.NewStore("test-store", "/tmp/test")

	key := []byte("test-key")
	value := []byte("test-value")

	// Test put
	err := store.Set(key, value)
	assert.NoError(t, err)

	// Test get
	retrievedValue, err := store.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value, retrievedValue)

	// Test has
	exists, err := store.Has(key)
	assert.NoError(t, err)
	assert.True(t, exists)
}

func TestMemoryStorePutOverwrite(t *testing.T) {
	store := memoryStorage.NewStore("test-store", "/tmp/test")

	key := []byte("test-key")
	value1 := []byte("value1")
	value2 := []byte("value2")

	// Put first value
	err := store.Set(key, value1)
	assert.NoError(t, err)

	// Put second value (overwrite)
	err = store.Set(key, value2)
	assert.NoError(t, err)

	// Get should return second value
	retrievedValue, err := store.Get(key)
	assert.NoError(t, err)
	assert.Equal(t, value2, retrievedValue)
}

func TestMemoryStoreDelete(t *testing.T) {
	store := memoryStorage.NewStore("test-store", "/tmp/test")

	key := []byte("test-key")
	value := []byte("test-value")

	// Put value
	err := store.Set(key, value)
	assert.NoError(t, err)

	// Verify it exists
	exists, err := store.Has(key)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Delete
	err = store.Delete(key)
	assert.NoError(t, err)

	// Verify it's gone
	exists, err = store.Has(key)
	assert.NoError(t, err)
	assert.False(t, exists)

	// Get should return error
	_, err = store.Get(key)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "storage.key_not_found")
}

func TestMemoryStoreDeleteNonExistent(t *testing.T) {
	store := memoryStorage.NewStore("test-store", "/tmp/test")

	// Delete non-existent key (should return error)
	err := store.Delete([]byte("non-existent"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "storage.key_not_found")
}

func TestMemoryStoreWithDifferentDataTypes(t *testing.T) {
	store := memoryStorage.NewStore("test-store", "/tmp/test")

	tests := []struct {
		name  string
		key   []byte
		value []byte
	}{
		{"string data", []byte("string-key"), []byte("string value")},
		{"binary data", []byte("binary-key"), []byte{0x00, 0x01, 0x02, 0xFF}},
		{"empty value", []byte("empty-key"), []byte{}},
		{"unicode key", []byte("unicode-🔑"), []byte("unicode value")},
		{"large data", []byte("large-key"), make([]byte, 1024)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Put
			err := store.Set(tt.key, tt.value)
			assert.NoError(t, err)

			// Get and verify
			retrievedValue, err := store.Get(tt.key)
			assert.NoError(t, err)
			assert.Equal(t, tt.value, retrievedValue)

			// Has
			exists, err := store.Has(tt.key)
			assert.NoError(t, err)
			assert.True(t, exists)
		})
	}
}

func TestMemoryStoreNilKeyHandling(t *testing.T) {
	store := memoryStorage.NewStore("test-store", "/tmp/test")

	// Test with nil key (should handle gracefully)
	err := store.Set(nil, []byte("value"))
	assert.NoError(t, err)

	// Get with nil key
	value, err := store.Get(nil)
	assert.NoError(t, err)
	assert.Equal(t, []byte("value"), value)

	// Has with nil key
	exists, err := store.Has(nil)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Delete with nil key
	err = store.Delete(nil)
	assert.NoError(t, err)
}

func TestMemoryStoreClose(t *testing.T) {
	store := memoryStorage.NewStore("test-store", "/tmp/test")

	// Put some data
	err := store.Set([]byte("key"), []byte("value"))
	assert.NoError(t, err)

	// Close should not error
	err = store.Close()
	assert.NoError(t, err)

	// Multiple closes should not error
	err = store.Close()
	assert.NoError(t, err)
}
