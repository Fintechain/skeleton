package mocks

import (
	"github.com/ebanfa/skeleton/internal/domain/component"
)

// MockComponent implements the component.Component interface for testing
type MockComponent struct {
	// Function fields for controlling behavior
	IDFunc         func() string
	NameFunc       func() string
	TypeFunc       func() component.ComponentType
	MetadataFunc   func() component.Metadata
	InitializeFunc func(ctx component.Context) error
	DisposeFunc    func() error

	// Call tracking for verification
	IDCalls         int
	NameCalls       int
	TypeCalls       int
	MetadataCalls   int
	InitializeCalls []component.Context
	DisposeCalls    int
}

// ID mocks the ID method
func (m *MockComponent) ID() string {
	m.IDCalls++
	if m.IDFunc != nil {
		return m.IDFunc()
	}
	return "mock-component-id"
}

// Name mocks the Name method
func (m *MockComponent) Name() string {
	m.NameCalls++
	if m.NameFunc != nil {
		return m.NameFunc()
	}
	return "Mock Component"
}

// Type mocks the Type method
func (m *MockComponent) Type() component.ComponentType {
	m.TypeCalls++
	if m.TypeFunc != nil {
		return m.TypeFunc()
	}
	return component.TypeBasic
}

// Metadata mocks the Metadata method
func (m *MockComponent) Metadata() component.Metadata {
	m.MetadataCalls++
	if m.MetadataFunc != nil {
		return m.MetadataFunc()
	}
	return component.Metadata{}
}

// Initialize mocks the Initialize method
func (m *MockComponent) Initialize(ctx component.Context) error {
	m.InitializeCalls = append(m.InitializeCalls, ctx)
	if m.InitializeFunc != nil {
		return m.InitializeFunc(ctx)
	}
	return nil
}

// Dispose mocks the Dispose method
func (m *MockComponent) Dispose() error {
	m.DisposeCalls++
	if m.DisposeFunc != nil {
		return m.DisposeFunc()
	}
	return nil
}
