package installers

import (
	"fmt"
	"runtime"
)

type Installer interface {
	Install(packages []string, repoPath string) error
	FilterSupportedPackages(packages []string) []string
	installPackage(packageName, repoPath string) error
}

func GetInstaller() (Installer, error) {
	switch runtime.GOOS {
	case "darwin":
		return macOSInstaller{}, nil
	case "linux":
		return nil, fmt.Errorf("Linux is not currently supported")
	case "windows":
		return nil, fmt.Errorf("Windows is not currently supported")
	default:
		return nil, fmt.Errorf("Unsupported Platform: %s", runtime.GOOS)
	}
}
