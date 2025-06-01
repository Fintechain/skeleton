// Package event provides public APIs for the event system.
package event

import (
	"time"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/infrastructure/event"
)

// ===== EVENT INTERFACES =====

// EventBus facilitates component communication.
type EventBus = event.EventBus

// EventHandler is a callback function that processes events.
type EventHandler = event.EventHandler

// Subscription represents a subscription to events of a specific topic.
type Subscription = event.Subscription

// ===== EVENT TYPES =====

// Event represents an event in the system.
type Event = event.Event

// ===== EVENT ERROR CONSTANTS =====

// Common event error codes
const (
	ErrEventPublishFailed   = "event.publish_failed"
	ErrEventSubscribeFailed = "event.subscribe_failed"
	ErrInvalidEventTopic    = "event.invalid_topic"
	ErrEventHandlerFailed   = "event.handler_failed"
	ErrEventBusUnavailable  = "event.bus_unavailable"
)

// ===== ERROR HANDLING =====

// Error represents a domain-specific error from the event system.
type Error = component.Error

// NewError creates a new event error with the given code, message, and optional cause.
func NewError(code, message string, cause error) *Error {
	return component.NewError(code, message, cause)
}

// IsEventError checks if an error is an event error with the given code.
func IsEventError(err error, code string) bool {
	return component.IsComponentError(err, code)
}

// ===== EVENT CONSTRUCTORS =====

// NewEventBus creates a new event bus with default configuration.
// This is the primary way to create an EventBus instance for component communication.
func NewEventBus() EventBus {
	return event.NewEventBus()
}

// ===== EVENT UTILITIES =====

// NewEvent creates a new event with the given topic, source, and payload.
func NewEvent(topic, source string, payload map[string]interface{}) *Event {
	return &Event{
		Topic:   topic,
		Source:  source,
		Time:    time.Now(),
		Payload: payload,
	}
}

// Publish publishes an event to the given event bus.
func Publish(bus EventBus, topic string, data interface{}) {
	bus.Publish(topic, data)
}

// Subscribe subscribes to events of a specific topic.
func Subscribe(bus EventBus, topic string, handler EventHandler) Subscription {
	return bus.Subscribe(topic, handler)
}

// SubscribeAsync subscribes to events of a specific topic asynchronously.
func SubscribeAsync(bus EventBus, topic string, handler EventHandler) Subscription {
	return bus.SubscribeAsync(topic, handler)
}
