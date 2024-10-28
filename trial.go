package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

// Get the block size used by the file
func getBlocks(info os.FileInfo) int64 {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		return stat.Blocks // Get blocks used by the file
	}
	// Defaulting to size / 512 for fallback (1 block = 512 bytes)
	return (info.Size() + 511) / 512
}

func mainn() {
	dir := "." // Starting directory

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// We only care about regular files and directories for block totals
		if info.Mode().IsRegular() || info.IsDir() {
			blocks := getBlocks(info)
			fmt.Printf("%-40s %6d blocks\n", path, blocks)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path: %v\n", err)
	}
}
