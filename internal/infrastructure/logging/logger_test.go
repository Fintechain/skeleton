package logging

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// MockLogger is a mock implementation of LoggerBackend
type MockLogger struct {
	messages []string
}

func (m *MockLogger) Printf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	m.messages = append(m.messages, message)
}

// MockExit is a mock implementation of OSExitFunc
type MockExit struct {
	exitCalled bool
	exitCode   int
}

func (m *MockExit) Exit(code int) {
	m.exitCalled = true
	m.exitCode = code
}

// TestStandardLoggerCreation tests the creation of a StandardLogger
func TestStandardLoggerCreation(t *testing.T) {
	// Test with options
	mockLogger := &MockLogger{}
	mockExit := &MockExit{}

	logger := NewStandardLogger(StandardLoggerOptions{
		Backend: mockLogger,
		Level:   Info,
		Exit:    mockExit.Exit,
	})

	if logger == nil {
		t.Fatal("NewStandardLogger returned nil")
	}

	if logger.backend != mockLogger {
		t.Error("Backend not properly set")
	}

	if logger.level != Info {
		t.Errorf("Level not properly set, expected %v, got %v", Info, logger.level)
	}

	if logger.exit == nil {
		t.Error("Exit function not set")
	}

	// Test factory method
	factoryLogger := CreateStandardLogger(Warn)

	if factoryLogger == nil {
		t.Fatal("CreateStandardLogger returned nil")
	}

	if factoryLogger.level != Warn {
		t.Errorf("Level not properly set by factory, expected %v, got %v", Warn, factoryLogger.level)
	}

	// Test with writer
	var buf bytes.Buffer
	writerLogger := CreateStandardLoggerWithWriter(Error, &buf)

	if writerLogger == nil {
		t.Fatal("CreateStandardLoggerWithWriter returned nil")
	}

	if writerLogger.level != Error {
		t.Errorf("Level not properly set by writer factory, expected %v, got %v", Error, writerLogger.level)
	}
}

// TestStandardLoggerLevelFiltering tests that messages below the set level are filtered out
func TestStandardLoggerLevelFiltering(t *testing.T) {
	mockLogger := &MockLogger{}
	logger := NewStandardLogger(StandardLoggerOptions{
		Backend: mockLogger,
		Level:   Warn,
	})

	// These should be filtered out
	logger.Debug("debug message")
	logger.Info("info message")

	// These should pass through
	logger.Warn("warn message")
	logger.Error("error message")

	if len(mockLogger.messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(mockLogger.messages))
	}

	if !strings.Contains(mockLogger.messages[0], "[WARN]") {
		t.Errorf("Expected warn message, got: %s", mockLogger.messages[0])
	}

	if !strings.Contains(mockLogger.messages[1], "[ERROR]") {
		t.Errorf("Expected error message, got: %s", mockLogger.messages[1])
	}
}

// TestStandardLoggerMethods tests each logging method
func TestStandardLoggerMethods(t *testing.T) {
	mockLogger := &MockLogger{}
	logger := NewStandardLogger(StandardLoggerOptions{
		Backend: mockLogger,
		Level:   Debug,
	})

	// Test each method
	logger.Debug("debug %s", "test")
	logger.Info("info %s", "test")
	logger.Warn("warn %s", "test")
	logger.Error("error %s", "test")

	expectedPrefixes := []string{"[DEBUG]", "[INFO]", "[WARN]", "[ERROR]"}

	if len(mockLogger.messages) != 4 {
		t.Errorf("Expected 4 messages, got %d", len(mockLogger.messages))
	}

	for i, prefix := range expectedPrefixes {
		if !strings.Contains(mockLogger.messages[i], prefix) {
			t.Errorf("Message %d expected to contain prefix %s, got: %s", i, prefix, mockLogger.messages[i])
		}

		if !strings.Contains(mockLogger.messages[i], "test") {
			t.Errorf("Message %d expected to contain 'test', got: %s", i, mockLogger.messages[i])
		}
	}
}

// TestStandardLoggerFatal tests the Fatal method
func TestStandardLoggerFatal(t *testing.T) {
	mockLogger := &MockLogger{}
	mockExit := &MockExit{}

	logger := NewStandardLogger(StandardLoggerOptions{
		Backend: mockLogger,
		Level:   Debug,
		Exit:    mockExit.Exit,
	})

	// Test Fatal method
	logger.Fatal("fatal %s", "error")

	if !mockExit.exitCalled {
		t.Error("Exit not called on Fatal")
	}

	if mockExit.exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", mockExit.exitCode)
	}

	if len(mockLogger.messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(mockLogger.messages))
	}

	if !strings.Contains(mockLogger.messages[0], "[FATAL]") {
		t.Errorf("Expected fatal message, got: %s", mockLogger.messages[0])
	}

	if !strings.Contains(mockLogger.messages[0], "fatal error") {
		t.Errorf("Expected 'fatal error' in message, got: %s", mockLogger.messages[0])
	}
}

// TestStandardLoggerLevelPrefix tests the levelPrefix method
func TestStandardLoggerLevelPrefix(t *testing.T) {
	logger := CreateStandardLogger(Debug)

	testCases := []struct {
		level    LogLevel
		expected string
	}{
		{Debug, "[DEBUG]"},
		{Info, "[INFO]"},
		{Warn, "[WARN]"},
		{Error, "[ERROR]"},
		{Fatal, "[FATAL]"},
		{LogLevel(99), "[UNKNOWN]"},
	}

	for _, tc := range testCases {
		prefix := logger.levelPrefix(tc.level)
		if prefix != tc.expected {
			t.Errorf("For level %v, expected prefix %s, got %s", tc.level, tc.expected, prefix)
		}
	}
}
