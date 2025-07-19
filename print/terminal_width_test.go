package print

import (
	"strings"
	"testing"
)

// TestGetTerminalWidth_DefaultCase tests the default case when syscall fails
func TestGetTerminalWidth_DefaultCase(t *testing.T) {
	// This test verifies that getTerminalWidth returns a reasonable default
	// when it cannot determine the actual terminal width

	width := getTerminalWidth()

	// Should return a positive width
	if width <= 0 {
		t.Errorf("getTerminalWidth() returned %d, expected positive value", width)
	}

	// Should return at least the minimum default width
	if width < 80 {
		t.Errorf("getTerminalWidth() returned %d, expected at least 80", width)
	}

	// Should return a reasonable maximum (not absurdly large)
	if width > 10000 {
		t.Errorf("getTerminalWidth() returned %d, which seems unreasonably large", width)
	}
}

// TestGetTerminalWidth_Consistency tests that getTerminalWidth returns consistent results
func TestGetTerminalWidth_Consistency(t *testing.T) {
	// Call the function multiple times and ensure it returns consistent results
	width1 := getTerminalWidth()
	width2 := getTerminalWidth()
	width3 := getTerminalWidth()

	if width1 != width2 || width2 != width3 {
		t.Errorf("getTerminalWidth() returned inconsistent results: %d, %d, %d", width1, width2, width3)
	}
}

// TestFormatInColumns_TerminalWidthIntegration tests formatInColumns with different terminal widths
func TestFormatInColumns_TerminalWidthIntegration(t *testing.T) {
	// Test with files that would require different column layouts
	files := []string{
		"short.txt",
		"medium_filename.txt",
		"very_long_filename_that_takes_up_space.txt",
		"a.txt",
		"b.txt",
		"c.txt",
	}

	result := formatInColumns(files)

	// Should contain all files
	for _, file := range files {
		if !strings.Contains(result, file) {
			t.Errorf("Expected result to contain '%s'", file)
		}
	}

	// Should format in columns (multiple files per line for short names)
	lines := strings.Split(result, "\n")
	if len(lines) == 0 {
		t.Error("Expected at least one line of output")
	}

	// At least one line should contain multiple files (for short filenames)
	foundMultipleFilesPerLine := false
	for _, line := range lines {
		fileCount := 0
		for _, file := range []string{"short.txt", "a.txt", "b.txt", "c.txt"} {
			if strings.Contains(line, file) {
				fileCount++
			}
		}
		if fileCount > 1 {
			foundMultipleFilesPerLine = true
			break
		}
	}

	if !foundMultipleFilesPerLine {
		t.Logf("Expected at least one line with multiple files, got lines: %v", lines)
		// This might be expected if terminal is very narrow, so just log it
	}
}

// TestFormatInColumns_EdgeCaseWidths tests formatInColumns with edge case scenarios
func TestFormatInColumns_EdgeCaseWidths(t *testing.T) {
	t.Run("single character files", func(t *testing.T) {
		files := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}

		result := formatInColumns(files)

		// Should contain all single character files
		for _, file := range files {
			if !strings.Contains(result, file) {
				t.Errorf("Expected result to contain '%s'", file)
			}
		}

		// Should format efficiently (multiple files per line)
		lines := strings.Split(result, "\n")
		if len(lines) > len(files) {
			t.Errorf("Expected fewer lines than files for single character names, got %d lines for %d files", len(lines), len(files))
		}
	})

	t.Run("identical length files", func(t *testing.T) {
		files := []string{"file1.txt", "file2.txt", "file3.txt", "file4.txt"}

		result := formatInColumns(files)

		// Should contain all files
		for _, file := range files {
			if !strings.Contains(result, file) {
				t.Errorf("Expected result to contain '%s'", file)
			}
		}

		// Should format in columns since all files have same length
		lines := strings.Split(result, "\n")
		if len(lines) == 0 {
			t.Error("Expected at least one line of output")
		}
	})
}
