package system

import (
	"testing"
	"time"

	"github.com/ebanfa/skeleton/internal/domain/storage"
	"github.com/ebanfa/skeleton/internal/infrastructure/system"
	pkgSystem "github.com/ebanfa/skeleton/pkg/system"
)

func TestConfigurationIntegration_DefaultBehavior(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test that system works with completely default configuration
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
			t.Errorf("StartSystem() with defaults failed: %v", err)
		}
	case <-ctx.Done():
		t.Error("Default configuration test timed out")
	}
}

func TestConfigurationIntegration_CustomOverrides(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Test custom configuration override
	customConfig := &system.Config{
		ServiceID: "custom-integration-test",
		StorageConfig: storage.MultiStoreConfig{
			RootPath:      env.TempDir,
			DefaultEngine: "memory",
		},
	}

	ctx, cancel := createTestContext(30 * time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		err := pkgSystem.StartSystem(
			pkgSystem.WithConfig(customConfig),
		)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("StartSystem() with custom config failed: %v", err)
		}
	case <-ctx.Done():
		t.Error("Custom configuration test timed out")
	}
}

func TestConfigurationIntegration_ServiceIDVariations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	testCases := []struct {
		name      string
		serviceID string
		shouldErr bool
	}{
		{
			name:      "normal service ID",
			serviceID: "test-service",
			shouldErr: false,
		},
		{
			name:      "service ID with numbers",
			serviceID: "test-service-123",
			shouldErr: false,
		},
		{
			name:      "service ID with underscores",
			serviceID: "test_service_integration",
			shouldErr: false,
		},
		{
			name:      "empty service ID",
			serviceID: "",
			shouldErr: false, // System should handle this gracefully
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			env := setupTestEnvironment(t)
			defer env.Cleanup()

			config := &system.Config{
				ServiceID: tc.serviceID,
				StorageConfig: storage.MultiStoreConfig{
					RootPath:      env.TempDir,
					DefaultEngine: "memory",
				},
			}

			ctx, cancel := createTestContext(30 * time.Second)
			defer cancel()

			done := make(chan error, 1)
			go func() {
				err := pkgSystem.StartSystem(
					pkgSystem.WithConfig(config),
				)
				done <- err
			}()

			select {
			case err := <-done:
				if tc.shouldErr && err == nil {
					t.Errorf("Expected error for service ID %q, but got none", tc.serviceID)
				}
				if !tc.shouldErr && err != nil {
					t.Errorf("Unexpected error for service ID %q: %v", tc.serviceID, err)
				}
			case <-ctx.Done():
				t.Errorf("Test timed out for service ID %q", tc.serviceID)
			}
		})
	}
}

func TestConfigurationIntegration_StorageEngineVariations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	testCases := []struct {
		name          string
		defaultEngine string
		shouldErr     bool
	}{
		{
			name:          "memory engine",
			defaultEngine: "memory",
			shouldErr:     false,
		},
		{
			name:          "empty engine",
			defaultEngine: "",
			shouldErr:     false, // Should use system default
		},
		{
			name:          "unknown engine",
			defaultEngine: "unknown-engine",
			shouldErr:     false, // System should handle gracefully
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			env := setupTestEnvironment(t)
			defer env.Cleanup()

			config := &system.Config{
				ServiceID: "storage-test",
				StorageConfig: storage.MultiStoreConfig{
					RootPath:      env.TempDir,
					DefaultEngine: tc.defaultEngine,
				},
			}

			ctx, cancel := createTestContext(30 * time.Second)
			defer cancel()

			done := make(chan error, 1)
			go func() {
				err := pkgSystem.StartSystem(
					pkgSystem.WithConfig(config),
				)
				done <- err
			}()

			select {
			case err := <-done:
				if tc.shouldErr && err == nil {
					t.Errorf("Expected error for engine %q, but got none", tc.defaultEngine)
				}
				if !tc.shouldErr && err != nil {
					t.Errorf("Unexpected error for engine %q: %v", tc.defaultEngine, err)
				}
			case <-ctx.Done():
				t.Errorf("Test timed out for engine %q", tc.defaultEngine)
			}
		})
	}
}

func TestConfigurationIntegration_StoragePathVariations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	testCases := []struct {
		name      string
		setupPath func(*TestEnvironment) string
		shouldErr bool
	}{
		{
			name: "valid temp directory",
			setupPath: func(env *TestEnvironment) string {
				return env.TempDir
			},
			shouldErr: false,
		},
		{
			name: "relative path",
			setupPath: func(env *TestEnvironment) string {
				return "./test-data"
			},
			shouldErr: false,
		},
		{
			name: "empty path",
			setupPath: func(env *TestEnvironment) string {
				return ""
			},
			shouldErr: false, // Should use system default
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			env := setupTestEnvironment(t)
			defer env.Cleanup()

			rootPath := tc.setupPath(env)

			config := &system.Config{
				ServiceID: "path-test",
				StorageConfig: storage.MultiStoreConfig{
					RootPath:      rootPath,
					DefaultEngine: "memory",
				},
			}

			ctx, cancel := createTestContext(30 * time.Second)
			defer cancel()

			done := make(chan error, 1)
			go func() {
				err := pkgSystem.StartSystem(
					pkgSystem.WithConfig(config),
				)
				done <- err
			}()

			select {
			case err := <-done:
				if tc.shouldErr && err == nil {
					t.Errorf("Expected error for path %q, but got none", rootPath)
				}
				if !tc.shouldErr && err != nil {
					t.Errorf("Unexpected error for path %q: %v", rootPath, err)
				}
			case <-ctx.Done():
				t.Errorf("Test timed out for path %q", rootPath)
			}
		})
	}
}

func TestConfigurationIntegration_NilConfiguration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test system behavior with nil configuration (should use defaults)
	ctx, cancel := createTestContext(30 * time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		err := pkgSystem.StartSystem(
			pkgSystem.WithConfig(nil),
		)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("StartSystem() with nil config failed: %v", err)
		}
	case <-ctx.Done():
		t.Error("Nil configuration test timed out")
	}
}

func TestConfigurationIntegration_ComplexConfiguration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Test with a complex configuration
	complexConfig := &system.Config{
		ServiceID: "complex-integration-test-with-long-name",
		StorageConfig: storage.MultiStoreConfig{
			RootPath:      env.TempDir + "/complex/nested/path",
			DefaultEngine: "memory",
		},
	}

	ctx, cancel := createTestContext(30 * time.Second)
	defer cancel()

	done := make(chan error, 1)
	go func() {
		err := pkgSystem.StartSystem(
			pkgSystem.WithConfig(complexConfig),
		)
		done <- err
	}()

	select {
	case err := <-done:
		if err != nil {
			t.Errorf("StartSystem() with complex config failed: %v", err)
		}
	case <-ctx.Done():
		t.Error("Complex configuration test timed out")
	}
}
