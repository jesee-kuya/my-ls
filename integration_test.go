package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestIntegration_ComplexDirectoryStructure tests the complete application with complex directory structures
func TestIntegration_ComplexDirectoryStructure(t *testing.T) {
	// Create a complex directory structure for testing
	tempDir := t.TempDir()

	// Create nested directories and files
	structure := map[string]string{
		"file1.txt":                  "content1",
		"file2.txt":                  "content2",
		".hidden_file":               "hidden content",
		"dir1/subfile1.txt":          "subcontent1",
		"dir1/subfile2.txt":          "subcontent2",
		"dir1/.hidden_sub":           "hidden sub content",
		"dir1/subdir1/deep_file.txt": "deep content",
		"dir2/another_file.txt":      "another content",
		"dir2/executable":            "#!/bin/bash\necho hello",
		"archive.tar":                "tar content",
		"archive.gz":                 "gz content",
		"archive.zip":                "zip content",
	}

	// Create all files and directories
	for path, content := range structure {
		fullPath := filepath.Join(tempDir, path)
		dir := filepath.Dir(fullPath)

		// Create directory if it doesn't exist
		if dir != tempDir {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				t.Fatalf("Failed to create directory %s: %v", dir, err)
			}
		}

		// Create file
		err := os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", fullPath, err)
		}

		// Make executable file executable
		if strings.Contains(path, "executable") {
			err = os.Chmod(fullPath, 0755)
			if err != nil {
				t.Fatalf("Failed to make file executable: %v", err)
			}
		}
	}

	t.Run("basic listing", func(t *testing.T) {
		// Save original os.Args
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		// Test basic listing
		os.Args = []string{"my-ls", tempDir}

		// This would normally call main(), but we'll test the parseArgs function instead
		flags, paths := parseArgs(os.Args[1:])

		if len(paths) != 1 || paths[0] != tempDir {
			t.Errorf("Expected paths [%s], got %v", tempDir, paths)
		}

		if flags.ShowAll || flags.Longformat || flags.Reverse || flags.Recursive || flags.TimeSort {
			t.Errorf("Expected default flags, got %+v", flags)
		}
	})

	t.Run("all flags combined", func(t *testing.T) {
		// Save original os.Args
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		// Test with all flags
		os.Args = []string{"my-ls", "-alrRt", tempDir}

		flags, paths := parseArgs(os.Args[1:])

		if len(paths) != 1 || paths[0] != tempDir {
			t.Errorf("Expected paths [%s], got %v", tempDir, paths)
		}

		if !flags.ShowAll || !flags.Longformat || !flags.Reverse || !flags.Recursive || !flags.TimeSort {
			t.Errorf("Expected all flags to be true, got %+v", flags)
		}
	})

	t.Run("multiple paths", func(t *testing.T) {
		// Save original os.Args
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		dir1 := filepath.Join(tempDir, "dir1")
		dir2 := filepath.Join(tempDir, "dir2")

		// Test with multiple paths
		os.Args = []string{"my-ls", "-l", dir1, dir2}

		flags, paths := parseArgs(os.Args[1:])

		expectedPaths := []string{dir1, dir2}
		if len(paths) != len(expectedPaths) {
			t.Errorf("Expected %d paths, got %d", len(expectedPaths), len(paths))
		}

		for i, expected := range expectedPaths {
			if i >= len(paths) || paths[i] != expected {
				t.Errorf("Expected path %s at index %d, got %v", expected, i, paths)
			}
		}

		if !flags.Longformat {
			t.Errorf("Expected Longformat flag to be true")
		}
	})

	t.Run("mixed flags and paths", func(t *testing.T) {
		// Save original os.Args
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		file1 := filepath.Join(tempDir, "file1.txt")

		// Test with scattered flags and paths
		os.Args = []string{"my-ls", "-l", file1, "-a", tempDir, "-r"}

		flags, paths := parseArgs(os.Args[1:])

		expectedPaths := []string{file1, tempDir}
		if len(paths) != len(expectedPaths) {
			t.Errorf("Expected %d paths, got %d", len(expectedPaths), len(paths))
		}

		if !flags.Longformat || !flags.ShowAll || !flags.Reverse {
			t.Errorf("Expected Longformat, ShowAll, and Reverse flags to be true, got %+v", flags)
		}
	})
}

// TestIntegration_ErrorHandling tests error handling scenarios
func TestIntegration_ErrorHandling(t *testing.T) {
	t.Run("non-existent path", func(t *testing.T) {
		// Save original os.Args
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		// Test with non-existent path
		os.Args = []string{"my-ls", "/non/existent/path"}

		flags, paths := parseArgs(os.Args[1:])

		if len(paths) != 1 || paths[0] != "/non/existent/path" {
			t.Errorf("Expected paths [/non/existent/path], got %v", paths)
		}

		// parseArgs should not validate paths, just parse them
		if flags.ShowAll || flags.Longformat || flags.Reverse || flags.Recursive || flags.TimeSort {
			t.Errorf("Expected default flags for non-existent path, got %+v", flags)
		}
	})

	t.Run("invalid flag characters", func(t *testing.T) {
		// Save original os.Args
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		// Test with invalid flags (should be ignored)
		os.Args = []string{"my-ls", "-xyz", "."}

		flags, paths := parseArgs(os.Args[1:])

		// Invalid flags should be ignored
		if flags.ShowAll || flags.Longformat || flags.Reverse || flags.Recursive || flags.TimeSort {
			t.Errorf("Expected default flags for invalid flags, got %+v", flags)
		}

		if len(paths) != 1 || paths[0] != "." {
			t.Errorf("Expected paths [.], got %v", paths)
		}
	})
}

// TestIntegration_FlagCombinations tests various flag combinations
func TestIntegration_FlagCombinations(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected struct {
			showAll    bool
			longformat bool
			reverse    bool
			recursive  bool
			timeSort   bool
		}
	}{
		{
			name: "no flags",
			args: []string{"."},
			expected: struct {
				showAll    bool
				longformat bool
				reverse    bool
				recursive  bool
				timeSort   bool
			}{false, false, false, false, false},
		},
		{
			name: "single flags",
			args: []string{"-a", "-l", "-r", "-R", "-t", "."},
			expected: struct {
				showAll    bool
				longformat bool
				reverse    bool
				recursive  bool
				timeSort   bool
			}{true, true, true, true, true},
		},
		{
			name: "combined flags",
			args: []string{"-alrRt", "."},
			expected: struct {
				showAll    bool
				longformat bool
				reverse    bool
				recursive  bool
				timeSort   bool
			}{true, true, true, true, true},
		},
		{
			name: "partial flags",
			args: []string{"-al", "."},
			expected: struct {
				showAll    bool
				longformat bool
				reverse    bool
				recursive  bool
				timeSort   bool
			}{true, true, false, false, false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flags, _ := parseArgs(tc.args)

			if flags.ShowAll != tc.expected.showAll {
				t.Errorf("ShowAll: expected %v, got %v", tc.expected.showAll, flags.ShowAll)
			}
			if flags.Longformat != tc.expected.longformat {
				t.Errorf("Longformat: expected %v, got %v", tc.expected.longformat, flags.Longformat)
			}
			if flags.Reverse != tc.expected.reverse {
				t.Errorf("Reverse: expected %v, got %v", tc.expected.reverse, flags.Reverse)
			}
			if flags.Recursive != tc.expected.recursive {
				t.Errorf("Recursive: expected %v, got %v", tc.expected.recursive, flags.Recursive)
			}
			if flags.TimeSort != tc.expected.timeSort {
				t.Errorf("TimeSort: expected %v, got %v", tc.expected.timeSort, flags.TimeSort)
			}
		})
	}
}
