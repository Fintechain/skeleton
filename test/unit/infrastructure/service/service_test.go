package service

import (
	"sync"
	"testing"
	"time"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/service"
	componentImpl "github.com/fintechain/skeleton/internal/infrastructure/component"
	serviceImpl "github.com/fintechain/skeleton/internal/infrastructure/service"
	"github.com/fintechain/skeleton/pkg/context"
	"github.com/fintechain/skeleton/test/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	tests := []struct {
		name          string
		baseComponent component.Component
		expectNil     bool
		description   string
	}{
		{
			name:          "valid component",
			baseComponent: mocks.NewFactory().ComponentInterface(),
			expectNil:     false,
			description:   "Should create service with valid component",
		},
		{
			name:          "nil component",
			baseComponent: nil,
			expectNil:     true,
			description:   "Should return nil with nil component",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			svc := serviceImpl.NewService(tt.baseComponent)

			// Verify
			if tt.expectNil {
				assert.Nil(t, svc, tt.description)
			} else {
				assert.NotNil(t, svc, tt.description)

				// Verify interface compliance
				var _ service.Service = svc

				// Verify component delegation
				assert.Equal(t, tt.baseComponent.ID(), svc.ID())
				assert.Equal(t, tt.baseComponent.Name(), svc.Name())
				assert.Equal(t, tt.baseComponent.Type(), svc.Type())

				// Verify initial status
				assert.Equal(t, service.StatusStopped, svc.Status())
			}
		})
	}
}

func TestServiceInterfaceCompliance(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()

	// Test interface compliance
	svc := serviceImpl.NewService(mockComponent)
	require.NotNil(t, svc)

	// Verify service interface
	var _ service.Service = svc

	// Verify component interface (through embedding)
	var _ component.Component = svc
}

func TestServiceLifecycle(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.Component().
		WithID("test-service").
		WithName("Test Service").
		Build()

	svc := serviceImpl.NewService(mockComponent)
	require.NotNil(t, svc)

	ctx := context.NewContext()

	// Test initial state
	assert.Equal(t, service.StatusStopped, svc.Status())

	// Test start
	err := svc.Start(ctx)
	assert.NoError(t, err)
	assert.Equal(t, service.StatusRunning, svc.Status())

	// Test stop
	err = svc.Stop(ctx)
	assert.NoError(t, err)
	assert.Equal(t, service.StatusStopped, svc.Status())
}

func TestServiceStart(t *testing.T) {
	tests := []struct {
		name           string
		setupService   func() service.Service
		setupContext   func() context.Context
		expectError    bool
		errorContains  string
		expectedStatus service.ServiceStatus
		description    string
	}{
		{
			name: "successful start",
			setupService: func() service.Service {
				factory := mocks.NewFactory()
				mockComponent := factory.Component().
					WithID("start-service").
					Build()
				return serviceImpl.NewService(mockComponent)
			},
			setupContext: func() context.Context {
				return context.NewContext()
			},
			expectError:    false,
			expectedStatus: service.StatusRunning,
			description:    "Should start successfully with valid context",
		},
		{
			name: "nil context",
			setupService: func() service.Service {
				factory := mocks.NewFactory()
				mockComponent := factory.ComponentInterface()
				return serviceImpl.NewService(mockComponent)
			},
			setupContext: func() context.Context {
				return nil
			},
			expectError:    true,
			errorContains:  "context is required",
			expectedStatus: service.StatusStopped,
			description:    "Should fail with nil context",
		},
		{
			name: "start already running service",
			setupService: func() service.Service {
				factory := mocks.NewFactory()
				mockComponent := factory.ComponentInterface()
				svc := serviceImpl.NewService(mockComponent)

				// Start the service first
				ctx := context.NewContext()
				svc.Start(ctx)

				return svc
			},
			setupContext: func() context.Context {
				return context.NewContext()
			},
			expectError:    false, // Should be idempotent
			expectedStatus: service.StatusRunning,
			description:    "Should handle starting already running service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			svc := tt.setupService()
			ctx := tt.setupContext()

			// Execute
			err := svc.Start(ctx)

			// Verify
			if tt.expectError {
				assert.Error(t, err, tt.description)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err, tt.description)
			}

			assert.Equal(t, tt.expectedStatus, svc.Status(), tt.description)
		})
	}
}

func TestServiceStop(t *testing.T) {
	tests := []struct {
		name           string
		setupService   func() service.Service
		setupContext   func() context.Context
		expectError    bool
		errorContains  string
		expectedStatus service.ServiceStatus
		description    string
	}{
		{
			name: "successful stop",
			setupService: func() service.Service {
				factory := mocks.NewFactory()
				mockComponent := factory.ComponentInterface()
				svc := serviceImpl.NewService(mockComponent)

				// Start the service first
				ctx := context.NewContext()
				svc.Start(ctx)

				return svc
			},
			setupContext: func() context.Context {
				return context.NewContext()
			},
			expectError:    false,
			expectedStatus: service.StatusStopped,
			description:    "Should stop successfully when running",
		},
		{
			name: "nil context",
			setupService: func() service.Service {
				factory := mocks.NewFactory()
				mockComponent := factory.ComponentInterface()
				svc := serviceImpl.NewService(mockComponent)

				// Start the service first
				ctx := context.NewContext()
				svc.Start(ctx)

				return svc
			},
			setupContext: func() context.Context {
				return nil
			},
			expectError:    true,
			errorContains:  "context is required",
			expectedStatus: service.StatusRunning, // Should remain running on error
			description:    "Should fail with nil context",
		},
		{
			name: "stop already stopped service",
			setupService: func() service.Service {
				factory := mocks.NewFactory()
				mockComponent := factory.ComponentInterface()
				return serviceImpl.NewService(mockComponent)
			},
			setupContext: func() context.Context {
				return context.NewContext()
			},
			expectError:    false, // Should be idempotent
			expectedStatus: service.StatusStopped,
			description:    "Should handle stopping already stopped service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			svc := tt.setupService()
			ctx := tt.setupContext()

			// Execute
			err := svc.Stop(ctx)

			// Verify
			if tt.expectError {
				assert.Error(t, err, tt.description)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err, tt.description)
			}

			assert.Equal(t, tt.expectedStatus, svc.Status(), tt.description)
		})
	}
}

func TestServiceStatus(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()
	svc := serviceImpl.NewService(mockComponent)
	require.NotNil(t, svc)

	ctx := context.NewContext()

	// Test status transitions
	assert.Equal(t, service.StatusStopped, svc.Status())

	// Start service
	err := svc.Start(ctx)
	require.NoError(t, err)
	assert.Equal(t, service.StatusRunning, svc.Status())

	// Stop service
	err = svc.Stop(ctx)
	require.NoError(t, err)
	assert.Equal(t, service.StatusStopped, svc.Status())
}

func TestServiceConcurrency(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.Component().
		WithID("concurrent-service").
		Build()

	svc := serviceImpl.NewService(mockComponent)
	require.NotNil(t, svc)

	// Test concurrent start/stop operations
	numGoroutines := 10
	var wg sync.WaitGroup
	errors := make([]error, numGoroutines*2) // start + stop for each goroutine

	for i := 0; i < numGoroutines; i++ {
		wg.Add(2) // One for start, one for stop

		go func(index int) {
			defer wg.Done()
			ctx := context.NewContext()
			errors[index*2] = svc.Start(ctx)
		}(i)

		go func(index int) {
			defer wg.Done()
			ctx := context.NewContext()
			// Add small delay to ensure start happens first
			time.Sleep(time.Millisecond)
			errors[index*2+1] = svc.Stop(ctx)
		}(i)
	}

	wg.Wait()

	// Verify no errors occurred (operations should be thread-safe)
	for i, err := range errors {
		assert.NoError(t, err, "Operation %d should succeed", i)
	}

	// Final status should be deterministic
	finalStatus := svc.Status()
	assert.True(t, finalStatus == service.StatusRunning || finalStatus == service.StatusStopped)
}

func TestServiceErrorHandling(t *testing.T) {
	factory := mocks.NewFactory()

	t.Run("service with nil component", func(t *testing.T) {
		// Create service with component that becomes nil (simulating error condition)
		svc := serviceImpl.NewService(factory.ComponentInterface())

		// Use type assertion to access internal fields for testing
		if baseSvc, ok := svc.(*serviceImpl.BaseService); ok {
			// This tests the error path where component is nil during operations
			baseSvc.Component = nil

			ctx := context.NewContext()

			// Test start with nil component
			err := baseSvc.Start(ctx)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), service.ErrServiceNotFound)

			// Test stop with nil component
			err = baseSvc.Stop(ctx)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), service.ErrServiceNotFound)
		}
	})
}

func TestServiceComponentDelegation(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.Component().
		WithID("test-id").
		WithName("Test Service").
		WithType(component.TypeService).
		WithDescription("Test service description").
		Build()

	svc := serviceImpl.NewService(mockComponent)
	require.NotNil(t, svc)

	// Test component method delegation
	assert.Equal(t, "test-id", svc.ID())
	assert.Equal(t, "Test Service", svc.Name())
	assert.Equal(t, component.TypeService, svc.Type())
	assert.Equal(t, "Test service description", svc.Description())

	// Test metadata delegation
	metadata := svc.Metadata()
	assert.NotNil(t, metadata)
}

func TestServiceLifecycleIntegration(t *testing.T) {
	factory := mocks.NewFactory()
	mockSystem := factory.SystemInterface()

	// Test service initialization with system
	mockComponent := factory.Component().
		WithID("lifecycle-service").
		Build()

	svc := serviceImpl.NewService(mockComponent)
	require.NotNil(t, svc)

	// Test initialization
	ctx := context.NewContext()
	err := svc.Initialize(ctx, mockSystem)
	assert.NoError(t, err)

	// Test service lifecycle after initialization
	err = svc.Start(ctx)
	assert.NoError(t, err)
	assert.Equal(t, service.StatusRunning, svc.Status())

	err = svc.Stop(ctx)
	assert.NoError(t, err)
	assert.Equal(t, service.StatusStopped, svc.Status())

	// Test disposal
	err = svc.Dispose()
	assert.NoError(t, err)
}

func TestServiceWithLifecycleAwareComponent(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()
	lifecycleAware := componentImpl.NewLifecycleAwareComponent(mockComponent)

	svc := serviceImpl.NewService(lifecycleAware)
	require.NotNil(t, svc)

	ctx := context.NewContext()

	// Test start with lifecycle aware component
	err := svc.Start(ctx)
	assert.NoError(t, err)
	assert.Equal(t, service.StatusRunning, svc.Status())
	assert.Equal(t, component.StateActive, lifecycleAware.State())

	// Test stop with lifecycle aware component
	err = svc.Stop(ctx)
	assert.NoError(t, err)
	assert.Equal(t, service.StatusStopped, svc.Status())
	assert.Equal(t, component.StateDisposed, lifecycleAware.State())
}

func TestServiceStatusConstants(t *testing.T) {
	// Test that status constants are properly defined
	assert.Equal(t, "stopped", string(service.StatusStopped))
	assert.Equal(t, "starting", string(service.StatusStarting))
	assert.Equal(t, "running", string(service.StatusRunning))
	assert.Equal(t, "stopping", string(service.StatusStopping))
	assert.Equal(t, "failed", string(service.StatusFailed))
}

func TestServiceMultipleStartStop(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.Component().
		WithID("multi-lifecycle-service").
		Build()

	svc := serviceImpl.NewService(mockComponent)
	require.NotNil(t, svc)

	ctx := context.NewContext()

	// Test multiple start/stop cycles
	for i := 0; i < 3; i++ {
		// Start
		err := svc.Start(ctx)
		assert.NoError(t, err, "Start cycle %d should succeed", i)
		assert.Equal(t, service.StatusRunning, svc.Status(), "Should be running after start %d", i)

		// Stop
		err = svc.Stop(ctx)
		assert.NoError(t, err, "Stop cycle %d should succeed", i)
		assert.Equal(t, service.StatusStopped, svc.Status(), "Should be stopped after stop %d", i)
	}
}

func TestServiceContextPropagation(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()

	svc := serviceImpl.NewService(mockComponent)
	require.NotNil(t, svc)

	// Test with context containing metadata
	ctx := context.NewContext()
	ctx = ctx.WithValue("service_metadata", map[string]interface{}{
		"environment": "test",
		"version":     "1.0.0",
	})

	// Start with context
	err := svc.Start(ctx)
	assert.NoError(t, err)
	assert.Equal(t, service.StatusRunning, svc.Status())

	// Stop with context
	err = svc.Stop(ctx)
	assert.NoError(t, err)
	assert.Equal(t, service.StatusStopped, svc.Status())
}

func TestServiceErrorStates(t *testing.T) {
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()

	svc := serviceImpl.NewService(mockComponent)
	require.NotNil(t, svc)

	// Test error handling during start
	t.Run("start with invalid context", func(t *testing.T) {
		err := svc.Start(nil)
		assert.Error(t, err)
		assert.Equal(t, service.StatusStopped, svc.Status())
	})

	// Test error handling during stop
	t.Run("stop with invalid context after successful start", func(t *testing.T) {
		ctx := context.NewContext()

		// Start successfully
		err := svc.Start(ctx)
		require.NoError(t, err)
		assert.Equal(t, service.StatusRunning, svc.Status())

		// Try to stop with nil context
		err = svc.Stop(nil)
		assert.Error(t, err)
		assert.Equal(t, service.StatusRunning, svc.Status()) // Should remain running

		// Clean stop
		err = svc.Stop(ctx)
		assert.NoError(t, err)
		assert.Equal(t, service.StatusStopped, svc.Status())
	})
}
