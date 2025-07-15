package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRecursiveFlagIntegration(t *testing.T) {
	// Build the binary first
	cmd := exec.Command("go", "build", "-o", "my-ls-test")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}
	defer os.Remove("my-ls-test")

	// Create test directory structure
	tempDir := t.TempDir()

	// Create nested structure:
	// tempDir/
	//   ├── file1.txt
	//   ├── .hidden_file.txt
	//   ├── dir1/
	//   │   ├── file2.txt
	//   │   └── subdir/
	//   │       └── file3.txt
	//   └── .hidden_dir/
	//       └── hidden_content.txt

	err = os.WriteFile(filepath.Join(tempDir, "file1.txt"), []byte("content1"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file1.txt: %v", err)
	}

	err = os.WriteFile(filepath.Join(tempDir, ".hidden_file.txt"), []byte("hidden"), 0644)
	if err != nil {
		t.Fatalf("Failed to create .hidden_file.txt: %v", err)
	}

	dir1 := filepath.Join(tempDir, "dir1")
	err = os.Mkdir(dir1, 0755)
	if err != nil {
		t.Fatalf("Failed to create dir1: %v", err)
	}

	err = os.WriteFile(filepath.Join(dir1, "file2.txt"), []byte("content2"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file2.txt: %v", err)
	}

	subdir := filepath.Join(dir1, "subdir")
	err = os.Mkdir(subdir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdir: %v", err)
	}

	err = os.WriteFile(filepath.Join(subdir, "file3.txt"), []byte("content3"), 0644)
	if err != nil {
		t.Fatalf("Failed to create file3.txt: %v", err)
	}

	hiddenDir := filepath.Join(tempDir, ".hidden_dir")
	err = os.Mkdir(hiddenDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .hidden_dir: %v", err)
	}

	err = os.WriteFile(filepath.Join(hiddenDir, "hidden_content.txt"), []byte("hidden content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create hidden_content.txt: %v", err)
	}

	t.Run("recursive flag alone", func(t *testing.T) {
		cmd := exec.Command("./my-ls-test", "-R", tempDir)
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Command failed: %v", err)
		}

		outputStr := string(output)

		// Should show all directories
		if !strings.Contains(outputStr, tempDir+":") {
			t.Error("Root directory header not found")
		}
		if !strings.Contains(outputStr, filepath.Join(tempDir, "dir1")+":") {
			t.Error("dir1 header not found")
		}
		if !strings.Contains(outputStr, filepath.Join(tempDir, "dir1", "subdir")+":") {
			t.Error("subdir header not found")
		}

		// Should show files in each directory
		if !strings.Contains(outputStr, "file1.txt") {
			t.Error("file1.txt not found")
		}
		if !strings.Contains(outputStr, "file2.txt") {
			t.Error("file2.txt not found")
		}
		if !strings.Contains(outputStr, "file3.txt") {
			t.Error("file3.txt not found")
		}

		// Should not show hidden files/directories without -a
		if strings.Contains(outputStr, ".hidden_file.txt") {
			t.Error("Hidden file should not be shown without -a flag")
		}
		if strings.Contains(outputStr, ".hidden_dir") {
			t.Error("Hidden directory should not be shown without -a flag")
		}
	})

	t.Run("recursive with show all (-aR)", func(t *testing.T) {
		cmd := exec.Command("./my-ls-test", "-aR", tempDir)
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Command failed: %v", err)
		}

		outputStr := string(output)

		// Should show hidden files and directories
		if !strings.Contains(outputStr, ".hidden_file.txt") {
			t.Error("Hidden file not found with -a flag")
		}
		if !strings.Contains(outputStr, ".hidden_dir:") {
			t.Error("Hidden directory header not found with -a flag")
		}
		if !strings.Contains(outputStr, "hidden_content.txt") {
			t.Error("Content of hidden directory not found")
		}

		// Should show . and .. entries (they appear as ".  .." in the output)
		// The output format shows them as separate entries, so check for both
		if !strings.Contains(outputStr, ".  ..") {
			// Alternative check: look for . and .. as separate entries
			if !strings.Contains(outputStr, ".") || !strings.Contains(outputStr, "..") {
				t.Errorf(". and .. entries not found with -a flag. Output was: %s", outputStr)
			}
		}
	})

	t.Run("recursive with long format (-lR)", func(t *testing.T) {
		cmd := exec.Command("./my-ls-test", "-lR", tempDir)
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Command failed: %v", err)
		}

		outputStr := string(output)

		// Should show long format with permissions, sizes, etc.
		if !strings.Contains(outputStr, "-rw-") {
			t.Error("File permissions not found in long format")
		}
		if !strings.Contains(outputStr, "drwx") {
			t.Error("Directory permissions not found in long format")
		}
		if !strings.Contains(outputStr, "total") {
			t.Error("Total block count not found in long format")
		}
	})

	t.Run("recursive with reverse (-rR)", func(t *testing.T) {
		cmd := exec.Command("./my-ls-test", "-rR", tempDir)
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Command failed: %v", err)
		}

		outputStr := string(output)

		// Look for the pattern where directory header and content are on the same line
		// e.g., "test_dir:subdir1  file1.txt"
		if !strings.Contains(outputStr, tempDir+":") {
			t.Error("Root directory header not found")
		}

		// Find the line that contains both the directory header and content
		lines := strings.Split(outputStr, "\n")
		var rootContentLine string
		for _, line := range lines {
			if strings.HasPrefix(line, tempDir+":") && strings.Contains(line, "file1.txt") {
				rootContentLine = line
				break
			}
		}

		if rootContentLine == "" {
			t.Errorf("Could not find root directory content line. Output was: %s", outputStr)
		} else {
			// In reverse order, file1.txt should come before dir1 (reverse alphabetical)
			// Normal alphabetical: dir1, file1.txt
			// Reverse alphabetical: file1.txt, dir1
			// Remove the directory header part to check just the content
			contentPart := strings.TrimPrefix(rootContentLine, tempDir+":")
			dir1Index := strings.Index(contentPart, "dir1")
			file1Index := strings.Index(contentPart, "file1.txt")
			if dir1Index == -1 || file1Index == -1 {
				t.Errorf("Could not find both dir1 and file1.txt in content: %s", contentPart)
			} else if file1Index > dir1Index {
				t.Error("Reverse order not working: file1.txt should come before dir1 in reverse order")
			}
		}
	})

	t.Run("all flags combined (-alrR)", func(t *testing.T) {
		cmd := exec.Command("./my-ls-test", "-alrR", tempDir)
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Command failed: %v", err)
		}

		outputStr := string(output)

		// Should have all features: recursive, show all, long format, reverse
		if !strings.Contains(outputStr, ".hidden_dir:") {
			t.Error("Hidden directory not shown with -a")
		}
		if !strings.Contains(outputStr, "-rw-") {
			t.Error("Long format not working")
		}
		if !strings.Contains(outputStr, "total") {
			t.Error("Total block count not found")
		}

		// Check that all directories are processed
		expectedDirs := []string{
			tempDir + ":",
			filepath.Join(tempDir, "dir1") + ":",
			filepath.Join(tempDir, "dir1", "subdir") + ":",
			filepath.Join(tempDir, ".hidden_dir") + ":",
		}

		for _, expectedDir := range expectedDirs {
			if !strings.Contains(outputStr, expectedDir) {
				t.Errorf("Expected directory header %s not found", expectedDir)
			}
		}
	})

	t.Run("recursive with multiple paths", func(t *testing.T) {
		cmd := exec.Command("./my-ls-test", "-R", dir1, hiddenDir)
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("Command failed: %v", err)
		}

		outputStr := string(output)

		// Should process both paths recursively
		if !strings.Contains(outputStr, dir1+":") {
			t.Error("dir1 header not found")
		}
		if !strings.Contains(outputStr, hiddenDir+":") {
			t.Error("hidden_dir header not found")
		}
		if !strings.Contains(outputStr, filepath.Join(dir1, "subdir")+":") {
			t.Error("subdir header not found")
		}
	})
}
