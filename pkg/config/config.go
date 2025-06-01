// Package config provides configuration management interfaces and types.
package config

import (
	"github.com/fintechain/skeleton/internal/domain/config"
	configImpl "github.com/fintechain/skeleton/internal/infrastructure/config"
)

// Re-export configuration interfaces
type Configuration = config.Configuration
type ConfigurationSource = config.ConfigurationSource

// Re-export configuration error constants
const (
	ErrConfigNotFound   = config.ErrConfigNotFound
	ErrConfigWrongType  = config.ErrConfigWrongType
	ErrConfigLoadFailed = config.ErrConfigLoadFailed
)

// NewConfiguration creates a new Configuration instance with the provided configuration sources.
// This factory function provides access to the concrete configuration implementation.
func NewConfiguration(sources ...ConfigurationSource) Configuration {
	return configImpl.NewConfiguration(sources...)
}

// NewMemoryConfigurationSource creates a new in-memory configuration source.
// This factory function provides access to the memory configuration source implementation.
func NewMemoryConfigurationSource(values map[string]interface{}) ConfigurationSource {
	return configImpl.NewMemoryConfigurationSource(values)
}
