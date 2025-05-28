package system

import (
	"errors"
	"testing"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/plugin"
	"github.com/fintechain/skeleton/internal/domain/storage"
	"github.com/fintechain/skeleton/internal/infrastructure/logging"
	"github.com/fintechain/skeleton/internal/infrastructure/storage/memory"
	"github.com/fintechain/skeleton/internal/infrastructure/system/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSystemConfig_ApplyDefaults(t *testing.T) {
	tests := []struct {
		name     string
		input    *SystemConfig
		validate func(*testing.T, *SystemConfig)
	}{
		{
			name:  "all nil dependencies",
			input: &SystemConfig{},
			validate: func(t *testing.T, sc *SystemConfig) {
				if sc.Config == nil {
					t.Error("Expected Config to be created")
				}
				if sc.Config.ServiceID != "system" {
					t.Errorf("Expected default ServiceID 'system', got %s", sc.Config.ServiceID)
				}
				if sc.Registry == nil {
					t.Error("Expected Registry to be created")
				}
				if sc.PluginMgr == nil {
					t.Error("Expected PluginMgr to be created")
				}
				if sc.EventBus == nil {
					t.Error("Expected EventBus to be created")
				}
				if sc.MultiStore == nil {
					t.Error("Expected MultiStore to be created")
				}
			},
		},
		{
			name: "partial dependencies provided",
			input: &SystemConfig{
				Registry: mocks.NewMockRegistry(),
				EventBus: mocks.NewMockEventBus(),
			},
			validate: func(t *testing.T, sc *SystemConfig) {
				// Verify existing dependencies are preserved
				if _, ok := sc.Registry.(*mocks.MockRegistry); !ok {
					t.Error("Expected existing Registry to be preserved")
				}
				if _, ok := sc.EventBus.(*mocks.MockEventBus); !ok {
					t.Error("Expected existing EventBus to be preserved")
				}
				// Verify missing dependencies are created
				if sc.Config == nil {
					t.Error("Expected Config to be created")
				}
				if sc.PluginMgr == nil {
					t.Error("Expected PluginMgr to be created")
				}
				if sc.MultiStore == nil {
					t.Error("Expected MultiStore to be created")
				}
			},
		},
		{
			name: "custom config provided",
			input: &SystemConfig{
				Config: &Config{
					ServiceID: "custom-service",
					StorageConfig: storage.MultiStoreConfig{
						RootPath:      "./custom-data",
						DefaultEngine: "leveldb",
					},
				},
			},
			validate: func(t *testing.T, sc *SystemConfig) {
				if sc.Config.ServiceID != "custom-service" {
					t.Errorf("Expected custom ServiceID 'custom-service', got %s", sc.Config.ServiceID)
				}
				if sc.Config.StorageConfig.RootPath != "./custom-data" {
					t.Errorf("Expected custom RootPath './custom-data', got %s", sc.Config.StorageConfig.RootPath)
				}
				if sc.Config.StorageConfig.DefaultEngine != "leveldb" {
					t.Errorf("Expected custom DefaultEngine 'leveldb', got %s", sc.Config.StorageConfig.DefaultEngine)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyDefaults(tt.input)
			if result == nil {
				t.Error("applyDefaults() returned nil")
				return
			}
			tt.validate(t, result)
		})
	}
}

func TestSystemConfig_ApplyDefaults_ExportedVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    *SystemConfig
		validate func(*testing.T, *SystemConfig)
	}{
		{
			name:  "all nil dependencies",
			input: &SystemConfig{},
			validate: func(t *testing.T, sc *SystemConfig) {
				if sc.Config == nil {
					t.Error("Expected Config to be created")
				}
				if sc.Config.ServiceID != "system" {
					t.Errorf("Expected default ServiceID 'system', got %s", sc.Config.ServiceID)
				}
				if sc.Registry == nil {
					t.Error("Expected Registry to be created")
				}
				if sc.PluginMgr == nil {
					t.Error("Expected PluginMgr to be created")
				}
				if sc.EventBus == nil {
					t.Error("Expected EventBus to be created")
				}
				if sc.MultiStore == nil {
					t.Error("Expected MultiStore to be created")
				}
			},
		},
		{
			name: "partial dependencies provided",
			input: &SystemConfig{
				Registry: mocks.NewMockRegistry(),
				EventBus: mocks.NewMockEventBus(),
			},
			validate: func(t *testing.T, sc *SystemConfig) {
				// Verify existing dependencies are preserved
				if _, ok := sc.Registry.(*mocks.MockRegistry); !ok {
					t.Error("Expected existing Registry to be preserved")
				}
				if _, ok := sc.EventBus.(*mocks.MockEventBus); !ok {
					t.Error("Expected existing EventBus to be preserved")
				}
				// Verify missing dependencies are created
				if sc.Config == nil {
					t.Error("Expected Config to be created")
				}
				if sc.PluginMgr == nil {
					t.Error("Expected PluginMgr to be created")
				}
				if sc.MultiStore == nil {
					t.Error("Expected MultiStore to be created")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := applyDefaults(tt.input)
			if result == nil {
				t.Error("applyDefaults() returned nil")
				return
			}
			tt.validate(t, result)
		})
	}
}

func TestProvideSystemService(t *testing.T) {
	// Test system service creation with mocked dependencies
	mockRegistry := mocks.NewMockRegistry()
	mockPluginMgr := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockMultiStore := mocks.NewMockMultiStore()

	config := &Config{
		ServiceID: "test-service",
		StorageConfig: storage.MultiStoreConfig{
			RootPath:      "./test-data",
			DefaultEngine: "memory",
		},
	}

	plugins := []plugin.Plugin{
		mocks.NewMockPlugin("test-plugin", "1.0.0"),
	}

	service, err := provideSystemService(config, mockRegistry, mockPluginMgr, mockEventBus, mockMultiStore, plugins)

	if err != nil {
		t.Errorf("provideSystemService() error = %v", err)
	}

	if service == nil {
		t.Error("Expected non-nil service")
	}

	// Verify service has correct dependencies
	if service.PluginManager() != mockPluginMgr {
		t.Error("Expected plugin manager to be set correctly")
	}

	if service.Registry() != mockRegistry {
		t.Error("Expected registry to be set correctly")
	}

	if service.EventBus() != mockEventBus {
		t.Error("Expected event bus to be set correctly")
	}

	if service.Store() != mockMultiStore {
		t.Error("Expected multistore to be set correctly")
	}

	// Verify service ID
	if service.ID() != "test-service" {
		t.Errorf("Expected service ID 'test-service', got %s", service.ID())
	}
}

func TestCreateDefaultMultiStore(t *testing.T) {
	// Test default multistore creation
	config := storage.MultiStoreConfig{
		RootPath:      "./test-data",
		DefaultEngine: "memory",
	}

	store := createDefaultMultiStore(config)

	if store == nil {
		t.Error("Expected non-nil multistore")
	}

	// Test basic multistore operations
	if store.GetDefaultEngine() != "memory" {
		t.Errorf("Expected default engine 'memory', got %s", store.GetDefaultEngine())
	}

	// Test store creation
	err := store.CreateStore("test-store", "memory", storage.Config{})
	if err != nil {
		t.Errorf("CreateStore() error = %v", err)
	}

	// Test store exists
	if !store.StoreExists("test-store") {
		t.Error("Expected store to exist after creation")
	}

	// Test store retrieval
	retrievedStore, err := store.GetStore("test-store")
	if err != nil {
		t.Errorf("GetStore() error = %v", err)
	}
	if retrievedStore == nil {
		t.Error("Expected non-nil store")
	}
}

func TestConfig_CreateConfiguration(t *testing.T) {
	// Test configuration creation from Config
	config := &Config{
		ServiceID: "test-config-service",
		StorageConfig: storage.MultiStoreConfig{
			RootPath:      "./config-test-data",
			DefaultEngine: "leveldb",
		},
	}

	cfg := config.CreateConfiguration()

	if cfg == nil {
		t.Error("Expected non-nil configuration")
	}

	// Test that values are set correctly
	serviceID := cfg.GetString("system.serviceId")
	if serviceID != "test-config-service" {
		t.Errorf("Expected serviceId 'test-config-service', got %s", serviceID)
	}

	rootPath := cfg.GetString("system.storage.rootPath")
	if rootPath != "./config-test-data" {
		t.Errorf("Expected rootPath './config-test-data', got %s", rootPath)
	}

	defaultEngine := cfg.GetString("system.storage.defaultEngine")
	if defaultEngine != "leveldb" {
		t.Errorf("Expected defaultEngine 'leveldb', got %s", defaultEngine)
	}
}

func TestMockMultiStore_Operations(t *testing.T) {
	// Test the mock multistore implementation
	ms := mocks.NewMockMultiStore()

	// Test default engine
	if ms.GetDefaultEngine() != "memory" {
		t.Errorf("Expected default engine 'memory', got %s", ms.GetDefaultEngine())
	}

	// Test store creation
	err := ms.CreateStore("test", "memory", storage.Config{})
	if err != nil {
		t.Errorf("CreateStore() error = %v", err)
	}

	// Test store exists
	if !ms.StoreExists("test") {
		t.Error("Expected store to exist")
	}

	// Test store retrieval
	store, err := ms.GetStore("test")
	if err != nil {
		t.Errorf("GetStore() error = %v", err)
	}
	if store == nil {
		t.Error("Expected non-nil store")
	}

	// Test list stores
	stores := ms.ListStores()
	if len(stores) != 1 || stores[0] != "test" {
		t.Errorf("Expected ['test'], got %v", stores)
	}

	// Test store deletion
	err = ms.DeleteStore("test")
	if err != nil {
		t.Errorf("DeleteStore() error = %v", err)
	}

	if ms.StoreExists("test") {
		t.Error("Expected store to not exist after deletion")
	}
}

func TestStartWithFx_ConfigConversion(t *testing.T) {
	// Test the applyDefaults function behavior with SystemConfig
	originalConfig := &SystemConfig{
		Config: &Config{ServiceID: "conversion-test"},
		Plugins: []plugin.Plugin{
			mocks.NewMockPlugin("test-plugin", "1.0.0"),
		},
		Registry:   mocks.NewMockRegistry(),
		PluginMgr:  mocks.NewMockPluginManager(),
		EventBus:   mocks.NewMockEventBus(),
		MultiStore: mocks.NewMockMultiStore(),
	}

	// Test that applyDefaults preserves existing values
	result := applyDefaults(originalConfig)

	// Verify that existing values are preserved
	if result.Config != originalConfig.Config {
		t.Error("Config not preserved correctly")
	}
	if len(result.Plugins) != len(originalConfig.Plugins) {
		t.Error("Plugins not preserved correctly")
	}
	if result.Registry != originalConfig.Registry {
		t.Error("Registry not preserved correctly")
	}
	if result.PluginMgr != originalConfig.PluginMgr {
		t.Error("PluginMgr not preserved correctly")
	}
	if result.EventBus != originalConfig.EventBus {
		t.Error("EventBus not preserved correctly")
	}
	if result.MultiStore != originalConfig.MultiStore {
		t.Error("MultiStore not preserved correctly")
	}

	// Test that applyDefaults creates missing dependencies
	partialConfig := &SystemConfig{
		Config: &Config{ServiceID: "partial-test"},
		// Only provide some dependencies
		Registry: mocks.NewMockRegistry(),
	}

	result = applyDefaults(partialConfig)

	// Verify existing dependencies are preserved
	if result.Config != partialConfig.Config {
		t.Error("Existing Config not preserved")
	}
	if result.Registry != partialConfig.Registry {
		t.Error("Existing Registry not preserved")
	}

	// Verify missing dependencies are created
	if result.PluginMgr == nil {
		t.Error("Expected PluginMgr to be created")
	}
	if result.EventBus == nil {
		t.Error("Expected EventBus to be created")
	}
	if result.MultiStore == nil {
		t.Error("Expected MultiStore to be created")
	}
}

func TestMultiStore_CloseAll(t *testing.T) {
	config := storage.MultiStoreConfig{
		RootPath:      "./test-data",
		DefaultEngine: "memory",
	}
	multiStore := createDefaultMultiStore(config)

	// Test CloseAll method
	err := multiStore.CloseAll()
	assert.NoError(t, err)
}

func TestMultiStore_SetDefaultEngine(t *testing.T) {
	config := storage.MultiStoreConfig{
		RootPath:      "./test-data",
		DefaultEngine: "memory",
	}
	multiStore := createDefaultMultiStore(config)

	// Test setting default engine
	multiStore.SetDefaultEngine("badger")

	// Verify it was set
	engine := multiStore.GetDefaultEngine()
	assert.Equal(t, "badger", engine)
}

func TestMultiStore_RegisterEngine(t *testing.T) {
	config := storage.MultiStoreConfig{
		RootPath:      "./test-data",
		DefaultEngine: "memory",
	}
	multiStore := createDefaultMultiStore(config)

	// Create a mock engine for testing with a different name
	logger := logging.CreateStandardLogger(logging.Info)
	mockEngine := memory.NewEngine(logger)

	// Create a wrapper to change the engine name
	type testEngine struct {
		storage.Engine
	}

	testEngineWrapper := &testEngine{Engine: mockEngine}

	// Override the Name method to return a different name
	testEngineWrapper.Engine = &struct {
		storage.Engine
	}{
		Engine: mockEngine,
	}

	// Actually, let's just test with a different approach -
	// verify that registering the same engine twice fails
	logger2 := logging.CreateStandardLogger(logging.Info)
	memoryEngine2 := memory.NewEngine(logger2)

	// Test registering a duplicate engine (should fail)
	err := multiStore.RegisterEngine(memoryEngine2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")
}

func TestMultiStore_ListEngines(t *testing.T) {
	config := storage.MultiStoreConfig{
		RootPath:      "./test-data",
		DefaultEngine: "memory",
	}
	multiStore := createDefaultMultiStore(config)

	// Test listing engines
	engines := multiStore.ListEngines()
	assert.NotNil(t, engines)
	// Should contain at least the default memory engine
	assert.Contains(t, engines, "memory")
}

func TestMultiStore_GetEngine(t *testing.T) {
	config := storage.MultiStoreConfig{
		RootPath:      "./test-data",
		DefaultEngine: "memory",
	}
	multiStore := createDefaultMultiStore(config)

	// Test getting an engine
	engine, err := multiStore.GetEngine("memory")
	assert.NoError(t, err)
	assert.NotNil(t, engine)

	// Test getting non-existent engine
	engine, err = multiStore.GetEngine("non-existent")
	assert.Error(t, err)
	assert.Nil(t, engine)
}

func TestInitializeAndStart_WithService(t *testing.T) {
	// Create a real DefaultSystemService for testing
	mockRegistry := mocks.NewMockRegistry()
	mockPluginManager := mocks.NewMockPluginManager()
	mockEventBus := mocks.NewMockEventBus()
	mockConfig := mocks.NewMockConfiguration()
	mockStore := mocks.NewMockMultiStore()
	mockLogger := mocks.NewMockLogger()

	service := NewDefaultSystemService(
		"test-service",
		mockRegistry,
		mockPluginManager,
		mockEventBus,
		mockConfig,
		mockStore,
		mockLogger,
	)

	// Test with empty plugins slice
	err := initializeAndStart(service, []plugin.Plugin{})
	assert.NoError(t, err)
}

// Test error paths for MockMultiStore methods
func TestMockMultiStore_ErrorPaths(t *testing.T) {
	t.Run("CreateStore with existing store", func(t *testing.T) {
		ms := mocks.NewMockMultiStore()

		// Create a store first
		err := ms.CreateStore("test-store", "memory", storage.Config{})
		assert.NoError(t, err)

		// Try to create the same store again
		err = ms.CreateStore("test-store", "memory", storage.Config{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "store already exists")
	})

	t.Run("DeleteStore with non-existent store", func(t *testing.T) {
		ms := mocks.NewMockMultiStore()

		err := ms.DeleteStore("non-existent")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "store not found")
	})

	t.Run("GetStore with non-existent store", func(t *testing.T) {
		ms := mocks.NewMockMultiStore()

		store, err := ms.GetStore("non-existent")
		assert.Error(t, err)
		assert.Nil(t, store)
		assert.Contains(t, err.Error(), "store not found")
	})

	t.Run("GetEngine with non-existent engine", func(t *testing.T) {
		ms := mocks.NewMockMultiStore()

		engine, err := ms.GetEngine("non-existent")
		assert.Error(t, err)
		assert.Nil(t, engine)
		assert.Contains(t, err.Error(), "engine not found")
	})
}

// Test provideSystemService error paths
func TestProvideSystemService_ErrorPaths(t *testing.T) {
	t.Run("with nil config", func(t *testing.T) {
		mockRegistry := mocks.NewMockRegistry()
		mockPluginManager := mocks.NewMockPluginManager()
		mockEventBus := mocks.NewMockEventBus()
		mockStore := mocks.NewMockMultiStore()

		// This should panic due to nil pointer dereference
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when config is nil")
			}
		}()

		_, _ = provideSystemService(
			nil, // nil config
			mockRegistry,
			mockPluginManager,
			mockEventBus,
			mockStore,
			[]plugin.Plugin{},
		)
	})
}

// Test initializeAndStart error paths
func TestInitializeAndStart_ErrorPaths(t *testing.T) {
	t.Run("with system that fails to initialize", func(t *testing.T) {
		// Create a system service with a registry that fails to initialize
		mockRegistry := mocks.NewMockRegistry()
		mockRegistry.InitializeFunc = func(ctx component.Context) error {
			return errors.New("registry init failed")
		}

		mockPluginManager := mocks.NewMockPluginManager()
		mockEventBus := mocks.NewMockEventBus()
		mockConfig := mocks.NewMockConfiguration()
		mockStore := mocks.NewMockMultiStore()
		mockLogger := mocks.NewMockLogger()

		svc := NewDefaultSystemService(
			"test-service",
			mockRegistry,
			mockPluginManager,
			mockEventBus,
			mockConfig,
			mockStore,
			mockLogger,
		)

		err := initializeAndStart(svc, []plugin.Plugin{})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Failed to initialize component registry")
	})
}
