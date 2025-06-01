# Temporary Domain Model Implementation

This package contains the redesigned domain model that solves the circular dependency issue in the current `skeleton/internal/domain` package.

## Problem Statement

The current domain design has a major flaw: components cannot access the system to access resources due to circular dependencies:

- Component (domain) needs System interface
- System interface needs Component interface (through Registry)
- Registry needs Component interface

## Solution

Break the circular dependency by making the Registry generic infrastructure and having System not depend on Component interface.

## New Architecture

### Fundamental Primitives
1. **Component** - Basic building block
2. **Operation** - Component that executes discrete work
3. **Service** - Component that provides ongoing functionality

### Infrastructure/Resources
1. **EventBus** - Event communication
2. **Storage** - Data persistence
3. **Logger** - Logging facility
4. **Registry** - Generic component registry

### System
- Provides access to all resources
- Can execute system functions like `executeOperation`
- Components get System access without circular dependencies

## Key Design Principles

1. **No Circular Dependencies**: System doesn't import Component
2. **Generic Registry**: Registry uses minimal `Identifiable` interface
3. **System Access**: Components get full System access for resources
4. **Clean Separation**: Domain primitives separate from infrastructure

## Implementation Status

- [ ] Generic Registry interface
- [ ] System interface
- [ ] Component interface with System access
- [ ] Operation interface
- [ ] Service interface
- [ ] Base implementations 