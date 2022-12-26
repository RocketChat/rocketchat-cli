package docker

import (
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"rocketchat-cli/filesystem"
)

type ComposeDef struct {
	Services struct {
		Traefik struct {
			Image    string   `yaml:"image"`
			Restart  string   `yaml:"restart"`
			Ports    []string `yaml:"ports"`
			Volumes  []string `yaml:"volumes"`
			Labels   []string `yaml:"labels"`
			Networks []string `yaml:"networks"`
		} `yaml:"traefik"`
		Postgres struct {
			Image       string `yaml:"image"`
			Restart     string `yaml:"restart"`
			Environment struct {
				PostgresPassword   string `yaml:"POSTGRES_PASSWORD"`
				PostgresUser       string `yaml:"POSTGRES_USER"`
				PostgresDB         string `yaml:"POSTGRES_DB"`
				PostgresInitDbArgs string `yaml:"POSTGRES_INITDB_ARGS"`
			} `yaml:"environment"`
			Volumes  []string `yaml:"volumes"`
			Networks []string `yaml:"networks"`
		} `yaml:"postgres"`
		Redis struct {
			Image    string   `yaml:"image"`
			Restart  string   `yaml:"restart"`
			Networks []string `yaml:"networks"`
		} `yaml:"redis"`
		Synapse struct {
			Image       string `yaml:"image"`
			Restart     string `yaml:"restart"`
			Environment struct {
				SynapseConfigDid  string `yaml:"SYNAPSE_CONFIG_DIR"`
				SynapseConfigPath string `yaml:"SYNAPSE_CONFIG_PATH"`
				UID               string `yaml:"UID"`
				GID               string `yaml:"GID"`
				TZ                string `yaml:"TZ"`
			} `yaml:"environment"`
			Volumes  []string `yaml:"volumes"`
			Ports    []string `yaml:"ports"`
			Labels   []string `yaml:"labels"`
			Networks []string `yaml:"networks"`
		} `yaml:"synapse"`
		Nginx struct {
			Image    string   `yaml:"image"`
			Restart  string   `yaml:"restart"`
			Volumes  []string `yaml:"volumes"`
			Labels   []string `yaml:"labels"`
			Networks []string `yaml:"networks"`
		} `yaml:"nginx"`
		Element struct {
			Image    string   `yaml:"image"`
			Restart  string   `yaml:"restart"`
			Volumes  []string `yaml:"volumes"`
			Labels   []string `yaml:"labels"`
			Networks []string `yaml:"networks"`
		} `yaml:"element"`
		Rocketchat struct {
			Image       string   `yaml:"image"`
			Command     string   `yaml:"command"`
			Restart     string   `yaml:"restart"`
			Volumes     []string `yaml:"volumes"`
			Environment []string `yaml:"environment"`
			DependsOn   []string `yaml:"depends_on"`
			Ports       []string `yaml:"ports"`
			Networks    []string `yaml:"networks"`
		} `yaml:"rocketchat"`
		Mongodb struct {
			Image    string   `yaml:"image"`
			Restart  string   `yaml:"restart"`
			Volumes  []string `yaml:"volumes"`
			Command  string   `yaml:"command"`
			Labels   []string `yaml:"labels"`
			Networks []string `yaml:"networks"`
		} `yaml:"mongodb"`
		MongoInitReplica struct {
			Image     string   `yaml:"image"`
			Command   string   `yaml:"command"`
			DependsOn []string `yaml:"depends_on"`
			Networks  []string `yaml:"networks"`
		} `yaml:"mongo-init-replica"`
	} `yaml:"services"`
	Networks struct {
		Internal struct {
			Attachable bool `yaml:"attachable"`
		} `yaml:"internal"`
	} `yaml:"networks"`
}

var Compose ComposeDef

func ReadComposeFile(basePath string) {
	_, filePath := getComposeFilePath(basePath)

	// Unmarshal defaults
	err := yaml.Unmarshal([]byte(composeDefaults), &Compose)
	if err != nil {
		log.Fatal("[Docker][ReadComposeFile] Could not unmarshal defaults: ", err)
	}

	// Load current file if exists
	if _, err := os.Stat(filePath); err == nil {
		currentFile, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal("[Docker][ReadComposeFile] Could not read file: ", err)
		}

		// Unmarshal current file
		err = yaml.Unmarshal(currentFile, &Compose)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func WriteComposeFile(basePath string) {
	dirPath, filePath := getComposeFilePath(basePath)

	// Ensure directory
	err := filesystem.EnsureDir(dirPath)
	if err != nil {
		log.Fatal("[Docker][WriteComposeFile] Could not create directories: ", err)
	}

	// Marshal struct
	result, err := yaml.Marshal(Compose)
	if err != nil {
		log.Fatal("[Docker][WriteComposeFile] Could not marshal struct: ", err)
	}

	// Write file
	err = os.WriteFile(filePath, result, 0644)
	if err != nil {
		log.Fatal("[Docker][WriteComposeFile] Could not write file: ", err)
	}
}
