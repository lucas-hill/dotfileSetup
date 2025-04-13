package setup

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TODO:
// Might want to consider returning an error here
func SetupZsh(repoDir string) {
	err := linkDotfiles(filepath.Join(repoDir, "zsh"))
	if err != nil {
		fmt.Printf("Error linking zsh files: %v\n", err)
	}

	zshrc := filepath.Join(os.Getenv("HOME"), ".zshrc")
	sourceLine := `for file in ~/.zsh/*.zsh; do [ -r "$file" ] && source "$file"; done`

	data, err := os.ReadFile(zshrc)
	if err == nil && !strings.Contains(string(data), sourceLine) {
		f, _ := os.OpenFile(zshrc, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		defer f.Close()

		f.WriteString("\n" + sourceLine + "\n")
		fmt.Println("Added source line to .zshrc")
	} else if err != nil {
		fmt.Printf("Could not read .zshrc: %v\n", err)
	}
}
