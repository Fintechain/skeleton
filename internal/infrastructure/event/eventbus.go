package event

import (
	"sync"
	"sync/atomic"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/event"
	infraComponent "github.com/fintechain/skeleton/internal/infrastructure/component"
)

// subscription represents a single event subscription
type subscription struct {
	handler   event.EventHandler
	topic     string
	cancelled atomic.Bool
}

// Cancel cancels the subscription
func (s *subscription) Cancel() {
	s.cancelled.Store(true)
}

// Topic returns the topic this subscription is registered for
func (s *subscription) Topic() string {
	return s.topic
}

// EventBus implements the EventBusService interface
type EventBus struct {
	*infraComponent.BaseService
	subscribers map[string][]*subscription
	mu          sync.RWMutex
	wg          sync.WaitGroup
}

// NewEventBus creates a new event bus
func NewEventBus(config component.ComponentConfig) *EventBus {
	return &EventBus{
		BaseService: infraComponent.NewBaseService(config),
		subscribers: make(map[string][]*subscription),
	}
}

// Publish publishes an event synchronously to all subscribers
func (eb *EventBus) Publish(evt *event.Event) error {
	eb.mu.RLock()
	subs := eb.subscribers[evt.Topic]
	eb.mu.RUnlock()

	for _, sub := range subs {
		if sub.cancelled.Load() {
			continue
		}

		func() {
			defer func() {
				if r := recover(); r != nil {
					// Log panic but continue to next handler
				}
			}()
			sub.handler(evt)
		}()
	}

	return nil
}

// PublishAsync publishes an event asynchronously to all subscribers
func (eb *EventBus) PublishAsync(evt *event.Event) error {
	eb.mu.RLock()
	subs := eb.subscribers[evt.Topic]
	eb.mu.RUnlock()

	for _, sub := range subs {
		if sub.cancelled.Load() {
			continue
		}

		eb.wg.Add(1)
		go func(s *subscription) {
			defer eb.wg.Done()
			defer func() {
				if r := recover(); r != nil {
					// Log panic but continue
				}
			}()
			s.handler(evt)
		}(sub)
	}

	return nil
}

// Subscribe subscribes to events of a specific type
func (eb *EventBus) Subscribe(eventType string, handler event.EventHandler) event.Subscription {
	eb.mu.Lock()
	defer eb.mu.Unlock()

	sub := &subscription{
		handler: handler,
		topic:   eventType,
	}

	eb.subscribers[eventType] = append(eb.subscribers[eventType], sub)
	return sub
}

// SubscribeAsync subscribes to events of a specific type (same as Subscribe)
func (eb *EventBus) SubscribeAsync(eventType string, handler event.EventHandler) event.Subscription {
	return eb.Subscribe(eventType, handler)
}

// WaitAsync waits for all async operations to complete
func (eb *EventBus) WaitAsync() {
	eb.wg.Wait()
}
