package native

import (
	"encoding/json"
	"net/http"
)

/*
install, uninstall, get binary path
*/

type VersionList struct {
	Versions []struct {
		Version          string             `json:"version"`
		Lts              bool               `json:"lts_release"`
		ReleaseCandidate bool               `json:"release_candidate"`
		Production       bool               `json:"production_release"`
		Downloads        []DownloadArtifact `json:"downloads"`
	} `json:"versions"`
}

type DownloadArtifact struct {
	Arch    string `json:"arch"`
	Edition string `json:"edition"`
	Target  string `json:"target"`
	Archive struct {
		SHA256 string `json:"sha256"`
		Url    string `json:"url"`
	} `json:"archive"`
}

const versionsUrl = "https://downloads.mongodb.org/full.json"

func AvailableVersions() (*VersionList, error) {
	response, err := http.Get(versionsUrl)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var versions VersionList

	err = json.NewDecoder(response.Body).Decode(&versions)
	if err != nil {
		return nil, err
	}

	return &versions, nil
}

func Download(artifact DownloadArtifact) error {
	return nil
}
