package util

import "regexp"

var ansiPrefix = regexp.MustCompile(`^\x1b\[[0-9;]*m`)

func HasANSIPrefix(s string) bool {
	return ansiPrefix.MatchString(s)
}
