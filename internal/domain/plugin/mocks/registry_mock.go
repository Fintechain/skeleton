package mocks

import (
	"github.com/ebanfa/skeleton/internal/domain/component"
)

// MockRegistry implements the component.Registry interface for testing
type MockRegistry struct {
	// Function fields for controlling behavior
	RegisterFunc       func(component.Component) error
	UnregisterFunc     func(id string) error
	GetFunc            func(id string) (component.Component, error)
	FindByTypeFunc     func(componentType component.ComponentType) []component.Component
	FindByMetadataFunc func(key string, value interface{}) []component.Component
	InitializeFunc     func(ctx component.Context) error
	ShutdownFunc       func() error

	// Call tracking for verification
	RegisterCalls       []component.Component
	UnregisterCalls     []string
	GetCalls            []string
	FindByTypeCalls     []component.ComponentType
	FindByMetadataCalls []struct {
		Key   string
		Value interface{}
	}
	InitializeCalls int
	ShutdownCalls   int
}

// Register mocks the Register method
func (m *MockRegistry) Register(comp component.Component) error {
	m.RegisterCalls = append(m.RegisterCalls, comp)
	if m.RegisterFunc != nil {
		return m.RegisterFunc(comp)
	}
	return nil
}

// Unregister mocks the Unregister method
func (m *MockRegistry) Unregister(id string) error {
	m.UnregisterCalls = append(m.UnregisterCalls, id)
	if m.UnregisterFunc != nil {
		return m.UnregisterFunc(id)
	}
	return nil
}

// Get mocks the Get method
func (m *MockRegistry) Get(id string) (component.Component, error) {
	m.GetCalls = append(m.GetCalls, id)
	if m.GetFunc != nil {
		return m.GetFunc(id)
	}
	return nil, nil
}

// FindByType mocks the FindByType method
func (m *MockRegistry) FindByType(componentType component.ComponentType) []component.Component {
	m.FindByTypeCalls = append(m.FindByTypeCalls, componentType)
	if m.FindByTypeFunc != nil {
		return m.FindByTypeFunc(componentType)
	}
	return []component.Component{}
}

// FindByMetadata mocks the FindByMetadata method
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

// Initialize mocks the Initialize method
func (m *MockRegistry) Initialize(ctx component.Context) error {
	m.InitializeCalls++
	if m.InitializeFunc != nil {
		return m.InitializeFunc(ctx)
	}
	return nil
}

// Shutdown mocks the Shutdown method
func (m *MockRegistry) Shutdown() error {
	m.ShutdownCalls++
	if m.ShutdownFunc != nil {
		return m.ShutdownFunc()
	}
	return nil
}
