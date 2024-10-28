package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"my-ls/flags"
	"my-ls/format"
	"my-ls/listing"
)

func main() {
	// Parse flags
	options, err := flags.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "ls: %v\n", err)
		os.Exit(1)
	}

	// If no paths specified, use current directory
	if len(options.Paths) == 0 {
		options.Paths = append(options.Paths, ".")
	}

	// Handle multiple paths
	printMultiple := len(options.Paths) > 1
	hadError := false

	for i, path := range options.Paths {
		// Clean the path (handle multiple slashes)
		path = filepath.Clean(path)

		// Print path header if we're printing multiple paths
		if printMultiple {
			if i > 0 {
				fmt.Println()
			}
			fmt.Printf("%s:\n", path)
		}

		err := processPath(path, options, printMultiple)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ls: cannot access '%s': %v\n", path, err)
			hadError = true
			continue
		}
	}

	if hadError {
		os.Exit(1)
	}
}

func processPath(path string, options flags.Options, printHeader bool) error {
	// Get file info
	info, err := os.Lstat(path)
	if err != nil {
		return err
	}

	// Handle single file
	if !info.IsDir() {
		entries := []fs.DirEntry{dirEntryFromFileInfo(info, path)}
		fmt.Print(format.FormatOutput(entries, options))
		return nil
	}

	// Handle directory
	if options.Recursive {
		contents, err := listing.GetRecursiveDirectoryContents(path, options)
		if err != nil {
			return err
		}
		fmt.Print(format.FormatRecursiveOutput(contents, options))
	} else {
		entries, err := listing.GetDirectoryContents(path, options)
		if err != nil {
			return err
		}
		fmt.Print(format.FormatOutput(entries, options))
	}

	return nil
}

// Helper function to create DirEntry from FileInfo
func dirEntryFromFileInfo(info os.FileInfo, path string) fs.DirEntry {
	return &listing.FileInfoDirEntry{
		FileInfo: info,
		Path:     path,
	}
}
