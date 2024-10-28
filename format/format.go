package format
import (
	"fmt"
	"io/fs"
	"os"
	"os/user"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
	"my-ls/flags"
)
func IsExecutable(file os.DirEntry) bool {
	// Get the file info
	info, err := file.Info()
	if err != nil {
		return false
	}
	// Check if it's a regular file and if it's executable
	return info.Mode().Perm()&0o111 != 0 && !info.IsDir()
}
func FormatOutput(entries []fs.DirEntry, options flags.Options) string {
	var output strings.Builder
	if options.Long {
		for _, entry := range entries {
			output.WriteString(formatLongListing(entry))
			output.WriteString("\n")
		}
	} else {
		for i, entry := range entries {
			if entry.IsDir() {
				output.WriteString("\033[1m\033[34m" + entry.Name() + "\033[0m")
			} else {
				if strings.HasSuffix(entry.Name(), ".zip") || strings.HasSuffix(entry.Name(), ".tar") || strings.HasSuffix(entry.Name(), ".gz") || strings.HasSuffix(entry.Name(), ".tar.gz") || strings.HasSuffix(entry.Name(), ".tgz") || strings.HasSuffix(entry.Name(), ".bz2") || strings.HasSuffix(entry.Name(), ".tar.bz2") || strings.HasSuffix(entry.Name(), ".tbz") {
					output.WriteString("\033[1m\033[31m" + entry.Name() + "\033[0m")
				} else if IsExecutable(entry) {
					output.WriteString("\033[1m\033[32m" + entry.Name() + "\033[0m")
				} else {
					output.WriteString(entry.Name())
				}
			}
			if i < len(entries)-1 {
				output.WriteString("  ")
			}
		}
		output.WriteString("\n")
	}
	return output.String()
}
func FormatRecursiveOutput(contents map[string][]fs.DirEntry, options flags.Options) string {
	var output strings.Builder
	type dirInfo struct {
		path    string
		modTime time.Time
	}
	dirs := make([]dirInfo, 0, len(contents))
	for dir := range contents {
		info, err := os.Stat(dir)
		if err != nil {
			// Handle error or use a default time
			dirs = append(dirs, dirInfo{path: dir, modTime: time.Time{}})
		} else {
			dirs = append(dirs, dirInfo{path: dir, modTime: info.ModTime()})
		}
	}
	// Sort directories
	sort.Slice(dirs, func(i, j int) bool {
		if options.SortTime {
			return dirs[i].modTime.After(dirs[j].modTime)
		}
		return dirs[i].path < dirs[j].path
	})
	// Reverse directory order if -r flag is set
	if options.Reverse {
		for i, j := 0, len(dirs)-1; i < j; i, j = i+1, j-1 {
			dirs[i], dirs[j] = dirs[j], dirs[i]
		}
	}
	// Process root directory first
	if rootEntries, ok := contents["."]; ok {
		output.WriteString(".:")
		sortAndWriteEntries(&output, rootEntries, options)
	}
	// Process other directories
	for _, dir := range dirs {
		if dir.path == "." {
			continue // We've already processed the root directory
		}
		output.WriteString("\n")
		output.WriteString("\n")
		output.WriteString(dir.path + ":")
		entries := contents[dir.path]
		sortAndWriteEntries(&output, entries, options)
	}
	return strings.TrimSpace(output.String()) + "\n"
}
func sortAndWriteEntries(output *strings.Builder, entries []fs.DirEntry, options flags.Options) {
	// Sort entries
	sort.Slice(entries, func(i, j int) bool {
		if options.SortTime {
			infoI, _ := entries[i].Info()
			infoJ, _ := entries[j].Info()
			return infoI.ModTime().After(infoJ.ModTime())
		}
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})
	// Reverse entry order if -r flag is set
	if options.Reverse {
		for i, j := 0, len(entries)-1; i < j; i, j = i+1, j-1 {
			entries[i], entries[j] = entries[j], entries[i]
		}
	}
	// Write entries
	output.WriteString("\n")
	for i, entry := range entries {
		if i > 0 {
			output.WriteString("  ")
		}
		if entry.IsDir() {
			output.WriteString("\033[1m\033[34m" + entry.Name() + "\033[0m")
		} else {
			name := entry.Name()
			if strings.HasSuffix(name, ".zip") || strings.HasSuffix(name, ".tar") ||
				strings.HasSuffix(name, ".gz") || strings.HasSuffix(name, ".tar.gz") ||
				strings.HasSuffix(name, ".tgz") || strings.HasSuffix(name, ".bz2") ||
				strings.HasSuffix(name, ".tar.bz2") || strings.HasSuffix(name, ".tbz") {
				output.WriteString("\033[1m\033[31m" + name + "\033[0m")
			} else if IsExecutable(entry) {
				output.WriteString("\033[1m\033[32m" + name + "\033[0m")
			} else {
				output.WriteString(name)
			}
		}
	}
}
// Keep the existing formatLongListing function...
func formatLongListing(entry fs.DirEntry) string {
	info, err := entry.Info()
	if err != nil {
		return fmt.Sprintf("Error getting file info: %v", err)
	}
	mode := info.Mode().String()
	links := getLinks(info) // Placeholder for now
	owner := getOwner(info)
	group := getGroup(info)
	size := strconv.FormatInt(info.Size(), 10)
	modTime := info.ModTime().Format("Jan _2 15:04")
	name := entry.Name()
	return fmt.Sprintf("%s %d %s %s %6s %s %s", mode, links, owner, group, size, modTime, name)
}
func getOwner(info fs.FileInfo) string {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		if user, err := user.LookupId(strconv.Itoa(int(stat.Uid))); err == nil {
			return user.Username
		}
	}
	return "?"
}
func getLinks(info os.FileInfo) uint64 {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		return stat.Nlink
	}
	return 1
}
func getGroup(info fs.FileInfo) string {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		if group, err := user.LookupGroupId(strconv.Itoa(int(stat.Gid))); err == nil {
			return group.Name
		}
	}
	return "?"
}

func calculateTotalBlocks(entries []fs.DirEntry) int64 {
	var total int64
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		if stat, ok := info.Sys().(*syscall.Stat_t); ok {
			total += stat.Blocks / 2 // Convert 512-byte blocks to 1024-byte blocks
		}
	}
	return total
}
