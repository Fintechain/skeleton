// Package component provides the core interfaces and types for the component system.
package component

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/registry"
	"github.com/fintechain/skeleton/internal/domain/system"
)

// Re-export component types
type ComponentType = component.ComponentType
type Metadata = component.Metadata
type Component = component.Component

// Re-export registry types for convenience
type Identifiable = registry.Identifiable

// Re-export component type constants
const (
	TypeBasic       = component.TypeBasic
	TypeOperation   = component.TypeOperation
	TypeService     = component.TypeService
	TypeSystem      = component.TypeSystem
	TypeApplication = component.TypeApplication
)

// Re-export context and system types for convenience
type Context = context.Context
type System = system.System
