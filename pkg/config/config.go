// Package config provides configuration management exports.
package config

import (
	"github.com/fintechain/skeleton/internal/domain/config"
	infraConfig "github.com/fintechain/skeleton/internal/infrastructure/config"
)

// Configuration provides type-safe access to application configuration values.
type Configuration = config.Configuration

// ConfigurationSource provides configuration values from a specific source.
type ConfigurationSource = config.ConfigurationSource

// Error constants
const (
	ErrConfigKeyNotFound      = config.ErrConfigKeyNotFound
	ErrInvalidConfigValue     = config.ErrInvalidConfigValue
	ErrInvalidConfigType      = config.ErrInvalidConfigType
	ErrConfigReadOnly         = config.ErrConfigReadOnly
	ErrConfigSaveFailed       = config.ErrConfigSaveFailed
	ErrInvalidConfigFormat    = config.ErrInvalidConfigFormat
	ErrConfigValidationFailed = config.ErrConfigValidationFailed
)

// Factory functions for memory-based configuration

// NewMemoryConfiguration creates a new memory-based configuration.
func NewMemoryConfiguration() Configuration {
	return infraConfig.NewMemoryConfiguration()
}

// NewMemoryConfigurationWithData creates a new memory-based configuration with initial data.
func NewMemoryConfigurationWithData(data map[string]interface{}) Configuration {
	return infraConfig.NewMemoryConfigurationWithData(data)
}

// NewMemorySource creates a new memory-based configuration source.
func NewMemorySource() ConfigurationSource {
	return infraConfig.NewMemorySource()
}

// NewMemorySourceWithData creates a new memory-based configuration source with initial data.
func NewMemorySourceWithData(data map[string]interface{}) ConfigurationSource {
	return infraConfig.NewMemorySourceWithData(data)
}
