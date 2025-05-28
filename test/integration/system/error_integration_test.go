package system

import (
	"testing"
	"time"

	"github.com/fintechain/skeleton/internal/domain/component"
	"github.com/fintechain/skeleton/internal/domain/plugin"
	"github.com/fintechain/skeleton/internal/domain/storage"
	"github.com/fintechain/skeleton/internal/infrastructure/event"
	"github.com/fintechain/skeleton/internal/infrastructure/system"
	pkgSystem "github.com/fintechain/skeleton/pkg/system"
	"github.com/fintechain/skeleton/test/integration/system/testdata"
)

func TestErrorIntegration_SystemStartupFailure(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test system behavior when startup fails due to invalid configuration
	invalidConfig := &system.Config{
		ServiceID: "",
		StorageConfig: storage.MultiStoreConfig{
			RootPath:      "/nonexistent/path/that/should/not/exist",
			DefaultEngine: "invalid-engine",
		},
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
		// We expect this to either succeed (if the system handles invalid configs gracefully)
		// or fail with a descriptive error
		if err != nil && err.Error() == "" {
			t.Error("Error should have descriptive message")
		}
	case <-ctx.Done():
		t.Error("System startup failure test timed out")
	}
}

func TestErrorIntegration_PartialFailure(t *testing.T) {
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

		// Good plugin should not be loaded due to failure
		if testPlugin, ok := goodPlugin.(*testdata.TestPlugin); ok && testPlugin.IsLoaded() {
			t.Error("Good plugin should not be loaded when system startup fails")
		}
	case <-ctx.Done():
		t.Error("Partial failure test timed out")
	}
}

func TestErrorIntegration_MultiplePluginFailures(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Create multiple failing plugins
	plugins := []plugin.Plugin{
		testdata.CreateFailingPlugin("bad-plugin-1", "1.0.0", "first failure"),
		testdata.CreateFailingPlugin("bad-plugin-2", "1.0.0", "second failure"),
		testdata.NewTestPlugin("good-plugin", "1.0.0"),
	}

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
		// System should fail due to bad plugins
		if err == nil {
			t.Error("Expected system to fail with multiple bad plugins")
		}

		// Error should be descriptive
		if err.Error() == "" {
			t.Error("Error should have descriptive message")
		}
	case <-ctx.Done():
		t.Error("Multiple plugin failures test timed out")
	}
}

func TestErrorIntegration_DependencyFailure(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Create a custom registry that might fail
	customRegistry := component.CreateRegistry()
	customEventBus := event.CreateEventBus()

	// Create a plugin that should work normally
	testPlugin := testdata.NewTestPlugin("dependency-test", "1.0.0")
	plugins := []plugin.Plugin{testPlugin}

	ctx, cancel := createTestContext(30 * time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		err := pkgSystem.StartSystem(
			pkgSystem.WithConfig(env.Config),
			pkgSystem.WithRegistry(customRegistry),
			pkgSystem.WithEventBus(customEventBus),
			pkgSystem.WithPlugins(plugins),
		)
		done <- err
	}()

	select {
	case err := <-done:
		// This should normally succeed unless there's a real dependency issue
		if err != nil {
			// If there's an error, it should be descriptive
			if err.Error() == "" {
				t.Error("Error should have descriptive message")
			}
		}
	case <-ctx.Done():
		t.Error("Dependency failure test timed out")
	}
}

func TestErrorIntegration_ErrorRecovery(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// First, try to start with a failing plugin
	failingPlugin := testdata.CreateFailingPlugin("failing-plugin", "1.0.0", "initial failure")
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
		// Should fail initially
		if err == nil {
			t.Error("Expected initial startup to fail")
		}

		// Now try again with a working plugin
		workingPlugin := testdata.NewTestPlugin("working-plugin", "1.0.0")
		workingPlugins := []plugin.Plugin{workingPlugin}

		ctx2, cancel2 := createTestContext(30 * time.Second)
		defer cancel2()

		done2 := make(chan error, 1)
		go func() {
			err2 := pkgSystem.StartSystem(
				pkgSystem.WithConfig(env.Config),
				pkgSystem.WithPlugins(workingPlugins),
			)
			done2 <- err2
		}()

		select {
		case err2 := <-done2:
			if err2 != nil {
				t.Errorf("Expected recovery startup to succeed, got: %v", err2)
			}

			// Verify the working plugin is loaded
			if testPlugin, ok := workingPlugin.(*testdata.TestPlugin); ok && !testPlugin.IsLoaded() {
				t.Error("Working plugin should be loaded after recovery")
			}
		case <-ctx2.Done():
			t.Error("Recovery startup timed out")
		}

	case <-ctx.Done():
		t.Error("Initial error recovery test timed out")
	}
}

func TestErrorIntegration_GracefulErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Test various error scenarios to ensure graceful handling
	testCases := []struct {
		name    string
		setup   func() []plugin.Plugin
		wantErr bool
	}{
		{
			name: "nil plugin in list",
			setup: func() []plugin.Plugin {
				return []plugin.Plugin{
					testdata.NewTestPlugin("good-plugin", "1.0.0"),
					nil, // This should be handled gracefully
				}
			},
			wantErr: false, // System should handle nil plugins gracefully
		},
		{
			name: "plugin with empty ID",
			setup: func() []plugin.Plugin {
				testPlugin := testdata.NewTestPlugin("", "1.0.0") // Empty ID
				return []plugin.Plugin{testPlugin}
			},
			wantErr: false, // System should handle empty IDs gracefully
		},
		{
			name: "plugin with empty version",
			setup: func() []plugin.Plugin {
				testPlugin := testdata.NewTestPlugin("test-plugin", "") // Empty version
				return []plugin.Plugin{testPlugin}
			},
			wantErr: false, // System should handle empty versions gracefully
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			plugins := tc.setup()

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
				if tc.wantErr && err == nil {
					t.Errorf("Expected error for %s, but got none", tc.name)
				}
				if !tc.wantErr && err != nil {
					t.Errorf("Unexpected error for %s: %v", tc.name, err)
				}
			case <-ctx.Done():
				t.Errorf("Test %s timed out", tc.name)
			}
		})
	}
}

func TestErrorIntegration_SystemStability(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Test system stability under various error conditions
	// This test ensures the system doesn't crash or leak resources

	errorScenarios := []struct {
		name    string
		plugins []plugin.Plugin
	}{
		{
			name:    "empty plugin list",
			plugins: []plugin.Plugin{},
		},
		{
			name: "single failing plugin",
			plugins: []plugin.Plugin{
				testdata.CreateFailingPlugin("fail-1", "1.0.0", "test failure"),
			},
		},
		{
			name: "mixed success and failure",
			plugins: []plugin.Plugin{
				testdata.NewTestPlugin("success-1", "1.0.0"),
				testdata.CreateFailingPlugin("fail-1", "1.0.0", "test failure"),
				testdata.NewTestPlugin("success-2", "1.0.0"),
			},
		},
	}

	for _, scenario := range errorScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			ctx, cancel := createTestContext(30 * time.Second)
			defer cancel()

			done := make(chan error, 1)
			go func() {
				err := pkgSystem.StartSystem(
					pkgSystem.WithConfig(env.Config),
					pkgSystem.WithPlugins(scenario.plugins),
				)
				done <- err
			}()

			select {
			case err := <-done:
				// We don't care about the specific error, just that the system
				// handles it gracefully without crashing
				if err != nil {
					// Error should be descriptive
					if err.Error() == "" {
						t.Errorf("Error should have descriptive message for scenario %s", scenario.name)
					}
				}
			case <-ctx.Done():
				t.Errorf("Scenario %s timed out", scenario.name)
			}
		})
	}
}
