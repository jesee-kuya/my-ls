package util

import (
	"fmt"
	"os"
	"strings"
)

const (
	reset = "\033[0m"

	dirColour     = "\033[01;34m"    // bold blue
	exeColour     = "\033[01;32m"    // bold green
	symlinkColour = "\033[01;36m"    // bold cyan
	socketColour  = "\033[01;35m"    // bold magenta
	pipeColour    = "\033[40;33m"    // yellow on black background
	deviceColour  = "\033[40;33;01m" // bold yellow on black (block/char dev)
	archiveColour = "\033[01;31m"    // bold red
)

// ReadDirNames returns a list of file and directory names in dirPath
func ReadDirNames(dirPath string, showAll bool) ([]string, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	entries, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var names []string

	// Add . and .. entries when showAll is true
	if showAll {
		names = append(names, fmt.Sprintf("%s.%s", dirColour, reset))
		names = append(names, fmt.Sprintf("%s..%s", dirColour, reset))
	}

	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden files unless showAll is true
		if !showAll && strings.HasPrefix(name, ".") {
			continue
		}

		mode := entry.Mode()
		colour := reset

		switch {
		case mode.IsDir():
			colour = dirColour

		case mode&os.ModeSymlink != 0:
			colour = symlinkColour

		case mode&os.ModeSocket != 0:
			colour = socketColour

		case mode&os.ModeNamedPipe != 0:
			colour = pipeColour

		case mode&os.ModeDevice != 0:
			colour = deviceColour

		case mode&0o111 != 0:
			colour = exeColour

		case strings.HasSuffix(name, ".tar") ||
			strings.HasSuffix(name, ".gz") ||
			strings.HasSuffix(name, ".tgz") ||
			strings.HasSuffix(name, ".zip") ||
			strings.HasSuffix(name, ".bz2") ||
			strings.HasSuffix(name, ".xz"):
			colour = archiveColour
		}

		names = append(names, fmt.Sprintf("%s%s%s", colour, name, reset))
	}

	return names, nil
}
