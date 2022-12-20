package docker_images

import (
	"log"
	"regexp"
	"strings"
)

func GetSynapseImageVersions(dev bool) (string, []string) {
	versions := GetImageVersions("matrixdotorg/synapse")

	var latestStableVersion string
	var versionsList []string

	for _, version := range versions.Results {
		var err error

		// Skip latest and develop
		if strings.Contains(version.Name, "latest") || strings.Contains(version.Name, "develop") {
			continue
		}

		// Skip dev versions
		rcMatched := false
		devMatched := false
		if !dev {
			rcMatched, err = regexp.MatchString("rc[0-9]", version.Name)
			if err != nil {
				log.Fatal(err)
			}
			if rcMatched {
				continue
			}

			devMatched, err = regexp.MatchString("dev[0-9]", version.Name)
			if err != nil {
				log.Fatal(err)
			}
			if devMatched {
				continue
			}
		}

		// Set latest stable imageTag
		if latestStableVersion == "" && !rcMatched && !devMatched {
			latestStableVersion = version.Name
		}

		versionsList = append(versionsList, version.Name)
	}

	return latestStableVersion, versionsList
}
