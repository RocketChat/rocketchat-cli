package pages

import (
	"github.com/gookit/event"
	"github.com/rivo/tview"
	"rocketchat-cli/cli"
)

type Page string

const (
	HomePage          Page = "home"
	CurrentConfigPage Page = "current-config"
	WelcomePage       Page = "welcome"
	ConfigPage        Page = "config"
)

var pages = tview.NewPages()

func Initialize() {
	// Listen for page switching events
	event.On("switchPage", event.ListenerFunc(func(e event.Event) error {
		if e.Get("name").(Page) == HomePage {
			if cli.HasConfig() {
				SwitchPage(CurrentConfigPage)
			} else {
				SwitchPage(WelcomePage)
			}
		} else {
			SwitchPage(e.Get("name").(Page))
		}
		
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
