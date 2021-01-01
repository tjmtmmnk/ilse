package ilse

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tjmtmmnk/ilse/filter"
)

var (
	searchBar *tview.InputField
)

func initSearchBar() {
	searchBar = tview.NewInputField().
		SetLabel(">>> ").
		SetFieldBackgroundColor(tcell.ColorBlack)

	searchBar.SetBackgroundColor(tcell.ColorBlack)

	fl := filter.NewFilter(filter.FuzzySearch)
	searchBar.SetChangedFunc(func(text string) {
		results, err := fl.Search(text, app.state.searchOption)
		if err != nil {
			log.Fatalf("search error : %v", err)
		}
		app.state.mutex.Lock()
		app.state.matched = results
		app.state.mutex.Unlock()
		if len(results) > 0 {
			items := convertToListItems(results)
			updateList(items)
		} else {
			list.Clear()
			preview.Clear()
		}
	})

	searchBar.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			frame.SetFocus(list)
		case tcell.KeyEsc:
			frame.Stop()
		}
	})
}
