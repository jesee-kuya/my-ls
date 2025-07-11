package print

import (
	"fmt"

	"github.com/jesee-kuya/my-ls/util"
)

// Flags represents command-line flags for my-ls
type Flags struct {
	ShowAll    bool
	Longformat bool
}

func Print(paths []string, flags Flags) {
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
		var files []string

		if flags.Longformat {
			files, err = util.ReadDirNamesLong(dirPath, flags.ShowAll)
			if err != nil {
				outErrors = append(outErrors, fmt.Sprintf("Error reading directory: %v\n", err.Error()))
				continue
			}
		} else {
			files, err = util.ReadDirNames(dirPath, flags.ShowAll)
			if err != nil {
				outErrors = append(outErrors, fmt.Sprintf("Error reading directory: %v\n", err.Error()))
				continue
			}
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
		if flags.Longformat {
			fmt.Println(file)
			continue
		}

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
			if flags.Longformat {
				fmt.Println(line)
				continue
			}

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
