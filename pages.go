package ilse

import (
	"github.com/rivo/tview"
)

var (
	pages *tview.Pages
)

func initPages() {
	initLayout()
	pages = tview.NewPages().AddPage("main", mainLayout, true, true)
}
