package filter

import (
	"os/exec"
	"strconv"
	"strings"
)

type SearchResult struct {
	FileName string
	LineNum  int
	Text     string
}

type SearchOption struct {
	Command SearchCommand
	Mode    SearchMode
}

type SearchMode int
type SearchCommand int

const (
	Regex SearchMode = iota
	FirstMatch
	FirstMatchCase
	WordMatch
	WordMatchCase
)

const (
	RipGrep SearchCommand = iota
	FuzzyFind
)

func isValidQuery(q string) bool {
	return q != ""
}

func isValidRegex(q string) bool {
	return true
}

func Search(q string, option *SearchOption) ([]SearchResult, error) {
	var (
		cmd []string
	)
	if !isValidQuery(q) {
		return []SearchResult{}, nil
	}
	switch option.Mode {
	case Regex:
		if isValidRegex(q) {
			cmd = []string{
				"rg", "--color=always", "--line-number", "--with-filename",
				"--colors=path:none", "--colors=line:none", "-e", q,
			}
		} else {
			return []SearchResult{}, nil
		}
	case FirstMatch:
		cmd = []string{
			"rg", "--color=always", "--line-number", "--with-filename",
			"--colors=path:none", "--colors=line:none", "-i", q,
		}
	case FirstMatchCase:
		cmd = []string{
			"rg", "--color=always", "--line-number", "--with-filename",
			"--colors=path:none", "--colors=line:none", q,
		}
	case WordMatch:
		cmd = []string{
			"rg", "--color=always", "--line-number", "--with-filename",
			"--colors=path:none", "--colors=line:none", "-wi", q,
		}
	case WordMatchCase:
		cmd = []string{
			"rg", "--color=always", "--line-number", "--with-filename",
			"--colors=path:none", "--colors=line:none", "-w", q,
		}
	}
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if len(out) == 0 {
		return []SearchResult{}, nil
	}

	if err != nil {
		return []SearchResult{}, err
	}

	return convert(string(out), option), nil
}

func convert(result string, option *SearchOption) []SearchResult {
	var results []SearchResult
	for _, s := range strings.Split(string(result), "\n") {
		ignore, result := split(s, option)
		if !ignore {
			results = append(results, result)
		}
	}
	return results
}

func split(str string, option *SearchOption) (ignore bool, result SearchResult) {
	switch option.Command {
	case RipGrep:
		// first remove reset flag included in path, line
		str = strings.Replace(str, "\x1b[0m", "", 4)
		// change reset flag included in text to black foreground
		str = strings.ReplaceAll(str, "\x1b[0m", "\x1b[39;40m")
		splitted := strings.Split(str, ":")
		if len(splitted) < 3 {
			ignore = true
			return
		}

		fileName := splitted[0]
		lineNum, _ := strconv.Atoi(splitted[1])
		text := splitted[2]
		result = SearchResult{fileName, lineNum, text}
	}
	return
}
