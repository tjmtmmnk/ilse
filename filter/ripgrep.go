package filter

import (
	"os/exec"

	"github.com/tjmtmmnk/ilse/util"
)

type rg struct{}

func newRg() *rg {
	return &rg{}
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
