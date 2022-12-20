// Rocket.Chat CLI
// Author: Alan Sikora <alan.sikora@rocket.chat>
// https://github.com/alansikora
package main

import (
	"flag"
	"rc-cli/app"
	"rc-cli/cli"
	"rc-cli/docker"
	"rc-cli/element"
	"rc-cli/filesystem"
	"rc-cli/matrix"
	"rc-cli/nginx"
	"rc-cli/traefik"
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
