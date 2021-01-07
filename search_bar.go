package ilse

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tjmtmmnk/ilse/filter"
)

var (
	searchBar *tview.InputField
)

func searchBarHeader() string {
	filterName := func() string {
		switch app.searchOption.Command {
		case filter.RipGrep:
			return "Rg"
		case filter.FuzzySearch:
			return "Fs"
		default:
			return ""
		}
	}()

	modeName := func() string {
		switch opt := app.searchOption; {
		case opt.Mode == filter.HeadMatch && !opt.Case:
			return "HM"
		case opt.Mode == filter.HeadMatch && opt.Case:
			return "HM,C"
		case opt.Mode == filter.WordMatch && !opt.Case:
			return "WM"
		case opt.Mode == filter.WordMatch && opt.Case:
			return "WM,C"
		case opt.Mode == filter.Regex:
			return "Re"
		default:
			return ""
		}
	}()
	return fmt.Sprintf("(%s|%s >>>)", filterName, modeName)
}

func updateSearchBarHeader() {
	searchBar.SetLabel(searchBarHeader())
}

func initSearchBar() {
	searchBar = tview.NewInputField().
		SetLabel(searchBarHeader()).
		SetFieldBackgroundColor(tcell.ColorBlack)

	searchBar.SetBackgroundColor(tcell.ColorBlack)

	searchBar.SetChangedFunc(func(text string) {
		ftr := filter.NewFilter(app.searchOption.Command)
		results, err := ftr.Search(text, app.searchOption)
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
