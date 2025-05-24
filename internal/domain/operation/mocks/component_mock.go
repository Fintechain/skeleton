// Package mocks provides mock implementations for testing
package mocks

import (
	"github.com/ebanfa/skeleton/internal/domain/component"
)

// MockExternalComponent is a mock implementation of component.Component for testing
// It mocks external component dependencies, not components within the package
type MockExternalComponent struct {
	// Fields to control mock behavior
	IDFunc         func() string
	NameFunc       func() string
	TypeFunc       func() component.ComponentType
	MetadataFunc   func() component.Metadata
	InitializeFunc func(ctx component.Context) error
	DisposeFunc    func() error

	// Track method calls for verification
	IDCalls         int
	NameCalls       int
	TypeCalls       int
	MetadataCalls   int
	InitializeCalls []component.Context
	DisposeCalls    int
}

// ID returns the component ID
func (m *MockExternalComponent) ID() string {
	m.IDCalls++
	if m.IDFunc != nil {
		return m.IDFunc()
	}
	return "mock-id"
}

// Name returns the component name
func (m *MockExternalComponent) Name() string {
	m.NameCalls++
	if m.NameFunc != nil {
		return m.NameFunc()
	}
	return "Mock Component"
}

// Type returns the component type
func (m *MockExternalComponent) Type() component.ComponentType {
	m.TypeCalls++
	if m.TypeFunc != nil {
		return m.TypeFunc()
	}
	return component.TypeOperation
}

// Metadata returns the component metadata
func (m *MockExternalComponent) Metadata() component.Metadata {
	m.MetadataCalls++
	if m.MetadataFunc != nil {
		return m.MetadataFunc()
	}
	return make(component.Metadata)
}

// Initialize initializes the component
func (m *MockExternalComponent) Initialize(ctx component.Context) error {
	m.InitializeCalls = append(m.InitializeCalls, ctx)
	if m.InitializeFunc != nil {
		return m.InitializeFunc(ctx)
	}
	return nil
}

// Dispose releases component resources
func (m *MockExternalComponent) Dispose() error {
	m.DisposeCalls++
	if m.DisposeFunc != nil {
		return m.DisposeFunc()
	}
	return nil
}
