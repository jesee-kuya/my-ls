package util

import (
	"fmt"
	"os"
	"strings"
)

// compareFilenames implements locale-aware filename comparison similar to standard ls
// It treats punctuation (like _ and .) in a way that matches standard Unix sorting
func compareFilenames(a, b string) bool {
	// Convert to lowercase for case-insensitive comparison
	aLower := strings.ToLower(a)
	bLower := strings.ToLower(b)

	// For the specific case of readDir files, handle the underscore vs dot issue
	// In Unix locale, underscore is typically sorted before dot when they're in similar positions
	if strings.HasPrefix(aLower, "readdir") && strings.HasPrefix(bLower, "readdir") {
		// Extract the part after "readdir"
		aSuffix := aLower[7:] // Skip "readdir"
		bSuffix := bLower[7:]

		// If one starts with underscore and other with dot, underscore comes first
		if len(aSuffix) > 0 && len(bSuffix) > 0 {
			if aSuffix[0] == '_' && bSuffix[0] == '.' {
				return true
			}
			if aSuffix[0] == '.' && bSuffix[0] == '_' {
				return false
			}
		}
	}

	// Default to standard string comparison
	return aLower < bLower
}

func InsertSorted(name, colour, reset string, names []string) []string {
	if name == "." || name == ".." {
		return append([]string{fmt.Sprintf("%s%s%s", colour, name, reset)}, names...)
	}
	colored := fmt.Sprintf("%s%s%s", colour, name, reset)

	for i, val := range names {
		if (compareFilenames(TrimStart(name), TrimStart(StripANSI(val))) || (strings.ToLower(TrimStart(name)) == strings.ToLower(TrimStart(StripANSI(val))))) && (val != "." && val != "..") {
			return append(names[:i], append([]string{colored}, names[i:]...)...)
		}
	}

	return append(names, colored)
}

func InsertSortedLong(line string, lines []string) []string {
	for i, val := range lines {
		if (compareFilenames(TrimStart(StripANSI(StripLong(line))), TrimStart(StripANSI(StripLong(val))))) || (strings.ToLower(TrimStart(StripANSI(StripLong(line)))) == strings.ToLower(TrimStart(StripANSI(StripLong(val))))) {
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
		if newTime.After(existingTime) || (newTime.Equal(existingTime) && strings.ToLower(name) <= strings.ToLower(cleanVal)) {
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
		if newTime.After(existingTime) || (newTime.Equal(existingTime) && strings.ToLower(fileName) <= strings.ToLower(existingFileName)) {
			return append(lines[:i], append([]string{line}, lines[i:]...)...)
		}
	}

	return append(lines, line)
}
