// Package logging provides infrastructure implementations for structured logging.
package logging

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/fintechain/skeleton/internal/domain/logging"
)

// LogrusLogger implements the logging.Logger interface using Logrus.
type LogrusLogger struct {
	logger *logrus.Logger
	fields logrus.Fields
}

// LogrusConfig holds configuration for the Logrus logger.
type LogrusConfig struct {
	Level  string
	Format string
	Output io.Writer
	Fields map[string]interface{}
}

// NewLogrusLogger creates a new Logrus-based logger with dependency injection.
// All dependencies are injected via constructor parameters for testability.
func NewLogrusLogger(config LogrusConfig) (*LogrusLogger, error) {
	if config.Output == nil {
		config.Output = os.Stdout
	}

	logger := logrus.New()
	logger.SetOutput(config.Output)

	// Set log level
	if config.Level != "" {
		level, err := logrus.ParseLevel(config.Level)
		if err != nil {
			return nil, fmt.Errorf("%s: invalid log level '%s': %w", logging.ErrInvalidLogLevel, config.Level, err)
		}
		logger.SetLevel(level)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	// Set log format
	switch config.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	case "text", "":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	default:
		return nil, fmt.Errorf("%s: unsupported format '%s'", logging.ErrInvalidLogFormat, config.Format)
	}

	// Convert config fields to logrus.Fields
	fields := make(logrus.Fields)
	if config.Fields != nil {
		for k, v := range config.Fields {
			fields[k] = v
		}
	}

	return &LogrusLogger{
		logger: logger,
		fields: fields,
	}, nil
}

// Debug logs a debug-level message with optional structured data.
func (l *LogrusLogger) Debug(msg string, args ...interface{}) {
	entry := l.logger.WithFields(l.fields)
	if len(args) > 0 {
		entry = entry.WithFields(l.parseArgs(args...))
	}
	entry.Debug(msg)
}

// Info logs an info-level message with optional structured data.
func (l *LogrusLogger) Info(msg string, args ...interface{}) {
	entry := l.logger.WithFields(l.fields)
	if len(args) > 0 {
		entry = entry.WithFields(l.parseArgs(args...))
	}
	entry.Info(msg)
}

// Warn logs a warning-level message with optional structured data.
func (l *LogrusLogger) Warn(msg string, args ...interface{}) {
	entry := l.logger.WithFields(l.fields)
	if len(args) > 0 {
		entry = entry.WithFields(l.parseArgs(args...))
	}
	entry.Warn(msg)
}

// Error logs an error-level message with optional structured data.
func (l *LogrusLogger) Error(msg string, args ...interface{}) {
	entry := l.logger.WithFields(l.fields)
	if len(args) > 0 {
		entry = entry.WithFields(l.parseArgs(args...))
	}
	entry.Error(msg)
}

// WithFields creates a new logger instance with additional fields.
// This allows for contextual logging without modifying the original logger.
func (l *LogrusLogger) WithFields(fields map[string]interface{}) *LogrusLogger {
	newFields := make(logrus.Fields)

	// Copy existing fields
	for k, v := range l.fields {
		newFields[k] = v
	}

	// Add new fields
	for k, v := range fields {
		newFields[k] = v
	}

	return &LogrusLogger{
		logger: l.logger,
		fields: newFields,
	}
}

// parseArgs converts variadic arguments into logrus.Fields.
// Supports key-value pairs and maps.
func (l *LogrusLogger) parseArgs(args ...interface{}) logrus.Fields {
	fields := make(logrus.Fields)

	for i := 0; i < len(args); i++ {
		switch arg := args[i].(type) {
		case map[string]interface{}:
			// Handle map arguments
			for k, v := range arg {
				fields[k] = v
			}
		case string:
			// Handle key-value pairs
			if i+1 < len(args) {
				fields[arg] = args[i+1]
				i++ // Skip the value
			} else {
				fields["arg"] = arg
			}
		default:
			// Handle other types as indexed arguments
			fields[fmt.Sprintf("arg%d", i)] = arg
		}
	}

	return fields
}
