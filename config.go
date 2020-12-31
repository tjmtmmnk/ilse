package ilse

import (
	"os"
	"path/filepath"
)

type config struct {
	theme   string
	workDir string
}

func newConfig() (*config, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	evaledWorkDir, err := filepath.EvalSymlinks(workDir)
	if err != nil {
		return nil, err
	}

	return &config{
		theme:   "OneHalfDark",
		workDir: evaledWorkDir,
	}, nil
}
