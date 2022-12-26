// Rocket.Chat CLI
// Author: Alan Sikora <alan.sikora@rocket.chat>
// https://github.com/alansikora
package main

import (
	"flag"
	"github.com/gookit/event"
	"rocketchat-cli/app"
	"rocketchat-cli/app/pages"
	"rocketchat-cli/cli"
	"rocketchat-cli/docker"
	"rocketchat-cli/element"
	"rocketchat-cli/filesystem"
	"rocketchat-cli/matrix"
	"rocketchat-cli/nginx"
	"rocketchat-cli/traefik"
)

// Definitions
var devMode = flag.Bool("dev", false, "Enable development mode")

//var synapseLatestImageVersion string
//var synapseImageVersions []string
//var dendriteLatestImageVersion string
//var dendriteImageVersions []string
//var rocketChatLatestImageVersion string
//var rocketChatImageVersions []string

func main() {
	// Load the cli config file
	cli.ReadConfigFile(filesystem.RootPath)

	// Listen to reload config
	event.On("reloadConfig", event.ListenerFunc(func(e event.Event) error {
		cli.ReadConfigFile(filesystem.RootPath)
		event.MustFire("refreshPage", event.M{"name": pages.CurrentConfigPage})
		return nil
	}), event.AboveNormal)

	// Parse and set the flags
	flag.Parse()
	cli.Config.Config.DevMode = *devMode

	//// Load the docker image versions
	//fmt.Println("Loading Docker image versions...")
	//synapseLatestImageVersion, synapseImageVersions = GetSynapseImageVersions(false)
	//dendriteLatestImageVersion, dendriteImageVersions = GetDendriteImageVersions(false)
	//rocketChatLatestImageVersion, rocketChatImageVersions = GetRocketChatImageVersions(false)

	// Load all config files
	docker.ReadComposeFile(filesystem.DataPath)
	element.ReadConfigFile(filesystem.DataPath)
	matrix.ReadHomeserverFile(matrix.SynapseType, filesystem.DataPath)
	matrix.ReadRegistrationFile(matrix.SynapseType, filesystem.DataPath)
	nginx.ReadMatrixConfFile(filesystem.DataPath)
	traefik.ReadConfigFile(filesystem.DataPath)

	app.InitializeApp()
}
