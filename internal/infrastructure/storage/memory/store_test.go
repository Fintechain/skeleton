package memory

import (
	"testing"

	"github.com/ebanfa/skeleton/internal/domain/storage"
	storagetesting "github.com/ebanfa/skeleton/internal/domain/storage/testing"
	"github.com/ebanfa/skeleton/internal/infrastructure/storage/memory/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to create a test store with minimal dependencies
func createTestStore() storage.Store {
	logger := mocks.NewMockLogger()
	options := Options{
		MaxSize:     -1,
		MaxVersions: 10,
	}
	return NewStore("test-store", "/test/path", options, logger)
}

// TestStoreImplementsInterfaces tests that the memory store implements the expected interfaces
func TestStoreImplementsInterfaces(t *testing.T) {
	store := createTestStore()
	defer store.Close()

	// Test Store interface
	_, ok := store.(storage.Store)
	assert.True(t, ok, "Store should implement storage.Store")

	// Test Transactional interface
	_, ok = store.(storage.Transactional)
	assert.True(t, ok, "Store should implement storage.Transactional")

	// Test Versioned interface
	_, ok = store.(storage.Versioned)
	assert.True(t, ok, "Store should implement storage.Versioned")

	// Test RangeQueryable interface
	_, ok = store.(storage.RangeQueryable)
	assert.True(t, ok, "Store should implement storage.RangeQueryable")
}

// TestEngineImplementsInterface tests that the memory engine implements the Engine interface
func TestEngineImplementsInterface(t *testing.T) {
	logger := mocks.NewMockLogger()
	engine := NewEngine(logger)
	_, ok := engine.(storage.Engine)
	assert.True(t, ok, "Engine should implement storage.Engine")
}

// TestBasicStoreOperations tests basic store operations
func TestBasicStoreOperations(t *testing.T) {
	storagetesting.TestStoreCompliance(t, createTestStore)
}

// TestTransactionSupport tests transaction support
func TestTransactionSupport(t *testing.T) {
	store := createTestStore()
	defer store.Close()
	storagetesting.TestTransactionalCompliance(t, store)
}

// TestVersioningSupport tests versioning support
func TestVersioningSupport(t *testing.T) {
	store := createTestStore()
	defer store.Close()
	storagetesting.TestVersionedCompliance(t, store)
}

// TestRangeQueries tests range query functionality
func TestRangeQueries(t *testing.T) {
	store := createTestStore()
	defer store.Close()
	storagetesting.TestRangeQueryable(t, store)
}

// TestStoreSizeLimit tests the store size limit functionality
func TestStoreSizeLimit(t *testing.T) {
	logger := mocks.NewMockLogger()

	tests := []struct {
		name       string
		maxSize    int64
		keyValues  map[string]string
		shouldFail map[string]bool
	}{
		{
			name:    "No size limit",
			maxSize: -1,
			keyValues: map[string]string{
				"key1": "small value",
				"key2": string(make([]byte, 200)),
			},
			shouldFail: map[string]bool{
				"key1": false,
				"key2": false,
			},
		},
		{
			name:    "Small size limit",
			maxSize: 100,
			keyValues: map[string]string{
				"key1": "small value",
				"key2": string(make([]byte, 200)),
			},
			shouldFail: map[string]bool{
				"key1": false,
				"key2": true,
			},
		},
		{
			name:    "Exact size limit",
			maxSize: 20,
			keyValues: map[string]string{
				"key1": "12345678901234567890",  // 20 bytes
				"key2": "123456789012345678901", // 21 bytes
			},
			shouldFail: map[string]bool{
				"key1": false,
				"key2": true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a store with configured size limit
			options := Options{
				MaxSize: tc.maxSize,
			}
			store := NewStore("test-store", "/test/path", options, logger)
			defer store.Close()

			for key, value := range tc.keyValues {
				err := store.Set([]byte(key), []byte(value))

				if tc.shouldFail[key] {
					assert.Error(t, err, "Setting %s should fail with size limit %d", key, tc.maxSize)
					assert.True(t, storage.IsInvalidConfig(err), "Error should be InvalidConfig")

					// Verify key was not set
					has, err := store.Has([]byte(key))
					assert.NoError(t, err)
					assert.False(t, has, "Key %s should not exist after failed Set", key)
				} else {
					assert.NoError(t, err, "Setting %s should succeed with size limit %d", key, tc.maxSize)

					// Verify key was set correctly
					has, err := store.Has([]byte(key))
					assert.NoError(t, err)
					assert.True(t, has, "Key %s should exist after successful Set", key)

					val, err := store.Get([]byte(key))
					assert.NoError(t, err)
					assert.Equal(t, []byte(value), val)
				}
			}
		})
	}
}

// TestMaxVersionsLimit tests the max versions limit functionality
func TestMaxVersionsLimit(t *testing.T) {
	logger := mocks.NewMockLogger()

	tests := []struct {
		name           string
		maxVersions    int
		numVersions    int
		expectedKept   []int64
		expectedPurged []int64
	}{
		{
			name:           "Default versions limit",
			maxVersions:    3,
			numVersions:    5,
			expectedKept:   []int64{3, 4, 5},
			expectedPurged: []int64{1, 2},
		},
		{
			name:           "High versions limit",
			maxVersions:    10,
			numVersions:    5,
			expectedKept:   []int64{1, 2, 3, 4, 5},
			expectedPurged: []int64{},
		},
		{
			name:           "No versions limit",
			maxVersions:    0,
			numVersions:    5,
			expectedKept:   []int64{1, 2, 3, 4, 5},
			expectedPurged: []int64{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a store with configured version limit
			options := Options{
				MaxVersions: tc.maxVersions,
			}
			store := NewStore("test-store", "/test/path", options, logger)
			defer store.Close()

			versioned, ok := interface{}(store).(storage.Versioned)
			require.True(t, ok)

			// Create versions
			for i := 0; i < tc.numVersions; i++ {
				err := store.Set([]byte("key"), []byte{byte(i)})
				require.NoError(t, err)

				v, _, err := versioned.SaveVersion()
				require.NoError(t, err)
				require.Equal(t, int64(i+1), v)
			}

			// Check versions
			versions := versioned.ListVersions()
			require.Len(t, versions, len(tc.expectedKept))

			for _, v := range tc.expectedKept {
				require.Contains(t, versions, v)
			}

			for _, v := range tc.expectedPurged {
				require.NotContains(t, versions, v)
			}
		})
	}
}

// TestEngineFunctionality tests engine creation and opening of stores
func TestEngineFunctionality(t *testing.T) {
	// Use mock logger
	mockLogger := mocks.NewMockLogger()
	engine := NewEngine(mockLogger)

	// Test capabilities
	capabilities := engine.Capabilities()
	assert.True(t, capabilities.Transactions)
	assert.True(t, capabilities.Versioning)
	assert.True(t, capabilities.RangeQueries)
	assert.False(t, capabilities.Persistence)
	assert.False(t, capabilities.Compression)

	// Test store creation
	store, err := engine.Create("test-store", "/test/path", nil)
	require.NoError(t, err)
	require.NotNil(t, store)

	// Test creating a store with the same name
	_, err = engine.Create("test-store", "/test/path", nil)
	require.Error(t, err)
	require.True(t, storage.IsStoreExists(err))

	// Test opening a store
	openedStore, err := engine.Open("test-store", "/test/path")
	require.NoError(t, err)
	require.NotNil(t, openedStore)

	// Test opening a non-existent store
	_, err = engine.Open("non-existent-store", "/test/path")
	require.Error(t, err)
	require.True(t, storage.IsStoreNotFound(err))

	// Test closing and reopening
	err = store.Close()
	require.NoError(t, err)

	_, err = engine.Open("test-store", "/test/path")
	require.Error(t, err)
	require.True(t, storage.IsStoreClosed(err))
}

// TestCloseStoreWithActiveTransaction tests closing a store with an active transaction
func TestCloseStoreWithActiveTransaction(t *testing.T) {
	logger := mocks.NewMockLogger()
	store := NewStore("test-store", "/test/path", Options{MaxVersions: 10}, logger)

	// Start a transaction
	txStore, ok := interface{}(store).(storage.Transactional)
	require.True(t, ok)

	transaction, err := txStore.BeginTx()
	require.NoError(t, err)

	// Set a value in the transaction
	err = transaction.Set([]byte("tx-key"), []byte("tx-value"))
	require.NoError(t, err)

	// Close the store without committing the transaction
	err = store.Close()
	require.NoError(t, err)

	// Try to commit the transaction - should fail
	err = transaction.Commit()
	require.Error(t, err)
	require.True(t, storage.IsStoreClosed(err))
}

// TestStoreDelete tests the delete functionality of the store
func TestStoreDelete(t *testing.T) {
	logger := mocks.NewMockLogger()
	store := NewStore("test-store", "/test/path", Options{}, logger)
	defer store.Close()

	// Set some key-value pairs
	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for k, v := range testData {
		err := store.Set([]byte(k), []byte(v))
		require.NoError(t, err)
	}

	// Test deleting existing key
	err := store.Delete([]byte("key1"))
	require.NoError(t, err)

	// Verify key was deleted
	has, err := store.Has([]byte("key1"))
	require.NoError(t, err)
	require.False(t, has)

	// Test deleting non-existent key
	err = store.Delete([]byte("non-existent-key"))
	require.NoError(t, err) // Delete should not error on non-existent keys

	// Verify other keys are still intact
	for k, v := range testData {
		if k == "key1" {
			continue // Already deleted
		}

		val, err := store.Get([]byte(k))
		require.NoError(t, err)
		require.Equal(t, []byte(v), val)
	}
}

// TestStoreClose tests closing a store properly discards all data
func TestStoreClose(t *testing.T) {
	logger := mocks.NewMockLogger()
	store := NewStore("test-store", "/test/path", Options{}, logger)

	// Set some key-value pairs
	err := store.Set([]byte("key1"), []byte("value1"))
	require.NoError(t, err)

	// Close the store
	err = store.Close()
	require.NoError(t, err)

	// Verify operations fail after close
	_, err = store.Get([]byte("key1"))
	require.Error(t, err)
	require.True(t, storage.IsStoreClosed(err))

	err = store.Set([]byte("key2"), []byte("value2"))
	require.Error(t, err)
	require.True(t, storage.IsStoreClosed(err))

	err = store.Delete([]byte("key1"))
	require.Error(t, err)
	require.True(t, storage.IsStoreClosed(err))

	_, err = store.Has([]byte("key1"))
	require.Error(t, err)
	require.True(t, storage.IsStoreClosed(err))
}

// TestStoreLogger tests that the logger is used
func TestStoreLogger(t *testing.T) {
	mockLogger := mocks.NewMockLogger()

	infoCount := 0
	mockLogger.InfoFunc = func(format string, args ...interface{}) {
		infoCount++
	}

	debugCount := 0
	mockLogger.DebugFunc = func(format string, args ...interface{}) {
		debugCount++
	}

	store := NewStore("test-store", "/test/path", Options{}, mockLogger)
	defer store.Close()

	// Perform operations that should trigger logging
	store.Set([]byte("key1"), []byte("value1"))
	store.Get([]byte("key1"))
	store.Delete([]byte("key1"))

	// Verify logging occurred
	assert.Greater(t, len(mockLogger.InfoCalls)+len(mockLogger.DebugCalls), 0)
	assert.Greater(t, infoCount+debugCount, 0)
}
