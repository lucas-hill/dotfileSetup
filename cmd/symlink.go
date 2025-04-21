package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lucas-hill/dotfileSetup/internal/config"
	"github.com/lucas-hill/dotfileSetup/internal/linker"
	"github.com/lucas-hill/dotfileSetup/internal/setup"
	"github.com/lucas-hill/dotfileSetup/internal/tui/selectlist"
	"github.com/spf13/cobra"
)

var symlinkCmd = &cobra.Command{
	Use:   "link",
	Short: "Create a symlink for the packages you select",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		//TODO:
		//here I want to show a display for the packages that need to be linked.
		//From there I want the user to do a multi-select to pick what to link
		//Then read the config file to get the locations to create the symlinks
		//Create the acutal symlinks and handle conflicts

		configPath := "config/linker.yaml"
		symlinkData, err := config.LoadLinkerConfig(configPath)
		if err != nil {
			os.Exit(1)
		}

		//TODO:
		//Come up with a better way to handle this.
		//Hard coding for now...
		repoDir := filepath.Join(os.Getenv("HOME"), "dotfiles")

		packagesInDir, err := setup.ListManagedPackages(repoDir)
		if err != nil {
			fmt.Printf("Selection failed: %v\n", err)
		}

		linkMap := map[string]config.SymlinkItem{}
		for _, p := range symlinkData {
			linkMap[p.Name] = p
		}

		var availableLinks []string
		for _, pkg := range packagesInDir {
			if _, ok := linkMap[pkg]; ok {
				availableLinks = append(availableLinks, pkg)
			}
		}

		selected, err := selectlist.MultiSelect("Select packages to setup", availableLinks)
		if err != nil {
			fmt.Printf("Selection failed: %v\n", err)
			os.Exit(1)
		}

		//TODO:
		//Should this have the linkMap passed in???
		err = linker.CreateSymlinks(selected, &linkMap)
		if err != nil {
			fmt.Printf("Symlink failed: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(symlinkCmd)
}
