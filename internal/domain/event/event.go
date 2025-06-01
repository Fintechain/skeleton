// Package event provides event bus functionality for component communication.
package event

import (
	"time"
)

// Event represents an event in the system.
type Event struct {
	Topic   string                 // Event topic/type
	Source  string                 // Component that generated the event
	Time    time.Time              // When the event occurred
	Payload map[string]interface{} // Event data
}

// EventHandler is a callback function that processes events.
type EventHandler func(event *Event)

// Subscription represents a subscription to events of a specific topic.
type Subscription interface {
	// Cancel cancels the subscription.
	Cancel()

	// Topic returns the topic of the subscription.
	Topic() string
}

// EventBus facilitates component communication.
type EventBus interface {
	// Publication
	Publish(topic string, data interface{})

	// Subscription
	Subscribe(topic string, handler EventHandler) Subscription
	SubscribeAsync(topic string, handler EventHandler) Subscription

	// Control
	WaitAsync()
}
