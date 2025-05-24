package system

import (
	"github.com/ebanfa/skeleton/internal/domain/component"
	"github.com/ebanfa/skeleton/internal/domain/plugin"
	"github.com/ebanfa/skeleton/internal/domain/storage"
	"github.com/ebanfa/skeleton/internal/infrastructure/config"
	"github.com/ebanfa/skeleton/internal/infrastructure/event"
)

// Re-export commonly used types for client convenience
type Plugin = plugin.Plugin
type Registry = component.Registry
type PluginManager = plugin.PluginManager
type EventBus = event.EventBus
type MultiStore = storage.MultiStore
type Configuration = config.Configuration

// Config represents system configuration
type Config struct {
	ServiceID     string                   `json:"serviceId"`
	StorageConfig storage.MultiStoreConfig `json:"storage"`
}
