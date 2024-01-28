package pkg

import (
	"context"
)

// Driver "drives" an install via environment details and modifications
// "install" as in a runnable rocket.chat installation, pre or post install
type Driver interface {
	// MongoDbVersion returns the currently active mongo version for the given rocket.chat version
	MongoDbVersion(ctx context.Context) (string, error)

	// NodeJsVersion returns the currently active nodejs version for the given rocket.chat version
	NodeJsVersion(ctx context.Context) (string, error)

	// InstallNodeJS installs the provided nodejs version
	InstallNodeJS(ctx context.Context, version string) error

	// InstallMongoDB install the provided mongodb version
	InstallMongoDB(ctx context.Context, version string) error

	// NodeExecl is a runner for node binary
	NodeExecl(ctx context.Context, args ...string) error

	// MongoshExecl is a runner for mongosh binary
	MongoshExecl(ctx context.Context, args ...string) error

	// MongodExecl is a runner for mongod binary
	MongodExecl(ctx context.Context, args ...string) error
}
