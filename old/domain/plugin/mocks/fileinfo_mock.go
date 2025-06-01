package mocks

import (
	"os"
	"time"
)

// MockFileInfo implements the os.FileInfo interface for testing
type MockFileInfo struct {
	// Function fields for controlling behavior
	NameFunc    func() string
	SizeFunc    func() int64
	ModeFunc    func() os.FileMode
	ModTimeFunc func() time.Time
	IsDirFunc   func() bool
	SysFunc     func() interface{}
}

// Name mocks the Name method
func (m *MockFileInfo) Name() string {
	if m.NameFunc != nil {
		return m.NameFunc()
	}
	return "mock-file"
}

// Size mocks the Size method
func (m *MockFileInfo) Size() int64 {
	if m.SizeFunc != nil {
		return m.SizeFunc()
	}
	return 0
}

// Mode mocks the Mode method
func (m *MockFileInfo) Mode() os.FileMode {
	if m.ModeFunc != nil {
		return m.ModeFunc()
	}
	return os.ModePerm
}

// ModTime mocks the ModTime method
func (m *MockFileInfo) ModTime() time.Time {
	if m.ModTimeFunc != nil {
		return m.ModTimeFunc()
	}
	return time.Now()
}

// IsDir mocks the IsDir method
func (m *MockFileInfo) IsDir() bool {
	if m.IsDirFunc != nil {
		return m.IsDirFunc()
	}
	return false
}

// Sys mocks the Sys method
func (m *MockFileInfo) Sys() interface{} {
	if m.SysFunc != nil {
		return m.SysFunc()
	}
	return nil
}
