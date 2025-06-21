# Skeleton-Testkit Implementation Plan

## 🎯 **Revised Vision**

Based on comprehensive analysis of the Skeleton Framework, the skeleton-testkit should be a **simple, focused toolkit** that helps developers test their skeleton applications with real infrastructure dependencies. 

**Core Principle**: Enhance skeleton's already excellent testing capabilities with lightweight container helpers, don't reinvent or abstract over the framework.

## 🏗️ **Architecture Overview**

### The Right Approach
```
Your Skeleton App Testing Strategy:
├── Unit Tests (No Containers)
│   └── Use skeleton's built-in runtime.NewBuilder().BuildCommand() 
├── Integration Tests (Real Dependencies)
│   └── Use testcontainers + skeleton's configuration system
└── End-to-End Tests (Full Stack)
    └── Use testcontainers for app + dependencies
```

### Skeleton-Testkit Role
```
Skeleton Framework (Excellent testing built-in)
├── runtime.NewBuilder().BuildCommand() ✅ (Perfect for unit tests)
├── runtime.NewBuilder().BuildDaemon() ✅ (Perfect for integration)  
├── Builder API dependency injection ✅ (Perfect for custom dependencies)
└── Configuration system ✅ (Perfect for test configs)

Skeleton-Testkit (Simple helpers)
├── container.NewPostgres() ✅ (Setup real database)
├── testing.TestSkeletonOperation() ✅ (Wrapper for common patterns)
├── integration.NewSuite() ✅ (Full stack setup)
└── testing.GetJSON() ✅ (HTTP test helpers)
```

## 📋 **Implementation Phases**

### Phase 1: Foundation (Week 1)
**Goal**: Create simple, honest helpers that work with skeleton's patterns

#### 1.1 Core Container Helpers
```go
// pkg/container/postgres.go
type PostgresHelper struct {
    container testcontainers.Container
    host      string
    port      int
    database  string
    username  string
    password  string
}

func NewPostgres(opts ...PostgresOption) (*PostgresHelper, error) {
    cfg := &postgresConfig{
        database: "testdb",
        username: "test", 
        password: "test",
        image:    "postgres:15-alpine",
    }
    
    for _, opt := range opts {
        opt(cfg)
    }
    
    req := testcontainers.ContainerRequest{
        Image:        cfg.image,
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_DB":       cfg.database,
            "POSTGRES_USER":     cfg.username,
            "POSTGRES_PASSWORD": cfg.password,
        },
        WaitingFor: wait.ForListeningPort("5432/tcp"),
    }
    
    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          false,
    })
    if err != nil {
        return nil, err
    }
    
    return &PostgresHelper{
        container: container,
        database:  cfg.database,
        username:  cfg.username,
        password:  cfg.password,
    }, nil
}

func (p *PostgresHelper) Start(ctx context.Context) error {
    err := p.container.Start(ctx)
    if err != nil {
        return err
    }
    
    host, err := p.container.Host(ctx)
    if err != nil {
        return err
    }
    p.host = host
    
    port, err := p.container.MappedPort(ctx, "5432")
    if err != nil {
        return err
    }
    p.port = port.Int()
    
    return nil
}

func (p *PostgresHelper) ConnectionString() string {
    return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
        p.username, p.password, p.host, p.port, p.database)
}
```

#### 1.2 Skeleton Testing Utilities
```go
// pkg/testing/skeleton.go
func TestSkeletonOperation(t *testing.T, operationID string, input map[string]interface{}, 
    plugins []plugin.Plugin, options ...runtime.Option) map[string]interface{} {
    
    builder := runtime.NewBuilder().WithPlugins(plugins...)
    
    // Apply additional configuration if provided
    if len(options) > 0 {
        // Build custom configuration from options
        for _, opt := range options {
            // Apply configuration options to builder
            // (This would need implementation based on actual option types)
        }
    }
    
    // Use Builder API
    result, err := builder.BuildCommand(operationID, input)
    require.NoError(t, err)
    return result
}

func NewTestConfig(data map[string]interface{}) config.Configuration {
    return infraConfig.NewMemoryConfigurationWithData(data)
}

func WithDatabaseConfig(url string) runtime.Option {
    return runtime.WithConfig(NewTestConfig(map[string]interface{}{
        "database.url": url,
    }))
}
```

#### 1.3 HTTP Testing Helpers
```go
// pkg/testing/http.go
func GetJSON(t *testing.T, url string) map[string]interface{} {
    resp, err := http.Get(url)
    require.NoError(t, err)
    defer resp.Body.Close()
    
    var result map[string]interface{}
    err = json.NewDecoder(resp.Body).Decode(&result)
    require.NoError(t, err)
    
    return result
}

func PostJSON(t *testing.T, url string, payload interface{}) *http.Response {
    data, err := json.Marshal(payload)
    require.NoError(t, err)
    
    resp, err := http.Post(url, "application/json", bytes.NewReader(data))
    require.NoError(t, err)
    
    return resp
}

func AssertHTTPStatus(t *testing.T, resp *http.Response, expectedStatus int) {
    assert.Equal(t, expectedStatus, resp.StatusCode)
}
```

### Phase 2: Application Containers (Week 2)
**Goal**: Simple helpers for containerized skeleton apps

#### 2.1 App Container Helper
```go
// pkg/container/app.go
type AppContainer struct {
    container testcontainers.Container
    image     string
    env       map[string]string
    ports     []string
    waitFor   wait.Strategy
}

func NewSkeletonApp(image string, opts ...AppOption) (*AppContainer, error) {
    app := &AppContainer{
        image: image,
        env:   make(map[string]string),
        ports: []string{},
    }
    
    for _, opt := range opts {
        opt(app)
    }
    
    req := testcontainers.ContainerRequest{
        Image:        app.image,
        Env:          app.env,
        ExposedPorts: app.ports,
        WaitingFor:   app.waitFor,
    }
    
    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          false,
    })
    if err != nil {
        return nil, err
    }
    
    app.container = container
    return app, nil
}

func WithEnvironment(env map[string]string) AppOption {
    return func(app *AppContainer) {
        for k, v := range env {
            app.env[k] = v
        }
    }
}

func WithPorts(ports ...string) AppOption {
    return func(app *AppContainer) {
        app.ports = append(app.ports, ports...)
    }
}

func WithWaitForHTTP(path string) AppOption {
    return func(app *AppContainer) {
        if len(app.ports) > 0 {
            port := strings.Split(app.ports[0], "/")[0]
            app.waitFor = wait.ForHTTP(path).WithPort(port)
        }
    }
}
```

#### 2.2 Usage Examples
```go
// Test with real skeleton app in container
func TestSkeletonAppEndToEnd(t *testing.T) {
    ctx := context.Background()
    
    // Start dependencies
    postgres, err := container.NewPostgres()
    require.NoError(t, err)
    err = postgres.Start(ctx)
    require.NoError(t, err)
    defer postgres.Stop(ctx)
    
    // Start skeleton app
    app, err := container.NewSkeletonApp("my-skeleton-app:latest",
        container.WithEnvironment(map[string]string{
            "DATABASE_URL": postgres.ConnectionString(),
            "LOG_LEVEL":    "debug",
        }),
        container.WithPorts("8080/tcp"),
        container.WithWaitForHTTP("/health"),
    )
    require.NoError(t, err)
    err = app.Start(ctx)
    require.NoError(t, err)
    defer app.Stop(ctx)
    
    // Test endpoints
    baseURL := app.BaseURL()
    health := testing.GetJSON(t, baseURL+"/health")
    assert.Equal(t, "healthy", health["status"])
}
```

### Phase 3: Integration Suites (Week 3)
**Goal**: Convenient full-stack test setup

#### 3.1 Integration Suite
```go
// pkg/integration/suite.go
type IntegrationSuite struct {
    appImage   string
    postgres   *container.PostgresHelper
    redis      *container.RedisHelper
    app        *container.AppContainer
    appConfig  map[string]string
    cleanup    []func()
}

func NewIntegrationSuite(appImage string) *IntegrationSuite {
    return &IntegrationSuite{
        appImage:  appImage,
        appConfig: make(map[string]string),
        cleanup:   []func(){},
    }
}

func (s *IntegrationSuite) WithPostgres(opts ...container.PostgresOption) *IntegrationSuite {
    postgres, err := container.NewPostgres(opts...)
    if err != nil {
        panic(err) // In real impl, handle better
    }
    s.postgres = postgres
    return s
}

func (s *IntegrationSuite) WithAppConfig(env map[string]string) *IntegrationSuite {
    for k, v := range env {
        s.appConfig[k] = v
    }
    return s
}

func (s *IntegrationSuite) Start(ctx context.Context) error {
    // Start dependencies first
    if s.postgres != nil {
        err := s.postgres.Start(ctx)
        if err != nil {
            return err
        }
        s.cleanup = append(s.cleanup, func() { s.postgres.Stop(ctx) })
        s.appConfig["DATABASE_URL"] = s.postgres.ConnectionString()
    }
    
    if s.redis != nil {
        err := s.redis.Start(ctx)
        if err != nil {
            return err
        }
        s.cleanup = append(s.cleanup, func() { s.redis.Stop(ctx) })
        s.appConfig["REDIS_URL"] = s.redis.ConnectionString()
    }
    
    // Start app with all config
    app, err := container.NewSkeletonApp(s.appImage,
        container.WithEnvironment(s.appConfig),
        container.WithPorts("8080/tcp"),
        container.WithWaitForHTTP("/health"),
    )
    if err != nil {
        return err
    }
    
    err = app.Start(ctx)
    if err != nil {
        return err
    }
    
    s.app = app
    s.cleanup = append(s.cleanup, func() { s.app.Stop(ctx) })
    
    return nil
}

func (s *IntegrationSuite) Stop() error {
    // Cleanup in reverse order
    for i := len(s.cleanup) - 1; i >= 0; i-- {
        s.cleanup[i]()
    }
    return nil
}
```

### Phase 4: Documentation & Examples (Week 4)
**Goal**: Clear documentation with real patterns

#### 4.1 Complete Examples
```go
// examples/calculator/test/integration_test.go
func TestCalculatorIntegration(t *testing.T) {
    // Unit test - no containers needed
    result := testing.TestSkeletonOperation(t, "add",
        map[string]interface{}{"a": 5, "b": 3},
        []plugin.Plugin{calculator.NewPlugin()},
    )
    assert.Equal(t, 8.0, result["result"])
}

// examples/userservice/test/integration_test.go  
func TestUserServiceWithDatabase(t *testing.T) {
    ctx := context.Background()
    
    postgres, err := container.NewPostgres()
    require.NoError(t, err)
    err = postgres.Start(ctx)
    require.NoError(t, err)
    defer postgres.Stop(ctx)
    
    result := testing.TestSkeletonOperation(t, "create-user",
        map[string]interface{}{
            "name": "Alice", "email": "alice@example.com",
        },
        []plugin.Plugin{database.NewPlugin(), user.NewPlugin()},
        testing.WithDatabaseConfig(postgres.ConnectionString()),
    )
    
    assert.Equal(t, "created", result["status"])
}

// examples/webapp/test/e2e_test.go
func TestWebAppEndToEnd(t *testing.T) {
    ctx := context.Background()
    
    suite := integration.NewIntegrationSuite("webapp:latest").
        WithPostgres().
        WithAppConfig(map[string]string{"LOG_LEVEL": "debug"})
    
    err := suite.Start(ctx)
    require.NoError(t, err)
    defer suite.Stop()
    
    baseURL := suite.AppURL()
    
    // Test user creation flow
    userResp := testing.PostJSON(t, baseURL+"/api/users", map[string]interface{}{
        "name": "Bob", "email": "bob@example.com",
    })
    testing.AssertHTTPStatus(t, userResp, 201)
    
    // Test user retrieval
    users := testing.GetJSON(t, baseURL+"/api/users")
    assert.Len(t, users["users"], 1)
}
```

## 🔧 **Technical Implementation**

### Project Structure
```
skeleton-testkit/
├── pkg/
│   ├── container/          # Container helpers
│   │   ├── postgres.go     # PostgreSQL helper
│   │   ├── redis.go        # Redis helper
│   │   ├── app.go          # App container helper
│   │   └── options.go      # Option patterns
│   ├── testing/            # Skeleton testing utilities
│   │   ├── skeleton.go     # Operation testing helpers
│   │   ├── http.go         # HTTP testing helpers
│   │   └── config.go       # Configuration helpers
│   └── integration/        # Integration suite
│       └── suite.go        # Full stack setup
├── examples/              # Real examples
│   ├── calculator/        # Simple operation testing
│   ├── userservice/       # Database integration
│   └── webapp/           # Full stack E2E
├── internal/             # Internal helpers
└── docs/                 # Documentation
```

### Dependencies
```go
// go.mod
module github.com/fintechain/skeleton-testkit

require (
    github.com/fintechain/skeleton v0.2.0
    github.com/testcontainers/testcontainers-go v0.26.0
    github.com/stretchr/testify v1.8.4
)
```

### Key Interfaces
```go
// Keep interfaces minimal and focused
type ContainerHelper interface {
    Start(ctx context.Context) error
    Stop(ctx context.Context) error
    ConnectionString() string
}

type TestHelper interface {
    TestOperation(t *testing.T, operationID string, input map[string]interface{}, 
        plugins []plugin.Plugin, options ...runtime.Option) map[string]interface{}
}
```

## ✅ **Success Criteria**

### Technical Metrics
- [ ] Unit tests run in < 100ms
- [ ] Integration tests complete in < 10s
- [ ] E2E tests complete in < 30s
- [ ] 95%+ test coverage
- [ ] Zero dependencies on non-existent skeleton APIs

### Developer Experience
- [ ] Getting started guide takes < 5 minutes
- [ ] Common patterns have one-liner helpers
- [ ] Works with existing skeleton projects without changes
- [ ] Clear migration path from complex approaches

### Quality Standards
- [ ] All examples build and run successfully
- [ ] Documentation matches actual implementation
- [ ] No abstractions over simple patterns
- [ ] Real integration with actual skeleton apps

## 🚀 **Migration Strategy**

### For Existing Users
```go
// Before: Complex, fake APIs
app := testkit.NewSkeletonApp("app:latest").
    WithSkeletonServices(config) // ❌ Doesn't exist

// After: Simple, real patterns  
postgres, _ := container.NewPostgres()
postgres.Start(ctx)

result := testing.TestSkeletonOperation(t, "my-op", input,
    []plugin.Plugin{myPlugin},
    testing.WithDatabaseConfig(postgres.ConnectionString()),
)
```

### For New Users
```go
// Start simple with skeleton's built-in testing
func TestMyOperation(t *testing.T) {
    result, err := runtime.NewBuilder().
        WithPlugins(calculator.NewPlugin()).
        BuildCommand("add", 
        map[string]interface{}{"a": 1, "b": 2},
        runtime.WithPlugins(calculator.NewPlugin()),
    )
    require.NoError(t, err)
    assert.Equal(t, 3.0, result["result"])
}

// Add containers when you need real dependencies
func TestWithDatabase(t *testing.T) {
    postgres, _ := container.NewPostgres()
    defer postgres.Stop(ctx)
    postgres.Start(ctx)
    
    result := testing.TestSkeletonOperation(t, "query-users", input,
        []plugin.Plugin{database.NewPlugin()},
        testing.WithDatabaseConfig(postgres.ConnectionString()),
    )
}
```

## 📚 **Documentation Plan**

### Core Documentation
1. **Getting Started** - 5-minute quickstart
2. **Testing Patterns** - Unit, integration, E2E examples
3. **Container Helpers** - PostgreSQL, Redis, etc.
4. **Skeleton Integration** - How to work with skeleton's patterns
5. **Migration Guide** - From complex to simple approaches

### Example Projects
1. **Calculator** - Simple operation testing
2. **User Service** - Database integration
3. **Web App** - Full stack testing
4. **CLI Tool** - Command mode testing

This implementation plan focuses on **enhancing** skeleton's excellent built-in testing capabilities with **simple container helpers**, rather than trying to abstract or replace them. The result will be a toolkit that skeleton developers actually want to use because it makes their testing easier, not more complex.

---

**Version**: 3.0.0  
**Philosophy**: Enhance, Don't Replace  
**Timeline**: 4 weeks  
**Breaking Changes**: Complete rewrite focusing on simplicity 