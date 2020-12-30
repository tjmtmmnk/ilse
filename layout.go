package ilse

import "github.com/rivo/tview"

var (
	mainLayout *tview.Flex
)

func initLayout() {
	initSearchBar()
	initList()
	initPreview()
	mainLayout = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(searchBar, 0, 3, true).SetDirection(tview.FlexRow).
		AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
			AddItem(list, 0, 1, false).
			AddItem(preview, 0, 1, false), 0, 30, false)
}
