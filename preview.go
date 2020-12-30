package ilse

import (
	"fmt"
	"os/exec"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/tjmtmmnk/ilse/filter"
)

var (
	preview *tview.TextView
)

func initPreview() {
	preview = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true)

	preview.SetBackgroundColor(tcell.ColorDefault)
}

func getPreviewContent(item filter.SearchResult) (string, error) {
	_, h := app.screen.Size()
	from := item.Highlight - (h/2 - 1)
	if from < 0 {
		from = 0
	}
	to := item.Highlight + (h/2 + 1)
	lineRange := fmt.Sprintf("%d:%d", from, to)
	cmd := []string{"bat", "--line-range", lineRange, "--highlight-line", strconv.Itoa(item.Highlight), "--color=always", "--theme", app.config.theme, "--style=numbers", item.FileName}

	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	return string(out), err
}
