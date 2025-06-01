// Package logging provides public access to logging utilities for the skeleton framework.
// This package re-exports types and constructor functions from the internal logging implementation.
package logging

import (
	"io"

	"github.com/fintechain/skeleton/internal/infrastructure/logging"
)

// Re-export types from internal package
type (
	// Logger provides logging functionality.
	Logger = logging.Logger

	// LogLevel represents the severity level of a log message.
	LogLevel = logging.LogLevel

	// LoggerBackend is an interface for the backend implementation of logging.
	LoggerBackend = logging.LoggerBackend

	// OSExitFunc defines a function type for os.Exit to allow for testing.
	OSExitFunc = logging.OSExitFunc

	// StandardLogger is a simple implementation of Logger that uses the standard log package.
	StandardLogger = logging.StandardLogger

	// StandardLoggerOptions contains options for creating a StandardLogger.
	StandardLoggerOptions = logging.StandardLoggerOptions
)

// Re-export log level constants
const (
	// Debug level for detailed troubleshooting
	Debug = logging.Debug
	// Info level for general operational information
	Info = logging.Info
	// Warn level for warnings
	Warn = logging.Warn
	// Error level for errors
	Error = logging.Error
	// Fatal level for fatal errors that cause the program to exit
	Fatal = logging.Fatal
)

// Re-export constructor functions

// NewStandardLogger creates a new standard logger with the specified options.
func NewStandardLogger(options StandardLoggerOptions) *StandardLogger {
	return logging.NewStandardLogger(options)
}

// CreateStandardLogger is a factory method for backward compatibility.
func CreateStandardLogger(level LogLevel) *StandardLogger {
	return logging.CreateStandardLogger(level)
}

// CreateStandardLoggerWithWriter is a factory method that creates a logger with a specific writer.
func CreateStandardLoggerWithWriter(level LogLevel, writer io.Writer) *StandardLogger {
	return logging.CreateStandardLoggerWithWriter(level, writer)
}
