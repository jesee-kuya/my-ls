package util

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

func TestCompareStrings(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		expected bool
	}{
		// Basic alphabetical comparison
		{
			name:     "basic alphabetical a < b",
			a:        "a",
			b:        "b",
			expected: true,
		},
		{
			name:     "basic alphabetical b > a",
			a:        "b",
			b:        "a",
			expected: false,
		},
		{
			name:     "equal strings",
			a:        "test",
			b:        "test",
			expected: false,
		},
		// Case sensitivity tests
		{
			name:     "lowercase before uppercase",
			a:        "a",
			b:        "A",
			expected: true,
		},
		{
			name:     "uppercase after lowercase",
			a:        "A",
			b:        "a",
			expected: false,
		},
		// Special characters only
		{
			name:     "special chars only - ASCII order",
			a:        "!",
			b:        "@",
			expected: true,
		},
		{
			name:     "special chars only - reverse ASCII order",
			a:        "@",
			b:        "!",
			expected: false,
		},
		// Mixed strings
		{
			name:     "mixed strings - alphabetic part comparison",
			a:        "a1",
			b:        "b2",
			expected: true,
		},
		{
			name:     "mixed strings - same alphabetic part, different special",
			a:        "a!",
			b:        "a@",
			expected: true,
		},
		// Length differences
		{
			name:     "shorter string first",
			a:        "a",
			b:        "ab",
			expected: true,
		},
		{
			name:     "longer string second",
			a:        "ab",
			b:        "a",
			expected: false,
		},
		// Complex cases
		{
			name:     "complex mixed case",
			a:        "File1",
			b:        "file2",
			expected: true, // Significant parts: file1 vs. file2, 1 < 2
		},
		{
			name:     "numbers and letters",
			a:        "file10",
			b:        "file2",
			expected: true, // Significant parts: file10 vs. file2, 1 < 2
		},
		// Numeric vs. alphabetic
		{
			name:     "numeric before alphabetic",
			a:        "a1b",
			b:        "aab",
			expected: true, // Significant parts: a1b vs. aab, 1 < a
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

func TestInsertSorted(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		colour   string
		reset    string
		names    []string
		expected []string
	}{
		{
			name:     "insert into empty list",
			filename: "file.txt",
			colour:   "\033[0m",
			reset:    "\033[0m",
			names:    []string{},
			expected: []string{"\033[0mfile.txt\033[0m"},
		},
		{
			name:     "insert at beginning",
			filename: "a.txt",
			colour:   "\033[0m",
			reset:    "\033[0m",
			names:    []string{"\033[0mb.txt\033[0m"},
			expected: []string{"\033[0ma.txt\033[0m", "\033[0mb.txt\033[0m"},
		},
		{
			name:     "insert at end",
			filename: "z.txt",
			colour:   "\033[0m",
			reset:    "\033[0m",
			names:    []string{"\033[0ma.txt\033[0m"},
			expected: []string{"\033[0ma.txt\033[0m", "\033[0mz.txt\033[0m"},
		},
		{
			name:     "insert in middle",
			filename: "b.txt",
			colour:   "\033[0m",
			reset:    "\033[0m",
			names:    []string{"\033[0ma.txt\033[0m", "\033[0mc.txt\033[0m"},
			expected: []string{"\033[0ma.txt\033[0m", "\033[0mb.txt\033[0m", "\033[0mc.txt\033[0m"},
		},
		{
			name:     "insert dot directory",
			filename: ".",
			colour:   "\033[01;34m",
			reset:    "\033[0m",
			names:    []string{"\033[0ma.txt\033[0m"},
			expected: []string{"\033[01;34m.\033[0m", "\033[0ma.txt\033[0m"},
		},
		{
			name:     "insert dotdot directory",
			filename: "..",
			colour:   "\033[01;34m",
			reset:    "\033[0m",
			names:    []string{"\033[0ma.txt\033[0m"},
			expected: []string{"\033[01;34m..\033[0m", "\033[0ma.txt\033[0m"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InsertSorted(tt.filename, tt.colour, tt.reset, tt.names)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("InsertSorted() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInsertSortedLong(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		lines    []string
		expected []string
	}{
		{
			name:     "insert into empty list",
			line:     "-rw-r--r--  1 user group  100 Jan  1 12:00 file.txt",
			lines:    []string{},
			expected: []string{"-rw-r--r--  1 user group  100 Jan  1 12:00 file.txt"},
		},
		{
			name:     "insert at beginning",
			line:     "-rw-r--r--  1 user group  100 Jan  1 12:00 a.txt",
			lines:    []string{"-rw-r--r--  1 user group  100 Jan  1 12:00 b.txt"},
			expected: []string{"-rw-r--r--  1 user group  100 Jan  1 12:00 a.txt", "-rw-r--r--  1 user group  100 Jan  1 12:00 b.txt"},
		},
		{
			name:     "insert at end",
			line:     "-rw-r--r--  1 user group  100 Jan  1 12:00 z.txt",
			lines:    []string{"-rw-r--r--  1 user group  100 Jan  1 12:00 a.txt"},
			expected: []string{"-rw-r--r--  1 user group  100 Jan  1 12:00 a.txt", "-rw-r--r--  1 user group  100 Jan  1 12:00 z.txt"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InsertSortedLong(tt.line, tt.lines)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("InsertSortedLong() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTrimStart(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no dots",
			input:    "file.txt",
			expected: "file.txt",
		},
		{
			name:     "single leading dot",
			input:    ".hidden",
			expected: "hidden",
		},
		{
			name:     "multiple leading dots",
			input:    "...file",
			expected: "file",
		},
		{
			name:     "dots in middle",
			input:    "file.name.txt",
			expected: "file.name.txt",
		},
		{
			name:     "only dots",
			input:    "...",
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TrimStart(tt.input)
			if result != tt.expected {
				t.Errorf("TrimStart(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestInsertSortedByTime(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create test files with different modification times
	file1 := filepath.Join(tempDir, "old.txt")
	file2 := filepath.Join(tempDir, "new.txt")

	// Create files
	err := os.WriteFile(file1, []byte("old content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Wait a bit to ensure different timestamps
	time.Sleep(10 * time.Millisecond)

	err = os.WriteFile(file2, []byte("new content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name     string
		filename string
		colour   string
		reset    string
		dirPath  string
		names    []string
		expected []string
	}{
		{
			name:     "insert newer file first",
			filename: "new.txt",
			colour:   "\033[0m",
			reset:    "\033[0m",
			dirPath:  tempDir,
			names:    []string{"\033[0mold.txt\033[0m"},
			expected: []string{"\033[0mnew.txt\033[0m", "\033[0mold.txt\033[0m"},
		},
		{
			name:     "insert older file after newer",
			filename: "old.txt",
			colour:   "\033[0m",
			reset:    "\033[0m",
			dirPath:  tempDir,
			names:    []string{"\033[0mnew.txt\033[0m"},
			expected: []string{"\033[0mnew.txt\033[0m", "\033[0mold.txt\033[0m"},
		},
		{
			name:     "insert dot directory",
			filename: ".",
			colour:   "\033[01;34m",
			reset:    "\033[0m",
			dirPath:  tempDir,
			names:    []string{"\033[0mnew.txt\033[0m"},
			expected: []string{"\033[01;34m.\033[0m", "\033[0mnew.txt\033[0m"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := InsertSortedByTime(tt.filename, tt.colour, tt.reset, tt.dirPath, tt.names)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("InsertSortedByTime() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInsertSortedByTime_ErrorHandling(t *testing.T) {
	// Test with non-existent file - should fall back to alphabetical sorting
	result := InsertSortedByTime("nonexistent.txt", "\033[0m", "\033[0m", "/nonexistent", []string{})
	expected := []string{"\033[0mnonexistent.txt\033[0m"}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("InsertSortedByTime() with non-existent file = %v, want %v", result, expected)
	}
}

func TestInsertSortedLongByTime(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create test files with different modification times
	file1 := filepath.Join(tempDir, "old.txt")
	file2 := filepath.Join(tempDir, "new.txt")

	// Create files
	err := os.WriteFile(file1, []byte("old content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Wait a bit to ensure different timestamps
	time.Sleep(10 * time.Millisecond)

	err = os.WriteFile(file2, []byte("new content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	newLine := fmt.Sprintf("-rw-r--r--  1 user group  100 Jan  1 12:00 new.txt")
	oldLine := fmt.Sprintf("-rw-r--r--  1 user group  100 Jan  1 12:00 old.txt")

	result := InsertSortedLongByTime(newLine, tempDir, []string{oldLine})
	expected := []string{newLine, oldLine}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("InsertSortedLongByTime() = %v, want %v", result, expected)
	}
}

func TestInsertSortedLongByTime_ErrorHandling(t *testing.T) {
	// Test with non-existent file - should fall back to alphabetical sorting
	line := "-rw-r--r--  1 user group  100 Jan  1 12:00 nonexistent.txt"
	result := InsertSortedLongByTime(line, "/nonexistent", []string{})
	expected := []string{line}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("InsertSortedLongByTime() with non-existent file = %v, want %v", result, expected)
	}
}
