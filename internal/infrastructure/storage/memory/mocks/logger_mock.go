// Package mocks provides mock implementations for external dependencies
package mocks

import (
	"github.com/fintechain/skeleton/internal/infrastructure/logging"
)

// MockLogger is a mock implementation of the logging.Logger interface
type MockLogger struct {
	// Function implementations
	LogFunc   func(level logging.LogLevel, format string, args ...interface{})
	DebugFunc func(format string, args ...interface{})
	InfoFunc  func(format string, args ...interface{})
	WarnFunc  func(format string, args ...interface{})
	ErrorFunc func(format string, args ...interface{})
	FatalFunc func(format string, args ...interface{})

	// Call tracking
	LogCalls   []LogCall
	DebugCalls []DebugCall
	InfoCalls  []InfoCall
	WarnCalls  []WarnCall
	ErrorCalls []ErrorCall
	FatalCalls []FatalCall
}

// LogCall contains info about a Log call
type LogCall struct {
	Level  logging.LogLevel
	Format string
	Args   []interface{}
}

// DebugCall contains info about a Debug call
type DebugCall struct {
	Format string
	Args   []interface{}
}

// InfoCall contains info about an Info call
type InfoCall struct {
	Format string
	Args   []interface{}
}

// WarnCall contains info about a Warn call
type WarnCall struct {
	Format string
	Args   []interface{}
}

// ErrorCall contains info about an Error call
type ErrorCall struct {
	Format string
	Args   []interface{}
}

// FatalCall contains info about a Fatal call
type FatalCall struct {
	Format string
	Args   []interface{}
}

// Log implements logging.Logger.Log
func (m *MockLogger) Log(level logging.LogLevel, format string, args ...interface{}) {
	m.LogCalls = append(m.LogCalls, LogCall{Level: level, Format: format, Args: args})
	if m.LogFunc != nil {
		m.LogFunc(level, format, args...)
	}
}

// Debug implements logging.Logger.Debug
func (m *MockLogger) Debug(format string, args ...interface{}) {
	m.DebugCalls = append(m.DebugCalls, DebugCall{Format: format, Args: args})
	if m.DebugFunc != nil {
		m.DebugFunc(format, args...)
	}
}

// Info implements logging.Logger.Info
func (m *MockLogger) Info(format string, args ...interface{}) {
	m.InfoCalls = append(m.InfoCalls, InfoCall{Format: format, Args: args})
	if m.InfoFunc != nil {
		m.InfoFunc(format, args...)
	}
}

// Warn implements logging.Logger.Warn
func (m *MockLogger) Warn(format string, args ...interface{}) {
	m.WarnCalls = append(m.WarnCalls, WarnCall{Format: format, Args: args})
	if m.WarnFunc != nil {
		m.WarnFunc(format, args...)
	}
}

// Error implements logging.Logger.Error
func (m *MockLogger) Error(format string, args ...interface{}) {
	m.ErrorCalls = append(m.ErrorCalls, ErrorCall{Format: format, Args: args})
	if m.ErrorFunc != nil {
		m.ErrorFunc(format, args...)
	}
}

// Fatal implements logging.Logger.Fatal
func (m *MockLogger) Fatal(format string, args ...interface{}) {
	m.FatalCalls = append(m.FatalCalls, FatalCall{Format: format, Args: args})
	if m.FatalFunc != nil {
		m.FatalFunc(format, args...)
	}
}

// ClearCalls clears all tracked calls
func (m *MockLogger) ClearCalls() {
	m.LogCalls = nil
	m.DebugCalls = nil
	m.InfoCalls = nil
	m.WarnCalls = nil
	m.ErrorCalls = nil
	m.FatalCalls = nil
}

// NewMockLogger creates a new MockLogger
func NewMockLogger() *MockLogger {
	return &MockLogger{}
}
