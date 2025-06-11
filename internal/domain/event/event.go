// Package event provides event-driven communication capabilities for Fintechain Skeleton.
package event

import (
	"time"

	"github.com/fintechain/skeleton/internal/domain/component"
)

// Event represents a structured message in the event-driven communication system.
type Event struct {
	// Topic identifies the type and category of the event.
	Topic string

	// Source identifies the component that generated this event.
	Source string

	// Time indicates when the event occurred.
	Time time.Time

	// Payload contains event-specific data as key-value pairs.
	Payload map[string]interface{}
}

// EventHandler defines the callback function signature for processing events.
type EventHandler func(event *Event)

// Subscription represents a managed event listener with lifecycle control.
type Subscription interface {
	// Cancel cancels the subscription and stops event delivery.
	Cancel()

	// Topic returns the topic pattern this subscription is registered for.
	Topic() string
}

// EventBus provides publish-subscribe event functionality.
// This is the core event bus interface without lifecycle management.
type EventBus interface {
	// Publication - strongly typed
	Publish(event *Event) error
	PublishAsync(event *Event) error

	// Subscription
	Subscribe(eventType string, handler EventHandler) Subscription
	SubscribeAsync(eventType string, handler EventHandler) Subscription

	// Control
	WaitAsync()
}

// EventBusService provides publish-subscribe event functionality as an infrastructure service.
// It combines the core event bus functionality with service lifecycle management.
type EventBusService interface {
	component.Service
	EventBus
}
