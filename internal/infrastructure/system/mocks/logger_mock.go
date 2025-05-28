package mocks

import (
	"github.com/fintechain/skeleton/internal/infrastructure/logging"
)

// MockLogger is a mock implementation of logging.Logger for testing
type MockLogger struct {
	// Function fields for customizing behavior
	LogFunc   func(logging.LogLevel, string, ...interface{})
	DebugFunc func(string, ...interface{})
	InfoFunc  func(string, ...interface{})
	WarnFunc  func(string, ...interface{})
	ErrorFunc func(string, ...interface{})
	FatalFunc func(string, ...interface{})

	// Call tracking
	LogCalls   []LogCall
	DebugCalls []string
	InfoCalls  []string
	WarnCalls  []string
	ErrorCalls []string
	FatalCalls []string

	// State
	Messages []LogMessage
}

type LogCall struct {
	Level  logging.LogLevel
	Format string
	Args   []interface{}
}

type LogMessage struct {
	Level   logging.LogLevel
	Message string
}

// NewMockLogger creates a new mock logger
func NewMockLogger() *MockLogger {
	return &MockLogger{}
}

// Log implements logging.Logger
func (m *MockLogger) Log(level logging.LogLevel, format string, args ...interface{}) {
	m.LogCalls = append(m.LogCalls, LogCall{Level: level, Format: format, Args: args})
	if m.LogFunc != nil {
		m.LogFunc(level, format, args...)
		return
	}

	// Default behavior: store message
	message := format
	if len(args) > 0 {
		// Simple formatting for testing
		message = format
	}
	m.Messages = append(m.Messages, LogMessage{Level: level, Message: message})
}

// Debug implements logging.Logger
func (m *MockLogger) Debug(format string, args ...interface{}) {
	m.DebugCalls = append(m.DebugCalls, format)
	if m.DebugFunc != nil {
		m.DebugFunc(format, args...)
		return
	}
	m.Log(logging.Debug, format, args...)
}

// Info implements logging.Logger
func (m *MockLogger) Info(format string, args ...interface{}) {
	m.InfoCalls = append(m.InfoCalls, format)
	if m.InfoFunc != nil {
		m.InfoFunc(format, args...)
		return
	}
	m.Log(logging.Info, format, args...)
}

// Warn implements logging.Logger
func (m *MockLogger) Warn(format string, args ...interface{}) {
	m.WarnCalls = append(m.WarnCalls, format)
	if m.WarnFunc != nil {
		m.WarnFunc(format, args...)
		return
	}
	m.Log(logging.Warn, format, args...)
}

// Error implements logging.Logger
func (m *MockLogger) Error(format string, args ...interface{}) {
	m.ErrorCalls = append(m.ErrorCalls, format)
	if m.ErrorFunc != nil {
		m.ErrorFunc(format, args...)
		return
	}
	m.Log(logging.Error, format, args...)
}

// Fatal implements logging.Logger
func (m *MockLogger) Fatal(format string, args ...interface{}) {
	m.FatalCalls = append(m.FatalCalls, format)
	if m.FatalFunc != nil {
		m.FatalFunc(format, args...)
		return
	}
	m.Log(logging.Fatal, format, args...)
}
