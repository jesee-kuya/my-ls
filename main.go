package main

import (
	"os"

	"github.com/jesee-kuya/my-ls/print"
)

// Flags represents command-line flags for my-ls
type Flags struct {
	ShowAll bool // -a flag: show all files including hidden ones
}

// parseArgs parses command-line arguments and returns flags and paths
func parseArgs(args []string) (Flags, []string) {
	flags := Flags{}
	var paths []string

	// TODO: implement flag parsing logic

	return flags, paths
}

func main() {
	paths := []string{"."}

	if len(os.Args) > 1 {
		paths = os.Args[1:]
	}

	print.Print(paths)
}
