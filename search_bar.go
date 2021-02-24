package ilse

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tjmtmmnk/ilse/filter"
)

var (
	searchBar *tview.InputField
)

func searchBarHeader() string {
	filterName := func() string {
		name := ""
		switch app.searchOption.Command {
		case filter.RipGrep:
			name = "Rg"
		case filter.FuzzySearch:
			name = "Fs"
		}
		return name
	}()

	modeName := func() string {
		if app.searchOption.Command == filter.FuzzySearch {
			return ""
		}

		name := ""
		switch opt := app.searchOption; {
		case opt.Mode == filter.HeadMatch:
			name = "HM"
		case opt.Mode == filter.WordMatch:
			name = "WM"
		case opt.Mode == filter.Regex:
			name = "Re"
		}

		if app.searchOption.Mode != filter.Regex && app.searchOption.Case {
			name += ",C"
		}

		return name
	}()

	isOverMax := len(app.state.matched) > conf.MaxSearchResults

	header := filterName
	if modeName != "" {
		header += ("|" + modeName)
	}
	if isOverMax {
		header += (" " + fmt.Sprintf("(%d+)", conf.MaxSearchResults))
	}
	header = fmt.Sprintf("(%s) >>> ", header)

	return header
}

func updateSearchBarHeader() {
	searchBar.SetLabel(searchBarHeader())
}

func clearResult() {
	list.Clear()
	preview.Clear()
	app.state.mutex.Lock()
	app.state.matched = []filter.SearchResult{}
	app.state.mutex.Unlock()
	updateSearchBarHeader()
}

func clearAll() {
	searchBar.SetText("")
	clearResult()
}

func initSearchBar() {
	searchBar = tview.NewInputField().
		SetLabel(searchBarHeader()).
		SetFieldBackgroundColor(tcell.ColorBlack)

	searchBar.SetBackgroundColor(tcell.ColorBlack)

	searchBar.SetChangedFunc(func(text string) {
		if len(text) < 2 {
			clearResult()
			return
		}
		ftr := filter.NewFilter(app.searchOption)
		results, err := ftr.Search(text, app.searchOption)
		if err != nil {
			logger.Error("search error : %v", err)
		}
		app.state.mutex.Lock()
		app.state.matched = results
		app.state.mutex.Unlock()
		updateSearchBarHeader()

		if len(results) > 0 {
			last := len(results)
			if len(results) > conf.MaxSearchResults {
				last = conf.MaxSearchResults
			}
			items := convertToListItems(results[0:last])
			updateList(items)
		} else {
			clearResult()
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
