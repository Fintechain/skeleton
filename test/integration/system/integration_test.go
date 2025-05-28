package system

import (
	"context"
	"io"
	"os"
	"testing"
	"time"

	"github.com/fintechain/skeleton/internal/domain/storage"
	"github.com/fintechain/skeleton/internal/infrastructure/system"
)

// TestEnvironment provides isolated test environment for each integration test
type TestEnvironment struct {
	TempDir   string
	Config    *system.Config
	Resources []io.Closer
	CleanupFn func()
	Timeout   time.Duration
}

// setupTestEnvironment creates an isolated test environment for integration tests
func setupTestEnvironment(t *testing.T) *TestEnvironment {
	// Create temporary directory for test data
	tempDir, err := os.MkdirTemp("", "fx-integration-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	config := &system.Config{
		ServiceID: "integration-test",
		StorageConfig: storage.MultiStoreConfig{
			RootPath:      tempDir,
			DefaultEngine: "memory",
		},
	}

	env := &TestEnvironment{
		TempDir: tempDir,
		Config:  config,
		Timeout: 30 * time.Second,
		CleanupFn: func() {
			os.RemoveAll(tempDir)
		},
	}

	return env
}

// AddResource adds a resource to be cleaned up when the test environment is disposed
func (env *TestEnvironment) AddResource(resource io.Closer) {
	env.Resources = append(env.Resources, resource)
}

// Cleanup performs cleanup of all resources and temporary files
func (env *TestEnvironment) Cleanup() {
	// Close all resources
	for _, resource := range env.Resources {
		if err := resource.Close(); err != nil {
			// Log error but don't fail test cleanup
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

// createTestContext creates a context with timeout for integration tests
func createTestContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}
