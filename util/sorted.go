package util

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

// CompareStrings compares two strings based on custom sorting rules.
// If both strings contain only special (non-alphabetic) characters, sort by ASCII values (lowest first).
// If strings contain alphabetic characters, compare case-insensitively, prioritizing lowercase over uppercase.
// For mixed strings, compare alphabetic parts first (ignoring special characters, shorter comes first),
// then compare original strings character by character: lowercase before uppercase for alphabetic pairs,
// ASCII order for all other pairs.
func CompareStrings(a, b string) bool {
	// Convert strings to runes for proper Unicode handling
	ra, rb := []rune(a), []rune(b)

	// Check if both strings contain only special (non-alphabetic) characters
	hasLetterA, hasLetterB := false, false
	for _, r := range ra {
		if unicode.IsLetter(r) {
			hasLetterA = true
			break
		}
	}
	for _, r := range rb {
		if unicode.IsLetter(r) {
			hasLetterB = true
			break
		}
	}

	// If both strings have only special characters, compare by ASCII values
	if !hasLetterA && !hasLetterB {
		for i := 0; i < len(ra) && i < len(rb); i++ {
			if ra[i] != rb[i] {
				return ra[i] < rb[i]
			}
		}
		// If equal up to shorter length, shorter string comes first
		return len(ra) < len(rb)
	}

	// If either string has alphabetic characters, compare alphabetic parts first
	var alphaA, alphaB []rune
	for _, r := range ra {
		if unicode.IsLetter(r) {
			alphaA = append(alphaA, unicode.ToLower(r))
		}
	}
	for _, r := range rb {
		if unicode.IsLetter(r) {
			alphaB = append(alphaB, unicode.ToLower(r))
		}
	}

	// Compare alphabetic parts (case-insensitive)
	for i := 0; i < len(alphaA) && i < len(alphaB); i++ {
		if alphaA[i] != alphaB[i] {
			return alphaA[i] < alphaB[i]
		}
	}
	// If alphabetic parts are equal, shorter alphabetic part comes first
	if len(alphaA) != len(alphaB) {
		return len(alphaA) < len(alphaB)
	}

	// If alphabetic parts are equal, compare original strings character by character
	for i := 0; i < len(ra) && i < len(rb); i++ {
		ca, cb := ra[i], rb[i]
		if ca == cb {
			continue
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
