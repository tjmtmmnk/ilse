package ilse

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	frame *tview.Application
)

func initFrame() {
	initPages()
	frame = tview.NewApplication().SetScreen(app.screen)
	frame.SetRoot(pages, true).EnableMouse(true)
	frame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch key := event.Key(); {
		case key == tcell.KeyRune || key == tcell.KeyBackspace || key == tcell.KeyBackspace2:
			frame.SetFocus(searchBar)
		case key == tcell.KeyDown && searchBar.HasFocus():
			frame.SetFocus(list)
		}
		return event
	})
}
