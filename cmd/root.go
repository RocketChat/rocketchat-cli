package cmd

import (
	"github.com/rocketchat/booster/internal"
	"github.com/spf13/cobra"
)

type Opts struct {
	verbose bool
}

func RootCmd(generator *internal.LoggerGenerator) *cobra.Command {
	opts := Opts{}

	c := &cobra.Command{
		Use: "rocketchatctl",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	c.PersistentFlags().BoolVarP(&opts.verbose, "verbose", "v", false, "")

	c.AddCommand(installCommand(generator.NewLogger("install", false), &opts))
	c.AddCommand(doctorCommand(generator.NewLogger("doctor", false), &opts))
	c.AddCommand(nodeCommand())
	c.AddCommand(mongoCommand())

	return c
}
