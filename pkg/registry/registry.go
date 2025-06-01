// Package registry provides a generic registry for identifiable items.
package registry

import (
	"github.com/fintechain/skeleton/internal/domain/registry"
	registryImpl "github.com/fintechain/skeleton/internal/infrastructure/registry"
)

// Re-export registry interfaces
type Identifiable = registry.Identifiable
type IdentifiableFactory = registry.IdentifiableFactory
type Registry = registry.Registry

// Re-export registry types
type IdentifiableType = registry.IdentifiableType
type IdentifiableConfig = registry.IdentifiableConfig

// Re-export identifiable type constants
const (
	TypeBasic       = registry.TypeBasic
	TypeOperation   = registry.TypeOperation
	TypeService     = registry.TypeService
	TypeSystem      = registry.TypeSystem
	TypeApplication = registry.TypeApplication
)

// Re-export registry error constants
const (
	ErrItemNotFound      = registry.ErrItemNotFound
	ErrItemAlreadyExists = registry.ErrItemAlreadyExists
	ErrInvalidItem       = registry.ErrInvalidItem
)

// Factory functions

// NewRegistry creates a new thread-safe registry instance.
func NewRegistry() Registry {
	return registryImpl.NewRegistry()
}
