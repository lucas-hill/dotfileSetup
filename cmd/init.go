package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lucas-hill/dotfileSetup/internal/setup"
	"github.com/lucas-hill/dotfileSetup/internal/tui/textinput"
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
		if repoURL == "" {
			val, err := promptForInput("Enter Git repo URL to clone:", defaultRepo)
			if err != nil {
				fmt.Printf(" Failed to get repo input: %v", err)
				os.Exit(1)
			}
			repoURL = val
		}

		if directoryName == "" {
			val, err := promptForInput("Enter the local directory name to create:", defaultDirectoryName)
			if err != nil {
				fmt.Printf("Failed to get the directory name input %v\n", err)
				os.Exit(1)
			}
			directoryName = val
		}

		cloneDir := filepath.Join(os.Getenv("HOME"), directoryName)
		fmt.Printf("Cloning %s -> %s\n", repoURL, cloneDir)

		if _, err := os.Stat(cloneDir); err == nil {
			fmt.Println("Folder already exists. Skipping clone")
		} else {
			err := setup.CloneRepo(repoURL, cloneDir)
			if err != nil {
				fmt.Printf("Failed to clone repo: %v\n", err)
				os.Exit(1)
			}
		}

		fmt.Println("Dotfiles cloned and ready in:", cloneDir)
	},
}

// TODO:
// Extract out later if this is used more than just here
func promptForInput(prompt, defaultValue string) (string, error) {
	model := textinput.New(prompt, prompt, defaultValue)
	program := tea.NewProgram(model)
	finalModel, err := program.Run()
	if err != nil {
		return "", err
	}
	tiModel, ok := finalModel.(textinput.Model)
	if !ok {
		return "", fmt.Errorf("unexptedted model type")
	}

	if tiModel.Quit {
		fmt.Println("Cancelled by user.\nGoodbye...")
		os.Exit(0)
	}

	return tiModel.Entered, nil
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringVar(&repoURL, "repo", "", "Git repository URL")
	initCmd.Flags().StringVar(&directoryName, "dir", "", "Target directory name")
}
