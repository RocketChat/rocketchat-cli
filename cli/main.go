package cli

import "path/filepath"

const (
	configFile = ".cli_config.yaml"
)

func getConfigFilePath(basePath string) (string, string) {
	dirPath := basePath
	filePath := filepath.Join(dirPath, configFile)

	return dirPath, filePath
}
