package ilse

import (
	"fmt"
	"os/exec"
	"strconv"

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
}

func getPreviewContent(item filter.SearchResult) (string, error) {
	_, h := app.screen.Size()
	from := item.LineNum - (h/2 - 1)
	if from < 0 {
		from = 0
	}
	lineRange := fmt.Sprintf("%d:", from)
	cmd := []string{"bat", "--line-range", lineRange, "--highlight-line", strconv.Itoa(item.LineNum), "--color=always", "--theme", app.config.theme, "--style=numbers,changes", item.FileName}

	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	return string(out), err
}
