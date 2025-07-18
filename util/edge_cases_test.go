package util

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Test edge cases and error conditions to improve coverage

func TestJoinPath(t *testing.T) {
	tests := []struct {
		name     string
		dir      string
		file     string
		expected string
	}{
		{
			name:     "empty dir",
			dir:      "",
			file:     "file.txt",
			expected: "file.txt",
		},
		{
			name:     "dir with trailing slash",
			dir:      "/tmp/",
			file:     "file.txt",
			expected: "/tmp/file.txt",
		},
		{
			name:     "dir without trailing slash",
			dir:      "/tmp",
			file:     "file.txt",
			expected: "/tmp/file.txt",
		},
		{
			name:     "both empty",
			dir:      "",
			file:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := joinPath(tt.dir, tt.file)
			if result != tt.expected {
				t.Errorf("joinPath(%q, %q) = %q, want %q", tt.dir, tt.file, result, tt.expected)
			}
		})
	}
}

func TestGetAbsPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "absolute path",
			path:     "/tmp/test",
			expected: "/tmp/test",
		},
		{
			name:     "relative path",
			path:     "test",
			expected: "test",
		},
		{
			name:     "current directory",
			path:     ".",
			expected: ".",
		},
		{
			name:     "parent directory",
			path:     "..",
			expected: "..",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := getAbsPath(tt.path)
			if err != nil {
				t.Errorf("getAbsPath(%q) returned error: %v", tt.path, err)
			}
			if result != tt.expected {
				t.Errorf("getAbsPath(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestReadDirNames_EdgeCases(t *testing.T) {
	tempDir := t.TempDir()

	// Test with directory containing only hidden files
	hiddenOnlyDir := filepath.Join(tempDir, "hidden_only")
	err := os.Mkdir(hiddenOnlyDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	// Create only hidden files
	for _, name := range []string{".hidden1", ".hidden2", ".hidden3"} {
		err := os.WriteFile(filepath.Join(hiddenOnlyDir, name), []byte("content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create hidden file: %v", err)
		}
	}

	// Test without ShowAll - should return empty
	result, err := ReadDirNames(hiddenOnlyDir, Flags{ShowAll: false})
	if err != nil {
		t.Fatalf("ReadDirNames failed: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("Expected empty result for directory with only hidden files, got: %v", result)
	}

	// Test with ShowAll - should return all files plus . and ..
	result, err = ReadDirNames(hiddenOnlyDir, Flags{ShowAll: true})
	if err != nil {
		t.Fatalf("ReadDirNames failed: %v", err)
	}
	if len(result) != 5 { // 3 hidden files + . + ..
		t.Errorf("Expected 5 files with ShowAll, got: %d files", len(result))
	}
}

func TestReadDirNames_SpecialFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create files with different extensions to test archive detection
	testFiles := map[string]string{
		"archive.tar": "tar content",
		"archive.gz":  "gz content",
		"archive.tgz": "tgz content",
		"archive.zip": "zip content",
		"archive.bz2": "bz2 content",
		"archive.xz":  "xz content",
		"regular.txt": "regular content",
	}

	for name, content := range testFiles {
		err := os.WriteFile(filepath.Join(tempDir, name), []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", name, err)
		}
	}

	result, err := ReadDirNames(tempDir, Flags{ShowAll: false})
	if err != nil {
		t.Fatalf("ReadDirNames failed: %v", err)
	}

	// Check that all files are present
	if len(result) != len(testFiles) {
		t.Errorf("Expected %d files, got %d", len(testFiles), len(result))
	}

	// Check that archive files have the correct color (should contain archive color)
	archiveFiles := []string{"archive.tar", "archive.gz", "archive.tgz", "archive.zip", "archive.bz2", "archive.xz"}
	for _, file := range result {
		cleanName := StripANSI(file)
		for _, archiveFile := range archiveFiles {
			if cleanName == archiveFile {
				if !strings.Contains(file, archiveColour) {
					t.Errorf("Archive file %s should have archive color, got: %s", archiveFile, file)
				}
			}
		}
	}
}

func TestReadDirNamesLong_EdgeCases(t *testing.T) {
	tempDir := t.TempDir()

	// Create a file with a very long name
	longName := strings.Repeat("a", 100) + ".txt"
	err := os.WriteFile(filepath.Join(tempDir, longName), []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create long-named file: %v", err)
	}

	result, err := ReadDirNamesLong(tempDir, Flags{ShowAll: false})
	if err != nil {
		t.Fatalf("ReadDirNamesLong failed: %v", err)
	}

	// Should have total line + 1 file
	if len(result) != 2 {
		t.Errorf("Expected 2 lines (total + file), got %d", len(result))
	}

	// Check that the long filename is present
	found := false
	for _, line := range result {
		if strings.Contains(line, longName) {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Long filename not found in output: %v", result)
	}
}

func TestCollectDirectoriesRecursively_EdgeCases(t *testing.T) {
	tempDir := t.TempDir()

	// Create a deep directory structure
	deepDir := filepath.Join(tempDir, "level1", "level2", "level3")
	err := os.MkdirAll(deepDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create deep directory: %v", err)
	}

	// Create files at different levels
	for _, dir := range []string{tempDir, filepath.Join(tempDir, "level1"), filepath.Join(tempDir, "level1", "level2"), deepDir} {
		err := os.WriteFile(filepath.Join(dir, "file.txt"), []byte("content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file in %s: %v", dir, err)
		}
	}

	result, err := CollectDirectoriesRecursively([]string{tempDir}, Flags{ShowAll: false})
	if err != nil {
		t.Fatalf("CollectDirectoriesRecursively failed: %v", err)
	}

	// Should include the root directory and all subdirectories
	expectedDirs := []string{
		tempDir,
		filepath.Join(tempDir, "level1"),
		filepath.Join(tempDir, "level1", "level2"),
		deepDir,
	}

	if len(result) != len(expectedDirs) {
		t.Errorf("Expected %d directories, got %d: %v", len(expectedDirs), len(result), result)
	}

	// Check that all expected directories are present
	for _, expectedDir := range expectedDirs {
		found := false
		for _, actualDir := range result {
			if actualDir == expectedDir {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected directory %s not found in result: %v", expectedDir, result)
		}
	}
}

func TestGetStat_EdgeCases(t *testing.T) {
	tempDir := t.TempDir()

	// Create a regular file
	testFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testFile, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test getStat function
	stat := getStat(testFile)

	// Basic checks - the function should return a valid stat structure
	if stat.Size == 0 {
		t.Errorf("Expected non-zero size for test file")
	}

	// Test with non-existent file - should still return a stat structure (may be zero)
	nonExistentFile := filepath.Join(tempDir, "nonexistent.txt")
	stat = getStat(nonExistentFile)
	// The function handles errors internally, so we just check it doesn't panic
}
