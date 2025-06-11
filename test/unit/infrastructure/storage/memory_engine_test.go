package storage

import (
	"testing"

	"github.com/fintechain/skeleton/internal/domain/storage"
	memoryStorage "github.com/fintechain/skeleton/internal/infrastructure/storage/memory"
	"github.com/stretchr/testify/assert"
)

func TestNewMemoryEngine(t *testing.T) {
	engine := memoryStorage.NewEngine()
	assert.NotNil(t, engine)

	// Verify interface compliance
	var _ storage.Engine = engine

	// Test basic properties
	assert.Equal(t, "memory", engine.Name())
}

func TestMemoryEngineCapabilities(t *testing.T) {
	engine := memoryStorage.NewEngine()

	capabilities := engine.Capabilities()
	assert.NotNil(t, capabilities)

	// Memory engine capabilities
	assert.False(t, capabilities.Transactions)
	assert.False(t, capabilities.Versioning)
	assert.False(t, capabilities.RangeQueries)
	assert.False(t, capabilities.Persistence)
	assert.False(t, capabilities.Compression)
}

func TestMemoryEngineCreateStore(t *testing.T) {
	engine := memoryStorage.NewEngine()

	// Test create store
	store, err := engine.Create("test-store", "/tmp/test", nil)
	assert.NoError(t, err)
	assert.NotNil(t, store)

	// Verify interface compliance
	var _ storage.Store = store

	// Test store properties
	assert.Equal(t, "test-store", store.Name())
	assert.Equal(t, "/tmp/test", store.Path())
}

func TestMemoryEngineCreateStoreWithConfig(t *testing.T) {
	engine := memoryStorage.NewEngine()

	// Test create store with config (should ignore config for memory engine)
	config := map[string]interface{}{
		"buffer_size": 1024,
		"cache_size":  512,
	}

	store, err := engine.Create("config-store", "/tmp/config", config)
	assert.NoError(t, err)
	assert.NotNil(t, store)
	assert.Equal(t, "config-store", store.Name())
	assert.Equal(t, "/tmp/config", store.Path())
}

func TestMemoryEngineCreateMultipleStores(t *testing.T) {
	engine := memoryStorage.NewEngine()

	// Create multiple stores
	store1, err := engine.Create("store1", "/tmp/store1", nil)
	assert.NoError(t, err)
	assert.NotNil(t, store1)

	store2, err := engine.Create("store2", "/tmp/store2", nil)
	assert.NoError(t, err)
	assert.NotNil(t, store2)

	// Stores should be independent
	assert.NotEqual(t, store1, store2)
	assert.Equal(t, "store1", store1.Name())
	assert.Equal(t, "store2", store2.Name())

	// Test that stores are independent
	err = store1.Set([]byte("key"), []byte("value1"))
	assert.NoError(t, err)

	err = store2.Set([]byte("key"), []byte("value2"))
	assert.NoError(t, err)

	// Values should be different
	value1, err := store1.Get([]byte("key"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("value1"), value1)

	value2, err := store2.Get([]byte("key"))
	assert.NoError(t, err)
	assert.Equal(t, []byte("value2"), value2)
}
