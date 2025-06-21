// Package config provides infrastructure implementations for configuration management.
package config

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fintechain/skeleton/internal/domain/config"
)

// MemorySource implements the ConfigurationSource interface using in-memory storage.
type MemorySource struct {
	data map[string]interface{}
	mu   sync.RWMutex
}

// NewMemorySource creates a new memory-based configuration source.
func NewMemorySource() *MemorySource {
	return &MemorySource{
		data: make(map[string]interface{}),
	}
}

// NewMemorySourceWithData creates a new memory-based configuration source with initial data.
func NewMemorySourceWithData(data map[string]interface{}) *MemorySource {
	source := NewMemorySource()
	if data != nil {
		// Use SetValue to properly handle dot notation keys
		for k, v := range data {
			source.SetValue(k, v)
		}
	}
	return source
}

// LoadConfig loads configuration data from the source (no-op for memory source).
func (s *MemorySource) LoadConfig() error {
	return nil // No-op for memory source
}

// GetValue retrieves a raw configuration value by key.
func (s *MemorySource) GetValue(key string) (interface{}, bool) {
	if key == "" {
		return nil, false
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	// Support nested keys using dot notation
	parts := strings.Split(key, ".")
	current := s.data

	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part - return the value
			value, exists := current[part]
			return value, exists
		}

		// Navigate deeper into nested structure
		if nested, ok := current[part].(map[string]interface{}); ok {
			current = nested
		} else {
			return nil, false
		}
	}

	return nil, false
}

// SetValue sets a configuration value by key.
func (s *MemorySource) SetValue(key string, value interface{}) {
	if key == "" {
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Support nested keys using dot notation
	parts := strings.Split(key, ".")
	current := s.data

	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part - set the value
			current[part] = value
			return
		}

		// Navigate or create nested structure
		if nested, ok := current[part].(map[string]interface{}); ok {
			current = nested
		} else {
			// Create new nested map
			newMap := make(map[string]interface{})
			current[part] = newMap
			current = newMap
		}
	}
}

// SetValues sets multiple configuration values.
func (s *MemorySource) SetValues(values map[string]interface{}) {
	if values == nil {
		return
	}

	for key, value := range values {
		s.SetValue(key, value)
	}
}

// Clear removes all configuration values.
func (s *MemorySource) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data = make(map[string]interface{})
}

// GetAllKeys returns all configuration keys.
func (s *MemorySource) GetAllKeys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var keys []string
	s.collectKeys("", s.data, &keys)
	return keys
}

// collectKeys recursively collects all keys from nested maps.
func (s *MemorySource) collectKeys(prefix string, data map[string]interface{}, keys *[]string) {
	for key, value := range data {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		if nested, ok := value.(map[string]interface{}); ok {
			s.collectKeys(fullKey, nested, keys)
		} else {
			*keys = append(*keys, fullKey)
		}
	}
}

// MemoryConfiguration implements the Configuration interface using in-memory storage.
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
func (c *MemoryConfiguration) GetString(key string) string {
	value, exists := c.source.GetValue(key)
	if !exists {
		return ""
	}

	return fmt.Sprintf("%v", value)
}

// GetStringDefault retrieves a string configuration value with a default fallback.
func (c *MemoryConfiguration) GetStringDefault(key, defaultValue string) string {
	value, exists := c.source.GetValue(key)
	if !exists {
		return defaultValue
	}

	return fmt.Sprintf("%v", value)
}

// GetInt retrieves an integer configuration value.
func (c *MemoryConfiguration) GetInt(key string) (int, error) {
	value, exists := c.source.GetValue(key)
	if !exists {
		return 0, fmt.Errorf(config.ErrConfigKeyNotFound)
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
		return 0, fmt.Errorf(config.ErrInvalidConfigType)
	default:
		return 0, fmt.Errorf(config.ErrInvalidConfigType)
	}
}

// GetIntDefault retrieves an integer configuration value with a default fallback.
func (c *MemoryConfiguration) GetIntDefault(key string, defaultValue int) int {
	value, err := c.GetInt(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetBool retrieves a boolean configuration value.
func (c *MemoryConfiguration) GetBool(key string) (bool, error) {
	value, exists := c.source.GetValue(key)
	if !exists {
		return false, fmt.Errorf(config.ErrConfigKeyNotFound)
	}

	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		switch strings.ToLower(v) {
		case "true", "1", "yes", "on":
			return true, nil
		case "false", "0", "no", "off":
			return false, nil
		default:
			return false, fmt.Errorf(config.ErrInvalidConfigType)
		}
	default:
		return false, fmt.Errorf(config.ErrInvalidConfigType)
	}
}

// GetBoolDefault retrieves a boolean configuration value with a default fallback.
func (c *MemoryConfiguration) GetBoolDefault(key string, defaultValue bool) bool {
	value, err := c.GetBool(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetDuration retrieves a duration configuration value.
func (c *MemoryConfiguration) GetDuration(key string) (time.Duration, error) {
	value, exists := c.source.GetValue(key)
	if !exists {
		return 0, fmt.Errorf(config.ErrConfigKeyNotFound)
	}

	switch v := value.(type) {
	case time.Duration:
		return v, nil
	case string:
		if parsed, err := time.ParseDuration(v); err == nil {
			return parsed, nil
		}
		return 0, fmt.Errorf(config.ErrInvalidConfigType)
	case int64:
		return time.Duration(v), nil
	case float64:
		return time.Duration(v), nil
	default:
		return 0, fmt.Errorf(config.ErrInvalidConfigType)
	}
}

// GetDurationDefault retrieves a duration configuration value with a default fallback.
func (c *MemoryConfiguration) GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	value, err := c.GetDuration(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetObject deserializes a configuration section into a struct.
func (c *MemoryConfiguration) GetObject(key string, result interface{}) error {
	if result == nil {
		return fmt.Errorf(config.ErrInvalidConfigValue)
	}

	value, exists := c.source.GetValue(key)
	if !exists {
		return fmt.Errorf(config.ErrConfigKeyNotFound)
	}

	// Use JSON marshaling/unmarshaling for object conversion
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf(config.ErrInvalidConfigType)
	}

	if err := json.Unmarshal(jsonData, result); err != nil {
		return fmt.Errorf(config.ErrInvalidConfigType)
	}

	return nil
}

// Exists checks whether a configuration key exists.
func (c *MemoryConfiguration) Exists(key string) bool {
	_, exists := c.source.GetValue(key)
	return exists
}

// SetValue sets a configuration value (helper method for testing).
func (c *MemoryConfiguration) SetValue(key string, value interface{}) {
	c.source.SetValue(key, value)
}

// SetValues sets multiple configuration values (helper method for testing).
func (c *MemoryConfiguration) SetValues(values map[string]interface{}) {
	c.source.SetValues(values)
}

// Clear removes all configuration values (helper method for testing).
func (c *MemoryConfiguration) Clear() {
	c.source.Clear()
}

// GetSource returns the underlying configuration source (helper method for testing).
func (c *MemoryConfiguration) GetSource() *MemorySource {
	return c.source
}
