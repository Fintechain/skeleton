// Package logging provides infrastructure implementations for structured logging.
package logging

// NoOpLogger implements the logging.Logger interface with no-op operations.
// This is useful for testing scenarios where logging should be disabled.
type NoOpLogger struct{}

// NewNoOpLogger creates a new no-op logger.
// This follows dependency injection patterns by providing a constructor.
func NewNoOpLogger() *NoOpLogger {
	return &NoOpLogger{}
}

// Debug logs a debug-level message (no-op implementation).
func (l *NoOpLogger) Debug(msg string, args ...interface{}) {
	// No-op: silently discard the log message
}

// Info logs an info-level message (no-op implementation).
func (l *NoOpLogger) Info(msg string, args ...interface{}) {
	// No-op: silently discard the log message
}

// Warn logs a warning-level message (no-op implementation).
func (l *NoOpLogger) Warn(msg string, args ...interface{}) {
	// No-op: silently discard the log message
}

// Error logs an error-level message (no-op implementation).
func (l *NoOpLogger) Error(msg string, args ...interface{}) {
	// No-op: silently discard the log message
}
