# Unit Testing Implementation for FX Integration

## Task Overview
Following the successful FX integration implementation (both Phase 1: PluginManager integration and Phase 2: FX Integration), we now need comprehensive unit tests using mock implementations. This focuses on creating appropriate mocks and unit tests to achieve 90% test coverage for our FX integration components.

## Target Components for Testing

### Primary Test Targets
1. **Public API Layer** (`pkg/system/system.go`)
   - Functional options behavior
   - Configuration validation
   - Error handling at API boundary

2. **FX Bootstrap Layer** (`internal/infrastructure/system/fx_bootstrap.go`)
   - Internal configuration logic
   - Default creation functions
   - Provider functions
   - System initialization logic

3. **Updated SystemService** (`internal/infrastructure/system/default_system_service.go`)
   - PluginManager integration
   - Constructor behavior with new dependency

## Mocking Strategy and Patterns

### 1. External Dependency Mocking
**CRITICAL**: Only create mocks for external dependencies from other packages.

**Required Mocks** (place in `internal/infrastructure/system/mocks/`):
- `component_registry_mock.go` - Mock for `component.Registry`
- `plugin_manager_mock.go` - Mock for `plugin.PluginManager`
- `event_bus_mock.go` - Mock for `event.EventBus`
- `multistore_mock.go` - Mock for `storage.MultiStore`
- `config_mock.go` - Mock for `config.Configuration`
- `logger_mock.go` - Mock for `logging.Logger`
- `plugin_mock.go` - Mock for `plugin.Plugin`

Example Mock Pattern:
```go
// internal/infrastructure/system/mocks/plugin_manager_mock.go
package mocks

import "github.com/your-org/skeleton/internal/domain/plugin"

type MockPluginManager struct {
    RegisterPluginFunc func(plugin.Plugin) error
    GetPluginFunc      func(string) (plugin.Plugin, error)
    ListPluginsFunc    func() []plugin.Plugin
    
    // Call tracking
    RegisterPluginCalls []plugin.Plugin
    GetPluginCalls      []string
    ListPluginsCalls    int
}

func (m *MockPluginManager) RegisterPlugin(p plugin.Plugin) error {
    m.RegisterPluginCalls = append(m.RegisterPluginCalls, p)
    if m.RegisterPluginFunc != nil {
        return m.RegisterPluginFunc(p)
    }
    return nil
}

func (m *MockPluginManager) GetPlugin(id string) (plugin.Plugin, error) {
    m.GetPluginCalls = append(m.GetPluginCalls, id)
    if m.GetPluginFunc != nil {
        return m.GetPluginFunc(id)
    }
    return nil, nil
}

func (m *MockPluginManager) ListPlugins() []plugin.Plugin {
    m.ListPluginsCalls++
    if m.ListPluginsFunc != nil {
        return m.ListPluginsFunc()
    }
    return []plugin.Plugin{}
}
```

### 2. Internal Component Testing
For interfaces within the same package, create test doubles within test files themselves.

## Required Test Files

### 1. Public API Tests (`pkg/system/system_test.go`)
```go
package system

import (
    "testing"
    "github.com/your-org/skeleton/internal/domain/plugin"
    "github.com/your-org/skeleton/internal/domain/storage"
)

// Test functional options
func TestWithConfig(t *testing.T) {
    config := &Config{ServiceID: "test"}
    option := WithConfig(config)
    
    sc := &systemConfig{}
    option(sc)
    
    if sc.config != config {
        t.Errorf("Expected config to be set")
    }
}

func TestWithPlugins(t *testing.T) {
    // Test plugin option setting
}

func TestWithRegistry(t *testing.T) {
    // Test registry option setting
}

func TestWithPluginManager(t *testing.T) {
    // Test plugin manager option setting
}

func TestWithEventBus(t *testing.T) {
    // Test event bus option setting
}

func TestWithMultiStore(t *testing.T) {
    // Test multistore option setting
}

// Test StartSystem behavior
func TestStartSystem_DefaultConfiguration(t *testing.T) {
    // Test that StartSystem works with no options
    // Mock fx.New and verify default creation
}

func TestStartSystem_WithCustomOptions(t *testing.T) {
    // Test StartSystem with various option combinations
}

func TestStartSystem_ErrorHandling(t *testing.T) {
    // Test error propagation from fx layer
}
```

### 2. FX Bootstrap Tests (`internal/infrastructure/system/fx_bootstrap_test.go`)
```go
package system

import (
    "testing"
    "github.com/your-org/skeleton/internal/infrastructure/system/mocks"
)

func TestSystemConfig_ApplyDefaults(t *testing.T) {
    tests := []struct {
        name     string
        input    *systemConfig
        validate func(*testing.T, *systemConfig)
    }{
        {
            name:  "all nil dependencies",
            input: &systemConfig{},
            validate: func(t *testing.T, sc *systemConfig) {
                if sc.registry == nil {
                    t.Error("Expected registry to be created")
                }
                if sc.pluginMgr == nil {
                    t.Error("Expected plugin manager to be created")
                }
                // Validate all defaults are created
            },
        },
        {
            name: "partial dependencies provided",
            input: &systemConfig{
                registry: &mocks.MockRegistry{},
            },
            validate: func(t *testing.T, sc *systemConfig) {
                // Verify existing dependencies are preserved
                // Verify missing dependencies are created
            },
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.input.applyDefaults()
            if err != nil {
                t.Errorf("applyDefaults() error = %v", err)
            }
            tt.validate(t, tt.input)
        })
    }
}

func TestProvideSystemService(t *testing.T) {
    // Test system service creation with mocked dependencies
    mockRegistry := &mocks.MockRegistry{}
    mockPluginMgr := &mocks.MockPluginManager{}
    mockEventBus := &mocks.MockEventBus{}
    mockMultiStore := &mocks.MockMultiStore{}
    
    config := &Config{ServiceID: "test"}
    
    service, err := provideSystemService(config, mockRegistry, mockPluginMgr, mockEventBus, mockMultiStore)
    
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
}

func TestInitializeAndStart(t *testing.T) {
    // Test system initialization and plugin registration
    mockService := &MockSystemService{}
    testPlugins := []plugin.Plugin{
        &mocks.MockPlugin{IDValue: "test-plugin-1"},
        &mocks.MockPlugin{IDValue: "test-plugin-2"},
    }
    
    err := initializeAndStart(mockService, testPlugins)
    
    if err != nil {
        t.Errorf("initializeAndStart() error = %v", err)
    }
    
    // Verify Initialize was called
    if !mockService.InitializeCalled {
        t.Error("Expected Initialize to be called")
    }
    
    // Verify Start was called
    if !mockService.StartCalled {
        t.Error("Expected Start to be called")
    }
    
    // Verify plugins were registered
    if len(mockService.RegisteredPlugins) != 2 {
        t.Errorf("Expected 2 plugins registered, got %d", len(mockService.RegisteredPlugins))
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
}

func TestCreateDefaultConfig(t *testing.T) {
    // Test default configuration creation
    config := createDefaultConfig()
    
    if config.ServiceID == "" {
        t.Error("Expected non-empty ServiceID")
    }
    
    if config.StorageConfig.DefaultEngine == "" {
        t.Error("Expected non-empty DefaultEngine")
    }
}
```

### 3. SystemService Integration Tests (`internal/infrastructure/system/default_system_service_test.go`)
```go
package system

import (
    "testing"
    "github.com/your-org/skeleton/internal/infrastructure/system/mocks"
)

func TestNewDefaultSystemService_WithPluginManager(t *testing.T) {
    // Test that constructor properly accepts PluginManager
    mockRegistry := &mocks.MockRegistry{}
    mockPluginMgr := &mocks.MockPluginManager{}
    mockEventBus := &mocks.MockEventBus{}
    mockMultiStore := &mocks.MockMultiStore{}
    mockConfig := &mocks.MockConfig{}
    mockLogger := &mocks.MockLogger{}
    
    service := NewDefaultSystemService(
        "test-service",
        mockRegistry,
        mockPluginMgr,
        mockEventBus,
        mockConfig,
        mockMultiStore,
        mockLogger,
    )
    
    if service == nil {
        t.Error("Expected non-nil service")
    }
    
    // Test PluginManager getter
    if service.PluginManager() != mockPluginMgr {
        t.Error("Expected PluginManager to be accessible")
    }
}

func TestDefaultSystemService_PluginManagerGetter(t *testing.T) {
    // Test the PluginManager() getter method
    mockPluginMgr := &mocks.MockPluginManager{}
    
    service := &DefaultSystemService{
        pluginManager: mockPluginMgr,
    }
    
    result := service.PluginManager()
    
    if result != mockPluginMgr {
        t.Error("Expected PluginManager getter to return correct instance")
    }
}
```

## Test Coverage Requirements

### 1. Target: 90% Test Coverage
Use `go test -cover` to verify coverage:
```bash
# Test individual packages
go test -cover ./pkg/system/...
go test -cover ./internal/infrastructure/system/...

# Generate detailed coverage report
go test -coverprofile=coverage.out ./pkg/system/... ./internal/infrastructure/system/...
go tool cover -html=coverage.out -o coverage.html
```

### 2. Coverage Focus Areas
- **All public functions** in `pkg/system/system.go`
- **All internal functions** in `fx_bootstrap.go`
- **Constructor and getter methods** in SystemService
- **Error handling paths**
- **Edge cases and nil checks**

## Mock File Organization

```
internal/infrastructure/system/
├── mocks/
│   ├── component_registry_mock.go
│   ├── plugin_manager_mock.go
│   ├── event_bus_mock.go
│   ├── multistore_mock.go
│   ├── config_mock.go
│   ├── logger_mock.go
│   └── plugin_mock.go
├── fx_bootstrap.go
├── fx_bootstrap_test.go
├── default_system_service.go
└── default_system_service_test.go

pkg/system/
├── system.go
└── system_test.go
```

## Specific Testing Requirements

### 1. Unit Tests Only
- Focus on testing business logic in isolation
- Mock all external dependencies
- No real fx.New() calls in unit tests
- Fast execution (under 1 second per test)

### 2. Error Scenario Testing
- Test error propagation from dependencies
- Test invalid configuration scenarios
- Test nil pointer handling
- Test error recovery mechanisms

### 3. Behavioral Verification
- Verify mock calls with correct parameters
- Test call ordering where important
- Verify side effects (like plugin registration)
- Test state changes

## Implementation Steps

1. **Analyze Dependencies**: Review fx_bootstrap.go and identify all external dependencies
2. **Create Mocks**: Implement mocks for each external interface in mocks/ directory
3. **Write Public API Tests**: Test all functional options and StartSystem behavior
4. **Write Bootstrap Tests**: Test internal configuration and provider logic
5. **Write SystemService Tests**: Test PluginManager integration
6. **Verify Coverage**: Ensure 90% coverage target is met
7. **Document Approach**: Explain test organization and challenging areas

## Success Criteria
- [ ] 90%+ test coverage achieved
- [ ] All external dependencies properly mocked
- [ ] All public API functions tested
- [ ] All internal bootstrap logic tested
- [ ] PluginManager integration verified
- [ ] Error scenarios comprehensively tested
- [ ] Tests execute quickly (unit test characteristic)
- [ ] Clear test organization and documentation

This unit testing implementation ensures our FX integration is thoroughly verified through automated tests while maintaining the isolation and speed characteristics of proper unit tests.