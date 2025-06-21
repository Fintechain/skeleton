# Skeleton-Testkit Development Context

You are working on the **skeleton-testkit** - a companion testing framework for applications built with the Fintechain Skeleton framework. The testkit provides container-based integration testing capabilities specifically designed for skeleton-based applications.

## 📚 **Required Reading - Skeleton Framework Context**

**Step 1: Review skeleton framework documentation**
- ../skeleton/docs/PLUGIN_DEVELOPMENT_GUIDE.md
- ../skeleton/docs/RUNTIME_DEVELOPMENT_GUIDE.md  
- ../skeleton/docs/SERVICE_OPERATIONS_DEVELOPMENT_GUIDE.md
- ../skeleton/README.md (framework overview and dual runtime modes)
- ../skeleton/CHANGELOG.md (v0.3.0 breaking changes)

**Step 2: Understand testkit architecture**
- docs/architecture/skeleton-testkit-specification.md
- docs/architecture/skeleton-testkit.md
- pkg/testkit/testkit.go (main API)
- test/integration/ (usage examples)

## 🎯 **Project Context**

### **What skeleton-testkit IS:**
- **Integration testing framework** for skeleton-based applications
- **Container management** using Docker/Testcontainers
- **Database testing support** (PostgreSQL, Redis containers)
- **Health checking and verification** utilities
- **Skeleton-aware testing** (understands plugins, services, operations)

### **What skeleton-testkit is NOT:**
- ❌ Not a replacement for the skeleton framework
- ❌ Not for unit testing (skeleton has test/unit/README.md for that)
- ❌ Not a general-purpose testing framework
- ❌ Not ready for promotion yet (still in development)

### **Relationship to skeleton framework:**
```
skeleton-testkit (v0.1.0 - In Development)
    │
    │ Tests applications built with
    ▼
skeleton (v0.3.0 - Stable)
    │
    │ Builds applications using
    ▼
Your Application
```

## 🏗️ **Current Architecture**

### **Package Structure:**
```
pkg/
├── testkit/        # Main API entry point
├── container/      # Container management (app, postgres, redis)
├── health/         # Health checking utilities
└── verification/   # System verification tools

internal/
├── domain/         # Domain interfaces
├── infrastructure/ # Docker/Testcontainers implementations
└── ...
```

### **Key APIs:**
```go
// Main testkit API
app := testkit.NewSkeletonApp("my-app:latest")
postgres := testkit.NewPostgresContainer()
redis := testkit.NewRedisContainer()

// Container management
app.WithSkeletonConfig(&container.SkeletonConfig{...})
app.Start(ctx)
app.Stop(ctx)

// Health and verification
health.CheckContainer(app)
verification.VerifySystemState(app)
```

## 🎯 **Development Guidelines**

### **✅ Focus Areas:**
- **Container-based testing** for skeleton applications
- **Integration test patterns** and utilities
- **Database testing support** (Postgres, Redis, etc.)
- **Skeleton framework integration** (plugins, services, operations)
- **Health checking and verification** capabilities
- **CI/CD pipeline integration** support

### **🔄 Current Status:**
- **Version**: v0.1.0 (early development)
- **Stability**: Not ready for production use
- **Goal**: Mature to production-ready state
- **Priority**: Internal development before external promotion

### **📋 Version Strategy:**
- **Future goal**: Synchronize versions with skeleton framework
- **Current**: Independent versioning during development
- **Target**: skeleton v0.3.x ↔ skeleton-testkit v0.3.x

## 🧪 **Testing Philosophy**

### **Integration Testing Focus:**
- **End-to-end testing** of complete skeleton applications
- **Database integration** with real database containers
- **Multi-service testing** (app + dependencies)
- **Performance testing** under realistic conditions
- **CI/CD integration** for automated testing

### **Skeleton-Aware Testing:**
- **Plugin configuration** testing
- **Service lifecycle** verification
- **Operation execution** testing
- **Component interaction** validation
- **Runtime mode testing** (daemon vs command)

## 🎯 **Current Task Context**

[Include specific task details here]

## 🔧 **Development Environment**

- **Go Version**: 1.21+
- **Dependencies**: Docker, Testcontainers
- **Testing**: Integration tests in test/integration/
- **Documentation**: docs/ directory structure
- **Examples**: examples/ directory

## 📝 **Key Considerations**

1. **Skeleton Framework Compatibility**: Ensure testkit works with skeleton v0.3.0+ API
2. **Container Management**: Use Docker/Testcontainers for realistic testing environments
3. **Documentation**: Keep docs updated as testkit matures
4. **Version Sync**: Plan for future version synchronization with skeleton
5. **Testing Patterns**: Establish best practices for skeleton app testing

## 🔄 **Skeleton Framework API Reference (v0.3.0)**

### **Runtime Package (Builder API):**
```go
import "github.com/fintechain/skeleton/pkg/runtime"

// Daemon mode - long-running services
runtime.NewBuilder().WithPlugins(plugins...).BuildDaemon()

// Command mode - execute and exit
result, err := runtime.NewBuilder().WithPlugins(plugins...).
    BuildCommand("operation-id", input)

// Custom dependencies
runtime.NewBuilder().
    WithPlugins(plugins...).
    WithConfig(customConfig).
    WithLogger(customLogger).
    BuildDaemon()
```

### **Component System:**
```go
// Component types
type Component interface {
    ID() ComponentID
    Name() string
    Initialize(ctx context.Context, system System) error
    Dispose() error
}

// Three types: Services, Operations, Components
```

### **Plugin Architecture:**
```go
type MyPlugin struct {
    *component.BaseService
    runtime runtime.RuntimeEnvironment
}

func (p *MyPlugin) Initialize(ctx context.Context, system component.System) error {
    p.runtime = system.(runtime.RuntimeEnvironment)
    registry := system.Registry()
    return registry.Register(myService, myOperation)
}
```

## ⚠️ **Breaking Changes (v0.3.0)**

- **pkg/fx package completely removed** - replaced with Builder API
- **FX dependency injection removed** - now uses simple Builder pattern
- **All FX functions removed**: No more `StartDaemon()` or `ExecuteCommand()` legacy functions
- **Only Builder API available**: Use `runtime.NewBuilder()` for all applications

## 🚀 **Current Builder API (Only Option)**

- **Builder Pattern**: `runtime.NewBuilder().WithPlugins().BuildDaemon()`
- **Command Pattern**: `runtime.NewBuilder().WithPlugins().BuildCommand()`
- **Custom Dependencies**: `runtime.NewBuilder().WithConfig().WithLogger().WithEventBus()`
- **Simplified Architecture**: No complex dependency injection, direct constructor pattern

## 🎯 **Usage Instructions**

When starting work on skeleton-testkit:

1. **Copy this template** into your AI conversation
2. **Add specific task details** in the "Current Task Context" section
3. **Reference the skeleton docs** as needed during development
4. **Keep the context** throughout the conversation for consistency

---

**Remember**: The testkit is specifically designed for skeleton-based applications. Always consider how skeleton's dual runtime modes (daemon/command), plugin architecture, and component system should be tested through the testkit's container-based approach.

**Version Compatibility**: This template is for skeleton v0.3.0+ with the unified runtime API. Ensure testkit development aligns with the current skeleton framework architecture and API. 