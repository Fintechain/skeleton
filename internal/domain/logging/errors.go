// Package logging provides interfaces and types for the logging system.
package logging

// Standard logging error codes
const (
	// ErrLoggerNotFound is returned when a logger doesn't exist
	ErrLoggerNotFound = "logging.logger_not_found"

	// ErrLoggerExists is returned when creating a logger that already exists
	ErrLoggerExists = "logging.logger_exists"

	// ErrLoggerNotInitialized is returned when operations are performed on a non-initialized logger
	ErrLoggerNotInitialized = "logging.logger_not_initialized"

	// ErrLoggerClosed is returned when operations are performed on a closed logger
	ErrLoggerClosed = "logging.logger_closed"

	// ErrInvalidLogFormat is returned when an invalid log format is provided
	ErrInvalidLogFormat = "logging.invalid_log_format"

	// ErrLogWriteFailed is returned when writing to log fails
	ErrLogWriteFailed = "logging.log_write_failed"

	// ErrInvalidLogConfig is returned when invalid logging configuration is provided
	ErrInvalidLogConfig = "logging.invalid_log_config"
)
