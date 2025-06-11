package component

import (
	"testing"

	"github.com/fintechain/skeleton/internal/domain/component"
	infraComponent "github.com/fintechain/skeleton/internal/infrastructure/component"
	infraContext "github.com/fintechain/skeleton/internal/infrastructure/context"
	"github.com/fintechain/skeleton/test/unit/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewBaseService(t *testing.T) {
	config := component.ComponentConfig{
		ID:          "test-service",
		Name:        "Test Service",
		Type:        component.TypeService,
		Description: "Test service description",
		Version:     "1.0.0",
	}

	service := infraComponent.NewBaseService(config)
	assert.NotNil(t, service)

	// Verify interface compliance
	var _ component.Service = service
	var _ component.Component = service

	// Test basic properties
	assert.Equal(t, component.ComponentID("test-service"), service.ID())
	assert.Equal(t, "Test Service", service.Name())
	assert.Equal(t, component.TypeService, service.Type())
	assert.Equal(t, "1.0.0", service.Version())
}

func TestBaseServiceInitialState(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "test-service",
		Name: "Test Service",
		Type: component.TypeService,
	}

	service := infraComponent.NewBaseService(config)

	// Test initial state
	assert.False(t, service.IsRunning())
	assert.Equal(t, component.StatusStopped, service.Status())
}

func TestBaseServiceLifecycle(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "test-service",
		Name: "Test Service",
		Type: component.TypeService,
	}

	service := infraComponent.NewBaseService(config)
	ctx := infraContext.NewContext()

	// Test start
	err := service.Start(ctx)
	assert.NoError(t, err)
	assert.True(t, service.IsRunning())
	assert.Equal(t, component.StatusRunning, service.Status())

	// Test stop
	err = service.Stop(ctx)
	assert.NoError(t, err)
	assert.False(t, service.IsRunning())
	assert.Equal(t, component.StatusStopped, service.Status())
}

func TestBaseServiceIdempotentOperations(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "test-service",
		Name: "Test Service",
		Type: component.TypeService,
	}

	service := infraComponent.NewBaseService(config)
	ctx := infraContext.NewContext()

	// Test multiple starts (should be idempotent)
	err := service.Start(ctx)
	assert.NoError(t, err)
	assert.True(t, service.IsRunning())

	err = service.Start(ctx)
	assert.NoError(t, err)
	assert.True(t, service.IsRunning())

	// Test multiple stops (should be idempotent)
	err = service.Stop(ctx)
	assert.NoError(t, err)
	assert.False(t, service.IsRunning())

	err = service.Stop(ctx)
	assert.NoError(t, err)
	assert.False(t, service.IsRunning())
}

func TestBaseServiceInitializeAndDispose(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "test-service",
		Name: "Test Service",
		Type: component.TypeService,
	}

	service := infraComponent.NewBaseService(config)
	factory := mocks.NewFactory()
	mockSystem := factory.SystemInterface()
	ctx := infraContext.NewContext()

	// Test initialization
	err := service.Initialize(ctx, mockSystem)
	assert.NoError(t, err)

	// Test disposal
	err = service.Dispose()
	assert.NoError(t, err)
}

func TestBaseServiceStatusTransitions(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "test-service",
		Name: "Test Service",
		Type: component.TypeService,
	}

	service := infraComponent.NewBaseService(config)
	ctx := infraContext.NewContext()

	// Initial state
	assert.Equal(t, component.StatusStopped, service.Status())
	assert.False(t, service.IsRunning())

	// Start service
	err := service.Start(ctx)
	assert.NoError(t, err)
	assert.Equal(t, component.StatusRunning, service.Status())
	assert.True(t, service.IsRunning())

	// Stop service
	err = service.Stop(ctx)
	assert.NoError(t, err)
	assert.Equal(t, component.StatusStopped, service.Status())
	assert.False(t, service.IsRunning())
}
