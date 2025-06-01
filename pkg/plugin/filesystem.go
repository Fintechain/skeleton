// Package plugin provides filesystem plugin interfaces and types.
package plugin

import (
	"github.com/fintechain/skeleton/internal/domain/plugin"
)

// Re-export filesystem interfaces
type FileSystem = plugin.FileSystem

// Re-export filesystem implementations
type StandardFileSystem = plugin.StandardFileSystem

// Re-export filesystem constructor
var NewStandardFileSystem = plugin.NewStandardFileSystem
