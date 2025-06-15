package providers

import (
	"fmt"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/context"
	"github.com/fintechain/skeleton/internal/domain/event"
)

// CustomEventBus demonstrates a custom event bus implementation
type CustomEventBus struct {
	subscribers map[string][]func(interface{})
	status      component.ServiceStatus
}

// NewCustomEventBus creates a new custom event bus instance
func NewCustomEventBus() event.EventBusService {
	return &CustomEventBus{
		subscribers: make(map[string][]func(interface{})),
		status:      component.StatusStopped,
	}
}

func (e *CustomEventBus) ID() component.ComponentID {
	return "custom-event-bus"
}

func (e *CustomEventBus) Name() string {
	return "Custom Event Bus"
}

func (e *CustomEventBus) Type() component.ComponentType {
	return component.TypeService
}

func (e *CustomEventBus) Description() string {
	return "Custom event bus with enhanced features"
}

func (e *CustomEventBus) Version() string {
	return "1.0.0"
}

func (e *CustomEventBus) Metadata() component.Metadata {
	return component.Metadata{
		"provider": "custom",
		"features": []string{"async", "persistent"},
	}
}

func (e *CustomEventBus) Initialize(ctx context.Context, system component.System) error {
	fmt.Println("[CUSTOM] Custom event bus initialized")
	return nil
}

func (e *CustomEventBus) Dispose() error {
	fmt.Println("[CUSTOM] Custom event bus disposed")
	return nil
}

func (e *CustomEventBus) Start(ctx context.Context) error {
	fmt.Println("[CUSTOM] Custom event bus started")
	e.status = component.StatusRunning
	return nil
}

func (e *CustomEventBus) Stop(ctx context.Context) error {
	fmt.Println("[CUSTOM] Custom event bus stopped")
	e.status = component.StatusStopped
	return nil
}

func (e *CustomEventBus) IsRunning() bool {
	return e.status == component.StatusRunning
}

func (e *CustomEventBus) Status() component.ServiceStatus {
	return e.status
}

func (e *CustomEventBus) Publish(event *event.Event) error {
	fmt.Printf("[CUSTOM] Publishing event to topic '%s': %+v\n", event.Topic, event.Payload)

	if subscribers, ok := e.subscribers[event.Topic]; ok {
		for _, subscriber := range subscribers {
			go subscriber(event) // Async delivery
		}
	}

	return nil
}

func (e *CustomEventBus) PublishAsync(event *event.Event) error {
	go e.Publish(event)
	return nil
}

func (e *CustomEventBus) Subscribe(eventType string, handler event.EventHandler) event.Subscription {
	fmt.Printf("[CUSTOM] New subscription to topic '%s'\n", eventType)
	e.subscribers[eventType] = append(e.subscribers[eventType], func(data interface{}) {
		if evt, ok := data.(*event.Event); ok {
			handler(evt)
		}
	})
	return &customSubscription{topic: eventType}
}

func (e *CustomEventBus) SubscribeAsync(eventType string, handler event.EventHandler) event.Subscription {
	return e.Subscribe(eventType, func(evt *event.Event) {
		go handler(evt)
	})
}

func (e *CustomEventBus) WaitAsync() {
	// Simple implementation - in real world you'd wait for async operations
}

// customSubscription implements the Subscription interface
type customSubscription struct {
	topic string
}

func (s *customSubscription) Cancel() {
	fmt.Printf("[CUSTOM] Cancelling subscription to topic '%s'\n", s.topic)
}

func (s *customSubscription) Topic() string {
	return s.topic
}
