// Package component provides the core interfaces and types for the component system.
package component

import (
	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/registry"
	sys "github.com/fintechain/skeleton/internal/domain/system"
)

// ComponentType is an alias to registry.IdentifiableType to unify the type system
type ComponentType = registry.IdentifiableType

// Re-export component type constants for convenience
const (
	TypeBasic       = registry.TypeBasic
	TypeOperation   = registry.TypeOperation
	TypeService     = registry.TypeService
	TypeSystem      = registry.TypeSystem
	TypeApplication = registry.TypeApplication
)

// Metadata is a map of key-value pairs for component metadata.
type Metadata map[string]interface{}

// Component is the fundamental building block of the system.
// It composes Identifiable to inherit ID, Name, Description, and Version methods.
type Component interface {
	registry.Identifiable

	// Component-specific properties
	Type() ComponentType
	Metadata() Metadata

	// Lifecycle
	Initialize(ctx context.Context, system sys.System) error
	Dispose() error
}
