package main

import (
	"os"

	"github.com/jesee-kuya/my-ls/print"
)

// Flags represents command-line flags for my-ls
type Flags struct {
	ShowAll bool // -a flag: show all files including hidden ones
}

func main() {
	paths := []string{"."}

	if len(os.Args) > 1 {
		paths = os.Args[1:]
	}

	print.Print(paths)
}
