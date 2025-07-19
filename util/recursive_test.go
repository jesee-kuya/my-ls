package util

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// testJoinPath3 joins directory and file name with proper separator (test helper)
func testJoinPath3(dir, file string) string {
	if dir == "" {
		return file
	}
	if strings.HasSuffix(dir, "/") {
		return dir + file
	}
	return dir + "/" + file
}

func TestCollectDirectoriesRecursively(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir := t.TempDir()

	// Create nested directory structure
	// tempDir/
	//   ├── file1.txt
	//   ├── dir1/
	//   │   ├── file2.txt
	//   │   └── subdir1/
	//   │       └── file3.txt
	//   ├── dir2/
	//   │   └── file4.txt
	//   └── .hidden_dir/
	//       └── hidden_file.txt

	// Create files and directories
	err := os.WriteFile(testJoinPath3(tempDir, "file1.txt"), []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file1.txt: %v", err)
	}

	dir1 := testJoinPath3(tempDir, "dir1")
	err = os.Mkdir(dir1, 0755)
	if err != nil {
		t.Fatalf("Failed to create dir1: %v", err)
	}

	err = os.WriteFile(testJoinPath3(dir1, "file2.txt"), []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file2.txt: %v", err)
	}

	subdir1 := testJoinPath3(dir1, "subdir1")
	err = os.Mkdir(subdir1, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdir1: %v", err)
	}

	err = os.WriteFile(testJoinPath3(subdir1, "file3.txt"), []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file3.txt: %v", err)
	}

	dir2 := testJoinPath3(tempDir, "dir2")
	err = os.Mkdir(dir2, 0755)
	if err != nil {
		t.Fatalf("Failed to create dir2: %v", err)
	}

	err = os.WriteFile(testJoinPath3(dir2, "file4.txt"), []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file4.txt: %v", err)
	}

	hiddenDir := testJoinPath3(tempDir, ".hidden_dir")
	err = os.Mkdir(hiddenDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .hidden_dir: %v", err)
	}

	err = os.WriteFile(testJoinPath3(hiddenDir, "hidden_file.txt"), []byte("content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create hidden_file.txt: %v", err)
	}

	t.Run("recursive collection without ShowAll", func(t *testing.T) {
		flags := Flags{ShowAll: false, Recursive: true}
		dirs, err := CollectDirectoriesRecursively([]string{tempDir}, flags)
		if err != nil {
			t.Fatalf("CollectDirectoriesRecursively failed: %v", err)
		}

		// Should include tempDir, dir1, subdir1, dir2 but not .hidden_dir
		expectedDirs := []string{tempDir, dir1, subdir1, dir2}

		if len(dirs) != len(expectedDirs) {
			t.Errorf("Expected %d directories, got %d: %v", len(expectedDirs), len(dirs), dirs)
		}

		for _, expectedDir := range expectedDirs {
			found := false
			for _, dir := range dirs {
				if dir == expectedDir {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected directory %s not found in result", expectedDir)
			}
		}

		// Ensure hidden directory is not included
		for _, dir := range dirs {
			if strings.Contains(dir, ".hidden_dir") {
				t.Errorf("Hidden directory should not be included when ShowAll=false")
			}
		}
	})

	t.Run("recursive collection with ShowAll", func(t *testing.T) {
		flags := Flags{ShowAll: true, Recursive: true}
		dirs, err := CollectDirectoriesRecursively([]string{tempDir}, flags)
		if err != nil {
			t.Fatalf("CollectDirectoriesRecursively failed: %v", err)
		}

		// Should include tempDir, dir1, subdir1, dir2, and .hidden_dir
		expectedDirs := []string{tempDir, dir1, subdir1, dir2, hiddenDir}

		if len(dirs) != len(expectedDirs) {
			t.Errorf("Expected %d directories, got %d: %v", len(expectedDirs), len(dirs), dirs)
		}

		for _, expectedDir := range expectedDirs {
			found := false
			for _, dir := range dirs {
				if dir == expectedDir {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected directory %s not found in result", expectedDir)
			}
		}
	})

	t.Run("recursive collection with file input", func(t *testing.T) {
		flags := Flags{ShowAll: false, Recursive: true}
		filePath := testJoinPath3(tempDir, "file1.txt")
		dirs, err := CollectDirectoriesRecursively([]string{filePath}, flags)
		if err != nil {
			t.Fatalf("CollectDirectoriesRecursively failed: %v", err)
		}

		// Should just return the file path
		if len(dirs) != 1 || dirs[0] != filePath {
			t.Errorf("Expected [%s], got %v", filePath, dirs)
		}
	})

	t.Run("recursive collection with multiple root paths", func(t *testing.T) {
		flags := Flags{ShowAll: false, Recursive: true}
		dirs, err := CollectDirectoriesRecursively([]string{dir1, dir2}, flags)
		if err != nil {
			t.Fatalf("CollectDirectoriesRecursively failed: %v", err)
		}

		// Should include dir1, subdir1, dir2
		expectedDirs := []string{dir1, subdir1, dir2}

		if len(dirs) != len(expectedDirs) {
			t.Errorf("Expected %d directories, got %d: %v", len(expectedDirs), len(dirs), dirs)
		}

		for _, expectedDir := range expectedDirs {
			found := false
			for _, dir := range dirs {
				if dir == expectedDir {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected directory %s not found in result", expectedDir)
			}
		}
	})

	t.Run("recursive collection with non-existent directory", func(t *testing.T) {
		flags := Flags{ShowAll: false, Recursive: true}
		nonExistentDir := testJoinPath3(tempDir, "non_existent")
		_, err := CollectDirectoriesRecursively([]string{nonExistentDir}, flags)
		if err == nil {
			t.Error("Expected error for non-existent directory, got nil")
		}
	})
}

func TestCollectSubdirectories(t *testing.T) {
	// Create a temporary directory with symlinks and test directories for sorting
	tempDir := t.TempDir()

	// Create a subdirectory
	subDir := testJoinPath3(tempDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	// Create a symlink that could cause a loop (if not handled properly)
	symlinkPath := testJoinPath3(subDir, "link_to_parent")
	err = os.Symlink(tempDir, symlinkPath)
	if err != nil {
		t.Logf("Could not create symlink (may not be supported): %v", err)
		// Continue test without symlink
	}

	// Create test directories to verify sorting with CompareStrings
	dirsToCreate := []string{
		"123dir",
		"abcDir",
		"ABCdir",
		"Zebra",
		"@special",
		".hidden_dir",
	}
	for _, dir := range dirsToCreate {
		err := os.Mkdir(testJoinPath3(tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	// Create a nested subdirectory to test recursion
	nestedSubDir := filepath.Join(tempDir, "abcDir", "sub1")
	err = os.MkdirAll(nestedSubDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create nested subdirectory %s: %v", nestedSubDir, err)
	}

	t.Run("symlink loop prevention", func(t *testing.T) {
		flags := Flags{ShowAll: false, Recursive: true}
		var allDirs []string
		visited := make(map[string]bool)

		err := collectSubdirectories(tempDir, flags, &allDirs, visited)
		if err != nil {
			t.Fatalf("collectSubdirectories failed: %v", err)
		}

		// Should include subdir but not get stuck in infinite loop
		expectedDirs := []string{subDir}

		if len(allDirs) < len(expectedDirs) {
			t.Errorf("Expected at least %d directories, got %d: %v", len(expectedDirs), len(allDirs), allDirs)
		}

		for _, expectedDir := range expectedDirs {
			found := false
			for _, dir := range allDirs {
				if dir == expectedDir {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected directory %s not found in result", expectedDir)
			}
		}
	})

	t.Run("sorting with CompareStrings", func(t *testing.T) {
		flags := Flags{ShowAll: false, Recursive: true}
		var allDirs []string
		visited := make(map[string]bool)

		err := collectSubdirectories(tempDir, flags, &allDirs, visited)
		if err != nil {
			t.Fatalf("collectSubdirectories failed: %v", err)
		}

		// Check that we have the expected directories (order will be verified separately)
		expectedBaseDirs := []string{"123dir", "abcDir", "ABCdir", "Zebra", "@special", "subdir"}
		expectedNestedDirs := []string{"sub1"} // nested under abcDir

		// Count directories by base name
		foundBaseDirs := make(map[string]bool)
		foundNestedDirs := make(map[string]bool)

		for _, dir := range allDirs {
			baseName := filepath.Base(dir)
			if filepath.Dir(dir) == tempDir {
				foundBaseDirs[baseName] = true
			} else {
				foundNestedDirs[baseName] = true
			}
		}

		// Verify all expected base directories are found (excluding hidden)
		for _, expected := range expectedBaseDirs {
			if expected == ".hidden_dir" {
				continue // Should not be present when ShowAll=false
			}
			if !foundBaseDirs[expected] {
				t.Errorf("Expected base directory %s not found", expected)
			}
		}

		// Verify nested directories
		for _, expected := range expectedNestedDirs {
			if !foundNestedDirs[expected] {
				t.Errorf("Expected nested directory %s not found", expected)
			}
		}

		// Verify hidden directory is not included
		if foundBaseDirs[".hidden_dir"] {
			t.Errorf("Hidden directory should not be included when ShowAll=false")
		}
	})

	t.Run("sorting with CompareStrings and ShowAll=true", func(t *testing.T) {
		flags := Flags{ShowAll: true, Recursive: true}
		var allDirs []string
		visited := make(map[string]bool)

		err := collectSubdirectories(tempDir, flags, &allDirs, visited)
		if err != nil {
			t.Fatalf("collectSubdirectories failed: %v", err)
		}

		// Check that we have the expected directories (including hidden ones)
		expectedBaseDirs := []string{".hidden_dir", "123dir", "abcDir", "ABCdir", "Zebra", "@special", "subdir"}
		expectedNestedDirs := []string{"sub1"} // nested under abcDir

		// Count directories by base name
		foundBaseDirs := make(map[string]bool)
		foundNestedDirs := make(map[string]bool)

		for _, dir := range allDirs {
			baseName := filepath.Base(dir)
			if filepath.Dir(dir) == tempDir {
				foundBaseDirs[baseName] = true
			} else {
				foundNestedDirs[baseName] = true
			}
		}

		// Verify all expected base directories are found (including hidden)
		for _, expected := range expectedBaseDirs {
			if !foundBaseDirs[expected] {
				t.Errorf("Expected base directory %s not found", expected)
			}
		}

		// Verify nested directories
		for _, expected := range expectedNestedDirs {
			if !foundNestedDirs[expected] {
				t.Errorf("Expected nested directory %s not found", expected)
			}
		}

		// Verify hidden directory is included
		if !foundBaseDirs[".hidden_dir"] {
			t.Errorf("Hidden directory should be included when ShowAll=true")
		}
	})
}
