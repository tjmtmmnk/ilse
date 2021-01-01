package filter

import (
	"os/exec"
)

type rg struct{}

func (r *rg) Search(q string, option *SearchOption) ([]SearchResult, error) {
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
