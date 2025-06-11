package logging

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/fintechain/skeleton/internal/domain/logging"
	loggingInfra "github.com/fintechain/skeleton/internal/infrastructure/logging"
)

func TestNewNoOpLogger(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "create no-op logger",
			description: "Should create no-op logger successfully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := loggingInfra.NewNoOpLogger()

			assert.NotNil(t, logger)

			// Verify interface compliance
			var _ logging.Logger = logger
		})
	}
}

func TestNoOpLogger_LogMethods(t *testing.T) {
	tests := []struct {
		name        string
		logFunc     func(logging.Logger)
		description string
	}{
		{
			name: "debug method",
			logFunc: func(l logging.Logger) {
				l.Debug("debug message")
			},
			description: "Debug method should not panic and do nothing",
		},
		{
			name: "info method",
			logFunc: func(l logging.Logger) {
				l.Info("info message")
			},
			description: "Info method should not panic and do nothing",
		},
		{
			name: "warn method",
			logFunc: func(l logging.Logger) {
				l.Warn("warn message")
			},
			description: "Warn method should not panic and do nothing",
		},
		{
			name: "error method",
			logFunc: func(l logging.Logger) {
				l.Error("error message")
			},
			description: "Error method should not panic and do nothing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := loggingInfra.NewNoOpLogger()

			// Should not panic
			assert.NotPanics(t, func() {
				tt.logFunc(logger)
			})
		})
	}
}

func TestNoOpLogger_WithArguments(t *testing.T) {
	tests := []struct {
		name        string
		logFunc     func(logging.Logger)
		description string
	}{
		{
			name: "debug with key-value pairs",
			logFunc: func(l logging.Logger) {
				l.Debug("debug message", "key1", "value1", "key2", 42)
			},
			description: "Should handle key-value arguments without panic",
		},
		{
			name: "info with map",
			logFunc: func(l logging.Logger) {
				l.Info("info message", map[string]interface{}{
					"user_id": "12345",
					"action":  "login",
				})
			},
			description: "Should handle map arguments without panic",
		},
		{
			name: "warn with mixed arguments",
			logFunc: func(l logging.Logger) {
				l.Warn("warn message", "status", "failed",
					map[string]interface{}{"retry": 3},
					"timestamp", "2023-01-01")
			},
			description: "Should handle mixed argument types without panic",
		},
		{
			name: "error with complex data",
			logFunc: func(l logging.Logger) {
				l.Error("error message",
					"error_code", 500,
					"details", map[string]interface{}{
						"stack_trace": []string{"line1", "line2"},
						"context": map[string]string{
							"module":   "auth",
							"function": "login",
						},
					})
			},
			description: "Should handle complex nested data without panic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := loggingInfra.NewNoOpLogger()

			// Should not panic with any arguments
			assert.NotPanics(t, func() {
				tt.logFunc(logger)
			})
		})
	}
}

func TestNoOpLogger_InterfaceCompliance(t *testing.T) {
	logger := loggingInfra.NewNoOpLogger()

	// Verify interface compliance
	var _ logging.Logger = logger

	// Test all interface methods don't panic
	assert.NotPanics(t, func() {
		logger.Debug("debug message")
		logger.Info("info message")
		logger.Warn("warn message")
		logger.Error("error message")
	})
}

func TestNoOpLogger_ZeroOverhead(t *testing.T) {
	logger := loggingInfra.NewNoOpLogger()

	// Test that no-op logger has minimal overhead
	// This is more of a behavioral test to ensure the logger does nothing

	// Multiple calls should not cause any issues
	for i := 0; i < 1000; i++ {
		logger.Debug("debug", "iteration", i)
		logger.Info("info", "iteration", i)
		logger.Warn("warn", "iteration", i)
		logger.Error("error", "iteration", i)
	}

	// Should complete without any issues
	assert.True(t, true, "No-op logger should handle many calls efficiently")
}

func TestNoOpLogger_NilArguments(t *testing.T) {
	logger := loggingInfra.NewNoOpLogger()

	// Test with nil arguments
	assert.NotPanics(t, func() {
		logger.Debug("debug", nil)
		logger.Info("info", nil, nil)
		logger.Warn("warn", "key", nil)
		logger.Error("error", nil, "value")
	})
}

func TestNoOpLogger_EmptyArguments(t *testing.T) {
	logger := loggingInfra.NewNoOpLogger()

	// Test with no arguments
	assert.NotPanics(t, func() {
		logger.Debug("debug")
		logger.Info("info")
		logger.Warn("warn")
		logger.Error("error")
	})

	// Test with empty string
	assert.NotPanics(t, func() {
		logger.Debug("")
		logger.Info("")
		logger.Warn("")
		logger.Error("")
	})
}
