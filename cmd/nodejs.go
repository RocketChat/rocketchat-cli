package cmd

import (
	"log"

	"github.com/rocketchat/booster/pkg/native"
	"github.com/spf13/cobra"
)

func nodeCommand() *cobra.Command {
	c := &cobra.Command{
		Use: "nodejs",
		RunE: func(cmd *cobra.Command, args []string) error {
			n, err := native.NewNodeManager(false, "21.5.0", "/opt/booster")
			if err != nil {
				log.Fatal(err)
			}

			log.Print(n.Version())

			err = n.EnsureInstalled()
			if err != nil {
				log.Fatal(err)
			}

			out, serr, err := n.Npm().Run("--versions")
			log.Print(string(out), string(serr), err)

			return nil
		},
	}

	return c
}
