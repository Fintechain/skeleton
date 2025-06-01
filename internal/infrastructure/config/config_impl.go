package config

import (
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/fintechain/skeleton/internal/domain/config"
)

// DefaultConfiguration provides a concrete implementation of the Configuration interface.
type DefaultConfiguration struct {
	sources []config.ConfigurationSource

	// Cached values for performance
	mu     sync.RWMutex
	values map[string]interface{}
}

// NewConfiguration creates a new Configuration instance with the provided configuration sources.
// This constructor accepts configuration source interface dependencies for testability.
func NewConfiguration(sources ...config.ConfigurationSource) config.Configuration {
	cfg := &DefaultConfiguration{
		sources: sources,
		values:  make(map[string]interface{}),
	}

	// Load configuration from all sources
	cfg.loadFromSources()

	return cfg
}

// loadFromSources loads configuration values from all sources.
func (c *DefaultConfiguration) loadFromSources() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Load from each source in order (later sources override earlier ones)
	for _, source := range c.sources {
		if err := source.LoadConfig(); err != nil {
			// Log error but continue with other sources
			continue
		}

		// This is a simplified approach - in a real implementation,
		// you would iterate through all keys from the source
		// For now, we'll just mark that sources are loaded
	}
}

// getValue retrieves a value from sources, checking cache first.
func (c *DefaultConfiguration) getValue(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Check cache first
	if value, exists := c.values[key]; exists {
		return value, true
	}

	// Check sources in reverse order (last source wins)
	for i := len(c.sources) - 1; i >= 0; i-- {
		if value, exists := c.sources[i].GetValue(key); exists {
			// Cache the value
			c.values[key] = value
			return value, true
		}
	}

	return nil, false
}

// GetString retrieves a string configuration value.
func (c *DefaultConfiguration) GetString(key string) string {
	value, exists := c.getValue(key)
	if !exists {
		return ""
	}

	if str, ok := value.(string); ok {
		return str
	}

	// Try to convert to string
	return fmt.Sprintf("%v", value)
}

// GetStringDefault retrieves a string configuration value with a default value.
func (c *DefaultConfiguration) GetStringDefault(key, defaultValue string) string {
	if !c.Exists(key) {
		return defaultValue
	}
	return c.GetString(key)
}

// GetInt retrieves an integer configuration value.
func (c *DefaultConfiguration) GetInt(key string) (int, error) {
	value, exists := c.getValue(key)
	if !exists {
		return 0, fmt.Errorf("%s: key %s not found", config.ErrConfigNotFound, key)
	}

	switch v := value.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return 0, fmt.Errorf("%s: cannot convert %s to int: %w", config.ErrConfigWrongType, key, err)
		}
		return parsed, nil
	default:
		return 0, fmt.Errorf("%s: key %s is not an integer", config.ErrConfigWrongType, key)
	}
}

// GetIntDefault retrieves an integer configuration value with a default value.
func (c *DefaultConfiguration) GetIntDefault(key string, defaultValue int) int {
	value, err := c.GetInt(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetBool retrieves a boolean configuration value.
func (c *DefaultConfiguration) GetBool(key string) (bool, error) {
	value, exists := c.getValue(key)
	if !exists {
		return false, fmt.Errorf("%s: key %s not found", config.ErrConfigNotFound, key)
	}

	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		parsed, err := strconv.ParseBool(v)
		if err != nil {
			return false, fmt.Errorf("%s: cannot convert %s to bool: %w", config.ErrConfigWrongType, key, err)
		}
		return parsed, nil
	default:
		return false, fmt.Errorf("%s: key %s is not a boolean", config.ErrConfigWrongType, key)
	}
}

// GetBoolDefault retrieves a boolean configuration value with a default value.
func (c *DefaultConfiguration) GetBoolDefault(key string, defaultValue bool) bool {
	value, err := c.GetBool(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetDuration retrieves a duration configuration value.
func (c *DefaultConfiguration) GetDuration(key string) (time.Duration, error) {
	value, exists := c.getValue(key)
	if !exists {
		return 0, fmt.Errorf("%s: key %s not found", config.ErrConfigNotFound, key)
	}

	switch v := value.(type) {
	case time.Duration:
		return v, nil
	case string:
		parsed, err := time.ParseDuration(v)
		if err != nil {
			return 0, fmt.Errorf("%s: cannot convert %s to duration: %w", config.ErrConfigWrongType, key, err)
		}
		return parsed, nil
	case int64:
		return time.Duration(v), nil
	case float64:
		return time.Duration(v), nil
	default:
		return 0, fmt.Errorf("%s: key %s is not a duration", config.ErrConfigWrongType, key)
	}
}

// GetDurationDefault retrieves a duration configuration value with a default value.
func (c *DefaultConfiguration) GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	value, err := c.GetDuration(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetObject deserializes a configuration section into a struct.
func (c *DefaultConfiguration) GetObject(key string, result interface{}) error {
	value, exists := c.getValue(key)
	if !exists {
		return fmt.Errorf("%s: key %s not found", config.ErrConfigNotFound, key)
	}

	// Convert value to JSON and then unmarshal into result
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("%s: cannot marshal %s to JSON: %w", config.ErrConfigWrongType, key, err)
	}

	err = json.Unmarshal(jsonData, result)
	if err != nil {
		return fmt.Errorf("%s: cannot unmarshal %s into target type: %w", config.ErrConfigWrongType, key, err)
	}

	return nil
}

// Exists checks if a configuration key exists.
func (c *DefaultConfiguration) Exists(key string) bool {
	_, exists := c.getValue(key)
	return exists
}

// MemoryConfigurationSource provides an in-memory configuration source for testing and simple use cases.
type MemoryConfigurationSource struct {
	values map[string]interface{}
	loaded bool
}

// NewMemoryConfigurationSource creates a new in-memory configuration source.
func NewMemoryConfigurationSource(values map[string]interface{}) config.ConfigurationSource {
	return &MemoryConfigurationSource{
		values: values,
		loaded: false,
	}
}

// LoadConfig loads configuration from memory.
func (m *MemoryConfigurationSource) LoadConfig() error {
	m.loaded = true
	return nil
}

// GetValue retrieves a raw configuration value.
func (m *MemoryConfigurationSource) GetValue(key string) (interface{}, bool) {
	if !m.loaded {
		return nil, false
	}

	value, exists := m.values[key]
	return value, exists
}
