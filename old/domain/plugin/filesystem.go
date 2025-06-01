package plugin

import (
	"io/fs"
	"os"
	"path/filepath"
)

// FileSystem abstracts file system operations for better testability.
type FileSystem interface {
	// Stat returns a FileInfo describing the named file.
	Stat(name string) (os.FileInfo, error)

	// WalkDir walks the file tree rooted at root, calling fn for each file or directory.
	WalkDir(root string, fn fs.WalkDirFunc) error
}

// StandardFileSystem implements the FileSystem interface using the standard library.
type StandardFileSystem struct{}

// NewStandardFileSystem creates a new filesystem implementation that uses the standard library.
func NewStandardFileSystem() *StandardFileSystem {
	return &StandardFileSystem{}
}

// Stat returns a FileInfo describing the named file.
func (fs *StandardFileSystem) Stat(name string) (os.FileInfo, error) {
	return os.Stat(name)
}

// WalkDir walks the file tree rooted at root, calling fn for each file or directory.
func (fs *StandardFileSystem) WalkDir(root string, fn fs.WalkDirFunc) error {
	return filepath.WalkDir(root, fn)
}
