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

func initSearchBar() error {
	searchBar = tview.NewInputField().SetLabel(">>> ")

	searchBar.SetChangedFunc(func(text string) {
		results, err := filter.Search(text, nil)
		if err != nil {
			log.Fatal(err)
		}
		app.state.matched = results
	})

	searchBar.SetDoneFunc(func(key tcell.Key) {
		switch key {
		case tcell.KeyEnter:
			frame.SetFocus(list)
		case tcell.KeyEsc:
			frame.Stop()
		}
	})

	return nil
}
