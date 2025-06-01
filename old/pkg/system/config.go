package system

import (
	"github.com/fintechain/skeleton/internal/domain/storage"
)

// DefaultConfig creates a default system configuration
func DefaultConfig() *Config {
	return &Config{
		ServiceID: "system",
		StorageConfig: storage.MultiStoreConfig{
			RootPath:      "./data",
			DefaultEngine: "memory",
		},
	}
}
