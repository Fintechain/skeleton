package mocks

import (
	"github.com/fintechain/skeleton/internal/domain/component"
)

// MockComponent is a mock implementation of component.Component
type MockComponent struct {
	IDValue       string
	NameValue     string
	TypeValue     component.ComponentType
	MetadataValue component.Metadata

	// Call tracking
	InitializeCalls []component.Context
	DisposeCalls    int

	// Function overrides
	InitializeFunc func(component.Context) error
	DisposeFunc    func() error
}

// NewMockComponent creates a new mock component
func NewMockComponent(id, name string, componentType component.ComponentType) *MockComponent {
	return &MockComponent{
		IDValue:       id,
		NameValue:     name,
		TypeValue:     componentType,
		MetadataValue: make(component.Metadata),
	}
}

// ID returns the component ID
func (m *MockComponent) ID() string {
	return m.IDValue
}

// Name returns the component name
func (m *MockComponent) Name() string {
	return m.NameValue
}

// Type returns the component type
func (m *MockComponent) Type() component.ComponentType {
	return m.TypeValue
}

// Metadata returns the component metadata
func (m *MockComponent) Metadata() component.Metadata {
	return m.MetadataValue
}

// Initialize initializes the component
func (m *MockComponent) Initialize(ctx component.Context) error {
	m.InitializeCalls = append(m.InitializeCalls, ctx)
	if m.InitializeFunc != nil {
		return m.InitializeFunc(ctx)
	}
	return nil
}

// Dispose disposes the component
func (m *MockComponent) Dispose() error {
	m.DisposeCalls++
	if m.DisposeFunc != nil {
		return m.DisposeFunc()
	}
	return nil
}
