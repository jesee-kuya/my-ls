package util

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsValidDir(t *testing.T) {
	tempDir := t.TempDir()

	// Create a temp file (not a directory)
	tempFile := filepath.Join(tempDir, "file.txt")
	err := os.WriteFile(tempFile, []byte("test"), 0o644)
	if err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	tests := []struct {
		name    string
		input   string
		wantErr bool
		isDir   bool
	}{
		{
			name:    "Valid directory",
			input:   tempDir,
			wantErr: false,
			isDir:   true,
		},
		{
			name:    "Non-existent path",
			input:   filepath.Join(tempDir, "nope"),
			wantErr: true,
		},
		{
			name:    "Existing file but not a directory",
			input:   tempFile,
			wantErr: false,
			isDir:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := IsValidDir(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && info.IsDir() != tt.isDir {
				t.Errorf("IsValidDir() isDir = %v, want %v", info.IsDir(), tt.isDir)
			}
		})
	}
}
