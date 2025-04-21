package linker

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func CreateSymlinks(packages []string, repoDir string) error {
	//NOTE:
	//Need access to the source and the destination.
	//source comes from the available packages

	for _, pkg := range packages {
		packageDir := filepath.Join(repoDir, pkg)
	}

	return nil
}

// NOTE:
// For now I'm just skipping a file if the user already has something in this directory...
// The user will be responsible for handling a conflict for now
func create(source, destination string) error {

	_, err := os.Lstat(source)
	if err != nil {
		fmt.Printf("This filepath is already taken %s. Skipping...", source)
		return err
	}

	if err := os.Symlink(source, destination); err != nil {
		fmt.Printf("Failed to link %s: %v\n", destination, err)
		return err
	}

	fmt.Printf("Linked %s → %s\n", destination, source)
	return nil
}

func listAllFilesRecursive(packageDir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(packageDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, err
}

func splitPathParts(path string) []string {
	cleaned := filepath.Clean(path)
	return strings.Split(cleaned, string(filepath.Separator))
}

func commonPrefixParts(a, b []string) []string {
	length := min(len(a), len(b))
	var result []string
	for i := 0; i < length; i++ {
		if a[i] != b[i] {
			break
		}
		result = append(result, a[i])
	}
	return result
}

// TODO:
// This needs to start from the repodir/packageName
// No reason to start before this
func splitByDeepestCommonDirectories(files []string) ([]string, error) {
	filesSplitByDirectory := make(map[string][]string, len(files))

	// Step 1: Pre-split all files
	for i, file := range files {
		filesSplitByDirectory[strconv.Itoa(i)] = splitPathParts(file)
	}

	sharedDirs := map[string]struct{}{}

	for i := 0; i < len(files); i++ {
		for j := i + 1; j < len(files); j++ {
			partsA := filesSplitByDirectory[strconv.Itoa(i)]
			partsB := filesSplitByDirectory[strconv.Itoa(j)]

			shared := commonPrefixParts(partsA, partsB)
			if len(shared) > 0 {
				sharedPath := filepath.Join(shared...)
				sharedDirs[sharedPath] = struct{}{}
			}
		}
	}

	// Step 3: Return as a sorted slice
	var result []string
	for dir := range sharedDirs {
		result = append(result, dir)
	}
	sort.Strings(result)
	return result, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func commonPathPrefix(a, b string) string {
	aParts := strings.Split(filepath.Clean(a), string(filepath.Separator))
	bParts := strings.Split(filepath.Clean(b), string(filepath.Separator))

	length := min(len(aParts), len(bParts))
	var shared []string
	for i := 0; i < length; i++ {
		if aParts[i] != bParts[i] {
			break
		}
		shared = append(shared, aParts[i])
	}
	return filepath.Join(shared...)
}

// TODO:
// Not sure what should happen here.
// Could be on a per path basis?
// Don't want to completely nuke someone's zsh file...
// Could be they are appended???
// Would need to check for duplicate aliases and funciton names to avoid conflicts
// func MergeConflictingFiles() {
//
// }
