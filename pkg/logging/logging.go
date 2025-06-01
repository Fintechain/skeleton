// Package logging provides logging interfaces and types.
package logging

import (
	"github.com/fintechain/skeleton/internal/domain/logging"
	loggingImpl "github.com/fintechain/skeleton/internal/infrastructure/logging"
)

// Re-export logging interface
type Logger = logging.Logger

// Re-export logging error constants
const (
	ErrLoggerNotAvailable = logging.ErrLoggerNotAvailable
	ErrInvalidLogLevel    = logging.ErrInvalidLogLevel
)

// Re-export logging types
type LogLevel = loggingImpl.LogLevel

// Re-export log level constants
const (
	DebugLevel = loggingImpl.DebugLevel
	InfoLevel  = loggingImpl.InfoLevel
	WarnLevel  = loggingImpl.WarnLevel
	ErrorLevel = loggingImpl.ErrorLevel
)

// NewLogger creates a new Logger instance.
// This factory function provides access to the concrete logger implementation.
func NewLogger() Logger {
	return loggingImpl.NewLogger()
}

// NewLoggerWithLevel creates a new Logger with a specific log level.
func NewLoggerWithLevel(level LogLevel) Logger {
	return loggingImpl.NewLoggerWithLevel(level)
}

// NewLoggerWithPrefix creates a new Logger with a prefix.
func NewLoggerWithPrefix(prefix string) Logger {
	return loggingImpl.NewLoggerWithPrefix(prefix)
}

// NewConsoleLogger creates a new console logger.
func NewConsoleLogger() Logger {
	return loggingImpl.NewConsoleLogger()
}

// NewStructuredLogger creates a new structured logger.
func NewStructuredLogger() *loggingImpl.StructuredLogger {
	return loggingImpl.NewStructuredLogger()
}
