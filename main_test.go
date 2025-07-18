package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/jesee-kuya/my-ls/util"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedFlags util.Flags
		expectedPaths []string
	}{
		{
			name:          "no arguments",
			args:          []string{},
			expectedFlags: util.Flags{ShowAll: false, Longformat: false, Reverse: false, Recursive: false, TimeSort: false},
			expectedPaths: []string{"."},
		},
		{
			name:          "only -a flag",
			args:          []string{"-a"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: false, Reverse: false, Recursive: false, TimeSort: false},
			expectedPaths: []string{"."},
		},
		{
			name:          "only path",
			args:          []string{"/tmp"},
			expectedFlags: util.Flags{ShowAll: false, Longformat: false, Reverse: false, Recursive: false, TimeSort: false},
			expectedPaths: []string{"/tmp"},
		},
		{
			name:          "-a flag with path",
			args:          []string{"-a", "/tmp"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: false, Reverse: false, Recursive: false, TimeSort: false},
			expectedPaths: []string{"/tmp"},
		},
		{
			name:          "path with -a flag",
			args:          []string{"/tmp", "-a"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: false, Reverse: false, Recursive: false, TimeSort: false},
			expectedPaths: []string{"/tmp"},
		},
		{
			name:          "multiple paths with -a flag",
			args:          []string{"-a", "/tmp", "/home"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: false, Reverse: false, Recursive: false, TimeSort: false},
			expectedPaths: []string{"/tmp", "/home"},
		},
		{
			name:          "multiple flags including -a",
			args:          []string{"-al"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: false, Recursive: false, TimeSort: false},
			expectedPaths: []string{"."},
		},
		{
			name:          "multiple flags with -a and paths",
			args:          []string{"-al", "/tmp", "/home"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: false, Recursive: false, TimeSort: false},
			expectedPaths: []string{"/tmp", "/home"},
		},
		{
			name:          "flags without -a",
			args:          []string{"-l"},
			expectedFlags: util.Flags{ShowAll: false, Longformat: true, Reverse: false, Recursive: false, TimeSort: false},
			expectedPaths: []string{"."},
		},
		{
			name:          "mixed flags and paths",
			args:          []string{"/tmp", "-a", "/home", "-l"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: false, Recursive: false, TimeSort: false},
			expectedPaths: []string{"/tmp", "/home"},
		},
		{
			name:          "only -R flag",
			args:          []string{"-R"},
			expectedFlags: util.Flags{ShowAll: false, Longformat: false, Reverse: false, Recursive: true, TimeSort: false},
			expectedPaths: []string{"."},
		},
		{
			name:          "-R flag with path",
			args:          []string{"-R", "/tmp"},
			expectedFlags: util.Flags{ShowAll: false, Longformat: false, Reverse: false, Recursive: true, TimeSort: false},
			expectedPaths: []string{"/tmp"},
		},
		{
			name:          "combined flags with -R",
			args:          []string{"-alR"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: false, Recursive: true, TimeSort: false},
			expectedPaths: []string{"."},
		},
		{
			name:          "all flags combined",
			args:          []string{"-alrR"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: true, Recursive: true, TimeSort: false},
			expectedPaths: []string{"."},
		},
		{
			name:          "only -t flag",
			args:          []string{"-t"},
			expectedFlags: util.Flags{ShowAll: false, Longformat: false, Reverse: false, Recursive: false, TimeSort: true},
			expectedPaths: []string{"."},
		},
		{
			name:          "-t flag with path",
			args:          []string{"-t", "/tmp"},
			expectedFlags: util.Flags{ShowAll: false, Longformat: false, Reverse: false, Recursive: false, TimeSort: true},
			expectedPaths: []string{"/tmp"},
		},
		{
			name:          "combined flags with -t",
			args:          []string{"-alt"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: false, Recursive: false, TimeSort: true},
			expectedPaths: []string{"."},
		},
		{
			name:          "all flags including -t",
			args:          []string{"-alrRt"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: true, Recursive: true, TimeSort: true},
			expectedPaths: []string{"."},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags, paths := parseArgs(tt.args)

			if !reflect.DeepEqual(flags, tt.expectedFlags) {
				t.Errorf("parseArgs() flags = %v, want %v", flags, tt.expectedFlags)
			}

			if !reflect.DeepEqual(paths, tt.expectedPaths) {
				t.Errorf("parseArgs() paths = %v, want %v", paths, tt.expectedPaths)
			}
		})
	}
}

func TestParseArgs_EdgeCases(t *testing.T) {
	t.Run("empty flag", func(t *testing.T) {
		flags, paths := parseArgs([]string{"-"})
		expectedFlags := util.Flags{ShowAll: false, Longformat: false, Reverse: false, Recursive: false, TimeSort: false}
		expectedPaths := []string{"."}

		if !reflect.DeepEqual(flags, expectedFlags) {
			t.Errorf("parseArgs() flags = %v, want %v", flags, expectedFlags)
		}

		if !reflect.DeepEqual(paths, expectedPaths) {
			t.Errorf("parseArgs() paths = %v, want %v", paths, expectedPaths)
		}
	})

	t.Run("unknown flag", func(t *testing.T) {
		flags, paths := parseArgs([]string{"-x"})
		expectedFlags := util.Flags{ShowAll: false, Longformat: false, Reverse: false, Recursive: false, TimeSort: false}
		expectedPaths := []string{"."}

		if !reflect.DeepEqual(flags, expectedFlags) {
			t.Errorf("parseArgs() flags = %v, want %v", flags, expectedFlags)
		}

		if !reflect.DeepEqual(paths, expectedPaths) {
			t.Errorf("parseArgs() paths = %v, want %v", paths, expectedPaths)
		}
	})

	t.Run("multiple -a flags", func(t *testing.T) {
		flags, paths := parseArgs([]string{"-a", "-a"})
		expectedFlags := util.Flags{ShowAll: true, Longformat: false, Reverse: false, Recursive: false, TimeSort: false}
		expectedPaths := []string{"."}

		if !reflect.DeepEqual(flags, expectedFlags) {
			t.Errorf("parseArgs() flags = %v, want %v", flags, expectedFlags)
		}

		if !reflect.DeepEqual(paths, expectedPaths) {
			t.Errorf("parseArgs() paths = %v, want %v", paths, expectedPaths)
		}
	})
}

// TestParseArgs_CombinedVsSeparateFlags tests that combined flags (-la) and separate flags (-l -a)
// produce identical results, ensuring Unix convention compatibility
func TestParseArgs_CombinedVsSeparateFlags(t *testing.T) {
	testCases := []struct {
		name          string
		combined      []string
		separate      []string
		expectedFlags util.Flags
		expectedPaths []string
	}{
		{
			name:          "la flags",
			combined:      []string{"-la"},
			separate:      []string{"-l", "-a"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: false, Recursive: false, TimeSort: false},
			expectedPaths: []string{"."},
		},
		{
			name:          "al flags (different order)",
			combined:      []string{"-al"},
			separate:      []string{"-a", "-l"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: false, Recursive: false, TimeSort: false},
			expectedPaths: []string{"."},
		},
		{
			name:          "lar flags",
			combined:      []string{"-lar"},
			separate:      []string{"-l", "-a", "-r"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: true, Recursive: false, TimeSort: false},
			expectedPaths: []string{"."},
		},
		{
			name:          "lart flags (all flags)",
			combined:      []string{"-lart"},
			separate:      []string{"-l", "-a", "-r", "-t"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: true, Recursive: false, TimeSort: true},
			expectedPaths: []string{"."},
		},
		{
			name:          "laRt flags with recursive",
			combined:      []string{"-laRt"},
			separate:      []string{"-l", "-a", "-R", "-t"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: false, Recursive: true, TimeSort: true},
			expectedPaths: []string{"."},
		},
		{
			name:          "all flags combined vs separate",
			combined:      []string{"-alrRt"},
			separate:      []string{"-a", "-l", "-r", "-R", "-t"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: true, Recursive: true, TimeSort: true},
			expectedPaths: []string{"."},
		},
		{
			name:          "flags with path - combined",
			combined:      []string{"-la", "/tmp"},
			separate:      []string{"-l", "-a", "/tmp"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: false, Recursive: false, TimeSort: false},
			expectedPaths: []string{"/tmp"},
		},
		{
			name:          "flags with multiple paths",
			combined:      []string{"-la", "/tmp", "/home"},
			separate:      []string{"-l", "-a", "/tmp", "/home"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: false, Recursive: false, TimeSort: false},
			expectedPaths: []string{"/tmp", "/home"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test combined format
			combinedFlags, combinedPaths := parseArgs(tc.combined)
			if !reflect.DeepEqual(combinedFlags, tc.expectedFlags) {
				t.Errorf("Combined format: parseArgs(%v) flags = %v, want %v", tc.combined, combinedFlags, tc.expectedFlags)
			}
			if !reflect.DeepEqual(combinedPaths, tc.expectedPaths) {
				t.Errorf("Combined format: parseArgs(%v) paths = %v, want %v", tc.combined, combinedPaths, tc.expectedPaths)
			}

			// Test separate format
			separateFlags, separatePaths := parseArgs(tc.separate)
			if !reflect.DeepEqual(separateFlags, tc.expectedFlags) {
				t.Errorf("Separate format: parseArgs(%v) flags = %v, want %v", tc.separate, separateFlags, tc.expectedFlags)
			}
			if !reflect.DeepEqual(separatePaths, tc.expectedPaths) {
				t.Errorf("Separate format: parseArgs(%v) paths = %v, want %v", tc.separate, separatePaths, tc.expectedPaths)
			}

			// Most importantly, verify that combined and separate produce identical results
			if !reflect.DeepEqual(combinedFlags, separateFlags) {
				t.Errorf("Flag mismatch: combined %v != separate %v", combinedFlags, separateFlags)
			}
			if !reflect.DeepEqual(combinedPaths, separatePaths) {
				t.Errorf("Path mismatch: combined %v != separate %v", combinedPaths, separatePaths)
			}
		})
	}
}

// TestParseArgs_MixedFormats tests mixed flag formats like "-la -r" or "-l -ar"
func TestParseArgs_MixedFormats(t *testing.T) {
	testCases := []struct {
		name          string
		args          []string
		expectedFlags util.Flags
		expectedPaths []string
	}{
		{
			name:          "mixed: -la -r",
			args:          []string{"-la", "-r"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: true, Recursive: false, TimeSort: false},
			expectedPaths: []string{"."},
		},
		{
			name:          "mixed: -l -ar",
			args:          []string{"-l", "-ar"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: true, Recursive: false, TimeSort: false},
			expectedPaths: []string{"."},
		},
		{
			name:          "mixed: -a -lrt",
			args:          []string{"-a", "-lrt"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: true, Recursive: false, TimeSort: true},
			expectedPaths: []string{"."},
		},
		{
			name:          "mixed: -al -R -t",
			args:          []string{"-al", "-R", "-t"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: false, Recursive: true, TimeSort: true},
			expectedPaths: []string{"."},
		},
		{
			name:          "mixed with paths: -la /tmp -r /home",
			args:          []string{"-la", "/tmp", "-r", "/home"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: true, Recursive: false, TimeSort: false},
			expectedPaths: []string{"/tmp", "/home"},
		},
		{
			name:          "scattered flags: -l /tmp -a /home -r",
			args:          []string{"-l", "/tmp", "-a", "/home", "-r"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true, Reverse: true, Recursive: false, TimeSort: false},
			expectedPaths: []string{"/tmp", "/home"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			flags, paths := parseArgs(tc.args)

			if !reflect.DeepEqual(flags, tc.expectedFlags) {
				t.Errorf("parseArgs(%v) flags = %v, want %v", tc.args, flags, tc.expectedFlags)
			}

			if !reflect.DeepEqual(paths, tc.expectedPaths) {
				t.Errorf("parseArgs(%v) paths = %v, want %v", tc.args, paths, tc.expectedPaths)
			}
		})
	}
}

// TestFlagFormatEquivalence tests that combined and separate flag formats produce identical output
func TestFlagFormatEquivalence(t *testing.T) {
	// Create a temporary directory with test files for consistent testing
	tempDir := t.TempDir()

	// Create test files
	testFiles := []string{"file1.txt", "file2.txt", ".hidden", "dir1"}
	for _, file := range testFiles[:3] {
		err := os.WriteFile(filepath.Join(tempDir, file), []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Create test directory
	err := os.Mkdir(filepath.Join(tempDir, "dir1"), 0755)
	if err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	testCases := []struct {
		name     string
		combined []string
		separate []string
	}{
		{
			name:     "la flags",
			combined: []string{"-la", tempDir},
			separate: []string{"-l", "-a", tempDir},
		},
		{
			name:     "lar flags",
			combined: []string{"-lar", tempDir},
			separate: []string{"-l", "-a", "-r", tempDir},
		},
		{
			name:     "lat flags",
			combined: []string{"-lat", tempDir},
			separate: []string{"-l", "-a", "-t", tempDir},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Parse arguments for both formats
			combinedFlags, combinedPaths := parseArgs(tc.combined)
			separateFlags, separatePaths := parseArgs(tc.separate)

			// Verify flags are identical
			if !reflect.DeepEqual(combinedFlags, separateFlags) {
				t.Errorf("Flags differ: combined=%v, separate=%v", combinedFlags, separateFlags)
			}

			// Verify paths are identical
			if !reflect.DeepEqual(combinedPaths, separatePaths) {
				t.Errorf("Paths differ: combined=%v, separate=%v", combinedPaths, separatePaths)
			}

			// Test that the actual output would be identical by calling the same functions
			// This ensures the entire pipeline works consistently
			if combinedFlags.Longformat {
				combinedOutput, err1 := util.ReadDirNamesLong(tempDir, combinedFlags)
				separateOutput, err2 := util.ReadDirNamesLong(tempDir, separateFlags)

				if err1 != err2 {
					t.Errorf("Errors differ: combined=%v, separate=%v", err1, err2)
				}

				if err1 == nil && !reflect.DeepEqual(combinedOutput, separateOutput) {
					t.Errorf("Long format outputs differ")
				}
			} else {
				combinedOutput, err1 := util.ReadDirNames(tempDir, combinedFlags)
				separateOutput, err2 := util.ReadDirNames(tempDir, separateFlags)

				if err1 != err2 {
					t.Errorf("Errors differ: combined=%v, separate=%v", err1, err2)
				}

				if err1 == nil && !reflect.DeepEqual(combinedOutput, separateOutput) {
					t.Errorf("Short format outputs differ")
				}
			}
		})
	}
}

// captureOutput captures stdout during function execution
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestMain_NoArgs(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Change to the temp directory
	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldDir)

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Save original os.Args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test main with no arguments
	os.Args = []string{"my-ls"}

	output := captureOutput(func() {
		main()
	})

	if !strings.Contains(output, "test.txt") {
		t.Errorf("Expected output to contain 'test.txt', got: %s", output)
	}
}

func TestMain_WithArgs(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create test files
	testFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	hiddenFile := filepath.Join(tempDir, ".hidden")
	err = os.WriteFile(hiddenFile, []byte("hidden content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create hidden file: %v", err)
	}

	// Save original os.Args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test main with -a flag
	os.Args = []string{"my-ls", "-a", tempDir}

	output := captureOutput(func() {
		main()
	})

	if !strings.Contains(output, "test.txt") {
		t.Errorf("Expected output to contain 'test.txt', got: %s", output)
	}

	if !strings.Contains(output, ".hidden") {
		t.Errorf("Expected output to contain '.hidden' with -a flag, got: %s", output)
	}
}

func TestMain_LongFormat(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Save original os.Args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test main with -l flag
	os.Args = []string{"my-ls", "-l", tempDir}

	output := captureOutput(func() {
		main()
	})

	if !strings.Contains(output, "test.txt") {
		t.Errorf("Expected output to contain 'test.txt', got: %s", output)
	}

	if !strings.Contains(output, "total") {
		t.Errorf("Expected output to contain 'total' in long format, got: %s", output)
	}
}

func TestMain_CombinedFlags(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create test files
	testFile := filepath.Join(tempDir, "test.txt")
	err := os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	hiddenFile := filepath.Join(tempDir, ".hidden")
	err = os.WriteFile(hiddenFile, []byte("hidden content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create hidden file: %v", err)
	}

	// Save original os.Args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test main with combined flags
	os.Args = []string{"my-ls", "-la", tempDir}

	output := captureOutput(func() {
		main()
	})

	if !strings.Contains(output, "test.txt") {
		t.Errorf("Expected output to contain 'test.txt', got: %s", output)
	}

	if !strings.Contains(output, ".hidden") {
		t.Errorf("Expected output to contain '.hidden' with -a flag, got: %s", output)
	}

	if !strings.Contains(output, "total") {
		t.Errorf("Expected output to contain 'total' in long format, got: %s", output)
	}
}
