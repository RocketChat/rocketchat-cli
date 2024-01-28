package cmd

import (
	"github.com/rocketchat/booster/internal"
	"github.com/spf13/cobra"
)

func doctorCommand(logger *internal.Logger, opts *Opts) *cobra.Command {
	c := &cobra.Command{
		Use:     "doctor",
		Short:   "Run and fix health of your install",
		Aliases: []string{"diagnostics", "diagnose", "diagnostic", "d"},
		Long:    "Doctor command is used to check health of your install and optionally patch them automagically. Additionally, the command can also be used to generate dumps for others to see and help with troubleshooting any misbehaviors.",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger.Debugf("starting checks for install health")

			return nil
		},
	}

	return c
}
