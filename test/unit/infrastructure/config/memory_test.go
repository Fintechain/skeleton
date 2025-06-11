package config

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/fintechain/skeleton/internal/domain/config"
	infraConfig "github.com/fintechain/skeleton/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
)

// TestMemorySource tests the MemorySource implementation
func TestMemorySource(t *testing.T) {
	t.Run("NewMemorySource", func(t *testing.T) {
		source := infraConfig.NewMemorySource()
		assert.NotNil(t, source)

		// Should be empty initially
		keys := source.GetAllKeys()
		assert.Empty(t, keys)
	})

	t.Run("NewMemorySourceWithData", func(t *testing.T) {
		data := map[string]interface{}{
			"key1": "value1",
			"key2": 42,
			"nested": map[string]interface{}{
				"key": "nested_value",
			},
		}

		source := infraConfig.NewMemorySourceWithData(data)
		assert.NotNil(t, source)

		// Should contain the initial data
		value, exists := source.GetValue("key1")
		assert.True(t, exists)
		assert.Equal(t, "value1", value)

		value, exists = source.GetValue("key2")
		assert.True(t, exists)
		assert.Equal(t, 42, value)

		// Test data isolation - modifying original map shouldn't affect source
		data["key1"] = "modified"
		value, exists = source.GetValue("key1")
		assert.True(t, exists)
		assert.Equal(t, "value1", value) // Should still be original value
	})

	t.Run("LoadConfig", func(t *testing.T) {
		source := infraConfig.NewMemorySource()
		err := source.LoadConfig()
		assert.NoError(t, err) // Should be no-op for memory source
	})

	t.Run("GetValue_SetValue", func(t *testing.T) {
		source := infraConfig.NewMemorySource()

		// Test non-existent key
		value, exists := source.GetValue("nonexistent")
		assert.False(t, exists)
		assert.Nil(t, value)

		// Test setting and getting values
		source.SetValue("string_key", "string_value")
		value, exists = source.GetValue("string_key")
		assert.True(t, exists)
		assert.Equal(t, "string_value", value)

		source.SetValue("int_key", 123)
		value, exists = source.GetValue("int_key")
		assert.True(t, exists)
		assert.Equal(t, 123, value)

		source.SetValue("bool_key", true)
		value, exists = source.GetValue("bool_key")
		assert.True(t, exists)
		assert.Equal(t, true, value)

		// Test empty key
		source.SetValue("", "should_not_set")
		value, exists = source.GetValue("")
		assert.False(t, exists)
	})

	t.Run("NestedKeys", func(t *testing.T) {
		source := infraConfig.NewMemorySource()

		// Set nested values using dot notation
		source.SetValue("database.host", "localhost")
		source.SetValue("database.port", 5432)
		source.SetValue("database.credentials.username", "admin")
		source.SetValue("database.credentials.password", "secret")

		// Test retrieval
		value, exists := source.GetValue("database.host")
		assert.True(t, exists)
		assert.Equal(t, "localhost", value)

		value, exists = source.GetValue("database.port")
		assert.True(t, exists)
		assert.Equal(t, 5432, value)

		value, exists = source.GetValue("database.credentials.username")
		assert.True(t, exists)
		assert.Equal(t, "admin", value)

		value, exists = source.GetValue("database.credentials.password")
		assert.True(t, exists)
		assert.Equal(t, "secret", value)

		// Test non-existent nested key
		value, exists = source.GetValue("database.nonexistent")
		assert.False(t, exists)

		value, exists = source.GetValue("nonexistent.key")
		assert.False(t, exists)
	})

	t.Run("SetValues", func(t *testing.T) {
		source := infraConfig.NewMemorySource()

		values := map[string]interface{}{
			"key1":       "value1",
			"key2":       42,
			"nested.key": "nested_value",
		}

		source.SetValues(values)

		// Verify all values were set
		for key, expectedValue := range values {
			value, exists := source.GetValue(key)
			assert.True(t, exists, "Key %s should exist", key)
			assert.Equal(t, expectedValue, value, "Value for key %s should match", key)
		}

		// Test with nil map
		source.SetValues(nil) // Should not panic

		// Test with empty key in map
		source.SetValues(map[string]interface{}{
			"":     "should_not_set",
			"key3": "should_set",
		})

		value, exists := source.GetValue("")
		assert.False(t, exists)

		value, exists = source.GetValue("key3")
		assert.True(t, exists)
		assert.Equal(t, "should_set", value)
	})

	t.Run("Clear", func(t *testing.T) {
		source := infraConfig.NewMemorySource()

		// Set some values
		source.SetValue("key1", "value1")
		source.SetValue("key2", "value2")

		// Verify they exist
		keys := source.GetAllKeys()
		assert.Len(t, keys, 2)

		// Clear all values
		source.Clear()

		// Verify they're gone
		keys = source.GetAllKeys()
		assert.Empty(t, keys)

		value, exists := source.GetValue("key1")
		assert.False(t, exists)
		assert.Nil(t, value)
	})

	t.Run("GetAllKeys", func(t *testing.T) {
		source := infraConfig.NewMemorySource()

		// Set various keys including nested ones
		source.SetValue("simple", "value")
		source.SetValue("nested.key1", "value1")
		source.SetValue("nested.key2", "value2")
		source.SetValue("deeply.nested.key", "deep_value")

		keys := source.GetAllKeys()

		expectedKeys := []string{
			"simple",
			"nested.key1",
			"nested.key2",
			"deeply.nested.key",
		}

		assert.Len(t, keys, len(expectedKeys))
		for _, expectedKey := range expectedKeys {
			assert.Contains(t, keys, expectedKey)
		}
	})

	t.Run("ThreadSafety", func(t *testing.T) {
		source := infraConfig.NewMemorySource()

		const numGoroutines = 100
		const numOperations = 100

		var wg sync.WaitGroup

		// Concurrent writes
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					key := fmt.Sprintf("key_%d_%d", id, j)
					value := fmt.Sprintf("value_%d_%d", id, j)
					source.SetValue(key, value)
				}
			}(i)
		}

		// Concurrent reads
		for i := 0; i < numGoroutines; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				for j := 0; j < numOperations; j++ {
					key := fmt.Sprintf("key_%d_%d", id, j)
					source.GetValue(key) // Don't care about result, just testing for races
				}
			}(i)
		}

		wg.Wait()

		// Verify some data was written
		keys := source.GetAllKeys()
		assert.NotEmpty(t, keys)
	})
}

// TestMemoryConfiguration tests the MemoryConfiguration implementation
func TestMemoryConfiguration(t *testing.T) {
	t.Run("NewMemoryConfiguration", func(t *testing.T) {
		cfg := infraConfig.NewMemoryConfiguration()
		assert.NotNil(t, cfg)
		assert.NotNil(t, cfg.GetSource())
	})

	t.Run("NewMemoryConfigurationWithData", func(t *testing.T) {
		data := map[string]interface{}{
			"string_key": "string_value",
			"int_key":    42,
			"bool_key":   true,
		}

		cfg := infraConfig.NewMemoryConfigurationWithData(data)
		assert.NotNil(t, cfg)

		// Verify data is accessible
		assert.Equal(t, "string_value", cfg.GetString("string_key"))
		assert.True(t, cfg.Exists("int_key"))
		assert.True(t, cfg.Exists("bool_key"))
	})

	t.Run("GetString", func(t *testing.T) {
		cfg := infraConfig.NewMemoryConfiguration()

		// Test non-existent key
		assert.Equal(t, "", cfg.GetString("nonexistent"))

		// Test string value
		cfg.SetValue("string_key", "test_string")
		assert.Equal(t, "test_string", cfg.GetString("string_key"))

		// Test conversion from other types
		cfg.SetValue("int_key", 123)
		assert.Equal(t, "123", cfg.GetString("int_key"))

		cfg.SetValue("bool_key", true)
		assert.Equal(t, "true", cfg.GetString("bool_key"))
	})

	t.Run("GetStringDefault", func(t *testing.T) {
		cfg := infraConfig.NewMemoryConfiguration()

		// Test non-existent key with default
		assert.Equal(t, "default_value", cfg.GetStringDefault("nonexistent", "default_value"))

		// Test existing key
		cfg.SetValue("existing_key", "existing_value")
		assert.Equal(t, "existing_value", cfg.GetStringDefault("existing_key", "default_value"))

		// Test conversion with default
		cfg.SetValue("int_key", 456)
		assert.Equal(t, "456", cfg.GetStringDefault("int_key", "default_value"))
	})

	t.Run("GetInt", func(t *testing.T) {
		cfg := infraConfig.NewMemoryConfiguration()

		// Test non-existent key
		value, err := cfg.GetInt("nonexistent")
		assert.Error(t, err)
		assert.Equal(t, config.ErrConfigKeyNotFound, err.Error())
		assert.Equal(t, 0, value)

		// Test int value
		cfg.SetValue("int_key", 123)
		value, err = cfg.GetInt("int_key")
		assert.NoError(t, err)
		assert.Equal(t, 123, value)

		// Test int64 value
		cfg.SetValue("int64_key", int64(456))
		value, err = cfg.GetInt("int64_key")
		assert.NoError(t, err)
		assert.Equal(t, 456, value)

		// Test float64 value
		cfg.SetValue("float_key", 789.0)
		value, err = cfg.GetInt("float_key")
		assert.NoError(t, err)
		assert.Equal(t, 789, value)

		// Test string conversion
		cfg.SetValue("string_int_key", "999")
		value, err = cfg.GetInt("string_int_key")
		assert.NoError(t, err)
		assert.Equal(t, 999, value)

		// Test invalid string
		cfg.SetValue("invalid_string_key", "not_a_number")
		value, err = cfg.GetInt("invalid_string_key")
		assert.Error(t, err)
		assert.Equal(t, config.ErrInvalidConfigType, err.Error())
		assert.Equal(t, 0, value)

		// Test invalid type
		cfg.SetValue("invalid_type_key", []string{"array"})
		value, err = cfg.GetInt("invalid_type_key")
		assert.Error(t, err)
		assert.Equal(t, config.ErrInvalidConfigType, err.Error())
		assert.Equal(t, 0, value)
	})

	t.Run("GetIntDefault", func(t *testing.T) {
		cfg := infraConfig.NewMemoryConfiguration()

		// Test non-existent key with default
		assert.Equal(t, 42, cfg.GetIntDefault("nonexistent", 42))

		// Test existing valid key
		cfg.SetValue("int_key", 123)
		assert.Equal(t, 123, cfg.GetIntDefault("int_key", 42))

		// Test existing invalid key with default
		cfg.SetValue("invalid_key", "not_a_number")
		assert.Equal(t, 42, cfg.GetIntDefault("invalid_key", 42))
	})

	t.Run("GetBool", func(t *testing.T) {
		cfg := infraConfig.NewMemoryConfiguration()

		// Test non-existent key
		value, err := cfg.GetBool("nonexistent")
		assert.Error(t, err)
		assert.Equal(t, config.ErrConfigKeyNotFound, err.Error())
		assert.False(t, value)

		// Test bool value
		cfg.SetValue("bool_true", true)
		value, err = cfg.GetBool("bool_true")
		assert.NoError(t, err)
		assert.True(t, value)

		cfg.SetValue("bool_false", false)
		value, err = cfg.GetBool("bool_false")
		assert.NoError(t, err)
		assert.False(t, value)

		// Test string conversion
		cfg.SetValue("string_true", "true")
		value, err = cfg.GetBool("string_true")
		assert.NoError(t, err)
		assert.True(t, value)

		cfg.SetValue("string_false", "false")
		value, err = cfg.GetBool("string_false")
		assert.NoError(t, err)
		assert.False(t, value)

		cfg.SetValue("string_1", "1")
		value, err = cfg.GetBool("string_1")
		assert.NoError(t, err)
		assert.True(t, value)

		cfg.SetValue("string_0", "0")
		value, err = cfg.GetBool("string_0")
		assert.NoError(t, err)
		assert.False(t, value)

		// Test invalid string
		cfg.SetValue("invalid_string", "maybe")
		value, err = cfg.GetBool("invalid_string")
		assert.Error(t, err)
		assert.Equal(t, config.ErrInvalidConfigType, err.Error())
		assert.False(t, value)

		// Test invalid type
		cfg.SetValue("invalid_type", 123)
		value, err = cfg.GetBool("invalid_type")
		assert.Error(t, err)
		assert.Equal(t, config.ErrInvalidConfigType, err.Error())
		assert.False(t, value)
	})

	t.Run("GetBoolDefault", func(t *testing.T) {
		cfg := infraConfig.NewMemoryConfiguration()

		// Test non-existent key with default
		assert.True(t, cfg.GetBoolDefault("nonexistent", true))
		assert.False(t, cfg.GetBoolDefault("nonexistent", false))

		// Test existing valid key
		cfg.SetValue("bool_key", true)
		assert.True(t, cfg.GetBoolDefault("bool_key", false))

		// Test existing invalid key with default
		cfg.SetValue("invalid_key", "maybe")
		assert.True(t, cfg.GetBoolDefault("invalid_key", true))
	})

	t.Run("GetDuration", func(t *testing.T) {
		cfg := infraConfig.NewMemoryConfiguration()

		// Test non-existent key
		value, err := cfg.GetDuration("nonexistent")
		assert.Error(t, err)
		assert.Equal(t, config.ErrConfigKeyNotFound, err.Error())
		assert.Equal(t, time.Duration(0), value)

		// Test duration value
		expectedDuration := 5 * time.Minute
		cfg.SetValue("duration_key", expectedDuration)
		value, err = cfg.GetDuration("duration_key")
		assert.NoError(t, err)
		assert.Equal(t, expectedDuration, value)

		// Test string conversion
		cfg.SetValue("string_duration", "10s")
		value, err = cfg.GetDuration("string_duration")
		assert.NoError(t, err)
		assert.Equal(t, 10*time.Second, value)

		cfg.SetValue("string_duration_complex", "1h30m45s")
		value, err = cfg.GetDuration("string_duration_complex")
		assert.NoError(t, err)
		assert.Equal(t, 1*time.Hour+30*time.Minute+45*time.Second, value)

		// Test int64 conversion
		cfg.SetValue("int64_duration", int64(1000000000)) // 1 second in nanoseconds
		value, err = cfg.GetDuration("int64_duration")
		assert.NoError(t, err)
		assert.Equal(t, time.Second, value)

		// Test float64 conversion
		cfg.SetValue("float64_duration", float64(2000000000)) // 2 seconds in nanoseconds
		value, err = cfg.GetDuration("float64_duration")
		assert.NoError(t, err)
		assert.Equal(t, 2*time.Second, value)

		// Test invalid string
		cfg.SetValue("invalid_string", "not_a_duration")
		value, err = cfg.GetDuration("invalid_string")
		assert.Error(t, err)
		assert.Equal(t, config.ErrInvalidConfigType, err.Error())
		assert.Equal(t, time.Duration(0), value)

		// Test invalid type
		cfg.SetValue("invalid_type", []string{"array"})
		value, err = cfg.GetDuration("invalid_type")
		assert.Error(t, err)
		assert.Equal(t, config.ErrInvalidConfigType, err.Error())
		assert.Equal(t, time.Duration(0), value)
	})

	t.Run("GetDurationDefault", func(t *testing.T) {
		cfg := infraConfig.NewMemoryConfiguration()
		defaultDuration := 30 * time.Second

		// Test non-existent key with default
		assert.Equal(t, defaultDuration, cfg.GetDurationDefault("nonexistent", defaultDuration))

		// Test existing valid key
		cfg.SetValue("duration_key", 1*time.Minute)
		assert.Equal(t, 1*time.Minute, cfg.GetDurationDefault("duration_key", defaultDuration))

		// Test existing invalid key with default
		cfg.SetValue("invalid_key", "not_a_duration")
		assert.Equal(t, defaultDuration, cfg.GetDurationDefault("invalid_key", defaultDuration))
	})

	t.Run("GetObject", func(t *testing.T) {
		cfg := infraConfig.NewMemoryConfiguration()

		// Test with nil result
		err := cfg.GetObject("key", nil)
		assert.Error(t, err)
		assert.Equal(t, config.ErrInvalidConfigValue, err.Error())

		// Test non-existent key
		var result map[string]interface{}
		err = cfg.GetObject("nonexistent", &result)
		assert.Error(t, err)
		assert.Equal(t, config.ErrConfigKeyNotFound, err.Error())

		// Test valid object
		testData := map[string]interface{}{
			"name":    "test",
			"value":   42,
			"enabled": true,
		}
		cfg.SetValue("object_key", testData)

		var retrievedData map[string]interface{}
		err = cfg.GetObject("object_key", &retrievedData)
		assert.NoError(t, err)
		// Note: JSON marshaling/unmarshaling converts int to float64
		expectedData := map[string]interface{}{
			"name":    "test",
			"value":   float64(42), // JSON converts int to float64
			"enabled": true,
		}
		assert.Equal(t, expectedData, retrievedData)

		// Test struct deserialization
		type TestStruct struct {
			Name    string `json:"name"`
			Value   int    `json:"value"`
			Enabled bool   `json:"enabled"`
		}

		var testStruct TestStruct
		err = cfg.GetObject("object_key", &testStruct)
		assert.NoError(t, err)
		assert.Equal(t, "test", testStruct.Name)
		assert.Equal(t, 42, testStruct.Value)
		assert.True(t, testStruct.Enabled)

		// Test invalid object (non-serializable)
		cfg.SetValue("invalid_object", make(chan int))
		var invalidResult map[string]interface{}
		err = cfg.GetObject("invalid_object", &invalidResult)
		assert.Error(t, err)
		assert.Equal(t, config.ErrInvalidConfigType, err.Error())
	})

	t.Run("Exists", func(t *testing.T) {
		cfg := infraConfig.NewMemoryConfiguration()

		// Test non-existent key
		assert.False(t, cfg.Exists("nonexistent"))

		// Test existing key
		cfg.SetValue("existing_key", "value")
		assert.True(t, cfg.Exists("existing_key"))

		// Test nested key
		cfg.SetValue("nested.key", "nested_value")
		assert.True(t, cfg.Exists("nested.key"))
	})

	t.Run("HelperMethods", func(t *testing.T) {
		cfg := infraConfig.NewMemoryConfiguration()

		// Test SetValue
		cfg.SetValue("test_key", "test_value")
		assert.Equal(t, "test_value", cfg.GetString("test_key"))

		// Test SetValues
		values := map[string]interface{}{
			"key1": "value1",
			"key2": 42,
			"key3": true,
		}
		cfg.SetValues(values)

		assert.Equal(t, "value1", cfg.GetString("key1"))
		assert.Equal(t, 42, cfg.GetIntDefault("key2", 0))
		assert.True(t, cfg.GetBoolDefault("key3", false))

		// Test Clear
		cfg.Clear()
		assert.False(t, cfg.Exists("test_key"))
		assert.False(t, cfg.Exists("key1"))
		assert.False(t, cfg.Exists("key2"))
		assert.False(t, cfg.Exists("key3"))

		// Test GetSource
		source := cfg.GetSource()
		assert.NotNil(t, source)
		assert.IsType(t, &infraConfig.MemorySource{}, source)
	})
}

// TestMemoryConfigurationIntegration tests integration scenarios
func TestMemoryConfigurationIntegration(t *testing.T) {
	t.Run("ComplexConfiguration", func(t *testing.T) {
		// Create a complex configuration structure
		configData := map[string]interface{}{
			"app": map[string]interface{}{
				"name":    "test-app",
				"version": "1.0.0",
				"debug":   true,
			},
			"database": map[string]interface{}{
				"host":    "localhost",
				"port":    5432,
				"name":    "testdb",
				"timeout": "30s",
				"ssl":     false,
				"pool": map[string]interface{}{
					"min_connections": 5,
					"max_connections": 20,
				},
			},
			"cache": map[string]interface{}{
				"enabled": true,
				"ttl":     "5m",
				"size":    1000,
			},
		}

		cfg := infraConfig.NewMemoryConfigurationWithData(configData)

		// Test nested string access
		assert.Equal(t, "test-app", cfg.GetString("app.name"))
		assert.Equal(t, "1.0.0", cfg.GetString("app.version"))
		assert.Equal(t, "localhost", cfg.GetString("database.host"))
		assert.Equal(t, "testdb", cfg.GetString("database.name"))

		// Test nested int access
		port, err := cfg.GetInt("database.port")
		assert.NoError(t, err)
		assert.Equal(t, 5432, port)

		minConn, err := cfg.GetInt("database.pool.min_connections")
		assert.NoError(t, err)
		assert.Equal(t, 5, minConn)

		maxConn, err := cfg.GetInt("database.pool.max_connections")
		assert.NoError(t, err)
		assert.Equal(t, 20, maxConn)

		cacheSize, err := cfg.GetInt("cache.size")
		assert.NoError(t, err)
		assert.Equal(t, 1000, cacheSize)

		// Test nested bool access
		debug, err := cfg.GetBool("app.debug")
		assert.NoError(t, err)
		assert.True(t, debug)

		ssl, err := cfg.GetBool("database.ssl")
		assert.NoError(t, err)
		assert.False(t, ssl)

		cacheEnabled, err := cfg.GetBool("cache.enabled")
		assert.NoError(t, err)
		assert.True(t, cacheEnabled)

		// Test nested duration access
		timeout, err := cfg.GetDuration("database.timeout")
		assert.NoError(t, err)
		assert.Equal(t, 30*time.Second, timeout)

		ttl, err := cfg.GetDuration("cache.ttl")
		assert.NoError(t, err)
		assert.Equal(t, 5*time.Minute, ttl)

		// Test object deserialization
		type DatabaseConfig struct {
			Host    string `json:"host"`
			Port    int    `json:"port"`
			Name    string `json:"name"`
			Timeout string `json:"timeout"`
			SSL     bool   `json:"ssl"`
			Pool    struct {
				MinConnections int `json:"min_connections"`
				MaxConnections int `json:"max_connections"`
			} `json:"pool"`
		}

		var dbConfig DatabaseConfig
		err = cfg.GetObject("database", &dbConfig)
		assert.NoError(t, err)
		assert.Equal(t, "localhost", dbConfig.Host)
		assert.Equal(t, 5432, dbConfig.Port)
		assert.Equal(t, "testdb", dbConfig.Name)
		assert.Equal(t, "30s", dbConfig.Timeout)
		assert.False(t, dbConfig.SSL)
		assert.Equal(t, 5, dbConfig.Pool.MinConnections)
		assert.Equal(t, 20, dbConfig.Pool.MaxConnections)
	})

	t.Run("DynamicConfiguration", func(t *testing.T) {
		cfg := infraConfig.NewMemoryConfiguration()

		// Start with basic configuration
		cfg.SetValue("feature.enabled", false)
		cfg.SetValue("feature.threshold", 100)

		assert.False(t, cfg.GetBoolDefault("feature.enabled", true))
		assert.Equal(t, 100, cfg.GetIntDefault("feature.threshold", 0))

		// Dynamically update configuration
		cfg.SetValue("feature.enabled", true)
		cfg.SetValue("feature.threshold", 200)
		cfg.SetValue("feature.timeout", "10s")

		assert.True(t, cfg.GetBoolDefault("feature.enabled", false))
		assert.Equal(t, 200, cfg.GetIntDefault("feature.threshold", 0))

		timeout, err := cfg.GetDuration("feature.timeout")
		assert.NoError(t, err)
		assert.Equal(t, 10*time.Second, timeout)

		// Bulk update
		updates := map[string]interface{}{
			"feature.max_retries":               3,
			"feature.backoff":                   "1s",
			"feature.circuit_breaker.enabled":   true,
			"feature.circuit_breaker.threshold": 5,
		}
		cfg.SetValues(updates)

		assert.Equal(t, 3, cfg.GetIntDefault("feature.max_retries", 0))

		backoff, err := cfg.GetDuration("feature.backoff")
		assert.NoError(t, err)
		assert.Equal(t, time.Second, backoff)

		assert.True(t, cfg.GetBoolDefault("feature.circuit_breaker.enabled", false))
		assert.Equal(t, 5, cfg.GetIntDefault("feature.circuit_breaker.threshold", 0))
	})
}

// BenchmarkMemoryConfiguration benchmarks the memory configuration performance
func BenchmarkMemoryConfiguration(b *testing.B) {
	cfg := infraConfig.NewMemoryConfiguration()

	// Setup test data
	cfg.SetValue("string_key", "test_value")
	cfg.SetValue("int_key", 42)
	cfg.SetValue("bool_key", true)
	cfg.SetValue("duration_key", "5s")
	cfg.SetValue("nested.key", "nested_value")

	b.Run("GetString", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cfg.GetString("string_key")
		}
	})

	b.Run("GetInt", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cfg.GetInt("int_key")
		}
	})

	b.Run("GetBool", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cfg.GetBool("bool_key")
		}
	})

	b.Run("GetDuration", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cfg.GetDuration("duration_key")
		}
	})

	b.Run("GetNestedKey", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cfg.GetString("nested.key")
		}
	})

	b.Run("SetValue", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cfg.SetValue("benchmark_key", i)
		}
	})

	b.Run("Exists", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			cfg.Exists("string_key")
		}
	})
}

// BenchmarkMemoryConfigurationConcurrent benchmarks concurrent access
func BenchmarkMemoryConfigurationConcurrent(b *testing.B) {
	cfg := infraConfig.NewMemoryConfiguration()

	// Setup test data
	for i := 0; i < 100; i++ {
		cfg.SetValue(fmt.Sprintf("key_%d", i), fmt.Sprintf("value_%d", i))
	}

	b.Run("ConcurrentReads", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				key := fmt.Sprintf("key_%d", i%100)
				cfg.GetString(key)
				i++
			}
		})
	})

	b.Run("ConcurrentWrites", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				key := fmt.Sprintf("bench_key_%d", i)
				value := fmt.Sprintf("bench_value_%d", i)
				cfg.SetValue(key, value)
				i++
			}
		})
	})

	b.Run("ConcurrentMixed", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				if i%2 == 0 {
					// Read operation
					key := fmt.Sprintf("key_%d", i%100)
					cfg.GetString(key)
				} else {
					// Write operation
					key := fmt.Sprintf("mixed_key_%d", i)
					value := fmt.Sprintf("mixed_value_%d", i)
					cfg.SetValue(key, value)
				}
				i++
			}
		})
	})
}
