package util

import (
	"fmt"
	"os"
)

// IsValidDir ensures dirPath is a valid directory
func IsValidDir(dirPath string) (os.FileInfo, error) {
	info, err := os.Stat(dirPath)
	if err != nil {
		return nil, fmt.Errorf("cannot access '%v': No such file or directory", dirPath)
	}

	return info, nil
}
