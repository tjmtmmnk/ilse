package ilse

import (
	"github.com/tjmtmmnk/ilse/util"
)

type Config struct {
	Theme            string
	userWorkDir      string
	homeDir          string
	MaxSearchResults int
	SearchCommand    string
	SearchMode       string
	CaseSensitive    bool
}

func NewConfig() (*Config, error) {
	userWorkDir, err := util.GetUserWorkDir()
	if err != nil {
		return nil, err
	}
	homeDir := util.GetHomeDir()

	return &Config{
		Theme:            "OneHalfDark",
		userWorkDir:      userWorkDir,
		homeDir:          homeDir,
		MaxSearchResults: 100,
		SearchCommand:    "rg",
		SearchMode:       "head",
		CaseSensitive:    false,
	}, nil
}
