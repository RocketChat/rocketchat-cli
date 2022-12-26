package config

import (
	"fmt"
	"github.com/gookit/event"
	"github.com/rivo/tview"
	"rocketchat-cli/app/pages"
	"rocketchat-cli/cli"
	"rocketchat-cli/docker"
	"rocketchat-cli/element"
	"rocketchat-cli/filesystem"
	"rocketchat-cli/matrix"
	"rocketchat-cli/nginx"
	"rocketchat-cli/traefik"
)

var configFormLayout = tview.NewFlex()
var configForm = tview.NewForm()

func AddFormItem(item tview.FormItem) {
	configForm.AddFormItem(item)
}

func BuildConfigForm() *tview.Flex {
	configFormLayout.SetDirection(tview.FlexRow).SetBorder(true)

	AddFormItem(buildRocketChatTagField())
	AddFormItem(buildSynapseTagField())
	AddFormItem(buildTraefikTagField())
	AddFormItem(buildRedisTagField())
	AddFormItem(buildNginxTagField())
	AddFormItem(buildElementTagField())
	AddFormItem(buildHostnameField())
	AddFormItem(buildTraefikEmailField())
	configForm.AddButton("Save", func() {
		cli.WriteConfigFile(filesystem.RootPath)

		docker.Compose.Services.Element.Image = fmt.Sprintf("vectorim/element-web:%s", cli.Config.Element.Tag)
		docker.Compose.Services.Element.Labels[2] = fmt.Sprintf("traefik.http.routers.traefik.rule=Host(`element.%s`)", cli.Config.Hostname)
		docker.Compose.Services.Rocketchat.Image = fmt.Sprintf("rocketchat/rocket.chat:%s", cli.Config.RocketChat.Tag)
		docker.Compose.Services.Synapse.Image = fmt.Sprintf("matrixdotorg/synapse:%s", cli.Config.Synapse.Tag)
		docker.Compose.Services.Synapse.Labels[2] = fmt.Sprintf("traefik.http.routers.traefik.rule=Host(`synapse.%s`)", cli.Config.Hostname)
		docker.Compose.Services.Traefik.Image = fmt.Sprintf("traefik:%s", cli.Config.Traefik.Tag)
		docker.Compose.Services.Traefik.Labels[2] = fmt.Sprintf("traefik.http.routers.traefik.rule=Host(`traefik.%s`)", cli.Config.Hostname)
		docker.Compose.Services.Redis.Image = fmt.Sprintf("redis:%s", cli.Config.Redis.Tag)
		docker.Compose.Services.Nginx.Image = fmt.Sprintf("nginx:%s", cli.Config.Nginx.Tag)
		docker.Compose.Services.Nginx.Labels[2] = fmt.Sprintf("traefik.http.routers.traefik.rule=Host(`%s`)", cli.Config.Hostname)
		docker.WriteComposeFile(filesystem.RootPath)

		element.Config.DefaultServerConfig.Homeserver.BaseUrl = fmt.Sprintf("https://synapse.%s", cli.Config.Hostname)
		element.Config.DefaultServerConfig.Homeserver.ServerName = cli.Config.Hostname
		element.WriteConfigFile(filesystem.DataPath)

		matrix.Synapse.ServerName = cli.Config.Hostname
		matrix.Synapse.LogConfig = fmt.Sprintf("/data/%s.log.config", cli.Config.Hostname)
		matrix.Synapse.SigningKeyPath = fmt.Sprintf("/data/%s.signing.key", cli.Config.Hostname)
		matrix.WriteHomeserverFile(matrix.SynapseType, filesystem.DataPath)
		matrix.WriteRegistrationFile(matrix.SynapseType, filesystem.DataPath)

		nginx.MatrixConf.ServerName = cli.Config.Hostname
		nginx.MatrixConf.Locations.WellknownClient.ReturnClause.Content = fmt.Sprintf("'{\"m.homeserver\": {\"base_url\": \"https://synapse.%s\"}}'", cli.Config.Hostname)
		nginx.MatrixConf.Locations.WellknownServer.ReturnClause.Content = fmt.Sprintf("'{\"m.homeserver\": {\"base_url\": \"https://synapse.%s\"}}'", cli.Config.Hostname)
		nginx.WriteMatrixConfFile(filesystem.DataPath)

		traefik.Config.CertificatesResolvers.LetsEncrypt.ACME.Email = cli.Config.Traefik.Email
		traefik.WriteACMEFileIfNeeded(filesystem.DataPath)
		traefik.WriteMiddlewaresFileIfNeeded(filesystem.DataPath)
		traefik.WriteRoutersFileIfNeeded(filesystem.DataPath)
		traefik.WriteConfigFile(filesystem.DataPath)

		event.MustFire("switchPage", event.M{"name": pages.HomePage})
		event.MustFire("switchPage", event.M{"name": pages.HomePage})
	})

	configFormLayout.AddItem(configForm, 0, 1, true)

	return configFormLayout
}
