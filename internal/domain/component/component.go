// Package component provides the core component system for Fintechain Skeleton.
//
// The component system implements a Domain-Driven Design (DDD) approach with
// three core types: Components (basic entities), Operations (executable tasks),
// and Services (long-running processes).
package component

import (
	"github.com/fintechain/skeleton/internal/domain/context"
)

// ComponentType represents the type classification of a component.
type ComponentType string

// ComponentID represents a unique identifier for a component within the system.
type ComponentID string

// StoreID represents a unique identifier for a storage instance.
type StoreID string

// EngineID represents a unique identifier for a storage engine.
type EngineID string

const (
	// TypeComponent represents a basic managed entity.
	TypeComponent ComponentType = "component"

	// TypeOperation represents an executable instruction.
	TypeOperation ComponentType = "operation"

	// TypeService represents a long-running process.
	TypeService ComponentType = "service"

	// Legacy types for backward compatibility (deprecated)
	TypeBasic       ComponentType = "basic"
	TypeSystem      ComponentType = "system"
	TypeApplication ComponentType = "application"
)

// Metadata represents component metadata as key-value pairs.
type Metadata map[string]interface{}

// Identifiable provides core identity properties that all components must implement.
type Identifiable interface {
	// ID returns the unique identifier for this component.
	ID() ComponentID

	// Name returns a human-readable name for this component.
	Name() string

	// Description returns a detailed description of the component's purpose.
	Description() string

	// Version returns the version string for this component.
	Version() string
}

// Component represents the fundamental building block of the system.
// Components follow a lifecycle: Creation -> Registration -> Initialization -> Active -> Disposal.
type Component interface {
	Identifiable

	// Type returns the component type classification.
	Type() ComponentType

	// Metadata returns component metadata as key-value pairs.
	// Can return nil if no metadata is available.
	Metadata() Metadata

	// Initialize prepares the component for use within the system.
	// Should be idempotent - safe to call multiple times.
	Initialize(ctx context.Context, system System) error

	// Dispose cleans up component resources and prepares for shutdown.
	// Should be idempotent - safe to call multiple times.
	Dispose() error
}
