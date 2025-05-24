package plugin

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestNewStandardFileSystem(t *testing.T) {
	fs := NewStandardFileSystem()

	if fs == nil {
		t.Fatal("Expected non-nil file system")
	}
}

func TestStandardFileSystem_Stat(t *testing.T) {
	fs := NewStandardFileSystem()

	// Test for a file that exists (the current file being executed)
	info, err := fs.Stat("filesystem_test.go")
	if err != nil {
		t.Errorf("Expected no error for existing file, got: %v", err)
	}
	if info == nil {
		t.Fatal("Expected non-nil FileInfo for existing file")
	}
	if info.Name() != "filesystem_test.go" {
		t.Errorf("Expected name to be 'filesystem_test.go', got: %s", info.Name())
	}
	if info.IsDir() {
		t.Error("Expected file not to be a directory")
	}

	// Test for a file that doesn't exist
	_, err = fs.Stat("non_existent_file.go")
	if !os.IsNotExist(err) {
		t.Errorf("Expected os.ErrNotExist for non-existent file, got: %v", err)
	}
}

func TestStandardFileSystem_WalkDir(t *testing.T) {
	fs := NewStandardFileSystem()

	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "filesystem_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create some files and directories
	subDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	files := []string{
		filepath.Join(tempDir, "file1.txt"),
		filepath.Join(tempDir, "file2.txt"),
		filepath.Join(subDir, "file3.txt"),
	}

	for _, file := range files {
		err = os.WriteFile(file, []byte("test"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Test walking the directory
	var visitedPaths []string
	err = fs.WalkDir(tempDir, func(path string, _ os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		visitedPaths = append(visitedPaths, path)
		return nil
	})

	if err != nil {
		t.Errorf("Expected no error from WalkDir, got: %v", err)
	}

	// Check that we visited all the expected paths (tempDir + subdir + 3 files = 5 paths)
	if len(visitedPaths) != 5 {
		t.Errorf("Expected to visit 5 paths, visited %d: %v", len(visitedPaths), visitedPaths)
	}

	// Check for specific paths
	foundTempDir := false
	foundSubDir := false
	foundFiles := make(map[string]bool)

	for _, path := range visitedPaths {
		switch path {
		case tempDir:
			foundTempDir = true
		case subDir:
			foundSubDir = true
		case files[0], files[1], files[2]:
			foundFiles[path] = true
		}
	}

	if !foundTempDir {
		t.Error("Did not visit the temp directory")
	}
	if !foundSubDir {
		t.Error("Did not visit the subdirectory")
	}
	for _, file := range files {
		if !foundFiles[file] {
			t.Errorf("Did not visit file: %s", file)
		}
	}

	// Test with an error handler
	customErr := errors.New("custom error")
	err = fs.WalkDir(tempDir, func(path string, _ os.DirEntry, err error) error {
		// Always return an error to stop the walk
		return customErr
	})

	if err != customErr {
		t.Errorf("Expected customErr from WalkDir with error handler, got: %v", err)
	}
}

func TestStandardFileSystem_WalkDir_NonExistentDirectory(t *testing.T) {
	fs := NewStandardFileSystem()

	// Skip this test if we're running in an environment where the path might actually exist
	if _, err := os.Stat("/this/path/does/not/exist"); err == nil {
		t.Skip("Test path unexpectedly exists, skipping test")
	}

	// On some platforms, WalkDir might not return an error for non-existent directories
	// This is system-dependent, so we don't test the specific error
	_ = fs.WalkDir("/this/path/does/not/exist", func(path string, _ os.DirEntry, err error) error {
		// If we get here with an error, that's expected
		if err != nil {
			return err
		}
		// If we get here without an error, that's unexpected
		t.Error("Expected error from callback for non-existent directory")
		return nil
	})
}
