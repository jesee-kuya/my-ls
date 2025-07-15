package util

import (
	"fmt"
	"strings"
)

func InsertSorted(name, colour, reset string, names []string) []string {
	if name == "." || name == ".." {
		return append([]string{fmt.Sprintf("%s%s%s", colour, name, reset)}, names...)
	}
	colored := fmt.Sprintf("%s%s%s", colour, name, reset)

	for i, val := range names {
		if ((strings.ToLower(TrimStart(name)) < strings.ToLower(TrimStart(StripANSI(val)))) || (strings.ToLower(TrimStart(name)) == strings.ToLower(TrimStart(StripANSI(val))))) && (val != "." && val != "..") {
			return append(names[:i], append([]string{colored}, names[i:]...)...)
		}
	}

	return append(names, colored)
}

func InsertSortedLong(line string, lines []string) []string {
	for i, val := range lines {
		if (strings.ToLower(TrimStart(StripANSI(StripLong(line)))) < strings.ToLower(TrimStart(StripANSI(StripLong(val))))) || (strings.ToLower(TrimStart(StripANSI(StripLong(line)))) == strings.ToLower(TrimStart(StripANSI(StripLong(val))))) {
			return append(lines[:i], append([]string{line}, lines[i:]...)...)
		}
	}

	return append(lines, line)
}

func TrimStart(name string) string {
	return strings.TrimLeft(name, ".")
}
