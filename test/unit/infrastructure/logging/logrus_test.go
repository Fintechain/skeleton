package logging

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/fintechain/skeleton/internal/domain/logging"
	loggingInfra "github.com/fintechain/skeleton/internal/infrastructure/logging"
)

func TestNewLogrusLogger(t *testing.T) {
	tests := []struct {
		name        string
		config      loggingInfra.LogrusConfig
		expectError bool
		errorMsg    string
		description string
	}{
		{
			name: "valid config with defaults",
			config: loggingInfra.LogrusConfig{
				Level:  "info",
				Format: "text",
			},
			expectError: false,
			description: "Should create logger with valid configuration",
		},
		{
			name: "valid config with json format",
			config: loggingInfra.LogrusConfig{
				Level:  "debug",
				Format: "json",
				Fields: map[string]interface{}{
					"service": "test",
					"version": "1.0.0",
				},
			},
			expectError: false,
			description: "Should create logger with JSON format and fields",
		},
		{
			name:        "empty config with defaults",
			config:      loggingInfra.LogrusConfig{},
			expectError: false,
			description: "Should create logger with default configuration",
		},
		{
			name: "invalid log level",
			config: loggingInfra.LogrusConfig{
				Level:  "invalid",
				Format: "text",
			},
			expectError: true,
			errorMsg:    logging.ErrInvalidLogLevel,
			description: "Should reject invalid log level",
		},
		{
			name: "invalid format",
			config: loggingInfra.LogrusConfig{
				Level:  "info",
				Format: "invalid",
			},
			expectError: true,
			errorMsg:    logging.ErrInvalidLogFormat,
			description: "Should reject invalid format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Use buffer for output capture
			var buf bytes.Buffer
			if tt.config.Output == nil {
				tt.config.Output = &buf
			}

			logger, err := loggingInfra.NewLogrusLogger(tt.config)

			if tt.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
				assert.Nil(t, logger)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, logger)

				// Verify interface compliance
				var _ logging.Logger = logger
			}
		})
	}
}

func TestLogrusLogger_LogLevels(t *testing.T) {
	tests := []struct {
		name        string
		level       string
		logFunc     func(logging.Logger, string, ...interface{})
		shouldLog   bool
		description string
	}{
		{
			name:        "debug level logs debug",
			level:       "debug",
			logFunc:     func(l logging.Logger, msg string, args ...interface{}) { l.Debug(msg, args...) },
			shouldLog:   true,
			description: "Debug level should log debug messages",
		},
		{
			name:        "debug level logs info",
			level:       "debug",
			logFunc:     func(l logging.Logger, msg string, args ...interface{}) { l.Info(msg, args...) },
			shouldLog:   true,
			description: "Debug level should log info messages",
		},
		{
			name:        "info level skips debug",
			level:       "info",
			logFunc:     func(l logging.Logger, msg string, args ...interface{}) { l.Debug(msg, args...) },
			shouldLog:   false,
			description: "Info level should skip debug messages",
		},
		{
			name:        "info level logs info",
			level:       "info",
			logFunc:     func(l logging.Logger, msg string, args ...interface{}) { l.Info(msg, args...) },
			shouldLog:   true,
			description: "Info level should log info messages",
		},
		{
			name:        "warn level logs warn",
			level:       "warn",
			logFunc:     func(l logging.Logger, msg string, args ...interface{}) { l.Warn(msg, args...) },
			shouldLog:   true,
			description: "Warn level should log warn messages",
		},
		{
			name:        "error level logs error",
			level:       "error",
			logFunc:     func(l logging.Logger, msg string, args ...interface{}) { l.Error(msg, args...) },
			shouldLog:   true,
			description: "Error level should log error messages",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			config := loggingInfra.LogrusConfig{
				Level:  tt.level,
				Format: "text",
				Output: &buf,
			}

			logger, err := loggingInfra.NewLogrusLogger(config)
			require.NoError(t, err)

			// Log a test message
			testMsg := "test message"
			tt.logFunc(logger, testMsg)

			output := buf.String()
			if tt.shouldLog {
				assert.Contains(t, output, testMsg)
			} else {
				assert.Empty(t, output)
			}
		})
	}
}

func TestLogrusLogger_StructuredLogging(t *testing.T) {
	tests := []struct {
		name        string
		format      string
		logFunc     func(logging.Logger)
		checkFunc   func(t *testing.T, output string)
		description string
	}{
		{
			name:   "key-value pairs",
			format: "json",
			logFunc: func(l logging.Logger) {
				l.Info("test message", "key1", "value1", "key2", 42)
			},
			checkFunc: func(t *testing.T, output string) {
				assert.Contains(t, output, "test message")
				assert.Contains(t, output, "key1")
				assert.Contains(t, output, "value1")
				assert.Contains(t, output, "key2")
				assert.Contains(t, output, "42")
			},
			description: "Should handle key-value pairs in structured logging",
		},
		{
			name:   "map arguments",
			format: "json",
			logFunc: func(l logging.Logger) {
				l.Error("error occurred", map[string]interface{}{
					"error_code": 500,
					"user_id":    "12345",
					"action":     "create_user",
				})
			},
			checkFunc: func(t *testing.T, output string) {
				assert.Contains(t, output, "error occurred")
				assert.Contains(t, output, "error_code")
				assert.Contains(t, output, "500")
				assert.Contains(t, output, "user_id")
				assert.Contains(t, output, "12345")
			},
			description: "Should handle map arguments in structured logging",
		},
		{
			name:   "mixed arguments",
			format: "json",
			logFunc: func(l logging.Logger) {
				l.Warn("warning", "status", "failed", map[string]interface{}{
					"retry_count": 3,
				}, "timestamp", "2023-01-01")
			},
			checkFunc: func(t *testing.T, output string) {
				assert.Contains(t, output, "warning")
				assert.Contains(t, output, "status")
				assert.Contains(t, output, "failed")
				assert.Contains(t, output, "retry_count")
				assert.Contains(t, output, "3")
				assert.Contains(t, output, "timestamp")
			},
			description: "Should handle mixed argument types",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			config := loggingInfra.LogrusConfig{
				Level:  "debug",
				Format: tt.format,
				Output: &buf,
			}

			logger, err := loggingInfra.NewLogrusLogger(config)
			require.NoError(t, err)

			tt.logFunc(logger)

			output := buf.String()
			tt.checkFunc(t, output)
		})
	}
}

func TestLogrusLogger_WithFields(t *testing.T) {
	var buf bytes.Buffer
	config := loggingInfra.LogrusConfig{
		Level:  "info",
		Format: "json",
		Output: &buf,
		Fields: map[string]interface{}{
			"service": "test-service",
			"version": "1.0.0",
		},
	}

	logger, err := loggingInfra.NewLogrusLogger(config)
	require.NoError(t, err)

	// Create logger with additional fields
	contextLogger := logger.WithFields(map[string]interface{}{
		"request_id": "req-123",
		"user_id":    "user-456",
	})

	contextLogger.Info("processing request")

	output := buf.String()
	assert.Contains(t, output, "processing request")
	assert.Contains(t, output, "service")
	assert.Contains(t, output, "test-service")
	assert.Contains(t, output, "version")
	assert.Contains(t, output, "1.0.0")
	assert.Contains(t, output, "request_id")
	assert.Contains(t, output, "req-123")
	assert.Contains(t, output, "user_id")
	assert.Contains(t, output, "user-456")
}

func TestLogrusLogger_Formats(t *testing.T) {
	tests := []struct {
		name        string
		format      string
		checkFunc   func(t *testing.T, output string)
		description string
	}{
		{
			name:   "text format",
			format: "text",
			checkFunc: func(t *testing.T, output string) {
				// Text format should contain timestamp and level
				assert.Contains(t, output, "level=info")
				assert.Contains(t, output, "test message")
			},
			description: "Should format logs as text",
		},
		{
			name:   "json format",
			format: "json",
			checkFunc: func(t *testing.T, output string) {
				// JSON format should be valid JSON
				assert.Contains(t, output, `"level":"info"`)
				assert.Contains(t, output, `"msg":"test message"`)
				assert.True(t, strings.HasPrefix(output, "{"))
			},
			description: "Should format logs as JSON",
		},
		{
			name:   "default format (text)",
			format: "",
			checkFunc: func(t *testing.T, output string) {
				// Default should be text format
				assert.Contains(t, output, "level=info")
				assert.Contains(t, output, "test message")
			},
			description: "Should default to text format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			config := loggingInfra.LogrusConfig{
				Level:  "info",
				Format: tt.format,
				Output: &buf,
			}

			logger, err := loggingInfra.NewLogrusLogger(config)
			require.NoError(t, err)

			logger.Info("test message")

			output := buf.String()
			tt.checkFunc(t, output)
		})
	}
}

func TestLogrusLogger_InterfaceCompliance(t *testing.T) {
	config := loggingInfra.LogrusConfig{
		Level:  "info",
		Format: "text",
		Output: &bytes.Buffer{},
	}

	logger, err := loggingInfra.NewLogrusLogger(config)
	require.NoError(t, err)

	// Verify interface compliance
	var _ logging.Logger = logger

	// Test all interface methods
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")
}

func TestLogrusLogger_DefaultOutput(t *testing.T) {
	// Test that default output is set when not provided
	config := loggingInfra.LogrusConfig{
		Level:  "info",
		Format: "text",
		// Output is nil, should default to os.Stdout
	}

	logger, err := loggingInfra.NewLogrusLogger(config)
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}

func TestLogrusLogger_LogLevelParsing(t *testing.T) {
	validLevels := []string{"panic", "fatal", "error", "warn", "info", "debug", "trace"}

	for _, level := range validLevels {
		t.Run("valid_level_"+level, func(t *testing.T) {
			config := loggingInfra.LogrusConfig{
				Level:  level,
				Format: "text",
				Output: &bytes.Buffer{},
			}

			logger, err := loggingInfra.NewLogrusLogger(config)
			assert.NoError(t, err)
			assert.NotNil(t, logger)
		})
	}
}
