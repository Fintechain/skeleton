// Package component provides public APIs for the component system.
package component

import (
	stdctx "context"
	"time"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/infrastructure/context"
)

// ===== CORE COMPONENT INTERFACES =====

// Component is the fundamental building block of the system.
type Component = component.Component

// Context represents the execution context for components.
type Context = component.Context

// Registry manages component registration and discovery.
type Registry = component.Registry

// Factory creates components from configuration.
type Factory = component.Factory

// ===== DEPENDENCY INJECTION INTERFACES =====

// DependencyAware interface for components that manage dependencies
type DependencyAware = component.DependencyAware

// DependencyAwareComponent provides dependency management capabilities
type DependencyAwareComponent = component.DependencyAwareComponent

// DependencyAwareComponentOptions contains options for creating a DependencyAwareComponent
type DependencyAwareComponentOptions = component.DependencyAwareComponentOptions

// BaseComponent provides a basic component implementation
type BaseComponent = component.BaseComponent

// BaseComponentOptions contains options for creating a BaseComponent
type BaseComponentOptions = component.BaseComponentOptions

// DefaultFactory provides a configurable factory implementation
type DefaultFactory = component.DefaultFactory

// DefaultFactoryOptions contains options for creating a DefaultFactory
type DefaultFactoryOptions = component.DefaultFactoryOptions

// ===== COMPONENT TYPES =====

// ComponentType represents the type of a component.
type ComponentType = component.ComponentType

// Metadata is a map of key-value pairs for component metadata.
type Metadata = component.Metadata

// ComponentConfig defines the configuration for creating a component.
type ComponentConfig = component.ComponentConfig

// ===== COMPONENT TYPE CONSTANTS =====

// Component type constants
const (
	TypeBasic       = component.TypeBasic
	TypeOperation   = component.TypeOperation
	TypeService     = component.TypeService
	TypeSystem      = component.TypeSystem
	TypeApplication = component.TypeApplication
)

// ===== COMPONENT ERROR CONSTANTS =====

// Common component error codes
const (
	ErrComponentNotFound    = component.ErrComponentNotFound
	ErrComponentExists      = component.ErrComponentExists
	ErrInvalidComponent     = component.ErrInvalidComponent
	ErrComponentCreation    = component.ErrComponentCreation
	ErrDependencyNotFound   = component.ErrDependencyNotFound
	ErrInitializationFailed = component.ErrInitializationFailed
	ErrDisposalFailed       = component.ErrDisposalFailed
)

// ===== COMPONENT UTILITIES =====

// Error represents a domain-specific error from the component system.
type Error = component.Error

// NewError creates a new component error with the given code, message, and optional cause.
func NewError(code, message string, cause error) *Error {
	return component.NewError(code, message, cause)
}

// IsComponentError checks if an error is a component error with the given code.
func IsComponentError(err error, code string) bool {
	return component.IsComponentError(err, code)
}

// ===== COMPONENT CONSTRUCTORS =====

// NewRegistry creates a new component registry
func NewRegistry() Registry {
	return component.CreateRegistry()
}

// NewBaseComponent creates a new base component
func NewBaseComponent(id, name string, componentType ComponentType) Component {
	return component.NewBaseComponent(id, name, componentType)
}

// NewBaseComponentWithOptions creates a base component with options
func NewBaseComponentWithOptions(options BaseComponentOptions) Component {
	return component.NewBaseComponentWithOptions(options)
}

// NewDependencyAwareComponent creates a dependency-aware component
func NewDependencyAwareComponent(base Component, dependencies []string) DependencyAware {
	return component.NewDependencyAwareComponent(base, dependencies)
}

// NewDependencyAwareComponentWithOptions creates a dependency-aware component with options
func NewDependencyAwareComponentWithOptions(options DependencyAwareComponentOptions) DependencyAware {
	return component.NewDependencyAwareComponentWithOptions(options)
}

// NewFactory creates a new component factory
func NewFactory() Factory {
	return component.CreateFactory()
}

// NewFactoryWithOptions creates a factory with specific options
func NewFactoryWithOptions(options DefaultFactoryOptions) Factory {
	return component.NewFactory(options)
}

// ===== CONTEXT UTILITIES =====

// NewContext creates a new component context from Go's standard context.
// This is the public API for creating component contexts from standard Go contexts.
func NewContext(ctx stdctx.Context) Context {
	return context.NewContext(ctx)
}

// WrapContext creates a component.Context from a standard Go context.
// This is the preferred factory method for creating a context and is the
// public API equivalent of the internal WrapContext function.
func WrapContext(ctx stdctx.Context) Context {
	return context.WrapContext(ctx)
}

// Background returns a new context with no values or cancellation.
// This is the public API equivalent of context.Background().
func Background() Context {
	return context.Background()
}

// TODO returns a new context that is never canceled.
// This is the public API equivalent of context.TODO().
func TODO() Context {
	return context.TODO()
}

// WithCancel returns a new context and a cancel function.
// This is the public API for creating cancelable contexts.
func WithCancel(parent Context) (Context, func()) {
	return context.WithCancel(parent)
}

// WithTimeout returns a new context with a timeout and a cancel function.
// This is the public API for creating contexts with timeouts.
func WithTimeout(parent Context, timeout time.Duration) (Context, func()) {
	return context.WithTimeout(parent, timeout)
}

// WithDeadline returns a new context with a deadline and a cancel function.
// This is the public API for creating contexts with deadlines.
func WithDeadline(parent Context, deadline time.Time) (Context, func()) {
	return context.WithDeadline(parent, deadline)
}
