// Package config provides concrete implementations for the configuration domain interfaces.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fintechain/skeleton/internal/domain/config"
)

// MemorySource implements the ConfigurationSource interface with in-memory storage.
// It provides thread-safe programmatic configuration setting and retrieval,
// making it ideal for testing scenarios.
type MemorySource struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

// NewMemorySource creates a new in-memory configuration source.
// This is the primary constructor for creating memory-based configuration sources.
func NewMemorySource() *MemorySource {
	return &MemorySource{
		data: make(map[string]interface{}),
	}
}

// NewMemorySourceWithData creates a new in-memory configuration source with initial data.
// The data map is copied to ensure isolation from external modifications.
func NewMemorySourceWithData(data map[string]interface{}) *MemorySource {
	source := &MemorySource{
		data: make(map[string]interface{}, len(data)),
	}

	// Deep copy the data to ensure isolation
	for k, v := range data {
		source.data[k] = v
	}

	return source
}

// LoadConfig loads configuration data from the source.
// For memory source, this is a no-op since data is already in memory.
func (m *MemorySource) LoadConfig() error {
	// Memory source doesn't need to load from external source
	return nil
}

// GetValue retrieves a raw configuration value by key.
// Returns the value and true if found, nil and false if not found.
func (m *MemorySource) GetValue(key string) (interface{}, bool) {
	if key == "" {
		return nil, false
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	// Support nested keys with dot notation (e.g., "database.host")
	value, exists := m.getNestedValue(key)
	return value, exists
}

// SetValue sets a configuration value by key.
// This is an additional helper method for programmatic configuration.
func (m *MemorySource) SetValue(key string, value interface{}) {
	if key == "" {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Support nested keys with dot notation
	m.setNestedValue(key, value)
}

// SetValues sets multiple configuration values at once.
// This is useful for bulk configuration updates.
func (m *MemorySource) SetValues(values map[string]interface{}) {
	if values == nil {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for key, value := range values {
		if key != "" {
			m.setNestedValue(key, value)
		}
	}
}

// Clear removes all configuration values.
// This is useful for test cleanup scenarios.
func (m *MemorySource) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[string]interface{})
}

// GetAllKeys returns all configuration keys.
// This is useful for debugging and introspection.
func (m *MemorySource) GetAllKeys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.getAllKeysRecursive("", m.data)
}

// getNestedValue retrieves a value using dot notation for nested keys.
func (m *MemorySource) getNestedValue(key string) (interface{}, bool) {
	parts := strings.Split(key, ".")
	current := m.data

	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part - return the value
			value, exists := current[part]
			return value, exists
		}

		// Intermediate part - navigate deeper
		next, exists := current[part]
		if !exists {
			return nil, false
		}

		// Convert to map[string]interface{} if possible
		if nextMap, ok := next.(map[string]interface{}); ok {
			current = nextMap
		} else {
			// Can't navigate deeper
			return nil, false
		}
	}

	return nil, false
}

// setNestedValue sets a value using dot notation for nested keys.
func (m *MemorySource) setNestedValue(key string, value interface{}) {
	parts := strings.Split(key, ".")
	current := m.data

	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part - set the value
			current[part] = value
			return
		}

		// Intermediate part - ensure nested map exists
		next, exists := current[part]
		if !exists {
			// Create new nested map
			next = make(map[string]interface{})
			current[part] = next
		}

		// Convert to map[string]interface{} if possible
		if nextMap, ok := next.(map[string]interface{}); ok {
			current = nextMap
		} else {
			// Overwrite with new map
			nextMap := make(map[string]interface{})
			current[part] = nextMap
			current = nextMap
		}
	}
}

// getAllKeysRecursive recursively collects all keys with dot notation.
func (m *MemorySource) getAllKeysRecursive(prefix string, data map[string]interface{}) []string {
	var keys []string

	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		if valueMap, ok := value.(map[string]interface{}); ok {
			// Recursively get keys from nested map
			nestedKeys := m.getAllKeysRecursive(fullKey, valueMap)
			keys = append(keys, nestedKeys...)
		} else {
			// Leaf value
			keys = append(keys, fullKey)
		}
	}

	return keys
}

// MemoryConfiguration implements the Configuration interface using a memory source.
// This provides a complete configuration implementation for testing scenarios.
type MemoryConfiguration struct {
	source *MemorySource
}

// NewMemoryConfiguration creates a new memory-based configuration.
func NewMemoryConfiguration() *MemoryConfiguration {
	return &MemoryConfiguration{
		source: NewMemorySource(),
	}
}

// NewMemoryConfigurationWithData creates a new memory-based configuration with initial data.
func NewMemoryConfigurationWithData(data map[string]interface{}) *MemoryConfiguration {
	return &MemoryConfiguration{
		source: NewMemorySourceWithData(data),
	}
}

// GetString retrieves a string configuration value.
func (m *MemoryConfiguration) GetString(key string) string {
	value, exists := m.source.GetValue(key)
	if !exists {
		return ""
	}

	if str, ok := value.(string); ok {
		return str
	}

	// Try to convert to string
	return fmt.Sprintf("%v", value)
}

// GetStringDefault retrieves a string configuration value with a default fallback.
func (m *MemoryConfiguration) GetStringDefault(key, defaultValue string) string {
	value, exists := m.source.GetValue(key)
	if !exists {
		return defaultValue
	}

	if str, ok := value.(string); ok {
		return str
	}

	// Try to convert to string
	return fmt.Sprintf("%v", value)
}

// GetInt retrieves an integer configuration value.
func (m *MemoryConfiguration) GetInt(key string) (int, error) {
	value, exists := m.source.GetValue(key)
	if !exists {
		return 0, errors.New(config.ErrConfigKeyNotFound)
	}

	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed, nil
		}
		return 0, errors.New(config.ErrInvalidConfigType)
	default:
		return 0, errors.New(config.ErrInvalidConfigType)
	}
}

// GetIntDefault retrieves an integer configuration value with a default fallback.
func (m *MemoryConfiguration) GetIntDefault(key string, defaultValue int) int {
	value, err := m.GetInt(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetBool retrieves a boolean configuration value.
func (m *MemoryConfiguration) GetBool(key string) (bool, error) {
	value, exists := m.source.GetValue(key)
	if !exists {
		return false, errors.New(config.ErrConfigKeyNotFound)
	}

	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		if parsed, err := strconv.ParseBool(v); err == nil {
			return parsed, nil
		}
		return false, errors.New(config.ErrInvalidConfigType)
	default:
		return false, errors.New(config.ErrInvalidConfigType)
	}
}

// GetBoolDefault retrieves a boolean configuration value with a default fallback.
func (m *MemoryConfiguration) GetBoolDefault(key string, defaultValue bool) bool {
	value, err := m.GetBool(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetDuration retrieves a duration configuration value.
func (m *MemoryConfiguration) GetDuration(key string) (time.Duration, error) {
	value, exists := m.source.GetValue(key)
	if !exists {
		return 0, errors.New(config.ErrConfigKeyNotFound)
	}

	switch v := value.(type) {
	case time.Duration:
		return v, nil
	case string:
		if parsed, err := time.ParseDuration(v); err == nil {
			return parsed, nil
		}
		return 0, errors.New(config.ErrInvalidConfigType)
	case int64:
		return time.Duration(v), nil
	case float64:
		return time.Duration(v), nil
	default:
		return 0, errors.New(config.ErrInvalidConfigType)
	}
}

// GetDurationDefault retrieves a duration configuration value with a default fallback.
func (m *MemoryConfiguration) GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	value, err := m.GetDuration(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetObject deserializes a configuration section into a struct.
func (m *MemoryConfiguration) GetObject(key string, result interface{}) error {
	if result == nil {
		return errors.New(config.ErrInvalidConfigValue)
	}

	value, exists := m.source.GetValue(key)
	if !exists {
		return errors.New(config.ErrConfigKeyNotFound)
	}

	// Use JSON marshaling/unmarshaling for object conversion
	jsonData, err := json.Marshal(value)
	if err != nil {
		return errors.New(config.ErrInvalidConfigType)
	}

	if err := json.Unmarshal(jsonData, result); err != nil {
		return errors.New(config.ErrInvalidConfigType)
	}

	return nil
}

// Exists checks whether a configuration key exists.
func (m *MemoryConfiguration) Exists(key string) bool {
	_, exists := m.source.GetValue(key)
	return exists
}

// SetValue sets a configuration value (additional helper for testing).
func (m *MemoryConfiguration) SetValue(key string, value interface{}) {
	m.source.SetValue(key, value)
}

// SetValues sets multiple configuration values (additional helper for testing).
func (m *MemoryConfiguration) SetValues(values map[string]interface{}) {
	m.source.SetValues(values)
}

// Clear removes all configuration values (additional helper for testing).
func (m *MemoryConfiguration) Clear() {
	m.source.Clear()
}

// GetSource returns the underlying memory source for advanced operations.
func (m *MemoryConfiguration) GetSource() *MemorySource {
	return m.source
}
