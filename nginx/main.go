package nginx

import "path/filepath"

const (
	matrixConfFile = "matrix.conf"
)

func getMatrixConfFilePath(basePath string) (string, string) {
	dirPath := filepath.Join(basePath, "nginx")
	filePath := filepath.Join(dirPath, matrixConfFile)

	return dirPath, filePath
}
