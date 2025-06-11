// Package webserver provides a complete web server plugin for the Fintechain Skeleton framework.
package webserver

import (
	"github.com/fintechain/skeleton/pkg/component"
	"github.com/fintechain/skeleton/pkg/context"
	"github.com/fintechain/skeleton/pkg/plugin"
	"github.com/fintechain/skeleton/pkg/runtime"
)

// WebServerPlugin provides HTTP server functionality as a plugin.
// This demonstrates the plugin-as-orchestrator pattern: plugins register and coordinate components.
type WebServerPlugin struct {
	*component.BaseService
	runtime        runtime.RuntimeEnvironment // Store runtime reference for framework services
	httpService    *HTTPService
	routeOperation *RouteOperation
	port           int
}

// NewWebServerPlugin creates a new web server plugin.
func NewWebServerPlugin(port int) *WebServerPlugin {
	config := component.ComponentConfig{
		ID:          "webserver",
		Name:        "Web Server Plugin",
		Description: "Provides HTTP server functionality with routing",
		Version:     "1.0.0",
	}

	return &WebServerPlugin{
		BaseService:    component.NewBaseService(config),
		httpService:    NewHTTPService(port),
		routeOperation: NewRouteOperation(),
		port:           port,
	}
}

// Author returns the plugin author.
func (w *WebServerPlugin) Author() string {
	return "Fintechain Team"
}

// PluginType returns the plugin type.
func (w *WebServerPlugin) PluginType() plugin.PluginType {
	return plugin.TypeIntegration
}

// Initialize sets up the plugin and registers its components.
// This is where the plugin orchestrates its components.
func (w *WebServerPlugin) Initialize(ctx context.Context, system component.System) error {
	if err := w.BaseService.Initialize(ctx, system); err != nil {
		return err
	}

	// Store runtime reference for framework services access
	w.runtime = system.(runtime.RuntimeEnvironment)
	logger := w.runtime.Logger()
	logger.Info("Initializing Web Server Plugin", "plugin_id", w.ID())

	// 1. Initialize the components this plugin provides
	if err := w.httpService.Initialize(ctx, system); err != nil {
		return err
	}

	if err := w.routeOperation.Initialize(ctx, system); err != nil {
		return err
	}

	// 2. Register components with the system registry
	registry := system.Registry()
	if err := registry.Register(w.httpService); err != nil {
		return err
	}

	if err := registry.Register(w.routeOperation); err != nil {
		return err
	}

	logger.Info("Web Server Plugin initialized successfully",
		"components_registered", 2,
		"http_service_id", w.httpService.ID(),
		"route_operation_id", w.routeOperation.ID())

	return nil
}

// Start activates the plugin and starts its services.
// Plugin directly manages service lifecycle.
func (w *WebServerPlugin) Start(ctx context.Context) error {
	if err := w.BaseService.Start(ctx); err != nil {
		return err
	}

	logger := w.runtime.Logger()
	logger.Info("Starting Web Server Plugin", "plugin_id", w.ID())

	// Plugin's responsibility: Start the services it manages
	if err := w.httpService.Start(ctx); err != nil {
		return err
	}

	logger.Info("Web Server Plugin started successfully")
	return nil
}

// Stop deactivates the plugin and stops its services.
// Plugin directly manages service lifecycle.
func (w *WebServerPlugin) Stop(ctx context.Context) error {
	logger := w.runtime.Logger()
	logger.Info("Stopping Web Server Plugin", "plugin_id", w.ID())

	// Plugin's responsibility: Stop the services it manages
	if err := w.httpService.Stop(ctx); err != nil {
		return err
	}

	if err := w.BaseService.Stop(ctx); err != nil {
		return err
	}

	logger.Info("Web Server Plugin stopped successfully")
	return nil
}

// GetHTTPService returns the HTTP service for external access.
// This allows other plugins to interact with the web server.
func (w *WebServerPlugin) GetHTTPService() *HTTPService {
	return w.httpService
}

// GetRouteOperation returns the route operation for external access.
func (w *WebServerPlugin) GetRouteOperation() *RouteOperation {
	return w.routeOperation
}

// AddCustomRoute allows other plugins to register routes with this web server.
// This demonstrates how plugins can expose functionality to other plugins.
func (w *WebServerPlugin) AddCustomRoute(pattern string, handler func(map[string]interface{}) map[string]interface{}) {
	// This would integrate with the HTTP service to add custom routes
	// Implementation would depend on the specific routing requirements
}

// Configuration keys used by this plugin:
//
// server.port (int): HTTP server port (default: constructor value)
// server.host (string): HTTP server host (default: "0.0.0.0")
// server.read_timeout (duration): HTTP read timeout (default: 10s)
// server.write_timeout (duration): HTTP write timeout (default: 10s)
//
// Example configuration (config.json):
// {
//   "server": {
//     "port": 8080,
//     "host": "localhost",
//     "read_timeout": "30s",
//     "write_timeout": "30s"
//   }
// }
