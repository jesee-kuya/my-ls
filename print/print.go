package print

import (
	"fmt"

	"github.com/jesee-kuya/my-ls/util"
)

func Print(paths []string, flags util.Flags) {
	outErrors := []string{}
	singleFiles := []string{}
	dirContents := []string{}
	content := []any{}

	// Handle recursive listing
	if flags.Recursive {
		allPaths, err := util.CollectDirectoriesRecursively(paths, flags)
		if err != nil {
			outErrors = append(outErrors, fmt.Sprintf("Error during recursive traversal: %v\n", err.Error()))
		} else {
			paths = allPaths
		}
	}

	multipleDirs := false
	if len(paths) > 1 || flags.Recursive {
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
			files, err = util.ReadDirNamesLong(dirPath, flags)
			if err != nil {
				outErrors = append(outErrors, fmt.Sprintf("Error reading directory: %v\n", err.Error()))
				continue
			}
		} else {
			files, err = util.ReadDirNames(dirPath, flags)
			if err != nil {
				outErrors = append(outErrors, fmt.Sprintf("Error reading directory: %v\n", err.Error()))
				continue
			}
		}

		if multipleDirs {
			dirContents = append(dirContents, fmt.Sprintf("%v:", dirPath))
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

		lines := c.([]string)
		if len(lines) == 0 {
			continue
		}

		// Print directory header (if present)
		if len(lines) > 0 && len(lines[0]) > 0 && lines[0][len(lines[0])-1] == ':' {
			fmt.Println(lines[0])
			lines = lines[1:] // Skip the header for content printing
		}

		// Print the directory contents
		for j, line := range lines {
			if flags.Longformat {
				fmt.Println(line)
				continue
			}

			if j == len(lines)-1 {
				fmt.Println(line)
				continue
			}
			if util.HasANSIPrefix(line) {
				fmt.Print(line + "  ")
				continue
			} else {
				fmt.Print(line + "  ")
			}
		}
	}
}
