// Package event provides concrete implementations of the event system.
package event

import (
	"sync"
	"time"

	"github.com/fintechain/skeleton/internal/domain/event"
)

// DefaultEventBus is a concrete implementation of the EventBus interface.
type DefaultEventBus struct {
	mu            sync.RWMutex
	subscriptions map[string][]*DefaultSubscription
	asyncWG       sync.WaitGroup
}

// DefaultSubscription is a concrete implementation of the Subscription interface.
type DefaultSubscription struct {
	topic     string
	handler   event.EventHandler
	cancelled bool
	mu        sync.RWMutex
	eventBus  *DefaultEventBus
}

// NewEventBus creates a new EventBus instance with minimal dependencies.
func NewEventBus() event.EventBus {
	return &DefaultEventBus{
		subscriptions: make(map[string][]*DefaultSubscription),
	}
}

// Publish publishes an event to all subscribers of the given topic.
func (eb *DefaultEventBus) Publish(topic string, data interface{}) {
	eb.mu.RLock()
	subs := eb.subscriptions[topic]
	eb.mu.RUnlock()

	if len(subs) == 0 {
		return
	}

	// Create the event
	evt := &event.Event{
		Topic:   topic,
		Source:  "system", // Could be enhanced to track actual source
		Time:    time.Now(),
		Payload: make(map[string]interface{}),
	}

	// Convert data to payload
	if data != nil {
		if payload, ok := data.(map[string]interface{}); ok {
			evt.Payload = payload
		} else {
			evt.Payload["data"] = data
		}
	}

	// Notify all subscribers
	for _, sub := range subs {
		sub.mu.RLock()
		if !sub.cancelled {
			// Call handler synchronously
			sub.handler(evt)
		}
		sub.mu.RUnlock()
	}
}

// Subscribe creates a synchronous subscription to the given topic.
func (eb *DefaultEventBus) Subscribe(topic string, handler event.EventHandler) event.Subscription {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	sub := &DefaultSubscription{
		topic:    topic,
		handler:  handler,
		eventBus: eb,
	}

	eb.subscriptions[topic] = append(eb.subscriptions[topic], sub)
	return sub
}

// SubscribeAsync creates an asynchronous subscription to the given topic.
func (eb *DefaultEventBus) SubscribeAsync(topic string, handler event.EventHandler) event.Subscription {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	// Wrap the handler to run asynchronously
	asyncHandler := func(evt *event.Event) {
		eb.asyncWG.Add(1)
		go func() {
			defer eb.asyncWG.Done()
			handler(evt)
		}()
	}

	sub := &DefaultSubscription{
		topic:    topic,
		handler:  asyncHandler,
		eventBus: eb,
	}

	eb.subscriptions[topic] = append(eb.subscriptions[topic], sub)
	return sub
}

// WaitAsync waits for all asynchronous event handlers to complete.
func (eb *DefaultEventBus) WaitAsync() {
	eb.asyncWG.Wait()
}

// Cancel cancels the subscription.
func (s *DefaultSubscription) Cancel() {
	s.mu.Lock()
	s.cancelled = true
	s.mu.Unlock()

	// Remove from event bus
	s.eventBus.mu.Lock()
	defer s.eventBus.mu.Unlock()

	subs := s.eventBus.subscriptions[s.topic]
	for i, sub := range subs {
		if sub == s {
			// Remove this subscription
			s.eventBus.subscriptions[s.topic] = append(subs[:i], subs[i+1:]...)
			break
		}
	}

	// Clean up empty topic entries
	if len(s.eventBus.subscriptions[s.topic]) == 0 {
		delete(s.eventBus.subscriptions, s.topic)
	}
}

// Topic returns the topic of the subscription.
func (s *DefaultSubscription) Topic() string {
	return s.topic
}
