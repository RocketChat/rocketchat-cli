package traefik

import (
	"path/filepath"
)

const (
	acmeFileName        = "acme.json"
	middlewaresFileName = "middlewares.yml"
	routersFileName     = "routers.yml"
	traefikFileName     = "traefik.yml"
)

func getACMEFilePath(basePath string) (string, string) {
	dirPath := filepath.Join(basePath, "traefik")
	filePath := filepath.Join(dirPath, acmeFileName)

	return dirPath, filePath
}

func getMiddlewaresFilePath(basePath string) (string, string) {
	dirPath := filepath.Join(basePath, "traefik", "config")
	filePath := filepath.Join(dirPath, middlewaresFileName)

	return dirPath, filePath
}

func getRoutersFilePath(basePath string) (string, string) {
	dirPath := filepath.Join(basePath, "traefik", "config")
	filePath := filepath.Join(dirPath, routersFileName)

	return dirPath, filePath
}

func getTraefikFilePath(basePath string) (string, string) {
	dirPath := filepath.Join(basePath, "traefik")
	filePath := filepath.Join(dirPath, traefikFileName)

	return dirPath, filePath
}
