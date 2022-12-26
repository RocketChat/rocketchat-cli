package matrix

import (
	"crypto/sha256"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"rocketchat-cli/filesystem"
	"time"
)

type RegistrationDef struct {
	ID                             string `yaml:"id"`
	HSToken                        string `yaml:"hs_token"`
	ASToken                        string `yaml:"as_token"`
	URL                            string `yaml:"url"`
	SenderLocalpart                string `yaml:"sender_localpart"`
	DeSorunomeMsc2409PushEphemeral bool   `yaml:"de.sorunome.msc2409.push_ephemeral"`
	Namespaces                     struct {
		Users []struct {
			Exclusive bool   `yaml:"exclusive"`
			Regex     string `yaml:"regex"`
		} `yaml:"users"`
		Rooms []struct {
			Exclusive bool   `yaml:"exclusive"`
			Regex     string `yaml:"regex"`
		} `yaml:"rooms"`
		Aliases []struct {
			Exclusive bool   `yaml:"exclusive"`
			Regex     string `yaml:"regex"`
		} `yaml:"aliases"`
	} `yaml:"namespaces"`
	Rocketchat struct {
		HomeserverURL    string `yaml:"homeserver_url"`
		HomeserverDomain string `yaml:"homeserver_domain"`
	} `yaml:"rocketchat"`
}

var Registration RegistrationDef

func generateSHA256(prefix string) string {
	h := sha256.New()

	h.Write([]byte(fmt.Sprintf("%s_%d", prefix, time.Now().Unix())))

	bs := h.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

func GenerateID() string {
	return fmt.Sprintf("rocketchat_%s", generateSHA256("id"))
}

func GenerateHSToken() string {
	return generateSHA256("hs")
}

func GenerateASToken() string {
	return generateSHA256("as")
}

func ReadRegistrationFile(homeserverType HomeserverType, basePath string) {
	_, filePath := getRegistrationFilePath(homeserverType, basePath)

	// Unmarshal defaults
	err := yaml.Unmarshal([]byte(registrationDefaults), &Registration)
	if err != nil {
		log.Fatal("[ReadRegistrationFile] Could not unmarshal defaults: ", err)
	}

	// Load current file if exists
	if _, err := os.Stat(filePath); err == nil {
		currentFile, err := os.ReadFile(filePath)
		if err != nil {
			log.Fatal("[ReadRegistrationFile] Could not read file: ", err)
		}

		// Unmarshal current file
		err = yaml.Unmarshal(currentFile, &Registration)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func WriteRegistrationFile(homeserverType HomeserverType, basePath string) {
	dirPath, filePath := getRegistrationFilePath(homeserverType, basePath)

	// Ensure directory
	err := filesystem.EnsureDir(dirPath)
	if err != nil {
		log.Fatal("[WriteRegistrationFile] Could not create directories: ", err)
	}

	// Marshal struct
	result, err := yaml.Marshal(Registration)
	if err != nil {
		log.Fatal("[WriteRegistrationFile] Could not marshal struct: ", err)
	}

	// Write file
	err = os.WriteFile(filePath, result, 0644)
	if err != nil {
		log.Fatal("[WriteRegistrationFile] Could not write file: ", err)
	}
}
