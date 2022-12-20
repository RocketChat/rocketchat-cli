package traefik

import (
	"errors"
	"log"
	"os"
	"rc-cli/filesystem"
)

func WriteACMEFileIfNeeded(basePath string) {
	dirPath, filePath := getACMEFilePath(basePath)

	// Create the file if needed
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		// Ensure directory
		err := filesystem.EnsureDir(dirPath)
		if err != nil {
			log.Fatal("[WriteACMEFileIfNeeded] Could not create directories: ", err)
		}

		// Create the file
		err = os.WriteFile(filePath, []byte(acmeDefaults), 0644)
		if err != nil {
			log.Fatal("[WriteACMEFileIfNeeded] Could not write file: ", err)
		}
	}
}
