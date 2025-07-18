package util

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

// CompareStrings compares two strings based on custom sorting rules.
// Significant parts (numeric and alphabetic characters) are compared first, with numeric characters before alphabetic.
// Numeric characters are compared by ASCII values, alphabetic characters case-insensitively with lowercase prioritized.
// If significant parts are equal, compare original strings character by character:
// - Numeric vs. anything: numeric comes first.
// - Alphabetic vs. alphabetic: lowercase before uppercase, case-insensitive otherwise.
// - Other pairs (special vs. special, alphabetic vs. special): ASCII order.
func CompareStrings(a, b string) bool {
	// Convert strings to runes for proper Unicode handling
	ra, rb := []rune(a), []rune(b)

	// Check if both strings contain only special (non-alphabetic, non-numeric) characters
	hasSignificantA, hasSignificantB := false, false
	for _, r := range ra {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			hasSignificantA = true
			break
		}
	}
	for _, r := range rb {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
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
		if unicode.IsLetter(r) {
			sigA = append(sigA, unicode.ToLower(r))
		} else if unicode.IsDigit(r) {
			sigA = append(sigA, r)
		}
	}
	for _, r := range rb {
		if unicode.IsLetter(r) {
			sigB = append(sigB, unicode.ToLower(r))
		} else if unicode.IsDigit(r) {
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
		if unicode.IsDigit(ca) && unicode.IsLetter(cb) {
			return true
		}
		if unicode.IsLetter(ca) && unicode.IsDigit(cb) {
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
		if unicode.IsDigit(ca) && !unicode.IsDigit(cb) {
			return true
		}
		if !unicode.IsDigit(ca) && unicode.IsDigit(cb) {
			return false
		}
		if unicode.IsDigit(ca) && unicode.IsDigit(cb) {
			return ca < cb
		}
		// If both are alphabetic, prioritize lowercase
		if unicode.IsLetter(ca) && unicode.IsLetter(cb) {
			caLower, cbLower := unicode.ToLower(ca), unicode.ToLower(cb)
			if caLower != cbLower {
				return caLower < cbLower
			}
			// If lowercase versions are equal, lowercase comes before uppercase
			if unicode.IsLower(ca) && unicode.IsUpper(cb) {
				return true
			}
			if unicode.IsUpper(ca) && unicode.IsLower(cb) {
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
