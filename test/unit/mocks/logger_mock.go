package mocks

import (
	"fmt"
	"sync"
)

// LogLevel represents the severity level of a log message.
type LogLevel int

const (
	// LogLevelDebug represents debug level logging.
	LogLevelDebug LogLevel = iota
	// LogLevelInfo represents info level logging.
	LogLevelInfo
	// LogLevelWarn represents warning level logging.
	LogLevelWarn
	// LogLevelError represents error level logging.
	LogLevelError
)

// String returns the string representation of the log level.
func (l LogLevel) String() string {
	switch l {
	case LogLevelDebug:
		return "DEBUG"
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarn:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// LogEntry represents a single log entry.
type LogEntry struct {
	Level   LogLevel
	Message string
	Args    []interface{}
}

// MockLogger provides a configurable mock implementation of a logger interface.
type MockLogger struct {
	mu sync.RWMutex

	// Configuration
	level       LogLevel
	shouldPanic bool

	// State tracking
	entries   []LogEntry
	callCount map[string]int
}

// NewMockLogger creates a new configurable logger mock.
func NewMockLogger() *MockLogger {
	return &MockLogger{
		level:     LogLevelDebug,
		entries:   make([]LogEntry, 0),
		callCount: make(map[string]int),
	}
}

// Logger Interface Implementation

// Debug logs a debug message.
func (l *MockLogger) Debug(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.trackCall("Debug")
	if l.level <= LogLevelDebug {
		l.addEntry(LogLevelDebug, format, args...)
	}
}

// Info logs an info message.
func (l *MockLogger) Info(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.trackCall("Info")
	if l.level <= LogLevelInfo {
		l.addEntry(LogLevelInfo, format, args...)
	}
}

// Warn logs a warning message.
func (l *MockLogger) Warn(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.trackCall("Warn")
	if l.level <= LogLevelWarn {
		l.addEntry(LogLevelWarn, format, args...)
	}
}

// Error logs an error message.
func (l *MockLogger) Error(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.trackCall("Error")
	if l.shouldPanic {
		panic(fmt.Sprintf(format, args...))
	}
	if l.level <= LogLevelError {
		l.addEntry(LogLevelError, format, args...)
	}
}

// Mock Configuration Methods

// SetLevel sets the minimum log level.
func (l *MockLogger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
}

// SetShouldPanic configures the mock to panic on error logs.
func (l *MockLogger) SetShouldPanic(shouldPanic bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.shouldPanic = shouldPanic
}

// State Verification Methods

// GetEntries returns all logged entries.
func (l *MockLogger) GetEntries() []LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	entries := make([]LogEntry, len(l.entries))
	copy(entries, l.entries)
	return entries
}

// GetEntriesByLevel returns all logged entries for a specific level.
func (l *MockLogger) GetEntriesByLevel(level LogLevel) []LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var filtered []LogEntry
	for _, entry := range l.entries {
		if entry.Level == level {
			filtered = append(filtered, entry)
		}
	}
	return filtered
}

// GetCallCount returns the number of times a method was called.
func (l *MockLogger) GetCallCount(method string) int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.callCount[method]
}

// GetLastEntry returns the last logged entry, or nil if no entries exist.
func (l *MockLogger) GetLastEntry() *LogEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if len(l.entries) == 0 {
		return nil
	}
	return &l.entries[len(l.entries)-1]
}

// HasEntry checks if an entry with the given level and message exists.
func (l *MockLogger) HasEntry(level LogLevel, message string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	for _, entry := range l.entries {
		if entry.Level == level && entry.Message == message {
			return true
		}
	}
	return false
}

// HasEntryContaining checks if an entry with the given level and message substring exists.
func (l *MockLogger) HasEntryContaining(level LogLevel, substring string) bool {
	l.mu.RLock()
	defer l.mu.RUnlock()

	for _, entry := range l.entries {
		if entry.Level == level {
			formatted := fmt.Sprintf(entry.Message, entry.Args...)
			if contains(formatted, substring) {
				return true
			}
		}
	}
	return false
}

// Clear removes all logged entries and resets call counts.
func (l *MockLogger) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.entries = make([]LogEntry, 0)
	l.callCount = make(map[string]int)
}

// Reset clears all mock state and configuration.
func (l *MockLogger) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.level = LogLevelDebug
	l.shouldPanic = false
	l.entries = make([]LogEntry, 0)
	l.callCount = make(map[string]int)
}

// Helper Methods

// addEntry adds a new log entry.
func (l *MockLogger) addEntry(level LogLevel, format string, args ...interface{}) {
	entry := LogEntry{
		Level:   level,
		Message: format,
		Args:    args,
	}
	l.entries = append(l.entries, entry)
}

// trackCall records a method call for verification.
func (l *MockLogger) trackCall(method string) {
	l.callCount[method]++
}

// contains checks if a string contains a substring (simple implementation).
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (len(substr) == 0 || findSubstring(s, substr))
}

// findSubstring finds a substring in a string.
func findSubstring(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// LoggerMockBuilder provides a fluent interface for configuring logger mocks.
type LoggerMockBuilder struct {
	mock *MockLogger
}

// NewLoggerMockBuilder creates a new logger mock builder.
func NewLoggerMockBuilder() *LoggerMockBuilder {
	return &LoggerMockBuilder{
		mock: NewMockLogger(),
	}
}

// WithLevel sets the minimum log level.
func (b *LoggerMockBuilder) WithLevel(level LogLevel) *LoggerMockBuilder {
	b.mock.SetLevel(level)
	return b
}

// WithPanicOnError configures the mock to panic on error logs.
func (b *LoggerMockBuilder) WithPanicOnError(shouldPanic bool) *LoggerMockBuilder {
	b.mock.SetShouldPanic(shouldPanic)
	return b
}

// Build returns the configured logger mock.
func (b *LoggerMockBuilder) Build() *MockLogger {
	return b.mock
}
