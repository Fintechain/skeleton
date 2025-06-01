package service

import (
	"testing"

	"github.com/fintechain/skeleton/internal/domain/service/mocks"
)

func TestBaseService_CreateBaseService(t *testing.T) {
	// Setup
	mockComp := mocks.NewMockComponent()

	// Execute
	baseService := CreateBaseService(mockComp)

	// Verify
	if baseService == nil {
		t.Error("CreateBaseService returned nil")
	}

	if baseService.Component != mockComp {
		t.Error("BaseService not using provided component")
	}

	if baseService.Status() != StatusStopped {
		t.Errorf("BaseService initial status is %v, expected %v", baseService.Status(), StatusStopped)
	}
}

func TestBaseService_NewBaseService(t *testing.T) {
	// Setup
	mockComp := mocks.NewMockComponent()
	options := BaseServiceOptions{
		Component: mockComp,
	}

	// Execute
	baseService := NewBaseService(options)

	// Verify
	if baseService == nil {
		t.Error("NewBaseService returned nil")
	}

	if baseService.Component != mockComp {
		t.Error("BaseService not using provided component")
	}

	if baseService.Status() != StatusStopped {
		t.Errorf("BaseService initial status is %v, expected %v", baseService.Status(), StatusStopped)
	}
}

func TestBaseService_Start(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus ServiceStatus
		expectedError bool
		expectedFinal ServiceStatus
	}{
		{
			name:          "Start when stopped",
			initialStatus: StatusStopped,
			expectedError: false,
			expectedFinal: StatusRunning,
		},
		{
			name:          "Start when already running",
			initialStatus: StatusRunning,
			expectedError: false,
			expectedFinal: StatusRunning,
		},
		{
			name:          "Start when starting",
			initialStatus: StatusStarting,
			expectedError: false,
			expectedFinal: StatusStarting,
		},
		{
			name:          "Start when stopping",
			initialStatus: StatusStopping,
			expectedError: true,
			expectedFinal: StatusStopping,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockComp := mocks.NewMockComponent()
			baseService := CreateBaseService(mockComp)
			baseService.SetStatus(tt.initialStatus)
			mockCtx := mocks.NewMockContext()

			// Execute
			err := baseService.Start(mockCtx)

			// Verify
			if (err != nil) != tt.expectedError {
				t.Errorf("Start() error = %v, expectedError %v", err, tt.expectedError)
			}

			if baseService.Status() != tt.expectedFinal {
				t.Errorf("Status after Start() = %v, want %v", baseService.Status(), tt.expectedFinal)
			}
		})
	}
}

func TestBaseService_Stop(t *testing.T) {
	tests := []struct {
		name          string
		initialStatus ServiceStatus
		expectedError bool
		expectedFinal ServiceStatus
	}{
		{
			name:          "Stop when running",
			initialStatus: StatusRunning,
			expectedError: false,
			expectedFinal: StatusStopped,
		},
		{
			name:          "Stop when already stopped",
			initialStatus: StatusStopped,
			expectedError: false,
			expectedFinal: StatusStopped,
		},
		{
			name:          "Stop when stopping",
			initialStatus: StatusStopping,
			expectedError: false,
			expectedFinal: StatusStopping,
		},
		{
			name:          "Stop when starting",
			initialStatus: StatusStarting,
			expectedError: true,
			expectedFinal: StatusStarting,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockComp := mocks.NewMockComponent()
			baseService := CreateBaseService(mockComp)
			baseService.SetStatus(tt.initialStatus)
			mockCtx := mocks.NewMockContext()

			// Execute
			err := baseService.Stop(mockCtx)

			// Verify
			if (err != nil) != tt.expectedError {
				t.Errorf("Stop() error = %v, expectedError %v", err, tt.expectedError)
			}

			if baseService.Status() != tt.expectedFinal {
				t.Errorf("Status after Stop() = %v, want %v", baseService.Status(), tt.expectedFinal)
			}
		})
	}
}

func TestBaseService_Status(t *testing.T) {
	// Setup
	mockComp := mocks.NewMockComponent()
	baseService := CreateBaseService(mockComp)

	tests := []ServiceStatus{
		StatusStopped,
		StatusStarting,
		StatusRunning,
		StatusStopping,
		StatusFailed,
	}

	for _, status := range tests {
		t.Run(string(status), func(t *testing.T) {
			// Execute
			baseService.SetStatus(status)

			// Verify
			if baseService.Status() != status {
				t.Errorf("Status() = %v, want %v", baseService.Status(), status)
			}
		})
	}
}

func TestBaseService_SetStatus(t *testing.T) {
	// Setup
	mockComp := mocks.NewMockComponent()
	baseService := CreateBaseService(mockComp)

	tests := []ServiceStatus{
		StatusStopped,
		StatusStarting,
		StatusRunning,
		StatusStopping,
		StatusFailed,
	}

	for _, status := range tests {
		t.Run(string(status), func(t *testing.T) {
			// Execute
			baseService.SetStatus(status)

			// Verify
			if baseService.Status() != status {
				t.Errorf("Status() = %v, want %v", baseService.Status(), status)
			}
		})
	}
}

func TestBaseService_MarkAsFailed(t *testing.T) {
	// Setup
	mockComp := mocks.NewMockComponent()
	baseService := CreateBaseService(mockComp)

	// Initial status should be stopped
	if baseService.Status() != StatusStopped {
		t.Errorf("Initial status = %v, want %v", baseService.Status(), StatusStopped)
	}

	// Execute
	baseService.MarkAsFailed()

	// Verify
	if baseService.Status() != StatusFailed {
		t.Errorf("Status after MarkAsFailed() = %v, want %v", baseService.Status(), StatusFailed)
	}
}
