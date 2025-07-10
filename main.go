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
		paths = os.Args[1:]
	}

	if len(paths) > 1 {
		multipleDirs = true
	}

	for i, dirPath := range paths {
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
		if i != len(paths)-1 {
			dirContents = append(dirContents, "\n\n")
		}
	}
	for _, err := range outErrors {
		fmt.Println(err)
	}

	for _, file := range singleFiles {
		fmt.Println(file)
	}

	for i, content := range dirContents {
		if i != 1 {
			fmt.Print(content + "  ")
		} else {
			fmt.Print(content)
		}
		if i == len(dirContents)-1 {
			fmt.Println()
		}
	}
}
