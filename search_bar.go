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
	searchBar = tview.NewInputField().SetLabel(">>> ")

	searchBar.SetChangedFunc(func(text string) {
		results, err := filter.Search(text, app.state.searchOption)
		if err != nil {
			log.Fatal("search")
			log.Fatal(err)
		}
		app.state.mutex.Lock()
		app.state.matched = results
		app.state.mutex.Unlock()
		if len(results) > 0 {
			var texts []string
			for _, r := range results {
				texts = append(texts, r.Text)
			}
			updateList(texts)
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
