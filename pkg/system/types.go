package system

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/plugin"
	"github.com/fintechain/skeleton/internal/domain/storage"
	"github.com/fintechain/skeleton/internal/domain/system"
	"github.com/fintechain/skeleton/internal/infrastructure/config"
	"github.com/fintechain/skeleton/internal/infrastructure/event"
	infraSystem "github.com/fintechain/skeleton/internal/infrastructure/system"
)

// Re-export commonly used types for client convenience
type Plugin = plugin.Plugin
type Registry = component.Registry
type PluginManager = plugin.PluginManager
type EventBus = event.EventBus
type MultiStore = storage.MultiStore
type Configuration = config.Configuration

// Re-export component interfaces for external plugin development
type Component = component.Component
type Context = component.Context
type ComponentType = component.ComponentType
type Metadata = component.Metadata

// Re-export system service interface
type SystemService = system.SystemService

// Re-export storage types
type MultiStoreConfig = storage.MultiStoreConfig
type StorageConfig = storage.Config
type Store = storage.Store

// Config represents system configuration (public version)
type Config struct {
	ServiceID     string                   `json:"serviceId"`
	StorageConfig storage.MultiStoreConfig `json:"storage"`
}

// InternalConfig represents the internal system configuration
type InternalConfig = infraSystem.Config
