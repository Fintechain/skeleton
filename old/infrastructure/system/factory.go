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

// Factory creates system service instances
type Factory struct {
	registry      component.Registry
	pluginManager plugin.PluginManager
	eventBus      event.EventBus
	configuration config.Configuration
	multiStore    storage.MultiStore
	logger        logging.Logger
}

// NewFactory creates a new system service factory
func NewFactory(
	registry component.Registry,
	pluginManager plugin.PluginManager,
	eventBus event.EventBus,
	configuration config.Configuration,
	multiStore storage.MultiStore,
	logger logging.Logger,
) *Factory {
	if logger == nil {
		logger = logging.CreateStandardLogger(logging.Info)
	}

	return &Factory{
		registry:      registry,
		pluginManager: pluginManager,
		eventBus:      eventBus,
		configuration: configuration,
		multiStore:    multiStore,
		logger:        logger,
	}
}

// CreateSystemService creates a new system service instance
func (f *Factory) CreateSystemService(config *system.SystemServiceConfig) (system.SystemService, error) {
	// Create a new system service
	svc := NewDefaultSystemService(
		config.ServiceID,
		f.registry,
		f.pluginManager,
		f.eventBus,
		f.configuration,
		f.multiStore,
		f.logger,
	)

	// Set configuration options in the service via metadata
	// Note: We need to cast to access SetMetadata method through the BaseService
	if baseComponent, ok := svc.BaseService.Component.(*component.BaseComponent); ok {
		baseComponent.SetMetadata("enableOperations", config.EnableOperations)
		baseComponent.SetMetadata("enableServices", config.EnableServices)
		baseComponent.SetMetadata("enablePlugins", config.EnablePlugins)
		baseComponent.SetMetadata("enableEventLog", config.EnableEventLog)
	}

	f.logger.Info("Created system service with ID: %s", config.ServiceID)
	return svc, nil
}
