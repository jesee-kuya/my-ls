package util

import (
	"errors"
	"os"
	"strings"
)

// ReadDirNames returns a list of file and directory names in dirPath
func ReadDirNames(dirPath string, showHidden bool) ([]string, error) {
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
		name := entry.Name()
		if !showHidden && strings.HasPrefix(name, ".") {
			continue
		}
		names = append(names, name)
	}
	if names == nil {
		return nil, errors.New("no entries found")
	}
	return names, nil
}
