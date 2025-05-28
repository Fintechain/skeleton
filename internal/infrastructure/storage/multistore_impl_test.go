package storage

import (
	"os"
	"testing"

	"github.com/fintechain/skeleton/internal/domain/component/mocks"
	"github.com/fintechain/skeleton/internal/domain/storage"
	"github.com/fintechain/skeleton/internal/infrastructure/logging"
)

func TestNewMultiStore(t *testing.T) {
	logger := logging.CreateStandardLogger(logging.Info)
	eventBus := mocks.NewMockEventBus()

	// Test with nil config
	ms := NewMultiStore(nil, logger, eventBus)
	if ms == nil {
		t.Fatal("NewMultiStore returned nil")
	}

	if ms.GetDefaultEngine() != "memory" {
		t.Errorf("Expected default engine 'memory', got '%s'", ms.GetDefaultEngine())
	}

	// Test with custom config
	config := &storage.MultiStoreConfig{
		RootPath:      "./test-data",
		DefaultEngine: "memory",
		EngineConfigs: make(map[string]storage.Config),
	}

	ms2 := NewMultiStore(config, logger, eventBus)
	if ms2 == nil {
		t.Fatal("NewMultiStore with config returned nil")
	}

	if ms2.GetDefaultEngine() != "memory" {
		t.Errorf("Expected default engine 'memory', got '%s'", ms2.GetDefaultEngine())
	}

	// Test with nil event bus (should work)
	ms3 := NewMultiStore(config, logger, nil)
	if ms3 == nil {
		t.Fatal("NewMultiStore with nil event bus returned nil")
	}
}

func TestNewMultiStore_PanicsWithNilLogger(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with nil logger")
		}
	}()

	NewMultiStore(nil, nil, nil)
}

func TestMultiStore_EngineManagement(t *testing.T) {
	logger := logging.CreateStandardLogger(logging.Info)
	eventBus := mocks.NewMockEventBus()
	ms := NewMultiStore(nil, logger, eventBus)

	// Test listing engines
	engines := ms.ListEngines()
	if len(engines) == 0 {
		t.Error("Expected at least one engine to be registered")
	}

	// Memory engine should be registered by default
	found := false
	for _, engine := range engines {
		if engine == "memory" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Memory engine should be registered by default")
	}

	// Test getting engine
	engine, err := ms.GetEngine("memory")
	if err != nil {
		t.Errorf("Failed to get memory engine: %v", err)
	}
	if engine == nil {
		t.Error("Got nil engine")
	}
	if engine.Name() != "memory" {
		t.Errorf("Expected engine name 'memory', got '%s'", engine.Name())
	}

	// Test getting non-existent engine
	_, err = ms.GetEngine("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent engine")
	}
	if !storage.IsEngineNotFound(err) {
		t.Errorf("Expected ErrEngineNotFound, got %v", err)
	}

	// Test getting engine with empty name
	_, err = ms.GetEngine("")
	if err == nil {
		t.Error("Expected error for empty engine name")
	}
	if !storage.IsInvalidConfig(err) {
		t.Errorf("Expected ErrInvalidConfig, got %v", err)
	}
}

func TestMultiStore_StoreManagement(t *testing.T) {
	logger := logging.CreateStandardLogger(logging.Info)
	eventBus := mocks.NewMockEventBus()

	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "multistore_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := &storage.MultiStoreConfig{
		RootPath:      tempDir,
		DefaultEngine: "memory",
		EngineConfigs: make(map[string]storage.Config),
	}

	ms := NewMultiStore(config, logger, eventBus)

	// Test creating store
	err = ms.CreateStore("test-store", "", nil)
	if err != nil {
		t.Errorf("Failed to create store: %v", err)
	}

	// Verify store created event was published
	events := eventBus.GetEventsByTopic(storage.TopicStoreCreated)
	if len(events) != 1 {
		t.Errorf("Expected 1 store created event, got %d", len(events))
	}

	// Test store exists
	if !ms.StoreExists("test-store") {
		t.Error("Store should exist after creation")
	}

	// Test listing stores
	stores := ms.ListStores()
	if len(stores) != 1 {
		t.Errorf("Expected 1 store, got %d", len(stores))
	}
	if stores[0] != "test-store" {
		t.Errorf("Expected store name 'test-store', got '%s'", stores[0])
	}

	// Test getting store
	store, err := ms.GetStore("test-store")
	if err != nil {
		t.Errorf("Failed to get store: %v", err)
	}
	if store == nil {
		t.Error("Got nil store")
	}
	if store.Name() != "test-store" {
		t.Errorf("Expected store name 'test-store', got '%s'", store.Name())
	}

	// Test creating duplicate store
	err = ms.CreateStore("test-store", "", nil)
	if err == nil {
		t.Error("Expected error when creating duplicate store")
	}
	if !storage.IsStoreExists(err) {
		t.Errorf("Expected ErrStoreExists, got %v", err)
	}

	// Test deleting store
	err = ms.DeleteStore("test-store")
	if err != nil {
		t.Errorf("Failed to delete store: %v", err)
	}

	// Verify store deleted event was published
	events = eventBus.GetEventsByTopic(storage.TopicStoreDeleted)
	if len(events) != 1 {
		t.Errorf("Expected 1 store deleted event, got %d", len(events))
	}

	// Test store no longer exists
	if ms.StoreExists("test-store") {
		t.Error("Store should not exist after deletion")
	}

	// Test getting deleted store
	_, err = ms.GetStore("test-store")
	if err == nil {
		t.Error("Expected error when getting deleted store")
	}
	if !storage.IsStoreNotFound(err) {
		t.Errorf("Expected ErrStoreNotFound, got %v", err)
	}
}

func TestMultiStore_EventPublishing(t *testing.T) {
	logger := logging.CreateStandardLogger(logging.Info)
	eventBus := mocks.NewMockEventBus()

	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "multistore_event_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := &storage.MultiStoreConfig{
		RootPath:      tempDir,
		DefaultEngine: "memory",
		EngineConfigs: make(map[string]storage.Config),
	}

	ms := NewMultiStore(config, logger, eventBus)

	// Create multiple stores
	err = ms.CreateStore("store1", "", nil)
	if err != nil {
		t.Errorf("Failed to create store1: %v", err)
	}

	err = ms.CreateStore("store2", "", nil)
	if err != nil {
		t.Errorf("Failed to create store2: %v", err)
	}

	// Verify store created events
	createdEvents := eventBus.GetEventsByTopic(storage.TopicStoreCreated)
	if len(createdEvents) != 2 {
		t.Errorf("Expected 2 store created events, got %d", len(createdEvents))
	}

	// Test CloseAll publishes events
	err = ms.CloseAll()
	if err != nil {
		t.Errorf("Failed to close all stores: %v", err)
	}

	// Verify store closed events
	closedEvents := eventBus.GetEventsByTopic(storage.TopicStoreClosed)
	if len(closedEvents) != 2 {
		t.Errorf("Expected 2 store closed events, got %d", len(closedEvents))
	}
}

func TestMultiStore_WithoutEventBus(t *testing.T) {
	logger := logging.CreateStandardLogger(logging.Info)

	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "multistore_no_event_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := &storage.MultiStoreConfig{
		RootPath:      tempDir,
		DefaultEngine: "memory",
		EngineConfigs: make(map[string]storage.Config),
	}

	// Test with nil event bus - should work without errors
	ms := NewMultiStore(config, logger, nil)

	// Test creating store without event bus
	err = ms.CreateStore("test-store", "", nil)
	if err != nil {
		t.Errorf("Failed to create store without event bus: %v", err)
	}

	// Test deleting store without event bus
	err = ms.DeleteStore("test-store")
	if err != nil {
		t.Errorf("Failed to delete store without event bus: %v", err)
	}
}

func TestMultiStore_ErrorCases(t *testing.T) {
	logger := logging.CreateStandardLogger(logging.Info)
	eventBus := mocks.NewMockEventBus()
	ms := NewMultiStore(nil, logger, eventBus)

	// Test creating store with empty name
	err := ms.CreateStore("", "", nil)
	if err == nil {
		t.Error("Expected error for empty store name")
	}
	if !storage.IsInvalidConfig(err) {
		t.Errorf("Expected ErrInvalidConfig, got %v", err)
	}

	// Test getting store with empty name
	_, err = ms.GetStore("")
	if err == nil {
		t.Error("Expected error for empty store name")
	}
	if !storage.IsInvalidConfig(err) {
		t.Errorf("Expected ErrInvalidConfig, got %v", err)
	}

	// Test deleting store with empty name
	err = ms.DeleteStore("")
	if err == nil {
		t.Error("Expected error for empty store name")
	}
	if !storage.IsInvalidConfig(err) {
		t.Errorf("Expected ErrInvalidConfig, got %v", err)
	}

	// Test getting non-existent store
	_, err = ms.GetStore("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent store")
	}
	if !storage.IsStoreNotFound(err) {
		t.Errorf("Expected ErrStoreNotFound, got %v", err)
	}

	// Test deleting non-existent store
	err = ms.DeleteStore("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent store")
	}
	if !storage.IsStoreNotFound(err) {
		t.Errorf("Expected ErrStoreNotFound, got %v", err)
	}

	// Test creating store with non-existent engine
	err = ms.CreateStore("test-store", "nonexistent", nil)
	if err == nil {
		t.Error("Expected error for non-existent engine")
	}
	if !storage.IsEngineNotFound(err) {
		t.Errorf("Expected ErrEngineNotFound, got %v", err)
	}
}

func TestMultiStore_DefaultEngine(t *testing.T) {
	logger := logging.CreateStandardLogger(logging.Info)
	ms := NewMultiStore(nil, logger, nil)

	// Test getting default engine
	defaultEngine := ms.GetDefaultEngine()
	if defaultEngine != "memory" {
		t.Errorf("Expected default engine 'memory', got '%s'", defaultEngine)
	}

	// Test setting default engine
	ms.SetDefaultEngine("custom")
	if ms.GetDefaultEngine() != "custom" {
		t.Errorf("Expected default engine 'custom', got '%s'", ms.GetDefaultEngine())
	}
}

func TestMultiStore_CloseAll(t *testing.T) {
	logger := logging.CreateStandardLogger(logging.Info)

	// Create temporary directory for test
	tempDir, err := os.MkdirTemp("", "multistore_closeall_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	config := &storage.MultiStoreConfig{
		RootPath:      tempDir,
		DefaultEngine: "memory",
		EngineConfigs: make(map[string]storage.Config),
	}

	ms := NewMultiStore(config, logger, nil)

	// Create multiple stores
	err = ms.CreateStore("store1", "", nil)
	if err != nil {
		t.Errorf("Failed to create store1: %v", err)
	}

	err = ms.CreateStore("store2", "", nil)
	if err != nil {
		t.Errorf("Failed to create store2: %v", err)
	}

	// Verify stores exist
	if len(ms.ListStores()) != 2 {
		t.Errorf("Expected 2 stores, got %d", len(ms.ListStores()))
	}

	// Close all stores
	err = ms.CloseAll()
	if err != nil {
		t.Errorf("Failed to close all stores: %v", err)
	}

	// Verify stores are cleared
	if len(ms.ListStores()) != 0 {
		t.Errorf("Expected 0 stores after CloseAll, got %d", len(ms.ListStores()))
	}
}
