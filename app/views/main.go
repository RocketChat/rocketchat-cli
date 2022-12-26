package views

import (
	"github.com/gookit/event"
	"github.com/rivo/tview"
	"rocketchat-cli/app/pages"
)

func InitializeMainView() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	// Initialize pages
	pages.Initialize()

	// Build pages
	pages.AddPage(pages.HomePage, buildHomePage(), true, false)
	pages.AddPage(pages.SetupPage, buildSetupPage(), true, false)

	// Layout
	layout.
		AddItem(tview.NewFlex().
			AddItem(BuildTitle(), 22, 1, false), 3, 1, false).
		AddItem(tview.NewFlex().
			AddItem(pages.GetPages(), 0, 1, true), 0, 1, true)

	event.MustFire("switchPage", event.M{"name": pages.HomePage})

	return layout
}
