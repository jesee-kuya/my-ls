package util

import (
	"errors"
	"os"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// testJoinPath2 joins directory and file name with proper separator (test helper)
func testJoinPath2(dir, file string) string {
	if dir == "" {
		return file
	}
	if strings.HasSuffix(dir, "/") {
		return dir + file
	}
	return dir + "/" + file
}

func TestReadDirNames(t *testing.T) {
	t.Run("Table-driven ReadDirNames tests", func(t *testing.T) {
		base := t.TempDir()

		dirWithFiles := testJoinPath2(base, "with_files")
		os.Mkdir(dirWithFiles, 0o755)
		os.WriteFile(testJoinPath2(dirWithFiles, "a.txt"), []byte("a"), 0o644)
		os.WriteFile(testJoinPath2(dirWithFiles, "b.txt"), []byte("b"), 0o644)

		emptyDir := testJoinPath2(base, "empty")
		os.Mkdir(emptyDir, 0o755)

		missingDir := testJoinPath2(base, "missing")

		tests := []struct {
			name      string
			input     string
			want      []string
			expectErr error
		}{
			{
				name:      "Valid directory with files",
				input:     dirWithFiles,
				want:      []string{"\033[0ma.txt\033[0m", "\033[0mb.txt\033[0m"},
				expectErr: nil,
			},
			{
				name:      "Empty directory",
				input:     emptyDir,
				want:      nil,
				expectErr: nil,
			},
			{
				name:      "Non-existent directory",
				input:     missingDir,
				want:      nil,
				expectErr: os.ErrNotExist,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := ReadDirNames(tt.input, Flags{ShowAll: false})

				if tt.expectErr != nil {
					if err == nil {
						t.Errorf("expected error %v, got nil", tt.expectErr)
						return
					}

					if errors.Is(tt.expectErr, os.ErrNotExist) {
						if !os.IsNotExist(err) {
							t.Errorf("expected a file-not-exist error, got: %v", err)
						}
						return
					}

					if err.Error() != tt.expectErr.Error() {
						t.Errorf("expected error %q, got %q", tt.expectErr.Error(), err.Error())
					}
					return
				}

				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}

				sort.Strings(got)
				sort.Strings(tt.want)

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("expected %v, got %v", tt.want, got)
				}
			})
		}
	})
}

func TestReadDirNamesLong(t *testing.T) {
	tempDir := t.TempDir()

	// Create test files
	testFile := testJoinPath2(tempDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	hiddenFile := testJoinPath2(tempDir, ".hidden")
	err = os.WriteFile(hiddenFile, []byte("hidden content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create hidden file: %v", err)
	}

	tests := []struct {
		name  string
		flags Flags
	}{
		{
			name:  "long format without ShowAll",
			flags: Flags{Longformat: true, ShowAll: false},
		},
		{
			name:  "long format with ShowAll",
			flags: Flags{Longformat: true, ShowAll: true},
		},
		{
			name:  "long format with TimeSort",
			flags: Flags{Longformat: true, TimeSort: true},
		},
		{
			name:  "long format with Reverse",
			flags: Flags{Longformat: true, Reverse: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ReadDirNamesLong(tempDir, tt.flags)
			if err != nil {
				t.Fatalf("ReadDirNamesLong() error = %v", err)
			}

			// Should have at least the "total" line
			if len(result) == 0 {
				t.Errorf("ReadDirNamesLong() returned empty result")
			}

			// First line should be "total X"
			if !strings.HasPrefix(result[0], "total ") {
				t.Errorf("ReadDirNamesLong() first line should start with 'total ', got: %s", result[0])
			}

			// Should contain test.txt
			found := false
			for _, line := range result {
				if strings.Contains(line, "test.txt") {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("ReadDirNamesLong() should contain test.txt, got: %v", result)
			}

			// Hidden file should only appear with ShowAll
			foundHidden := false
			for _, line := range result {
				if strings.Contains(line, ".hidden") {
					foundHidden = true
					break
				}
			}
			if tt.flags.ShowAll && !foundHidden {
				t.Errorf("ReadDirNamesLong() with ShowAll should contain .hidden, got: %v", result)
			}
			if !tt.flags.ShowAll && foundHidden {
				t.Errorf("ReadDirNamesLong() without ShowAll should not contain .hidden, got: %v", result)
			}
		})
	}
}

func TestGetFileColor(t *testing.T) {
	tests := []struct {
		name     string
		mode     os.FileMode
		filename string
		expected string
	}{
		{
			name:     "regular file",
			mode:     0644,
			filename: "file.txt",
			expected: reset,
		},
		{
			name:     "directory",
			mode:     os.ModeDir | 0755,
			filename: "dirname",
			expected: dirColour,
		},
		{
			name:     "executable file",
			mode:     0755,
			filename: "executable",
			expected: exeColour,
		},
		{
			name:     "symlink",
			mode:     os.ModeSymlink | 0777,
			filename: "symlink",
			expected: symlinkColour,
		},
		{
			name:     "socket",
			mode:     os.ModeSocket | 0755,
			filename: "socket",
			expected: socketColour,
		},
		{
			name:     "named pipe",
			mode:     os.ModeNamedPipe | 0644,
			filename: "pipe",
			expected: pipeColour,
		},
		{
			name:     "device",
			mode:     os.ModeDevice | 0644,
			filename: "device",
			expected: deviceColour,
		},
		{
			name:     "tar archive",
			mode:     0644,
			filename: "archive.tar",
			expected: archiveColour,
		},
		{
			name:     "gz archive",
			mode:     0644,
			filename: "archive.gz",
			expected: archiveColour,
		},
		{
			name:     "zip archive",
			mode:     0644,
			filename: "archive.zip",
			expected: archiveColour,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getFileColor(tt.mode, tt.filename)
			if result != tt.expected {
				t.Errorf("getFileColor(%v, %q) = %q, want %q", tt.mode, tt.filename, result, tt.expected)
			}
		})
	}
}
