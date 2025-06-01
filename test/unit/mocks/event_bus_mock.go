package mocks

import (
	"fmt"
	"sync"
	"time"

	"github.com/fintechain/skeleton/pkg/event"
)

// MockEventBus provides a configurable mock implementation of the event.EventBus interface.
type MockEventBus struct {
	mu sync.RWMutex

	// Event storage
	publishedEvents []event.Event
	subscriptions   map[string][]event.EventHandler

	// Behavior configuration
	shouldFail   bool
	failureError string

	// Call tracking
	callCount map[string]int
	lastCalls map[string][]interface{}

	// Async control
	asyncWaiting bool
}

// NewMockEventBus creates a new configurable event bus mock.
func NewMockEventBus() *MockEventBus {
	return &MockEventBus{
		publishedEvents: make([]event.Event, 0),
		subscriptions:   make(map[string][]event.EventHandler),
		callCount:       make(map[string]int),
		lastCalls:       make(map[string][]interface{}),
	}
}

// EventBus Interface Implementation

// Publish publishes an event to all subscribers.
func (m *MockEventBus) Publish(topic string, data interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("Publish", topic, data)

	// Create event
	evt := event.Event{
		Topic:   topic,
		Source:  "mock",
		Time:    time.Now(),
		Payload: map[string]interface{}{"data": data},
	}

	// Store the published event
	m.publishedEvents = append(m.publishedEvents, evt)

	// Notify subscribers (in mock, we just track this)
	if handlers, exists := m.subscriptions[topic]; exists {
		for _, handler := range handlers {
			// In a real implementation, we would call handler(&evt)
			// For mock purposes, we just track that it would be called
			_ = handler
		}
	}
}

// Subscribe subscribes a handler to events of a specific topic.
func (m *MockEventBus) Subscribe(topic string, handler event.EventHandler) event.Subscription {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("Subscribe", topic, handler)

	if m.subscriptions[topic] == nil {
		m.subscriptions[topic] = make([]event.EventHandler, 0)
	}
	m.subscriptions[topic] = append(m.subscriptions[topic], handler)

	return &MockSubscription{
		topic:    topic,
		eventBus: m,
	}
}

// SubscribeAsync subscribes a handler to events of a specific topic asynchronously.
func (m *MockEventBus) SubscribeAsync(topic string, handler event.EventHandler) event.Subscription {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("SubscribeAsync", topic, handler)

	if m.subscriptions[topic] == nil {
		m.subscriptions[topic] = make([]event.EventHandler, 0)
	}
	m.subscriptions[topic] = append(m.subscriptions[topic], handler)

	return &MockSubscription{
		topic:    topic,
		eventBus: m,
	}
}

// WaitAsync waits for all async operations to complete.
func (m *MockEventBus) WaitAsync() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.trackCall("WaitAsync")
	m.asyncWaiting = true
}

// Mock Configuration Methods

// GetPublishedEvents returns all published events.
func (m *MockEventBus) GetPublishedEvents() []event.Event {
	m.mu.RLock()
	defer m.mu.RUnlock()

	events := make([]event.Event, len(m.publishedEvents))
	copy(events, m.publishedEvents)
	return events
}

// GetSubscriberCount returns the number of subscribers for a topic.
func (m *MockEventBus) GetSubscriberCount(topic string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if handlers, exists := m.subscriptions[topic]; exists {
		return len(handlers)
	}
	return 0
}

// RemoveSubscription removes a subscription (used by MockSubscription).
func (m *MockEventBus) RemoveSubscription(topic string, handler event.EventHandler) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if handlers, exists := m.subscriptions[topic]; exists {
		// Remove the handler (simplified for mock)
		for i, h := range handlers {
			if fmt.Sprintf("%p", h) == fmt.Sprintf("%p", handler) {
				m.subscriptions[topic] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}

// SetShouldFail configures the mock to fail operations.
func (m *MockEventBus) SetShouldFail(fail bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.shouldFail = fail
}

// SetFailureError sets the error message for failed operations.
func (m *MockEventBus) SetFailureError(err string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.failureError = err
}

// State Verification Methods

// GetCallCount returns the number of times a method was called.
func (m *MockEventBus) GetCallCount(method string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.callCount[method]
}

// GetLastCall returns the parameters of the last call to a method.
func (m *MockEventBus) GetLastCall(method string) []interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastCalls[method]
}

// Reset clears all mock state and configuration.
func (m *MockEventBus) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.publishedEvents = make([]event.Event, 0)
	m.subscriptions = make(map[string][]event.EventHandler)
	m.shouldFail = false
	m.failureError = ""
	m.callCount = make(map[string]int)
	m.lastCalls = make(map[string][]interface{})
	m.asyncWaiting = false
}

// Helper Methods

// trackCall records a method call for verification.
func (m *MockEventBus) trackCall(method string, args ...interface{}) {
	m.callCount[method]++
	m.lastCalls[method] = args
}

// getFailureError returns the configured failure error or a default.
func (m *MockEventBus) getFailureError(method string) string {
	if m.failureError != "" {
		return m.failureError
	}
	return fmt.Sprintf("mock_event_bus.%s_failed", method)
}

// MockSubscription provides a mock implementation of event.Subscription.
type MockSubscription struct {
	topic    string
	eventBus *MockEventBus
	handler  event.EventHandler
}

// Cancel cancels the subscription.
func (s *MockSubscription) Cancel() {
	s.eventBus.RemoveSubscription(s.topic, s.handler)
}

// Topic returns the topic of the subscription.
func (s *MockSubscription) Topic() string {
	return s.topic
}

// EventBusMockBuilder provides a fluent interface for configuring event bus mocks.
type EventBusMockBuilder struct {
	mock *MockEventBus
}

// NewEventBusMockBuilder creates a new event bus mock builder.
func NewEventBusMockBuilder() *EventBusMockBuilder {
	return &EventBusMockBuilder{
		mock: NewMockEventBus(),
	}
}

// WithFailure configures the mock to fail operations.
func (b *EventBusMockBuilder) WithFailure(fail bool) *EventBusMockBuilder {
	b.mock.SetShouldFail(fail)
	return b
}

// WithFailureError sets the error message for failed operations.
func (b *EventBusMockBuilder) WithFailureError(err string) *EventBusMockBuilder {
	b.mock.SetFailureError(err)
	return b
}

// Build returns the configured mock event bus.
func (b *EventBusMockBuilder) Build() event.EventBus {
	return b.mock
}
