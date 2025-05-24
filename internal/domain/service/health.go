package service

import (
	"sync"
	"time"

	"github.com/ebanfa/skeleton/internal/domain/component"
	"github.com/ebanfa/skeleton/internal/infrastructure/context"
)

// HealthStatus represents the health status of a service.
type HealthStatus string

const (
	// HealthStatusUnknown indicates the health status is unknown.
	HealthStatusUnknown HealthStatus = "unknown"

	// HealthStatusHealthy indicates the service is healthy.
	HealthStatusHealthy HealthStatus = "healthy"

	// HealthStatusUnhealthy indicates the service is unhealthy.
	HealthStatusUnhealthy HealthStatus = "unhealthy"

	// HealthStatusDegraded indicates the service is running but degraded.
	HealthStatusDegraded HealthStatus = "degraded"
)

// HealthCheck is an interface that services can implement to indicate their health.
type HealthCheck interface {
	// IsHealthy returns true if the service is healthy.
	IsHealthy() bool
}

// HealthResult contains the result of a health check.
type HealthResult struct {
	// ServiceID is the ID of the service.
	ServiceID string

	// Status is the health status of the service.
	Status HealthStatus

	// Timestamp is when the check was performed.
	Timestamp time.Time

	// Message is an optional message about the health check.
	Message string
}

// HealthMonitor monitors the health of services.
type HealthMonitor struct {
	registry     component.Registry
	results      map[string]HealthResult
	checkTimeout time.Duration
	mu           sync.RWMutex
}

// HealthMonitorOptions contains options for creating a HealthMonitor.
type HealthMonitorOptions struct {
	Registry     component.Registry
	CheckTimeout time.Duration
}

// NewHealthMonitor creates a new health monitor with dependency injection.
func NewHealthMonitor(options HealthMonitorOptions) *HealthMonitor {
	timeout := options.CheckTimeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	return &HealthMonitor{
		registry:     options.Registry,
		results:      make(map[string]HealthResult),
		checkTimeout: timeout,
		mu:           sync.RWMutex{},
	}
}

// CreateHealthMonitor is a factory method for backward compatibility.
// Creates a HealthMonitor with default settings.
func CreateHealthMonitor(registry component.Registry) *HealthMonitor {
	return NewHealthMonitor(HealthMonitorOptions{
		Registry: registry,
	})
}

// SetCheckTimeout sets the timeout for health checks.
func (m *HealthMonitor) SetCheckTimeout(timeout time.Duration) {
	m.checkTimeout = timeout
}

// CheckService checks the health of a specific service.
func (m *HealthMonitor) CheckService(ctx component.Context, serviceID string) HealthResult {
	// Get the service from the registry
	comp, err := m.registry.Get(serviceID)
	if err != nil {
		return HealthResult{
			ServiceID: serviceID,
			Status:    HealthStatusUnknown,
			Timestamp: time.Now(),
			Message:   "Service not found in registry",
		}
	}

	// Check if the component is a service with health check
	service, ok := comp.(Service)
	if !ok {
		return HealthResult{
			ServiceID: serviceID,
			Status:    HealthStatusUnknown,
			Timestamp: time.Now(),
			Message:   "Component is not a service",
		}
	}

	// Check if the service is running
	if service.Status() != StatusRunning {
		return HealthResult{
			ServiceID: serviceID,
			Status:    HealthStatusUnhealthy,
			Timestamp: time.Now(),
			Message:   "Service is not running",
		}
	}

	// Check if the service implements HealthCheck
	healthCheck, ok := service.(HealthCheck)
	if !ok {
		// If not, assume it's healthy if it's running
		return HealthResult{
			ServiceID: serviceID,
			Status:    HealthStatusHealthy,
			Timestamp: time.Now(),
			Message:   "Service is running",
		}
	}

	// Create a context with timeout
	timeoutCtx, cancel := context.WithTimeout(ctx, m.checkTimeout)
	defer cancel()

	// Channel for the health check result
	resultCh := make(chan bool, 1) // Use buffered channel to avoid goroutine leak

	// Perform the health check in a goroutine
	go func() {
		select {
		case <-timeoutCtx.Done():
			// Don't do anything if the context is already done
			return
		default:
			// Check health and send result to channel
			resultCh <- healthCheck.IsHealthy()
		}
	}()

	// Wait for the result or timeout
	var healthy bool
	select {
	case healthy = <-resultCh:
		// Got a result
	case <-timeoutCtx.Done():
		// Timeout
		return HealthResult{
			ServiceID: serviceID,
			Status:    HealthStatusDegraded,
			Timestamp: time.Now(),
			Message:   "Health check timed out",
		}
	}

	// Create and return the result
	result := HealthResult{
		ServiceID: serviceID,
		Timestamp: time.Now(),
	}

	if healthy {
		result.Status = HealthStatusHealthy
		result.Message = "Service is healthy"
	} else {
		result.Status = HealthStatusUnhealthy
		result.Message = "Service reported as unhealthy"
	}

	// Store the result
	m.mu.Lock()
	m.results[serviceID] = result
	m.mu.Unlock()

	return result
}

// CheckAllServices checks the health of all services registered in the registry.
func (m *HealthMonitor) CheckAllServices(ctx component.Context) map[string]HealthResult {
	// Get all services from the registry
	services := m.registry.FindByType(component.TypeService)

	// Create a map for the results
	results := make(map[string]HealthResult)

	// Check each service
	for _, comp := range services {
		result := m.CheckService(ctx, comp.ID())
		results[comp.ID()] = result
	}

	// Store the results
	m.mu.Lock()
	for id, result := range results {
		m.results[id] = result
	}
	m.mu.Unlock()

	return results
}

// GetResult gets the most recent health check result for a service.
func (m *HealthMonitor) GetResult(serviceID string) (HealthResult, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result, ok := m.results[serviceID]
	return result, ok
}

// GetAllResults gets all health check results.
func (m *HealthMonitor) GetAllResults() map[string]HealthResult {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Create a copy of the results
	results := make(map[string]HealthResult)
	for id, result := range m.results {
		results[id] = result
	}

	return results
}
