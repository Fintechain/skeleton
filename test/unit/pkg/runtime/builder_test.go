package runtime_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/pkg/runtime"
	"github.com/fintechain/skeleton/test/unit/mocks"
)

// TestRuntimeBuilder tests the new Builder API
func TestRuntimeBuilder(t *testing.T) {
	factory := mocks.NewFactory()

	t.Run("NewBuilder creates builder", func(t *testing.T) {
		builder := runtime.NewBuilder()
		assert.NotNil(t, builder)
	})

	t.Run("WithPlugins adds plugins", func(t *testing.T) {
		mockPlugin1 := factory.PluginInterface()
		mockPlugin2 := factory.PluginInterface()

		builder := runtime.NewBuilder().
			WithPlugins(mockPlugin1, mockPlugin2)

		assert.NotNil(t, builder)
		// Builder should be chainable
	})

	t.Run("BuildCommand with no plugins", func(t *testing.T) {
		// This should work with default configuration
		// Note: This will fail because no operation is registered, but it tests the setup
		result, err := runtime.NewBuilder().
			BuildCommand("non-existent-operation", map[string]interface{}{"test": "data"})

		// We expect an error because the operation doesn't exist
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "operation execution failed")
	})

	t.Run("BuildCommand with plugins", func(t *testing.T) {
		mockPlugin := factory.PluginInterface()

		// Configure mock plugin expectations
		mockPlugin.On("ID").Return(component.ComponentID("test-plugin"))
		mockPlugin.On("Initialize", mock.Anything, mock.Anything).Return(nil)

		// This should work with the plugin but still fail because no operation is registered
		result, err := runtime.NewBuilder().
			WithPlugins(mockPlugin).
			BuildCommand("test-operation", map[string]interface{}{"input": "test"})

		// We expect an error because the operation doesn't exist
		assert.Error(t, err)
		assert.Nil(t, result)

		// Verify plugin was called
		mockPlugin.AssertExpectations(t)
	})

	t.Run("Builder is chainable", func(t *testing.T) {
		mockPlugin := factory.PluginInterface()

		// Test method chaining
		builder := runtime.NewBuilder().
			WithPlugins(mockPlugin)

		assert.NotNil(t, builder)
	})
}

// TestBuilderWithCustomDependencies tests custom dependency injection
func TestBuilderWithCustomDependencies(t *testing.T) {
	t.Run("WithConfig sets custom configuration", func(t *testing.T) {
		factory := mocks.NewFactory()
		mockConfig := factory.ConfigurationInterface()

		builder := runtime.NewBuilder().
			WithConfig(mockConfig)

		assert.NotNil(t, builder)
		// Builder should accept custom configuration
	})

	t.Run("Builder method chaining", func(t *testing.T) {
		factory := mocks.NewFactory()
		mockPlugin := factory.PluginInterface()
		mockConfig := factory.ConfigurationInterface()

		// Configure mock plugin expectations
		mockPlugin.On("ID").Return(component.ComponentID("test-plugin"))
		mockPlugin.On("Initialize", mock.Anything, mock.Anything).Return(nil)

		builder := runtime.NewBuilder().
			WithPlugins(mockPlugin).
			WithConfig(mockConfig)

		assert.NotNil(t, builder)

		// Test that it can attempt to build (will fail due to no operation, but tests wiring)
		result, err := builder.BuildCommand("test-operation", map[string]interface{}{"input": "test"})

		// We expect an error because the operation doesn't exist
		assert.Error(t, err)
		assert.Nil(t, result)

		// Verify plugin was called
		mockPlugin.AssertExpectations(t)
	})
}

// BenchmarkBuilderAPI benchmarks the Builder API
func BenchmarkBuilderAPI(b *testing.B) {
	factory := mocks.NewFactory()
	mockPlugin := factory.PluginInterface()

	// Configure mock plugin expectations
	mockPlugin.On("ID").Return(component.ComponentID("test-plugin"))
	mockPlugin.On("Initialize", mock.Anything, mock.Anything).Return(nil)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Benchmark builder creation and command execution attempt
		runtime.NewBuilder().
			WithPlugins(mockPlugin).
			BuildCommand("non-existent", map[string]interface{}{"test": "data"})
	}
}
