package util

import (
	"errors"
	"os"
)

// IsValidDir ensures dirPath is a valid directory
func IsValidDir(dirPath string) (bool, error) {
	info, err := os.Stat(dirPath)
	if err != nil {
		return false, errors.New("directory does not exist")
	}
	if !info.IsDir() {
		return false, errors.New("not a directory")
	}
	return true, nil
}
