// Package webserver provides HTTP routing components for the Fintechain Skeleton framework.
package webserver

import (
	"fmt"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
	infraComponent "github.com/fintechain/skeleton/internal/infrastructure/component"
)

// RouteOperation handles HTTP route processing.
// This demonstrates an Operation component with proper framework integration.
type RouteOperation struct {
	*infraComponent.BaseOperation
	system component.System // Store system reference for framework services
}

// NewRouteOperation creates a new route operation.
func NewRouteOperation() *RouteOperation {
	config := component.ComponentConfig{
		ID:          "http-route",
		Name:        "HTTP Route Processor",
		Description: "Processes HTTP route requests",
		Version:     "1.0.0",
	}

	return &RouteOperation{
		BaseOperation: infraComponent.NewBaseOperation(config),
	}
}

// Initialize stores the system reference for framework services access.
func (r *RouteOperation) Initialize(ctx context.Context, system component.System) error {
	if err := r.BaseOperation.Initialize(ctx, system); err != nil {
		return err
	}

	// Store system reference for framework services access
	r.system = system

	return nil
}

// Execute processes HTTP route requests.
func (r *RouteOperation) Execute(ctx context.Context, input component.Input) (component.Output, error) {
	data, ok := input.Data.(map[string]interface{})
	if !ok {
		return component.Output{}, fmt.Errorf("invalid input data format")
	}

	method, _ := data["method"].(string)
	path, _ := data["path"].(string)

	// Simulate route processing (focus on framework patterns, not real HTTP)
	response := map[string]interface{}{
		"status":       "success",
		"method":       method,
		"path":         path,
		"message":      fmt.Sprintf("Route %s %s processed successfully", method, path),
		"operation_id": string(r.ID()),
	}

	return component.Output{
		Data: response,
	}, nil
}
