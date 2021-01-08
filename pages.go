package ilse

import (
	"github.com/rivo/tview"
)

var (
	pages    *tview.Pages
	mainPage = "main"
	treePage = "tree"
)

func initPages() error {
	initLayout()
	if err := initTree(); err != nil {
		return err
	}
	pages = tview.NewPages().
		AddPage(mainPage, mainLayout, true, true).
		AddPage(treePage, tree, true, false)
	return nil
}
