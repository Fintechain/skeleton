package mocks

import (
	"sync"
	"time"

	"github.com/ebanfa/skeleton/internal/infrastructure/event"
)

// MockEventBus implements the event.EventBus interface for testing
type MockEventBus struct {
	events      []event.Event
	subscribers map[string][]event.EventHandler
	mu          sync.RWMutex
}

// MockSubscription implements the event.Subscription interface
type MockSubscription struct {
	eventBus *MockEventBus
	topic    string
	index    int
}

// NewMockEventBus creates a new mock event bus
func NewMockEventBus() *MockEventBus {
	return &MockEventBus{
		events:      make([]event.Event, 0),
		subscribers: make(map[string][]event.EventHandler),
	}
}

// Publish publishes an event to the event bus
func (m *MockEventBus) Publish(topic string, data interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Create the event
	newEvent := event.Event{
		Topic:   topic,
		Source:  "test-source",
		Time:    time.Now(),
		Payload: make(map[string]interface{}),
	}

	// Add data to payload if it's a map
	if payloadMap, ok := data.(map[string]interface{}); ok {
		for k, v := range payloadMap {
			newEvent.Payload[k] = v
		}
	} else {
		// Otherwise just add it as a "data" key
		newEvent.Payload["data"] = data
	}

	// Store the event
	m.events = append(m.events, newEvent)

	// Notify subscribers
	if handlers, ok := m.subscribers[topic]; ok {
		for _, handler := range handlers {
			handler(&newEvent)
		}
	}
}

// Subscribe subscribes to events with the given topic
func (m *MockEventBus) Subscribe(topic string, handler event.EventHandler) event.Subscription {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, ok := m.subscribers[topic]; !ok {
		m.subscribers[topic] = make([]event.EventHandler, 0)
	}

	index := len(m.subscribers[topic])
	m.subscribers[topic] = append(m.subscribers[topic], handler)

	return &MockSubscription{
		eventBus: m,
		topic:    topic,
		index:    index,
	}
}

// SubscribeAsync subscribes asynchronously to events
func (m *MockEventBus) SubscribeAsync(topic string, handler event.EventHandler) event.Subscription {
	return m.Subscribe(topic, handler)
}

// WaitAsync waits for all async events to be processed (no-op in this mock)
func (m *MockEventBus) WaitAsync() {
	// No-op in mock implementation
}

// GetEvents returns all events published to the bus
func (m *MockEventBus) GetEvents() []event.Event {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to prevent concurrent modification
	events := make([]event.Event, len(m.events))
	copy(events, m.events)
	return events
}

// GetEventsByTopic returns events for a specific topic
func (m *MockEventBus) GetEventsByTopic(topic string) []event.Event {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var filtered []event.Event
	for _, evt := range m.events {
		if evt.Topic == topic {
			filtered = append(filtered, evt)
		}
	}
	return filtered
}

// ClearEvents clears all events
func (m *MockEventBus) ClearEvents() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.events = make([]event.Event, 0)
}

// Topic returns the topic of the subscription
func (s *MockSubscription) Topic() string {
	return s.topic
}

// Cancel cancels the subscription
func (s *MockSubscription) Cancel() {
	s.eventBus.mu.Lock()
	defer s.eventBus.mu.Unlock()

	if handlers, ok := s.eventBus.subscribers[s.topic]; ok {
		if s.index < len(handlers) {
			// Remove the handler by replacing it with the last one and reducing slice length
			lastIdx := len(handlers) - 1
			handlers[s.index] = handlers[lastIdx]
			s.eventBus.subscribers[s.topic] = handlers[:lastIdx]
		}
	}
}
