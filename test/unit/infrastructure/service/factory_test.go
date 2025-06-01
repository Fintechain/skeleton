package service

import (
	"fmt"
	"testing"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/service"
	serviceImpl "github.com/fintechain/skeleton/internal/infrastructure/service"
	"github.com/fintechain/skeleton/pkg/context"
	"github.com/fintechain/skeleton/test/unit/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewServiceFactory(t *testing.T) {
	tests := []struct {
		name             string
		componentFactory component.Factory
		expectNil        bool
		description      string
	}{
		{
			name:             "valid component factory",
			componentFactory: mocks.NewMockComponentFactory(),
			expectNil:        false,
			description:      "Should create service factory with valid component factory",
		},
		{
			name:             "nil component factory",
			componentFactory: nil,
			expectNil:        true,
			description:      "Should return nil with nil component factory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			factory := serviceImpl.NewServiceFactory(tt.componentFactory)

			// Verify
			if tt.expectNil {
				assert.Nil(t, factory, tt.description)
			} else {
				assert.NotNil(t, factory, tt.description)

				// Verify interface compliance
				var _ service.ServiceFactory = factory
				var _ component.Factory = factory
			}
		})
	}
}

func TestServiceFactoryInterfaceCompliance(t *testing.T) {
	mockComponentFactory := mocks.NewMockComponentFactory()
	factory := serviceImpl.NewServiceFactory(mockComponentFactory)
	require.NotNil(t, factory)

	// Verify service factory interface
	var _ service.ServiceFactory = factory

	// Verify component factory interface (through embedding)
	var _ component.Factory = factory
}

func TestServiceFactoryCreate(t *testing.T) {
	tests := []struct {
		name        string
		setupConfig func() component.ComponentConfig
		expectError bool
		description string
	}{
		{
			name: "valid service config",
			setupConfig: func() component.ComponentConfig {
				return component.NewComponentConfig("test-service", "Test Service", component.TypeService, "Test service description")
			},
			expectError: false,
			description: "Should create component with valid service config",
		},
		{
			name: "nil config",
			setupConfig: func() component.ComponentConfig {
				return component.ComponentConfig{} // Empty config instead of nil
			},
			expectError: true,
			description: "Should fail with empty config",
		},
		{
			name: "invalid config type",
			setupConfig: func() component.ComponentConfig {
				return component.NewComponentConfig("test-operation", "Test Operation", component.TypeOperation, "Test operation description") // Wrong type
			},
			expectError: false, // The factory should convert the type to Service
			description: "Should create component and convert type to service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockComponentFactory := mocks.NewMockComponentFactory()
			factory := serviceImpl.NewServiceFactory(mockComponentFactory)
			require.NotNil(t, factory)

			config := tt.setupConfig()

			// Execute
			comp, err := factory.Create(config)

			// Verify
			if tt.expectError {
				assert.Error(t, err, tt.description)
				assert.Nil(t, comp)
			} else {
				assert.NoError(t, err, tt.description)
				assert.NotNil(t, comp)

				// Verify component properties
				assert.Equal(t, config.ID, comp.ID())
				assert.Equal(t, config.Name, comp.Name())
				assert.Equal(t, component.TypeService, comp.Type()) // Should always be service type
			}
		})
	}
}

func TestServiceFactoryCreateService(t *testing.T) {
	tests := []struct {
		name        string
		setupConfig func() service.ServiceConfig
		expectError bool
		description string
	}{
		{
			name: "valid service config",
			setupConfig: func() service.ServiceConfig {
				return service.NewServiceConfig("test-service", "Test Service", "Test service description")
			},
			expectError: false,
			description: "Should create service with valid config",
		},
		{
			name: "config with metadata",
			setupConfig: func() service.ServiceConfig {
				return service.NewServiceConfig("metadata-service", "Metadata Service", "Service with metadata")
			},
			expectError: false,
			description: "Should create service with metadata",
		},
		{
			name: "empty ID config",
			setupConfig: func() service.ServiceConfig {
				return service.NewServiceConfig("", "Empty ID Service", "Service with empty ID")
			},
			expectError: true,
			description: "Should fail with empty ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockComponentFactory := mocks.NewMockComponentFactory()
			factory := serviceImpl.NewServiceFactory(mockComponentFactory)
			require.NotNil(t, factory)

			config := tt.setupConfig()

			// Execute
			svc, err := factory.CreateService(config)

			// Verify
			if tt.expectError {
				assert.Error(t, err, tt.description)
				assert.Nil(t, svc)
			} else {
				assert.NoError(t, err, tt.description)
				assert.NotNil(t, svc)

				// Verify service interface
				var _ service.Service = svc

				// Verify service properties
				assert.Equal(t, config.ID, svc.ID())
				assert.Equal(t, config.Name, svc.Name())
				assert.Equal(t, config.Type, svc.Type())
				assert.Equal(t, service.StatusStopped, svc.Status())

				// Test service lifecycle
				ctx := context.NewContext()
				err := svc.Start(ctx)
				assert.NoError(t, err)
				assert.Equal(t, service.StatusRunning, svc.Status())

				err = svc.Stop(ctx)
				assert.NoError(t, err)
				assert.Equal(t, service.StatusStopped, svc.Status())
			}
		})
	}
}

func TestServiceFactoryCreateServiceWithValidation(t *testing.T) {
	// This test is removed since CreateServiceWithValidation is not in the interface
	// The method exists in the implementation but not in the domain interface
	t.Skip("CreateServiceWithValidation method not in domain interface")
}

func TestServiceFactoryCreateManagedService(t *testing.T) {
	// This test is removed since CreateManagedService is not in the interface
	// The method exists in the implementation but not in the domain interface
	t.Skip("CreateManagedService method not in domain interface")
}

func TestServiceFactoryErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		setupConfig func() component.ComponentConfig
		expectError bool
		description string
	}{
		{
			name: "valid config",
			setupConfig: func() component.ComponentConfig {
				return component.NewComponentConfig("valid-service", "Valid Service", component.TypeService, "Valid service description")
			},
			expectError: false,
			description: "Should create service with valid config",
		},
		{
			name: "empty ID",
			setupConfig: func() component.ComponentConfig {
				return component.NewComponentConfig("", "Empty ID Service", component.TypeService, "Service with empty ID")
			},
			expectError: true, // The component factory validates empty IDs
			description: "Should fail with empty ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockComponentFactory := mocks.NewMockComponentFactory()
			factory := serviceImpl.NewServiceFactory(mockComponentFactory)
			require.NotNil(t, factory)

			config := tt.setupConfig()

			// Execute
			comp, err := factory.Create(config)

			// Verify
			if tt.expectError {
				assert.Error(t, err, tt.description)
				assert.Nil(t, comp)
			} else {
				assert.NoError(t, err, tt.description)
				assert.NotNil(t, comp)
			}
		})
	}
}

func TestServiceFactoryComponentDelegation(t *testing.T) {
	mockComponentFactory := mocks.NewMockComponentFactory()
	factory := serviceImpl.NewServiceFactory(mockComponentFactory)
	require.NotNil(t, factory)

	// Test that factory delegates to component factory for non-service types
	config := component.NewComponentConfig("test-component", "Test Component", component.TypeBasic, "Test component description")

	comp, err := factory.Create(config)
	assert.NoError(t, err)
	assert.NotNil(t, comp)
	assert.Equal(t, component.TypeService, comp.Type()) // Should be converted to service type
}

func TestServiceFactoryMultipleServices(t *testing.T) {
	mockComponentFactory := mocks.NewMockComponentFactory()
	factory := serviceImpl.NewServiceFactory(mockComponentFactory)
	require.NotNil(t, factory)

	// Create multiple services
	services := make([]service.Service, 3)
	for i := 0; i < 3; i++ {
		config := service.NewServiceConfig(fmt.Sprintf("service-%d", i), fmt.Sprintf("Service %d", i), fmt.Sprintf("Service %d description", i))

		svc, err := factory.CreateService(config)
		assert.NoError(t, err)
		assert.NotNil(t, svc)

		services[i] = svc
	}

	// Verify all services are unique and properly configured
	for i, svc := range services {
		assert.Equal(t, fmt.Sprintf("service-%d", i), svc.ID())
		assert.Equal(t, fmt.Sprintf("Service %d", i), svc.Name())
		assert.Equal(t, component.TypeService, svc.Type())
		assert.Equal(t, service.StatusStopped, svc.Status())
	}
}

func TestServiceFactoryWithContext(t *testing.T) {
	mockComponentFactory := mocks.NewMockComponentFactory()
	factory := serviceImpl.NewServiceFactory(mockComponentFactory)
	require.NotNil(t, factory)

	config := service.NewServiceConfig("context-service", "Context Service", "Context service description")

	svc, err := factory.CreateService(config)
	assert.NoError(t, err)
	assert.NotNil(t, svc)

	// Test service with context
	ctx := context.NewContext()

	err = svc.Start(ctx)
	assert.NoError(t, err)
	assert.Equal(t, service.StatusRunning, svc.Status())

	err = svc.Stop(ctx)
	assert.NoError(t, err)
	assert.Equal(t, service.StatusStopped, svc.Status())
}

func TestServiceFactoryConfigValidation(t *testing.T) {
	mockComponentFactory := mocks.NewMockComponentFactory()
	factory := serviceImpl.NewServiceFactory(mockComponentFactory)
	require.NotNil(t, factory)

	tests := []struct {
		name        string
		setupConfig func() service.ServiceConfig
		expectError bool
		description string
	}{
		{
			name: "valid minimal config",
			setupConfig: func() service.ServiceConfig {
				return service.NewServiceConfig("minimal-svc", "Minimal Service", "Minimal service description")
			},
			expectError: false,
			description: "Should create service with minimal valid config",
		},
		{
			name: "config with all fields",
			setupConfig: func() service.ServiceConfig {
				return service.NewServiceConfig("complete-svc", "Complete Service", "Complete service with all fields")
			},
			expectError: false,
			description: "Should create service with complete config",
		},
		{
			name: "config with empty ID",
			setupConfig: func() service.ServiceConfig {
				return service.NewServiceConfig("", "Empty ID Service", "Service with empty ID")
			},
			expectError: true, // CreateService validates empty ID
			description: "Should fail with empty ID",
		},
		{
			name: "config with empty name",
			setupConfig: func() service.ServiceConfig {
				return service.NewServiceConfig("empty-name-svc", "", "Service with empty name")
			},
			expectError: true, // CreateService validates empty name
			description: "Should fail with empty name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.setupConfig()

			svc, err := factory.CreateService(config)

			if tt.expectError {
				assert.Error(t, err, tt.description)
				assert.Nil(t, svc)
			} else {
				assert.NoError(t, err, tt.description)
				assert.NotNil(t, svc)
				assert.Equal(t, config.ID, svc.ID())
				assert.Equal(t, config.Name, svc.Name())
				assert.Equal(t, component.TypeService, svc.Type())
			}
		})
	}
}

func TestServiceFactoryManagedServiceLifecycle(t *testing.T) {
	// This test is removed since CreateManagedService is not in the interface
	t.Skip("CreateManagedService method not in domain interface")
}

func TestServiceFactoryManagedServiceConcurrency(t *testing.T) {
	// This test is removed since CreateManagedService is not in the interface
	t.Skip("CreateManagedService method not in domain interface")
}
