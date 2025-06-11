// Package component provides interfaces and types for the component system.
package component

// Standard component error codes
const (
	// Component lifecycle errors
	ErrComponentNotFound           = "component.component_not_found"
	ErrComponentExists             = "component.component_exists"
	ErrInvalidComponentType        = "component.invalid_component_type"
	ErrComponentNotInitialized     = "component.component_not_initialized"
	ErrComponentAlreadyInitialized = "component.component_already_initialized"
	ErrComponentDisposed           = "component.component_disposed"
	ErrInvalidComponentConfig      = "component.invalid_component_config"

	// Factory errors
	ErrFactoryNotFound = "component.factory_not_found"

	// Registry errors
	ErrRegistryFull      = "component.registry_full"
	ErrItemNotFound      = "component.item_not_found"
	ErrItemAlreadyExists = "component.item_already_exists"
	ErrInvalidItem       = "component.invalid_item"

	// Dependency errors
	ErrDependencyNotFound = "component.dependency_not_found"
	ErrCircularDependency = "component.circular_dependency"

	// Service errors
	ErrServiceNotFound       = "component.service_not_found"
	ErrServiceNotRunning     = "component.service_not_running"
	ErrServiceAlreadyRunning = "component.service_already_running"
	ErrServiceStartFailed    = "component.service_start_failed"
	ErrServiceStopFailed     = "component.service_stop_failed"

	// Plugin errors
	ErrPluginNotFound        = "component.plugin_not_found"
	ErrPluginAlreadyLoaded   = "component.plugin_already_loaded"
	ErrPluginLoadFailed      = "component.plugin_load_failed"
	ErrPluginUnloadFailed    = "component.plugin_unload_failed"
	ErrPluginDiscoveryFailed = "component.plugin_discovery_failed"

	// System errors
	ErrSystemNotInitialized = "component.system_not_initialized"
	ErrSystemNotStarted     = "component.system_not_started"
	ErrSystemAlreadyStarted = "component.system_already_started"
	ErrOperationNotFound    = "component.operation_not_found"
	ErrOperationFailed      = "component.operation_failed"

	// Infrastructure errors
	ErrEventBusNotAvailable     = "component.event_bus_not_available"
	ErrStoreManagerNotAvailable = "component.store_manager_not_available"
)
