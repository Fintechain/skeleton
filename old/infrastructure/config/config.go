// Package config provides configuration handling for the application.
package config

import (
	"time"
)

// Configuration provides access to application configuration values
type Configuration interface {
	// GetString retrieves a string configuration value
	GetString(key string) string

	// GetStringDefault retrieves a string configuration value with a default value
	GetStringDefault(key, defaultValue string) string

	// GetInt retrieves an integer configuration value
	GetInt(key string) (int, error)

	// GetIntDefault retrieves an integer configuration value with a default value
	GetIntDefault(key string, defaultValue int) int

	// GetBool retrieves a boolean configuration value
	GetBool(key string) (bool, error)

	// GetBoolDefault retrieves a boolean configuration value with a default value
	GetBoolDefault(key string, defaultValue bool) bool

	// GetDuration retrieves a duration configuration value
	GetDuration(key string) (time.Duration, error)

	// GetDurationDefault retrieves a duration configuration value with a default value
	GetDurationDefault(key string, defaultValue time.Duration) time.Duration

	// GetObject deserializes a configuration section into a struct
	GetObject(key string, result interface{}) error

	// Exists checks if a configuration key exists
	Exists(key string) bool
}

// ConfigurationSource provides configuration values from a source
type ConfigurationSource interface {
	// LoadConfig loads configuration from the source
	LoadConfig() error

	// GetValue retrieves a raw configuration value
	GetValue(key string) (interface{}, bool)
}

// Common error codes for configuration operations
const (
	ErrConfigNotFound   = "config.not_found"
	ErrConfigWrongType  = "config.wrong_type"
	ErrConfigLoadFailed = "config.load_failed"
)

// DefaultConfig represents a basic configuration implementation.
type DefaultConfig struct {
	// Name is the configuration name
	Name string

	// Properties is a map of configuration properties
	Properties map[string]interface{}

	// Source is the configuration source
	Source ConfigurationSource
}

// DefaultConfigOptions contains options for creating a DefaultConfig
type DefaultConfigOptions struct {
	Name   string
	Source ConfigurationSource
}

// NewDefaultConfig creates a new configuration with the given options.
// It follows the constructor injection pattern for the ConfigurationSource dependency.
func NewDefaultConfig(options DefaultConfigOptions) *DefaultConfig {
	return &DefaultConfig{
		Name:       options.Name,
		Properties: make(map[string]interface{}),
		Source:     options.Source,
	}
}

// CreateDefaultConfig is a factory method for backward compatibility.
// Creates a config with a default name and no source.
func CreateDefaultConfig() *DefaultConfig {
	return NewDefaultConfig(DefaultConfigOptions{
		Name: "default",
	})
}

// CreateConfigWithSource is a factory method that creates a configuration with the given source.
func CreateConfigWithSource(name string, source ConfigurationSource) *DefaultConfig {
	return NewDefaultConfig(DefaultConfigOptions{
		Name:   name,
		Source: source,
	})
}

// Set sets a configuration property.
func (c *DefaultConfig) Set(key string, value interface{}) {
	c.Properties[key] = value
}

// Get gets a configuration property.
func (c *DefaultConfig) Get(key string) interface{} {
	// Try to get from properties first
	if value, exists := c.Properties[key]; exists {
		return value
	}

	// If we have a source, try to get from there
	if c.Source != nil {
		if value, exists := c.Source.GetValue(key); exists {
			return value
		}
	}

	return nil
}

// Exists checks if a configuration key exists.
func (c *DefaultConfig) Exists(key string) bool {
	// Check in properties
	if _, exists := c.Properties[key]; exists {
		return true
	}

	// Check in source if available
	if c.Source != nil {
		if _, exists := c.Source.GetValue(key); exists {
			return true
		}
	}

	return false
}

// GetString gets a string configuration property.
func (c *DefaultConfig) GetString(key string) string {
	return c.GetStringDefault(key, "")
}

// GetStringDefault gets a string configuration property with default.
func (c *DefaultConfig) GetStringDefault(key string, defaultValue string) string {
	value := c.Get(key)
	if value == nil {
		return defaultValue
	}

	if strValue, ok := value.(string); ok {
		return strValue
	}

	return defaultValue
}

// GetInt gets an int configuration property.
func (c *DefaultConfig) GetInt(key string) (int, error) {
	value := c.Get(key)
	if value == nil {
		return 0, &ConfigError{
			Code:    ErrConfigNotFound,
			Message: "Configuration key not found",
			Details: map[string]interface{}{"key": key},
		}
	}

	if intValue, ok := value.(int); ok {
		return intValue, nil
	}

	return 0, &ConfigError{
		Code:    ErrConfigWrongType,
		Message: "Configuration value is not an integer",
		Details: map[string]interface{}{"key": key, "value": value},
	}
}

// GetIntDefault gets an int configuration property with default.
func (c *DefaultConfig) GetIntDefault(key string, defaultValue int) int {
	value, err := c.GetInt(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetBool gets a boolean configuration property.
func (c *DefaultConfig) GetBool(key string) (bool, error) {
	value := c.Get(key)
	if value == nil {
		return false, &ConfigError{
			Code:    ErrConfigNotFound,
			Message: "Configuration key not found",
			Details: map[string]interface{}{"key": key},
		}
	}

	if boolValue, ok := value.(bool); ok {
		return boolValue, nil
	}

	return false, &ConfigError{
		Code:    ErrConfigWrongType,
		Message: "Configuration value is not a boolean",
		Details: map[string]interface{}{"key": key, "value": value},
	}
}

// GetBoolDefault gets a boolean configuration property with default.
func (c *DefaultConfig) GetBoolDefault(key string, defaultValue bool) bool {
	value, err := c.GetBool(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetDuration gets a duration configuration property.
func (c *DefaultConfig) GetDuration(key string) (time.Duration, error) {
	value := c.Get(key)
	if value == nil {
		return 0, &ConfigError{
			Code:    ErrConfigNotFound,
			Message: "Configuration key not found",
			Details: map[string]interface{}{"key": key},
		}
	}

	// Handle different types that can represent duration
	switch v := value.(type) {
	case time.Duration:
		return v, nil
	case string:
		duration, err := time.ParseDuration(v)
		if err != nil {
			return 0, &ConfigError{
				Code:    ErrConfigWrongType,
				Message: "Failed to parse duration from string",
				Details: map[string]interface{}{"key": key, "value": v, "error": err.Error()},
			}
		}
		return duration, nil
	case int:
		return time.Duration(v) * time.Millisecond, nil
	case int64:
		return time.Duration(v) * time.Millisecond, nil
	default:
		return 0, &ConfigError{
			Code:    ErrConfigWrongType,
			Message: "Configuration value cannot be converted to duration",
			Details: map[string]interface{}{"key": key, "value": value},
		}
	}
}

// GetDurationDefault gets a duration configuration property with default.
func (c *DefaultConfig) GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	value, err := c.GetDuration(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetObject deserializes a configuration section into a struct.
func (c *DefaultConfig) GetObject(key string, result interface{}) error {
	// This is a stub implementation to be replaced with actual deserialization logic
	// For example, using JSON or other serialization format
	return &ConfigError{
		Code:    ErrConfigNotFound,
		Message: "GetObject not implemented",
		Details: map[string]interface{}{"key": key},
	}
}

// ConfigError represents a configuration error.
type ConfigError struct {
	Code    string
	Message string
	Details map[string]interface{}
}

// Error returns the error message.
func (e *ConfigError) Error() string {
	return e.Message
}
