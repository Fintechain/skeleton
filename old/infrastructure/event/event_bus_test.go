package event

import (
	"sync"
	"testing"
	"time"
)

func TestEventBusSubscribeAndPublish(t *testing.T) {
	bus := NewEventBus()

	// Track received events
	var receivedEvents []*Event
	var mu sync.Mutex

	// Subscribe to events
	_ = bus.Subscribe("test-topic", func(e *Event) {
		mu.Lock()
		defer mu.Unlock()
		receivedEvents = append(receivedEvents, e)
	})

	// Publish an event
	testData := "test-data"
	bus.Publish("test-topic", testData)

	// Check if the event was received
	mu.Lock()
	defer mu.Unlock()

	if len(receivedEvents) != 1 {
		t.Fatalf("Expected 1 event, got %d", len(receivedEvents))
	}

	event := receivedEvents[0]
	if event.Topic != "test-topic" {
		t.Errorf("Expected topic 'test-topic', got '%s'", event.Topic)
	}

	if event.Payload["data"] != testData {
		t.Errorf("Expected payload data '%v', got '%v'", testData, event.Payload["data"])
	}
}

func TestEventBusSubscriptionCancellation(t *testing.T) {
	bus := NewEventBus()

	counter := 0
	subscription := bus.Subscribe("test-topic", func(e *Event) {
		counter++
	})

	// First publish - should be received
	bus.Publish("test-topic", "data1")

	if counter != 1 {
		t.Fatalf("Expected counter to be 1, got %d", counter)
	}

	// Cancel subscription
	subscription.Cancel()

	// Second publish - should not be received
	bus.Publish("test-topic", "data2")

	if counter != 1 {
		t.Errorf("Expected counter to remain 1 after cancellation, got %d", counter)
	}
}

func TestEventBusAsyncSubscription(t *testing.T) {
	bus := NewEventBus()

	var wg sync.WaitGroup
	wg.Add(1)

	// Use a channel to synchronize the test
	done := make(chan struct{})

	// Subscribe asynchronously
	_ = bus.SubscribeAsync("async-topic", func(e *Event) {
		// Simulate some work
		time.Sleep(50 * time.Millisecond)
		close(done)
		wg.Done()
	})

	// Publish an event
	bus.Publish("async-topic", "async-data")

	// Wait for all async handlers to complete
	select {
	case <-done:
		// Event was processed
	case <-time.After(200 * time.Millisecond):
		t.Fatal("Async handler timed out")
	}

	// Make sure WaitAsync works
	bus.WaitAsync()

	// Check that wg.Wait() doesn't block (all tasks completed)
	wgChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(wgChan)
	}()

	select {
	case <-wgChan:
		// Success - wait group completed
	case <-time.After(100 * time.Millisecond):
		t.Fatal("WaitAsync did not wait for all handlers")
	}
}

func TestEventBusMultipleSubscribers(t *testing.T) {
	bus := NewEventBus()

	counters := make([]int, 3)

	// Add multiple subscribers to the same topic
	for i := 0; i < 3; i++ {
		idx := i // Capture loop variable
		bus.Subscribe("multi-topic", func(e *Event) {
			counters[idx]++
		})
	}

	// Publish an event
	bus.Publish("multi-topic", "multi-data")

	// All subscribers should have received the event
	for i, count := range counters {
		if count != 1 {
			t.Errorf("Subscriber %d: expected count 1, got %d", i, count)
		}
	}
}

func TestEventBusMultipleTopics(t *testing.T) {
	bus := NewEventBus()

	topic1Count := 0
	topic2Count := 0

	bus.Subscribe("topic1", func(e *Event) {
		topic1Count++
	})

	bus.Subscribe("topic2", func(e *Event) {
		topic2Count++
	})

	// Publish to topic1
	bus.Publish("topic1", "data1")

	if topic1Count != 1 {
		t.Errorf("Topic1: expected count 1, got %d", topic1Count)
	}

	if topic2Count != 0 {
		t.Errorf("Topic2: expected count 0, got %d", topic2Count)
	}

	// Publish to topic2
	bus.Publish("topic2", "data2")

	if topic1Count != 1 {
		t.Errorf("Topic1: expected count to remain 1, got %d", topic1Count)
	}

	if topic2Count != 1 {
		t.Errorf("Topic2: expected count 1, got %d", topic2Count)
	}
}
