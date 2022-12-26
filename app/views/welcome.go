package views

import (
	"github.com/gookit/event"
	"github.com/rivo/tview"
	"rocketchat-cli/app/pages"
)

func buildWelcomePage() *tview.Flex {
	page := tview.NewFlex().SetDirection(tview.FlexRow)

	uiWelcome := tview.NewTextView()
	uiWelcome.SetText("Welcome to the Rocket.Chat CLI!")

	page.AddItem(uiWelcome, 2, 1, false)

	// Menu Items
	uiMenuList := tview.NewList().
		ShowSecondaryText(false).
		AddItem("Configure", "", 'r', func() {
			event.MustFire("switchPage", event.M{"name": pages.ConfigPage})
		}).
		AddItem("Quit", "", 'q', func() {
			event.MustFire("quit", nil)
		})

	page.AddItem(uiMenuList, 0, 1, true)

	return page
}
