package util

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"syscall"
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

		names = InsertSorted(name, colour, reset, names)
	}

	return names, nil
}

func ReadDirNamesLong(dirPath string, showAll bool) ([]string, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	entries, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var lines []string
	var totalBlocks int64
	filesToShow := []os.FileInfo{}

	// Optionally include . and ..
	if showAll {
		for _, special := range []string{".", ".."} {
			fullPath := filepath.Join(dirPath, special)
			info, err := os.Lstat(fullPath)
			if err == nil {
				filesToShow = append(filesToShow, info)
				stat := getStat(fullPath)
				totalBlocks += int64(stat.Blocks)
			}
		}
	}

	for _, entry := range entries {
		name := entry.Name()
		if !showAll && strings.HasPrefix(name, ".") {
			continue
		}
		filesToShow = append(filesToShow, entry)
		stat := getStat(filepath.Join(dirPath, name))
		totalBlocks += int64(stat.Blocks)
	}

	for _, entry := range filesToShow {
		line := formatLongEntry(entry.Name(), dirPath)
		lines = InsertSortedLong(line, lines)
	}

	lines = append([]string{fmt.Sprintf("total %d", totalBlocks/2)}, lines...)

	return lines, nil
}

// formatLongEntry builds an ls -l style line for a file
func formatLongEntry(name, dirPath string) string {
	fullPath := filepath.Join(dirPath, name)

	info, err := os.Lstat(fullPath)
	if err != nil {
		return fmt.Sprintf("Error reading %s: %v", name, err)
	}

	var stat syscall.Stat_t
	if err := syscall.Stat(fullPath, &stat); err != nil {
		return fmt.Sprintf("Error stat %s: %v", name, err)
	}

	// Permissions
	mode := info.Mode().String()

	// Link count
	links := stat.Nlink

	// Owner and group
	uid := fmt.Sprint(stat.Uid)
	gid := fmt.Sprint(stat.Gid)

	owner := uid
	if u, err := user.LookupId(uid); err == nil {
		owner = u.Username
	}
	group := gid
	if g, err := user.LookupGroupId(gid); err == nil {
		group = g.Name
	}

	// Size
	size := info.Size()

	// Time
	modTime := info.ModTime().Format("Jan _2 15:04")

	// Color
	color := getFileColor(info.Mode(), name)

	// Format full line
	return fmt.Sprintf("%s %2d %-8s %-8s %6d %s %s%s%s",
		mode, links, owner, group, size, modTime, color, name, reset)
}

// getFileColor returns the ANSI color based on file mode or suffix
func getFileColor(mode os.FileMode, name string) string {
	switch {
	case mode.IsDir():
		return dirColour
	case mode&os.ModeSymlink != 0:
		return symlinkColour
	case mode&os.ModeSocket != 0:
		return socketColour
	case mode&os.ModeNamedPipe != 0:
		return pipeColour
	case mode&os.ModeDevice != 0:
		return deviceColour
	case mode&0o111 != 0:
		return exeColour
	case strings.HasSuffix(name, ".tar"),
		strings.HasSuffix(name, ".gz"),
		strings.HasSuffix(name, ".tgz"),
		strings.HasSuffix(name, ".zip"),
		strings.HasSuffix(name, ".bz2"),
		strings.HasSuffix(name, ".xz"):
		return archiveColour
	default:
		return reset
	}
}

func getStat(path string) syscall.Stat_t {
	var stat syscall.Stat_t
	_ = syscall.Stat(path, &stat)
	return stat
}
