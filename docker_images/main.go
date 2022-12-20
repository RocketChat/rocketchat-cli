package docker_images

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type ImageVersion struct {
	Name string
}

type ImageVersions struct {
	Count   int
	Results []ImageVersion
}

func GetImageVersions(imageName string) ImageVersions {
	requestURL := fmt.Sprintf("https://registry.hub.docker.com/v2/repositories/%s/tags?page_size=100", imageName)
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		fmt.Printf("client: could not create request: %s\n", err)
		os.Exit(1)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("client: error making http request: %s\n", err)
		os.Exit(1)
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("client: could not read response body: %s\n", err)
		os.Exit(1)
	}

	var versions ImageVersions
	err = json.Unmarshal(resBody, &versions)
	if err != nil {
		fmt.Printf("client: could not unmarshal response body: %s\n", err)
		os.Exit(1)
	}

	return versions
}
