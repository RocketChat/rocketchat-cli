package traefik

import (
	"errors"
	"log"
	"os"
	"rc-cli/filesystem"
)

func WriteMiddlewaresFileIfNeeded(basePath string) {
	dirPath, filePath := getMiddlewaresFilePath(basePath)

	// Create the file if needed
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		// Ensure directory
		err := filesystem.EnsureDir(dirPath)
		if err != nil {
			log.Fatal("[WriteMiddlewaresFileIfNeeded] Could not create directories: ", err)
		}

		// Create the file
		err = os.WriteFile(filePath, []byte(middlewaresDefaults), 0644)
		if err != nil {
			log.Fatal("[WriteMiddlewaresFileIfNeeded] Could not write file: ", err)
		}
	}
}
