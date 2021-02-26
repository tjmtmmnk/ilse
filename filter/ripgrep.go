package filter

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"github.com/tjmtmmnk/ilse/util"
)

type rg struct{}

func newRg() *rg {
	return &rg{}
}

func isValidQuery(q string) bool {
	return q != ""
}

func isValidRegex(q string) bool {
	return true
}

func convert(result string, option *SearchOption) ([]SearchResult, error) {
	results := make([]SearchResult, 0, option.Limit)
	for i, s := range strings.Split(result, "\n") {
		if i > option.Limit {
			break
		}
		res, err := split(s, option)
		if err != nil {
			return []SearchResult{}, err
		}
		if res != nil {
			results = append(results, *res)
		}
	}
	return results, nil
}

func split(str string, option *SearchOption) (*SearchResult, error) {
	// first remove reset flag included in path, line
	str = strings.Replace(str, "\x1b[0m", "", 4)
	splitted := strings.Split(str, ":")
	if len(splitted) < 3 {
		return nil, nil
	}

	fileName := splitted[0]
	lineNum, err := strconv.Atoi(splitted[1])
	if err != nil {
		return nil, errors.New("line number wrong format")
	}
	// change reset flag included in text to black foreground
	text := strings.ReplaceAll(splitted[2], "\x1b[0m", "\x1b[39;40m")
	return &SearchResult{fileName, lineNum, text}, nil
}

func (r *rg) Search(q string, option *SearchOption) ([]SearchResult, error) {
	if !isValidQuery(q) {
		return []SearchResult{}, nil
	}

	cmd := []string{
		"rg", "--color=always", "--line-number", "--with-filename",
		"--colors=path:none", "--colors=line:none",
	}

	switch option.Mode {
	case Regex:
		if isValidRegex(q) {
			cmd = append(cmd, "-e")
		} else {
			return []SearchResult{}, nil
		}
	case HeadMatch:
	case WordMatch:
		cmd = append(cmd, "-w")
	}
	if option.Mode != Regex && !option.Case {
		cmd = append(cmd, "-i")
	}

	cmd = append(cmd, q)

	if option.TargetDir != "" {
		cmd = append(cmd, option.TargetDir)
	}

	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if len(out) == 0 {
		return []SearchResult{}, nil
	}

	if err != nil {
		util.Logger.Warn(err)
		return []SearchResult{}, err
	}

	results, err := convert(string(out), option)
	if err != nil {
		util.Logger.Warn(err)
		return []SearchResult{}, err
	}

	return results, nil
}
