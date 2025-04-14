package setup

import (
	"fmt"
	"os"
	"path/filepath"
)

func LinkDotfiles(sourceDir string) error {
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		return err
	}

	for _, f := range files {
		src := filepath.Join(sourceDir, f.Name())
		dest := filepath.Join(os.Getenv("HOME"), f.Name())

		if _, err := os.Lstat(dest); err == nil {
			fmt.Printf("⚠️  %s already exists, skipping\n", dest)
			continue
		}

		if err := os.Symlink(src, dest); err != nil {
			fmt.Printf("Failed to link %s: %v\n", dest, err)
		} else {
			fmt.Printf("Linked %s → %s\n", dest, src)
		}
	}

	return nil
}
