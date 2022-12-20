package traefik

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"rc-cli/filesystem"
)

type ConfigDef struct {
	EntryPoints struct {
		Web struct {
			Address string `yaml:"address"`
		} `yaml:"web"`
		WebSecure struct {
			Address string `yaml:"address"`
		} `yaml:"web-secure"`
	} `yaml:"entryPoints"`
	API struct {
		Dashboard bool `yaml:"dashboard"`
		Insecure  bool `yaml:"insecure"`
	} `yaml:"api"`
	Providers struct {
		File struct {
			Directory string `yaml:"directory"`
			Watch     bool   `yaml:"watch"`
		} `yaml:"file"`
		Docker struct {
			Endpoint         string `yaml:"endpoint"`
			Network          string `yaml:"network"`
			Watch            bool   `yaml:"watch"`
			ExposedByDefault bool   `yaml:"exposedByDefault"`
		} `yaml:"docker"`
	} `yaml:"providers"`
	CertificatesResolvers struct {
		LetsEncrypt struct {
			ACME struct {
				CAServer      string `yaml:"caServer"`
				Email         string `yaml:"email"`
				Storage       string `yaml:"storage"`
				HttpChallenge struct {
					EntryPoint string `yaml:"entryPoint"`
				} `yaml:"httpChallenge"`
			} `yaml:"acme"`
		} `yaml:"letsencrypt"`
	} `yaml:"certificatesResolvers"`
}

var Config ConfigDef

func ReadConfigFile(basePath string) {
	_, filePath := getTraefikFilePath(basePath)

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
	dirPath, filePath := getTraefikFilePath(basePath)

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
