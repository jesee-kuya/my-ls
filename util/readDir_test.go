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
