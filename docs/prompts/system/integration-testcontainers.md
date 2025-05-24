## 🐳 **Production-Like Integration Testing with Testcontainers - Discussion**

Based on my review of the codebase and the existing integration testing infrastructure, here's how we can implement a production-like container-based integration test that mimics real-world startup/shutdown behavior.

### 📋 **Current Architecture Context**

From the codebase analysis, I can see:

1. **Existing System**: Clean FX-based system with `pkg/system.StartSystem()` API
2. **Current Integration Tests**: File-based isolation in `test/integration/system/`
3. **Container Infrastructure**: Already exists in `gate-services` and `fiat-onramp` projects
4. **System Service**: Comprehensive lifecycle management with proper shutdown handling
5. **Missing Public Shutdown API**: No `StopSystem()` function exposed in public API

### 🎯 **Production-Like Testing Strategy**

#### **1. Component TestKit Module (Reusable Testing Framework)**

**A. Proposed Module Structure**
```
component-testkit/
├── go.mod                           # Independent Go module
├── container/
│   ├── manager.go                   # Container lifecycle management
│   ├── verification.go              # Startup/shutdown verification engine
│   ├── health.go                    # Health check utilities
│   └── config.go                    # Container configuration
├── fixtures/
│   ├── plugins.go                   # Test plugin implementations
│   ├── components.go                # Test component implementations
│   └── configs.go                   # Test configuration generators
├── patterns/
│   ├── startup_patterns.go          # Common startup verification patterns
│   ├── shutdown_patterns.go         # Common shutdown verification patterns
│   └── health_patterns.go           # Health check patterns
└── examples/
    ├── basic_app_test.go            # Example usage
    └── advanced_app_test.go         # Advanced patterns
```

**B. Reusable Container Manager**
```go
// Package container provides reusable container testing utilities
package container

import (
    "context"
    "time"
    "github.com/testcontainers/testcontainers-go"
)

// AppContainer represents a containerized application under test
type AppContainer struct {
    testcontainers.Container
    Config      *AppConfig
    Ports       map[string]int
    HealthURL   string
    MetricsURL  string
    ShutdownURL string
}

// AppConfig defines the configuration for the application container
type AppConfig struct {
    // Container settings
    DockerfilePath   string
    BuildContext     string
    ExposedPorts     []string
    Environment      map[string]string
    
    // Application settings
    ConfigFile       string
    HealthEndpoint   string
    MetricsEndpoint  string
    ShutdownEndpoint string
    
    // Verification settings
    Verification     *VerificationConfig
}

// VerificationConfig defines how to verify startup/shutdown
type VerificationConfig struct {
    // Health-based verification (primary)
    HealthChecks     []HealthCheckConfig
    MetricsChecks    []MetricsConfig
    
    // Log-based verification (optional/fallback)
    LogPatterns      []LogPattern
    LogRequired      bool  // If false, log failures won't fail the test
    
    // Timing configuration
    StartupTimeout   time.Duration
    ShutdownTimeout  time.Duration
    HealthTimeout    time.Duration
}

// LogPattern defines an optional log pattern to check
type LogPattern struct {
    Pattern     string
    Required    bool        // If false, missing pattern won't fail test
    Timeout     time.Duration
    Description string      // Human-readable description
}
```

**C. Verification Engine**
```go
// verification.go - Reusable verification logic
package container

// VerificationEngine handles startup and shutdown verification
type VerificationEngine struct {
    config *VerificationConfig
    logger Logger
}

// VerifyStartup verifies that the application started correctly
func (ve *VerificationEngine) VerifyStartup(container *AppContainer) error {
    // 1. Primary verification: Health endpoints (always required)
    if err := ve.verifyHealthEndpoints(container); err != nil {
        return fmt.Errorf("health verification failed: %w", err)
    }
    
    // 2. Secondary verification: Metrics endpoints
    if err := ve.verifyMetricsEndpoints(container); err != nil {
        return fmt.Errorf("metrics verification failed: %w", err)
    }
    
    // 3. Optional verification: Log patterns (non-blocking)
    if err := ve.verifyLogPatterns(container); err != nil {
        if ve.config.LogRequired {
            return fmt.Errorf("log verification failed: %w", err)
        } else {
            ve.logger.Warn("Log verification failed (non-critical): %v", err)
        }
    }
    
    return nil
}

// VerifyShutdown verifies that the application shut down gracefully
func (ve *VerificationEngine) VerifyShutdown(container *AppContainer) error {
    // 1. Verify health endpoints become unavailable
    if err := ve.waitForHealthUnavailable(container); err != nil {
        return fmt.Errorf("health endpoint still available: %w", err)
    }
    
    // 2. Verify container exits cleanly
    exitCode, err := container.Container.Wait(context.Background())
    if err != nil {
        return fmt.Errorf("container wait failed: %w", err)
    }
    
    if exitCode != 0 {
        return fmt.Errorf("unexpected exit code: got %d, want 0", exitCode)
    }
    
    // 3. Optional: Verify shutdown log patterns
    if err := ve.verifyShutdownLogs(container); err != nil {
        if ve.config.LogRequired {
            return fmt.Errorf("shutdown log verification failed: %w", err)
        } else {
            ve.logger.Warn("Shutdown log verification failed (non-critical): %v", err)
        }
    }
    
    return nil
}
```

#### **2. Enhanced Public API with Shutdown Support**

**A. Missing StopSystem Function**
```go
// pkg/system/system.go - Enhanced public API

// SystemHandle represents a running system that can be stopped
type SystemHandle struct {
    systemService system.SystemService
    shutdownFn    func() error
    ctx           context.Context
    cancel        context.CancelFunc
}

// StartSystem starts the system with the given options and returns a handle
func StartSystem(options ...Option) (*SystemHandle, error) {
    config := &system.SystemConfig{}
    for _, option := range options {
        option(config)
    }
    
    // Create context for system lifecycle
    ctx, cancel := context.WithCancel(context.Background())
    
    // Start system with FX
    systemService, shutdownFn, err := system.StartWithFxAndContext(ctx, config)
    if err != nil {
        cancel()
        return nil, err
    }
    
    return &SystemHandle{
        systemService: systemService,
        shutdownFn:    shutdownFn,
        ctx:           ctx,
        cancel:        cancel,
    }, nil
}

// StopSystem gracefully stops the system
func (h *SystemHandle) Stop() error {
    if h.shutdownFn != nil {
        if err := h.shutdownFn(); err != nil {
            return fmt.Errorf("system shutdown failed: %w", err)
        }
    }
    
    h.cancel()
    return nil
}

// SystemService returns the underlying system service for advanced usage
func (h *SystemHandle) SystemService() system.SystemService {
    return h.systemService
}

// Context returns the system context
func (h *SystemHandle) Context() context.Context {
    return h.ctx
}

// StopSystem provides a simple shutdown function for backward compatibility
func StopSystem(handle *SystemHandle) error {
    if handle == nil {
        return fmt.Errorf("system handle is nil")
    }
    return handle.Stop()
}
```

**B. Enhanced FX Bootstrap with Context Support**
```go
// internal/infrastructure/system/fx_bootstrap.go - Enhanced with shutdown

// StartWithFxAndContext starts the system with FX and returns shutdown function
func StartWithFxAndContext(ctx context.Context, config *SystemConfig) (system.SystemService, func() error, error) {
    // Apply defaults
    if err := config.applyDefaults(); err != nil {
        return nil, nil, err
    }
    
    var systemService *DefaultSystemService
    var shutdownFn func() error
    
    app := fx.New(
        // ... existing fx configuration ...
        
        // Capture system service and shutdown function
        fx.Populate(&systemService),
        fx.Invoke(func(lc fx.Lifecycle, sys *DefaultSystemService) {
            shutdownFn = func() error {
                // Graceful shutdown sequence
                shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
                defer cancel()
                
                return sys.Stop(component.WrapContext(shutdownCtx))
            }
        }),
    )
    
    // Start the application
    startCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
    defer cancel()
    
    if err := app.Start(startCtx); err != nil {
        return nil, nil, fmt.Errorf("failed to start fx application: %w", err)
    }
    
    return systemService, shutdownFn, nil
}
```

#### **3. Container-Based Application Testing with TestKit**

**A. Application Container Setup Using TestKit**
```go
// Using the component-testkit module
import (
    "github.com/ebanfa/component-testkit/container"
    "github.com/ebanfa/component-testkit/patterns"
)

func createSkeletonContainer(ctx context.Context) (*container.AppContainer, error) {
    config := &container.AppConfig{
        DockerfilePath:   "test/integration/containers/Dockerfile.skeleton",
        BuildContext:     "../../..",
        ExposedPorts:     []string{"8080/tcp", "9090/tcp"},
        Environment: map[string]string{
            "LOG_LEVEL": "info",
            "ENV":       "test",
        },
        
        HealthEndpoint:   "/health",
        MetricsEndpoint:  "/metrics", 
        ShutdownEndpoint: "/admin/shutdown",
        
        Verification: &container.VerificationConfig{
            HealthChecks: []container.HealthCheckConfig{
                {
                    Endpoint:     "/health",
                    ExpectedCode: 200,
                    Timeout:      5 * time.Second,
                    Retries:      10,
                },
                {
                    Endpoint:     "/ready",
                    ExpectedCode: 200,
                    Timeout:      5 * time.Second,
                    Retries:      5,
                },
            },
            
            MetricsChecks: []container.MetricsConfig{
                {
                    Endpoint: "/metrics",
                    RequiredMetrics: []string{
                        "system_startup_time",
                        "system_components_total",
                        "system_plugins_loaded",
                    },
                },
            },
            
            // Log patterns are optional and non-blocking
            LogPatterns: []container.LogPattern{
                {
                    Pattern:     "System initialized",
                    Required:    false,  // Won't fail test if missing
                    Timeout:     10 * time.Second,
                    Description: "System initialization completion",
                },
                {
                    Pattern:     "System started successfully",
                    Required:    false,  // Won't fail test if missing
                    Timeout:     15 * time.Second,
                    Description: "System startup completion",
                },
            },
            LogRequired: false,  // Log verification failures won't fail tests
            
            StartupTimeout:  60 * time.Second,
            ShutdownTimeout: 30 * time.Second,
            HealthTimeout:   30 * time.Second,
        },
    }
    
    return container.CreateAppContainer(ctx, config)
}
```

#### **4. Production-Like Application Structure with Shutdown API**

**A. Skeleton Application Entry Point with Shutdown Support**
```go
// cmd/skeleton-app/main.go
func main() {
    // Production-like configuration loading
    config := loadConfiguration()
    
    // Setup signal handling for graceful shutdown
    ctx, cancel := signal.NotifyContext(context.Background(), 
        syscall.SIGINT, syscall.SIGTERM)
    defer cancel()
    
    // Setup health and metrics endpoints
    healthServer := setupHealthServer(config.HealthPort)
    metricsServer := setupMetricsServer(config.MetricsPort)
    
    // Start the system using the enhanced public API
    systemOptions := []system.Option{
        system.WithConfig(config.SystemConfig),
        system.WithPlugins(loadPlugins(config.PluginPaths)),
    }
    
    // Start system and get handle for shutdown
    systemHandle, err := system.StartSystem(systemOptions...)
    if err != nil {
        log.Fatalf("Failed to start system: %v", err)
    }
    
    // Setup shutdown endpoint for API-based shutdown
    setupShutdownEndpoint(healthServer, systemHandle)
    
    log.Info("System started successfully")
    
    // Wait for shutdown signal or system error
    select {
    case <-ctx.Done():
        log.Info("Shutdown signal received")
    case <-systemHandle.Context().Done():
        log.Info("System context cancelled")
    }
    
    // Graceful shutdown sequence
    log.Info("Starting graceful shutdown...")
    
    // 1. Stop accepting new requests
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer shutdownCancel()
    
    // 2. Stop HTTP servers
    if err := healthServer.Shutdown(shutdownCtx); err != nil {
        log.Error("Health server shutdown error: %v", err)
    }
    if err := metricsServer.Shutdown(shutdownCtx); err != nil {
        log.Error("Metrics server shutdown error: %v", err)
    }
    
    // 3. Stop the system
    if err := systemHandle.Stop(); err != nil {
        log.Error("System shutdown error: %v", err)
        os.Exit(1)
    }
    
    log.Info("Application shutdown complete")
}

// setupShutdownEndpoint adds a shutdown endpoint for API-based testing
func setupShutdownEndpoint(server *http.Server, systemHandle *system.SystemHandle) {
    if mux, ok := server.Handler.(*http.ServeMux); ok {
        mux.HandleFunc("/admin/shutdown", func(w http.ResponseWriter, r *http.Request) {
            if r.Method != http.MethodPost {
                w.WriteHeader(http.StatusMethodNotAllowed)
                return
            }
            
            w.WriteHeader(http.StatusAccepted)
            w.Write([]byte(`{"status":"shutdown_initiated"}`))
            
            // Trigger shutdown in a goroutine
            go func() {
                time.Sleep(100 * time.Millisecond) // Allow response to be sent
                if err := systemHandle.Stop(); err != nil {
                    log.Error("API shutdown error: %v", err)
                }
            }()
        })
    }
}
```

### 🔄 **Graceful Shutdown Testing with Enhanced Verification**

#### **1. Multiple Shutdown Mechanisms Using TestKit**

**A. Signal-Based Shutdown**
```go
func testSignalShutdown(t *testing.T, container *container.AppContainer) {
    // Use the testkit verification engine
    verifier := container.NewVerificationEngine(container.Config.Verification)
    
    // Send SIGTERM to container
    err := container.Container.SendSignal(ctx, syscall.SIGTERM)
    require.NoError(t, err)
    
    // Use testkit's shutdown verification (health-based, logs optional)
    err = verifier.VerifyShutdown(container)
    require.NoError(t, err)
}
```

**B. API-Based Shutdown**
```go
func testAPIShutdown(t *testing.T, container *container.AppContainer) {
    // Use testkit verification engine
    verifier := container.NewVerificationEngine(container.Config.Verification)
    
    // Call shutdown endpoint
    resp, err := http.Post(container.ShutdownURL, "application/json", nil)
    require.NoError(t, err)
    require.Equal(t, http.StatusAccepted, resp.StatusCode)
    
    // Verify shutdown using testkit (health-based verification)
    err = verifier.VerifyShutdown(container)
    require.NoError(t, err)
}
```

**C. File-Based Shutdown**
```go
func testFileShutdown(t *testing.T, container *container.AppContainer) {
    // Use testkit verification engine
    verifier := container.NewVerificationEngine(container.Config.Verification)
    
    // Create shutdown file in container
    err := container.Container.CopyToContainer(ctx, 
        []byte("shutdown"), "/tmp/shutdown.signal", 644)
    require.NoError(t, err)
    
    // Application should detect file and shutdown
    err = verifier.VerifyShutdown(container)
    require.NoError(t, err)
}
```

#### **2. Enhanced Shutdown Verification with Optional Logs**

**A. Health-First Verification Strategy**
```go
// verification.go in component-testkit
func (ve *VerificationEngine) VerifyShutdown(container *AppContainer) error {
    startTime := time.Now()
    
    // 1. PRIMARY: Health endpoint verification (required)
    if err := ve.waitForHealthUnavailable(container); err != nil {
        return fmt.Errorf("health endpoint verification failed: %w", err)
    }
    
    // 2. PRIMARY: Container exit verification (required)
    exitCode, err := container.Container.Wait(context.Background())
    if err != nil {
        return fmt.Errorf("container wait failed: %w", err)
    }
    
    if exitCode != 0 {
        return fmt.Errorf("unexpected exit code: got %d, want 0", exitCode)
    }
    
    // 3. SECONDARY: Shutdown timing verification (required)
    shutdownTime := time.Since(startTime)
    if shutdownTime > ve.config.ShutdownTimeout {
        return fmt.Errorf("shutdown took too long: %v > %v", shutdownTime, ve.config.ShutdownTimeout)
    }
    
    // 4. OPTIONAL: Log pattern verification (non-blocking)
    if len(ve.config.LogPatterns) > 0 {
        if err := ve.verifyShutdownLogs(container); err != nil {
            if ve.config.LogRequired {
                return fmt.Errorf("log verification failed: %w", err)
            } else {
                ve.logger.Warn("Log verification failed (non-critical): %v", err)
            }
        }
    }
    
    return nil
}

// waitForHealthUnavailable waits for health endpoint to become unavailable
func (ve *VerificationEngine) waitForHealthUnavailable(container *AppContainer) error {
    timeout := time.After(ve.config.ShutdownTimeout)
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-timeout:
            return fmt.Errorf("health endpoint still available after %v", ve.config.ShutdownTimeout)
        case <-ticker.C:
            resp, err := http.Get(container.HealthURL)
            if err != nil || resp.StatusCode >= 500 {
                // Health endpoint is unavailable - shutdown successful
                return nil
            }
            resp.Body.Close()
        }
    }
}
```

### 🏗️ **Test Structure and Organization with TestKit**

#### **1. Enhanced Test File Organization**
```
skeleton/test/integration/containers/
├── container_integration_test.go    # Main test suite using testkit
├── startup_test.go                  # Startup verification tests
├── shutdown_test.go                 # Shutdown verification tests
├── health_test.go                   # Health check tests
├── testdata/
│   ├── Dockerfile.skeleton          # Application container
│   ├── configs/
│   │   ├── test-production.yaml     # Production-like config
│   │   └── test-minimal.yaml        # Minimal config
│   └── plugins/                     # Test plugins
└── helpers/
    ├── testkit_helpers.go           # TestKit integration helpers
    └── skeleton_helpers.go          # Skeleton-specific helpers

# Component TestKit Module (separate repository/module)
component-testkit/
├── container/                       # Reusable container management
├── patterns/                        # Common verification patterns
├── fixtures/                        # Test fixtures and utilities
└── examples/                        # Usage examples
```

#### **2. Test Suite Structure Using TestKit**
```go
func TestSkeletonContainerIntegration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping container integration tests in short mode")
    }
    
    ctx := context.Background()
    
    t.Run("ProductionStartup", func(t *testing.T) {
        container := setupSkeletonContainerWithTestKit(t, ctx, "test-production.yaml")
        defer cleanupContainer(t, container)
        
        // Use testkit verification engine
        verifier := container.NewVerificationEngine(container.Config.Verification)
        err := verifier.VerifyStartup(container)
        require.NoError(t, err)
    })
    
    t.Run("GracefulShutdown", func(t *testing.T) {
        t.Run("SignalShutdown", func(t *testing.T) {
            container := setupSkeletonContainerWithTestKit(t, ctx, "test-production.yaml")
            defer cleanupContainer(t, container)
            
            testSignalShutdown(t, container)
        })
        
        t.Run("APIShutdown", func(t *testing.T) {
            container := setupSkeletonContainerWithTestKit(t, ctx, "test-production.yaml")
            defer cleanupContainer(t, container)
            
            testAPIShutdown(t, container)
        })
    })
    
    t.Run("FailureScenarios", func(t *testing.T) {
        // Test startup failures, resource constraints, etc.
        // Using testkit patterns for common failure scenarios
    })
}

// setupSkeletonContainerWithTestKit creates a container using the testkit
func setupSkeletonContainerWithTestKit(t *testing.T, ctx context.Context, configFile string) *container.AppContainer {
    config := &container.AppConfig{
        // ... configuration using testkit patterns
        Verification: patterns.DefaultWebAppVerification(), // Reusable pattern
    }
    
    // Override log requirements for flexible testing
    config.Verification.LogRequired = false
    config.Verification.LogPatterns = patterns.SkeletonAppLogPatterns() // Skeleton-specific patterns
    
    container, err := container.CreateAppContainer(ctx, config)
    require.NoError(t, err)
    
    return container
}
```

### 📊 **Success Criteria with Enhanced Verification**

#### **1. Startup Verification (Health-First Approach)**
- ✅ Container starts within 60 seconds
- ✅ **PRIMARY**: Health endpoints respond correctly (required)
- ✅ **PRIMARY**: Metrics endpoints expose expected metrics (required)
- ✅ **SECONDARY**: Expected log messages appear (optional, non-blocking)
- ✅ System components initialize properly (verified via health checks)
- ✅ Plugins load successfully (verified via metrics/health)

#### **2. Shutdown Verification (Health-First Approach)**
- ✅ **PRIMARY**: Graceful shutdown completes within 30 seconds (required)
- ✅ **PRIMARY**: Health endpoints become unavailable (required)
- ✅ **PRIMARY**: Exit code is 0 for normal shutdown (required)
- ✅ **SECONDARY**: Expected shutdown log messages appear (optional, non-blocking)
- ✅ All resources are cleaned up properly (verified via container exit)

#### **3. Production Simulation with TestKit**
- ✅ Uses realistic configuration
- ✅ Handles resource constraints
- ✅ Responds to production-like signals
- ✅ Maintains observability throughout lifecycle
- ✅ **NEW**: Reusable testing patterns across applications
- ✅ **NEW**: Consistent verification strategies
- ✅ **NEW**: Optional log verification for flexibility

### 🚀 **Implementation Benefits**

#### **1. Component TestKit Module Benefits**
- **Reusability**: Same testing patterns across all applications using the component system
- **Consistency**: Standardized verification strategies
- **Maintainability**: Single place to improve testing patterns
- **Documentation**: TestKit serves as testing documentation/examples
- **Community**: Other teams can contribute testing patterns

#### **2. Enhanced Public API Benefits**
- **Production Readiness**: Applications can implement proper shutdown handling
- **Testing Flexibility**: Both in-process and container tests can use shutdown API
- **Signal Handling**: Applications can respond to OS signals properly
- **Resource Management**: Proper cleanup and resource release

#### **3. Optional Log Verification Benefits**
- **Robustness**: Tests don't fail due to log format changes
- **Flexibility**: Applications can customize logging without breaking tests
- **Primary Focus**: Health endpoints provide more reliable verification
- **Debugging**: Log patterns still available for debugging when needed

This enhanced approach provides a comprehensive, reusable, and production-ready testing strategy that addresses all the suggested improvements while maintaining the clean architecture principles of the existing codebase.
