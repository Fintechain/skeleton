// Package mocks provides centralized mock implementations for the Skeleton Framework.
// This file contains mocks for component factory interfaces.
package mocks

import (
	"fmt"
	"sync"

	"github.com/fintechain/skeleton/internal/domain/component"
)

// MockComponentFactory implements the component.Factory interface for testing.
type MockComponentFactory struct {
	mu sync.RWMutex

	// Factory fields
	createFunc   func(config component.ComponentConfig) (component.Component, error)
	createResult component.Component
	createError  error
	createCalled bool
	createConfig component.ComponentConfig
}

// NewMockComponentFactory creates a new mock component factory.
func NewMockComponentFactory() *MockComponentFactory {
	return &MockComponentFactory{}
}

// Create implements the component.Factory interface.
func (m *MockComponentFactory) Create(config component.ComponentConfig) (component.Component, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.createCalled = true
	m.createConfig = config

	if m.createFunc != nil {
		return m.createFunc(config)
	}

	if m.createError != nil {
		return nil, m.createError
	}

	// Add validation to match real component factory behavior
	if config.ID == "" {
		return nil, fmt.Errorf("component ID cannot be empty")
	}
	if config.Name == "" {
		return nil, fmt.Errorf("component name cannot be empty")
	}

	if m.createResult != nil {
		return m.createResult, nil
	}

	// Default behavior: create a mock component
	mockComp := NewMockComponent()
	mockComp.SetID(config.ID)
	mockComp.SetName(config.Name)
	mockComp.SetDescription(config.Description)
	mockComp.SetType(component.ComponentType(config.Type))
	return mockComp, nil
}

// Mock configuration methods

// SetCreateFunc sets a custom function for Create.
func (m *MockComponentFactory) SetCreateFunc(fn func(config component.ComponentConfig) (component.Component, error)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createFunc = fn
}

// SetCreateResult sets the result to return from Create.
func (m *MockComponentFactory) SetCreateResult(result component.Component) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createResult = result
}

// SetCreateError sets the error to return from Create.
func (m *MockComponentFactory) SetCreateError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.createError = err
}

// Verification methods

// WasCreateCalled returns true if Create was called.
func (m *MockComponentFactory) WasCreateCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.createCalled
}

// GetCreateConfig returns the config passed to the last Create call.
func (m *MockComponentFactory) GetCreateConfig() component.ComponentConfig {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.createConfig
}
