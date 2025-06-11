package config

import (
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/fintechain/skeleton/internal/domain/config"
	infraConfig "github.com/fintechain/skeleton/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileSource(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "config_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	t.Run("NewFileSource", func(t *testing.T) {
		source := infraConfig.NewFileSource("test.json")
		assert.NotNil(t, source)
	})

	t.Run("LoadConfig_EmptyPath", func(t *testing.T) {
		source := infraConfig.NewFileSource("")
		err := source.LoadConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), config.ErrInvalidConfigValue)
	})

	t.Run("LoadConfig_NonexistentFile", func(t *testing.T) {
		source := infraConfig.NewFileSource("nonexistent.json")
		err := source.LoadConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), config.ErrConfigSaveFailed)
	})

	t.Run("LoadConfig_InvalidJSON", func(t *testing.T) {
		// Create invalid JSON file
		invalidPath := filepath.Join(tempDir, "invalid.json")
		err := os.WriteFile(invalidPath, []byte("invalid json"), 0644)
		require.NoError(t, err)

		source := infraConfig.NewFileSource(invalidPath)
		err = source.LoadConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), config.ErrInvalidConfigFormat)
	})

	t.Run("LoadConfig_ValidJSON", func(t *testing.T) {
		// Create valid JSON file
		validPath := filepath.Join(tempDir, "valid.json")
		validJSON := `{
			"string": "value",
			"number": 42,
			"nested": {
				"key": "nested_value"
			}
		}`
		err := os.WriteFile(validPath, []byte(validJSON), 0644)
		require.NoError(t, err)

		source := infraConfig.NewFileSource(validPath)
		err = source.LoadConfig()
		assert.NoError(t, err)

		// Test value retrieval
		value, exists := source.GetValue("string")
		assert.True(t, exists)
		assert.Equal(t, "value", value)

		value, exists = source.GetValue("number")
		assert.True(t, exists)
		assert.Equal(t, float64(42), value) // JSON numbers are float64

		value, exists = source.GetValue("nested.key")
		assert.True(t, exists)
		assert.Equal(t, "nested_value", value)
	})

	t.Run("GetValue_EmptyKey", func(t *testing.T) {
		source := infraConfig.NewFileSource("test.json")
		value, exists := source.GetValue("")
		assert.False(t, exists)
		assert.Nil(t, value)
	})

	t.Run("GetValue_NonexistentKey", func(t *testing.T) {
		validPath := filepath.Join(tempDir, "empty.json")
		err := os.WriteFile(validPath, []byte("{}"), 0644)
		require.NoError(t, err)

		source := infraConfig.NewFileSource(validPath)
		err = source.LoadConfig()
		require.NoError(t, err)

		value, exists := source.GetValue("nonexistent")
		assert.False(t, exists)
		assert.Nil(t, value)
	})

	t.Run("GetValue_NestedKeys", func(t *testing.T) {
		// Create JSON with nested structure
		validPath := filepath.Join(tempDir, "nested.json")
		validJSON := `{
			"database": {
				"host": "localhost",
				"port": 5432,
				"credentials": {
					"username": "admin",
					"password": "secret"
				}
			}
		}`
		err := os.WriteFile(validPath, []byte(validJSON), 0644)
		require.NoError(t, err)

		source := infraConfig.NewFileSource(validPath)
		err = source.LoadConfig()
		require.NoError(t, err)

		// Test nested key access
		value, exists := source.GetValue("database.host")
		assert.True(t, exists)
		assert.Equal(t, "localhost", value)

		value, exists = source.GetValue("database.credentials.username")
		assert.True(t, exists)
		assert.Equal(t, "admin", value)

		// Test invalid nested paths
		value, exists = source.GetValue("database.host.invalid")
		assert.False(t, exists)
		assert.Nil(t, value)

		value, exists = source.GetValue("database.nonexistent")
		assert.False(t, exists)
		assert.Nil(t, value)
	})

	t.Run("ThreadSafety", func(t *testing.T) {
		// Create initial JSON file
		validPath := filepath.Join(tempDir, "concurrent.json")
		validJSON := `{"key": "value"}`
		err := os.WriteFile(validPath, []byte(validJSON), 0644)
		require.NoError(t, err)

		source := infraConfig.NewFileSource(validPath)
		err = source.LoadConfig()
		require.NoError(t, err)

		var wg sync.WaitGroup
		numGoroutines := 10

		// Concurrent reads
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < 100; j++ {
					source.GetValue("key")
					source.GetValue("database.host")
				}
			}()
		}

		// Concurrent reloads
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				source.LoadConfig()
			}()
		}

		wg.Wait()
	})
}
