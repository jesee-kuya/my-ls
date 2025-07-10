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
	paths := []string{"."}
	if len(os.Args) > 1 {
		paths = os.Args[1:]
	}

	for _, dirPath := range paths {
		info, err := util.IsValidDir(dirPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			return
		}

		if !info.IsDir() {
			fmt.Println(info.Name())
			return
		}

		files, err := util.ReadDirNames(dirPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading directory: %v\n", err)
			return
		}

		for _, name := range files {
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
