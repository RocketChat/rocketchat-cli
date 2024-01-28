package main

import (
	"io"
	"os"

	"github.com/rocketchat/booster/cmd"
	"github.com/rocketchat/booster/internal"
	"github.com/sirupsen/logrus"
)

const LOG_FILE = "/tmp/booster.log"

func main() {
	var writer io.Writer

	if f, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600); err != nil {
		logrus.Errorf("failed to open logfile %s, error: %v", LOG_FILE, err)
		writer = io.Discard
	} else {
		writer = f
		defer f.Close()
	}

	generator := internal.NewLoggerGenerator(writer)


	if err := cmd.RootCmd(generator).Execute(); err != nil {
		panic(err)
	}
}
