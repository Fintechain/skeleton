package mocks

import (
	"github.com/ebanfa/skeleton/internal/domain/component"
)

// MockRegistry implements the component.Registry interface for testing
type MockRegistry struct {
	// Fields to control mock behavior
	RegisterFunc       func(component.Component) error
	UnregisterFunc     func(id string) error
	GetFunc            func(id string) (component.Component, error)
	FindByTypeFunc     func(componentType component.ComponentType) []component.Component
	FindByMetadataFunc func(key string, value interface{}) []component.Component
	InitializeFunc     func(ctx component.Context) error
	ShutdownFunc       func() error

	// Fields to track method calls
	RegisterCalls       []component.Component
	UnregisterCalls     []string
	GetCalls            []string
	FindByTypeCalls     []component.ComponentType
	FindByMetadataCalls []struct {
		Key   string
		Value interface{}
	}
	InitializeCalls []component.Context
	ShutdownCalls   int
}

// Register adds a component to the registry
func (m *MockRegistry) Register(comp component.Component) error {
	m.RegisterCalls = append(m.RegisterCalls, comp)
	if m.RegisterFunc != nil {
		return m.RegisterFunc(comp)
	}
	return nil
}

// Unregister removes a component from the registry
func (m *MockRegistry) Unregister(id string) error {
	m.UnregisterCalls = append(m.UnregisterCalls, id)
	if m.UnregisterFunc != nil {
		return m.UnregisterFunc(id)
	}
	return nil
}

// Get retrieves a component by ID
func (m *MockRegistry) Get(id string) (component.Component, error) {
	m.GetCalls = append(m.GetCalls, id)
	if m.GetFunc != nil {
		return m.GetFunc(id)
	}
	return nil, component.NewError(component.ErrComponentNotFound, "component not found in mock", nil)
}

// FindByType finds all components of a specific type
func (m *MockRegistry) FindByType(componentType component.ComponentType) []component.Component {
	m.FindByTypeCalls = append(m.FindByTypeCalls, componentType)
	if m.FindByTypeFunc != nil {
		return m.FindByTypeFunc(componentType)
	}
	return []component.Component{}
}

// FindByMetadata finds all components with a specific metadata key-value pair
func (m *MockRegistry) FindByMetadata(key string, value interface{}) []component.Component {
	m.FindByMetadataCalls = append(m.FindByMetadataCalls, struct {
		Key   string
		Value interface{}
	}{key, value})
	if m.FindByMetadataFunc != nil {
		return m.FindByMetadataFunc(key, value)
	}
	return []component.Component{}
}

// Initialize initializes all registered components
func (m *MockRegistry) Initialize(ctx component.Context) error {
	m.InitializeCalls = append(m.InitializeCalls, ctx)
	if m.InitializeFunc != nil {
		return m.InitializeFunc(ctx)
	}
	return nil
}

// Shutdown disposes all registered components
func (m *MockRegistry) Shutdown() error {
	m.ShutdownCalls++
	if m.ShutdownFunc != nil {
		return m.ShutdownFunc()
	}
	return nil
}

// NewMockRegistry creates a new mock registry with default behavior
func NewMockRegistry() *MockRegistry {
	return &MockRegistry{
		RegisterCalls:   make([]component.Component, 0),
		UnregisterCalls: make([]string, 0),
		GetCalls:        make([]string, 0),
		FindByTypeCalls: make([]component.ComponentType, 0),
		FindByMetadataCalls: make([]struct {
			Key   string
			Value interface{}
		}, 0),
		InitializeCalls: make([]component.Context, 0),
	}
}
