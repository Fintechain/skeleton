package system

import (
	"strings"
	"testing"
	"time"

	"github.com/ebanfa/skeleton/internal/domain/plugin"
	pkgSystem "github.com/ebanfa/skeleton/pkg/system"
	"github.com/ebanfa/skeleton/test/integration/system/testdata"
)

func TestPluginIntegration_Registration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Create test plugin with specific behavior
	testPlugin := testdata.NewTestPlugin("registration-test", "1.0.0")

	// Track plugin registration
	plugins := []plugin.Plugin{testPlugin}

	// Start system and verify plugin registration
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
			t.Fatalf("StartSystem() failed: %v", err)
		}

		// Verify plugin was properly registered and loaded
		if testPlugin, ok := testPlugin.(*testdata.TestPlugin); ok && !testPlugin.IsLoaded() {
			t.Error("Plugin should be loaded after successful registration")
		}
	case <-ctx.Done():
		t.Error("Plugin registration test timed out")
	}
}

func TestPluginIntegration_MultiplePlugins(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Create multiple test plugins
	plugins := []plugin.Plugin{
		testdata.NewTestPlugin("multi-test-1", "1.0.0"),
		testdata.NewTestPlugin("multi-test-2", "2.0.0"),
		testdata.NewTestPlugin("multi-test-3", "3.0.0"),
	}

	// Start system with multiple plugins
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
	case <-ctx.Done():
		t.Error("Multiple plugins test timed out")
	}
}

func TestPluginIntegration_FailureHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Create a plugin that will fail to load
	failingPlugin := testdata.CreateFailingPlugin("failing-plugin", "1.0.0", "simulated plugin load failure")
	plugins := []plugin.Plugin{failingPlugin}

	// Start system and expect failure
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
		if !strings.Contains(err.Error(), "plugin") && !strings.Contains(err.Error(), "simulated") {
			t.Errorf("Error should mention plugin failure, got: %v", err)
		}
	case <-ctx.Done():
		t.Error("Plugin failure handling test timed out")
	}
}

func TestPluginIntegration_PartialFailure(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Create mix of good and bad plugins
	goodPlugin := testdata.NewTestPlugin("good-plugin", "1.0.0")
	badPlugin := testdata.CreateFailingPlugin("bad-plugin", "1.0.0", "simulated failure")

	plugins := []plugin.Plugin{goodPlugin, badPlugin}

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
		// System should fail due to bad plugin
		if err == nil {
			t.Error("Expected system to fail with bad plugin")
		}

		// Verify good plugin is not loaded when system startup fails
		if testPlugin, ok := goodPlugin.(*testdata.TestPlugin); ok && testPlugin.IsLoaded() {
			t.Error("Good plugin should not be loaded when system startup fails")
		}
	case <-ctx.Done():
		t.Error("Partial failure test timed out")
	}
}

func TestPluginIntegration_ComponentRegistration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Create plugin with components
	testPlugin := testdata.NewTestPlugin("component-test", "1.0.0")
	plugins := []plugin.Plugin{testPlugin}

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
			t.Fatalf("StartSystem() failed: %v", err)
		}

		// Verify plugin is loaded
		if testPlugin, ok := testPlugin.(*testdata.TestPlugin); ok && !testPlugin.IsLoaded() {
			t.Error("Plugin should be loaded after successful registration")
		}

		// Verify plugin has components
		components := testPlugin.Components()
		if len(components) == 0 {
			t.Error("Plugin should have components")
		}

		// Verify components are properly structured
		for _, comp := range components {
			if comp.ID() == "" {
				t.Error("Component should have non-empty ID")
			}
			if comp.Name() == "" {
				t.Error("Component should have non-empty name")
			}
		}
	case <-ctx.Done():
		t.Error("Component registration test timed out")
	}
}

func TestPluginIntegration_SlowLoading(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Create plugin that takes time to load
	slowPlugin := testdata.CreateSlowPlugin("slow-plugin", "1.0.0", 2*time.Second)
	plugins := []plugin.Plugin{slowPlugin}

	// Use longer timeout for this test
	ctx, cancel := createTestContext(60 * time.Second)
	defer cancel()

	start := time.Now()
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
		elapsed := time.Since(start)

		if err != nil {
			t.Fatalf("StartSystem() with slow plugin failed: %v", err)
		}

		// Verify it took at least the expected delay
		if elapsed < 2*time.Second {
			t.Errorf("Expected startup to take at least 2 seconds, took %v", elapsed)
		}

		// Verify plugin is loaded despite being slow
		if testPlugin, ok := slowPlugin.(*testdata.TestPlugin); ok && !testPlugin.IsLoaded() {
			t.Error("Slow plugin should still be loaded")
		}
	case <-ctx.Done():
		t.Error("Slow loading test timed out")
	}
}

func TestPluginIntegration_EmptyPluginList(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Test with empty plugin list
	var plugins []plugin.Plugin

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
			t.Errorf("StartSystem() with empty plugin list failed: %v", err)
		}
	case <-ctx.Done():
		t.Error("Empty plugin list test timed out")
	}
}
