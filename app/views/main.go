package views

import (
	"github.com/rivo/tview"
	"rc-cli/app/pages"
)

func InitializeMainView() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	// Build pages
	pages.AddPage(pages.InitialSetupPage, BuildInitialSetupPage(), true, true)

	// Layout
	layout.
		AddItem(tview.NewFlex().
			AddItem(BuildTitle(), 22, 1, false), 3, 1, false).
		AddItem(tview.NewFlex().
			AddItem(pages.GetPages(), 0, 1, true), 0, 1, true)

	return layout
}
