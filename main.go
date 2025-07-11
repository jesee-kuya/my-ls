package main

import (
	"os"

	"github.com/jesee-kuya/my-ls/print"
)

func main() {
	// Parse command-line arguments using parseArgs
	flags, paths := parseArgs(os.Args[1:])

	// Use the print package with flags
	print.Print(paths, flags)
}

func ParseArgs(args []string) (print.Flags, []string) {
	var flags print.Flags
	var paths []string

	for _, arg := range args {
		if strings.HasPrefix(arg, "-") {
			// Handle flags - loop through each character after the dash
			for _, char := range arg[1:] {
				switch char {
				case 'a':
					flags.ShowHidden = true
				}
			}
		} else {
			// It's a path
			paths = append(paths, arg)
		}
	}

	// Default to current directory if no paths specified
	if len(paths) == 0 {
		paths = []string{"."}
	}

	return flags, paths
}
