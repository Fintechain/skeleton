// Package logging provides structured logging infrastructure for Fintechain Skeleton.
package logging

import (
	"github.com/fintechain/skeleton/internal/domain/component"
)

// Logger represents the logging facility for structured application logging.
type Logger interface {
	// Debug logs a debug-level message with optional structured data.
	Debug(msg string, args ...interface{})

	// Info logs an info-level message with optional structured data.
	Info(msg string, args ...interface{})

	// Warn logs a warning-level message with optional structured data.
	Warn(msg string, args ...interface{})

	// Error logs an error-level message with optional structured data.
	Error(msg string, args ...interface{})
}

// LoggerService provides structured logging functionality as an infrastructure service.
// It combines the core logging functionality with service lifecycle management.
type LoggerService interface {
	component.Service
	Logger
}

// Common logging error codes.
const (
	// ErrLoggerNotAvailable indicates that the logging system is not available.
	ErrLoggerNotAvailable = "logging.logger_not_available"

	// ErrInvalidLogLevel indicates that an invalid log level was specified.
	ErrInvalidLogLevel = "logging.invalid_log_level"
)
