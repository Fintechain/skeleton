// Package config provides interfaces and types for the configuration system.
package config

// Standard config error codes
const (
	// ErrConfigKeyNotFound is returned when a configuration key doesn't exist
	ErrConfigKeyNotFound = "config.config_key_not_found"

	// ErrInvalidConfigValue is returned when an invalid configuration value is provided
	ErrInvalidConfigValue = "config.invalid_config_value"

	// ErrInvalidConfigType is returned when an invalid configuration type is provided
	ErrInvalidConfigType = "config.invalid_config_type"

	// ErrConfigReadOnly is returned when attempting to modify a read-only configuration
	ErrConfigReadOnly = "config.config_read_only"

	// ErrConfigSaveFailed is returned when configuration saving fails
	ErrConfigSaveFailed = "config.config_save_failed"

	// ErrInvalidConfigFormat is returned when an invalid configuration format is provided
	ErrInvalidConfigFormat = "config.invalid_config_format"

	// ErrConfigValidationFailed is returned when configuration validation fails
	ErrConfigValidationFailed = "config.config_validation_failed"
)
