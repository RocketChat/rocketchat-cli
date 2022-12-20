package matrix

import "path/filepath"

const (
	homeserverFile   = "homeserver.yaml"
	registrationFile = "registration.yaml"
)

type HomeserverType string

const (
	SynapseType  HomeserverType = "synapse"
	DendriteType                = "dendrite"
)

func getHomeserverFilePath(homeserverType HomeserverType, basePath string) (string, string) {
	dirPath := filepath.Join(basePath, "matrix", string(homeserverType))
	filePath := filepath.Join(dirPath, homeserverFile)

	return dirPath, filePath
}

func getRegistrationFilePath(homeserverType HomeserverType, basePath string) (string, string) {
	dirPath := filepath.Join(basePath, "matrix", string(homeserverType))
	filePath := filepath.Join(dirPath, registrationFile)

	return dirPath, filePath
}
