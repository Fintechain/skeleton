package event

import (
	"sync"
	"time"
)

// subscription implements the Subscription interface.
type subscription struct {
	topic    string
	handler  EventHandler
	bus      *DefaultEventBus
	isClosed bool
	mu       sync.Mutex // Protects isClosed
}

// Topic returns the topic of the subscription.
func (s *subscription) Topic() string {
	return s.topic
}

// Cancel cancels the subscription.
func (s *subscription) Cancel() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isClosed {
		return
	}

	s.isClosed = true
	s.bus.unsubscribe(s)
}

// DefaultEventBus provides a standard implementation of the EventBus interface.
type DefaultEventBus struct {
	subscribers map[string][]*subscription
	mu          sync.RWMutex   // Protects subscribers
	wg          sync.WaitGroup // Tracks async event handlers
	// Any future dependencies like logger, metrics would go here
}

// DefaultEventBusConfig holds configuration options for DefaultEventBus
type DefaultEventBusConfig struct {
	// Add configuration options here as needed
}

// NewEventBus creates a new event bus with default configuration.
func NewEventBus() *DefaultEventBus {
	return NewEventBusWithConfig(DefaultEventBusConfig{})
}

// NewEventBusWithConfig creates a new event bus with the specified configuration.
// This constructor provides a place to inject dependencies in the future.
func NewEventBusWithConfig(config DefaultEventBusConfig) *DefaultEventBus {
	return &DefaultEventBus{
		subscribers: make(map[string][]*subscription),
		// Initialize any injected dependencies here when added
	}
}

// CreateEventBus is a factory method for backward compatibility
func CreateEventBus() EventBus {
	return NewEventBus()
}

// Subscribe registers a handler for a specific topic.
func (b *DefaultEventBus) Subscribe(topic string, handler EventHandler) Subscription {
	b.mu.Lock()
	defer b.mu.Unlock()

	sub := &subscription{
		topic:   topic,
		handler: handler,
		bus:     b,
	}

	b.subscribers[topic] = append(b.subscribers[topic], sub)
	return sub
}

// SubscribeAsync registers an async handler for a specific topic.
func (b *DefaultEventBus) SubscribeAsync(topic string, handler EventHandler) Subscription {
	// Wrap the handler to make it async
	asyncHandler := func(event *Event) {
		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			handler(event)
		}()
	}

	return b.Subscribe(topic, asyncHandler)
}

// Publish publishes an event to all subscribers of the topic.
func (b *DefaultEventBus) Publish(topic string, data interface{}) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// Create the event
	event := &Event{
		Topic:   topic,
		Time:    time.Now(),
		Payload: make(map[string]interface{}),
	}

	// Set the payload data
	if data != nil {
		// For simplicity, we just add the data directly
		// In practice, we might want to structure this better
		event.Payload["data"] = data
	}

	// Notify all subscribers
	subs := b.subscribers[topic]
	for _, sub := range subs {
		sub.handler(event)
	}
}

// WaitAsync waits for all async handlers to complete.
func (b *DefaultEventBus) WaitAsync() {
	b.wg.Wait()
}

// unsubscribe removes a subscription.
func (b *DefaultEventBus) unsubscribe(sub *subscription) {
	b.mu.Lock()
	defer b.mu.Unlock()

	topic := sub.topic
	subs := b.subscribers[topic]

	// Find and remove the subscription
	for i, s := range subs {
		if s == sub {
			// Remove the subscription by replacing it with the last one
			// and reducing the slice length by 1 (more efficient than append)
			lastIdx := len(subs) - 1
			subs[i] = subs[lastIdx]
			b.subscribers[topic] = subs[:lastIdx]

			// If there are no more subscribers for this topic, remove the topic
			if len(b.subscribers[topic]) == 0 {
				delete(b.subscribers, topic)
			}

			break
		}
	}
}
