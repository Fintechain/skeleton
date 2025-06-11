package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/config"
	"github.com/fintechain/skeleton/internal/domain/logging"
	infraconfig "github.com/fintechain/skeleton/internal/infrastructure/config"
	infraruntime "github.com/fintechain/skeleton/internal/infrastructure/runtime"
	"github.com/fintechain/skeleton/pkg/event"
	"github.com/fintechain/skeleton/pkg/plugin"
	"github.com/fintechain/skeleton/test/unit/mocks"
)

// createTestConfiguration creates a simple memory configuration for testing
func createTestConfiguration() config.Configuration {
	return infraconfig.NewMemoryConfiguration()
}

// createTestDependencies creates test dependencies for runtime construction
func createTestDependencies() (
	component.Registry,
	config.Configuration,
	*mocks.MockPluginManager,
	*mocks.MockEventBusService,
	*mocks.MockLoggerService,
) {
	factory := mocks.NewFactory()
	return factory.RegistryInterface(),
		createTestConfiguration(),
		factory.PluginManagerInterface(),
		factory.EventBusServiceInterface(),
		factory.LoggerServiceInterface()
}

// TestRuntimeInterfaceCompliance verifies that Runtime implements the RuntimeEnvironment interface
func TestRuntimeInterfaceCompliance(t *testing.T) {
	registry, config, pluginManager, eventBus, logger := createTestDependencies()

	// Create runtime
	runtime, err := infraruntime.NewRuntime(registry, config, pluginManager, eventBus, logger)
	assert.NoError(t, err)
	assert.NotNil(t, runtime)
}

// TestNewRuntime tests the constructor for Runtime
func TestNewRuntime(t *testing.T) {
	factory := mocks.NewFactory()
	validRegistry := factory.RegistryInterface()
	validConfig := createTestConfiguration()
	validPluginManager := factory.PluginManagerInterface()
	validEventBus := factory.EventBusServiceInterface()
	validLogger := factory.LoggerServiceInterface()

	tests := []struct {
		name          string
		registry      component.Registry
		config        config.Configuration
		pluginManager plugin.PluginManager
		eventBus      event.EventBusService
		logger        logging.LoggerService
		expectError   bool
		errorMsg      string
		description   string
	}{
		{
			name:          "valid dependencies",
			registry:      validRegistry,
			config:        validConfig,
			pluginManager: validPluginManager,
			eventBus:      validEventBus,
			logger:        validLogger,
			expectError:   false,
			description:   "Should create runtime with valid dependencies",
		},
		{
			name:          "nil registry",
			registry:      nil,
			config:        validConfig,
			pluginManager: validPluginManager,
			eventBus:      validEventBus,
			logger:        validLogger,
			expectError:   true,
			errorMsg:      infraruntime.ErrNilRegistry,
			description:   "Should return error for nil registry",
		},
		{
			name:          "nil configuration",
			registry:      validRegistry,
			config:        nil,
			pluginManager: validPluginManager,
			eventBus:      validEventBus,
			logger:        validLogger,
			expectError:   true,
			errorMsg:      infraruntime.ErrNilConfiguration,
			description:   "Should return error for nil configuration",
		},
		{
			name:          "nil plugin manager",
			registry:      validRegistry,
			config:        validConfig,
			pluginManager: nil,
			eventBus:      validEventBus,
			logger:        validLogger,
			expectError:   true,
			errorMsg:      infraruntime.ErrNilPluginManager,
			description:   "Should return error for nil plugin manager",
		},
		{
			name:          "nil event bus",
			registry:      validRegistry,
			config:        validConfig,
			pluginManager: validPluginManager,
			eventBus:      nil,
			logger:        validLogger,
			expectError:   true,
			errorMsg:      infraruntime.ErrNilEventBus,
			description:   "Should return error for nil event bus",
		},
		{
			name:          "nil logger",
			registry:      validRegistry,
			config:        validConfig,
			pluginManager: validPluginManager,
			eventBus:      validEventBus,
			logger:        nil,
			expectError:   true,
			errorMsg:      infraruntime.ErrNilLogger,
			description:   "Should return error for nil logger",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runtime, err := infraruntime.NewRuntime(tt.registry, tt.config, tt.pluginManager, tt.eventBus, tt.logger)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, runtime)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, runtime)
			}
		})
	}
}

// TestRuntimeDirectAccess tests direct access to dependencies
func TestRuntimeDirectAccess(t *testing.T) {
	registry, config, pluginManager, eventBus, logger := createTestDependencies()

	// Create runtime
	runtime, err := infraruntime.NewRuntime(registry, config, pluginManager, eventBus, logger)
	assert.NoError(t, err)

	// Test direct access to dependencies
	assert.Equal(t, registry, runtime.Registry())
	assert.Equal(t, config, runtime.Configuration())

	// Test service accessors (direct access with dependency injection)
	pm := runtime.PluginManager()
	assert.Equal(t, pluginManager, pm)

	eb := runtime.EventBus()
	assert.Equal(t, eventBus, eb)

	lg := runtime.Logger()
	assert.Equal(t, logger, lg)
}

// TestRuntimeState tests runtime state management
func TestRuntimeState(t *testing.T) {
	registry, config, pluginManager, eventBus, logger := createTestDependencies()

	// Set up service mocks for start/stop
	pluginManager.On("Start", mock.Anything).Return(nil)
	eventBus.On("Start", mock.Anything).Return(nil)
	logger.On("Start", mock.Anything).Return(nil)
	pluginManager.On("Stop", mock.Anything).Return(nil)
	eventBus.On("Stop", mock.Anything).Return(nil)
	logger.On("Stop", mock.Anything).Return(nil)

	// Create runtime
	runtime, err := infraruntime.NewRuntime(registry, config, pluginManager, eventBus, logger)
	assert.NoError(t, err)

	// Test initial state
	assert.False(t, runtime.IsRunning())

	// Test start
	err = runtime.Start(nil)
	assert.NoError(t, err)
	assert.True(t, runtime.IsRunning())

	// Test stop
	err = runtime.Stop(nil)
	assert.NoError(t, err)
	assert.False(t, runtime.IsRunning())

	// Verify mock expectations
	pluginManager.AssertExpectations(t)
	eventBus.AssertExpectations(t)
	logger.AssertExpectations(t)
}
