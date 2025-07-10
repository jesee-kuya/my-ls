package print

import (
	"fmt"

	"github.com/jesee-kuya/my-ls/util"
)

// Flags represents the command line flags
type Flags struct {
	ShowHidden bool // -a flag
}

func Print(paths []string, flags any) {
	outErrors := []string{}
	singleFiles := []string{}
	dirContents := []string{}
	content := []any{}
	multipleDirs := false
	if len(paths) > 1 {
		multipleDirs = true
	}

	for _, dirPath := range paths {
		info, err := util.IsValidDir(dirPath)
		if err != nil {
			outErrors = append(outErrors, fmt.Sprintf("Error: %v\n", err.Error()))
			continue
		}

		if !info.IsDir() {
			singleFiles = append(singleFiles, dirPath)
			continue
		}

		files, err := util.ReadDirNames(dirPath, true)
		if err != nil {
			outErrors = append(outErrors, fmt.Sprintf("Error reading directory: %v\n", err.Error()))
			continue
		}

		if multipleDirs {
			dirContents = append(dirContents, fmt.Sprintf("%v:\n", dirPath))
		}

		dirContents = append(dirContents, files...)
		content = append(content, dirContents)
		dirContents = []string{}
	}
	for _, err := range outErrors {
		fmt.Println(err)
	}

	for i, file := range singleFiles {
		if i == len(singleFiles)-1 {
			fmt.Print(file + "\n\n")
			continue
		}
		fmt.Print(file + "  ")
	}

	for i, c := range content {
		if i != 0 {
			fmt.Println()
		}
		for i, line := range c.([]string) {
			if i == len(c.([]string))-1 {
				fmt.Println(line)
				continue
			}
			if util.HasANSIPrefix(line) {
				fmt.Print(line + "  ")
				continue
			} else {
				fmt.Print(line)
			}
		}
	}
}
