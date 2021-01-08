package ilse

import (
	"github.com/rivo/tview"
)

var (
	pages *tview.Pages
)

func initPages() error {
	initLayout()
	if err := initTree(); err != nil {
		return err
	}
	pages = tview.NewPages().
		AddPage("main", mainLayout, true, true).
		AddPage("tree", tree, true, false)
	return nil
}
