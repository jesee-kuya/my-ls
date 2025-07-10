package main

import (
	"fmt"
	"os"

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
