package main

import (
	"reflect"
	"testing"

	"github.com/jesee-kuya/my-ls/print"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedFlags print.Flags
		expectedPaths []string
	}{
		{
			name:          "no arguments",
			args:          []string{},
			expectedFlags: print.Flags{ShowAll: false},
			expectedPaths: []string{"."},
		},
		{
			name:          "only -a flag",
			args:          []string{"-a"},
			expectedFlags: print.Flags{ShowAll: true},
			expectedPaths: []string{"."},
		},
		{
			name:          "only path",
			args:          []string{"/tmp"},
			expectedFlags: print.Flags{ShowAll: false},
			expectedPaths: []string{"/tmp"},
		},
		{
			name:          "-a flag with path",
			args:          []string{"-a", "/tmp"},
			expectedFlags: print.Flags{ShowAll: true},
			expectedPaths: []string{"/tmp"},
		},
		{
			name:          "path with -a flag",
			args:          []string{"/tmp", "-a"},
			expectedFlags: print.Flags{ShowAll: true},
			expectedPaths: []string{"/tmp"},
		},
		{
			name:          "multiple paths with -a flag",
			args:          []string{"-a", "/tmp", "/home"},
			expectedFlags: print.Flags{ShowAll: true},
			expectedPaths: []string{"/tmp", "/home"},
		},
		{
			name:          "multiple flags including -a",
			args:          []string{"-al"},
			expectedFlags: print.Flags{ShowAll: true},
			expectedPaths: []string{"."},
		},
		{
			name:          "multiple flags with -a and paths",
			args:          []string{"-al", "/tmp", "/home"},
			expectedFlags: print.Flags{ShowAll: true},
			expectedPaths: []string{"/tmp", "/home"},
		},
		{
			name:          "flags without -a",
			args:          []string{"-l"},
			expectedFlags: print.Flags{ShowAll: false},
			expectedPaths: []string{"."},
		},
		{
			name:          "mixed flags and paths",
			args:          []string{"/tmp", "-a", "/home", "-l"},
			expectedFlags: print.Flags{ShowAll: true},
			expectedPaths: []string{"/tmp", "/home"},
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
		expectedFlags := print.Flags{ShowAll: false}
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
		expectedFlags := print.Flags{ShowAll: false}
		expectedPaths := []string{"."}

		if !reflect.DeepEqual(flags, expectedFlags) {
			t.Errorf("parseArgs() flags = %v, want %v", flags, expectedFlags)
		}

		if !reflect.DeepEqual(paths, expectedPaths) {
			t.Errorf("parseArgs() paths = %v, want %v", paths, expectedPaths)
		}
	})

	t.Run("flag with unknown characters", func(t *testing.T) {
		flags, paths := parseArgs([]string{"-axl"})
		expectedFlags := print.Flags{ShowAll: true} // should still parse 'a'
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
		expectedFlags := print.Flags{ShowAll: true}
		expectedPaths := []string{"."}

		if !reflect.DeepEqual(flags, expectedFlags) {
			t.Errorf("parseArgs() flags = %v, want %v", flags, expectedFlags)
		}

		if !reflect.DeepEqual(paths, expectedPaths) {
			t.Errorf("parseArgs() paths = %v, want %v", paths, expectedPaths)
		}
	})
}
