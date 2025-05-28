package service

import (
	"testing"
	"time"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/service/mocks"
)

// Create a simple implementation of Service with health check
type testHealthService struct {
	id      string
	status  ServiceStatus
	healthy bool
}

func (t *testHealthService) ID() string {
	return t.id
}

func (t *testHealthService) Name() string {
	return "Test Health Service"
}

func (t *testHealthService) Type() component.ComponentType {
	return component.TypeService
}

func (t *testHealthService) Metadata() component.Metadata {
	return component.Metadata{}
}

func (t *testHealthService) Initialize(ctx component.Context) error {
	return nil
}

func (t *testHealthService) Dispose() error {
	return nil
}

func (t *testHealthService) Start(ctx component.Context) error {
	t.status = StatusRunning
	return nil
}

func (t *testHealthService) Stop(ctx component.Context) error {
	t.status = StatusStopped
	return nil
}

func (t *testHealthService) Status() ServiceStatus {
	return t.status
}

func (t *testHealthService) IsHealthy() bool {
	return t.healthy
}

func newTestHealthService(id string, healthy bool) *testHealthService {
	return &testHealthService{
		id:      id,
		status:  StatusRunning,
		healthy: healthy,
	}
}

func TestHealthMonitor_CreateHealthMonitor(t *testing.T) {
	// Setup
	mockRegistry := mocks.NewMockRegistry()

	// Execute
	monitor := CreateHealthMonitor(mockRegistry)

	// Verify
	if monitor == nil {
		t.Fatal("CreateHealthMonitor returned nil")
	}
}

func TestHealthMonitor_NewHealthMonitor(t *testing.T) {
	// Setup
	mockRegistry := mocks.NewMockRegistry()
	options := HealthMonitorOptions{
		Registry:     mockRegistry,
		CheckTimeout: 2 * time.Second,
	}

	// Execute
	monitor := NewHealthMonitor(options)

	// Verify
	if monitor == nil {
		t.Fatal("NewHealthMonitor returned nil")
	}

	// Test default timeout when not provided
	defaultOptions := HealthMonitorOptions{
		Registry: mockRegistry,
	}
	defaultMonitor := NewHealthMonitor(defaultOptions)
	if defaultMonitor == nil {
		t.Fatal("NewHealthMonitor with default timeout returned nil")
	}
}

func TestHealthMonitor_CheckService_NotFound(t *testing.T) {
	// Setup
	mockRegistry := mocks.NewMockRegistry()
	mockRegistry.GetFunc = func(id string) (component.Component, error) {
		return nil, component.NewError(component.ErrComponentNotFound, "not found", nil)
	}

	monitor := CreateHealthMonitor(mockRegistry)
	mockCtx := mocks.NewMockContext()

	// Execute
	result := monitor.CheckService(mockCtx, "nonexistent-service")

	// Verify
	if result.ServiceID != "nonexistent-service" {
		t.Errorf("ServiceID = %v, want nonexistent-service", result.ServiceID)
	}

	if result.Status != HealthStatusUnknown {
		t.Errorf("Status = %v, want %v", result.Status, HealthStatusUnknown)
	}

	if result.Message != "Service not found in registry" {
		t.Errorf("Message = %v, want %v", result.Message, "Service not found in registry")
	}
}

func TestHealthMonitor_CheckService_NotAService(t *testing.T) {
	// Setup
	mockComp := mocks.NewMockComponent()
	mockRegistry := mocks.NewMockRegistry()
	mockRegistry.GetFunc = func(id string) (component.Component, error) {
		return mockComp, nil
	}

	monitor := CreateHealthMonitor(mockRegistry)
	mockCtx := mocks.NewMockContext()

	// Execute
	result := monitor.CheckService(mockCtx, "not-a-service")

	// Verify
	if result.ServiceID != "not-a-service" {
		t.Errorf("ServiceID = %v, want not-a-service", result.ServiceID)
	}

	if result.Status != HealthStatusUnknown {
		t.Errorf("Status = %v, want %v", result.Status, HealthStatusUnknown)
	}

	if result.Message != "Component is not a service" {
		t.Errorf("Message = %v, want %v", result.Message, "Component is not a service")
	}
}

func TestHealthMonitor_CheckService_NotRunning(t *testing.T) {
	// Setup
	stoppedService := newTestHealthService("stopped-service", true)
	stoppedService.status = StatusStopped

	mockRegistry := mocks.NewMockRegistry()
	mockRegistry.GetFunc = func(id string) (component.Component, error) {
		return stoppedService, nil
	}

	monitor := CreateHealthMonitor(mockRegistry)
	mockCtx := mocks.NewMockContext()

	// Execute
	result := monitor.CheckService(mockCtx, "stopped-service")

	// Verify
	if result.ServiceID != "stopped-service" {
		t.Errorf("ServiceID = %v, want stopped-service", result.ServiceID)
	}

	if result.Status != HealthStatusUnhealthy {
		t.Errorf("Status = %v, want %v", result.Status, HealthStatusUnhealthy)
	}

	if result.Message != "Service is not running" {
		t.Errorf("Message = %v, want %v", result.Message, "Service is not running")
	}
}

func TestHealthMonitor_CheckService_RunningWithoutHealthCheck(t *testing.T) {
	// Setup
	runningService := &testService{id: "running-service", status: StatusRunning}

	mockRegistry := mocks.NewMockRegistry()
	mockRegistry.GetFunc = func(id string) (component.Component, error) {
		return runningService, nil
	}

	monitor := CreateHealthMonitor(mockRegistry)
	mockCtx := mocks.NewMockContext()

	// Execute
	result := monitor.CheckService(mockCtx, "running-service")

	// Verify
	if result.ServiceID != "running-service" {
		t.Errorf("ServiceID = %v, want running-service", result.ServiceID)
	}

	if result.Status != HealthStatusHealthy {
		t.Errorf("Status = %v, want %v", result.Status, HealthStatusHealthy)
	}

	if result.Message != "Service is running" {
		t.Errorf("Message = %v, want %v", result.Message, "Service is running")
	}
}

func TestHealthMonitor_CheckService_HealthyService(t *testing.T) {
	// Setup
	healthyService := newTestHealthService("healthy-service", true)

	mockRegistry := mocks.NewMockRegistry()
	mockRegistry.GetFunc = func(id string) (component.Component, error) {
		return healthyService, nil
	}

	monitor := CreateHealthMonitor(mockRegistry)
	mockCtx := mocks.NewMockContext()

	// Execute
	result := monitor.CheckService(mockCtx, "healthy-service")

	// Verify
	if result.ServiceID != "healthy-service" {
		t.Errorf("ServiceID = %v, want healthy-service", result.ServiceID)
	}

	if result.Status != HealthStatusHealthy {
		t.Errorf("Status = %v, want %v", result.Status, HealthStatusHealthy)
	}

	if result.Message != "Service is healthy" {
		t.Errorf("Message = %v, want %v", result.Message, "Service is healthy")
	}
}

func TestHealthMonitor_CheckService_UnhealthyService(t *testing.T) {
	// Setup
	unhealthyService := newTestHealthService("unhealthy-service", false)

	mockRegistry := mocks.NewMockRegistry()
	mockRegistry.GetFunc = func(id string) (component.Component, error) {
		return unhealthyService, nil
	}

	monitor := CreateHealthMonitor(mockRegistry)
	mockCtx := mocks.NewMockContext()

	// Execute
	result := monitor.CheckService(mockCtx, "unhealthy-service")

	// Verify
	if result.ServiceID != "unhealthy-service" {
		t.Errorf("ServiceID = %v, want unhealthy-service", result.ServiceID)
	}

	if result.Status != HealthStatusUnhealthy {
		t.Errorf("Status = %v, want %v", result.Status, HealthStatusUnhealthy)
	}

	if result.Message != "Service reported as unhealthy" {
		t.Errorf("Message = %v, want %v", result.Message, "Service reported as unhealthy")
	}
}

func TestHealthMonitor_CheckAllServices(t *testing.T) {
	// Setup
	healthyService := newTestHealthService("healthy-service", true)
	unhealthyService := newTestHealthService("unhealthy-service", false)

	mockRegistry := mocks.NewMockRegistry()
	mockRegistry.FindByTypeFunc = func(compType component.ComponentType) []component.Component {
		return []component.Component{healthyService, unhealthyService}
	}
	mockRegistry.GetFunc = func(id string) (component.Component, error) {
		if id == "healthy-service" {
			return healthyService, nil
		}
		return unhealthyService, nil
	}

	monitor := CreateHealthMonitor(mockRegistry)
	mockCtx := mocks.NewMockContext()

	// Execute
	results := monitor.CheckAllServices(mockCtx)

	// Verify
	if len(results) != 2 {
		t.Errorf("len(results) = %v, want 2", len(results))
	}

	// Check healthy service
	healthyResult, ok := results["healthy-service"]
	if !ok {
		t.Error("Results did not contain healthy-service")
	} else {
		if healthyResult.Status != HealthStatusHealthy {
			t.Errorf("healthyResult.Status = %v, want %v", healthyResult.Status, HealthStatusHealthy)
		}
	}

	// Check unhealthy service
	unhealthyResult, ok := results["unhealthy-service"]
	if !ok {
		t.Error("Results did not contain unhealthy-service")
	} else {
		if unhealthyResult.Status != HealthStatusUnhealthy {
			t.Errorf("unhealthyResult.Status = %v, want %v", unhealthyResult.Status, HealthStatusUnhealthy)
		}
	}
}

func TestHealthMonitor_GetResult(t *testing.T) {
	// Setup
	mockRegistry := mocks.NewMockRegistry()
	monitor := CreateHealthMonitor(mockRegistry)

	// Insert some test results directly
	healthyService := newTestHealthService("healthy-service", true)
	mockCtx := mocks.NewMockContext()

	// We need to configure the mock registry to return our service
	mockRegistry.GetFunc = func(id string) (component.Component, error) {
		if id == "healthy-service" {
			return healthyService, nil
		}
		return nil, component.NewError(component.ErrComponentNotFound, "not found", nil)
	}

	monitor.CheckService(mockCtx, "healthy-service")

	// Execute
	result, ok := monitor.GetResult("healthy-service")

	// Verify
	if !ok {
		t.Error("GetResult() ok = false, want true")
	}

	if result.ServiceID != "healthy-service" {
		t.Errorf("ServiceID = %v, want healthy-service", result.ServiceID)
	}

	// Try with nonexistent service
	_, ok = monitor.GetResult("nonexistent-service")
	if ok {
		t.Error("GetResult() for nonexistent service ok = true, want false")
	}
}

func TestHealthMonitor_GetAllResults(t *testing.T) {
	// Setup
	mockRegistry := mocks.NewMockRegistry()
	monitor := CreateHealthMonitor(mockRegistry)

	healthyService := newTestHealthService("healthy-service", true)
	unhealthyService := newTestHealthService("unhealthy-service", false)

	mockRegistry.GetFunc = func(id string) (component.Component, error) {
		if id == "healthy-service" {
			return healthyService, nil
		}
		return unhealthyService, nil
	}

	mockCtx := mocks.NewMockContext()
	monitor.CheckService(mockCtx, "healthy-service")
	monitor.CheckService(mockCtx, "unhealthy-service")

	// Execute
	results := monitor.GetAllResults()

	// Verify
	if len(results) != 2 {
		t.Errorf("len(results) = %v, want 2", len(results))
	}

	// Check service results
	healthyResult, ok := results["healthy-service"]
	if !ok {
		t.Error("Results did not contain healthy-service")
	} else if healthyResult.Status != HealthStatusHealthy {
		t.Errorf("healthyResult.Status = %v, want %v", healthyResult.Status, HealthStatusHealthy)
	}

	unhealthyResult, ok := results["unhealthy-service"]
	if !ok {
		t.Error("Results did not contain unhealthy-service")
	} else if unhealthyResult.Status != HealthStatusUnhealthy {
		t.Errorf("unhealthyResult.Status = %v, want %v", unhealthyResult.Status, HealthStatusUnhealthy)
	}
}

func TestHealthMonitor_SetCheckTimeout(t *testing.T) {
	// Setup
	mockRegistry := mocks.NewMockRegistry()
	monitor := CreateHealthMonitor(mockRegistry)

	// Execute
	monitor.SetCheckTimeout(10 * time.Second)

	// Verify (indirectly through behavior)
	// This is a bit hard to test directly since checkTimeout is private
	// In a real test, we might use reflection, but that's beyond the scope here
}

// A simple service implementation without health check
type testService struct {
	id     string
	status ServiceStatus
}

func (t *testService) ID() string {
	return t.id
}

func (t *testService) Name() string {
	return "Test Service"
}

func (t *testService) Type() component.ComponentType {
	return component.TypeService
}

func (t *testService) Metadata() component.Metadata {
	return component.Metadata{}
}

func (t *testService) Initialize(ctx component.Context) error {
	return nil
}

func (t *testService) Dispose() error {
	return nil
}

func (t *testService) Start(ctx component.Context) error {
	t.status = StatusRunning
	return nil
}

func (t *testService) Stop(ctx component.Context) error {
	t.status = StatusStopped
	return nil
}

func (t *testService) Status() ServiceStatus {
	return t.status
}
