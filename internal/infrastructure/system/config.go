// Package system provides the implementation of the system service.
package system

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/plugin"
	"github.com/fintechain/skeleton/internal/domain/storage"
	"github.com/fintechain/skeleton/internal/domain/system"
	"github.com/fintechain/skeleton/internal/infrastructure/config"
	"github.com/fintechain/skeleton/internal/infrastructure/event"
	"github.com/fintechain/skeleton/internal/infrastructure/logging"
)

// SystemServiceConfig extends the domain config with infrastructure dependencies
type SystemServiceConfig struct {
	*system.SystemServiceConfig
	Registry      component.Registry
	PluginManager plugin.PluginManager
	EventBus      event.EventBus
	MultiStore    storage.MultiStore
	Logger        logging.Logger
}

// ConfigLoader loads system service configuration
type ConfigLoader struct {
	config config.Configuration
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader(config config.Configuration) *ConfigLoader {
	return &ConfigLoader{
		config: config,
	}
}

// LoadSystemConfig loads the system service configuration
func (l *ConfigLoader) LoadSystemConfig() (*system.SystemServiceConfig, error) {
	// Create default configuration
	cfg := &system.SystemServiceConfig{
		ServiceID:        "system",
		EnableOperations: true,
		EnableServices:   true,
		EnablePlugins:    true,
		EnableEventLog:   true,
	}

	// Override with configuration values if present
	if l.config != nil {
		// Load service ID
		if serviceID := l.config.GetStringDefault("system.serviceId", ""); serviceID != "" {
			cfg.ServiceID = serviceID
		}

		// Load feature flags
		enableOps, err := l.config.GetBool("system.enableOperations")
		if err == nil {
			cfg.EnableOperations = enableOps
		}

		enableSvc, err := l.config.GetBool("system.enableServices")
		if err == nil {
			cfg.EnableServices = enableSvc
		}

		enablePlugins, err := l.config.GetBool("system.enablePlugins")
		if err == nil {
			cfg.EnablePlugins = enablePlugins
		}

		enableEventLog, err := l.config.GetBool("system.enableEventLog")
		if err == nil {
			cfg.EnableEventLog = enableEventLog
		}

		// Load storage configuration
		var storageConfig storage.MultiStoreConfig

		// Load root path
		rootPath := l.config.GetStringDefault("system.storage.rootPath", "data")
		storageConfig.RootPath = rootPath

		// Load default engine
		defaultEngine := l.config.GetStringDefault("system.storage.defaultEngine", "memory")
		storageConfig.DefaultEngine = defaultEngine

		// Load engine configs
		// This would typically require more complex deserialization
		// For now, we'll create an empty map
		storageConfig.EngineConfigs = make(map[string]storage.Config)

		// Try to load additional engine configurations from the config
		// This is a placeholder for a more complete implementation
		cfg.StorageConfig = storageConfig
	}

	return cfg, nil
}
