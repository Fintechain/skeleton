// Package webserver provides HTTP server components for the Fintechain Skeleton framework.
package webserver

import (
	"fmt"
	"time"

	"github.com/fintechain/skeleton/pkg/component"
	"github.com/fintechain/skeleton/pkg/context"
	"github.com/fintechain/skeleton/pkg/runtime"
)

// RouteOperation handles HTTP route processing.
// This demonstrates an Operation component with proper framework integration.
type RouteOperation struct {
	*component.BaseOperation
	runtime runtime.RuntimeEnvironment // Store runtime reference for framework services
}

// NewRouteOperation creates a new route operation.
func NewRouteOperation() *RouteOperation {
	config := component.ComponentConfig{
		ID:          "http-route",
		Name:        "HTTP Route Handler",
		Description: "Processes HTTP route requests",
		Version:     "1.0.0",
	}

	return &RouteOperation{
		BaseOperation: component.NewBaseOperation(config),
	}
}

// Initialize stores the runtime reference for framework services access.
func (r *RouteOperation) Initialize(ctx context.Context, system component.System) error {
	if err := r.BaseOperation.Initialize(ctx, system); err != nil {
		return err
	}

	// Store runtime reference - this is the key pattern to demonstrate
	r.runtime = system.(runtime.RuntimeEnvironment)

	// Access framework services to show the pattern
	logger := r.runtime.Logger()
	logger.Info("Route Operation initialized", "component_id", r.ID())

	return nil
}

// Execute processes an HTTP route request.
// Input should contain route data, returns processed response data.
func (r *RouteOperation) Execute(ctx context.Context, input component.Input) (component.Output, error) {
	// Access framework services through stored runtime reference
	logger := r.runtime.Logger()

	// Parse input data
	data, ok := input.Data.(map[string]interface{})
	if !ok {
		return component.Output{}, fmt.Errorf("invalid input data format")
	}

	// Extract route information (simplified)
	method, _ := data["method"].(string)
	path, _ := data["path"].(string)

	logger.Info("Processing HTTP route",
		"method", method,
		"path", path,
		"operation_id", r.ID())

	// Simple route processing (focus on framework patterns, not HTTP logic)
	response := r.processRoute(method, path)

	return component.Output{
		Data: response,
	}, nil
}

// processRoute handles simple route processing logic.
func (r *RouteOperation) processRoute(method, path string) map[string]interface{} {
	// Simple route processing - demonstrate operation logic without complexity
	switch path {
	case "/api/health":
		return map[string]interface{}{
			"status":    "healthy",
			"service":   "route-operation",
			"timestamp": time.Now().Format(time.RFC3339),
		}
	case "/api/users":
		return r.handleUsersRoute(method)
	default:
		return map[string]interface{}{
			"message": "Route processed",
			"method":  method,
			"path":    path,
			"status":  "success",
		}
	}
}

// handleUsersRoute demonstrates simple endpoint processing.
func (r *RouteOperation) handleUsersRoute(method string) map[string]interface{} {
	switch method {
	case "GET":
		return map[string]interface{}{
			"users": []map[string]interface{}{
				{"id": 1, "name": "Alice"},
				{"id": 2, "name": "Bob"},
			},
			"count": 2,
		}
	case "POST":
		return map[string]interface{}{
			"message": "User created",
			"user_id": 123,
			"status":  "created",
		}
	default:
		return map[string]interface{}{
			"error":  "Method not supported",
			"method": method,
			"status": "error",
		}
	}
}
