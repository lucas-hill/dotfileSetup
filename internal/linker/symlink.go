package linker

import (
	"fmt"
	"os"
)

// NOTE:
// For now I'm just skipping a file if the user already has something in this directory...
// The user will be responsible for handling a conflict for now
func Create(source, destination string) error {

	_, err := os.Lstat(source)
	if err != nil {
		fmt.Printf("This filepath is already taken %s. Skipping...", source)
		return err
	}

	if err := os.Symlink(source, destination); err != nil {
		fmt.Printf("Failed to link %s: %v\n", destination, err)
		return err
	}

	fmt.Printf("Linked %s → %s\n", destination, source)
	return nil
}

// TODO:
// Not sure what should happen here.
// Could be on a per path basis?
// Don't want to completely nuke someone's zsh file...
// Could be they are appended???
// Would need to check for duplicate aliases and funciton names to avoid conflicts
// func MergeConflictingFiles() {
//
// }

// NOTE:
// Not sure this should be here...
// func LoadOSConfig(path string, currentOS string) (map[string]string, error) {
//     fileContents, err := os.ReadFile(path)
//     if(err != nil){
//         fmt.Errorf("Failed to read the config file")
//         return nil, err
//     }
//
// }

// func CreateIfMissing(source, destination string) error {
//
// }
//
// func ForceCreate(source, destination string) error {
//
// }
//
// func Validate(destination, expectedSource string) (bool, error) {
//
// }
