// Package event provides event bus functionality for component communication.
package event

import (
	"github.com/fintechain/skeleton/internal/domain/event"
	eventImpl "github.com/fintechain/skeleton/internal/infrastructure/event"
)

// Re-export event interfaces
type EventHandler = event.EventHandler
type Subscription = event.Subscription
type EventBus = event.EventBus

// Re-export event types
type Event = event.Event

// NewEventBus creates a new EventBus instance.
// This factory function provides access to the concrete event bus implementation.
func NewEventBus() EventBus {
	return eventImpl.NewEventBus()
}
