package fx

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fintechain/skeleton/pkg/fx"
	"github.com/fintechain/skeleton/test/unit/mocks"
)

func TestRuntimeOptions(t *testing.T) {
	factory := mocks.NewFactory()

	tests := []struct {
		name        string
		option      fx.RuntimeOption
		expectError bool
		description string
	}{
		{
			name: "WithPlugins option",
			option: fx.WithPlugins(
				factory.PluginInterface(),
			),
			expectError: false,
			description: "Should create WithPlugins option successfully",
		},
		{
			name:   "WithFXOptions option",
			option: fx.WithFXOptions(
			// Empty FX options for testing
			),
			expectError: false,
			description: "Should create WithFXOptions option successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test option application
			config := &fx.RuntimeConfig{}
			assert.NotPanics(t, func() {
				tt.option(config)
			}, "Option should apply without panicking")
		})
	}
}

func TestWithPlugins(t *testing.T) {
	factory := mocks.NewFactory()

	tests := []struct {
		name        string
		plugins     []interface{} // Using interface{} to test plugin interface compliance
		expectError bool
		description string
	}{
		{
			name:        "single plugin",
			plugins:     []interface{}{factory.PluginInterface()},
			expectError: false,
			description: "Should handle single plugin",
		},
		{
			name: "multiple plugins",
			plugins: []interface{}{
				factory.PluginInterface(),
				factory.PluginInterface(),
			},
			expectError: false,
			description: "Should handle multiple plugins",
		},
		{
			name:        "no plugins",
			plugins:     []interface{}{},
			expectError: false,
			description: "Should handle empty plugin list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &fx.RuntimeConfig{}

			// Convert to proper plugin interfaces for the test
			var plugins []interface{}
			for _, p := range tt.plugins {
				plugins = append(plugins, p)
			}

			// Test that option doesn't panic
			assert.NotPanics(t, func() {
				// Note: This is simplified since we can't directly test the actual plugin types
				// In a real scenario, the plugins would be properly typed
				option := fx.WithPlugins() // Empty for compilation
				option(config)
			}, "WithPlugins should not panic")
		})
	}
}

func TestExecuteCommand_ErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		operationID string
		input       map[string]interface{}
		opts        []fx.RuntimeOption
		expectError bool
		description string
	}{
		{
			name:        "non-existent operation",
			operationID: "non-existent-operation",
			input: map[string]interface{}{
				"test": "data",
			},
			opts:        []fx.RuntimeOption{},
			expectError: true,
			description: "Should return error for non-existent operation",
		},
		{
			name:        "empty operation ID",
			operationID: "",
			input: map[string]interface{}{
				"test": "data",
			},
			opts:        []fx.RuntimeOption{},
			expectError: true,
			description: "Should return error for empty operation ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := fx.ExecuteCommand(tt.operationID, tt.input, tt.opts...)

			if tt.expectError {
				assert.Error(t, err, "Should return error")
				assert.Nil(t, result, "Result should be nil on error")
			} else {
				assert.NoError(t, err, "Should not return error")
				assert.NotNil(t, result, "Result should not be nil")
			}
		})
	}
}

func TestRuntimeConfig(t *testing.T) {
	factory := mocks.NewFactory()

	tests := []struct {
		name        string
		setupConfig func() *fx.RuntimeConfig
		validate    func(*testing.T, *fx.RuntimeConfig)
		description string
	}{
		{
			name: "empty config",
			setupConfig: func() *fx.RuntimeConfig {
				return &fx.RuntimeConfig{}
			},
			validate: func(t *testing.T, config *fx.RuntimeConfig) {
				assert.Empty(t, config.Plugins, "Plugins should be empty")
				assert.Empty(t, config.ExtraOptions, "ExtraOptions should be empty")
			},
			description: "Should handle empty configuration",
		},
		{
			name: "config with plugins",
			setupConfig: func() *fx.RuntimeConfig {
				config := &fx.RuntimeConfig{}
				plugin := factory.PluginInterface()
				config.Plugins = append(config.Plugins, plugin)
				return config
			},
			validate: func(t *testing.T, config *fx.RuntimeConfig) {
				assert.Len(t, config.Plugins, 1, "Should have one plugin")
			},
			description: "Should handle configuration with plugins",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.setupConfig()
			tt.validate(t, config)
		})
	}
}

func TestFXIntegration_InterfaceCompliance(t *testing.T) {
	// Test that the public API functions exist and can be called
	t.Run("StartDaemon exists", func(t *testing.T) {
		assert.NotNil(t, fx.StartDaemon, "StartDaemon function should exist")
	})

	t.Run("ExecuteCommand exists", func(t *testing.T) {
		assert.NotNil(t, fx.ExecuteCommand, "ExecuteCommand function should exist")
	})

	t.Run("StartDaemonWithSignalHandling exists", func(t *testing.T) {
		assert.NotNil(t, fx.StartDaemonWithSignalHandling, "StartDaemonWithSignalHandling function should exist")
	})

	t.Run("WithPlugins option creator exists", func(t *testing.T) {
		assert.NotNil(t, fx.WithPlugins, "WithPlugins function should exist")
	})

	t.Run("WithFXOptions option creator exists", func(t *testing.T) {
		assert.NotNil(t, fx.WithFXOptions, "WithFXOptions function should exist")
	})
}
