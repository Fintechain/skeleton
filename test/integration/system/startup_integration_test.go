package system

import (
	"strings"
	"testing"
	"time"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/plugin"
	"github.com/fintechain/skeleton/internal/infrastructure/event"
	"github.com/fintechain/skeleton/internal/infrastructure/system"
	pkgSystem "github.com/fintechain/skeleton/pkg/system"
	"github.com/fintechain/skeleton/test/integration/system/testdata"
)

func TestSystemStartup_AllDefaults(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test complete system startup with all defaults
	ctx, cancel := createTestContext(30 * time.Second)
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
	defer env.Cleanup()

	// Test system startup with custom configuration
	ctx, cancel := createTestContext(30 * time.Second)
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
	defer env.Cleanup()

	// Create test plugins
	plugins := []plugin.Plugin{
		testdata.NewTestPlugin("test-plugin-1", "1.0.0"),
		testdata.NewTestPlugin("test-plugin-2", "2.0.0"),
	}

	// Test system startup with plugins
	ctx, cancel := createTestContext(30 * time.Second)
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
	defer env.Cleanup()

	// Create custom dependencies using real implementations
	customRegistry := component.CreateRegistry()
	customEventBus := event.CreateEventBus()

	// Test system startup with custom dependencies
	ctx, cancel := createTestContext(30 * time.Second)
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

func TestSystemStartup_WithMixedDependencies(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Create some custom dependencies, leave others as defaults
	customRegistry := component.CreateRegistry()
	plugins := []plugin.Plugin{
		testdata.NewTestPlugin("mixed-test-plugin", "1.0.0"),
	}

	// Test system startup with mixed dependencies
	ctx, cancel := createTestContext(30 * time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		err := pkgSystem.StartSystem(
			pkgSystem.WithConfig(env.Config),
			pkgSystem.WithRegistry(customRegistry),
			pkgSystem.WithPlugins(plugins),
			// EventBus and MultiStore will use defaults
		)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("StartSystem() with mixed dependencies failed: %v", err)
		}

		// Verify plugin was loaded
		if testPlugin, ok := plugins[0].(*testdata.TestPlugin); ok {
			if !testPlugin.IsLoaded() {
				t.Error("Plugin should be loaded with mixed dependencies")
			}
		}
	case <-ctx.Done():
		t.Error("StartSystem() with mixed dependencies timed out")
	}
}

func TestSystemStartup_InvalidConfiguration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test system startup with invalid configuration
	invalidConfig := &system.Config{
		ServiceID: "", // Invalid: empty service ID
	}

	ctx, cancel := createTestContext(30 * time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		err := pkgSystem.StartSystem(
			pkgSystem.WithConfig(invalidConfig),
		)
		done <- err
	}()

	select {
	case err := <-done:
		// We expect this to either succeed (if empty ServiceID is handled)
		// or fail gracefully with a descriptive error
		if err != nil && err.Error() == "" {
			t.Error("Error should have descriptive message")
		}
	case <-ctx.Done():
		t.Error("StartSystem() with invalid config timed out")
	}
}

func TestSystemStartup_FailureScenarios(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Create a plugin that will fail to load
	failingPlugin := testdata.CreateFailingPlugin("failing-plugin", "1.0.0", "simulated plugin load failure")
	plugins := []plugin.Plugin{failingPlugin}

	ctx, cancel := createTestContext(30 * time.Second)
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
		if err == nil {
			t.Error("Expected StartSystem() to fail with failing plugin")
		}

		// Verify error contains plugin failure information
		if !strings.Contains(err.Error(), "plugin") {
			t.Errorf("Error should mention plugin failure, got: %v", err)
		}
	case <-ctx.Done():
		t.Error("StartSystem() with failing plugin timed out")
	}
}
