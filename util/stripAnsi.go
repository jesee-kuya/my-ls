package util

import "regexp"

var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*m`)
var longFormat = regexp.MustCompile(`^([\-d][rwx\-]{9})\s+\d+\s+\S+\s+\S+\s+\d+\s+[A-Z][a-z]{2}\s+\d{1,2}\s+\d{2}:\d{2}\s+`)

// StripAnsi removes ANSI escape codes from a string
func StripANSI(s string) string {
	return ansiEscape.ReplaceAllString(s, "")
}

func StripLong(s string) string {
	return longFormat.ReplaceAllString(s, "")
}
