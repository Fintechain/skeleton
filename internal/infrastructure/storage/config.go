// Package storage provides infrastructure implementations for storage configuration.
package storage

import (
	"fmt"

	"github.com/fintechain/skeleton/internal/domain/storage"
)

// StoreConfigImpl implements the storage.StoreConfig interface.
type StoreConfigImpl struct {
	engine  string
	path    string
	options storage.Config
}

// NewStoreConfig creates a new store configuration.
// This follows dependency injection by accepting all configuration parameters.
func NewStoreConfig(engine, path string, options storage.Config) *StoreConfigImpl {
	if options == nil {
		options = make(storage.Config)
	}

	return &StoreConfigImpl{
		engine:  engine,
		path:    path,
		options: options,
	}
}

// GetEngine returns the engine to use for this store.
func (c *StoreConfigImpl) GetEngine() string {
	return c.engine
}

// GetPath returns the path for this store.
func (c *StoreConfigImpl) GetPath() string {
	return c.path
}

// GetOptions returns engine-specific options.
func (c *StoreConfigImpl) GetOptions() storage.Config {
	return c.options
}

// Validate validates the store configuration.
func (c *StoreConfigImpl) Validate() error {
	if c.engine == "" {
		return fmt.Errorf("%s: engine cannot be empty", storage.ErrInvalidConfig)
	}

	if c.path == "" {
		return fmt.Errorf("%s: path cannot be empty", storage.ErrInvalidConfig)
	}

	return nil
}

// MultiStoreConfigImpl implements configuration for MultiStore.
type MultiStoreConfigImpl struct {
	rootPath      string
	defaultEngine string
	engineConfigs map[string]storage.Config
}

// NewMultiStoreConfig creates a new multi-store configuration.
// This follows dependency injection by accepting all configuration parameters.
func NewMultiStoreConfig(rootPath, defaultEngine string, engineConfigs map[string]storage.Config) *MultiStoreConfigImpl {
	if engineConfigs == nil {
		engineConfigs = make(map[string]storage.Config)
	}

	return &MultiStoreConfigImpl{
		rootPath:      rootPath,
		defaultEngine: defaultEngine,
		engineConfigs: engineConfigs,
	}
}

// GetRootPath returns the base directory for all stores.
func (c *MultiStoreConfigImpl) GetRootPath() string {
	return c.rootPath
}

// GetDefaultEngine returns the default engine to use.
func (c *MultiStoreConfigImpl) GetDefaultEngine() string {
	return c.defaultEngine
}

// GetEngineConfigs returns engine-specific configurations.
func (c *MultiStoreConfigImpl) GetEngineConfigs() map[string]storage.Config {
	return c.engineConfigs
}

// GetEngineConfig returns configuration for a specific engine.
func (c *MultiStoreConfigImpl) GetEngineConfig(engine string) (storage.Config, bool) {
	config, exists := c.engineConfigs[engine]
	return config, exists
}

// SetEngineConfig sets configuration for a specific engine.
func (c *MultiStoreConfigImpl) SetEngineConfig(engine string, config storage.Config) {
	if c.engineConfigs == nil {
		c.engineConfigs = make(map[string]storage.Config)
	}
	c.engineConfigs[engine] = config
}

// Validate validates the multi-store configuration.
func (c *MultiStoreConfigImpl) Validate() error {
	if c.rootPath == "" {
		return fmt.Errorf("%s: root path cannot be empty", storage.ErrInvalidConfig)
	}

	if c.defaultEngine == "" {
		return fmt.Errorf("%s: default engine cannot be empty", storage.ErrInvalidConfig)
	}

	return nil
}
