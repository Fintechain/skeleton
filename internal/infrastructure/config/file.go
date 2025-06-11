// Package config provides concrete implementations for the configuration domain interfaces.
package config

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
	"sync"

	"github.com/fintechain/skeleton/internal/domain/config"
)

// FileSource implements the ConfigurationSource interface with JSON file parsing.
// It provides thread-safe access to configuration values loaded from a JSON file.
type FileSource struct {
	path string
	data map[string]interface{}
	mu   sync.RWMutex
}

// NewFileSource creates a new file-based configuration source.
// The path parameter specifies the JSON file to load.
func NewFileSource(path string) *FileSource {
	return &FileSource{
		path: path,
		data: make(map[string]interface{}),
	}
}

// LoadConfig loads configuration data from the JSON file.
// This method is thread-safe and can be called multiple times to reload configuration.
func (f *FileSource) LoadConfig() error {
	if f.path == "" {
		return errors.New(config.ErrInvalidConfigValue)
	}

	// Read file contents
	data, err := os.ReadFile(f.path)
	if err != nil {
		return errors.New(config.ErrConfigSaveFailed)
	}

	// Parse JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return errors.New(config.ErrInvalidConfigFormat)
	}

	// Update configuration atomically
	f.mu.Lock()
	defer f.mu.Unlock()

	f.data = jsonData
	return nil
}

// GetValue retrieves a raw configuration value by key.
// Returns the value and true if found, nil and false if not found.
// This method supports nested keys with dot notation (e.g., "database.host").
func (f *FileSource) GetValue(key string) (interface{}, bool) {
	if key == "" {
		return nil, false
	}

	f.mu.RLock()
	defer f.mu.RUnlock()

	// Support nested keys with dot notation
	return f.getNestedValue(key)
}

// getNestedValue retrieves a value using dot notation for nested keys.
// This is an internal helper method that assumes the lock is already held.
func (f *FileSource) getNestedValue(key string) (interface{}, bool) {
	parts := strings.Split(key, ".")
	current := f.data

	for i, part := range parts {
		if i == len(parts)-1 {
			// Last part - return the value
			value, exists := current[part]
			return value, exists
		}

		// Intermediate part - navigate deeper
		next, exists := current[part]
		if !exists {
			return nil, false
		}

		// Convert to map[string]interface{} if possible
		if nextMap, ok := next.(map[string]interface{}); ok {
			current = nextMap
		} else {
			// Can't navigate deeper
			return nil, false
		}
	}

	return nil, false
}
