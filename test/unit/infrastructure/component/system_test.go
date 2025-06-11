package component

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fintechain/skeleton/internal/domain/component"
	infraComponent "github.com/fintechain/skeleton/internal/infrastructure/component"
	"github.com/fintechain/skeleton/test/unit/mocks"
)

// TestNewSystem tests the constructor for System
func TestNewSystem(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()

	system := infraComponent.NewSystem(mockRegistry)
	assert.NotNil(t, system)
	assert.Equal(t, mockRegistry, system.Registry())
	assert.False(t, system.IsRunning())
}

// TestSystemInterfaceCompliance verifies that System implements the System interface
func TestSystemInterfaceCompliance(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()

	system := infraComponent.NewSystem(mockRegistry)
	var _ component.System = system
	assert.NotNil(t, system)
}

// TestSystemLifecycle tests system start and stop
func TestSystemLifecycle(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()

	system := infraComponent.NewSystem(mockRegistry)

	// Test initial state
	assert.False(t, system.IsRunning())

	// Test start
	err := system.Start(nil)
	assert.NoError(t, err)
	assert.True(t, system.IsRunning())

	// Test idempotent start
	err = system.Start(nil)
	assert.NoError(t, err)
	assert.True(t, system.IsRunning())

	// Test stop
	err = system.Stop(nil)
	assert.NoError(t, err)
	assert.False(t, system.IsRunning())

	// Test idempotent stop
	err = system.Stop(nil)
	assert.NoError(t, err)
	assert.False(t, system.IsRunning())
}

// TestSystemExecuteOperation tests operation execution
func TestSystemExecuteOperation(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()
	mockOperation := factory.OperationInterface()

	system := infraComponent.NewSystem(mockRegistry)

	// Configure mocks
	operationID := component.ComponentID("test-operation")
	input := component.Input{Data: "test"}
	expectedOutput := component.Output{Data: "result"}

	mockRegistry.On("Get", operationID).Return(mockOperation, nil)
	mockOperation.On("Execute", nil, input).Return(expectedOutput, nil)

	// Test operation execution
	output, err := system.ExecuteOperation(nil, operationID, input)
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output)

	// Verify mock expectations
	mockRegistry.AssertExpectations(t)
	mockOperation.AssertExpectations(t)
}

// TestSystemExecuteOperationNotFound tests operation execution with non-existent operation
func TestSystemExecuteOperationNotFound(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()

	system := infraComponent.NewSystem(mockRegistry)

	// Configure mock to return error
	operationID := component.ComponentID("non-existent")
	input := component.Input{Data: "test"}

	mockRegistry.On("Get", operationID).Return(nil, assert.AnError)

	// Test operation execution
	_, err := system.ExecuteOperation(nil, operationID, input)
	assert.Error(t, err)

	mockRegistry.AssertExpectations(t)
}

// TestSystemExecuteOperationInvalidType tests operation execution with wrong component type
func TestSystemExecuteOperationInvalidType(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()
	mockComponent := factory.ComponentInterface()

	system := infraComponent.NewSystem(mockRegistry)

	// Configure mocks
	operationID := component.ComponentID("not-operation")
	input := component.Input{Data: "test"}

	mockRegistry.On("Get", operationID).Return(mockComponent, nil)

	// Test operation execution with wrong type
	_, err := system.ExecuteOperation(nil, operationID, input)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), component.ErrInvalidComponentType)

	mockRegistry.AssertExpectations(t)
}

// TestSystemStartService tests service starting
func TestSystemStartService(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()
	mockService := factory.ServiceInterface()

	system := infraComponent.NewSystem(mockRegistry)

	// Configure mocks
	serviceID := component.ComponentID("test-service")

	mockRegistry.On("Get", serviceID).Return(mockService, nil)
	mockService.On("Start", nil).Return(nil)

	// Test service start
	err := system.StartService(nil, serviceID)
	assert.NoError(t, err)

	// Verify mock expectations
	mockRegistry.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

// TestSystemStartServiceNotFound tests service starting with non-existent service
func TestSystemStartServiceNotFound(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()

	system := infraComponent.NewSystem(mockRegistry)

	// Configure mock to return error
	serviceID := component.ComponentID("non-existent")

	mockRegistry.On("Get", serviceID).Return(nil, assert.AnError)

	// Test service start
	err := system.StartService(nil, serviceID)
	assert.Error(t, err)

	mockRegistry.AssertExpectations(t)
}

// TestSystemStartServiceInvalidType tests service starting with wrong component type
func TestSystemStartServiceInvalidType(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()
	mockComponent := factory.ComponentInterface()

	system := infraComponent.NewSystem(mockRegistry)

	// Configure mocks
	serviceID := component.ComponentID("not-service")

	mockRegistry.On("Get", serviceID).Return(mockComponent, nil)

	// Test service start with wrong type
	err := system.StartService(nil, serviceID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), component.ErrInvalidComponentType)

	mockRegistry.AssertExpectations(t)
}

// TestSystemStopService tests service stopping
func TestSystemStopService(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()
	mockService := factory.ServiceInterface()

	system := infraComponent.NewSystem(mockRegistry)

	// Configure mocks
	serviceID := component.ComponentID("test-service")

	mockRegistry.On("Get", serviceID).Return(mockService, nil)
	mockService.On("Stop", nil).Return(nil)

	// Test service stop
	err := system.StopService(nil, serviceID)
	assert.NoError(t, err)

	// Verify mock expectations
	mockRegistry.AssertExpectations(t)
	mockService.AssertExpectations(t)
}

// TestSystemStopServiceNotFound tests service stopping with non-existent service
func TestSystemStopServiceNotFound(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()

	system := infraComponent.NewSystem(mockRegistry)

	// Configure mock to return error
	serviceID := component.ComponentID("non-existent")

	mockRegistry.On("Get", serviceID).Return(nil, assert.AnError)

	// Test service stop
	err := system.StopService(nil, serviceID)
	assert.Error(t, err)

	mockRegistry.AssertExpectations(t)
}

// TestSystemStopServiceInvalidType tests service stopping with wrong component type
func TestSystemStopServiceInvalidType(t *testing.T) {
	factory := mocks.NewFactory()
	mockRegistry := factory.RegistryInterface()
	mockComponent := factory.ComponentInterface()

	system := infraComponent.NewSystem(mockRegistry)

	// Configure mocks
	serviceID := component.ComponentID("not-service")

	mockRegistry.On("Get", serviceID).Return(mockComponent, nil)

	// Test service stop with wrong type
	err := system.StopService(nil, serviceID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), component.ErrInvalidComponentType)

	mockRegistry.AssertExpectations(t)
}
