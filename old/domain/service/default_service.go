package service

import (
	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/infrastructure/context"
)

// DefaultService provides a standard implementation of the Service interface.
type DefaultService struct {
	baseService *BaseService // Use composition instead of embedding
	startFunc   func(ctx component.Context) error
	stopFunc    func(ctx component.Context) error
	healthFunc  func() bool
}

// DefaultServiceOptions contains options for creating a DefaultService.
type DefaultServiceOptions struct {
	Component  component.Component
	StartFunc  func(ctx component.Context) error
	StopFunc   func(ctx component.Context) error
	HealthFunc func() bool
}

// NewDefaultService creates a new default service with dependency injection.
func NewDefaultService(options DefaultServiceOptions) *DefaultService {
	return &DefaultService{
		baseService: NewBaseService(BaseServiceOptions{
			Component: options.Component,
		}),
		startFunc:  options.StartFunc,
		stopFunc:   options.StopFunc,
		healthFunc: options.HealthFunc,
	}
}

// CreateDefaultService is a factory method for backward compatibility.
func CreateDefaultService(comp component.Component) *DefaultService {
	return NewDefaultService(DefaultServiceOptions{
		Component: comp,
	})
}

// ID delegates to the base component
func (s *DefaultService) ID() string {
	return s.baseService.ID()
}

// Name delegates to the base component
func (s *DefaultService) Name() string {
	return s.baseService.Name()
}

// Type delegates to the base component
func (s *DefaultService) Type() component.ComponentType {
	return s.baseService.Type()
}

// Metadata delegates to the base component
func (s *DefaultService) Metadata() component.Metadata {
	return s.baseService.Metadata()
}

// Initialize initializes the service.
func (s *DefaultService) Initialize(ctx component.Context) error {
	return s.baseService.Initialize(ctx)
}

// Dispose releases resources used by the service.
func (s *DefaultService) Dispose() error {
	return s.baseService.Dispose()
}

// WithStartFunc sets the start function for this service.
func (s *DefaultService) WithStartFunc(fn func(ctx component.Context) error) *DefaultService {
	s.startFunc = fn
	return s
}

// WithStopFunc sets the stop function for this service.
func (s *DefaultService) WithStopFunc(fn func(ctx component.Context) error) *DefaultService {
	s.stopFunc = fn
	return s
}

// WithHealthFunc sets the health check function for this service.
func (s *DefaultService) WithHealthFunc(fn func() bool) *DefaultService {
	s.healthFunc = fn
	return s
}

// Start starts the service.
func (s *DefaultService) Start(ctx component.Context) error {
	// If we have a start function, use it
	if s.startFunc != nil {
		err := s.startFunc(ctx)
		if err != nil {
			return component.NewError(
				ErrServiceStart,
				"service start function failed",
				err,
			).WithDetail("service_id", s.ID())
		}
	}

	// Call the base service start method
	return s.baseService.Start(ctx)
}

// Stop stops the service.
func (s *DefaultService) Stop(ctx component.Context) error {
	// If we have a stop function, use it
	if s.stopFunc != nil {
		err := s.stopFunc(ctx)
		if err != nil {
			return component.NewError(
				ErrServiceStop,
				"service stop function failed",
				err,
			).WithDetail("service_id", s.ID())
		}
	}

	// Call the base service stop method
	return s.baseService.Stop(ctx)
}

// Status returns the current service status.
func (s *DefaultService) Status() ServiceStatus {
	return s.baseService.Status()
}

// IsHealthy checks if the service is healthy.
func (s *DefaultService) IsHealthy() bool {
	// If we have a health function, use it
	if s.healthFunc != nil {
		return s.healthFunc()
	}

	// Default to returning true if the service is running
	return s.Status() == StatusRunning
}

// BackgroundService creates a service that runs a function in the background.
func BackgroundService(id string, fn func(ctx component.Context) error) Service {
	// Create a base component for this service
	comp := component.CreateDefaultComponent(id, "Background: "+id, component.TypeService, "Runs a function in the background")

	// Create a default service
	service := CreateDefaultService(comp)

	// Add a cancel function
	var cancelFunc func()

	// Set up the start function
	service.WithStartFunc(func(ctx component.Context) error {
		// Create a cancelable context
		cancelCtx, cancel := context.WithCancel(ctx)
		cancelFunc = cancel

		// Start the function in the background
		go func() {
			_ = fn(cancelCtx)
		}()

		return nil
	})

	// Set up the stop function
	service.WithStopFunc(func(ctx component.Context) error {
		// Cancel the context
		if cancelFunc != nil {
			cancelFunc()
		}
		return nil
	})

	return service
}
