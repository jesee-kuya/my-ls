package main

import (
	"os"

	"github.com/jesee-kuya/my-ls/print"
)

func main() {
	paths := []string{"."}

	if len(os.Args) > 1 {
		paths = os.Args[1:]
	}

	print.Print(paths)
}
