package installers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lucas-hill/dotfileSetup/internal/setup"
)

// TODO:
// Figure out a better way to keep these in sync with the installPacke method
var supportedPackagesMac = map[string]bool{
	"zsh":       true,
	"docker":    true,
	"nvim":      true,
	"tmux":      true,
	"vscode":    true,
	"aerospace": true,
}

type macOSInstaller struct{}

func (i macOSInstaller) FilterSupportedPackages(packages []string) []string {

	var filtered []string
	for _, pkg := range packages {
		if _, ok := supportedPackagesMac[pkg]; ok {
			filtered = append(filtered, pkg)
		}
	}
	return filtered
}

func (i macOSInstaller) Install(packages []string, repoPath string) error {
	for _, pkg := range packages {
		err := i.installPackage(pkg, repoPath)
		if err != nil {
			return fmt.Errorf("failed to install '%s': %w", pkg, err)
		}
	}
	return nil
}

func (i macOSInstaller) installPackage(packageName, repoPath string) error {
	if !supportedPackagesMac[packageName] {
		return fmt.Errorf("package '%s' is not supported on macOS", packageName)
	}

	fmt.Printf("Installing %s on macOS\n", packageName)

	switch packageName {
	case "zsh":
		return setupZsh(repoPath)
	case "docker":
		return setupDocker()
	case "nvim":
		return setupNeovim(repoPath)
	case "tmux":
		return setupTmux(repoPath)
	case "aerospace":
		return setupAerospace(repoPath)
	default:
		return fmt.Errorf("installer not implemented for '%s'", packageName)
	}
}

func setupZsh(repoDir string) error {
	err := setup.LinkDotfiles(filepath.Join(repoDir, "zsh"))
	if err != nil {
		fmt.Printf("Error linking zsh files: %v\n", err)
		return err
	}

	zshrc := filepath.Join(os.Getenv("HOME"), ".zshrc")
	sourceLine := `for file in ~/.zsh/*.zsh; do [ -r "$file" ] && source "$file"; done`

	data, err := os.ReadFile(zshrc)
	if err == nil && !strings.Contains(string(data), sourceLine) {
		//TODO:
		//Find out what this is doing...
		f, _ := os.OpenFile(zshrc, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		defer f.Close()

		f.WriteString("\n" + sourceLine + "\n")
		fmt.Println("Added source line to .zshrc")
	} else if err != nil {
		fmt.Printf("Could not read .zshrc: %v\n", err)
		return err
	}
	return nil
}

func setupDocker() error {
	dockerPath := filepath.Join(os.Getenv("HOME"), "Downloads", "Docker.dmg")

	fmt.Println("Downloading Docker Desktop...")
	if err := runCommand("curl", "-L", "-o", dockerPath, "https://desktop.docker.com/mac/main/arm64/Docker.dmg"); err != nil {
		return fmt.Errorf("failed to download Docker.dmg: %w", err)
	}

	fmt.Println("Mounting the DMG...")
	if err := runCommand("hdiutil", "attach", dockerPath); err != nil {
		return fmt.Errorf("failed to mount Docker.dmg: %w", err)
	}

	fmt.Println("Running Docker Desktop installer...")
	if err := runCommand("sudo", "/Volumes/Docker/Docker.app/Contents/MacOS/install"); err != nil {
		return fmt.Errorf("failed to install Docker: %w", err)
	}

	fmt.Println("Unmounting the DMG...")
	if err := runCommand("hdiutil", "detach", "/Volumes/Docker"); err != nil {
		return fmt.Errorf("failed to unmount Docker.dmg: %w", err)
	}

	fmt.Println("Docker Desktop installed successfully.")
	return nil
}

func setupNeovim(repoPath string) error {
	fmt.Println("Downloading Neovim...")

	neovimURL := "https://github.com/neovim/neovim/releases/latest/download/nvim-macos-arm64.tar.gz"
	archivePath := filepath.Join(os.TempDir(), "nvim-macos.tar.gz")

	if err := runCommand("curl", "-L", "-o", archivePath, neovimURL); err != nil {
		return fmt.Errorf("Failed to download Neovim: %w", err)
	}

	//NOTE:
	//This is needed to make nvim available globally
	// installPath := filepath.Join(os.Getenv("HOME"), "bin", "nvim-macos-arm64")
	// symlinkPath := "/usr/local/bin/nvim"
	// nvimBinary := filepath.Join(installPath, "bin", "nvim")
	// if _, err := os.Lstat(symlinkPath); err == nil {
	// 	fmt.Println("Existing symlink found, skipping...")
	// } else {
	// 	if err := os.Symlink(nvimBinary, symlinkPath); err != nil {
	// 		return fmt.Errorf("Failed to symlink nvim to /user/local/bin: %w", err)
	// 	}
	// }
	// fmt.Println("Symlinked nvim -> /user/local/bin/nvim")

	nvimConfigSrc := filepath.Join(repoPath, "nvim")
	nvimConfigDest := filepath.Join(os.Getenv("HOME"), ".config", "nvim")
	if err := os.MkdirAll(filepath.Dir(nvimConfigDest), 0755); err != nil {
		return err
	}
	if err := os.Symlink(nvimConfigSrc, nvimConfigDest); err != nil && !os.IsExist(err) {
		return fmt.Errorf("Failed to symlink Neovim config: %w", err)
	}

	fmt.Println("Neovim installed\n")

	//NOTE:
	//Having the user manually sorce this on their own
	fmt.Println("🔧 Add the following to your ~/.zshrc or ~/.bashrc:")
	fmt.Println(`export PATH="$HOME/bin/nvim-macos/bin:$PATH"`)
	return nil
}

func setupVSCode(repoPath string) error {
	fmt.Println("Installing VSCode...")

	vscodeZip := filepath.Join(os.TempDir(), "vscode.zip")
	err := runCommand("curl", "-L", "-o", vscodeZip, "https://code.visualstudio.com/sha/download?build=stable&os=darwin-arm64")
	if err != nil {
		return fmt.Errorf("failed to download VSCode: %w", err)
	}

	err = runCommand("unzip", "-q", vscodeZip, "-d", "/Applications")
	if err != nil {
		return fmt.Errorf("failed to unzip VSCode: %w", err)
	}

	//NOTE:
	//Not really sure if this is the correct directory...
	vscodeBin := "/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code"
	if _, err := os.Stat(vscodeBin); err == nil {
		runCommand("ln", "-s", vscodeBin, "/usr/local/bin/code")
	}

	fmt.Println("✅ VSCode installed successfully.")
	return nil
}

func setupTmux(repoPath string) error {
	fmt.Println("Downloading tmux...")

	tmuxConfigDestPath := filepath.Join(os.Getenv("HOME"), "tmux.conf")
	tmuxConfigSrc := filepath.Join(repoPath, "tmux")

	if err := os.Symlink(tmuxConfigSrc, tmuxConfigDestPath); err != nil {
		return fmt.Errorf("Failed to symlink Tmux config: %w", err)
	}

	if err := runCommand("brew", "install", "tmux"); err != nil {
		return fmt.Errorf("Failed to download Tmux: %w", err)
	}
	fmt.Println("Tmux installed successfully.")
	return nil
}

func setupAerospace(repoPath string) error {
	fmt.Println("Downloading aerospace...")

	tmuxConfigDestPath := filepath.Join(os.Getenv("HOME"), ".config")
	tmuxConfigSrc := filepath.Join(repoPath, "aerospace")

	if err := os.Symlink(tmuxConfigSrc, tmuxConfigDestPath); err != nil {
		return fmt.Errorf("Failed to symlink Aerospace config: %w", err)
	}

	if err := runCommand("brew", "install", "aerospace"); err != nil {
		return fmt.Errorf("Failed to download Aerospace: %w", err)
	}
	fmt.Println("Aerospace installed successfully.")
	return nil
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin // so sudo can prompt if needed
	return cmd.Run()
}
