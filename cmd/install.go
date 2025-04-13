package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lucas-hill/dotfileSetup/internal/setup"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [components]",
	Short: "Install one or more dotfiles components",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		for _, component := range args {
			fmt.Printf("Installing %s...\n", component)
			repoDir := filepath.Join(os.Getenv("HOME"), "dotfiles")

			switch component {
			case "zsh":
				setup.SetupZsh(repoDir)
			case "docker":
				setup.SetupDocker()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
