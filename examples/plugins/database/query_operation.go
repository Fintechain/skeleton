// Package database provides database connectivity components for the Fintechain Skeleton framework.
package database

import (
	"fmt"
	"strings"
	"time"

	"github.com/fintechain/skeleton/pkg/component"
	"github.com/fintechain/skeleton/pkg/context"
	"github.com/fintechain/skeleton/pkg/runtime"
)

// QueryOperation handles database query processing.
// This demonstrates an Operation component with proper framework integration.
type QueryOperation struct {
	*component.BaseOperation
	runtime runtime.RuntimeEnvironment // Store runtime reference for framework services
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
		BaseOperation: component.NewBaseOperation(config),
	}
}

// Initialize stores the runtime reference for framework services access.
func (q *QueryOperation) Initialize(ctx context.Context, system component.System) error {
	if err := q.BaseOperation.Initialize(ctx, system); err != nil {
		return err
	}

	// Store runtime reference - this is the key pattern to demonstrate
	q.runtime = system.(runtime.RuntimeEnvironment)

	// Access framework services to show the pattern
	logger := q.runtime.Logger()
	logger.Info("Database Query Operation initialized", "component_id", q.ID())

	return nil
}

// Execute processes a database query request.
// Input should contain query information, returns processed query data.
func (q *QueryOperation) Execute(ctx context.Context, input component.Input) (component.Output, error) {
	// Access framework services through stored runtime reference
	logger := q.runtime.Logger()

	// Parse input data
	data, ok := input.Data.(map[string]interface{})
	if !ok {
		return component.Output{}, fmt.Errorf("invalid input data format")
	}

	// Extract query information (simplified)
	queryType, _ := data["type"].(string)
	querySQL, _ := data["sql"].(string)

	logger.Info("Processing database query",
		"query_type", queryType,
		"sql_length", len(querySQL),
		"operation_id", q.ID())

	// Simple query processing (focus on framework patterns, not SQL logic)
	result := q.processQuery(queryType, querySQL)

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
