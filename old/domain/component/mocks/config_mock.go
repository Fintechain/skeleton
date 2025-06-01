package mocks

import (
	"fmt"
	"time"
)

// MockConfig implements a mock configuration for testing
type MockConfig struct {
	values map[string]interface{}
}

// NewMockConfig creates a new mock configuration
func NewMockConfig() *MockConfig {
	return &MockConfig{
		values: make(map[string]interface{}),
	}
}

// SetValue sets a configuration value
func (m *MockConfig) SetValue(key string, value interface{}) {
	m.values[key] = value
}

// GetString retrieves a string configuration value
func (m *MockConfig) GetString(key string) string {
	return m.GetStringDefault(key, "")
}

// GetStringDefault retrieves a string configuration value with a default value
func (m *MockConfig) GetStringDefault(key, defaultValue string) string {
	if value, exists := m.values[key]; exists {
		if strValue, ok := value.(string); ok {
			return strValue
		}
	}
	return defaultValue
}

// GetInt retrieves an integer configuration value
func (m *MockConfig) GetInt(key string) (int, error) {
	if value, exists := m.values[key]; exists {
		if intValue, ok := value.(int); ok {
			return intValue, nil
		}
		return 0, fmt.Errorf("value for key %s is not an integer", key)
	}
	return 0, fmt.Errorf("key %s not found", key)
}

// GetIntDefault retrieves an integer configuration value with a default value
func (m *MockConfig) GetIntDefault(key string, defaultValue int) int {
	value, err := m.GetInt(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetBool retrieves a boolean configuration value
func (m *MockConfig) GetBool(key string) (bool, error) {
	if value, exists := m.values[key]; exists {
		if boolValue, ok := value.(bool); ok {
			return boolValue, nil
		}
		return false, fmt.Errorf("value for key %s is not a boolean", key)
	}
	return false, fmt.Errorf("key %s not found", key)
}

// GetBoolDefault retrieves a boolean configuration value with a default value
func (m *MockConfig) GetBoolDefault(key string, defaultValue bool) bool {
	value, err := m.GetBool(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetDuration retrieves a duration configuration value
func (m *MockConfig) GetDuration(key string) (time.Duration, error) {
	if value, exists := m.values[key]; exists {
		switch v := value.(type) {
		case time.Duration:
			return v, nil
		case string:
			duration, err := time.ParseDuration(v)
			if err != nil {
				return 0, fmt.Errorf("failed to parse duration from string: %w", err)
			}
			return duration, nil
		case int:
			return time.Duration(v) * time.Millisecond, nil
		case int64:
			return time.Duration(v) * time.Millisecond, nil
		}
		return 0, fmt.Errorf("value for key %s cannot be converted to duration", key)
	}
	return 0, fmt.Errorf("key %s not found", key)
}

// GetDurationDefault retrieves a duration configuration value with a default value
func (m *MockConfig) GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	value, err := m.GetDuration(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetObject deserializes a configuration section into a struct
func (m *MockConfig) GetObject(key string, result interface{}) error {
	// In a real implementation, this would deserialize the value to the result object
	// For this mock, we'll just return an error if the key doesn't exist
	if _, exists := m.values[key]; !exists {
		return fmt.Errorf("key %s not found", key)
	}
	return nil
}

// Exists checks if a configuration key exists
func (m *MockConfig) Exists(key string) bool {
	_, exists := m.values[key]
	return exists
}
