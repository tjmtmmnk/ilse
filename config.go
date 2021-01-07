package ilse

import (
	"github.com/tjmtmmnk/ilse/util"
)

type config struct {
	theme   string
	workDir string
}

func newConfig() (*config, error) {
	workDir, err := util.GetWorkDir()
	if err != nil {
		return nil, err
	}

	return &config{
		theme:   "OneHalfDark",
		workDir: workDir,
	}, nil
}
