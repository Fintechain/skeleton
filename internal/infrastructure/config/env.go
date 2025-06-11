// Package config provides infrastructure implementations for configuration management.
package config

import (
	"os"
	"strings"
	"sync"
)

// EnvSource implements config.ConfigurationSource for environment variables.
type EnvSource struct {
	prefix string
	cache  map[string]interface{}
	mu     sync.RWMutex
}

// NewEnvSource creates a new environment variable configuration source.
// The prefix parameter is optional and will be used to filter environment variables.
// If prefix is empty, all environment variables will be accessible.
func NewEnvSource(prefix string) *EnvSource {
	return &EnvSource{
		prefix: prefix,
		cache:  make(map[string]interface{}),
	}
}

// LoadConfig loads all environment variables into the cache.
// Environment variables are filtered by prefix if one was specified.
// Keys are normalized by:
// 1. Removing the prefix
// 2. Converting to lowercase
// 3. Converting underscores to dots for nested keys
func (e *EnvSource) LoadConfig() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Clear existing cache
	e.cache = make(map[string]interface{})

	// Get all environment variables
	for _, env := range os.Environ() {
		key, value, found := strings.Cut(env, "=")
		if !found {
			continue
		}

		// Skip if doesn't match prefix
		if e.prefix != "" && !strings.HasPrefix(key, e.prefix) {
			continue
		}

		// Remove prefix if present
		if e.prefix != "" {
			key = strings.TrimPrefix(key, e.prefix)
		}

		// Convert to lowercase and replace underscores with dots
		key = strings.ToLower(key)
		key = strings.ReplaceAll(key, "_", ".")

		// Remove leading dot if present
		key = strings.TrimPrefix(key, ".")

		// Store in cache
		e.cache[key] = value
	}

	return nil
}

// GetValue retrieves a configuration value by key.
// The key should be in dot notation format (e.g., "database.host").
// Returns the value and true if found, nil and false if not found.
func (e *EnvSource) GetValue(key string) (interface{}, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Convert key to lowercase for case-insensitive lookup
	key = strings.ToLower(key)

	value, ok := e.cache[key]
	return value, ok
}
