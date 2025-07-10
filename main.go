package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/jesee-kuya/my-ls/util"
)

type Flags struct {
	ShowHidden bool // -a flag
}

func main() {
	// Parse command-line arguments using parseArgs
	flags, paths := parseArgs(os.Args[1:]) // Pass os.Args[1:] to skip the program name

	for _, dirPath := range paths {
		// Validate if the path is a valid directory
		info, err := util.IsValidDir(dirPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			continue // Continue with the next path instead of exiting
		}

		// If it's not a directory, print the file name and continue
		if !info.IsDir() {
			fmt.Println(info.Name())
			continue
		}

		// Read directory contents
		files, err := util.ReadDirNames(dirPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading directory %s: %v\n", dirPath, err)
			continue
		}

		// Filter files based on ShowHidden flag
		for _, name := range files {
			// Skip hidden files (starting with '.') unless ShowHidden is true
			if !flags.ShowHidden && strings.HasPrefix(name, ".") {
				continue
			}
			fmt.Println(name)
		}
	}
}

func parseArgs(args []string) (Flags, []string) {
	var flags Flags
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
