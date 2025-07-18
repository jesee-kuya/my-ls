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
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Complex ANSI codes",
			input:    "\033[40;33;01mcomplex\033[0m",
			expected: "complex",
		},
		{
			name:     "ANSI codes with numbers",
			input:    "\033[1;2;3;4;5mtext\033[0m",
			expected: "text",
		},
		{
			name:     "Multiple resets",
			input:    "\033[01;34mtext\033[0m\033[0m\033[0m",
			expected: "text",
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

func TestStripLong(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Standard long format",
			input:    "-rw-r--r--  1 user group  100 Jan  1 12:00 file.txt",
			expected: "file.txt",
		},
		{
			name:     "Directory long format",
			input:    "drwxr-xr-x  2 user group  4096 Jan 15 14:30 dirname",
			expected: "dirname",
		},
		{
			name:     "Executable long format",
			input:    "-rwxr-xr-x  1 user group  8192 Feb  5 09:15 executable",
			expected: "executable",
		},
		{
			name:     "No long format prefix",
			input:    "just a filename",
			expected: "just a filename",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Long format with ANSI codes",
			input:    "-rw-r--r--  1 user group  100 Jan  1 12:00 \033[01;34mfile.txt\033[0m",
			expected: "\033[01;34mfile.txt\033[0m",
		},
		{
			name:     "Long format with single digit day",
			input:    "-rw-r--r--  1 user group  100 Jan  5 12:00 file.txt",
			expected: "file.txt",
		},
		{
			name:     "Symlink long format",
			input:    "lrwxrwxrwx  1 user group  10 Mar  1 12:00 symlink",
			expected: "lrwxrwxrwx  1 user group  10 Mar  1 12:00 symlink", // regex doesn't match symlinks
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripLong(tt.input)
			if result != tt.expected {
				t.Errorf("StripLong(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
