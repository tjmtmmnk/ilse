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

	isOverMax := len(app.state.matched) > cfg.maxSearchResults

	header := fmt.Sprintf("(%s|%s) >>> ", filterName, modeName)
	if isOverMax {
		header = fmt.Sprintf("(%s|%s (%d+)) >>> ", filterName, modeName, cfg.maxSearchResults)
	}

	return header
}

func updateSearchBarHeader() {
	searchBar.SetLabel(searchBarHeader())
}

func clear() {
	list.Clear()
	preview.Clear()
	app.state.mutex.Lock()
	app.state.matched = []filter.SearchResult{}
	app.state.mutex.Unlock()
	updateSearchBarHeader()
}

func initSearchBar() {
	searchBar = tview.NewInputField().
		SetLabel(searchBarHeader()).
		SetFieldBackgroundColor(tcell.ColorBlack)

	searchBar.SetBackgroundColor(tcell.ColorBlack)

	searchBar.SetChangedFunc(func(text string) {
		if len(text) < 2 {
			clear()
			return
		}
		ftr := filter.NewFilter(app.searchOption.Command)
		results, err := ftr.Search(text, app.searchOption)
		if err != nil {
			log.Fatalf("search error : %v", err)
		}
		app.state.mutex.Lock()
		app.state.matched = results
		app.state.mutex.Unlock()
		updateSearchBarHeader()

		if len(results) > 0 {
			last := len(results)
			if len(results) > cfg.maxSearchResults {
				last = cfg.maxSearchResults
			}
			items := convertToListItems(results[0:last])
			updateList(items)
		} else {
			clear()
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
