package util

import "testing"

func TestHasANSIPrefix(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Starts with bold blue ANSI",
			input:    "\033[01;34mdir",
			expected: true,
		},
		{
			name:     "Starts with reset ANSI",
			input:    "\033[0mfile",
			expected: true,
		},
		{
			name:     "Starts with non-ANSI character",
			input:    "normal.txt",
			expected: false,
		},
		{
			name:     "Contains ANSI but not at start",
			input:    "log\033[01;31merror\033[0m",
			expected: false,
		},
		{
			name:     "Multiple ANSI, first at start",
			input:    "\033[01;35m\033[0mscript.sh",
			expected: true,
		},
		{
			name:     "Empty string",
			input:    "",
			expected: false,
		},
		{
			name:     "Malformed escape sequence",
			input:    "\033[XYZdir",
			expected: false,
		},
		{
			name:     "Escaped manually typed sequence",
			input:    "\\033[01;34mdir",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HasANSIPrefix(tt.input)
			if result != tt.expected {
				t.Errorf("HasANSIPrefix(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
