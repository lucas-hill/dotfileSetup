package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/lucas-hill/dotfileSetup/internal/installers"
	"github.com/lucas-hill/dotfileSetup/internal/setup"
	"github.com/lucas-hill/dotfileSetup/internal/tui/selectlist"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [components]",
	Short: "Install one or more dotfiles components",
	Args:  cobra.MinimumNArgs(0),
	Run: func(cmd *cobra.Command, args []string) {

		repoDir := filepath.Join(os.Getenv("HOME"), "dotfiles")

		//NOTE:
		//placing the installer here so we can filter out the packages we can access
		installer, err := installers.GetInstaller()
		if err != nil {
			os.Exit(1)
		}

		var selected []string
		if len(args) > 0 {
			selected = installer.FilterSupportedPackages(args)
		} else {
			packagesInRepo, err := setup.ListManagedPackages(repoDir)
			if err != nil {
				fmt.Printf("Packages to install could not be found %v\n", err)
			}
			availablePackages := installer.FilterSupportedPackages(packagesInRepo)
			selected, err = selectlist.MultiSelect("Select packages to setup", availablePackages)
			if err != nil {
				fmt.Printf("Selection failed: %v\n", err)
				os.Exit(1)
			}
		}

		err = installer.Install(selected, repoDir)
		if err != nil {
			fmt.Printf("Installation failed %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
}
