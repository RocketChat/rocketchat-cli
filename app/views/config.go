package views

import (
	"github.com/rivo/tview"
	"rocketchat-cli/app/views/config"
)

func buildConfigPage() *tview.Flex {
	page := tview.NewFlex().SetDirection(tview.FlexRow)

	page.AddItem(config.BuildConfigForm(), 0, 1, true)

	return page
}
