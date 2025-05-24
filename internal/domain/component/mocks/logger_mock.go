package mocks

import (
	"fmt"

	"github.com/ebanfa/skeleton/internal/infrastructure/logging"
)

// MockLogger implements logging.Logger interface for testing
type MockLogger struct {
	LogEntries []LogEntry
}

// LogEntry represents a logged message
type LogEntry struct {
	Level   logging.LogLevel
	Message string
}

// NewMockLogger creates a new mock logger
func NewMockLogger() *MockLogger {
	return &MockLogger{
		LogEntries: make([]LogEntry, 0),
	}
}

// Log records a log entry
func (m *MockLogger) Log(level logging.LogLevel, format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	m.LogEntries = append(m.LogEntries, LogEntry{
		Level:   level,
		Message: message,
	})
}

// Debug logs a debug message
func (m *MockLogger) Debug(format string, args ...interface{}) {
	m.Log(logging.Debug, format, args...)
}

// Info logs an informational message
func (m *MockLogger) Info(format string, args ...interface{}) {
	m.Log(logging.Info, format, args...)
}

// Warn logs a warning message
func (m *MockLogger) Warn(format string, args ...interface{}) {
	m.Log(logging.Warn, format, args...)
}

// Error logs an error message
func (m *MockLogger) Error(format string, args ...interface{}) {
	m.Log(logging.Error, format, args...)
}

// Fatal logs a fatal message
func (m *MockLogger) Fatal(format string, args ...interface{}) {
	m.Log(logging.Fatal, format, args...)
}

// GetLogEntries returns all logged entries
func (m *MockLogger) GetLogEntries() []LogEntry {
	return m.LogEntries
}

// ClearLogEntries clears all log entries
func (m *MockLogger) ClearLogEntries() {
	m.LogEntries = make([]LogEntry, 0)
}

// GetLogEntriesByLevel returns log entries for a specific level
func (m *MockLogger) GetLogEntriesByLevel(level logging.LogLevel) []LogEntry {
	var filtered []LogEntry
	for _, entry := range m.LogEntries {
		if entry.Level == level {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

// ContainsMessage checks if any log entry contains the given substring
func (m *MockLogger) ContainsMessage(substring string) bool {
	for _, entry := range m.LogEntries {
		if fmt.Sprintf("%s", entry.Message) == substring {
			return true
		}
	}
	return false
}
