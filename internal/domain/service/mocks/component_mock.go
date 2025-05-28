package mocks

import (
	"github.com/fintechain/skeleton/internal/domain/component"
)

// MockComponent implements the component.Component interface for testing
type MockComponent struct {
	// Fields to control mock behavior
	IDFunc         func() string
	NameFunc       func() string
	TypeFunc       func() component.ComponentType
	MetadataFunc   func() component.Metadata
	InitializeFunc func(ctx component.Context) error
	DisposeFunc    func() error

	// Fields to track method calls
	IDCalls         int
	NameCalls       int
	TypeCalls       int
	MetadataCalls   int
	InitializeCalls []component.Context
	DisposeCalls    int
}

// ID returns the component's ID
func (m *MockComponent) ID() string {
	m.IDCalls++
	if m.IDFunc != nil {
		return m.IDFunc()
	}
	return "mock-component-id"
}

// Name returns the component's name
func (m *MockComponent) Name() string {
	m.NameCalls++
	if m.NameFunc != nil {
		return m.NameFunc()
	}
	return "Mock Component"
}

// Type returns the component's type
func (m *MockComponent) Type() component.ComponentType {
	m.TypeCalls++
	if m.TypeFunc != nil {
		return m.TypeFunc()
	}
	return component.TypeBasic
}

// Metadata returns the component's metadata
func (m *MockComponent) Metadata() component.Metadata {
	m.MetadataCalls++
	if m.MetadataFunc != nil {
		return m.MetadataFunc()
	}
	return component.Metadata{
		"mock": true,
	}
}

// Initialize initializes the component
func (m *MockComponent) Initialize(ctx component.Context) error {
	m.InitializeCalls = append(m.InitializeCalls, ctx)
	if m.InitializeFunc != nil {
		return m.InitializeFunc(ctx)
	}
	return nil
}

// Dispose releases resources used by the component
func (m *MockComponent) Dispose() error {
	m.DisposeCalls++
	if m.DisposeFunc != nil {
		return m.DisposeFunc()
	}
	return nil
}

// NewMockComponent creates a new mock component with default behavior
func NewMockComponent() *MockComponent {
	return &MockComponent{}
}
