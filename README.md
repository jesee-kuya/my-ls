# my-ls

A Go implementation of the Unix/Linux `ls` command with support for various display options and formatting.

## Overview

`my-ls` is a command-line utility that lists directory contents with similar functionality to the standard Unix/Linux `ls` command. It provides colored output for different file types and supports various flags to customize the display.

## Features

- List files and directories with proper formatting
- Colored output for different file types:
  - Directories: bold blue
  - Executable files: bold green
  - Symbolic links: bold cyan
  - Socket files: bold magenta
  - Pipes: yellow on black background
  - Device files: bold yellow on black
  - Archive files: bold red
- Support for various display options:
  - `-a`: Show all files, including hidden files (those starting with a dot)
  - `-l`: Use long listing format with detailed file information
  - `-r`: Reverse the order of the sort
  - `-R`: List subdirectories recursively
  - `-t`: Sort by modification time, newest first

## Installation

### Prerequisites

- Go 1.24.3 or higher

### Building from Source

1. Clone the repository:
   ```
   git clone https://github.com/jesee-kuya/my-ls.git
   cd my-ls
   ```

2. Build the project:
   ```
   go build
   ```

3. Install the binary (optional):
   ```
   go install
   ```

## Usage

```
my-ls [OPTIONS] [FILE/DIRECTORY...]
```

### Examples

List files in the current directory:
```
my-ls
```

Show all files including hidden ones:
```
my-ls -a
```

Use long listing format:
```
my-ls -l
```

Combine options:
```
my-ls -la
```

List files in a specific directory:
```
my-ls /path/to/directory
```

List files recursively:
```
my-ls -R /path/to/directory
```

Sort files by modification time (newest first):
```
my-ls -t
```

Reverse the sort order:
```
my-ls -r
```

## Project Structure

- `main.go`: Entry point of the application, handles command-line arguments
- `print/`: Contains code for displaying file listings
  - `print.go`: Handles the formatting and printing of file listings
- `util/`: Contains utility functions
  - `readDir.go`: Core functionality for reading directory contents
  - `sorted.go`: Functions for sorting file listings
  - `time.go`: Time-related utilities
  - `stripAnsi.go`: Functions for handling ANSI color codes
  - `isValidDir.go`: Directory validation
  - `hasAnsi.go`: Detection of ANSI escape sequences
  - `reverse.go`: Functions for reversing lists

## Testing

Run the tests with:
```
go test ./...
```

## License

This project is licensed under the [MIT License](LICENSE).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.