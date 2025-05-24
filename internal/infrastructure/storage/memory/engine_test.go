package memory

import (
	"testing"

	"github.com/ebanfa/skeleton/internal/domain/storage"
	"github.com/ebanfa/skeleton/internal/infrastructure/storage/memory/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEngineCreation(t *testing.T) {
	// Test with valid logger
	mockLogger := mocks.NewMockLogger()
	engine := NewEngine(mockLogger)
	assert.NotNil(t, engine)

	// Check engine name
	assert.Equal(t, "memory", engine.Name())

	// Check engine capabilities
	capabilities := engine.Capabilities()
	assert.True(t, capabilities.Transactions)
	assert.True(t, capabilities.Versioning)
	assert.True(t, capabilities.RangeQueries)
	assert.False(t, capabilities.Persistence)
	assert.False(t, capabilities.Compression)
}

func TestEngineCreateStore(t *testing.T) {
	mockLogger := mocks.NewMockLogger()
	engine := NewEngine(mockLogger)

	tests := []struct {
		name         string
		storeName    string
		storePath    string
		config       storage.Config
		expectError  bool
		errorMatcher func(error) bool
	}{
		{
			name:        "Valid store creation",
			storeName:   "test-store-1",
			storePath:   "/test/path",
			config:      nil,
			expectError: false,
		},
		{
			name:        "Store with config",
			storeName:   "test-store-2",
			storePath:   "/test/path",
			config:      storage.Config{storage.ConfigCacheSize: int64(1000)},
			expectError: false,
		},
		{
			name:         "Duplicate store name",
			storeName:    "test-store-1", // Same as first test case
			storePath:    "/test/path",
			config:       nil,
			expectError:  true,
			errorMatcher: storage.IsStoreExists,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			store, err := engine.Create(tc.storeName, tc.storePath, tc.config)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMatcher != nil {
					assert.True(t, tc.errorMatcher(err), "Error doesn't match expected type")
				}
				assert.Nil(t, store)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, store)

				// Verify store properties
				assert.Equal(t, tc.storeName, store.Name())
				assert.Equal(t, tc.storePath, store.Path())

				// Verify store implements required interfaces
				_, ok := store.(storage.Store)
				assert.True(t, ok, "Store should implement storage.Store")

				_, ok = store.(storage.Transactional)
				assert.True(t, ok, "Store should implement storage.Transactional")

				_, ok = store.(storage.Versioned)
				assert.True(t, ok, "Store should implement storage.Versioned")

				_, ok = store.(storage.RangeQueryable)
				assert.True(t, ok, "Store should implement storage.RangeQueryable")
			}
		})
	}
}

func TestEngineOpenStore(t *testing.T) {
	mockLogger := mocks.NewMockLogger()
	engine := NewEngine(mockLogger)

	// Create a test store
	storeName := "test-open-store"
	storePath := "/test/open/path"
	store, err := engine.Create(storeName, storePath, nil)
	require.NoError(t, err)
	require.NotNil(t, store)

	tests := []struct {
		name         string
		storeName    string
		storePath    string
		closeFirst   bool
		expectError  bool
		errorMatcher func(error) bool
	}{
		{
			name:        "Open existing store",
			storeName:   storeName,
			storePath:   storePath,
			closeFirst:  false,
			expectError: false,
		},
		{
			name:         "Open non-existent store",
			storeName:    "non-existent-store",
			storePath:    storePath,
			closeFirst:   false,
			expectError:  true,
			errorMatcher: storage.IsStoreNotFound,
		},
		{
			name:         "Open closed store",
			storeName:    storeName,
			storePath:    storePath,
			closeFirst:   true,
			expectError:  true,
			errorMatcher: storage.IsStoreClosed,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.closeFirst && tc.storeName == storeName {
				// Close the store first if required by the test
				err := store.Close()
				require.NoError(t, err)
			}

			openedStore, err := engine.Open(tc.storeName, tc.storePath)

			if tc.expectError {
				assert.Error(t, err)
				if tc.errorMatcher != nil {
					assert.True(t, tc.errorMatcher(err), "Error doesn't match expected type")
				}
				assert.Nil(t, openedStore)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, openedStore)
				assert.Equal(t, tc.storeName, openedStore.Name())
				assert.Equal(t, tc.storePath, openedStore.Path())
			}
		})
	}
}

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name         string
		config       storage.Config
		expectedSize int64
		expectedVers int
	}{
		{
			name:         "Nil config",
			config:       nil,
			expectedSize: -1,
			expectedVers: 100,
		},
		{
			name:         "Empty config",
			config:       storage.Config{},
			expectedSize: -1,
			expectedVers: 100,
		},
		{
			name:         "Int64 max size",
			config:       storage.Config{storage.ConfigCacheSize: int64(2048)},
			expectedSize: 2048,
			expectedVers: 100,
		},
		{
			name:         "Int max size",
			config:       storage.Config{storage.ConfigCacheSize: 1024},
			expectedSize: 1024,
			expectedVers: 100,
		},
		{
			name:         "Max versions",
			config:       storage.Config{storage.ConfigMaxVersions: 50},
			expectedSize: -1,
			expectedVers: 50,
		},
		{
			name: "Combined config",
			config: storage.Config{
				storage.ConfigCacheSize:   int64(4096),
				storage.ConfigMaxVersions: 20,
			},
			expectedSize: 4096,
			expectedVers: 20,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			options := parseConfig(tc.config)
			assert.Equal(t, tc.expectedSize, options.MaxSize)
			assert.Equal(t, tc.expectedVers, options.MaxVersions)
		})
	}
}

func TestEngineNilLogger(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r, "Creating an engine with nil logger should panic")
	}()

	// This should panic
	_ = NewEngine(nil)
}

func TestEngineLoggerUsage(t *testing.T) {
	mockLogger := mocks.NewMockLogger()
	engine := NewEngine(mockLogger)

	// Create a store to trigger log messages
	_, err := engine.Create("test-log-store", "/test/path", nil)
	assert.NoError(t, err)

	// Verify debug log was called
	assert.Greater(t, len(mockLogger.DebugCalls), 0)
}
