// Package component provides component system interfaces and implementations.
package component

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	infraComponent "github.com/fintechain/skeleton/internal/infrastructure/component"
)

// Core interfaces
type Component = component.Component
type Registry = component.Registry
type System = component.System
type Factory = component.Factory
type Operation = component.Operation
type Service = component.Service

// Types and constants
type ComponentID = component.ComponentID
type ComponentType = component.ComponentType
type ComponentConfig = component.ComponentConfig
type Input = component.Input
type Output = component.Output
type Metadata = component.Metadata
type ServiceStatus = component.ServiceStatus

// Component types
const (
	TypeComponent   = component.TypeComponent
	TypeOperation   = component.TypeOperation
	TypeService     = component.TypeService
	TypeSystem      = component.TypeSystem
	TypeApplication = component.TypeApplication
)

// Service status constants
const (
	StatusStopped  = component.StatusStopped
	StatusStarting = component.StatusStarting
	StatusRunning  = component.StatusRunning
	StatusStopping = component.StatusStopping
	StatusError    = component.StatusError
)

// Error constants
const (
	ErrComponentNotFound           = component.ErrComponentNotFound
	ErrComponentExists             = component.ErrComponentExists
	ErrInvalidComponentType        = component.ErrInvalidComponentType
	ErrComponentNotInitialized     = component.ErrComponentNotInitialized
	ErrComponentAlreadyInitialized = component.ErrComponentAlreadyInitialized
	ErrComponentDisposed           = component.ErrComponentDisposed
	ErrInvalidComponentConfig      = component.ErrInvalidComponentConfig
	ErrFactoryNotFound             = component.ErrFactoryNotFound
	ErrRegistryFull                = component.ErrRegistryFull
	ErrItemNotFound                = component.ErrItemNotFound
	ErrItemAlreadyExists           = component.ErrItemAlreadyExists
	ErrInvalidItem                 = component.ErrInvalidItem
	ErrDependencyNotFound          = component.ErrDependencyNotFound
	ErrCircularDependency          = component.ErrCircularDependency
	ErrServiceNotFound             = component.ErrServiceNotFound
	ErrServiceNotRunning           = component.ErrServiceNotRunning
	ErrServiceAlreadyRunning       = component.ErrServiceAlreadyRunning
	ErrServiceStartFailed          = component.ErrServiceStartFailed
	ErrServiceStopFailed           = component.ErrServiceStopFailed
	ErrSystemNotInitialized        = component.ErrSystemNotInitialized
	ErrSystemNotStarted            = component.ErrSystemNotStarted
	ErrSystemAlreadyStarted        = component.ErrSystemAlreadyStarted
	ErrOperationNotFound           = component.ErrOperationNotFound
	ErrOperationFailed             = component.ErrOperationFailed
)

// Base implementations
type BaseComponent = infraComponent.BaseComponent
type BaseOperation = infraComponent.BaseOperation
type BaseService = infraComponent.BaseService

// Factory functions
var NewRegistry = infraComponent.NewRegistry
var NewSystem = infraComponent.NewSystem
var NewBaseComponent = infraComponent.NewBaseComponent
var NewBaseOperation = infraComponent.NewBaseOperation
var NewBaseService = infraComponent.NewBaseService

// Note: ComponentConfig is a struct, not created by a constructor function
