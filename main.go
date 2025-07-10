package main

import (
	"fmt"
	"os"

	"github.com/jesee-kuya/my-ls/util"
)

func main() {
	paths := []string{"."}
	outErrors := []string{}
	singleFiles := []string{}
	dirContents := []string{}
	multipleDirs := false

	if len(os.Args) > 1 {
		multipleDirs = true
		paths = os.Args[1:]
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

		files, err := util.ReadDirNames(dirPath)
		if err != nil {
			outErrors = append(outErrors, fmt.Sprintf("Error reading directory: %v\n", err.Error()))
			continue
		}

		if multipleDirs {
			dirContents = append(dirContents, fmt.Sprintf("%v:\n", dirPath))
		}

		dirContents = append(dirContents, files...)
	}
	for _, err := range outErrors {
		fmt.Println(err)
	}

	for _, file := range singleFiles {
		fmt.Println(file)
	}

	for _, content := range dirContents {
		fmt.Print(content + "  ")
	}
	fmt.Println()
}
