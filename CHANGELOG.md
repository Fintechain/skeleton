# Changelog

## [0.3.0] - 2024-12-19

### Major Framework Consolidation and Cleanup

#### Structural Changes
- **Removed deprecated pkg/fx package** (533 lines) - functionality merged into pkg/runtime
- **Removed old internal/fx implementation** (516 lines) - replaced with fx_core.go and fx_runtime.go
- **Consolidated runtime API** - all FX functionality now accessible via pkg/runtime
- **Removed outdated example projects** (fx-demo, traditional-runtime, individual plugins)
- **Restructured documentation** from docs/development/ to docs/ root level
- **Removed unused config implementations** (composite, env, file) - simplified to memory-based
- **Removed old runtime builder pattern** - replaced with FX-based approach

#### Documentation Overhaul
- **Complete README.md rewrite** with accurate API references and working examples
- **Fixed all pkg/fx references** to use correct pkg/runtime API throughout codebase
- **Simplified examples/README.md** with focus on complete-app example
- **Updated examples/complete-app** with proper structure (modes/, plugins/, providers/)
- **Added comprehensive pkg/runtime/README.md** documentation
- **Moved and updated plugin development guides**

#### Code Quality Improvements
- **Enhanced pkg/runtime/runtime.go** with proper FX integration (159 lines added)
- **Refactored internal/infrastructure/config/memory.go** for better maintainability
- **Added new pkg/config/ package** for clean configuration exports
- **Removed 6,614 lines of outdated/duplicate code**
- **Added comprehensive test structure** for new runtime package

#### API Consolidation
- **Unified API surface**: runtime.StartDaemon() and runtime.ExecuteCommand()
- **Removed API confusion** between fx.* and runtime.* functions
- **Maintained backward compatibility** while simplifying public interface
- **Internal FX usage preserved** but hidden from public API

#### Impact
- **Files changed**: 34 files, +606/-6,614 lines
- **Major cleanup and consolidation** effort, removing technical debt
- **Cleaner, more maintainable API surface** while maintaining all functionality

### Breaking Changes
- All pkg/fx references must be updated to pkg/runtime
- Import paths changed from "github.com/fintechain/skeleton/pkg/fx" to "github.com/fintechain/skeleton/pkg/runtime"
- Function calls changed from fx.StartDaemon() to runtime.StartDaemon()
- Function calls changed from fx.ExecuteCommand() to runtime.ExecuteCommand()

### Migration Guide
- Update imports: `"github.com/fintechain/skeleton/pkg/fx"` → `"github.com/fintechain/skeleton/pkg/runtime"`
- Update function calls: `fx.StartDaemon()` → `runtime.StartDaemon()`
- Update function calls: `fx.ExecuteCommand()` → `runtime.ExecuteCommand()`
- All functionality remains the same, only the package name changed

## [0.2.0] - 2024-12-19

### Major Changes
- Complete framework rewrite and architecture overhaul
- Breaking changes across all components
- New component lifecycle management
- Enhanced plugin system and runtime environment

## [2.0.0] - 2024-12-19

### Major Improvements

#### Package Reorganization
- **Removed monolithic `common` package**: Split utilities into focused, single-purpose packages
- **Created focused utility packages**:
  - `pkg/logging/` - Logging utilities and interfaces
  - `pkg/id/` - ID generation utilities  
  - `pkg/crypto/` - Cryptographic utilities
- **Improved package independence**: Each package is now self-contained with minimal cross-dependencies

#### Factory Pattern Implementation
- **Added factory constructors to service package**:
  - `NewServiceFactory()` - Creates service factories
  - `CreateService(config)` - Creates services from configuration
  - `NewService(name, type)` - Convenience function for simple service creation
- **Added factory constructors to operation package**:
  - `NewOperationFactory()` - Creates operation factories
  - `CreateOperation(config)` - Creates operations from configuration
  - `NewOperation(name, type)` - Convenience function for simple operation creation

#### Enhanced Documentation
- **Updated README.md** with comprehensive examples of factory usage
- **Added factory pattern to design principles**
- **Improved code examples** showing both factory and direct usage patterns
- **Enhanced package descriptions** with clear separation of concerns

### Technical Improvements

#### Service Package
- Added `ServiceFactory` interface for consistent service creation
- Added `DefaultServiceFactory` implementation with component factory integration
- Enhanced service configuration with `ServiceConfig` struct
- Improved error handling with domain-specific error codes

#### Operation Package  
- Added `OperationFactory` interface for consistent operation creation
- Added `DefaultOperationFactory` implementation with component factory integration
- Enhanced operation configuration with `OperationConfig` struct
- Improved error handling with domain-specific error codes

#### Package Structure
```
pkg/
├── system/          # System lifecycle management
├── component/       # Component management and context
├── service/         # Service lifecycle with factory creation
├── operation/       # Executable operations with factory creation
├── plugin/          # Plugin system
├── config/          # Configuration management
├── event/           # Event system
├── storage/         # Storage abstractions
├── logging/         # Logging utilities (moved from common)
├── id/              # ID generation (moved from common)
└── crypto/          # Cryptographic utilities (moved from common)
```

### Breaking Changes
- Removed `pkg/common/` package - utilities moved to focused packages
- Import paths changed for logging, ID generation, and crypto utilities
- Service and operation creation now uses factory pattern (old direct creation still supported)

### Migration Guide
- Update imports from `pkg/common/logging` to `pkg/logging`
- Update imports from `pkg/common/id` to `pkg/id`  
- Update imports from `pkg/common/crypto` to `pkg/crypto`
- Consider using new factory constructors for service and operation creation

### Benefits
- **Better maintainability**: Focused packages are easier to understand and modify
- **Improved testability**: Factory pattern enables better dependency injection
- **Enhanced extensibility**: Clear interfaces and factory patterns support customization
- **Reduced coupling**: Elimination of monolithic common package reduces dependencies
- **Consistent patterns**: Factory constructors provide standardized creation patterns 