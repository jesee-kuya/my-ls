package util

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"syscall"
)

type Flags struct {
	ShowAll    bool
	Longformat bool
	Reverse    bool
	Recursive  bool
	TimeSort   bool
}

type fileDisplayInfo struct {
	os.FileInfo

	mode    string
	links   string
	user    string
	group   string
	size    string
	modTime string
}

type maxWidths struct {
	links int
	user  int
	group int
	size  int
}

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
func ReadDirNames(dirPath string, flag Flags) ([]string, error) {
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
	if flag.ShowAll {
		names = append(names, fmt.Sprintf("%s.%s", dirColour, reset))
		names = append(names, fmt.Sprintf("%s..%s", dirColour, reset))
	}

	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden files unless showAll is true
		if !flag.ShowAll && strings.HasPrefix(name, ".") {
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

		if flag.TimeSort {
			names = InsertSortedByTime(name, colour, reset, dirPath, names)
		} else {
			names = InsertSorted(name, colour, reset, names)
		}
	}

	if flag.Reverse {
		Reverse(names)
	}

	return names, nil
}

func ReadDirNamesLong(dirPath string, flag Flags) ([]string, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	entries, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	var displayInfos []fileDisplayInfo
	var widths maxWidths
	var totalBlocks int64

	// Combine initial filtering and data gathering into a single list
	filesToProcess := []os.FileInfo{}
	if flag.ShowAll {
		for _, special := range []string{".", ".."} {
			info, err := os.Lstat(joinPath(dirPath, special))
			if err == nil {
				filesToProcess = append(filesToProcess, info)
			}
		}
	}
	for _, entry := range entries {
		if !flag.ShowAll && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		filesToProcess = append(filesToProcess, entry)
	}

	for _, info := range filesToProcess {
		fullPath := joinPath(dirPath, info.Name())
		stat := getStat(fullPath)
		totalBlocks += int64(stat.Blocks)

		// Owner and group
		uid := fmt.Sprint(stat.Uid)
		owner := uid
		if u, err := user.LookupId(uid); err == nil {
			owner = u.Username
		}
		gid := fmt.Sprint(stat.Gid)
		group := gid
		if g, err := user.LookupGroupId(gid); err == nil {
			group = g.Name
		}

		// Convert numeric fields to strings for length calculation
		linksStr := fmt.Sprint(stat.Nlink)
		sizeStr := fmt.Sprint(info.Size())

		// Update max widths
		if len(linksStr) > widths.links {
			widths.links = len(linksStr)
		}
		if len(owner) > widths.user {
			widths.user = len(owner)
		}
		if len(group) > widths.group {
			widths.group = len(group)
		}
		if len(sizeStr) > widths.size {
			widths.size = len(sizeStr)
		}

		// Store the processed info
		displayInfos = append(displayInfos, fileDisplayInfo{
			FileInfo: info,
			mode:     info.Mode().String(),
			links:    linksStr,
			user:     owner,
			group:    group,
			size:     sizeStr,
			modTime:  info.ModTime().Format("Jan _2 15:04"),
		})
	}

	if flag.Reverse {
		for i, j := 0, len(displayInfos)-1; i < j; i, j = i+1, j-1 {
			displayInfos[i], displayInfos[j] = displayInfos[j], displayInfos[i]
		}
	}

	var lines []string

	for _, di := range displayInfos {
		color := getFileColor(di.Mode(), di.Name())
		fileName := fmt.Sprintf("%s%s%s", color, di.Name(), reset)

		// Use the calculated max widths to format the line perfectly
		line := fmt.Sprintf("%-10s %*s %-*s %-*s %*s %s %s",
			di.mode,
			widths.links, di.links,
			widths.user, di.user,
			widths.group, di.group,
			widths.size, di.size,
			di.modTime,
			fileName,
		)
		if flag.TimeSort {
			lines = InsertSortedLongByTime(line, dirPath, lines)
		} else {
			lines = InsertSortedLong(line, lines)
		}
	}
	lines = append([]string{fmt.Sprintf("total %d", totalBlocks/2)}, lines...)

	return lines, nil
}

// getFileColor remains the same as it's a perfect helper function
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

// getStat also remains the same
func getStat(path string) syscall.Stat_t {
	var stat syscall.Stat_t
	if err := syscall.Lstat(path, &stat); err != nil {
		_ = syscall.Stat(path, &stat)
	}
	return stat
}

// joinPath joins directory and file name with proper separator
func joinPath(dir, file string) string {
	if dir == "" {
		return file
	}
	if strings.HasSuffix(dir, "/") {
		return dir + file
	}
	return dir + "/" + file
}

// getAbsPath returns absolute path (simplified version)
func getAbsPath(path string) (string, error) {
	if strings.HasPrefix(path, "/") {
		return path, nil
	}
	// For relative paths, we'll use a simple approach
	// In a real implementation, you'd need to get the current working directory
	// but since we can't use filepath, we'll return the path as-is for relative paths
	return path, nil
}

// CollectDirectoriesRecursively traverses directories recursively and returns all directory paths
func CollectDirectoriesRecursively(rootPaths []string, flags Flags) ([]string, error) {
	var allDirs []string
	visited := make(map[string]bool)

	for _, rootPath := range rootPaths {
		info, err := IsValidDir(rootPath)
		if err != nil {
			return nil, err
		}

		if !info.IsDir() {
			// If it's a file, just add it to the list
			allDirs = append(allDirs, rootPath)
			continue
		}

		// Add the root directory first
		allDirs = append(allDirs, rootPath)

		// Recursively collect subdirectories
		err = collectSubdirectories(rootPath, flags, &allDirs, visited)
		if err != nil {
			return nil, err
		}
	}

	return allDirs, nil
}

// collectSubdirectories is a helper function that recursively collects subdirectories
func collectSubdirectories(dirPath string, flags Flags, allDirs *[]string, visited map[string]bool) error {
	// Prevent infinite loops with symlinks
	absPath, err := getAbsPath(dirPath)
	if err != nil {
		return err
	}

	if visited[absPath] {
		return nil
	}
	visited[absPath] = true

	dir, err := os.Open(dirPath)
	if err != nil {
		return err
	}
	defer dir.Close()

	entries, err := dir.Readdir(-1)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden files unless showAll is true
		if !flags.ShowAll && strings.HasPrefix(name, ".") {
			continue
		}

		// Skip . and .. entries
		if name == "." || name == ".." {
			continue
		}

		if entry.IsDir() {
			subDirPath := joinPath(dirPath, name)
			*allDirs = append(*allDirs, subDirPath)

			// Recursively process subdirectory
			err = collectSubdirectories(subDirPath, flags, allDirs, visited)
			if err != nil {
				// Continue processing other directories even if one fails
				continue
			}
		}
	}

	return nil
}
