# PluginManager Integration Prompt (Execute FIRST)

## Task Overview
You are tasked with safely integrating PluginManager as a dependency into the existing SystemService implementation. This is **Phase 1** and must be completed before any fx integration work.

## CRITICAL CONSTRAINTS

### File Modification Rules
**ONLY modify existing files. DO NOT create any new files.**

**Files you may modify (and ONLY these):**
- `internal/infrastructure/system/default_system_service.go`
- `internal/infrastructure/system/config.go`
- `internal/infrastructure/system/factory.go`
- `internal/infrastructure/system/builder.go`

### Mandatory Analysis Phase
**BEFORE making ANY changes, you MUST:**

1. **Review and understand these packages completely:**
   - `internal/domain/system/system.go` - System interfaces and types
   - `internal/infrastructure/system/` - Current SystemService implementation
   - `internal/infrastructure/event/` - Event bus patterns
   - `internal/infrastructure/config/` - Configuration patterns
   - `internal/infrastructure/context/` - Context management
   - `internal/infrastructure/storage/` - Storage patterns
   - `internal/infrastructure/logging/` - Logging patterns

2. **Locate and understand:**
   - Current SystemService constructor signature
   - Existing SystemServiceConfig structure
   - PluginManager interface definition and location
   - Current dependency injection patterns
   - Factory and builder patterns for SystemService

3. **Verify compatibility:**
   - Ensure no circular dependencies
   - Confirm PluginManager interface is stable
   - Check existing dependency patterns

## Implementation Requirements

### Step 1: Analysis and Discovery
```markdown
MANDATORY: Document your findings before proceeding:
- Current SystemService constructor signature
- Current SystemServiceConfig fields
- PluginManager interface location and definition
- Existing dependency injection pattern
```

### Step 2: Update SystemServiceConfig
Add PluginManager as a dependency to the configuration struct:
```go
type SystemServiceConfig struct {
    Registry      component.Registry
    PluginManager plugin.PluginManager  // ← ADD THIS
    EventBus      event.EventBus
    MultiStore    storage.MultiStore
    Config        *Config
}
```

### Step 3: Update SystemService Constructor
Modify the constructor to accept and store the PluginManager dependency:
- Accept PluginManager from config
- Store it as a field on SystemService
- Ensure proper initialization

### Step 4: Update Factory and Builder
Ensure Factory and Builder patterns provide PluginManager:
- Update factory to create/provide PluginManager
- Update builder to accept PluginManager configuration
- Maintain existing patterns and interfaces

### Step 5: Add PluginManager Access Method
Add a method to SystemService to access the PluginManager:
```go
func (s *DefaultSystemService) PluginManager() plugin.PluginManager
```

## Safety Guidelines

### Conservative Approach
- **Minimal Changes**: Only modify what's absolutely necessary
- **Pattern Preservation**: Follow existing dependency injection patterns exactly
- **Interface Consistency**: Use existing PluginManager interface as-is
- **No Breaking Changes**: Existing code must continue to work

### Validation Steps
After each modification:
1. Ensure code compiles
2. Verify no circular imports
3. Check all interfaces are satisfied
4. Confirm existing patterns are preserved

## Expected Outcome

After completion, the SystemService should:
1. Accept PluginManager as a constructor dependency
2. Provide access to PluginManager via getter method
3. Maintain all existing functionality
4. Follow established dependency injection patterns
5. Compile without errors

## Error Handling
- Use existing error patterns in the codebase
- Don't introduce new error types
- Follow established error wrapping conventions

## Success Criteria
- [ ] SystemService constructor accepts PluginManager
- [ ] SystemService provides PluginManager() getter method
- [ ] Factory creates PluginManager appropriately
- [ ] Builder supports PluginManager configuration
- [ ] All existing code continues to work
- [ ] No new files created
- [ ] Code compiles successfully
- [ ] No circular dependencies introduced

## Reference Materials
- Review the provided specification documents
- Follow patterns established in existing infrastructure code
- Use existing interfaces without modification
- Maintain consistency with current architecture

**REMEMBER: This is Phase 1. Do not implement any fx integration or public API. Focus solely on integrating PluginManager into the existing SystemService.**