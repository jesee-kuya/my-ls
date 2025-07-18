package util

import (
	"testing"
	"time"
)

func TestFormatTime(t *testing.T) {
	tests := []struct {
		name     string
		input    time.Time
		expected string
	}{
		{
			name:     "January date",
			input:    time.Date(2023, time.January, 15, 14, 30, 0, 0, time.UTC),
			expected: "Jan 15 14:30",
		},
		{
			name:     "February date",
			input:    time.Date(2023, time.February, 5, 9, 15, 0, 0, time.UTC),
			expected: "Feb  5 09:15",
		},
		{
			name:     "December date",
			input:    time.Date(2023, time.December, 25, 23, 59, 0, 0, time.UTC),
			expected: "Dec 25 23:59",
		},
		{
			name:     "Single digit day",
			input:    time.Date(2023, time.March, 1, 12, 0, 0, 0, time.UTC),
			expected: "Mar  1 12:00",
		},
		{
			name:     "Double digit day",
			input:    time.Date(2023, time.April, 10, 8, 45, 0, 0, time.UTC),
			expected: "Apr 10 08:45",
		},
		{
			name:     "Midnight",
			input:    time.Date(2023, time.May, 20, 0, 0, 0, 0, time.UTC),
			expected: "May 20 00:00",
		},
		{
			name:     "Noon",
			input:    time.Date(2023, time.June, 15, 12, 0, 0, 0, time.UTC),
			expected: "Jun 15 12:00",
		},
		{
			name:     "Late evening",
			input:    time.Date(2023, time.July, 4, 23, 30, 0, 0, time.UTC),
			expected: "Jul  4 23:30",
		},
		{
			name:     "Early morning",
			input:    time.Date(2023, time.August, 8, 6, 15, 0, 0, time.UTC),
			expected: "Aug  8 06:15",
		},
		{
			name:     "September date",
			input:    time.Date(2023, time.September, 30, 18, 45, 0, 0, time.UTC),
			expected: "Sep 30 18:45",
		},
		{
			name:     "October date",
			input:    time.Date(2023, time.October, 31, 21, 30, 0, 0, time.UTC),
			expected: "Oct 31 21:30",
		},
		{
			name:     "November date",
			input:    time.Date(2023, time.November, 11, 11, 11, 0, 0, time.UTC),
			expected: "Nov 11 11:11",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatTime(tt.input)
			if result != tt.expected {
				t.Errorf("FormatTime(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFormatTime_ConsistentFormat(t *testing.T) {
	// Test that the format is consistent with Go's time formatting
	testTime := time.Date(2023, time.January, 15, 14, 30, 0, 0, time.UTC)

	result := FormatTime(testTime)
	expected := testTime.Format("Jan _2 15:04")

	if result != expected {
		t.Errorf("FormatTime() format inconsistent with Go's time.Format(), got %q, want %q", result, expected)
	}
}

func TestFormatTime_EdgeCases(t *testing.T) {
	// Test with zero time
	zeroTime := time.Time{}
	result := FormatTime(zeroTime)
	expected := zeroTime.Format("Jan _2 15:04")

	if result != expected {
		t.Errorf("FormatTime() with zero time = %q, want %q", result, expected)
	}

	// Test with Unix epoch
	epochTime := time.Unix(0, 0).UTC()
	result = FormatTime(epochTime)
	expected = epochTime.Format("Jan _2 15:04")

	if result != expected {
		t.Errorf("FormatTime() with Unix epoch = %q, want %q", result, expected)
	}
}

func TestFormatTime_DifferentTimezones(t *testing.T) {
	// Test that the function works with different timezones
	baseTime := time.Date(2023, time.January, 15, 14, 30, 0, 0, time.UTC)

	// Convert to different timezone
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Skip("Could not load timezone for test")
	}

	nyTime := baseTime.In(loc)
	result := FormatTime(nyTime)
	expected := nyTime.Format("Jan _2 15:04")

	if result != expected {
		t.Errorf("FormatTime() with timezone = %q, want %q", result, expected)
	}
}
