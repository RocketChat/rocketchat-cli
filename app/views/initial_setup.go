package views

import (
	"github.com/rivo/tview"
	"rc-cli/app/views/config"
)

func BuildInitialSetupPage() *tview.Flex {
	page := tview.NewFlex().SetDirection(tview.FlexRow)

	page.AddItem(config.BuildConfigForm(), 17, 1, true)

	return page
}
