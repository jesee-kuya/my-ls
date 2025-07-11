package main

import (
	"reflect"
	"sort"
	"testing"

	"github.com/jesee-kuya/my-ls/print"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantFlags print.Flags
		wantPaths []string
	}{
		{
			name:      "No arguments",
			args:      []string{},
			wantFlags: print.Flags{ShowHidden: false},
			wantPaths: []string{"."},
		},
		{
			name:      "Only -a flag",
			args:      []string{"-a"},
			wantFlags: print.Flags{ShowHidden: true},
			wantPaths: []string{"."},
		},
		{
			name:      "Flag with path",
			args:      []string{"-a", "/tmp"},
			wantFlags: print.Flags{ShowHidden: true},
			wantPaths: []string{"/tmp"},
		},
		{
			name:      "Path with flag",
			args:      []string{"/tmp", "-a"},
			wantFlags: print.Flags{ShowHidden: true},
			wantPaths: []string{"/tmp"},
		},
		{
			name:      "Multiple paths",
			args:      []string{"dir1", "dir2", "dir3"},
			wantFlags: print.Flags{ShowHidden: false},
			wantPaths: []string{"dir1", "dir2", "dir3"},
		},
		{
			name:      "Multiple paths with flag",
			args:      []string{"-a", "dir1", "dir2"},
			wantFlags: print.Flags{ShowHidden: true},
			wantPaths: []string{"dir1", "dir2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlags, gotPaths := parseArgs(tt.args)

			if !reflect.DeepEqual(gotFlags, tt.wantFlags) {
				t.Errorf("parseArgs() flags = %v, want %v", gotFlags, tt.wantFlags)
			}

			sort.Strings(gotPaths)
			sort.Strings(tt.wantPaths)
			if !reflect.DeepEqual(gotPaths, tt.wantPaths) {
				t.Errorf("parseArgs() paths = %v, want %v", gotPaths, tt.wantPaths)
			}
		})
	}
}

func TestParseArgsEdgeCases(t *testing.T) {
	t.Run("Combined flags (future-proofing)", func(t *testing.T) {
		// Test for when more flags are added later
		flags, paths := parseArgs([]string{"-al", "testdir"})

		if !flags.ShowHidden {
			t.Error("Expected ShowHidden to be true for -al flag")
		}

		if len(paths) != 1 || paths[0] != "testdir" {
			t.Errorf("Expected paths to be [testdir], got %v", paths)
		}
	})

	t.Run("Flag-like path names", func(t *testing.T) {
		// Test directory names that start with dash
		flags, paths := parseArgs([]string{"-a", "--", "-weird-dir-name"})

		if !flags.ShowHidden {
			t.Error("Expected ShowHidden to be true")
		}

		// Note: This test shows current behavior - it treats --weird-dir-name as flags
		// In a real implementation, you might want to handle -- as a flag terminator
		expectedPaths := []string{"."} // Since all args are treated as flags
		if !reflect.DeepEqual(paths, expectedPaths) {
			t.Errorf("Expected paths %v, got %v", expectedPaths, paths)
		}
	})
}
