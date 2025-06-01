package mocks

import (
	"fmt"
	"sync"
	"time"

	"github.com/fintechain/skeleton/pkg/config"
)

// MockConfiguration provides a configurable mock implementation of the config.Configuration interface.
type MockConfiguration struct {
	mu sync.RWMutex

	// Configuration data
	values map[string]interface{}

	// Behavior configuration
	shouldFail   bool
	failureError string

	// Call tracking
	callCount map[string]int
	lastCalls map[string][]interface{}
}

// NewMockConfiguration creates a new configurable configuration mock.
func NewMockConfiguration() *MockConfiguration {
	return &MockConfiguration{
		values:    make(map[string]interface{}),
		callCount: make(map[string]int),
		lastCalls: make(map[string][]interface{}),
	}
}

// Configuration Interface Implementation

// GetString retrieves a string configuration value.
func (m *MockConfiguration) GetString(key string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("GetString", key)

	if value, exists := m.values[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", value)
	}

	return ""
}

// GetStringDefault retrieves a string configuration value with a default value.
func (m *MockConfiguration) GetStringDefault(key, defaultValue string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("GetStringDefault", key, defaultValue)

	if value, exists := m.values[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
		return fmt.Sprintf("%v", value)
	}

	return defaultValue
}

// GetInt retrieves an integer configuration value.
func (m *MockConfiguration) GetInt(key string) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("GetInt", key)

	if m.shouldFail {
		return 0, fmt.Errorf("%s", m.getFailureError("GetInt"))
	}

	if value, exists := m.values[key]; exists {
		if i, ok := value.(int); ok {
			return i, nil
		}
		return 0, fmt.Errorf("value is not an integer")
	}

	return 0, fmt.Errorf("key not found")
}

// GetIntDefault retrieves an integer configuration value with a default value.
func (m *MockConfiguration) GetIntDefault(key string, defaultValue int) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("GetIntDefault", key, defaultValue)

	if value, exists := m.values[key]; exists {
		if i, ok := value.(int); ok {
			return i
		}
	}

	return defaultValue
}

// GetBool retrieves a boolean configuration value.
func (m *MockConfiguration) GetBool(key string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("GetBool", key)

	if m.shouldFail {
		return false, fmt.Errorf("%s", m.getFailureError("GetBool"))
	}

	if value, exists := m.values[key]; exists {
		if b, ok := value.(bool); ok {
			return b, nil
		}
		return false, fmt.Errorf("value is not a boolean")
	}

	return false, fmt.Errorf("key not found")
}

// GetBoolDefault retrieves a boolean configuration value with a default value.
func (m *MockConfiguration) GetBoolDefault(key string, defaultValue bool) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("GetBoolDefault", key, defaultValue)

	if value, exists := m.values[key]; exists {
		if b, ok := value.(bool); ok {
			return b
		}
	}

	return defaultValue
}

// GetDuration retrieves a duration configuration value.
func (m *MockConfiguration) GetDuration(key string) (time.Duration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("GetDuration", key)

	if m.shouldFail {
		return 0, fmt.Errorf("%s", m.getFailureError("GetDuration"))
	}

	if value, exists := m.values[key]; exists {
		if d, ok := value.(time.Duration); ok {
			return d, nil
		}
		return 0, fmt.Errorf("value is not a duration")
	}

	return 0, fmt.Errorf("key not found")
}

// GetDurationDefault retrieves a duration configuration value with a default value.
func (m *MockConfiguration) GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("GetDurationDefault", key, defaultValue)

	if value, exists := m.values[key]; exists {
		if d, ok := value.(time.Duration); ok {
			return d
		}
	}

	return defaultValue
}

// GetObject deserializes a configuration section into a struct.
func (m *MockConfiguration) GetObject(key string, result interface{}) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("GetObject", key, result)

	if m.shouldFail {
		return fmt.Errorf("%s", m.getFailureError("GetObject"))
	}

	// For mock purposes, we'll just return success
	// In a real implementation, this would deserialize the value
	return nil
}

// Exists checks if a configuration key exists.
func (m *MockConfiguration) Exists(key string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	m.trackCall("Exists", key)

	_, exists := m.values[key]
	return exists
}

// Mock Configuration Methods

// SetValue sets a configuration value directly.
func (m *MockConfiguration) SetValue(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.values[key] = value
}

// SetShouldFail configures the mock to fail operations.
func (m *MockConfiguration) SetShouldFail(fail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail = fail
}

// SetFailureError sets the error message for failed operations.
func (m *MockConfiguration) SetFailureError(err string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failureError = err
}

// State Verification Methods

// GetCallCount returns the number of times a method was called.
func (m *MockConfiguration) GetCallCount(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.callCount[method]
}

// GetLastCall returns the parameters of the last call to a method.
func (m *MockConfiguration) GetLastCall(method string) []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastCalls[method]
}

// Reset clears all mock state and configuration.
func (m *MockConfiguration) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.values = make(map[string]interface{})
	m.shouldFail = false
	m.failureError = ""
	m.callCount = make(map[string]int)
	m.lastCalls = make(map[string][]interface{})
}

// Helper Methods

// trackCall records a method call for verification.
func (m *MockConfiguration) trackCall(method string, args ...interface{}) {
	m.callCount[method]++
	m.lastCalls[method] = args
}

// getFailureError returns the configured failure error or a default.
func (m *MockConfiguration) getFailureError(method string) string {
	if m.failureError != "" {
		return m.failureError
	}
	return fmt.Sprintf("mock_configuration.%s_failed", method)
}

// ConfigurationMockBuilder provides a fluent interface for configuring configuration mocks.
type ConfigurationMockBuilder struct {
	mock *MockConfiguration
}

// NewConfigurationMockBuilder creates a new configuration mock builder.
func NewConfigurationMockBuilder() *ConfigurationMockBuilder {
	return &ConfigurationMockBuilder{
		mock: NewMockConfiguration(),
	}
}

// WithValue sets a configuration value.
func (b *ConfigurationMockBuilder) WithValue(key string, value interface{}) *ConfigurationMockBuilder {
	b.mock.SetValue(key, value)
	return b
}

// WithFailure configures the mock to fail operations.
func (b *ConfigurationMockBuilder) WithFailure(fail bool) *ConfigurationMockBuilder {
	b.mock.SetShouldFail(fail)
	return b
}

// WithFailureError sets the error message for failed operations.
func (b *ConfigurationMockBuilder) WithFailureError(err string) *ConfigurationMockBuilder {
	b.mock.SetFailureError(err)
	return b
}

// Build returns the configured mock configuration.
func (b *ConfigurationMockBuilder) Build() config.Configuration {
	return b.mock
}
