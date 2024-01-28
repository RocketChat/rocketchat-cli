package cmd

import (
	"github.com/rocketchat/booster/internal"
	"github.com/rocketchat/booster/pkg/actions"

	"github.com/spf13/cobra"
)

func installCommand(logger *internal.Logger, opts *Opts) *cobra.Command {
	c := &cobra.Command{
		Use: "install",
		Run: func(cmd *cobra.Command, args []string) {

			logger.Debug(opts.verbose)

			if err := actions.Installrocketchat(cmd.Context(), logger, nil, actions.InstallOptions{
				//NodeVersion:  "12.22.1",
				MongoVersion: "4.4",
				Version:      "6.5.0",
			}); err != nil {
				logger.Fatalf(err, "failed to install rocketchat")
			}
		},
	}

	return c
}
