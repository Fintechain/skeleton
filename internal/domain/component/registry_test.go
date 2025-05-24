package component

import (
	"testing"

	"github.com/ebanfa/skeleton/internal/domain/component/mocks"
)

func TestRegistryRegister(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create registry with mock dependencies
	registry := NewRegistry(DefaultRegistryOptions{
		Logger: mockLogger,
	})

	// Create a test component
	testComp := NewBaseComponent("test-id", "Test Component", TypeBasic)

	// Test registration
	err := registry.Register(testComp)
	if err != nil {
		t.Errorf("Failed to register component: %s", err)
	}

	// Test retrieval
	comp, err := registry.Get("test-id")
	if err != nil {
		t.Errorf("Failed to get component: %s", err)
	}

	if comp.ID() != "test-id" {
		t.Errorf("Retrieved component has wrong ID: got %s, expected test-id", comp.ID())
	}
}

func TestRegistryRegisterDuplicate(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create registry with mock dependencies
	registry := NewRegistry(DefaultRegistryOptions{
		Logger: mockLogger,
	})

	// Create a test component
	testComp := NewBaseComponent("test-id", "Test Component", TypeBasic)

	// Register component
	err := registry.Register(testComp)
	if err != nil {
		t.Errorf("Failed to register component: %s", err)
	}

	// Try to register again with the same ID
	err = registry.Register(testComp)
	if err == nil {
		t.Error("Expected error when registering duplicate component, but got nil")
	}

	// Check if error is the expected type
	if !IsComponentError(err, ErrComponentExists) {
		t.Errorf("Expected component exists error, but got: %s", err)
	}
}

func TestRegistryUnregister(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create registry with mock dependencies
	registry := NewRegistry(DefaultRegistryOptions{
		Logger: mockLogger,
	})

	// Create a test component
	testComp := NewBaseComponent("test-id", "Test Component", TypeBasic)

	// Register component
	err := registry.Register(testComp)
	if err != nil {
		t.Errorf("Failed to register component: %s", err)
	}

	// Unregister component
	err = registry.Unregister("test-id")
	if err != nil {
		t.Errorf("Failed to unregister component: %s", err)
	}

	// Try to get the unregistered component
	_, err = registry.Get("test-id")
	if err == nil {
		t.Error("Expected error getting unregistered component, but got nil")
	}

	// Check if error is the expected type
	if !IsComponentError(err, ErrComponentNotFound) {
		t.Errorf("Expected component not found error, but got: %s", err)
	}
}

func TestRegistryFindByType(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create registry with mock dependencies
	registry := NewRegistry(DefaultRegistryOptions{
		Logger: mockLogger,
	})

	// Create test components of different types
	comp1 := NewBaseComponent("test-id-1", "Test Component 1", TypeBasic)
	comp2 := NewBaseComponent("test-id-2", "Test Component 2", TypeService)
	comp3 := NewBaseComponent("test-id-3", "Test Component 3", TypeBasic)

	// Register components
	registry.Register(comp1)
	registry.Register(comp2)
	registry.Register(comp3)

	// Find components by type
	basicComps := registry.FindByType(TypeBasic)
	if len(basicComps) != 2 {
		t.Errorf("Expected 2 basic components, but got %d", len(basicComps))
	}

	serviceComps := registry.FindByType(TypeService)
	if len(serviceComps) != 1 {
		t.Errorf("Expected 1 service component, but got %d", len(serviceComps))
	}

	systemComps := registry.FindByType(TypeSystem)
	if len(systemComps) != 0 {
		t.Errorf("Expected 0 system components, but got %d", len(systemComps))
	}
}

func TestRegistryFindByMetadata(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create registry with mock dependencies
	registry := NewRegistry(DefaultRegistryOptions{
		Logger: mockLogger,
	})

	// Create test components with different metadata
	comp1 := NewBaseComponent("test-id-1", "Test Component 1", TypeBasic)
	comp1.SetMetadata("category", "database")

	comp2 := NewBaseComponent("test-id-2", "Test Component 2", TypeService)
	comp2.SetMetadata("category", "api")

	comp3 := NewBaseComponent("test-id-3", "Test Component 3", TypeBasic)
	comp3.SetMetadata("category", "database")

	// Register components
	registry.Register(comp1)
	registry.Register(comp2)
	registry.Register(comp3)

	// Find components by metadata
	databaseComps := registry.FindByMetadata("category", "database")
	if len(databaseComps) != 2 {
		t.Errorf("Expected 2 database components, but got %d", len(databaseComps))
	}

	apiComps := registry.FindByMetadata("category", "api")
	if len(apiComps) != 1 {
		t.Errorf("Expected 1 api component, but got %d", len(apiComps))
	}

	uiComps := registry.FindByMetadata("category", "ui")
	if len(uiComps) != 0 {
		t.Errorf("Expected 0 ui components, but got %d", len(uiComps))
	}
}

func TestRegistryCreateRegistry(t *testing.T) {
	// Test the factory method for backward compatibility
	registry := CreateRegistry()

	if registry == nil {
		t.Error("CreateRegistry returned nil")
	}

	// Verify we can use the registry normally
	comp := NewBaseComponent("test-id", "Test Component", TypeBasic)
	err := registry.Register(comp)
	if err != nil {
		t.Errorf("Failed to register component using CreateRegistry-created registry: %s", err)
	}
}

func TestRegistryFactoryMethods(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create registry with mock dependencies
	registry := NewRegistry(DefaultRegistryOptions{
		Logger: mockLogger,
	})

	// Create a test factory
	factory := CreateFactory()

	// Test RegisterFactory method
	err := registry.(*DefaultRegistry).RegisterFactory("factory-id", factory)
	if err != nil {
		t.Errorf("Failed to register factory: %s", err)
	}

	// Test GetFactory method
	retrievedFactory, err := registry.(*DefaultRegistry).GetFactory("factory-id")
	if err != nil {
		t.Errorf("Failed to get factory: %s", err)
	}

	if retrievedFactory != factory {
		t.Error("Retrieved factory is not the same as registered factory")
	}

	// Test duplicate factory registration
	err = registry.(*DefaultRegistry).RegisterFactory("factory-id", factory)
	if err == nil {
		t.Error("Expected error when registering duplicate factory, but got nil")
	}

	// Test getting non-existent factory
	_, err = registry.(*DefaultRegistry).GetFactory("non-existent")
	if err == nil {
		t.Error("Expected error when getting non-existent factory, but got nil")
	}
}

func TestRegistryInitializeAndShutdown(t *testing.T) {
	// Create mock dependencies
	mockLogger := mocks.NewMockLogger()

	// Create registry with mock dependencies
	registry := NewRegistry(DefaultRegistryOptions{
		Logger: mockLogger,
	})

	// Create test components
	comp1 := NewBaseComponent("test-id-1", "Test Component 1", TypeBasic)
	comp2 := NewBaseComponent("test-id-2", "Test Component 2", TypeService)

	// Register components
	registry.Register(comp1)
	registry.Register(comp2)

	// Create a context
	ctx := &mockContext{}

	// Test Initialize method
	err := registry.Initialize(ctx)
	if err != nil {
		t.Errorf("Failed to initialize registry: %s", err)
	}

	// Test Shutdown method
	err = registry.Shutdown()
	if err != nil {
		t.Errorf("Failed to shutdown registry: %s", err)
	}

	// Test initialization failure
	// Create a registry with a failing component
	failingRegistry := NewRegistry(DefaultRegistryOptions{
		Logger: mockLogger,
	})

	failingComponent := &mockFailingComponent{id: "failing-comp"}
	failingRegistry.Register(failingComponent)

	// Initialize should fail
	err = failingRegistry.Initialize(ctx)
	if err == nil {
		t.Error("Expected error initializing registry with failing component, but got nil")
	}

	// Check error type
	if !IsComponentError(err, ErrInitializationFailed) {
		t.Errorf("Expected initialization failed error, but got: %s", err)
	}

	// Test shutdown failure
	// Register a component that fails on dispose
	failingRegistry = NewRegistry(DefaultRegistryOptions{
		Logger: mockLogger,
	})

	// Register a regular component first (that won't fail)
	regularComp := NewBaseComponent("regular-comp", "Regular Component", TypeBasic)
	failingRegistry.Register(regularComp)

	// Now register a failing component
	failingRegistry.Register(failingComponent)

	// Shutdown should fail
	err = failingRegistry.Shutdown()
	if err == nil {
		t.Error("Expected error shutting down registry with failing component, but got nil")
	}

	// Check error type
	if !IsComponentError(err, ErrDisposalFailed) {
		t.Errorf("Expected disposal failed error, but got: %s", err)
	}
}
