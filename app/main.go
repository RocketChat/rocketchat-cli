package app

import (
	"github.com/gookit/event"
	"github.com/rivo/tview"
	"rc-cli/app/views"
)

var App = tview.NewApplication()

func InitializeApp() {
	layout := views.InitializeMainView()

	// Listen to focus change
	event.On("setFocus", event.ListenerFunc(func(e event.Event) error {
		App.SetFocus(e.Get("element").(tview.Primitive))
		return nil
	}), event.Normal)

	// Initialize UI
	if err := App.SetRoot(layout, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
