package pages

import (
	"github.com/rivo/tview"
)

type page string

const (
	InitialSetupPage page = "initial-setup"
)

var pages = tview.NewPages()

func GetPages() *tview.Pages {
	return pages
}

func AddPage(pageName page, element *tview.Flex, resize bool, visible bool) {
	pages.AddPage(string(pageName), element, resize, visible)
}
