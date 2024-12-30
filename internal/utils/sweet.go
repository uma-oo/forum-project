// internal/utils/sweet.go
package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetFolderPath(baseDir, folderName string) (string, error) {
	var folderPath string

	// Walk through the directory hierarchy
	err := filepath.WalkDir(baseDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Check if the current item is a directory and matches the folder name
		if d.IsDir() && d.Name() == folderName {
			folderPath = path
			return filepath.SkipDir // Stop further traversal once found
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	if folderPath == "" {
		return "", fmt.Errorf("folder %q not found", folderName)
	}

	return folderPath, nil
}
