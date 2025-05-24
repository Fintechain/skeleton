# System Integration Tests

This directory contains comprehensive integration tests for the Fx-based system startup and plugin management functionality.

## Test Structure

### Core Test Files

- `integration_test.go` - Test environment setup and utilities
- `startup_integration_test.go` - System startup scenarios
- `plugin_integration_test.go` - Plugin loading and management
- `config_integration_test.go` - Configuration handling
- `error_integration_test.go` - Error recovery and handling
- `performance_integration_test.go` - Performance and stress tests
- `suite_test.go` - Test suite runner and common utilities

### Test Data

- `testdata/` - Test plugins and components
- `testdata/configs/` - Test configuration files

## Running Tests

### All Integration Tests
```bash
cd skeleton
go test ./test/integration/system/...
```

### Specific Test Categories
```bash
# Startup tests only
go test ./test/integration/system/ -run TestSystemStartup

# Plugin tests only
go test ./test/integration/system/ -run TestPluginIntegration

# Configuration tests only
go test ./test/integration/system/ -run TestConfigurationIntegration

# Error handling tests only
go test ./test/integration/system/ -run TestErrorIntegration
```

### Performance Tests
```bash
# Run with performance tests (longer duration)
go test ./test/integration/system/ -run TestPerformanceIntegration

# Skip performance tests
go test ./test/integration/system/ -short
```

### Verbose Output
```bash
go test ./test/integration/system/ -v
```

## Test Categories

### 1. System Startup Tests (`startup_integration_test.go`)

Tests various system startup scenarios:
- Default configuration startup
- Custom configuration startup
- Plugin loading during startup
- Mixed dependency scenarios
- Invalid configuration handling
- Failure scenarios

### 2. Plugin Integration Tests (`plugin_integration_test.go`)

Tests plugin lifecycle and management:
- Plugin registration and loading
- Multiple plugin handling
- Plugin failure scenarios
- Component registration
- Slow-loading plugins
- Empty plugin lists

### 3. Configuration Tests (`config_integration_test.go`)

Tests configuration handling:
- Default behavior
- Custom overrides
- Service ID variations
- Storage engine configurations
- Path variations
- Nil configuration handling
- Complex configurations

### 4. Error Recovery Tests (`error_integration_test.go`)

Tests error handling and recovery:
- System startup failures
- Partial failures
- Multiple plugin failures
- Dependency failures
- Error recovery scenarios
- Graceful error handling
- System stability under errors

### 5. Performance Tests (`performance_integration_test.go`)

Tests system performance characteristics:
- Startup time measurement
- Concurrent plugin registration
- Memory usage monitoring
- Shutdown time measurement
- High-frequency operations
- Stress testing
- Resource cleanup verification

## Test Environment

Each test uses an isolated test environment that:
- Creates temporary directories for test data
- Provides clean configuration
- Manages resource cleanup
- Handles timeouts appropriately
- Isolates tests from each other

## Test Data

### Test Plugins (`testdata/test_plugin.go`)

Provides configurable test plugins that can:
- Simulate successful loading
- Simulate loading failures
- Introduce loading delays
- Track loading state
- Register test components

### Test Components (`testdata/test_component.go`)

Provides test components that:
- Implement the component interface
- Track initialization state
- Support metadata operations
- Simulate component lifecycle

### Configuration Files (`testdata/configs/`)

Provides various configuration scenarios:
- `valid_config.json` - Standard valid configuration
- `invalid_config.json` - Invalid configuration for error testing
- `minimal_config.json` - Minimal configuration
- `complex_config.json` - Complex configuration with all options

## Environment Variables

- `TEST_MODE=integration` - Set during test execution
- `CI` - Detected for CI environment adjustments
- `GITHUB_ACTIONS` - Detected for GitHub Actions
- `SKIP_PERFORMANCE_TESTS` - Skip performance tests if set

## Timeouts

- Default test timeout: 30 seconds
- Performance test timeout: 60 seconds
- Stress test duration: 30 seconds

## Best Practices

1. **Isolation**: Each test uses its own test environment
2. **Cleanup**: All tests clean up resources properly
3. **Timeouts**: All tests have appropriate timeouts
4. **Error Handling**: Tests verify both success and failure scenarios
5. **Performance**: Performance tests have reasonable thresholds
6. **Documentation**: Each test is well-documented with clear intent

## Troubleshooting

### Tests Timing Out
- Check if system startup is hanging
- Verify plugin implementations don't have infinite loops
- Ensure proper cleanup in test teardown

### Memory Issues
- Run tests with `-v` to see memory usage logs
- Check for resource leaks in plugin implementations
- Verify proper cleanup of goroutines

### Plugin Loading Failures
- Check plugin implementation for errors
- Verify component registration logic
- Review error messages for specific issues

### Configuration Issues
- Verify test configuration files are valid JSON
- Check file permissions for test directories
- Ensure storage paths are accessible 