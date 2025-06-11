package component

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fintechain/skeleton/internal/domain/component"
	infraComponent "github.com/fintechain/skeleton/internal/infrastructure/component"
	"github.com/fintechain/skeleton/test/unit/mocks"
)

// TestNewRegistry tests the constructor for Registry
func TestNewRegistry(t *testing.T) {
	registry := infraComponent.NewRegistry()
	assert.NotNil(t, registry)
	assert.Equal(t, 0, registry.Count())
}

// TestRegistryInterfaceCompliance verifies that Registry implements the Registry interface
func TestRegistryInterfaceCompliance(t *testing.T) {
	registry := infraComponent.NewRegistry()
	var _ component.Registry = registry
	assert.NotNil(t, registry)
}

// TestRegistryRegister tests component registration
func TestRegistryRegister(t *testing.T) {
	registry := infraComponent.NewRegistry()
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()

	// Configure mock
	mockComponent.On("ID").Return(component.ComponentID("test-component"))

	// Test registration
	err := registry.Register(mockComponent)
	assert.NoError(t, err)
	assert.Equal(t, 1, registry.Count())

	// Test duplicate registration
	err = registry.Register(mockComponent)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), component.ErrItemAlreadyExists)

	mockComponent.AssertExpectations(t)
}

// TestRegistryRegisterNilComponent tests registering nil component
func TestRegistryRegisterNilComponent(t *testing.T) {
	registry := infraComponent.NewRegistry()

	err := registry.Register(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), component.ErrInvalidItem)
}

// TestRegistryGet tests component retrieval
func TestRegistryGet(t *testing.T) {
	registry := infraComponent.NewRegistry()
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()

	// Configure mock
	componentID := component.ComponentID("test-component")
	mockComponent.On("ID").Return(componentID)

	// Register component
	err := registry.Register(mockComponent)
	assert.NoError(t, err)

	// Test retrieval
	retrieved, err := registry.Get(componentID)
	assert.NoError(t, err)
	assert.Equal(t, mockComponent, retrieved)

	// Test non-existent component
	_, err = registry.Get("non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), component.ErrItemNotFound)

	mockComponent.AssertExpectations(t)
}

// TestRegistryGetByType tests component retrieval by type
func TestRegistryGetByType(t *testing.T) {
	registry := infraComponent.NewRegistry()
	factory := mocks.NewFactory()

	// Create mock components
	mockComponent1 := factory.ComponentInterface()
	mockComponent2 := factory.ComponentInterface()
	mockComponent3 := factory.ComponentInterface()

	// Configure mocks
	mockComponent1.On("ID").Return(component.ComponentID("comp1"))
	mockComponent1.On("Type").Return(component.TypeService)
	mockComponent2.On("ID").Return(component.ComponentID("comp2"))
	mockComponent2.On("Type").Return(component.TypeService)
	mockComponent3.On("ID").Return(component.ComponentID("comp3"))
	mockComponent3.On("Type").Return(component.TypeOperation)

	// Register components
	registry.Register(mockComponent1)
	registry.Register(mockComponent2)
	registry.Register(mockComponent3)

	// Test retrieval by type
	services, err := registry.GetByType(component.TypeService)
	assert.NoError(t, err)
	assert.Len(t, services, 2)

	operations, err := registry.GetByType(component.TypeOperation)
	assert.NoError(t, err)
	assert.Len(t, operations, 1)

	// Verify mock expectations
	mockComponent1.AssertExpectations(t)
	mockComponent2.AssertExpectations(t)
	mockComponent3.AssertExpectations(t)
}

// TestRegistryHas tests component existence check
func TestRegistryHas(t *testing.T) {
	registry := infraComponent.NewRegistry()
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()

	// Configure mock
	componentID := component.ComponentID("test-component")
	mockComponent.On("ID").Return(componentID)

	// Test non-existent component
	assert.False(t, registry.Has(componentID))

	// Register component
	registry.Register(mockComponent)

	// Test existing component
	assert.True(t, registry.Has(componentID))

	mockComponent.AssertExpectations(t)
}

// TestRegistryList tests listing all component IDs
func TestRegistryList(t *testing.T) {
	registry := infraComponent.NewRegistry()
	factory := mocks.NewFactory()

	// Test empty registry
	ids := registry.List()
	assert.Empty(t, ids)

	// Add components
	mockComponent1 := factory.ComponentInterface()
	mockComponent2 := factory.ComponentInterface()

	mockComponent1.On("ID").Return(component.ComponentID("comp1"))
	mockComponent2.On("ID").Return(component.ComponentID("comp2"))

	registry.Register(mockComponent1)
	registry.Register(mockComponent2)

	// Test list
	ids = registry.List()
	assert.Len(t, ids, 2)
	assert.Contains(t, ids, component.ComponentID("comp1"))
	assert.Contains(t, ids, component.ComponentID("comp2"))

	mockComponent1.AssertExpectations(t)
	mockComponent2.AssertExpectations(t)
}

// TestRegistryUnregister tests component removal
func TestRegistryUnregister(t *testing.T) {
	registry := infraComponent.NewRegistry()
	factory := mocks.NewFactory()
	mockComponent := factory.ComponentInterface()

	// Configure mock
	componentID := component.ComponentID("test-component")
	mockComponent.On("ID").Return(componentID)

	// Test unregistering non-existent component
	err := registry.Unregister(componentID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), component.ErrItemNotFound)

	// Register component
	registry.Register(mockComponent)
	assert.Equal(t, 1, registry.Count())

	// Test unregistering existing component
	err = registry.Unregister(componentID)
	assert.NoError(t, err)
	assert.Equal(t, 0, registry.Count())

	mockComponent.AssertExpectations(t)
}

// TestRegistryClear tests clearing all components
func TestRegistryClear(t *testing.T) {
	registry := infraComponent.NewRegistry()
	factory := mocks.NewFactory()

	// Add components
	mockComponent1 := factory.ComponentInterface()
	mockComponent2 := factory.ComponentInterface()

	mockComponent1.On("ID").Return(component.ComponentID("comp1"))
	mockComponent2.On("ID").Return(component.ComponentID("comp2"))

	registry.Register(mockComponent1)
	registry.Register(mockComponent2)
	assert.Equal(t, 2, registry.Count())

	// Test clear
	err := registry.Clear()
	assert.NoError(t, err)
	assert.Equal(t, 0, registry.Count())

	mockComponent1.AssertExpectations(t)
	mockComponent2.AssertExpectations(t)
}

// TestRegistryFind tests finding components with predicate
func TestRegistryFind(t *testing.T) {
	registry := infraComponent.NewRegistry()
	factory := mocks.NewFactory()

	// Create mock components
	mockComponent1 := factory.ComponentInterface()
	mockComponent2 := factory.ComponentInterface()

	// Configure mocks
	mockComponent1.On("ID").Return(component.ComponentID("service-comp"))
	mockComponent1.On("Type").Return(component.TypeService)
	mockComponent2.On("ID").Return(component.ComponentID("operation-comp"))
	mockComponent2.On("Type").Return(component.TypeOperation)

	// Register components
	registry.Register(mockComponent1)
	registry.Register(mockComponent2)

	// Test find with predicate
	services, err := registry.Find(func(comp component.Component) bool {
		return comp.Type() == component.TypeService
	})
	assert.NoError(t, err)
	assert.Len(t, services, 1)

	// Test find with nil predicate
	_, err = registry.Find(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), component.ErrInvalidItem)

	mockComponent1.AssertExpectations(t)
	mockComponent2.AssertExpectations(t)
}
