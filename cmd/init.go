package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	repoURL       string
	directoryName string
)

const (
	defaultRepo          = "https://github.com/lucas-hill/dotfiles"
	defaultDirectoryName = "dotfiles"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Clone and in itialize dotfiles from a repository",
	Run: func(cmd *cobra.Command, args []string) {
		reader := bufio.NewReader(os.Stdin)

		if repoURL == "" {
			fmt.Printf("Enter Git repo URL [%s]: ", defaultRepo)
			input, _ := reader.ReadString('\n')
			repoURL = strings.TrimSpace(input)
			if repoURL == "" {
				repoURL = defaultRepo
			}
		}

		if directoryName == "" {
			fmt.Printf("Enter local folder name [%s] ", defaultDirectoryName)
			input, _ := reader.ReadString('\n')
			directoryName = strings.TrimSpace(input)
			if directoryName == "" {
				directoryName = defaultDirectoryName
			}
		}

		cloneDir := filepath.Join(os.Getenv("HOME"), directoryName)
		fmt.Printf("Cloning %s -> %s\n", repoURL, cloneDir)

		if _, err := os.Stat(cloneDir); err == nil {
			fmt.Println("Folder already exists. Skipping clone")
		} else {
			cmd := exec.Command("git", "clone", repoURL, cloneDir)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				fmt.Printf("Failed to clone repo: %v\n", err)
				os.Exit(1)
			}
		}

		fmt.Println("Dotfiles cloned and ready in:", cloneDir)
	},
}
