package print

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/jesee-kuya/my-ls/util"
)

// getTerminalWidth returns the terminal width, defaulting to 80 if unable to determine
func getTerminalWidth() int {
	type winsize struct {
		Row    uint16
		Col    uint16
		Xpixel uint16
		Ypixel uint16
	}

	ws := &winsize{}
	retCode, _, _ := syscall.Syscall(syscall.SYS_IOCTL,
		uintptr(syscall.Stdin),
		uintptr(syscall.TIOCGWINSZ),
		uintptr(unsafe.Pointer(ws)))

	if int(retCode) == -1 {
		return 80 // Default width
	}
	return int(ws.Col)
}

// formatInColumns formats a list of files in columns like standard ls
func formatInColumns(files []string) string {
	if len(files) == 0 {
		return ""
	}

	// Strip ANSI codes to calculate actual display width
	displayFiles := make([]string, len(files))
	fileLengths := make([]int, len(files))
	maxLen := 0
	for i, file := range files {
		displayFiles[i] = util.StripANSI(file)
		fileLengths[i] = len(displayFiles[i])
		if fileLengths[i] > maxLen {
			maxLen = fileLengths[i]
		}
	}

	termWidth := getTerminalWidth()

	// Try different numbers of columns to find the optimal layout
	bestCols := 1
	bestRows := len(files)

	for numCols := 1; numCols <= len(files); numCols++ {
		numRows := (len(files) + numCols - 1) / numCols

		// Calculate column widths for this layout (column-major ordering)
		colWidths := make([]int, numCols)
		for col := 0; col < numCols; col++ {
			maxColWidth := 0
			for row := 0; row < numRows; row++ {
				idx := col*numRows + row // Column-major indexing
				if idx < len(files) {
					if fileLengths[idx] > maxColWidth {
						maxColWidth = fileLengths[idx]
					}
				}
			}
			colWidths[col] = maxColWidth
		}

		// Calculate total width needed
		totalWidth := 0
		for i, width := range colWidths {
			totalWidth += width
			if i < len(colWidths)-1 {
				totalWidth += 2 // Space between columns
			}
		}

		// If this layout fits and uses fewer rows, use it
		if totalWidth <= termWidth && numRows < bestRows {
			bestCols = numCols
			bestRows = numRows
		}
	}

	// Format using the best layout
	numRows := (len(files) + bestCols - 1) / bestCols

	// Calculate column widths for the best layout (column-major ordering)
	colWidths := make([]int, bestCols)
	for col := 0; col < bestCols; col++ {
		maxColWidth := 0
		for row := 0; row < numRows; row++ {
			idx := col*numRows + row // Column-major indexing
			if idx < len(files) {
				if fileLengths[idx] > maxColWidth {
					maxColWidth = fileLengths[idx]
				}
			}
		}
		colWidths[col] = maxColWidth
	}

	var result strings.Builder
	for row := 0; row < numRows; row++ {
		for col := 0; col < bestCols; col++ {
			idx := col*numRows + row // Column-major indexing
			if idx < len(files) {
				result.WriteString(files[idx])

				// Add padding if not the last column in the row
				if col < bestCols-1 {
					padding := colWidths[col] - fileLengths[idx] + 2
					if padding < 2 {
						padding = 2 // Minimum 2 spaces
					}
					result.WriteString(strings.Repeat(" ", padding))
				}
			}
		}
		if row < numRows-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

func Print(paths []string, flags util.Flags) {
	outErrors := []string{}
	singleFiles := []string{}
	dirContents := []string{}
	content := []any{}

	// Handle recursive listing
	if flags.Recursive {
		allPaths, err := util.CollectDirectoriesRecursively(paths, flags)
		if err != nil {
			outErrors = append(outErrors, fmt.Sprintf("Error during recursive traversal: %v\n", err.Error()))
		} else {
			paths = allPaths
		}
	}

	multipleDirs := false
	if len(paths) > 1 || flags.Recursive {
		multipleDirs = true
	}

	for _, dirPath := range paths {
		info, err := util.IsValidDir(dirPath)
		if err != nil {
			outErrors = append(outErrors, fmt.Sprintf("Error: %v\n", err.Error()))
			continue
		}

		if !info.IsDir() {
			singleFiles = append(singleFiles, dirPath)
			continue
		}
		var files []string

		if flags.Longformat {
			files, err = util.ReadDirNamesLong(dirPath, flags)
			if err != nil {
				outErrors = append(outErrors, fmt.Sprintf("Error reading directory: %v\n", err.Error()))
				continue
			}
		} else {
			files, err = util.ReadDirNames(dirPath, flags)
			if err != nil {
				outErrors = append(outErrors, fmt.Sprintf("Error reading directory: %v\n", err.Error()))
				continue
			}
		}

		if multipleDirs {
			dirContents = append(dirContents, fmt.Sprintf("%v:", dirPath))
		}

		dirContents = append(dirContents, files...)
		content = append(content, dirContents)
		dirContents = []string{}
	}
	for _, err := range outErrors {
		fmt.Println(err)
	}

	for i, file := range singleFiles {
		if flags.Longformat {
			fmt.Println(file)
			continue
		}

		if i == len(singleFiles)-1 {
			fmt.Print(file + "\n\n")
			continue
		}
		fmt.Print(file + "  ")
	}

	for i, c := range content {
		if i != 0 {
			fmt.Println()
		}

		lines := c.([]string)
		if len(lines) == 0 {
			continue
		}

		// Print directory header (if present)
		if len(lines) > 0 && len(lines[0]) > 0 && lines[0][len(lines[0])-1] == ':' {
			fmt.Println(lines[0])
			lines = lines[1:] // Skip the header for content printing
		}

		// Print the directory contents
		if flags.Longformat {
			for _, line := range lines {
				fmt.Println(line)
			}
		} else {
			// Use column formatting for short format
			formatted := formatInColumns(lines)
			fmt.Print(formatted)
			if len(lines) > 0 {
				fmt.Println() // Add newline after the formatted output
			}
		}
	}
}
