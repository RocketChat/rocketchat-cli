package views

import "github.com/rivo/tview"

func BuildTitle() *tview.TextView {
	uiTitle := tview.NewTextView().
		SetText("Rocket.Chat CLI v0.1")
	uiTitle.SetBorder(true)
	return uiTitle
}
