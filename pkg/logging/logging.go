// Package logging provides logging system interfaces and implementations.
package logging

import (
	"github.com/fintechain/skeleton/internal/domain/logging"
	infraLogging "github.com/fintechain/skeleton/internal/infrastructure/logging"
)

// Core interfaces
type Logger = logging.Logger
type LoggerService = logging.LoggerService

// Factory functions
var NewNoOpLogger = infraLogging.NewNoOpLogger
var NewLogrusLogger = infraLogging.NewLogrusLogger

// LogrusConfig for creating Logrus loggers
type LogrusConfig = infraLogging.LogrusConfig

// Error constants
const (
	ErrLoggerNotAvailable   = logging.ErrLoggerNotAvailable
	ErrInvalidLogLevel      = logging.ErrInvalidLogLevel
	ErrLoggerNotFound       = logging.ErrLoggerNotFound
	ErrLoggerExists         = logging.ErrLoggerExists
	ErrLoggerNotInitialized = logging.ErrLoggerNotInitialized
	ErrLoggerClosed         = logging.ErrLoggerClosed
	ErrInvalidLogFormat     = logging.ErrInvalidLogFormat
	ErrLogWriteFailed       = logging.ErrLogWriteFailed
	ErrInvalidLogConfig     = logging.ErrInvalidLogConfig
)
