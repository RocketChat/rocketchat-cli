package docker

import "path/filepath"

const (
	composeFile = "docker-compose.yml"
)

func getComposeFilePath(basePath string) (string, string) {
	dirPath := basePath
	filePath := filepath.Join(dirPath, composeFile)

	return dirPath, filePath
}
