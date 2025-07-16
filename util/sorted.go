package util

import (
	"fmt"
	"os"
	"strings"
	"unicode"
)

// compareFilenames implements locale-aware filename comparison similar to ls in en_US.UTF-8
// It prioritizes lowercase over uppercase for all letters, deprioritizes non-alphabetic characters,
// and is case-sensitive. After a common prefix, if both characters at the same position are
// non-alphabetic, they are ignored until an alphabetic character is found in at least one string.
// Alphabetic characters are prioritized over non-alphabetic ones, with lowercase before uppercase.
// If both strings have the same alphabetic character (case-insensitive) at the same position,
// the non-alphabetic characters immediately following the prefix are compared using ASCII values.
func compareFilenames(a, b string) bool {
	// Convert strings to runes for proper Unicode handling
	ra, rb := []rune(a), []rune(b)

	// Track the first position after the common prefix where both are non-alphabetic
	var firstNonAlphaPos int = -1
	prefixLen := 0
	for i := 0; i < len(ra) && i < len(rb); i++ {
		if ra[i] != rb[i] {
			prefixLen = i
			break
		}
	}
	if prefixLen == len(ra) && prefixLen == len(rb) {
		return false // Equal strings
	}

	// Compare strings character by character after the prefix
	for i := prefixLen; i < len(ra) || i < len(rb); i++ {
		// Handle cases where one string is longer than the other
		if i >= len(ra) {
			return !unicode.IsLetter(rb[i]) // Non-letter in b means a comes first, else b
		}
		if i >= len(rb) {
			return unicode.IsLetter(ra[i]) // Letter in a means a comes first, else b
		}

		ca, cb := ra[i], rb[i]
		if ca == cb {
			continue
		}
		// If both characters are non-alphabetic, store the first such position and skip
		if !unicode.IsLetter(ca) && !unicode.IsLetter(cb) {
			if firstNonAlphaPos == -1 {
				firstNonAlphaPos = i
			}
			continue
		}
		// If one is alphabetic and the other is not, prioritize the alphabetic character
		if unicode.IsLetter(ca) && !unicode.IsLetter(cb) {
			return true
		}
		if !unicode.IsLetter(ca) && unicode.IsLetter(cb) {
			return false
		}
		// Both are alphabetic, use case-insensitive comparison
		caLower, cbLower := unicode.ToLower(ca), unicode.ToLower(cb)
		if caLower != cbLower {
			// If lowercase versions differ, use them for ordering
			return caLower < cbLower
		}
		// If lowercase versions are equal, prioritize lowercase over uppercase
		if unicode.IsLower(ca) && unicode.IsUpper(cb) {
			return true
		}
		if unicode.IsUpper(ca) && unicode.IsLower(cb) {
			return false
		}
		// If alphabetic characters are the same (case-insensitive), backtrack to compare
		// the first non-alphabetic characters after the prefix using ASCII
		if firstNonAlphaPos != -1 {
			return ra[firstNonAlphaPos] < rb[firstNonAlphaPos]
		}
	}
	return false // Equal up to the compared length
}

func InsertSorted(name, colour, reset string, names []string) []string {
	if name == "." || name == ".." {
		return append([]string{fmt.Sprintf("%s%s%s", colour, name, reset)}, names...)
	}
	colored := fmt.Sprintf("%s%s%s", colour, name, reset)

	for i, val := range names {
		cleanVal := TrimStart(StripANSI(val))
		if (compareFilenames(TrimStart(name), cleanVal) || (TrimStart(name) == cleanVal)) && (val != "." && val != "..") {
			return append(names[:i], append([]string{colored}, names[i:]...)...)
		}
	}

	return append(names, colored)
}

func InsertSortedLong(line string, lines []string) []string {
	for i, val := range lines {
		lineName := TrimStart(StripANSI(StripLong(line)))
		valName := TrimStart(StripANSI(StripLong(val)))
		if compareFilenames(lineName, valName) || (lineName == valName) {
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
			// If we can't get the time, insert burada
			return append(names[:i], append([]string{colored}, names[i:]...)...)
		}
		existingTime := existingInfo.ModTime()

		// Insert if new file is newer (or same time, then alphabetical)
		if newTime.After(existingTime) || (newTime.Equal(existingTime) && compareFilenames(name, cleanVal)) {
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
		if newTime.After(existingTime) || (newTime.Equal(existingTime) && compareFilenames(fileName, existingFileName)) {
			return append(lines[:i], append([]string{line}, lines[i:]...)...)
		}
	}

	return append(lines, line)
}
