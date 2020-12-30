package ilse

import (
	"log"

	"github.com/rivo/tview"
)

var (
	list *tview.List
)

func initList() {
	list = tview.NewList().ShowSecondaryText(false)

	list.SetChangedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		item := app.state.matched[index]
		if index != item.Index {
			log.Fatal("not match index")
		}

		text, err := getPreviewContent(item)

		text = tview.TranslateANSI(text)
		if err != nil {
			panic(err)
		}
		preview.Clear().SetText(text)
	})

	list.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		item := app.state.matched[index]
		if index != item.Index {
			log.Fatal("not match index")
		}
	})

	list.SetDoneFunc(func() {
		frame.SetFocus(searchBar)
	})
}
