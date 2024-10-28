package listing

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetDirEntryForPath(path string) (os.DirEntry, error) {
	// Get the directory of the file
	dir := filepath.Dir(path)

	// Get the file name
	fileName := filepath.Base(path)

	// Read the directory to get all directory entries
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// Loop through the directory entries to find the one matching the file name
	for _, entry := range entries {
		if entry.Name() == fileName {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("DirEntry not found for path: %s", path)
}
