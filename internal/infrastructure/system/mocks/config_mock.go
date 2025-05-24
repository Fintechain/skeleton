package mocks

import (
	"time"

	"github.com/ebanfa/skeleton/internal/infrastructure/config"
)

// MockConfiguration is a mock implementation of config.Configuration for testing
type MockConfiguration struct {
	// Function fields for customizing behavior
	GetStringFunc          func(string) string
	GetStringDefaultFunc   func(string, string) string
	GetIntFunc             func(string) (int, error)
	GetIntDefaultFunc      func(string, int) int
	GetBoolFunc            func(string) (bool, error)
	GetBoolDefaultFunc     func(string, bool) bool
	GetDurationFunc        func(string) (time.Duration, error)
	GetDurationDefaultFunc func(string, time.Duration) time.Duration
	GetObjectFunc          func(string, interface{}) error
	ExistsFunc             func(string) bool

	// Call tracking
	GetStringCalls          []string
	GetStringDefaultCalls   []GetStringDefaultCall
	GetIntCalls             []string
	GetIntDefaultCalls      []GetIntDefaultCall
	GetBoolCalls            []string
	GetBoolDefaultCalls     []GetBoolDefaultCall
	GetDurationCalls        []string
	GetDurationDefaultCalls []GetDurationDefaultCall
	GetObjectCalls          []GetObjectCall
	ExistsCalls             []string

	// State
	Values map[string]interface{}
}

type GetStringDefaultCall struct {
	Key          string
	DefaultValue string
}

type GetIntDefaultCall struct {
	Key          string
	DefaultValue int
}

type GetBoolDefaultCall struct {
	Key          string
	DefaultValue bool
}

type GetDurationDefaultCall struct {
	Key          string
	DefaultValue time.Duration
}

type GetObjectCall struct {
	Key    string
	Result interface{}
}

// NewMockConfiguration creates a new mock configuration
func NewMockConfiguration() *MockConfiguration {
	return &MockConfiguration{
		Values: make(map[string]interface{}),
	}
}

// GetString implements config.Configuration
func (m *MockConfiguration) GetString(key string) string {
	m.GetStringCalls = append(m.GetStringCalls, key)
	if m.GetStringFunc != nil {
		return m.GetStringFunc(key)
	}
	if value, exists := m.Values[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return ""
}

// GetStringDefault implements config.Configuration
func (m *MockConfiguration) GetStringDefault(key, defaultValue string) string {
	m.GetStringDefaultCalls = append(m.GetStringDefaultCalls, GetStringDefaultCall{Key: key, DefaultValue: defaultValue})
	if m.GetStringDefaultFunc != nil {
		return m.GetStringDefaultFunc(key, defaultValue)
	}
	if value, exists := m.Values[key]; exists {
		if str, ok := value.(string); ok {
			return str
		}
	}
	return defaultValue
}

// GetInt implements config.Configuration
func (m *MockConfiguration) GetInt(key string) (int, error) {
	m.GetIntCalls = append(m.GetIntCalls, key)
	if m.GetIntFunc != nil {
		return m.GetIntFunc(key)
	}
	if value, exists := m.Values[key]; exists {
		if i, ok := value.(int); ok {
			return i, nil
		}
	}
	return 0, &config.ConfigError{Code: config.ErrConfigNotFound, Message: "key not found"}
}

// GetIntDefault implements config.Configuration
func (m *MockConfiguration) GetIntDefault(key string, defaultValue int) int {
	m.GetIntDefaultCalls = append(m.GetIntDefaultCalls, GetIntDefaultCall{Key: key, DefaultValue: defaultValue})
	if m.GetIntDefaultFunc != nil {
		return m.GetIntDefaultFunc(key, defaultValue)
	}
	if value, exists := m.Values[key]; exists {
		if i, ok := value.(int); ok {
			return i
		}
	}
	return defaultValue
}

// GetBool implements config.Configuration
func (m *MockConfiguration) GetBool(key string) (bool, error) {
	m.GetBoolCalls = append(m.GetBoolCalls, key)
	if m.GetBoolFunc != nil {
		return m.GetBoolFunc(key)
	}
	if value, exists := m.Values[key]; exists {
		if b, ok := value.(bool); ok {
			return b, nil
		}
	}
	return false, &config.ConfigError{Code: config.ErrConfigNotFound, Message: "key not found"}
}

// GetBoolDefault implements config.Configuration
func (m *MockConfiguration) GetBoolDefault(key string, defaultValue bool) bool {
	m.GetBoolDefaultCalls = append(m.GetBoolDefaultCalls, GetBoolDefaultCall{Key: key, DefaultValue: defaultValue})
	if m.GetBoolDefaultFunc != nil {
		return m.GetBoolDefaultFunc(key, defaultValue)
	}
	if value, exists := m.Values[key]; exists {
		if b, ok := value.(bool); ok {
			return b
		}
	}
	return defaultValue
}

// GetDuration implements config.Configuration
func (m *MockConfiguration) GetDuration(key string) (time.Duration, error) {
	m.GetDurationCalls = append(m.GetDurationCalls, key)
	if m.GetDurationFunc != nil {
		return m.GetDurationFunc(key)
	}
	if value, exists := m.Values[key]; exists {
		if d, ok := value.(time.Duration); ok {
			return d, nil
		}
	}
	return 0, &config.ConfigError{Code: config.ErrConfigNotFound, Message: "key not found"}
}

// GetDurationDefault implements config.Configuration
func (m *MockConfiguration) GetDurationDefault(key string, defaultValue time.Duration) time.Duration {
	m.GetDurationDefaultCalls = append(m.GetDurationDefaultCalls, GetDurationDefaultCall{Key: key, DefaultValue: defaultValue})
	if m.GetDurationDefaultFunc != nil {
		return m.GetDurationDefaultFunc(key, defaultValue)
	}
	if value, exists := m.Values[key]; exists {
		if d, ok := value.(time.Duration); ok {
			return d
		}
	}
	return defaultValue
}

// GetObject implements config.Configuration
func (m *MockConfiguration) GetObject(key string, result interface{}) error {
	m.GetObjectCalls = append(m.GetObjectCalls, GetObjectCall{Key: key, Result: result})
	if m.GetObjectFunc != nil {
		return m.GetObjectFunc(key, result)
	}
	return &config.ConfigError{Code: config.ErrConfigNotFound, Message: "key not found"}
}

// Exists implements config.Configuration
func (m *MockConfiguration) Exists(key string) bool {
	m.ExistsCalls = append(m.ExistsCalls, key)
	if m.ExistsFunc != nil {
		return m.ExistsFunc(key)
	}
	_, exists := m.Values[key]
	return exists
}

// Set is a helper method for testing (not part of the interface)
func (m *MockConfiguration) Set(key string, value interface{}) {
	if m.Values == nil {
		m.Values = make(map[string]interface{})
	}
	m.Values[key] = value
}
