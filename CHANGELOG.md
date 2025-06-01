# Changelog

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