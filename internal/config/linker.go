package config

import (
	"fmt"
	"os"
	"runtime"

	"gopkg.in/yaml.v3"
)

type Link struct {
	Destination string `yaml:"destination"`
}

type PackageLinker struct {
	Name  string   `yaml:"name"`
	Tags  []string `yaml:"tags"`
	Links []Link   `yaml:"links"`
}

type LinkerOSConfig struct {
	Darwin  []PackageLinker `yaml:"darwin"`
	Linux   []PackageLinker `yaml:"linux"`
	Windows []PackageLinker `yaml:"windows"`
	Default []PackageLinker `yaml:"default"`
}

func LoadLinkerConfig(path string) ([]PackageLinker, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var configData LinkerOSConfig
	if err = yaml.Unmarshal(data, &configData); err != nil {
		return nil, fmt.Errorf("failed to parse yaml: %w", err)
	}

	switch runtime.GOOS {
	case "darwin":
		return configData.Darwin, nil
	case "linux":
		return configData.Linux, nil
	case "windows":
		return configData.Windows, nil
	default:
		return configData.Default, nil
	}

}
