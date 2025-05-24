package mocks

import (
	"io/fs"
	"os"
)

// MockFileSystem implements the plugin.FileSystem interface for testing
type MockFileSystem struct {
	// Function fields for controlling behavior
	StatFunc    func(name string) (os.FileInfo, error)
	WalkDirFunc func(root string, fn fs.WalkDirFunc) error

	// Call tracking for verification
	StatCalls    []string
	WalkDirCalls []string
}

// Stat mocks the Stat method
func (m *MockFileSystem) Stat(name string) (os.FileInfo, error) {
	m.StatCalls = append(m.StatCalls, name)
	if m.StatFunc != nil {
		return m.StatFunc(name)
	}
	return nil, os.ErrNotExist
}

// WalkDir mocks the WalkDir method
func (m *MockFileSystem) WalkDir(root string, fn fs.WalkDirFunc) error {
	m.WalkDirCalls = append(m.WalkDirCalls, root)
	if m.WalkDirFunc != nil {
		return m.WalkDirFunc(root, fn)
	}
	return nil
}
