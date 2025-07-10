package util

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestReadDirNames(t *testing.T) {
	t.Run("Table-driven ReadDirNames tests", func(t *testing.T) {
		base := t.TempDir()

		dirWithFiles := filepath.Join(base, "with_files")
		os.Mkdir(dirWithFiles, 0o755)
		os.WriteFile(filepath.Join(dirWithFiles, "a.txt"), []byte("a"), 0o644)
		os.WriteFile(filepath.Join(dirWithFiles, "b.txt"), []byte("b"), 0o644)

		emptyDir := filepath.Join(base, "empty")
		os.Mkdir(emptyDir, 0o755)

		missingDir := filepath.Join(base, "missing")

		tests := []struct {
			name      string
			input     string
			want      []string
			expectErr error
		}{
			{
				name:      "Valid directory with files",
				input:     dirWithFiles,
				want:      []string{"a.txt", "b.txt"},
				expectErr: nil,
			},
			{
				name:      "Empty directory",
				input:     emptyDir,
				want:      nil,
				expectErr: errors.New("no entries found"),
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
				got, err := ReadDirNames(tt.input, false)

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

func TestReadDirNamesHiddenFiles(t *testing.T) {
	t.Run("Hidden files visibility", func(t *testing.T) {
		// Setup test directory with hidden and non-hidden files
		base := t.TempDir()
		testDir := filepath.Join(base, "test_dir")
		os.Mkdir(testDir, 0o755)

		// Create test files
		os.WriteFile(filepath.Join(testDir, "normal.txt"), []byte("normal"), 0o644)
		os.WriteFile(filepath.Join(testDir, ".hidden.txt"), []byte("hidden"), 0o644)
		os.WriteFile(filepath.Join(testDir, "another.txt"), []byte("another"), 0o644)
		os.WriteFile(filepath.Join(testDir, ".hidden2.txt"), []byte("hidden2"), 0o644)

		tests := []struct {
			name       string
			showHidden bool
			want       []string
			expectErr  error
		}{
			{
				name:       "Show hidden files",
				showHidden: true,
				want:       []string{"normal.txt", ".hidden.txt", "another.txt", ".hidden2.txt"},
				expectErr:  nil,
			},
			{
				name:       "Hide hidden files",
				showHidden: false,
				want:       []string{"normal.txt", "another.txt"},
				expectErr:  nil,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := ReadDirNames(testDir, tt.showHidden)

				if tt.expectErr != nil {
					if err == nil {
						t.Errorf("expected error %v, got nil", tt.expectErr)
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
