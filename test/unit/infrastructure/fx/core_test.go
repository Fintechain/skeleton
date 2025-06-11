package fx

import (
	"testing"

	"github.com/stretchr/testify/assert"

	internalFX "github.com/fintechain/skeleton/internal/fx"
	"github.com/fintechain/skeleton/test/unit/mocks"
)

func TestCoreModule_InterfaceCompliance(t *testing.T) {
	// Verify that CoreModule is a valid FX module
	assert.NotNil(t, internalFX.CoreModule, "CoreModule should not be nil")
}

func TestNewConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		params      internalFX.ConfigurationParams
		expectError bool
		description string
	}{
		{
			name:        "valid configuration creation",
			params:      internalFX.ConfigurationParams{},
			expectError: false,
			description: "Should create configuration with default sources",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internalFX.NewConfiguration(tt.params)

			assert.NotNil(t, result.Configuration, "Configuration should not be nil")

			// Test basic configuration operations
			exists := result.Configuration.Exists("nonexistent.key")
			assert.False(t, exists, "Non-existent key should return false")
		})
	}
}

func TestNewRegistry(t *testing.T) {
	tests := []struct {
		name        string
		params      internalFX.RegistryParams
		expectError bool
		description string
	}{
		{
			name:        "valid registry creation",
			params:      internalFX.RegistryParams{},
			expectError: false,
			description: "Should create registry successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internalFX.NewRegistry(tt.params)

			assert.NotNil(t, result.Registry, "Registry should not be nil")

			// Test basic registry operations
			count := result.Registry.Count()
			assert.Equal(t, 0, count, "New registry should be empty")
		})
	}
}

func TestNewEventBus(t *testing.T) {
	tests := []struct {
		name        string
		params      internalFX.EventBusParams
		expectError bool
		description string
	}{
		{
			name:        "valid event bus creation",
			params:      internalFX.EventBusParams{},
			expectError: false,
			description: "Should create event bus successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internalFX.NewEventBus(tt.params)

			assert.NotNil(t, result.EventBus, "EventBus should not be nil")
			assert.Equal(t, "event_bus", string(result.EventBus.ID()), "EventBus should have correct ID")
		})
	}
}

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name        string
		params      internalFX.LoggerParams
		expectError bool
		description string
	}{
		{
			name:        "valid logger creation",
			params:      internalFX.LoggerParams{},
			expectError: false,
			description: "Should create logger successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internalFX.NewLogger(tt.params)

			assert.NotNil(t, result.Logger, "Logger should not be nil")
			assert.Equal(t, "logger", string(result.Logger.ID()), "Logger should have correct ID")
		})
	}
}

func TestNewPluginManager(t *testing.T) {
	tests := []struct {
		name        string
		params      internalFX.PluginManagerParams
		expectError bool
		description string
	}{
		{
			name:        "valid plugin manager creation",
			params:      internalFX.PluginManagerParams{},
			expectError: false,
			description: "Should create plugin manager successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := internalFX.NewPluginManager(tt.params)

			assert.NotNil(t, result.PluginManager, "PluginManager should not be nil")
			assert.Equal(t, "plugin_manager", string(result.PluginManager.ID()), "PluginManager should have correct ID")
		})
	}
}

func TestNewRuntimeEnvironment(t *testing.T) {
	factory := mocks.NewFactory()

	tests := []struct {
		name        string
		params      internalFX.RuntimeEnvironmentParams
		expectError bool
		description string
	}{
		{
			name: "valid runtime environment creation",
			params: internalFX.RuntimeEnvironmentParams{
				Registry:      factory.RegistryInterface(),
				Configuration: factory.ConfigurationInterface(),
				PluginManager: factory.PluginManagerInterface(),
				EventBus:      factory.EventBusServiceInterface(),
				Logger:        factory.LoggerServiceInterface(),
			},
			expectError: false,
			description: "Should create runtime environment with valid dependencies",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := internalFX.NewRuntimeEnvironment(tt.params)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result.Runtime)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result.Runtime, "Runtime should not be nil")
			}
		})
	}
}
