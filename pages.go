package ilse

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	pages *tview.Pages
)

func initPages() {
	initLayout()
	pages = tview.NewPages().AddPage("main", mainLayout, true, true)

	pages.SetBackgroundColor(tcell.ColorDefault)
}
