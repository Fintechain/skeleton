// Package storage provides storage interfaces and types.
package storage

import (
	"github.com/fintechain/skeleton/internal/domain/storage"
)

// Re-export config types
type Config = storage.Config
type StoreConfig = storage.StoreConfig
type DefaultStoreConfig = storage.DefaultStoreConfig
type MultiStoreConfig = storage.MultiStoreConfig

// Re-export config constants
const (
	ConfigCacheSize   = storage.ConfigCacheSize
	ConfigCompression = storage.ConfigCompression
	ConfigMaxVersions = storage.ConfigMaxVersions
	ConfigSyncWrites  = storage.ConfigSyncWrites
	ConfigFileFormat  = storage.ConfigFileFormat
	ConfigReadOnly    = storage.ConfigReadOnly
)
