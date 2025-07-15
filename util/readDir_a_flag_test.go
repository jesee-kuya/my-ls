package util

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadDirNames_AFlag(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create some regular files
	regularFiles := []string{"file1.txt", "file2.txt", "README.md"}
	for _, file := range regularFiles {
		err := os.WriteFile(filepath.Join(tempDir, file), []byte("content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Create some hidden files
	hiddenFiles := []string{".hidden1", ".hidden2", ".gitignore"}
	for _, file := range hiddenFiles {
		err := os.WriteFile(filepath.Join(tempDir, file), []byte("content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create hidden file %s: %v", file, err)
		}
	}

	// Create a hidden directory
	hiddenDir := filepath.Join(tempDir, ".hidden_dir")
	err := os.Mkdir(hiddenDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create hidden directory: %v", err)
	}

	t.Run("showAll=false should not show hidden files", func(t *testing.T) {
		names, err := ReadDirNames(tempDir, Flags{ShowAll: false})
		if err != nil {
			t.Fatalf("ReadDirNames failed: %v", err)
		}

		// Should only contain regular files
		expectedCount := len(regularFiles)
		if len(names) != expectedCount {
			t.Errorf("Expected %d files, got %d", expectedCount, len(names))
		}

		// Check that no hidden files are present
		for _, name := range names {
			cleanName := StripANSI(name)
			if strings.HasPrefix(cleanName, ".") {
				t.Errorf("Found hidden file %s when showAll=false", cleanName)
			}
		}

		// Check that all regular files are present
		for _, expectedFile := range regularFiles {
			found := false
			for _, name := range names {
				cleanName := StripANSI(name)
				if cleanName == expectedFile {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected file %s not found", expectedFile)
			}
		}
	})

	t.Run("showAll=true should show all files including hidden ones", func(t *testing.T) {
		names, err := ReadDirNames(tempDir, Flags{ShowAll: true})
		if err != nil {
			t.Fatalf("ReadDirNames failed: %v", err)
		}

		// Should contain . and .. plus regular files plus hidden files
		expectedCount := 2 + len(regularFiles) + len(hiddenFiles) + 1 // +1 for hidden directory
		if len(names) != expectedCount {
			t.Errorf("Expected %d files, got %d", expectedCount, len(names))
		}

		// Check that . and .. are first
		if len(names) >= 2 {
			if StripANSI(names[0]) != "." {
				t.Errorf("Expected first entry to be '.', got %s", StripANSI(names[0]))
			}
			if StripANSI(names[1]) != ".." {
				t.Errorf("Expected second entry to be '..', got %s", StripANSI(names[1]))
			}
		}

		// Check that all regular files are present
		for _, expectedFile := range regularFiles {
			found := false
			for _, name := range names {
				cleanName := StripANSI(name)
				if cleanName == expectedFile {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected regular file %s not found", expectedFile)
			}
		}

		// Check that all hidden files are present
		for _, expectedFile := range hiddenFiles {
			found := false
			for _, name := range names {
				cleanName := StripANSI(name)
				if cleanName == expectedFile {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected hidden file %s not found", expectedFile)
			}
		}

		// Check that hidden directory is present
		found := false
		for _, name := range names {
			cleanName := StripANSI(name)
			if cleanName == ".hidden_dir" {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected hidden directory .hidden_dir not found")
		}
	})

	t.Run("files should be sorted correctly with showAll=true", func(t *testing.T) {
		names, err := ReadDirNames(tempDir, Flags{ShowAll: true})
		if err != nil {
			t.Fatalf("ReadDirNames failed: %v", err)
		}

		// Skip . and .. and check if the rest are sorted
		if len(names) > 2 {
			sortedPortion := names[2:]
			for i := 1; i < len(sortedPortion); i++ {
				prev := strings.ToLower(TrimStart(StripANSI(sortedPortion[i-1])))
				curr := strings.ToLower(TrimStart(StripANSI(sortedPortion[i])))
				if prev > curr {
					t.Errorf("Files not sorted correctly: %s should come before %s (prev: %s, curr: %s)", curr, prev, prev, curr)
				}
			}
		}
	})

	t.Run("files should be sorted correctly with showAll=false", func(t *testing.T) {
		names, err := ReadDirNames(tempDir, Flags{ShowAll: false})
		if err != nil {
			t.Fatalf("ReadDirNames failed: %v", err)
		}

		// Check if all files are sorted
		for i := 1; i < len(names); i++ {
			prev := strings.ToLower(TrimStart(StripANSI(names[i-1])))
			curr := strings.ToLower(TrimStart(StripANSI(names[i])))
			if prev > curr {
				t.Errorf("Files not sorted correctly: %s should come before %s", curr, prev)
			}
		}
	})
}
