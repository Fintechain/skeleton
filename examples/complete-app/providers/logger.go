// Package providers contains custom implementations of framework services
// for demonstration purposes.
package providers

import (
	"fmt"
	"time"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/logging"
)

// CustomLogger demonstrates a custom logging implementation
type CustomLogger struct {
	prefix string
	status component.ServiceStatus
}

// NewCustomLogger creates a new custom logger instance
func NewCustomLogger() logging.LoggerService {
	return &CustomLogger{
		prefix: "[CUSTOM]",
		status: component.StatusStopped,
	}
}

func (l *CustomLogger) ID() component.ComponentID {
	return "custom-logger"
}

func (l *CustomLogger) Name() string {
	return "Custom Logger"
}

func (l *CustomLogger) Type() component.ComponentType {
	return component.TypeService
}

func (l *CustomLogger) Description() string {
	return "Custom structured logger implementation"
}

func (l *CustomLogger) Version() string {
	return "1.0.0"
}

func (l *CustomLogger) Metadata() component.Metadata {
	return component.Metadata{
		"provider": "custom",
		"features": []string{"structured", "timestamped"},
	}
}

func (l *CustomLogger) Initialize(ctx context.Context, system component.System) error {
	fmt.Printf("%s Custom logger initialized\n", l.prefix)
	return nil
}

func (l *CustomLogger) Dispose() error {
	fmt.Printf("%s Custom logger disposed\n", l.prefix)
	return nil
}

func (l *CustomLogger) Start(ctx context.Context) error {
	fmt.Printf("%s Custom logger started\n", l.prefix)
	l.status = component.StatusRunning
	return nil
}

func (l *CustomLogger) Stop(ctx context.Context) error {
	fmt.Printf("%s Custom logger stopped\n", l.prefix)
	l.status = component.StatusStopped
	return nil
}

func (l *CustomLogger) IsRunning() bool {
	return l.status == component.StatusRunning
}

func (l *CustomLogger) Status() component.ServiceStatus {
	return l.status
}

func (l *CustomLogger) Debug(msg string, fields ...interface{}) {
	l.log("DEBUG", msg, fields...)
}

func (l *CustomLogger) Info(msg string, fields ...interface{}) {
	l.log("INFO", msg, fields...)
}

func (l *CustomLogger) Warn(msg string, fields ...interface{}) {
	l.log("WARN", msg, fields...)
}

func (l *CustomLogger) Error(msg string, fields ...interface{}) {
	l.log("ERROR", msg, fields...)
}

func (l *CustomLogger) log(level, msg string, fields ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s [%s] %s %s", l.prefix, level, timestamp, msg)
	if len(fields) > 0 {
		fmt.Printf(" %+v", fields)
	}
	fmt.Println()
}
