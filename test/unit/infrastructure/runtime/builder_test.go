package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/fintechain/skeleton/internal/domain/component"
	pkgruntime "github.com/fintechain/skeleton/internal/domain/runtime"
	infraRuntime "github.com/fintechain/skeleton/internal/infrastructure/runtime"
	"github.com/fintechain/skeleton/test/unit/mocks"
)

// TestNewRuntimeWithOptions tests the runtime builder
func TestNewRuntimeWithOptions(t *testing.T) {
	tests := []struct {
		name        string
		options     []infraRuntime.RuntimeOption
		expectError bool
		description string
	}{
		{
			name:        "default runtime",
			options:     []infraRuntime.RuntimeOption{},
			expectError: false,
			description: "Should create runtime with default settings",
		},

		{
			name: "runtime with custom registry",
			options: func() []infraRuntime.RuntimeOption {
				mockRegistry := mocks.NewFactory().RegistryInterface()
				mockRegistry.On("Register", mock.Anything).Return(nil)
				return []infraRuntime.RuntimeOption{
					infraRuntime.WithRegistry(mockRegistry),
				}
			}(),
			expectError: false,
			description: "Should create runtime with custom registry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runtime, err := infraRuntime.NewRuntimeWithOptions(tt.options...)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, runtime)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, runtime)

				// Verify interface compliance
				var _ pkgruntime.RuntimeEnvironment = runtime
			}
		})
	}
}

// TestRuntimeBuilderWithPlugins tests plugin loading during runtime creation
func TestRuntimeBuilderWithPlugins(t *testing.T) {
	factory := mocks.NewFactory()
	mockPlugin := factory.PluginInterface()

	// Configure mock plugin - only expect Initialize since system is not running
	mockPlugin.On("ID").Return(component.ComponentID("test-plugin"))
	mockPlugin.On("Initialize", mock.Anything, mock.Anything).Return(nil)

	runtime, err := infraRuntime.NewRuntimeWithOptions(
		infraRuntime.WithPlugins(mockPlugin),
	)

	assert.NoError(t, err)
	assert.NotNil(t, runtime)

	// Verify plugin manager is accessible
	pluginManager := runtime.PluginManager()
	assert.NotNil(t, pluginManager)

	// Verify mock expectations
	mockPlugin.AssertExpectations(t)
}

// TestRuntimeBuilderInterfaceCompliance verifies the builder creates compliant runtimes
func TestRuntimeBuilderInterfaceCompliance(t *testing.T) {
	runtime, err := infraRuntime.NewRuntimeWithOptions()
	assert.NoError(t, err)

	// Verify interface compliance
	var _ pkgruntime.RuntimeEnvironment = runtime
	assert.NotNil(t, runtime)

	// Test basic functionality
	assert.NotNil(t, runtime.Registry())
	assert.False(t, runtime.IsRunning())
}
