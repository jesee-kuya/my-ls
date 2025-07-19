package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestReadDirNames_SpecialCharacters tests handling of files with special characters
func TestReadDirNames_SpecialCharacters(t *testing.T) {
	tempDir := t.TempDir()

	// Create files with special characters in names
	specialFiles := []string{
		"file with spaces.txt",
		"file-with-dashes.txt",
		"file_with_underscores.txt",
		"file.with.dots.txt",
		"file@with@symbols.txt",
		"file#with#hash.txt",
		"file$with$dollar.txt",
		"file%with%percent.txt",
		"file&with&ampersand.txt",
		"file(with)parentheses.txt",
		"file[with]brackets.txt",
		"file{with}braces.txt",
		"file+with+plus.txt",
		"file=with=equals.txt",
	}

	for _, filename := range specialFiles {
		err := os.WriteFile(filepath.Join(tempDir, filename), []byte("content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", filename, err)
		}
	}

	// Test reading directory with special character files
	names, err := ReadDirNames(tempDir, Flags{})
	if err != nil {
		t.Fatalf("ReadDirNames failed: %v", err)
	}

	// Should contain all files
	if len(names) != len(specialFiles) {
		t.Errorf("Expected %d files, got %d", len(specialFiles), len(names))
	}

	// Verify all files are present (they should be colored)
	for _, expectedFile := range specialFiles {
		found := false
		for _, name := range names {
			if strings.Contains(name, expectedFile) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("File %s not found in output", expectedFile)
		}
	}
}

// TestReadDirNamesLong_SpecialPermissions tests long format with special permissions
func TestReadDirNamesLong_SpecialPermissions(t *testing.T) {
	tempDir := t.TempDir()

	// Create files with different permissions
	testFiles := map[string]os.FileMode{
		"readonly.txt":   0444,
		"writeonly.txt":  0222,
		"executable.txt": 0755,
		"noread.txt":     0333,
		"nowrite.txt":    0555,
		"noexec.txt":     0666,
	}

	for filename, mode := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		err := os.WriteFile(filePath, []byte("content"), mode)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", filename, err)
		}
	}

	// Test long format listing
	lines, err := ReadDirNamesLong(tempDir, Flags{})
	if err != nil {
		t.Fatalf("ReadDirNamesLong failed: %v", err)
	}

	// Should have total line plus all files
	if len(lines) < len(testFiles)+1 {
		t.Errorf("Expected at least %d lines, got %d", len(testFiles)+1, len(lines))
	}

	// First line should be total
	if !strings.HasPrefix(lines[0], "total ") {
		t.Errorf("First line should start with 'total ', got: %s", lines[0])
	}

	// Verify all files are present in long format
	for filename := range testFiles {
		found := false
		for _, line := range lines[1:] {
			if strings.Contains(line, filename) {
				found = true
				// Verify line contains permission information
				if !strings.Contains(line, "-") && !strings.Contains(line, "r") && !strings.Contains(line, "w") && !strings.Contains(line, "x") {
					t.Errorf("Line for %s should contain permission information: %s", filename, line)
				}
				break
			}
		}
		if !found {
			t.Errorf("File %s not found in long format output", filename)
		}
	}
}

// TestInsertSortedByTime_EdgeCases tests edge cases in time-based sorting
func TestInsertSortedByTime_EdgeCases(t *testing.T) {
	tempDir := t.TempDir()

	// Create files with specific timestamps
	now := time.Now()
	files := []struct {
		name string
		time time.Time
	}{
		{"newest.txt", now},
		{"older.txt", now.Add(-1 * time.Hour)},
		{"oldest.txt", now.Add(-2 * time.Hour)},
		{"same_time1.txt", now.Add(-30 * time.Minute)},
		{"same_time2.txt", now.Add(-30 * time.Minute)}, // Same time as previous
	}

	for _, file := range files {
		filePath := filepath.Join(tempDir, file.name)
		err := os.WriteFile(filePath, []byte("content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", file.name, err)
		}

		// Set the modification time
		err = os.Chtimes(filePath, file.time, file.time)
		if err != nil {
			t.Fatalf("Failed to set time for %s: %v", file.name, err)
		}
	}

	// Test time-based sorting
	names, err := ReadDirNames(tempDir, Flags{TimeSort: true})
	if err != nil {
		t.Fatalf("ReadDirNames with TimeSort failed: %v", err)
	}

	// Should contain all files
	if len(names) != len(files) {
		t.Errorf("Expected %d files, got %d", len(files), len(names))
	}

	// Verify newest file comes first (time sort is newest first)
	if !strings.Contains(names[0], "newest.txt") {
		t.Errorf("Expected newest.txt to be first, got: %s", names[0])
	}
}

// TestCollectDirectoriesRecursively_DeepNesting tests very deep directory nesting
func TestCollectDirectoriesRecursively_DeepNesting(t *testing.T) {
	tempDir := t.TempDir()

	// Create a deeply nested directory structure
	currentDir := tempDir
	depth := 10
	for i := 0; i < depth; i++ {
		currentDir = filepath.Join(currentDir, fmt.Sprintf("level%d", i))
		err := os.Mkdir(currentDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory at level %d: %v", i, err)
		}

		// Create a file at each level
		testFile := filepath.Join(currentDir, fmt.Sprintf("file%d.txt", i))
		err = os.WriteFile(testFile, []byte("content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create file at level %d: %v", i, err)
		}
	}

	// Test recursive collection
	dirs, err := CollectDirectoriesRecursively([]string{tempDir}, Flags{Recursive: true})
	if err != nil {
		t.Fatalf("CollectDirectoriesRecursively failed: %v", err)
	}

	// Should include root directory plus all nested directories
	expectedDirCount := depth + 1 // root + all levels
	if len(dirs) != expectedDirCount {
		t.Errorf("Expected %d directories, got %d: %v", expectedDirCount, len(dirs), dirs)
	}

	// Verify all levels are present
	for i := 0; i < depth; i++ {
		levelName := fmt.Sprintf("level%d", i)
		found := false
		for _, dir := range dirs {
			if strings.Contains(dir, levelName) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Level %d directory not found in recursive collection", i)
		}
	}
}

// TestCompareStrings_SpecialCases tests special cases in string comparison
func TestCompareStrings_SpecialCases(t *testing.T) {
	tests := []struct {
		name     string
		a, b     string
		expected bool
	}{
		{
			name:     "numbers vs letters",
			a:        "123",
			b:        "abc",
			expected: true, // numbers come before letters
		},
		{
			name:     "mixed numbers and letters - same prefix",
			a:        "file1",
			b:        "file2",
			expected: true, // 1 < 2
		},
		{
			name:     "case sensitivity with same letters",
			a:        "File",
			b:        "file",
			expected: false, // uppercase comes after lowercase
		},
		{
			name:     "special characters vs alphanumeric",
			a:        "@file",
			b:        "afile",
			expected: false, // special chars have different ordering
		},
		{
			name:     "whitespace handling",
			a:        " file",
			b:        "file",
			expected: true, // space comes before letters
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareStrings(tt.a, tt.b)
			if result != tt.expected {
				t.Errorf("CompareStrings(%q, %q) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}
