package mocks

import (
	"github.com/fintechain/skeleton/internal/domain/component"
)

// MockRegistry is a mock implementation of component.Registry for testing
type MockRegistry struct {
	// Function fields for customizing behavior
	RegisterFunc       func(component.Component) error
	UnregisterFunc     func(string) error
	GetFunc            func(string) (component.Component, error)
	FindByTypeFunc     func(component.ComponentType) []component.Component
	FindByMetadataFunc func(string, interface{}) []component.Component
	InitializeFunc     func(component.Context) error
	ShutdownFunc       func() error

	// Call tracking
	RegisterCalls       []component.Component
	UnregisterCalls     []string
	GetCalls            []string
	FindByTypeCalls     []component.ComponentType
	FindByMetadataCalls []FindByMetadataCall
	InitializeCalls     []component.Context
	ShutdownCalls       int

	// State
	Components map[string]component.Component
}

type FindByMetadataCall struct {
	Key   string
	Value interface{}
}

// NewMockRegistry creates a new mock registry
func NewMockRegistry() *MockRegistry {
	return &MockRegistry{
		Components: make(map[string]component.Component),
	}
}

// Register implements component.Registry
func (m *MockRegistry) Register(comp component.Component) error {
	m.RegisterCalls = append(m.RegisterCalls, comp)
	if m.RegisterFunc != nil {
		return m.RegisterFunc(comp)
	}
	if m.Components != nil {
		m.Components[comp.ID()] = comp
	}
	return nil
}

// Unregister implements component.Registry
func (m *MockRegistry) Unregister(id string) error {
	m.UnregisterCalls = append(m.UnregisterCalls, id)
	if m.UnregisterFunc != nil {
		return m.UnregisterFunc(id)
	}
	if m.Components != nil {
		delete(m.Components, id)
	}
	return nil
}

// Get implements component.Registry
func (m *MockRegistry) Get(id string) (component.Component, error) {
	m.GetCalls = append(m.GetCalls, id)
	if m.GetFunc != nil {
		return m.GetFunc(id)
	}
	if m.Components != nil {
		if comp, exists := m.Components[id]; exists {
			return comp, nil
		}
	}
	return nil, component.NewError(component.ErrComponentNotFound, "component not found", nil)
}

// FindByType implements component.Registry
func (m *MockRegistry) FindByType(componentType component.ComponentType) []component.Component {
	m.FindByTypeCalls = append(m.FindByTypeCalls, componentType)
	if m.FindByTypeFunc != nil {
		return m.FindByTypeFunc(componentType)
	}
	var result []component.Component
	if m.Components != nil {
		for _, comp := range m.Components {
			if comp.Type() == componentType {
				result = append(result, comp)
			}
		}
	}
	return result
}

// FindByMetadata implements component.Registry
func (m *MockRegistry) FindByMetadata(key string, value interface{}) []component.Component {
	m.FindByMetadataCalls = append(m.FindByMetadataCalls, FindByMetadataCall{Key: key, Value: value})
	if m.FindByMetadataFunc != nil {
		return m.FindByMetadataFunc(key, value)
	}
	var result []component.Component
	if m.Components != nil {
		for _, comp := range m.Components {
			metadata := comp.Metadata()
			if val, exists := metadata[key]; exists && val == value {
				result = append(result, comp)
			}
		}
	}
	return result
}

// Initialize implements component.Registry
func (m *MockRegistry) Initialize(ctx component.Context) error {
	m.InitializeCalls = append(m.InitializeCalls, ctx)
	if m.InitializeFunc != nil {
		return m.InitializeFunc(ctx)
	}
	return nil
}

// Shutdown implements component.Registry
func (m *MockRegistry) Shutdown() error {
	m.ShutdownCalls++
	if m.ShutdownFunc != nil {
		return m.ShutdownFunc()
	}
	return nil
}
