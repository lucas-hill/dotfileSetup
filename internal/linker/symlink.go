package linker

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SymlinkTarget struct {
	Local   string
	Dotfile string
}

// TODO:
// Break up the parts of this function
func CreateSymlinks(packages []string, repoDir string) error {

	for _, pkg := range packages {
		packageDir := filepath.Join(repoDir, pkg)
		fmt.Printf("packageDirectory: %s\n", packageDir)

		files, err := listAllFilesRecursive(packageDir)
		if err != nil {
			return err
		}
		fmt.Printf("files: %s'\n", files)

		var filesToLink []SymlinkTarget
		for _, file := range files {
			localPath, err := repoPathToLocalPath(file, packageDir)
			if err != nil {
				return err
			}

			localPath = filepath.Join(os.Getenv("HOME"), localPath)
			if fileExists(localPath) {
				fmt.Printf("File already exists on this machine. Skipping: %s...\n", localPath)
				continue
			}
			filesToLink = append(filesToLink, SymlinkTarget{Local: localPath, Dotfile: file})
		}

		//NOTE:
		//Here we know the files don't exist on the local machine.
		//Now we need to find the first directory for each that doesn't exist
		//Don't want to try and link files twice...
		symlinkSet := make(map[string]struct{})
		for _, target := range filesToLink {
			fmt.Println("local: ", target.Local)
			path, err := deepestNonExistentPath(target.Local)
			if err != nil {
				symlinkSet[target.Local] = struct{}{}
				continue
			}
			symlinkSet[path] = struct{}{}
		}

		fmt.Println("symlinkSet:", symlinkSet)

		//Now we have unique list of destination paths we can use to create our symlinks
		// For our craete statement we need a source and a destination
		// the destination is our symlinkSet
		// the source is from our dotfiles that need to be converted back from the local version
		type symlink struct {
			Source      string
			Destination string
		}
		var allSymlinks []symlink
		for link, _ := range symlinkSet {
			relPath, err := filepath.Rel(os.Getenv("HOME"), link)
			if err != nil {
				fmt.Println("Error getting the relative path")
			}
			newLink := symlink{}
			newLink.Destination = link
			newLink.Source = filepath.Join(packageDir, relPath)
			allSymlinks = append(allSymlinks, newLink)
		}
		fmt.Println("symlinks:", allSymlinks)
	}

	return nil
}

// NOTE:
// Might rename this...
// It is actually getting the deepest non existant path + 1 to preserve the directry that needs to be sourced
func deepestNonExistentPath(path string) (string, error) {
	if fileExists(path) {
		return path, nil
	}
	current := path
	previous := path
	for {
		parent := filepath.Dir(current)
		if fileExists(current) {
			return previous, nil
		}
		previous = current
		current = parent
		fmt.Println("moving to parent directory", current)
	}
}

// func deepestNonExistentPath(path string) (string, error) {
// 	current := path
// 	for {
// 		if !fileExists(current) {
// 			current = filepath.Dir(current)
// 			fmt.Println("moving to parent directory", current)
// 			continue
// 		}
// 		break
// 	}
// 	return current, nil
// }

func fileExists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil || !os.IsNotExist(err)
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

func repoPathToLocalPath(fullFilePath, packageDir string) (string, error) {
	fullFilePath = filepath.Clean(fullFilePath)
	packageDir = filepath.Clean(packageDir)

	if !strings.HasPrefix(fullFilePath, packageDir) {
		return "", fmt.Errorf("file pah %q does not start with expeted root %q", fullFilePath, packageDir)
	}
	rel := strings.TrimPrefix(fullFilePath, packageDir)

	return rel, nil
}

// NOTE:
// This will get the first available file on the local machine that we can create a symlink for
// or it will return an error if the file already exists on the machine.
// func getCommonLocalFiles(paths []string, packagePath string) ([]string, error) {
// 	localDir := filepath.Join(os.Getenv("HOME"), packagePath)
//
// 	return existingLocalFiles, nil
// }

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
// func splitByDeepestCommonDirectories(files []string) ([]string, error) {
// 	filesSplitByDirectory := make(map[string][]string, len(files))
//
// 	// Step 1: Pre-split all files
// 	for i, file := range files {
// 		filesSplitByDirectory[strconv.Itoa(i)] = splitPathParts(file)
// 	}
//
// 	sharedDirs := map[string]struct{}{}
//
// 	for i := 0; i < len(files); i++ {
// 		for j := i + 1; j < len(files); j++ {
// 			partsA := filesSplitByDirectory[strconv.Itoa(i)]
// 			partsB := filesSplitByDirectory[strconv.Itoa(j)]
//
// 			shared := commonPrefixParts(partsA, partsB)
// 			if len(shared) > 0 {
// 				sharedPath := filepath.Join(shared...)
// 				sharedDirs[sharedPath] = struct{}{}
// 			}
// 		}
// 	}
//
// 	// Step 3: Return as a sorted slice
// 	var result []string
// 	for dir := range sharedDirs {
// 		result = append(result, dir)
// 	}
// 	sort.Strings(result)
// 	return result, nil
// }

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
