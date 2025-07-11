package util

import (
	"testing"
)

func TestStripANSI(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No ANSI codes",
			input:    "hello.txt",
			expected: "hello.txt",
		},
		{
			name:     "Bold blue directory",
			input:    "\033[01;34mdir\033[0m",
			expected: "dir",
		},
		{
			name:     "Bold green executable",
			input:    "\033[01;32mrun.sh\033[0m",
			expected: "run.sh",
		},
		{
			name:     "Mixed ANSI codes and text",
			input:    "file-\033[01;31merror\033[0m.log",
			expected: "file-error.log",
		},
		{
			name:     "Multiple ANSI sequences",
			input:    "\033[01;34mdir\033[0m/\033[01;32mfile\033[0m",
			expected: "dir/file",
		},
		{
			name:     "Only ANSI codes",
			input:    "\033[01;35m\033[0m",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := StripANSI(tt.input)
			if output != tt.expected {
				t.Errorf("StripANSI(%q) = %q; want %q", tt.input, output, tt.expected)
			}
		})
	}
}
