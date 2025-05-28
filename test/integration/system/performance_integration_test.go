package system

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fintechain/skeleton/internal/domain/plugin"
	pkgSystem "github.com/fintechain/skeleton/pkg/system"
	"github.com/fintechain/skeleton/test/integration/system/testdata"
)

func TestPerformanceIntegration_SystemStartupTime(t *testing.T) {
	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Measure startup time
	start := time.Now()

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
			t.Errorf("StartSystem() failed: %v", err)
		}

		startupTime := time.Since(start)

		// Startup should complete within reasonable time (adjust threshold as needed)
		if startupTime > 5*time.Second {
			t.Errorf("System startup took too long: %v", startupTime)
		}

		t.Logf("System startup time: %v", startupTime)
	case <-ctx.Done():
		t.Error("System startup time test timed out")
	}
}

func TestPerformanceIntegration_ConcurrentPluginRegistration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	const numPlugins = 100
	var wg sync.WaitGroup
	errors := make(chan error, numPlugins)

	start := time.Now()

	// Create plugins concurrently
	plugins := make([]plugin.Plugin, numPlugins)
	for i := 0; i < numPlugins; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			testPlugin := testdata.NewTestPlugin(fmt.Sprintf("concurrent-plugin-%d", id), "1.0.0")
			plugins[id] = testPlugin
		}(i)
	}

	wg.Wait()

	// Now test system startup with all plugins
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
			t.Errorf("StartSystem() with concurrent plugins failed: %v", err)
		}

		registrationTime := time.Since(start)

		t.Logf("Concurrent plugin registration time for %d plugins: %v", numPlugins, registrationTime)

		// Performance assertion - adjust threshold as needed
		if registrationTime > 10*time.Second {
			t.Errorf("Concurrent plugin registration took too long: %v", registrationTime)
		}
	case <-ctx.Done():
		t.Error("Concurrent plugin registration test timed out")
	}

	close(errors)

	// Check for errors
	for err := range errors {
		t.Errorf("Plugin registration error: %v", err)
	}
}

func TestPerformanceIntegration_MemoryUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Get initial memory stats
	var m1 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	// Register multiple plugins
	plugins := make([]plugin.Plugin, 50)
	for i := 0; i < 50; i++ {
		plugins[i] = testdata.NewTestPlugin(fmt.Sprintf("memory-test-plugin-%d", i), "1.0.0")
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
		if err != nil {
			t.Errorf("StartSystem() failed: %v", err)
		}

		// Force garbage collection and get final memory stats
		runtime.GC()
		var m2 runtime.MemStats
		runtime.ReadMemStats(&m2)

		memoryIncrease := m2.Alloc - m1.Alloc

		t.Logf("Memory increase: %d bytes", memoryIncrease)

		// Memory increase should be reasonable (adjust threshold as needed)
		if memoryIncrease > 50*1024*1024 { // 50MB threshold
			t.Errorf("Memory usage increased too much: %d bytes", memoryIncrease)
		}
	case <-ctx.Done():
		t.Error("Memory usage test timed out")
	}
}

func TestPerformanceIntegration_SystemShutdownTime(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Register some plugins to make shutdown more realistic
	plugins := make([]plugin.Plugin, 10)
	for i := 0; i < 10; i++ {
		plugins[i] = testdata.NewTestPlugin(fmt.Sprintf("shutdown-test-plugin-%d", i), "1.0.0")
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
		if err != nil {
			t.Errorf("StartSystem() failed: %v", err)
		}

		// Measure shutdown time (note: current system doesn't have explicit shutdown)
		start := time.Now()

		// Since there's no explicit shutdown in the current API, we'll just measure cleanup
		env.Cleanup()

		shutdownTime := time.Since(start)

		// Shutdown should complete within reasonable time
		if shutdownTime > 3*time.Second {
			t.Errorf("System shutdown took too long: %v", shutdownTime)
		}

		t.Logf("System shutdown time: %v", shutdownTime)
	case <-ctx.Done():
		t.Error("System shutdown time test timed out")
	}
}

func TestPerformanceIntegration_HighFrequencyOperations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Register a plugin first
	testPlugin := testdata.NewTestPlugin("high-frequency-test-plugin", "1.0.0")
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
			t.Errorf("StartSystem() failed: %v", err)
		}

		const numOperations = 1000
		start := time.Now()

		// Perform high-frequency operations (e.g., accessing plugin info)
		for i := 0; i < numOperations; i++ {
			// Since we don't have direct access to the system instance,
			// we'll simulate high-frequency operations by accessing plugin data
			_ = testPlugin.ID()
			_ = testPlugin.Version()
			_ = testPlugin.Components()
		}

		operationTime := time.Since(start)
		avgTimePerOp := operationTime / numOperations

		t.Logf("High-frequency operations: %d ops in %v (avg: %v per op)",
			numOperations, operationTime, avgTimePerOp)

		// Operations should be fast
		if avgTimePerOp > 1*time.Millisecond {
			t.Errorf("Operations are too slow: %v per operation", avgTimePerOp)
		}
	case <-ctx.Done():
		t.Error("High-frequency operations test timed out")
	}
}

func TestPerformanceIntegration_StressTest(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	const (
		numWorkers = 10
		duration   = 30 * time.Second
	)

	ctx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()

	var wg sync.WaitGroup
	errorCount := int64(0)
	operationCount := int64(0)

	// Start multiple workers performing various operations
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				default:
					// Perform various operations (create plugins and test system startup)
					testPlugin := testdata.NewTestPlugin(
						fmt.Sprintf("stress-plugin-%d-%d", workerID, time.Now().UnixNano()),
						"1.0.0",
					)

					// Test system startup with the plugin
					testCtx, testCancel := createTestContext(5 * time.Second)
					testDone := make(chan error, 1)
					go func() {
						err := pkgSystem.StartSystem(
							pkgSystem.WithConfig(env.Config),
							pkgSystem.WithPlugins([]plugin.Plugin{testPlugin}),
						)
						testDone <- err
					}()

					select {
					case err := <-testDone:
						if err != nil {
							atomic.AddInt64(&errorCount, 1)
						} else {
							atomic.AddInt64(&operationCount, 1)
						}
					case <-testCtx.Done():
						atomic.AddInt64(&errorCount, 1)
					}
					testCancel()

					// Small delay to prevent overwhelming the system
					time.Sleep(10 * time.Millisecond)
				}
			}
		}(i)
	}

	wg.Wait()

	t.Logf("Stress test completed: %d operations, %d errors in %v",
		operationCount, errorCount, duration)

	// Error rate should be low
	totalOps := operationCount + errorCount
	if totalOps > 0 {
		errorRate := float64(errorCount) / float64(totalOps)
		if errorRate > 0.01 { // Less than 1% error rate
			t.Errorf("Error rate too high during stress test: %.2f%%", errorRate*100)
		}
	}
}

func TestPerformanceIntegration_ResourceCleanup(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	env := setupTestEnvironment(t)
	defer env.Cleanup()

	// Get initial goroutine count
	initialGoroutines := runtime.NumGoroutine()

	// Register and test plugins multiple times
	for cycle := 0; cycle < 5; cycle++ {
		plugins := make([]plugin.Plugin, 10)
		for i := 0; i < 10; i++ {
			plugins[i] = testdata.NewTestPlugin(fmt.Sprintf("cleanup-test-plugin-%d-%d", cycle, i), "1.0.0")
		}

		ctx, cancel := createTestContext(30 * time.Second)
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
				t.Errorf("StartSystem() failed in cycle %d: %v", cycle, err)
			}
		case <-ctx.Done():
			t.Errorf("StartSystem() timed out in cycle %d", cycle)
		}
		cancel()

		// Simulate some work
		time.Sleep(100 * time.Millisecond)
	}

	// Allow some time for cleanup
	time.Sleep(1 * time.Second)
	runtime.GC()

	finalGoroutines := runtime.NumGoroutine()

	t.Logf("Goroutines: initial=%d, final=%d", initialGoroutines, finalGoroutines)

	// Should not have significant goroutine leaks
	if finalGoroutines > initialGoroutines+5 {
		t.Errorf("Potential goroutine leak detected: initial=%d, final=%d", initialGoroutines, finalGoroutines)
	}
}
