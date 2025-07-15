package main

import (
	"reflect"
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
			expectedFlags: util.Flags{ShowAll: false, Longformat: false},
			expectedPaths: []string{"."},
		},
		{
			name:          "only -a flag",
			args:          []string{"-a"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: false},
			expectedPaths: []string{"."},
		},
		{
			name:          "only path",
			args:          []string{"/tmp"},
			expectedFlags: util.Flags{ShowAll: false, Longformat: false},
			expectedPaths: []string{"/tmp"},
		},
		{
			name:          "-a flag with path",
			args:          []string{"-a", "/tmp"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: false},
			expectedPaths: []string{"/tmp"},
		},
		{
			name:          "path with -a flag",
			args:          []string{"/tmp", "-a"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: false},
			expectedPaths: []string{"/tmp"},
		},
		{
			name:          "multiple paths with -a flag",
			args:          []string{"-a", "/tmp", "/home"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: false},
			expectedPaths: []string{"/tmp", "/home"},
		},
		{
			name:          "multiple flags including -a",
			args:          []string{"-al"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true},
			expectedPaths: []string{"."},
		},
		{
			name:          "multiple flags with -a and paths",
			args:          []string{"-al", "/tmp", "/home"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true},
			expectedPaths: []string{"/tmp", "/home"},
		},
		{
			name:          "flags without -a",
			args:          []string{"-l"},
			expectedFlags: util.Flags{ShowAll: false, Longformat: true},
			expectedPaths: []string{"."},
		},
		{
			name:          "mixed flags and paths",
			args:          []string{"/tmp", "-a", "/home", "-l"},
			expectedFlags: util.Flags{ShowAll: true, Longformat: true},
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
		expectedFlags := util.Flags{ShowAll: false, Longformat: false}
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
		expectedFlags := util.Flags{ShowAll: false, Longformat: false}
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
		expectedFlags := util.Flags{ShowAll: true, Longformat: false}
		expectedPaths := []string{"."}

		if !reflect.DeepEqual(flags, expectedFlags) {
			t.Errorf("parseArgs() flags = %v, want %v", flags, expectedFlags)
		}

		if !reflect.DeepEqual(paths, expectedPaths) {
			t.Errorf("parseArgs() paths = %v, want %v", paths, expectedPaths)
		}
	})
}
