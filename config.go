package ilse

import (
	"fmt"
	"os"

	"github.com/tjmtmmnk/ilse/util"
)

type config struct {
	theme            string
	userWorkDir      string
	homeDir          string
	maxSearchResults int
}

func newConfig() (*config, error) {
	userWorkDir, err := util.GetUserWorkDir()
	if err != nil {
		return nil, err
	}

	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	homeDir := fmt.Sprintf("%s/.ilse", userHomeDir)

	return &config{
		theme:            "OneHalfDark",
		userWorkDir:      userWorkDir,
		homeDir:          homeDir,
		maxSearchResults: 100,
	}, nil
}
