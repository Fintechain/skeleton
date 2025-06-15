// Package database provides a complete database plugin for the Fintechain Skeleton framework.
package database

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/plugin"
	infraComponent "github.com/fintechain/skeleton/internal/infrastructure/component"
)

// DatabasePlugin provides database functionality as a plugin.
// This demonstrates the plugin-as-orchestrator pattern with database components.
type DatabasePlugin struct {
	*infraComponent.BaseService
	system            component.System // Store system reference for framework services
	connectionService *DatabaseConnectionService
	queryOperation    *QueryOperation
	dbType            string
	connectionString  string
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
		BaseService:       infraComponent.NewBaseService(config),
		connectionService: NewDatabaseConnectionService(driverName, dataSource),
		queryOperation:    NewQueryOperation(),
		dbType:            driverName,
		connectionString:  dataSource,
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

// Initialize sets up the plugin and registers its components.
func (d *DatabasePlugin) Initialize(ctx context.Context, system component.System) error {
	if err := d.BaseService.Initialize(ctx, system); err != nil {
		return err
	}

	// Store system reference for framework services access
	d.system = system

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

	return nil
}

// Start activates the plugin and starts its services.
func (d *DatabasePlugin) Start(ctx context.Context) error {
	if err := d.BaseService.Start(ctx); err != nil {
		return err
	}

	// Plugin's responsibility: Start the services it manages
	if err := d.connectionService.Start(ctx); err != nil {
		return err
	}

	return nil
}

// Stop deactivates the plugin and stops its services.
func (d *DatabasePlugin) Stop(ctx context.Context) error {
	// Plugin's responsibility: Stop the services it manages
	if err := d.connectionService.Stop(ctx); err != nil {
		return err
	}

	if err := d.BaseService.Stop(ctx); err != nil {
		return err
	}

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
