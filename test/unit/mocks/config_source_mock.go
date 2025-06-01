// Package mocks provides centralized mock implementations for the Skeleton Framework.
// This file contains mocks for configuration source interfaces.
package mocks

import (
	"sync"

	"github.com/fintechain/skeleton/internal/domain/config"
)

// MockConfigurationSource implements the config.ConfigurationSource interface for testing.
type MockConfigurationSource struct {
	mu sync.RWMutex

	// Configuration data
	data map[string]interface{}

	// Behavior configuration
	shouldFail   bool
	failureError string

	// Call tracking
	loadConfigCalls int
	getValueCalls   int

	// Custom functions for advanced testing
	loadConfigFunc func() error
	getValueFunc   func(key string) (interface{}, bool)
}

// NewMockConfigurationSource creates a new mock configuration source with default behavior.
func NewMockConfigurationSource() *MockConfigurationSource {
	return &MockConfigurationSource{
		data: make(map[string]interface{}),
	}
}

// ConfigurationSource interface implementation

// LoadConfig implements the config.ConfigurationSource interface.
func (m *MockConfigurationSource) LoadConfig() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.loadConfigCalls++

	if m.shouldFail {
		if m.failureError != "" {
			return &ConfigError{Code: m.failureError}
		}
		return &ConfigError{Code: "mock.load_config_failed"}
	}

	if m.loadConfigFunc != nil {
		return m.loadConfigFunc()
	}

	return nil
}

// GetValue implements the config.ConfigurationSource interface.
func (m *MockConfigurationSource) GetValue(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.getValueCalls++

	if m.getValueFunc != nil {
		return m.getValueFunc(key)
	}

	value, exists := m.data[key]
	return value, exists
}

// Mock configuration methods

// SetValue sets a configuration value.
func (m *MockConfigurationSource) SetValue(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

// SetValues sets multiple configuration values.
func (m *MockConfigurationSource) SetValues(values map[string]interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, v := range values {
		m.data[k] = v
	}
}

// RemoveValue removes a configuration value.
func (m *MockConfigurationSource) RemoveValue(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
}

// SetShouldFail configures whether operations should fail.
func (m *MockConfigurationSource) SetShouldFail(fail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail = fail
}

// SetFailureError sets the error to return when operations fail.
func (m *MockConfigurationSource) SetFailureError(err string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failureError = err
}

// SetLoadConfigFunc sets a custom function for LoadConfig operations.
func (m *MockConfigurationSource) SetLoadConfigFunc(fn func() error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.loadConfigFunc = fn
}

// SetGetValueFunc sets a custom function for GetValue operations.
func (m *MockConfigurationSource) SetGetValueFunc(fn func(key string) (interface{}, bool)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.getValueFunc = fn
}

// Verification methods

// LoadConfigCallCount returns the number of LoadConfig calls.
func (m *MockConfigurationSource) LoadConfigCallCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.loadConfigCalls
}

// GetValueCallCount returns the number of GetValue calls.
func (m *MockConfigurationSource) GetValueCallCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.getValueCalls
}

// GetData returns a copy of the internal data.
func (m *MockConfigurationSource) GetData() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]interface{})
	for k, v := range m.data {
		result[k] = v
	}
	return result
}

// Reset clears all data and resets call counters.
func (m *MockConfigurationSource) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[string]interface{})
	m.shouldFail = false
	m.failureError = ""
	m.loadConfigCalls = 0
	m.getValueCalls = 0
	m.loadConfigFunc = nil
	m.getValueFunc = nil
}

// ConfigError implements error interface for configuration-specific errors.
type ConfigError struct {
	Code    string
	Message string
}

// Error implements the error interface.
func (e *ConfigError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Code
}

// MockConfigurationSourceBuilder provides a fluent interface for configuring MockConfigurationSource.
type MockConfigurationSourceBuilder struct {
	source *MockConfigurationSource
}

// NewMockConfigurationSourceBuilder creates a new builder for MockConfigurationSource.
func NewMockConfigurationSourceBuilder() *MockConfigurationSourceBuilder {
	return &MockConfigurationSourceBuilder{
		source: NewMockConfigurationSource(),
	}
}

// WithValue sets a configuration value.
func (b *MockConfigurationSourceBuilder) WithValue(key string, value interface{}) *MockConfigurationSourceBuilder {
	b.source.SetValue(key, value)
	return b
}

// WithValues sets multiple configuration values.
func (b *MockConfigurationSourceBuilder) WithValues(values map[string]interface{}) *MockConfigurationSourceBuilder {
	b.source.SetValues(values)
	return b
}

// WithFailure configures the source to fail operations.
func (b *MockConfigurationSourceBuilder) WithFailure(fail bool) *MockConfigurationSourceBuilder {
	b.source.SetShouldFail(fail)
	return b
}

// WithFailureError sets the error to return when operations fail.
func (b *MockConfigurationSourceBuilder) WithFailureError(err string) *MockConfigurationSourceBuilder {
	b.source.SetFailureError(err)
	return b
}

// Build returns the configured MockConfigurationSource.
func (b *MockConfigurationSourceBuilder) Build() config.ConfigurationSource {
	return b.source
}
