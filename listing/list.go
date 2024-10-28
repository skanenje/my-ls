package listing

import (
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"my-ls/flags"
)

// FileInfoDirEntry implements fs.DirEntry using a FileInfo
type FileInfoDirEntry struct {
	FileInfo os.FileInfo
	Path     string
}

func (e *FileInfoDirEntry) Name() string {
	return e.FileInfo.Name()
}

func (e *FileInfoDirEntry) IsDir() bool {
	return e.FileInfo.IsDir()
}

func (e *FileInfoDirEntry) Type() fs.FileMode {
	return e.FileInfo.Mode().Type()
}

func (e *FileInfoDirEntry) Info() (fs.FileInfo, error) {
	return e.FileInfo, nil
}

func GetDirectoryContents(dir string, options flags.Options) ([]fs.DirEntry, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var result []fs.DirEntry

	if options.All {
		// Add . and .. to the beginning of the list
		currentDir, err := os.Stat(dir)
		if err != nil {
			return nil, err
		}
		parentDir, err := os.Stat(filepath.Dir(dir))
		if err != nil {
			return nil, err
		}

		dotEntry := &FileInfoDirEntry{FileInfo: currentDir, Path: filepath.Join(dir, ".")}
		dotDotEntry := &FileInfoDirEntry{FileInfo: parentDir, Path: filepath.Join(dir, "..")}
		result = append(result, dotEntry, dotDotEntry)
	}

	for _, entry := range entries {
		if options.All || entry.Name()[0] != '.' {
			result = append(result, entry)
		}
	}

	// Sort the entries, keeping . and .. at the beginning
	if options.SortTime {
		output := SorTimeSortt(result)
		return output, nil

	} else if options.Reverse {
		output := ReverseSort(result)
		return output, nil
	}
	output := Sort(result)

	return output, nil
}

func GetRecursiveDirectoryContents(rootDir string, options flags.Options) (map[string][]fs.DirEntry, error) {
	result := make(map[string][]fs.DirEntry)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Handle hidden files/directories
		if !options.All && info.Name()[0] == '.' && path != rootDir {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			entries, err := GetDirectoryContents(path, options)
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(rootDir, path)
			if err != nil {
				return err
			}

			// Apply sorting based on options
			if options.SortTime {
				entries = SorTimeSortt(entries)
			} else {
				entries = Sort(entries)
			}

			// Apply reverse sorting if needed
			if options.Reverse {
				entries = ReverseSort(entries)
			}

			if relPath == "." {
				result["."] = entries
			} else {
				result["./"+relPath] = entries
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	// Sort the keys to ensure consistent output order
	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	sortedResult := make(map[string][]fs.DirEntry)
	for _, k := range keys {
		sortedResult[k] = result[k]
	}

	return sortedResult, nil
}

func ReverseSort(result []fs.DirEntry) []fs.DirEntry {
	sort.Slice(result, func(i, j int) bool {
		// Always keep . first
		if result[i].Name() == "." {
			return true
		}
		if result[j].Name() == "." {
			return false
		}
		// Keep .. second
		if result[i].Name() == ".." {
			return true
		}
		if result[j].Name() == ".." {
			return false
		}
		// For all other entries, sort alphabetically ignoring case
		return strings.ToLower(strings.TrimPrefix(result[i].Name(), ".")) < strings.ToLower(strings.TrimPrefix(result[j].Name(), "."))
	})

	output := []fs.DirEntry{}
	for i := len(result) - 1; i >= 0; i-- {
		output = append(output, result[i])
	}
	return output
}

func SorTimeSortt(result []fs.DirEntry) []fs.DirEntry {
	sort.Slice(result, func(i, j int) bool {
		info1, _ := result[i].Info()
		info2, _ := result[j].Info()
		return info1.ModTime().String() > info2.ModTime().String()
	})
	return result
}

func Sort(result []fs.DirEntry) []fs.DirEntry {
	sort.Slice(result, func(i, j int) bool {
		// Always keep . first
		if result[i].Name() == "." {
			return true
		}
		if result[j].Name() == "." {
			return false
		}
		// Keep .. second
		if result[i].Name() == ".." {
			return true
		}
		if result[j].Name() == ".." {
			return false
		}
		// For all other entries, sort alphabetically ignoring case
		return strings.ToLower(strings.TrimPrefix(result[i].Name(), ".")) < strings.ToLower(strings.TrimPrefix(result[j].Name(), "."))
	})
	return result
}
