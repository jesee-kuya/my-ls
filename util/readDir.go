package util

import (
	"fmt"
	"os"
)

// ReadDirNames returns a list of file and directory names in dirPath
func ReadDirNames(dirPath string) ([]string, error) {
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
	for _, entry := range entries {
		if entry.IsDir() {
			names = append(names, fmt.Sprintf("%v%v", "\x1b[34m", entry.Name()))
		} else {
			names = append(names, fmt.Sprintf("%v%v", "\x1b[0m", entry.Name()))
		}
	}

	return names, nil
}
