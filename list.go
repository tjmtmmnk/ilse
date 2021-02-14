package ilse

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tjmtmmnk/ilse/filter"
	"github.com/tjmtmmnk/ilse/util"
)

var (
	list *tview.List
)

func initList() {
	initPreview()
	list = tview.NewList().ShowSecondaryText(false)

	list.SetBackgroundColor(tcell.ColorBlack)

	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if index >= len(app.state.matched) {
			return
		}
		item := app.state.matched[index]

		text, err := getPreviewContent(item)
		if err != nil {
			logger.Error("fail to fetch preview content : %v", err)
		}
		preview.SetText(text)
	})

	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		item := app.state.matched[index]
		f := func() {
			openFile(item.FileName, item.LineNum)
		}
		frame.Suspend(f)
	})

	list.SetDoneFunc(func() {
		frame.SetFocus(searchBar)
	})
}

func convertToListItems(results []filter.SearchResult) []string {
	var items []string
	for _, r := range results {
		item := fmt.Sprintf("[purple:black:-]%s[-]:[green]%d[-] %s[-:black]", util.ShortFileName(r.FileName), r.LineNum, r.Text)
		items = append(items, item)
	}
	return items
}

func updateList(items []string) {
	list.Clear()
	for _, item := range items {
		text := tview.TranslateANSI(item)
		list.AddItem(text, "", 0, nil)
	}
}
