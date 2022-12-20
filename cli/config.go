package cli

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"rc-cli/filesystem"
)

type ConfigDef struct {
	Config struct {
		DevMode bool `yaml:"dev_mode"`
	} `yaml:"config"`
	RocketChat struct {
		Tag string `yaml:"tag"`
	} `yaml:"rocketchat"`
	Synapse struct {
		Tag string `yaml:"tag"`
	} `yaml:"synapse"`
	Traefik struct {
		Email string `yaml:"email"`
		Tag   string `yaml:"tag"`
	} `yaml:"traefik"`
	Redis struct {
		Tag string `yaml:"tag"`
	} `yaml:"redis"`
	Nginx struct {
		Tag string `yaml:"tag"`
	} `yaml:"nginx"`
	Element struct {
		Tag string `yaml:"tag"`
	} `yaml:"element"`
	Hostname string `yaml:"hostname"`
}

var Config ConfigDef

func GetExpectedURLs() []string {
	expectedUrls := []string{
		Config.Hostname,
		fmt.Sprintf("traefik.%s", Config.Hostname),
	}

	expectedUrls = append(expectedUrls, fmt.Sprintf("matrix.%s", Config.Hostname))
	expectedUrls = append(expectedUrls, fmt.Sprintf("synapse.%s", Config.Hostname))
	expectedUrls = append(expectedUrls, fmt.Sprintf("element.%s", Config.Hostname))

	return expectedUrls
}

func HasConfig() bool {
	return Config.Hostname != ""
}

func ReadConfigFile(basePath string) {
	_, filePath := getConfigFilePath(basePath)

	// Unmarshal defaults
	err := yaml.Unmarshal([]byte(configDefaults), &Config)
	if err != nil {
		log.Fatal("[ReadConfigFile] Could not unmarshal defaults: ", err)
	}

	// Load current file if exists
	if _, err := os.Stat(filePath); err == nil {
		currentFile, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal("[ReadConfigFile] Could not read file: ", err)
		}

		// Unmarshal current file
		err = yaml.Unmarshal(currentFile, &Config)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func WriteConfigFile(basePath string) {
	dirPath, filePath := getConfigFilePath(basePath)

	// Ensure directory
	err := filesystem.EnsureDir(dirPath)
	if err != nil {
		log.Fatal("[WriteConfigFile] Could not create directories: ", err)
	}

	// Marshal struct
	result, err := yaml.Marshal(Config)
	if err != nil {
		log.Fatal("[WriteConfigFile] Could not marshal struct: ", err)
	}

	// Write file
	err = os.WriteFile(filePath, result, 0644)
	if err != nil {
		log.Fatal("[WriteConfigFile] Could not write file: ", err)
	}
}
