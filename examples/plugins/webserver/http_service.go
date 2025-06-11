// Package webserver provides HTTP server components for the Fintechain Skeleton framework.
package webserver

import (
	"github.com/fintechain/skeleton/pkg/component"
	"github.com/fintechain/skeleton/pkg/context"
	"github.com/fintechain/skeleton/pkg/runtime"
)

// HTTPService provides HTTP server functionality.
// This demonstrates a Service component with proper framework integration.
type HTTPService struct {
	*component.BaseService
	runtime runtime.RuntimeEnvironment // Store runtime reference for framework services
	port    int
}

// NewHTTPService creates a new HTTP service with the specified port.
func NewHTTPService(port int) *HTTPService {
	config := component.ComponentConfig{
		ID:          "http-server",
		Name:        "HTTP Server",
		Description: "Provides HTTP server functionality",
		Version:     "1.0.0",
	}

	return &HTTPService{
		BaseService: component.NewBaseService(config),
		port:        port,
	}
}

// Initialize stores the runtime reference for framework services access.
func (h *HTTPService) Initialize(ctx context.Context, system component.System) error {
	if err := h.BaseService.Initialize(ctx, system); err != nil {
		return err
	}

	// Store runtime reference - this is the key pattern to demonstrate
	h.runtime = system.(runtime.RuntimeEnvironment)

	// Access framework services to show the pattern
	logger := h.runtime.Logger()
	logger.Info("HTTP Service initialized", "component_id", h.ID())

	return nil
}

// Start begins the HTTP server operation.
func (h *HTTPService) Start(ctx context.Context) error {
	if err := h.BaseService.Start(ctx); err != nil {
		return err
	}

	// Access framework services through stored runtime reference
	logger := h.runtime.Logger()
	config := h.runtime.Configuration()

	// Get port from configuration with fallback to constructor value
	port := h.port
	if configPort := config.GetIntDefault("http.port", 0); configPort > 0 {
		port = configPort
	}

	host := config.GetStringDefault("http.host", "0.0.0.0")

	// Simulate HTTP server start (focus on framework patterns, not real HTTP)
	logger.Info("HTTP server started",
		"host", host,
		"port", port,
		"status", "running",
		"service_id", h.ID())

	return nil
}

// Stop gracefully shuts down the HTTP server.
func (h *HTTPService) Stop(ctx context.Context) error {
	// Access framework services through stored runtime reference
	logger := h.runtime.Logger()
	logger.Info("Stopping HTTP server", "service_id", h.ID())

	// Simulate graceful shutdown
	logger.Info("HTTP server stopped successfully")

	return h.BaseService.Stop(ctx)
}

// GetPort returns the configured port for external access.
func (h *HTTPService) GetPort() int {
	if h.runtime != nil {
		config := h.runtime.Configuration()
		return config.GetIntDefault("http.port", h.port)
	}
	return h.port
}
