package native

import (
		"context"
)

type nativeInstaller struct {
	// path is where rocketchat will be installed
	path string
}

type systemctlHandler struct {
	path string
}


func (n nativeInstaller) InstallNodeJS(ctx context.Context, version string) error {
	//

	return nil
}

func (n nativeInstaller) InstallMongoDB(ctx context.Context, version string) error {
	//

	return nil
}

func (n nativeInstaller) Installrocketchat(ctx context.Context, version string) error {
	//

	return nil
}
