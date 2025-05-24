// Package system provides the implementation of the system service.
package system

import (
	"fmt"

	"github.com/ebanfa/skeleton/internal/domain/component"
	"github.com/ebanfa/skeleton/internal/domain/plugin"
	"github.com/ebanfa/skeleton/internal/domain/storage"
	"github.com/ebanfa/skeleton/internal/domain/system"
	"github.com/ebanfa/skeleton/internal/infrastructure/config"
	"github.com/ebanfa/skeleton/internal/infrastructure/event"
	"github.com/ebanfa/skeleton/internal/infrastructure/logging"
)

// Builder helps construct a complete system with default components
type Builder struct {
	serviceID     string
	logger        logging.Logger
	configuration config.Configuration
	registry      component.Registry
	pluginManager plugin.PluginManager
	eventBus      event.EventBus
	multiStore    storage.MultiStore
}

// NewBuilder creates a new system builder
func NewBuilder(serviceID string) *Builder {
	return &Builder{
		serviceID: serviceID,
		logger:    logging.CreateStandardLogger(logging.Info),
	}
}

// WithLogger sets the logger
func (b *Builder) WithLogger(logger logging.Logger) *Builder {
	b.logger = logger
	return b
}

// WithConfiguration sets the configuration
func (b *Builder) WithConfiguration(configuration config.Configuration) *Builder {
	b.configuration = configuration
	return b
}

// WithRegistry sets the component registry
func (b *Builder) WithRegistry(registry component.Registry) *Builder {
	b.registry = registry
	return b
}

// WithPluginManager sets the plugin manager
func (b *Builder) WithPluginManager(pluginManager plugin.PluginManager) *Builder {
	b.pluginManager = pluginManager
	return b
}

// WithEventBus sets the event bus
func (b *Builder) WithEventBus(eventBus event.EventBus) *Builder {
	b.eventBus = eventBus
	return b
}

// WithMultiStore sets the multi-store
func (b *Builder) WithMultiStore(multiStore storage.MultiStore) *Builder {
	b.multiStore = multiStore
	return b
}

// Build creates the system service
func (b *Builder) Build() (system.SystemService, error) {
	// Validate required dependencies
	if b.registry == nil {
		return nil, fmt.Errorf("registry is required")
	}

	if b.pluginManager == nil {
		return nil, fmt.Errorf("plugin manager is required")
	}

	if b.eventBus == nil {
		return nil, fmt.Errorf("event bus is required")
	}

	if b.configuration == nil {
		return nil, fmt.Errorf("configuration is required")
	}

	if b.multiStore == nil {
		return nil, fmt.Errorf("multi-store is required")
	}

	// Create system config
	systemConfig := &system.SystemServiceConfig{
		ServiceID:        b.serviceID,
		EnableOperations: true,
		EnableServices:   true,
		EnablePlugins:    true,
		EnableEventLog:   true,
	}

	// Create factory and system service
	factory := NewFactory(
		b.registry,
		b.pluginManager,
		b.eventBus,
		b.configuration,
		b.multiStore,
		b.logger,
	)

	service, err := factory.CreateSystemService(systemConfig)
	if err != nil {
		return nil, err
	}

	return service, nil
}
