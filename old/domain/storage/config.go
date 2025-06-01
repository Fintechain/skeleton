// Package storage provides interfaces and types for the storage system.
package storage

// Config type alias to the infrastructure layer's config type for backward compatibility
type Config map[string]interface{}

// StoreConfig provides access to store-specific configuration
type StoreConfig interface {
	// GetEngine returns the engine to use for this store
	GetEngine() string

	// GetPath returns the path for this store
	GetPath() string

	// GetOptions returns engine-specific options
	GetOptions() Config
}

// DefaultStoreConfig implements StoreConfig interface
type DefaultStoreConfig struct {
	// Engine is the store engine (memory, file, etc.)
	Engine string

	// Path is the store path/location
	Path string

	// Options contains engine-specific options
	Options Config
}

// GetEngine returns the engine to use for this store
func (c *DefaultStoreConfig) GetEngine() string {
	return c.Engine
}

// GetPath returns the path for this store
func (c *DefaultStoreConfig) GetPath() string {
	return c.Path
}

// GetOptions returns engine-specific options
func (c *DefaultStoreConfig) GetOptions() Config {
	return c.Options
}

// MultiStoreConfig defines configuration for MultiStore.
// This should be used with the infrastructure's configuration system.
type MultiStoreConfig struct {
	// RootPath is the base directory for all stores managed by this MultiStore.
	RootPath string `json:"rootPath"`

	// DefaultEngine is the engine to use when no engine is specified.
	DefaultEngine string `json:"defaultEngine"`

	// EngineConfigs contains engine-specific configuration.
	// The keys are engine names, and the values are engine-specific configs.
	EngineConfigs map[string]Config `json:"engineConfigs"`
}

// Common configuration keys that can be used across different engines.
const (
	// ConfigCacheSize is the size of the cache in bytes.
	ConfigCacheSize = "cache_size"

	// ConfigCompression enables or disables compression.
	ConfigCompression = "compression"

	// ConfigMaxVersions is the maximum number of versions to keep.
	ConfigMaxVersions = "max_versions"

	// ConfigSyncWrites enables or disables synchronized writes.
	ConfigSyncWrites = "sync_writes"

	// ConfigFileFormat specifies the serialization format (json, gob, etc.).
	ConfigFileFormat = "format"

	// ConfigReadOnly opens the store in read-only mode.
	ConfigReadOnly = "read_only"
)
