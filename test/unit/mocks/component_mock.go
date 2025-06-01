package mocks

import (
	"fmt"
	"sync"

	"github.com/fintechain/skeleton/pkg/component"
	"github.com/fintechain/skeleton/pkg/context"
	"github.com/fintechain/skeleton/pkg/system"
)

// MockComponent provides a configurable mock implementation of the component.Component interface.
// It supports behavior configuration, error injection, call tracking, and state verification
// for comprehensive testing of components that depend on component functionality.
type MockComponent struct {
	mu sync.RWMutex

	// Component identity
	id          string
	name        string
	description string
	version     string
	compType    component.ComponentType
	metadata    component.Metadata

	// Behavior configuration
	shouldFail       bool
	failureError     string
	initializeError  string
	disposeError     string
	forceInitialized bool
	forceDisposed    bool

	// Call tracking
	callCount map[string]int
	lastCalls map[string][]interface{}

	// State verification
	initializeCalls []InitializeCall
	disposeCalls    int
	isInitialized   bool
	isDisposed      bool
}

// InitializeCall represents a call to Initialize method.
type InitializeCall struct {
	Context context.Context
	System  system.System
}

// NewMockComponent creates a new configurable component mock.
func NewMockComponent() *MockComponent {
	return &MockComponent{
		id:          "mock-component",
		name:        "Mock Component",
		description: "A configurable mock component for testing",
		version:     "1.0.0",
		compType:    component.TypeBasic,
		metadata:    make(component.Metadata),
		callCount:   make(map[string]int),
		lastCalls:   make(map[string][]interface{}),
	}
}

// Component Interface Implementation (registry.Identifiable)

// ID returns the component's unique identifier.
func (m *MockComponent) ID() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.trackCall("ID")
	return m.id
}

// Name returns the component's human-readable name.
func (m *MockComponent) Name() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.trackCall("Name")
	return m.name
}

// Description returns the component's description.
func (m *MockComponent) Description() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.trackCall("Description")
	return m.description
}

// Version returns the component's version.
func (m *MockComponent) Version() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.trackCall("Version")
	return m.version
}

// Component Interface Implementation (component.Component)

// Type returns the component's type.
func (m *MockComponent) Type() component.ComponentType {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.trackCall("Type")
	return m.compType
}

// Metadata returns the component's metadata.
func (m *MockComponent) Metadata() component.Metadata {
	m.mu.RLock()
	defer m.mu.RUnlock()
	m.trackCall("Metadata")

	// Return a copy to prevent external modification
	metadata := make(component.Metadata)
	for k, v := range m.metadata {
		metadata[k] = v
	}
	return metadata
}

// Initialize initializes the component with the given context and system.
func (m *MockComponent) Initialize(ctx context.Context, sys system.System) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("Initialize", ctx, sys)
	m.initializeCalls = append(m.initializeCalls, InitializeCall{
		Context: ctx,
		System:  sys,
	})

	if m.shouldFail {
		return fmt.Errorf("%s", m.getFailureError("Initialize"))
	}

	if m.initializeError != "" {
		return fmt.Errorf("%s", m.initializeError)
	}

	m.isInitialized = true
	return nil
}

// Dispose disposes of the component and cleans up resources.
func (m *MockComponent) Dispose() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("Dispose")
	m.disposeCalls++

	if m.shouldFail {
		return fmt.Errorf("%s", m.getFailureError("Dispose"))
	}

	if m.disposeError != "" {
		return fmt.Errorf("%s", m.disposeError)
	}

	m.isDisposed = true
	return nil
}

// Mock Configuration Methods

// SetID sets the component's ID.
func (m *MockComponent) SetID(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.id = id
}

// SetName sets the component's name.
func (m *MockComponent) SetName(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.name = name
}

// SetDescription sets the component's description.
func (m *MockComponent) SetDescription(description string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.description = description
}

// SetVersion sets the component's version.
func (m *MockComponent) SetVersion(version string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.version = version
}

// SetType sets the component's type.
func (m *MockComponent) SetType(compType component.ComponentType) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.compType = compType
}

// SetMetadata sets the component's metadata.
func (m *MockComponent) SetMetadata(metadata component.Metadata) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.metadata = make(component.Metadata)
	for k, v := range metadata {
		m.metadata[k] = v
	}
}

// AddMetadata adds a key-value pair to the component's metadata.
func (m *MockComponent) AddMetadata(key string, value interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.metadata == nil {
		m.metadata = make(component.Metadata)
	}
	m.metadata[key] = value
}

// SetShouldFail configures the mock to fail all operations.
func (m *MockComponent) SetShouldFail(fail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail = fail
}

// SetFailureError sets the error message for failed operations.
func (m *MockComponent) SetFailureError(err string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failureError = err
}

// SetInitializeError sets a specific error for Initialize method.
func (m *MockComponent) SetInitializeError(err string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.initializeError = err
}

// SetDisposeError sets a specific error for Dispose method.
func (m *MockComponent) SetDisposeError(err string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.disposeError = err
}

// SetForceInitialized forces the component to appear initialized.
func (m *MockComponent) SetForceInitialized(initialized bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.forceInitialized = initialized
	m.isInitialized = initialized
}

// SetForceDisposed forces the component to appear disposed.
func (m *MockComponent) SetForceDisposed(disposed bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.forceDisposed = disposed
	m.isDisposed = disposed
}

// State Verification Methods

// GetCallCount returns the number of times a method was called.
func (m *MockComponent) GetCallCount(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.callCount[method]
}

// GetLastCall returns the parameters of the last call to a method.
func (m *MockComponent) GetLastCall(method string) []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastCalls[method]
}

// GetInitializeCalls returns all calls to Initialize method.
func (m *MockComponent) GetInitializeCalls() []InitializeCall {
	m.mu.RLock()
	defer m.mu.RUnlock()
	calls := make([]InitializeCall, len(m.initializeCalls))
	copy(calls, m.initializeCalls)
	return calls
}

// GetDisposeCalls returns the number of times Dispose was called.
func (m *MockComponent) GetDisposeCalls() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.disposeCalls
}

// WasInitializeCalled returns true if Initialize was called.
func (m *MockComponent) WasInitializeCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.initializeCalls) > 0
}

// WasDisposeCalled returns true if Dispose was called.
func (m *MockComponent) WasDisposeCalled() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.disposeCalls > 0
}

// IsInitialized returns true if the component is initialized.
func (m *MockComponent) IsInitialized() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isInitialized || m.forceInitialized
}

// IsDisposed returns true if the component is disposed.
func (m *MockComponent) IsDisposed() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.isDisposed || m.forceDisposed
}

// Reset clears all mock state and configuration.
func (m *MockComponent) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.id = "mock-component"
	m.name = "Mock Component"
	m.description = "A configurable mock component for testing"
	m.version = "1.0.0"
	m.compType = component.TypeBasic
	m.metadata = make(component.Metadata)
	m.shouldFail = false
	m.failureError = ""
	m.initializeError = ""
	m.disposeError = ""
	m.forceInitialized = false
	m.forceDisposed = false
	m.callCount = make(map[string]int)
	m.lastCalls = make(map[string][]interface{})
	m.initializeCalls = nil
	m.disposeCalls = 0
	m.isInitialized = false
	m.isDisposed = false
}

// Helper Methods

// trackCall records a method call for verification.
func (m *MockComponent) trackCall(method string, args ...interface{}) {
	m.callCount[method]++
	m.lastCalls[method] = args
}

// getFailureError returns the configured failure error or a default.
func (m *MockComponent) getFailureError(method string) string {
	if m.failureError != "" {
		return m.failureError
	}
	return fmt.Sprintf("mock_component.%s_failed", method)
}

// ComponentMockBuilder provides a fluent interface for configuring component mocks.
type ComponentMockBuilder struct {
	mock *MockComponent
}

// NewComponentMockBuilder creates a new component mock builder.
func NewComponentMockBuilder() *ComponentMockBuilder {
	return &ComponentMockBuilder{
		mock: NewMockComponent(),
	}
}

// WithID sets the component's ID.
func (b *ComponentMockBuilder) WithID(id string) *ComponentMockBuilder {
	b.mock.SetID(id)
	return b
}

// WithName sets the component's name.
func (b *ComponentMockBuilder) WithName(name string) *ComponentMockBuilder {
	b.mock.SetName(name)
	return b
}

// WithDescription sets the component's description.
func (b *ComponentMockBuilder) WithDescription(description string) *ComponentMockBuilder {
	b.mock.SetDescription(description)
	return b
}

// WithVersion sets the component's version.
func (b *ComponentMockBuilder) WithVersion(version string) *ComponentMockBuilder {
	b.mock.SetVersion(version)
	return b
}

// WithType sets the component's type.
func (b *ComponentMockBuilder) WithType(compType component.ComponentType) *ComponentMockBuilder {
	b.mock.SetType(compType)
	return b
}

// WithMetadata sets the component's metadata.
func (b *ComponentMockBuilder) WithMetadata(metadata component.Metadata) *ComponentMockBuilder {
	b.mock.SetMetadata(metadata)
	return b
}

// WithMetadataEntry adds a metadata entry.
func (b *ComponentMockBuilder) WithMetadataEntry(key string, value interface{}) *ComponentMockBuilder {
	b.mock.AddMetadata(key, value)
	return b
}

// WithFailure configures the mock to fail operations.
func (b *ComponentMockBuilder) WithFailure(fail bool) *ComponentMockBuilder {
	b.mock.SetShouldFail(fail)
	return b
}

// WithFailureError sets the error message for failed operations.
func (b *ComponentMockBuilder) WithFailureError(err string) *ComponentMockBuilder {
	b.mock.SetFailureError(err)
	return b
}

// WithInitializeError sets a specific error for Initialize method.
func (b *ComponentMockBuilder) WithInitializeError(err string) *ComponentMockBuilder {
	b.mock.SetInitializeError(err)
	return b
}

// WithDisposeError sets a specific error for Dispose method.
func (b *ComponentMockBuilder) WithDisposeError(err string) *ComponentMockBuilder {
	b.mock.SetDisposeError(err)
	return b
}

// Build returns the configured mock component.
func (b *ComponentMockBuilder) Build() component.Component {
	return b.mock
}
