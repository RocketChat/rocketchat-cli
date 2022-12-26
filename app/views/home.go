package views

import (
	"fmt"
	"github.com/gookit/event"
	"github.com/rivo/tview"
	"rocketchat-cli/app/pages"
	"rocketchat-cli/cli"
	"strings"
)

func buildHomePage() *tview.Flex {
	page := tview.NewFlex().SetDirection(tview.FlexRow)

	welcomeText := fmt.Sprintf(`
You current install:
 - Hostname: %s

The following entries must exist on your DNS, pointing to your server's IP address:
 - %s
`, cli.Config.Hostname, strings.Join(cli.GetExpectedURLs(), "\n - "))

	uiWelcome := tview.NewTextView()
	uiWelcome.SetText(welcomeText)

	page.AddItem(uiWelcome, 11, 1, false)

	// Menu Items
	uiMenuList := tview.NewList().
		ShowSecondaryText(false).
		AddItem("Reconfigure", "", 'r', func() {
			event.MustFire("switchPage", event.M{"name": pages.SetupPage})
		}).
		AddItem("Quit", "", 'q', func() {
			event.MustFire("quit", nil)
		})

	page.AddItem(uiMenuList, 0, 1, true)

	return page
}
