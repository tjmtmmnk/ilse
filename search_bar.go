package ilse

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tjmtmmnk/ilse/filter"
	"github.com/tjmtmmnk/ilse/util"
)

var (
	searchBar *tview.InputField
)

func searchBarHeader() string {
	var sb strings.Builder
	sb.Grow(5)

	switch app.searchOption.Command {
	case filter.RipGrep:
		sb.WriteString("Rg")
	case filter.FuzzySearch:
		sb.WriteString("Fs")
	}

	if app.searchOption.Command == filter.RipGrep {
		sb.WriteString("|")
		switch opt := app.searchOption; {
		case opt.Mode == filter.HeadMatch:
			sb.WriteString("HM")
		case opt.Mode == filter.WordMatch:
			sb.WriteString("WM")
		case opt.Mode == filter.Regex:
			sb.WriteString("Re")
		}

		if app.searchOption.Mode != filter.Regex && app.searchOption.Case {
			sb.WriteString(",C")
		}
	}

	isOverMax := len(app.state.matched) > conf.MaxSearchResults

	if isOverMax {
		sb.WriteString(" " + fmt.Sprintf("(%d+)", conf.MaxSearchResults))
	}
	header := fmt.Sprintf("(%s) >>> ", sb.String())

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
			util.Logger.Error("search error : %v", err)
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
