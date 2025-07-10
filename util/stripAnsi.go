package util

import "regexp"

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// StripAnsi removes ANSI escape codes from a string
func StripANSI(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}
