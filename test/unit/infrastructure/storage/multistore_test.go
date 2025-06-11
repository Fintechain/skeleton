package storage

import (
	"testing"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/storage"
	infraContext "github.com/fintechain/skeleton/internal/infrastructure/context"
	infraStorage "github.com/fintechain/skeleton/internal/infrastructure/storage"
	memoryStorage "github.com/fintechain/skeleton/internal/infrastructure/storage/memory"
	"github.com/fintechain/skeleton/test/unit/mocks"
	"github.com/stretchr/testify/assert"
)

func TestNewMultiStore(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "multi-store",
		Name: "Multi Store",
		Type: component.TypeService,
	}

	multiStore := infraStorage.NewMultiStore(config, "/tmp/multistore")
	assert.NotNil(t, multiStore)

	// Verify interface compliance
	var _ storage.MultiStore = multiStore
	var _ component.Service = multiStore
	var _ component.Component = multiStore

	// Test basic properties
	assert.Equal(t, component.ComponentID("multi-store"), multiStore.ID())
	assert.Equal(t, "Multi Store", multiStore.Name())
	assert.Equal(t, component.TypeService, multiStore.Type())
}

func TestMultiStoreInitialState(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "multi-store",
		Name: "Multi Store",
		Type: component.TypeService,
	}

	multiStore := infraStorage.NewMultiStore(config, "/tmp/multistore")

	// Test initial state
	assert.False(t, multiStore.IsRunning())
	assert.Equal(t, component.StatusStopped, multiStore.Status())

	// Test no stores initially
	stores := multiStore.ListStores()
	assert.Empty(t, stores)
}

func TestMultiStoreRegisterEngine(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "multi-store",
		Name: "Multi Store",
		Type: component.TypeService,
	}

	multiStore := infraStorage.NewMultiStore(config, "/tmp/multistore")
	engine := memoryStorage.NewEngine()

	// Test register engine
	err := multiStore.RegisterEngine(engine)
	assert.NoError(t, err)

	// Test register duplicate engine (should return error)
	err = multiStore.RegisterEngine(engine)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "storage.engine_exists")
}

func TestMultiStoreRegisterNilEngine(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "multi-store",
		Name: "Multi Store",
		Type: component.TypeService,
	}

	multiStore := infraStorage.NewMultiStore(config, "/tmp/multistore")

	// Test register nil engine
	err := multiStore.RegisterEngine(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "storage.invalid_config")
}

func TestMultiStoreCreateStore(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "multi-store",
		Name: "Multi Store",
		Type: component.TypeService,
	}

	multiStore := infraStorage.NewMultiStore(config, "/tmp/multistore")
	engine := memoryStorage.NewEngine()

	// Register engine
	err := multiStore.RegisterEngine(engine)
	assert.NoError(t, err)

	// Test create store
	err = multiStore.CreateStore("test-store", "memory", nil)
	assert.NoError(t, err)

	// Test store is listed
	stores := multiStore.ListStores()
	assert.Len(t, stores, 1)
	assert.Contains(t, stores, "test-store")
}

func TestMultiStoreCreateStoreWithoutEngine(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "multi-store",
		Name: "Multi Store",
		Type: component.TypeService,
	}

	multiStore := infraStorage.NewMultiStore(config, "/tmp/multistore")

	// Test create store without registering engine
	err := multiStore.CreateStore("test-store", "non-existent-engine", nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "storage.engine_not_found")
}

func TestMultiStoreGetStore(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "multi-store",
		Name: "Multi Store",
		Type: component.TypeService,
	}

	multiStore := infraStorage.NewMultiStore(config, "/tmp/multistore")
	engine := memoryStorage.NewEngine()

	// Register engine and create store
	err := multiStore.RegisterEngine(engine)
	assert.NoError(t, err)

	err = multiStore.CreateStore("test-store", "memory", nil)
	assert.NoError(t, err)

	// Test get store
	store, err := multiStore.GetStore("test-store")
	assert.NoError(t, err)
	assert.NotNil(t, store)
	assert.Equal(t, "test-store", store.Name())

	// Test get non-existent store
	_, err = multiStore.GetStore("non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "storage.store_not_found")
}

func TestMultiStoreDeleteStore(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "multi-store",
		Name: "Multi Store",
		Type: component.TypeService,
	}

	multiStore := infraStorage.NewMultiStore(config, "/tmp/multistore")
	engine := memoryStorage.NewEngine()

	// Register engine and create store
	err := multiStore.RegisterEngine(engine)
	assert.NoError(t, err)

	err = multiStore.CreateStore("test-store", "memory", nil)
	assert.NoError(t, err)

	// Verify store exists
	stores := multiStore.ListStores()
	assert.Len(t, stores, 1)

	// Test delete store
	err = multiStore.DeleteStore("test-store")
	assert.NoError(t, err)

	// Verify store is removed
	stores = multiStore.ListStores()
	assert.Empty(t, stores)

	// Test delete non-existent store
	err = multiStore.DeleteStore("non-existent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "storage.store_not_found")
}

func TestMultiStoreLifecycle(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "multi-store",
		Name: "Multi Store",
		Type: component.TypeService,
	}

	multiStore := infraStorage.NewMultiStore(config, "/tmp/multistore")
	ctx := infraContext.NewContext()

	// Test start
	err := multiStore.Start(ctx)
	assert.NoError(t, err)
	assert.True(t, multiStore.IsRunning())
	assert.Equal(t, component.StatusRunning, multiStore.Status())

	// Test stop
	err = multiStore.Stop(ctx)
	assert.NoError(t, err)
	assert.False(t, multiStore.IsRunning())
	assert.Equal(t, component.StatusStopped, multiStore.Status())
}

func TestMultiStoreInitializeAndDispose(t *testing.T) {
	config := component.ComponentConfig{
		ID:   "multi-store",
		Name: "Multi Store",
		Type: component.TypeService,
	}

	multiStore := infraStorage.NewMultiStore(config, "/tmp/multistore")
	factory := mocks.NewFactory()
	mockSystem := factory.SystemInterface()
	ctx := infraContext.NewContext()

	// Test initialization
	err := multiStore.Initialize(ctx, mockSystem)
	assert.NoError(t, err)

	// Test disposal
	err = multiStore.Dispose()
	assert.NoError(t, err)
}
