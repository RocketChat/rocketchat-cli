package views

import (
	"github.com/rivo/tview"
	"rc-cli/app/views/config"
)

func buildSetupPage() *tview.Flex {
	page := tview.NewFlex().SetDirection(tview.FlexRow)

	page.AddItem(config.BuildConfigForm(), 0, 1, true)

	return page
}