package util

import (
	"reflect"
	"testing"
)

func TestReverse(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "single element",
			input:    []string{"a"},
			expected: []string{"a"},
		},
		{
			name:     "two elements",
			input:    []string{"a", "b"},
			expected: []string{"b", "a"},
		},
		{
			name:     "three elements",
			input:    []string{"a", "b", "c"},
			expected: []string{"c", "b", "a"},
		},
		{
			name:     "four elements",
			input:    []string{"a", "b", "c", "d"},
			expected: []string{"d", "c", "b", "a"},
		},
		{
			name:     "five elements",
			input:    []string{"first", "second", "third", "fourth", "fifth"},
			expected: []string{"fifth", "fourth", "third", "second", "first"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy of input to avoid modifying the test data
			input := make([]string, len(tt.input))
			copy(input, tt.input)

			Reverse(input)

			if !reflect.DeepEqual(input, tt.expected) {
				t.Errorf("Reverse() = %v, want %v", input, tt.expected)
			}
		})
	}
}

func TestReverse_IntSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    []int
		expected []int
	}{
		{
			name:     "empty int slice",
			input:    []int{},
			expected: []int{},
		},
		{
			name:     "single int",
			input:    []int{1},
			expected: []int{1},
		},
		{
			name:     "multiple ints",
			input:    []int{1, 2, 3, 4, 5},
			expected: []int{5, 4, 3, 2, 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make a copy of input to avoid modifying the test data
			input := make([]int, len(tt.input))
			copy(input, tt.input)

			Reverse(input)

			if !reflect.DeepEqual(input, tt.expected) {
				t.Errorf("Reverse() = %v, want %v", input, tt.expected)
			}
		})
	}
}

func TestReverse_ModifiesOriginalSlice(t *testing.T) {
	// Test that Reverse modifies the original slice in place
	original := []string{"a", "b", "c"}
	expected := []string{"c", "b", "a"}

	Reverse(original)

	if !reflect.DeepEqual(original, expected) {
		t.Errorf("Reverse() should modify original slice in place, got %v, want %v", original, expected)
	}
}

func TestReverse_GenericTypes(t *testing.T) {
	// Test with different types to ensure generic functionality works

	// Test with float64
	floats := []float64{1.1, 2.2, 3.3}
	expectedFloats := []float64{3.3, 2.2, 1.1}
	Reverse(floats)
	if !reflect.DeepEqual(floats, expectedFloats) {
		t.Errorf("Reverse() with float64 = %v, want %v", floats, expectedFloats)
	}

	// Test with bool
	bools := []bool{true, false, true}
	expectedBools := []bool{true, false, true}
	Reverse(bools)
	if !reflect.DeepEqual(bools, expectedBools) {
		t.Errorf("Reverse() with bool = %v, want %v", bools, expectedBools)
	}
}
