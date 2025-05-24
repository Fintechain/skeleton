package mocks

import (
	"github.com/ebanfa/skeleton/internal/infrastructure/event"
)

// MockEventBus is a mock implementation of event.EventBus for testing
type MockEventBus struct {
	// Function fields for customizing behavior
	PublishFunc        func(string, interface{})
	SubscribeFunc      func(string, event.EventHandler) event.Subscription
	SubscribeAsyncFunc func(string, event.EventHandler) event.Subscription
	WaitAsyncFunc      func()

	// Call tracking
	PublishCalls        []PublishCall
	SubscribeCalls      []SubscribeCall
	SubscribeAsyncCalls []SubscribeCall
	WaitAsyncCalls      int

	// State
	Events        []event.Event
	Subscriptions map[string][]event.EventHandler
}

type PublishCall struct {
	Topic string
	Data  interface{}
}

type SubscribeCall struct {
	Topic   string
	Handler event.EventHandler
}

// MockSubscription is a mock implementation of event.Subscription
type MockSubscription struct {
	topic     string
	cancelled bool
}

// NewMockEventBus creates a new mock event bus
func NewMockEventBus() *MockEventBus {
	return &MockEventBus{
		Subscriptions: make(map[string][]event.EventHandler),
	}
}

// Publish implements event.EventBus
func (m *MockEventBus) Publish(topic string, data interface{}) {
	m.PublishCalls = append(m.PublishCalls, PublishCall{Topic: topic, Data: data})
	if m.PublishFunc != nil {
		m.PublishFunc(topic, data)
		return
	}

	// Default behavior: store event and notify handlers
	evt := event.Event{
		Topic:   topic,
		Payload: map[string]interface{}{"data": data},
	}
	m.Events = append(m.Events, evt)

	// Notify handlers
	if handlers, exists := m.Subscriptions[topic]; exists {
		for _, handler := range handlers {
			handler(&evt)
		}
	}
}

// Subscribe implements event.EventBus
func (m *MockEventBus) Subscribe(topic string, handler event.EventHandler) event.Subscription {
	m.SubscribeCalls = append(m.SubscribeCalls, SubscribeCall{Topic: topic, Handler: handler})
	if m.SubscribeFunc != nil {
		return m.SubscribeFunc(topic, handler)
	}

	// Default behavior: store handler
	if m.Subscriptions[topic] == nil {
		m.Subscriptions[topic] = []event.EventHandler{}
	}
	m.Subscriptions[topic] = append(m.Subscriptions[topic], handler)

	return &MockSubscription{topic: topic}
}

// SubscribeAsync implements event.EventBus
func (m *MockEventBus) SubscribeAsync(topic string, handler event.EventHandler) event.Subscription {
	m.SubscribeAsyncCalls = append(m.SubscribeAsyncCalls, SubscribeCall{Topic: topic, Handler: handler})
	if m.SubscribeAsyncFunc != nil {
		return m.SubscribeAsyncFunc(topic, handler)
	}

	// Default behavior: same as Subscribe for testing
	return m.Subscribe(topic, handler)
}

// WaitAsync implements event.EventBus
func (m *MockEventBus) WaitAsync() {
	m.WaitAsyncCalls++
	if m.WaitAsyncFunc != nil {
		m.WaitAsyncFunc()
	}
}

// Cancel implements event.Subscription
func (s *MockSubscription) Cancel() {
	s.cancelled = true
}

// Topic implements event.Subscription
func (s *MockSubscription) Topic() string {
	return s.topic
}

// IsCancelled is a helper method for testing
func (s *MockSubscription) IsCancelled() bool {
	return s.cancelled
}
