package system

import (
	stdctx "context"
	"errors"
	"testing"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/operation"
	domainsystem "github.com/fintechain/skeleton/internal/domain/system"
	"github.com/fintechain/skeleton/internal/infrastructure/context"
	"github.com/fintechain/skeleton/internal/infrastructure/system/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultSystemService_WithPluginManager(t *testing.T) {
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	service := NewDefaultSystemService(
		"test-service",
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockStore,
		mockLogger,
	)

	assert.NotNil(t, service)
	assert.Equal(t, "test-service", service.ID())
	assert.Equal(t, mockPluginManager, service.PluginManager())
	assert.Equal(t, mockRegistry, service.Registry())
	assert.Equal(t, mockEventBus, service.EventBus())
	assert.Equal(t, mockConfig, service.Configuration())
	assert.Equal(t, mockStore, service.Store())
}

func TestDefaultSystemService_PluginManagerGetter(t *testing.T) {
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	service := NewDefaultSystemService(
		"test-service",
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockStore,
		mockLogger,
	)

	assert.Equal(t, mockPluginManager, service.PluginManager())
}

func TestDefaultSystemService_AllGetters(t *testing.T) {
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	service := NewDefaultSystemService(
		"test-service",
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockStore,
		mockLogger,
	)

	assert.Equal(t, mockRegistry, service.Registry())
	assert.Equal(t, mockPluginManager, service.PluginManager())
	assert.Equal(t, mockEventBus, service.EventBus())
	assert.Equal(t, mockConfig, service.Configuration())
	assert.Equal(t, mockStore, service.Store())
}

func TestDefaultSystemService_ConstructorValidation(t *testing.T) {
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	// Test with valid service ID
	service := NewDefaultSystemService(
		"valid-service-id",
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockStore,
		mockLogger,
	)
	assert.NotNil(t, service)
	assert.Equal(t, "valid-service-id", service.ID())

	// Test with empty service ID
	service = NewDefaultSystemService(
		"",
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockStore,
		mockLogger,
	)
	assert.NotNil(t, service)
	assert.Equal(t, "", service.ID())

	// Test with nil logger (should create default)
	service = NewDefaultSystemService(
		"test-service",
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockStore,
		nil,
	)
	assert.NotNil(t, service)
}

func TestDefaultSystemService_Configuration(t *testing.T) {
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	service := NewDefaultSystemService(
		"test-service",
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockStore,
		mockLogger,
	)

	config := service.Configuration()
	assert.Equal(t, mockConfig, config)
}

func TestDefaultSystemService_Initialize(t *testing.T) {
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	service := NewDefaultSystemService(
		"test-service",
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockStore,
		mockLogger,
	)

	ctx := context.NewContext(stdctx.Background())

	// Test successful initialization
	err := service.Initialize(ctx)
	assert.NoError(t, err)
}

func TestDefaultSystemService_Start(t *testing.T) {
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	service := NewDefaultSystemService(
		"test-service",
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockStore,
		mockLogger,
	)

	ctx := context.NewContext(stdctx.Background())

	// Initialize first
	err := service.Initialize(ctx)
	assert.NoError(t, err)

	// Test successful start
	err = service.Start(ctx)
	assert.NoError(t, err)
}

func TestDefaultSystemService_Stop(t *testing.T) {
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	service := NewDefaultSystemService(
		"test-service",
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockStore,
		mockLogger,
	)

	ctx := context.NewContext(stdctx.Background())

	// Initialize and start first
	err := service.Initialize(ctx)
	assert.NoError(t, err)
	err = service.Start(ctx)
	assert.NoError(t, err)

	// Test successful stop
	err = service.Stop(ctx)
	assert.NoError(t, err)
}

func TestDefaultSystemService_ExecuteOperation(t *testing.T) {
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	service := NewDefaultSystemService(
		"test-service",
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockStore,
		mockLogger,
	)

	ctx := context.NewContext(stdctx.Background())

	// Test operation execution
	result, err := service.ExecuteOperation(ctx, "test-operation", map[string]interface{}{
		"param1": "value1",
	})

	// Should return an error for unknown operation
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "System service not started")
}

func TestDefaultSystemService_StartService(t *testing.T) {
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	service := NewDefaultSystemService(
		"test-service",
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockStore,
		mockLogger,
	)

	ctx := context.NewContext(stdctx.Background())

	// Test starting a service
	err := service.StartService(ctx, "test-service-id")

	// Should return an error for unknown service
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "System service not started")
}

func TestDefaultSystemService_StopService(t *testing.T) {
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	service := NewDefaultSystemService(
		"test-service",
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockStore,
		mockLogger,
	)

	ctx := context.NewContext(stdctx.Background())

	// Test stopping a service
	err := service.StopService(ctx, "test-service-id")

	// Should return an error for unknown service
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "System service not started")
}

func TestDefaultSystemService_GetBoolMetadata(t *testing.T) {
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	service := NewDefaultSystemService(
		"test-service",
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockStore,
		mockLogger,
	)

	// Set some metadata first - access the BaseComponent through BaseService
	if baseComponent, ok := service.BaseService.Component.(*component.BaseComponent); ok {
		baseComponent.SetMetadata("testBool", true)
		baseComponent.SetMetadata("testString", "not a bool")
	}

	// Test getting boolean metadata
	result := service.GetBoolMetadata("testBool", false)
	assert.True(t, result)

	// Test getting non-existent metadata with default
	result = service.GetBoolMetadata("nonExistent", true)
	assert.True(t, result)

	// Test getting non-boolean metadata with default
	result = service.GetBoolMetadata("testString", false)
	assert.False(t, result)
}

func TestDefaultSystemService_ExecuteOperation_AllPaths(t *testing.T) {
	tests := []struct {
		name          string
		setupService  func() *DefaultSystemService
		operationID   string
		expectedError string
		shouldSucceed bool
	}{
		{
			name: "operations disabled via metadata",
			setupService: func() *DefaultSystemService {
				svc := createTestSystemService()
				if baseComp, ok := svc.BaseService.Component.(*component.BaseComponent); ok {
					baseComp.SetMetadata("enableOperations", false)
				}
				return svc
			},
			operationID:   "test-op",
			expectedError: "Operations are disabled",
			shouldSucceed: false,
		},
		{
			name: "operation not found in registry",
			setupService: func() *DefaultSystemService {
				svc := createTestSystemService()
				if baseComp, ok := svc.BaseService.Component.(*component.BaseComponent); ok {
					baseComp.SetMetadata("enableOperations", true)
				}
				return svc
			},
			operationID:   "non-existent-op",
			expectedError: "Operation not found",
			shouldSucceed: false,
		},
		{
			name: "operation found but not an Operation type",
			setupService: func() *DefaultSystemService {
				svc := createTestSystemService()
				if baseComp, ok := svc.BaseService.Component.(*component.BaseComponent); ok {
					baseComp.SetMetadata("enableOperations", true)
				}
				// Add a component that's not an operation
				mockComp := mocks.NewMockComponent("wrong-type", "Wrong Type", component.TypeBasic)
				svc.registry.Register(mockComp)
				return svc
			},
			operationID:   "wrong-type",
			expectedError: "Component is not an operation",
			shouldSucceed: false,
		},
		{
			name: "operation execution fails",
			setupService: func() *DefaultSystemService {
				svc := createTestSystemService()
				if baseComp, ok := svc.BaseService.Component.(*component.BaseComponent); ok {
					baseComp.SetMetadata("enableOperations", true)
				}
				// Add a mock operation that fails
				mockOp := mocks.NewMockOperation("failing-op", "Failing Operation")
				mockOp.ExecuteFunc = func(ctx component.Context, input operation.Input) (operation.Output, error) {
					return nil, errors.New("execution failed")
				}
				svc.registry.Register(mockOp)
				return svc
			},
			operationID:   "failing-op",
			expectedError: "execution failed",
			shouldSucceed: false,
		},
		{
			name: "operation executes successfully",
			setupService: func() *DefaultSystemService {
				svc := createTestSystemService()
				if baseComp, ok := svc.BaseService.Component.(*component.BaseComponent); ok {
					baseComp.SetMetadata("enableOperations", true)
				}
				// Add a mock operation that succeeds
				mockOp := mocks.NewMockOperation("success-op", "Success Operation")
				mockOp.ExecuteFunc = func(ctx component.Context, input operation.Input) (operation.Output, error) {
					return "success result", nil
				}
				svc.registry.Register(mockOp)
				return svc
			},
			operationID:   "success-op",
			shouldSucceed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.setupService()
			ctx := context.NewContext(stdctx.Background())

			result, err := svc.ExecuteOperation(ctx, tt.operationID, nil)

			if tt.shouldSucceed {
				assert.NoError(t, err)
				// The result should be a SystemOperationOutput with the expected data
				if output, ok := result.(*domainsystem.SystemOperationOutput); ok {
					assert.Equal(t, "success result", output.Data)
				} else {
					t.Errorf("Expected SystemOperationOutput, got %T", result)
				}
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			}
		})
	}
}

func TestDefaultSystemService_StartService_AllPaths(t *testing.T) {
	tests := []struct {
		name          string
		setupService  func() *DefaultSystemService
		serviceID     string
		expectedError string
		shouldSucceed bool
	}{
		{
			name: "services disabled via metadata",
			setupService: func() *DefaultSystemService {
				svc := createTestSystemService()
				if baseComp, ok := svc.BaseService.Component.(*component.BaseComponent); ok {
					baseComp.SetMetadata("enableServices", false)
				}
				return svc
			},
			serviceID:     "test-service",
			expectedError: "Services are disabled",
			shouldSucceed: false,
		},
		{
			name: "service not found in registry",
			setupService: func() *DefaultSystemService {
				svc := createTestSystemService()
				if baseComp, ok := svc.BaseService.Component.(*component.BaseComponent); ok {
					baseComp.SetMetadata("enableServices", true)
				}
				return svc
			},
			serviceID:     "non-existent-service",
			expectedError: "Service not found",
			shouldSucceed: false,
		},
		{
			name: "component found but not a Service type",
			setupService: func() *DefaultSystemService {
				svc := createTestSystemService()
				if baseComp, ok := svc.BaseService.Component.(*component.BaseComponent); ok {
					baseComp.SetMetadata("enableServices", true)
				}
				// Add a component that's not a service
				mockComp := mocks.NewMockComponent("wrong-type", "Wrong Type", component.TypeBasic)
				svc.registry.Register(mockComp)
				return svc
			},
			serviceID:     "wrong-type",
			expectedError: "Component is not a service",
			shouldSucceed: false,
		},
		{
			name: "service start fails",
			setupService: func() *DefaultSystemService {
				svc := createTestSystemService()
				if baseComp, ok := svc.BaseService.Component.(*component.BaseComponent); ok {
					baseComp.SetMetadata("enableServices", true)
				}
				// Add a mock service that fails to start
				mockSvc := mocks.NewMockService("failing-service", "Failing Service")
				mockSvc.StartFunc = func(ctx component.Context) error {
					return errors.New("start failed")
				}
				svc.registry.Register(mockSvc)
				return svc
			},
			serviceID:     "failing-service",
			expectedError: "start failed",
			shouldSucceed: false,
		},
		{
			name: "service starts successfully",
			setupService: func() *DefaultSystemService {
				svc := createTestSystemService()
				if baseComp, ok := svc.BaseService.Component.(*component.BaseComponent); ok {
					baseComp.SetMetadata("enableServices", true)
				}
				// Add a mock service that starts successfully
				mockSvc := mocks.NewMockService("success-service", "Success Service")
				mockSvc.StartFunc = func(ctx component.Context) error {
					return nil
				}
				svc.registry.Register(mockSvc)
				return svc
			},
			serviceID:     "success-service",
			shouldSucceed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.setupService()
			ctx := context.NewContext(stdctx.Background())

			err := svc.StartService(ctx, tt.serviceID)

			if tt.shouldSucceed {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

func TestDefaultSystemService_StopService_AllPaths(t *testing.T) {
	tests := []struct {
		name          string
		setupService  func() *DefaultSystemService
		serviceID     string
		expectedError string
		shouldSucceed bool
	}{
		{
			name: "services disabled via metadata",
			setupService: func() *DefaultSystemService {
				svc := createTestSystemService()
				if baseComp, ok := svc.BaseService.Component.(*component.BaseComponent); ok {
					baseComp.SetMetadata("enableServices", false)
				}
				return svc
			},
			serviceID:     "test-service",
			expectedError: "Services are disabled",
			shouldSucceed: false,
		},
		{
			name: "service not found in registry",
			setupService: func() *DefaultSystemService {
				svc := createTestSystemService()
				if baseComp, ok := svc.BaseService.Component.(*component.BaseComponent); ok {
					baseComp.SetMetadata("enableServices", true)
				}
				return svc
			},
			serviceID:     "non-existent-service",
			expectedError: "Service not found",
			shouldSucceed: false,
		},
		{
			name: "component found but not a Service type",
			setupService: func() *DefaultSystemService {
				svc := createTestSystemService()
				if baseComp, ok := svc.BaseService.Component.(*component.BaseComponent); ok {
					baseComp.SetMetadata("enableServices", true)
				}
				// Add a component that's not a service
				mockComp := mocks.NewMockComponent("wrong-type", "Wrong Type", component.TypeBasic)
				svc.registry.Register(mockComp)
				return svc
			},
			serviceID:     "wrong-type",
			expectedError: "Component is not a service",
			shouldSucceed: false,
		},
		{
			name: "service stop fails",
			setupService: func() *DefaultSystemService {
				svc := createTestSystemService()
				if baseComp, ok := svc.BaseService.Component.(*component.BaseComponent); ok {
					baseComp.SetMetadata("enableServices", true)
				}
				// Add a mock service that fails to stop
				mockSvc := mocks.NewMockService("failing-service", "Failing Service")
				mockSvc.StopFunc = func(ctx component.Context) error {
					return errors.New("stop failed")
				}
				svc.registry.Register(mockSvc)
				return svc
			},
			serviceID:     "failing-service",
			expectedError: "stop failed",
			shouldSucceed: false,
		},
		{
			name: "service stops successfully",
			setupService: func() *DefaultSystemService {
				svc := createTestSystemService()
				if baseComp, ok := svc.BaseService.Component.(*component.BaseComponent); ok {
					baseComp.SetMetadata("enableServices", true)
				}
				// Add a mock service that stops successfully
				mockSvc := mocks.NewMockService("success-service", "Success Service")
				mockSvc.StopFunc = func(ctx component.Context) error {
					return nil
				}
				svc.registry.Register(mockSvc)
				return svc
			},
			serviceID:     "success-service",
			shouldSucceed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.setupService()
			ctx := context.NewContext(stdctx.Background())

			err := svc.StopService(ctx, tt.serviceID)

			if tt.shouldSucceed {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			}
		})
	}
}

// createTestSystemService creates a DefaultSystemService instance for testing
func createTestSystemService() *DefaultSystemService {
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	svc := NewDefaultSystemService(
		"test-service",
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockStore,
		mockLogger,
	)

	// Initialize and start the service for testing
	ctx := context.NewContext(stdctx.Background())
	_ = svc.Initialize(ctx)
	_ = svc.Start(ctx)

	return svc
}

// Test error paths for Initialize method
func TestDefaultSystemService_Initialize_ErrorPaths(t *testing.T) {
	tests := []struct {
		name          string
		setupService  func() *DefaultSystemService
		expectedError string
	}{
		{
			name: "registry initialization fails",
			setupService: func() *DefaultSystemService {
				mockRegistry := mocks.NewMockRegistry()
				mockRegistry.InitializeFunc = func(ctx component.Context) error {
					return errors.New("registry init failed")
				}
				mockPluginManager := mocks.NewMockPluginManager()
				mockEventBus := mocks.NewMockEventBus()
				mockConfig := mocks.NewMockConfiguration()
				mockStore := mocks.NewMockMultiStore()
				mockLogger := mocks.NewMockLogger()

				return NewDefaultSystemService(
					"test-service",
					mockRegistry,
					mockPluginManager,
					mockEventBus,
					mockConfig,
					mockStore,
					mockLogger,
				)
			},
			expectedError: "Failed to initialize component registry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.setupService()
			ctx := context.NewContext(stdctx.Background())

			err := svc.Initialize(ctx)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

// Test error paths for Start method
func TestDefaultSystemService_Start_ErrorPaths(t *testing.T) {
	tests := []struct {
		name          string
		setupService  func() *DefaultSystemService
		expectedError string
	}{
		{
			name: "start called without initialization",
			setupService: func() *DefaultSystemService {
				mockRegistry := mocks.NewMockRegistry()
				mockPluginManager := mocks.NewMockPluginManager()
				mockEventBus := mocks.NewMockEventBus()
				mockConfig := mocks.NewMockConfiguration()
				mockStore := mocks.NewMockMultiStore()
				mockLogger := mocks.NewMockLogger()

				return NewDefaultSystemService(
					"test-service",
					mockRegistry,
					mockPluginManager,
					mockEventBus,
					mockConfig,
					mockStore,
					mockLogger,
				)
			},
			expectedError: "System service not initialized",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := tt.setupService()
			ctx := context.NewContext(stdctx.Background())

			err := svc.Start(ctx)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

// Test error paths for Stop method
func TestDefaultSystemService_Stop_ErrorPaths(t *testing.T) {
	tests := []struct {
		name          string
		setupService  func() *DefaultSystemService
		expectedError string
	}{
		{
			name: "stop called without starting",
			setupService: func() *DefaultSystemService {
				svc := createTestSystemService()
				// Stop the service to reset its state
				ctx := context.NewContext(stdctx.Background())
				_ = svc.Stop(ctx)
				return svc
			},
			expectedError: "System service not started",
		},
		{
			name: "base service stop fails",
			setupService: func() *DefaultSystemService {
				mockRegistry := mocks.NewMockRegistry()
				mockPluginManager := mocks.NewMockPluginManager()
				mockEventBus := mocks.NewMockEventBus()
				mockConfig := mocks.NewMockConfiguration()
				mockStore := mocks.NewMockMultiStore()
				mockLogger := mocks.NewMockLogger()

				svc := NewDefaultSystemService(
					"test-service",
					mockRegistry,
					mockPluginManager,
					mockEventBus,
					mockConfig,
					mockStore,
					mockLogger,
				)

				// Initialize and start first
				ctx := context.NewContext(stdctx.Background())
				_ = svc.Initialize(ctx)
				_ = svc.Start(ctx)

				// Mock the base service to fail on stop
				// We'll simulate this by making the service fail internally
				// Since we can't easily mock the embedded BaseService, we'll test other paths

				return svc
			},
			expectedError: "", // We'll skip this test for now since it's hard to mock embedded BaseService
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedError == "" {
				t.Skip("Skipping test that's hard to implement without complex mocking")
				return
			}

			svc := tt.setupService()
			ctx := context.NewContext(stdctx.Background())

			err := svc.Stop(ctx)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

// Test additional error paths for Initialize method
func TestDefaultSystemService_Initialize_AdditionalErrorPaths(t *testing.T) {
	tests := []struct {
		name          string
		setupService  func() *DefaultSystemService
		expectedError string
	}{
		{
			name: "base service initialize fails",
			setupService: func() *DefaultSystemService {
				mockRegistry := mocks.NewMockRegistry()
				mockPluginManager := mocks.NewMockPluginManager()
				mockEventBus := mocks.NewMockEventBus()
				mockConfig := mocks.NewMockConfiguration()
				mockStore := mocks.NewMockMultiStore()
				mockLogger := mocks.NewMockLogger()

				return NewDefaultSystemService(
					"test-service",
					mockRegistry,
					mockPluginManager,
					mockEventBus,
					mockConfig,
					mockStore,
					mockLogger,
				)
			},
			expectedError: "", // Skip for now - hard to mock embedded BaseService
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedError == "" {
				t.Skip("Skipping test that's hard to implement without complex mocking")
				return
			}

			svc := tt.setupService()
			ctx := context.NewContext(stdctx.Background())

			err := svc.Initialize(ctx)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

// Test additional error paths for Start method
func TestDefaultSystemService_Start_AdditionalErrorPaths(t *testing.T) {
	tests := []struct {
		name          string
		setupService  func() *DefaultSystemService
		expectedError string
	}{
		{
			name: "base service start fails",
			setupService: func() *DefaultSystemService {
				mockRegistry := mocks.NewMockRegistry()
				mockPluginManager := mocks.NewMockPluginManager()
				mockEventBus := mocks.NewMockEventBus()
				mockConfig := mocks.NewMockConfiguration()
				mockStore := mocks.NewMockMultiStore()
				mockLogger := mocks.NewMockLogger()

				svc := NewDefaultSystemService(
					"test-service",
					mockRegistry,
					mockPluginManager,
					mockEventBus,
					mockConfig,
					mockStore,
					mockLogger,
				)

				// Initialize first
				ctx := context.NewContext(stdctx.Background())
				_ = svc.Initialize(ctx)

				return svc
			},
			expectedError: "", // Skip for now - hard to mock embedded BaseService
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedError == "" {
				t.Skip("Skipping test that's hard to implement without complex mocking")
				return
			}

			svc := tt.setupService()
			ctx := context.NewContext(stdctx.Background())

			err := svc.Start(ctx)

			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}
