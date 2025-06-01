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
