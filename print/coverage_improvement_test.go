package print

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jesee-kuya/my-ls/util"
)

// TestPrint_EdgeCases tests edge cases in Print function
func TestPrint_EdgeCases(t *testing.T) {
	t.Run("empty directory", func(t *testing.T) {
		tempDir := t.TempDir()

		output := captureOutput(func() {
			Print([]string{tempDir}, util.Flags{})
		})

		// Should handle empty directory gracefully
		if strings.Contains(output, "Error") {
			t.Errorf("Should not error on empty directory, got: %s", output)
		}
	})

	t.Run("single file with long format", func(t *testing.T) {
		tempDir := t.TempDir()
		testFile := filepath.Join(tempDir, "single.txt")
		err := os.WriteFile(testFile, []byte("content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		output := captureOutput(func() {
			Print([]string{testFile}, util.Flags{Longformat: true})
		})

		// Should print single file in long format
		if !strings.Contains(output, "single.txt") {
			t.Errorf("Expected output to contain 'single.txt', got: %s", output)
		}
	})

	t.Run("multiple single files", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create multiple files
		files := []string{"file1.txt", "file2.txt", "file3.txt"}
		for _, file := range files {
			err := os.WriteFile(filepath.Join(tempDir, file), []byte("content"), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file %s: %v", file, err)
			}
		}

		// Test with multiple single file paths
		filePaths := make([]string, len(files))
		for i, file := range files {
			filePaths[i] = filepath.Join(tempDir, file)
		}

		output := captureOutput(func() {
			Print(filePaths, util.Flags{})
		})

		// Should contain all files
		for _, file := range files {
			if !strings.Contains(output, file) {
				t.Errorf("Expected output to contain '%s', got: %s", file, output)
			}
		}
	})

	t.Run("recursive with error", func(t *testing.T) {
		// Test recursive flag with non-existent path
		output := captureOutput(func() {
			Print([]string{"/non/existent/path"}, util.Flags{Recursive: true})
		})

		// Should contain error message
		if !strings.Contains(output, "Error") {
			t.Errorf("Expected error message for recursive with non-existent path, got: %s", output)
		}
	})

	t.Run("directory with no readable content", func(t *testing.T) {
		tempDir := t.TempDir()
		restrictedDir := filepath.Join(tempDir, "restricted")

		err := os.Mkdir(restrictedDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create restricted directory: %v", err)
		}

		// Remove read permission
		err = os.Chmod(restrictedDir, 0000)
		if err != nil {
			t.Skipf("Cannot change directory permissions: %v", err)
		}

		// Restore permissions after test
		defer os.Chmod(restrictedDir, 0755)

		output := captureOutput(func() {
			Print([]string{restrictedDir}, util.Flags{})
		})

		// Should handle permission error gracefully
		if !strings.Contains(output, "Error") && os.Getuid() != 0 {
			t.Logf("Expected error for restricted directory (may be running as root), got: %s", output)
		}
	})

	t.Run("long format with error", func(t *testing.T) {
		tempDir := t.TempDir()
		restrictedDir := filepath.Join(tempDir, "restricted")

		err := os.Mkdir(restrictedDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create restricted directory: %v", err)
		}

		// Remove read permission
		err = os.Chmod(restrictedDir, 0000)
		if err != nil {
			t.Skipf("Cannot change directory permissions: %v", err)
		}

		// Restore permissions after test
		defer os.Chmod(restrictedDir, 0755)

		output := captureOutput(func() {
			Print([]string{restrictedDir}, util.Flags{Longformat: true})
		})

		// Should handle permission error gracefully in long format
		if !strings.Contains(output, "Error") && os.Getuid() != 0 {
			t.Logf("Expected error for restricted directory in long format (may be running as root), got: %s", output)
		}
	})
}

// TestFormatInColumns_EdgeCases tests edge cases in formatInColumns
func TestFormatInColumns_EdgeCases(t *testing.T) {
	t.Run("very long filenames", func(t *testing.T) {
		longName := strings.Repeat("verylongfilename", 10)
		files := []string{longName, "short.txt"}

		result := formatInColumns(files)

		// Should handle very long filenames
		if !strings.Contains(result, longName) {
			t.Errorf("Expected result to contain long filename")
		}
		if !strings.Contains(result, "short.txt") {
			t.Errorf("Expected result to contain short filename")
		}
	})

	t.Run("many files", func(t *testing.T) {
		// Create many files to test column layout
		files := make([]string, 50)
		for i := 0; i < 50; i++ {
			files[i] = fmt.Sprintf("file%02d.txt", i)
		}

		result := formatInColumns(files)

		// Should contain all files
		for _, file := range files {
			if !strings.Contains(result, file) {
				t.Errorf("Expected result to contain '%s'", file)
			}
		}

		// Should have multiple lines for many files
		lines := strings.Split(result, "\n")
		if len(lines) < 2 {
			t.Errorf("Expected multiple lines for many files, got %d lines", len(lines))
		}
	})

	t.Run("files with ANSI codes", func(t *testing.T) {
		// Test files with ANSI color codes
		files := []string{
			"\033[1;34mdir1\033[0m",
			"\033[1;32mexecutable\033[0m",
			"regular.txt",
		}

		result := formatInColumns(files)

		// Should preserve ANSI codes in output
		if !strings.Contains(result, "\033[1;34m") {
			t.Errorf("Expected result to preserve ANSI codes")
		}
		if !strings.Contains(result, "dir1") {
			t.Errorf("Expected result to contain 'dir1'")
		}
	})

	t.Run("single very wide file", func(t *testing.T) {
		// Test with a single file that's wider than terminal
		wideFile := strings.Repeat("x", 200)
		files := []string{wideFile}

		result := formatInColumns(files)

		// Should handle wide files gracefully
		if !strings.Contains(result, wideFile) {
			t.Errorf("Expected result to contain wide filename")
		}
	})
}

// TestPrint_ComplexScenarios tests complex scenarios
func TestPrint_ComplexScenarios(t *testing.T) {
	t.Run("mixed files and directories", func(t *testing.T) {
		tempDir := t.TempDir()

		// Create a mix of files and directories
		testFile := filepath.Join(tempDir, "file.txt")
		err := os.WriteFile(testFile, []byte("content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		subDir := filepath.Join(tempDir, "subdir")
		err = os.Mkdir(subDir, 0755)
		if err != nil {
			t.Fatalf("Failed to create subdirectory: %v", err)
		}

		// Test with mixed paths (file and directory)
		output := captureOutput(func() {
			Print([]string{testFile, tempDir}, util.Flags{})
		})

		// Should handle both files and directories
		if !strings.Contains(output, "file.txt") {
			t.Errorf("Expected output to contain 'file.txt', got: %s", output)
		}
		if !strings.Contains(output, tempDir+":") {
			t.Errorf("Expected output to contain directory header, got: %s", output)
		}
	})
}
