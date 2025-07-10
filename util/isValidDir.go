package util

import (
	"fmt"
	"os"
)

// IsValidDir ensures dirPath is a valid directory
func IsValidDir(dirPath string) (bool, os.FileInfo, error) {
	info, err := os.Stat(dirPath)
	if err != nil {
		return false, nil, fmt.Errorf("cannot access '%v': No such file or directory", dirPath)
	}

	return true, info, nil
}
