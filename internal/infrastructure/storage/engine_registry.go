// Package storage provides infrastructure implementations for storage management.
package storage

import (
	"fmt"
	"sync"

	"github.com/fintechain/skeleton/internal/domain/storage"
)

// EngineRegistry manages storage engines and provides engine discovery.
type EngineRegistry struct {
	engines map[string]storage.Engine
	mu      sync.RWMutex
}

// NewEngineRegistry creates a new engine registry.
// This follows dependency injection by accepting no external dependencies
// since the registry is self-contained.
func NewEngineRegistry() *EngineRegistry {
	return &EngineRegistry{
		engines: make(map[string]storage.Engine),
	}
}

// Register registers a storage engine with the registry.
func (r *EngineRegistry) Register(engine storage.Engine) error {
	if engine == nil {
		return fmt.Errorf("%s: engine cannot be nil", storage.ErrInvalidConfig)
	}

	name := engine.Name()
	if name == "" {
		return fmt.Errorf("%s: engine name cannot be empty", storage.ErrInvalidConfig)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if engine already exists
	if _, exists := r.engines[name]; exists {
		return fmt.Errorf("%s: engine '%s' already exists", storage.ErrEngineExists, name)
	}

	r.engines[name] = engine
	return nil
}

// Get retrieves a storage engine by name.
func (r *EngineRegistry) Get(name string) (storage.Engine, error) {
	if name == "" {
		return nil, fmt.Errorf("%s: engine name cannot be empty", storage.ErrInvalidConfig)
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	engine, exists := r.engines[name]
	if !exists {
		return nil, fmt.Errorf("%s: engine '%s' not found", storage.ErrEngineNotFound, name)
	}

	return engine, nil
}

// Unregister removes a storage engine from the registry.
func (r *EngineRegistry) Unregister(name string) error {
	if name == "" {
		return fmt.Errorf("%s: engine name cannot be empty", storage.ErrInvalidConfig)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.engines[name]; !exists {
		return fmt.Errorf("%s: engine '%s' not found", storage.ErrEngineNotFound, name)
	}

	delete(r.engines, name)
	return nil
}

// List returns a list of all registered engine names.
func (r *EngineRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.engines))
	for name := range r.engines {
		names = append(names, name)
	}

	return names
}

// Has checks if an engine with the given name is registered.
func (r *EngineRegistry) Has(name string) bool {
	if name == "" {
		return false
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.engines[name]
	return exists
}

// Count returns the number of registered engines.
func (r *EngineRegistry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.engines)
}

// Clear removes all engines from the registry.
func (r *EngineRegistry) Clear() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.engines = make(map[string]storage.Engine)
	return nil
}
