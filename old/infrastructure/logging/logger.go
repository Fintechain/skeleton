// Package logging provides logging facilities for the component system.
package logging

import (
	"fmt"
	"io"
	"log"
	"os"
)

// LogLevel represents the severity level of a log message.
type LogLevel int

const (
	// Debug level for detailed troubleshooting
	Debug LogLevel = iota
	// Info level for general operational information
	Info
	// Warn level for warnings
	Warn
	// Error level for errors
	Error
	// Fatal level for fatal errors that cause the program to exit
	Fatal
)

// Logger provides logging functionality.
type Logger interface {
	// Log logs a message at the specified level.
	Log(level LogLevel, format string, args ...interface{})

	// Debug logs a debug message.
	Debug(format string, args ...interface{})

	// Info logs an informational message.
	Info(format string, args ...interface{})

	// Warn logs a warning message.
	Warn(format string, args ...interface{})

	// Error logs an error message.
	Error(format string, args ...interface{})

	// Fatal logs a fatal message and exits the program.
	Fatal(format string, args ...interface{})
}

// LoggerBackend is an interface for the backend implementation of logging.
// This isolates our code from the external dependency (standard log package).
type LoggerBackend interface {
	// Printf logs a formatted message.
	Printf(format string, v ...interface{})
}

// OSExitFunc defines a function type for os.Exit to allow for testing.
type OSExitFunc func(code int)

// StandardLogger is a simple implementation of Logger that uses the standard log package.
type StandardLogger struct {
	backend LoggerBackend
	level   LogLevel
	exit    OSExitFunc
}

// StandardLoggerOptions contains options for creating a StandardLogger.
type StandardLoggerOptions struct {
	Backend LoggerBackend
	Level   LogLevel
	Exit    OSExitFunc
}

// NewStandardLogger creates a new standard logger with the specified options.
// This follows the constructor injection pattern for dependencies.
func NewStandardLogger(options StandardLoggerOptions) *StandardLogger {
	// Set defaults for optional dependencies
	backend := options.Backend
	if backend == nil {
		backend = log.New(os.Stdout, "", log.LstdFlags)
	}

	exit := options.Exit
	if exit == nil {
		exit = os.Exit
	}

	return &StandardLogger{
		backend: backend,
		level:   options.Level,
		exit:    exit,
	}
}

// CreateStandardLogger is a factory method for backward compatibility.
// Creates a standard logger with default settings.
func CreateStandardLogger(level LogLevel) *StandardLogger {
	return NewStandardLogger(StandardLoggerOptions{
		Level: level,
	})
}

// CreateStandardLoggerWithWriter is a factory method that creates a logger with a specific writer.
func CreateStandardLoggerWithWriter(level LogLevel, writer io.Writer) *StandardLogger {
	return NewStandardLogger(StandardLoggerOptions{
		Level:   level,
		Backend: log.New(writer, "", log.LstdFlags),
	})
}

// Log logs a message at the specified level.
func (l *StandardLogger) Log(level LogLevel, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	prefix := l.levelPrefix(level)
	message := fmt.Sprintf(format, args...)
	l.backend.Printf("%s %s", prefix, message)

	if level == Fatal {
		l.exit(1)
	}
}

// Debug logs a debug message.
func (l *StandardLogger) Debug(format string, args ...interface{}) {
	l.Log(Debug, format, args...)
}

// Info logs an informational message.
func (l *StandardLogger) Info(format string, args ...interface{}) {
	l.Log(Info, format, args...)
}

// Warn logs a warning message.
func (l *StandardLogger) Warn(format string, args ...interface{}) {
	l.Log(Warn, format, args...)
}

// Error logs an error message.
func (l *StandardLogger) Error(format string, args ...interface{}) {
	l.Log(Error, format, args...)
}

// Fatal logs a fatal message and exits the program.
func (l *StandardLogger) Fatal(format string, args ...interface{}) {
	l.Log(Fatal, format, args...)
}

// levelPrefix returns the prefix for the log level.
func (l *StandardLogger) levelPrefix(level LogLevel) string {
	switch level {
	case Debug:
		return "[DEBUG]"
	case Info:
		return "[INFO]"
	case Warn:
		return "[WARN]"
	case Error:
		return "[ERROR]"
	case Fatal:
		return "[FATAL]"
	default:
		return "[UNKNOWN]"
	}
}
