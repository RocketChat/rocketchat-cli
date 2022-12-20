package pages

import (
	"github.com/gookit/event"
	"github.com/rivo/tview"
)

type Page string

const (
	SetupPage Page = "initial-setup"
	HomePage  Page = "existing-setup"
)

var pages = tview.NewPages()

func Initialize() {
	// Listen for page switching events
	event.On("switchPage", event.ListenerFunc(func(e event.Event) error {
		SwitchPage(e.Get("name").(Page))
		return nil
	}), event.Normal)

}

func GetPages() *tview.Pages {
	return pages
}

func AddPage(pageName Page, element *tview.Flex, resize bool, visible bool) {
	pages.AddPage(string(pageName), element, resize, visible)
}

func SwitchPage(pageName Page) {
	pages.SwitchToPage(string(pageName))
}
