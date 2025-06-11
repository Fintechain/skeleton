// Package database provides database connectivity components for the Fintechain Skeleton framework.
package database

import (
	"github.com/fintechain/skeleton/pkg/component"
	"github.com/fintechain/skeleton/pkg/context"
	"github.com/fintechain/skeleton/pkg/runtime"
)

// DatabaseConnectionService manages database connections.
// This demonstrates a Service component with proper framework integration.
type DatabaseConnectionService struct {
	*component.BaseService
	runtime    runtime.RuntimeEnvironment // Store runtime reference for framework services
	driverName string
	dataSource string
	connected  bool
}

// NewDatabaseConnectionService creates a new database connection service.
func NewDatabaseConnectionService(driverName, dataSource string) *DatabaseConnectionService {
	config := component.ComponentConfig{
		ID:          "database-connection",
		Name:        "Database Connection",
		Description: "Manages database connectivity",
		Version:     "1.0.0",
	}

	return &DatabaseConnectionService{
		BaseService: component.NewBaseService(config),
		driverName:  driverName,
		dataSource:  dataSource,
	}
}

// Initialize stores the runtime reference for framework services access.
func (d *DatabaseConnectionService) Initialize(ctx context.Context, system component.System) error {
	if err := d.BaseService.Initialize(ctx, system); err != nil {
		return err
	}

	// Store runtime reference - this is the key pattern to demonstrate
	d.runtime = system.(runtime.RuntimeEnvironment)

	// Access framework services to show the pattern
	logger := d.runtime.Logger()
	logger.Info("Database Connection Service initialized", "component_id", d.ID())

	return nil
}

// Start establishes the database connection.
func (d *DatabaseConnectionService) Start(ctx context.Context) error {
	if err := d.BaseService.Start(ctx); err != nil {
		return err
	}

	// Access framework services through stored runtime reference
	logger := d.runtime.Logger()
	config := d.runtime.Configuration()

	// Get configuration values with fallbacks
	driverName := config.GetStringDefault("database.driver", d.driverName)
	maxConns := config.GetIntDefault("database.max_connections", 10)

	// Simulate database connection (focus on framework patterns, not real DB)
	logger.Info("Connecting to database",
		"driver", driverName,
		"max_connections", maxConns,
		"service_id", d.ID())

	// Simulate successful connection
	d.connected = true
	logger.Info("Database connection established successfully")

	return nil
}

// Stop closes the database connection gracefully.
func (d *DatabaseConnectionService) Stop(ctx context.Context) error {
	// Access framework services through stored runtime reference
	logger := d.runtime.Logger()
	logger.Info("Closing database connection", "service_id", d.ID())

	// Simulate connection cleanup
	d.connected = false
	logger.Info("Database connection closed successfully")

	return d.BaseService.Stop(ctx)
}

// IsConnected returns whether the database connection is active.
func (d *DatabaseConnectionService) IsConnected() bool {
	return d.connected
}

// GetConnectionInfo returns connection information for other components.
func (d *DatabaseConnectionService) GetConnectionInfo() map[string]interface{} {
	return map[string]interface{}{
		"driver":     d.driverName,
		"connected":  d.connected,
		"service_id": d.ID(),
	}
}
