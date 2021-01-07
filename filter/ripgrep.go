package filter

import (
	"os/exec"
)

type rg struct{}

func newRg() *rg {
	return &rg{}
}

func (r *rg) Search(q string, option *SearchOption) ([]SearchResult, error) {
	var (
		cmd []string
	)
	if !isValidQuery(q) {
		return []SearchResult{}, nil
	}
	switch mode := option.Mode; {
	case mode == Regex:
		if isValidRegex(q) {
			cmd = []string{
				"rg", "--color=always", "--line-number", "--with-filename",
				"--colors=path:none", "--colors=line:none", "-e", q,
			}
		} else {
			return []SearchResult{}, nil
		}
	case mode == HeadMatch && !option.Case:
		cmd = []string{
			"rg", "--color=always", "--line-number", "--with-filename",
			"--colors=path:none", "--colors=line:none", "-i", q,
		}
	case mode == HeadMatch && option.Case:
		cmd = []string{
			"rg", "--color=always", "--line-number", "--with-filename",
			"--colors=path:none", "--colors=line:none", q,
		}
	case mode == WordMatch && !option.Case:
		cmd = []string{
			"rg", "--color=always", "--line-number", "--with-filename",
			"--colors=path:none", "--colors=line:none", "-wi", q,
		}
	case mode == WordMatch && option.Case:
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
