package service

import (
	"errors"
	"testing"

	"github.com/ebanfa/skeleton/internal/domain/component"
	"github.com/ebanfa/skeleton/internal/domain/service/mocks"
)

func TestDefaultService_CreateDefaultService(t *testing.T) {
	// Setup
	mockComp := mocks.NewMockComponent()

	// Execute
	defaultService := CreateDefaultService(mockComp)

	// Verify
	if defaultService == nil {
		t.Error("CreateDefaultService returned nil")
	}

	if defaultService.ID() != mockComp.ID() {
		t.Errorf("ID() = %v, want %v", defaultService.ID(), mockComp.ID())
	}

	if defaultService.Status() != StatusStopped {
		t.Errorf("Status = %v, want %v", defaultService.Status(), StatusStopped)
	}
}

func TestDefaultService_NewDefaultService(t *testing.T) {
	// Setup
	mockComp := mocks.NewMockComponent()

	startCalled := false
	stopCalled := false
	healthyCalled := false

	options := DefaultServiceOptions{
		Component: mockComp,
		StartFunc: func(ctx component.Context) error {
			startCalled = true
			return nil
		},
		StopFunc: func(ctx component.Context) error {
			stopCalled = true
			return nil
		},
		HealthFunc: func() bool {
			healthyCalled = true
			return true
		},
	}

	// Execute
	defaultService := NewDefaultService(options)
	mockCtx := mocks.NewMockContext()

	err1 := defaultService.Start(mockCtx)
	err2 := defaultService.Stop(mockCtx)
	healthy := defaultService.IsHealthy()

	// Verify
	if defaultService == nil {
		t.Error("NewDefaultService returned nil")
	}

	if !startCalled {
		t.Error("Start function was not called")
	}

	if err1 != nil {
		t.Errorf("Start() error = %v, want nil", err1)
	}

	if !stopCalled {
		t.Error("Stop function was not called")
	}

	if err2 != nil {
		t.Errorf("Stop() error = %v, want nil", err2)
	}

	if !healthyCalled {
		t.Error("IsHealthy function was not called")
	}

	if !healthy {
		t.Error("IsHealthy() = false, want true")
	}
}

func TestDefaultService_DelegationToBaseService(t *testing.T) {
	// Setup
	mockComp := mocks.NewMockComponent()
	defaultService := CreateDefaultService(mockComp)
	mockCtx := mocks.NewMockContext()

	// Test delegation to base component methods
	if defaultService.ID() != mockComp.ID() {
		t.Errorf("ID() = %v, want %v", defaultService.ID(), mockComp.ID())
	}

	if defaultService.Name() != mockComp.Name() {
		t.Errorf("Name() = %v, want %v", defaultService.Name(), mockComp.Name())
	}

	if defaultService.Type() != mockComp.Type() {
		t.Errorf("Type() = %v, want %v", defaultService.Type(), mockComp.Type())
	}

	// Test delegation to base service methods
	if defaultService.Status() != StatusStopped {
		t.Errorf("Status() = %v, want %v", defaultService.Status(), StatusStopped)
	}

	// Initialize and Dispose should delegate to the base component
	defaultService.Initialize(mockCtx)
	if mockComp.InitializeCalls == nil || len(mockComp.InitializeCalls) != 1 {
		t.Error("Initialize() did not delegate to component")
	}

	defaultService.Dispose()
	if mockComp.DisposeCalls != 1 {
		t.Error("Dispose() did not delegate to component")
	}
}

func TestDefaultService_StartFunc(t *testing.T) {
	tests := []struct {
		name        string
		startFunc   func(ctx component.Context) error
		expectError bool
	}{
		{
			name: "Start function successful",
			startFunc: func(ctx component.Context) error {
				return nil
			},
			expectError: false,
		},
		{
			name: "Start function returns error",
			startFunc: func(ctx component.Context) error {
				return errors.New("start error")
			},
			expectError: true,
		},
		{
			name:        "No start function",
			startFunc:   nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockComp := mocks.NewMockComponent()
			defaultService := CreateDefaultService(mockComp)

			if tt.startFunc != nil {
				defaultService.WithStartFunc(tt.startFunc)
			}

			mockCtx := mocks.NewMockContext()

			// Execute
			err := defaultService.Start(mockCtx)

			// Verify
			if (err != nil) != tt.expectError {
				t.Errorf("Start() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestDefaultService_StopFunc(t *testing.T) {
	tests := []struct {
		name        string
		stopFunc    func(ctx component.Context) error
		expectError bool
	}{
		{
			name: "Stop function successful",
			stopFunc: func(ctx component.Context) error {
				return nil
			},
			expectError: false,
		},
		{
			name: "Stop function returns error",
			stopFunc: func(ctx component.Context) error {
				return errors.New("stop error")
			},
			expectError: true,
		},
		{
			name:        "No stop function",
			stopFunc:    nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockComp := mocks.NewMockComponent()
			defaultService := CreateDefaultService(mockComp)

			if tt.stopFunc != nil {
				defaultService.WithStopFunc(tt.stopFunc)
			}

			mockCtx := mocks.NewMockContext()

			// Execute
			err := defaultService.Stop(mockCtx)

			// Verify
			if (err != nil) != tt.expectError {
				t.Errorf("Stop() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

func TestDefaultService_HealthFunc(t *testing.T) {
	tests := []struct {
		name       string
		healthFunc func() bool
		want       bool
	}{
		{
			name: "Health function returns true",
			healthFunc: func() bool {
				return true
			},
			want: true,
		},
		{
			name: "Health function returns false",
			healthFunc: func() bool {
				return false
			},
			want: false,
		},
		{
			name:       "No health function, service running",
			healthFunc: nil,
			want:       true, // Default is true when running
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockComp := mocks.NewMockComponent()
			defaultService := CreateDefaultService(mockComp)

			if tt.healthFunc != nil {
				defaultService.WithHealthFunc(tt.healthFunc)
			}

			// For the third test case where we rely on Status
			if tt.name == "No health function, service running" {
				mockCtx := mocks.NewMockContext()
				_ = defaultService.Start(mockCtx)
			}

			// Execute
			got := defaultService.IsHealthy()

			// Verify
			if got != tt.want {
				t.Errorf("IsHealthy() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBackgroundService(t *testing.T) {
	// Setup
	bgFunction := func(ctx component.Context) error {
		// Just simulate a background process
		return nil
	}

	// Execute
	bgService := BackgroundService("test-bg", bgFunction)
	if bgService == nil {
		t.Fatal("BackgroundService() returned nil")
	}

	if bgService.ID() != "test-bg" {
		t.Errorf("ID() = %v, want test-bg", bgService.ID())
	}

	mockCtx := mocks.NewMockContext()

	// Start the service
	err1 := bgService.Start(mockCtx)
	if err1 != nil {
		t.Errorf("Start() error = %v", err1)
	}

	// Stop the service
	err2 := bgService.Stop(mockCtx)
	if err2 != nil {
		t.Errorf("Stop() error = %v", err2)
	}

	// Verify status
	if bgService.Status() != StatusStopped {
		t.Errorf("Status after Stop() = %v, want %v", bgService.Status(), StatusStopped)
	}
}
