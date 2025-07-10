package main

import (
	"fmt"
	"os"

	"github.com/jesee-kuya/my-ls/util"
)

func main() {
	dirPath := "."
	if len(os.Args) > 1 {
		dirPath = os.Args[1]
	}

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
