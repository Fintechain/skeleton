// Package logging provides logging infrastructure.
package logging

// Logger represents the logging facility.
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// Common logging error codes
const (
	ErrLoggerNotAvailable = "logging.logger_not_available"
	ErrInvalidLogLevel    = "logging.invalid_log_level"
)
