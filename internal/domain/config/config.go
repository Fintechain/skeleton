// Package config provides centralized configuration management for Fintechain Skeleton.
package config

import (
	"time"
)

// Configuration provides type-safe access to application configuration values.
type Configuration interface {
	// GetString retrieves a string configuration value.
	GetString(key string) string

	// GetStringDefault retrieves a string configuration value with a default fallback.
	GetStringDefault(key, defaultValue string) string

	// GetInt retrieves an integer configuration value.
	GetInt(key string) (int, error)

	// GetIntDefault retrieves an integer configuration value with a default fallback.
	GetIntDefault(key string, defaultValue int) int

	// GetBool retrieves a boolean configuration value.
	GetBool(key string) (bool, error)

	// GetBoolDefault retrieves a boolean configuration value with a default fallback.
	GetBoolDefault(key string, defaultValue bool) bool

	// GetDuration retrieves a duration configuration value.
	GetDuration(key string) (time.Duration, error)

	// GetDurationDefault retrieves a duration configuration value with a default fallback.
	GetDurationDefault(key string, defaultValue time.Duration) time.Duration

	// GetObject deserializes a configuration section into a struct.
	GetObject(key string, result interface{}) error

	// Exists checks whether a configuration key exists.
	Exists(key string) bool
}

// ConfigurationSource provides configuration values from a specific source.
type ConfigurationSource interface {
	// LoadConfig loads configuration data from the source.
	LoadConfig() error

	// GetValue retrieves a raw configuration value by key.
	GetValue(key string) (interface{}, bool)
}

// Common error codes for configuration operations.
const (
	// ErrConfigNotFound indicates that a requested configuration key does not exist.
	ErrConfigNotFound = "config.not_found"

	// ErrConfigWrongType indicates that a configuration value cannot be converted to the requested type.
	ErrConfigWrongType = "config.wrong_type"

	// ErrConfigLoadFailed indicates that a configuration source failed to load.
	ErrConfigLoadFailed = "config.load_failed"
)
