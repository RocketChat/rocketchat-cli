package config

import (
	"github.com/rivo/tview"
	"rc-cli/cli"
)

func buildRocketChatTagField() *tview.InputField {
	return tview.NewInputField().
		SetLabel("Rocket.Chat Image Tag").
		SetText(cli.Config.RocketChat.Tag).
		SetFieldWidth(20).
		SetChangedFunc(func(value string) {
			cli.Config.RocketChat.Tag = value
		})
}

func buildSynapseTagField() *tview.InputField {
	return tview.NewInputField().
		SetLabel("Synapse Image Tag").
		SetText(cli.Config.Synapse.Tag).
		SetFieldWidth(20).
		SetChangedFunc(func(value string) {
			cli.Config.Synapse.Tag = value
		})
}

func buildTraefikTagField() *tview.InputField {
	return tview.NewInputField().
		SetLabel("Traefik Image Tag").
		SetText(cli.Config.Traefik.Tag).
		SetFieldWidth(20).
		SetChangedFunc(func(value string) {
			cli.Config.Traefik.Tag = value
		})
}

func buildRedisTagField() *tview.InputField {
	return tview.NewInputField().
		SetLabel("Redis Image Tag").
		SetText(cli.Config.Redis.Tag).
		SetFieldWidth(20).
		SetChangedFunc(func(value string) {
			cli.Config.Redis.Tag = value
		})
}

func buildNginxTagField() *tview.InputField {
	return tview.NewInputField().
		SetLabel("Nginx Image Tag").
		SetText(cli.Config.Nginx.Tag).
		SetFieldWidth(20).
		SetChangedFunc(func(value string) {
			cli.Config.Nginx.Tag = value
		})
}

func buildElementTagField() *tview.InputField {
	return tview.NewInputField().
		SetLabel("Element Image Tag").
		SetText(cli.Config.Element.Tag).
		SetFieldWidth(20).
		SetChangedFunc(func(value string) {
			cli.Config.Element.Tag = value
		})
}

func buildHostnameField() *tview.InputField {
	return tview.NewInputField().
		SetLabel("Hostname").
		SetText(cli.Config.Hostname).
		SetFieldWidth(20).
		SetChangedFunc(func(value string) {
			cli.Config.Hostname = value
		})
}

func buildTraefikEmailField() *tview.InputField {
	return tview.NewInputField().
		SetLabel("Email").
		SetText(cli.Config.Traefik.Email).
		SetFieldWidth(20).
		SetChangedFunc(func(value string) {
			cli.Config.Traefik.Email = value
		})
}
