package element

import "path/filepath"

const (
	configFile = "config.json"
)

func getConfigFilePath(basePath string) (string, string) {
	dirPath := filepath.Join(basePath, "matrix", "element")
	filePath := filepath.Join(dirPath, configFile)

	return dirPath, filePath
}
