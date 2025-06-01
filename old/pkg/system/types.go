package system

import (
	"github.com/fintechain/skeleton/internal/domain/storage"
	"github.com/fintechain/skeleton/internal/domain/system"
	infraSystem "github.com/fintechain/skeleton/internal/infrastructure/system"
)

// ===== SYSTEM-SPECIFIC TYPES =====

// SystemService represents the main system service interface
type SystemService = system.SystemService

// SystemOperationInput represents input for system operations
type SystemOperationInput = system.SystemOperationInput

// SystemOperationOutput represents output from system operations
type SystemOperationOutput = system.SystemOperationOutput

// Config represents system configuration (public version)
type Config struct {
	ServiceID     string                   `json:"serviceId"`
	StorageConfig storage.MultiStoreConfig `json:"storage"`
}

// InternalConfig represents the internal system configuration
type InternalConfig = infraSystem.Config
