// Package logging provides concrete implementations of the logging system.
package logging

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/fintechain/skeleton/internal/domain/logging"
)

// LogLevel represents the logging level.
type LogLevel int

const (
	// DebugLevel for debug messages
	DebugLevel LogLevel = iota
	// InfoLevel for informational messages
	InfoLevel
	// WarnLevel for warning messages
	WarnLevel
	// ErrorLevel for error messages
	ErrorLevel
)

// DefaultLogger is a concrete implementation of the Logger interface.
type DefaultLogger struct {
	mu      sync.RWMutex
	level   LogLevel
	logger  *log.Logger
	prefix  string
	enabled bool
}

// NewLogger creates a new Logger instance with minimal dependencies.
func NewLogger() logging.Logger {
	return &DefaultLogger{
		level:   InfoLevel,
		logger:  log.New(os.Stdout, "", log.LstdFlags),
		enabled: true,
	}
}

// NewLoggerWithLevel creates a new Logger with a specific log level.
func NewLoggerWithLevel(level LogLevel) logging.Logger {
	return &DefaultLogger{
		level:   level,
		logger:  log.New(os.Stdout, "", log.LstdFlags),
		enabled: true,
	}
}

// NewLoggerWithPrefix creates a new Logger with a prefix.
func NewLoggerWithPrefix(prefix string) logging.Logger {
	return &DefaultLogger{
		level:   InfoLevel,
		logger:  log.New(os.Stdout, "", log.LstdFlags),
		prefix:  prefix,
		enabled: true,
	}
}

// SetLevel sets the logging level.
func (l *DefaultLogger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetEnabled enables or disables logging.
func (l *DefaultLogger) SetEnabled(enabled bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.enabled = enabled
}

// Debug logs a debug message.
func (l *DefaultLogger) Debug(msg string, args ...interface{}) {
	l.log(DebugLevel, "DEBUG", msg, args...)
}

// Info logs an informational message.
func (l *DefaultLogger) Info(msg string, args ...interface{}) {
	l.log(InfoLevel, "INFO", msg, args...)
}

// Warn logs a warning message.
func (l *DefaultLogger) Warn(msg string, args ...interface{}) {
	l.log(WarnLevel, "WARN", msg, args...)
}

// Error logs an error message.
func (l *DefaultLogger) Error(msg string, args ...interface{}) {
	l.log(ErrorLevel, "ERROR", msg, args...)
}

// log is the internal logging method.
func (l *DefaultLogger) log(level LogLevel, levelStr, msg string, args ...interface{}) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Check if logging is enabled and level is appropriate
	if !l.enabled || level < l.level {
		return
	}

	// Format the message
	formattedMsg := msg
	if len(args) > 0 {
		formattedMsg = fmt.Sprintf(msg, args...)
	}

	// Create the log entry
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s: %s", timestamp, levelStr, formattedMsg)

	// Add prefix if configured
	if l.prefix != "" {
		logEntry = fmt.Sprintf("[%s] %s", l.prefix, logEntry)
	}

	// Output the log entry
	l.logger.Println(logEntry)
}

// ConsoleLogger is a simple console-based logger implementation.
type ConsoleLogger struct {
	*DefaultLogger
}

// NewConsoleLogger creates a new console logger.
func NewConsoleLogger() logging.Logger {
	return &ConsoleLogger{
		DefaultLogger: &DefaultLogger{
			level:   InfoLevel,
			logger:  log.New(os.Stdout, "", 0), // No default timestamp since we add our own
			enabled: true,
		},
	}
}

// StructuredLogger provides structured logging capabilities.
type StructuredLogger struct {
	*DefaultLogger
	fields map[string]interface{}
	mu     sync.RWMutex
}

// NewStructuredLogger creates a new structured logger.
func NewStructuredLogger() *StructuredLogger {
	return &StructuredLogger{
		DefaultLogger: &DefaultLogger{
			level:   InfoLevel,
			logger:  log.New(os.Stdout, "", 0),
			enabled: true,
		},
		fields: make(map[string]interface{}),
	}
}

// WithField adds a field to the structured logger.
func (sl *StructuredLogger) WithField(key string, value interface{}) *StructuredLogger {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	newLogger := &StructuredLogger{
		DefaultLogger: sl.DefaultLogger,
		fields:        make(map[string]interface{}),
	}

	// Copy existing fields
	for k, v := range sl.fields {
		newLogger.fields[k] = v
	}

	// Add new field
	newLogger.fields[key] = value
	return newLogger
}

// WithFields adds multiple fields to the structured logger.
func (sl *StructuredLogger) WithFields(fields map[string]interface{}) *StructuredLogger {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	newLogger := &StructuredLogger{
		DefaultLogger: sl.DefaultLogger,
		fields:        make(map[string]interface{}),
	}

	// Copy existing fields
	for k, v := range sl.fields {
		newLogger.fields[k] = v
	}

	// Add new fields
	for k, v := range fields {
		newLogger.fields[k] = v
	}

	return newLogger
}

// log overrides the default log method to include structured fields.
func (sl *StructuredLogger) log(level LogLevel, levelStr, msg string, args ...interface{}) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	// Check if logging is enabled and level is appropriate
	if !sl.enabled || level < sl.level {
		return
	}

	// Format the message
	formattedMsg := msg
	if len(args) > 0 {
		formattedMsg = fmt.Sprintf(msg, args...)
	}

	// Create structured log entry
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logEntry := fmt.Sprintf("[%s] %s: %s", timestamp, levelStr, formattedMsg)

	// Add structured fields
	if len(sl.fields) > 0 {
		fieldsStr := ""
		for k, v := range sl.fields {
			if fieldsStr != "" {
				fieldsStr += ", "
			}
			fieldsStr += fmt.Sprintf("%s=%v", k, v)
		}
		logEntry += fmt.Sprintf(" {%s}", fieldsStr)
	}

	// Add prefix if configured
	if sl.prefix != "" {
		logEntry = fmt.Sprintf("[%s] %s", sl.prefix, logEntry)
	}

	// Output the log entry
	sl.logger.Println(logEntry)
}
