package setup

import "os"

func ListManagedPackages(repoDir string) ([]string, error) {
	entries, err := os.ReadDir(repoDir)
	if err != nil {
		return nil, err
	}

	var packages []string

	for _, entry := range entries {
		if entry.IsDir() {
			packages = append(packages, entry.Name())
		}
	}

	return packages, nil
}
