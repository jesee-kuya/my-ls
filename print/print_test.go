package print

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jesee-kuya/my-ls/util"
)

// captureOutput captures stdout during function execution
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestGetTerminalWidth(t *testing.T) {
	width := getTerminalWidth()
	if width <= 0 {
		t.Errorf("getTerminalWidth() returned %d, expected positive value", width)
	}
	// Should return at least the default width of 80
	if width < 80 {
		t.Errorf("getTerminalWidth() returned %d, expected at least 80", width)
	}
}

func TestFormatInColumns(t *testing.T) {
	tests := []struct {
		name     string
		files    []string
		expected string
	}{
		{
			name:     "empty list",
			files:    []string{},
			expected: "",
		},
		{
			name:     "single file",
			files:    []string{"file1.txt"},
			expected: "file1.txt",
		},
		{
			name:     "two files",
			files:    []string{"file1.txt", "file2.txt"},
			expected: "file1.txt  file2.txt",
		},
		{
			name:     "files with ANSI codes",
			files:    []string{"\033[01;34mdir1\033[0m", "\033[0mfile1.txt\033[0m"},
			expected: "\033[01;34mdir1\033[0m  \033[0mfile1.txt\033[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatInColumns(tt.files)
			if result != tt.expected {
				t.Errorf("formatInColumns() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestPrint_SingleFile(t *testing.T) {
	// Create a temporary file
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name  string
		flags util.Flags
		paths []string
	}{
		{
			name:  "single file normal format",
			flags: util.Flags{},
			paths: []string{testFile},
		},
		{
			name:  "single file long format",
			flags: util.Flags{Longformat: true},
			paths: []string{testFile},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				Print(tt.paths, tt.flags)
			})

			if !strings.Contains(output, "test.txt") {
				t.Errorf("Expected output to contain 'test.txt', got: %s", output)
			}
		})
	}
}

func TestPrint_Directory(t *testing.T) {
	// Create a temporary directory with files
	tempDir := t.TempDir()

	// Create test files
	testFiles := []string{"file1.txt", "file2.txt", ".hidden"}
	for _, file := range testFiles {
		err := os.WriteFile(filepath.Join(tempDir, file), []byte("content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	tests := []struct {
		name  string
		flags util.Flags
		paths []string
	}{
		{
			name:  "directory normal format",
			flags: util.Flags{},
			paths: []string{tempDir},
		},
		{
			name:  "directory with -a flag",
			flags: util.Flags{ShowAll: true},
			paths: []string{tempDir},
		},
		{
			name:  "directory long format",
			flags: util.Flags{Longformat: true},
			paths: []string{tempDir},
		},
		{
			name:  "directory long format with -a",
			flags: util.Flags{Longformat: true, ShowAll: true},
			paths: []string{tempDir},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				Print(tt.paths, tt.flags)
			})

			// Should contain visible files
			if !strings.Contains(output, "file1.txt") {
				t.Errorf("Expected output to contain 'file1.txt', got: %s", output)
			}

			// Hidden files should only appear with ShowAll flag
			if tt.flags.ShowAll {
				if !strings.Contains(output, ".hidden") {
					t.Errorf("Expected output to contain '.hidden' with ShowAll flag, got: %s", output)
				}
			} else {
				if strings.Contains(output, ".hidden") {
					t.Errorf("Expected output to NOT contain '.hidden' without ShowAll flag, got: %s", output)
				}
			}
		})
	}
}

func TestPrint_MultipleDirectories(t *testing.T) {
	// Create two temporary directories
	tempDir1 := t.TempDir()
	tempDir2 := t.TempDir()

	// Create test files in each directory
	os.WriteFile(filepath.Join(tempDir1, "file1.txt"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(tempDir2, "file2.txt"), []byte("content"), 0644)

	output := captureOutput(func() {
		Print([]string{tempDir1, tempDir2}, util.Flags{})
	})

	// Should contain directory headers
	if !strings.Contains(output, tempDir1+":") {
		t.Errorf("Expected output to contain directory header '%s:', got: %s", tempDir1, output)
	}
	if !strings.Contains(output, tempDir2+":") {
		t.Errorf("Expected output to contain directory header '%s:', got: %s", tempDir2, output)
	}

	// Should contain files from both directories
	if !strings.Contains(output, "file1.txt") {
		t.Errorf("Expected output to contain 'file1.txt', got: %s", output)
	}
	if !strings.Contains(output, "file2.txt") {
		t.Errorf("Expected output to contain 'file2.txt', got: %s", output)
	}
}

func TestPrint_RecursiveFlag(t *testing.T) {
	// Create a nested directory structure
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create files in both directories
	os.WriteFile(filepath.Join(tempDir, "root.txt"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(subDir, "sub.txt"), []byte("content"), 0644)

	output := captureOutput(func() {
		Print([]string{tempDir}, util.Flags{Recursive: true})
	})

	// Should contain both directory headers
	if !strings.Contains(output, tempDir+":") {
		t.Errorf("Expected output to contain root directory header, got: %s", output)
	}
	if !strings.Contains(output, subDir+":") {
		t.Errorf("Expected output to contain subdirectory header, got: %s", output)
	}

	// Should contain files from both directories
	if !strings.Contains(output, "root.txt") {
		t.Errorf("Expected output to contain 'root.txt', got: %s", output)
	}
	if !strings.Contains(output, "sub.txt") {
		t.Errorf("Expected output to contain 'sub.txt', got: %s", output)
	}
}

func TestPrint_ErrorHandling(t *testing.T) {
	// Test with non-existent path
	output := captureOutput(func() {
		Print([]string{"/non/existent/path"}, util.Flags{})
	})

	if !strings.Contains(output, "Error") {
		t.Errorf("Expected error message for non-existent path, got: %s", output)
	}
}

func TestPrint_ReverseFlag(t *testing.T) {
	// Create a temporary directory with files
	tempDir := t.TempDir()

	// Create test files with different names to test sorting
	testFiles := []string{"a.txt", "b.txt", "c.txt"}
	for _, file := range testFiles {
		err := os.WriteFile(filepath.Join(tempDir, file), []byte("content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Test normal order
	normalOutput := captureOutput(func() {
		Print([]string{tempDir}, util.Flags{})
	})

	// Test reverse order
	reverseOutput := captureOutput(func() {
		Print([]string{tempDir}, util.Flags{Reverse: true})
	})

	// The outputs should be different
	if normalOutput == reverseOutput {
		t.Errorf("Expected different output with reverse flag, but got same output")
	}
}
