// Package event provides event system interfaces and implementations.
package event

import (
	"time"

	"github.com/fintechain/skeleton/internal/domain/event"
	infraEvent "github.com/fintechain/skeleton/internal/infrastructure/event"
)

// Core interfaces
type Event = event.Event
type EventHandler = event.EventHandler
type Subscription = event.Subscription
type EventBus = event.EventBus
type EventBusService = event.EventBusService

// Factory functions
var NewEventBus = infraEvent.NewEventBus

// Event constructor helper
func NewEvent(topic, source string, payload map[string]interface{}) *Event {
	return &Event{
		Topic:   topic,
		Source:  source,
		Time:    time.Now(),
		Payload: payload,
	}
}
