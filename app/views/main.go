package views

import (
	"github.com/gookit/event"
	"github.com/rivo/tview"
	"rocketchat-cli/app/pages"
)

type PageBuildFunc func() *tview.Flex

var pageBuildFunctions = map[pages.Page]PageBuildFunc{
	pages.WelcomePage:       buildWelcomePage,
	pages.CurrentConfigPage: buildCurrentConfigPage,
	pages.ConfigPage:        buildConfigPage,
}

func addPage(pageName pages.Page, element *tview.Flex) {
	pages.AddPage(pageName, element, true, false)
}

func InitializeMainView() *tview.Flex {
	layout := tview.NewFlex().SetDirection(tview.FlexRow)

	// Initialize pages
	pages.Initialize()

	// Build pages
	for pageName, buildFunc := range pageBuildFunctions {
		addPage(pageName, buildFunc())
	}

	// Listen to refresh page events
	event.On("refreshPage", event.ListenerFunc(func(e event.Event) error {
		page := e.Get("name").(pages.Page)
		addPage(page, pageBuildFunctions[page]())
		return nil
	}), event.AboveNormal)

	// Layout
	layout.
		AddItem(tview.NewFlex().
			AddItem(BuildTitle(), 22, 1, false), 3, 1, false).
		AddItem(tview.NewFlex().
			AddItem(pages.GetPages(), 0, 1, true), 0, 1, true)

	event.MustFire("switchPage", event.M{"name": pages.HomePage})

	return layout
}
