package event

import (
	"testing"
	"time"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/event"
	infraContext "github.com/fintechain/skeleton/internal/infrastructure/context"
	infraEvent "github.com/fintechain/skeleton/internal/infrastructure/event"
	"github.com/fintechain/skeleton/test/unit/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewEventBus(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "event-bus",
		Name: "Event Bus",
		Type: component.TypeService,
	}

	eventBus := infraEvent.NewEventBus(config)
	assert.NotNil(t, eventBus)

	// Verify interface compliance
	var _ event.EventBusService = eventBus
	var _ component.Service = eventBus
	var _ component.Component = eventBus

	// Test basic properties
	assert.Equal(t, component.ComponentID("event-bus"), eventBus.ID())
	assert.Equal(t, "Event Bus", eventBus.Name())
	assert.Equal(t, component.TypeService, eventBus.Type())
}

func TestEventBusInitialState(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "event-bus",
		Name: "Event Bus",
		Type: component.TypeService,
	}

	eventBus := infraEvent.NewEventBus(config)

	// Test initial state
	assert.False(t, eventBus.IsRunning())
	assert.Equal(t, component.StatusStopped, eventBus.Status())
}

func TestEventBusPublish(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "event-bus",
		Name: "Event Bus",
		Type: component.TypeService,
	}

	eventBus := infraEvent.NewEventBus(config)
	ctx := infraContext.NewContext()

	// Start the event bus
	err := eventBus.Start(ctx)
	assert.NoError(t, err)

	// Test publish (should not error even with no subscribers)
	testEvent := &event.Event{
		Topic:   "test.topic",
		Source:  "test",
		Time:    time.Now(),
		Payload: map[string]interface{}{"message": "test event", "value": 42},
	}

	err = eventBus.Publish(testEvent)
	assert.NoError(t, err)
}

func TestEventBusPublishAsync(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "event-bus",
		Name: "Event Bus",
		Type: component.TypeService,
	}

	eventBus := infraEvent.NewEventBus(config)
	ctx := infraContext.NewContext()

	// Start the event bus
	err := eventBus.Start(ctx)
	assert.NoError(t, err)

	// Test async publish
	asyncEvent := &event.Event{
		Topic:   "test.async.topic",
		Source:  "test",
		Time:    time.Now(),
		Payload: map[string]interface{}{"message": "async test event", "value": 100},
	}

	err = eventBus.PublishAsync(asyncEvent)
	assert.NoError(t, err)

	// Give async operations time to complete
	time.Sleep(10 * time.Millisecond)
}

func TestEventBusSubscribe(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "event-bus",
		Name: "Event Bus",
		Type: component.TypeService,
	}

	eventBus := infraEvent.NewEventBus(config)
	ctx := infraContext.NewContext()

	// Start the event bus
	err := eventBus.Start(ctx)
	assert.NoError(t, err)

	// Test subscribe
	received := false
	handler := func(e *event.Event) {
		received = true
		assert.Equal(t, "test.subscribe.topic", e.Topic)
		assert.NotNil(t, e.Payload)
	}

	subscription := eventBus.Subscribe("test.subscribe.topic", handler)
	assert.NotNil(t, subscription)

	// Publish an event
	publishEvent := &event.Event{
		Topic:   "test.subscribe.topic",
		Source:  "test",
		Time:    time.Now(),
		Payload: map[string]interface{}{"test": "data"},
	}
	err = eventBus.Publish(publishEvent)
	assert.NoError(t, err)

	// Give handler time to execute
	time.Sleep(10 * time.Millisecond)
	assert.True(t, received)
}

func TestEventBusSubscriptionCancel(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "event-bus",
		Name: "Event Bus",
		Type: component.TypeService,
	}

	eventBus := infraEvent.NewEventBus(config)
	ctx := infraContext.NewContext()

	// Start the event bus
	err := eventBus.Start(ctx)
	assert.NoError(t, err)

	// Test subscribe and cancel
	callCount := 0
	handler := func(e *event.Event) {
		callCount++
	}

	subscription := eventBus.Subscribe("test.cancel.topic", handler)

	// Publish first event
	firstEvent := &event.Event{
		Topic:   "test.cancel.topic",
		Source:  "test",
		Time:    time.Now(),
		Payload: map[string]interface{}{"test": "data1"},
	}
	err = eventBus.Publish(firstEvent)
	assert.NoError(t, err)
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 1, callCount)

	// Cancel subscription
	subscription.Cancel()

	// Publish second event (should not be received)
	secondEvent := &event.Event{
		Topic:   "test.cancel.topic",
		Source:  "test",
		Time:    time.Now(),
		Payload: map[string]interface{}{"test": "data2"},
	}
	err = eventBus.Publish(secondEvent)
	assert.NoError(t, err)
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 1, callCount) // Should still be 1
}

func TestEventBusLifecycle(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "event-bus",
		Name: "Event Bus",
		Type: component.TypeService,
	}

	eventBus := infraEvent.NewEventBus(config)
	ctx := infraContext.NewContext()

	// Test start
	err := eventBus.Start(ctx)
	assert.NoError(t, err)
	assert.True(t, eventBus.IsRunning())
	assert.Equal(t, component.StatusRunning, eventBus.Status())

	// Test stop
	err = eventBus.Stop(ctx)
	assert.NoError(t, err)
	assert.False(t, eventBus.IsRunning())
	assert.Equal(t, component.StatusStopped, eventBus.Status())
}

func TestEventBusInitializeAndDispose(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "event-bus",
		Name: "Event Bus",
		Type: component.TypeService,
	}

	eventBus := infraEvent.NewEventBus(config)
	factory := mocks.NewFactory()
	mockSystem := factory.SystemInterface()
	ctx := infraContext.NewContext()

	// Test initialization
	err := eventBus.Initialize(ctx, mockSystem)
	assert.NoError(t, err)

	// Test disposal
	err = eventBus.Dispose()
	assert.NoError(t, err)
}
