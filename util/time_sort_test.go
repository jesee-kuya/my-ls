package util

import (
	"os"
	"strings"
	"testing"
	"time"
)

// testJoinPath5 joins directory and file name with proper separator (test helper)
func testJoinPath5(dir, file string) string {
	if dir == "" {
		return file
	}
	if strings.HasSuffix(dir, "/") {
		return dir + file
	}
	return dir + "/" + file
}

func TestTimeSortFunctionality(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create files with different modification times
	// We'll create them in reverse chronological order to test sorting
	files := []struct {
		name    string
		content string
		delay   time.Duration
	}{
		{"oldest.txt", "oldest content", 0},
		{"middle.txt", "middle content", 100 * time.Millisecond},
		{"newest.txt", "newest content", 200 * time.Millisecond},
	}

	baseTime := time.Now()
	for i, file := range files {
		time.Sleep(file.delay)
		filePath := testJoinPath5(tempDir, file.name)
		err := os.WriteFile(filePath, []byte(file.content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file.name, err)
		}

		// Set specific modification times to ensure consistent ordering
		modTime := baseTime.Add(time.Duration(i) * time.Second)
		err = os.Chtimes(filePath, modTime, modTime)
		if err != nil {
			t.Fatalf("Failed to set modification time for %s: %v", file.name, err)
		}
	}

	t.Run("ReadDirNames with TimeSort=false should sort alphabetically", func(t *testing.T) {
		names, err := ReadDirNames(tempDir, Flags{TimeSort: false})
		if err != nil {
			t.Fatalf("ReadDirNames failed: %v", err)
		}

		// Should be in alphabetical order
		expectedOrder := []string{"middle.txt", "newest.txt", "oldest.txt"}
		if len(names) != len(expectedOrder) {
			t.Errorf("Expected %d files, got %d", len(expectedOrder), len(names))
		}

		for i, expected := range expectedOrder {
			if i >= len(names) {
				break
			}
			actual := StripANSI(names[i])
			if actual != expected {
				t.Errorf("Position %d: expected %s, got %s", i, expected, actual)
			}
		}
	})

	t.Run("ReadDirNames with TimeSort=true should sort by modification time", func(t *testing.T) {
		names, err := ReadDirNames(tempDir, Flags{TimeSort: true})
		if err != nil {
			t.Fatalf("ReadDirNames failed: %v", err)
		}

		// Should be in time order (newest first)
		expectedOrder := []string{"newest.txt", "middle.txt", "oldest.txt"}
		if len(names) != len(expectedOrder) {
			t.Errorf("Expected %d files, got %d", len(expectedOrder), len(names))
		}

		for i, expected := range expectedOrder {
			if i >= len(names) {
				break
			}
			actual := StripANSI(names[i])
			if actual != expected {
				t.Errorf("Position %d: expected %s, got %s", i, expected, actual)
			}
		}
	})

	t.Run("ReadDirNames with TimeSort=true and Reverse=true", func(t *testing.T) {
		names, err := ReadDirNames(tempDir, Flags{TimeSort: true, Reverse: true})
		if err != nil {
			t.Fatalf("ReadDirNames failed: %v", err)
		}

		// Should be in reverse time order (oldest first)
		expectedOrder := []string{"oldest.txt", "middle.txt", "newest.txt"}
		if len(names) != len(expectedOrder) {
			t.Errorf("Expected %d files, got %d", len(expectedOrder), len(names))
		}

		for i, expected := range expectedOrder {
			if i >= len(names) {
				break
			}
			actual := StripANSI(names[i])
			if actual != expected {
				t.Errorf("Position %d: expected %s, got %s", i, expected, actual)
			}
		}
	})

	t.Run("ReadDirNamesLong with TimeSort=true should sort by modification time", func(t *testing.T) {
		lines, err := ReadDirNamesLong(tempDir, Flags{TimeSort: true})
		if err != nil {
			t.Fatalf("ReadDirNamesLong failed: %v", err)
		}

		// Skip the "total" line
		if len(lines) < 1 || !strings.HasPrefix(lines[0], "total") {
			t.Fatalf("Expected first line to start with 'total', got: %v", lines)
		}

		fileLines := lines[1:] // Skip the total line
		expectedOrder := []string{"newest.txt", "middle.txt", "oldest.txt"}

		if len(fileLines) != len(expectedOrder) {
			t.Errorf("Expected %d file lines, got %d", len(expectedOrder), len(fileLines))
		}

		for i, expected := range expectedOrder {
			if i >= len(fileLines) {
				break
			}
			// Extract filename from long format line
			actual := StripANSI(StripLong(fileLines[i]))
			if actual != expected {
				t.Errorf("Position %d: expected %s, got %s", i, expected, actual)
			}
		}
	})

	t.Run("ReadDirNamesLong with TimeSort=true and Reverse=true", func(t *testing.T) {
		lines, err := ReadDirNamesLong(tempDir, Flags{TimeSort: true, Reverse: true})
		if err != nil {
			t.Fatalf("ReadDirNamesLong failed: %v", err)
		}

		// Skip the "total" line
		if len(lines) < 1 || !strings.HasPrefix(lines[0], "total") {
			t.Fatalf("Expected first line to start with 'total', got: %v", lines)
		}

		fileLines := lines[1:] // Skip the total line
		expectedOrder := []string{"oldest.txt", "middle.txt", "newest.txt"}

		if len(fileLines) != len(expectedOrder) {
			t.Errorf("Expected %d file lines, got %d", len(expectedOrder), len(fileLines))
		}

		for i, expected := range expectedOrder {
			if i >= len(fileLines) {
				break
			}
			// Extract filename from long format line
			actual := StripANSI(StripLong(fileLines[i]))
			if actual != expected {
				t.Errorf("Position %d: expected %s, got %s", i, expected, actual)
			}
		}
	})

	t.Run("TimeSort with ShowAll=true should handle . and .. correctly", func(t *testing.T) {
		names, err := ReadDirNames(tempDir, Flags{TimeSort: true, ShowAll: true})
		if err != nil {
			t.Fatalf("ReadDirNames failed: %v", err)
		}

		// . and .. should still be first
		if len(names) < 2 {
			t.Fatalf("Expected at least 2 entries (. and ..), got %d", len(names))
		}

		if StripANSI(names[0]) != "." {
			t.Errorf("Expected first entry to be '.', got %s", StripANSI(names[0]))
		}
		if StripANSI(names[1]) != ".." {
			t.Errorf("Expected second entry to be '..', got %s", StripANSI(names[1]))
		}

		// The rest should be in time order
		fileEntries := names[2:]
		expectedOrder := []string{"newest.txt", "middle.txt", "oldest.txt"}

		if len(fileEntries) != len(expectedOrder) {
			t.Errorf("Expected %d file entries, got %d", len(expectedOrder), len(fileEntries))
		}

		for i, expected := range expectedOrder {
			if i >= len(fileEntries) {
				break
			}
			actual := StripANSI(fileEntries[i])
			if actual != expected {
				t.Errorf("Position %d: expected %s, got %s", i, expected, actual)
			}
		}
	})
}

func TestInsertSortedByTimeIntegration(t *testing.T) {
	tempDir := t.TempDir()

	// Create test files with known modification times
	file1 := testJoinPath5(tempDir, "file1.txt")
	file2 := testJoinPath5(tempDir, "file2.txt")

	baseTime := time.Now()

	// Create file1 (older)
	err := os.WriteFile(file1, []byte("content1"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file1: %v", err)
	}
	err = os.Chtimes(file1, baseTime, baseTime)
	if err != nil {
		t.Fatalf("Failed to set time for file1: %v", err)
	}

	// Create file2 (newer)
	err = os.WriteFile(file2, []byte("content2"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file2: %v", err)
	}
	err = os.Chtimes(file2, baseTime.Add(time.Second), baseTime.Add(time.Second))
	if err != nil {
		t.Fatalf("Failed to set time for file2: %v", err)
	}

	t.Run("InsertSortedByTime should insert newer files first", func(t *testing.T) {
		var names []string

		// Insert older file first
		names = InsertSortedByTime("file1.txt", "", "", tempDir, names)
		// Insert newer file
		names = InsertSortedByTime("file2.txt", "", "", tempDir, names)

		if len(names) != 2 {
			t.Fatalf("Expected 2 files, got %d", len(names))
		}

		// file2.txt should be first (newer)
		if names[0] != "file2.txt" {
			t.Errorf("Expected file2.txt first, got %s", names[0])
		}
		if names[1] != "file1.txt" {
			t.Errorf("Expected file1.txt second, got %s", names[1])
		}
	})
}
