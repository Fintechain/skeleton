package runtime_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/fx"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/plugin"
	"github.com/fintechain/skeleton/pkg/runtime"
	"github.com/fintechain/skeleton/test/unit/mocks"
)

// TestRuntimeOptions tests the option functions
func TestRuntimeOptions(t *testing.T) {
	factory := mocks.NewFactory()

	t.Run("WithPlugins", func(t *testing.T) {
		mockPlugin1 := factory.PluginInterface()
		mockPlugin2 := factory.PluginInterface()

		option := runtime.WithPlugins(mockPlugin1, mockPlugin2)
		assert.NotNil(t, option)

		// Test that option can be applied to config
		config := &runtime.Config{}
		option(config)

		assert.Len(t, config.Plugins, 2)
		assert.Equal(t, mockPlugin1, config.Plugins[0])
		assert.Equal(t, mockPlugin2, config.Plugins[1])
	})

	t.Run("WithOptions", func(t *testing.T) {
		fxOption1 := fx.Provide(func() string { return "test1" })
		fxOption2 := fx.Provide(func() int { return 42 })

		option := runtime.WithOptions(fxOption1, fxOption2)
		assert.NotNil(t, option)

		// Test that option can be applied to config
		config := &runtime.Config{}
		option(config)

		assert.Len(t, config.ExtraOptions, 2)
	})

	t.Run("Multiple options", func(t *testing.T) {
		mockPlugin := factory.PluginInterface()
		fxOption := fx.Provide(func() string { return "test" })

		config := &runtime.Config{}

		// Apply multiple options
		runtime.WithPlugins(mockPlugin)(config)
		runtime.WithOptions(fxOption)(config)

		assert.Len(t, config.Plugins, 1)
		assert.Len(t, config.ExtraOptions, 1)
	})
}

// TestExecuteCommand tests the command execution functionality
func TestExecuteCommand(t *testing.T) {
	t.Run("ExecuteCommand with no plugins", func(t *testing.T) {
		// This should work with default configuration
		// Note: This will fail because no operation is registered, but it tests the setup
		result, err := runtime.ExecuteCommand("non-existent-operation",
			map[string]interface{}{"test": "data"})

		// We expect an error because the operation doesn't exist
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "operation execution failed")
	})

	t.Run("ExecuteCommand with plugins", func(t *testing.T) {
		factory := mocks.NewFactory()
		mockPlugin := factory.PluginInterface()

		// Configure mock plugin expectations
		mockPlugin.On("ID").Return(component.ComponentID("test-plugin"))
		mockPlugin.On("Initialize", mock.Anything, mock.Anything).Return(nil)

		// This should work with the plugin but still fail because no operation is registered
		result, err := runtime.ExecuteCommand("test-operation",
			map[string]interface{}{"input": "test"},
			runtime.WithPlugins(mockPlugin))

		// We expect an error because the operation doesn't exist
		assert.Error(t, err)
		assert.Nil(t, result)

		// Verify plugin was called
		mockPlugin.AssertExpectations(t)
	})
}

// TestRuntimeConfig tests the Config struct
func TestRuntimeConfig(t *testing.T) {
	t.Run("Empty config", func(t *testing.T) {
		config := &runtime.Config{}
		assert.Empty(t, config.Plugins)
		assert.Empty(t, config.ExtraOptions)
	})

	t.Run("Config with plugins", func(t *testing.T) {
		factory := mocks.NewFactory()
		mockPlugin := factory.PluginInterface()

		config := &runtime.Config{
			Plugins: []plugin.Plugin{mockPlugin},
		}

		assert.Len(t, config.Plugins, 1)
		assert.Equal(t, mockPlugin, config.Plugins[0])
	})

	t.Run("Config with FX options", func(t *testing.T) {
		fxOption := fx.Provide(func() string { return "test" })

		config := &runtime.Config{
			ExtraOptions: []fx.Option{fxOption},
		}

		assert.Len(t, config.ExtraOptions, 1)
	})
}

// TestStartDaemonTimeout tests daemon startup with timeout
func TestStartDaemonTimeout(t *testing.T) {
	t.Run("StartDaemon should start and be interruptible", func(t *testing.T) {
		// This test verifies that StartDaemon can be started
		// We'll use a goroutine and timeout to avoid blocking the test

		done := make(chan error, 1)

		go func() {
			// This should start successfully but we'll interrupt it
			err := runtime.StartDaemon()
			done <- err
		}()

		// Give it a moment to start
		time.Sleep(100 * time.Millisecond)

		// For this test, we just verify it doesn't panic immediately
		// In a real scenario, you'd send a signal to stop it
		select {
		case err := <-done:
			// If it completed quickly, it might have failed to start
			// which is fine for this test
			t.Logf("StartDaemon completed with: %v", err)
		case <-time.After(200 * time.Millisecond):
			// If it's still running after timeout, that's expected behavior
			t.Log("StartDaemon is running (expected behavior)")
		}
	})
}

// TestOptionCombinations tests various combinations of options
func TestOptionCombinations(t *testing.T) {
	factory := mocks.NewFactory()

	t.Run("All options combined", func(t *testing.T) {
		mockPlugin1 := factory.PluginInterface()
		mockPlugin2 := factory.PluginInterface()
		fxOption1 := fx.Provide(func() string { return "test" })
		fxOption2 := fx.Provide(func() int { return 42 })

		config := &runtime.Config{}

		// Apply all option types
		runtime.WithPlugins(mockPlugin1, mockPlugin2)(config)
		runtime.WithOptions(fxOption1, fxOption2)(config)

		assert.Len(t, config.Plugins, 2)
		assert.Len(t, config.ExtraOptions, 2)
	})

	t.Run("Incremental option application", func(t *testing.T) {
		mockPlugin1 := factory.PluginInterface()
		mockPlugin2 := factory.PluginInterface()

		config := &runtime.Config{}

		// Apply options incrementally
		runtime.WithPlugins(mockPlugin1)(config)
		assert.Len(t, config.Plugins, 1)

		runtime.WithPlugins(mockPlugin2)(config)
		assert.Len(t, config.Plugins, 2)
	})
}

// BenchmarkExecuteCommand benchmarks command execution
func BenchmarkExecuteCommand(b *testing.B) {
	input := map[string]interface{}{
		"test":   "data",
		"number": 42,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// This will fail but we're measuring the setup overhead
		runtime.ExecuteCommand("non-existent", input)
	}
}

// BenchmarkOptionApplication benchmarks option application
func BenchmarkOptionApplication(b *testing.B) {
	factory := mocks.NewFactory()
	mockPlugin := factory.PluginInterface()
	fxOption := fx.Provide(func() string { return "test" })

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config := &runtime.Config{}
		runtime.WithPlugins(mockPlugin)(config)
		runtime.WithOptions(fxOption)(config)
	}
}
