// Package database provides database connection services for the Fintechain Skeleton framework.
package database

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
	infraComponent "github.com/fintechain/skeleton/internal/infrastructure/component"
)

// DatabaseConnectionService manages database connections.
// This demonstrates a Service component with proper framework integration.
type DatabaseConnectionService struct {
	*infraComponent.BaseService
	system           component.System // Store system reference for framework services
	driverName       string
	connectionString string
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
		BaseService:      infraComponent.NewBaseService(config),
		driverName:       driverName,
		connectionString: dataSource,
	}
}

// Initialize stores the system reference for framework services access.
func (d *DatabaseConnectionService) Initialize(ctx context.Context, system component.System) error {
	if err := d.BaseService.Initialize(ctx, system); err != nil {
		return err
	}

	// Store system reference for framework services access
	d.system = system

	return nil
}

// Start establishes the database connection.
func (d *DatabaseConnectionService) Start(ctx context.Context) error {
	if err := d.BaseService.Start(ctx); err != nil {
		return err
	}

	// Simulate database connection (focus on framework patterns, not real DB)
	// In a real implementation, you would establish actual database connections here

	return nil
}

// Stop closes the database connection.
func (d *DatabaseConnectionService) Stop(ctx context.Context) error {
	// Simulate connection cleanup
	// In a real implementation, you would close database connections here

	return d.BaseService.Stop(ctx)
}

// IsConnected returns whether the database connection is active.
func (d *DatabaseConnectionService) IsConnected() bool {
	return true // This is a placeholder implementation. Actual implementation should check the connection status.
}

// GetConnectionInfo returns connection information for other components.
func (d *DatabaseConnectionService) GetConnectionInfo() map[string]interface{} {
	return map[string]interface{}{
		"driver":     d.driverName,
		"connected":  d.IsConnected(),
		"service_id": d.ID(),
	}
}
