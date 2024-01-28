package native

import (
	"context"
	"os/exec"

	"github.com/rocketchat/booster/pkg"
)

type nativeDriver struct {
	opts NativeDriverOpts
}

type NativeDriverOpts struct {
	InstallDirectory  string
	rocketchatVersion string
	NodeJsVersion     string
	MongoVersion      string
}

func NewDriver(opts NativeDriverOpts) pkg.Driver {
	return nil
}

func (d *nativeDriver) MongoDbVersion(ctx context.Context) (string, error) {
	return "", nil
}

func (d *nativeDriver) NodeJsVersion(ctx context.Context) (string, error) {
	if b, err := exec.Command("node", "--version").CombinedOutput(); err != nil {
		return "", err
	} else {
		return string(b), nil
	}
}

func (d *nativeDriver) InstallNodeJs(ctx context.Context, version string) error {
	return nil
}

func (d *nativeDriver) InstallMongoDb(ctx context.Context, version string) error {
	return nil
}

func (d *nativeDriver) NodeRun(ctx context.Context, commands []string) error {
	return nil
}

func (d *nativeDriver) MongoshRun(ctx context.Context, commands []string) error {
	return nil
}
