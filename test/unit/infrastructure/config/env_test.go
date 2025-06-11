package config_test

import (
	"testing"

	"github.com/fintechain/skeleton/internal/domain/config"
	infraconfig "github.com/fintechain/skeleton/internal/infrastructure/config"
	"github.com/stretchr/testify/assert"
)

func TestEnvSource_InterfaceCompliance(t *testing.T) {
	// Verify interface compliance
	var _ config.ConfigurationSource = (*infraconfig.EnvSource)(nil)
}

func TestEnvSource_Constructor(t *testing.T) {
	tests := []struct {
		name        string
		prefix      string
		description string
	}{
		{
			name:        "with prefix",
			prefix:      "APP_",
			description: "Should create source with prefix",
		},
		{
			name:        "without prefix",
			prefix:      "",
			description: "Should create source without prefix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source := infraconfig.NewEnvSource(tt.prefix)
			assert.NotNil(t, source, "Constructor should return non-nil source")
		})
	}
}

func TestEnvSource_LoadConfig(t *testing.T) {
	tests := []struct {
		name         string
		setupEnv     map[string]string
		prefix       string
		expectedKeys []string
		expectError  bool
		description  string
	}{
		{
			name: "with prefix",
			setupEnv: map[string]string{
				"APP_DATABASE_HOST": "localhost",
				"APP_DATABASE_PORT": "5432",
				"RANDOM_VAR":        "value",
			},
			prefix: "APP_",
			expectedKeys: []string{
				"database.host",
				"database.port",
			},
			expectError: false,
			description: "Should load only prefixed variables",
		},
		{
			name: "without prefix",
			setupEnv: map[string]string{
				"APP_DATABASE_HOST": "localhost",
				"RANDOM_VAR":        "value",
			},
			prefix: "",
			expectedKeys: []string{
				"app.database.host",
				"random.var",
			},
			expectError: false,
			description: "Should load all variables",
		},
		{
			name:         "empty environment",
			setupEnv:     map[string]string{},
			prefix:       "",
			expectedKeys: []string{},
			expectError:  false,
			description:  "Should handle empty environment",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			for k, v := range tt.setupEnv {
				t.Setenv(k, v)
			}

			source := infraconfig.NewEnvSource(tt.prefix)
			err := source.LoadConfig()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Verify expected keys are present
				for _, key := range tt.expectedKeys {
					value, ok := source.GetValue(key)
					assert.True(t, ok, "Expected key %s not found", key)
					assert.NotNil(t, value)
				}
			}
		})
	}
}

func TestEnvSource_GetValue(t *testing.T) {
	tests := []struct {
		name        string
		setupEnv    map[string]string
		prefix      string
		key         string
		expectFound bool
		expectValue string
		description string
	}{
		{
			name: "existing key",
			setupEnv: map[string]string{
				"APP_DATABASE_HOST": "localhost",
			},
			prefix:      "APP_",
			key:         "database.host",
			expectFound: true,
			expectValue: "localhost",
			description: "Should return value for existing key",
		},
		{
			name: "case insensitive key",
			setupEnv: map[string]string{
				"APP_DATABASE_HOST": "localhost",
			},
			prefix:      "APP_",
			key:         "DATABASE.HOST",
			expectFound: true,
			expectValue: "localhost",
			description: "Should handle case-insensitive keys",
		},
		{
			name: "non-existent key",
			setupEnv: map[string]string{
				"APP_DATABASE_HOST": "localhost",
			},
			prefix:      "APP_",
			key:         "non.existent",
			expectFound: false,
			expectValue: "",
			description: "Should handle non-existent keys",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			for k, v := range tt.setupEnv {
				t.Setenv(k, v)
			}

			source := infraconfig.NewEnvSource(tt.prefix)
			err := source.LoadConfig()
			assert.NoError(t, err)

			value, ok := source.GetValue(tt.key)
			assert.Equal(t, tt.expectFound, ok)
			if tt.expectFound {
				assert.Equal(t, tt.expectValue, value)
			}
		})
	}
}

func TestEnvSource_ThreadSafety(t *testing.T) {
	// Setup test environment
	t.Setenv("APP_TEST_KEY", "test_value")

	source := infraconfig.NewEnvSource("APP_")
	err := source.LoadConfig()
	assert.NoError(t, err)

	// Test concurrent reads
	concurrentReads := 100
	done := make(chan bool)

	for i := 0; i < concurrentReads; i++ {
		go func() {
			value, ok := source.GetValue("test.key")
			assert.True(t, ok)
			assert.Equal(t, "test_value", value)
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < concurrentReads; i++ {
		<-done
	}

	// Test concurrent load and read
	for i := 0; i < 10; i++ {
		go func() {
			err := source.LoadConfig()
			assert.NoError(t, err)
			done <- true
		}()
		go func() {
			value, ok := source.GetValue("test.key")
			assert.True(t, ok)
			assert.Equal(t, "test_value", value)
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 20; i++ {
		<-done
	}
}
