// Package config provides public APIs for the configuration system.
package config

import (
	"time"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/infrastructure/config"
)

// ===== CONFIGURATION INTERFACES =====

// Configuration provides access to application configuration values.
type Configuration = config.Configuration

// ConfigurationSource provides configuration values from a source.
type ConfigurationSource = config.ConfigurationSource

// ===== CONFIGURATION ERROR CONSTANTS =====

// Common configuration error codes
const (
	ErrConfigNotFound   = "config.not_found"
	ErrConfigWrongType  = "config.wrong_type"
	ErrConfigLoadFailed = "config.load_failed"
	ErrConfigInvalid    = "config.invalid"
	ErrConfigMissing    = "config.missing_required"
)

// ===== ERROR HANDLING =====

// Error represents a domain-specific error from the configuration system.
type Error = component.Error

// NewError creates a new configuration error with the given code, message, and optional cause.
func NewError(code, message string, cause error) *Error {
	return component.NewError(code, message, cause)
}

// IsConfigError checks if an error is a configuration error with the given code.
func IsConfigError(err error, code string) bool {
	return component.IsComponentError(err, code)
}

// ===== CONFIGURATION CONSTRUCTORS =====

// NewConfiguration creates a new configuration with the given name and optional source.
// This is the primary way to create a Configuration instance.
func NewConfiguration(name string, source ...ConfigurationSource) Configuration {
	var configSource ConfigurationSource
	if len(source) > 0 {
		configSource = source[0]
	}
	return config.CreateConfigWithSource(name, configSource)
}

// NewDefaultConfiguration creates a new configuration with default settings.
// This provides a convenient way to create a basic configuration instance.
func NewDefaultConfiguration() Configuration {
	return config.CreateDefaultConfig()
}

// ===== CONFIGURATION UTILITIES =====

// GetString retrieves a string configuration value.
func GetString(cfg Configuration, key string) string {
	return cfg.GetString(key)
}

// GetStringDefault retrieves a string configuration value with a default.
func GetStringDefault(cfg Configuration, key, defaultValue string) string {
	return cfg.GetStringDefault(key, defaultValue)
}

// GetInt retrieves an integer configuration value.
func GetInt(cfg Configuration, key string) (int, error) {
	return cfg.GetInt(key)
}

// GetIntDefault retrieves an integer configuration value with a default.
func GetIntDefault(cfg Configuration, key string, defaultValue int) int {
	return cfg.GetIntDefault(key, defaultValue)
}

// GetBool retrieves a boolean configuration value.
func GetBool(cfg Configuration, key string) (bool, error) {
	return cfg.GetBool(key)
}

// GetBoolDefault retrieves a boolean configuration value with a default.
func GetBoolDefault(cfg Configuration, key string, defaultValue bool) bool {
	return cfg.GetBoolDefault(key, defaultValue)
}

// GetDuration retrieves a duration configuration value.
func GetDuration(cfg Configuration, key string) (time.Duration, error) {
	return cfg.GetDuration(key)
}

// GetDurationDefault retrieves a duration configuration value with a default.
func GetDurationDefault(cfg Configuration, key string, defaultValue time.Duration) time.Duration {
	return cfg.GetDurationDefault(key, defaultValue)
}

// GetObject deserializes a configuration section into a struct.
func GetObject(cfg Configuration, key string, result interface{}) error {
	return cfg.GetObject(key, result)
}

// Exists checks if a configuration key exists.
func Exists(cfg Configuration, key string) bool {
	return cfg.Exists(key)
}
