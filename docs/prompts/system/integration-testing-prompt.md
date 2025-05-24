# Integration Testing Implementation for FX Integration

## Task Overview
Following the unit testing implementation, we now need comprehensive integration tests that verify the end-to-end behavior of our FX integration system. This focuses on testing real component interactions, actual fx behavior, and complete system startup scenarios.

## Integration Testing Scope

### What Integration Tests Should Cover
1. **End-to-End System Startup**: Real fx.New() calls with actual dependency injection
2. **Plugin Integration**: Real plugin registration and lifecycle management
3. **Component Interactions**: Actual communication between system components
4. **Configuration Handling**: Real configuration parsing and application
5. **Error Recovery**: System behavior under real failure conditions

### What Integration Tests Should NOT Cover
- Individual function logic (covered by unit tests)
- Mocked dependency behavior (covered by unit tests)
- Isolated component testing (covered by unit tests)

## Target Integration Scenarios

### 1. System Startup Integration
- Complete system startup with default configurations
- System startup with custom configurations
- System startup with mixed default/custom dependencies
- System startup failure scenarios

### 2. Plugin Integration
- Real plugin registration and activation
- Multiple plugin coordination
- Plugin failure handling and recovery
- Plugin lifecycle management

### 3. Dependency Integration
- Real component communication through interfaces
- Event propagation between components
- Storage operations with real MultiStore
- Configuration changes affecting multiple components

### 4. Error Integration
- Graceful error handling across component boundaries
- Error recovery and system stability
- Partial failure scenarios

## Test Infrastructure Setup

### 1. Test Environment Isolation
Create isolated test environments for each integration test:

```go
// internal/infrastructure/system/integration_test.go
package system

import (
    "os"
    "path/filepath"
    "testing"
    "context"
)

type TestEnvironment struct {
    TempDir    string
    Config     *Config
    CleanupFn  func()
}

func setupTestEnvironment(t *testing.T) *TestEnvironment {
    // Create temporary directory for test data
    tempDir, err := os.MkdirTemp("", "fx-integration-test-*")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    
    config := &Config{
        ServiceID: "integration-test",
        StorageConfig: storage.MultiStoreConfig{
            RootPath:      tempDir,
            DefaultEngine: "memory",
        },
    }
    
    return &TestEnvironment{
        TempDir: tempDir,
        Config:  config,
        CleanupFn: func() {
            os.RemoveAll(tempDir)
        },
    }
}
```

### 2. Test Plugin Implementation
Create real test plugins for integration scenarios:

```go
// internal/infrastructure/system/testdata/test_plugin.go
package testdata

import (
    "github.com/your-org/skeleton/internal/domain/component"
    "github.com/your-org/skeleton/internal/domain/plugin"
)

type TestPlugin struct {
    id         string
    version    string
    components []component.Component
    loaded     bool
    loadError  error
}

func NewTestPlugin(id, version string) *TestPlugin {
    return &TestPlugin{
        id:      id,
        version: version,
        components: []component.Component{
            &TestComponent{id: id + "-component"},
        },
    }
}

func (p *TestPlugin) ID() string { return p.id }
func (p *TestPlugin) Version() string { return p.version }

func (p *TestPlugin) Load(ctx component.Context, registry component.Registry) error {
    if p.loadError != nil {
        return p.loadError
    }
    
    // Register components with registry
    for _, comp := range p.components {
        if err := registry.Register(comp); err != nil {
            return err
        }
    }
    
    p.loaded = true
    return nil
}

func (p *TestPlugin) Unload(ctx component.Context) error {
    p.loaded = false
    return nil
}

func (p *TestPlugin) Components() []component.Component {
    return p.components
}

// Test helper methods
func (p *TestPlugin) IsLoaded() bool { return p.loaded }
func (p *TestPlugin) SetLoadError(err error) { p.loadError = err }
```

## Required Integration Test Files

### 1. System Startup Integration (`internal/infrastructure/system/startup_integration_test.go`)
```go
package system

import (
    "testing"
    "time"
    
    pkgSystem "github.com/your-org/skeleton/pkg/system"
    "github.com/your-org/skeleton/internal/infrastructure/system/testdata"
)

func TestSystemStartup_AllDefaults(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Test complete system startup with all defaults
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    done := make(chan error, 1)
    go func() {
        err := pkgSystem.StartSystem()
        done <- err
    }()
    
    select {
    case err := <-done:
        if err != nil {
            t.Errorf("StartSystem() failed: %v", err)
        }
    case <-ctx.Done():
        t.Error("StartSystem() timed out")
    }
}

func TestSystemStartup_WithCustomConfig(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    env := setupTestEnvironment(t)
    defer env.CleanupFn()
    
    // Test system startup with custom configuration
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    done := make(chan error, 1)
    go func() {
        err := pkgSystem.StartSystem(
            pkgSystem.WithConfig(env.Config),
        )
        done <- err
    }()
    
    select {
    case err := <-done:
        if err != nil {
            t.Errorf("StartSystem() with custom config failed: %v", err)
        }
    case <-ctx.Done():
        t.Error("StartSystem() with custom config timed out")
    }
}

func TestSystemStartup_WithPlugins(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    env := setupTestEnvironment(t)
    defer env.CleanupFn()
    
    // Create test plugins
    plugins := []plugin.Plugin{
        testdata.NewTestPlugin("test-plugin-1", "1.0.0"),
        testdata.NewTestPlugin("test-plugin-2", "2.0.0"),
    }
    
    // Test system startup with plugins
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    done := make(chan error, 1)
    go func() {
        err := pkgSystem.StartSystem(
            pkgSystem.WithConfig(env.Config),
            pkgSystem.WithPlugins(plugins),
        )
        done <- err
    }()
    
    select {
    case err := <-done:
        if err != nil {
            t.Errorf("StartSystem() with plugins failed: %v", err)
        }
        
        // Verify plugins were loaded
        for _, plugin := range plugins {
            if testPlugin, ok := plugin.(*testdata.TestPlugin); ok {
                if !testPlugin.IsLoaded() {
                    t.Errorf("Plugin %s was not loaded", plugin.ID())
                }
            }
        }
    case <-ctx.Done():
        t.Error("StartSystem() with plugins timed out")
    }
}

func TestSystemStartup_WithCustomDependencies(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    env := setupTestEnvironment(t)
    defer env.CleanupFn()
    
    // Create custom dependencies using real implementations
    customRegistry := component.NewDefaultRegistry()
    customEventBus := event.NewDefaultEventBus()
    
    // Test system startup with custom dependencies
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    done := make(chan error, 1)
    go func() {
        err := pkgSystem.StartSystem(
            pkgSystem.WithConfig(env.Config),
            pkgSystem.WithRegistry(customRegistry),
            pkgSystem.WithEventBus(customEventBus),
        )
        done <- err
    }()
    
    select {
    case err := <-done:
        if err != nil {
            t.Errorf("StartSystem() with custom dependencies failed: %v", err)
        }
    case <-ctx.Done():
        t.Error("StartSystem() with custom dependencies timed out")
    }
}
```

### 2. Plugin Integration Tests (`internal/infrastructure/system/plugin_integration_test.go`)
```go
package system

import (
    "testing"
    
    pkgSystem "github.com/your-org/skeleton/pkg/system"
    "github.com/your-org/skeleton/internal/infrastructure/system/testdata"
)

func TestPluginIntegration_Registration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    env := setupTestEnvironment(t)
    defer env.CleanupFn()
    
    // Create test plugin with specific behavior
    testPlugin := testdata.NewTestPlugin("registration-test", "1.0.0")
    
    // Track plugin registration
    plugins := []plugin.Plugin{testPlugin}
    
    // Start system and verify plugin registration
    err := pkgSystem.StartSystem(
        pkgSystem.WithConfig(env.Config),
        pkgSystem.WithPlugins(plugins),
    )
    
    if err != nil {
        t.Fatalf("StartSystem() failed: %v", err)
    }
    
    // Verify plugin was properly registered and loaded
    if !testPlugin.IsLoaded() {
        t.Error("Plugin was not loaded during system startup")
    }
}

func TestPluginIntegration_MultiplePlugins(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    env := setupTestEnvironment(t)
    defer env.CleanupFn()
    
    // Create multiple test plugins
    plugins := []plugin.Plugin{
        testdata.NewTestPlugin("multi-test-1", "1.0.0"),
        testdata.NewTestPlugin("multi-test-2", "2.0.0"),
        testdata.NewTestPlugin("multi-test-3", "3.0.0"),
    }
    
    // Start system with multiple plugins
    err := pkgSystem.StartSystem(
        pkgSystem.WithConfig(env.Config),
        pkgSystem.WithPlugins(plugins),
    )
    
    if err != nil {
        t.Fatalf("StartSystem() with multiple plugins failed: %v", err)
    }
    
    // Verify all plugins were loaded
    for _, plugin := range plugins {
        if testPlugin, ok := plugin.(*testdata.TestPlugin); ok {
            if !testPlugin.IsLoaded() {
                t.Errorf("Plugin %s was not loaded", plugin.ID())
            }
        }
    }
}

func TestPluginIntegration_FailureHandling(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    env := setupTestEnvironment(t)
    defer env.CleanupFn()
    
    // Create a plugin that will fail to load
    failingPlugin := testdata.NewTestPlugin("failing-plugin", "1.0.0")
    failingPlugin.SetLoadError(errors.New("simulated plugin load failure"))
    
    plugins := []plugin.Plugin{failingPlugin}
    
    // Start system and expect failure
    err := pkgSystem.StartSystem(
        pkgSystem.WithConfig(env.Config),
        pkgSystem.WithPlugins(plugins),
    )
    
    if err == nil {
        t.Error("Expected StartSystem() to fail with failing plugin")
    }
    
    // Verify error contains plugin failure information
    if !strings.Contains(err.Error(), "plugin load failure") {
        t.Errorf("Error should mention plugin failure, got: %v", err)
    }
}
```

### 3. Configuration Integration Tests (`internal/infrastructure/system/config_integration_test.go`)
```go
package system

import (
    "testing"
    
    pkgSystem "github.com/your-org/skeleton/pkg/system"
)

func TestConfigurationIntegration_DefaultBehavior(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Test that system works with completely default configuration
    err := pkgSystem.StartSystem()
    
    if err != nil {
        t.Errorf("StartSystem() with defaults failed: %v", err)
    }
}

func TestConfigurationIntegration_CustomOverrides(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    env := setupTestEnvironment(t)
    defer env.CleanupFn()
    
    // Test custom configuration override
    customConfig := &pkgSystem.Config{
        ServiceID: "custom-integration-test",
        StorageConfig: storage.MultiStoreConfig{
            RootPath:      env.TempDir,
            DefaultEngine: "memory",
        },
    }
    
    err := pkgSystem.StartSystem(
        pkgSystem.WithConfig(customConfig),
    )
    
    if err != nil {
        t.Errorf("StartSystem() with custom config failed: %v", err)
    }
}

func TestConfigurationIntegration_Validation(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Test configuration validation
    invalidConfig := &pkgSystem.Config{
        ServiceID: "", // Invalid: empty service ID
        StorageConfig: storage.MultiStoreConfig{
            RootPath:      "/invalid/path/that/does/not/exist",
            DefaultEngine: "unknown-engine",
        },
    }
    
    err := pkgSystem.StartSystem(
        pkgSystem.WithConfig(invalidConfig),
    )
    
    if err == nil {
        t.Error("Expected StartSystem() to fail with invalid config")
    }
}
```

### 4. Error Recovery Integration Tests (`internal/infrastructure/system/error_integration_test.go`)
```go
package system

import (
    "testing"
    "errors"
    
    pkgSystem "github.com/your-org/skeleton/pkg/system"
    "github.com/your-org/skeleton/internal/infrastructure/system/testdata"
)

func TestErrorIntegration_SystemStartupFailure(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Test system behavior when startup fails
    invalidConfig := &pkgSystem.Config{
        ServiceID: "",
        StorageConfig: storage.MultiStoreConfig{
            RootPath:      "/nonexistent/path",
            DefaultEngine: "invalid",
        },
    }
    
    err := pkgSystem.StartSystem(
        pkgSystem.WithConfig(invalidConfig),
    )
    
    if err == nil {
        t.Error("Expected system startup to fail with invalid configuration")
    }
    
    // Verify error is descriptive
    if err.Error() == "" {
        t.Error("Error should have descriptive message")
    }
}

func TestErrorIntegration_PartialFailure(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    env := setupTestEnvironment(t)
    defer env.CleanupFn()
    
    // Create mix of good and bad plugins
    goodPlugin := testdata.NewTestPlugin("good-plugin", "1.0.0")
    badPlugin := testdata.NewTestPlugin("bad-plugin", "1.0.0")
    badPlugin.SetLoadError(errors.New("simulated failure"))
    
    plugins := []plugin.Plugin{goodPlugin, badPlugin}
    
    err := pkgSystem.StartSystem(
        pkgSystem.WithConfig(env.Config),
        pkgSystem.WithPlugins(plugins),
    )
    
    // System should fail due to bad plugin
    if err == nil {
        t.Error("Expected system to fail with bad plugin")
    }
    
    // Good plugin should not be loaded due to failure
    if goodPlugin.IsLoaded() {
        t.Error("Good plugin should not be loaded when system startup fails")
    }
}
```

## Test Execution and Organization

### Test File Structure
```
internal/infrastructure/system/
├── integration_test.go              # Test environment setup
├── startup_integration_test.go      # System startup scenarios
├── plugin_integration_test.go       # Plugin integration scenarios
├── config_integration_test.go       # Configuration scenarios
├── error_integration_test.go        # Error recovery scenarios
└── testdata/
    ├── test_plugin.go               # Test plugin implementation
    ├── test_component.go            # Test component implementation
    └── test_configs/
        ├── valid_config.json        # Valid test configurations
        └── invalid_config.json      # Invalid test configurations
```

### Running Integration Tests

#### Development Execution
```bash
# Run all integration tests
go test -v ./internal/infrastructure/system/...

# Run integration tests with timeout
go test -timeout 60s ./internal/infrastructure/system/...

# Skip integration tests during development
go test -short ./internal/infrastructure/system/...

# Run specific integration test
go test -v ./internal/infrastructure/system/ -run TestSystemStartup_AllDefaults
```

#### CI/CD Execution
```bash
# Run integration tests in CI with extended timeout
go test -v -timeout 300s ./internal/infrastructure/system/...

# Generate coverage for integration tests
go test -coverprofile=integration.out ./internal/infrastructure/system/...
go tool cover -html=integration.out -o integration_coverage.html
```

## Performance and Resource Considerations

### 1. Test Isolation
```go
func TestResourceIsolation(t *testing.T) {
    // Each test should clean up its resources
    env := setupTestEnvironment(t)
    defer env.CleanupFn()
    
    // Test execution
    err := pkgSystem.StartSystem(
        pkgSystem.WithConfig(env.Config),
    )
    
    if err != nil {
        t.Errorf("Test failed: %v", err)
    }
    
    // Cleanup is automatic via defer
}
```

### 2. Timeout Management
```go
func TestTimeoutHandling(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Set reasonable timeouts for integration tests
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    done := make(chan error, 1)
    go func() {
        err := pkgSystem.StartSystem()
        done <- err
    }()
    
    select {
    case err := <-done:
        if err != nil {
            t.Errorf("System startup failed: %v", err)
        }
    case <-ctx.Done():
        t.Error("System startup timed out")
    }
}
```

### 3. Resource Cleanup
```go
type TestEnvironment struct {
    TempDir     string
    Config      *Config
    Resources   []io.Closer
    CleanupFn   func()
}

func (env *TestEnvironment) AddResource(resource io.Closer) {
    env.Resources = append(env.Resources, resource)
}

func (env *TestEnvironment) Cleanup() {
    // Close all resources
    for _, resource := range env.Resources {
        if err := resource.Close(); err != nil {
            log.Printf("Failed to close resource: %v", err)
        }
    }
    
    // Remove temporary files
    if env.TempDir != "" {
        os.RemoveAll(env.TempDir)
    }
    
    // Additional cleanup
    if env.CleanupFn != nil {
        env.CleanupFn()
    }
}
```

## Verification and Validation

### 1. System State Verification
```go
func verifySystemState(t *testing.T, expectedPluginCount int) {
    // This would require access to the actual system instance
    // For integration tests, we verify through observable behavior
    
    // Example: Check that expected number of plugins are loaded
    // Example: Verify configuration is applied correctly
    // Example: Check that services are running
}
```

### 2. Plugin Integration Verification
```go
func verifyPluginIntegration(t *testing.T, plugins []plugin.Plugin) {
    for _, plugin := range plugins {
        if testPlugin, ok := plugin.(*testdata.TestPlugin); ok {
            // Verify plugin is loaded
            if !testPlugin.IsLoaded() {
                t.Errorf("Plugin %s should be loaded", plugin.ID())
            }
            
            // Verify plugin components are registered
            components := testPlugin.Components()
            if len(components) == 0 {
                t.Errorf("Plugin %s should have components", plugin.ID())
            }
        }
    }
}
```

### 3. Error Scenario Verification
```go
func verifyErrorScenario(t *testing.T, err error, expectedErrorType string) {
    if err == nil {
        t.Errorf("Expected error of type %s, got nil", expectedErrorType)
        return
    }
    
    if !strings.Contains(err.Error(), expectedErrorType) {
        t.Errorf("Expected error containing %s, got: %v", expectedErrorType, err)
    }
}
```

## Success Criteria and Coverage Goals

### 1. Integration Test Coverage
- **System Startup**: All configuration combinations work correctly
- **Plugin Integration**: Plugins load, register components, and function properly
- **Error Handling**: System gracefully handles and reports failures
- **Configuration**: All config options are properly applied
- **Resource Management**: Proper cleanup and resource management

### 2. Performance Expectations
- **Startup Time**: System should start within 30 seconds with default config
- **Memory Usage**: No memory leaks during test execution
- **Resource Cleanup**: All temporary resources properly cleaned up
- **Test Execution**: Integration test suite completes within 5 minutes

### 3. Reliability Verification
- **Consistency**: Tests pass consistently across multiple runs
- **Environment Independence**: Tests work in different environments
- **Isolation**: Tests don't interfere with each other
- **Error Recovery**: System properly handles and reports error conditions

## Implementation Steps

1. **Setup Test Infrastructure**
   - Create test environment helpers
   - Implement test plugin framework
   - Set up resource management

2. **Implement System Startup Tests**
   - Test default configuration startup
   - Test custom configuration scenarios
   - Test mixed dependency scenarios

3. **Implement Plugin Integration Tests**
   - Test single plugin integration
   - Test multiple plugin coordination
   - Test plugin failure scenarios

4. **Implement Configuration Tests**
   - Test configuration validation
   - Test configuration override behavior
   - Test invalid configuration handling

5. **Implement Error Integration Tests**
   - Test system startup failures
   - Test partial failure scenarios
   - Test error recovery mechanisms

6. **Verify and Optimize**
   - Run complete test suite
   - Verify performance characteristics
   - Optimize test execution time
   - Document test approach and findings

## Documentation Requirements

### Test Organization Documentation
- Explain integration test structure and purpose
- Document test environment setup and cleanup
- Describe plugin testing framework

### Challenging Areas Documentation
- Complex error scenarios and how they're tested
- Plugin coordination testing approach
- Configuration validation testing strategy

### Performance and Resource Documentation
- Test execution time expectations
- Resource usage patterns
- Cleanup procedures and verification

This integration testing implementation ensures our FX integration works correctly in real-world scenarios while maintaining proper test isolation and resource management.