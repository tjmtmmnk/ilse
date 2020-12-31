package ilse

import (
	"log"

	"github.com/rivo/tview"
)

var (
	list *tview.List
)

func initList() {
	initPreview()
	list = tview.NewList().ShowSecondaryText(false)

	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		if index >= len(app.state.matched) {
			return
		}
		item := app.state.matched[index]

		text, err := getPreviewContent(item)
		if err != nil {
			log.Fatalf("fail to fetch preview content : %v", err)
		}
		text = tview.TranslateANSI(text)
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

func updateList(items []string) {
	list.Clear()
	for _, item := range items {
		text := tview.TranslateANSI(item)
		list.AddItem(text, "", 0, nil)
	}
}
