// Package event provides interfaces and types for the event system.
package event

// Standard event error codes
const (
	// ErrEventNotFound is returned when an event doesn't exist
	ErrEventNotFound = "event.event_not_found"

	// ErrInvalidEventType is returned when an invalid event type is provided
	ErrInvalidEventType = "event.invalid_event_type"

	// ErrInvalidEventData is returned when invalid event data is provided
	ErrInvalidEventData = "event.invalid_event_data"

	// ErrEventBusNotStarted is returned when operations are performed on a non-started event bus
	ErrEventBusNotStarted = "event.event_bus_not_started"

	// ErrEventBusAlreadyStarted is returned when starting an already started event bus
	ErrEventBusAlreadyStarted = "event.event_bus_already_started"

	// ErrSubscriberNotFound is returned when a subscriber doesn't exist
	ErrSubscriberNotFound = "event.subscriber_not_found"

	// ErrSubscriberExists is returned when creating a subscriber that already exists
	ErrSubscriberExists = "event.subscriber_exists"

	// ErrPublishFailed is returned when event publishing fails
	ErrPublishFailed = "event.publish_failed"

	// ErrSubscriptionFailed is returned when event subscription fails
	ErrSubscriptionFailed = "event.subscription_failed"

	// ErrInvalidEventConfig is returned when invalid event configuration is provided
	ErrInvalidEventConfig = "event.invalid_event_config"
)
