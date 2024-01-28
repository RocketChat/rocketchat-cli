package actions

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/rocketchat/booster/internal"
	"github.com/rocketchat/booster/pkg"
)

type InstallOptions struct {
	MongoVersion string
	NodeVersion  string
	Version      string
}

type InstallDetails []InstallOptions

// TODO: add some form of options
func Installrocketchat(ctx context.Context, logger *internal.Logger, driver pkg.Driver, opts InstallOptions) error {
	if opts.Version == "" {
		return errors.New("no rocketchat version detected")
	}

	logger.Debugf("fetching dependency requirements")
	info, err := getReleaseInfo(opts.Version, logger)
	if err != nil {
		return fmt.Errorf("failed to fetch release information: %v", err)
	}

	if opts.NodeVersion != "" && opts.NodeVersion != info.NodeVersion {
		logger.Debugf("asked nodejs version to install %s", opts.NodeVersion)
		return fmt.Errorf("passed node version does not match what is required, try not passing the parameter as it can be automatically calculated")
	} else {
		opts.NodeVersion = info.Tag
		logger.Infof("installing nodejs version %s", opts.NodeVersion)
	}

	if opts.MongoVersion != "" {
		logger.Debugf("asked mongodb version to install %s", opts.MongoVersion)
		for i, v := range info.CompatibleMongoVersions {
			if v == opts.MongoVersion {
				logger.Debugf("mongodb versions matched")
				if i == 0 {
					logger.Warnf("your mongodb version may go out of life soon, it is recommended to not use this for this install phase if possible. remove the version requirement and the appropriate version will be auto detected")
				}
				goto install
			}
		}
		return fmt.Errorf("passed mongodb version does not match any of the supported ones, try not passing the parameter as it can be automatically calculated")
	} else {
		logger.Debugf("calculating mongodb version based on the latest version of mongodb we support")
		opts.MongoVersion = info.CompatibleMongoVersions[len(info.CompatibleMongoVersions)-1]
		logger.Infof("installing mongodb version %s", opts.MongoVersion)
	}

install:
	logger.Debugf("install config %v", opts)

	return nil
}

type releaseInfo struct {
	Tag                     string   `json:"tag"`
	NodeVersion             string   `json:"nodeVersion"`
	CompatibleMongoVersions []string `json:"compatibleMongoVersions"`
}

func getReleaseInfo(version string, logger *internal.Logger) (releaseInfo, error) {
	var info releaseInfo

	r, err := http.Get(fmt.Sprintf("https://releases.rocket.chat/%s/info", version))
	if err != nil {
		return releaseInfo{}, err
	}

	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return releaseInfo{}, err
	}

	logger.Debugf("server returned release info: %v", string(b))

	err = json.Unmarshal(b, &info)
	if err != nil {
		return releaseInfo{}, err
	}

	return info, nil
}

func findInstallation(location string) (*InstallDetails, error) {
	i := &InstallDetails{}

	dirs, err := os.ReadDir(location)
	if err != nil {
		return nil, err
	}

	for _, entry := range dirs {
		_ = entry
	}

	return i, nil
}
