// Package database provides database connectivity components for the Fintechain Skeleton framework.
package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
	domainRuntime "github.com/fintechain/skeleton/internal/domain/runtime"
	infraComponent "github.com/fintechain/skeleton/internal/infrastructure/component"
	"github.com/fintechain/skeleton/pkg/event"
)

// QueryOperation handles database query processing.
// This demonstrates an Operation component with proper framework integration.
type QueryOperation struct {
	*infraComponent.BaseOperation
	system component.System // Store system reference for framework services
}

// NewQueryOperation creates a new query operation.
func NewQueryOperation() *QueryOperation {
	config := component.ComponentConfig{
		ID:          "database-query",
		Name:        "Database Query Processor",
		Description: "Processes database query requests",
		Version:     "1.0.0",
	}

	return &QueryOperation{
		BaseOperation: infraComponent.NewBaseOperation(config),
	}
}

// Initialize stores the system reference for framework services access.
func (q *QueryOperation) Initialize(ctx context.Context, system component.System) error {
	if err := q.BaseOperation.Initialize(ctx, system); err != nil {
		return err
	}

	// Store system reference for framework services access
	q.system = system

	return nil
}

// Execute processes database queries.
func (q *QueryOperation) Execute(ctx context.Context, input component.Input) (component.Output, error) {
	data, ok := input.Data.(map[string]interface{})
	if !ok {
		return component.Output{}, fmt.Errorf("invalid input data format")
	}

	query, _ := data["query"].(string)

	// Access configuration through the RuntimeEnvironment interface
	var appName string
	var appVersion string
	var loggerTest string
	var eventBusID string

	// Try to cast to RuntimeEnvironment to access configuration
	if runtimeEnv, ok := q.system.(domainRuntime.RuntimeEnvironment); ok {
		config := runtimeEnv.Configuration()
		appName = config.GetStringDefault("app.name", "Default App")
		appVersion = config.GetStringDefault("app.version", "1.0.0")

		// Test custom logger (we can't get ID from Logger interface, but we can verify it works)
		logger := runtimeEnv.Logger()
		logger.Info("🔍 Testing custom logger from database query operation")
		loggerTest = "Custom Logger Working"

		// Test custom event bus
		eventBus := runtimeEnv.EventBus()
		eventBusID = string(eventBus.ID())
		// Publish a test event to verify event bus is working
		testEvent := event.NewEvent("test.custom.dependencies", "database-query-operation", map[string]interface{}{
			"source":  "database-query-operation",
			"message": "Testing custom event bus",
		})
		eventBus.Publish(testEvent)
	} else {
		appName = "System Interface Failed"
		appVersion = "System Interface Failed"
		loggerTest = "System Interface Failed"
		eventBusID = "System Interface Failed"
	}

	// Simulate query processing (focus on framework patterns, not real SQL)
	result := map[string]interface{}{
		"status":        "success",
		"query":         query,
		"rows_affected": 1,
		"operation_id":  string(q.ID()),
		"message":       fmt.Sprintf("Query executed successfully: %s", query),
		"app_name":      appName,    // This will show if custom config is used
		"app_version":   appVersion, // This will show if custom config is used
		"logger_test":   loggerTest, // This will show if custom logger is used
		"event_bus_id":  eventBusID, // This will show if custom event bus is used
	}

	return component.Output{
		Data: result,
	}, nil
}

// processQuery handles simple query processing logic.
func (q *QueryOperation) processQuery(queryType, querySQL string) map[string]interface{} {
	// Simple query processing - demonstrate operation logic without complexity
	switch queryType {
	case "validate":
		return q.validateQuery(querySQL)
	case "parse":
		return q.parseQuery(querySQL)
	default:
		return map[string]interface{}{
			"message":    "Query processed",
			"query_type": queryType,
			"sql_length": len(querySQL),
			"status":     "success",
			"timestamp":  time.Now().Format(time.RFC3339),
		}
	}
}

// validateQuery performs simple query validation.
func (q *QueryOperation) validateQuery(query string) map[string]interface{} {
	if query == "" {
		return map[string]interface{}{
			"valid":   false,
			"message": "Empty query",
			"status":  "error",
		}
	}

	// Simple validation (demonstrate operation logic)
	valid := len(query) > 0 && len(query) < 1000
	operation := q.getOperationType(query)

	return map[string]interface{}{
		"valid":     valid,
		"operation": operation,
		"length":    len(query),
		"status":    "validated",
	}
}

// parseQuery analyzes simple query structure.
func (q *QueryOperation) parseQuery(query string) map[string]interface{} {
	if query == "" {
		return map[string]interface{}{
			"status":  "error",
			"message": "Empty query",
		}
	}

	operation := q.getOperationType(query)

	return map[string]interface{}{
		"sql":       query,
		"operation": operation,
		"length":    len(query),
		"status":    "parsed",
		"timestamp": time.Now().Format(time.RFC3339),
	}
}

// getOperationType determines the SQL operation type (simplified).
func (q *QueryOperation) getOperationType(query string) string {
	upper := strings.ToUpper(strings.TrimSpace(query))

	if strings.HasPrefix(upper, "SELECT") {
		return "SELECT"
	} else if strings.HasPrefix(upper, "INSERT") {
		return "INSERT"
	} else if strings.HasPrefix(upper, "UPDATE") {
		return "UPDATE"
	} else if strings.HasPrefix(upper, "DELETE") {
		return "DELETE"
	}

	return "UNKNOWN"
}
