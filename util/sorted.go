package util

import (
	"fmt"
	"os"
	"strings"
)

// CompareStrings compares two strings based on custom sorting rules.
// Significant parts (numeric and alphabetic characters) are compared first, with numeric characters before alphabetic.
// Numeric characters (0-9) are compared by ASCII values, alphabetic characters (a-z, A-Z) case-insensitively with lowercase prioritized.
// If significant parts are equal, compare original strings character by character:
// - Numeric vs. anything: numeric comes first.
// - Alphabetic vs. alphabetic: lowercase before uppercase, case-insensitive otherwise.
// - Other pairs (special vs. special, alphabetic vs. special): ASCII order.
func CompareStrings(a, b string) bool {
	// Convert strings to runes for proper handling
	ra, rb := []rune(a), []rune(b)

	// Check if both strings contain only special (non-alphabetic, non-numeric) characters
	hasSignificantA, hasSignificantB := false, false
	for _, r := range ra {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			hasSignificantA = true
			break
		}
	}
	for _, r := range rb {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			hasSignificantB = true
			break
		}
	}

	// If both strings have only special characters (no letters or digits), compare by ASCII values
	if !hasSignificantA && !hasSignificantB {
		for i := 0; i < len(ra) && i < len(rb); i++ {
			if ra[i] != rb[i] {
				return ra[i] < rb[i]
			}
		}
		// If equal up to shorter length, shorter string comes first
		return len(ra) < len(rb)
	}

	// If either string has significant characters (letters or digits), compare significant parts first
	var sigA, sigB []rune
	for _, r := range ra {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			sigA = append(sigA, []rune(strings.ToLower(string(r)))[0])
		} else if r >= '0' && r <= '9' {
			sigA = append(sigA, r)
		}
	}
	for _, r := range rb {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			sigB = append(sigB, []rune(strings.ToLower(string(r)))[0])
		} else if r >= '0' && r <= '9' {
			sigB = append(sigB, r)
		}
	}

	// Compare significant parts (numeric before alphabetic, alphabetic case-insensitive)
	for i := 0; i < len(sigA) && i < len(sigB); i++ {
		ca, cb := sigA[i], sigB[i]
		if ca == cb {
			continue
		}
		// If one is numeric and the other alphabetic, numeric comes first
		if (ca >= '0' && ca <= '9') && ((cb >= 'a' && cb <= 'z') || (cb >= 'A' && cb <= 'Z')) {
			return true
		}
		if ((ca >= 'a' && ca <= 'z') || (ca >= 'A' && ca <= 'Z')) && (cb >= '0' && cb <= '9') {
			return false
		}
		// If both are alphabetic or both numeric, compare directly
		return ca < cb
	}
	// If significant parts are equal, shorter significant part comes first
	if len(sigA) != len(sigB) {
		return len(sigA) < len(sigB)
	}

	// If significant parts are equal, compare original strings character by character
	for i := 0; i < len(ra) && i < len(rb); i++ {
		ca, cb := ra[i], rb[i]
		if ca == cb {
			continue
		}
		// Prioritize numeric characters over all others
		if (ca >= '0' && ca <= '9') && !(cb >= '0' && cb <= '9') {
			return true
		}
		if !(ca >= '0' && ca <= '9') && (cb >= '0' && cb <= '9') {
			return false
		}
		if (ca >= '0' && ca <= '9') && (cb >= '0' && cb <= '9') {
			return ca < cb
		}
		// If both are alphabetic, prioritize lowercase
		if ((ca >= 'a' && ca <= 'z') || (ca >= 'A' && ca <= 'Z')) && ((cb >= 'a' && cb <= 'z') || (cb >= 'A' && cb <= 'Z')) {
			caLower := []rune(strings.ToLower(string(ca)))[0]
			cbLower := []rune(strings.ToLower(string(cb)))[0]
			if caLower != cbLower {
				return caLower < cbLower
			}
			// If lowercase versions are equal, lowercase comes before uppercase
			if (ca >= 'a' && ca <= 'z') && (cb >= 'A' && cb <= 'Z') {
				return true
			}
			if (ca >= 'A' && ca <= 'Z') && (cb >= 'a' && cb <= 'z') {
				return false
			}
			continue
		}
		// For any other comparison (special vs. special or alphabetic vs. special), use ASCII
		return ca < cb
	}
	// If equal up to shorter length, shorter string comes first
	return len(ra) < len(rb)
}

func InsertSorted(name, colour, reset string, names []string) []string {
	if name == "." || name == ".." {
		return append([]string{fmt.Sprintf("%s%s%s", colour, name, reset)}, names...)
	}
	colored := fmt.Sprintf("%s%s%s", colour, name, reset)

	for i, val := range names {
		cleanVal := TrimStart(StripANSI(val))
		if (CompareStrings(TrimStart(name), cleanVal) || (TrimStart(name) == cleanVal)) && (val != "." && val != "..") {
			return append(names[:i], append([]string{colored}, names[i:]...)...)
		}
	}

	return append(names, colored)
}

func InsertSortedLong(line string, lines []string) []string {
	for i, val := range lines {
		lineName := TrimStart(StripANSI(StripLong(line)))
		valName := TrimStart(StripANSI(StripLong(val)))
		if CompareStrings(lineName, valName) || (lineName == valName) {
			return append(lines[:i], append([]string{line}, lines[i:]...)...)
		}
	}

	return append(lines, line)
}

func TrimStart(name string) string {
	return strings.TrimLeft(name, ".")
}

// InsertSortedByTime inserts a file into a list sorted by modification time (newest first)
func InsertSortedByTime(name, colour, reset, dirPath string, names []string) []string {
	if name == "." || name == ".." {
		return append([]string{fmt.Sprintf("%s%s%s", colour, name, reset)}, names...)
	}
	colored := fmt.Sprintf("%s%s%s", colour, name, reset)

	// Get modification time for the new file
	newFilePath := joinPath(dirPath, name)
	newInfo, err := os.Stat(newFilePath)
	if err != nil {
		// If we can't get the time, fall back to alphabetical sorting
		return InsertSorted(name, colour, reset, names)
	}
	newTime := newInfo.ModTime()

	for i, val := range names {
		cleanVal := StripANSI(val)
		if cleanVal == "." || cleanVal == ".." {
			continue
		}

		// Get modification time for the existing file
		existingFilePath := joinPath(dirPath, cleanVal)
		existingInfo, err := os.Stat(existingFilePath)
		if err != nil {
			// If we can't get the time, insert here
			return append(names[:i], append([]string{colored}, names[i:]...)...)
		}
		existingTime := existingInfo.ModTime()

		// Insert if new file is newer (or same time, then alphabetical)
		if newTime.After(existingTime) || (newTime.Equal(existingTime) && CompareStrings(name, cleanVal)) {
			return append(names[:i], append([]string{colored}, names[i:]...)...)
		}
	}

	return append(names, colored)
}

// InsertSortedLongByTime inserts a long format line sorted by modification time (newest first)
func InsertSortedLongByTime(line, dirPath string, lines []string) []string {
	fileName := StripANSI(StripLong(line))

	// Get modification time for the new file
	newFilePath := joinPath(dirPath, fileName)
	newInfo, err := os.Stat(newFilePath)
	if err != nil {
		// If we can't get the time, fall back to alphabetical sorting
		return InsertSortedLong(line, lines)
	}
	newTime := newInfo.ModTime()

	for i, val := range lines {
		existingFileName := StripANSI(StripLong(val))

		// Get modification time for the existing file
		existingFilePath := joinPath(dirPath, existingFileName)
		existingInfo, err := os.Stat(existingFilePath)
		if err != nil {
			// If we can't get the time, insert here
			return append(lines[:i], append([]string{line}, lines[i:]...)...)
		}
		existingTime := existingInfo.ModTime()

		// Insert if new file is newer (or same time, then alphabetical)
		if newTime.After(existingTime) || (newTime.Equal(existingTime) && CompareStrings(fileName, existingFileName)) {
			return append(lines[:i], append([]string{line}, lines[i:]...)...)
		}
	}

	return append(lines, line)
}
