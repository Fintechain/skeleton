package config

import (
	"testing"
	"time"
)

// TestDefaultConfigCreate tests the creation of a DefaultConfig
func TestDefaultConfigCreate(t *testing.T) {
	// Test with options
	options := DefaultConfigOptions{
		Name: "test-config",
	}
	config := NewDefaultConfig(options)

	if config == nil {
		t.Fatal("NewDefaultConfig returned nil")
	}

	if config.Name != "test-config" {
		t.Errorf("Expected name 'test-config', got '%s'", config.Name)
	}

	if config.Properties == nil {
		t.Error("Properties map should be initialized")
	}

	// Test factory method
	defaultConfig := CreateDefaultConfig()
	if defaultConfig == nil {
		t.Fatal("CreateDefaultConfig returned nil")
	}

	if defaultConfig.Name != "default" {
		t.Errorf("Expected name 'default', got '%s'", defaultConfig.Name)
	}

	// Test with source
	mockSource := &MockConfigSource{values: map[string]interface{}{}}
	sourceConfig := CreateConfigWithSource("source-config", mockSource)

	if sourceConfig == nil {
		t.Fatal("CreateConfigWithSource returned nil")
	}

	if sourceConfig.Name != "source-config" {
		t.Errorf("Expected name 'source-config', got '%s'", sourceConfig.Name)
	}

	if sourceConfig.Source != mockSource {
		t.Error("Source was not correctly set")
	}
}

// TestDefaultConfigSetAndGet tests setting and getting config values
func TestDefaultConfigSetAndGet(t *testing.T) {
	config := CreateDefaultConfig()

	// Set and get a string
	config.Set("string-key", "string-value")
	value := config.Get("string-key")

	if value != "string-value" {
		t.Errorf("Expected 'string-value', got '%v'", value)
	}

	// Set and get an int
	config.Set("int-key", 42)
	value = config.Get("int-key")

	if value != 42 {
		t.Errorf("Expected 42, got '%v'", value)
	}

	// Get a non-existent key
	value = config.Get("non-existent")
	if value != nil {
		t.Errorf("Expected nil for non-existent key, got '%v'", value)
	}
}

// TestDefaultConfigGetFromSource tests getting values from a source
func TestDefaultConfigGetFromSource(t *testing.T) {
	mockSource := &MockConfigSource{
		values: map[string]interface{}{
			"source-key": "source-value",
		},
	}

	config := CreateConfigWithSource("test", mockSource)

	// Get from source
	value := config.Get("source-key")
	if value != "source-value" {
		t.Errorf("Expected 'source-value' from source, got '%v'", value)
	}

	// Properties should take precedence over source
	config.Set("source-key", "properties-value")
	value = config.Get("source-key")

	if value != "properties-value" {
		t.Errorf("Expected 'properties-value' from properties, got '%v'", value)
	}
}

// TestDefaultConfigExists tests checking if a key exists
func TestDefaultConfigExists(t *testing.T) {
	mockSource := &MockConfigSource{
		values: map[string]interface{}{
			"source-key": "source-value",
		},
	}

	config := CreateConfigWithSource("test", mockSource)
	config.Set("properties-key", "properties-value")

	// Check properties key
	if !config.Exists("properties-key") {
		t.Error("Expected properties-key to exist")
	}

	// Check source key
	if !config.Exists("source-key") {
		t.Error("Expected source-key to exist")
	}

	// Check non-existent key
	if config.Exists("non-existent") {
		t.Error("Expected non-existent key to not exist")
	}
}

// TestDefaultConfigGetString tests getting string values
func TestDefaultConfigGetString(t *testing.T) {
	config := CreateDefaultConfig()
	config.Set("string-key", "string-value")
	config.Set("int-key", 42)

	// Get existing string
	value := config.GetString("string-key")
	if value != "string-value" {
		t.Errorf("Expected 'string-value', got '%s'", value)
	}

	// Get non-string value (should return default)
	value = config.GetString("int-key")
	if value != "" {
		t.Errorf("Expected empty string for non-string value, got '%s'", value)
	}

	// Get non-existent key (should return default)
	value = config.GetString("non-existent")
	if value != "" {
		t.Errorf("Expected empty string for non-existent key, got '%s'", value)
	}

	// Get string with default
	value = config.GetStringDefault("non-existent", "default-value")
	if value != "default-value" {
		t.Errorf("Expected 'default-value', got '%s'", value)
	}
}

// TestDefaultConfigGetInt tests getting integer values
func TestDefaultConfigGetInt(t *testing.T) {
	config := CreateDefaultConfig()
	config.Set("int-key", 42)
	config.Set("string-key", "not-an-int")

	// Get existing int
	value, err := config.GetInt("int-key")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if value != 42 {
		t.Errorf("Expected 42, got %d", value)
	}

	// Get non-int value (should return error)
	_, err = config.GetInt("string-key")
	if err == nil {
		t.Error("Expected error for non-int value")
	}

	// Get non-existent key (should return error)
	_, err = config.GetInt("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent key")
	}

	// Get int with default
	value = config.GetIntDefault("non-existent", 99)
	if value != 99 {
		t.Errorf("Expected 99, got %d", value)
	}
}

// TestDefaultConfigGetBool tests getting boolean values
func TestDefaultConfigGetBool(t *testing.T) {
	config := CreateDefaultConfig()
	config.Set("bool-key-true", true)
	config.Set("bool-key-false", false)
	config.Set("string-key", "not-a-bool")

	// Get existing bool (true)
	value, err := config.GetBool("bool-key-true")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !value {
		t.Error("Expected true, got false")
	}

	// Get existing bool (false)
	value, err = config.GetBool("bool-key-false")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if value {
		t.Error("Expected false, got true")
	}

	// Get non-bool value (should return error)
	_, err = config.GetBool("string-key")
	if err == nil {
		t.Error("Expected error for non-bool value")
	}

	// Get bool with default
	value = config.GetBoolDefault("non-existent", true)
	if !value {
		t.Error("Expected true default, got false")
	}
}

// TestDefaultConfigGetDuration tests getting duration values
func TestDefaultConfigGetDuration(t *testing.T) {
	config := CreateDefaultConfig()

	duration := 5 * time.Second
	config.Set("duration-key", duration)
	config.Set("string-duration-key", "10s")
	config.Set("int-duration-key", 15000) // 15 seconds in milliseconds
	config.Set("invalid-string-key", "not-a-duration")

	// Get existing duration
	value, err := config.GetDuration("duration-key")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if value != duration {
		t.Errorf("Expected %v, got %v", duration, value)
	}

	// Get duration from string
	value, err = config.GetDuration("string-duration-key")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if value != 10*time.Second {
		t.Errorf("Expected 10s, got %v", value)
	}

	// Get duration from int
	value, err = config.GetDuration("int-duration-key")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if value != 15*time.Second {
		t.Errorf("Expected 15s, got %v", value)
	}

	// Get invalid string duration (should return error)
	_, err = config.GetDuration("invalid-string-key")
	if err == nil {
		t.Error("Expected error for invalid string duration")
	}

	// Get duration with default
	value = config.GetDurationDefault("non-existent", 20*time.Second)
	if value != 20*time.Second {
		t.Errorf("Expected 20s, got %v", value)
	}
}

// TestConfigError tests the ConfigError type
func TestConfigError(t *testing.T) {
	err := &ConfigError{
		Code:    "test.error",
		Message: "Test error message",
		Details: map[string]interface{}{"key": "value"},
	}

	errStr := err.Error()
	if errStr != "Test error message" {
		t.Errorf("Expected 'Test error message', got '%s'", errStr)
	}
}

// MockConfigSource is a simple implementation of ConfigurationSource for testing
type MockConfigSource struct {
	values     map[string]interface{}
	loadError  error
	loadCalled bool
}

func (m *MockConfigSource) LoadConfig() error {
	m.loadCalled = true
	return m.loadError
}

func (m *MockConfigSource) GetValue(key string) (interface{}, bool) {
	value, exists := m.values[key]
	return value, exists
}
