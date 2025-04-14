package setup

import (
	"os"
	"os/exec"
)

func CloneRepo(repoURL, storageDirectory string) error {
	cmd := exec.Command("git", "clone", repoURL, storageDirectory)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
