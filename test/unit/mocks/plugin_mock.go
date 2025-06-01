// Package mocks provides centralized mock implementations for the Skeleton Framework.
// This file contains mocks for plugin interfaces.
package mocks

import (
	"sync"

	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/registry"
)

// MockPlugin implements the plugin.Plugin interface for testing.
type MockPlugin struct {
	mu sync.RWMutex

	// Identifiable fields
	id          string
	name        string
	description string
	version     string

	// Plugin-specific fields
	loadFunc     func(ctx context.Context, registrar registry.Registry) error
	unloadFunc   func(ctx context.Context) error
	loadError    error
	unloadError  error
	loadCalled   bool
	unloadCalled bool
	isLoaded     bool
}

// NewMockPlugin creates a new mock plugin.
func NewMockPlugin() *MockPlugin {
	return &MockPlugin{
		id:          "mock-plugin",
		name:        "Mock Plugin",
		description: "A mock plugin for testing",
		version:     "1.0.0",
	}
}

// Identifiable interface implementation

// ID implements the registry.Identifiable interface.
func (m *MockPlugin) ID() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.id
}

// Name implements the registry.Identifiable interface.
func (m *MockPlugin) Name() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.name
}

// Description implements the registry.Identifiable interface.
func (m *MockPlugin) Description() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.description
}

// Version implements the registry.Identifiable interface.
func (m *MockPlugin) Version() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.version
}

// Plugin interface implementation

// Load implements the plugin.Plugin interface.
func (m *MockPlugin) Load(ctx context.Context, registrar registry.Registry) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.loadCalled = true

	if m.loadFunc != nil {
		err := m.loadFunc(ctx, registrar)
		if err == nil {
			m.isLoaded = true
		}
		return err
	}

	if m.loadError != nil {
		return m.loadError
	}

	m.isLoaded = true
	return nil
}

// Unload implements the plugin.Plugin interface.
func (m *MockPlugin) Unload(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.unloadCalled = true

	if m.unloadFunc != nil {
		err := m.unloadFunc(ctx)
		if err == nil {
			m.isLoaded = false
		}
		return err
	}

	if m.unloadError != nil {
		return m.unloadError
	}

	m.isLoaded = false
	return nil
}

// Mock configuration methods

// SetID sets the plugin ID.
func (m *MockPlugin) SetID(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.id = id
}

// SetName sets the plugin name.
func (m *MockPlugin) SetName(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.name = name
}

// SetDescription sets the plugin description.
func (m *MockPlugin) SetDescription(description string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.description = description
}

// SetVersion sets the plugin version.
func (m *MockPlugin) SetVersion(version string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.version = version
}

// SetLoadFunc sets a custom function for Load.
func (m *MockPlugin) SetLoadFunc(fn func(ctx context.Context, registrar registry.Registry) error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.loadFunc = fn
}

// SetUnloadFunc sets a custom function for Unload.
func (m *MockPlugin) SetUnloadFunc(fn func(ctx context.Context) error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.unloadFunc = fn
}

// SetLoadError sets the error to return from Load.
func (m *MockPlugin) SetLoadError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.loadError = err
}

// SetUnloadError sets the error to return from Unload.
func (m *MockPlugin) SetUnloadError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.unloadError = err
}

// Verification methods

// WasLoadCalled returns true if Load was called.
func (m *MockPlugin) WasLoadCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.loadCalled
}

// WasUnloadCalled returns true if Unload was called.
func (m *MockPlugin) WasUnloadCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.unloadCalled
}

// IsLoaded returns true if the plugin is currently loaded.
func (m *MockPlugin) IsLoaded() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isLoaded
}
