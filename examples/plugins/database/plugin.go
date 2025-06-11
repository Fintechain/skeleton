// Package database provides a complete database plugin for the Fintechain Skeleton framework.
package database

import (
	"github.com/fintechain/skeleton/pkg/component"
	"github.com/fintechain/skeleton/pkg/context"
	"github.com/fintechain/skeleton/pkg/plugin"
	"github.com/fintechain/skeleton/pkg/runtime"
)

// DatabasePlugin orchestrates database components.
// This demonstrates the plugin-as-orchestrator pattern: plugins register and coordinate components.
type DatabasePlugin struct {
	*component.BaseService
	runtime           runtime.RuntimeEnvironment // Store runtime reference for framework services
	connectionService *DatabaseConnectionService
	queryOperation    *QueryOperation
	driverName        string
	dataSource        string
}

// NewDatabasePlugin creates a new database plugin.
func NewDatabasePlugin(driverName, dataSource string) *DatabasePlugin {
	config := component.ComponentConfig{
		ID:          "database-plugin",
		Name:        "Database Plugin",
		Description: "Provides database connectivity and query processing capabilities",
		Version:     "1.0.0",
	}

	return &DatabasePlugin{
		BaseService:       component.NewBaseService(config),
		connectionService: NewDatabaseConnectionService(driverName, dataSource),
		queryOperation:    NewQueryOperation(),
		driverName:        driverName,
		dataSource:        dataSource,
	}
}

// Author returns the plugin author.
func (d *DatabasePlugin) Author() string {
	return "Fintechain Team"
}

// PluginType returns the plugin type.
func (d *DatabasePlugin) PluginType() plugin.PluginType {
	return plugin.TypeIntegration
}

// Initialize is where the plugin registers its components with the system.
// This is the plugin's main responsibility: component orchestration.
func (d *DatabasePlugin) Initialize(ctx context.Context, system component.System) error {
	if err := d.BaseService.Initialize(ctx, system); err != nil {
		return err
	}

	// Store runtime reference for framework services access
	d.runtime = system.(runtime.RuntimeEnvironment)
	logger := d.runtime.Logger()
	logger.Info("Initializing Database Plugin", "plugin_id", d.ID())

	// 1. Initialize the components this plugin provides
	if err := d.connectionService.Initialize(ctx, system); err != nil {
		return err
	}

	if err := d.queryOperation.Initialize(ctx, system); err != nil {
		return err
	}

	// 2. Register components with the system registry
	registry := system.Registry()
	if err := registry.Register(d.connectionService); err != nil {
		return err
	}

	if err := registry.Register(d.queryOperation); err != nil {
		return err
	}

	logger.Info("Database Plugin initialized successfully",
		"components_registered", 2,
		"connection_service_id", d.connectionService.ID(),
		"query_operation_id", d.queryOperation.ID())

	return nil
}

// Start starts the plugin and its services (called in daemon mode).
// Plugin directly manages service lifecycle.
func (d *DatabasePlugin) Start(ctx context.Context) error {
	if err := d.BaseService.Start(ctx); err != nil {
		return err
	}

	logger := d.runtime.Logger()
	logger.Info("Starting Database Plugin", "plugin_id", d.ID())

	// Plugin's responsibility: Start the services it manages
	if err := d.connectionService.Start(ctx); err != nil {
		return err
	}

	logger.Info("Database Plugin started successfully")
	return nil
}

// Stop stops the plugin and its services (called during shutdown).
// Plugin directly manages service lifecycle.
func (d *DatabasePlugin) Stop(ctx context.Context) error {
	logger := d.runtime.Logger()
	logger.Info("Stopping Database Plugin", "plugin_id", d.ID())

	// Plugin's responsibility: Stop the services it manages
	if err := d.connectionService.Stop(ctx); err != nil {
		return err
	}

	if err := d.BaseService.Stop(ctx); err != nil {
		return err
	}

	logger.Info("Database Plugin stopped successfully")
	return nil
}

// GetConnectionService provides access to the database connection service.
func (d *DatabasePlugin) GetConnectionService() *DatabaseConnectionService {
	return d.connectionService
}

// GetQueryOperation provides access to the query operation.
func (d *DatabasePlugin) GetQueryOperation() *QueryOperation {
	return d.queryOperation
}
