package traefik

import (
	"errors"
	"log"
	"os"
	"rc-cli/filesystem"
)

func WriteRoutersFileIfNeeded(basePath string) {
	dirPath, filePath := getRoutersFilePath(basePath)

	// Create the file if needed
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		// Ensure directory
		err := filesystem.EnsureDir(dirPath)
		if err != nil {
			log.Fatal("[WriteRoutersFileIfNeeded] Could not create directories: ", err)
		}

		// Create the file
		err = os.WriteFile(filePath, []byte(routersDefaults), 0644)
		if err != nil {
			log.Fatal("[WriteRoutersFileIfNeeded] Could not write file: ", err)
		}
	}
}
