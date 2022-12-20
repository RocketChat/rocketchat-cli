package docker_images

import "strings"

func GetDendriteImageVersions(dev bool) (string, []string) {
	versions := GetImageVersions("matrixdotorg/dendrite-monolith")

	var latestStableVersion string
	var versionsList []string

	for _, version := range versions.Results {
		// Skip latest
		if strings.Contains(version.Name, "latest") {
			continue
		}

		// Skip develop if needed
		devMatched := strings.Contains(version.Name, "develop")
		if !dev && devMatched {
			continue
		}

		// Set latest stable imageTag
		if latestStableVersion == "" && !devMatched {
			latestStableVersion = version.Name
		}

		versionsList = append(versionsList, version.Name)
	}

	return latestStableVersion, versionsList
}
