package main

import (
	"os"

	"github.com/jesee-kuya/my-ls/print"
)

// parseArgs parses command-line arguments and returns flags and paths
func parseArgs(args []string) (print.Flags, []string) {
	flags := print.Flags{}
	var paths []string

	for _, arg := range args {
		if len(arg) > 0 && arg[0] == '-' {
			// Parse flags
			for _, char := range arg[1:] {
				switch char {
				case 'a':
					flags.ShowAll = true
				case 'l':
					flags.Longformat = true
				}
			}
		} else {
			// It's a path
			paths = append(paths, arg)
		}
	}

	// If no paths specified, use current directory
	if len(paths) == 0 {
		paths = []string{"."}
	}

	return flags, paths
}

func main() {
	var flags print.Flags
	var paths []string

	if len(os.Args) > 1 {
		flags, paths = parseArgs(os.Args[1:])
	} else {
		paths = []string{"."}
	}

	print.Print(paths, flags)
}
