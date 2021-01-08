package ilse

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tjmtmmnk/ilse/filter"
)

var (
	frame *tview.Application
)

func initFrame() error {
	if err := initPages(); err != nil {
		return err
	}
	frame = tview.NewApplication().SetScreen(app.screen)
	frame.SetRoot(pages, true).EnableMouse(true)

	frame.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		currentPage, _ := pages.GetFrontPage()
		if currentPage != mainPage {
			return event
		}
		switch key := event.Key(); {
		case key == tcell.KeyRune || key == tcell.KeyBackspace || key == tcell.KeyBackspace2:
			frame.SetFocus(searchBar)
		case key == tcell.KeyDown && searchBar.HasFocus():
			frame.SetFocus(list)
		case key == tcell.KeyCtrlW:
			app.searchOption.Mode = filter.WordMatch
			updateSearchBarHeader()
		case key == tcell.KeyCtrlE:
			app.searchOption.Mode = filter.HeadMatch
			updateSearchBarHeader()
		case key == tcell.KeyCtrlT:
			app.searchOption.Case = !app.searchOption.Case
			updateSearchBarHeader()
		case key == tcell.KeyCtrlR:
			app.searchOption.Mode = filter.Regex
			updateSearchBarHeader()
		case key == tcell.KeyCtrlG:
			app.searchOption.Command = filter.RipGrep
			updateSearchBarHeader()
		case key == tcell.KeyCtrlF:
			app.searchOption.Command = filter.FuzzySearch
			updateSearchBarHeader()
		case key == tcell.KeyCtrlD:
			pages.SwitchToPage("tree")
		}
		return event
	})
	return nil
}
