// Package main provides a simple test plugin for demonstrating the Builder API.
package main

import (
	"fmt"
	"time"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/plugin"
	"github.com/fintechain/skeleton/internal/domain/runtime"
	infraComponent "github.com/fintechain/skeleton/internal/infrastructure/component"
)

// TestPlugin demonstrates the plugin-as-orchestrator pattern
type TestPlugin struct {
	*infraComponent.BaseService
	runtime   runtime.RuntimeEnvironment // Store runtime reference
	service   *TestService
	operation *TestOperation
}

// NewTestPlugin creates a new test plugin for demonstration purposes
func NewTestPlugin() *TestPlugin {
	config := component.ComponentConfig{
		ID:          "test-plugin",
		Name:        "Test Plugin",
		Description: "Simple plugin for Builder API demonstration",
		Version:     "1.0.0",
	}

	return &TestPlugin{
		BaseService: infraComponent.NewBaseService(config),
		service:     NewTestService(),
		operation:   NewTestOperation(),
	}
}

// Author returns the plugin author
func (p *TestPlugin) Author() string {
	return "Fintechain Team"
}

// PluginType returns the plugin type
func (p *TestPlugin) PluginType() plugin.PluginType {
	return plugin.TypeProcessor
}

// Initialize sets up the plugin and registers its components
func (p *TestPlugin) Initialize(ctx context.Context, system component.System) error {
	if err := p.BaseService.Initialize(ctx, system); err != nil {
		return err
	}

	// Store runtime reference
	p.runtime = system.(runtime.RuntimeEnvironment)
	logger := p.runtime.Logger()
	logger.Info("Initializing test plugin", "plugin_id", p.ID())

	// Initialize components
	if err := p.service.Initialize(ctx, system); err != nil {
		return err
	}

	if err := p.operation.Initialize(ctx, system); err != nil {
		return err
	}

	// Register components with system registry
	registry := system.Registry()
	if err := registry.Register(p.service); err != nil {
		return err
	}

	if err := registry.Register(p.operation); err != nil {
		return err
	}

	logger.Info("Test plugin initialized successfully",
		"components_registered", 2,
		"service_id", p.service.ID(),
		"operation_id", p.operation.ID())

	return nil
}

// Start manages service lifecycle (daemon mode)
func (p *TestPlugin) Start(ctx context.Context) error {
	if err := p.BaseService.Start(ctx); err != nil {
		return err
	}

	logger := p.runtime.Logger()
	logger.Info("Starting test plugin", "plugin_id", p.ID())

	// Start the service
	if err := p.service.Start(ctx); err != nil {
		return err
	}

	logger.Info("Test plugin started successfully")
	return nil
}

// Stop manages service cleanup
func (p *TestPlugin) Stop(ctx context.Context) error {
	logger := p.runtime.Logger()
	logger.Info("Stopping test plugin", "plugin_id", p.ID())

	// Stop the service
	if err := p.service.Stop(ctx); err != nil {
		return err
	}

	if err := p.BaseService.Stop(ctx); err != nil {
		return err
	}

	logger.Info("Test plugin stopped successfully")
	return nil
}

// TestService demonstrates a simple service component
type TestService struct {
	*infraComponent.BaseService
	runtime runtime.RuntimeEnvironment
}

// NewTestService creates a new test service
func NewTestService() *TestService {
	config := component.ComponentConfig{
		ID:          "test-service",
		Name:        "Test Service",
		Description: "Simple service for demonstration",
		Version:     "1.0.0",
	}
	return &TestService{
		BaseService: infraComponent.NewBaseService(config),
	}
}

// Initialize initializes the test service
func (s *TestService) Initialize(ctx context.Context, system component.System) error {
	if err := s.BaseService.Initialize(ctx, system); err != nil {
		return err
	}

	// Store runtime reference
	s.runtime = system.(runtime.RuntimeEnvironment)
	logger := s.runtime.Logger()
	logger.Info("Test service initialized", "service_id", s.ID())

	return nil
}

// Start starts the test service
func (s *TestService) Start(ctx context.Context) error {
	if err := s.BaseService.Start(ctx); err != nil {
		return err
	}

	logger := s.runtime.Logger()
	config := s.runtime.Configuration()

	// Get configuration with defaults
	enabled := config.GetBoolDefault("test_service.enabled", true)
	interval := config.GetDurationDefault("test_service.interval", 30*time.Second)

	if !enabled {
		logger.Info("Test service disabled by configuration")
		return nil
	}

	logger.Info("Test service started",
		"service_id", s.ID(),
		"interval", interval,
		"status", "running")

	return nil
}

// Stop stops the test service
func (s *TestService) Stop(ctx context.Context) error {
	logger := s.runtime.Logger()
	logger.Info("Test service stopping", "service_id", s.ID())

	logger.Info("Test service stopped successfully")
	return s.BaseService.Stop(ctx)
}

// TestOperation demonstrates a simple operation component
type TestOperation struct {
	*infraComponent.BaseOperation
	runtime runtime.RuntimeEnvironment
}

// NewTestOperation creates a new test operation
func NewTestOperation() *TestOperation {
	config := component.ComponentConfig{
		ID:          "test-operation",
		Name:        "Test Operation",
		Description: "Simple operation for demonstration",
		Version:     "1.0.0",
	}
	return &TestOperation{
		BaseOperation: infraComponent.NewBaseOperation(config),
	}
}

// Initialize initializes the test operation
func (o *TestOperation) Initialize(ctx context.Context, system component.System) error {
	if err := o.BaseOperation.Initialize(ctx, system); err != nil {
		return err
	}

	// Store runtime reference
	o.runtime = system.(runtime.RuntimeEnvironment)
	logger := o.runtime.Logger()
	logger.Info("Test operation initialized", "operation_id", o.ID())

	return nil
}

// Execute executes the test operation
func (o *TestOperation) Execute(ctx context.Context, input component.Input) (component.Output, error) {
	logger := o.runtime.Logger()

	// Parse input
	data, ok := input.Data.(map[string]interface{})
	if !ok {
		return component.Output{}, fmt.Errorf("invalid input data format")
	}

	// Extract input values
	message, _ := data["message"].(string)
	number, _ := data["number"].(float64)

	logger.Info("Processing test operation",
		"operation_id", o.ID(),
		"input_message", message,
		"input_number", number)

	// Simple processing
	result := map[string]interface{}{
		"processed_message": fmt.Sprintf("Processed: %s", message),
		"processed_number":  number * 2,
		"status":            "success",
		"timestamp":         time.Now().Format(time.RFC3339),
		"operation_id":      string(o.ID()),
	}

	logger.Info("Test operation completed successfully",
		"result_message", result["processed_message"],
		"result_number", result["processed_number"])

	return component.Output{Data: result}, nil
}
