// Package logging provides infrastructure implementations for structured logging.
package logging

import (
	"errors"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/logging"
	infraComponent "github.com/fintechain/skeleton/internal/infrastructure/component"
)

// Logger implements the LoggerService interface by wrapping a Logger
// and providing service lifecycle management.
type Logger struct {
	*infraComponent.BaseService
	logger logging.Logger
}

// NewLogger creates a new logger service with the provided logger.
func NewLogger(config component.ComponentConfig, logger logging.Logger) (*Logger, error) {
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}

	return &Logger{
		BaseService: infraComponent.NewBaseService(config),
		logger:      logger,
	}, nil
}

// Debug logs a debug-level message with optional structured data.
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.logger.Debug(msg, args...)
}

// Info logs an info-level message with optional structured data.
func (l *Logger) Info(msg string, args ...interface{}) {
	l.logger.Info(msg, args...)
}

// Warn logs a warning-level message with optional structured data.
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.logger.Warn(msg, args...)
}

// Error logs an error-level message with optional structured data.
func (l *Logger) Error(msg string, args ...interface{}) {
	l.logger.Error(msg, args...)
}

// Start begins the logger service operation.
func (l *Logger) Start(ctx context.Context) error {
	return l.BaseService.Start(ctx)
}

// Stop ends the logger service operation.
func (l *Logger) Stop(ctx context.Context) error {
	return l.BaseService.Stop(ctx)
}
