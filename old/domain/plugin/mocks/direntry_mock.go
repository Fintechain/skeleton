package mocks

import (
	"io/fs"
)

// MockDirEntry implements the fs.DirEntry interface for testing
type MockDirEntry struct {
	// Function fields for controlling behavior
	NameFunc        func() string
	IsDirectoryFunc func() bool
	TypeFunc        func() fs.FileMode
	InfoFunc        func() (fs.FileInfo, error)
}

// Name mocks the Name method
func (m *MockDirEntry) Name() string {
	if m.NameFunc != nil {
		return m.NameFunc()
	}
	return "mock-entry"
}

// IsDir mocks the IsDir method
func (m *MockDirEntry) IsDir() bool {
	if m.IsDirectoryFunc != nil {
		return m.IsDirectoryFunc()
	}
	return false
}

// Type mocks the Type method
func (m *MockDirEntry) Type() fs.FileMode {
	if m.TypeFunc != nil {
		return m.TypeFunc()
	}
	return 0644
}

// Info mocks the Info method
func (m *MockDirEntry) Info() (fs.FileInfo, error) {
	if m.InfoFunc != nil {
		return m.InfoFunc()
	}
	// Return a mock FileInfo by default
	return &MockFileInfo{}, nil
}
