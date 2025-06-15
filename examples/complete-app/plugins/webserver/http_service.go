// Package webserver provides HTTP service components for the Fintechain Skeleton framework.
package webserver

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
	infraComponent "github.com/fintechain/skeleton/internal/infrastructure/component"
)

// HTTPService manages HTTP server functionality.
// This demonstrates a Service component with proper framework integration.
type HTTPService struct {
	*infraComponent.BaseService
	system component.System // Store system reference for framework services
	port   int
}

// NewHTTPService creates a new HTTP service.
func NewHTTPService(port int) *HTTPService {
	config := component.ComponentConfig{
		ID:          "http-service",
		Name:        "HTTP Service",
		Description: "Manages HTTP server functionality",
		Version:     "1.0.0",
	}

	return &HTTPService{
		BaseService: infraComponent.NewBaseService(config),
		port:        port,
	}
}

// Initialize stores the system reference for framework services access.
func (h *HTTPService) Initialize(ctx context.Context, system component.System) error {
	if err := h.BaseService.Initialize(ctx, system); err != nil {
		return err
	}

	// Store system reference - this is the key pattern to demonstrate
	h.system = system

	return nil
}

// Start begins the HTTP server operation.
func (h *HTTPService) Start(ctx context.Context) error {
	if err := h.BaseService.Start(ctx); err != nil {
		return err
	}

	// Simulate HTTP server start (focus on framework patterns, not real HTTP)
	// In a real implementation, you would start an actual HTTP server here

	return nil
}

// Stop gracefully shuts down the HTTP server.
func (h *HTTPService) Stop(ctx context.Context) error {
	// Simulate graceful shutdown
	// In a real implementation, you would stop the HTTP server here

	return h.BaseService.Stop(ctx)
}

// GetPort returns the configured port for external access.
func (h *HTTPService) GetPort() int {
	return h.port
}
