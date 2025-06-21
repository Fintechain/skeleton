# Skeleton-Testkit Specification

## 1. Overview

The **skeleton-testkit** is a simple integration testing toolkit for Skeleton Framework applications. It provides lightweight helpers for testing skeleton apps with real infrastructure dependencies using Docker containers.

### 1.1 Purpose
- **Simple Integration Testing**: Test real skeleton applications with real dependencies
- **Container Helpers**: Lightweight wrappers around testcontainers-go
- **Skeleton-Aware Utilities**: Helpers that understand skeleton framework patterns
- **Developer Experience**: Minimal setup for common testing scenarios

### 1.2 Design Philosophy
- ✅ **Simple over Complex**: Lightweight helpers, not heavy abstractions
- ✅ **Real over Fake**: Test with real databases, real skeleton apps, real endpoints
- ✅ **Framework-Aware**: Understand skeleton's runtime modes and configuration patterns
- ✅ **Standard Tools**: Build on testcontainers-go, not reinvent container orchestration

### 1.3 What It Does NOT Do
- ❌ Abstract over skeleton framework APIs (skeleton is already simple)
- ❌ Provide fake skeleton services (test with real skeleton components)
- ❌ Control skeleton runtime modes externally (apps control their own modes)
- ❌ Create complex container orchestration (testcontainers-go is sufficient)

### 1.4 Architecture
```
┌─────────────────────────────────────────────────────────┐
│              Your Skeleton Application                  │
│       (Real skeleton app with real plugins)            │
│        runtime.NewBuilder().BuildDaemon()               │
│        runtime.NewBuilder().BuildCommand()              │
└───────────────────────────┬─────────────────────────────┘
                            │ tested with
┌───────────────────────────▼─────────────────────────────┐
│                  Skeleton-Testkit                       │
│            (Simple Testing Helpers)                     │
├─────────────────┬─────────────────┬─────────────────────┤
│   Container     │   Skeleton      │    HTTP             │
│   Helpers       │   Utilities     │   Helpers           │
│   (testcontainers) │ (config, testing) │ (endpoint testing) │
└─────────────────┴─────────────────┴─────────────────────┘
```

## 2. Core Interfaces

### 2.1 Container Helpers
```go
// Simple container setup helpers
type ContainerHelper interface {
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    ConnectionString() string
    Host() string
    Port() int
}

type PostgresHelper struct {
    container testcontainers.Container
    host      string
    port      int
    database  string
    username  string
    password  string
}

type RedisHelper struct {
    container testcontainers.Container
    host      string
    port      int
}
```

### 2.2 Skeleton Testing Utilities
```go
// Skeleton-specific testing helpers
type SkeletonTester struct {
    baseURL string
    client  *http.Client
}

type SkeletonConfig struct {
    Data map[string]interface{}
}

// Helper for creating test configurations
func NewTestConfig(data map[string]interface{}) config.Configuration

// Helper for testing skeleton operations
func TestOperation(t *testing.T, operationID string, input map[string]interface{}, 
    plugins []plugin.Plugin, options ...runtime.Option) (map[string]interface{}, error)

// Helper for testing skeleton daemon endpoints  
func TestDaemonEndpoint(t *testing.T, baseURL, path string, expectedStatus int) *http.Response
```

### 2.3 Application Container (Simple)
```go
type AppContainer struct {
    container testcontainers.Container
    image     string
    env       map[string]string
    ports     []string
}

func (a *AppContainer) Start(ctx context.Context) error
func (a *AppContainer) Stop(ctx context.Context) error
func (a *AppContainer) Host() string
func (a *AppContainer) Port(containerPort string) (string, error)
func (a *AppContainer) BaseURL() string
func (a *AppContainer) Logs(ctx context.Context) (string, error)
```

## 3. Public API

### 3.1 Container Package (pkg/container)
```go
// Infrastructure containers
func NewPostgres(opts ...PostgresOption) (*PostgresHelper, error)
func NewRedis(opts ...RedisOption) (*RedisHelper, error)
func NewMongoDB(opts ...MongoOption) (*MongoHelper, error)

// Application containers
func NewSkeletonApp(image string, opts ...AppOption) (*AppContainer, error)

// PostgresOption functions
func WithDatabase(db string) PostgresOption
func WithCredentials(user, password string) PostgresOption
func WithInitScript(script string) PostgresOption

// AppOption functions  
func WithEnvironment(env map[string]string) AppOption
func WithPorts(ports ...string) AppOption
func WithWaitForHTTP(path string) AppOption
func WithWaitForLog(message string) AppOption
```

### 3.2 Testing Package (pkg/testing)
```go
// Skeleton-specific test helpers
func TestSkeletonOperation(t *testing.T, operationID string, input map[string]interface{}, 
    plugins []plugin.Plugin, options ...runtime.Option) map[string]interface{}

func TestSkeletonDaemon(t *testing.T, plugins []plugin.Plugin, 
    testFunc func(baseURL string), options ...runtime.Option)

// Configuration helpers
func NewTestConfig(data map[string]interface{}) config.Configuration
func WithDatabaseURL(url string) ConfigOption
func WithLogLevel(level string) ConfigOption
func WithCustomConfig(key string, value interface{}) ConfigOption

// HTTP testing helpers
func GetJSON(t *testing.T, url string) map[string]interface{}
func PostJSON(t *testing.T, url string, payload interface{}) *http.Response
func AssertHTTPStatus(t *testing.T, resp *http.Response, expectedStatus int)
func AssertJSONResponse(t *testing.T, resp *http.Response, expected map[string]interface{})
```

### 3.3 Integration Package (pkg/integration)
```go
// Full integration test setup
type IntegrationSuite struct {
    postgres *PostgresHelper
    redis    *RedisHelper
    app      *AppContainer
    cleanup  []func()
}

func NewIntegrationSuite(appImage string) *IntegrationSuite
func (s *IntegrationSuite) WithPostgres(opts ...PostgresOption) *IntegrationSuite
func (s *IntegrationSuite) WithRedis(opts ...RedisOption) *IntegrationSuite
func (s *IntegrationSuite) WithAppConfig(env map[string]string) *IntegrationSuite
func (s *IntegrationSuite) Start(ctx context.Context) error
func (s *IntegrationSuite) Stop() error
func (s *IntegrationSuite) AppURL() string
func (s *IntegrationSuite) DatabaseURL() string
```

## 4. Usage Patterns

### 4.1 Unit Testing (No Containers)
```go
func TestCalculatorOperation(t *testing.T) {
    // ✅ Simple skeleton operation testing
    result := testing.TestSkeletonOperation(t, "add", 
        map[string]interface{}{"a": 5, "b": 3},
        []plugin.Plugin{calculator.NewPlugin()},
    )
    
    assert.Equal(t, 8.0, result["result"])
}
```

### 4.2 Integration Testing with Real Database
```go
func TestUserServiceWithDatabase(t *testing.T) {
    ctx := context.Background()
    
    // Start real PostgreSQL
    postgres, err := container.NewPostgres(
        container.WithDatabase("testdb"),
        container.WithCredentials("test", "test"),
    )
    require.NoError(t, err)
    err = postgres.Start(ctx)
    require.NoError(t, err)
    defer postgres.Stop(ctx)
    
    // Test skeleton operation with real database
    result := testing.TestSkeletonOperation(t, "create-user",
        map[string]interface{}{
            "name":  "John Doe", 
            "email": "john@example.com",
        },
        []plugin.Plugin{
            database.NewPlugin(),
            user.NewPlugin(),
        },
        runtime.WithConfig(testing.NewTestConfig(map[string]interface{}{
            "database.url": postgres.ConnectionString(),
        })),
    )
    
    assert.Equal(t, "created", result["status"])
    assert.NotEmpty(t, result["user_id"])
}
```

### 4.3 End-to-End Testing with Real Skeleton App
```go
func TestFullSkeletonApplication(t *testing.T) {
    ctx := context.Background()
    
    // Setup infrastructure
    postgres, err := container.NewPostgres()
    require.NoError(t, err)
    err = postgres.Start(ctx)
    require.NoError(t, err)
    defer postgres.Stop(ctx)
    
    // Start your skeleton app as a container
    app, err := container.NewSkeletonApp("my-skeleton-app:latest",
        container.WithEnvironment(map[string]string{
            "DATABASE_URL": postgres.ConnectionString(),
            "LOG_LEVEL":    "debug",
        }),
        container.WithPorts("8080"),
        container.WithWaitForHTTP("/health"),
    )
    require.NoError(t, err)
    err = app.Start(ctx)
    require.NoError(t, err)
    defer app.Stop(ctx)
    
    // Test real HTTP endpoints
    baseURL := app.BaseURL()
    
    // Test health endpoint
    resp := testing.GetJSON(t, baseURL+"/health")
    assert.Equal(t, "healthy", resp["status"])
    
    // Test business endpoint
    userResp := testing.PostJSON(t, baseURL+"/api/users", map[string]interface{}{
        "name":  "Alice Smith",
        "email": "alice@example.com",
    })
    testing.AssertHTTPStatus(t, userResp, 201)
    
    // Verify user was created
    usersResp := testing.GetJSON(t, baseURL+"/api/users")
    users := usersResp["users"].([]interface{})
    assert.Len(t, users, 1)
}
```

### 4.4 Integration Suite Pattern
```go
func TestCompleteUserFlow(t *testing.T) {
    ctx := context.Background()
    
    // Setup complete integration environment
    suite := integration.NewIntegrationSuite("user-service:latest").
        WithPostgres(
            container.WithInitScript("testdata/schema.sql"),
        ).
        WithRedis().
        WithAppConfig(map[string]string{
            "LOG_LEVEL": "debug",
        })
    
    err := suite.Start(ctx)
    require.NoError(t, err)
    defer suite.Stop()
    
    baseURL := suite.AppURL()
    
    // Test complete user workflow
    t.Run("create user", func(t *testing.T) {
        resp := testing.PostJSON(t, baseURL+"/api/users", map[string]interface{}{
            "name": "Bob Wilson", "email": "bob@example.com",
        })
        testing.AssertHTTPStatus(t, resp, 201)
    })
    
    t.Run("get users", func(t *testing.T) {
        resp := testing.GetJSON(t, baseURL+"/api/users")
        users := resp["users"].([]interface{})
        assert.Len(t, users, 1)
    })
    
    t.Run("update user", func(t *testing.T) {
        resp := testing.PostJSON(t, baseURL+"/api/users/1", map[string]interface{}{
            "name": "Robert Wilson",
        })
        testing.AssertHTTPStatus(t, resp, 200)
    })
}
```

## 5. Configuration

### 5.1 Environment Variables
| Variable | Description | Default |
|----------|-------------|---------|
| `TESTKIT_LOG_LEVEL` | Logging level for testkit | `info` |
| `TESTKIT_TIMEOUT` | Default container timeout | `30s` |
| `TESTKIT_KEEP_CONTAINERS` | Keep containers after tests (debugging) | `false` |

### 5.2 Skeleton App Structure
Your skeleton applications should be structured to work well with container testing:

```go
// cmd/myapp/main.go - Your skeleton application
package main

import (
    "os"
    "github.com/fintechain/skeleton/pkg/runtime"
    "github.com/mycompany/myapp/internal/plugins"
)

func main() {
    // Determine mode from command line args
    if len(os.Args) > 1 && os.Args[1] == "command" {
        // Command mode
        if len(os.Args) < 3 {
            panic("command mode requires operation name")
        }
        
        operation := os.Args[2]
        result, err := runtime.NewBuilder().
            WithPlugins(plugins.AllPlugins()...).
            BuildCommand(operation, getInputFromArgsOrEnv()) // Your input parsing logic
        if err != nil {
            panic(err)
        }
        
        // Output result and exit
        fmt.Printf("Result: %v\n", result)
        
    } else {
        // Daemon mode (default)
        err := runtime.NewBuilder().
            WithPlugins(plugins.AllPlugins()...).
            BuildDaemon()
        if err != nil {
            panic(err)
        }
    }
        if err != nil {
            panic(err)
        }
    }
}
```

### 5.3 Dockerfile Example
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp cmd/myapp/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/myapp .

# Default to daemon mode, but allow command mode
ENTRYPOINT ["./myapp"]
CMD []
```

## 6. Best Practices

### 6.1 Test Organization
```go
// Structure your tests clearly
func TestUnit_CalculatorOperations(t *testing.T) {
    // Fast unit tests with no containers
}

func TestIntegration_UserServiceWithDatabase(t *testing.T) {
    // Integration tests with real database
}

func TestE2E_CompleteUserWorkflow(t *testing.T) {
    // End-to-end tests with full application
}
```

### 6.2 Container Lifecycle Management
```go
func TestWithProperCleanup(t *testing.T) {
    ctx := context.Background()
    
    // Always use defer for cleanup
    postgres, err := container.NewPostgres()
    require.NoError(t, err)
    defer postgres.Stop(ctx) // Always clean up
    
    err = postgres.Start(ctx)
    require.NoError(t, err)
    
    // Your test logic here
}
```

### 6.3 Configuration Management
```go
func TestWithConfiguration(t *testing.T) {
    // Use skeleton's configuration system properly
    testConfig := testing.NewTestConfig(map[string]interface{}{
        "database.url":      postgres.ConnectionString(),
        "cache.url":         redis.ConnectionString(),
        "log.level":         "debug",
        "app.environment":   "test",
    })
    
    result := testing.TestSkeletonOperation(t, "my-operation", input,
        plugins,
        runtime.WithConfig(testConfig),
    )
}
```

## 7. Implementation Requirements

### 7.1 Dependencies
- Go 1.21+
- testcontainers-go
- Skeleton Framework
- Docker (for integration tests)

### 7.2 Performance Goals
- Unit tests: < 100ms per test
- Integration tests: < 10s per test (including container startup)
- End-to-end tests: < 30s per test

### 7.3 Reliability
- Automatic container cleanup on test failure
- Proper error context and debugging information
- Support for parallel test execution
- Configurable timeouts

## 8. Migration from Complex Version

### 8.1 Before (Complex)
```go
// Old complex approach
app := testkit.NewSkeletonApp("my-app:latest").
    WithSkeletonServices([]SkeletonServiceConfig{...}).  // ❌ Doesn't exist
    WithSkeletonOperations([]SkeletonOperationConfig{...}) // ❌ Doesn't exist

verifier := verification.NewSkeletonVerifier(app)
err = verifier.VerifyPluginLoaded("my-plugin") // ❌ Unnecessary
```

### 8.2 After (Simple)
```go
// New simple approach
app, err := container.NewSkeletonApp("my-app:latest",
    container.WithEnvironment(map[string]string{
        "DATABASE_URL": postgres.ConnectionString(),
    }),
)

// Test real endpoints directly
resp := testing.GetJSON(t, app.BaseURL()+"/health")
assert.Equal(t, "healthy", resp["status"])
```

---

**Version**: 3.0.0  
**Philosophy**: Simple, Real, Framework-Aware  
**Last Updated**: 2024-12-27 