package cmd

import (
	"log"
	"time"

	"github.com/rocketchat/booster/pkg/native"
	"github.com/spf13/cobra"
)

func mongoCommand() *cobra.Command {
	c := &cobra.Command{
		Use: "mongo",
		RunE: func(cmd *cobra.Command, args []string) error {
			m, err := native.NewMongoManager("6.0.4", "/opt/booster")
			if err != nil {
				log.Fatal(err)
			}

			err = m.Install()
			if err != nil {
				return err
			}

			v, err := m.Mongod().Version()
			if err != nil {
				return err
			}

			log.Println(v)

			log.Println("starting mongod")

			mongod := m.Mongod()


			err = mongod.Start()
			if err != nil {
				return err
			}

			time.Sleep(time.Second * 10)
			err = m.ReplStart()
			if err != nil {
				return err
			}
			log.Println("repl started")

			m.Repl("print('DATABASE VERSION :', db.version())\n")

			time.Sleep(time.Second * 5)


			log.Println("stopping mongod")

			return mongod.Stop()
		},
	}

	return c
}
