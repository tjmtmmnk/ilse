package ilse

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

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

	preview.SetBackgroundColor(tcell.ColorBlack)
}

func getPreviewContent(item filter.SearchResult) (string, error) {
	_, h := app.screen.Size()
	from := item.LineNum - (h/2 - 1)
	if from < 0 {
		from = 0
	}
	to := item.LineNum + (h/2 - 1)
	lineRange := fmt.Sprintf("%d:%d", from, to)
	cmd := []string{"bat", "--line-range", lineRange, "--highlight-line", strconv.Itoa(item.LineNum), "--color=always", "--theme", cfg.theme, "--style=numbers,changes", item.FileName}

	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	text := string(out)
	text = strings.ReplaceAll(text, "\x1b[0m", "\x1b[39;40m")
	return tview.TranslateANSI(text), err
}
